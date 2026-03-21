package mcu_test

import (
	"bytes"
	"go32u4/pkg/mcu"
	"go32u4/pkg/peripherals"
	"testing"
)

// TestUSBHIDReports verifies that different types of HID reports (Keyboard, Mouse, Consumer)
// are correctly captured by the simulation when the firmware "sends" them.
func TestUSBHIDReports(t *testing.T) {
	m := mcu.NewATmega32u4()
	m.Periph.USBConfigured = true

	// Helper to send a report on a given endpoint
	sendReport := func(epNum uint8, data []byte) {
		m.WriteIO(peripherals.UENUM, epNum)
		for _, b := range data {
			m.WriteIO(peripherals.UEDATX, b)
		}
		// Set TXINI to simulate hardware ready
		m.Periph.USBEndpoints[epNum].Interrupt |= 0x01
		// Firmware clears TXINI to "send"
		m.WriteIO(peripherals.UEINTX, 0xFE) // Clear bit 0
	}

	tests := []struct {
		name     string
		endpoint uint8
		report   []byte
	}{
		{
			name:     "Keyboard Report",
			endpoint: 1,
			report:   []byte{0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00}, // 'a'
		},
		{
			name:     "Mouse Report",
			endpoint: 2,
			report:   []byte{0x01, 0x10, 0xEF, 0x00}, // Left click, X=16, Y=-17
		},
		{
			name:     "Consumer Control Report",
			endpoint: 3,
			report:   []byte{0xE9, 0x00}, // Volume Up (0x00E9)
		},
		{
			name:     "System Control Report",
			endpoint: 3,
			report:   []byte{0x01}, // Power Down
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sendReport(tt.endpoint, tt.report)

			if len(m.Periph.CapturedHIDReports) <= i {
				t.Fatalf("Expected report %d to be captured, but only %d reports found", i+1, len(m.Periph.CapturedHIDReports))
			}

			captured := m.Periph.CapturedHIDReports[i]
			if !bytes.Equal(captured, tt.report) {
				t.Errorf("Captured report mismatch.\nExpected: %02X\nGot:      %02X", tt.report, captured)
			}
		})
	}
}

// TestUSBHIDMultipleReports verifies that multiple reports on the same endpoint are captured correctly.
func TestUSBHIDMultipleReports(t *testing.T) {
	m := mcu.NewATmega32u4()
	m.Periph.USBConfigured = true

	m.WriteIO(peripherals.UENUM, 1)

	// Send 3 reports
	reports := [][]byte{
		{0x02, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00}, // Shift + 'a'
		{0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // Shift only
		{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // All up
	}

	for _, r := range reports {
		for _, b := range r {
			m.WriteIO(peripherals.UEDATX, b)
		}
		m.Periph.USBEndpoints[1].Interrupt |= 0x01
		m.WriteIO(peripherals.UEINTX, 0xFE)
	}

	if len(m.Periph.CapturedHIDReports) != 3 {
		t.Fatalf("Expected 3 reports, got %d", len(m.Periph.CapturedHIDReports))
	}

	for i, r := range reports {
		if !bytes.Equal(m.Periph.CapturedHIDReports[i], r) {
			t.Errorf("Report %d mismatch", i)
		}
	}
}

// TestUSBEndpointSelection verifies that selecting different endpoints works correctly for data writing.
func TestUSBEndpointSelection(t *testing.T) {
	m := mcu.NewATmega32u4()
	
	// Select EP1
	m.WriteIO(peripherals.UENUM, 1)
	m.WriteIO(peripherals.UEDATX, 0xAA)
	
	// Select EP2
	m.WriteIO(peripherals.UENUM, 2)
	m.WriteIO(peripherals.UEDATX, 0xBB)
	
	// Check EP1 FIFO
	if len(m.Periph.USBEndpoints[1].FIFO) != 1 || m.Periph.USBEndpoints[1].FIFO[0] != 0xAA {
		t.Errorf("EP1 FIFO mismatch")
	}
	
	// Check EP2 FIFO
	if len(m.Periph.USBEndpoints[2].FIFO) != 1 || m.Periph.USBEndpoints[2].FIFO[0] != 0xBB {
		t.Errorf("EP2 FIFO mismatch")
	}
	
	// Check UEBCLX (Byte Count Register)
	m.WriteIO(peripherals.UENUM, 1)
	if m.ReadIO(peripherals.UEBCLX) != 1 {
		t.Errorf("UEBCLX for EP1 should be 1")
	}
	
	m.WriteIO(peripherals.UENUM, 2)
	if m.ReadIO(peripherals.UEBCLX) != 1 {
		t.Errorf("UEBCLX for EP2 should be 1")
	}
}
