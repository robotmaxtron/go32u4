package loader_test

import (
	"go32u4/pkg/loader"
	"os"
	"testing"
)

type MockFlash struct {
	data [32768]uint16
}

func (m *MockFlash) Flash() []uint16 {
	return m.data[:]
}

func TestLoadHex(t *testing.T) {
	// Create a simple hex file
	// :020000000F00EF
	// :020002001F00DD
	// :00000001FF
	hexContent := ":020000000F00EF\n:020002001F00DD\n:00000001FF\n"
	tmpFile := "test_load.hex"
	err := os.WriteFile(tmpFile, []byte(hexContent), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Remove(tmpFile)
	}()

	mock := &MockFlash{}
	err = loader.LoadHex(tmpFile, mock)
	if err != nil {
		t.Fatalf("LoadHex failed: %v", err)
	}

	if mock.data[0] != 0x000F {
		t.Errorf("Expected data[0] 0x000F, got %04X", mock.data[0])
	}
	if mock.data[1] != 0x001F {
		t.Errorf("Expected data[1] 0x001F, got %04X", mock.data[1])
	}
}

func TestLoadHexInvalid(t *testing.T) {
	tmpFile := "invalid.hex"
	err := os.WriteFile(tmpFile, []byte("invalid content"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Remove(tmpFile)
	}()

	mock := &MockFlash{}
	err = loader.LoadHex(tmpFile, mock)
	if err == nil {
		t.Error("Expected error for invalid hex file, got nil")
	}
}

func TestLoadHexFileNotFound(t *testing.T) {
	mock := &MockFlash{}
	err := loader.LoadHex("non_existent.hex", mock)
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}
