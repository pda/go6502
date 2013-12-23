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

  ram := &go6502.Ram{}

  via := &go6502.Via6522{}

  addressBus, _ := go6502.CreateBus()
  addressBus.Attach(ram, "ram", 0x0000)
  addressBus.Attach(via, "VIA", 0xD000)
  addressBus.Attach(kernal, "kernal", 0xE000)
  fmt.Println(addressBus)

  cpu := &go6502.Cpu{Bus: addressBus}
  cpu.Reset()
  fmt.Println(cpu)

  for i := 0; i < 32; i++ {
    fmt.Println("\n--- Step", i)
    cpu.Step()
  }

  fmt.Println("Dumping RAM into core file")
  ram.Dump("core")

}
