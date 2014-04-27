/*
	Package debugger provides an interactive stepping debugger for go6502 with
	breakpoints on instruction type, register values and memory location.

	Example

	An example interactive debugging session:

		$ go run go6502.go --via-ssd1306 --debug
		CPU PC:0xF31F AC:0x00 X:0x00 Y:0x00 SP:0x00 SR:--_b-i--
		Next: SEI implied
		$F31F> step
		CPU PC:0xF320 AC:0x00 X:0x00 Y:0x00 SP:0x00 SR:--_b----
		Next: LDX immediate $FF
		$F320> break-register X $FF
		Breakpoint set: X = $FF (255)
		$F320> run
		Breakpoint for X = $FF (255)
		CPU PC:0xF322 AC:0x00 X:0xFF Y:0x00 SP:0x00 SR:n-_b----
		Next: TXS implied
		$F322> step
		Breakpoint for X = $FF (255)
		CPU PC:0xF323 AC:0x00 X:0xFF Y:0x00 SP:0xFF SR:n-_b----
		Next: CLI implied
		$F323>
		Breakpoint for X = $FF (255)
		CPU PC:0xF324 AC:0x00 X:0xFF Y:0x00 SP:0xFF SR:n-_b-i--
		Next: CLD implied
		$F324>
		Breakpoint for X = $FF (255)
		CPU PC:0xF325 AC:0x00 X:0xFF Y:0x00 SP:0xFF SR:n-_b-i--
		Next: JMP absolute $F07B
		$F325> break-instruction nop
		$F325> r
		Breakpoint for X = $FF (255)
		CPU PC:0xF07B AC:0x00 X:0xFF Y:0x00 SP:0xFF SR:n-_b-i--
		Next: LDA immediate $00
		$F07B> q
*/
package debugger

/**
 * TODO:
 * -  `step n` e.g. `step 100` to step 100 instructions.
 * -  Read and write CLI history file.
 * -  Ability to label addresses, persist+load.
 * -  Tab completion.
 * -  Command argument validation.
 */

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pda/go6502/go6502"
	"github.com/peterh/liner"
)

const (
	debugCmdNone = iota
	debugCmdBreakAddress
	debugCmdBreakInstruction
	debugCmdBreakRegister
	debugCmdExit
	debugCmdHelp
	debugCmdInvalid
	debugCmdRead
	debugCmdRead16
	debugCmdRun
	debugCmdStep
)

type Debugger struct {
	inputQueue        []string
	cpu               *go6502.Cpu
	liner             *liner.State
	lastCommand       *DebuggerCommand
	run               bool
	breakAddress      bool
	breakAddressValue go6502.Address
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

// NewDebugger creates a debugger.
// Be sure to defer a call to Debugger.Close() afterwards, or your terminal
// will be left in a broken state.
func NewDebugger(cpu *go6502.Cpu) *Debugger {
	d := &Debugger{liner: liner.NewLiner(), cpu: cpu}
	return d
}

// Close the debugger session, including resetting the terminal to its previous
// state.
func (d *Debugger) Close() {
	d.liner.Close()
}

// Queue a list of commands to be executed at the next prompt(s).
// This is useful for accepting a list of commands as a CLI parameter.
func (d *Debugger) QueueCommands(cmds []string) {
	d.inputQueue = append(d.inputQueue, cmds...)
}

func (d *Debugger) checkRegBreakpoint(regStr string, on bool, expect byte, actual byte) {
	if on && actual == expect {
		fmt.Printf("Breakpoint for %s = $%02X (%d)\n", regStr, expect, expect)
		d.run = false
	}
}

func (d *Debugger) doBreakpoints(in *go6502.Instruction) {
	inName := in.Name()

	if inName == d.breakInstruction {
		fmt.Printf("Breakpoint for instruction %s\n", inName)
		d.run = false
	}

	if d.breakAddress && d.cpu.PC == d.breakAddressValue {
		fmt.Printf("Breakpoint for PC address = $%04X\n", d.breakAddressValue)
		d.run = false
	}

	d.checkRegBreakpoint("A", d.breakRegA, d.breakRegAValue, d.cpu.AC)
	d.checkRegBreakpoint("X", d.breakRegX, d.breakRegXValue, d.cpu.X)
	d.checkRegBreakpoint("Y", d.breakRegY, d.breakRegYValue, d.cpu.Y)
}

// BeforeExecute receives each go6502.Instruction just before the program
// counter is incremented and the instruction executed.
func (d *Debugger) BeforeExecute(in *go6502.Instruction) {

	d.doBreakpoints(in)

	if d.run {
		return
	}

	fmt.Println(d.cpu)
	fmt.Println("Next:", in)

	for !d.commandLoop(in) {
		// next
	}
}

// Returns true when control is to be released.
func (d *Debugger) commandLoop(in *go6502.Instruction) (release bool) {
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
	case debugCmdBreakAddress:
		d.commandBreakAddress(cmd)
	case debugCmdBreakInstruction:
		d.breakInstruction = strings.ToUpper(cmd.arguments[0])
	case debugCmdBreakRegister:
		d.commandBreakRegister(cmd)
	case debugCmdExit:
		d.cpu.ExitChan <- 0
	case debugCmdHelp:
		d.commandHelp(cmd)
	case debugCmdNone:
		// pass
	case debugCmdRead:
		d.commandRead(cmd)
	case debugCmdRead16:
		d.commandRead16(cmd)
	case debugCmdRun:
		d.run = true
		release = true
	case debugCmdStep:
		release = true
	case debugCmdInvalid:
		fmt.Println("Invalid command.")
	default:
		panic("Unknown command code.")
	}

	return
}

func (d *Debugger) commandRead(cmd *DebuggerCommand) {
	addr64, err := d.parseUint(cmd.arguments[0], 16)
	if err != nil {
		panic(err)
	}
	addr := go6502.Address(addr64)
	v := d.cpu.Bus.Read(addr)
	fmt.Printf("$%04X => $%02X 0b%08b %d %q\n", addr, v, v, v, v)
}

func (d *Debugger) commandRead16(cmd *DebuggerCommand) {
	addr64, err := d.parseUint(cmd.arguments[0], 16)
	if err != nil {
		panic(err)
	}
	addrLo := go6502.Address(addr64)
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
	fmt.Println("Hex input formats: 0x1234 $1234")
	fmt.Println("Commands expecting uint16 treat . as current address (PC).")
}

func (d *Debugger) commandBreakAddress(cmd *DebuggerCommand) {
	value64, err := d.parseUint(cmd.arguments[0], 16)
	if err != nil {
		panic(err)
	}
	addr := go6502.Address(value64)
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

	value64, err := d.parseUint(valueStr, 8)
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
		input     string
		err       error
	)

	if len(d.inputQueue) > 0 {
		input = d.inputQueue[0]
		d.inputQueue = d.inputQueue[1:]
		fmt.Printf("%s%s\n", d.prompt(), input)
	} else {
		input, err = d.readInput()
		if err != nil {
			return nil, err
		}
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
		id = debugCmdNone
	case "break-address", "break-addr", "ba":
		id = debugCmdBreakAddress
	case "break-instruction", "bi":
		id = debugCmdBreakInstruction
	case "break-register", "break-reg", "br":
		id = debugCmdBreakRegister
	case "exit", "quit", "q":
		id = debugCmdExit
	case "help", "h", "?":
		id = debugCmdHelp
	case "read":
		id = debugCmdRead
	case "read16":
		id = debugCmdRead16
	case "run", "r":
		id = debugCmdRun
	case "step", "st", "s":
		id = debugCmdStep
	default:
		id = debugCmdInvalid
	}

	if id == debugCmdNone && d.lastCommand != nil {
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
	return fmt.Sprintf("$%04X> ", d.cpu.PC)
}

func (d *Debugger) parseUint(s string, bits int) (uint64, error) {
	if s == "." && bits == 16 {
		return uint64(d.cpu.PC), nil
	}
	s = strings.Replace(s, "$", "0x", 1)
	return strconv.ParseUint(s, 0, bits)
}
