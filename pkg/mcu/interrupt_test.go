package mcu

import (
	"go32u4/pkg/cpu"
	"testing"
)

func TestInterruptPrioritization(t *testing.T) {
	m := NewATmega32u4()
	m.CPU.SetFlag(cpu.SREG_I, true)

	// Trigger multiple interrupts
	// Vector 1 (INT0) - Priority 2
	// Vector 4 (INT3) - Priority 5
	// Vector 10 (TIMER1 COMPA) - Priority 11
	// Vector 40 (USB) - Priority 41
	
	m.TriggerInterrupt(40)
	m.TriggerInterrupt(4)
	m.TriggerInterrupt(10)
	m.TriggerInterrupt(1)

	// Step 1: Vector 1 should be handled first
	err := m.Step()
	if err != nil { t.Fatal(err) }
	if m.CPU.PC != 1*2 {
		t.Errorf("Expected PC to be at vector 1, got %d", m.CPU.PC)
	}
	
	// Set PC back to main and re-enable interrupts (normally RETI does this)
	m.CPU.PC = 100
	m.CPU.SetFlag(cpu.SREG_I, true)

	// Step 2: Vector 4 should be next
	err = m.Step()
	if err != nil { t.Fatal(err) }
	if m.CPU.PC != 4*2 {
		t.Errorf("Expected PC to be at vector 4, got %d", m.CPU.PC)
	}

	m.CPU.PC = 100
	m.CPU.SetFlag(cpu.SREG_I, true)

	// Step 3: Vector 10 should be next
	err = m.Step()
	if err != nil { t.Fatal(err) }
	if m.CPU.PC != 10*2 {
		t.Errorf("Expected PC to be at vector 10, got %d", m.CPU.PC)
	}

	m.CPU.PC = 100
	m.CPU.SetFlag(cpu.SREG_I, true)

	// Step 4: Vector 40 should be last
	err = m.Step()
	if err != nil { t.Fatal(err) }
	if m.CPU.PC != 40*2 {
		t.Errorf("Expected PC to be at vector 40, got %d", m.CPU.PC)
	}
}
