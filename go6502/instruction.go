package go6502

import (
	"fmt"
)

// addressing modes
const (
	_ = iota
	absolute
	absoluteX
	absoluteY
	accumulator
	immediate
	implied
	indirect
	indirectX
	indirectY
	relative
	zeropage
	zeropageX
	zeropageY
)

var addressingNames = [...]string{
	"",
	"absolute",
	"absoluteX",
	"absoluteY",
	"accumulator",
	"immediate",
	"implied",
	"(indirect)",
	"(indirect,X)",
	"(indirect),Y",
	"relative",
	"zeropage",
	"zeropageX",
	"zeropageY",
}

// adc..tya represent the 6502 instruction set mnemonics. Each mnemonic maps to
// a number of different opcodes, depending on the addressing mode.
const (
	_ = iota
	adc
	and
	asl
	bcc
	bcs
	beq
	bit
	bmi
	bne
	bpl
	brk
	bvc
	bvs
	clc
	cld
	cli
	clv
	cmp
	cpx
	cpy
	dec
	dex
	dey
	eor
	inc
	inx
	iny
	jmp
	jsr
	lda
	ldx
	ldy
	lsr
	nop
	ora
	pha
	php
	pla
	plp
	rol
	ror
	rti
	rts
	sbc
	sec
	sed
	sei
	sta
	stx
	sty
	tax
	tay
	tsx
	txa
	txs
	tya
	_end
)

// Instruction represents a 6502 op-code instruction, including the addressing
// mode encoded into the op-code, but not the operand value following the
// opcode in memory.
type Instruction struct {
	id         uint8 // the const identifier of the instruction type, e.g. ADC
	opcode     byte
	addressing int
	bytes      int
	cycles     int
	flags      int
}

// Iop is an instruction with its operand.
// One or both of the operand types will be zero.
// This is determined by (in.bytes - 1) / 8
type Iop struct {
	in   *Instruction
	op8  uint8
	op16 address
}

func (iop *Iop) String() (s string) {
	switch iop.in.bytes {
	case 3:
		s = fmt.Sprintf("%v $%04X", iop.in, iop.op16)
	case 2:
		s = fmt.Sprintf("%v $%02X", iop.in, iop.op8)
	case 1:
		s = iop.in.String()
	}
	return
}

func findInstruction(opcode byte) *Instruction {
	// TODO: singleton instructions; they're immutable.
	var i Instruction
	switch opcode {
	default:
		panic(fmt.Sprintf("Unknown opcode: 0x%02X", opcode))
	case 0x69:
		i = Instruction{adc, opcode, immediate, 2, 2, 0}
	case 0x65:
		i = Instruction{adc, opcode, zeropage, 2, 3, 0}
	case 0x75:
		i = Instruction{adc, opcode, zeropageX, 2, 4, 0}
	case 0x6D:
		i = Instruction{adc, opcode, absolute, 3, 4, 0}
	case 0x7D:
		i = Instruction{adc, opcode, absoluteX, 3, 4, 0}
	case 0x79:
		i = Instruction{adc, opcode, absoluteY, 3, 4, 0}
	case 0x61:
		i = Instruction{adc, opcode, indirectX, 2, 6, 0}
	case 0x71:
		i = Instruction{adc, opcode, indirectY, 2, 5, 0}
	case 0x29:
		i = Instruction{and, opcode, immediate, 2, 2, 0}
	case 0x25:
		i = Instruction{and, opcode, zeropage, 2, 3, 0}
	case 0x35:
		i = Instruction{and, opcode, zeropageX, 2, 4, 0}
	case 0x2D:
		i = Instruction{and, opcode, absolute, 3, 4, 0}
	case 0x3D:
		i = Instruction{and, opcode, absoluteX, 3, 4, 0}
	case 0x39:
		i = Instruction{and, opcode, absoluteY, 3, 4, 0}
	case 0x21:
		i = Instruction{and, opcode, indirectX, 2, 6, 0}
	case 0x31:
		i = Instruction{and, opcode, indirectY, 2, 5, 0}
	case 0x0A:
		i = Instruction{asl, opcode, accumulator, 1, 2, 0}
	case 0x06:
		i = Instruction{asl, opcode, zeropage, 2, 5, 0}
	case 0x16:
		i = Instruction{asl, opcode, zeropageX, 2, 6, 0}
	case 0x0E:
		i = Instruction{asl, opcode, absolute, 3, 6, 0}
	case 0x1E:
		i = Instruction{asl, opcode, absoluteX, 3, 7, 0}
	case 0x90:
		i = Instruction{bcc, opcode, relative, 2, 2, 0}
	case 0xB0:
		i = Instruction{bcs, opcode, relative, 2, 2, 0}
	case 0xF0:
		i = Instruction{beq, opcode, relative, 2, 2, 0}
	case 0x24:
		i = Instruction{bit, opcode, zeropage, 2, 3, 0}
	case 0x2C:
		i = Instruction{bit, opcode, absolute, 3, 4, 0}
	case 0x30:
		i = Instruction{bmi, opcode, relative, 2, 2, 0}
	case 0xD0:
		i = Instruction{bne, opcode, relative, 2, 2, 0}
	case 0x10:
		i = Instruction{bpl, opcode, relative, 2, 2, 0}
	case 0x00:
		i = Instruction{brk, opcode, implied, 1, 7, 0}
	case 0x50:
		i = Instruction{bvc, opcode, relative, 2, 2, 0}
	case 0x70:
		i = Instruction{bvs, opcode, relative, 2, 2, 0}
	case 0x18:
		i = Instruction{clc, opcode, implied, 1, 2, 0}
	case 0xD8:
		i = Instruction{cld, opcode, implied, 1, 2, 0}
	case 0x58:
		i = Instruction{cli, opcode, implied, 1, 2, 0}
	case 0xB8:
		i = Instruction{clv, opcode, implied, 1, 2, 0}
	case 0xC9:
		i = Instruction{cmp, opcode, immediate, 2, 2, 0}
	case 0xC5:
		i = Instruction{cmp, opcode, zeropage, 2, 3, 0}
	case 0xD5:
		i = Instruction{cmp, opcode, zeropageX, 2, 4, 0}
	case 0xCD:
		i = Instruction{cmp, opcode, absolute, 3, 4, 0}
	case 0xDD:
		i = Instruction{cmp, opcode, absoluteX, 3, 4, 0}
	case 0xD9:
		i = Instruction{cmp, opcode, absoluteY, 3, 4, 0}
	case 0xC1:
		i = Instruction{cmp, opcode, indirectX, 2, 6, 0}
	case 0xD1:
		i = Instruction{cmp, opcode, indirectY, 2, 5, 0}
	case 0xE0:
		i = Instruction{cpx, opcode, immediate, 2, 2, 0}
	case 0xE4:
		i = Instruction{cpx, opcode, zeropage, 2, 3, 0}
	case 0xEC:
		i = Instruction{cpx, opcode, absolute, 3, 4, 0}
	case 0xC0:
		i = Instruction{cpy, opcode, immediate, 2, 2, 0}
	case 0xC4:
		i = Instruction{cpy, opcode, zeropage, 2, 3, 0}
	case 0xCC:
		i = Instruction{cpy, opcode, absolute, 3, 4, 0}
	case 0xC6:
		i = Instruction{dec, opcode, zeropage, 2, 5, 0}
	case 0xD6:
		i = Instruction{dec, opcode, zeropageX, 2, 6, 0}
	case 0xCE:
		i = Instruction{dec, opcode, absolute, 3, 3, 0}
	case 0xDE:
		i = Instruction{dec, opcode, absoluteX, 3, 7, 0}
	case 0xCA:
		i = Instruction{dex, opcode, implied, 1, 2, 0}
	case 0x88:
		i = Instruction{dey, opcode, implied, 1, 2, 0}
	case 0x49:
		i = Instruction{eor, opcode, immediate, 2, 2, 0}
	case 0x45:
		i = Instruction{eor, opcode, zeropage, 2, 3, 0}
	case 0x55:
		i = Instruction{eor, opcode, zeropageX, 2, 4, 0}
	case 0x4D:
		i = Instruction{eor, opcode, absolute, 3, 4, 0}
	case 0x5D:
		i = Instruction{eor, opcode, absoluteX, 3, 4, 0}
	case 0x59:
		i = Instruction{eor, opcode, absoluteY, 3, 4, 0}
	case 0x41:
		i = Instruction{eor, opcode, indirectX, 2, 6, 0}
	case 0x51:
		i = Instruction{eor, opcode, indirectY, 2, 5, 0}
	case 0xE6:
		i = Instruction{inc, opcode, zeropage, 2, 5, 0}
	case 0xF6:
		i = Instruction{inc, opcode, zeropageX, 2, 6, 0}
	case 0xEE:
		i = Instruction{inc, opcode, absolute, 3, 6, 0}
	case 0xFE:
		i = Instruction{inc, opcode, absoluteX, 3, 7, 0}
	case 0xE8:
		i = Instruction{inx, opcode, implied, 1, 2, 0}
	case 0xC8:
		i = Instruction{iny, opcode, implied, 1, 2, 0}
	case 0x4C:
		i = Instruction{jmp, opcode, absolute, 3, 3, 0}
	case 0x6C:
		i = Instruction{jmp, opcode, indirect, 3, 5, 0}
	case 0x20:
		i = Instruction{jsr, opcode, absolute, 3, 6, 0}
	case 0xA9:
		i = Instruction{lda, opcode, immediate, 2, 2, 0}
	case 0xA5:
		i = Instruction{lda, opcode, zeropage, 2, 3, 0}
	case 0xB5:
		i = Instruction{lda, opcode, zeropageX, 2, 4, 0}
	case 0xAD:
		i = Instruction{lda, opcode, absolute, 3, 4, 0}
	case 0xBD:
		i = Instruction{lda, opcode, absoluteX, 3, 4, 0}
	case 0xB9:
		i = Instruction{lda, opcode, absoluteY, 3, 4, 0}
	case 0xA1:
		i = Instruction{lda, opcode, indirectX, 2, 6, 0}
	case 0xB1:
		i = Instruction{lda, opcode, indirectY, 2, 5, 0}
	case 0xA2:
		i = Instruction{ldx, opcode, immediate, 2, 2, 0}
	case 0xA6:
		i = Instruction{ldx, opcode, zeropage, 2, 3, 0}
	case 0xB6:
		i = Instruction{ldx, opcode, zeropageY, 2, 4, 0}
	case 0xAE:
		i = Instruction{ldx, opcode, absolute, 3, 4, 0}
	case 0xBE:
		i = Instruction{ldx, opcode, absoluteY, 3, 4, 0}
	case 0xA0:
		i = Instruction{ldy, opcode, immediate, 2, 2, 0}
	case 0xA4:
		i = Instruction{ldy, opcode, zeropage, 2, 3, 0}
	case 0xB4:
		i = Instruction{ldy, opcode, zeropageX, 2, 4, 0}
	case 0xAC:
		i = Instruction{ldy, opcode, absolute, 3, 4, 0}
	case 0xBC:
		i = Instruction{ldy, opcode, absoluteX, 3, 4, 0}
	case 0x4A:
		i = Instruction{lsr, opcode, accumulator, 1, 2, 0}
	case 0x46:
		i = Instruction{lsr, opcode, zeropage, 2, 5, 0}
	case 0x56:
		i = Instruction{lsr, opcode, zeropageX, 2, 6, 0}
	case 0x4E:
		i = Instruction{lsr, opcode, absolute, 3, 6, 0}
	case 0x5E:
		i = Instruction{lsr, opcode, absoluteX, 3, 7, 0}
	case 0xEA:
		i = Instruction{nop, opcode, implied, 1, 2, 0}
	case 0x09:
		i = Instruction{ora, opcode, immediate, 2, 2, 0}
	case 0x05:
		i = Instruction{ora, opcode, zeropage, 2, 3, 0}
	case 0x15:
		i = Instruction{ora, opcode, zeropageX, 2, 4, 0}
	case 0x0D:
		i = Instruction{ora, opcode, absolute, 3, 4, 0}
	case 0x1D:
		i = Instruction{ora, opcode, absoluteX, 3, 4, 0}
	case 0x19:
		i = Instruction{ora, opcode, absoluteY, 3, 4, 0}
	case 0x01:
		i = Instruction{ora, opcode, indirectX, 2, 6, 0}
	case 0x11:
		i = Instruction{ora, opcode, indirectY, 2, 5, 0}
	case 0x48:
		i = Instruction{pha, opcode, implied, 1, 3, 0}
	case 0x08:
		i = Instruction{php, opcode, implied, 1, 3, 0}
	case 0x68:
		i = Instruction{pla, opcode, implied, 1, 4, 0}
	case 0x28:
		i = Instruction{php, opcode, implied, 1, 4, 0}
	case 0x2A:
		i = Instruction{rol, opcode, accumulator, 1, 2, 0}
	case 0x26:
		i = Instruction{rol, opcode, zeropage, 2, 5, 0}
	case 0x36:
		i = Instruction{rol, opcode, zeropageX, 2, 6, 0}
	case 0x2E:
		i = Instruction{rol, opcode, absolute, 3, 6, 0}
	case 0x3E:
		i = Instruction{rol, opcode, absoluteX, 3, 7, 0}
	case 0x6A:
		i = Instruction{ror, opcode, accumulator, 1, 2, 0}
	case 0x66:
		i = Instruction{ror, opcode, zeropage, 2, 5, 0}
	case 0x76:
		i = Instruction{ror, opcode, zeropageX, 2, 6, 0}
	case 0x6E:
		i = Instruction{ror, opcode, absolute, 3, 6, 0}
	case 0x7E:
		i = Instruction{ror, opcode, absoluteX, 3, 7, 0}
	case 0x40:
		i = Instruction{rti, opcode, implied, 1, 6, 0}
	case 0x60:
		i = Instruction{rts, opcode, implied, 1, 6, 0}
	case 0xE9:
		i = Instruction{sbc, opcode, immediate, 2, 2, 0}
	case 0xE5:
		i = Instruction{sbc, opcode, zeropage, 2, 3, 0}
	case 0xF5:
		i = Instruction{sbc, opcode, zeropageX, 2, 4, 0}
	case 0xED:
		i = Instruction{sbc, opcode, absolute, 3, 4, 0}
	case 0xFD:
		i = Instruction{sbc, opcode, absoluteX, 3, 4, 0}
	case 0xF9:
		i = Instruction{sbc, opcode, absoluteY, 3, 4, 0}
	case 0xE1:
		i = Instruction{sbc, opcode, indirectX, 2, 6, 0}
	case 0xF1:
		i = Instruction{sbc, opcode, indirectY, 2, 5, 0}
	case 0x38:
		i = Instruction{sec, opcode, implied, 1, 2, 0}
	case 0xF8:
		i = Instruction{sed, opcode, implied, 1, 2, 0}
	case 0x78:
		i = Instruction{sei, opcode, implied, 1, 2, 0}
	case 0x85:
		i = Instruction{sta, opcode, zeropage, 2, 3, 0}
	case 0x95:
		i = Instruction{sta, opcode, zeropageX, 2, 4, 0}
	case 0x8D:
		i = Instruction{sta, opcode, absolute, 3, 4, 0}
	case 0x9D:
		i = Instruction{sta, opcode, absoluteX, 3, 5, 0}
	case 0x99:
		i = Instruction{sta, opcode, absoluteY, 3, 5, 0}
	case 0x81:
		i = Instruction{sta, opcode, indirectX, 2, 6, 0}
	case 0x91:
		i = Instruction{sta, opcode, indirectY, 2, 6, 0}
	case 0x86:
		i = Instruction{stx, opcode, zeropage, 2, 3, 0}
	case 0x96:
		i = Instruction{stx, opcode, zeropageY, 2, 4, 0}
	case 0x8E:
		i = Instruction{stx, opcode, absolute, 3, 4, 0}
	case 0x84:
		i = Instruction{sty, opcode, zeropage, 2, 3, 0}
	case 0x94:
		i = Instruction{sty, opcode, zeropageX, 2, 4, 0}
	case 0x8C:
		i = Instruction{sty, opcode, absolute, 3, 4, 0}
	case 0xAA:
		i = Instruction{tax, opcode, implied, 1, 2, 0}
	case 0xA8:
		i = Instruction{tay, opcode, implied, 1, 2, 0}
	case 0xBA:
		i = Instruction{tsx, opcode, implied, 1, 2, 0}
	case 0x8A:
		i = Instruction{txa, opcode, implied, 1, 2, 0}
	case 0x9A:
		i = Instruction{txs, opcode, implied, 1, 2, 0}
	case 0x98:
		i = Instruction{tya, opcode, implied, 1, 2, 0}
	case 0xFF:
		i = Instruction{_end, opcode, implied, 1, 1, 0}
	}
	return &i
}

func (i *Instruction) String() string {
	return fmt.Sprintf("%s %s", i.name(), addressingNames[i.addressing])
}

func (i *Instruction) name() (s string) {
	switch i.id {
	case _end:
		s = "_END"
	case adc:
		s = "ADC"
	case and:
		s = "AND"
	case asl:
		s = "ASL"
	case bcc:
		s = "BCC"
	case bcs:
		s = "BCS"
	case beq:
		s = "BEQ"
	case bit:
		s = "BIT"
	case bmi:
		s = "BMI"
	case bne:
		s = "BNE"
	case bpl:
		s = "BPL"
	case brk:
		s = "BRK"
	case bvc:
		s = "BVC"
	case bvs:
		s = "BVS"
	case clc:
		s = "CLC"
	case cld:
		s = "CLD"
	case cli:
		s = "CLI"
	case clv:
		s = "CLV"
	case cmp:
		s = "CMP"
	case cpx:
		s = "CPX"
	case cpy:
		s = "CPY"
	case dec:
		s = "DEC"
	case dex:
		s = "DEX"
	case dey:
		s = "DEY"
	case eor:
		s = "EOR"
	case inc:
		s = "INC"
	case inx:
		s = "INX"
	case iny:
		s = "INY"
	case jmp:
		s = "JMP"
	case jsr:
		s = "JSR"
	case lda:
		s = "LDA"
	case ldx:
		s = "LDX"
	case ldy:
		s = "LDY"
	case lsr:
		s = "LSR"
	case nop:
		s = "NOP"
	case ora:
		s = "ORA"
	case pha:
		s = "PHA"
	case php:
		s = "PHP"
	case pla:
		s = "PLA"
	case plp:
		s = "PLP"
	case rol:
		s = "ROL"
	case ror:
		s = "ROR"
	case rti:
		s = "RTI"
	case rts:
		s = "RTS"
	case sbc:
		s = "SBC"
	case sec:
		s = "SEC"
	case sed:
		s = "SED"
	case sei:
		s = "SEI"
	case sta:
		s = "STA"
	case stx:
		s = "STX"
	case sty:
		s = "STY"
	case tax:
		s = "TAX"
	case tay:
		s = "TAY"
	case tsx:
		s = "TSX"
	case txa:
		s = "TXA"
	case txs:
		s = "TXS"
	case tya:
		s = "TYA"
	}
	return
}
