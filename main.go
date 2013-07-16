package main

import(
  "fmt"
  "c64"
)

func main() {

  kernal := c64.RomFromFile("rom/kernal.rom")

  addressBus := &c64.Bus{Ram: &c64.Ram{}, Kernal: kernal}
  fmt.Println(addressBus)

  cpu := &c64.Cpu{Bus: addressBus}
  cpu.Reset()
  fmt.Println(cpu)

  for i := 0;; i++ {
    fmt.Println("\n--- Step", i)
    cpu.Step()
  }

}
