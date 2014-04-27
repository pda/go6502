package bus

import (
	"fmt"

	"github.com/pda/go6502/memory"
)

// OffsetMemory wraps a Memory object, rewriting read/write addresses by the
// given offset. This makes it possible to mount memory into a larger address
// space at a given base address.
type OffsetMemory struct {
	Offset uint16
	memory.Memory
}

// Read returns a byte from the underlying Memory after rewriting the address
// using the offset.
func (om OffsetMemory) Read(a uint16) byte {
	return om.Memory.Read(a - om.Offset)
}

func (om OffsetMemory) String() string {
	return fmt.Sprintf("OffsetMemory(%v)", om.Memory)
}

// Write stores a byte in the underlying Memory after rewriting the address
// using the offset.
func (om OffsetMemory) Write(a uint16, value byte) {
	om.Memory.Write(a-om.Offset, value)
}
