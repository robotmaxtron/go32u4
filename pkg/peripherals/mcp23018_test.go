package peripherals

import (
	"testing"
)

type mockSys struct {
	ioRegs []uint8
}

func (s *mockSys) IORegs() []uint8 {
	return s.ioRegs
}

func (s *mockSys) TriggerInterrupt(vector uint8) {}
func (s *mockSys) Cycles() uint64                { return 0 }
func (s *mockSys) SaveEEPROM() error             { return nil }
func (s *mockSys) PinCallback(port int8, mask uint8, value uint8) {}
func (s *mockSys) FlashWrite(address uint16, value uint16)        {}
func (s *mockSys) FlashErase(address uint16)                     {}

func setupMCP23018() (*Manager, *MCP23018, []uint8) {
	ioRegs := make([]uint8, 256)
	sys := &mockSys{ioRegs: ioRegs}
	m := NewManager(sys)
	mcp := NewMCP23018(0x20, m)
	m.RegisterTWIClient(mcp)
	return m, mcp, ioRegs
}

func (m *Manager) updateTWI(ioRegs []uint8) {
	m.updateTWIState(ioRegs)
}

func TestMCP23018_BasicReadWrite(t *testing.T) {
	m, mcp, ioRegs := setupMCP23018()

	// 1. Start I2C (SLA+W)
	m.TWIState = 0x08 // START transmitted
	ioRegs[TWDR] = (mcp.addr << 1) | 0
	m.updateTWI(ioRegs)
	if m.TWIState != 0x18 {
		t.Errorf("Expected TWIState 0x18 (SLA+W ACK), got 0x%02X", m.TWIState)
	}

	// 2. Select Register (IODIRA)
	ioRegs[TWDR] = MCP23018_IODIRA
	m.updateTWI(ioRegs)
	if m.TWIState != 0x28 {
		t.Errorf("Expected TWIState 0x28 (Data ACK), got 0x%02X", m.TWIState)
	}
	if mcp.selected != MCP23018_IODIRA {
		t.Errorf("Expected MCP23018_Selected %d, got %d", MCP23018_IODIRA, mcp.selected)
	}

	// 3. Write to Register (IODIRA = 0x55)
	ioRegs[TWDR] = 0x55
	m.updateTWI(ioRegs)
	if m.TWIState != 0x28 {
		t.Errorf("Expected TWIState 0x28 (Data ACK), got 0x%02X", m.TWIState)
	}
	if mcp.regs[MCP23018_IODIRA] != 0x55 {
		t.Errorf("Expected IODIRA 0x55, got 0x%02X", mcp.regs[MCP23018_IODIRA])
	}

	// 4. Auto-increment check (next write should be to IODIRB)
	ioRegs[TWDR] = 0xAA
	m.updateTWI(ioRegs)
	if mcp.regs[MCP23018_IODIRB] != 0xAA {
		t.Errorf("Expected IODIRB 0xAA, got 0x%02X", mcp.regs[MCP23018_IODIRB])
	}
}

func TestMCP23018_BankSwitching(t *testing.T) {
	m, mcp, ioRegs := setupMCP23018()

	// 1. Write IOCON.BANK = 1
	// SLA+W
	m.TWIState = 0x08
	ioRegs[TWDR] = (mcp.addr << 1) | 0
	m.updateTWI(ioRegs)
	// Select IOCON
	ioRegs[TWDR] = MCP23018_IOCON
	m.updateTWI(ioRegs)
	// Set BANK=1 (0x80)
	ioRegs[TWDR] = 0x80
	m.updateTWI(ioRegs)
	
	mcp.OnStop() // Reset selected address for next transaction

	if (mcp.regs[MCP23018_IOCON] & 0x80) == 0 {
		t.Fatal("Failed to set IOCON.BANK")
	}

	// 2. Test Bank 1 Mapping
	// Write to selected 0x01 in Bank 1
	m.TWIState = 0x08
	ioRegs[TWDR] = (mcp.addr << 1) | 0
	m.updateTWI(ioRegs)
	ioRegs[TWDR] = 0x01
	m.updateTWI(ioRegs) // MCP23018_Selected = 0x01
	ioRegs[TWDR] = 0xAA
	m.updateTWI(ioRegs) // Should write to IPOLA (reg 0x02)

	if mcp.regs[MCP23018_IPOLA] != 0xAA {
		t.Errorf("Expected IPOLA 0xAA (Bank 1), got 0x%02X", mcp.regs[MCP23018_IPOLA])
	}

	// Write to selected 0x11 in Bank 1
	mcp.OnStop() // Reset selected for next transaction
	m.TWIState = 0x08
	ioRegs[TWDR] = (mcp.addr << 1) | 0
	m.updateTWI(ioRegs)
	ioRegs[TWDR] = 0x11
	m.updateTWI(ioRegs)
	ioRegs[TWDR] = 0xBB
	m.updateTWI(ioRegs) // Should write to IPOLB (reg 0x03)

	if mcp.regs[MCP23018_IPOLB] != 0xBB {
		t.Errorf("Expected IPOLB 0xBB (Bank 1), got 0x%02X", mcp.regs[MCP23018_IPOLB])
	}
}

func TestMCP23018_ReadOtherRegisters(t *testing.T) {
	_, mcp, _ := setupMCP23018()

	// 1. Write to IODIRA
	mcp.regs[MCP23018_IODIRA] = 0x12
	// 2. Read from IODIRA
	mcp.selected = MCP23018_IODIRA
	val := mcp.OnRead()
	if val != 0x12 {
		t.Errorf("Expected IODIRA 0x12, got 0x%02X", val)
	}

	// 3. Read from invalid register
	mcp.selected = 0xFF
	val = mcp.OnRead()
	if val != 0xFF {
		t.Errorf("Expected 0xFF for invalid register, got 0x%02X", val)
	}
}

func TestMCP23018_PinInput(t *testing.T) {
	m, mcp, ioRegs := setupMCP23018()

	// 1. Configure IODIRA as input (0xFF, default)
	// 2. Set external state (bit 0 low)
	mcp.External = 0xFFFE
	// 3. Set GPPUA bit 0 (internal pull-up)
	mcp.regs[MCP23018_GPPUA] = 0x01

	// SLA+R
	m.TWIState = 0x08
	ioRegs[TWDR] = (mcp.addr << 1) | 1
	mcp.selected = MCP23018_GPIOA
	m.updateTWI(ioRegs)

	if ioRegs[TWDR] != 0xFE {
		t.Errorf("Expected GPIOA 0xFE, got 0x%02X", ioRegs[TWDR])
	}

	// 4. Set external state (bit 0 high-Z, emulated by 1)
	mcp.External = 0xFFFF
	
	mcp.selected = MCP23018_GPIOA
	val := mcp.OnRead()
	// Pin 0 is external 1 AND internal pull-up is 1 -> 1
	if val != 0xFF {
		t.Errorf("Expected GPIOA 0xFF, got 0x%02X", val)
	}

	// 5. Disable pull-up, but keep external high (high-Z)
	mcp.regs[MCP23018_GPPUA] = 0x00
	m.PullUpResistor = 1000001.0 // Disable external pull-up emulator
	mcp.External = 0xFFFF
	val = mcp.OnRead()
	// Pin 0 is external 1 BUT no pull-up -> 0 (floating case in implementation)
	if val != 0x00 {
		t.Errorf("Expected GPIOA 0x00 (all floating), got 0x%02X", val)
	}

	// 6. Enable external pull-up emulator
	m.PullUpResistor = 2200.0 // Default 2.2k
	mcp.External = 0xFFFF
	mcp.regs[MCP23018_IODIRA] = 0xFF
	mcp.selected = MCP23018_GPIOA
	val = mcp.OnRead()
	// Now all bits have external pull-up.
	if val != 0xFF {
		t.Errorf("Expected GPIOA 0xFF (external pull-up), got 0x%02X, IODIRA: 0x%02X, PullUp: %f", val, mcp.regs[MCP23018_IODIRA], m.PullUpResistor)
	}
}

func TestMCP23018_Interrupts(t *testing.T) {
	_, mcp, _ := setupMCP23018()

	// Enable interrupt-on-change for bit 0
	mcp.regs[MCP23018_GPINTENA] = 0x01
	// INTCON = 0 (compare against previous value)
	mcp.regs[MCP23018_INTCONA] = 0x00
	
	// Initial state: bit 0 is 1 (pull-up)
	mcp.External = 0xFFFF
	mcp.regs[MCP23018_GPPUA] = 0xFF
	mcp.selected = MCP23018_GPIOA
	_ = mcp.OnRead() // Current value becomes 1
	
	// Change bit 0 to 0
	mcp.External = 0xFFFE
	mcp.selected = MCP23018_GPIOA
	_ = mcp.OnRead()
	
	if (mcp.regs[MCP23018_INTFA] & 0x01) == 0 {
		t.Errorf("Expected INTFA bit 0 to be set")
	}
	if mcp.regs[MCP23018_INTCAPA] != 0xFE {
		t.Errorf("Expected INTCAPA 0xFE, got 0x%02X", mcp.regs[MCP23018_INTCAPA])
	}
	
	// Test INTCON = 1 (compare against DEFVAL)
	mcp.regs[MCP23018_INTFA] = 0
	mcp.regs[MCP23018_INTCONA] = 0x01
	mcp.regs[MCP23018_DEFVALA] = 0xFF // Expect all bits 1
	
	// bit 0 is currently 0. DEFVAL is 1. Should trigger.
	mcp.selected = MCP23018_GPIOA
	_ = mcp.OnRead()
	if (mcp.regs[MCP23018_INTFA] & 0x01) == 0 {
		t.Errorf("Expected INTFA bit 0 to be set (DEFVAL mismatch)")
	}
}

func TestMCP23018_Output(t *testing.T) {
	_, mcp, _ := setupMCP23018()

	// Set Port A as output
	mcp.regs[MCP23018_IODIRA] = 0x00
	// Set OLATA = 0x55 (drive low for bits 1, 3, 5, 7)
	mcp.regs[MCP23018_OLATA] = 0x55
	
	mcp.selected = MCP23018_GPIOA
	mcp.regs[MCP23018_GPPUA] = 0xFF // Pull-ups enabled
	val := mcp.OnRead()
	
	if val != 0x55 {
		t.Errorf("Expected GPIOA 0x55, got 0x%02X", val)
	}
	
	// Set OLATA = 0x00 (all driven low)
	mcp.regs[MCP23018_OLATA] = 0x00
	mcp.selected = MCP23018_GPIOA
	val = mcp.OnRead()
	if val != 0x00 {
		t.Errorf("Expected GPIOA 0x00, got 0x%02X", val)
	}
}

func TestMCP23018_I2C_ErrorStates(t *testing.T) {
	m, mcp, ioRegs := setupMCP23018()

	// 1. Invalid Address SLA+W
	m.TWIState = 0x08 // START transmitted
	invalidAddr := mcp.addr + 1
	ioRegs[TWDR] = (invalidAddr << 1) | 0
	m.updateTWI(ioRegs)
	if m.TWIState != 0x20 {
		t.Errorf("Expected TWIState 0x20 (SLA+W NACK), got 0x%02X", m.TWIState)
	}

	// 2. Invalid Address SLA+R
	m.TWIState = 0x08 // START transmitted
	ioRegs[TWDR] = (invalidAddr << 1) | 1
	m.updateTWI(ioRegs)
	if m.TWIState != 0x48 {
		t.Errorf("Expected TWIState 0x48 (SLA+R NACK), got 0x%02X", m.TWIState)
	}

	// 3. No Pull-Up Resistor
	m.PullUpResistor = 1000001.0 // Above 1M threshold
	m.TWIState = 0x08
	ioRegs[TWDR] = (mcp.addr << 1) | 0
	m.updateTWI(ioRegs)
	if m.TWIState != 0xF8 {
		t.Errorf("Expected TWIState 0xF8 (No Pull-Up Error), got 0x%02X", m.TWIState)
	}

	// 4. Write to Read-Only Register (INTFA)
	m.TWIState = 0x08
	ioRegs[TWDR] = (mcp.addr << 1) | 0
	m.updateTWI(ioRegs)
	ioRegs[TWDR] = MCP23018_INTFA
	m.updateTWI(ioRegs)
	mcp.regs[MCP23018_INTFA] = 0x00
	ioRegs[TWDR] = 0xFF
	m.updateTWI(ioRegs) // Should NOT write to MCP23018_Regs
	if mcp.regs[MCP23018_INTFA] != 0x00 {
		t.Errorf("Expected INTFA 0x00 (read-only), got 0x%02X", mcp.regs[MCP23018_INTFA])
	}

	// 5. SLA+W NACK state
	m.TWIState = 0x20
	m.updateTWI(ioRegs)
	if m.TWIState != 0xF8 {
		t.Errorf("Expected TWIState 0xF8 after SLA+W NACK, got 0x%02X", m.TWIState)
	}

	// 6. SLA+R NACK state
	m.TWIState = 0x48
	m.updateTWI(ioRegs)
	if m.TWIState != 0xF8 {
		t.Errorf("Expected TWIState 0xF8 after SLA+R NACK, got 0x%02X", m.TWIState)
	}
}

func TestMCP23018_I2C_Stop(t *testing.T) {
	m, _, ioRegs := setupMCP23018()

	// 1. Simulate STOP condition via TWCR
	// Bit 4 is TWSTO (STOP condition)
	m.IOCallback(TWCR, (1<<2)|(1<<4), true) // TWEN=1, TWSTO=1
	if m.TWIState != 0xF8 {
		t.Errorf("Expected TWIState 0xF8 after STOP, got 0x%02X", m.TWIState)
	}
	if ioRegs[TWSR]&0xF8 != 0xF8 {
		t.Errorf("Expected TWSR 0xF8, got 0x%02X", ioRegs[TWSR])
	}
	// TWSTO should be cleared by hardware
	if ioRegs[TWCR]&(1<<4) != 0 {
		t.Errorf("Expected TWSTO bit to be cleared")
	}
}

func TestMCP23018_Regs_Selected(t *testing.T) {
	_, mcp, _ := setupMCP23018()

	if mcp.Address() != 0x20 {
		t.Errorf("Expected address 0x20, got %02X", mcp.Address())
	}

	if mcp.Selected() != 0xFF {
		t.Errorf("Expected selected 0xFF, got %02X", mcp.Selected())
	}

	// Test Regs access
	mcp.regs[MCP23018_IODIRA] = 0x55
	if mcp.Regs(MCP23018_IODIRA) != 0x55 {
		t.Errorf("Expected Regs(IODIRA) 0x55, got %02X", mcp.Regs(MCP23018_IODIRA))
	}
}
