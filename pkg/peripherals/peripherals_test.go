package peripherals_test

import (
	"go32u4/pkg/peripherals"
	"testing"
)

type MockSystem struct {
	ioRegs []uint8
	ints   uint64
	cycles uint64
}

func (m *MockSystem) IORegs() []uint8                                { return m.ioRegs }
func (m *MockSystem) TriggerInterrupt(vector uint8)                  { m.ints |= 1 << vector }
func (m *MockSystem) Cycles() uint64                                 { return m.cycles }
func (m *MockSystem) SaveEEPROM() error                              { return nil }
func (m *MockSystem) PinCallback(port int8, mask uint8, value uint8) {}
func (m *MockSystem) FlashWrite(address uint16, value uint16)        {}
func (m *MockSystem) FlashErase(address uint16)                      {}

func TestTimer0(t *testing.T) {
	sys := &MockSystem{ioRegs: make([]uint8, 256)}
	p := peripherals.NewManager(sys)

	// Configure Timer0: Prescaler 8
	p.IOCallback(peripherals.TCCR0B, 2, true)
	p.IOCallback(peripherals.TIMSK0, 1, true) // TOIE0

	p.Timer0Counter = 255
	sys.cycles = 8 // Divisor is 8

	p.Tick(1)
	if p.Timer0Counter != 0 {
		t.Errorf("Expected overflow to 0, got %d", p.Timer0Counter)
	}
	if (sys.ints & (1 << 23)) == 0 {
		t.Error("Expected Timer0 Overflow interrupt triggered")
	}
}

func TestEEPROM(t *testing.T) {
	sys := &MockSystem{ioRegs: make([]uint8, 256)}
	p := peripherals.NewManager(sys)

	// Write value to EEPROM
	p.IOCallback(peripherals.EEARL, 0x10, true)
	p.IOCallback(peripherals.EEDR, 0x55, true)
	p.IOCallback(peripherals.EECR, 0x04, true) // EEMPE
	p.IOCallback(peripherals.EECR, 0x06, true) // EEPE (EEMPE | EEPE)

	if p.EEPROM[0x10] != 0x55 {
		t.Errorf("Expected EEPROM[0x10] 0x55, got %02X", p.EEPROM[0x10])
	}

	// Read value from EEPROM
	p.IOCallback(peripherals.EEARL, 0x10, true)
	p.IOCallback(peripherals.EECR, 0x01, true) // EERE

	val := p.IOCallback(peripherals.EEDR, 0, false)
	if val != 0x55 {
		t.Errorf("Expected EEDR 0x55, got %02X", val)
	}
}

func TestTWIState(t *testing.T) {
	sys := &MockSystem{ioRegs: make([]uint8, 256)}
	p := peripherals.NewManager(sys)
	mcp := peripherals.NewMCP23018(0x20, p)
	p.RegisterTWIClient(mcp)

	// Start TWI
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2)|(1<<5), true) // TWINT, TWEN, TWSTA
	if p.TWIState != 0x08 {
		t.Errorf("Expected TWIState 0x08 (START), got %02X", p.TWIState)
	}

	// Write Address
	sys.ioRegs[peripherals.TWDR] = 0x20 << 1             // Slave 0x20, write
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2), true) // TWINT, TWEN (trigger state update)

	if p.TWIState != 0x18 {
		t.Errorf("Expected TWIState 0x18 (SLA+W ACK), got %02X", p.TWIState)
	}
}

type MCP23018Wrapper struct {
	*peripherals.MCP23018
}

func (w *MCP23018Wrapper) Regs(idx uint8) uint8 {
	// Since we can't access private 'regs' from another package, 
	// we use TWI to read it or provide a helper in mcp23018.go.
	// Actually, let's just make 'regs' exported in mcp23018.go for testing or provide a GetReg method.
	return 0 // TODO: fix
}

func TestMCP23018(t *testing.T) {
	sys := &MockSystem{ioRegs: make([]uint8, 256)}
	p := peripherals.NewManager(sys)
	mcp := peripherals.NewMCP23018(0x20, p)
	p.RegisterTWIClient(mcp)

	// 1. Start TWI and Address MCP23018 (0x20) for Write
	sys.ioRegs[peripherals.TWDR] = 0x20 << 1                    // 0x40 (SLA+W)
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2)|(1<<5), true) // START
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2), true)        // Trigger updateTWIState (SLA+W)

	if p.TWIState != 0x18 {
		t.Errorf("Expected TWIState 0x18, got %02X", p.TWIState)
	}

	// 2. Write Register Address (IODIRA)
	sys.ioRegs[peripherals.TWDR] = peripherals.MCP23018_IODIRA
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2), true) // Data ACK

	if p.TWIState != 0x28 {
		t.Errorf("Expected TWIState 0x28, got %02X", p.TWIState)
	}

	// 3. Write Value to IODIRA (0x00 - all outputs)
	sys.ioRegs[peripherals.TWDR] = 0x00
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2), true) // Data ACK

	// We can't access mcp.regs directly, so we'll read it back via TWI in step 5.
	// Or we could have exported Regs in MCP23018.

	// 4. STOP
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2)|(1<<4), true) // STOP

	// 5. Read back IODIRA
	// START + SLA+W
	sys.ioRegs[peripherals.TWDR] = 0x20 << 1
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2)|(1<<5), true)
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2), true)
	// Write register addr
	sys.ioRegs[peripherals.TWDR] = peripherals.MCP23018_IODIRA
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2), true)
	// Repeated START + SLA+R
	sys.ioRegs[peripherals.TWDR] = (0x20 << 1) | 0x01          // 0x41 (SLA+R)
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2)|(1<<5), true) // Repeated START
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2), true)        // SLA+R

	if p.TWIState != 0x40 {
		t.Errorf("Expected TWIState 0x40, got %02X", p.TWIState)
	}
	if sys.ioRegs[peripherals.TWDR] != 0x00 {
		t.Errorf("Expected TWDR to have IODIRA value 0x00, got %02X", sys.ioRegs[peripherals.TWDR])
	}
}

func TestErgoDoxMatrixScan(t *testing.T) {
	sys := &MockSystem{ioRegs: make([]uint8, 256)}
	p := peripherals.NewManager(sys)
	mcp := peripherals.NewMCP23018(0x20, p)
	p.RegisterTWIClient(mcp)

	// ErgoDox Left Hand (MCP23018) Initialization
	// 1. Set IODIRA (columns) as outputs, IODIRB (rows) as inputs
	// SLA+W
	sys.ioRegs[peripherals.TWDR] = 0x20 << 1
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2)|(1<<5), true)
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2), true)

	// Reg Addr 0x00 (IODIRA)
	sys.ioRegs[peripherals.TWDR] = 0x00
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2), true)

	// IODIRA = 0x00 (outputs)
	sys.ioRegs[peripherals.TWDR] = 0x00
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2), true)
	// IODIRB = 0xFF (inputs)
	sys.ioRegs[peripherals.TWDR] = 0xFF
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2), true)

	mcp.OnStop() // Reset selected for next transaction

	// 2. Set GPPU B (rows) pull-ups enabled
	// STOP then START
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2)|(1<<4), true)
	sys.ioRegs[peripherals.TWDR] = 0x20 << 1
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2)|(1<<5), true)
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2), true)

	// Reg Addr 0x0D (GPPUB)
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2)|(1<<4), true)
	sys.ioRegs[peripherals.TWDR] = 0x20 << 1
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2)|(1<<5), true)
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2), true)
	sys.ioRegs[peripherals.TWDR] = 0x0D
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2), true)
	// GPPUB = 0xFF
	sys.ioRegs[peripherals.TWDR] = 0xFF
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2), true)

	mcp.OnStop()

	// 3. Simulate Key Press: Row 0, Col 0
	// In ErgoDox, Columns are often driven LOW.
	// So we set OLATA to 0xFE (Col 0 low, others high/high-Z)
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2)|(1<<4), true)
	sys.ioRegs[peripherals.TWDR] = 0x20 << 1
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2)|(1<<5), true)
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2), true)
	sys.ioRegs[peripherals.TWDR] = 0x14 // OLATA
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2), true)
	sys.ioRegs[peripherals.TWDR] = 0xFE
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2), true)

	mcp.OnStop()

	// External world: User presses key at Row 0, Col 0.
	// This connects Row 0 (PB0) to Col 0 (PA0).
	// Since PA0 is driven LOW, PB0 will also go LOW.
	mcp.External = 0xFEFE // PA0 (bit 0) and PB0 (bit 8) both pulled low by key connection

	// 4. Scan: Read GPIOB
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2)|(1<<4), true)
	sys.ioRegs[peripherals.TWDR] = 0x20 << 1
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2)|(1<<5), true)
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2), true)
	sys.ioRegs[peripherals.TWDR] = 0x13 // GPIOB
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2), true)

	// SLA+R
	sys.ioRegs[peripherals.TWDR] = (0x20 << 1) | 0x01
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2)|(1<<5), true)
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2), true)

	// Verify result
	res := sys.ioRegs[peripherals.TWDR]
	if res&0x01 != 0 {
		t.Errorf("Expected GPIOB Bit 0 to be LOW (key pressed), got %02X", res)
	}
	if res&0x02 == 0 {
		t.Errorf("Expected GPIOB Bit 1 to be HIGH (key not pressed), got %02X", res)
	}
}

func TestWatchdog(t *testing.T) {
	sys := &MockSystem{ioRegs: make([]uint8, 256)}
	p := peripherals.NewManager(sys)

	// 1. Enable Timed Change
	p.IOCallback(peripherals.WDTCSR, (1<<peripherals.WDCE)|(1<<peripherals.WDE), true)
	if p.WatchdogTimedChange != 4 {
		t.Errorf("Expected Timed Change window of 4, got %d", p.WatchdogTimedChange)
	}

	// 2. Set Prescaler and Enable Interrupt Mode
	// No WDP bits set = prescaler 0 = 2048 ticks
	p.IOCallback(peripherals.WDTCSR, 1<<peripherals.WDIE, true)
	if p.WatchdogTimedChange != 0 {
		t.Errorf("Expected Timed Change window closed, got %d", p.WatchdogTimedChange)
	}

	// Check timeout: 2048 * 125 = 256000
	if p.WatchdogTimeout != 256000 {
		t.Errorf("Expected Timeout 256000 cycles, got %d", p.WatchdogTimeout)
	}

	// 3. Tick until interrupt
	p.Tick(255999)
	if (sys.ints & (1 << 4)) != 0 {
		t.Error("Watchdog interrupt triggered too early")
	}

	p.Tick(1)
	if (sys.ints & (1 << 4)) == 0 {
		t.Error("Watchdog interrupt NOT triggered at timeout")
	}

	// 4. Verify WDIE is cleared after interrupt (since WDE is not set)
	if (sys.ioRegs[peripherals.WDTCSR] & (1 << peripherals.WDIE)) != 0 {
		t.Error("Expected WDIE to be cleared by hardware after interrupt")
	}

	// 5. Enable System Reset Mode
	// Timed change again
	p.IOCallback(peripherals.WDTCSR, (1<<peripherals.WDCE)|(1<<peripherals.WDE), true)
	// Set WDE and WDP0 (prescaler 1 = 4096 ticks)
	p.IOCallback(peripherals.WDTCSR, (1<<peripherals.WDE)|(1<<peripherals.WDP0), true)

	// Check timeout: 4096 * 125 = 512000
	if p.WatchdogTimeout != 512000 {
		t.Errorf("Expected Timeout 512000 cycles, got %d", p.WatchdogTimeout)
	}

	p.Tick(511999)
	if p.WatchdogReset {
		t.Error("Watchdog reset triggered too early")
	}

	p.Tick(1)
	if !p.WatchdogReset {
		t.Error("Watchdog reset NOT triggered at timeout")
	}
}

func TestTimer4(t *testing.T) {
	sys := &MockSystem{ioRegs: make([]uint8, 256)}
	p := peripherals.NewManager(sys)

	// 1. Test 10-bit register access via TC4H
	p.IOCallback(peripherals.TC4H, 0x02, true)  // High bits = 2 (bit 9 set)
	p.IOCallback(peripherals.TCNT4, 0x55, true) // Low bits = 0x55
	if p.Timer4Counter != 0x255 {
		t.Errorf("Expected TCNT4 0x255, got %03X", p.Timer4Counter)
	}

	p.IOCallback(peripherals.TC4H, 0x01, true)
	p.IOCallback(peripherals.OCR4C, 0xAA, true)
	if p.Timer4OCR4C != 0x1AA {
		t.Errorf("Expected OCR4C 0x1AA, got %03X", p.Timer4OCR4C)
	}

	// 2. Test reading 10-bit register
	p.Timer4Counter = 0x3FF
	val := p.IOCallback(peripherals.TCNT4, 0, false)
	if val != 0xFF {
		t.Errorf("Expected TCNT4 low byte 0xFF, got %02X", val)
	}
	high := p.IOCallback(peripherals.TC4H, 0, false)
	if high != 0x03 {
		t.Errorf("Expected TC4H 0x03, got %02X", high)
	}

	// 3. Test Timer 4 counting and overflow with OCR4C as TOP
	p.Timer4Counter = 0x1A9
	p.Timer4OCR4C = 0x1AA
	p.Timer4ControlB = 0x01                      // Prescaler 1
	p.IOCallback(peripherals.TIMSK4, 1<<2, true) // TOIE4

	p.Tick(1)
	if p.Timer4Counter != 0x1AA {
		t.Errorf("Expected TCNT4 0x1AA, got %03X", p.Timer4Counter)
	}

	p.Tick(1)
	if p.Timer4Counter != 0 {
		t.Errorf("Expected TCNT4 overflow to 0, got %03X", p.Timer4Counter)
	}
	if (sys.ints & (1 << 39)) == 0 {
		t.Error("Expected Timer4 Overflow interrupt triggered (vector 39)")
	}

	// 4. Test PLL Clocking (PCKE)
	sys.ints = 0
	p.Timer4Counter = 0
	p.Timer4OCR4C = 0x3FF // Default TOP
	p.PLLControl = 1 << 2 // PCKE set

	// With PCKE, 1 system cycle = 4 Timer 4 cycles
	p.Tick(1)
	if p.Timer4Counter != 4 {
		t.Errorf("Expected TCNT4 4 after 1 system cycle with PCKE, got %d", p.Timer4Counter)
	}

	// 5. Test OCR4A compare match
	p.Timer4Counter = 0x100
	p.Timer4OCR4A = 0x102
	p.IOCallback(peripherals.TIMSK4, 1<<6, true) // OCIE4A

	p.Tick(1) // TCNT4 becomes 0x104 (with PCKE)
	if (sys.ints & (1 << 38)) == 0 {
		t.Error("Expected Timer4 OCR4A interrupt triggered (vector 38)")
	}
}

func TestUSBFullEmulation(t *testing.T) {
	sys := &MockSystem{ioRegs: make([]uint8, 256)}
	p := peripherals.NewManager(sys)

	// 1. Test Endpoint Selection
	p.IOCallback(peripherals.UENUM, 3, true)
	if p.USBSelectedEP != 3 {
		t.Errorf("Expected USBSelectedEP 3, got %d", p.USBSelectedEP)
	}

	// 2. Test FIFO write/read
	p.IOCallback(peripherals.UEDATX, 0xAA, true)
	p.IOCallback(peripherals.UEDATX, 0xBB, true)

	if len(p.USBEndpoints[3].FIFO) != 2 {
		t.Errorf("Expected FIFO length 2, got %d", len(p.USBEndpoints[3].FIFO))
	}

	val1 := p.IOCallback(peripherals.UEDATX, 0, false)
	if val1 != 0xAA {
		t.Errorf("Expected 0xAA, got %02X", val1)
	}

	val2 := p.IOCallback(peripherals.UEDATX, 0, false)
	if val2 != 0xBB {
		t.Errorf("Expected 0xBB, got %02X", val2)
	}

	// 3. Test EP0 Setup handling
	p.USBEndpoints[0].SetupFIFO = []byte{0x00, 0x05, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00} // SET_ADDRESS 0x10
	p.Tick(1)

	if p.USBAddress != 0x10 {
		t.Errorf("Expected USBAddress 0x10, got %02X", p.USBAddress)
	}
	if (p.USBEndpoints[0].Interrupt & (1 << 3)) == 0 {
		t.Error("Expected RXSTPI interrupt flag set")
	}
	if (p.USBEndpoints[0].Interrupt & (1 << 0)) == 0 {
		t.Error("Expected TXINI interrupt flag set (Status Stage)")
	}

	// 4. Test HID Class request
	p.USBEndpoints[0].SetupFIFO = []byte{0x21, 0x09, 0x01, 0x02, 0x00, 0x00, 0x00, 0x00} // SET_REPORT (HID)
	p.Tick(1)
	if (p.USBEndpoints[0].Interrupt & (1 << 0)) == 0 {
		t.Error("Expected TXINI for HID SET_REPORT")
	}

	// 5. Test HID Data Injection
	p.USBConfigured = true
	p.HIDKeyMap[0] = 0x04 // 'a' key
	p.Tick(1)
	if len(p.USBEndpoints[1].FIFO) != 8 {
		t.Errorf("Expected 8 bytes in EP1 FIFO, got %d", len(p.USBEndpoints[1].FIFO))
	}
	if p.USBEndpoints[1].FIFO[0] != 0x04 {
		t.Errorf("Expected key code 0x04, got %02X", p.USBEndpoints[1].FIFO[0])
	}
	if (p.USBEndpoints[1].Interrupt & (1 << 0)) == 0 {
		t.Error("Expected TXINI for injected HID report")
	}
}

func TestUSBPullUp(t *testing.T) {
	sys := &MockSystem{ioRegs: make([]uint8, 256)}
	p := peripherals.NewManager(sys)
	mcp := peripherals.NewMCP23018(0x20, p)
	p.RegisterTWIClient(mcp)

	// 1. Verify default pull-up
	if p.PullUpResistor != 2200.0 {
		t.Errorf("Expected default pull-up 2200.0, got %f", p.PullUpResistor)
	}

	// 2. Test I2C communication with default pull-up (should succeed)
	sys.ioRegs[peripherals.TWDR] = 0x20 << 1
	p.TWIState = 0x08
	sys.ioRegs[peripherals.TWCR] = 1 << 7               // Set TWINT bit to simulate hardware ready
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2), true) // trigger updateTWIState by "clearing" TWINT
	if p.TWIState != 0x18 {
		t.Errorf("Expected TWIState 0x18 (ACK), got %02X", p.TWIState)
	}

	// 3. Test I2C communication without pull-up (should fail)
	p.PullUpResistor = 2000000.0 // 2M ohm (too high)
	p.TWIState = 0x08
	sys.ioRegs[peripherals.TWCR] = 1 << 7
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2), true)
	if p.TWIState != 0xF8 {
		t.Errorf("Expected TWIState 0xF8 (Error), got %02X", p.TWIState)
	}

	// 4. Test MCP23018 pin read with and without pull-up
	p.PullUpResistor = 2200.0
	// addr already 0x20
	// IODIRA is 0xFF by default
	// GPPUA is 0x00 by default
	mcp.External = 0xFFFF                        // External high-Z

	// Should read 1 because of external pull-up
	sys.ioRegs[peripherals.TWDR] = (0x20 << 1) | 1 // SLA+R
	p.TWIState = 0x08
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2), true)
	
	// Data state 0x40 (SLA+R ACK)
	if p.TWIState != 0x40 {
		t.Errorf("Expected TWIState 0x40, got %02X", p.TWIState)
	}
	if sys.ioRegs[peripherals.TWDR] != 0xFF {
		t.Errorf("Expected 0xFF with 2.2k pull-up, got %02X", sys.ioRegs[peripherals.TWDR])
	}

	p.PullUpResistor = 2000000.0 // No pull-up
	p.TWIState = 0x08
	sys.ioRegs[peripherals.TWDR] = (0x20 << 1) | 1 // SLA+R
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2), true)
	if p.TWIState != 0xF8 {
		t.Errorf("Expected TWIState 0xF8 (No Pull-Up), got %02X", p.TWIState)
	}
}

func TestMacros(t *testing.T) {
	sys := &MockSystem{ioRegs: make([]uint8, 256)}
	p := peripherals.NewManager(sys)
	p.USBConfigured = true

	// 1. Define a macro: 'a', then delay 10 cycles, then 'b'
	macro := peripherals.MacroTable{
		Records: []peripherals.MacroRecord{
			{KeyMap: [8]uint8{0x04}, Delay: 10},
			{KeyMap: [8]uint8{0x05}, Delay: 5},
		},
	}

	// 2. Trigger the macro
	p.TriggerMacro(macro)
	if !p.MacroActive {
		t.Error("Expected MacroActive to be true after TriggerMacro")
	}

	// 3. First Tick: should process the first record ('a')
	p.Tick(1)
	if p.HIDKeyMap[0] != 0x04 {
		t.Errorf("Expected HIDKeyMap[0] 0x04, got %02X", p.HIDKeyMap[0])
	}
	if p.MacroDelayCounter != 10 {
		t.Errorf("Expected MacroDelayCounter 10, got %d", p.MacroDelayCounter)
	}

	// 4. Tick through delay (9 cycles)
	p.Tick(9)
	if p.HIDKeyMap[0] != 0x04 {
		t.Errorf("Expected HIDKeyMap[0] still 0x04 during delay, got %02X", p.HIDKeyMap[0])
	}
	if p.MacroDelayCounter != 1 {
		t.Errorf("Expected MacroDelayCounter 1, got %d", p.MacroDelayCounter)
	}

	// 5. Next Tick: should finish delay and process second record ('b')
	p.Tick(1)
	if p.HIDKeyMap[0] != 0x05 {
		t.Errorf("Expected HIDKeyMap[0] 0x05 after delay, got %02X", p.HIDKeyMap[0])
	}
	if p.MacroDelayCounter != 5 {
		t.Errorf("Expected MacroDelayCounter 5, got %d", p.MacroDelayCounter)
	}

	// 6. Final Tick: finish delay and clear MacroActive
	p.Tick(5)
	if p.MacroActive {
		t.Error("Expected MacroActive to be false after finishing macro")
	}
}
