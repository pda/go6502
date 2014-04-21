package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/pda/go6502/go6502"
)

const (
	kernalPath = "rom/kernal.rom"
)

func main() {
	os.Exit(mainReturningStatus())
}

func mainReturningStatus() int {

	options := go6502.ParseOptions()

	logFile, err := os.Create(options.LogFile)
	if err != nil {
		panic(err)
	}
	logger := log.New(logFile, "", log.Ltime|log.Lmicroseconds|log.Lshortfile)
	logger.Println("Logger initialized")

	kernal, err := go6502.RomFromFile(kernalPath)
	if err != nil {
		panic(err)
	} else {
		logger.Printf("Loaded %s: %d bytes\n", kernalPath, kernal.Size)
	}

	ram := &go6502.Ram{}

	via := go6502.NewVia6522(logger, options)
	if options.ViaSsd1306 {
		ssd1306 := go6502.NewSsd1306()
		defer ssd1306.Close()
		via.AttachToPortB(ssd1306)
	}

	via.Reset()

	addressBus, _ := go6502.CreateBus(logger)
	addressBus.Attach(ram, "ram", 0x0000)
	addressBus.Attach(via, "VIA", 0xC000)
	addressBus.Attach(kernal, "kernal", 0xE000)
	logger.Println(addressBus)

	exitChan := make(chan int, 0)

	cpu := &go6502.Cpu{Bus: addressBus, ExitChan: exitChan}
	if options.Debug {
		debugger := go6502.NewDebugger(cpu)
		defer debugger.Close()
		debugger.QueueCommands(options.DebugCmds)
		cpu.AttachDebugger(debugger)
	}
	cpu.Reset()
	logger.Println(cpu)

	// Dispatch CPU in a goroutine.
	go func() {
		i := 0
		for {
			logger.Println("\n--- Step", i)
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
