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
	pc       Address
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
	in := ReadInstruction(c.pc, c.Bus)
	if c.debugger != nil {
		c.debugger.BeforeExecute(in)
	}
	c.pc += Address(in.bytes)
	c.Execute(in)
}

func (c *Cpu) String() string {
	return fmt.Sprintf(
		"CPU pc:0x%04X ac:0x%02X x:0x%02X y:0x%02X sp:0x%02X sr:%s",
		c.pc, c.ac, c.x, c.y, c.sp,
		c.statusString(),
	)
}

func (c *Cpu) StackHead(offset int8) Address {
	return Address(0x0100) + Address(c.sp) + Address(offset)
}

func (c *Cpu) resolveOperand(in *Instruction) uint8 {
	switch in.addressing {
	case immediate:
		return in.op8
	default:
		return c.Bus.Read(c.memoryAddress(in))
	}
}

func (c *Cpu) memoryAddress(in *Instruction) Address {
	switch in.addressing {
	case absolute:
		return in.op16
	case absoluteX:
		return in.op16 + Address(c.x)
	case absoluteY:
		return in.op16 + Address(c.y)

	// Indexed Indirect (X)
	// Operand is the zero-page location of a little-endian 16-bit base address.
	// The X register is added (wrapping; discarding overflow) before loading.
	// The resulting address loaded from (base+X) becomes the effective operand.
	// (base + X) must be in zero-page.
	case indirectX:
		location := Address(in.op8 + c.x)
		if location == 0xFF {
			panic("Indexed indirect high-byte not on zero page.")
		}
		return c.Bus.Read16(location)

	// Indirect Indexed (Y)
	// Operand is the zero-page location of a little-endian 16-bit address.
	// The address is loaded, and then the Y register is added to it.
	// The resulting loaded_address + Y becomes the effective operand.
	case indirectY:
		return c.Bus.Read16(Address(in.op8)) + Address(c.y)

	case zeropage:
		return Address(in.op8)
	case zeropageX:
		return Address(in.op8 + c.x)
	case zeropageY:
		return Address(in.op8 + c.y)
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

func (c *Cpu) branch(in *Instruction) {
	relative := int8(in.op8) // signed
	if relative >= 0 {
		c.pc += Address(relative)
	} else {
		c.pc -= Address(-relative)
	}
}

func (c *Cpu) Execute(in *Instruction) {
	switch in.id {
	case adc:
		c.ADC(in)
	case and:
		c.AND(in)
	case asl:
		c.ASL(in)
	case bcc:
		c.BCC(in)
	case bcs:
		c.BCS(in)
	case beq:
		c.BEQ(in)
	case bmi:
		c.BMI(in)
	case bne:
		c.BNE(in)
	case bpl:
		c.BPL(in)
	case clc:
		c.CLC(in)
	case cld:
		c.CLD(in)
	case cli:
		c.CLI(in)
	case cmp:
		c.CMP(in)
	case cpx:
		c.CPX(in)
	case cpy:
		c.CPY(in)
	case dec:
		c.DEC(in)
	case dex:
		c.DEX(in)
	case dey:
		c.DEY(in)
	case eor:
		c.EOR(in)
	case inc:
		c.INC(in)
	case inx:
		c.INX(in)
	case iny:
		c.INY(in)
	case jmp:
		c.JMP(in)
	case jsr:
		c.JSR(in)
	case lda:
		c.LDA(in)
	case ldx:
		c.LDX(in)
	case ldy:
		c.LDY(in)
	case lsr:
		c.LSR(in)
	case nop:
		c.NOP(in)
	case ora:
		c.ORA(in)
	case pha:
		c.PHA(in)
	case pla:
		c.PLA(in)
	case rol:
		c.ROL(in)
	case rts:
		c.RTS(in)
	case sbc:
		c.SBC(in)
	case sei:
		c.SEI(in)
	case sta:
		c.STA(in)
	case stx:
		c.STX(in)
	case sty:
		c.STY(in)
	case tax:
		c.TAX(in)
	case tay:
		c.TAY(in)
	case txa:
		c.TXA(in)
	case txs:
		c.TXS(in)
	case tya:
		c.TYA(in)
	case _end:
		c._END(in)
	default:
		panic(fmt.Sprintf("unhandled instruction: %v", in))
	}
}

// ADC: Add memory and carry to accumulator.
func (c *Cpu) ADC(in *Instruction) {
	value16 := uint16(c.ac) + uint16(c.resolveOperand(in)) + uint16(c.getStatusInt(sCarry))
	c.setStatus(sCarry, value16 > 0xFF)
	c.ac = uint8(value16)
	c.updateStatus(c.ac)
}

// AND: And accumulator with memory.
func (c *Cpu) AND(in *Instruction) {
	c.ac &= c.resolveOperand(in)
	c.updateStatus(c.ac)
}

// ASL: Shift memory or accumulator left one bit.
func (c *Cpu) ASL(in *Instruction) {
	switch in.addressing {
	case accumulator:
		c.setStatus(sCarry, (c.ac>>7) == 1) // carry = old bit 7
		c.ac <<= 1
		c.updateStatus(c.ac)
	default:
		address := c.memoryAddress(in)
		value := c.Bus.Read(address)
		c.setStatus(sCarry, (value>>7) == 1) // carry = old bit 7
		value <<= 1
		c.Bus.Write(address, value)
		c.updateStatus(value)
	}
}

// BCC: Branch if carry clear.
func (c *Cpu) BCC(in *Instruction) {
	if !c.getStatus(sCarry) {
		c.branch(in)
	}
}

// BCS: Branch if carry set.
func (c *Cpu) BCS(in *Instruction) {
	if c.getStatus(sCarry) {
		c.branch(in)
	}
}

// BEQ: Branch if equal (z=1).
func (c *Cpu) BEQ(in *Instruction) {
	if c.getStatus(sZero) {
		c.branch(in)
	}
}

// BMI: Branch if negative.
func (c *Cpu) BMI(in *Instruction) {
	if c.getStatus(sNegative) {
		c.branch(in)
	}
}

// BNE: Branch if not equal.
func (c *Cpu) BNE(in *Instruction) {
	if !c.getStatus(sZero) {
		c.branch(in)
	}
}

// BPL: Branch if positive.
func (c *Cpu) BPL(in *Instruction) {
	if !c.getStatus(sNegative) {
		c.branch(in)
	}
}

// CLC: Clear carry flag.
func (c *Cpu) CLC(in *Instruction) {
	c.setStatus(sCarry, false)
}

// CLD: Clear decimal mode flag.
func (c *Cpu) CLD(in *Instruction) {
	c.setStatus(sDecimal, false)
}

// CLI: Clear interrupt-disable flag.
func (c *Cpu) CLI(in *Instruction) {
	c.setStatus(sInterrupt, true)
}

// CMP: Compare accumulator with memory.
func (c *Cpu) CMP(in *Instruction) {
	value := c.resolveOperand(in)
	c.setStatus(sCarry, c.ac >= value)
	c.updateStatus(c.ac - value)
}

// CPX: Compare index register X with memory.
func (c *Cpu) CPX(in *Instruction) {
	value := c.resolveOperand(in)
	c.setStatus(sCarry, c.x >= value)
	c.updateStatus(c.x - value)
}

// CPY: Compare index register Y with memory.
func (c *Cpu) CPY(in *Instruction) {
	value := c.resolveOperand(in)
	c.setStatus(sCarry, c.y >= value)
	c.updateStatus(c.y - value)
}

// DEC: Decrement.
func (c *Cpu) DEC(in *Instruction) {
	address := c.memoryAddress(in)
	value := c.Bus.Read(address) - 1
	c.Bus.Write(address, value)
	c.updateStatus(value)
}

// DEX: Decrement index register X.
func (c *Cpu) DEX(in *Instruction) {
	c.x--
	c.updateStatus(c.x)
}

// DEY: Decrement index register Y.
func (c *Cpu) DEY(in *Instruction) {
	c.y--
	c.updateStatus(c.y)
}

// EOR: Exclusive-OR accumulator with memory.
func (c *Cpu) EOR(in *Instruction) {
	value := c.resolveOperand(in)
	c.ac ^= value
	c.updateStatus(c.ac)
}

// INC: Increment.
func (c *Cpu) INC(in *Instruction) {
	address := c.memoryAddress(in)
	value := c.Bus.Read(address) + 1
	c.Bus.Write(address, value)
	c.updateStatus(value)
}

// INX: Increment index register X.
func (c *Cpu) INX(in *Instruction) {
	c.x++
	c.updateStatus(c.x)
}

// INY: Increment index register Y.
func (c *Cpu) INY(in *Instruction) {
	c.y++
	c.updateStatus(c.y)
}

// JMP: Jump.
func (c *Cpu) JMP(in *Instruction) {
	c.pc = c.memoryAddress(in)
}

// JSR: Jump to subroutine.
func (c *Cpu) JSR(in *Instruction) {
	c.Bus.Write16(c.StackHead(-1), c.pc-1)
	c.sp -= 2
	c.pc = in.op16
}

// LDA: Load accumulator from memory.
func (c *Cpu) LDA(in *Instruction) {
	c.ac = c.resolveOperand(in)
	c.updateStatus(c.ac)
}

// LDX: Load index register X from memory.
func (c *Cpu) LDX(in *Instruction) {
	c.x = c.resolveOperand(in)
	c.updateStatus(c.x)
}

// LDY: Load index register Y from memory.
func (c *Cpu) LDY(in *Instruction) {
	c.y = c.resolveOperand(in)
	c.updateStatus(c.y)
}

// LSR: Logical shift memory or accumulator right.
func (c *Cpu) LSR(in *Instruction) {
	switch in.addressing {
	case accumulator:
		c.setStatus(sCarry, c.ac&1 == 1)
		c.ac >>= 1
		c.updateStatus(c.ac)
	default:
		address := c.memoryAddress(in)
		value := c.Bus.Read(address)
		c.setStatus(sCarry, value&1 == 1)
		value >>= 1
		c.Bus.Write(address, value)
		c.updateStatus(value)
	}
}

// NOP: No operation.
func (c *Cpu) NOP(in *Instruction) {
}

// ORA: OR accumulator with memory.
func (c *Cpu) ORA(in *Instruction) {
	c.ac |= c.resolveOperand(in)
	c.updateStatus(c.ac)
}

// PHA: Push accumulator onto stack.
func (c *Cpu) PHA(in *Instruction) {
	c.Bus.Write(0x0100+Address(c.sp), c.ac)
	c.sp--
}

// PLA: Pull accumulator from stack.
func (c *Cpu) PLA(in *Instruction) {
	c.sp++
	c.ac = c.Bus.Read(0x0100 + Address(c.sp))
}

// ROL: Rotate memory or accumulator left one bit.
func (c *Cpu) ROL(in *Instruction) {
	carry := c.getStatusInt(sCarry)
	switch in.addressing {
	case accumulator:
		c.setStatus(sCarry, (c.ac>>7) == 1)
		c.ac = (c.ac << 1) | carry
		c.updateStatus(c.ac)
	default:
		address := c.memoryAddress(in)
		value := c.Bus.Read(address)
		c.setStatus(sCarry, (value>>7) == 1)
		value = (value << 1) | carry
		c.Bus.Write(address, value)
		c.updateStatus(value)
	}
}

// RTS: Return from subroutine.
func (c *Cpu) RTS(in *Instruction) {
	c.pc = c.Bus.Read16(c.StackHead(1))
	c.sp += 2
	c.pc += 1
}

// SBC: Subtract memory with borrow from accumulator.
func (c *Cpu) SBC(in *Instruction) {
	valueSigned := int16(c.ac) - int16(c.resolveOperand(in))
	if !c.getStatus(sCarry) {
		valueSigned--
	}
	c.setStatus(sCarry, valueSigned < 0)
	c.ac = uint8(valueSigned)
}

// SEI: Set interrupt-disable flag.
func (c *Cpu) SEI(in *Instruction) {
	c.setStatus(sInterrupt, false)
}

// STA: Store accumulator to memory.
func (c *Cpu) STA(in *Instruction) {
	c.Bus.Write(c.memoryAddress(in), c.ac)
}

// STX: Store index register X to memory.
func (c *Cpu) STX(in *Instruction) {
	c.Bus.Write(c.memoryAddress(in), c.x)
}

// STY: Store index register Y to memory.
func (c *Cpu) STY(in *Instruction) {
	c.Bus.Write(c.memoryAddress(in), c.y)
}

// TAX: Transfer accumulator to index register X.
func (c *Cpu) TAX(in *Instruction) {
	c.x = c.ac
	c.updateStatus(c.x)
}

// TAY: Transfer accumulator to index register Y.
func (c *Cpu) TAY(in *Instruction) {
	c.y = c.ac
	c.updateStatus(c.y)
}

// TXA: Transfer index register X to accumulator.
func (c *Cpu) TXA(in *Instruction) {
	c.ac = c.x
	c.updateStatus(c.ac)
}

// TXS: Transfer index register X to stack pointer.
func (c *Cpu) TXS(in *Instruction) {
	c.sp = c.x
	c.updateStatus(c.sp)
}

// TYA: Transfer index register Y to accumulator.
func (c *Cpu) TYA(in *Instruction) {
	c.ac = c.y
	c.updateStatus(c.ac)
}

// _END: Custom go6502 instruction.
// Exit, with contents of X register as exit status.
func (c *Cpu) _END(in *Instruction) {
	c.ExitChan <- int(c.x)
}
