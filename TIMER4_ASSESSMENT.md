# ATmega32u4 Timer 4 Implementation Assessment

Timer 4 on the ATmega32u4 is a high-speed, 10-bit Timer/Counter with unique features compared to the standard 8-bit and 16-bit timers (Timer 0, 1, 3) found in the AVR architecture. This document assesses the current implementation gaps and the technical challenges for full emulation.

## 1. Key Features of ATmega32u4 Timer 4
- **High-Speed Operation**: Can be clocked from a Phase-Locked Loop (PLL) at 64 MHz (or 96 MHz depending on configuration), allowing for very high PWM frequencies.
- **10-bit Resolution**: Unlike Timer 0 (8-bit), Timer 4 supports up to 10-bit resolution (0-1023).
- **Multiple Compare Units**: Three independent output compare units (A, B, D) and one unit (C) typically used to define the TOP value.
- **Enhanced PWM Modes**: Supports "Enhanced PWM" with up to 11-bit effective resolution and a "PWM6" mode for brushless DC motor control.
- **Dead-Time Generator**: Allows programmable delay between complementary PWM outputs.
- **Asynchronous Clocking**: Can run independently of the CPU clock using the PLL.

## 2. Current Implementation Gaps

The current implementation in `pkg/peripherals/peripherals.go` is a minimal 8-bit approximation of Timer 4.

### 2.1 Register Gaps
- **Missing Registers**: `TCCR4C`, `TCCR4D`, `TCCR4E`, `TC4H` (High Byte Temp), `OCR4B`, `OCR4C`, `OCR4D`, `DT4` (Dead Time) are defined as constants but not fully utilized in logic.
- **High Byte Handling**: The `TC4H` register is critical for accessing 10-bit values in `TCNT4`, `OCR4A/B/C/D`. The current code only tracks an 8-bit `Timer4Counter`.
- **Top Value**: The implementation assumes Timer 4 always wraps at 255 (8-bit). In reality, `OCR4C` defines the TOP value in most modes.

### 2.2 Functional Gaps
- **Clock Source / PLL**: There is no emulation of the `PLLCSR` (PLL Control and Status Register) or the asynchronous clocking logic. The timer currently ticks based on the CPU cycle count.
- **10-bit Logic**: The counter and compare logic are strictly 8-bit.
- **Waveform Generation**: Standard modes like "Phase and Frequency Correct PWM" and "PWM6" are not implemented.
- **Interrupts**: Only `TOIE4` (Overflow) and `OCIE4A` (Compare Match A) are partially implemented. `OCIE4B` and `OCIE4D` are missing.

## 3. Implementation Challenges

### 3.1 Asynchronous Clocking (PLL)
Timer 4 can run at 64MHz while the CPU runs at 16MHz. 
- **Challenge**: The `Tick(cycles uint64)` function currently assumes all peripherals advance based on CPU cycles. For Timer 4, the number of "ticks" per CPU cycle depends on the PLL state and the `TCCR4B` prescaler.
- **Solution**: The `Manager` needs to track the PLL state and calculate a fractional or multi-tick update for Timer 4 during each CPU cycle.

### 3.2 10-bit Register Access (TC4H)
Accessing 10-bit registers on an 8-bit bus requires a temporary high-byte register (`TC4H`). 
- **Challenge**: Writing to `TC4H` stores data that is used in the *next* write to an 8-bit register (like `OCR4C` or `TCNT4`). This state must be maintained correctly in the `Manager`.
- **Solution**: Implement a `TempHighByte` field in the `Manager` and update the 10-bit internal registers only when the low-byte part of the pair is written.

### 3.3 Complexity of PWM Modes
Timer 4 has specialized PWM modes (like PWM6) that interact with multiple pins and have complex timing requirements.
- **Challenge**: Implementing these modes requires a more sophisticated state machine than the simple increment-and-compare used for Timer 0.

## 4. Proposed Roadmap for Timer 4
1. **Infrastructure**: Add `PLLCSR` emulation and 10-bit register state to `Manager`.
2. **Register Logic**: Update `IOCallback` to handle `TC4H` and 10-bit writes/reads for Timer 4 registers.
3. **Core Logic**: Update `updateTimer4` to use `OCR4C` as TOP and support 10-bit counting.
4. **Clocking**: Implement PLL-based clocking multipliers.
5. **Advanced Features**: Add Compare B/D, Dead-Time, and PWM modes.

## 5. Summary
Implementation Difficulty: **High**
The high-speed and asynchronous nature of Timer 4 makes it significantly more complex to emulate accurately than the standard AVR timers. Full implementation will require careful cycle-accurate logic and state management for the PLL and 10-bit register access.
