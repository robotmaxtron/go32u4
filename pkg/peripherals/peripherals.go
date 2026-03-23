package peripherals

const (
	SYS_FREQ = 16000000 // 16 MHz
	WDT_FREQ = 128000   // 128 kHz
)

// ATmega32u4 IO Addresses (0x00-0x3F) and Extended IO (0x60-0xFF)
const (
	PINB  = 0x03
	DDRB  = 0x04
	PORTB = 0x05

	PINC  = 0x06
	DDRC  = 0x07
	PORTC = 0x08

	PIND  = 0x09
	DDRD  = 0x0A
	PORTD = 0x0B

	PINE  = 0x0C
	DDRE  = 0x0D
	PORTE = 0x0E

	PINF  = 0x0F
	DDRF  = 0x10
	PORTF = 0x11

	TCCR0A = 0x24
	TCCR0B = 0x25
	TCNT0  = 0x26
	OCR0A  = 0x27
	OCR0B  = 0x28

	// Timer1
	TCCR1A = 0x80
	TCCR1B = 0x81
	TCCR1C = 0x82
	TCNT1L = 0x84
	TCNT1H = 0x85
	ICR1L  = 0x86
	ICR1H  = 0x87
	OCR1AL = 0x88
	OCR1AH = 0x89
	OCR1BL = 0x8A
	OCR1BH = 0x8B
	OCR1CL = 0x8C
	OCR1CH = 0x8D

	// Timer3
	TCCR3A = 0x90
	TCCR3B = 0x91
	TCCR3C = 0x92
	TCNT3L = 0x94
	TCNT3H = 0x95
	ICR3L  = 0x96
	ICR3H  = 0x97
	OCR3AL = 0x98
	OCR3AH = 0x99
	OCR3BL = 0x9A
	OCR3BH = 0x9B
	OCR3CL = 0x9C
	OCR3CH = 0x9D

	// Interrupt Flags and Masks
	EICRA  = 0x69 // External Interrupt Control Register A
	EICRB  = 0x6A // External Interrupt Control Register B
	EIMSK  = 0x1D // External Interrupt Mask Register
	EIFR   = 0x1C // External Interrupt Flag Register
	PCICR  = 0x68 // Pin Change Interrupt Control Register
	PCIFR  = 0x1B // Pin Change Interrupt Flag Register
	PCMSK0 = 0x6B // Pin Change Mask Register 0
	TIMSK0 = 0x6E // Timer/Counter0 Interrupt Mask Register
	TIFR0  = 0x15 // Timer/Counter0 Interrupt Flag Register
	TIMSK1 = 0x6F // Timer/Counter1 Interrupt Mask Register
	TIFR1  = 0x16 // Timer/Counter1 Interrupt Flag Register
	TIMSK3 = 0x71 // Timer/Counter3 Interrupt Mask Register
	TIFR3  = 0x18 // Timer/Counter3 Interrupt Flag Register
	TIMSK4 = 0x72 // Timer/Counter4 Interrupt Mask Register
	TIFR4  = 0x19 // Timer/Counter4 Interrupt Flag Register

	// EEPROM Registers
	EEARL = 0x21 // EEPROM Address Register Low
	EEARH = 0x22 // EEPROM Address Register High
	EEDR  = 0x20 // EEPROM Data Register
	EECR  = 0x1F // EEPROM Control Register

	// USART1 Registers
	UDR1   = 0xCE
	UCSR1A = 0xC8
	UCSR1B = 0xC9
	UCSR1C = 0xCA
	UBRR1L = 0xCC
	UBRR1H = 0xCD

	// SPI Registers
	SPCR = 0x2C
	SPSR = 0x2D
	SPDR = 0x2E

	// TWI Registers
	TWBR = 0xB8
	TWSR = 0xB9
	TWAR = 0xBA
	TWDR = 0xBB
	TWCR = 0xBC
	TWAMR = 0xBD

	// ADC Registers
	ADMUX  = 0x7C
	ADCSRA = 0x7A
	ADCSRB = 0x7B
	ADCL   = 0x78
	ADCH   = 0x79
	DIDR0  = 0x7E
	DIDR2  = 0x7D

	// Timer 3 Mask Register
	// TIFR3  = 0x18 // Timer/Counter3 Interrupt Flag Register

	// Timer4 Registers
	TCCR4A = 0xC0
	TCCR4B = 0xC1
	TCCR4C = 0xC2
	TCCR4D = 0xC3
	TCCR4E = 0xC4
	TCNT4  = 0xBE
	TC4H   = 0xBF
	OCR4A  = 0xCF
	OCR4B  = 0xD0
	OCR4C  = 0xD1
	OCR4D  = 0xD2
	DT4    = 0xD4

	// PLL Registers
	PLLCSR = 0x29
	PLLE   = 1
	PLOCK  = 0
	PCKE   = 2

	// USB Registers (Simplified CDC mapping)
	UHWCON  = 0x47
	USBCON  = 0x48
	USBSTA  = 0x49
	USBINT  = 0x4A
	UDCON   = 0x4B
	UDINT   = 0x4C
	UDIEN   = 0x4D
	UDADDR  = 0xE8
	UENUM   = 0xE9
	UEINT   = 0xEA
	UEINTX  = 0xEB
	UERST   = 0xEC
	UECONX  = 0x53
	UECFG0X = 0x54
	UECFG1X = 0x55
	UESTA0X = 0x56
	UESTA1X = 0x57
	UEDATX  = 0x58
	UEBCLX  = 0x59
	UEBCHX  = 0x5A

	// Power Management
	SMCR = 0x5B // Sleep Mode Control Register

	// SPM Control and Status Register
	SPMCSR = 0x37
	SPMIE  = 7
	RWWSB  = 6
	SIGRD  = 5
	RWWSRE = 4
	BLBSET = 3
	PGWRT  = 2
	PGERS  = 1
	SPMEN  = 0

	// Watchdog Registers
	WDTCSR = 0x60 // Watchdog Timer Control Register

	// Watchdog Timer Bits
	WDP0 = 0
	WDP1 = 1
	WDP2 = 2
	WDE  = 3
	WDCE = 4
	WDP3 = 5
	WDIE = 6
	WDIF = 7

	ADIF = 4
	
	// Power Reduction Registers
	PRR0 = 0x64
	PRR1 = 0x65
	
	// PRR0 bits
	PRTWI    = 6
	PRTIM0   = 5
	PRTIM1   = 3
	PRSPI    = 2
	PRUSART1 = 1
	PRADC    = 0
	
	// PRR1 bits
	PRUSB  = 7
	PRTIM4 = 4
	PRTIM3 = 3
)

// RegisterMasks defines write masks for I/O registers to protect read-only and reserved bits.
// A '1' in the mask means the bit is writable.
// Any address not in this map will have a default mask of 0xFF.
var RegisterMasks = map[uint16]uint8{
	// PIN/DDR/PORT are handled specially or fully writable
	PINB: 0xFF, DDRB: 0xFF, PORTB: 0xFF,
	PINC: 0xFF, DDRC: 0xFF, PORTC: 0xFF,
	PIND: 0xFF, DDRD: 0xFF, PORTD: 0xFF,
	PINE: 0xFF, DDRE: 0xFF, PORTE: 0xFF,
	PINF: 0xFF, DDRF: 0xFF, PORTF: 0xFF,

	// Timer 0
	TCCR0A: 0xFF, TCCR0B: 0xFF, TCNT0: 0xFF, OCR0A: 0xFF, OCR0B: 0xFF, TIMSK0: 0x07, TIFR0: 0x07,

	// Timer 1
	TCCR1A: 0xFF, TCCR1B: 0xFF, TCCR1C: 0xFF, TCNT1L: 0xFF, TCNT1H: 0xFF,
	ICR1L: 0xFF, ICR1H: 0xFF, OCR1AL: 0xFF, OCR1AH: 0xFF, OCR1BL: 0xFF, OCR1BH: 0xFF,
	OCR1CL: 0xFF, OCR1CH: 0xFF, TIMSK1: 0x2F, TIFR1: 0x2F,

	// Timer 3
	TCCR3A: 0xFF, TCCR3B: 0xFF, TCCR3C: 0xFF, TCNT3L: 0xFF, TCNT3H: 0xFF,
	ICR3L: 0xFF, ICR3H: 0xFF, OCR3AL: 0xFF, OCR3AH: 0xFF, OCR3BL: 0xFF, OCR3BH: 0xFF,
	OCR3CL: 0xFF, OCR3CH: 0xFF, TIMSK3: 0x2F, TIFR3: 0x2F,

	// Timer 4
	TCCR4A: 0xFF, TCCR4B: 0x7F, TCCR4C: 0xFF, TCCR4D: 0xFF, TCCR4E: 0xFF,
	TCNT4: 0xFF, TC4H: 0x03, OCR4A: 0xFF, OCR4B: 0xFF, OCR4C: 0xFF, OCR4D: 0xFF,
	DT4: 0xFF, TIMSK4: 0xEE, TIFR4: 0xEE,

	// TWI
	TWBR: 0xFF, TWSR: 0x03, TWAR: 0xFF, TWDR: 0xFF, TWCR: 0xFF, TWAMR: 0xFE,

	// ADC
	ADMUX: 0xFF, ADCSRA: 0xEF, ADCSRB: 0x7F, ADCL: 0xFF, ADCH: 0xFF, DIDR0: 0xFF, DIDR2: 0xFF,

	// SPI
	SPCR: 0xFF, SPSR: 0x01,

	// UART1
	UCSR1A: 0x03, UCSR1B: 0xFF, UCSR1C: 0xFF, UBRR1L: 0xFF, UBRR1H: 0x0F, UDR1: 0xFF,

	// EEPROM
	EEARL: 0xFF, EEARH: 0x03, EEDR: 0xFF, EECR: 0x3F,

	// PLL
	PLLCSR: 0x1E,

	// USB
	UHWCON: 0x81, USBCON: 0xBF, USBSTA: 0x00, USBINT: 0x00,
	UDCON: 0x03, UDINT: 0xFF, UDIEN: 0xFF, UDADDR: 0xFF,
	UENUM: 0x07, UERST: 0x7F, UECONX: 0xFF, UECFG0X: 0xFF, UECFG1X: 0xFF,
	UESTA0X: 0x00, UESTA1X: 0x00, UEDATX: 0xFF, UEBCLX: 0x00, UEBCHX: 0x00,
	UEINT: 0x00, UEINTX: 0xFF,

	// Misc
	SMCR: 0x0F, SPMCSR: 0xBF, WDTCSR: 0xFF, PRR0: 0xFF, PRR1: 0xFF,
	EICRA: 0xFF, EICRB: 0x33, EIMSK: 0xFF, EIFR: 0xFF, PCICR: 0x01, PCIFR: 0x01, PCMSK0: 0xFF,
}

// MacroRecord represents a single HID report in a macro sequence.
type MacroRecord struct {
	KeyMap [8]uint8 // The HID report to send
	Delay  uint32   // Delay after sending this report (in cycles)
}

// MacroTable defines a sequence of HID reports.
type MacroTable struct {
	Records []MacroRecord
}

const EEPROMSize = 1024

// Endpoint represents a USB endpoint.
type Endpoint struct {
	FIFO      []byte
	SetupFIFO []byte
	Config0   uint8
	Config1   uint8
	Status0   uint8
	Status1   uint8
	Control   uint8
	Interrupt uint8
}

// System represents the core system that peripherals can interact with.
type System interface {
	IORegs() []uint8
	TriggerInterrupt(vector uint8)
	Cycles() uint64
	SaveEEPROM() error
	PinCallback(port int8, mask uint8, value uint8)
	FlashWrite(address uint16, value uint16)
	FlashErase(address uint16)
}

// Manager manages the state and updates of hardware peripherals.
type Manager struct {
	Sys System

	// Timer0
	Timer0Counter  uint8
	Timer0ControlA uint8
	Timer0ControlB uint8
	Timer0CompareA uint8
	Timer0CompareB uint8

	// Timer1
	Timer1Counter  uint16
	Timer1ControlA uint8
	Timer1ControlB uint8
	Timer1CompareA uint16
	Timer1CompareB uint16
	Timer1CompareC uint16
	Timer1InputCap uint16

	// Timer3
	Timer3Counter  uint16
	Timer3ControlA uint8
	Timer3ControlB uint8
	Timer3CompareA uint16
	Timer3CompareB uint16
	Timer3CompareC uint16
	Timer3InputCap uint16

	// EEPROM
	EEPROM           [EEPROMSize]uint8
	EEPROMWriteTimer uint64

	// USART1
	UART1TXBuffer []byte
	UART1RXBuffer []byte

	// SPI
	SPIBuffer uint8

	// TWI
	TWIBuffer      uint8
	TWIState       uint8
	TWIPendingStop bool
	TWIClients     []TWIClient
	ActiveTWI      TWIClient

	// ADC
	ADCValue uint16

	// USB Full Emulation
	USBEndpoints   [7]Endpoint
	USBDeviceState uint8
	USBAddress     uint8
	USBSelectedEP  uint8
	USBConfigured  bool

	// HID State
	HIDKeyMap          [8]uint8 // Simple keymap for emulation
	CapturedHIDReports [][]byte // Reports sent by the firmware

	// SPM State
	SPMBuffer   [64]uint16
	SPMPageAddr uint16
	SPMTimeout  uint8

	// Sleep
	SleepEnabled bool

	// Timer4
	Timer4Counter  uint16
	Timer4HighByte uint8
	Timer4ControlA uint8
	Timer4ControlB uint8
	Timer4ControlC uint8
	Timer4ControlD uint8
	Timer4ControlE uint8
	Timer4OCR4A    uint16
	Timer4OCR4B    uint16
	Timer4OCR4C    uint16
	Timer4OCR4D    uint16
	Timer4DT4      uint8

	// PLL
	PLLControl uint8

	// Watchdog
	WatchdogCycles      uint64
	WatchdogTimeout     uint64
	WatchdogTimedChange uint64
	WatchdogReset       bool

	// Macro State
	MacroActive       bool
	MacroQueue        []MacroRecord
	MacroCurrentIndex int
	MacroDelayCounter uint64

	// I2C/TWI Pull-up
	PullUpResistor float64 // Pull-up resistor value in ohms (default 2200)

	Timer4SubCycles uint64 // For precise clock scaling
}

type TWIClient interface {
	Address() uint8
	OnStart(isRead bool) bool // Returns ACK/NACK
	OnWrite(data uint8) bool  // Returns ACK/NACK
	OnRead() uint8            // Returns data byte
	OnStop()
}

func (m *Manager) RegisterTWIClient(client TWIClient) {
	m.TWIClients = append(m.TWIClients, client)
}

// TriggerMacro queues a sequence of HID reports to be sent.
func (m *Manager) TriggerMacro(table MacroTable) {
	m.MacroQueue = append(m.MacroQueue, table.Records...)
	m.MacroActive = true
}

func (m *Manager) Reset() {
	m.Timer0Counter = 0
	m.Timer0ControlA = 0
	m.Timer0ControlB = 0
	m.Timer1Counter = 0
	m.Timer1ControlA = 0
	m.Timer1ControlB = 0
	m.Timer3Counter = 0
	m.Timer3ControlA = 0
	m.Timer3ControlB = 0
	m.Timer4Counter = 0
	m.Timer4ControlA = 0
	m.Timer4ControlB = 0
	m.WatchdogCycles = 0
	m.WatchdogReset = false
	m.EEPROMWriteTimer = 0
	m.SPMTimeout = 0
	m.SleepEnabled = false
	m.USBSelectedEP = 0
	m.USBConfigured = false
	// Reset FIFOs
	for i := range m.USBEndpoints {
		m.USBEndpoints[i].FIFO = nil
		m.USBEndpoints[i].SetupFIFO = nil
		m.USBEndpoints[i].Interrupt = 0
	}
}

func NewManager(sys System) *Manager {
	m := &Manager{
		Sys: sys,
	}
	// Default I2C pull-up for ErgoDox
	m.PullUpResistor = 2200.0
	m.Reset()
	return m
}

// IOCallback handles the core I/O behavior for ATmega32u4.
func (m *Manager) IOCallback(address uint16, value uint8, isWrite bool) uint8 {
	ioRegs := m.Sys.IORegs()
	if isWrite {
		// Apply register mask
		mask, ok := RegisterMasks[address]
		if !ok {
			mask = 0xFF // Default to fully writable if not in mask table
		}
		
		if mask == 0 && address != TIFR0 && address != TIFR1 && address != TIFR3 && address != TIFR4 &&
			address != EIFR && address != PCIFR && address != ADCSRA && address != UEINTX {
			// If mask is 0 and it's not a special register handled below, it's read-only
			return ioRegs[address]
		}
		
		// For masked registers, apply the mask unless it's handled specially below
		switch address {
		case TIFR0, TIFR1, TIFR3, TIFR4, EIFR, PCIFR, ADCSRA, UEINTX:
			// These are handled specially below (e.g. W1C)
		default:
			if ok {
				value = (value & mask) | (ioRegs[address] & ^mask)
			}
		}

		switch address {
		case TIFR0, TIFR1, TIFR3, TIFR4, EIFR, PCIFR:
			// Write-1-to-Clear
			ioRegs[address] &= ^value
			return ioRegs[address]
		case PORTB, DDRB, PORTC, DDRC, PORTD, DDRD, PORTE, DDRE, PORTF, DDRF:
			oldVal := ioRegs[address]
			ioRegs[address] = value
			portName := int8('B')
			switch address {
			case PORTB, DDRB:
				portName = 'B'
			case PORTC, DDRC:
				portName = 'C'
			case PORTD, DDRD:
				portName = 'D'
			case PORTE, DDRE:
				portName = 'E'
			case PORTF, DDRF:
				portName = 'F'
			}
			m.Sys.PinCallback(portName, oldVal^value, value)
		case PINB, PINC, PIND, PINE, PINF:
			portAddr := address + 2
			oldPort := ioRegs[portAddr]
			newPort := oldPort ^ value
			ioRegs[portAddr] = newPort
			portIdx := (address - PINB) / 3
			portName := int8('B' + portIdx)
			m.Sys.PinCallback(portName, value, newPort)
		case TCNT0:
			m.Timer0Counter = value
		case TCCR0A:
			m.Timer0ControlA = value
		case TCCR0B:
			m.Timer0ControlB = value
		case OCR0A:
			m.Timer0CompareA = value
		case OCR0B:
			m.Timer0CompareB = value
		case TCNT1L:
			m.Timer1Counter = (m.Timer1Counter & 0xFF00) | uint16(value)
		case TCNT1H:
			m.Timer1Counter = (m.Timer1Counter & 0x00FF) | (uint16(value) << 8)
		case OCR1AL:
			m.Timer1CompareA = (m.Timer1CompareA & 0xFF00) | uint16(value)
		case OCR1AH:
			m.Timer1CompareA = (m.Timer1CompareA & 0x00FF) | (uint16(value) << 8)
		case OCR1BL:
			m.Timer1CompareB = (m.Timer1CompareB & 0xFF00) | uint16(value)
		case OCR1BH:
			m.Timer1CompareB = (m.Timer1CompareB & 0x00FF) | (uint16(value) << 8)
		case OCR1CL:
			m.Timer1CompareC = (m.Timer1CompareC & 0xFF00) | uint16(value)
		case OCR1CH:
			m.Timer1CompareC = (m.Timer1CompareC & 0x00FF) | (uint16(value) << 8)
		case TCCR1A:
			m.Timer1ControlA = value
		case TCCR1B:
			m.Timer1ControlB = value
		case TCNT3L:
			m.Timer3Counter = (m.Timer3Counter & 0xFF00) | uint16(value)
		case TCNT3H:
			m.Timer3Counter = (m.Timer3Counter & 0x00FF) | (uint16(value) << 8)
		case OCR3AL:
			m.Timer3CompareA = (m.Timer3CompareA & 0xFF00) | uint16(value)
		case OCR3AH:
			m.Timer3CompareA = (m.Timer3CompareA & 0x00FF) | (uint16(value) << 8)
		case OCR3BL:
			m.Timer3CompareB = (m.Timer3CompareB & 0xFF00) | uint16(value)
		case OCR3BH:
			m.Timer3CompareB = (m.Timer3CompareB & 0x00FF) | (uint16(value) << 8)
		case OCR3CL:
			m.Timer3CompareC = (m.Timer3CompareC & 0xFF00) | uint16(value)
		case OCR3CH:
  	m.Timer3CompareC = (m.Timer3CompareC & 0x00FF) | (uint16(value) << 8)
		case EICRA:
			ioRegs[EICRA] = value
		case EICRB:
			ioRegs[EICRB] = value
		case EIMSK:
			ioRegs[EIMSK] = value
		case PRR0:
			ioRegs[PRR0] = value
		case PRR1:
			ioRegs[PRR1] = value
		case TCCR3A:
			m.Timer3ControlA = value
		case TCCR3B:
			m.Timer3ControlB = value
		case UDR1:
			m.UART1TXBuffer = append(m.UART1TXBuffer, value)
			ioRegs[UCSR1A] |= (1 << 5) | (1 << 6) // UDRE1 and TXC1
			if (ioRegs[UCSR1B] & (1 << 6)) != 0 {
				m.Sys.TriggerInterrupt(27)
			}
			if (ioRegs[UCSR1B] & (1 << 5)) != 0 {
				m.Sys.TriggerInterrupt(26)
			}
		case SPDR:
			m.SPIBuffer = value
			ioRegs[SPSR] |= 1 << 7 // SPIF
			if (ioRegs[SPCR] & (1 << 7)) != 0 {
				m.Sys.TriggerInterrupt(18)
			}
		case ADCSRA:
			// W1C for ADIF (bit 4)
			if value&(1<<ADIF) != 0 {
				ioRegs[ADCSRA] &= ^uint8(1 << ADIF)
			}
			// Normal bits (ADEN, ADSC, ADATE, ADIE, ADPS)
			// ADEN=7, ADSC=6, ADATE=5, ADIE=3, ADPS=2:0
			// mask includes all these bits.
			mask := uint8((1 << 7) | (1 << 6) | (1 << 5) | (1 << 3) | 0x07)
			ioRegs[ADCSRA] = (ioRegs[ADCSRA] & ^mask) | (value & mask)

			if value&(1<<6) != 0 {
				// Simulate conversion start
				ioRegs[ADCSRA] &= ^uint8(1 << 6)
				ioRegs[ADCSRA] |= 1 << ADIF
			}
			if ioRegs[ADCSRA]&(1<<4) != 0 && ioRegs[ADCSRA]&(1<<3) != 0 {
				m.Sys.TriggerInterrupt(29)
			}
		case EECR:
			// Block writes if EEPROM is busy
			if m.EEPROMWriteTimer > 0 {
				return ioRegs[EECR]
			}
			if value&0x01 != 0 {
				addr := (uint16(ioRegs[EEARH]) << 8) | uint16(ioRegs[EEARL])
				if addr < uint16(len(m.EEPROM)) {
					ioRegs[EEDR] = m.EEPROM[addr]
				}
			}
			if value&0x02 != 0 {
				// EEPE: EEPROM Program Enable
				if ioRegs[EECR]&0x04 != 0 { // EEMPE: EEPROM Master Program Enable
					addr := (uint16(ioRegs[EEARH]) << 8) | uint16(ioRegs[EEARL])
					if addr < uint16(len(m.EEPROM)) {
						m.EEPROM[addr] = ioRegs[EEDR]
						_ = m.Sys.SaveEEPROM()
					}
					// Start the write timer (simulate 3.3ms @ SYS_FREQ)
					m.EEPROMWriteTimer = (33 * SYS_FREQ) / 10000 // 3.3ms
					// Clear EEMPE bit as per datasheet (cleared by hardware after 4 cycles, but here we clear it after EEPE is set)
					value &= ^uint8(0x04)
				} else {
					// EEPE is ignored if EEMPE is not set
					value &= ^uint8(0x02)
				}
			}
			ioRegs[EECR] = value
		case EEDR:
			// Block writes if EEPROM is busy
			if m.EEPROMWriteTimer > 0 {
				return ioRegs[EEDR]
			}
			ioRegs[address] = value
		case TWSR:
			// Bits 7:3 are status (read-only), 2 is reserved, 1:0 are prescaler (RW)
			value = (value & 0x03) | (ioRegs[TWSR] & 0xF8)
			ioRegs[address] = value
		case TWDR:
			m.TWIBuffer = value
			ioRegs[address] = value
		case TWCR:
			oldTWCR := ioRegs[TWCR]
			ioRegs[TWCR] = value
			if value&(1<<7) != 0 && oldTWCR&(1<<7) != 0 {
				ioRegs[TWCR] &= ^uint8(1 << 7)
			}
			if value&(1<<2) != 0 {
				if value&(1<<5) != 0 {
					m.TWIState = 0x08
					ioRegs[TWSR] = (ioRegs[TWSR] & 0x07) | 0x08
					ioRegs[TWCR] |= 1 << 7
					ioRegs[TWCR] &= ^uint8(1 << 5)
				} else if value&(1<<4) != 0 {
					m.TWIState = 0xF8
					ioRegs[TWSR] = (ioRegs[TWSR] & 0x07) | 0xF8
					ioRegs[TWCR] &= ^uint8(1 << 4)
				} else if (value&(1<<7)) != 0 && (oldTWCR&(1<<7)) != 0 {
					m.updateTWIState(ioRegs)
				}
			}
			if (ioRegs[TWCR]&(1<<7)) != 0 && (ioRegs[TWCR]&(1<<0)) != 0 {
				m.Sys.TriggerInterrupt(24)
			}
		case USBSTA:
			// VBUS (bit 0) is read-only
			value = (value & ^uint8(0x01)) | (ioRegs[USBSTA] & 0x01)
			ioRegs[address] = value
		case UENUM:
			m.USBSelectedEP = value & 0x07
			if m.USBSelectedEP < 7 {
				ep := &m.USBEndpoints[m.USBSelectedEP]
				ioRegs[UECFG0X] = ep.Config0
				ioRegs[UECFG1X] = ep.Config1
				ioRegs[UESTA0X] = ep.Status0
				ioRegs[UESTA1X] = ep.Status1
				ioRegs[UECONX] = ep.Control
				ioRegs[UEINTX] = ep.Interrupt
				ioRegs[UEBCLX] = uint8(len(ep.FIFO) & 0xFF)
				ioRegs[UEBCHX] = uint8(len(ep.FIFO) >> 8)
			}
			ioRegs[address] = value
		case UEDATX:
			if m.USBSelectedEP < 7 {
				ep := &m.USBEndpoints[m.USBSelectedEP]
				ep.FIFO = append(ep.FIFO, value)
				ioRegs[UEBCLX] = uint8(len(ep.FIFO) & 0xFF)
				ioRegs[UEBCHX] = uint8(len(ep.FIFO) >> 8)
			}
			ioRegs[address] = value
		case UEINTX:
			if m.USBSelectedEP < 7 {
				ep := &m.USBEndpoints[m.USBSelectedEP]
				// Check if TXINI is being cleared (bit 0)
				// Clearing TXINI in hardware means the FIFO is ready to be sent.
				if (ep.Interrupt&0x01) != 0 && (value&0x01) == 0 {
					// Capture the report if it's on a HID endpoint (assume EP1 for now, or any non-EP0)
					if m.USBSelectedEP != 0 && len(ep.FIFO) > 0 {
						report := make([]byte, len(ep.FIFO))
						copy(report, ep.FIFO)
						m.CapturedHIDReports = append(m.CapturedHIDReports, report)
						ep.FIFO = nil // Clear FIFO after "sending"
					}
				}
				// Clear flags by writing 0 as per datasheet
				ep.Interrupt &= value
				ioRegs[UEINTX] = ep.Interrupt
			}
		case UECFG0X:
			if m.USBSelectedEP < 7 {
				m.USBEndpoints[m.USBSelectedEP].Config0 = value
			}
			ioRegs[address] = value
		case UECFG1X:
			if m.USBSelectedEP < 7 {
				m.USBEndpoints[m.USBSelectedEP].Config1 = value
			}
			ioRegs[address] = value
		case UECONX:
			if m.USBSelectedEP < 7 {
				m.USBEndpoints[m.USBSelectedEP].Control = value
			}
			ioRegs[address] = value
		case UDADDR:
			m.USBAddress = value & 0x7F
			ioRegs[address] = value
		case TC4H:
			m.Timer4HighByte = value & 0x03
		case TCNT4:
			m.Timer4Counter = (uint16(m.Timer4HighByte) << 8) | uint16(value)
		case OCR4A:
			m.Timer4OCR4A = (uint16(m.Timer4HighByte) << 8) | uint16(value)
		case OCR4B:
			m.Timer4OCR4B = (uint16(m.Timer4HighByte) << 8) | uint16(value)
		case OCR4C:
			m.Timer4OCR4C = (uint16(m.Timer4HighByte) << 8) | uint16(value)
		case OCR4D:
			m.Timer4OCR4D = (uint16(m.Timer4HighByte) << 8) | uint16(value)
		case TCCR4A:
			m.Timer4ControlA = value
		case TCCR4B:
			m.Timer4ControlB = value
		case TCCR4C:
			m.Timer4ControlC = value
		case TCCR4D:
			m.Timer4ControlD = value
		case TCCR4E:
			m.Timer4ControlE = value
		case DT4:
			m.Timer4DT4 = value
		case PLLCSR:
			m.PLLControl = value
			// Simulate PLOCK bit being set immediately after PLLE is set
			if (value & (1 << PLLE)) != 0 {
				ioRegs[PLLCSR] |= 1 << PLOCK
			} else {
				ioRegs[PLLCSR] &= ^uint8(1 << PLOCK)
			}
		case SMCR:
			m.SleepEnabled = (value & 0x01) != 0
			ioRegs[address] = value
		case SPMCSR:
			// Bit 6 (RWWSB) is read-only
			value = (value & ^uint8(1<<RWWSB)) | (ioRegs[SPMCSR] & (1 << RWWSB))
			// SPMCSR handling
			// Bits: SPMIE, RWWSB, SIGRD, RWWSRE, BLBSET, PGWRT, PGERS, SPMEN
			// SPMEN is automatically cleared after 4 cycles if not used by SPM instruction.
			if (value & (1 << SPMEN)) != 0 {
				m.SPMTimeout = 4
			}
			if (value & (1 << PGWRT)) != 0 {
				m.SPMTimeout = 4 // PGWRT is also cleared after 4 cycles or SPM execution
			}
			ioRegs[address] = value
		case WDTCSR:
			oldVal := ioRegs[WDTCSR]
			// The WDCE and WDE bits must be set to 1 to enable changes to the other bits.
			// Setting these bits initiates a 4-cycle window for further changes.
			if (value&(1<<WDCE)) != 0 && (value&(1<<WDE)) != 0 {
				m.WatchdogTimedChange = 4
				ioRegs[address] = value
				return ioRegs[address]
			}

			// If we are in the timed change window or the WDE bit is being cleared
			if m.WatchdogTimedChange > 0 || (oldVal&(1<<WDE) != 0 && value&(1<<WDE) == 0) {
				// Update WDP bits and WDE/WDIE
				newVal := value & ^uint8(1<<WDCE) // WDCE is always cleared
				ioRegs[WDTCSR] = newVal
				m.WatchdogTimedChange = 0
				m.WatchdogCycles = 0 // Reset counter on any update to WDP
				m.updateWatchdogTimeout(newVal)
			}
		default:
			ioRegs[address] = value
		}
		return ioRegs[address]
	} else {
		switch address {
		case TCNT0:
			return m.Timer0Counter
		case TCCR0A:
			return m.Timer0ControlA
		case TCCR0B:
			return m.Timer0ControlB
		case OCR0A:
			return m.Timer0CompareA
		case OCR0B:
			return m.Timer0CompareB
		case TCNT1L:
			return uint8(m.Timer1Counter & 0xFF)
		case TCNT1H:
			return uint8(m.Timer1Counter >> 8)
		case OCR1AL:
			return uint8(m.Timer1CompareA & 0xFF)
		case OCR1AH:
			return uint8(m.Timer1CompareA >> 8)
		case TCCR1A:
			return m.Timer1ControlA
		case TCCR1B:
			return m.Timer1ControlB
		case TCNT3L:
			return uint8(m.Timer3Counter & 0xFF)
		case TCNT3H:
			return uint8(m.Timer3Counter >> 8)
		case OCR3AL:
			return uint8(m.Timer3CompareA & 0xFF)
		case OCR3AH:
			return uint8(m.Timer3CompareA >> 8)
		case TCCR3A:
			return m.Timer3ControlA
		case TCCR3B:
			return m.Timer3ControlB
		case UDR1:
			if len(m.UART1RXBuffer) > 0 {
				val := m.UART1RXBuffer[0]
				m.UART1RXBuffer = m.UART1RXBuffer[1:]
				if len(m.UART1RXBuffer) == 0 {
					ioRegs[UCSR1A] &= ^uint8(1 << 7)
				}
				return val
			}
			return 0
		case SPDR:
			ioRegs[SPSR] &= ^uint8(1 << 7)
			return m.SPIBuffer
		case ADCL:
			return uint8(m.ADCValue & 0xFF)
		case ADCH:
			return uint8(m.ADCValue >> 8)
		case UEDATX:
			if m.USBSelectedEP < 7 {
				ep := &m.USBEndpoints[m.USBSelectedEP]
				if len(ep.FIFO) > 0 {
					val := ep.FIFO[0]
 				ep.FIFO = ep.FIFO[1:]
					ioRegs[UEBCLX] = uint8(len(ep.FIFO) & 0xFF)
					ioRegs[UEBCHX] = uint8(len(ep.FIFO) >> 8)
					return val
				}
			}
			return 0
		case UEBCLX:
			if m.USBSelectedEP < 7 {
				return uint8(len(m.USBEndpoints[m.USBSelectedEP].FIFO) & 0xFF)
			}
			return 0
		case UEBCHX:
			if m.USBSelectedEP < 7 {
				return uint8(len(m.USBEndpoints[m.USBSelectedEP].FIFO) >> 8)
			}
			return 0
		case TC4H:
			return m.Timer4HighByte
		case TCNT4:
			m.Timer4HighByte = uint8(m.Timer4Counter >> 8)
			return uint8(m.Timer4Counter & 0xFF)
		case OCR4A:
			m.Timer4HighByte = uint8(m.Timer4OCR4A >> 8)
			return uint8(m.Timer4OCR4A & 0xFF)
		case OCR4B:
			m.Timer4HighByte = uint8(m.Timer4OCR4B >> 8)
			return uint8(m.Timer4OCR4B & 0xFF)
		case OCR4C:
			m.Timer4HighByte = uint8(m.Timer4OCR4C >> 8)
			return uint8(m.Timer4OCR4C & 0xFF)
		case OCR4D:
			m.Timer4HighByte = uint8(m.Timer4OCR4D >> 8)
			return uint8(m.Timer4OCR4D & 0xFF)
		case TCCR4A:
			return m.Timer4ControlA
		case TCCR4B:
			return m.Timer4ControlB
		case TCCR4C:
			return m.Timer4ControlC
		case TCCR4D:
			return m.Timer4ControlD
		case TCCR4E:
			return m.Timer4ControlE
		case DT4:
			return m.Timer4DT4
		case PLLCSR:
			return m.PLLControl | (ioRegs[PLLCSR] & (1 << PLOCK))
		default:
			return ioRegs[address]
		}
	}
}

func (m *Manager) Tick(cycles uint64) {
	ioRegs := m.Sys.IORegs()
	prr0 := ioRegs[PRR0]
	prr1 := ioRegs[PRR1]

	if (prr0 & (1 << PRTIM0)) == 0 {
		m.updateTimer0(cycles)
	}
	if (prr0 & (1 << PRTIM1)) == 0 {
		m.updateTimer1(cycles)
	}
	if (prr1 & (1 << PRTIM3)) == 0 {
		m.updateTimer3(cycles)
	}
	if (prr1 & (1 << PRTIM4)) == 0 {
		m.updateTimer4(cycles)
	}
	m.updateWatchdog(cycles)
	if (prr1 & (1 << PRUSB)) == 0 {
		m.updateUSB()
	}
	m.updateSPM(cycles)
	m.updateEEPROM(cycles)
	m.updateMacro(cycles)
}

func (m *Manager) HandlePinChange(port int8, mask uint8, value uint8) {
	ioRegs := m.Sys.IORegs()
	
	// External Interrupts INT0-INT3 are on Port D
	if port == 'D' {
		// INT0: PD0, INT1: PD1, INT2: PD2, INT3: PD3
		for i := uint8(0); i <= 3; i++ {
			if (mask & (1 << i)) != 0 {
				m.checkExternalInterrupt(i, (value>>i)&1)
			}
		}
	}
	// INT6 is on Port E bit 6
	if port == 'E' && (mask&(1<<6)) != 0 {
		m.checkExternalInterrupt(6, (value>>6)&1)
	}

	// PCINT0 on Port B
	if port == 'B' {
		if (ioRegs[PCICR] & 0x01) != 0 {
			pcmask0 := ioRegs[PCMSK0]
			if (mask & pcmask0) != 0 {
				ioRegs[PCIFR] |= 0x01
				m.Sys.TriggerInterrupt(9) // PCINT0 vector
			}
		}
	}
}

func (m *Manager) checkExternalInterrupt(intNum uint8, level uint8) {
	ioRegs := m.Sys.IORegs()
	if (ioRegs[EIMSK] & (1 << intNum)) == 0 {
		return
	}

	var mode uint8
	if intNum <= 3 {
		mode = (ioRegs[EICRA] >> (intNum * 2)) & 0x03
	} else if intNum == 6 {
		mode = (ioRegs[EICRB] >> 4) & 0x03 // INT6 is bits 4:5 of EICRB
	} else {
		return
	}

	trigger := false
	switch mode {
	case 0: // Low level
		if level == 0 {
			trigger = true
		}
	case 1: // Any edge
		trigger = true
	case 2: // Falling edge
		if level == 0 {
			trigger = true
		}
	case 3: // Rising edge
		if level == 1 {
			trigger = true
		}
	}

	if trigger {
		ioRegs[EIFR] |= (1 << intNum)
		// Vector numbers: INT0: 2, INT1: 3, INT2: 4, INT3: 5, INT6: 7
		vector := uint8(0)
		switch intNum {
		case 0, 1, 2, 3:
			vector = 2 + intNum
		case 6:
			vector = 7
		}
		if vector != 0 {
			m.Sys.TriggerInterrupt(vector)
		}
	}
}

func (m *Manager) updateMacro(cycles uint64) {
	if !m.MacroActive || len(m.MacroQueue) == 0 {
		return
	}

	if m.MacroDelayCounter > 0 {
		if cycles >= m.MacroDelayCounter {
			m.MacroDelayCounter = 0
		} else {
			m.MacroDelayCounter -= cycles
			return
		}
	}

	// Move current record to HIDKeyMap and trigger USB
	record := m.MacroQueue[0]
	m.HIDKeyMap = record.KeyMap
	m.MacroDelayCounter = uint64(record.Delay)
	m.MacroQueue = m.MacroQueue[1:]

	if len(m.MacroQueue) == 0 {
		m.MacroActive = false
	}
}

func (m *Manager) updateEEPROM(cycles uint64) {
	if m.EEPROMWriteTimer > 0 {
		if cycles >= m.EEPROMWriteTimer {
			m.EEPROMWriteTimer = 0
		} else {
			m.EEPROMWriteTimer -= cycles
		}

		if m.EEPROMWriteTimer == 0 {
			ioRegs := m.Sys.IORegs()
			ioRegs[EECR] &= ^uint8(0x02) // Clear EEPE (EEPROM Program Enable)
		}
	}
}

func (m *Manager) updateSPM(cycles uint64) {
	if m.SPMTimeout > 0 {
		if cycles >= uint64(m.SPMTimeout) {
			m.SPMTimeout = 0
		} else {
			m.SPMTimeout -= uint8(cycles)
		}

		if m.SPMTimeout == 0 {
			ioRegs := m.Sys.IORegs()
			ioRegs[SPMCSR] &= ^uint8((1 << SPMEN) | (1 << PGWRT) | (1 << PGERS) | (1 << RWWSRE) | (1 << BLBSET))
		}
	}
}

func (m *Manager) updateUSB() {
	ioRegs := m.Sys.IORegs()

	// Handle USB Reset
	if (ioRegs[USBCON] & (1 << 0)) != 0 { // USBE bit
		// Emulate a VBUS connection if USBE is set
		ioRegs[USBSTA] |= (1 << 0) // VBUS
	}

	// Simple HID injection (emulate keys being pressed)
	if m.USBConfigured {
		// Example: In a real implementation, we would check for HID events here.
		// For now, we just ensure that if HIDKeyMap is updated externally,
		// the data can be moved into endpoint FIFOs.
		hidEP := &m.USBEndpoints[1] // Assume EP1 is HID IN
		if len(hidEP.FIFO) == 0 && m.HIDKeyMap[0] != 0 {
			hidEP.FIFO = append(hidEP.FIFO, m.HIDKeyMap[:]...)
			hidEP.Interrupt |= (1 << 0) // TXINI
		}
	}

	// Handle endpoint 0 (Control) setup
	ep0 := &m.USBEndpoints[0]
	if len(ep0.SetupFIFO) >= 8 {
		ep0.Interrupt |= (1 << 3) // RXSTPI
		if m.USBSelectedEP == 0 {
			ioRegs[UEINTX] = ep0.Interrupt
		}
		// Trigger USB General interrupt
		m.Sys.TriggerInterrupt(10)
		m.handleEP0Setup()
	}

	// Basic endpoint interrupt handling
	for i := 0; i < 7; i++ {
		if m.USBEndpoints[i].Interrupt != 0 {
			ioRegs[UEINT] |= (1 << uint(i))
		} else {
			ioRegs[UEINT] &= ^(1 << uint(i))
		}
	}

	if ioRegs[UEINT] != 0 {
		m.Sys.TriggerInterrupt(11) // USB Endpoint interrupt
	}
}

func (m *Manager) handleEP0Setup() {
	ep0 := &m.USBEndpoints[0]
	if len(ep0.SetupFIFO) < 8 {
		return
	}
	bmRequestType := ep0.SetupFIFO[0]
	bRequest := ep0.SetupFIFO[1]
	wValue := uint16(ep0.SetupFIFO[2]) | (uint16(ep0.SetupFIFO[3]) << 8)
	// wIndex := uint16(ep0.SetupFIFO[4]) | (uint16(ep0.SetupFIFO[5]) << 8)
	wLength := uint16(ep0.SetupFIFO[6]) | (uint16(ep0.SetupFIFO[7]) << 8)

	// Standard Requests
	switch bmRequestType & 0x60 {
	case 0:
		switch bRequest {
		case 0x06: // GET_DESCRIPTOR
			descType := ep0.SetupFIFO[3]
			switch descType {
			case 0x01: // Device Descriptor
				m.sendDescriptor(0, []byte{
					18, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00, 0x40,
					0xEB, 0x03, 0x24, 0x20, 0x00, 0x01, 0x01, 0x02, 0x03, 0x01,
				}, wLength)
			case 0x02: // Configuration Descriptor
				m.sendDescriptor(0, []byte{
					9, 0x02, 34, 0, 1, 1, 0, 0xA0, 50,
					9, 0x04, 0, 0, 1, 0x03, 0, 0, 0,
					9, 0x21, 0x11, 0x01, 0, 1, 0x22, 47, 0,
					7, 0x05, 0x81, 0x03, 8, 0, 10,
				}, wLength)
			case 0x22: // HID Report Descriptor
				m.sendDescriptor(0, make([]byte, 47), wLength)
			}
		case 0x05: // SET_ADDRESS
			m.USBAddress = uint8(wValue & 0x7F)
			// Status Stage (In ACK)
			ep0.Interrupt |= (1 << 0) // TXINI
		case 0x09: // SET_CONFIGURATION
			m.USBConfigured = (wValue != 0)
			ep0.Interrupt |= (1 << 0) // TXINI
		}
	case 0x20: // Class-Specific Requests
		// Handle HID-specific requests
		switch bRequest {
		case 0x01: // GET_REPORT
			// Emulate successful response for HID GetReport
			ep0.Interrupt |= (1 << 1) // RXOUTI (ACK from host)
		case 0x09: // SET_REPORT
			// Emulate successful response for HID SetReport
			ep0.Interrupt |= (1 << 0) // TXINI (ACK to host)
		}
	}

	// Clear Setup FIFO after handling
	ep0.SetupFIFO = ep0.SetupFIFO[8:]
}

func (m *Manager) sendDescriptor(epNum uint8, data []byte, maxLen uint16) {
	ep := &m.USBEndpoints[epNum]
	length := uint16(len(data))
	if length > maxLen {
		length = maxLen
	}
	ep.FIFO = append(ep.FIFO, data[:length]...)
	ep.Interrupt |= (1 << 0) // TXINI
}

func (m *Manager) updateTimer0(cycles uint64) {
	prescaler := m.Timer0ControlB & 0x07
	if prescaler == 0 {
		return
	}
	divisor := uint64(1)
	switch prescaler {
	case 2:
		divisor = 8
	case 3:
		divisor = 64
	case 4:
		divisor = 256
	case 5:
		divisor = 1024
	}
	ioRegs := m.Sys.IORegs()
	for i := uint64(0); i < cycles; i++ {
		if m.Sys.Cycles()%divisor == 0 {
			oldVal := m.Timer0Counter
			m.Timer0Counter++
			if m.Timer0Counter < oldVal {
				ioRegs[TIFR0] |= 1 << 0
				if (ioRegs[TIMSK0] & (1 << 0)) != 0 {
					m.Sys.TriggerInterrupt(23)
				}
			}
			if m.Timer0Counter == m.Timer0CompareA {
				ioRegs[TIFR0] |= 1 << 1
				if (ioRegs[TIMSK0] & (1 << 1)) != 0 {
					m.Sys.TriggerInterrupt(21)
				}
			}
			if m.Timer0Counter == m.Timer0CompareB {
				ioRegs[TIFR0] |= 1 << 2
				if (ioRegs[TIMSK0] & (1 << 2)) != 0 {
					m.Sys.TriggerInterrupt(22)
				}
			}
		}
	}
}

func (m *Manager) updateTimer1(cycles uint64) {
	prescaler := m.Timer1ControlB & 0x07
	if prescaler == 0 {
		return
	}
	divisor := uint64(1)
	switch prescaler {
	case 2:
		divisor = 8
	case 3:
		divisor = 64
	case 4:
		divisor = 256
	case 5:
		divisor = 1024
	}
	ioRegs := m.Sys.IORegs()
	for i := uint64(0); i < cycles; i++ {
		if m.Sys.Cycles()%divisor == 0 {
			oldVal := m.Timer1Counter
			m.Timer1Counter++
			if m.Timer1Counter < oldVal {
				ioRegs[TIFR1] |= 1 << 0
				if (ioRegs[TIMSK1] & (1 << 0)) != 0 {
					m.Sys.TriggerInterrupt(20)
				}
			}
			if m.Timer1Counter == m.Timer1CompareA {
				ioRegs[TIFR1] |= 1 << 1
				if (ioRegs[TIMSK1] & (1 << 1)) != 0 {
					m.Sys.TriggerInterrupt(17)
				}
			}
			if m.Timer1Counter == m.Timer1CompareB {
				ioRegs[TIFR1] |= 1 << 2
				if (ioRegs[TIMSK1] & (1 << 2)) != 0 {
					m.Sys.TriggerInterrupt(18)
				}
			}
			if m.Timer1Counter == m.Timer1CompareC {
				ioRegs[TIFR1] |= 1 << 3
				if (ioRegs[TIMSK1] & (1 << 3)) != 0 {
					m.Sys.TriggerInterrupt(19)
				}
			}
		}
	}
}

func (m *Manager) updateTimer3(cycles uint64) {
	prescaler := m.Timer3ControlB & 0x07
	if prescaler == 0 {
		return
	}
	divisor := uint64(1)
	switch prescaler {
	case 2:
		divisor = 8
	case 3:
		divisor = 64
	case 4:
		divisor = 256
	case 5:
		divisor = 1024
	}
	ioRegs := m.Sys.IORegs()
	for i := uint64(0); i < cycles; i++ {
		if m.Sys.Cycles()%divisor == 0 {
			oldVal := m.Timer3Counter
			m.Timer3Counter++
			if m.Timer3Counter < oldVal {
				ioRegs[TIFR3] |= 1 << 0
				if (ioRegs[TIMSK3] & (1 << 0)) != 0 {
					m.Sys.TriggerInterrupt(35)
				}
			}
			if m.Timer3Counter == m.Timer3CompareA {
				ioRegs[TIFR3] |= 1 << 1
				if (ioRegs[TIMSK3] & (1 << 1)) != 0 {
					m.Sys.TriggerInterrupt(32)
				}
			}
			if m.Timer3Counter == m.Timer3CompareB {
				ioRegs[TIFR3] |= 1 << 2
				if (ioRegs[TIMSK3] & (1 << 2)) != 0 {
					m.Sys.TriggerInterrupt(33)
				}
			}
			if m.Timer3Counter == m.Timer3CompareC {
				ioRegs[TIFR3] |= 1 << 3
				if (ioRegs[TIMSK3] & (1 << 3)) != 0 {
					m.Sys.TriggerInterrupt(34)
				}
			}
		}
	}
}

func (m *Manager) updateTimer4(cycles uint64) {
	ioRegs := m.Sys.IORegs()
	prescaler := m.Timer4ControlB & 0x0F
	if prescaler == 0 {
		return
	}

	// Clock source
	usePLL := (m.PLLControl & (1 << PCKE)) != 0
	is96MHz := (m.PLLControl & (1 << 4)) != 0 // PINDIV bit for 32U4 (96MHz if 1)
	
	divisor := uint64(1)
	if prescaler >= 1 && prescaler <= 15 {
		divisor = 1 << (prescaler - 1)
	}

	// Adjust cycles based on clock source relative to system clock (16MHz)
	// PLL is either 64MHz or 96MHz.
	tickIncr := uint64(1)
	if usePLL {
		if is96MHz {
			tickIncr = 6 // 96MHz / 16MHz = 6
		} else {
			tickIncr = 4 // 64MHz / 16MHz = 4
		}
	}

	m.Timer4SubCycles += cycles * tickIncr
	for m.Timer4SubCycles >= divisor {
		m.Timer4SubCycles -= divisor
		m.Timer4Counter++

		// OCR4C acts as TOP in many modes
		top := m.Timer4OCR4C
		if top == 0 {
			top = 0x3FF // Default 10-bit TOP
		}

		if m.Timer4Counter > top {
			m.Timer4Counter = 0
			// Overflow interrupt
			ioRegs[TIFR4] |= 1 << 2
			if (ioRegs[TIMSK4] & (1 << 2)) != 0 {
				m.Sys.TriggerInterrupt(39)
			}
		}

		if m.Timer4Counter == m.Timer4OCR4A {
			ioRegs[TIFR4] |= 1 << 6
			if (ioRegs[TIMSK4] & (1 << 6)) != 0 {
				m.Sys.TriggerInterrupt(38)
			}
		}

		if m.Timer4Counter == m.Timer4OCR4B {
			ioRegs[TIFR4] |= 1 << 5
			if (ioRegs[TIMSK4] & (1 << 5)) != 0 {
				m.Sys.TriggerInterrupt(40)
			}
		}

		if m.Timer4Counter == m.Timer4OCR4D {
			ioRegs[TIFR4] |= 1 << 7
			if (ioRegs[TIMSK4] & (1 << 7)) != 0 {
				m.Sys.TriggerInterrupt(41)
			}
		}
	}
}

func (m *Manager) updateTWIState(ioRegs []uint8) {
	switch m.TWIState {
	case 0x08: // START transmitted
		sla := ioRegs[TWDR]
		addr := sla >> 1
		isRead := sla&0x01 != 0

		// Verify pull-up resistor is present (realistic I2C)
		if m.PullUpResistor > 1000000.0 {
			// Without pull-ups, I2C bus remains low, communication fails
			m.TWIState = 0xF8 // Error
			m.ActiveTWI = nil
		} else {
			m.ActiveTWI = nil
			for _, client := range m.TWIClients {
				if client.Address() == addr {
					if client.OnStart(isRead) {
						m.ActiveTWI = client
						if isRead {
							m.TWIState = 0x40 // SLA+R ACK
							ioRegs[TWDR] = client.OnRead()
						} else {
							m.TWIState = 0x18 // SLA+W ACK
						}
					}
					break
				}
			}

			if m.ActiveTWI == nil {
				if isRead {
					m.TWIState = 0x48 // SLA+R NACK
				} else {
					m.TWIState = 0x20 // SLA+W NACK
				}
			}
		}
	case 0x18: // SLA+W ACK
		if m.ActiveTWI != nil {
			if m.ActiveTWI.OnWrite(ioRegs[TWDR]) {
				m.TWIState = 0x28 // Data ACK
			} else {
				m.TWIState = 0x30 // Data NACK
			}
		} else {
			m.TWIState = 0x20
		}
	case 0x20: // SLA+W NACK
		m.TWIState = 0xF8
	case 0x28: // Data ACK (Write)
		if m.ActiveTWI != nil {
			if m.ActiveTWI.OnWrite(ioRegs[TWDR]) {
				m.TWIState = 0x28 // Data ACK
			} else {
				m.TWIState = 0x30 // Data NACK
			}
		}
	case 0x40: // SLA+R ACK
		if m.ActiveTWI != nil {
			m.TWIState = 0x50 // Data ACK
		} else {
			m.TWIState = 0x48
		}
	case 0x48: // SLA+R NACK
		m.TWIState = 0xF8
	case 0x50: // Data ACK (Read)
		if m.ActiveTWI != nil {
			ioRegs[TWDR] = m.ActiveTWI.OnRead()
		}
		m.TWIState = 0x50
	case 0xF8: // STOP or Error
		if m.ActiveTWI != nil {
			m.ActiveTWI.OnStop()
			m.ActiveTWI = nil
		}
	}
	ioRegs[TWSR] = (ioRegs[TWSR] & 0x07) | m.TWIState
	ioRegs[TWCR] |= 1 << 7
}

func (m *Manager) updateWatchdog(cycles uint64) {
	if m.WatchdogTimedChange > 0 {
		if cycles >= m.WatchdogTimedChange {
			m.WatchdogTimedChange = 0
			m.Sys.IORegs()[WDTCSR] &= ^uint8(1 << WDCE)
		} else {
			m.WatchdogTimedChange -= cycles
		}
	}

	ioRegs := m.Sys.IORegs()
	wdtcsr := ioRegs[WDTCSR]

	// Watchdog is enabled if either WDE or WDIE is set
	if (wdtcsr&(1<<WDE)) != 0 || (wdtcsr&(1<<WDIE)) != 0 {
		m.WatchdogCycles += cycles
		if m.WatchdogCycles >= m.WatchdogTimeout {
			// Timeout occurred!
			if (wdtcsr & (1 << WDIE)) != 0 {
				// Interrupt Mode
				ioRegs[WDTCSR] |= 1 << WDIF
				m.Sys.TriggerInterrupt(4) // WDT vector is 4 on ATmega32u4

				// If WDE is also set, the next timeout will cause a reset.
				// If not, WDIE is cleared by hardware.
				if (wdtcsr & (1 << WDE)) == 0 {
					ioRegs[WDTCSR] &= ^uint8(1 << WDIE)
				}
				// Reset cycles for next period (Interrupt mode might lead to Reset mode)
				m.WatchdogCycles = 0
			} else {
				// System Reset Mode (or both if WDIE was already cleared)
				m.WatchdogReset = true
				m.WatchdogCycles = 0
			}
		}
	} else {
		m.WatchdogCycles = 0
	}
}

func (m *Manager) updateWatchdogTimeout(wdtcsr uint8) {
	// Prescaler bits WDP3, WDP2, WDP1, WDP0
	prescaler := (wdtcsr & 0x07) | ((wdtcsr & (1 << WDP3)) >> 2)

	// Timeout in cycles (assuming 128kHz internal oscillator for WDT on AVR)
	// Typical ATmega32u4 has ~128kHz WDT.
	// 16MHz / 128kHz = 125 cycles of the main clock per WDT tick.
	// But usually we just count WDT ticks.
	// WDT Prescaler:
	// 0: 2K cycles (~16ms)
	// 1: 4K cycles (~32ms)
	// 2: 8K cycles (~64ms)
	// 3: 16K cycles (~0.125s)
	// 4: 32K cycles (~0.25s)
	// 5: 64K cycles (~0.5s)
	// 6: 128K cycles (~1.0s)
	// 7: 256K cycles (~2.0s)
	// 8: 512K cycles (~4.0s)
	// 9: 1024K cycles (~8.0s)

	wdtTicks := uint64(2048) << prescaler
	m.WatchdogTimeout = wdtTicks * (SYS_FREQ / WDT_FREQ)
}
