package main

import (
	"flag"
	"fmt"
	"go32u4/pkg/loader"
	"go32u4/pkg/mcu"
	"time"
)

func main() {
	hexFile := flag.String("hex", "", "Hex file to load")
	eepromFile := flag.String("eeprom", "", "EEPROM file for persistence")
	maxCycles := flag.Uint64("cycles", 1000000, "Maximum number of cycles to run")
	flag.Parse()

	if *hexFile == "" {
		fmt.Println("Please provide a hex file using -hex")
		return
	}

	m := mcu.NewATmega32u4()
	if *eepromFile != "" {
		err := m.LoadEEPROM(*eepromFile)
		if err != nil {
			fmt.Printf("Error loading EEPROM: %v\n", err)
			return
		}
	}
	err := loader.LoadHex(*hexFile, m)
	if err != nil {
		fmt.Printf("Error loading hex file: %v\n", err)
		return
	}

	start := time.Now()
	for m.CPU.Cycles < *maxCycles && !m.CPU.Halted {
		err := m.Step()
		if err != nil {
			fmt.Printf("CPU Halted: %v\n", err)
			break
		}
	}
	duration := time.Since(start)

	fmt.Printf("Executed %d cycles in %v\n", m.CPU.Cycles, duration)
	if duration > 0 {
		fmt.Printf("Performance: %.2f MHz\n", float64(m.CPU.Cycles)/duration.Seconds()/1e6)
	}
}
