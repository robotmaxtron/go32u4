package mcu_test

import (
	"go32u4/pkg/cpu"
	"go32u4/pkg/mcu"
	"testing"
)

func TestMCUMemoryMapping(t *testing.T) {
	m := mcu.NewATmega32u4()

	// Register mapping (0-31)
	m.WriteSRAM(0, 0xAA)
	if m.CPU.Reg[0] != 0xAA {
		t.Errorf("Expected Reg[0] 0xAA, got %02X", m.CPU.Reg[0])
	}

	// IO Mapping (0x20-0x5F, which is 32-95 in SRAM address)
	// SREG is at IO 0x3F (32 + 0x3F = 95)
	m.WriteSRAM(95, 0x55)
	if m.CPU.SREG != 0x55 {
		t.Errorf("Expected SREG 0x55, got %02X", m.CPU.SREG)
	}

	// SRAM Mapping
	// SRAM starts after IO registers (32 + 256 = 288)
	m.WriteSRAM(288, 0xBE)
	if m.SRAMData[0] != 0xBE {
		t.Errorf("Expected SRAMData[0] 0xBE, got %02X", m.SRAMData[0])
	}
}

func TestMCUInterrupts(t *testing.T) {
	m := mcu.NewATmega32u4()
	m.GlobalInterrupts = true
	m.PendingInterrupts = (1 << 1) // INT0

	_ = m.Step()
	// Should have executed interrupt 1: PC = 1 * 2 = 2
	if m.CPU.PC != 2 {
		t.Errorf("Expected PC 2 after interrupt, got %d", m.CPU.PC)
	}
	if (m.PendingInterrupts & (1 << 1)) != 0 {
		t.Error("Expected INT0 to be cleared")
	}
	if m.CPU.GetFlag(cpu.SREG_I) {
		t.Error("Expected Global Interrupts to be disabled during execution")
	}
}
