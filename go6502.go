/*
	go6502 emulates the pda6502 computer. This includes the MOS 6502
	processor, memory-mapping address bus, RAM and ROM, MOS 6522 VIA
	controller, SSD1306 OLED display, and perhaps more.

	Read more at https://github.com/pda/go6502 and https://github.com/pda/pda6502
*/
package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/pda/go6502/bus"
	"github.com/pda/go6502/cli"
	"github.com/pda/go6502/cpu"
	"github.com/pda/go6502/debugger"
	"github.com/pda/go6502/ili9340"
	"github.com/pda/go6502/memory"
	"github.com/pda/go6502/sd"
	"github.com/pda/go6502/speedometer"
	"github.com/pda/go6502/spi"
	"github.com/pda/go6502/ssd1306"
	"github.com/pda/go6502/via6522"
)

const (
	kernalPath  = "rom/kernal.rom"
	charRomPath = "rom/char.rom"
)

func main() {
	os.Exit(mainReturningStatus())
}

func mainReturningStatus() int {

	options := cli.ParseFlags()

	// Create addressable devices.

	kernal, err := memory.RomFromFile(kernalPath)
	if err != nil {
		panic(err)
	}

	charRom, err := memory.RomFromFile(charRomPath)
	if err != nil {
		panic(err)
	}

	ram := &memory.Ram{}

	via := via6522.NewVia6522(via6522.Options{
		DumpAscii:  options.ViaDumpAscii,
		DumpBinary: options.ViaDumpBinary,
	})

	if options.Ili9340 {
		ili9340, err := ili9340.NewDisplay(spi.PinMap{
			Sclk: 0,
			Mosi: 6,
			Miso: 7,
			Ss:   5,
		})
		if err != nil {
			panic(err)
		}
		via.AttachToPortB(ili9340)
	}

	if options.ViaSsd1306 {
		ssd1306 := ssd1306.NewSsd1306()
		via.AttachToPortA(ssd1306)
	}

	if len(options.SdCard) > 0 {
		sd, err := sd.NewSdCard(spi.PinMap{
			Sclk: 0,
			Mosi: 6,
			Miso: 7,
			Ss:   4,
		})
		if err != nil {
			panic(err)
		}
		err = sd.LoadFile(options.SdCard)
		if err != nil {
			panic(err)
		}
		via.AttachToPortB(sd)
	}

	via.Reset()

	// Attach devices to address bus.

	addressBus, _ := bus.CreateBus()
	addressBus.Attach(ram, "ram", 0x0000)
	addressBus.Attach(via, "VIA", 0x9000)
	addressBus.Attach(charRom, "char", 0xB000)
	addressBus.Attach(kernal, "kernal", 0xF000)

	exitChan := make(chan int, 0)

	cpu := &cpu.Cpu{Bus: addressBus, ExitChan: exitChan}
	defer cpu.Shutdown()
	if options.Debug {
		debugger := debugger.NewDebugger(cpu, options.DebugSymbolFile)
		debugger.QueueCommands(options.DebugCmds)
		cpu.AttachMonitor(debugger)
	} else if options.Speedometer {
		speedo := speedometer.NewSpeedometer()
		cpu.AttachMonitor(speedo)
	}
	cpu.Reset()

	// Dispatch CPU in a goroutine.
	go func() {
		for {
			cpu.Step()
		}
	}()

	var (
		sig        os.Signal
		exitStatus int
	)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	select {
	case exitStatus = <-exitChan:
		// pass
	case sig = <-sigChan:
		fmt.Println("\nGot signal:", sig)
		exitStatus = 1
	}

	fmt.Println(cpu)
	fmt.Println("Dumping RAM into core file")
	ram.Dump("core")

	return exitStatus
}
