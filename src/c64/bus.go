package c64

import(
  "fmt"
)

type Bus struct {
  DataPort *DataPort
  Ram *Ram
  Kernal *Rom
}

// CPU I/O data port.
// The Output Register is located at Address 0x0001.
// The Data Direction Register is at Address 0x0000.
type DataPort struct {
  data [2]byte
}

func (d *DataPort) Read(a address) byte {
  value := d.data[a]
  fmt.Printf("R DataPort[0x%04X] --> 0x%02X\n", a, value)
  return value
}

func (d *DataPort) String() string {
  return "DataPort"
}

func (d *DataPort) Write(a address, value byte) {
  fmt.Printf("W DataPort[0x%04X] <-- 0x%02X\n", a, value)
  d.data[a] = value
}

func (b *Bus) backendFor(a address) (mem Memory) {
  if a <= 0x0001 {
    return b.DataPort
  } else if a >= 0xE000 && a <= 0xFFFF {
    // TODO: cache instance
    return OffsetMemory{offset: 0xE000, memory: b.Kernal}
  } else {
    return b.Ram
  }
}

func (b *Bus) Read(a address) byte {
  mem := b.backendFor(a)
  value := mem.Read(a)
  fmt.Printf("R Bus[0x%04X] %v --> 0x%02X\n", a, mem, value)
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
  fmt.Printf("W Bus[0x%04X] <-- 0x%02X\n", a, value)
  mem.Write(a, value)
}

func (b *Bus) Write16(a address, value address) {
  b.Write(a, byte(value))
  b.Write(a + 1, byte(value >> 8))
}
