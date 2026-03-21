### Test Data and Assets

This directory contains resources used for testing and benchmarking the `go32u4` simulator.

#### Contents

- **`benchmark.go`**: A performance measurement tool that simulates real-world firmware behavior. It includes logic to bypass long initialization sequences (like USB PLL lock) and can track firmware "scans" via specific SRAM addresses.
- **`hex/`**: Directory containing compiled AVR hex files for testing.
    - `test.hex`: A basic test firmware used for general functional and performance testing.

#### Usage: Benchmarking

To run the performance benchmark:

```bash
go run testdata/benchmark.go -hex testdata/hex/test.hex -cycles 20000000
```

Common flags for `benchmark.go`:
- `-hex`: Path to the hex file to load.
- `-cycles`: Maximum number of CPU cycles to simulate (default: 20,000,000).
- `-ticks-addr`: SRAM address of a 16-bit "ticks" counter in the guest firmware to calculate average cycles per scan.
- `-usb-config-addr`: SRAM address of the USB configuration byte to bypass wait-for-enumeration loops.
