package main

import(
  "fmt"
  "go6502"
)

const(
  kernalPath = "rom/kernal.rom"
)

func main() {

  kernal, err := go6502.RomFromFile(kernalPath)
  if err != nil {
    panic(err)
  } else {
    fmt.Printf("Loaded %s: %d bytes\n", kernalPath, kernal.Size)
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
