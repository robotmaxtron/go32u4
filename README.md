# go32u4 - ATmega32u4 Simulator

`go32u4` is a high-performance, cycle-accurate (for core ISA) simulator for the ATmega32u4 microcontroller, 
specifically designed to simulate the **Adafruit ItsyBitsy 32u4 5V**. It serves as a testing target for validating and 
benchmarking `.hex` files.

## Status & Features

`go32u4` provides a comprehensive simulation of the ATmega32u4 microcontroller.

### CPU Core
- **Instruction Set**: Full AVR instruction set (130+ instructions) with cycle-accurate timing.
    - Includes advanced Load/Store modes (LDD/STD with displacement, post-increment/pre-decrement).
    - Full bit manipulation support (BST, BLD, SBR, CBR).
- **SREG**: Full flag support (`I`, `T`, `H`, `S`, `V`, `N`, `Z`, `C`).
- **Stack Pointer**: Fully operational for all stack-related instructions.
- **Interrupt Controller**: Full support for the interrupt vector table (43 vectors), prioritization according to ATmega32U4 datasheet, and nested interrupts.
- **Power Management**: Support for `SLEEP` modes and wake-up cycles.
- **SPM**: Store Program Memory (SPM) for flash page erase, page write, and buffer fill logic.

### Peripherals Implementation

| Component            | Status      | Implementation Details                                                                |
|:---------------------|:------------|:--------------------------------------------------------------------------------------|
| **Interrupts**       | Implemented | Vector table handling, hardware-accurate prioritization, and execution logic.         |
| **Timers (0, 1, 3)** | Implemented | Prescalers, overflow interrupts, and compare-match (OCR A/B/C) support.               |
| **Timer 4**          | Implemented | 10-bit high-speed timer with PLL-based 64MHz clocking support and OCR4C as TOP.       |
| **USB Controller**   | Implemented | Register-accurate USB 2.0 endpoint management, HID report capture, and state machine. |
| **GPIO**             | Implemented | Register-level simulation with `PinCallback` mechanism.                               |
| **EEPROM**           | Implemented | Fully functional with 3.4ms write timing simulation and disk-backed persistence.      |
| **USART/SPI/TWI**    | Implemented | USART1 (Serial), SPI transfer flags, and TWI (I2C) master state machine.              |
| **MCP23018**         | Implemented | Full emulation of I2C I/O expander with banked registers and pin logic.               |
| **ADC**              | Implemented | Basic conversion logic and interrupt triggering.                                      |
| **Watchdog Timer**   | Implemented | Watchdog state and system reset logic including full `WDTCSR` bit logic.              |
| **Macro Execution**  | Implemented | Sequence-based HID report generation for complex macros with delay support.           |
| **Sleep Modes**      | Implemented | `SLEEP` instruction and power reduction register (PRR) support.                       |

## Known Gaps & Missing Features

While the simulator is highly capable of running real-world firmware, some low-level hardware details are not yet fully 
simulated:

- **Fuse Bits**: Emulation of fuse bits for clock selection and hardware configuration.

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
| `pkg/bus`         | 0.0% (interface only)         |
| `pkg/cpu`         | 76.4% (core ISA coverage)     |
| `pkg/loader`      | 85.0% (comprehensive parsing) |
| `pkg/mcu`         | 91.7% (memory and interrupts) |
| `pkg/peripherals` | 82.7% (core I/O and timers)   |
| **Total**         | **81.7%**                     |

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

The simulator is optimized for speed, leveraging Go's efficient execution. Benchmarks show:

- **Instruction Execution**: ~22.9 ns/op (approx. 43.7 MHz simulated speed)
- **Overall Performance**: ~48.4 MHz (on Apple M4)

This performance exceeds the real ATmega32u4's 16 MHz clock, making it ideal for rapid firmware validation and CI pipelines.

To run benchmarks:
```bash
go test -bench=. ./pkg/...
```

## Compatible & Adaptable Boards

While `go32u4` is specifically tuned for the **ATmega32u4**, its modular architecture makes it highly compatible or easily adaptable to other AVR-based microcontrollers. The core AVR instruction set (AVR-6) is common across many Atmel boards.

### Highly Compatible Boards
These boards use MCUs with very similar core architectures and would require only minor configuration changes (e.g., memory mapping, interrupt vectors):

- **ATmega328P (Arduino Uno/Nano)**:
    - **Status**: Easily Adaptable.
    - **Differences**: Smaller SRAM (2KB), different I/O register offsets, lacks USB and Timer 4.
    - **Adaptation**: Update `pkg/mcu` to use 328P memory map and 26 interrupt vectors.
- **ATmega2560 (Arduino Mega)**:
    - **Status**: Adaptable.
    - **Differences**: Larger Flash (256KB) requiring 3-byte PC support (EIND), more I/O ports, and additional timers.
    - **Adaptation**: Extend `pkg/cpu` for 22-bit addressing and update `pkg/mcu` for the expanded I/O and SRAM.
- **ATmega16U4 / ATmega8U4**:
    - **Status**: Fully Compatible.
    - **Differences**: Smaller Flash/SRAM sizes.
    - **Adaptation**: Simply adjust `FlashSize` and `SRAMSize` constants in `pkg/mcu`.

### How to Adapt for a New Board
1. **Define Memory Map**: Create a new struct in `pkg/mcu` that implements the `bus.Bus` interface with the target MCU's SRAM and I/O layout.
2. **Configure Interrupts**: Update the interrupt vector table handling in `handleInterrupts` to match the target datasheet.
3. **Map Peripherals**: Link the target I/O addresses in `pkg/peripherals` to the existing `Manager` logic or add new peripheral simulations as needed.

## Architecture & Extensibility

`go32u4` is designed with modularity in mind, making it possible to extend support to other AVR-based MCUs or even 
different architectures:

- **CPU Core (`pkg/cpu`)**: The instruction execution logic is separated from the memory bus. To support a different 
AVR core (e.g., ATmega328P), you would update the flash size and register mappings.
- **Memory Bus (`pkg/bus`)**: Defines the interfaces for memory access and interrupts. Any new MCU must implement the 
`Bus` and `InterruptController` interfaces.
- **Peripherals (`pkg/peripherals`)**: The `Manager` handles I/O callbacks. New hardware features can be added by 
implementing new cases in the `IOCallback` and updating the `Tick` function.
- **MCU Orchestration (`pkg/mcu`)**: This package ties everything together. Adding a new chip involves creating a new 
struct that implements `bus.Bus` and configuring the CPU/Peripherals accordingly.

## License

This project is licensed under the MIT License – see the LICENSE file for details.
