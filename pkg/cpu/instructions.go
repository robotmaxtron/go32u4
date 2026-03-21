package cpu

// Execute decodes and performs the operation specified by the given opcode.
// It updates the CPU state (registers, memory, SREG, PC) and accounts for instruction cycles.
func (c *CPU) Execute(opcode uint16) {
	// Simple instruction decoding for demonstration
	// and common instructions.

	// NOP (0000 0000 0000 0000)
	if opcode == 0x0000 {
		return
	}

	// RJMP (1100 kkkk kkkk kkkk)
	if opcode&0xF000 == 0xC000 {
		offset := int16(opcode & 0x0FFF)
		if offset&0x0800 != 0 {
			offset |= -4096 // Sign extend
		}
		c.PC = uint16(int32(c.PC) + int32(offset))
		c.Cycles++ // RJMP takes 2 cycles
		c.TickPeripheralsHelper(1)
		return
	}

	// ADD (0000 11rd dddd rrrr)
	if opcode&0xFC00 == 0x0C00 {
		d := (opcode >> 4) & 0x001F
		r := (opcode & 0x000F) | ((opcode >> 5) & 0x0010)
		op1 := c.Reg[d]
		op2 := c.Reg[r]
		res16 := uint16(op1) + uint16(op2)
		res := uint8(res16)
		c.Reg[d] = res
		
		// Update SREG
		c.SetFlag(SREG_H, ((op1&op2)|(op2&^res)|(^res&op1))&0x08 != 0)
		c.SetFlag(SREG_V, ((op1&op2&^res)|(^op1&^op2&res))&0x80 != 0)
		c.SetFlag(SREG_N, (res&0x80) != 0)
		c.SetFlag(SREG_Z, res == 0)
		c.SetFlag(SREG_C, (res16&0x0100) != 0)
		c.SetFlag(SREG_S, c.GetFlag(SREG_N) != c.GetFlag(SREG_V))
		return
	}

	// ADC (0001 11rd dddd rrrr)
	if opcode&0xFC00 == 0x1C00 {
		d := (opcode >> 4) & 0x001F
		r := (opcode & 0x000F) | ((opcode >> 5) & 0x0010)
		op1 := c.Reg[d]
		op2 := c.Reg[r]
		carry := uint8(0)
		if c.GetFlag(SREG_C) {
			carry = 1
		}
		res16 := uint16(op1) + uint16(op2) + uint16(carry)
		res := uint8(res16)
		c.Reg[d] = res
		
		// Update SREG
		c.SetFlag(SREG_H, ((op1&op2)|(op2&^res)|(^res&op1))&0x08 != 0)
		c.SetFlag(SREG_V, ((op1&op2&^res)|(^op1&^op2&res))&0x80 != 0)
		c.SetFlag(SREG_N, (res&0x80) != 0)
		c.SetFlag(SREG_Z, res == 0)
		c.SetFlag(SREG_C, (res16&0x0100) != 0)
		c.SetFlag(SREG_S, c.GetFlag(SREG_N) != c.GetFlag(SREG_V))
		return
	}

	// SUB (0001 10rd dddd rrrr)
	if opcode&0xFC00 == 0x1800 {
		d := (opcode >> 4) & 0x1F
		r := (opcode & 0x0F) | ((opcode >> 5) & 0x10)
		op1 := c.Reg[d]
		op2 := c.Reg[r]
		res := op1 - op2
		
		// Update SREG
		c.SetFlag(SREG_H, ((^op1&op2)|(op2&res)|(res&^op1))&0x08 != 0)
		c.SetFlag(SREG_V, ((op1&^op2&^res)|(^op1&op2&res))&0x80 != 0)
		c.SetFlag(SREG_N, (res&0x80) != 0)
		c.SetFlag(SREG_Z, res == 0)
		c.SetFlag(SREG_C, ((^op1&op2)|(op2&res)|(res&^op1))&0x80 != 0)
		c.SetFlag(SREG_S, c.GetFlag(SREG_N) != c.GetFlag(SREG_V))
		
		c.Reg[d] = res
		return
	}

	// SUBI (0101 kkkk dddd kkkk)
	if opcode&0xF000 == 0x5000 {
		k := uint8(opcode&0x000F) | uint8((opcode>>4)&0x00F0)
		d := 16 + uint8((opcode>>4)&0x000F)
		op1 := c.Reg[d]
		op2 := k
		res := op1 - op2
		
		// Update SREG
		c.SetFlag(SREG_H, ((^op1&op2)|(op2&res)|(res&^op1))&0x08 != 0)
		c.SetFlag(SREG_V, ((op1&^op2&^res)|(^op1&op2&res))&0x80 != 0)
		c.SetFlag(SREG_N, (res&0x80) != 0)
		c.SetFlag(SREG_Z, res == 0)
		c.SetFlag(SREG_C, ((^op1&op2)|(op2&res)|(res&^op1))&0x80 != 0)
		c.SetFlag(SREG_S, c.GetFlag(SREG_N) != c.GetFlag(SREG_V))
		
		c.Reg[d] = res
		return
	}

	// SBC (0000 10rd dddd rrrr)
	if opcode&0xFC00 == 0x0800 {
		d := (opcode >> 4) & 0x1F
		r := (opcode & 0x0F) | ((opcode >> 5) & 0x10)
		op1 := c.Reg[d]
		op2 := c.Reg[r]
		carry := uint8(0)
		if c.GetFlag(SREG_C) {
			carry = 1
		}
		res := op1 - op2 - carry
		
		// Update SREG
		c.SetFlag(SREG_H, ((^op1&op2)|(op2&res)|(res&^op1))&0x08 != 0)
		c.SetFlag(SREG_V, ((op1&^op2&^res)|(^op1&op2&res))&0x80 != 0)
		c.SetFlag(SREG_N, (res&0x80) != 0)
		c.SetFlag(SREG_Z, res == 0 && c.GetFlag(SREG_Z))
		c.SetFlag(SREG_C, ((^op1&op2)|(op2&res)|(res&^op1))&0x80 != 0)
		c.SetFlag(SREG_S, c.GetFlag(SREG_N) != c.GetFlag(SREG_V))
		
		c.Reg[d] = res
		return
	}

	// AND (0010 00rd dddd rrrr)
	if opcode&0xFC00 == 0x2000 {
		d := (opcode >> 4) & 0x1F
		r := (opcode & 0x0F) | ((opcode >> 5) & 0x10)
		res := c.Reg[d] & c.Reg[r]
		c.Reg[d] = res
		c.SetFlag(SREG_V, false)
		c.SetFlag(SREG_N, (res&0x80) != 0)
		c.SetFlag(SREG_Z, res == 0)
		c.SetFlag(SREG_S, c.GetFlag(SREG_N))
		return
	}

	// OR (0010 10rd dddd rrrr)
	if opcode&0xFC00 == 0x2800 {
		d := (opcode >> 4) & 0x1F
		r := (opcode & 0x0F) | ((opcode >> 5) & 0x10)
		res := c.Reg[d] | c.Reg[r]
		c.Reg[d] = res
		c.SetFlag(SREG_V, false)
		c.SetFlag(SREG_N, (res&0x80) != 0)
		c.SetFlag(SREG_Z, res == 0)
		c.SetFlag(SREG_S, c.GetFlag(SREG_N))
		return
	}

	// EOR (0010 01rd dddd rrrr)
	if opcode&0xFC00 == 0x2400 {
		d := (opcode >> 4) & 0x1F
		r := (opcode & 0x0F) | ((opcode >> 5) & 0x10)
		res := c.Reg[d] ^ c.Reg[r]
		c.Reg[d] = res
		c.SetFlag(SREG_V, false)
		c.SetFlag(SREG_N, (res&0x80) != 0)
		c.SetFlag(SREG_Z, res == 0)
		c.SetFlag(SREG_S, c.GetFlag(SREG_N))
		return
	}

	// COM (1001 010d dddd 0000)
	if opcode&0xFE0F == 0x9400 {
		d := (opcode >> 4) & 0x1F
		res := 0xFF - c.Reg[d]
		c.Reg[d] = res
		c.SetFlag(SREG_V, false)
		c.SetFlag(SREG_N, (res&0x80) != 0)
		c.SetFlag(SREG_Z, res == 0)
		c.SetFlag(SREG_C, true)
		c.SetFlag(SREG_S, c.GetFlag(SREG_N))
		return
	}

	// NEG (1001 010d dddd 0001)
	if opcode&0xFE0F == 0x9401 {
		d := (opcode >> 4) & 0x1F
		op := c.Reg[d]
		res := 0 - op
		c.Reg[d] = res
		c.SetFlag(SREG_H, (res|op)&0x08 != 0)
		c.SetFlag(SREG_V, op == 0x80)
		c.SetFlag(SREG_N, (res&0x80) != 0)
		c.SetFlag(SREG_Z, res == 0)
		c.SetFlag(SREG_C, res != 0)
		c.SetFlag(SREG_S, c.GetFlag(SREG_N) != c.GetFlag(SREG_V))
		return
	}

	// MUL (1001 11rd dddd rrrr)
	if opcode&0xFC00 == 0x9C00 {
		d := (opcode >> 4) & 0x1F
		r := (opcode & 0x0F) | ((opcode >> 5) & 0x10)
		res := uint16(c.Reg[d]) * uint16(c.Reg[r])
		c.Reg[0] = uint8(res & 0xFF)
		c.Reg[1] = uint8(res >> 8)
		c.SetFlag(SREG_C, (res&0x8000) != 0)
		c.SetFlag(SREG_Z, res == 0)
		c.Cycles++ // MUL takes 2 cycles
		c.TickPeripheralsHelper(1)
		return
	}

	// MULS (0000 0010 dddd rrrr) - only R16-R31
	if opcode&0xFF00 == 0x0200 {
		d := 16 + ((opcode >> 4) & 0x0F)
		r := 16 + (opcode & 0x0F)
		res := int16(int8(c.Reg[d])) * int16(int8(c.Reg[r]))
		c.Reg[0] = uint8(uint16(res) & 0xFF)
		c.Reg[1] = uint8(uint16(res) >> 8)
		c.SetFlag(SREG_C, (uint16(res)&0x8000) != 0)
		c.SetFlag(SREG_Z, res == 0)
		c.Cycles++ // MULS takes 2 cycles
		c.TickPeripheralsHelper(1)
		return
	}

	// LDI (1110 kkkk dddd kkkk)
	if opcode&0xF000 == 0xE000 {
		k := uint8(opcode&0x000F) | uint8((opcode>>4)&0x00F0)
		d := 16 + uint8((opcode>>4)&0x000F)
		c.Reg[d] = k
		return
	}

	// OUT (1011 1AAr rrrr AAAA)
	if opcode&0xF800 == 0xB800 {
		r := (opcode >> 4) & 0x1F
		a := (opcode & 0x0F) | ((opcode >> 5) & 0x30)
		c.WriteIO(a, c.Reg[r])
		return
	}

	// IN (1011 0AAr rrrr AAAA)
	if opcode&0xF800 == 0xB000 {
		r := (opcode >> 4) & 0x1F
		a := (opcode & 0x0F) | ((opcode >> 5) & 0x30)
		c.Reg[r] = c.ReadIO(a)
		return
	}

	// INC (1001 010d dddd 0011)
	if opcode&0xFE0F == 0x9403 {
		d := (opcode >> 4) & 0x1F
		op := c.Reg[d]
		res := op + 1
		c.Reg[d] = res
		c.SetFlag(SREG_V, op == 0x7F)
		c.SetFlag(SREG_N, (res&0x80) != 0)
		c.SetFlag(SREG_Z, res == 0)
		c.SetFlag(SREG_S, c.GetFlag(SREG_N) != c.GetFlag(SREG_V))
		return
	}

	// DEC (1001 010d dddd 1010)
	if opcode&0xFE0F == 0x940A {
		d := (opcode >> 4) & 0x1F
		op := c.Reg[d]
		res := op - 1
		c.Reg[d] = res
		c.SetFlag(SREG_V, op == 0x80)
		c.SetFlag(SREG_N, (res&0x80) != 0)
		c.SetFlag(SREG_Z, res == 0)
		c.SetFlag(SREG_S, c.GetFlag(SREG_N) != c.GetFlag(SREG_V))
		return
	}

	// PUSH (1001 001r rrrr 1111)
	if opcode&0xFE0F == 0x920F {
		r := (opcode >> 4) & 0x1F
		c.Push(c.Reg[r])
		c.Cycles++ // PUSH takes 2 cycles
		c.TickPeripheralsHelper(1)
		return
	}

	// POP (1001 000d dddd 1111)
	if opcode&0xFE0F == 0x900F {
		d := (opcode >> 4) & 0x1F
		c.Reg[d] = c.Pop()
		c.Cycles++ // POP takes 2 cycles
		c.TickPeripheralsHelper(1)
		return
	}

	// BSET (1001 0100 0sss 1000) - SEC, SEZ, etc.
	if opcode&0xFF8F == 0x9408 {
		s := uint8((opcode >> 4) & 0x07)
		c.SREG |= (1 << s)
		return
	}

	// BCLR (1001 0100 1sss 1000) - CLC, CLZ, etc.
	if opcode&0xFF8F == 0x9488 {
		s := uint8((opcode >> 4) & 0x07)
		c.SREG &= ^(1 << s)
		return
	}

	// SBI (1001 1010 AAAA Abbb)
	if opcode&0xFF00 == 0x9A00 {
		a := (opcode >> 3) & 0x1F
		b := opcode & 0x07
		val := c.ReadIO(a)
		c.WriteIO(a, val|(1<<b))
		c.Cycles++ // SBI takes 2 cycles
		c.TickPeripheralsHelper(1)
		return
	}

	// CBI (1001 1000 AAAA Abbb)
	if opcode&0xFF00 == 0x9800 {
		a := (opcode >> 3) & 0x1F
		b := opcode & 0x07
		val := c.ReadIO(a)
		c.WriteIO(a, val&^(1<<b))
		c.Cycles++ // CBI takes 2 cycles
		c.TickPeripheralsHelper(1)
		return
	}

	// LSL (is ADD Rd, Rd)
	// LSR (1001 010d dddd 0110)
	if opcode&0xFE0F == 0x9406 {
		d := (opcode >> 4) & 0x1F
		val := c.Reg[d]
		c.SetFlag(SREG_C, (val&0x01) != 0)
		res := val >> 1
		c.Reg[d] = res
		c.SetFlag(SREG_Z, res == 0)
		c.SetFlag(SREG_N, false)
		c.SetFlag(SREG_V, c.GetFlag(SREG_N) != c.GetFlag(SREG_C))
		c.SetFlag(SREG_S, c.GetFlag(SREG_N) != c.GetFlag(SREG_V))
		return
	}

	// ASR (1001 010d dddd 0101)
	if opcode&0xFE0F == 0x9405 {
		d := (opcode >> 4) & 0x1F
		val := c.Reg[d]
		c.SetFlag(SREG_C, (val&0x01) != 0)
		res := uint8(int8(val) >> 1)
		c.Reg[d] = res
		c.SetFlag(SREG_Z, res == 0)
		c.SetFlag(SREG_N, (res&0x80) != 0)
		c.SetFlag(SREG_V, c.GetFlag(SREG_N) != c.GetFlag(SREG_C))
		c.SetFlag(SREG_S, c.GetFlag(SREG_N) != c.GetFlag(SREG_V))
		return
	}

	// ROR (1001 010d dddd 0111)
	if opcode&0xFE0F == 0x9407 {
		d := (opcode >> 4) & 0x1F
		val := c.Reg[d]
		oldC := uint8(0)
		if c.GetFlag(SREG_C) {
			oldC = 0x80
		}
		c.SetFlag(SREG_C, (val&0x01) != 0)
		res := (val >> 1) | oldC
		c.Reg[d] = res
		c.SetFlag(SREG_Z, res == 0)
		c.SetFlag(SREG_N, (res&0x80) != 0)
		c.SetFlag(SREG_V, c.GetFlag(SREG_N) != c.GetFlag(SREG_C))
		c.SetFlag(SREG_S, c.GetFlag(SREG_N) != c.GetFlag(SREG_V))
		return
	}

	// SWAP (1001 010d dddd 0010)
	if opcode&0xFE0F == 0x9402 {
		d := (opcode >> 4) & 0x1F
		val := c.Reg[d]
		res := (val << 4) | (val >> 4)
		c.Reg[d] = res
		return
	}

	// CP (0001 01rd dddd rrrr)
	if opcode&0xFC00 == 0x1400 {
		d := (opcode >> 4) & 0x1F
		r := (opcode & 0x0F) | ((opcode >> 5) & 0x10)
		op1 := c.Reg[d]
		op2 := c.Reg[r]
		res := op1 - op2
		c.SetFlag(SREG_H, ((^op1&op2)|(op2&res)|(res&^op1))&0x08 != 0)
		c.SetFlag(SREG_V, ((op1&^op2&^res)|(^op1&op2&res))&0x80 != 0)
		c.SetFlag(SREG_N, (res&0x80) != 0)
		c.SetFlag(SREG_Z, res == 0)
		c.SetFlag(SREG_C, ((^op1&op2)|(op2&res)|(res&^op1))&0x80 != 0)
		c.SetFlag(SREG_S, c.GetFlag(SREG_N) != c.GetFlag(SREG_V))
		return
	}

	// CPC (0000 01rd dddd rrrr)
	if opcode&0xFC00 == 0x0400 {
		d := (opcode >> 4) & 0x1F
		r := (opcode & 0x0F) | ((opcode >> 5) & 0x10)
		op1 := c.Reg[d]
		op2 := c.Reg[r]
		carry := uint8(0)
		if c.GetFlag(SREG_C) {
			carry = 1
		}
		res := op1 - op2 - carry
		c.SetFlag(SREG_H, ((^op1&op2)|(op2&res)|(res&^op1))&0x08 != 0)
		c.SetFlag(SREG_V, ((op1&^op2&^res)|(^op1&op2&res))&0x80 != 0)
		c.SetFlag(SREG_N, (res&0x80) != 0)
		c.SetFlag(SREG_Z, res == 0 && c.GetFlag(SREG_Z))
		c.SetFlag(SREG_C, ((^op1&op2)|(op2&res)|(res&^op1))&0x80 != 0)
		c.SetFlag(SREG_S, c.GetFlag(SREG_N) != c.GetFlag(SREG_V))
		return
	}

	// CPI (0011 kkkk dddd kkkk)
	if opcode&0xF000 == 0x3000 {
		k := uint8(opcode&0x000F) | uint8((opcode>>4)&0x00F0)
		d := 16 + uint8((opcode>>4)&0x000F)
		op1 := c.Reg[d]
		op2 := k
		res := op1 - op2
		c.SetFlag(SREG_H, ((^op1&op2)|(op2&res)|(res&^op1))&0x08 != 0)
		c.SetFlag(SREG_V, ((op1&^op2&^res)|(^op1&op2&res))&0x80 != 0)
		c.SetFlag(SREG_N, (res&0x80) != 0)
		c.SetFlag(SREG_Z, res == 0)
		c.SetFlag(SREG_C, ((^op1&op2)|(op2&res)|(res&^op1))&0x80 != 0)
		c.SetFlag(SREG_S, c.GetFlag(SREG_N) != c.GetFlag(SREG_V))
		return
	}

	// BRBS (1111 00kk kkkk ksss)
	if opcode&0xFC00 == 0xF000 {
		s := uint8(opcode & 0x07)
		k := int8((opcode >> 3) & 0x7F)
		if k&0x40 != 0 {
			k |= -128 // sign extend 7-bit to 8-bit
		}
		if c.GetFlag(s) {
			c.PC = uint16(int32(c.PC) + int32(k))
			c.Cycles++ // Branch taken: 2 cycles total
			c.TickPeripherals(1)
		}
		return
	}

	// BRBC (1111 01kk kkkk ksss)
	if opcode&0xFC00 == 0xF400 {
		s := uint8(opcode & 0x07)
		k := int8((opcode >> 3) & 0x7F)
		if k&0x40 != 0 {
			k |= -128
		}
		if !c.GetFlag(s) {
			c.PC = uint16(int32(c.PC) + int32(k))
			c.Cycles++ // Branch taken: 2 cycles total
			c.TickPeripherals(1)
		}
		return
	}

	flash := c.Bus.Flash()
	// CALL (1001 010k kkkk 111k) (kkkk kkkk kkkk kkkk) - 32-bit
	if opcode&0xFE0E == 0x940E {
		k1 := (opcode >> 4) & 0x01F
		k2 := opcode & 0x01
		nextOp := flash[c.PC]
		c.PC++
		target := (uint32(k1) << 17) | (uint32(k2) << 16) | uint32(nextOp)
		// Push PC
		c.Push(uint8(c.PC & 0xFF))
		c.Push(uint8((c.PC >> 8) & 0xFF))
		c.PC = uint16(target)
		c.Cycles += 3 // CALL takes 4 cycles
		c.TickPeripheralsHelper(3)
		return
	}

	// RET (1001 0101 0000 1000)
	if opcode == 0x9508 {
		pch := uint16(c.Pop())
		pcl := uint16(c.Pop())
		c.PC = (pch << 8) | pcl
		c.Cycles += 3 // RET takes 4 cycles
		c.TickPeripheralsHelper(3)
		return
	}

	// RETI (1001 0101 0001 1000)
	if opcode == 0x9518 {
		pch := uint16(c.Pop())
		pcl := uint16(c.Pop())
		c.PC = (pch << 8) | pcl
		c.SetFlag(SREG_I, true)
		c.Cycles += 3 // RETI takes 4 cycles
		c.TickPeripheralsHelper(3)
		return
	}

	// RCALL (1101 kkkk kkkk kkkk)
	if opcode&0xF000 == 0xD000 {
		k := int16(opcode & 0x0FFF)
		if k&0x0800 != 0 {
			k |= -4096
		}
		c.Push(uint8(c.PC & 0xFF))
		c.Push(uint8((c.PC >> 8) & 0xFF))
		c.PC = uint16(int32(c.PC) + int32(k))
		c.Cycles += 2 // RCALL takes 3 cycles
		c.TickPeripheralsHelper(2)
		return
	}

	// JMP (1001 010k kkkk 110k) (kkkk kkkk kkkk kkkk)
	if opcode&0xFE0E == 0x940C {
		k1 := (opcode >> 4) & 0x01F
		k2 := opcode & 0x01
		nextOp := flash[c.PC]
		c.PC++
		target := (uint32(k1) << 17) | (uint32(k2) << 16) | uint32(nextOp)
		c.PC = uint16(target)
		c.Cycles++ // JMP takes 3 cycles
		c.TickPeripheralsHelper(1)
		return
	}

	// IJMP (1001 0100 0000 1001)
	if opcode == 0x9409 {
		// Target is in Z register (R31:R30)
		target := (uint16(c.Reg[31]) << 8) | uint16(c.Reg[30])
		c.PC = target
		c.Cycles++ // IJMP takes 2 cycles
		c.TickPeripheralsHelper(1)
		return
	}

	// ICALL (1001 0101 0000 1001)
	if opcode == 0x9509 {
		target := (uint16(c.Reg[31]) << 8) | uint16(c.Reg[30])
		c.Push(uint8(c.PC & 0xFF))
		c.Push(uint8((c.PC >> 8) & 0xFF))
		c.PC = target
		c.Cycles += 2 // ICALL takes 3 cycles
		c.TickPeripherals(2)
		return
	}

	// LDS (1001 000d dddd 0000) (kkkk kkkk kkkk kkkk)
	if opcode&0xFE0F == 0x9000 {
		d := (opcode >> 4) & 0x1F
		addr := flash[c.PC]
		c.PC++
		c.Reg[d] = c.ReadSRAM(addr)
		c.Cycles++ // LDS takes 2 cycles
		c.TickPeripheralsHelper(1)
		return
	}

	// STS (1001 001r rrrr 0000) (kkkk kkkk kkkk kkkk)
	if opcode&0xFE0F == 0x9200 {
		r := (opcode >> 4) & 0x1F
		addr := flash[c.PC]
		c.PC++
		c.WriteSRAM(addr, c.Reg[r])
		c.Cycles++ // STS takes 2 cycles
		c.TickPeripheralsHelper(1)
		return
	}

	// LD Rd, X (1001 000d dddd 1100)
	if opcode&0xFE0F == 0x900C {
		d := (opcode >> 4) & 0x1F
		x := (uint16(c.Reg[27]) << 8) | uint16(c.Reg[26])
		c.Reg[d] = c.ReadSRAM(x)
		c.Cycles++ // LD takes 2 cycles
		c.TickPeripheralsHelper(1)
		return
	}

	// LD Rd, X+ (1001 000d dddd 1101)
	if opcode&0xFE0F == 0x900D {
		d := (opcode >> 4) & 0x1F
		x := (uint16(c.Reg[27]) << 8) | uint16(c.Reg[26])
		c.Reg[d] = c.ReadSRAM(x)
		x++
		c.Reg[26] = uint8(x & 0xFF)
		c.Reg[27] = uint8(x >> 8)
		c.Cycles++ // LD takes 2 cycles
		c.TickPeripheralsHelper(1)
		return
	}

	// ST X, Rr (1001 001r rrrr 1100)
	if opcode&0xFE0F == 0x920C {
		r := (opcode >> 4) & 0x1F
		x := (uint16(c.Reg[27]) << 8) | uint16(c.Reg[26])
		c.WriteSRAM(x, c.Reg[r])
		c.Cycles++ // ST takes 2 cycles
		c.TickPeripheralsHelper(1)
		return
	}

	// ST X+, Rr (1001 001r rrrr 1101)
	if opcode&0xFE0F == 0x920D {
		r := (opcode >> 4) & 0x1F
		x := (uint16(c.Reg[27]) << 8) | uint16(c.Reg[26])
		c.WriteSRAM(x, c.Reg[r])
		x++
		c.Reg[26] = uint8(x & 0xFF)
		c.Reg[27] = uint8(x >> 8)
		c.Cycles++ // ST takes 2 cycles
		c.TickPeripheralsHelper(1)
		return
	}

	// LD Rd, Y (1000 000d dddd 1000)
	if opcode&0xFE0F == 0x8008 {
		d := (opcode >> 4) & 0x1F
		y := (uint16(c.Reg[29]) << 8) | uint16(c.Reg[28])
		c.Reg[d] = c.ReadSRAM(y)
		c.Cycles++
		c.TickPeripheralsHelper(1)
		return
	}

	// ST Y, Rr (1000 001r rrrr 1000)
	if opcode&0xFE0F == 0x8208 {
		r := (opcode >> 4) & 0x1F
		y := (uint16(c.Reg[29]) << 8) | uint16(c.Reg[28])
		c.WriteSRAM(y, c.Reg[r])
		c.Cycles++
		c.TickPeripheralsHelper(1)
		return
	}

	// LD Rd, Z (1000 000d dddd 0000)
	if opcode&0xFE0F == 0x8000 {
		d := (opcode >> 4) & 0x1F
		z := (uint16(c.Reg[31]) << 8) | uint16(c.Reg[30])
		c.Reg[d] = c.ReadSRAM(z)
		c.Cycles++
		c.TickPeripheralsHelper(1)
		return
	}

	// ST Z, Rr (1000 001r rrrr 0000)
	if opcode&0xFE0F == 0x8200 {
		r := (opcode >> 4) & 0x1F
		z := (uint16(c.Reg[31]) << 8) | uint16(c.Reg[30])
		c.WriteSRAM(z, c.Reg[r])
		c.Cycles++
		c.TickPeripheralsHelper(1)
		return
	}

	// LPM (1001 0101 1100 1000) - Rd is R0, addr is Z
	if opcode == 0x95C8 {
		z := (uint16(c.Reg[31]) << 8) | uint16(c.Reg[30])
		word := flash[z>>1]
		if z&0x01 == 0 {
			c.Reg[0] = uint8(word & 0xFF)
		} else {
			c.Reg[0] = uint8(word >> 8)
		}
		c.Cycles += 2 // LPM takes 3 cycles
		c.TickPeripheralsHelper(2)
		return
	}

	// LPM Rd, Z (1001 000d dddd 0100)
	if opcode&0xFE0F == 0x9004 {
		d := (opcode >> 4) & 0x1F
		z := (uint16(c.Reg[31]) << 8) | uint16(c.Reg[30])
		word := flash[z>>1]
		if z&0x01 == 0 {
			c.Reg[d] = uint8(word & 0xFF)
		} else {
			c.Reg[d] = uint8(word >> 8)
		}
		c.Cycles += 2 // LPM takes 3 cycles
		c.TickPeripheralsHelper(2)
		return
	}
	
	// LPM Rd, Z+ (1001 000d dddd 0101)
	if opcode&0xFE0F == 0x9005 {
		d := (opcode >> 4) & 0x1F
		z := (uint16(c.Reg[31]) << 8) | uint16(c.Reg[30])
		word := flash[z>>1]
		if z&0x01 == 0 {
			c.Reg[d] = uint8(word & 0xFF)
		} else {
			c.Reg[d] = uint8(word >> 8)
		}
		z++
		c.Reg[30] = uint8(z & 0xFF)
		c.Reg[31] = uint8(z >> 8)
		c.Cycles += 2 // LPM takes 3 cycles
		c.TickPeripheralsHelper(2)
		return
	}

	// SLEEP (1001 0101 1000 1000)
	if opcode == 0x9588 {
		// This needs to be handled via the Bus/MCU since Periph is moved
		// For now we'll just skip it or add a Sleep method to CPU
		return
	}
	
	// Default to NOP for unimplemented instructions for now
	// fmt.Printf("Unimplemented opcode: %04X\n", opcode)
}

