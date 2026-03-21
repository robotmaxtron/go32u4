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
)

type ATmega32u4 struct {
	CPU        *cpu.CPU
	FlashData  [FlashSize]uint16
	SRAMData   [SRAMSize]uint8
	IORegData  [0x100]uint8
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
	return mcu
}

// ReadSRAM bus.Bus implementation
func (m *ATmega32u4) ReadSRAM(address uint16) uint8 {
	if address < 32 {
		return m.CPU.Reg[address]
	}
	if address < 32+uint16(len(m.IORegData)) {
		return m.ReadIO(address - 32)
	}
	sramAddr := address - (32 + uint16(len(m.IORegData)))
	if int(sramAddr) < len(m.SRAMData) {
		return m.SRAMData[sramAddr]
	}
	return 0
}

func (m *ATmega32u4) WriteSRAM(address uint16, value uint8) {
	if address < 32 {
		m.CPU.Reg[address] = value
	} else if address < 32+uint16(len(m.IORegData)) {
		m.WriteIO(address-32, value)
	} else {
		sramAddr := address - (32 + uint16(len(m.IORegData)))
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
	return m.Periph.IOCallback(address, 0, false)
}

func (m *ATmega32u4) WriteIO(address uint16, value uint8) {
	if address == 0x3F {
		m.CPU.SREG = value
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

// IORegs peripherals.System implementation
func (m *ATmega32u4) IORegs() []uint8               { return m.IORegData[:] }
func (m *ATmega32u4) TriggerInterrupt(vector uint8) { m.PendingInterrupts |= (1 << vector) }
func (m *ATmega32u4) Cycles() uint64                { return m.CPU.Cycles }
func (m *ATmega32u4) SaveEEPROM() error {
	if m.EEPROMFile == "" {
		return nil
	}
	return os.WriteFile(m.EEPROMFile, m.Periph.EEPROM[:], 0644)
}

func (m *ATmega32u4) PinCallback(port int8, mask uint8, value uint8) {
	if m.PinCallbackFunc != nil {
		m.PinCallbackFunc(port, mask, value)
	}
}

// SetGlobalInterrupts bus.InterruptController implementation
func (m *ATmega32u4) SetGlobalInterrupts(enabled bool) { m.GlobalInterrupts = enabled }
func (m *ATmega32u4) GetGlobalInterrupts() bool        { return m.GlobalInterrupts }
func (m *ATmega32u4) ClearInterrupt(vector uint8)      { m.PendingInterrupts &= ^(1 << vector) }

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

func (m *ATmega32u4) Step() error {
	if m.GlobalInterrupts && m.PendingInterrupts != 0 {
		m.Periph.SleepEnabled = false
		m.handleInterrupts()
		return nil
	}
	if m.Periph.SleepEnabled {
		m.CPU.Cycles++
		m.Periph.Tick(1)
		return nil
	}
	return m.CPU.Step()
}

func (m *ATmega32u4) handleInterrupts() {
	for i := uint8(1); i < 43; i++ {
		if (m.PendingInterrupts & (1 << i)) != 0 {
			m.executeInterrupt(i)
			m.PendingInterrupts &= ^(1 << i)
			break
		}
	}
}

func (m *ATmega32u4) executeInterrupt(vector uint8) {
	m.CPU.Push(uint8(m.CPU.PC & 0xFF))
	m.CPU.Push(uint8((m.CPU.PC >> 8) & 0xFF))
	m.CPU.SetFlag(cpu.SREG_I, false)
	m.CPU.PC = uint16(vector) * 2
	m.CPU.Cycles += 4
	m.Periph.Tick(4)
}
