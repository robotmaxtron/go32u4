package mcu_test

import (
	"go32u4/pkg/cpu"
	"go32u4/pkg/mcu"
	"go32u4/pkg/peripherals"
	"testing"
)

func TestADIW(t *testing.T) {
	m := mcu.NewATmega32u4()
	m.CPU.Reg[24] = 0xFE
	m.CPU.Reg[25] = 0xFF
	// ADIW R24, 1 (1001 0110 0000 0001) -> 0x9601
	m.FlashData[0] = 0x9601
	_ = m.Step()
	if m.CPU.Reg[24] != 0xFF || m.CPU.Reg[25] != 0xFF {
		t.Errorf("ADIW failed: expected 0xFFFF, got %02X%02X", m.CPU.Reg[25], m.CPU.Reg[24])
	}
	if m.CPU.GetFlag(cpu.SREG_Z) { t.Error("Z flag should be clear") }

	// ADIW R24, 1 again -> 0x0000
	m.FlashData[1] = 0x9601
	_ = m.Step()
	if m.CPU.Reg[24] != 0x00 || m.CPU.Reg[25] != 0x00 {
		t.Errorf("ADIW failed: expected 0x0000, got %02X%02X", m.CPU.Reg[25], m.CPU.Reg[24])
	}
	if !m.CPU.GetFlag(cpu.SREG_Z) { t.Error("Z flag should be set") }
	if !m.CPU.GetFlag(cpu.SREG_C) { t.Error("C flag should be set") }
}

func TestSBIW(t *testing.T) {
	m := mcu.NewATmega32u4()
	m.CPU.Reg[24] = 0x01
	m.CPU.Reg[25] = 0x00
	// SBIW R24, 1 (1001 0111 0000 0001) -> 0x9701
	m.FlashData[0] = 0x9701
	_ = m.Step()
	if m.CPU.Reg[24] != 0x00 || m.CPU.Reg[25] != 0x00 {
		t.Errorf("SBIW failed: expected 0x0000, got %02X%02X", m.CPU.Reg[25], m.CPU.Reg[24])
	}
	if !m.CPU.GetFlag(cpu.SREG_Z) { t.Error("Z flag should be set") }

	// SBIW R24, 1 again -> 0xFFFF
	m.FlashData[1] = 0x9701
	_ = m.Step()
	if m.CPU.Reg[24] != 0xFF || m.CPU.Reg[25] != 0xFF {
		t.Errorf("SBIW failed: expected 0xFFFF, got %02X%02X", m.CPU.Reg[25], m.CPU.Reg[24])
	}
	if !m.CPU.GetFlag(cpu.SREG_C) { t.Error("C flag should be set") }
}

func TestRCALL_ICALL_IJMP(t *testing.T) {
	m := mcu.NewATmega32u4()
	// RCALL +2 (PC=1 -> 3)
	m.FlashData[0] = 0xD002
	_ = m.Step()
	if m.CPU.PC != 3 { t.Errorf("RCALL failed: expected PC 3, got %d", m.CPU.PC) }
	
	// Check stack: return address should be 1
	pch := m.CPU.Pop()
	pcl := m.CPU.Pop()
	if pcl != 1 || pch != 0 { t.Errorf("RCALL stack failed: expected 1, got %d", (uint16(pch)<<8)|uint16(pcl)) }

	// ICALL (Z=0x0100)
	m.CPU.Reg[30] = 0x00
	m.CPU.Reg[31] = 0x01
	m.CPU.PC = 10
	m.FlashData[10] = 0x9509
	_ = m.Step()
	if m.CPU.PC != 0x0100 { t.Errorf("ICALL failed: expected PC 0x0100, got %04X", m.CPU.PC) }

	// IJMP (Z=0x0200)
	m.CPU.Reg[30] = 0x00
	m.CPU.Reg[31] = 0x02
	m.CPU.PC = 20
	m.FlashData[20] = 0x9409
	_ = m.Step()
	if m.CPU.PC != 0x0200 { t.Errorf("IJMP failed: expected PC 0x0200, got %04X", m.CPU.PC) }
}

func TestLPM(t *testing.T) {
	m := mcu.NewATmega32u4()
	m.FlashData[0x100] = 0xABCD
	m.CPU.Reg[30] = 0x00
	m.CPU.Reg[31] = 0x02 // Address 0x200 (byte address for 0x100 word)
	m.FlashData[0] = 0x95C8
	_ = m.Step()
	if m.CPU.Reg[0] != 0xCD { t.Errorf("LPM failed (low): expected 0xCD, got %02X", m.CPU.Reg[0]) }

	m.CPU.Reg[30] = 0x01
	m.CPU.Reg[31] = 0x02 // Address 0x201
	m.FlashData[1] = 0x95C8
	_ = m.Step()
	if m.CPU.Reg[0] != 0xAB { t.Errorf("LPM failed (high): expected 0xAB, got %02X", m.CPU.Reg[0]) }
}

func TestUSBRegisterAddresses(t *testing.T) {
	m := mcu.NewATmega32u4()
	// UDADDR = 0xE8
	m.WriteIO(peripherals.UDADDR, 0x42)
	if m.ReadIO(peripherals.UDADDR) != 0x42 {
		t.Errorf("UDADDR write failed: expected 0x42, got %02X", m.ReadIO(peripherals.UDADDR))
	}
	
	// UEINTX (0xEB) Write-0-to-Clear according to datasheet
	// Need to select an endpoint first
	m.WriteIO(peripherals.UENUM, 0)
	m.Periph.USBEndpoints[0].Interrupt = 0xFF
	m.WriteIO(peripherals.UEINTX, 0xFE) // Write 0 to bit 0 to clear it
	if m.ReadIO(peripherals.UEINTX) != 0xFE {
		t.Errorf("UEINTX Write-0-to-Clear failed: expected 0xFE, got %02X", m.ReadIO(peripherals.UEINTX))
	}
}
