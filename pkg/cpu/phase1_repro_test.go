package cpu_test

import (
	"go32u4/pkg/cpu"
	"go32u4/pkg/mcu"
	"testing"
)

func TestLDD_STD_Y(t *testing.T) {
	m := mcu.NewATmega32u4()
	m.CPU.Reg[28] = 0x00 // YL
	m.CPU.Reg[29] = 0x01 // YH
	m.CPU.Reg[16] = 0x55

	// STD Y+1, R16 (10q0 001r rrrr 1yyy)
	// r=16 (10000), q=1, y=1
	// 10 0 0 00 1 10000 1 001 = 0x8309
	m.FlashData[0] = 0x8309
	// LDD R17, Y+1 (10q0 000d dddd 1yyy)
	// d=17 (10001), q=1, y=1
	// 10 0 0 00 0 10001 1 001 = 0x8119
	m.FlashData[1] = 0x8119

	_ = m.Step()
	val := m.CPU.ReadSRAM(0x0101)
	if val != 0x55 {
		t.Errorf("Expected SRAM[0x0101] 0x55, got %02X (STD Y+1 failed)", val)
	}

	_ = m.Step()
	if m.CPU.Reg[17] != 0x55 {
		t.Errorf("Expected R17 0x55, got %02X (LDD Y+1 failed)", m.CPU.Reg[17])
	}
}

func TestLD_ST_Y_PreDec_PostInc(t *testing.T) {
	m := mcu.NewATmega32u4()
	m.CPU.Reg[28] = 0x00
	m.CPU.Reg[29] = 0x01
	m.CPU.Reg[16] = 0xAA

	// ST Y+, R16 (1001 001r rrrr 1001)
	// r=16 (10000)
	// 1001 001 10000 1001 = 0x9309
	m.FlashData[0] = 0x9309
	// LD R17, -Y (1001 000d dddd 1010)
	// d=17 (10001)
	// 1001 000 10001 1010 = 0x911A
	m.FlashData[1] = 0x911A

	_ = m.Step()
	val := m.CPU.ReadSRAM(0x0100)
	if val != 0xAA {
		t.Errorf("Expected SRAM[0x0100] 0xAA, got %02X (ST Y+ failed)", val)
	}
	y := (uint16(m.CPU.Reg[29]) << 8) | uint16(m.CPU.Reg[28])
	if y != 0x0101 {
		t.Errorf("Expected Y 0x0101, got %04X", y)
	}

	_ = m.Step()
	if m.CPU.Reg[17] != 0xAA {
		t.Errorf("Expected R17 0xAA, got %02X (LD -Y failed)", m.CPU.Reg[17])
	}
	y = (uint16(m.CPU.Reg[29]) << 8) | uint16(m.CPU.Reg[28])
	if y != 0x0100 {
		t.Errorf("Expected Y 0x0100, got %04X", y)
	}
}

func TestBST_BLD(t *testing.T) {
	m := mcu.NewATmega32u4()
	m.CPU.Reg[16] = 0x80
	// BST R16, 7 (1111 101r rrrr 0bbb)
	// r=16 (10000), b=7
	// 1111 101 10000 0 111 = 0xFB07
	m.FlashData[0] = 0xFB07
	// BLD R17, 0 (1111 100d dddd 0bbb)
	// d=17 (10001), b=0
	// 1111 100 10001 0 000 = 0xF910
	m.FlashData[1] = 0xF910

	_ = m.Step()
	if !m.CPU.GetFlag(cpu.SREG_T) {
		t.Errorf("Expected T flag to be set (BST R16, 7 failed)")
	}

	_ = m.Step()
	if m.CPU.Reg[17] != 0x01 {
		t.Errorf("Expected R17 0x01, got %02X (BLD R17, 0 failed)", m.CPU.Reg[17])
	}
}

func TestConditionalBranches(t *testing.T) {
	m := mcu.NewATmega32u4()
	m.CPU.SetFlag(cpu.SREG_Z, true)
	// BREQ +2 (1111 00 kkkkkk k 001)
	// k=2 (0000010)
	// 1111 00 0000010 001 = 0xF011
	m.FlashData[0] = 0xF011
	m.FlashData[1] = 0x0000
	m.FlashData[2] = 0x0000
	m.FlashData[3] = 0x0000

	_ = m.Step()
	if m.CPU.PC != 3 {
		t.Errorf("Expected PC 3 (BREQ taken), got %d", m.CPU.PC)
	}

	m.CPU.PC = 4
	m.CPU.SetFlag(cpu.SREG_Z, false)
	// BRNE +2 (1111 01 kkkkkk k 001)
	// 1111 01 0000010 001 = 0xF411
	m.FlashData[4] = 0xF411

	_ = m.Step()
	if m.CPU.PC != 7 {
		t.Errorf("Expected PC 7 (BRNE taken), got %d", m.CPU.PC)
	}
}

func TestSBR_CBR(t *testing.T) {
	m := mcu.NewATmega32u4()
	m.CPU.Reg[16] = 0x55
	// SBR R16, 0xAA (ORI R16, 0xAA)
	// 0110 1010 0000 1010 = 0x6A0A
	m.FlashData[0] = 0x6A0A
	// CBR R16, 0x0F (ANDI R16, 0xF0)
	// 0111 1111 0000 0000 = 0x7F00
	m.FlashData[1] = 0x7F00

	_ = m.Step()
	if m.CPU.Reg[16] != 0xFF {
		t.Errorf("Expected R16 0xFF after SBR, got %02X", m.CPU.Reg[16])
	}

	_ = m.Step()
	if m.CPU.Reg[16] != 0xF0 {
		t.Errorf("Expected R16 0xF0 after CBR, got %02X", m.CPU.Reg[16])
	}
}
