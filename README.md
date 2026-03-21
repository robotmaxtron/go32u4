# go32u4 - ATmega32u4 Simulator

`go32u4` is a high-performance, cycle-accurate (for core ISA) simulator for the ATmega32u4 microcontroller, 
specifically designed to simulate the **Adafruit ItsyBitsy 32u4 5V**. It serves as a testing target for validating and 
benchmarking `.hex` files.

## Status & Features

`go32u4` provides a comprehensive simulation of the ATmega32u4 microcontroller.

### CPU Core
- **Instruction Set**: Full AVR instruction set (130+ instructions) with cycle-accurate timing.
- **SREG**: Full flag support (`I`, `T`, `H`, `S`, `V`, `N`, `Z`, `C`).
- **Stack Pointer**: Fully operational for all stack-related instructions.
- **Interrupt Controller**: Full support for the interrupt vector table (43 vectors), prioritization, and nested 
interrupts.
- **Power Management**: Support for `SLEEP` modes and wake-up cycles.

### Peripherals Implementation

| Component            | Status      | Implementation Details                                                                             |
|:---------------------|:------------|:---------------------------------------------------------------------------------------------------|
| **Interrupts**       | Implemented | Vector table handling, prioritization, and execution logic.                                        |
| **Timers (0, 1, 3)** | Implemented | Prescalers, overflow interrupts, and compare-match (OCR).                                          |
| **Timer 4**          | Partial     | High-speed timer basics. See [TIMER4_ASSESSMENT.md](TIMER4_ASSESSMENT.md) for implementation gaps. |
| **USB Controller**   | Emulated    | High-level USB CDC (Serial) emulation via virtual buffers.                                         |
| **GPIO**             | Implemented | Register-level simulation with `PinCallback` mechanism.                                            |
| **EEPROM**           | Implemented | Fully functional with optional disk-backed persistence.                                            |
| **USART/SPI/TWI**    | Implemented | USART1 (Serial), SPI transfer flags, and TWI (I2C) master state machine.                           |
| **ADC**              | Implemented | Basic conversion logic and interrupt triggering.                                                   |
| **Watchdog Timer**   | Implemented | Watchdog state and system reset logic including full `WDTCSR` bit logic.                           |
| **Sleep Modes**      | Implemented | `SLEEP` instruction and power reduction register support.                                          |
| **Documentation**    | Updated     | Detailed assessment for [Timer 4 implementation](TIMER4_ASSESSMENT.md).                            |

## Known Gaps & Missing Features

While the simulator is highly capable of running real-world firmware, some low-level hardware details are not yet fully simulated:

- **SPM (Store Program Memory)**: The `SPM` instruction for self-programming/bootloader simulation is currently missing.
- **Advanced Timer 4**: Lacks PLL interaction and complex 10-bit PWM logic. See [TIMER4_ASSESSMENT.md](TIMER4_ASSESSMENT.md).
- **USB Hardware**: The simulator uses high-level CDC emulation instead of a full register-level USB 2.0 state machine.

## Project Structure

The project follows a modular architecture to ensure maintainability and scalability:

- `pkg/cpu`: Core instruction set execution logic.
- `pkg/bus`: Communication interfaces between CPU and peripherals/memory.
- `pkg/mcu`: System orchestration (e.g., `ATmega32u4` struct) and memory mapping.
- `pkg/peripherals`: Hardware emulation for timers, USB, serial, etc.
- `pkg/loader`: Intel Hex file loader.

## Installation

Ensure you have [Go](https://golang.org/doc/install) installed (version 1.16 or later recommended).

```bash
git clone https://github.com/robotmaxtron/go32u4.git
cd go32u4
go build -o go32u4 main.go
```

## Usage

Run a `.hex` file through the simulator using the CLI:

```bash
./go32u4 -hex path/to/your_firmware.hex -cycles 1000000
```

### Options

| Flag      | Description                                   | Default   |
|:----------|:----------------------------------------------|:----------|
| `-hex`    | Path to the Intel Hex file to load (Required) | `""`      |
| `-eeprom` | Path to a file for EEPROM data persistence    | `""`      |
| `-cycles` | Maximum number of cycles to execute           | `1000000` |

## Testing & Coverage

`go32u4` includes a comprehensive suite of unit tests for all critical components:

- **CPU Core**: Instruction decoding, flag updates, and stack operations.
- **MCU**: Memory mapping (Registers/IO/SRAM) and interrupt prioritization.
- **Peripherals**: Timer overflows, EEPROM read/write, and TWI state machine.
- **Loader**: Intel Hex record parsing and validation.

### Code Coverage

Current code coverage results (as of March 2026):

| Package           | Statement Coverage            |
|:------------------|:------------------------------|
| `pkg/cpu`         | 20.2% (focus on core timing)  |
| `pkg/loader`      | 84.6% (comprehensive parsing) |
| `pkg/mcu`         | 41.6% (memory and interrupts) |
| `pkg/peripherals` | 33.5% (core I/O and timers)   |
| **Total**         | **28.6%**                     |

To run tests and generate a coverage report:
```bash
go test -cover ./pkg/...
```

### Linting

The project uses `golangci-lint` to ensure code quality. To run the linter:

```bash
golangci-lint run ./...
```

## Performance

The simulator is optimized for speed, leveraging Go's efficient execution. Benchmarks on an Apple M4 silicon show:

- **Instruction Execution**: ~13.0 ns/op (approx. 76.9 MHz simulated speed)
- **Peripheral Ticking**: ~0.52 ns/op

This performance exceeds the real ATmega32u4's 16 MHz clock, making it ideal for rapid firmware validation and CI pipelines.

To run benchmarks:
```bash
go test -bench=. ./pkg/...
```

## Architecture & Extensibility

`go32u4` is designed with modularity in mind, making it possible to extend support to other AVR-based MCUs or even different architectures:

- **CPU Core (`pkg/cpu`)**: The instruction execution logic is separated from the memory bus. To support a different AVR core (e.g., ATmega328P), you would update the flash size and register mappings.
- **Memory Bus (`pkg/bus`)**: Defines the interfaces for memory access and interrupts. Any new MCU must implement the `Bus` and `InterruptController` interfaces.
- **Peripherals (`pkg/peripherals`)**: The `Manager` handles I/O callbacks. New hardware features can be added by implementing new cases in the `IOCallback` and updating the `Tick` function.
- **MCU Orchestration (`pkg/mcu`)**: This package ties everything together. Adding a new chip involves creating a new struct that implements `bus.Bus` and configuring the CPU/Peripherals accordingly.

## License

This project is licensed under the MIT License – see the LICENSE file for details.
