package cpu_test

import (
	"go32u4/pkg/cpu"
	"go32u4/pkg/mcu"
	"go32u4/pkg/peripherals"
	"testing"
)

func TestCPUInitialization(t *testing.T) {
	m := mcu.NewATmega32u4()
	if m.CPU.PC != 0 {
		t.Errorf("Expected PC 0, got %d", m.CPU.PC)
	}
	// 64 (IO) + 160 (ExtIO) + 2560 (SRAM) - 1 = 2783 = 0xADF
	if m.CPU.SP != 0x0ADF {
		t.Errorf("Expected SP 0x0ADF, got %04X", m.CPU.SP)
	}
}

func TestLDI(t *testing.T) {
	m := mcu.NewATmega32u4()
	m.FlashData[0] = 0xE00F // LDI R16, 0x0F
	m.FlashData[1] = 0xEF1F // LDI R17, 0xFF

	_ = m.Step()
	if m.CPU.Reg[16] != 0x0F {
		t.Errorf("Expected R16 0x0F, got %02X", m.CPU.Reg[16])
	}

	_ = m.Step()
	if m.CPU.Reg[17] != 0xFF {
		t.Errorf("Expected R17 0xFF, got %02X", m.CPU.Reg[17])
	}
}

func TestADD(t *testing.T) {
	m := mcu.NewATmega32u4()
	m.CPU.Reg[16] = 0x10
	m.CPU.Reg[17] = 0x20
	// 0000 11rd dddd rrrr
	// Rd = 16 (10000)
	// Rr = 17 (10001)
	// r = (opcode & 0x000F) | ((opcode >> 5) & 0x0010)
	// d = (opcode >> 4) & 0x001F
	// For d=16, (d<<4) = 0x0100
	// For r=17, (r&0x0F) = 1, (r&0x10)<<5 = 0x0200
	// Opcode = 0x0C00 | 0x0100 | 0x0200 | 0x0001 = 0x0F01
	m.FlashData[0] = 0x0F01

	_ = m.Step()
	if m.CPU.Reg[16] != 0x30 {
		t.Errorf("Expected R16 0x30, got %02X", m.CPU.Reg[16])
	}
}

func TestTimer0Overflow(t *testing.T) {
	m := mcu.NewATmega32u4()
	m.WriteIO(peripherals.TCCR0B, 1) // Prescaler 1
	m.WriteIO(peripherals.TIMSK0, 1) // TOIE0
	m.CPU.SetFlag(cpu.SREG_I, true)
	m.GlobalInterrupts = true

	m.Periph.Timer0Counter = 255
	m.FlashData[0] = 0x0000 // NOP

	_ = m.Step()
	// After NOP, Timer0 should overflow
	if m.Periph.Timer0Counter != 0 {
		t.Errorf("Expected Timer0 0, got %d", m.Periph.Timer0Counter)
	}
	// Interrupt should be pending (vector 24 is index 23)
	if (m.PendingInterrupts & (1 << 23)) == 0 {
		t.Errorf("Expected Timer0 Overflow interrupt pending, got %b", m.PendingInterrupts)
	}

	// Next step should execute interrupt
	_ = m.Step()
	// ATmega32u4 Timer0 Overflow vector is 24 (address (24-1)*2 = 46)
	if m.CPU.PC != 23*2 {
		t.Errorf("Expected PC %d, got %d", 23*2, m.CPU.PC)
	}
}

func TestSUB(t *testing.T) {
	m := mcu.NewATmega32u4()
	m.CPU.Reg[16] = 0x20
	m.CPU.Reg[17] = 0x10
	// SUB R16, R17 (0001 10rd dddd rrrr)
	// d=16 (10000), r=17 (10001)
	// Mask: 0xFC00 == 0x1800
	// d: (0x1921 >> 4) & 0x1F = 0x12 (18) -> WRONG.
	// We want d=16 (10000), r=17 (10001)
	// opcode: 0001 10 1 10000 10001
	//          0001 1011 0000 1001
	//          1    B    0    9
	m.FlashData[0] = 0x1B01 // Correction: 0x1800 | (1<<4 bit r) | (16<<4 bit d) | 1 (low 4 bits of r)
	// r = (opcode & 0x0F) | ((opcode >> 5) & 0x10)
	// d = (opcode >> 4) & 0x1F
	// For d=16, (16<<4) = 0x0100
	// For r=17, (17&0x0F)=1, (17&0x10)<<5 = 0x0200
	// Opcode = 0x1800 | 0x0100 | 0x0200 | 1 = 0x1B01

	_ = m.Step()
	if m.CPU.Reg[16] != 0x10 {
		t.Errorf("Expected R16 0x10, got %02X", m.CPU.Reg[16])
	}
	if m.CPU.GetFlag(cpu.SREG_Z) {
		t.Error("Zero flag should be false")
	}
}

func TestRJMP(t *testing.T) {
	m := mcu.NewATmega32u4()
	// RJMP +2 (1100 0000 0000 0010)
	m.FlashData[0] = 0xC002

	_ = m.Step()
	// PC started at 0, after opcode fetch PC=1. Then RJMP +2 -> PC=3
	if m.CPU.PC != 3 {
		t.Errorf("Expected PC 3, got %d", m.CPU.PC)
	}
}

func TestStack(t *testing.T) {
	m := mcu.NewATmega32u4()
	m.CPU.Reg[16] = 0xAA
	// PUSH R16 (1001 0011 0000 1111) -> 0x930F
	m.FlashData[0] = 0x930F
	// POP R17 (1001 0001 0001 1111) -> 0x911F
	m.FlashData[1] = 0x911F

	_ = m.Step()
	if m.CPU.SP != 0x0ADE {
		t.Errorf("Expected SP 0x0ADE, got %04X", m.CPU.SP)
	}

	_ = m.Step()
	if m.CPU.Reg[17] != 0xAA {
		t.Errorf("Expected R17 0xAA, got %02X", m.CPU.Reg[17])
	}
	if m.CPU.SP != 0x0ADF {
		t.Errorf("Expected SP 0x0ADF, got %04X", m.CPU.SP)
	}
}
