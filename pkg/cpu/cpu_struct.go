package cpu

import (
	"go32u4/pkg/bus"
)

const (
	SREG_C = 0
	SREG_Z = 1
	SREG_N = 2
	SREG_V = 3
	SREG_S = 4
	SREG_H = 5
	SREG_T = 6
	SREG_I = 7
)

const (
	IORegSize    = 64
	ExtIORegSize = 160
)

type CPU struct {
	Reg    [32]uint8
	PC     uint16
	SP     uint16
	SREG   uint8
	Cycles uint64

	Bus                 bus.Bus
	InterruptController bus.InterruptController
	TickPeripherals     func(cycles uint64)
	Halted              bool
}

func NewCPU(b bus.Bus, ic bus.InterruptController) *CPU {
	return &CPU{
		SP:                  uint16(IORegSize + ExtIORegSize + 2560 - 1), // Default top of SRAM for 32u4
		Bus:                 b,
		InterruptController: ic,
	}
}
