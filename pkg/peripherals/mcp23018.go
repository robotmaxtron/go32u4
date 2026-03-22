package peripherals

// MCP23018 Registers
const (
	MCP23018_IODIRA   = 0x00
	MCP23018_IODIRB   = 0x01
	MCP23018_IPOLA    = 0x02
	MCP23018_IPOLB    = 0x03
	MCP23018_GPINTENA = 0x04
	MCP23018_GPINTENB = 0x05
	MCP23018_DEFVALA  = 0x06
	MCP23018_DEFVALB  = 0x07
	MCP23018_INTCONA  = 0x08
	MCP23018_INTCONB  = 0x09
	MCP23018_IOCON    = 0x0A
	MCP23018_GPPUA    = 0x0C
	MCP23018_GPPUB    = 0x0D
	MCP23018_INTFA    = 0x0E
	MCP23018_INTFB    = 0x0F
	MCP23018_INTCAPA  = 0x10
	MCP23018_INTCAPB  = 0x11
	MCP23018_GPIOA    = 0x12
	MCP23018_GPIOB    = 0x13
	MCP23018_OLATA    = 0x14
	MCP23018_OLATB    = 0x15
)

type MCP23018 struct {
	addr     uint8
	regs     [0x16]uint8
	selected uint8
	External uint16 // External pin states (high/low)
	manager  *Manager
}

func NewMCP23018(addr uint8, manager *Manager) *MCP23018 {
	mcp := &MCP23018{
		addr:     addr,
		manager:  manager,
		selected: 0xFF,
	}
	mcp.regs[MCP23018_IODIRA] = 0xFF
	mcp.regs[MCP23018_IODIRB] = 0xFF
	mcp.External = 0xFFFF
	return mcp
}

func (mcp *MCP23018) Address() uint8 {
	return mcp.addr
}

func (mcp *MCP23018) OnStart(isRead bool) bool {
	// If it's a read, we should prepare the data in TWDR (handled by OnRead in this design, but OnStart can prepare)
	return true // Always ACK for its address
}

func (mcp *MCP23018) OnWrite(data uint8) bool {
	if mcp.selected == 0xFF {
		mcp.selected = data
	} else {
		reg := mcp.getRegAddr(mcp.selected)
		if reg < 0x16 {
			// Some registers are read-only or have read-only bits
			if reg == MCP23018_INTFA || reg == MCP23018_INTFB ||
				reg == MCP23018_INTCAPA || reg == MCP23018_INTCAPB ||
				reg == MCP23018_GPIOA || reg == MCP23018_GPIOB {
				// Read-only
			} else {
				mcp.regs[reg] = data
			}
			// Auto-increment register address
			if mcp.selected < 0xFF {
				mcp.selected++
			}
		}
	}
	return true
}

func (mcp *MCP23018) OnRead() uint8 {
	res := mcp.prepareRead()
	// Auto-increment
	if mcp.selected < 0xFF {
		mcp.selected++
	}
	return res
}

func (mcp *MCP23018) OnStop() {
	mcp.selected = 0xFF
}

func (mcp *MCP23018) getRegAddr(selected uint8) uint8 {
	bank := (mcp.regs[MCP23018_IOCON] >> 7) & 0x01
	if bank == 0 {
		return selected
	}
	// Bank 1 mapping
	if selected <= 0x0A {
		return selected * 2 // Port A registers
	}
	if selected >= 0x10 && selected <= 0x1A {
		return (selected-0x10)*2 + 1 // Port B registers
	}
	return selected
}

func (mcp *MCP23018) Regs(idx uint8) uint8 {
	return mcp.regs[idx]
}

func (mcp *MCP23018) Selected() uint8 {
	return mcp.selected
}

func (mcp *MCP23018) prepareRead() uint8 {
	reg := mcp.getRegAddr(mcp.selected)
	if reg == MCP23018_GPIOA || reg == MCP23018_GPIOB {
		// Read actual pin state
		var iodir, gppu, olat uint8
		var ext uint8
		if reg == MCP23018_GPIOA {
			iodir = mcp.regs[MCP23018_IODIRA]
			gppu = mcp.regs[MCP23018_GPPUA]
			olat = mcp.regs[MCP23018_OLATA]
			ext = uint8(mcp.External & 0xFF)
		} else {
			iodir = mcp.regs[MCP23018_IODIRB]
			gppu = mcp.regs[MCP23018_GPPUB]
			olat = mcp.regs[MCP23018_OLATB]
			ext = uint8(mcp.External >> 8)
		}

		res := uint8(0)
		for i := uint(0); i < 8; i++ {
			isInput := (iodir & (1 << i)) != 0
			pinExt := (ext & (1 << i)) != 0
			hasInternalPullUp := (gppu & (1 << i)) != 0
			hasExternalPullUp := mcp.manager.PullUpResistor < 1000000.0

			if isInput {
				// Input pin
				if pinExt && (hasInternalPullUp || hasExternalPullUp) {
					res |= (1 << i)
				} else if pinExt && !hasInternalPullUp && !hasExternalPullUp {
					// Floating pin, return 0 (as per original logic)
					res &= ^uint8(1 << i)
				} else if !pinExt {
					// Driven low externally
					res &= ^uint8(1 << i)
				}
			} else {
				// Output pin (Open-drain)
				isHigh := (olat & (1 << i)) != 0
				if isHigh {
					// High-Z, follows external and pull-up
					if pinExt || hasInternalPullUp || hasExternalPullUp {
						res |= (1 << i)
					}
				} else {
					// Driven low
					res &= ^uint8(1 << i)
				}
			}
		}

		var ipol uint8
		if reg == MCP23018_GPIOA {
			ipol = mcp.regs[MCP23018_IPOLA]
		} else {
			ipol = mcp.regs[MCP23018_IPOLB]
		}
		res ^= ipol

		oldGPIO := mcp.regs[reg]
		changed := oldGPIO ^ res
		var gpinten, defval, intcon uint8
		if reg == MCP23018_GPIOA {
			gpinten = mcp.regs[MCP23018_GPINTENA]
			defval = mcp.regs[MCP23018_DEFVALA]
			intcon = mcp.regs[MCP23018_INTCONA]
		} else {
			gpinten = mcp.regs[MCP23018_GPINTENB]
			defval = mcp.regs[MCP23018_DEFVALB]
			intcon = mcp.regs[MCP23018_INTCONB]
		}

		for i := uint(0); i < 8; i++ {
			if (gpinten & (1 << i)) != 0 {
				trigger := false
				if (intcon & (1 << i)) != 0 {
					bitVal := (res >> i) & 0x01
					defBit := (defval >> i) & 0x01
					if bitVal != defBit {
						trigger = true
					}
				} else {
					if (changed & (1 << i)) != 0 {
						trigger = true
					}
				}

				if trigger {
					if reg == MCP23018_GPIOA {
						mcp.regs[MCP23018_INTFA] |= (1 << i)
						mcp.regs[MCP23018_INTCAPA] = res
					} else {
						mcp.regs[MCP23018_INTFB] |= (1 << i)
						mcp.regs[MCP23018_INTCAPB] = res
					}
				}
			}
		}

		mcp.regs[reg] = res
		return res
	} else if reg < 0x16 {
		return mcp.regs[reg]
	}
	return 0xFF
}
