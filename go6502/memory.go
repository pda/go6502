package go6502

import (
	"fmt"
)

// Memory is a general interface for reading and writing bytes to and from
// 16-bit addresses.
type Memory interface {
	Read(Address) byte
	Write(Address, byte)
	Size() int
}

// OffsetMemory wraps a Memory object, rewriting read/write addresses by the
// given offset. This makes it possible to mount memory into a larger address
// space at a given base address.
type OffsetMemory struct {
	offset Address
	Memory
}

// Read returns a byte from the underlying Memory after rewriting the address
// using the offset.
func (om OffsetMemory) Read(a Address) byte {
	return om.Memory.Read(a - om.offset)
}

func (om OffsetMemory) String() string {
	return fmt.Sprintf("OffsetMemory(%v)", om.Memory)
}

// Write stores a byte in the underlying Memory after rewriting the address
// using the offset.
func (om OffsetMemory) Write(a Address, value byte) {
	om.Memory.Write(a-om.offset, value)
}
