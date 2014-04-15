package go6502

import (
	"fmt"
	"github.com/peterh/liner"
	"os"
)

const (
	DEBUG_CMD_NONE = iota
	DEBUG_CMD_STEP
	DEBUG_CMD_EXIT
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
		cmd int
		err error
	)

	for cmd == 0 && err == nil {
		cmd, err = d.getCommand()
	}

	switch cmd {
	case DEBUG_CMD_STEP:
		return
	case DEBUG_CMD_EXIT:
		os.Exit(0)
	default:
		panic("Invalid command")
	}
}

func (d *Debugger) getCommand() (int, error) {
	input, err := d.readInput()
	if err != nil {
		return DEBUG_CMD_NONE, err
	}
	if input == "" {
		input = d.lastInput
	}

	var cmd int
	switch input {
	case "":
		cmd = DEBUG_CMD_NONE
	case "step", "st", "s":
		cmd = DEBUG_CMD_STEP
	case "exit", "quit":
		cmd = DEBUG_CMD_EXIT
	default:
		fmt.Println("Invalid command.")
		cmd = DEBUG_CMD_NONE
	}

	d.lastInput = input

	return cmd, nil
}

func (d *Debugger) readInput() (string, error) {
	input, err := d.liner.Prompt(d.prompt())
	if err != nil {
		return "", err
	}
	d.liner.AppendHistory(input)
	return input, nil
}

func (d *Debugger) prompt() string {
	return fmt.Sprintf("$%04X> ", d.cpu.pc)
}
