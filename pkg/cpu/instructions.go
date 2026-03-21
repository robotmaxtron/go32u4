package cpu

func (c *CPU) Execute(opcode uint16) {
	if opcode == 0x0000 { return }

	// BST / BLD
	if (opcode & 0xFE08) == 0xFA00 { // BST
		r, b := (opcode>>4)&0x1F, uint8(opcode&0x07)
		c.SetFlag(SREG_T, (c.Reg[r]&(1<<b)) != 0); return
	}
	if (opcode & 0xFE08) == 0xF800 { // BLD
		d, b := (opcode>>4)&0x1F, uint8(opcode&0x07)
		if c.GetFlag(SREG_T) { c.Reg[d] |= (1 << b) } else { c.Reg[d] &= ^(1 << b) }; return
	}

	// LDD / STD Displacement (10q0 00?r rrrr ?yyy)
	if (opcode & 0xD208) == 0x8208 { // STD Y+q
		r, q := (opcode>>4)&0x1F, (opcode&0x07)|((opcode&0x0C00)>>7)|((opcode&0x2000)>>8)
		ptr := (uint16(c.Reg[29]) << 8) | uint16(c.Reg[28])
		c.WriteSRAM(ptr+q, c.Reg[r]); c.Cycles++; c.TickPeripheralsHelper(1); return
	}
	if (opcode & 0xD208) == 0x8008 { // LDD Y+q
		d, q := (opcode>>4)&0x1F, (opcode&0x07)|((opcode&0x0C00)>>7)|((opcode&0x2000)>>8)
		ptr := (uint16(c.Reg[29]) << 8) | uint16(c.Reg[28])
		c.Reg[d] = c.ReadSRAM(ptr+q); c.Cycles++; c.TickPeripheralsHelper(1); return
	}
	if (opcode & 0xD208) == 0x8200 { // STD Z+q
		r, q := (opcode>>4)&0x1F, (opcode&0x07)|((opcode&0x0C00)>>7)|((opcode&0x2000)>>8)
		ptr := (uint16(c.Reg[31]) << 8) | uint16(c.Reg[30])
		c.WriteSRAM(ptr+q, c.Reg[r]); c.Cycles++; c.TickPeripheralsHelper(1); return
	}
	if (opcode & 0xD208) == 0x8000 { // LDD Z+q
		d, q := (opcode>>4)&0x1F, (opcode&0x07)|((opcode&0x0C00)>>7)|((opcode&0x2000)>>8)
		ptr := (uint16(c.Reg[31]) << 8) | uint16(c.Reg[30])
		c.Reg[d] = c.ReadSRAM(ptr+q); c.Cycles++; c.TickPeripheralsHelper(1); return
	}

	// LDS / STS
	flash := c.Bus.Flash()
	if (opcode&0xFE0F) == 0x9000 { d := (opcode >> 4) & 0x1F; k := flash[c.PC]; c.PC++; c.Reg[d] = c.ReadSRAM(k); c.Cycles++; c.TickPeripheralsHelper(1); return }
	if (opcode&0xFE0F) == 0x9200 { r := (opcode >> 4) & 0x1F; k := flash[c.PC]; c.PC++; c.WriteSRAM(k, c.Reg[r]); c.Cycles++; c.TickPeripheralsHelper(1); return }

	// Indirect LD/ST
	if (opcode & 0xFE0F) == 0x900D { d,ptr := (opcode>>4)&0x1F, (uint16(c.Reg[27])<<8)|uint16(c.Reg[26]); c.Reg[d]=c.ReadSRAM(ptr); ptr++; c.Reg[26],c.Reg[27]=uint8(ptr),uint8(ptr>>8); c.Cycles++; c.TickPeripheralsHelper(1); return }
	if (opcode & 0xFE0F) == 0x920D { r,ptr := (opcode>>4)&0x1F, (uint16(c.Reg[27])<<8)|uint16(c.Reg[26]); c.WriteSRAM(ptr,c.Reg[r]); ptr++; c.Reg[26],c.Reg[27]=uint8(ptr),uint8(ptr>>8); c.Cycles++; c.TickPeripheralsHelper(1); return }
	if (opcode & 0xFE0F) == 0x900E { d,ptr := (opcode>>4)&0x1F, ((uint16(c.Reg[27])<<8)|uint16(c.Reg[26]))-1; c.Reg[26],c.Reg[27]=uint8(ptr),uint8(ptr>>8); c.Reg[d]=c.ReadSRAM(ptr); c.Cycles++; c.TickPeripheralsHelper(1); return }
	if (opcode & 0xFE0F) == 0x920E { r,ptr := (opcode>>4)&0x1F, ((uint16(c.Reg[27])<<8)|uint16(c.Reg[26]))-1; c.Reg[26],c.Reg[27]=uint8(ptr),uint8(ptr>>8); c.WriteSRAM(ptr,c.Reg[r]); c.Cycles++; c.TickPeripheralsHelper(1); return }
	if (opcode & 0xFE0F) == 0x9009 { d,ptr := (opcode>>4)&0x1F, (uint16(c.Reg[29])<<8)|uint16(c.Reg[28]); c.Reg[d]=c.ReadSRAM(ptr); ptr++; c.Reg[28],c.Reg[29]=uint8(ptr),uint8(ptr>>8); c.Cycles++; c.TickPeripheralsHelper(1); return }
	if (opcode & 0xFE0F) == 0x9209 { r,ptr := (opcode>>4)&0x1F, (uint16(c.Reg[29])<<8)|uint16(c.Reg[28]); c.WriteSRAM(ptr,c.Reg[r]); ptr++; c.Reg[28],c.Reg[29]=uint8(ptr),uint8(ptr>>8); c.Cycles++; c.TickPeripheralsHelper(1); return }
	if (opcode & 0xFE0F) == 0x900A { d,ptr := (opcode>>4)&0x1F, ((uint16(c.Reg[29])<<8)|uint16(c.Reg[28]))-1; c.Reg[28],c.Reg[29]=uint8(ptr),uint8(ptr>>8); c.Reg[d]=c.ReadSRAM(ptr); c.Cycles++; c.TickPeripheralsHelper(1); return }
	if (opcode & 0xFE0F) == 0x920A { r,ptr := (opcode>>4)&0x1F, ((uint16(c.Reg[29])<<8)|uint16(c.Reg[28]))-1; c.Reg[28],c.Reg[29]=uint8(ptr),uint8(ptr>>8); c.WriteSRAM(ptr,c.Reg[r]); c.Cycles++; c.TickPeripheralsHelper(1); return }
	if (opcode & 0xFE0F) == 0x9001 { d,ptr := (opcode>>4)&0x1F, (uint16(c.Reg[31])<<8)|uint16(c.Reg[30]); c.Reg[d]=c.ReadSRAM(ptr); ptr++; c.Reg[30],c.Reg[31]=uint8(ptr),uint8(ptr>>8); c.Cycles++; c.TickPeripheralsHelper(1); return }
	if (opcode & 0xFE0F) == 0x9201 { r,ptr := (opcode>>4)&0x1F, (uint16(c.Reg[31])<<8)|uint16(c.Reg[30]); c.WriteSRAM(ptr,c.Reg[r]); ptr++; c.Reg[30],c.Reg[31]=uint8(ptr),uint8(ptr>>8); c.Cycles++; c.TickPeripheralsHelper(1); return }
	if (opcode & 0xFE0F) == 0x9002 { d,ptr := (opcode>>4)&0x1F, ((uint16(c.Reg[31])<<8)|uint16(c.Reg[30]))-1; c.Reg[30],c.Reg[31]=uint8(ptr),uint8(ptr>>8); c.Reg[d]=c.ReadSRAM(ptr); c.Cycles++; c.TickPeripheralsHelper(1); return }
	if (opcode & 0xFE0F) == 0x9202 { r,ptr := (opcode>>4)&0x1F, ((uint16(c.Reg[31])<<8)|uint16(c.Reg[30]))-1; c.Reg[30],c.Reg[31]=uint8(ptr),uint8(ptr>>8); c.WriteSRAM(ptr,c.Reg[r]); c.Cycles++; c.TickPeripheralsHelper(1); return }

	// Conditional Branches (BREQ, BRNE, etc.) - Handled by BRBS/BRBC
	// BRBS
	if (opcode&0xFC00) == 0xF000 { s,k := uint8(opcode&7), int8((opcode>>3)&0x7F); if k&0x40!=0 {k|=-128}; if c.GetFlag(s) { c.PC = uint16(int32(c.PC)+int32(k)); c.Cycles++; c.TickPeripherals(1) }; return }
	// BRBC
	if (opcode&0xFC00) == 0xF400 { s,k := uint8(opcode&7), int8((opcode>>3)&0x7F); if k&0x40!=0 {k|=-128}; if !c.GetFlag(s) { c.PC = uint16(int32(c.PC)+int32(k)); c.Cycles++; c.TickPeripherals(1) }; return }

	// LDI, ANDI/CBR, ORI/SBR, SUBI, CPI
	if (opcode&0xF000) == 0xE000 { c.Reg[16+uint8((opcode>>4)&0x0F)] = uint8(opcode&0x0F)|uint8((opcode>>4)&0xF0); return }
	if (opcode&0xF000) == 0x7000 { d,k := 16+uint8((opcode>>4)&0x0F), uint8(opcode&0x0F)|uint8((opcode>>4)&0xF0); c.Reg[d]&=k; res := c.Reg[d]; c.SetFlag(SREG_V,false); c.SetFlag(SREG_N,(res&0x80)!=0); c.SetFlag(SREG_Z,res==0); c.SetFlag(SREG_S,c.GetFlag(SREG_N)); return }
	if (opcode&0xF000) == 0x6000 { d,k := 16+uint8((opcode>>4)&0x0F), uint8(opcode&0x0F)|uint8((opcode>>4)&0xF0); c.Reg[d]|=k; res := c.Reg[d]; c.SetFlag(SREG_V,false); c.SetFlag(SREG_N,(res&0x80)!=0); c.SetFlag(SREG_Z,res==0); c.SetFlag(SREG_S,c.GetFlag(SREG_N)); return }
	if (opcode&0xF000) == 0x5000 { d,k := 16+uint8((opcode>>4)&0x0F), uint8(opcode&0x0F)|uint8((opcode>>4)&0xF0); op1 := c.Reg[d]; res := op1-k; c.SetFlag(SREG_H,((^op1&k)|(k&res)|(res&^op1))&0x08!=0); c.SetFlag(SREG_V,((op1&^k&^res)|(^op1&k&res))&0x80!=0); c.SetFlag(SREG_N,(res&0x80)!=0); c.SetFlag(SREG_Z,res==0); c.SetFlag(SREG_C,((^op1&k)|(k&res)|(res&^op1))&0x80!=0); c.SetFlag(SREG_S,c.GetFlag(SREG_N)!=c.GetFlag(SREG_V)); c.Reg[d]=res; return }
	if (opcode&0xF000) == 0x3000 { d,k := 16+uint8((opcode>>4)&0x0F), uint8(opcode&0x0F)|uint8((opcode>>4)&0xF0); op1 := c.Reg[d]; res := op1-k; c.SetFlag(SREG_H,((^op1&k)|(k&res)|(res&^op1))&0x08!=0); c.SetFlag(SREG_V,((op1&^k&^res)|(^op1&k&res))&0x80!=0); c.SetFlag(SREG_N,(res&0x80)!=0); c.SetFlag(SREG_Z,res==0); c.SetFlag(SREG_C,((^op1&k)|(k&res)|(res&^op1))&0x80!=0); c.SetFlag(SREG_S,c.GetFlag(SREG_N)!=c.GetFlag(SREG_V)); return }

	// RJMP
	if (opcode&0xF000) == 0xC000 { offset := int16(opcode&0x0FFF); if offset&0x800!=0 {offset|=-4096}; c.PC = uint16(int32(c.PC)+int32(offset)); c.Cycles++; c.TickPeripheralsHelper(1); return }

	// Reg-Reg Ops
	if (opcode&0xFC00) == 0x0C00 { d,r := (opcode>>4)&0x1F, (opcode&0x0F)|((opcode>>5)&0x10); op1,op2 := c.Reg[d],c.Reg[r]; res16 := uint16(op1)+uint16(op2); res := uint8(res16); c.Reg[d]=res; c.SetFlag(SREG_H,((op1&op2)|(op2&^res)|(^res&op1))&0x08!=0); c.SetFlag(SREG_V,((op1&op2&^res)|(^op1&^op2&res))&0x80!=0); c.SetFlag(SREG_N,(res&0x80)!=0); c.SetFlag(SREG_Z,res==0); c.SetFlag(SREG_C,(res16&0x100)!=0); c.SetFlag(SREG_S,c.GetFlag(SREG_N)!=c.GetFlag(SREG_V)); return }
	if (opcode&0xFC00) == 0x1800 { d,r := (opcode>>4)&0x1F, (opcode&0x0F)|((opcode>>5)&0x10); op1,op2 := c.Reg[d],c.Reg[r]; res := op1-op2; c.SetFlag(SREG_H,((^op1&op2)|(op2&res)|(res&^op1))&0x08!=0); c.SetFlag(SREG_V,((op1&^op2&^res)|(^op1&op2&res))&0x80!=0); c.SetFlag(SREG_N,(res&0x80)!=0); c.SetFlag(SREG_Z,res==0); c.SetFlag(SREG_C,((^op1&op2)|(op2&res)|(res&^op1))&0x80!=0); c.SetFlag(SREG_S,c.GetFlag(SREG_N)!=c.GetFlag(SREG_V)); c.Reg[d]=res; return }
	if (opcode&0xFC00) == 0x1400 { d,r := (opcode>>4)&0x1F, (opcode&0x0F)|((opcode>>5)&0x10); op1,op2 := c.Reg[d],c.Reg[r]; res := op1-op2; c.SetFlag(SREG_H,((^op1&op2)|(op2&res)|(res&^op1))&0x08!=0); c.SetFlag(SREG_V,((op1&^op2&^res)|(^op1&op2&res))&0x80!=0); c.SetFlag(SREG_N,(res&0x80)!=0); c.SetFlag(SREG_Z,res==0); c.SetFlag(SREG_C,((^op1&op2)|(op2&res)|(res&^op1))&0x80!=0); c.SetFlag(SREG_S,c.GetFlag(SREG_N)!=c.GetFlag(SREG_V)); return }
	if (opcode&0xFC00) == 0x2000 { d,r := (opcode>>4)&0x1F, (opcode&0x0F)|((opcode>>5)&0x10); res := c.Reg[d]&c.Reg[r]; c.Reg[d]=res; c.SetFlag(SREG_V,false); c.SetFlag(SREG_N,(res&0x80)!=0); c.SetFlag(SREG_Z,res==0); c.SetFlag(SREG_S,c.GetFlag(SREG_N)); return }
	if (opcode&0xFC00) == 0x2800 { d,r := (opcode>>4)&0x1F, (opcode&0x0F)|((opcode>>5)&0x10); res := c.Reg[d]|c.Reg[r]; c.Reg[d]=res; c.SetFlag(SREG_V,false); c.SetFlag(SREG_N,(res&0x80)!=0); c.SetFlag(SREG_Z,res==0); c.SetFlag(SREG_S,c.GetFlag(SREG_N)); return }
	if (opcode&0xFC00) == 0x2400 { d,r := (opcode>>4)&0x1F, (opcode&0x0F)|((opcode>>5)&0x10); res := c.Reg[d]^c.Reg[r]; c.Reg[d]=res; c.SetFlag(SREG_V,false); c.SetFlag(SREG_N,(res&0x80)!=0); c.SetFlag(SREG_Z,res==0); c.SetFlag(SREG_S,c.GetFlag(SREG_N)); return }

	// Fallback Indirect (No Displacement)
	if (opcode&0xFE0F) == 0x900C { d := (opcode>>4)&0x1F; c.Reg[d]=c.ReadSRAM((uint16(c.Reg[27])<<8)|uint16(c.Reg[26])); c.Cycles++; c.TickPeripheralsHelper(1); return }
	if (opcode&0xFE0F) == 0x920C { r := (opcode>>4)&0x1F; c.WriteSRAM((uint16(c.Reg[27])<<8)|uint16(c.Reg[26]), c.Reg[r]); c.Cycles++; c.TickPeripheralsHelper(1); return }
	if (opcode&0xFE0F) == 0x8008 { d := (opcode>>4)&0x1F; c.Reg[d]=c.ReadSRAM((uint16(c.Reg[29])<<8)|uint16(c.Reg[28])); c.Cycles++; c.TickPeripheralsHelper(1); return }
	if (opcode&0xFE0F) == 0x8208 { r := (opcode>>4)&0x1F; c.WriteSRAM((uint16(c.Reg[29])<<8)|uint16(c.Reg[28]), c.Reg[r]); c.Cycles++; c.TickPeripheralsHelper(1); return }
	if (opcode&0xFE0F) == 0x8000 { d := (opcode>>4)&0x1F; c.Reg[d]=c.ReadSRAM((uint16(c.Reg[31])<<8)|uint16(c.Reg[30])); c.Cycles++; c.TickPeripheralsHelper(1); return }
	if (opcode&0xFE0F) == 0x8200 { r := (opcode>>4)&0x1F; c.WriteSRAM((uint16(c.Reg[31])<<8)|uint16(c.Reg[30]), c.Reg[r]); c.Cycles++; c.TickPeripheralsHelper(1); return }

	// PUSH / POP
	if (opcode&0xFE0F) == 0x920F { c.Push(c.Reg[(opcode>>4)&0x1F]); c.Cycles++; c.TickPeripheralsHelper(1); return }
	if (opcode&0xFE0F) == 0x900F { c.Reg[(opcode>>4)&0x1F] = c.Pop(); c.Cycles++; c.TickPeripheralsHelper(1); return }

	// OUT / IN
	if (opcode&0xF800) == 0xB800 { c.WriteIO((opcode&0x0F)|((opcode>>5)&0x30), c.Reg[(opcode>>4)&0x1F]); return }
	if (opcode&0xF800) == 0xB000 { c.Reg[(opcode>>4)&0x1F] = c.ReadIO((opcode&0x0F)|((opcode>>5)&0x30)); return }

	// CALL / RET / RETI
	if (opcode&0xFE0E) == 0x940E { k1,k2 := (opcode>>4)&0x1F, opcode&0x01; nextOp := flash[c.PC]; c.PC++; target := (uint32(k1)<<17)|(uint32(k2)<<16)|uint32(nextOp); c.Push(uint8(c.PC)); c.Push(uint8(c.PC>>8)); c.PC = uint16(target); c.Cycles += 3; c.TickPeripheralsHelper(3); return }
	if opcode == 0x9508 { pch,pcl := uint16(c.Pop()), uint16(c.Pop()); c.PC = (pch<<8)|pcl; c.Cycles+=3; c.TickPeripheralsHelper(3); return }
	if opcode == 0x9518 { pch,pcl := uint16(c.Pop()), uint16(c.Pop()); c.PC = (pch<<8)|pcl; c.SetFlag(SREG_I, true); c.Cycles+=3; c.TickPeripheralsHelper(3); return }

	// Misc
	if opcode == 0x95C8 { z := (uint16(c.Reg[31])<<8)|uint16(c.Reg[30]); word := flash[z>>1]; if z&0x01 == 0 { c.Reg[0] = uint8(word) } else { c.Reg[0] = uint8(word>>8) }; c.Cycles += 2; c.TickPeripheralsHelper(2); return }
	if (opcode&0xFE0F) == 0x9004 { d,z := (opcode>>4)&0x1F, (uint16(c.Reg[31])<<8)|uint16(c.Reg[30]); word := flash[z>>1]; if z&0x01 == 0 { c.Reg[d] = uint8(word) } else { c.Reg[d] = uint8(word>>8) }; c.Cycles += 2; c.TickPeripheralsHelper(2); return }
	if (opcode&0xFE0F) == 0x9005 { d,z := (opcode>>4)&0x1F, (uint16(c.Reg[31])<<8)|uint16(c.Reg[30]); word := flash[z>>1]; if z&0x01 == 0 { c.Reg[d] = uint8(word) } else { c.Reg[d] = uint8(word>>8) }; z++; c.Reg[30],c.Reg[31] = uint8(z),uint8(z>>8); c.Cycles += 2; c.TickPeripheralsHelper(2); return }
	if opcode == 0x9588 { c.Bus.SetSleep(true); if c.Halted { c.PC-- }; return }
	if opcode == 0x95E8 { // SPM
		spmcsr := c.ReadIO(0x37)
		if spmcsr&0x01 != 0 {
			z := (uint16(c.Reg[31]) << 8) | uint16(c.Reg[30])
			if spmcsr&0x02 != 0 { // PGERS
				c.Bus.FlashErase(z >> 1)
			} else if spmcsr&0x04 != 0 { // PGWRT
				c.Bus.FlashCommit(z >> 1)
			} else { // Buffer Write
				val := (uint16(c.Reg[1]) << 8) | uint16(c.Reg[0])
				c.Bus.FlashWrite(z>>1, val)
			}
		}
		c.Cycles++; c.TickPeripheralsHelper(1); return
	}
}
