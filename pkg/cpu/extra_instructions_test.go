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
	m.FlashData[0] = 0x9F01
	_ = m.Step()
	if !m.CPU.GetFlag(cpu.SREG_C) {
		t.Error("Expected Carry flag set for MUL 0xFF * 0xFF")
	}
	if m.CPU.GetFlag(cpu.SREG_Z) {
		t.Error("Expected Zero flag clear for MUL 0xFF * 0xFF")
	}

	m.CPU.PC = 1
	m.CPU.Reg[16] = 0
	m.CPU.Reg[17] = 0
	m.FlashData[1] = 0x9F01
	_ = m.Step()
	if !m.CPU.GetFlag(cpu.SREG_Z) {
		t.Error("Expected Zero flag set for MUL 0 * 0")
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
	// Step() fetches 0x1301 at 0, PC=1. CPSE skips CALL (2 words). PC becomes 1 + 2 = 3.
	if m.CPU.PC != 3 {
		t.Errorf("Expected PC 3 (skipped 2-word CALL), got %d", m.CPU.PC)
	}
}

func TestSkipInstructions(t *testing.T) {
	m := mcu.NewATmega32u4()
	// SBIC 0x1C, 3 -> 0x99E3
	m.FlashData[0] = 0x99E3
	m.FlashData[1] = 0x0000 // NOP
	m.WriteIO(0x1C, 0x00)   // Bit 3 is 0
	_ = m.Step()
	if m.CPU.PC != 2 {
		t.Errorf("SBIC failed to skip, PC=%d", m.CPU.PC)
	}

	// SBIS 0x1C, 3 -> 0x9BE3
	m.CPU.PC = 2
	m.FlashData[2] = 0x9BE3
	m.FlashData[3] = 0x0000
	m.WriteIO(0x1C, 0x08) // Bit 3 is 1
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

	// SBRS R16, 0 -> 0xFE00
	m.CPU.PC = 6
	m.CPU.Reg[16] = 0x01
	m.FlashData[6] = 0xFE00
	m.FlashData[7] = 0x940C // JMP (2-word)
	m.FlashData[8] = 0x0000
	_ = m.Step()
	// pc was 6, step fetches 0xFE00, pc=7. SBRS skips JMP (2 words). PC becomes 7+2=9.
	// Oh, I see: instructions.go says "c.PC++" for 2-word skip. 
	// If PC=7, c.PC++ -> 8. Then checks if 2-word. If yes, c.PC++ -> 9. 
	// Wait, it says "nextOp := flash[c.PC]; c.PC++". So it reads at PC, then increments.
	// In my case PC=7, it reads at 7 (0x940C), then c.PC becomes 8. Then it checks 0x940C & ... == 0x940C. 
	// Yes! So it should work. PC should be 9. 
	if m.CPU.PC != 9 {
		t.Errorf("SBRS failed to skip 2-word instruction, PC=%d", m.CPU.PC)
	}
}

func TestALUImmediates(t *testing.T) {
	m := mcu.NewATmega32u4()
	m.CPU.Reg[16] = 0xFF
	m.FlashData[0] = 0x700F // ANDI R16, 0x0F
	_ = m.Step()
	if m.CPU.Reg[16] != 0x0F {
		t.Errorf("ANDI failed, got %02X", m.CPU.Reg[16])
	}

	m.FlashData[1] = 0x6F00 // ORI R16, 0xF0
	_ = m.Step()
	if m.CPU.Reg[16] != 0xFF {
		t.Errorf("ORI failed, got %02X", m.CPU.Reg[16])
	}

	m.FlashData[2] = 0x3F0F // CPI R16, 0x0F
	_ = m.Step()
	if !m.CPU.GetFlag(cpu.SREG_Z) {
		t.Error("CPI failed to set Zero flag")
	}
}

func TestALURegReg(t *testing.T) {
	m := mcu.NewATmega32u4()
	m.CPU.Reg[16] = 0x80
	m.CPU.Reg[17] = 0x80
	m.CPU.SetFlag(cpu.SREG_C, true)
	m.FlashData[0] = 0x1F01 // ADC R16, R17
	_ = m.Step()
	if m.CPU.Reg[16] != 0x01 {
		t.Errorf("ADC failed, got %02X", m.CPU.Reg[16])
	}
	if !m.CPU.GetFlag(cpu.SREG_C) {
		t.Error("ADC failed to set Carry flag")
	}

	m.CPU.Reg[16] = 0x10
	m.CPU.Reg[17] = 0x05
	m.CPU.SetFlag(cpu.SREG_C, true)
	m.CPU.SetFlag(cpu.SREG_Z, true)
	m.FlashData[1] = 0x0B01 // SBC R16, R17
	_ = m.Step()
	if m.CPU.Reg[16] != 0x0A {
		t.Errorf("SBC failed, got %02X", m.CPU.Reg[16])
	}
	if m.CPU.GetFlag(cpu.SREG_Z) {
		t.Error("SBC should have cleared Zero flag")
	}

	m.CPU.Reg[18] = 0x05
	m.CPU.Reg[19] = 0x05
	m.CPU.SetFlag(cpu.SREG_C, false)
	m.CPU.SetFlag(cpu.SREG_Z, true)
	m.FlashData[2] = 0x0723 // CPC R18, R19
	_ = m.Step()
	if !m.CPU.GetFlag(cpu.SREG_Z) {
		t.Error("CPC should have kept Zero flag true")
	}
}

func TestDataTransferSRAM(t *testing.T) {
	m := mcu.NewATmega32u4()
	// SRAM starts at 256. (0x0100)
	m.FlashData[0] = 0x9000
	m.FlashData[1] = 0x0100 // LDS R16, 0x0100
	m.SRAMData[0] = 0x55
	_ = m.Step()
	if m.CPU.Reg[16] != 0x55 {
		t.Errorf("LDS failed, got %02X", m.CPU.Reg[16])
	}

	m.FlashData[2] = 0x9200
	m.FlashData[3] = 0x0101 // STS 0x0101, R16
	_ = m.Step()
	if m.SRAMData[1] != 0x55 {
		t.Errorf("STS failed, got %02X", m.SRAMData[1])
	}

	m.WriteIO(0x05, 0x00)
	m.CPU.Reg[17] = 0xAA
	m.FlashData[4] = 0xB915 // OUT 0x05, R17
	_ = m.Step()
	if m.ReadIO(0x05) != 0xAA {
		t.Errorf("OUT failed, got %02X", m.ReadIO(0x05))
	}

	m.WriteIO(0x03, 0x55)
	m.FlashData[5] = 0xB123 // IN R18, 0x03
	_ = m.Step()
	if m.CPU.Reg[18] != 0x55 {
		t.Errorf("IN failed, got %02X", m.CPU.Reg[18])
	}

	m.CPU.Reg[19] = 0x80
	m.FlashData[6] = 0xFB37 // BST R19, 7
	_ = m.Step()
	if !m.CPU.GetFlag(cpu.SREG_T) {
		t.Error("BST failed")
	}

	m.CPU.Reg[20] = 0x00
	m.FlashData[7] = 0xF940 // BLD R20, 0
	_ = m.Step()
	if m.CPU.Reg[20] != 0x01 {
		t.Errorf("BLD failed, got %02X", m.CPU.Reg[20])
	}
}

func TestMiscInstructionsExtra(t *testing.T) {
	m := mcu.NewATmega32u4()
	// RETI. pch,pcl := Pop(), Pop() -> High then Low.
	// To get PC=0x1234: Push 0x34 (Low), then 0x12 (High).
	m.CPU.Push(0x34)
	m.CPU.Push(0x12)
	m.FlashData[0] = 0x9518
	_ = m.Step()
	if m.CPU.PC != 0x1234 {
		t.Errorf("RETI failed, PC=%04X", m.CPU.PC)
	}

	// ASR R16. Opcode: 1001 010d dddd 0101 -> 0x9405
	m.CPU.PC = 0
	m.CPU.Reg[16] = 0x81
	m.FlashData[0] = 0x9405
	_ = m.Step()
	if m.CPU.Reg[16] != 0xC0 {
		t.Errorf("ASR failed, got %02X", m.CPU.Reg[16])
	}
	if !m.CPU.GetFlag(cpu.SREG_C) {
		t.Error("ASR Carry bit fail")
	}
}
