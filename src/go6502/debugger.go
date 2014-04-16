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
	DEBUG_CMD_BREAK_ADDRESS
	DEBUG_CMD_BREAK_INSTRUCTION
	DEBUG_CMD_BREAK_REGISTER
	DEBUG_CMD_EXIT
	DEBUG_CMD_HELP
	DEBUG_CMD_INVALID
	DEBUG_CMD_READ
	DEBUG_CMD_READ16
	DEBUG_CMD_RUN
	DEBUG_CMD_STEP
)

type Debugger struct {
	cpu               *Cpu
	liner             *liner.State
	lastCommand       *DebuggerCommand
	run               bool
	breakAddress      bool
	breakAddressValue address
	breakInstruction  string
	breakRegA         bool
	breakRegAValue    byte
	breakRegX         bool
	breakRegXValue    byte
	breakRegY         bool
	breakRegYValue    byte
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

func (d *Debugger) doBreakpoints(iop *Iop) {
	inName := iop.in.name()

	if inName == d.breakInstruction {
		fmt.Printf("Breakpoint for instruction %s\n", inName)
		d.run = false
	}

	if d.breakAddress && d.cpu.pc == d.breakAddressValue {
		fmt.Printf("Breakpoint for PC address = $%04X\n", d.breakAddressValue)
		d.run = false
	}

	d.checkRegBreakpoint("A", d.breakRegA, d.breakRegAValue, d.cpu.ac)
	d.checkRegBreakpoint("X", d.breakRegX, d.breakRegXValue, d.cpu.x)
	d.checkRegBreakpoint("Y", d.breakRegY, d.breakRegYValue, d.cpu.y)
}

func (d *Debugger) BeforeExecute(iop *Iop) {

	d.doBreakpoints(iop)

	if d.run {
		return
	}

	fmt.Println(d.cpu)
	fmt.Println(iop)

	for !d.commandLoop(iop) {
		// next
	}
}

// Returns true when control is to be released.
func (d *Debugger) commandLoop(iop *Iop) (release bool) {
	var (
		cmd *DebuggerCommand
		err error
	)

	for cmd == nil && err == nil {
		cmd, err = d.getCommand()
	}
	if err != nil {
		panic(err)
	}

	switch cmd.id {
	case DEBUG_CMD_BREAK_ADDRESS:
		d.commandBreakAddress(cmd)
	case DEBUG_CMD_BREAK_INSTRUCTION:
		d.breakInstruction = cmd.arguments[0]
	case DEBUG_CMD_BREAK_REGISTER:
		d.commandBreakRegister(cmd)
	case DEBUG_CMD_EXIT:
		os.Exit(0)
	case DEBUG_CMD_HELP:
		d.commandHelp(cmd)
	case DEBUG_CMD_NONE:
		// pass
	case DEBUG_CMD_READ:
		d.commandRead(cmd)
	case DEBUG_CMD_READ16:
		d.commandRead16(cmd)
	case DEBUG_CMD_RUN:
		d.run = true
		release = true
	case DEBUG_CMD_STEP:
		release = true
	case DEBUG_CMD_INVALID:
		fmt.Println("Invalid command.")
	default:
		panic("Unknown command code.")
	}

	return
}

func (d *Debugger) commandRead(cmd *DebuggerCommand) {
	addr64, err := strconv.ParseUint(cmd.arguments[0], 0, 16)
	if err != nil {
		panic(err)
	}
	addr := address(addr64)
	v := d.cpu.Bus.Read(addr)
	fmt.Printf("$%04X => $%02X 0b%08b %d %q\n", addr, v, v, v, v)
}

func (d *Debugger) commandRead16(cmd *DebuggerCommand) {
	addr64, err := strconv.ParseUint(cmd.arguments[0], 0, 16)
	if err != nil {
		panic(err)
	}
	addrLo := address(addr64)
	addrHi := addrLo + 1
	vLo := d.cpu.Bus.Read(addrLo)
	vHi := d.cpu.Bus.Read(addrHi)
	v := (uint16(vHi) << 8) | uint16(vLo)
	fmt.Printf("$%04X,%04X => $%04X 0b%016b %d\n", addrLo, addrHi, v, v, v)
}

func (d *Debugger) commandHelp(cmd *DebuggerCommand) {
	fmt.Println("")
	fmt.Println("pda6502 debuger")
	fmt.Println("---------------")
	fmt.Println("break-address <addr> (alias: ba) e.g. ba 0x1000")
	fmt.Println("break-instruction <mnemonic> (alias: bi) e.g. bi NOP")
	fmt.Println("break-register <x|y|a> <value> (alias: br) e.g. br x 128")
	fmt.Println("exit (alias: quit, q) Shut down the emulator.")
	fmt.Println("help (alias: h, ?) This help.")
	fmt.Println("read <address> - Read and display 8-bit integer at address.")
	fmt.Println("read16 <address> - Read and display 16-bit integer at address.")
	fmt.Println("run (alias: r) Run continuously until breakpoint.")
	fmt.Println("step (alias: s) Run only the current instruction.")
	fmt.Println("(blank) Repeat the previous command.")
	fmt.Println("")
}

func (d *Debugger) commandBreakAddress(cmd *DebuggerCommand) {
	value64, err := strconv.ParseUint(cmd.arguments[0], 0, 16)
	if err != nil {
		panic(err)
	}
	addr := address(value64)
	d.breakAddress = true
	d.breakAddressValue = addr
}

func (d *Debugger) commandBreakRegister(cmd *DebuggerCommand) {
	regStr := cmd.arguments[0]
	valueStr := cmd.arguments[1]

	var ptr *byte
	switch regStr {
	case "A", "a", "AC", "ac":
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
		cmdString = strings.ToLower(fields[0])
	}
	if len(fields) >= 2 {
		arguments = fields[1:]
	}

	switch cmdString {
	case "":
		id = DEBUG_CMD_NONE
	case "break-address", "break-addr", "ba":
		id = DEBUG_CMD_BREAK_ADDRESS
	case "break-instruction", "bi":
		id = DEBUG_CMD_BREAK_INSTRUCTION
	case "break-register", "break-reg", "br":
		id = DEBUG_CMD_BREAK_REGISTER
	case "exit", "quit", "q":
		id = DEBUG_CMD_EXIT
	case "help", "h", "?":
		id = DEBUG_CMD_HELP
	case "read":
		id = DEBUG_CMD_READ
	case "read16":
		id = DEBUG_CMD_READ16
	case "run", "r":
		id = DEBUG_CMD_RUN
	case "step", "st", "s":
		id = DEBUG_CMD_STEP
	default:
		id = DEBUG_CMD_INVALID
	}

	if id == DEBUG_CMD_NONE && d.lastCommand != nil {
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
