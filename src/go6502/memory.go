package go6502

import(
  "fmt"
)

// Memory

type Memory interface {
  Read(address) byte;
  Write(address, byte);
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
  om.Memory.Write(a - om.offset, value)
}

// Ram (32K)

type Ram [0x8000]byte
func (r *Ram) String() string {
  return "(RAM 32K)"
}

func (mem *Ram) Read(a address) byte {
  return mem[a]
}

func (mem *Ram) Write(a address, value byte) {
  mem[a] = value
}

func (mem *Ram) Size() int {
  return 0x8000 // 32K
}
