package peripherals_test

import (
	"go32u4/pkg/mcu"
	"go32u4/pkg/peripherals"
	"testing"
)

func TestEEPROMTiming(t *testing.T) {
	m := mcu.NewATmega32u4()
	
	// 1. Set Address and Data
	m.WriteIO(peripherals.EEARL, 0x10)
	m.WriteIO(peripherals.EEARH, 0x00)
	m.WriteIO(peripherals.EEDR, 0x42)
	
	// 2. Set EEMPE (Master Write Enable)
	m.WriteIO(peripherals.EECR, 0x04)
	
	// 3. Set EEPE (Write Enable)
	m.WriteIO(peripherals.EECR, 0x04 | 0x02)
	
	// Verify EEPE is set
	if (m.ReadIO(peripherals.EECR) & 0x02) == 0 {
		t.Fatal("EEPE bit should be set")
	}
	
	// 4. Tick some cycles, EEPE should still be set
	for i := 0; i < 1000; i++ {
		m.Periph.Tick(1)
	}
	
	if (m.ReadIO(peripherals.EECR) & 0x02) == 0 {
		t.Fatal("EEPE bit should still be set after 1000 cycles")
	}
	
	// 5. Tick enough cycles to finish write
	m.Periph.Tick(54400) // The timer we set
	
	if (m.ReadIO(peripherals.EECR) & 0x02) != 0 {
		t.Fatal("EEPE bit should be cleared after write time")
	}
	
	// 6. Verify data was written
	m.WriteIO(peripherals.EECR, 0x01) // EERE (Read Enable)
	if m.ReadIO(peripherals.EEDR) != 0x42 {
		t.Errorf("Expected EEDR 0x42, got %02X", m.ReadIO(peripherals.EEDR))
	}
}
