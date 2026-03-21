package cpu_test

import (
	"go32u4/pkg/cpu"
	"go32u4/pkg/mcu"
	"testing"
)

func TestADC(t *testing.T) {
	m := mcu.NewATmega32u4()
	m.CPU.Reg[16] = 0x80
	m.CPU.Reg[17] = 0x80
	m.CPU.SetFlag(cpu.SREG_C, true)
	// ADC R16, R17 (0001 11rd dddd rrrr) -> 0x1F01
	m.FlashData[0] = 0x1F01

	m.Step()
	if m.CPU.Reg[16] != 0x01 {
		t.Errorf("Expected R16 0x01, got %02X", m.CPU.Reg[16])
	}
	if !m.CPU.GetFlag(cpu.SREG_C) {
		t.Error("Expected Carry set")
	}
}

func TestSBC(t *testing.T) {
	m := mcu.NewATmega32u4()
	m.CPU.Reg[16] = 0x10
	m.CPU.Reg[17] = 0x05
	m.CPU.SetFlag(cpu.SREG_C, true)
	// SBC R16, R17 (0000 10rd dddd rrrr) -> 0x0B01
	m.FlashData[0] = 0x0B01

	m.Step()
	if m.CPU.Reg[16] != 0x0A {
		t.Errorf("Expected R16 0x0A, got %02X", m.CPU.Reg[16])
	}
}

func TestCPC(t *testing.T) {
	m := mcu.NewATmega32u4()
	m.CPU.Reg[16] = 0x05
	m.CPU.Reg[17] = 0x05
	m.CPU.SetFlag(cpu.SREG_Z, true)
	m.CPU.SetFlag(cpu.SREG_C, false)
	// CPC R16, R17 (0000 01rd dddd rrrr) -> 0x0701
	m.FlashData[0] = 0x0701
	m.Step()
	if !m.CPU.GetFlag(cpu.SREG_Z) {
		t.Error("Expected Z flag to remain set")
	}

	m.CPU.PC = 0
	m.CPU.Reg[16] = 0x06
	m.CPU.SetFlag(cpu.SREG_Z, true)
	m.Step()
	if m.CPU.GetFlag(cpu.SREG_Z) {
		t.Error("Expected Z flag to be cleared")
	}
}

func TestSBI_CBI(t *testing.T) {
	m := mcu.NewATmega32u4()
	// Using a simple IO register EIFR (0x1C)
	// SBI 0x1C, 3 (1001 1010 AAAA Abbb) -> 0x9AE3
	m.FlashData[0] = 0x9AE3
	// CBI 0x1C, 3 (1001 1000 AAAA Abbb) -> 0x98E3
	m.FlashData[1] = 0x98E3

	m.Step()
	if (m.CPU.ReadIO(0x1C) & 0x08) == 0 {
		t.Errorf("Expected bit 3 set, got %02X", m.CPU.ReadIO(0x1C))
	}

	m.Step()
	if (m.CPU.ReadIO(0x1C) & 0x08) != 0 {
		t.Errorf("Expected bit 3 cleared, got %02X", m.CPU.ReadIO(0x1C))
	}
}

func TestSBIC_SBIS(t *testing.T) {
	m := mcu.NewATmega32u4()
	// SBIC 0x1C, 3 -> 0x99E3
	m.FlashData[0] = 0x99E3
	m.FlashData[1] = 0x0000 // NOP (skipped)
	m.FlashData[2] = 0x0000 // NOP

	m.CPU.WriteIO(0x1C, 0x00) // Bit 3 is clear
	m.Step()
	// SBIC at PC=0 should skip PC=1. PC should be 2.
	if m.CPU.PC != 2 {
		t.Errorf("Expected PC 2, got %d", m.CPU.PC)
	}

	m.CPU.PC = 0
	// SBIS 0x1C, 3 -> 0x9BE3
	m.FlashData[0] = 0x9BE3
	m.CPU.WriteIO(0x1C, 0x08) // Bit 3 is set
	m.Step()
	// SBIS at PC=0 should skip PC=1. PC should be 2.
	if m.CPU.PC != 2 {
		t.Errorf("Expected PC 2, got %d", m.CPU.PC)
	}
}

func TestCPSE(t *testing.T) {
	m := mcu.NewATmega32u4()
	// CPSE R16, R17 (0001 00rd dddd rrrr)
	// r=17 (r_high=1, r_low=1), d=16
	// Opcode = 0001 00 1 10000 0001 = 0001 0011 0000 0001 = 0x1301
	m.FlashData[0] = 0x1301
	m.FlashData[1] = 0x0000 // NOP
	m.CPU.Reg[16] = 0x55
	m.CPU.Reg[17] = 0x55
	m.Step()
	if m.CPU.PC != 2 {
		t.Errorf("Expected PC 2, got %d", m.CPU.PC)
	}
}

func TestShifts(t *testing.T) {
	m := mcu.NewATmega32u4()
	m.CPU.Reg[16] = 0x03
	m.FlashData[0] = 0x9506 // LSR R16
	m.Step()
	if m.CPU.Reg[16] != 0x01 {
		t.Errorf("Expected 0x01, got %02X", m.CPU.Reg[16])
	}
	if !m.CPU.GetFlag(cpu.SREG_C) {
		t.Error("Expected Carry set")
	}
}

func TestMUL(t *testing.T) {
	m := mcu.NewATmega32u4()
	m.CPU.Reg[16] = 0x0A
	m.CPU.Reg[17] = 0x0B
	// MUL R16, R17 (1001 11rd dddd rrrr)
	// d = 16 (10000), r = 17 (10001)
	// r_low = 1, r_high = 1
	// rd dddd rrrr = 1 0000 1 0001 = 0x211
	// Opcode = 0x9C00 | 0x111 = 0x9D11  Wait, rd dddd rrrr.
	// r = bits 9, 3-0. d = bits 8-4.
	// For r=17 (10001), bit 9 is 1, bits 3-0 are 0001.
	// For d=16 (10000), bits 8-4 are 10000.
	// Opcode = 1001 11 1 10000 0001 = 1001 1111 0000 0001 = 0x9F01
	m.FlashData[0] = 0x9F01
	m.Step()
	if m.CPU.Reg[0] != 110 {
		t.Errorf("Expected R0 110, got %d", m.CPU.Reg[0])
	}
}
