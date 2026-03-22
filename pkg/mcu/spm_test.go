package mcu

import (
	"testing"
)

func TestSPM(t *testing.T) {
	m := NewATmega32u4()

	// 1. Erase Page at 0x1000 (byte address)
	// Word address = 0x800
	m.CPU.Reg[30] = 0x00 // ZL
	m.CPU.Reg[31] = 0x10 // ZH (Z = 0x1000)
	
	// Set PGERS and SPMEN
	m.WriteIO(0x37, 0x03) 
	
	// Execute SPM (0x95E8)
	m.FlashData[0] = 0x95E8
	m.CPU.PC = 0
	_ = m.Step()

	// Verify page is erased (all 0xFFFF)
	for i := 0x800; i < 0x840; i++ {
		if m.FlashData[i] != 0xFFFF {
			t.Errorf("Flash at %04X not erased: %04X", i, m.FlashData[i])
		}
	}

	// 2. Fill temporary buffer
	// Word 1: 0x1234 at index 0 (Z=0x1000)
	m.CPU.Reg[0] = 0x34 // R0
	m.CPU.Reg[1] = 0x12 // R1
	m.CPU.Reg[30] = 0x00
	m.CPU.Reg[31] = 0x10
	m.WriteIO(0x37, 0x01) // SPMEN
	m.CPU.PC = 0
	_ = m.Step()
	if m.Periph.SPMBuffer[0] != 0x1234 {
		t.Errorf("SPMBuffer[0] expected 0x1234, got %04X", m.Periph.SPMBuffer[0])
	}

	// Word 2: 0x5678 at index 1 (Z=0x1002)
	m.CPU.Reg[0] = 0x78
	m.CPU.Reg[1] = 0x56
	m.CPU.Reg[30] = 0x02
	m.CPU.Reg[31] = 0x10
	m.WriteIO(0x37, 0x01) // SPMEN
	m.CPU.PC = 0
	_ = m.Step()
	if m.Periph.SPMBuffer[1] != 0x5678 {
		t.Errorf("SPMBuffer[1] expected 0x5678, got %04X", m.Periph.SPMBuffer[1])
	}

	// 3. Write Page
	m.CPU.Reg[30] = 0x00
	m.CPU.Reg[31] = 0x10
	m.WriteIO(0x37, 0x05) // PGWRT | SPMEN
	m.CPU.PC = 0
	_ = m.Step()

	// Verify flash content
	if m.FlashData[0x800] != 0x1234 {
		t.Errorf("Flash at 0x800 expected 0x1234, got %04X", m.FlashData[0x800])
	}
	if m.FlashData[0x801] != 0x5678 {
		t.Errorf("Flash at 0x801 expected 0x5678, got %04X", m.FlashData[0x801])
	}
}
