package cpu

import "fmt"

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

	// 65C02 only
	zpindirect
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
	"(zeropage)",
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

	// 65{S}C02-only
	bra
	phx
	phy
	plx
	ply
	stz
	trb
	tsb

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

	// 65{S}C02-only
	"BRA",
	"PHX",
	"PHY",
	"PLX",
	"PLY",
	"STZ",
	"TRB",
	"TSB",

	"_END",
}

// OpType represents a 6502 op-code instruction, including the addressing
// mode encoded into the op-code, but not the operand value following the
// opcode in memory.
type OpType struct {

	// Opcode is a byte representing an instruction and its addressing mode.
	Opcode byte

	// id is an internal identifier of the instruction type/mnemonic, e.g. ADC
	id uint8

	// addressing is an internal identifier for the addressing mode.
	addressing uint8

	// Bytes is the size of the instruction with its operand.
	// Opcodes with implicit/null operand are 1 byte.
	// Opcodes with immediate or zeropage operand are 2 bytes.
	// Opcodes with adddress operands are 3 bytes.
	Bytes uint8

	// Cycles is the number of times the system clock signal will rise and fall
	// before the instruction is complete.
	Cycles uint8
}

func (ot OpType) String() string {
	return fmt.Sprintf("%s %s", ot.Name(), addressingNames[ot.addressing])
}

// Name returns the instruction mnemonic name, e.g. ADC or TYA.
func (ot OpType) Name() (s string) {
	return instructionNames[ot.id]
}

func (ot OpType) IsAbsolute() bool {
	return ot.addressing == absolute
}

var optypes = map[uint8]OpType{
	0x69: OpType{0x69, adc, immediate, 2, 2},
	0x65: OpType{0x65, adc, zeropage, 2, 3},
	0x75: OpType{0x75, adc, zeropageX, 2, 4},
	0x6D: OpType{0x6D, adc, absolute, 3, 4},
	0x7D: OpType{0x7D, adc, absoluteX, 3, 4},
	0x79: OpType{0x79, adc, absoluteY, 3, 4},
	0x61: OpType{0x61, adc, indirectX, 2, 6},
	0x71: OpType{0x71, adc, indirectY, 2, 5},
	0x29: OpType{0x29, and, immediate, 2, 2},
	0x25: OpType{0x25, and, zeropage, 2, 3},
	0x35: OpType{0x35, and, zeropageX, 2, 4},
	0x2D: OpType{0x2D, and, absolute, 3, 4},
	0x3D: OpType{0x3D, and, absoluteX, 3, 4},
	0x39: OpType{0x39, and, absoluteY, 3, 4},
	0x21: OpType{0x21, and, indirectX, 2, 6},
	0x31: OpType{0x31, and, indirectY, 2, 5},
	0x0A: OpType{0x0A, asl, accumulator, 1, 2},
	0x06: OpType{0x06, asl, zeropage, 2, 5},
	0x16: OpType{0x16, asl, zeropageX, 2, 6},
	0x0E: OpType{0x0E, asl, absolute, 3, 6},
	0x1E: OpType{0x1E, asl, absoluteX, 3, 7},
	0x90: OpType{0x90, bcc, relative, 2, 2},
	0xB0: OpType{0xB0, bcs, relative, 2, 2},
	0xF0: OpType{0xF0, beq, relative, 2, 2},
	0x24: OpType{0x24, bit, zeropage, 2, 3},
	0x2C: OpType{0x2C, bit, absolute, 3, 4},
	0x30: OpType{0x30, bmi, relative, 2, 2},
	0xD0: OpType{0xD0, bne, relative, 2, 2},
	0x10: OpType{0x10, bpl, relative, 2, 2},
	0x00: OpType{0x00, brk, implied, 1, 7},
	0x50: OpType{0x50, bvc, relative, 2, 2},
	0x70: OpType{0x70, bvs, relative, 2, 2},
	0x18: OpType{0x18, clc, implied, 1, 2},
	0xD8: OpType{0xD8, cld, implied, 1, 2},
	0x58: OpType{0x58, cli, implied, 1, 2},
	0xB8: OpType{0xB8, clv, implied, 1, 2},
	0xC9: OpType{0xC9, cmp, immediate, 2, 2},
	0xC5: OpType{0xC5, cmp, zeropage, 2, 3},
	0xD5: OpType{0xD5, cmp, zeropageX, 2, 4},
	0xCD: OpType{0xCD, cmp, absolute, 3, 4},
	0xDD: OpType{0xDD, cmp, absoluteX, 3, 4},
	0xD9: OpType{0xD9, cmp, absoluteY, 3, 4},
	0xC1: OpType{0xC1, cmp, indirectX, 2, 6},
	0xD1: OpType{0xD1, cmp, indirectY, 2, 5},
	0xE0: OpType{0xE0, cpx, immediate, 2, 2},
	0xE4: OpType{0xE4, cpx, zeropage, 2, 3},
	0xEC: OpType{0xEC, cpx, absolute, 3, 4},
	0xC0: OpType{0xC0, cpy, immediate, 2, 2},
	0xC4: OpType{0xC4, cpy, zeropage, 2, 3},
	0xCC: OpType{0xCC, cpy, absolute, 3, 4},
	0xC6: OpType{0xC6, dec, zeropage, 2, 5},
	0xD6: OpType{0xD6, dec, zeropageX, 2, 6},
	0xCE: OpType{0xCE, dec, absolute, 3, 3},
	0xDE: OpType{0xDE, dec, absoluteX, 3, 7},
	0xCA: OpType{0xCA, dex, implied, 1, 2},
	0x88: OpType{0x88, dey, implied, 1, 2},
	0x49: OpType{0x49, eor, immediate, 2, 2},
	0x45: OpType{0x45, eor, zeropage, 2, 3},
	0x55: OpType{0x55, eor, zeropageX, 2, 4},
	0x4D: OpType{0x4D, eor, absolute, 3, 4},
	0x5D: OpType{0x5D, eor, absoluteX, 3, 4},
	0x59: OpType{0x59, eor, absoluteY, 3, 4},
	0x41: OpType{0x41, eor, indirectX, 2, 6},
	0x51: OpType{0x51, eor, indirectY, 2, 5},
	0xE6: OpType{0xE6, inc, zeropage, 2, 5},
	0xF6: OpType{0xF6, inc, zeropageX, 2, 6},
	0xEE: OpType{0xEE, inc, absolute, 3, 6},
	0xFE: OpType{0xFE, inc, absoluteX, 3, 7},
	0xE8: OpType{0xE8, inx, implied, 1, 2},
	0xC8: OpType{0xC8, iny, implied, 1, 2},
	0x4C: OpType{0x4C, jmp, absolute, 3, 3},
	0x6C: OpType{0x6C, jmp, indirect, 3, 5},
	0x20: OpType{0x20, jsr, absolute, 3, 6},
	0xA9: OpType{0xA9, lda, immediate, 2, 2},
	0xA5: OpType{0xA5, lda, zeropage, 2, 3},
	0xB5: OpType{0xB5, lda, zeropageX, 2, 4},
	0xAD: OpType{0xAD, lda, absolute, 3, 4},
	0xBD: OpType{0xBD, lda, absoluteX, 3, 4},
	0xB9: OpType{0xB9, lda, absoluteY, 3, 4},
	0xA1: OpType{0xA1, lda, indirectX, 2, 6},
	0xB1: OpType{0xB1, lda, indirectY, 2, 5},
	0xA2: OpType{0xA2, ldx, immediate, 2, 2},
	0xA6: OpType{0xA6, ldx, zeropage, 2, 3},
	0xB6: OpType{0xB6, ldx, zeropageY, 2, 4},
	0xAE: OpType{0xAE, ldx, absolute, 3, 4},
	0xBE: OpType{0xBE, ldx, absoluteY, 3, 4},
	0xA0: OpType{0xA0, ldy, immediate, 2, 2},
	0xA4: OpType{0xA4, ldy, zeropage, 2, 3},
	0xB4: OpType{0xB4, ldy, zeropageX, 2, 4},
	0xAC: OpType{0xAC, ldy, absolute, 3, 4},
	0xBC: OpType{0xBC, ldy, absoluteX, 3, 4},
	0x4A: OpType{0x4A, lsr, accumulator, 1, 2},
	0x46: OpType{0x46, lsr, zeropage, 2, 5},
	0x56: OpType{0x56, lsr, zeropageX, 2, 6},
	0x4E: OpType{0x4E, lsr, absolute, 3, 6},
	0x5E: OpType{0x5E, lsr, absoluteX, 3, 7},
	0xEA: OpType{0xEA, nop, implied, 1, 2},
	0x09: OpType{0x09, ora, immediate, 2, 2},
	0x05: OpType{0x05, ora, zeropage, 2, 3},
	0x15: OpType{0x15, ora, zeropageX, 2, 4},
	0x0D: OpType{0x0D, ora, absolute, 3, 4},
	0x1D: OpType{0x1D, ora, absoluteX, 3, 4},
	0x19: OpType{0x19, ora, absoluteY, 3, 4},
	0x01: OpType{0x01, ora, indirectX, 2, 6},
	0x11: OpType{0x11, ora, indirectY, 2, 5},
	0x48: OpType{0x48, pha, implied, 1, 3},
	0x08: OpType{0x08, php, implied, 1, 3},
	0x68: OpType{0x68, pla, implied, 1, 4},
	0x28: OpType{0x28, plp, implied, 1, 4},
	0x2A: OpType{0x2A, rol, accumulator, 1, 2},
	0x26: OpType{0x26, rol, zeropage, 2, 5},
	0x36: OpType{0x36, rol, zeropageX, 2, 6},
	0x2E: OpType{0x2E, rol, absolute, 3, 6},
	0x3E: OpType{0x3E, rol, absoluteX, 3, 7},
	0x6A: OpType{0x6A, ror, accumulator, 1, 2},
	0x66: OpType{0x66, ror, zeropage, 2, 5},
	0x76: OpType{0x76, ror, zeropageX, 2, 6},
	0x6E: OpType{0x6E, ror, absolute, 3, 6},
	0x7E: OpType{0x7E, ror, absoluteX, 3, 7},
	0x40: OpType{0x40, rti, implied, 1, 6},
	0x60: OpType{0x60, rts, implied, 1, 6},
	0xE9: OpType{0xE9, sbc, immediate, 2, 2},
	0xE5: OpType{0xE5, sbc, zeropage, 2, 3},
	0xF5: OpType{0xF5, sbc, zeropageX, 2, 4},
	0xED: OpType{0xED, sbc, absolute, 3, 4},
	0xFD: OpType{0xFD, sbc, absoluteX, 3, 4},
	0xF9: OpType{0xF9, sbc, absoluteY, 3, 4},
	0xE1: OpType{0xE1, sbc, indirectX, 2, 6},
	0xF1: OpType{0xF1, sbc, indirectY, 2, 5},
	0x38: OpType{0x38, sec, implied, 1, 2},
	0xF8: OpType{0xF8, sed, implied, 1, 2},
	0x78: OpType{0x78, sei, implied, 1, 2},
	0x85: OpType{0x85, sta, zeropage, 2, 3},
	0x95: OpType{0x95, sta, zeropageX, 2, 4},
	0x8D: OpType{0x8D, sta, absolute, 3, 4},
	0x9D: OpType{0x9D, sta, absoluteX, 3, 5},
	0x99: OpType{0x99, sta, absoluteY, 3, 5},
	0x81: OpType{0x81, sta, indirectX, 2, 6},
	0x91: OpType{0x91, sta, indirectY, 2, 6},
	0x86: OpType{0x86, stx, zeropage, 2, 3},
	0x96: OpType{0x96, stx, zeropageY, 2, 4},
	0x8E: OpType{0x8E, stx, absolute, 3, 4},
	0x84: OpType{0x84, sty, zeropage, 2, 3},
	0x94: OpType{0x94, sty, zeropageX, 2, 4},
	0x8C: OpType{0x8C, sty, absolute, 3, 4},
	0xAA: OpType{0xAA, tax, implied, 1, 2},
	0xA8: OpType{0xA8, tay, implied, 1, 2},
	0xBA: OpType{0xBA, tsx, implied, 1, 2},
	0x8A: OpType{0x8A, txa, implied, 1, 2},
	0x9A: OpType{0x9A, txs, implied, 1, 2},
	0x98: OpType{0x98, tya, implied, 1, 2},

	// 65C02 only

	// Additional addressing modes
	0x12: OpType{0x12, ora, zpindirect, 2, 5},
	0x32: OpType{0x32, and, zpindirect, 2, 5},
	0x52: OpType{0x52, eor, zpindirect, 2, 5},
	0x72: OpType{0x72, adc, zpindirect, 2, 5},
	0x92: OpType{0x92, sta, zpindirect, 2, 5},
	0xB2: OpType{0xB2, lda, zpindirect, 2, 5},
	0xD2: OpType{0xD2, cmp, zpindirect, 2, 5},
	0xF2: OpType{0xF2, sbc, zpindirect, 2, 5},
	0x89: OpType{0x89, bit, immediate, 2, 2},
	0x34: OpType{0x34, bit, zeropageX, 2, 4},
	0x3C: OpType{0x3C, bit, absoluteX, 3, 4},
	0x3A: OpType{0x3A, dec, implied, 1, 2},
	0x1A: OpType{0x1A, inc, implied, 1, 2},
	0x7C: OpType{0x7C, jmp, absoluteX, 3, 6},

	// New instructions
	0x80: OpType{0x80, bra, relative, 2, 3},
	0xDA: OpType{0xDA, phx, implied, 1, 3},
	0x5A: OpType{0x5A, phy, implied, 1, 3},
	0xFA: OpType{0xFA, plx, implied, 1, 4},
	0x7A: OpType{0x7A, ply, implied, 1, 4},
	0x64: OpType{0x64, stz, zeropage, 2, 3},
	0x74: OpType{0x74, stz, zeropageX, 2, 4},
	0x9C: OpType{0x9C, stz, absolute, 3, 4},
	0x9E: OpType{0x9E, stz, absoluteX, 3, 5},
	0x14: OpType{0x14, trb, zeropage, 2, 5},
	0x1C: OpType{0x1C, trb, absolute, 3, 6},
	0x04: OpType{0x04, tsb, zeropage, 2, 5},
	0x0C: OpType{0x0C, tsb, absolute, 3, 5},

	0xFF: OpType{0xFF, _end, implied, 1, 1},
}
