package go6502

import (
	"fmt"
)

// Memory

type Memory interface {
	Read(Address) byte
	Write(Address, byte)
	Size() int
}

// OffsetMemory

type OffsetMemory struct {
	offset Address
	Memory
}

func (om OffsetMemory) Read(a Address) byte {
	return om.Memory.Read(a - om.offset)
}

func (om OffsetMemory) String() string {
	return fmt.Sprintf("OffsetMemory(%v)", om.Memory)
}

func (om OffsetMemory) Write(a Address, value byte) {
	om.Memory.Write(a-om.offset, value)
}
