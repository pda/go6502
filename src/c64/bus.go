package c64

import(
  "fmt"
)

type Bus struct {
  Ram *Ram
}

func (b *Bus) Read(a address) byte {
  value := b.Ram[a]
  fmt.Printf("Bus[0x%04X] --> 0x%02X\n", a, value)
  return value
}

func (b *Bus) Read16(a address) address {
  lo := address(b.Read(a))
  hi := address(b.Read(a + 1))
  return hi << 8 | lo
}

func (b *Bus) String() string {
  return fmt.Sprintf("Bus Ram:%v", b.Ram)
}

func (b *Bus) Write(a address, value byte) {
  fmt.Printf("Bus[0x%04X] <-- 0x%02X\n", a, value)
  b.Ram[a] = value
}

func (b *Bus) Write16(a address, value address) {
  b.Write(a, byte(value))
  b.Write(a + 1, byte(value >> 8))
}
