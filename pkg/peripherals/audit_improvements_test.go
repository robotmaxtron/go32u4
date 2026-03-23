package peripherals_test

import (
	"go32u4/pkg/mcu"
	"go32u4/pkg/peripherals"
	"testing"
)

func TestPRR(t *testing.T) {
	m := mcu.NewATmega32u4()
	ioRegs := m.IORegs()

	// 1. Disable Timer 0 via PRR0
	ioRegs[peripherals.PRR0] |= (1 << peripherals.PRTIM0)
	m.WriteIO(peripherals.TCCR0B, 0x01) // Start Timer 0
	m.WriteIO(peripherals.TCNT0, 0x00)
	
	m.Step()
	if m.ReadIO(peripherals.TCNT0) != 0 {
		t.Error("Timer 0 should not tick when disabled by PRR0")
	}

	// 2. Enable Timer 0 via PRR0
	ioRegs[peripherals.PRR0] &= ^uint8(1 << peripherals.PRTIM0)
	m.Step()
	if m.ReadIO(peripherals.TCNT0) == 0 {
		t.Error("Timer 0 should tick when enabled by PRR0")
	}
}

func TestUEBCHX(t *testing.T) {
	m := mcu.NewATmega32u4()
	
	m.WriteIO(peripherals.UENUM, 1)
	// Fill FIFO with 300 bytes
	for i := 0; i < 300; i++ {
		m.WriteIO(peripherals.UEDATX, uint8(i))
	}

	low := m.ReadIO(peripherals.UEBCLX)
	high := m.ReadIO(peripherals.UEBCHX)
	count := uint16(high)<<8 | uint16(low)

	if count != 300 {
		t.Errorf("Expected 300 bytes in FIFO, got %d (low: %d, high: %d)", count, low, high)
	}
	
	// Read one byte
	m.ReadIO(peripherals.UEDATX)
	low = m.ReadIO(peripherals.UEBCLX)
	high = m.ReadIO(peripherals.UEBCHX)
	count = uint16(high)<<8 | uint16(low)
	if count != 299 {
		t.Errorf("Expected 299 bytes in FIFO after read, got %d", count)
	}
}

func TestTimer4PLL96(t *testing.T) {
	m := mcu.NewATmega32u4()
	
	m.WriteIO(peripherals.PLLCSR, (1 << peripherals.PCKE) | (1 << 4)) // PLL on, 96MHz
	m.WriteIO(peripherals.TCCR4B, 0x01) // Prescaler 1
	m.WriteIO(peripherals.TCNT4, 0x00)
	
	m.Step() // 1 system cycle @ 16MHz
	// Should be 6 Timer 4 cycles
	val := m.ReadIO(peripherals.TCNT4)
	if val != 6 {
		t.Errorf("Expected TCNT4 to be 6 after 1 system cycle at 96MHz, got %d", val)
	}
}

func TestExternalInterrupts(t *testing.T) {
	m := mcu.NewATmega32u4()
	
	// Configure INT0 for falling edge
	m.WriteIO(peripherals.EICRA, 0x02) // Mode 2: Falling edge for INT0
	m.WriteIO(peripherals.EIMSK, 0x01) // Enable INT0
	
	// Simulate PD0 falling edge
	m.PinCallback('D', 0x01, 0x00) // Mask PD0, Value 0
	
	if (m.ReadIO(peripherals.EIFR) & 0x01) == 0 {
		t.Error("EIFR bit 0 (INTF0) should be set")
	}
	
	// Check if interrupt is pending
	// Vector 2 is INT0
	if (m.PendingInterrupts & (1 << 2)) == 0 {
		t.Errorf("INT0 interrupt (bit 2) should be pending, got %x", m.PendingInterrupts)
	}

	// Test PCINT0
	m.WriteIO(peripherals.PCICR, 0x01) // Enable PCIE0
	m.WriteIO(peripherals.PCMSK0, 0x01) // Enable PCINT0
	m.PinCallback('B', 0x01, 0x01) // Toggle PB0
	
	if (m.ReadIO(peripherals.PCIFR) & 0x01) == 0 {
		t.Error("PCIFR bit 0 should be set")
	}
}
