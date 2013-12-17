package main

import(
  "fmt"
  "go6502"
)

func main() {

  kernal, err := go6502.RomFromFile("rom/kernal.rom")
  if err != nil {
    panic(err)
  }

  addressBus := &go6502.Bus{Ram: &go6502.Ram{}, Kernal: kernal}
  fmt.Println(addressBus)

  cpu := &go6502.Cpu{Bus: addressBus}
  cpu.Reset()
  fmt.Println(cpu)

  for i := 0; i < 32; i++ {
    fmt.Println("\n--- Step", i)
    cpu.Step()
  }

}
