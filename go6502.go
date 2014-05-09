package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/pda/go6502/bus"
	"github.com/pda/go6502/cpu"
	"github.com/pda/go6502/debugger"
	"github.com/pda/go6502/go6502"
	"github.com/pda/go6502/memory"
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

	options := go6502.ParseOptions()

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
		defer ssd1306.Close()
		via.AttachToPortB(ssd1306)
	}

	via.Reset()

	addressBus, _ := bus.CreateBus()
	addressBus.Attach(ram, "ram", 0x0000)
	addressBus.Attach(via, "VIA", 0xC000)
	addressBus.Attach(kernal, "kernal", 0xE000)

	exitChan := make(chan int, 0)

	cpu := &cpu.Cpu{Bus: addressBus, ExitChan: exitChan}
	if options.Debug {
		debugger := debugger.NewDebugger(cpu)
		defer debugger.Close()
		debugger.QueueCommands(options.DebugCmds)
		cpu.AttachMonitor(debugger)
	} else if options.Speedometer {
		speedo := go6502.NewSpeedometer()
		defer speedo.Close()
		cpu.AttachMonitor(speedo)
	}
	cpu.Reset()

	// Dispatch CPU in a goroutine.
	go func() {
		i := 0
		for {
			cpu.Step()
			i++
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
