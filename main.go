package main

import(
  "fmt"
  "c64"
)

func main() {

  kernal := c64.RomFromFile("rom/kernal.rom")
  fmt.Println(kernal)

  addressBus := &c64.Bus{Ram: &c64.Ram{}}
  addressBus.Write16(0xFFFC, 0xDEAD) // Start address, normally on ROM.
  fmt.Println(addressBus)

  cpu := &c64.Cpu{Bus: addressBus}
  cpu.Reset()
  fmt.Println(cpu)

}
