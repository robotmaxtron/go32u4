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
	maxCycles := flag.Uint64("cycles", 20000000, "Maximum number of cycles to run")
	mainTicksAddr := flag.Uint("ticks-addr", 0x0, "SRAM address of main_ticks (0-indexed SRAM)")
	usbConfigAddr := flag.Uint("usb-config-addr", 0x0, "SRAM address of usb_configuration")
	flag.Parse()

	if *hexFile == "" {
		fmt.Println("Please provide a hex file using -hex")
		return
	}

	m := mcu.NewATmega32u4()
	err := loader.LoadHex(*hexFile, m)
	if err != nil {
		fmt.Printf("Error loading hex file: %v\n", err)
		return
	}

	// Initialize memory to bypass USB wait (same as avr8js benchmark)
	if *usbConfigAddr != 0 {
		m.WriteSRAM(uint16(*usbConfigAddr), 1)
	}
	// Bypass PLL wait in usb_init
	// In IOCallback: if (value & (1 << PLLE)) != 0 { ioRegs[PLLCSR] |= (1 << PLOCK) }
	// PLLCSR is 0x49.
	m.IORegData[0x49] |= 1 << 0 // PLOCK

	fmt.Println("--- go32u4 Performance Benchmarking ---")
	start := time.Now()
	for m.CPU.Cycles < *maxCycles && !m.CPU.Halted {
		err := m.Step()
		if err != nil {
			fmt.Printf("CPU Halted: %v\n", err)
			break
		}
		// In go32u4, we might need to bypass EEPROM hang too if implemented
		// For now let's see.
	}
	duration := time.Since(start)

	var finalTicks uint16
	if *mainTicksAddr != 0 {
		// SRAM in ATmega32u4 struct starts after Registers and IO
		// but ReadSRAM(addr) handles the mapping.
		// In AVR, SRAM (data space) is:
		// 0x00-0x1F: Registers
		// 0x20-0x5F: I/O Registers
		// 0x60-0xFF: Extended I/O Registers
		// 0x100+: Internal SRAM
		// symbols in ELF usually have 0x800000 offset.
		// For ATmega32u4, main_ticks is in Internal SRAM.

		low := m.ReadSRAM(uint16(*mainTicksAddr))
		high := m.ReadSRAM(uint16(*mainTicksAddr + 1))
		finalTicks = uint16(low) | (uint16(high) << 8)
	}

	fmt.Printf("Results after %d CPU cycles:\n", m.CPU.Cycles)
	fmt.Printf("Total Main Ticks (Completed Scans): %d\n", finalTicks)
	fmt.Printf("----------------------------------\n")
	if finalTicks > 0 {
		cyclesPerScan := float64(m.CPU.Cycles) / float64(finalTicks)
		simulatedTimeS := float64(m.CPU.Cycles) / 16000000.0
		scanRateHz := float64(finalTicks) / simulatedTimeS

		fmt.Printf("Average CPU Cycles per Scan: %.2f\n", cyclesPerScan)
		fmt.Printf("Estimated Scan Rate: %.2f Hz\n", scanRateHz)
		fmt.Printf("Simulated Time: %.4f s\n", simulatedTimeS)
	} else {
		fmt.Println("No scans completed.")
	}
	fmt.Printf("----------------------------------\n")
	fmt.Printf("Execution time: %v\n", duration)
	if duration > 0 {
		fmt.Printf("Performance: %.2f MHz\n", float64(m.CPU.Cycles)/duration.Seconds()/1e6)
	}
}
