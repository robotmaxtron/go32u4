package peripherals

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
	SPCR = 0x4C
	SPSR = 0x4D
	SPDR = 0x4E

	// TWI Registers
	TWBR  = 0x70
	TWSR  = 0x71
	TWAR  = 0x72
	TWDR  = 0x73
	TWCR  = 0x74
	TWAMR = 0x75

	// ADC Registers
	ADMUX  = 0x7C
	ADCSRA = 0x7A
	ADCSRB = 0x7B
	ADCL   = 0x78
	ADCH   = 0x79
	DIDR0  = 0x7E
	DIDR2  = 0x7D

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
	PLLCSR = 0x49
	PLLE   = 1
	PLOCK  = 0
	PCKE   = 2

	// USB Registers (Simplified CDC mapping)
	UHWCON  = 0xD7
	USBCON  = 0xD8
	USBSTA  = 0xD9
	USBINT  = 0xDA
	UDCON   = 0xE0
	UDINT   = 0xE1
	UDIEN   = 0xE2
	UDADDR  = 0xE3
	UEINTX  = 0xE8
	UENUM   = 0xE9
	UERST   = 0xEA
	UECONX  = 0xEB
	UECFG0X = 0xEC
	UECFG1X = 0xED
	UESTA0X = 0xEE
	UESTA1X = 0xEF
	UEINT   = 0xF4
	UEDATX  = 0xF1
	UEBCLX  = 0xF2

	// Power Management
	SMCR = 0x53 // Sleep Mode Control Register

	// Watchdog Registers
	WDTCSR = 0x60 // Watchdog Timer Control Register

	// Watchdog Timer Bits
	WDP0  = 0
	WDP1  = 1
	WDP2  = 2
	WDE   = 3
	WDCE  = 4
	WDP3  = 5
	WDIE  = 6
	WDIF  = 7
)

const EEPROMSize = 1024

// System represents the core system that peripherals can interact with.
type System interface {
	IORegs() []uint8
	TriggerInterrupt(vector uint8)
	Cycles() uint64
	SaveEEPROM() error
	PinCallback(port int8, mask uint8, value uint8)
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
	EEPROM [EEPROMSize]uint8

	// USART1
	UART1TXBuffer []byte
	UART1RXBuffer []byte

	// SPI
	SPIBuffer uint8

	// TWI
	TWIBuffer      uint8
	TWIState       uint8
	TWIPendingStop bool

	// ADC
	ADCValue uint16

	// USB CDC Emulation
	USBTXBuffer      []byte
	USBRXBuffer      []byte
	SelectedEndpoint uint8

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
}

func NewManager(sys System) *Manager {
	return &Manager{
		Sys: sys,
	}
}

// IOCallback handles the core I/O behavior for ATmega32u4.
func (m *Manager) IOCallback(address uint16, value uint8, isWrite bool) uint8 {
	ioRegs := m.Sys.IORegs()
	if isWrite {
		switch address {
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
			if value&(1<<6) != 0 {
				ioRegs[ADCSRA] &= ^uint8(1 << 6)
				ioRegs[ADCSRA] |= 1 << 4
				if value&(1<<3) != 0 {
					m.Sys.TriggerInterrupt(29)
				}
			}
		case EECR:
			if value&0x01 != 0 {
				addr := (uint16(ioRegs[EEARH]) << 8) | uint16(ioRegs[EEARL])
				if addr < uint16(len(m.EEPROM)) {
					ioRegs[EEDR] = m.EEPROM[addr]
				}
			}
			if value&0x02 != 0 {
				if ioRegs[EECR]&0x04 != 0 {
					addr := (uint16(ioRegs[EEARH]) << 8) | uint16(ioRegs[EEARL])
					if addr < uint16(len(m.EEPROM)) {
						m.EEPROM[addr] = ioRegs[EEDR]
						_ = m.Sys.SaveEEPROM()
					}
					value &= ^uint8(0x06)
				}
			}
			ioRegs[EECR] = value
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
		case UENUM:
			m.SelectedEndpoint = value & 0x07
			ioRegs[address] = value
		case UEDATX:
			if m.SelectedEndpoint == 3 || m.SelectedEndpoint == 4 {
				m.USBTXBuffer = append(m.USBTXBuffer, value)
				ioRegs[UEINTX] |= 1 << 0
			}
			ioRegs[address] = value
		case UEINTX:
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
				ioRegs[PLLCSR] |= (1 << PLOCK)
			} else {
				ioRegs[PLLCSR] &= ^uint8(1 << PLOCK)
			}
		case SMCR:
			m.SleepEnabled = (value & 0x01) != 0
			ioRegs[address] = value
		case WDTCSR:
			oldVal := ioRegs[WDTCSR]
			// The WDCE and WDE bits must be set to 1 to enable changes to the other bits.
			// Setting these bits initiates a 4-cycle window for further changes.
			if (value&(1<<WDCE)) != 0 && (value&(1<<WDE)) != 0 {
				m.WatchdogTimedChange = 4
				// We don't update WDTCSR with these bits yet, as per datasheet
				// actually the datasheet says "Within the next four clock cycles, write the WDE and Watchdog prescaler
				// bits (WDP) as desired, but with the WDCE bit cleared."
				// "The WDCE bit is always cleared by hardware after four clock cycles."
				ioRegs[WDTCSR] |= (1 << WDCE) | (1 << WDE)
				return value
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
		return value
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
			if (m.SelectedEndpoint == 3 || m.SelectedEndpoint == 4) && len(m.USBRXBuffer) > 0 {
				val := m.USBRXBuffer[0]
				m.USBRXBuffer = m.USBRXBuffer[1:]
				if len(m.USBRXBuffer) == 0 {
					ioRegs[UEINTX] &= ^uint8(1 << 1)
				}
				return val
			}
			return ioRegs[address]
		case UEBCLX:
			if m.SelectedEndpoint == 3 || m.SelectedEndpoint == 4 {
				return uint8(len(m.USBRXBuffer))
			}
			return ioRegs[address]
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
	m.updateTimer0(cycles)
	m.updateTimer1(cycles)
	m.updateTimer3(cycles)
	m.updateTimer4(cycles)
	m.updateWatchdog(cycles)
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
	
	divisor := uint64(1)
	if prescaler >= 1 && prescaler <= 15 {
		divisor = 1 << (prescaler - 1)
	}

	// Adjust divisor for PLL if needed (PLL is 64MHz, System is 16MHz)
	// If PCKE is set, Timer 4 runs at 64MHz.
	// Since Tick() is called with system cycles (16MHz), 
	// we need to process 4 Timer 4 cycles for every 1 system cycle if PCKE is set.
	t4Cycles := cycles
	if usePLL {
		t4Cycles = cycles * 4
	}

	for i := uint64(0); i < t4Cycles; i++ {
		// This is a simplification. For precise timing, we should track sub-cycle progress.
		// But given the current Tick(cycles) architecture, we process them in bulk.
		
		// If divisor is > 1, we skip ticks.
		// For simplicity, we only tick if i % divisor == 0
		if i % divisor != 0 {
			continue
		}
		
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
	case 0x08:
		sla := ioRegs[TWDR]
		if sla&0x01 == 0 {
			m.TWIState = 0x18
		} else {
			m.TWIState = 0x40
		}
	case 0x18, 0x28:
		m.TWIState = 0x28
	case 0x40:
		m.TWIState = 0x50
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
				ioRegs[WDTCSR] |= (1 << WDIF)
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
	// If the CPU is 16MHz, then 16ms is 256,000 cycles.
	// 128 ticks of 125 cycles each = 16000 cycles? No.
	// 2048 * 125 = 256,000. Correct.
	m.WatchdogTimeout = wdtTicks * 125
}
