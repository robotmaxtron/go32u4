package mcu

import (
	"go32u4/pkg/bus"
	"go32u4/pkg/cpu"
	"go32u4/pkg/peripherals"
	"os"
)

const (
	FlashSize    = 32768
	SRAMSize     = 2560
	IORegSize    = 64
	ExtIORegSize = 160
	TotalIORegs  = 0x100
)

type ATmega32u4 struct {
	CPU        *cpu.CPU
	FlashData  [FlashSize]uint16
	SRAMData   [SRAMSize]uint8
	IORegData  [TotalIORegs]uint8
	Periph     *peripherals.Manager
	EEPROMFile string

	PendingInterrupts uint64
	GlobalInterrupts  bool

	IOCallback      bus.IOCallback
	PinCallbackFunc bus.PinCallback
}

func NewATmega32u4() *ATmega32u4 {
	mcu := &ATmega32u4{}
	mcu.Periph = peripherals.NewManager(mcu)
	mcu.CPU = cpu.NewCPU(mcu, mcu)
	mcu.CPU.TickPeripherals = mcu.Periph.Tick
	mcu.Reset()
	return mcu
}

// ReadSRAM bus.Bus implementation
func (m *ATmega32u4) ReadSRAM(address uint16) uint8 {
	if address < 32 {
		return m.CPU.Reg[address]
	}
	if address < 32+TotalIORegs {
		return m.ReadIO(address - 32)
	}
	sramAddr := address - (32 + TotalIORegs)
	if int(sramAddr) < len(m.SRAMData) {
		return m.SRAMData[sramAddr]
	}
	return 0
}

func (m *ATmega32u4) WriteSRAM(address uint16, value uint8) {
	if address < 32 {
		m.CPU.Reg[address] = value
	} else if address < 32+TotalIORegs {
		m.WriteIO(address-32, value)
	} else {
		sramAddr := address - (32 + TotalIORegs)
		if int(sramAddr) < len(m.SRAMData) {
			m.SRAMData[sramAddr] = value
		}
	}
}

func (m *ATmega32u4) ReadIO(address uint16) uint8 {
	if address == 0x3F {
		return m.CPU.SREG
	}
	if address == 0x3E {
		return uint8(m.CPU.SP >> 8)
	}
	if address == 0x3D {
		return uint8(m.CPU.SP & 0xFF)
	}
	// Handle standard registers like PORTB (0x05)
	// through Periph.IOCallback, which returns ioRegs[address] if not handled.
	return m.Periph.IOCallback(address, 0, false)
}

func (m *ATmega32u4) WriteIO(address uint16, value uint8) {
	if address == 0x3F {
		m.CPU.SREG = value
		m.GlobalInterrupts = (value & 0x80) != 0
		return
	}
	if address == 0x3E {
		m.CPU.SP = (m.CPU.SP & 0x00FF) | (uint16(value) << 8)
		return
	}
	if address == 0x3D {
		m.CPU.SP = (m.CPU.SP & 0xFF00) | uint16(value)
		return
	}
	m.Periph.IOCallback(address, value, true)
}

func (m *ATmega32u4) Flash() []uint16 { return m.FlashData[:] }
func (m *ATmega32u4) FlashWrite(address uint16, value uint16) {
	if int(address) < len(m.FlashData) {
		// If SPMCSR.SPMEN is set and no other bits, it's a buffer write.
		spmcsr := m.ReadIO(0x37)
		if spmcsr&0x01 != 0 && spmcsr&0x1E == 0 {
			m.Periph.SPMBuffer[address%64] = value
		} else {
			m.FlashData[address] = value
		}
	}
}

func (m *ATmega32u4) FlashErase(address uint16) {
	// Address is a word address. Erase a page.
	// Page size for 32u4 is 128 bytes (64 words).
	pageSizeWords := uint16(64)
	pageStart := (address / pageSizeWords) * pageSizeWords
	for i := uint16(0); i < pageSizeWords; i++ {
		if int(pageStart+i) < len(m.FlashData) {
			m.FlashData[pageStart+i] = 0xFFFF
		}
	}
}

func (m *ATmega32u4) FlashCommit(address uint16) {
	// Commits SPM buffer to flash page
	// Page size for 32u4 is 128 bytes (64 words).
	pageSizeWords := uint16(64)
	pageStart := (address / pageSizeWords) * pageSizeWords
	for i := uint16(0); i < pageSizeWords; i++ {
		if int(pageStart+i) < len(m.FlashData) {
			m.FlashData[pageStart+i] = m.Periph.SPMBuffer[i]
		}
	}
}

func (m *ATmega32u4) SetSleep(enabled bool) {
	// SLEEP instruction only takes effect if SE (Sleep Enable) bit in SMCR is set.
	// SMCR is at IO address 0x53.
	if m.ReadIO(0x53)&0x01 != 0 {
		m.Periph.SleepEnabled = enabled
		m.CPU.Halted = enabled
	}
}

// IORegs peripherals.System implementation
func (m *ATmega32u4) IORegs() []uint8               { return m.IORegData[:] }
func (m *ATmega32u4) TriggerInterrupt(vector uint8) { m.PendingInterrupts |= 1 << uint64(vector) }
func (m *ATmega32u4) Cycles() uint64                { return m.CPU.Cycles }
func (m *ATmega32u4) SaveEEPROM() error {
	if m.EEPROMFile == "" {
		return nil
	}
	return os.WriteFile(m.EEPROMFile, m.Periph.EEPROM[:], 0644)
}

func (m *ATmega32u4) PinCallback(port int8, mask uint8, value uint8) {
	m.Periph.HandlePinChange(port, mask, value)
	if m.PinCallbackFunc != nil {
		m.PinCallbackFunc(port, mask, value)
	}
}

// SetGlobalInterrupts bus.InterruptController implementation
func (m *ATmega32u4) SetGlobalInterrupts(enabled bool) { m.GlobalInterrupts = enabled }
func (m *ATmega32u4) GetGlobalInterrupts() bool        { return m.GlobalInterrupts }
func (m *ATmega32u4) ClearInterrupt(vector uint8)      { m.PendingInterrupts &= ^(1 << uint64(vector)) }

func (m *ATmega32u4) LoadEEPROM(filename string) error {
	m.EEPROMFile = filename
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	copy(m.Periph.EEPROM[:], data)
	return nil
}

func (m *ATmega32u4) Reset() {
	// Reset CPU
	m.CPU.PC = 0
	m.CPU.SP = uint16(cpu.SRAMStart + SRAMSize - 1)
	m.CPU.SREG = 0
	m.CPU.Halted = false
	for i := range m.CPU.Reg {
		m.CPU.Reg[i] = 0
	}

	// Reset IO Registers
	for i := range m.IORegData {
		m.IORegData[i] = 0
	}

	// Reset Peripherals state
	m.Periph.Reset()

	m.PendingInterrupts = 0
	m.GlobalInterrupts = false
}

func (m *ATmega32u4) Step() error {
	if m.PendingInterrupts != 0 {
		// Wake up on any pending interrupt
		m.Periph.SleepEnabled = false
		m.CPU.Halted = false
	}

	if m.Periph.WatchdogReset {
		m.Reset()
		return nil
	}

	if m.GlobalInterrupts && m.PendingInterrupts != 0 {
		m.handleInterrupts()
		return nil
	}

	if m.Periph.SleepEnabled {
		m.CPU.Cycles++
		m.Periph.Tick(1)
		return nil
	}

	// Capture PC before Step to check for SPM
	pc := m.CPU.PC
	err := m.CPU.Step()
	if err != nil {
		return err
	}

	// Check if the instruction just executed was SPM (0x95E8)
	if int(pc) < len(m.FlashData) && m.FlashData[pc] == 0x95E8 {
		m.handleSPM()
	}

	return nil
}

func (m *ATmega32u4) handleSPM() {
	spmcsr := m.ReadIO(0x37)
	if (spmcsr & (1 << 0)) == 0 { // SPMEN must be set
		return
	}
	z := (uint16(m.CPU.Reg[31]) << 8) | uint16(m.CPU.Reg[30])
	if (spmcsr & (1 << 1)) != 0 { // PGERS (Page Erase)
		pageAddr := z / 2
		for i := uint16(0); i < 64; i++ {
			m.FlashData[pageAddr+i] = 0xFFFF
		}
	} else if (spmcsr & (1 << 2)) != 0 { // PGWRT (Page Write)
		pageAddr := z / 2
		for i := uint16(0); i < 64; i++ {
			m.FlashData[pageAddr+i] = m.Periph.SPMBuffer[i]
		}
	} else if (spmcsr & (1 << 0)) != 0 { // SPMEN set but no PGERS/PGWRT -> Buffer Fill
		word := (uint16(m.CPU.Reg[1]) << 8) | uint16(m.CPU.Reg[0])
		offset := (z & 0x7E) / 2
		m.Periph.SPMBuffer[offset] = word
	}
	// Clear SPMEN after SPM execution
	m.WriteIO(0x37, spmcsr&^uint8(1<<0))
}

func (m *ATmega32u4) handleInterrupts() {
	// ATmega32u4 has vectors 1 to 42.
	// Vector 1: RESET (handled separately usually)
	// Vector 2: INT0 (bit 1)
	// Vector 42: USB General (bit 41)
	// PendingInterrupts is a bitmask where bit i corresponds to Vector i+1.
	for i := uint8(1); i < 43; i++ {
		if (m.PendingInterrupts & (1 << uint64(i))) != 0 {
			m.executeInterrupt(i + 1)
			m.PendingInterrupts &= ^(1 << uint64(i))
			break
		}
	}
}

func (m *ATmega32u4) executeInterrupt(vector uint8) {
	// Address = (VectorNumber - 1) * 2
	m.CPU.Push(uint8((m.CPU.PC >> 8) & 0xFF))
	m.CPU.Push(uint8(m.CPU.PC & 0xFF))
	m.CPU.SetFlag(cpu.SREG_I, false)
	m.GlobalInterrupts = false
	m.CPU.PC = uint16(vector-1) * 2
}
