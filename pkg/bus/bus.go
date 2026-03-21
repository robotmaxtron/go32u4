package bus

// Bus represents the unified memory and I/O bus that connects the CPU to peripherals and memory.
type Bus interface {
	ReadSRAM(address uint16) uint8
	WriteSRAM(address uint16, value uint8)
	ReadIO(address uint16) uint8
	WriteIO(address uint16, value uint8)
	Flash() []uint16
	FlashWrite(address uint16, value uint16)
	FlashErase(address uint16)
	FlashCommit(address uint16)
}

// IOCallback is a function that handles specialized I/O peripheral behavior.
type IOCallback func(address uint16, value uint8, isWrite bool) uint8

// PinCallback is a function called when a GPIO pin state might have changed.
// It is called with the port name ('B', 'C', 'D', 'E', 'F'), the mask of changed bits, and the new PORT value.
type PinCallback func(port int8, mask uint8, value uint8)

// InterruptController defines the interface for triggering and managing interrupts.
type InterruptController interface {
	TriggerInterrupt(vector uint8)
	ClearInterrupt(vector uint8)
	SetGlobalInterrupts(enabled bool)
	GetGlobalInterrupts() bool
}
