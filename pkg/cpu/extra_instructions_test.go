package cpu_test

import (
	"go32u4/pkg/cpu"
	"go32u4/pkg/mcu"
	"testing"
)

func TestMULFlags(t *testing.T) {
	m := mcu.NewATmega32u4()
	m.CPU.Reg[16] = 0xFF
	m.CPU.Reg[17] = 0xFF
	m.FlashData[0] = 0x9F01 // MUL R16, R16 -> d=16, r=16
	_ = m.Step()
	if !m.CPU.GetFlag(cpu.SREG_C) {
		t.Error("Expected Carry flag set for MUL 0xFF * 0xFF")
	}
	if m.CPU.GetFlag(cpu.SREG_Z) {
		t.Error("Expected Zero flag clear for MUL 0xFF * 0xFF")
	}
}

func TestCPSE_2Word(t *testing.T) {
	m := mcu.NewATmega32u4()
	m.CPU.Reg[16] = 0x55
	m.CPU.Reg[17] = 0x55
	m.FlashData[0] = 0x1301 // CPSE R16, R17
	m.FlashData[1] = 0x940E // CALL (2-word)
	m.FlashData[2] = 0x0000
	m.FlashData[3] = 0x0000 // Next instruction

	_ = m.Step()
	if m.CPU.PC != 3 {
		t.Errorf("Expected PC 3 (skipped 2-word CALL), got %d", m.CPU.PC)
	}
}

func TestSkipInstructions(t *testing.T) {
	m := mcu.NewATmega32u4()
	// SBIC 0x04, 3 -> 0x9923. A=0x04 (DDRB). Bit 3.
	m.FlashData[0] = 0x9923
	m.FlashData[1] = 0x0000 // NOP
	m.WriteIO(0x04, 0x00)   // Bit 3 is 0
	_ = m.Step()
	if m.CPU.PC != 2 {
		t.Errorf("SBIC failed to skip, PC=%d", m.CPU.PC)
	}

	// SBIS 0x04, 3 -> 0x9B23
	m.CPU.PC = 2
	m.FlashData[2] = 0x9B23
	m.FlashData[3] = 0x0000
	m.WriteIO(0x04, 0x08) // Bit 3 is 1
	_ = m.Step()
	if m.CPU.PC != 4 {
		t.Errorf("SBIS failed to skip, PC=%d", m.CPU.PC)
	}

	// SBRC R16, 0 -> 0xFC00
	m.CPU.PC = 4
	m.CPU.Reg[16] = 0xFE
	m.FlashData[4] = 0xFC00
	m.FlashData[5] = 0x0000
	_ = m.Step()
	if m.CPU.PC != 6 {
		t.Errorf("SBRC failed to skip, PC=%d", m.CPU.PC)
	}

	// SBRS R16, 0 -> 0xFF00. 1111 111r rrrr 0bbb. r=16=10000. b=0.
	// 1111 111 1 0000 0 000 -> 0xFF00.
	m.CPU.PC = 6
	m.CPU.Reg[16] = 0x01
	m.FlashData[6] = 0xFF00
	m.FlashData[7] = 0x940C // JMP (2-word)
	m.FlashData[8] = 0x0000
	_ = m.Step()
	if m.CPU.PC != 9 {
		t.Errorf("SBRS failed to skip 2-word instruction, PC=%d", m.CPU.PC)
	}
}

func TestDataTransferSRAM(t *testing.T) {
	m := mcu.NewATmega32u4()
	// SRAM starts after 32 regs + 256 IO regs = 288 (0x0120)
	// LDS R16, 0x0120. Rd=16=0x10. 1001 0001 0000 0000 -> 0x9100.
	m.FlashData[0] = 0x9100
	m.FlashData[1] = 0x0120
	m.SRAMData[0] = 0x55
	_ = m.Step()
	if m.CPU.Reg[16] != 0x55 {
		t.Errorf("LDS failed, got %02X", m.CPU.Reg[16])
	}

	// STS 0x0121, R16. Rr=16=0x10. 1001 0011 0000 0000 -> 0x9300.
	m.FlashData[2] = 0x9300
	m.FlashData[3] = 0x0121
	_ = m.Step()
	if m.SRAMData[1] != 0x55 {
		t.Errorf("STS failed, got %02X", m.SRAMData[1])
	}

	// Use PORTB (0x05) for OUT test
	m.WriteIO(0x05, 0x00)
	m.CPU.Reg[17] = 0xAA
	// OUT 0x05, R17. 0xB915.
	m.FlashData[4] = 0xB915
	_ = m.Step()
	if m.ReadIO(0x05) != 0xAA {
		t.Errorf("OUT failed, got %02X", m.ReadIO(0x05))
	}

	// Use DDRB (0x04) for IN test (it's a simple read/write register)
	m.WriteIO(0x04, 0x55)
	// IN R18, 0x04. Rd=18, A=0x04.
	// 1011 0AA d dddd AAAA.
	// A=0x04 = 000100. bits 5:4 = 00. bits 3:0 = 0100.
	// Rd=18 = 10010.
	// 1011 000 1 0010 0100 -> 0xB124.
	m.FlashData[5] = 0xB124
	_ = m.Step()
	if m.CPU.Reg[18] != 0x55 {
		t.Errorf("IN failed, got %02X", m.CPU.Reg[18])
	}

	m.CPU.Reg[19] = 0x80
	// BST R19, 7. 0xFB37.
	m.FlashData[6] = 0xFB37
	_ = m.Step()
	if !m.CPU.GetFlag(cpu.SREG_T) {
		t.Error("BST failed")
	}

	m.CPU.Reg[20] = 0x00
	// BLD R20, 0. 0xF940.
	m.FlashData[7] = 0xF940
	_ = m.Step()
	if m.CPU.Reg[20] != 0x01 {
		t.Errorf("BLD failed, got %02X", m.CPU.Reg[20])
	}
}

func TestMiscInstructionsExtra(t *testing.T) {
	m := mcu.NewATmega32u4()
	m.CPU.Push(0x12)
	m.CPU.Push(0x34)
	m.FlashData[0] = 0x9518 // RETI
	_ = m.Step()
	if m.CPU.PC != 0x1234 {
		t.Errorf("RETI failed, PC=%04X", m.CPU.PC)
	}

	// ASR R16. 0x9505.
	m.CPU.PC = 0
	m.CPU.Reg[16] = 0x81
	m.FlashData[0] = 0x9505
	_ = m.Step()
	if m.CPU.Reg[16] != 0xC0 {
		t.Errorf("ASR failed, got %02X", m.CPU.Reg[16])
	}
	if !m.CPU.GetFlag(cpu.SREG_C) {
		t.Error("ASR Carry bit fail")
	}
}
