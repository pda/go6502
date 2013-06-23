package c64

import(
  "fmt"
)

type Cpu struct {
  pc address
  ac byte
  x byte
  y byte
  sp byte
  sr byte
  Bus *Bus
}

func (c *Cpu) Reset() {
  c.pc = c.Bus.Read16(0xFFFC)
  c.ac = 0x00
  c.x = 0x00
  c.y = 0x00
  c.sp = 0xFF // address relative to second page of memory (0x0100 ~ 0x01FF)
  c.sr = 0x00
}

func (c *Cpu) String() string {
  return fmt.Sprintf(
    "CPU pc:0x%04X ac:0x%02X x:0x%02X y:0x%02X sp:0x%02X sr:0x%02X",
    c.pc, c.ac, c.x, c.y, c.sp, c.sr)
}
