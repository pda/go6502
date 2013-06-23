package c64

import(
  "fmt"
)

// Address

type address uint16;

// Memory

type Memory interface {
  Read(address) byte;
  Write(address, byte);
}

// OffsetMemory

type OffsetMemory struct {
  offset address
  memory Memory
}

func (om OffsetMemory) Read(a address) byte {
  return om.memory.Read(a - om.offset)
}

func (om OffsetMemory) String() string {
  return fmt.Sprintf("OffsetMemory(%v)", om.memory)
}

func (om OffsetMemory) Write(a address, value byte) {
  om.memory.Write(a - om.offset, value)
}

// Ram

type Ram [0x10000]byte
func (r *Ram) String() string {
  return "(RAM 64K)"
}

func (mem *Ram) Read(a address) byte {
  return mem[a]
}

func (mem *Ram) Write(a address, value byte) {
  mem[a] = value
}
