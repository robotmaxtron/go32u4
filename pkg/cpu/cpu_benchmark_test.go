package cpu_test

import (
	"go32u4/pkg/mcu"
	"testing"
)

func BenchmarkInstructionExecution(b *testing.B) {
	m := mcu.NewATmega32u4()
	// Fill flash with NOPs (0x0000)
	for i := range m.FlashData {
		m.FlashData[i] = 0x0000
	}
	// Add some instructions
	m.FlashData[0] = 0xE00F // LDI R16, 0x0F
	m.FlashData[1] = 0xE01F // LDI R17, 0xFF
	m.FlashData[2] = 0x0F01 // ADD R16, R17
	m.FlashData[3] = 0xCFFC // RJMP -4 (jump back to 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := m.Step()
		if err != nil {
			b.Fatal(err)
		}
		if m.CPU.PC >= 4 {
			m.CPU.PC = 0
		}
	}
}

func BenchmarkTimerTick(b *testing.B) {
	m := mcu.NewATmega32u4()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Periph.Tick(1)
	}
}
