package cpu

import (
	"fmt"

	"github.com/pda/go6502/bus"
)

// Instruction is an OpType plus its operand.
// One or both of the operand types will be zero.
// This is determined by (ot.Bytes - 1) / 8
type Instruction struct {
	OpType

	// The single-byte operand, for 2-byte instructions.
	Op8 uint8

	// The 16-bit operand, for 3-byte instructions.
	Op16 uint16
}

func (in Instruction) String() (s string) {
	switch in.Bytes {
	case 3:
		s = fmt.Sprintf("%v $%04X", in.OpType, in.Op16)
	case 2:
		s = fmt.Sprintf("%v $%02X", in.OpType, in.Op8)
	case 1:
		s = in.OpType.String()
	}
	return
}

// ReadInstruction reads an instruction from the bus starting at the given
// address. An instruction may be 1, 2 or 3 bytes long, including its optional
// 8 or 16 bit operand.
func ReadInstruction(pc uint16, bus *bus.Bus) Instruction {
	opcode := bus.Read(pc)
	optype, ok := optypes[opcode]
	if !ok {
		panic(fmt.Sprintf("Illegal opcode $%02X at $%04X", opcode, pc))
	}
	in := Instruction{OpType: optype}
	switch in.Bytes {
	case 1: // no operand
	case 2:
		in.Op8 = bus.Read(pc + 1)
	case 3:
		in.Op16 = bus.Read16(pc + 1)
	default:
		panic(fmt.Sprintf("unhandled instruction length: %d", in.Bytes))
	}
	return in
}
