package peripherals_test

import (
	"go32u4/pkg/peripherals"
	"testing"
)

func TestW1CFlags(t *testing.T) {
	sys := &MockSystem{ioRegs: make([]uint8, 256)}
	p := peripherals.NewManager(sys)

	// Test TIFR0 (W1C)
	sys.ioRegs[peripherals.TIFR0] = 0x07 // All flags set
	p.IOCallback(peripherals.TIFR0, 0x01, true) // Write 1 to clear bit 0
	if sys.ioRegs[peripherals.TIFR0] != 0x06 {
		t.Errorf("Expected TIFR0 0x06, got %02X", sys.ioRegs[peripherals.TIFR0])
	}

	p.IOCallback(peripherals.TIFR0, 0x06, true) // Write 1 to clear bits 1 and 2
	if sys.ioRegs[peripherals.TIFR0] != 0x00 {
		t.Errorf("Expected TIFR0 0x00, got %02X", sys.ioRegs[peripherals.TIFR0])
	}

	// Test ADCSRA (ADIF is W1C)
	// ADCSRA bit 4 is ADIF, bit 3 is ADIE.
	sys.ioRegs[peripherals.ADCSRA] = (1 << 4) | (1 << 3) // ADIF and ADIE set
	// Write with ADIF set to 1 and ADIE set to 1 (W1C for ADIF)
	p.IOCallback(peripherals.ADCSRA, (1<<4)|(1<<3), true)
	if (sys.ioRegs[peripherals.ADCSRA] & (1 << 4)) != 0 {
		t.Error("Expected ADIF bit to be cleared in ADCSRA")
	}
	if (sys.ioRegs[peripherals.ADCSRA] & (1 << 3)) == 0 {
		t.Error("Expected ADIE bit to remain set in ADCSRA")
	}
}

func TestReadOnlyMasks(t *testing.T) {
	sys := &MockSystem{ioRegs: make([]uint8, 256)}
	p := peripherals.NewManager(sys)

	// Test TWSR (Bits 7:3 are read-only)
	sys.ioRegs[peripherals.TWSR] = 0xF8 // Status = 0x1F, Prescaler = 0
	p.IOCallback(peripherals.TWSR, 0x07, true) // Try to write to status bits and prescaler
	// Bits 7:3 should remain 0xF8, bits 1:0 should become 0x03 (0x07 & 0x03)
	// Bit 2 is reserved and should be 0.
	if sys.ioRegs[peripherals.TWSR] != 0xFB { // 0xF8 | 0x03
		t.Errorf("Expected TWSR 0xFB, got %02X", sys.ioRegs[peripherals.TWSR])
	}

	// Test USBSTA (Bit 0 VBUS is read-only)
	sys.ioRegs[peripherals.USBSTA] = 0x01 // VBUS = 1
	p.IOCallback(peripherals.USBSTA, 0x00, true) // Try to clear VBUS
	if (sys.ioRegs[peripherals.USBSTA] & 0x01) == 0 {
		t.Error("Expected VBUS bit to remain set in USBSTA")
	}

	// Test SPMCSR (Bit 6 RWWSB is read-only)
	sys.ioRegs[peripherals.SPMCSR] = 1 << 6
	p.IOCallback(peripherals.SPMCSR, 0, true)
	if (sys.ioRegs[peripherals.SPMCSR] & (1 << 6)) == 0 {
		t.Error("Expected RWWSB bit to remain set in SPMCSR")
	}
}

func TestSPMCSRAutoClear(t *testing.T) {
	sys := &MockSystem{ioRegs: make([]uint8, 256)}
	p := peripherals.NewManager(sys)

	p.IOCallback(peripherals.SPMCSR, (1 << 0) | (1 << 1), true) // SPMEN and PGERS
	if (sys.ioRegs[peripherals.SPMCSR] & (1 << 0)) == 0 {
		t.Error("SPMEN should be set")
	}

	p.Tick(3)
	if (sys.ioRegs[peripherals.SPMCSR] & (1 << 0)) == 0 {
		t.Error("SPMEN cleared too early")
	}
	
	p.Tick(1)
	if (sys.ioRegs[peripherals.SPMCSR] & (1 << 0)) != 0 {
		t.Error("SPMEN should be cleared after 4 cycles")
	}
	if (sys.ioRegs[peripherals.SPMCSR] & (1 << 1)) != 0 {
		t.Error("PGERS should be cleared after 4 cycles")
	}
}
