package cpu

func (c *CPU) Execute(opcode uint16) int {
	flash := c.Bus.Flash()

	// 1. MUL (1001 11rd dddd rrrr)
	if (opcode & 0xFC00) == 0x9C00 {
		d := (opcode >> 4) & 0x001F
		r := (opcode & 0x000F) | ((opcode >> 5) & 0x0010)
		res := uint16(c.Reg[d]) * uint16(c.Reg[r])
		c.Reg[0], c.Reg[1] = uint8(res), uint8(res>>8)
		c.SetFlag(SREG_Z, res == 0); c.SetFlag(SREG_C, (res&0x8000) != 0)
		return 1
	}

	// 2. CPSE (0001 00rd dddd rrrr)
	if (opcode & 0xFC00) == 0x1000 {
		d := (opcode >> 4) & 0x001F
		r := (opcode & 0x000F) | ((opcode >> 5) & 0x0010)
		if c.Reg[d] == c.Reg[r] {
			if int(c.PC) < len(flash) {
				nextOp := flash[c.PC]
				c.PC++ // Skip next instruction
				// Check for 2-word instructions: LDS, STS, CALL, JMP
				if (nextOp&0xFE0F) == 0x9000 || (nextOp&0xFE0F) == 0x9200 || (nextOp&0xFE0E) == 0x940E || (nextOp&0xFE0E) == 0x940C {
					if int(c.PC) < len(flash) {
						c.PC++
					}
				}
			}
		}
		return 1
	}

	// 3. Skip instructions SBIC / SBIS / SBRC / SBRS
	skip := false
	if (opcode & 0xFF00) == 0x9900 { // SBIC (1001 1001 AAAA Abbb)
		port, bit := (opcode>>3)&0x1F, uint8(opcode&0x07)
		if (c.Bus.ReadIO(uint16(port)) & (1 << bit)) == 0 { skip = true }
	} else if (opcode & 0xFF00) == 0x9B00 { // SBIS (1001 1011 AAAA Abbb)
		port, bit := (opcode>>3)&0x1F, uint8(opcode&0x07)
		if (c.Bus.ReadIO(uint16(port)) & (1 << bit)) != 0 { skip = true }
	} else if (opcode & 0xFE08) == 0xFC00 { // SBRC (1111 110r rrrr 0bbb)
		r, bit := (opcode>>4)&0x1F, uint8(opcode&0x07)
		if (c.Reg[r] & (1 << bit)) == 0 { skip = true }
	} else if (opcode & 0xFE08) == 0xFE00 { // SBRS (1111 111r rrrr 0bbb)
		r, bit := (opcode>>4)&0x1F, uint8(opcode&0x07)
		if (c.Reg[r] & (1 << bit)) != 0 { skip = true }
	}

	if skip {
		if int(c.PC) < len(flash) {
			nextOp := flash[c.PC]
			c.PC++
			// Check for 2-word instructions
			if (nextOp&0xFE0F) == 0x9000 || (nextOp&0xFE0F) == 0x9200 || (nextOp&0xFE0E) == 0x940E || (nextOp&0xFE0E) == 0x940C {
				if int(c.PC) < len(flash) {
					c.PC++
				}
			}
		}
	}

	// 4. I/O Bit
	if (opcode & 0xFE08) == 0x9A00 { // SBI
		port, bit := (opcode>>3)&0x1F, uint8(opcode&0x07)
		val := c.Bus.ReadIO(uint16(port)) | (1 << bit)
		c.Bus.WriteIO(uint16(port), val); return 1
	}
	if (opcode & 0xFE08) == 0x9800 { // CBI
		port, bit := (opcode>>3)&0x1F, uint8(opcode&0x07)
		val := c.Bus.ReadIO(uint16(port)) & ^(1 << bit)
		c.Bus.WriteIO(uint16(port), val); return 1
	}

	// 5. ALU Immediates
	if (opcode&0xF000) == 0xE000 { d,k := 16+uint8((opcode>>4)&0x0F), uint8(opcode&0x0F)|uint8((opcode>>4)&0xF0); c.Reg[d]=k; return 1 }
	if (opcode&0xF000) == 0x7000 { d,k := 16+uint8((opcode>>4)&0x0F), uint8(opcode&0x0F)|uint8((opcode>>4)&0xF0); c.Reg[d]&=k; res := c.Reg[d]; c.SetFlag(SREG_V,false); c.SetFlag(SREG_N,(res&0x80)!=0); c.SetFlag(SREG_Z,res==0); c.SetFlag(SREG_S,c.GetFlag(SREG_N)); return 1 }
	if (opcode&0xF000) == 0x6000 { d,k := 16+uint8((opcode>>4)&0x0F), uint8(opcode&0x0F)|uint8((opcode>>4)&0xF0); c.Reg[d]|=k; res := c.Reg[d]; c.SetFlag(SREG_V,false); c.SetFlag(SREG_N,(res&0x80)!=0); c.SetFlag(SREG_Z,res==0); c.SetFlag(SREG_S,c.GetFlag(SREG_N)); return 1 }
	if (opcode&0xF000) == 0x5000 { d,k := 16+uint8((opcode>>4)&0x0F), uint8(opcode&0x0F)|uint8((opcode>>4)&0xF0); op1 := c.Reg[d]; res := op1-k; c.SetFlag(SREG_H,((^op1&k)|(k&res)|(res&^op1))&0x08!=0); c.SetFlag(SREG_V,((op1&^k&^res)|(^op1&k&res))&0x80!=0); c.SetFlag(SREG_N,(res&0x80)!=0); c.SetFlag(SREG_Z,res==0); c.SetFlag(SREG_C,((^op1&k)|(k&res)|(res&^op1))&0x80!=0); c.SetFlag(SREG_S,c.GetFlag(SREG_N)!=c.GetFlag(SREG_V)); c.Reg[d]=res; return 1 }
	if (opcode&0xF000) == 0x3000 { d,k := 16+uint8((opcode>>4)&0x0F), uint8(opcode&0x0F)|uint8((opcode>>4)&0xF0); op1 := c.Reg[d]; res := op1-k; c.SetFlag(SREG_H,((^op1&k)|(k&res)|(res&^op1))&0x08!=0); c.SetFlag(SREG_V,((op1&^k&^res)|(^op1&k&res))&0x80!=0); c.SetFlag(SREG_N,(res&0x80)!=0); c.SetFlag(SREG_Z,res==0); c.SetFlag(SREG_C,((^op1&k)|(k&res)|(res&^op1))&0x80!=0); c.SetFlag(SREG_S,c.GetFlag(SREG_N)!=c.GetFlag(SREG_V)); return 1 }

	// 6. ALU Reg-Reg
	if (opcode&0xFC00) == 0x0C00 { d,r := (opcode>>4)&0x1F, (opcode&0x0F)|((opcode>>5)&0x10); op1,op2 := c.Reg[d],c.Reg[r]; res16 := uint16(op1)+uint16(op2); res := uint8(res16); c.Reg[d]=res; c.SetFlag(SREG_H,((op1&op2)|(op2&^res)|(^res&op1))&0x08!=0); c.SetFlag(SREG_V,((op1&op2&^res)|(^op1&^op2&res))&0x80!=0); c.SetFlag(SREG_N,(res&0x80)!=0); c.SetFlag(SREG_Z,res==0); c.SetFlag(SREG_C,(res16&0x100)!=0); c.SetFlag(SREG_S,c.GetFlag(SREG_N)!=c.GetFlag(SREG_V)); return 1 }
	if (opcode&0xFC00) == 0x1C00 { d,r := (opcode>>4)&0x1F, (opcode&0x0F)|((opcode>>5)&0x10); op1,op2 := c.Reg[d],c.Reg[r]; carry := uint8(0); if c.GetFlag(SREG_C) { carry = 1 }; res16 := uint16(op1)+uint16(op2)+uint16(carry); res := uint8(res16); c.Reg[d]=res; c.SetFlag(SREG_H,((op1&op2)|(op2&^res)|(^res&op1))&0x08!=0); c.SetFlag(SREG_V,((op1&^op2&^res)|(^op1&op2&res))&0x80!=0); c.SetFlag(SREG_N,(res&0x80)!=0); c.SetFlag(SREG_Z,res==0); c.SetFlag(SREG_C,(res16&0x100)!=0); c.SetFlag(SREG_S,c.GetFlag(SREG_N)!=c.GetFlag(SREG_V)); return 1 }
	if (opcode&0xFC00) == 0x1800 { d,r := (opcode>>4)&0x1F, (opcode&0x0F)|((opcode>>5)&0x10); op1,op2 := c.Reg[d],c.Reg[r]; res := op1-op2; c.SetFlag(SREG_H,((^op1&op2)|(op2&res)|(res&^op1))&0x08!=0); c.SetFlag(SREG_V,((op1&^op2&^res)|(^op1&op2&res))&0x80!=0); c.SetFlag(SREG_N,(res&0x80)!=0); c.SetFlag(SREG_Z,res==0); c.SetFlag(SREG_C,((^op1&op2)|(op2&res)|(res&^op1))&0x80!=0); c.SetFlag(SREG_S,c.GetFlag(SREG_N)!=c.GetFlag(SREG_V)); c.Reg[d]=res; return 1 }
	if (opcode&0xFC00) == 0x0400 { d,r := (opcode>>4)&0x1F, (opcode&0x0F)|((opcode>>5)&0x10); op1,op2 := c.Reg[d],c.Reg[r]; carry := uint8(0); if c.GetFlag(SREG_C) { carry = 1 }; res := op1-op2-carry; c.SetFlag(SREG_H,((^op1&op2)|(op2&res)|(res&^op1))&0x08!=0); c.SetFlag(SREG_V,((^op1&op2&res)|(op1&^op2&^res))&0x80!=0); c.SetFlag(SREG_N,(res&0x80)!=0); c.SetFlag(SREG_Z,res==0 && c.GetFlag(SREG_Z)); c.SetFlag(SREG_C,((^op1&op2)|(op2&res)|(res&^op1))&0x80!=0); c.SetFlag(SREG_S,c.GetFlag(SREG_N)!=c.GetFlag(SREG_V)); return 1 }
	if (opcode&0xFC00) == 0x0800 { d,r := (opcode>>4)&0x1F, (opcode&0x0F)|((opcode>>5)&0x10); op1,op2 := c.Reg[d],c.Reg[r]; carry := uint8(0); if c.GetFlag(SREG_C) { carry = 1 }; res := op1-op2-carry; c.SetFlag(SREG_H,((^op1&op2)|(op2&res)|(res&^op1))&0x08!=0); c.SetFlag(SREG_V,((^op1&op2&res)|(op1&^op2&^res))&0x80!=0); c.SetFlag(SREG_N,(res&0x80)!=0); c.SetFlag(SREG_Z,res==0 && c.GetFlag(SREG_Z)); c.SetFlag(SREG_C,((^op1&op2)|(op2&res)|(res&^op1))&0x80!=0); c.SetFlag(SREG_S,c.GetFlag(SREG_N)!=c.GetFlag(SREG_V)); c.Reg[d]=res; return 1 }

	// 7. Branch / Jump
	if (opcode&0xF000) == 0xC000 { offset := int16(opcode&0x0FFF); if offset&0x800!=0 {offset|=-4096}; c.PC = uint16(int32(c.PC)+int32(offset)); return 2 }
	if (opcode&0xFC00) == 0xF000 { s,k := uint8(opcode&7), int8((opcode>>3)&0x7F); if k&0x40!=0 {k|=-128}; if c.GetFlag(s) { c.PC = uint16(int32(c.PC)+int32(k)); return 2 }; return 1 }
	if (opcode&0xFC00) == 0xF400 { s,k := uint8(opcode&7), int8((opcode>>3)&0x7F); if k&0x40!=0 {k|=-128}; if !c.GetFlag(s) { c.PC = uint16(int32(c.PC)+int32(k)); return 2 }; return 1 }

	// 8. Data Transfer
	if (opcode&0xFE0F) == 0x9000 { d,k := (opcode>>4)&0x1F, flash[c.PC]; c.PC++; c.Reg[d] = c.Bus.ReadSRAM(k); return 2 }
	if (opcode&0xFE0F) == 0x9200 { r,k := (opcode>>4)&0x1F, flash[c.PC]; c.PC++; c.Bus.WriteSRAM(k, c.Reg[r]); return 2 }
	if (opcode&0xFE0E) == 0x940E { k1,k2,k3 := (opcode>>4)&0x1F, opcode&0x01, flash[c.PC]; c.PC++; target := (uint32(k1)<<17)|(uint32(k2)<<16)|uint32(k3); c.Push(uint8(c.PC)); c.Push(uint8(c.PC>>8)); c.PC = uint16(target); return 4 }
	if (opcode&0xFE0F) == 0x920F { c.Push(c.Reg[(opcode>>4)&0x1F]); return 2 }
	if (opcode&0xFE0F) == 0x900F { c.Reg[(opcode>>4)&0x1F] = c.Pop(); return 2 }

	if (opcode&0xF800) == 0xB800 { c.Bus.WriteIO((opcode&0x0F)|((opcode>>5)&0x30), c.Reg[(opcode>>4)&0x1F]); return 1 }
	if (opcode&0xF800) == 0xB000 { c.Reg[(opcode>>4)&0x1F] = c.Bus.ReadIO((opcode&0x0F)|((opcode>>5)&0x30)); return 1 }

	if (opcode & 0xFE08) == 0xFA00 { r, b := (opcode>>4)&0x1F, uint8(opcode&0x07); c.SetFlag(SREG_T, (c.Reg[r]&(1<<b)) != 0); return 1 }
	if (opcode & 0xFE08) == 0xF800 { d, b := (opcode>>4)&0x1F, uint8(opcode&0x07); if c.GetFlag(SREG_T) { c.Reg[d] |= (1 << b) } else { c.Reg[d] &= ^(1 << b) }; return 1 }

	if (opcode & 0xD208) == 0x8208 { r, q := (opcode>>4)&0x1F, (opcode&0x07)|((opcode&0x0C00)>>7)|((opcode&0x2000)>>8); ptr := (uint16(c.Reg[29]) << 8) | uint16(c.Reg[28]); c.Bus.WriteSRAM(ptr+q, c.Reg[r]); return 2 }
	if (opcode & 0xD208) == 0x8008 { d, q := (opcode>>4)&0x1F, (opcode&0x07)|((opcode&0x0C00)>>7)|((opcode&0x2000)>>8); ptr := (uint16(c.Reg[29]) << 8) | uint16(c.Reg[28]); c.Reg[d] = c.Bus.ReadSRAM(ptr+q); return 2 }
	if (opcode & 0xD208) == 0x8200 { r, q := (opcode>>4)&0x1F, (opcode&0x07)|((opcode&0x0C00)>>7)|((opcode&0x2000)>>8); ptr := (uint16(c.Reg[31]) << 8) | uint16(c.Reg[30]); c.Bus.WriteSRAM(ptr+q, c.Reg[r]); return 2 }
	if (opcode & 0xD208) == 0x8000 { d, q := (opcode>>4)&0x1F, (opcode&0x07)|((opcode&0x0C00)>>7)|((opcode&0x2000)>>8); ptr := (uint16(c.Reg[31]) << 8) | uint16(c.Reg[30]); c.Reg[d] = c.Bus.ReadSRAM(ptr+q); return 2 }

	if (opcode & 0xFE0F) == 0x9009 { d,ptr := (opcode>>4)&0x1F, (uint16(c.Reg[29])<<8)|uint16(c.Reg[28]); c.Reg[d]=c.Bus.ReadSRAM(ptr); ptr++; c.Reg[28],c.Reg[29]=uint8(ptr),uint8(ptr>>8); return 2 }
	if (opcode & 0xFE0F) == 0x900A { d,ptr := (opcode>>4)&0x1F, ((uint16(c.Reg[29])<<8)|uint16(c.Reg[28]))-1; c.Reg[28],c.Reg[29]=uint8(ptr),uint8(ptr>>8); c.Reg[d]=c.Bus.ReadSRAM(ptr); return 2 }
	if (opcode & 0xFE0F) == 0x9209 { r,ptr := (opcode>>4)&0x1F, (uint16(c.Reg[29])<<8)|uint16(c.Reg[28]); c.Bus.WriteSRAM(ptr,c.Reg[r]); ptr++; c.Reg[28],c.Reg[29]=uint8(ptr),uint8(ptr>>8); return 2 }
	if (opcode & 0xFE0F) == 0x920A { r,ptr := (opcode>>4)&0x1F, ((uint16(c.Reg[29])<<8)|uint16(c.Reg[28]))-1; c.Reg[28],c.Reg[29]=uint8(ptr),uint8(ptr>>8); c.Bus.WriteSRAM(ptr,c.Reg[r]); return 2 }

	if (opcode & 0xFE0F) == 0x9001 { d,ptr := (opcode>>4)&0x1F, (uint16(c.Reg[31])<<8)|uint16(c.Reg[30]); c.Reg[d]=c.Bus.ReadSRAM(ptr); ptr++; c.Reg[30],c.Reg[31]=uint8(ptr),uint8(ptr>>8); return 2 }
	if (opcode & 0xFE0F) == 0x9002 { d,ptr := (opcode>>4)&0x1F, ((uint16(c.Reg[31])<<8)|uint16(c.Reg[30]))-1; c.Reg[30],c.Reg[31]=uint8(ptr),uint8(ptr>>8); c.Reg[d]=c.Bus.ReadSRAM(ptr); return 2 }
	if (opcode & 0xFE0F) == 0x9201 { r,ptr := (opcode>>4)&0x1F, (uint16(c.Reg[31])<<8)|uint16(c.Reg[30]); c.Bus.WriteSRAM(ptr,c.Reg[r]); ptr++; c.Reg[30],c.Reg[31]=uint8(ptr),uint8(ptr>>8); return 2 }
	if (opcode & 0xFE0F) == 0x9202 { r,ptr := (opcode>>4)&0x1F, ((uint16(c.Reg[31])<<8)|uint16(c.Reg[30]))-1; c.Reg[30],c.Reg[31]=uint8(ptr),uint8(ptr>>8); c.Bus.WriteSRAM(ptr,c.Reg[r]); return 2 }

	if opcode == 0x9508 { pch,pcl := uint16(c.Pop()), uint16(c.Pop()); c.PC = (pch<<8)|pcl; return 4 }
	if opcode == 0x9518 { pch,pcl := uint16(c.Pop()), uint16(c.Pop()); c.PC = (pch<<8)|pcl; c.SetFlag(SREG_I, true); return 4 }
	if opcode == 0x95E8 { return 2 } // SPM
	if opcode == 0x9588 { c.Bus.SetSleep(true); if c.Halted { c.PC-- }; return 1 }
	if (opcode&0xFE0F) == 0x9406 { d := (opcode >> 4) & 0x1F; op1 := c.Reg[d]; res := op1 >> 1; c.Reg[d] = res; c.SetFlag(SREG_C, (op1&0x01) != 0); c.SetFlag(SREG_Z, res == 0); c.SetFlag(SREG_N, false); c.SetFlag(SREG_V, c.GetFlag(SREG_N) != c.GetFlag(SREG_C)); c.SetFlag(SREG_S, c.GetFlag(SREG_N) != c.GetFlag(SREG_V)); return 1 }

	return 1
}
