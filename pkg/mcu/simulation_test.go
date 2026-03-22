package mcu_test

import (
	"go32u4/pkg/mcu"
	"go32u4/pkg/peripherals"
	"testing"
)

func TestComplexMacroSimulation(t *testing.T) {
	m := mcu.NewATmega32u4()
	mcp := peripherals.NewMCP23018(0x20, m.Periph)
	m.Periph.RegisterTWIClient(mcp)

	// 1. Setup the environment
	// We'll simulate an ErgoDox-like matrix where Row 0, Col 0 is a macro key.
	// In ErgoDox, columns (PORTA on MCP23018) are driven LOW one by one.
	// Rows (PORTB on MCP23018) are read with pull-ups.

	// Mocking the firmware's I2C initialization
	// SLA+W for MCP23018 (0x20)
	m.WriteIO(peripherals.TWDR, 0x20<<1)
	m.WriteIO(peripherals.TWCR, (1<<7)|(1<<2)|(1<<5)) // START
	m.Step()
	m.WriteIO(peripherals.TWCR, (1<<7)|(1<<2)) // SLA+W
	m.Step()

	// Set IODIRA = 0x00 (outputs), IODIRB = 0xFF (inputs)
	m.WriteIO(peripherals.TWDR, 0x00) // Reg addr IODIRA
	m.WriteIO(peripherals.TWCR, (1<<7)|(1<<2))
	m.Step()
	m.WriteIO(peripherals.TWDR, 0x00) // Value for IODIRA
	m.WriteIO(peripherals.TWCR, (1<<7)|(1<<2))
	m.Step()
	m.WriteIO(peripherals.TWDR, 0xFF) // Value for IODIRB
	m.WriteIO(peripherals.TWCR, (1<<7)|(1<<2))
	m.Step()
	m.WriteIO(peripherals.TWCR, (1<<7)|(1<<2)|(1<<4)) // STOP
	m.Step()

	// 2. Simulate the Macro Trigger (Key Press)
	// User presses Row 0, Col 0.
	// We'll assume the firmware scans by setting PA0 low and reading PB0.
	mcp.External = 0xFEFE // Both PA0 and PB0 go low when connected and PA0 is driven low.

	// 3. Mock the firmware loop (Simplified)
	// Instead of running a real hex (which we don't have a full QMK-like one for),
	// we'll simulate what the firmware *would* do when it detects this press.
	
	// Firmware detects the press and decides to send a macro: "Hi"
	// "Hi" is Shift+h, i.
	
	// Report 1: Shift down
	m.Periph.USBConfigured = true
	m.WriteIO(peripherals.UENUM, 1) // EP1 (HID)
	m.WriteIO(peripherals.UEDATX, 0x02) // Modifier: Left Shift
	m.WriteIO(peripherals.UEDATX, 0x00) // Reserved
	m.WriteIO(peripherals.UEDATX, 0x00) // Key 1
	m.WriteIO(peripherals.UEDATX, 0x00) // Key 2
	m.WriteIO(peripherals.UEDATX, 0x00) // Key 3
	m.WriteIO(peripherals.UEDATX, 0x00) // Key 4
	m.WriteIO(peripherals.UEDATX, 0x00) // Key 5
	m.WriteIO(peripherals.UEDATX, 0x00) // Key 6
	
	// Firmware clears TXINI to send the report
	m.Periph.USBEndpoints[1].Interrupt |= 0x01 // TXINI set by hardware when ready
	m.WriteIO(peripherals.UEINTX, ^uint8(0x01)) // Firmware clears TXINI
	
	// Report 2: Shift+h
	m.WriteIO(peripherals.UEDATX, 0x02) // Modifier: Left Shift
	m.WriteIO(peripherals.UEDATX, 0x00)
	m.WriteIO(peripherals.UEDATX, 0x0B) // 'h'
	for i := 0; i < 5; i++ { m.WriteIO(peripherals.UEDATX, 0) }
	m.Periph.USBEndpoints[1].Interrupt |= 0x01
	m.WriteIO(peripherals.UEINTX, ^uint8(0x01))

	// Report 3: Shift down (key up 'h')
	m.WriteIO(peripherals.UEDATX, 0x02)
	for i := 0; i < 7; i++ { m.WriteIO(peripherals.UEDATX, 0) }
	m.Periph.USBEndpoints[1].Interrupt |= 0x01
	m.WriteIO(peripherals.UEINTX, ^uint8(0x01))

	// Report 4: 'i' (no shift)
	m.WriteIO(peripherals.UEDATX, 0x00)
	m.WriteIO(peripherals.UEDATX, 0x00)
	m.WriteIO(peripherals.UEDATX, 0x0C) // 'i'
	for i := 0; i < 5; i++ { m.WriteIO(peripherals.UEDATX, 0) }
	m.Periph.USBEndpoints[1].Interrupt |= 0x01
	m.WriteIO(peripherals.UEINTX, ^uint8(0x01))

	// Report 5: All up
	for i := 0; i < 8; i++ { m.WriteIO(peripherals.UEDATX, 0) }
	m.Periph.USBEndpoints[1].Interrupt |= 0x01
	m.WriteIO(peripherals.UEINTX, ^uint8(0x01))

	// 4. Verification
	if len(m.Periph.CapturedHIDReports) != 5 {
		t.Fatalf("Expected 5 captured HID reports, got %d", len(m.Periph.CapturedHIDReports))
	}

	// Verify 'h' was sent with Shift
	report2 := m.Periph.CapturedHIDReports[1]
	if report2[0] != 0x02 {
		t.Errorf("Report 2: Expected modifier 0x02 (Shift), got %02X", report2[0])
	}
	if report2[2] != 0x0B {
		t.Errorf("Report 2: Expected key 0x0B ('h'), got %02X", report2[2])
	}

	// Verify 'i' was sent without Shift
	report4 := m.Periph.CapturedHIDReports[3]
	if report4[0] != 0x00 {
		t.Errorf("Report 4: Expected modifier 0x00, got %02X", report4[0])
	}
	if report4[2] != 0x0C {
		t.Errorf("Report 4: Expected key 0x0C ('i'), got %02X", report4[2])
	}
}
