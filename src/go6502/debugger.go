package go6502

import (
	"fmt"
	"github.com/peterh/liner"
	"os"
	"strings"
)

const (
	DEBUG_CMD_NONE = iota
	DEBUG_CMD_INVALID
	DEBUG_CMD_EXIT
	DEBUG_CMD_RUN
	DEBUG_CMD_STEP
)

type Debugger struct {
	cpu         *Cpu
	liner       *liner.State
	lastCommand *DebuggerCommand
	run         bool
}

type DebuggerCommand struct {
	id        int
	input     string
	arguments []string
}

func NewDebugger(cpu *Cpu) *Debugger {
	d := &Debugger{liner: liner.NewLiner(), cpu: cpu}
	return d
}

func (d *Debugger) BeforeExecute(iop *Iop) {

	if d.run {
		return
	}

	fmt.Println(d.cpu)

	var (
		cmd *DebuggerCommand
		err error
	)

	for cmd == nil && err == nil {
		cmd, err = d.getCommand()
	}

	switch cmd.id {
	case DEBUG_CMD_EXIT:
		os.Exit(0)
	case DEBUG_CMD_RUN:
		d.run = true
	case DEBUG_CMD_STEP:
		// pass
	case DEBUG_CMD_NONE:
	default:
		panic("Invalid command")
	}
}

func (d *Debugger) getCommand() (*DebuggerCommand, error) {
	var (
		id        int
		cmdString string
		arguments []string
		cmd       *DebuggerCommand
	)

	input, err := d.readInput()
	if err != nil {
		return nil, err
	}

	fields := strings.Fields(input)

	if len(fields) >= 1 {
		cmdString = fields[0]
	}
	if len(fields) >= 2 {
		arguments = fields[1:]
	}

	switch cmdString {
	case "":
		id = DEBUG_CMD_NONE
	case "exit", "quit":
		id = DEBUG_CMD_EXIT
	case "run", "r":
		id = DEBUG_CMD_RUN
	case "step", "st", "s":
		id = DEBUG_CMD_STEP
	default:
		fmt.Println("Invalid command.")
		id = DEBUG_CMD_INVALID
	}

	if id == DEBUG_CMD_NONE {
		cmd = d.lastCommand
	} else {
		cmd = &DebuggerCommand{id, input, arguments}
		d.lastCommand = cmd
	}

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
