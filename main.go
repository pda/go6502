package main

import(
  "fmt"
)

func main() {

  addressBus := &Bus{Ram: &Ram{}}
  addressBus.Write16(0xFFFC, 0xDEAD) // Start address, normally on ROM.
  fmt.Println(addressBus)

  cpu := &Cpu{Bus: addressBus}
  cpu.Reset()
  fmt.Println(cpu)

}

// Address
type address uint16;

// Bus

type Bus struct {
  Ram *Ram
}
func (b *Bus) String() string {
  return fmt.Sprintf("Bus Ram:%v", b.Ram)
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
func (b *Bus) Write(a address, value byte) {
  fmt.Printf("Bus[0x%04X] <-- 0x%02X\n", a, value)
  b.Ram[a] = value
}
func (b *Bus) Write16(a address, value address) {
  b.Write(a, byte(value))
  b.Write(a + 1, byte(value >> 8))
}


// Cpu

type Cpu struct {
  pc address
  ac byte
  x byte
  y byte
  sp byte
  sr byte
  Bus *Bus
}
func (c *Cpu) String() string {
  return fmt.Sprintf(
    "CPU pc:0x%04X ac:0x%02X x:0x%02X y:0x%02X sp:0x%02X sr:0x%02X",
    c.pc, c.ac, c.x, c.y, c.sp, c.sr)
}
func (c *Cpu) Reset() {
  c.pc = c.Bus.Read16(0xFFFC)
  c.ac = 0x00
  c.x = 0x00
  c.y = 0x00
  c.sp = 0xFF // address relative to second page of memory (0x0100 ~ 0x01FF)
  c.sr = 0x00
}

// Ram

type Ram [0x10000]byte
func (r *Ram) String() string {
  return "(RAM 64K)"
}
