package mcu_test

import (
	"go32u4/pkg/mcu"
	"go32u4/pkg/cpu"
	"testing"
)

func TestRETReturnOrder(t *testing.T) {
	m := mcu.NewATmega32u4()
	// Target PC: 0x1234
	// Stack: [..., High(0x12), Low(0x34)] <- SP
	m.CPU.Push(0x12) // High byte first
	m.CPU.Push(0x34) // Low byte second
	
	// RET opcode: 0x9508
	m.FlashData[0] = 0x9508
	m.CPU.PC = 0
	
	err := m.CPU.Step()
	if err != nil { t.Fatal(err) }
	
	if m.CPU.PC != 0x1234 {
		t.Errorf("Expected PC 0x1234 after RET, got %04X", m.CPU.PC)
	}
}

func TestRETIReturnOrder(t *testing.T) {
	m := mcu.NewATmega32u4()
	// Target PC: 0x5678
	m.CPU.Push(0x56)
	m.CPU.Push(0x78)
	
	// RETI opcode: 0x9518
	m.FlashData[0] = 0x9518
	m.CPU.PC = 0
	
	err := m.CPU.Step()
	if err != nil { t.Fatal(err) }
	
	if m.CPU.PC != 0x5678 {
		t.Errorf("Expected PC 0x5678 after RETI, got %04X", m.CPU.PC)
	}
	if !m.CPU.GetFlag(cpu.SREG_I) {
		t.Error("RETI should set SREG_I")
	}
}

func TestWDTCSRTimedChangeUpdate(t *testing.T) {
	m := mcu.NewATmega32u4()
	// WDTCSR = 0x60
	// Enable timed change: set WDCE (bit 4) and WDE (bit 3)
	val := uint8((1 << 4) | (1 << 3))
	m.WriteIO(0x60, val)
	
	if m.ReadIO(0x60) != val {
		t.Errorf("Expected WDTCSR to be %02X after timed change enable, got %02X", val, m.ReadIO(0x60))
	}
}

func TestTWSRMasking(t *testing.T) {
	m := mcu.NewATmega32u4()
	// TWSR = 0xB9
	// Try to write to bits 3-7
	m.WriteIO(0xB9, 0xF8) // 1111 1000
	
	val := m.ReadIO(0xB9)
	if (val & 0xF8) != 0 {
		t.Errorf("Expected TWSR status bits (3-7) to be protected, got %02X", val)
	}
	
	// Try to write to bits 0-1
	m.WriteIO(0xB9, 0x03)
	val = m.ReadIO(0xB9)
	if (val & 0x03) != 0x03 {
		t.Errorf("Expected TWSR prescaler bits (0-1) to be writable, got %02X", val)
	}
}

func TestSkipCyclePenalty(t *testing.T) {
	m := mcu.NewATmega32u4()
	
	// CPSE R0, R1 (0x1001) - R0 == R1
	m.CPU.Reg[0] = 5
	m.CPU.Reg[1] = 5
	m.FlashData[0] = 0x1001
	
	// 2-word instruction: CALL 0x100 (0x940E, 0x0080)
	m.FlashData[1] = 0x940E
	m.FlashData[2] = 0x0080
	
	m.CPU.PC = 0
	m.CPU.Cycles = 0
	err := m.CPU.Step()
	if err != nil { t.Fatal(err) }
	
	if m.CPU.PC != 3 {
		t.Errorf("Expected PC 3 after skipping 2-word instruction, got %d", m.CPU.PC)
	}
	if m.CPU.Cycles != 3 {
		t.Errorf("Expected 3 cycles for skipping 2-word instruction, got %d", m.CPU.Cycles)
	}
	
	// Test skipping 1-word instruction
	m.CPU.PC = 10
	m.FlashData[10] = 0x1001 // CPSE R0, R1
	m.FlashData[11] = 0x2700 // CLR R0 (1-word)
	m.CPU.Cycles = 0
	err = m.CPU.Step()
	if err != nil { t.Fatal(err) }
	
	if m.CPU.PC != 12 {
		t.Errorf("Expected PC 12 after skipping 1-word instruction, got %d", m.CPU.PC)
	}
	if m.CPU.Cycles != 2 {
		t.Errorf("Expected 2 cycles for skipping 1-word instruction, got %d", m.CPU.Cycles)
	}
}
