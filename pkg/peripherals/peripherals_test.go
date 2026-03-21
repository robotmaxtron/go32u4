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

func (m *MockSystem) IORegs() []uint8                       { return m.ioRegs }
func (m *MockSystem) TriggerInterrupt(vector uint8)       { m.ints |= (1 << vector) }
func (m *MockSystem) Cycles() uint64                       { return m.cycles }
func (m *MockSystem) SaveEEPROM() error                    { return nil }
func (m *MockSystem) PinCallback(port int8, mask uint8, value uint8) {}

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

	// Start TWI
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2)|(1<<5), true) // TWINT, TWEN, TWSTA
	if p.TWIState != 0x08 {
		t.Errorf("Expected TWIState 0x08 (START), got %02X", p.TWIState)
	}

	// Write Address
	sys.ioRegs[peripherals.TWDR] = 0x20 // Slave 0x10, write
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2), true) // TWINT, TWEN (trigger state update)
	
	if p.TWIState != 0x18 {
		t.Errorf("Expected TWIState 0x18 (SLA+W ACK), got %02X", p.TWIState)
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
	p.IOCallback(peripherals.WDTCSR, (1<<peripherals.WDIE), true)
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
