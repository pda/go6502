package go6502

import (
	"fmt"
	"github.com/peterh/liner"
	"os"
	"strconv"
	"strings"
)

const (
	DEBUG_CMD_NONE = iota
	DEBUG_CMD_BREAK_INSTRUCTION
	DEBUG_CMD_BREAK_REGISTER
	DEBUG_CMD_EXIT
	DEBUG_CMD_INVALID
	DEBUG_CMD_RUN
	DEBUG_CMD_STEP
)

type Debugger struct {
	cpu              *Cpu
	liner            *liner.State
	lastCommand      *DebuggerCommand
	run              bool
	breakInstruction string
	breakRegA        bool
	breakRegAValue   byte
	breakRegX        bool
	breakRegXValue   byte
	breakRegY        bool
	breakRegYValue   byte
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

func (d *Debugger) checkRegBreakpoint(regStr string, on bool, expect byte, actual byte) {
	if on && actual == expect {
		fmt.Printf("Breakpoint for %s = $%02X (%d)\n", regStr, expect, expect)
		d.run = false
	}
}

func (d *Debugger) BeforeExecute(iop *Iop) {

	inName := iop.in.name()

	if inName == d.breakInstruction {
		fmt.Printf("Breakpoint for instruction %s\n", inName)
		d.run = false
	}

	d.checkRegBreakpoint("A", d.breakRegA, d.breakRegAValue, d.cpu.ac)
	d.checkRegBreakpoint("X", d.breakRegX, d.breakRegXValue, d.cpu.x)
	d.checkRegBreakpoint("Y", d.breakRegY, d.breakRegYValue, d.cpu.y)

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
	case DEBUG_CMD_BREAK_INSTRUCTION:
		d.breakInstruction = cmd.arguments[0]
	case DEBUG_CMD_BREAK_REGISTER:
		d.commandBreakRegister(cmd)
	case DEBUG_CMD_EXIT:
		os.Exit(0)
	case DEBUG_CMD_NONE:
		// pass
	case DEBUG_CMD_RUN:
		d.run = true
	case DEBUG_CMD_STEP:
		// pass
	default:
		panic("Invalid command")
	}
}

func (d *Debugger) commandBreakRegister(cmd *DebuggerCommand) {
	regStr := cmd.arguments[0]
	valueStr := cmd.arguments[1]

	var ptr *byte
	switch regStr {
	case "A", "a":
		d.breakRegA = true
		ptr = &d.breakRegAValue
	case "X", "x":
		d.breakRegX = true
		ptr = &d.breakRegXValue
	case "Y", "y":
		d.breakRegY = true
		ptr = &d.breakRegYValue
	default:
		panic(fmt.Errorf("Invalid register for break-register"))
	}

	value64, err := strconv.ParseUint(valueStr, 0, 8)
	if err != nil {
		panic(err)
	}
	value := byte(value64)

	fmt.Printf("Breakpoint set: %s = $%02X (%d)\n", regStr, value, value)

	*ptr = value
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
	case "break-instruction", "bi":
		id = DEBUG_CMD_BREAK_INSTRUCTION
	case "break-register", "break-reg", "br":
		id = DEBUG_CMD_BREAK_REGISTER
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
