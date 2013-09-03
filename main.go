package main

import(
  "fmt"
  "c64"
)

func main() {

  dataPort := &c64.DataPort{}
  ram := &c64.Ram{}
  kernal := c64.RomFromFile("rom/kernal.rom")

  addressBus := &c64.Bus{DataPort: dataPort, Ram: ram, Kernal: kernal}
  fmt.Println(addressBus)

  cpu := &c64.Cpu{Bus: addressBus}
  cpu.Reset()
  fmt.Println(cpu)

  for i := 0;; i++ {
    fmt.Println("\n--- Step", i)
    cpu.Step()
  }

}
