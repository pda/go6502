package c64

import(
  "fmt"
)

type Bus struct {
  Ram *Ram
  Kernal *Rom
}

func (b *Bus) backendFor(a address) (mem Memory) {
  if a >= 0xE000 && a <= 0xFFFF {
    return OffsetMemory{offset: 0xE000, memory: b.Kernal}
  } else {
    return b.Ram
  }
}

func (b *Bus) Read(a address) byte {
  mem := b.backendFor(a)
  value := mem.Read(a)
  fmt.Printf("Bus[0x%04X] %v --> 0x%02X\n", a, mem, value)
  return value
}

func (b *Bus) Read16(a address) address {
  lo := address(b.Read(a))
  hi := address(b.Read(a + 1))
  return hi << 8 | lo
}

func (b *Bus) String() string {
  return fmt.Sprintf("Bus Ram:%v Kernal:%v",
    b.Ram, b.Kernal)
}

func (b *Bus) Write(a address, value byte) {
  mem := b.backendFor(a)
  fmt.Printf("Bus[0x%04X] <-- 0x%02X\n", a, value)
  mem.Write(a, value)
}

func (b *Bus) Write16(a address, value address) {
  b.Write(a, byte(value))
  b.Write(a + 1, byte(value >> 8))
}
