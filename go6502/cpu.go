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

func (c *Cpu) AttachDebugger(d *Debugger) {
	c.debugger = d
}

// Reset the CPU, emulating triggering the RESB line.
// From 65C02 manual: All Registers are initialized by software except the
// Decimal and Interrupt disable mode select bits of the Processor Status
// Register (P) are initialized by hardware. ... The program counter is loaded
// with the reset vector from locations FFFC (low byte) and FFFD (high byte).
func (c *Cpu) Reset() {
	c.pc = c.Bus.Read16(0xFFFC)
	c.sr = 0x34 // Manual says xx1101xx, this sets 00110100.
}

func (c *Cpu) Step() {
	op := c.Bus.Read(c.pc)
	ot := findInstruction(op)
	iop := c.readOperand(ot)
	if c.debugger != nil {
		c.debugger.BeforeExecute(iop)
	}
	c.pc += address(ot.bytes)
	c.Execute(iop)
}

func (c *Cpu) readOperand(ot *OpType) *Iop {
	// read instruction with operand
	iop := &Iop{ot: ot}
	switch ot.bytes {
	case 1: // no operand
	case 2:
		iop.op8 = c.Bus.Read(c.pc + 1)
	case 3:
		iop.op16 = c.Bus.Read16(c.pc + 1)
	default:
		panic(fmt.Sprintf("unhandled instruction length: %d", ot.bytes))
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
	switch iop.ot.addressing {
	case immediate:
		return iop.op8
	default:
		return c.Bus.Read(c.memoryAddress(iop))
	}
}

func (c *Cpu) memoryAddress(iop *Iop) address {
	switch iop.ot.addressing {
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
	chars := "nv_bdizc"
	out := make([]string, 8)
	for i := 0; i < 8; i++ {
		if c.getStatus(uint8(7 - i)) {
			out[i] = string(chars[i])
		} else {
			out[i] = "-"
		}
	}
	return strings.Join(out, "")
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
	switch iop.ot.id {
	case adc:
		c.ADC(iop)
	case and:
		c.AND(iop)
	case asl:
		c.ASL(iop)
	case bcc:
		c.BCC(iop)
	case bcs:
		c.BCS(iop)
	case beq:
		c.BEQ(iop)
	case bmi:
		c.BMI(iop)
	case bne:
		c.BNE(iop)
	case bpl:
		c.BPL(iop)
	case clc:
		c.CLC(iop)
	case cld:
		c.CLD(iop)
	case cli:
		c.CLI(iop)
	case cmp:
		c.CMP(iop)
	case cpx:
		c.CPX(iop)
	case cpy:
		c.CPY(iop)
	case dec:
		c.DEC(iop)
	case dex:
		c.DEX(iop)
	case dey:
		c.DEY(iop)
	case eor:
		c.EOR(iop)
	case inc:
		c.INC(iop)
	case inx:
		c.INX(iop)
	case iny:
		c.INY(iop)
	case jmp:
		c.JMP(iop)
	case jsr:
		c.JSR(iop)
	case lda:
		c.LDA(iop)
	case ldx:
		c.LDX(iop)
	case ldy:
		c.LDY(iop)
	case lsr:
		c.LSR(iop)
	case nop:
		c.NOP(iop)
	case ora:
		c.ORA(iop)
	case pha:
		c.PHA(iop)
	case pla:
		c.PLA(iop)
	case rol:
		c.ROL(iop)
	case rts:
		c.RTS(iop)
	case sbc:
		c.SBC(iop)
	case sei:
		c.SEI(iop)
	case sta:
		c.STA(iop)
	case stx:
		c.STX(iop)
	case sty:
		c.STY(iop)
	case tax:
		c.TAX(iop)
	case tay:
		c.TAY(iop)
	case txa:
		c.TXA(iop)
	case txs:
		c.TXS(iop)
	case tya:
		c.TYA(iop)
	case _end:
		c._END(iop)
	default:
		panic(fmt.Sprintf("unhandled instruction: %v", iop))
	}
}

// ADC: Add memory and carry to accumulator.
func (c *Cpu) ADC(iop *Iop) {
	value16 := uint16(c.ac) + uint16(c.resolveOperand(iop)) + uint16(c.getStatusInt(sCarry))
	c.setStatus(sCarry, value16 > 0xFF)
	c.ac = uint8(value16)
	c.updateStatus(c.ac)
}

// AND: And accumulator with memory.
func (c *Cpu) AND(iop *Iop) {
	c.ac &= c.resolveOperand(iop)
	c.updateStatus(c.ac)
}

// ASL: Shift memory or accumulator left one bit.
func (c *Cpu) ASL(iop *Iop) {
	switch iop.ot.addressing {
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

// BCC: Branch if carry clear.
func (c *Cpu) BCC(iop *Iop) {
	if !c.getStatus(sCarry) {
		c.branch(iop)
	}
}

// BCS: Branch if carry set.
func (c *Cpu) BCS(iop *Iop) {
	if c.getStatus(sCarry) {
		c.branch(iop)
	}
}

// BEQ: Branch if equal (z=1).
func (c *Cpu) BEQ(iop *Iop) {
	if c.getStatus(sZero) {
		c.branch(iop)
	}
}

// BMI: Branch if negative.
func (c *Cpu) BMI(iop *Iop) {
	if c.getStatus(sNegative) {
		c.branch(iop)
	}
}

// BNE: Branch if not equal.
func (c *Cpu) BNE(iop *Iop) {
	if !c.getStatus(sZero) {
		c.branch(iop)
	}
}

// BPL: Branch if positive.
func (c *Cpu) BPL(iop *Iop) {
	if !c.getStatus(sNegative) {
		c.branch(iop)
	}
}

// CLC: Clear carry flag.
func (c *Cpu) CLC(iop *Iop) {
	c.setStatus(sCarry, false)
}

// CLD: Clear decimal mode flag.
func (c *Cpu) CLD(iop *Iop) {
	c.setStatus(sDecimal, false)
}

// CLI: Clear interrupt-disable flag.
func (c *Cpu) CLI(iop *Iop) {
	c.setStatus(sInterrupt, true)
}

// CMP: Compare accumulator with memory.
func (c *Cpu) CMP(iop *Iop) {
	value := c.resolveOperand(iop)
	c.setStatus(sCarry, c.ac >= value)
	c.updateStatus(c.ac - value)
}

// CPX: Compare index register X with memory.
func (c *Cpu) CPX(iop *Iop) {
	value := c.resolveOperand(iop)
	c.setStatus(sCarry, c.x >= value)
	c.updateStatus(c.x - value)
}

// CPY: Compare index register Y with memory.
func (c *Cpu) CPY(iop *Iop) {
	value := c.resolveOperand(iop)
	c.setStatus(sCarry, c.y >= value)
	c.updateStatus(c.y - value)
}

// DEC: Decrement.
func (c *Cpu) DEC(iop *Iop) {
	address := c.memoryAddress(iop)
	value := c.Bus.Read(address) - 1
	c.Bus.Write(address, value)
	c.updateStatus(value)
}

// DEX: Decrement index register X.
func (c *Cpu) DEX(iop *Iop) {
	c.x--
	c.updateStatus(c.x)
}

// DEY: Decrement index register Y.
func (c *Cpu) DEY(iop *Iop) {
	c.y--
	c.updateStatus(c.y)
}

// EOR: Exclusive-OR accumulator with memory.
func (c *Cpu) EOR(iop *Iop) {
	value := c.resolveOperand(iop)
	c.ac ^= value
	c.updateStatus(c.ac)
}

// INC: Increment.
func (c *Cpu) INC(iop *Iop) {
	address := c.memoryAddress(iop)
	value := c.Bus.Read(address) + 1
	c.Bus.Write(address, value)
	c.updateStatus(value)
}

// INX: Increment index register X.
func (c *Cpu) INX(iop *Iop) {
	c.x++
	c.updateStatus(c.x)
}

// INY: Increment index register Y.
func (c *Cpu) INY(iop *Iop) {
	c.y++
	c.updateStatus(c.y)
}

// JMP: Jump.
func (c *Cpu) JMP(iop *Iop) {
	c.pc = c.memoryAddress(iop)
}

// JSR: Jump to subroutine.
func (c *Cpu) JSR(iop *Iop) {
	c.Bus.Write16(c.StackHead(-1), c.pc-1)
	c.sp -= 2
	c.pc = iop.op16
}

// LDA: Load accumulator from memory.
func (c *Cpu) LDA(iop *Iop) {
	c.ac = c.resolveOperand(iop)
	c.updateStatus(c.ac)
}

// LDX: Load index register X from memory.
func (c *Cpu) LDX(iop *Iop) {
	c.x = c.resolveOperand(iop)
	c.updateStatus(c.x)
}

// LDY: Load index register Y from memory.
func (c *Cpu) LDY(iop *Iop) {
	c.y = c.resolveOperand(iop)
	c.updateStatus(c.y)
}

// LSR: Logical shift memory or accumulator right.
func (c *Cpu) LSR(iop *Iop) {
	switch iop.ot.addressing {
	case accumulator:
		c.setStatus(sCarry, c.ac&1 == 1)
		c.ac >>= 1
		c.updateStatus(c.ac)
	default:
		address := c.memoryAddress(iop)
		value := c.Bus.Read(address)
		c.setStatus(sCarry, value&1 == 1)
		value >>= 1
		c.Bus.Write(address, value)
		c.updateStatus(value)
	}
}

// NOP: No operation.
func (c *Cpu) NOP(iop *Iop) {
}

// ORA: OR accumulator with memory.
func (c *Cpu) ORA(iop *Iop) {
	c.ac |= c.resolveOperand(iop)
	c.updateStatus(c.ac)
}

// PHA: Push accumulator onto stack.
func (c *Cpu) PHA(iop *Iop) {
	c.Bus.Write(0x0100+address(c.sp), c.ac)
	c.sp--
}

// PLA: Pull accumulator from stack.
func (c *Cpu) PLA(iop *Iop) {
	c.sp++
	c.ac = c.Bus.Read(0x0100 + address(c.sp))
}

// ROL: Rotate memory or accumulator left one bit.
func (c *Cpu) ROL(iop *Iop) {
	carry := c.getStatusInt(sCarry)
	switch iop.ot.addressing {
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

// RTS: Return from subroutine.
func (c *Cpu) RTS(iop *Iop) {
	c.pc = c.Bus.Read16(c.StackHead(1))
	c.sp += 2
	c.pc += 1
}

// SBC: Subtract memory with borrow from accumulator.
func (c *Cpu) SBC(iop *Iop) {
	valueSigned := int16(c.ac) - int16(c.resolveOperand(iop))
	if !c.getStatus(sCarry) {
		valueSigned--
	}
	c.setStatus(sCarry, valueSigned < 0)
	c.ac = uint8(valueSigned)
}

// SEI: Set interrupt-disable flag.
func (c *Cpu) SEI(iop *Iop) {
	c.setStatus(sInterrupt, false)
}

// STA: Store accumulator to memory.
func (c *Cpu) STA(iop *Iop) {
	c.Bus.Write(c.memoryAddress(iop), c.ac)
}

// STX: Store index register X to memory.
func (c *Cpu) STX(iop *Iop) {
	c.Bus.Write(c.memoryAddress(iop), c.x)
}

// STY: Store index register Y to memory.
func (c *Cpu) STY(iop *Iop) {
	c.Bus.Write(c.memoryAddress(iop), c.y)
}

// TAX: Transfer accumulator to index register X.
func (c *Cpu) TAX(iop *Iop) {
	c.x = c.ac
	c.updateStatus(c.x)
}

// TAY: Transfer accumulator to index register Y.
func (c *Cpu) TAY(iop *Iop) {
	c.y = c.ac
	c.updateStatus(c.y)
}

// TXA: Transfer index register X to accumulator.
func (c *Cpu) TXA(iop *Iop) {
	c.ac = c.x
	c.updateStatus(c.ac)
}

// TXS: Transfer index register X to stack pointer.
func (c *Cpu) TXS(iop *Iop) {
	c.sp = c.x
	c.updateStatus(c.sp)
}

// TYA: Transfer index register Y to accumulator.
func (c *Cpu) TYA(iop *Iop) {
	c.ac = c.y
	c.updateStatus(c.ac)
}

// _END: Custom go6502 instruction.
// Exit, with contents of X register as exit status.
func (c *Cpu) _END(iop *Iop) {
	c.ExitChan <- int(c.x)
}
