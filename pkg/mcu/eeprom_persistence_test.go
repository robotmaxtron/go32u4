package mcu_test

import (
	"go32u4/pkg/mcu"
	"go32u4/pkg/peripherals"
	"os"
	"path/filepath"
	"testing"
)

func TestEEPROMPersistence(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "eeprom_test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	eepromFile := filepath.Join(tmpDir, "eeprom.bin")
	m := mcu.NewATmega32u4()

	// 1. Test LoadEEPROM with non-existent file
	err = m.LoadEEPROM(eepromFile)
	if err != nil {
		t.Errorf("LoadEEPROM failed for non-existent file: %v", err)
	}

	// 2. Write data to EEPROM via registers
	m.WriteIO(peripherals.EEARL, 0x05)
	m.WriteIO(peripherals.EEARH, 0x00)
	m.WriteIO(peripherals.EEDR, 0xDE)
	m.WriteIO(peripherals.EECR, 0x04) // EEMPE
	m.WriteIO(peripherals.EECR, 0x06) // EEPE | EEMPE

	// 3. Test SaveEEPROM (without setting filename explicitly in mcu, 
	// but LoadEEPROM already set it)
	err = m.SaveEEPROM()
	if err != nil {
		t.Errorf("SaveEEPROM failed: %v", err)
	}

	// Verify file exists and has correct content
	data, err := os.ReadFile(eepromFile)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 1024 {
		t.Errorf("Expected EEPROM file size 1024, got %d", len(data))
	}
	if data[0x05] != 0xDE {
		t.Errorf("Expected data at 0x05 to be 0xDE, got %02X", data[0x05])
	}

	// 4. Test LoadEEPROM with existing file
	m2 := mcu.NewATmega32u4()
	err = m2.LoadEEPROM(eepromFile)
	if err != nil {
		t.Errorf("LoadEEPROM failed: %v", err)
	}
	
	// Read back via registers
	m2.WriteIO(peripherals.EEARL, 0x05)
	m2.WriteIO(peripherals.EECR, 0x01) // EERE
	if val := m2.ReadIO(peripherals.EEDR); val != 0xDE {
		t.Errorf("Expected EEDR 0xDE after loading, got %02X", val)
	}
}

func TestSaveEEPROMNoFile(t *testing.T) {
	m := mcu.NewATmega32u4()
	// Should return nil if no file is set
	err := m.SaveEEPROM()
	if err != nil {
		t.Errorf("SaveEEPROM should return nil if no file set, got %v", err)
	}
}

func TestLoadEEPROMError(t *testing.T) {
	m := mcu.NewATmega32u4()
	// Try to load a directory as a file to trigger an error
	tmpDir, _ := os.MkdirTemp("", "eeprom_err")
	defer func() { _ = os.RemoveAll(tmpDir) }()
	
	err := m.LoadEEPROM(tmpDir)
	if err == nil {
		t.Error("Expected error when loading a directory as EEPROM file")
	}
}
