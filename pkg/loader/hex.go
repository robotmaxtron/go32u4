package loader

import (
	"fmt"
	"io"
	"os"
)

// Record types
const (
	RecordData                   = 0x00
	RecordEndOfFile              = 0x01
	RecordExtendedSegmentAddress = 0x02
	RecordStartSegmentAddress    = 0x03
	RecordExtendedLinearAddress  = 0x04
	RecordStartLinearAddress     = 0x05
)

// FlashMemory defines the interface for a memory that can be populated with flash data.
type FlashMemory interface {
	Flash() []uint16
}

// LoadHex parses an Intel Hex file and populates the Flash memory.
func LoadHex(filename string, target FlashMemory) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	flash := target.Flash()

	for {
		var colon byte
		_, err := fmt.Fscanf(file, "%c", &colon)
		if err == io.EOF {
			break
		}
		if colon == '\r' || colon == '\n' {
			continue
		}
		if colon != ':' {
			return fmt.Errorf("invalid hex file format (found %c)", colon)
		}

		var count, addressHigh, addressLow, recordType uint8
		n, err := fmt.Fscanf(file, "%02x%02x%02x%02x", &count, &addressHigh, &addressLow, &recordType)
		if err != nil {
			return fmt.Errorf("failed to read hex header (n=%d): %v", n, err)
		}

		address := uint16(addressHigh)<<8 | uint16(addressLow)

		data := make([]uint8, count)
		for i := 0; i < int(count); i++ {
			_, err = fmt.Fscanf(file, "%02x", &data[i])
			if err != nil {
				return err
			}
		}

		var fileChecksum uint8
		_, err = fmt.Fscanf(file, "%02x", &fileChecksum)
		if err != nil {
			return err
		}

		switch recordType {
		case RecordData:
			for i := 0; i < int(count); i += 2 {
				flashAddr := (address + uint16(i)) / 2
				if int(flashAddr) < len(flash) {
					var word uint16
					if i+1 < int(count) {
						word = uint16(data[i]) | uint16(data[i+1])<<8
					} else {
						word = uint16(data[i])
					}
					flash[flashAddr] = word
				}
			}
		case RecordEndOfFile:
			return nil
		case RecordExtendedLinearAddress:
		}
	}

	return nil
}
