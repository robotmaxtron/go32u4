package cpu

import (
	"fmt"
	"go32u4/pkg/bus"
)

// ATmega32u4 Constants
const (
	FlashSize    = 32768
	RegisterSize = 32
	IORegSize    = 64
	ExtIORegSize = 160

	SREG_C = 0
	SREG_Z = 1
	SREG_N = 2
	SREG_V = 3
	SREG_S = 4
	SREG_H = 5
	SREG_T = 6
	SREG_I = 7
)

// CPU represents the ATmega32u4 CPU state.
type CPU struct {
	Reg    [RegisterSize]uint8
	PC     uint16
	SREG   uint8
	SP     uint16
	Cycles uint64
	Halted bool

	// Bus is the connection to memory and peripherals.
	Bus bus.Bus

	// InterruptController handles interrupt dispatching.
	InterruptController bus.InterruptController

	// TickPeripherals is called to update peripheral state.
	TickPeripherals func(cycles uint64)
}

func NewCPU(b bus.Bus, ic bus.InterruptController) *CPU {
	return &CPU{
		PC:                  0,
		SP:                  uint16(IORegSize + ExtIORegSize + 2560 - 1), // Default top of SRAM for 32u4
		Bus:                 b,
		InterruptController: ic,
	}
}

func (c *CPU) Step() error {
	if c.Halted {
		return fmt.Errorf("CPU halted")
	}

	flash := c.Bus.Flash()
	if int(c.PC) >= len(flash) {
		c.Halted = true
		return fmt.Errorf("PC out of bounds: %04X", c.PC)
	}
	opcode := flash[c.PC]
	c.PC++
	c.Execute(opcode)
	if c.Halted {
		// If instruction caused halt (e.g. SLEEP), don't increment cycles here
		// as it might be handled by the caller (MCU)
		return nil
	}
	c.Cycles++
	if c.TickPeripherals != nil {
		c.TickPeripherals(1)
	}
	return nil
}

func (c *CPU) SetFlag(bit uint8, value bool) {
	if value {
		c.SREG |= 1 << bit
	} else {
		c.SREG &= ^(1 << bit)
	}
	if bit == SREG_I && c.InterruptController != nil {
		c.InterruptController.SetGlobalInterrupts(value)
	}
}

func (c *CPU) GetFlag(bit uint8) bool {
	return (c.SREG & (1 << bit)) != 0
}

func (c *CPU) WriteSRAM(address uint16, value uint8) {
	c.Bus.WriteSRAM(address, value)
}

func (c *CPU) ReadSRAM(address uint16) uint8 {
	return c.Bus.ReadSRAM(address)
}

func (c *CPU) ReadIO(address uint16) uint8 {
	return c.Bus.ReadIO(address)
}

func (c *CPU) WriteIO(address uint16, value uint8) {
	c.Bus.WriteIO(address, value)
}

func (c *CPU) Push(value uint8) {
	c.WriteSRAM(c.SP, value)
	c.SP--
}

func (c *CPU) Pop() uint8 {
	c.SP++
	return c.ReadSRAM(c.SP)
}

func (c *CPU) TickPeripheralsHelper(cycles uint64) {
	if c.TickPeripherals != nil {
		c.TickPeripherals(cycles)
	}
}
