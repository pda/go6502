package go6502

import (
	"fmt"
)

// Instruction is an OpType plus its operand.
// One or both of the operand types will be zero.
// This is determined by (ot.bytes - 1) / 8
type Instruction struct {
	OpType
	op8  uint8
	op16 uint16
}

func (in *Instruction) String() (s string) {
	switch in.bytes {
	case 3:
		s = fmt.Sprintf("%v $%04X", in.OpType, in.op16)
	case 2:
		s = fmt.Sprintf("%v $%02X", in.OpType, in.op8)
	case 1:
		s = in.OpType.String()
	}
	return
}

// ReadInstruction reads an instruction from the bus starting at the given
// address. An instruction may be 1, 2 or 3 bytes long, including its optional
// 8 or 16 bit operand.
func ReadInstruction(PC uint16, bus *Bus) *Instruction {
	in := &Instruction{OpType: optypes[bus.Read(PC)]}
	switch in.bytes {
	case 1: // no operand
	case 2:
		in.op8 = bus.Read(PC + 1)
	case 3:
		in.op16 = bus.Read16(PC + 1)
	default:
		panic(fmt.Sprintf("unhandled instruction length: %d", in.bytes))
	}
	return in
}
