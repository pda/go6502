/*
	Package cpu implements the MOS 6502 processor.

	cpu.Cpu requires a bus.Bus to read/write 8-bit data to 16-bit addresses.

	cpu.Cpu also provides a monitor hook, allowing external code to observe
	and block on instructions before they're executed.
*/
package cpu

import (
	"fmt"
	"strings"

	"github.com/pda/go6502/bus"
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

// StackBase is the base address of the stack, which begins at StackBase+0xFF
// and grows downwards towards this address.
const StackBase = 0x0100

// Cpu represents the internal state of the CPU.
type Cpu struct {

	// Program counter.
	PC uint16

	// Accumulator register.
	AC byte

	// X general purpose / index register.
	X byte

	// Y general purpose / index register.
	Y byte

	// Stack pointer (low byte of 0x0100..0x01FF).
	SP byte

	// Status register; carry, zero, interrupt, bcd, brk, _, overflow, sign.
	SR byte

	// Bus is the system address bus, mapping 64K of address space to
	// different back-end devices.
	Bus *bus.Bus

	monitor  Monitor
	ExitChan chan int
}

// A Monitor is a blocking observer of instruction execution.
type Monitor interface {
	BeforeExecute(Instruction)
	Shutdown()
}

// AttachMonitor sets the given Monitor to observe instructions before they
// execute, in a blocking manner. This allows for logging, analysis, and
// interactive debugging.
func (c *Cpu) AttachMonitor(m Monitor) {
	c.monitor = m
}

// Shutdown tells the CPU to shut-down, and to pass the message on
// to subordinates such as the address bus.
func (c *Cpu) Shutdown() {
	c.Bus.Shutdown()
	if c.monitor != nil {
		c.monitor.Shutdown()
	}
}

// Reset the CPU, emulating triggering the RESB line.
// From 65C02 manual: All Registers are initialized by software except the
// Decimal and Interrupt disable mode select bits of the Processor Status
// Register (P) are initialized by hardware. ... The program counter is loaded
// with the reset vector from locations FFFC (low byte) and FFFD (high byte).
func (c *Cpu) Reset() {
	c.PC = c.Bus.Read16(0xFFFC)
	c.SR = 0x34 // Manual says xx1101xx, this sets 00110100.
}

func (c *Cpu) Step() {
	in := ReadInstruction(c.PC, c.Bus)
	if c.monitor != nil {
		c.monitor.BeforeExecute(in)
	}
	c.PC += uint16(in.Bytes)
	c.execute(in)
}

func (c *Cpu) String() string {
	return fmt.Sprintf(
		"CPU PC:0x%04X AC:0x%02X X:0x%02X Y:0x%02X SP:0x%02X SR:%s",
		c.PC, c.AC, c.X, c.Y, c.SP,
		c.statusString(),
	)
}

func (c *Cpu) stackHead(offset int8) uint16 {
	return uint16(StackBase) + uint16(c.SP) + uint16(offset)
}

func (c *Cpu) resolveOperand(in Instruction) uint8 {
	switch in.addressing {
	case immediate:
		return in.op8
	default:
		return c.Bus.Read(c.memoryAddress(in))
	}
}

func (c *Cpu) memoryAddress(in Instruction) uint16 {
	switch in.addressing {
	case absolute:
		return in.op16
	case absoluteX:
		return in.op16 + uint16(c.X)
	case absoluteY:
		return in.op16 + uint16(c.Y)

	// Indexed Indirect (X)
	// Operand is the zero-page location of a little-endian 16-bit base address.
	// The X register is added (wrapping; discarding overflow) before loading.
	// The resulting address loaded from (base+X) becomes the effective operand.
	// (base + X) must be in zero-page.
	case indirectX:
		location := uint16(in.op8 + c.X)
		if location == 0xFF {
			panic("Indexed indirect high-byte not on zero page.")
		}
		return c.Bus.Read16(location)

	// Indirect Indexed (Y)
	// Operand is the zero-page location of a little-endian 16-bit address.
	// The address is loaded, and then the Y register is added to it.
	// The resulting loaded_address + Y becomes the effective operand.
	case indirectY:
		return c.Bus.Read16(uint16(in.op8)) + uint16(c.Y)

	case zeropage:
		return uint16(in.op8)
	case zeropageX:
		return uint16(in.op8 + c.X)
	case zeropageY:
		return uint16(in.op8 + c.Y)
	default:
		panic("unhandled addressing")
	}
}

func (c *Cpu) getStatus(bit uint8) bool {
	return c.getStatusInt(bit) == 1
}

func (c *Cpu) getStatusInt(bit uint8) uint8 {
	return (c.SR >> bit) & 1
}

func (c *Cpu) setStatus(bit uint8, state bool) {
	if state {
		c.SR |= 1 << bit
	} else {
		c.SR &^= 1 << bit
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

func (c *Cpu) branch(in Instruction) {
	relative := int8(in.op8) // signed
	if relative >= 0 {
		c.PC += uint16(relative)
	} else {
		c.PC -= uint16(-relative)
	}
}

func (c *Cpu) execute(in Instruction) {
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
	case brk:
		c.BRK(in)
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
	case sec:
		c.SEC(in)
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
func (c *Cpu) ADC(in Instruction) {
	value16 := uint16(c.AC) + uint16(c.resolveOperand(in)) + uint16(c.getStatusInt(sCarry))
	c.setStatus(sCarry, value16 > 0xFF)
	c.AC = uint8(value16)
	c.updateStatus(c.AC)
}

// AND: And accumulator with memory.
func (c *Cpu) AND(in Instruction) {
	c.AC &= c.resolveOperand(in)
	c.updateStatus(c.AC)
}

// ASL: Shift memory or accumulator left one bit.
func (c *Cpu) ASL(in Instruction) {
	switch in.addressing {
	case accumulator:
		c.setStatus(sCarry, (c.AC>>7) == 1) // carry = old bit 7
		c.AC <<= 1
		c.updateStatus(c.AC)
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
func (c *Cpu) BCC(in Instruction) {
	if !c.getStatus(sCarry) {
		c.branch(in)
	}
}

// BCS: Branch if carry set.
func (c *Cpu) BCS(in Instruction) {
	if c.getStatus(sCarry) {
		c.branch(in)
	}
}

// BEQ: Branch if equal (z=1).
func (c *Cpu) BEQ(in Instruction) {
	if c.getStatus(sZero) {
		c.branch(in)
	}
}

// BMI: Branch if negative.
func (c *Cpu) BMI(in Instruction) {
	if c.getStatus(sNegative) {
		c.branch(in)
	}
}

// BNE: Branch if not equal.
func (c *Cpu) BNE(in Instruction) {
	if !c.getStatus(sZero) {
		c.branch(in)
	}
}

// BPL: Branch if positive.
func (c *Cpu) BPL(in Instruction) {
	if !c.getStatus(sNegative) {
		c.branch(in)
	}
}

// BRK: software interrupt
func (c *Cpu) BRK(in Instruction) {
	// temporarily used to dump status
	fmt.Println("BRK:", c)
}

// CLC: Clear carry flag.
func (c *Cpu) CLC(in Instruction) {
	c.setStatus(sCarry, false)
}

// CLD: Clear decimal mode flag.
func (c *Cpu) CLD(in Instruction) {
	c.setStatus(sDecimal, false)
}

// CLI: Clear interrupt-disable flag.
func (c *Cpu) CLI(in Instruction) {
	c.setStatus(sInterrupt, true)
}

// CMP: Compare accumulator with memory.
func (c *Cpu) CMP(in Instruction) {
	value := c.resolveOperand(in)
	c.setStatus(sCarry, c.AC >= value)
	c.updateStatus(c.AC - value)
}

// CPX: Compare index register X with memory.
func (c *Cpu) CPX(in Instruction) {
	value := c.resolveOperand(in)
	c.setStatus(sCarry, c.X >= value)
	c.updateStatus(c.X - value)
}

// CPY: Compare index register Y with memory.
func (c *Cpu) CPY(in Instruction) {
	value := c.resolveOperand(in)
	c.setStatus(sCarry, c.Y >= value)
	c.updateStatus(c.Y - value)
}

// DEC: Decrement.
func (c *Cpu) DEC(in Instruction) {
	address := c.memoryAddress(in)
	value := c.Bus.Read(address) - 1
	c.Bus.Write(address, value)
	c.updateStatus(value)
}

// DEX: Decrement index register X.
func (c *Cpu) DEX(in Instruction) {
	c.X--
	c.updateStatus(c.X)
}

// DEY: Decrement index register Y.
func (c *Cpu) DEY(in Instruction) {
	c.Y--
	c.updateStatus(c.Y)
}

// EOR: Exclusive-OR accumulator with memory.
func (c *Cpu) EOR(in Instruction) {
	value := c.resolveOperand(in)
	c.AC ^= value
	c.updateStatus(c.AC)
}

// INC: Increment.
func (c *Cpu) INC(in Instruction) {
	address := c.memoryAddress(in)
	value := c.Bus.Read(address) + 1
	c.Bus.Write(address, value)
	c.updateStatus(value)
}

// INX: Increment index register X.
func (c *Cpu) INX(in Instruction) {
	c.X++
	c.updateStatus(c.X)
}

// INY: Increment index register Y.
func (c *Cpu) INY(in Instruction) {
	c.Y++
	c.updateStatus(c.Y)
}

// JMP: Jump.
func (c *Cpu) JMP(in Instruction) {
	c.PC = c.memoryAddress(in)
}

// JSR: Jump to subroutine.
func (c *Cpu) JSR(in Instruction) {
	c.Bus.Write16(c.stackHead(-1), c.PC-1)
	c.SP -= 2
	c.PC = in.op16
}

// LDA: Load accumulator from memory.
func (c *Cpu) LDA(in Instruction) {
	c.AC = c.resolveOperand(in)
	c.updateStatus(c.AC)
}

// LDX: Load index register X from memory.
func (c *Cpu) LDX(in Instruction) {
	c.X = c.resolveOperand(in)
	c.updateStatus(c.X)
}

// LDY: Load index register Y from memory.
func (c *Cpu) LDY(in Instruction) {
	c.Y = c.resolveOperand(in)
	c.updateStatus(c.Y)
}

// LSR: Logical shift memory or accumulator right.
func (c *Cpu) LSR(in Instruction) {
	switch in.addressing {
	case accumulator:
		c.setStatus(sCarry, c.AC&1 == 1)
		c.AC >>= 1
		c.updateStatus(c.AC)
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
func (c *Cpu) NOP(in Instruction) {
}

// ORA: OR accumulator with memory.
func (c *Cpu) ORA(in Instruction) {
	c.AC |= c.resolveOperand(in)
	c.updateStatus(c.AC)
}

// PHA: Push accumulator onto stack.
func (c *Cpu) PHA(in Instruction) {
	c.Bus.Write(0x0100+uint16(c.SP), c.AC)
	c.SP--
}

// PLA: Pull accumulator from stack.
func (c *Cpu) PLA(in Instruction) {
	c.SP++
	c.AC = c.Bus.Read(0x0100 + uint16(c.SP))
}

// ROL: Rotate memory or accumulator left one bit.
func (c *Cpu) ROL(in Instruction) {
	carry := c.getStatusInt(sCarry)
	switch in.addressing {
	case accumulator:
		c.setStatus(sCarry, (c.AC>>7) == 1)
		c.AC = (c.AC << 1) | carry
		c.updateStatus(c.AC)
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
func (c *Cpu) RTS(in Instruction) {
	c.PC = c.Bus.Read16(c.stackHead(1))
	c.SP += 2
	c.PC += 1
}

// SBC: Subtract memory with borrow from accumulator.
func (c *Cpu) SBC(in Instruction) {
	valueSigned := int16(c.AC) - int16(c.resolveOperand(in))
	if !c.getStatus(sCarry) {
		valueSigned--
	}
	c.AC = uint8(valueSigned)

	// v: Set if signed overflow; cleared if valid sign result.
	// TODO: c.setStatus(sOverflow, something)

	// c: Set if unsigned borrow not required; cleared if unsigned borrow.
	c.setStatus(sCarry, valueSigned >= 0)

	// n: Set if most significant bit of result is set; else cleared.
	// z: Set if result is zero; else cleared.
	c.updateStatus(c.AC)
}

// SEC: Set carry flag.
func (c *Cpu) SEC(in Instruction) {
	c.setStatus(sCarry, true)
}

// SEI: Set interrupt-disable flag.
func (c *Cpu) SEI(in Instruction) {
	c.setStatus(sInterrupt, false)
}

// STA: Store accumulator to memory.
func (c *Cpu) STA(in Instruction) {
	c.Bus.Write(c.memoryAddress(in), c.AC)
}

// STX: Store index register X to memory.
func (c *Cpu) STX(in Instruction) {
	c.Bus.Write(c.memoryAddress(in), c.X)
}

// STY: Store index register Y to memory.
func (c *Cpu) STY(in Instruction) {
	c.Bus.Write(c.memoryAddress(in), c.Y)
}

// TAX: Transfer accumulator to index register X.
func (c *Cpu) TAX(in Instruction) {
	c.X = c.AC
	c.updateStatus(c.X)
}

// TAY: Transfer accumulator to index register Y.
func (c *Cpu) TAY(in Instruction) {
	c.Y = c.AC
	c.updateStatus(c.Y)
}

// TXA: Transfer index register X to accumulator.
func (c *Cpu) TXA(in Instruction) {
	c.AC = c.X
	c.updateStatus(c.AC)
}

// TXS: Transfer index register X to stack pointer.
func (c *Cpu) TXS(in Instruction) {
	c.SP = c.X
	c.updateStatus(c.SP)
}

// TYA: Transfer index register Y to accumulator.
func (c *Cpu) TYA(in Instruction) {
	c.AC = c.Y
	c.updateStatus(c.AC)
}

// _END: Custom go6502 instruction.
// Exit, with contents of X register as exit status.
func (c *Cpu) _END(in Instruction) {
	c.ExitChan <- int(c.X)
}
