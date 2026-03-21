package mcu

import (
	"go32u4/pkg/cpu"
	"testing"
)

func TestSleepWakeup(t *testing.T) {
	m := NewATmega32u4()

	// 1. SLEEP without SE bit set should NOT sleep
	m.FlashData[0] = 0x9588 // SLEEP
	m.FlashData[1] = 0x0000 // NOP
	
	err := m.Step()
	if err != nil { t.Fatal(err) }
	if m.Periph.SleepEnabled {
		t.Errorf("CPU slept even though SE bit was NOT set")
	}
	if m.CPU.PC != 1 {
		t.Errorf("Expected PC to be 1, got %d", m.CPU.PC)
	}

	// 2. SLEEP with SE bit set SHOULD sleep
	m.WriteIO(0x53, 0x01) // SMCR = SE bit set
	m.CPU.PC = 0
	err = m.Step()
	if err != nil { t.Fatal(err) }
	if !m.Periph.SleepEnabled {
		t.Errorf("CPU did NOT sleep even though SE bit was set")
	}
	// PC should be at 0 (restored after SLEEP)
	if m.CPU.PC != 0 {
		t.Errorf("Expected PC to be 0 (restored after SLEEP), got %d", m.CPU.PC)
	}

	// 3. Step while sleeping should increment cycles but NOT PC
	cyclesBefore := m.CPU.Cycles
	err = m.Step()
	if err != nil { t.Fatal(err) }
	if !m.Periph.SleepEnabled {
		t.Errorf("CPU woke up without interrupt")
	}
	if m.CPU.PC != 0 {
		t.Errorf("PC changed while sleeping: %d", m.CPU.PC)
	}
	if m.CPU.Cycles != cyclesBefore+1 {
		t.Errorf("Cycles did not increment correctly while sleeping: %d -> %d", cyclesBefore, m.CPU.Cycles)
	}

	// 4. Trigger interrupt should wake up CPU
	m.CPU.SetFlag(cpu.SREG_I, true) // Enable global interrupts
	m.TriggerInterrupt(1) // Trigger INT0 (vector 1)
	
	err = m.Step()
	if err != nil { t.Fatal(err) }
	
	if m.Periph.SleepEnabled {
		t.Errorf("CPU is still sleeping after interrupt trigger")
	}
	
	// PC should be at vector 1 * 2 = 2
	// But it should have executed the instruction at 0 (the SLEEP instruction)
	// Wait, if it wakes up it should execute the instruction AFTER sleep.
	// In AVR, when it wakes up from an interrupt, it executes the interrupt handler
	// and then returns to the instruction AFTER sleep.
	if m.CPU.PC != 2 {
		t.Errorf("Expected PC to be at interrupt vector 2, got %d", m.CPU.PC)
	}
}
