package go6502

import (
	"fmt"
	"strings"
)

// status register bits
const (
	sCarry = iota
	sZero
	sInterrupt
	sDecimal
	sBreak
	_
	sOverflow
	sNegative
)

type Cpu struct {
	pc       address
	ac       byte
	x        byte
	y        byte
	sp       byte
	sr       byte
	Bus      *Bus
	debugger *Debugger
	ExitChan chan int
}

type Iop struct {
	in   *Instruction
	op8  uint8
	op16 address
}

func (iop *Iop) String() string {
	return fmt.Sprintf("%v op8:0x%02X op16:0x%04X", iop.in, iop.op8, iop.op16)
}

func (c *Cpu) AttachDebugger(d *Debugger) {
	c.debugger = d
}

func (c *Cpu) Reset() {
	c.pc = c.Bus.Read16(0xFFFC)
	c.ac = 0x00
	c.x = 0x00
	c.y = 0x00
	c.sp = 0xFF // address relative to second page of memory (0x0100 ~ 0x01FF)
	c.sr = 0x00
}

func (c *Cpu) Step() {
	op := c.Bus.Read(c.pc)
	in := findInstruction(op)
	iop := c.readOperand(in)
	if c.debugger != nil {
		c.debugger.BeforeExecute(iop)
	}
	c.pc += address(in.bytes)
	c.Execute(iop)
}

func (c *Cpu) readOperand(in *Instruction) *Iop {
	// read instruction with operand
	iop := &Iop{in: in}
	switch in.bytes {
	case 1: // no operand
	case 2:
		iop.op8 = c.Bus.Read(c.pc + 1)
	case 3:
		iop.op16 = c.Bus.Read16(c.pc + 1)
	default:
		panic(fmt.Sprintf("unhandled instruction length: %d", in.bytes))
	}
	return iop
}

func (c *Cpu) String() string {
	return fmt.Sprintf(
		"CPU pc:0x%04X ac:0x%02X x:0x%02X y:0x%02X sp:0x%02X sr:%s",
		c.pc, c.ac, c.x, c.y, c.sp,
		c.statusString(),
	)
}

func (c *Cpu) StackHead(offset int8) address {
	return address(0x0100) + address(c.sp) + address(offset)
}

func (c *Cpu) resolveOperand(iop *Iop) uint8 {
	switch iop.in.addressing {
	case immediate:
		return iop.op8
	default:
		return c.Bus.Read(c.memoryAddress(iop))
	}
}

func (c *Cpu) memoryAddress(iop *Iop) address {
	switch iop.in.addressing {
	case absolute:
		return iop.op16
	case absoluteX:
		return iop.op16 + address(c.x)
	case absoluteY:
		return iop.op16 + address(c.y)

	// Indexed Indirect (X)
	// Operand is the zero-page location of a little-endian 16-bit base address.
	// The X register is added (wrapping; discarding overflow) before loading.
	// The resulting address loaded from (base+X) becomes the effective operand.
	// (base + X) must be in zero-page.
	case indirectX:
		location := address(iop.op8 + c.x)
		if location == 0xFF {
			panic("Indexed indirect high-byte not on zero page.")
		}
		return c.Bus.Read16(location)

	// Indirect Indexed (Y)
	// Operand is the zero-page location of a little-endian 16-bit address.
	// The address is loaded, and then the Y register is added to it.
	// The resulting loaded_address + Y becomes the effective operand.
	case indirectY:
		return c.Bus.Read16(address(iop.op8)) + address(c.y)

	case zeropage:
		return address(iop.op8)
	case zeropageX:
		return address(iop.op8 + c.x)
	case zeropageY:
		return address(iop.op8 + c.y)
	default:
		panic("unhandled addressing")
	}
}

func (c *Cpu) getStatus(bit uint8) bool {
	return c.getStatusInt(bit) == 1
}

func (c *Cpu) getStatusInt(bit uint8) uint8 {
	return (c.sr >> bit) & 1
}

func (c *Cpu) setStatus(bit uint8, state bool) {
	if state {
		c.sr |= 1 << bit
	} else {
		c.sr &^= 1 << bit
	}
}

func (c *Cpu) updateStatus(value uint8) {
	c.setStatus(sZero, value == 0)
	c.setStatus(sNegative, (value>>7) == 1)
}

func (c *Cpu) statusString() string {
	// this is horrible, I think. Should be much simpler?
	var chars = "nv_bdizc"
	var out [8]string
	for i := 0; i < 8; i++ {
		if c.getStatus(uint8(7 - i)) {
			out[i] = string(chars[i])
		} else {
			out[i] = "-"
		}
	}
	return strings.Join(out[0:], "")
}

func (c *Cpu) branch(iop *Iop) {
	relative := int8(iop.op8) // signed
	if relative >= 0 {
		c.pc += address(relative)
	} else {
		c.pc -= address(-relative)
	}
}

func (c *Cpu) Execute(iop *Iop) {
	switch iop.in.id {
	case ADC:
		c.ADC(iop)
	case AND:
		c.AND(iop)
	case ASL:
		c.ASL(iop)
	case BCC:
		c.BCC(iop)
	case BCS:
		c.BCS(iop)
	case BEQ:
		c.BEQ(iop)
	case BMI:
		c.BMI(iop)
	case BNE:
		c.BNE(iop)
	case BPL:
		c.BPL(iop)
	case CLC:
		c.CLC(iop)
	case CLD:
		c.CLD(iop)
	case CLI:
		c.CLI(iop)
	case CMP:
		c.CMP(iop)
	case CPX:
		c.CPX(iop)
	case CPY:
		c.CPY(iop)
	case DEC:
		c.DEC(iop)
	case DEX:
		c.DEX(iop)
	case DEY:
		c.DEY(iop)
	case EOR:
		c.EOR(iop)
	case INC:
		c.INC(iop)
	case INX:
		c.INX(iop)
	case INY:
		c.INY(iop)
	case JMP:
		c.JMP(iop)
	case JSR:
		c.JSR(iop)
	case LDA:
		c.LDA(iop)
	case LDX:
		c.LDX(iop)
	case LDY:
		c.LDY(iop)
	case LSR:
		c.LSR(iop)
	case NOP:
		c.NOP(iop)
	case ORA:
		c.ORA(iop)
	case PHA:
		c.PHA(iop)
	case PLA:
		c.PLA(iop)
	case ROL:
		c.ROL(iop)
	case RTS:
		c.RTS(iop)
	case SBC:
		c.SBC(iop)
	case SEI:
		c.SEI(iop)
	case STA:
		c.STA(iop)
	case STX:
		c.STX(iop)
	case STY:
		c.STY(iop)
	case TAX:
		c.TAX(iop)
	case TAY:
		c.TAY(iop)
	case TXA:
		c.TXA(iop)
	case TXS:
		c.TXS(iop)
	case TYA:
		c.TYA(iop)
	case _END:
		c._END(iop)
	default:
		panic(fmt.Sprintf("unhandled instruction: %v", iop.in.name()))
	}
}

// add with carry
func (c *Cpu) ADC(iop *Iop) {
	value16 := uint16(c.ac) + uint16(c.resolveOperand(iop)) + uint16(c.getStatusInt(sCarry))
	c.setStatus(sCarry, value16 > 0xFF)
	c.ac = uint8(value16)
	c.updateStatus(c.ac)
}

// bitwise AND with accumulator
func (c *Cpu) AND(iop *Iop) {
	c.ac &= c.resolveOperand(iop)
	c.updateStatus(c.ac)
}

// arithmetic shift left
func (c *Cpu) ASL(iop *Iop) {
	switch iop.in.addressing {
	case accumulator:
		c.setStatus(sCarry, (c.ac>>7) == 1) // carry = old bit 7
		c.ac <<= 1
		c.updateStatus(c.ac)
	default:
		address := c.memoryAddress(iop)
		value := c.Bus.Read(address)
		c.setStatus(sCarry, (value>>7) == 1) // carry = old bit 7
		value <<= 1
		c.Bus.Write(address, value)
		c.updateStatus(value)
	}
}

// branch if carry clear
func (c *Cpu) BCC(iop *Iop) {
	if !c.getStatus(sCarry) {
		c.branch(iop)
	}
}

// branch on carry (when carry set)
func (c *Cpu) BCS(iop *Iop) {
	if c.getStatus(sCarry) {
		c.branch(iop)
	}
}

// branch on equal (zero set)
// (branch on z = 1)
func (c *Cpu) BEQ(iop *Iop) {
	if c.getStatus(sZero) {
		c.branch(iop)
	}
}

// branch on result minus (status negative)
func (c *Cpu) BMI(iop *Iop) {
	if c.getStatus(sNegative) {
		c.branch(iop)
	}
}

// branch on not-equal (zero clear)
func (c *Cpu) BNE(iop *Iop) {
	if !c.getStatus(sZero) {
		c.branch(iop)
	}
}

// branch on not-negative
func (c *Cpu) BPL(iop *Iop) {
	if !c.getStatus(sNegative) {
		c.branch(iop)
	}
}

// clear carry
func (c *Cpu) CLC(iop *Iop) {
	c.setStatus(sCarry, false)
}

// clear decimal
func (c *Cpu) CLD(iop *Iop) {
	c.setStatus(sDecimal, false)
}

// clear interrupt mask (enable maskable interrupts)
func (c *Cpu) CLI(iop *Iop) {
	c.setStatus(sInterrupt, true)
}

// compare (with accumulator)
func (c *Cpu) CMP(iop *Iop) {
	value := c.resolveOperand(iop)
	c.setStatus(sCarry, c.ac >= value)
	c.updateStatus(c.ac - value)
}

// compare X
func (c *Cpu) CPX(iop *Iop) {
	value := c.resolveOperand(iop)
	c.setStatus(sCarry, c.x >= value)
	c.updateStatus(c.x - value)
}

// compare Y
func (c *Cpu) CPY(iop *Iop) {
	value := c.resolveOperand(iop)
	c.setStatus(sCarry, c.y >= value)
	c.updateStatus(c.y - value)
}

// decrement memory
func (c *Cpu) DEC(iop *Iop) {
	address := c.memoryAddress(iop)
	value := c.Bus.Read(address) - 1
	c.Bus.Write(address, value)
	c.updateStatus(value)
}

// decrement x
func (c *Cpu) DEX(iop *Iop) {
	c.x--
	c.updateStatus(c.x)
}

// decrement y
func (c *Cpu) DEY(iop *Iop) {
	c.y--
	c.updateStatus(c.y)
}

// Exclusive OR (accumulator)
func (c *Cpu) EOR(iop *Iop) {
	value := c.resolveOperand(iop)
	c.ac ^= value
	c.updateStatus(c.ac)
}

// increment value in memory
func (c *Cpu) INC(iop *Iop) {
	address := c.memoryAddress(iop)
	value := c.Bus.Read(address) + 1
	c.Bus.Write(address, value)
	c.updateStatus(value)
}

// increment x
func (c *Cpu) INX(iop *Iop) {
	c.x++
	c.updateStatus(c.x)
}

// increment y
func (c *Cpu) INY(iop *Iop) {
	c.y++
	c.updateStatus(c.y)
}

// jump
func (c *Cpu) JMP(iop *Iop) {
	c.pc = c.memoryAddress(iop)
}

// jump to subroutine
func (c *Cpu) JSR(iop *Iop) {
	c.Bus.Write16(c.StackHead(-1), c.pc-1)
	c.sp -= 2
	c.pc = iop.op16
}

// load accumulator
func (c *Cpu) LDA(iop *Iop) {
	c.ac = c.resolveOperand(iop)
	c.updateStatus(c.ac)
}

// load Y
func (c *Cpu) LDY(iop *Iop) {
	c.y = c.resolveOperand(iop)
	c.updateStatus(c.y)
}

// load X
func (c *Cpu) LDX(iop *Iop) {
	c.x = c.resolveOperand(iop)
	c.updateStatus(c.x)
}

// logical shift right.
func (c *Cpu) LSR(iop *Iop) {
	// TODO: general support for memory-modifying instructions (ASL, LSR, ROL, ROR)
	switch iop.in.addressing {
	case accumulator:
		c.setStatus(sCarry, c.ac&1 == 1)
		c.ac >>= 1
		c.updateStatus(c.ac)
	case zeropageX:
		address := address(iop.op8 + c.x)
		value := c.Bus.Read(address)
		// TODO: carry?
		value >>= 1
		c.updateStatus(value)
		c.Bus.Write(address, value)
	default:
		panic("LSR addressing mode not implemented")
	}
}

// no operation
func (c *Cpu) NOP(iop *Iop) {
}

// bitwise OR memory with accumulator
func (c *Cpu) ORA(iop *Iop) {
	c.ac |= c.resolveOperand(iop)
	c.updateStatus(c.ac)
}

// push accumulator
func (c *Cpu) PHA(iop *Iop) {
	c.Bus.Write(0x0100+address(c.sp), c.ac)
	c.sp--
}

// pull accumulator
func (c *Cpu) PLA(iop *Iop) {
	c.sp++
	c.ac = c.Bus.Read(0x0100 + address(c.sp))
}

// bitwise rotate left
// SR carry into bit 0, original bit 7 into SR carry.
func (c *Cpu) ROL(iop *Iop) {
	// TODO: general support for memory-modifying instructions (ASL, LSR, ROL, ROR)
	carry := c.getStatusInt(sCarry)
	switch iop.in.addressing {
	case accumulator:
		c.setStatus(sCarry, (c.ac>>7) == 1)
		c.ac = (c.ac << 1) | carry
		c.updateStatus(c.ac)
	default:
		address := c.memoryAddress(iop)
		value := c.Bus.Read(address)
		c.setStatus(sCarry, (value>>7) == 1)
		value = (value << 1) | carry
		c.Bus.Write(address, value)
		c.updateStatus(value)
	}
}

// return from subroutine
func (c *Cpu) RTS(iop *Iop) {
	c.pc = c.Bus.Read16(c.StackHead(1))
	c.sp += 2
	c.pc += 1
}

// subract with carry
// TODO: overflow status
func (c *Cpu) SBC(iop *Iop) {
	valueSigned := int16(c.ac) - int16(c.resolveOperand(iop))
	if !c.getStatus(sCarry) {
		valueSigned--
	}
	c.setStatus(sCarry, valueSigned < 0)
	c.ac = uint8(valueSigned)
}

// set interrupt mask (disable maskable interrupts)
func (c *Cpu) SEI(iop *Iop) {
	c.setStatus(sInterrupt, false)
}

// store from accumulator
func (c *Cpu) STA(iop *Iop) {
	c.Bus.Write(c.memoryAddress(iop), c.ac)
}

// store from X
func (c *Cpu) STX(iop *Iop) {
	c.Bus.Write(c.memoryAddress(iop), c.x)
}

// store from Y
func (c *Cpu) STY(iop *Iop) {
	c.Bus.Write(c.memoryAddress(iop), c.y)
}

// transfer accumulator for index Y
func (c *Cpu) TAX(iop *Iop) {
	c.x = c.ac
	c.updateStatus(c.x)
}

// transfer accumulator for index Y
func (c *Cpu) TAY(iop *Iop) {
	c.y = c.ac
	c.updateStatus(c.y)
}

// Copies the current contents of the X register into the accumulator and sets
// the zero and negative flags as appropriate.
func (c *Cpu) TXA(iop *Iop) {
	c.ac = c.x
	c.updateStatus(c.ac)
}

// transfer X to stack pointer
func (c *Cpu) TXS(iop *Iop) {
	c.sp = c.x
	c.updateStatus(c.sp)
}

// Transfer Y to Accumulator
func (c *Cpu) TYA(iop *Iop) {
	c.ac = c.y
	c.updateStatus(c.ac)
}

// Custom go6502 instruction.
// Exit, with contents of X register as exit status.
func (c *Cpu) _END(iop *Iop) {
	c.ExitChan <- int(c.x)
}
