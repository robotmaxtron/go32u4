package cpu

import (
	"fmt"
)

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
	// Standard fetch-execute: opcode is fetched, PC increments, then execute.
	c.PC++

	cycles := c.Execute(opcode)
	c.Cycles += uint64(cycles)
	if c.TickPeripherals != nil {
		c.TickPeripherals(uint64(cycles))
	}
	return nil
}

func (c *CPU) SetFlag(bit uint8, value bool) {
	if value {
		c.SREG |= 1 << bit
	} else {
		c.SREG &= ^(1 << bit)
	}
}

func (c *CPU) GetFlag(bit uint8) bool {
	return (c.SREG & (1 << bit)) != 0
}

func (c *CPU) Push(val uint8) {
	c.Bus.WriteSRAM(c.SP, val)
	c.SP--
}

func (c *CPU) Pop() uint8 {
	c.SP++
	return c.Bus.ReadSRAM(c.SP)
}

func (c *CPU) ReadSRAM(address uint16) uint8 {
	return c.Bus.ReadSRAM(address)
}

func (c *CPU) WriteSRAM(address uint16, value uint8) {
	c.Bus.WriteSRAM(address, value)
}

func (c *CPU) ReadIO(address uint16) uint8 {
	return c.Bus.ReadIO(address)
}

func (c *CPU) WriteIO(address uint16, value uint8) {
	c.Bus.WriteIO(address, value)
}
