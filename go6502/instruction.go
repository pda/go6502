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

var instructionNames = [...]string{
	"",
	"ADC",
	"AND",
	"ASL",
	"BCC",
	"BCS",
	"BEQ",
	"BIT",
	"BMI",
	"BNE",
	"BPL",
	"BRK",
	"BVC",
	"BVS",
	"CLC",
	"CLD",
	"CLI",
	"CLV",
	"CMP",
	"CPX",
	"CPY",
	"DEC",
	"DEX",
	"DEY",
	"EOR",
	"INC",
	"INX",
	"INY",
	"JMP",
	"JSR",
	"LDA",
	"LDX",
	"LDY",
	"LSR",
	"NOP",
	"ORA",
	"PHA",
	"PHP",
	"PLA",
	"PLP",
	"ROL",
	"ROR",
	"RTI",
	"RTS",
	"SBC",
	"SEC",
	"SED",
	"SEI",
	"STA",
	"STX",
	"STY",
	"TAX",
	"TAY",
	"TSX",
	"TXA",
	"TXS",
	"TYA",
	"_END",
}

// OpType represents a 6502 op-code instruction, including the addressing
// mode encoded into the op-code, but not the operand value following the
// opcode in memory.
type OpType struct {
	id         uint8 // the const identifier of the instruction type, e.g. ADC
	opcode     byte
	addressing int
	bytes      int
	cycles     int
	flags      int
}

func (ot *OpType) String() string {
	return fmt.Sprintf("%s %s", ot.name(), addressingNames[ot.addressing])
}

func (ot *OpType) name() (s string) {
	return instructionNames[ot.id]
}

// Iop is an instruction with its operand.
// One or both of the operand types will be zero.
// This is determined by (ot.bytes - 1) / 8
type Iop struct {
	ot   *OpType
	op8  uint8
	op16 address
}

func (iop *Iop) String() (s string) {
	switch iop.ot.bytes {
	case 3:
		s = fmt.Sprintf("%v $%04X", iop.ot, iop.op16)
	case 2:
		s = fmt.Sprintf("%v $%02X", iop.ot, iop.op8)
	case 1:
		s = iop.ot.String()
	}
	return
}

func findInstruction(opcode byte) *OpType {
	// TODO: singleton instructions; they're immutable.
	var ot OpType
	switch opcode {
	default:
		panic(fmt.Sprintf("Unknown opcode: 0x%02X", opcode))
	case 0x69:
		ot = OpType{adc, opcode, immediate, 2, 2, 0}
	case 0x65:
		ot = OpType{adc, opcode, zeropage, 2, 3, 0}
	case 0x75:
		ot = OpType{adc, opcode, zeropageX, 2, 4, 0}
	case 0x6D:
		ot = OpType{adc, opcode, absolute, 3, 4, 0}
	case 0x7D:
		ot = OpType{adc, opcode, absoluteX, 3, 4, 0}
	case 0x79:
		ot = OpType{adc, opcode, absoluteY, 3, 4, 0}
	case 0x61:
		ot = OpType{adc, opcode, indirectX, 2, 6, 0}
	case 0x71:
		ot = OpType{adc, opcode, indirectY, 2, 5, 0}
	case 0x29:
		ot = OpType{and, opcode, immediate, 2, 2, 0}
	case 0x25:
		ot = OpType{and, opcode, zeropage, 2, 3, 0}
	case 0x35:
		ot = OpType{and, opcode, zeropageX, 2, 4, 0}
	case 0x2D:
		ot = OpType{and, opcode, absolute, 3, 4, 0}
	case 0x3D:
		ot = OpType{and, opcode, absoluteX, 3, 4, 0}
	case 0x39:
		ot = OpType{and, opcode, absoluteY, 3, 4, 0}
	case 0x21:
		ot = OpType{and, opcode, indirectX, 2, 6, 0}
	case 0x31:
		ot = OpType{and, opcode, indirectY, 2, 5, 0}
	case 0x0A:
		ot = OpType{asl, opcode, accumulator, 1, 2, 0}
	case 0x06:
		ot = OpType{asl, opcode, zeropage, 2, 5, 0}
	case 0x16:
		ot = OpType{asl, opcode, zeropageX, 2, 6, 0}
	case 0x0E:
		ot = OpType{asl, opcode, absolute, 3, 6, 0}
	case 0x1E:
		ot = OpType{asl, opcode, absoluteX, 3, 7, 0}
	case 0x90:
		ot = OpType{bcc, opcode, relative, 2, 2, 0}
	case 0xB0:
		ot = OpType{bcs, opcode, relative, 2, 2, 0}
	case 0xF0:
		ot = OpType{beq, opcode, relative, 2, 2, 0}
	case 0x24:
		ot = OpType{bit, opcode, zeropage, 2, 3, 0}
	case 0x2C:
		ot = OpType{bit, opcode, absolute, 3, 4, 0}
	case 0x30:
		ot = OpType{bmi, opcode, relative, 2, 2, 0}
	case 0xD0:
		ot = OpType{bne, opcode, relative, 2, 2, 0}
	case 0x10:
		ot = OpType{bpl, opcode, relative, 2, 2, 0}
	case 0x00:
		ot = OpType{brk, opcode, implied, 1, 7, 0}
	case 0x50:
		ot = OpType{bvc, opcode, relative, 2, 2, 0}
	case 0x70:
		ot = OpType{bvs, opcode, relative, 2, 2, 0}
	case 0x18:
		ot = OpType{clc, opcode, implied, 1, 2, 0}
	case 0xD8:
		ot = OpType{cld, opcode, implied, 1, 2, 0}
	case 0x58:
		ot = OpType{cli, opcode, implied, 1, 2, 0}
	case 0xB8:
		ot = OpType{clv, opcode, implied, 1, 2, 0}
	case 0xC9:
		ot = OpType{cmp, opcode, immediate, 2, 2, 0}
	case 0xC5:
		ot = OpType{cmp, opcode, zeropage, 2, 3, 0}
	case 0xD5:
		ot = OpType{cmp, opcode, zeropageX, 2, 4, 0}
	case 0xCD:
		ot = OpType{cmp, opcode, absolute, 3, 4, 0}
	case 0xDD:
		ot = OpType{cmp, opcode, absoluteX, 3, 4, 0}
	case 0xD9:
		ot = OpType{cmp, opcode, absoluteY, 3, 4, 0}
	case 0xC1:
		ot = OpType{cmp, opcode, indirectX, 2, 6, 0}
	case 0xD1:
		ot = OpType{cmp, opcode, indirectY, 2, 5, 0}
	case 0xE0:
		ot = OpType{cpx, opcode, immediate, 2, 2, 0}
	case 0xE4:
		ot = OpType{cpx, opcode, zeropage, 2, 3, 0}
	case 0xEC:
		ot = OpType{cpx, opcode, absolute, 3, 4, 0}
	case 0xC0:
		ot = OpType{cpy, opcode, immediate, 2, 2, 0}
	case 0xC4:
		ot = OpType{cpy, opcode, zeropage, 2, 3, 0}
	case 0xCC:
		ot = OpType{cpy, opcode, absolute, 3, 4, 0}
	case 0xC6:
		ot = OpType{dec, opcode, zeropage, 2, 5, 0}
	case 0xD6:
		ot = OpType{dec, opcode, zeropageX, 2, 6, 0}
	case 0xCE:
		ot = OpType{dec, opcode, absolute, 3, 3, 0}
	case 0xDE:
		ot = OpType{dec, opcode, absoluteX, 3, 7, 0}
	case 0xCA:
		ot = OpType{dex, opcode, implied, 1, 2, 0}
	case 0x88:
		ot = OpType{dey, opcode, implied, 1, 2, 0}
	case 0x49:
		ot = OpType{eor, opcode, immediate, 2, 2, 0}
	case 0x45:
		ot = OpType{eor, opcode, zeropage, 2, 3, 0}
	case 0x55:
		ot = OpType{eor, opcode, zeropageX, 2, 4, 0}
	case 0x4D:
		ot = OpType{eor, opcode, absolute, 3, 4, 0}
	case 0x5D:
		ot = OpType{eor, opcode, absoluteX, 3, 4, 0}
	case 0x59:
		ot = OpType{eor, opcode, absoluteY, 3, 4, 0}
	case 0x41:
		ot = OpType{eor, opcode, indirectX, 2, 6, 0}
	case 0x51:
		ot = OpType{eor, opcode, indirectY, 2, 5, 0}
	case 0xE6:
		ot = OpType{inc, opcode, zeropage, 2, 5, 0}
	case 0xF6:
		ot = OpType{inc, opcode, zeropageX, 2, 6, 0}
	case 0xEE:
		ot = OpType{inc, opcode, absolute, 3, 6, 0}
	case 0xFE:
		ot = OpType{inc, opcode, absoluteX, 3, 7, 0}
	case 0xE8:
		ot = OpType{inx, opcode, implied, 1, 2, 0}
	case 0xC8:
		ot = OpType{iny, opcode, implied, 1, 2, 0}
	case 0x4C:
		ot = OpType{jmp, opcode, absolute, 3, 3, 0}
	case 0x6C:
		ot = OpType{jmp, opcode, indirect, 3, 5, 0}
	case 0x20:
		ot = OpType{jsr, opcode, absolute, 3, 6, 0}
	case 0xA9:
		ot = OpType{lda, opcode, immediate, 2, 2, 0}
	case 0xA5:
		ot = OpType{lda, opcode, zeropage, 2, 3, 0}
	case 0xB5:
		ot = OpType{lda, opcode, zeropageX, 2, 4, 0}
	case 0xAD:
		ot = OpType{lda, opcode, absolute, 3, 4, 0}
	case 0xBD:
		ot = OpType{lda, opcode, absoluteX, 3, 4, 0}
	case 0xB9:
		ot = OpType{lda, opcode, absoluteY, 3, 4, 0}
	case 0xA1:
		ot = OpType{lda, opcode, indirectX, 2, 6, 0}
	case 0xB1:
		ot = OpType{lda, opcode, indirectY, 2, 5, 0}
	case 0xA2:
		ot = OpType{ldx, opcode, immediate, 2, 2, 0}
	case 0xA6:
		ot = OpType{ldx, opcode, zeropage, 2, 3, 0}
	case 0xB6:
		ot = OpType{ldx, opcode, zeropageY, 2, 4, 0}
	case 0xAE:
		ot = OpType{ldx, opcode, absolute, 3, 4, 0}
	case 0xBE:
		ot = OpType{ldx, opcode, absoluteY, 3, 4, 0}
	case 0xA0:
		ot = OpType{ldy, opcode, immediate, 2, 2, 0}
	case 0xA4:
		ot = OpType{ldy, opcode, zeropage, 2, 3, 0}
	case 0xB4:
		ot = OpType{ldy, opcode, zeropageX, 2, 4, 0}
	case 0xAC:
		ot = OpType{ldy, opcode, absolute, 3, 4, 0}
	case 0xBC:
		ot = OpType{ldy, opcode, absoluteX, 3, 4, 0}
	case 0x4A:
		ot = OpType{lsr, opcode, accumulator, 1, 2, 0}
	case 0x46:
		ot = OpType{lsr, opcode, zeropage, 2, 5, 0}
	case 0x56:
		ot = OpType{lsr, opcode, zeropageX, 2, 6, 0}
	case 0x4E:
		ot = OpType{lsr, opcode, absolute, 3, 6, 0}
	case 0x5E:
		ot = OpType{lsr, opcode, absoluteX, 3, 7, 0}
	case 0xEA:
		ot = OpType{nop, opcode, implied, 1, 2, 0}
	case 0x09:
		ot = OpType{ora, opcode, immediate, 2, 2, 0}
	case 0x05:
		ot = OpType{ora, opcode, zeropage, 2, 3, 0}
	case 0x15:
		ot = OpType{ora, opcode, zeropageX, 2, 4, 0}
	case 0x0D:
		ot = OpType{ora, opcode, absolute, 3, 4, 0}
	case 0x1D:
		ot = OpType{ora, opcode, absoluteX, 3, 4, 0}
	case 0x19:
		ot = OpType{ora, opcode, absoluteY, 3, 4, 0}
	case 0x01:
		ot = OpType{ora, opcode, indirectX, 2, 6, 0}
	case 0x11:
		ot = OpType{ora, opcode, indirectY, 2, 5, 0}
	case 0x48:
		ot = OpType{pha, opcode, implied, 1, 3, 0}
	case 0x08:
		ot = OpType{php, opcode, implied, 1, 3, 0}
	case 0x68:
		ot = OpType{pla, opcode, implied, 1, 4, 0}
	case 0x28:
		ot = OpType{php, opcode, implied, 1, 4, 0}
	case 0x2A:
		ot = OpType{rol, opcode, accumulator, 1, 2, 0}
	case 0x26:
		ot = OpType{rol, opcode, zeropage, 2, 5, 0}
	case 0x36:
		ot = OpType{rol, opcode, zeropageX, 2, 6, 0}
	case 0x2E:
		ot = OpType{rol, opcode, absolute, 3, 6, 0}
	case 0x3E:
		ot = OpType{rol, opcode, absoluteX, 3, 7, 0}
	case 0x6A:
		ot = OpType{ror, opcode, accumulator, 1, 2, 0}
	case 0x66:
		ot = OpType{ror, opcode, zeropage, 2, 5, 0}
	case 0x76:
		ot = OpType{ror, opcode, zeropageX, 2, 6, 0}
	case 0x6E:
		ot = OpType{ror, opcode, absolute, 3, 6, 0}
	case 0x7E:
		ot = OpType{ror, opcode, absoluteX, 3, 7, 0}
	case 0x40:
		ot = OpType{rti, opcode, implied, 1, 6, 0}
	case 0x60:
		ot = OpType{rts, opcode, implied, 1, 6, 0}
	case 0xE9:
		ot = OpType{sbc, opcode, immediate, 2, 2, 0}
	case 0xE5:
		ot = OpType{sbc, opcode, zeropage, 2, 3, 0}
	case 0xF5:
		ot = OpType{sbc, opcode, zeropageX, 2, 4, 0}
	case 0xED:
		ot = OpType{sbc, opcode, absolute, 3, 4, 0}
	case 0xFD:
		ot = OpType{sbc, opcode, absoluteX, 3, 4, 0}
	case 0xF9:
		ot = OpType{sbc, opcode, absoluteY, 3, 4, 0}
	case 0xE1:
		ot = OpType{sbc, opcode, indirectX, 2, 6, 0}
	case 0xF1:
		ot = OpType{sbc, opcode, indirectY, 2, 5, 0}
	case 0x38:
		ot = OpType{sec, opcode, implied, 1, 2, 0}
	case 0xF8:
		ot = OpType{sed, opcode, implied, 1, 2, 0}
	case 0x78:
		ot = OpType{sei, opcode, implied, 1, 2, 0}
	case 0x85:
		ot = OpType{sta, opcode, zeropage, 2, 3, 0}
	case 0x95:
		ot = OpType{sta, opcode, zeropageX, 2, 4, 0}
	case 0x8D:
		ot = OpType{sta, opcode, absolute, 3, 4, 0}
	case 0x9D:
		ot = OpType{sta, opcode, absoluteX, 3, 5, 0}
	case 0x99:
		ot = OpType{sta, opcode, absoluteY, 3, 5, 0}
	case 0x81:
		ot = OpType{sta, opcode, indirectX, 2, 6, 0}
	case 0x91:
		ot = OpType{sta, opcode, indirectY, 2, 6, 0}
	case 0x86:
		ot = OpType{stx, opcode, zeropage, 2, 3, 0}
	case 0x96:
		ot = OpType{stx, opcode, zeropageY, 2, 4, 0}
	case 0x8E:
		ot = OpType{stx, opcode, absolute, 3, 4, 0}
	case 0x84:
		ot = OpType{sty, opcode, zeropage, 2, 3, 0}
	case 0x94:
		ot = OpType{sty, opcode, zeropageX, 2, 4, 0}
	case 0x8C:
		ot = OpType{sty, opcode, absolute, 3, 4, 0}
	case 0xAA:
		ot = OpType{tax, opcode, implied, 1, 2, 0}
	case 0xA8:
		ot = OpType{tay, opcode, implied, 1, 2, 0}
	case 0xBA:
		ot = OpType{tsx, opcode, implied, 1, 2, 0}
	case 0x8A:
		ot = OpType{txa, opcode, implied, 1, 2, 0}
	case 0x9A:
		ot = OpType{txs, opcode, implied, 1, 2, 0}
	case 0x98:
		ot = OpType{tya, opcode, implied, 1, 2, 0}
	case 0xFF:
		ot = OpType{_end, opcode, implied, 1, 1, 0}
	}
	return &ot
}
