package go6502

import (
	"fmt"
)

// Instruction is an OpType plus its operand.
// One or both of the operand types will be zero.
// This is determined by (ot.bytes - 1) / 8
type Instruction struct {
	ot   *OpType
	op8  uint8
	op16 address
}

func (in *Instruction) String() (s string) {
	switch in.ot.bytes {
	case 3:
		s = fmt.Sprintf("%v $%04X", in.ot, in.op16)
	case 2:
		s = fmt.Sprintf("%v $%02X", in.ot, in.op8)
	case 1:
		s = in.ot.String()
	}
	return
}
