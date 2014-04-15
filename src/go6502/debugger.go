package go6502

import (
	"fmt"
	"github.com/peterh/liner"
)

const (
	DEBUG_CMD_STEP = iota
)

type Debugger struct {
	cpu       *Cpu
	liner     *liner.State
	lastInput string
}

func NewDebugger(cpu *Cpu) *Debugger {
	d := &Debugger{liner: liner.NewLiner(), cpu: cpu}
	return d
}

func (d *Debugger) Step() {

	fmt.Println(d.cpu)

	var (
		input   string
		command int
	)

	input = d.readInput()
	if input == "" {
		input = d.lastInput
	}

	switch input {
	case "step", "st", "s":
		command = DEBUG_CMD_STEP
	default:
		panic("que?")
	}

	d.lastInput = input

	switch command {
	case DEBUG_CMD_STEP:
	default:
		panic("Invalid command")
	}
}

func (d *Debugger) readInput() string {
	input, err := d.liner.Prompt(d.prompt())
	if err != nil {
		panic(err)
	}
	d.liner.AppendHistory(input)
	return input
}

func (d *Debugger) prompt() string {
	return fmt.Sprintf("$%04X> ", d.cpu.pc)
}
