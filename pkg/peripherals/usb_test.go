package peripherals_test

import (
	"go32u4/pkg/mcu"
	"testing"
)

func TestUSBGetDescriptor(t *testing.T) {
	m := mcu.NewATmega32u4()
	
	// Setup GET_DESCRIPTOR (Device) packet
	// bmRequestType: 0x80 (Device to Host, Standard, Device)
	// bRequest: 0x06 (GET_DESCRIPTOR)
	// wValue: 0x0100 (Type: 0x01 = Device, Index: 0)
	// wIndex: 0x0000
	// wLength: 0x0012 (18 bytes)
	setupPacket := []byte{0x80, 0x06, 0x00, 0x01, 0x00, 0x00, 0x12, 0x00}
	
	m.Periph.USBEndpoints[0].SetupFIFO = append(m.Periph.USBEndpoints[0].SetupFIFO, setupPacket...)
	
	// Trigger USB update
	m.Periph.Tick(1)
	
	// Check if TXINI is set on EP0
	if (m.Periph.USBEndpoints[0].Interrupt & 0x01) == 0 {
		t.Errorf("TXINI should be set on EP0 after GET_DESCRIPTOR")
	}
	
	// Check FIFO content
	fifo := m.Periph.USBEndpoints[0].FIFO
	if len(fifo) != 18 {
		t.Errorf("Expected 18 bytes in EP0 FIFO, got %d", len(fifo))
	} else if fifo[0] != 18 || fifo[1] != 0x01 {
		t.Errorf("FIFO doesn't look like a device descriptor: %v", fifo)
	}
}
