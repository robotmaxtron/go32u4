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

func (s *mockSys) TriggerInterrupt(vector uint8) {
	// Not needed for current tests as MCP23018 doesn't trigger MCU interrupts directly in peripherals.go
}

func (s *mockSys) Cycles() uint64 {
	return 0
}

func (s *mockSys) SaveEEPROM() error {
	return nil
}

func (s *mockSys) PinCallback(port int8, mask uint8, value uint8) {}
func (s *mockSys) FlashWrite(address uint16, value uint16)        {}
func (s *mockSys) FlashErase(address uint16)                     {}

func setupMCP23018() (*Manager, []uint8) {
	ioRegs := make([]uint8, 256)
	sys := &mockSys{ioRegs: ioRegs}
	m := NewManager(sys)
	return m, ioRegs
}

func (m *Manager) updateTWI(ioRegs []uint8) {
	m.updateTWIState(ioRegs)
}

func TestMCP23018_BasicReadWrite(t *testing.T) {
	m, ioRegs := setupMCP23018()

	// 1. Start I2C (SLA+W)
	m.TWIState = 0x08 // START transmitted
	ioRegs[TWDR] = (m.MCP23018_Addr << 1) | 0
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
	if m.MCP23018_Selected != MCP23018_IODIRA {
		t.Errorf("Expected MCP23018_Selected %d, got %d", MCP23018_IODIRA, m.MCP23018_Selected)
	}

	// 3. Write to Register (IODIRA = 0x55)
	ioRegs[TWDR] = 0x55
	m.updateTWI(ioRegs)
	if m.TWIState != 0x28 {
		t.Errorf("Expected TWIState 0x28 (Data ACK), got 0x%02X", m.TWIState)
	}
	if m.MCP23018_Regs[MCP23018_IODIRA] != 0x55 {
		t.Errorf("Expected IODIRA 0x55, got 0x%02X", m.MCP23018_Regs[MCP23018_IODIRA])
	}

	// 4. Auto-increment check (next write should be to IODIRB)
	ioRegs[TWDR] = 0xAA
	m.updateTWI(ioRegs)
	if m.MCP23018_Regs[MCP23018_IODIRB] != 0xAA {
		t.Errorf("Expected IODIRB 0xAA, got 0x%02X", m.MCP23018_Regs[MCP23018_IODIRB])
	}
	if m.MCP23018_Selected != MCP23018_IODIRB+1 {
		t.Errorf("Expected MCP23018_Selected %d, got %d", MCP23018_IODIRB+1, m.MCP23018_Selected)
	}
}

func TestMCP23018_BankSwitching(t *testing.T) {
	m, ioRegs := setupMCP23018()

	// 1. Write IOCON.BANK = 1
	// SLA+W
	m.TWIState = 0x08
	ioRegs[TWDR] = (m.MCP23018_Addr << 1) | 0
	m.updateTWI(ioRegs)
	// Select IOCON
	ioRegs[TWDR] = MCP23018_IOCON
	m.updateTWI(ioRegs)
	// Set BANK=1 (0x80)
	ioRegs[TWDR] = 0x80
	m.updateTWI(ioRegs)

	if (m.MCP23018_Regs[MCP23018_IOCON] & 0x80) == 0 {
		t.Fatal("Failed to set IOCON.BANK")
	}

	// 2. Test Bank 1 Mapping
	// Write to selected 0x01 in Bank 1
	m.TWIState = 0x08
	ioRegs[TWDR] = (m.MCP23018_Addr << 1) | 0
	m.updateTWI(ioRegs)
	ioRegs[TWDR] = 0x01
	m.updateTWI(ioRegs) // MCP23018_Selected = 0x01
	ioRegs[TWDR] = 0xAA
	m.updateTWI(ioRegs) // Should write to IPOLA (reg 0x02)

	if m.MCP23018_Regs[MCP23018_IPOLA] != 0xAA {
		t.Errorf("Expected IPOLA 0xAA (Bank 1), got 0x%02X", m.MCP23018_Regs[MCP23018_IPOLA])
	}

	// Write to selected 0x11 in Bank 1
	m.TWIState = 0x08
	ioRegs[TWDR] = (m.MCP23018_Addr << 1) | 0
	m.updateTWI(ioRegs)
	ioRegs[TWDR] = 0x11
	m.updateTWI(ioRegs)
	ioRegs[TWDR] = 0xBB
	m.updateTWI(ioRegs) // Should write to IPOLB (reg 0x03)

	if m.MCP23018_Regs[MCP23018_IPOLB] != 0xBB {
		t.Errorf("Expected IPOLB 0xBB (Bank 1), got 0x%02X", m.MCP23018_Regs[MCP23018_IPOLB])
	}

	// 3. Test Invalid Address in Bank 1 (should return selected as is)
	m.MCP23018_Selected = 0xFF
	addr := m.getMCP23018RegAddr(0xFF)
	if addr != 0xFF {
		t.Errorf("Expected invalid addr 0xFF to return 0xFF, got 0x%02X", addr)
	}
}

func TestMCP23018_ReadOtherRegisters(t *testing.T) {
	m, ioRegs := setupMCP23018()

	// 1. Write to IODIRA
	m.MCP23018_Regs[MCP23018_IODIRA] = 0x12
	// 2. Read from IODIRA
	m.MCP23018_Selected = MCP23018_IODIRA
	m.prepareMCP23018Read(ioRegs)
	if ioRegs[TWDR] != 0x12 {
		t.Errorf("Expected IODIRA 0x12, got 0x%02X", ioRegs[TWDR])
	}

	// 3. Read from invalid register
	m.MCP23018_Selected = 0xFF
	m.prepareMCP23018Read(ioRegs)
	if ioRegs[TWDR] != 0xFF {
		t.Errorf("Expected 0xFF for invalid register, got 0x%02X", ioRegs[TWDR])
	}
}

func TestMCP23018_PinInput(t *testing.T) {
	m, ioRegs := setupMCP23018()

	// 1. Configure IODIRA as input (0xFF, default)
	// 2. Set external state (bit 0 low)
	m.MCP23018_External = 0xFFFE
	// 3. Set GPPUA bit 0 (internal pull-up)
	m.MCP23018_Regs[MCP23018_GPPUA] = 0x01

	// SLA+R
	m.TWIState = 0x08
	ioRegs[TWDR] = (m.MCP23018_Addr << 1) | 1
	m.MCP23018_Selected = MCP23018_GPIOA
	m.updateTWI(ioRegs)

	// Pin 0 is external 0, so result should be 0 despite pull-up.
	// IPOL is default 0.
	// Pin 0 result = 0. Bits 1-7 result = 1 (external 1 AND external pull-up 10k).
	// So ioRegs[TWDR] should be 0xFE.
	if ioRegs[TWDR] != 0xFE {
		t.Errorf("Expected GPIOA 0xFE, got 0x%02X", ioRegs[TWDR])
	}

	// 4. Set external state (bit 0 high-Z, emulated by 1)
	m.MCP23018_External = 0xFFFF
	// updateTWI increments selected if it was 0x40/0x50 state.
	// Initial SLA+R set state to 0x40.
	// Next call to updateTWI (simulating Data ACK) will read next byte.
	m.updateTWI(ioRegs)
	// But wait, if TWIState is 0x40, it prepares read and sets state to 0x50.
	// If TWIState is 0x50, it increments Selected and prepares read.
	// Here Selected was GPIOA. After first updateTWI, it's still GPIOA (prepared).
	// After second updateTWI, Selected becomes GPIOB and prepares read for GPIOB.
	// We want to stay on GPIOA for this test, so let's reset Selected if needed, or just use GPIOB for next checks.
	
	m.MCP23018_Selected = MCP23018_GPIOA
	m.prepareMCP23018Read(ioRegs)
	// Pin 0 is external 1 AND internal pull-up is 1 -> 1
	if ioRegs[TWDR] != 0xFF {
		t.Errorf("Expected GPIOA 0xFF, got 0x%02X", ioRegs[TWDR])
	}

	// 5. Disable pull-up, but keep external high (high-Z)
	m.MCP23018_Regs[MCP23018_GPPUA] = 0x00
	m.PullUpResistor = 1000001.0 // Disable external pull-up emulator
	m.prepareMCP23018Read(ioRegs)
	// Pin 0 is external 1 BUT no pull-up -> 0 (floating case in implementation)
	if ioRegs[TWDR] != 0x00 {
		t.Errorf("Expected GPIOA 0x00 (all floating), got 0x%02X", ioRegs[TWDR])
	}

	// 6. Enable external pull-up emulator
	m.PullUpResistor = 10000.0 // 10k
	m.prepareMCP23018Read(ioRegs)
	// Now all bits have external pull-up.
	if ioRegs[TWDR] != 0xFF {
		t.Errorf("Expected GPIOA 0xFF (external pull-up), got 0x%02X", ioRegs[TWDR])
	}
}

func TestMCP23018_InputPolarity(t *testing.T) {
	m, ioRegs := setupMCP23018()

	// External bits: all high (0xFFFF)
	m.MCP23018_External = 0xFFFF
	m.MCP23018_Regs[MCP23018_GPPUA] = 0xFF
	
	// IPOLA = 0xAA (invert bits 1, 3, 5, 7)
	m.MCP23018_Regs[MCP23018_IPOLA] = 0xAA
	
	// SLA+R GPIOA
	m.MCP23018_Selected = MCP23018_GPIOA
	m.prepareMCP23018Read(ioRegs)
	
	// Input (1) XOR IPOL (0xAA) = 0x55
	if ioRegs[TWDR] != 0x55 {
		t.Errorf("Expected GPIOA 0x55 (IPOL), got 0x%02X", ioRegs[TWDR])
	}
}

func TestMCP23018_Interrupts(t *testing.T) {
	m, ioRegs := setupMCP23018()

	// Enable interrupt-on-change for bit 0
	m.MCP23018_Regs[MCP23018_GPINTENA] = 0x01
	// INTCON = 0 (compare against previous value)
	m.MCP23018_Regs[MCP23018_INTCONA] = 0x00
	
	// Initial state: bit 0 is 1 (pull-up)
	m.MCP23018_External = 0xFFFF
	m.MCP23018_Regs[MCP23018_GPPUA] = 0xFF
	m.MCP23018_Selected = MCP23018_GPIOA
	m.prepareMCP23018Read(ioRegs) // Current value becomes 1
	
	// Change bit 0 to 0
	m.MCP23018_External = 0xFFFE
	m.prepareMCP23018Read(ioRegs)
	
	if (m.MCP23018_Regs[MCP23018_INTFA] & 0x01) == 0 {
		t.Errorf("Expected INTFA bit 0 to be set")
	}
	if m.MCP23018_Regs[MCP23018_INTCAPA] != 0xFE {
		t.Errorf("Expected INTCAPA 0xFE, got 0x%02X", m.MCP23018_Regs[MCP23018_INTCAPA])
	}
	
	// Test INTCON = 1 (compare against DEFVAL)
	m.MCP23018_Regs[MCP23018_INTFA] = 0
	m.MCP23018_Regs[MCP23018_INTCONA] = 0x01
	m.MCP23018_Regs[MCP23018_DEFVALA] = 0xFF // Expect all bits 1
	
	// bit 0 is currently 0. DEFVAL is 1. Should trigger.
	m.prepareMCP23018Read(ioRegs)
	if (m.MCP23018_Regs[MCP23018_INTFA] & 0x01) == 0 {
		t.Errorf("Expected INTFA bit 0 to be set (DEFVAL mismatch)")
	}
}

func TestMCP23018_Output(t *testing.T) {
	m, ioRegs := setupMCP23018()

	// Set Port A as output
	m.MCP23018_Regs[MCP23018_IODIRA] = 0x00
	// Set OLATA = 0x55 (drive low for bits 1, 3, 5, 7)
	m.MCP23018_Regs[MCP23018_OLATA] = 0x55
	
	m.MCP23018_Selected = MCP23018_GPIOA
	m.prepareMCP23018Read(ioRegs)
	
	// MCP23018 is open-drain.
	// Output 1 -> High-Z
	// Output 0 -> Driven Low
	// If it's output 1 and we have pull-up -> 1
	// If it's output 0 -> 0
	
	m.MCP23018_Regs[MCP23018_GPPUA] = 0xFF // Pull-ups enabled
	m.prepareMCP23018Read(ioRegs)
	
	if ioRegs[TWDR] != 0x55 {
		t.Errorf("Expected GPIOA 0x55, got 0x%02X", ioRegs[TWDR])
	}
	
	// Set OLATA = 0x00 (all driven low)
	m.MCP23018_Regs[MCP23018_OLATA] = 0x00
	m.prepareMCP23018Read(ioRegs)
	if ioRegs[TWDR] != 0x00 {
		t.Errorf("Expected GPIOA 0x00, got 0x%02X", ioRegs[TWDR])
	}
}

func TestMCP23018_I2C_ErrorStates(t *testing.T) {
	m, ioRegs := setupMCP23018()

	// 1. Invalid Address SLA+W
	m.TWIState = 0x08 // START transmitted
	invalidAddr := m.MCP23018_Addr + 1
	ioRegs[TWDR] = (invalidAddr << 1) | 0
	m.updateTWI(ioRegs)
	if m.TWIState != 0x20 {
		t.Errorf("Expected TWIState 0x20 (SLA+W NACK), got 0x%02X", m.TWIState)
	}
	if m.MCP23018_Active {
		t.Errorf("Expected MCP23018_Active false for invalid address")
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
	ioRegs[TWDR] = (m.MCP23018_Addr << 1) | 0
	m.updateTWI(ioRegs)
	if m.TWIState != 0xF8 {
		t.Errorf("Expected TWIState 0xF8 (No Pull-Up Error), got 0x%02X", m.TWIState)
	}
	if m.MCP23018_Active {
		t.Errorf("Expected MCP23018_Active false when no pull-up")
	}

	// 4. Write to Read-Only Register (INTFA)
	m.TWIState = 0x08
	ioRegs[TWDR] = (m.MCP23018_Addr << 1) | 0
	m.updateTWI(ioRegs)
	ioRegs[TWDR] = MCP23018_INTFA
	m.updateTWI(ioRegs)
	m.MCP23018_Regs[MCP23018_INTFA] = 0x00
	ioRegs[TWDR] = 0xFF
	m.updateTWI(ioRegs) // Should NOT write to MCP23018_Regs
	if m.MCP23018_Regs[MCP23018_INTFA] != 0x00 {
		t.Errorf("Expected INTFA 0x00 (read-only), got 0x%02X", m.MCP23018_Regs[MCP23018_INTFA])
	}

	// 5. Bounds check on auto-increment
	m.MCP23018_Selected = 0xFF
	m.TWIState = 0x28 // Data ACK (Write)
	ioRegs[TWDR] = 0x55
	m.updateTWI(ioRegs)
	if m.MCP23018_Selected != 0xFF {
		t.Errorf("Expected MCP23018_Selected 0xFF (no increment), got 0x%02X", m.MCP23018_Selected)
	}

	// 6. SLA+W NACK state
	m.TWIState = 0x20
	m.updateTWI(ioRegs)
	if m.TWIState != 0xF8 {
		t.Errorf("Expected TWIState 0xF8 after SLA+W NACK, got 0x%02X", m.TWIState)
	}

	// 7. SLA+R NACK state
	m.TWIState = 0x48
	m.updateTWI(ioRegs)
	if m.TWIState != 0xF8 {
		t.Errorf("Expected TWIState 0xF8 after SLA+R NACK, got 0x%02X", m.TWIState)
	}

	// 8. SLA+W Active check
	m.MCP23018_Active = false
	m.TWIState = 0x18
	m.updateTWI(ioRegs)
	if m.TWIState != 0x20 {
		t.Errorf("Expected TWIState 0x20 (SLA+W NACK) for inactive, got 0x%02X", m.TWIState)
	}

	// 9. SLA+R Active check
	m.MCP23018_Active = false
	m.TWIState = 0x40
	m.updateTWI(ioRegs)
	if m.TWIState != 0x48 {
		t.Errorf("Expected TWIState 0x48 (SLA+R NACK) for inactive, got 0x%02X", m.TWIState)
	}
}

func TestMCP23018_I2C_Stop(t *testing.T) {
	m, ioRegs := setupMCP23018()

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
