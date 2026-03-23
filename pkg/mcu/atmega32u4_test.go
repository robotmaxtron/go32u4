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
	// SRAM starts at 32 + 256 = 288 (0x120)
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
	if val := m.ReadSRAM(388); val != 0xBB {
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

func TestMCUFlashOperations(t *testing.T) {
	m := mcu.NewATmega32u4()

	// FlashWrite (Standard)
	m.FlashWrite(10, 0x1234)
	if m.FlashData[10] != 0x1234 {
		t.Errorf("Expected FlashData[10] 0x1234, got %04X", m.FlashData[10])
	}

	// FlashWrite (SPM Buffer)
	// SPMCSR at 0x37. SPMEN is bit 0.
	m.WriteIO(0x37, 0x01)
	m.FlashWrite(10, 0x5678)
	if m.Periph.SPMBuffer[10%64] != 0x5678 {
		t.Errorf("Expected SPMBuffer[10] 0x5678, got %04X", m.Periph.SPMBuffer[10%64])
	}
	// Verify it didn't overwrite flash directly
	if m.FlashData[10] != 0x1234 {
		t.Errorf("Expected FlashData[10] still 0x1234, got %04X", m.FlashData[10])
	}

	// FlashErase
	m.FlashWrite(64, 0xAAAA)
	m.FlashWrite(65, 0xBBBB)
	m.FlashErase(64) // Erases page starting at 64
	if m.FlashData[64] != 0xFFFF || m.FlashData[127] != 0xFFFF {
		t.Errorf("Expected FlashData[64] and [127] 0xFFFF after erase, got %04X, %04X", m.FlashData[64], m.FlashData[127])
	}

	// FlashCommit
	for i := uint16(0); i < 64; i++ {
		m.Periph.SPMBuffer[i] = i
	}
	m.FlashCommit(128) // Commits to page starting at 128
	if m.FlashData[128] != 0 || m.FlashData[128+63] != 63 {
		t.Errorf("Expected FlashData[128] 0 and [191] 63 after commit, got %04X, %04X", m.FlashData[128], m.FlashData[191])
	}
}

func TestMCUOther(t *testing.T) {
	m := mcu.NewATmega32u4()

	// Cycles
	m.CPU.Cycles = 12345
	if m.Cycles() != 12345 {
		t.Errorf("Expected Cycles() 12345, got %d", m.Cycles())
	}

	// Global Interrupts
	m.SetGlobalInterrupts(true)
	if !m.GetGlobalInterrupts() {
		t.Error("Expected GlobalInterrupts to be true")
	}
	m.SetGlobalInterrupts(false)
	if m.GetGlobalInterrupts() {
		t.Error("Expected GlobalInterrupts to be false")
	}

	// Clear Interrupt
	m.TriggerInterrupt(5)
	if (m.PendingInterrupts & (1 << 5)) == 0 {
		t.Error("Expected bit 5 to be set in PendingInterrupts")
	}
	m.ClearInterrupt(5)
	if (m.PendingInterrupts & (1 << 5)) != 0 {
		t.Error("Expected bit 5 to be cleared in PendingInterrupts")
	}

	// PinCallback
	called := false
	m.PinCallbackFunc = func(port int8, mask uint8, value uint8) {
		called = true
		if port != 1 || mask != 0x01 || value != 0x01 {
			t.Errorf("PinCallback unexpected args: %d, %02X, %02X", port, mask, value)
		}
	}
	m.PinCallback(1, 0x01, 0x01)
	if !called {
		t.Error("Expected PinCallbackFunc to be called")
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
