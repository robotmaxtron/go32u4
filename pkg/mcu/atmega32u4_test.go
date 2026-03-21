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

func TestReadSRAM(t *testing.T) {
	m := mcu.NewATmega32u4()

	// Registers (0-31)
	m.CPU.Reg[0] = 0x11
	m.CPU.Reg[15] = 0x22
	m.CPU.Reg[31] = 0x33
	if val := m.ReadSRAM(0); val != 0x11 {
		t.Errorf("ReadSRAM(0) = %02X, expected 0x11", val)
	}
	if val := m.ReadSRAM(15); val != 0x22 {
		t.Errorf("ReadSRAM(15) = %02X, expected 0x22", val)
	}
	if val := m.ReadSRAM(31); val != 0x33 {
		t.Errorf("ReadSRAM(31) = %02X, expected 0x33", val)
	}

	// IO Mapping (32 to 32+255)
	// SREG is at IO 0x3F (32 + 0x3F = 95)
	m.CPU.SREG = 0x55
	if val := m.ReadSRAM(95); val != 0x55 {
		t.Errorf("ReadSRAM(95) = %02X, expected 0x55 (SREG)", val)
	}
	// Other IO register via Periph
	// TWDR is 0xBB (IO)
	m.WriteIO(0xBB, 0x44)
	if val := m.ReadSRAM(32+0xBB); val != 0x44 {
		t.Errorf("ReadSRAM(32+0xBB) = %02X, expected 0x44", val)
	}

	// SRAM Mapping (starts at 32 + 256 = 288)
	m.SRAMData[0] = 0xAA
	m.SRAMData[100] = 0xBB
	m.SRAMData[len(m.SRAMData)-1] = 0xCC
	if val := m.ReadSRAM(288); val != 0xAA {
		t.Errorf("ReadSRAM(288) = %02X, expected 0xAA", val)
	}
	if val := m.ReadSRAM(288+100); val != 0xBB {
		t.Errorf("ReadSRAM(388) = %02X, expected 0xBB", val)
	}
	if val := m.ReadSRAM(288+uint16(len(m.SRAMData))-1); val != 0xCC {
		t.Errorf("ReadSRAM(end) = %02X, expected 0xCC", val)
	}

	// Out of bounds
	if val := m.ReadSRAM(0xFFFF); val != 0 {
		t.Errorf("ReadSRAM(0xFFFF) = %02X, expected 0", val)
	}
}

func TestMCUInterrupts(t *testing.T) {
	m := mcu.NewATmega32u4()
	m.GlobalInterrupts = true
	// We use INT0 (Vector 2, index 1)
	m.PendingInterrupts = 1 << 1

	_ = m.Step()
	// Vector 2: Address (2-1) * 2 = 2
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
