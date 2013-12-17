package main

import(
  "fmt"
  "go6502"
)

func main() {

  kernal := go6502.RomFromFile("rom/kernal.rom")

  addressBus := &go6502.Bus{Ram: &go6502.Ram{}, Kernal: kernal}
  fmt.Println(addressBus)

  cpu := &go6502.Cpu{Bus: addressBus}
  cpu.Reset()
  fmt.Println(cpu)

  for i := 0;; i++ {
    fmt.Println("\n--- Step", i)
    cpu.Step()
  }

}
