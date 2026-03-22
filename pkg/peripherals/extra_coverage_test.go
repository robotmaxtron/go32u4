package peripherals_test

import (
	"go32u4/pkg/peripherals"
	"testing"
)

func TestSPM(t *testing.T) {
	sys := &MockSystem{ioRegs: make([]uint8, 256)}
	p := peripherals.NewManager(sys)

	// Set SPMEN bit in SPMCSR
	p.IOCallback(peripherals.SPMCSR, 1<<0, true) // 1<<0 is SPMEN

	if sys.ioRegs[peripherals.SPMCSR]&(1<<0) == 0 {
		t.Fatal("Expected SPMEN bit to be set")
	}

	// Tick 2 cycles
	p.Tick(2)
	if sys.ioRegs[peripherals.SPMCSR]&(1<<0) == 0 {
		t.Error("Expected SPMEN bit to still be set after 2 cycles")
	}

	// Tick 2 more cycles (total 4)
	p.Tick(2)
	if sys.ioRegs[peripherals.SPMCSR]&(1<<0) != 0 {
		t.Error("Expected SPMEN bit to be cleared after 4 cycles")
	}
}

func TestIOCallbackPinToggle(t *testing.T) {
	sys := &MockSystem{ioRegs: make([]uint8, 256)}
	p := peripherals.NewManager(sys)

	// Writing to PINB toggles PORTB
	sys.ioRegs[peripherals.PORTB] = 0xAA
	p.IOCallback(peripherals.PINB, 0xFF, true)
	if sys.ioRegs[peripherals.PORTB] != 0x55 {
		t.Errorf("Expected PORTB 0x55, got %02X", sys.ioRegs[peripherals.PORTB])
	}

	// Writing to PINC toggles PORTC
	sys.ioRegs[peripherals.PORTC] = 0x00
	p.IOCallback(peripherals.PINC, 0x01, true)
	if sys.ioRegs[peripherals.PORTC] != 0x01 {
		t.Errorf("Expected PORTC 0x01, got %02X", sys.ioRegs[peripherals.PORTC])
	}
}

func TestUART1Interrupts(t *testing.T) {
	sys := &MockSystem{ioRegs: make([]uint8, 256)}
	p := peripherals.NewManager(sys)

	// Enable UART1 TX interrupts
	sys.ioRegs[peripherals.UCSR1B] = (1 << 6) | (1 << 5) // TXCIE1 and UDRIE1

	// Writing to UDR1 should trigger interrupts
	p.IOCallback(peripherals.UDR1, 'A', true)

	if (sys.ints & (1 << 27)) == 0 {
		t.Error("Expected UART1 TX Complete interrupt")
	}
	if (sys.ints & (1 << 26)) == 0 {
		t.Error("Expected UART1 Data Register Empty interrupt")
	}
}

func TestSPIInterrupts(t *testing.T) {
	sys := &MockSystem{ioRegs: make([]uint8, 256)}
	p := peripherals.NewManager(sys)

	// Enable SPI interrupt
	sys.ioRegs[peripherals.SPCR] = (1 << 7) // SPIE

	// Writing to SPDR should trigger interrupt
	p.IOCallback(peripherals.SPDR, 0x55, true)

	if (sys.ints & (1 << 18)) == 0 {
		t.Error("Expected SPI interrupt")
	}
	if (sys.ioRegs[peripherals.SPSR] & (1 << 7)) == 0 {
		t.Error("Expected SPIF bit set")
	}
}

func TestADCInterrupts(t *testing.T) {
	sys := &MockSystem{ioRegs: make([]uint8, 256)}
	p := peripherals.NewManager(sys)

	// Set ADSC and ADIE
	p.IOCallback(peripherals.ADCSRA, (1<<6)|(1<<3), true)

	if (sys.ints & (1 << 29)) == 0 {
		t.Error("Expected ADC interrupt")
	}
	if (sys.ioRegs[peripherals.ADCSRA] & (1 << 4)) == 0 {
		t.Error("Expected ADIF bit set")
	}
	if (sys.ioRegs[peripherals.ADCSRA] & (1 << 6)) != 0 {
		t.Error("Expected ADSC bit cleared")
	}
}

func TestUSBRegisters(t *testing.T) {
	sys := &MockSystem{ioRegs: make([]uint8, 256)}
	p := peripherals.NewManager(sys)

	// Test UENUM and dependent registers
	p.IOCallback(peripherals.UENUM, 2, true)
	if p.USBSelectedEP != 2 {
		t.Errorf("Expected USBSelectedEP 2, got %d", p.USBSelectedEP)
	}

	p.IOCallback(peripherals.UECFG0X, 0x80, true)
	if p.USBEndpoints[2].Config0 != 0x80 {
		t.Errorf("Expected Config0 0x80, got %02X", p.USBEndpoints[2].Config0)
	}

	// Select EP2 again to verify read-back into IO regs
	p.IOCallback(peripherals.UENUM, 2, true)
	if sys.ioRegs[peripherals.UECFG0X] != 0x80 {
		t.Errorf("Expected IORegs[UECFG0X] 0x80, got %02X", sys.ioRegs[peripherals.UECFG0X])
	}
}

func TestWatchdogWindow(t *testing.T) {
	sys := &MockSystem{ioRegs: make([]uint8, 256)}
	p := peripherals.NewManager(sys)

	// Initiate timed change
	p.IOCallback(peripherals.WDTCSR, (1<<peripherals.WDCE)|(1<<peripherals.WDE), true)
	if p.WatchdogTimedChange != 4 {
		t.Errorf("Expected WatchdogTimedChange 4, got %d", p.WatchdogTimedChange)
	}

	p.Tick(4)
	if p.WatchdogTimedChange != 0 {
		t.Errorf("Expected WatchdogTimedChange 0 after 4 cycles, got %d", p.WatchdogTimedChange)
	}

	if sys.ioRegs[peripherals.WDTCSR]&(1<<peripherals.WDCE) != 0 {
		t.Error("Expected WDCE bit to be cleared after 4 cycles")
	}
}

func TestIOCallbackReads(t *testing.T) {
	sys := &MockSystem{ioRegs: make([]uint8, 256)}
	p := peripherals.NewManager(sys)

	// Timer 0
	p.Timer0Counter = 0x12
	if p.IOCallback(peripherals.TCNT0, 0, false) != 0x12 {
		t.Error("TCNT0 read mismatch")
	}
	p.Timer0ControlA = 0x34
	if p.IOCallback(peripherals.TCCR0A, 0, false) != 0x34 {
		t.Error("TCCR0A read mismatch")
	}
	p.Timer0ControlB = 0x56
	if p.IOCallback(peripherals.TCCR0B, 0, false) != 0x56 {
		t.Error("TCCR0B read mismatch")
	}
	p.Timer0CompareA = 0x78
	if p.IOCallback(peripherals.OCR0A, 0, false) != 0x78 {
		t.Error("OCR0A read mismatch")
	}
	p.Timer0CompareB = 0x9A
	if p.IOCallback(peripherals.OCR0B, 0, false) != 0x9A {
		t.Error("OCR0B read mismatch")
	}

	// Timer 1
	p.Timer1Counter = 0x1234
	if p.IOCallback(peripherals.TCNT1L, 0, false) != 0x34 {
		t.Error("TCNT1L read mismatch")
	}
	if p.IOCallback(peripherals.TCNT1H, 0, false) != 0x12 {
		t.Error("TCNT1H read mismatch")
	}
	p.Timer1CompareA = 0x5678
	if p.IOCallback(peripherals.OCR1AL, 0, false) != 0x78 {
		t.Error("OCR1AL read mismatch")
	}
	if p.IOCallback(peripherals.OCR1AH, 0, false) != 0x56 {
		t.Error("OCR1AH read mismatch")
	}
	p.Timer1ControlA = 0x9A
	if p.IOCallback(peripherals.TCCR1A, 0, false) != 0x9A {
		t.Error("TCCR1A read mismatch")
	}
	p.Timer1ControlB = 0xBC
	if p.IOCallback(peripherals.TCCR1B, 0, false) != 0xBC {
		t.Error("TCCR1B read mismatch")
	}

	// Timer 3
	p.Timer3Counter = 0x1122
	if p.IOCallback(peripherals.TCNT3L, 0, false) != 0x22 {
		t.Error("TCNT3L read mismatch")
	}
	if p.IOCallback(peripherals.TCNT3H, 0, false) != 0x11 {
		t.Error("TCNT3H read mismatch")
	}
	p.Timer3CompareA = 0x3344
	if p.IOCallback(peripherals.OCR3AL, 0, false) != 0x44 {
		t.Error("OCR3AL read mismatch")
	}
	if p.IOCallback(peripherals.OCR3AH, 0, false) != 0x33 {
		t.Error("OCR3AH read mismatch")
	}
	p.Timer3ControlA = 0x55
	if p.IOCallback(peripherals.TCCR3A, 0, false) != 0x55 {
		t.Error("TCCR3A read mismatch")
	}
	p.Timer3ControlB = 0x66
	if p.IOCallback(peripherals.TCCR3B, 0, false) != 0x66 {
		t.Error("TCCR3B read mismatch")
	}

	// UART1
	p.UART1RXBuffer = []byte{0xAA, 0xBB}
	sys.ioRegs[peripherals.UCSR1A] |= (1 << 7) // RXC1
	if p.IOCallback(peripherals.UDR1, 0, false) != 0xAA {
		t.Error("UDR1 read mismatch 1")
	}
	if (sys.ioRegs[peripherals.UCSR1A] & (1 << 7)) == 0 {
		t.Error("RXC1 should still be set")
	}
	if p.IOCallback(peripherals.UDR1, 0, false) != 0xBB {
		t.Error("UDR1 read mismatch 2")
	}
	if (sys.ioRegs[peripherals.UCSR1A] & (1 << 7)) != 0 {
		t.Error("RXC1 should be cleared")
	}
	if p.IOCallback(peripherals.UDR1, 0, false) != 0 {
		t.Error("UDR1 empty read mismatch")
	}

	// SPI
	p.SPIBuffer = 0xCC
	sys.ioRegs[peripherals.SPSR] |= (1 << 7) // SPIF
	if p.IOCallback(peripherals.SPDR, 0, false) != 0xCC {
		t.Error("SPDR read mismatch")
	}
	if (sys.ioRegs[peripherals.SPSR] & (1 << 7)) != 0 {
		t.Error("SPIF should be cleared")
	}

	// ADC
	p.ADCValue = 0x1234
	if p.IOCallback(peripherals.ADCL, 0, false) != 0x34 {
		t.Error("ADCL read mismatch")
	}
	if p.IOCallback(peripherals.ADCH, 0, false) != 0x12 {
		t.Error("ADCH read mismatch")
	}

	// USB
	p.USBSelectedEP = 1
	p.USBEndpoints[1].FIFO = []byte{0xDD}
	sys.ioRegs[peripherals.UEBCLX] = 1
	if p.IOCallback(peripherals.UEDATX, 0, false) != 0xDD {
		t.Error("UEDATX read mismatch")
	}
	if sys.ioRegs[peripherals.UEBCLX] != 0 {
		t.Error("UEBCLX update mismatch")
	}
	if p.IOCallback(peripherals.UEDATX, 0, false) != 0 {
		t.Error("UEDATX empty read mismatch")
	}
	p.USBSelectedEP = 8 // Invalid
	if p.IOCallback(peripherals.UEDATX, 0, false) != 0 {
		t.Error("UEDATX invalid EP read mismatch")
	}
	if p.IOCallback(peripherals.UEBCLX, 0, false) != 0 {
		t.Error("UEBCLX invalid EP read mismatch")
	}

	// Timer 4
	p.Timer4ControlA = 0x11
	if p.IOCallback(peripherals.TCCR4A, 0, false) != 0x11 {
		t.Error("TCCR4A read mismatch")
	}
	p.Timer4ControlB = 0x22
	if p.IOCallback(peripherals.TCCR4B, 0, false) != 0x22 {
		t.Error("TCCR4B read mismatch")
	}
	p.Timer4ControlC = 0x33
	if p.IOCallback(peripherals.TCCR4C, 0, false) != 0x33 {
		t.Error("TCCR4C read mismatch")
	}
	p.Timer4ControlD = 0x44
	if p.IOCallback(peripherals.TCCR4D, 0, false) != 0x44 {
		t.Error("TCCR4D read mismatch")
	}
	p.Timer4ControlE = 0x55
	if p.IOCallback(peripherals.TCCR4E, 0, false) != 0x55 {
		t.Error("TCCR4E read mismatch")
	}
	p.Timer4DT4 = 0x66
	if p.IOCallback(peripherals.DT4, 0, false) != 0x66 {
		t.Error("DT4 read mismatch")
	}

	// PLL
	p.PLLControl = 0x04
	sys.ioRegs[peripherals.PLLCSR] = 0x01 // PLOCK
	if p.IOCallback(peripherals.PLLCSR, 0, false) != 0x05 {
		t.Error("PLLCSR read mismatch")
	}

	// Default
	sys.ioRegs[0x50] = 0xAA
	if p.IOCallback(0x50, 0, false) != 0xAA {
		t.Error("Default read mismatch")
	}
}

func TestTWIStateErrors(t *testing.T) {
	sys := &MockSystem{ioRegs: make([]uint8, 256)}
	p := peripherals.NewManager(sys)

	// SLA+W NACK
	p.TWIState = 0x08 // START
	sys.ioRegs[peripherals.TWDR] = 0x30 << 1 // SLA+W to non-existent client
	sys.ioRegs[peripherals.TWCR] = (1 << 7)
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2), true)
	if p.TWIState != 0x20 {
		t.Errorf("Expected TWIState 0x20, got %02X", p.TWIState)
	}

	// SLA+W NACK -> Error
	p.TWIState = 0x20
	sys.ioRegs[peripherals.TWCR] = (1 << 7)
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2), true)
	if p.TWIState != 0xF8 {
		t.Errorf("Expected TWIState 0xF8, got %02X", p.TWIState)
	}

	// SLA+R NACK
	p.TWIState = 0x08 // START
	sys.ioRegs[peripherals.TWDR] = 0x30 << 1 | 1 // SLA+R to non-existent client
	sys.ioRegs[peripherals.TWCR] = (1 << 7)
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2), true)
	if p.TWIState != 0x48 {
		t.Errorf("Expected TWIState 0x48, got %02X", p.TWIState)
	}

	// SLA+R NACK -> Error
	p.TWIState = 0x48
	sys.ioRegs[peripherals.TWCR] = (1 << 7)
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2), true)
	if p.TWIState != 0xF8 {
		t.Errorf("Expected TWIState 0xF8, got %02X", p.TWIState)
	}

	// Data NACK (Write)
	client := &MockTWIClient{addr: 0x20}
	p.RegisterTWIClient(client)
	p.TWIState = 0x18
	p.ActiveTWI = client
	client.ackWrite = false
	sys.ioRegs[peripherals.TWCR] = (1 << 7)
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2), true)
	if p.TWIState != 0x30 {
		t.Errorf("Expected TWIState 0x30, got %02X", p.TWIState)
	}

	// Data NACK (Write) from Data ACK state
	p.TWIState = 0x28
	client.ackWrite = false
	sys.ioRegs[peripherals.TWCR] = (1 << 7)
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2), true)
	if p.TWIState != 0x30 {
		t.Errorf("Expected TWIState 0x30, got %02X", p.TWIState)
	}

	// SLA+R ACK with no client (should not happen normally but for coverage)
	p.TWIState = 0x40
	p.ActiveTWI = nil
	sys.ioRegs[peripherals.TWCR] = (1 << 7)
	p.IOCallback(peripherals.TWCR, (1<<7)|(1<<2), true)
	if p.TWIState != 0x48 {
		t.Errorf("Expected TWIState 0x48, got %02X", p.TWIState)
	}
}

type MockTWIClient struct {
	addr     uint8
	ackWrite bool
}

func (m *MockTWIClient) Address() uint8           { return m.addr }
func (m *MockTWIClient) OnStart(isRead bool) bool { return true }
func (m *MockTWIClient) OnWrite(data uint8) bool  { return m.ackWrite }
func (m *MockTWIClient) OnRead() uint8           { return 0 }
func (m *MockTWIClient) OnStop()                 {}

func TestEP0SetupStandardRequests(t *testing.T) {
	sys := &MockSystem{ioRegs: make([]uint8, 256)}
	p := peripherals.NewManager(sys)

	// GET_DESCRIPTOR Device
	p.USBEndpoints[0].SetupFIFO = []byte{0x80, 0x06, 0x00, 0x01, 0x00, 0x00, 0x40, 0x00}
	p.Tick(1)
	if len(p.USBEndpoints[0].FIFO) == 0 {
		t.Error("Expected device descriptor in FIFO")
	}

	// GET_DESCRIPTOR Config
	p.USBEndpoints[0].FIFO = nil
	p.USBEndpoints[0].SetupFIFO = []byte{0x80, 0x06, 0x00, 0x02, 0x00, 0x00, 0x40, 0x00}
	p.Tick(1)
	if len(p.USBEndpoints[0].FIFO) == 0 {
		t.Error("Expected config descriptor in FIFO")
	}

	// GET_DESCRIPTOR HID Report
	p.USBEndpoints[0].FIFO = nil
	p.USBEndpoints[0].SetupFIFO = []byte{0x81, 0x06, 0x00, 0x22, 0x00, 0x00, 0x40, 0x00}
	p.Tick(1)
	if len(p.USBEndpoints[0].FIFO) != 47 {
		t.Errorf("Expected HID report descriptor length 47, got %d", len(p.USBEndpoints[0].FIFO))
	}

	// SET_REPORT (GET_REPORT in ep0setup is 0x01)
	p.USBEndpoints[0].SetupFIFO = []byte{0xA1, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00}
	p.Tick(1)
	if (p.USBEndpoints[0].Interrupt & (1 << 1)) == 0 {
		t.Error("Expected RXOUTI for GET_REPORT")
	}
}

func TestWatchdogResetClear(t *testing.T) {
	sys := &MockSystem{ioRegs: make([]uint8, 256)}
	p := peripherals.NewManager(sys)

	// Enable Watchdog
	p.IOCallback(peripherals.WDTCSR, (1<<peripherals.WDCE)|(1<<peripherals.WDE), true)
	p.IOCallback(peripherals.WDTCSR, (1<<peripherals.WDE), true)

	// Clear WDE during timed change window
	p.IOCallback(peripherals.WDTCSR, (1<<peripherals.WDCE)|(1<<peripherals.WDE), true)
	p.IOCallback(peripherals.WDTCSR, 0, true)

	if (sys.ioRegs[peripherals.WDTCSR] & (1 << peripherals.WDE)) != 0 {
		t.Error("Expected WDE to be cleared")
	}
}
