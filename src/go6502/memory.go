package go6502

import (
	"fmt"
)

// Memory

type Memory interface {
	Read(address) byte
	Write(address, byte)
	Size() int
}

// OffsetMemory

type OffsetMemory struct {
	offset address
	Memory
}

func (om OffsetMemory) Read(a address) byte {
	return om.Memory.Read(a - om.offset)
}

func (om OffsetMemory) String() string {
	return fmt.Sprintf("OffsetMemory(%v)", om.Memory)
}

func (om OffsetMemory) Write(a address, value byte) {
	om.Memory.Write(a-om.offset, value)
}
