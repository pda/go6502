package main

import(
  "fmt"
  "go6502"
  "os"
  "os/signal"
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

  // Dispatch CPU in a goroutine.
  go func () {
    i := 0
    for {
      fmt.Println("\n--- Step", i)
      cpu.Step()
      i++
    }
  }()

  sigChan := make(chan os.Signal, 1)
  signal.Notify(sigChan, os.Interrupt)
  sig := <-sigChan
  fmt.Println("\nGot signal:", sig)

  fmt.Println("Dumping RAM into core file")
  ram.Dump("core")

  os.Exit(1)
}
