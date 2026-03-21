package peripherals_test

import (
	"go32u4/pkg/peripherals"
	"testing"
)

// TickMock advances cycles and calls Manager.Tick
func TickMock(p *peripherals.Manager, sys *MockSystem, cycles uint64) {
	for i := uint64(0); i < cycles; i++ {
		sys.cycles++
		p.Tick(1)
	}
}

func TestTimer1(t *testing.T) {
	sys := &MockSystem{ioRegs: make([]uint8, 256)}
	p := peripherals.NewManager(sys)

	// Interrupt vectors for Timer 1:
	// Capture: 16
	// CompA: 17
	// CompB: 18
	// CompC: 19
	// Overflow: 20

	// 1. Test Prescaler and Counter Increment
	// Set Prescaler to 8 (CS12:0 = 010)
	p.Timer1ControlB = 2
	sys.cycles = 0
	p.Timer1Counter = 100
	
	// Tick 8 cycles, should increment once
	// Condition: (sys.Cycles()) % 8 == 0
	// sys.Cycles() will go 1, 2, 3, 4, 5, 6, 7, 8.
	// At cycle 8, 8%8 == 0, so it increments.
	TickMock(p, sys, 8)
	if p.Timer1Counter != 101 {
		t.Errorf("Expected Timer1Counter 101, got %d", p.Timer1Counter)
	}

	// 2. Test Compare Match A
	p.Timer1CompareA = 105
	sys.ioRegs[peripherals.TIMSK1] = 1 << 1 // OCIE1A
	sys.ints = 0
	
	// Increment to 105 (4 more increments * 8 cycles = 32 cycles)
	TickMock(p, sys, 32)
	if p.Timer1Counter != 105 {
		t.Errorf("Expected Timer1Counter 105, got %d", p.Timer1Counter)
	}
	if (sys.ints & (1 << 17)) == 0 {
		t.Error("Expected Timer1 Compare Match A interrupt")
	}
	if (sys.ioRegs[peripherals.TIFR1] & (1 << 1)) == 0 {
		t.Error("Expected TIFR1 OCF1A flag set")
	}

	// 3. Test Compare Match B
	p.Timer1CompareB = 110
	sys.ioRegs[peripherals.TIMSK1] |= 1 << 2 // OCIE1B
	sys.ints = 0
	TickMock(p, sys, 40)
	if p.Timer1Counter != 110 {
		t.Errorf("Expected Timer1Counter 110, got %d", p.Timer1Counter)
	}
	if (sys.ints & (1 << 18)) == 0 {
		t.Error("Expected Timer1 Compare Match B interrupt")
	}
	if (sys.ioRegs[peripherals.TIFR1] & (1 << 2)) == 0 {
		t.Error("Expected TIFR1 OCF1B flag set")
	}

	// 4. Test Compare Match C
	p.Timer1CompareC = 115
	sys.ioRegs[peripherals.TIMSK1] |= 1 << 3 // OCIE1C
	sys.ints = 0
	TickMock(p, sys, 40)
	if p.Timer1Counter != 115 {
		t.Errorf("Expected Timer1Counter 115, got %d", p.Timer1Counter)
	}
	if (sys.ints & (1 << 19)) == 0 {
		t.Error("Expected Timer1 Compare Match C interrupt")
	}
	if (sys.ioRegs[peripherals.TIFR1] & (1 << 3)) == 0 {
		t.Error("Expected TIFR1 OCF1C flag set")
	}

	// 5. Test Overflow
	p.Timer1Counter = 0xFFFF
	sys.ioRegs[peripherals.TIMSK1] |= 1 << 0 // TOIE1
	sys.ints = 0
	TickMock(p, sys, 8)
	if p.Timer1Counter != 0 {
		t.Errorf("Expected Timer1Counter 0 after overflow, got %d", p.Timer1Counter)
	}
	if (sys.ints & (1 << 20)) == 0 {
		t.Error("Expected Timer1 Overflow interrupt")
	}
	if (sys.ioRegs[peripherals.TIFR1] & (1 << 0)) == 0 {
		t.Error("Expected TIFR1 TOV1 flag set")
	}

	// 6. Test other prescalers for coverage
	// Prescaler 64
	p.Timer1ControlB = 3
	p.Timer1Counter = 0
	sys.cycles = 0 // Reset cycles for easier alignment
	TickMock(p, sys, 64)
	if p.Timer1Counter != 1 {
		t.Errorf("Prescaler 64: Expected 1, got %d", p.Timer1Counter)
	}

	// Prescaler 256
	p.Timer1ControlB = 4
	p.Timer1Counter = 0
	sys.cycles = 0
	TickMock(p, sys, 256)
	if p.Timer1Counter != 1 {
		t.Errorf("Prescaler 256: Expected 1, got %d", p.Timer1Counter)
	}

	// Prescaler 1024
	p.Timer1ControlB = 5
	p.Timer1Counter = 0
	sys.cycles = 0
	TickMock(p, sys, 1024)
	if p.Timer1Counter != 1 {
		t.Errorf("Prescaler 1024: Expected 1, got %d", p.Timer1Counter)
	}

	// Prescaler 0 (Stopped)
	p.Timer1ControlB = 0
	p.Timer1Counter = 10
	TickMock(p, sys, 1024)
	if p.Timer1Counter != 10 {
		t.Errorf("Prescaler 0: Expected 10, got %d", p.Timer1Counter)
	}
	
	// Prescaler 1 (No prescaling)
	p.Timer1ControlB = 1
	p.Timer1Counter = 0
	sys.cycles = 0
	TickMock(p, sys, 1)
	if p.Timer1Counter != 1 {
		t.Errorf("Prescaler 1: Expected 1, got %d", p.Timer1Counter)
	}
}

func TestTimer3(t *testing.T) {
	sys := &MockSystem{ioRegs: make([]uint8, 256)}
	p := peripherals.NewManager(sys)

	// Interrupt vectors for Timer 3:
	// Capture: 31
	// CompA: 32
	// CompB: 33
	// CompC: 34
	// Overflow: 35

	// 1. Test Prescaler and Counter Increment
	p.Timer3ControlB = 2 // /8
	sys.cycles = 0
	p.Timer3Counter = 200
	TickMock(p, sys, 8)
	if p.Timer3Counter != 201 {
		t.Errorf("Expected Timer3Counter 201, got %d", p.Timer3Counter)
	}

	// 2. Test Compare Match A
	p.Timer3CompareA = 205
	sys.ioRegs[peripherals.TIMSK3] = 1 << 1 // OCIE3A
	sys.ints = 0
	TickMock(p, sys, 32)
	if (sys.ints & (1 << 32)) == 0 {
		t.Error("Expected Timer3 Compare Match A interrupt")
	}

	// 3. Test Compare Match B
	p.Timer3CompareB = 210
	sys.ioRegs[peripherals.TIMSK3] |= 1 << 2 // OCIE3B
	sys.ints = 0
	TickMock(p, sys, 40)
	if (sys.ints & (1 << 33)) == 0 {
		t.Error("Expected Timer3 Compare Match B interrupt")
	}

	// 4. Test Compare Match C
	p.Timer3CompareC = 215
	sys.ioRegs[peripherals.TIMSK3] |= 1 << 3 // OCIE3C
	sys.ints = 0
	TickMock(p, sys, 40)
	if (sys.ints & (1 << 34)) == 0 {
		t.Error("Expected Timer3 Compare Match C interrupt")
	}

	// 5. Test Overflow
	p.Timer3Counter = 0xFFFF
	sys.ioRegs[peripherals.TIMSK3] |= 1 << 0 // TOIE3
	sys.ints = 0
	TickMock(p, sys, 8)
	if p.Timer3Counter != 0 {
		t.Errorf("Expected Timer3Counter 0 after overflow, got %d", p.Timer3Counter)
	}
	if (sys.ints & (1 << 35)) == 0 {
		t.Error("Expected Timer3 Overflow interrupt")
	}
	
	// 6. Test other prescalers for coverage
	p.Timer3ControlB = 3 // /64
	p.Timer3Counter = 0
	sys.cycles = 0
	TickMock(p, sys, 64)
	if p.Timer3Counter != 1 {
		t.Errorf("Timer3 Prescaler 64: Expected 1, got %d", p.Timer3Counter)
	}

	p.Timer3ControlB = 4 // /256
	p.Timer3Counter = 0
	sys.cycles = 0
	TickMock(p, sys, 256)
	if p.Timer3Counter != 1 {
		t.Errorf("Timer3 Prescaler 256: Expected 1, got %d", p.Timer3Counter)
	}

	p.Timer3ControlB = 5 // /1024
	p.Timer3Counter = 0
	sys.cycles = 0
	TickMock(p, sys, 1024)
	if p.Timer3Counter != 1 {
		t.Errorf("Timer3 Prescaler 1024: Expected 1, got %d", p.Timer3Counter)
	}

	p.Timer3ControlB = 0 // Stopped
	p.Timer3Counter = 10
	TickMock(p, sys, 1024)
	if p.Timer3Counter != 10 {
		t.Errorf("Timer3 Prescaler 0: Expected 10, got %d", p.Timer3Counter)
	}
}
