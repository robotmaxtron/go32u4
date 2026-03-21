package mcu

import (
	"go32u4/pkg/cpu"
	"testing"
)

func TestInterruptPrioritization(t *testing.T) {
	m := NewATmega32u4()
	m.GlobalInterrupts = true
	m.CPU.SetFlag(cpu.SREG_I, true)

	// Trigger multiple interrupts
	// Bit i triggers Vector i+1
	// Vector 2 (INT0) - Bit 1
	// Vector 5 (INT3) - Bit 4
	// Vector 11 (TIMER1 COMPA) - Bit 10
	// Vector 41 (USB General) - Bit 40
	
	m.TriggerInterrupt(40)
	m.TriggerInterrupt(4)
	m.TriggerInterrupt(10)
	m.TriggerInterrupt(1)

	// Step 1: Vector 2 should be handled first
	err := m.Step()
	if err != nil { t.Fatal(err) }
	// Vector 2 is at address (2-1)*2 = 2
	if m.CPU.PC != 2 {
		t.Errorf("Expected PC to be at vector 2 (2), got %d", m.CPU.PC)
	}
	
	// Set PC back to main and re-enable interrupts (normally RETI does this)
	m.CPU.PC = 100
	m.CPU.SetFlag(cpu.SREG_I, true)
	m.GlobalInterrupts = true

	// Step 2: Vector 5 should be next
	err = m.Step()
	if err != nil { t.Fatal(err) }
	// Vector 5 is at address (5-1)*2 = 8
	if m.CPU.PC != 8 {
		t.Errorf("Expected PC to be at vector 5 (8), got %d", m.CPU.PC)
	}

	m.CPU.PC = 100
	m.CPU.SetFlag(cpu.SREG_I, true)
	m.GlobalInterrupts = true

	// Step 3: Vector 11 should be next
	err = m.Step()
	if err != nil { t.Fatal(err) }
	// Vector 11 is at address (11-1)*2 = 20
	if m.CPU.PC != 20 {
		t.Errorf("Expected PC to be at vector 11 (20), got %d", m.CPU.PC)
	}

	m.CPU.PC = 100
	m.CPU.SetFlag(cpu.SREG_I, true)
	m.GlobalInterrupts = true

	// Step 4: Vector 41 should be last
	err = m.Step()
	if err != nil { t.Fatal(err) }
	// Vector 41 is at address (41-1)*2 = 80
	if m.CPU.PC != 80 {
		t.Errorf("Expected PC to be at vector 41 (80), got %d", m.CPU.PC)
	}
}
