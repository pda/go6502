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
	"github.com/pda/go6502/memory"
	"github.com/pda/go6502/sd"
	"github.com/pda/go6502/speedometer"
	"github.com/pda/go6502/ssd1306"
	"github.com/pda/go6502/via6522"
)

const (
	kernalPath = "rom/kernal.rom"
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

	ram := &memory.Ram{}

	via := via6522.NewVia6522(via6522.Options{
		DumpAscii:  options.ViaDumpAscii,
		DumpBinary: options.ViaDumpBinary,
	})

	if options.ViaSsd1306 {
		ssd1306 := ssd1306.NewSsd1306()
		via.AttachToPortB(ssd1306)
	}

	if len(options.SdCard) > 0 {
		sd, err := sd.NewSdCard(sd.PinMap{
			Sclk: 4,
			Mosi: 5,
			Miso: 6,
			Ss:   7,
		})
		if err != nil {
			panic(err)
		}
		err = sd.LoadFile(options.SdCard)
		if err != nil {
			panic(err)
		}
		via.AttachToPortA(sd)
	}

	via.Reset()

	// Attach devices to address bus.

	addressBus, _ := bus.CreateBus()
	addressBus.Attach(ram, "ram", 0x0000)
	addressBus.Attach(via, "VIA", 0xC000)
	addressBus.Attach(kernal, "kernal", 0xE000)

	exitChan := make(chan int, 0)

	cpu := &cpu.Cpu{Bus: addressBus, ExitChan: exitChan}
	defer cpu.Shutdown()
	if options.Debug {
		debugger := debugger.NewDebugger(cpu)
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

	if exitStatus != 0 {
		fmt.Println(cpu)
		fmt.Println("Dumping RAM into core file")
		ram.Dump("core")
	}

	return exitStatus
}
