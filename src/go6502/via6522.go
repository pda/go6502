package go6502

import (
	"fmt"
)

// A partial emulation of MOS Technology 6522, or the modern WDC65C22
// incarnation.  This is a Versatile Interface Adapter (VIA) I/O controller
// designed for use with the 6502 microprocessor.
//
// The 4-bit RS (register select) is exposed as 16 bytes of address-space.  The
// processor chooses the register using four bits of the 16-bit address bus and
// reads/writes using the 8-bit data bus.
//
// Original 6522: http://en.wikipedia.org/wiki/MOS_Technology_6522
//
// WCD 65C22: http://www.westerndesigncenter.com/wdc/w65c22-chip.cfm
// Data sheet: http://www.westerndesigncenter.com/wdc/documentation/w65c22.pdf
//
// Peripheral ports
// ----------------
//
// The W65C22 includes functions for programmed control of two peripheral ports
// (Ports A and B). Two program controlled 8-bit bidirectional peripheral I/O
// ports allow direct interfacing between the microprocessor and selected
// peripheral units. Each port has input data latching capability. Two
// programmable Data Direction Registers (A and B) allow selection of data
// direction (input or output) on an individual line basis.
//
// RS registers relevant to peripheral ports:
// (a register is selected by setting an address to the 4-bit RS lines)
// 0x00: ORB/IRB; write: Output Register B, read: Input Register "B".
// 0x01: ORA/IRA; write: Output Register A, read: Input Register "A".
// 0x02: DDRB; Data Direction Register B
// 0x03: DDRA; Data Direction Register A
//
// External interface relevant to peripheral ports:
// PORTA: 8-bit independently bidirectional data to peripheral.
// PORTB: 8-bit independently bidirectional data to peripheral.
// DATA: 8-bit bidirectional data to microprocessor.
// RS: 4-bit register select.
// CA: 2-bit control lines for PORTA.
// CB: 2-bit control lines for PORTB.

const (
	VIA_ORB = 0x0
	VIA_IRB = 0x0
	VIA_ORA = 0x1
	VIA_IRA = 0x1

	VIA_DDRB = 0x2
	VIA_DDRA = 0x3
)

/**
 * Memory interface implementation.
 */

type Via6522 struct {
	ora  byte // output register port A
	orb  byte // output register port B
	ira  byte // input register port A
	irb  byte // input register port B
	ddra byte // data direction port A
	ddrb byte // data direction port B
}

func (via *Via6522) dumpDataDirectionRegisters() {
	fmt.Printf("VIA DDRA:%08b DDRB:%08b\n", via.ddra, via.ddrb)
}

func (via *Via6522) dumpDataRegisters() {
	fmt.Printf("VIA ORA:0x%02X ORB:0x%02X IRA:0x%02X IRB:0x%02X\n", via.ora, via.orb, via.ira, via.irb)
}

// Read the register specified by the given 4-bit address (0x00..0x0F).
func (via *Via6522) Read(a address) byte {
	switch a {
	default:
		panic(fmt.Sprintf("read from 0x%X not handled by Via6522", a))
	case 0x0:
		via.dumpDataRegisters()
		return via.irb
	case 0x1:
		via.dumpDataRegisters()
		return via.ira
	case 0x2:
		return via.ddra
	case 0x3:
		return via.ddra
	}
}

func (via *Via6522) Size() int {
	return 16 // 4-bit RS exposes 16 byte address space.
}

func (via *Via6522) String() string {
	return "VIA6522"
}

// Write to register specified by the given 4-bit address (0x00..0x0F).
func (via *Via6522) Write(a address, data byte) {
	switch a {
	default:
		panic(fmt.Sprintf("write to 0x%X not handled by Via6522", a))
	case 0x0:
		via.orb = data
		via.dumpDataRegisters()
	case 0x1:
		via.ora = data
		via.dumpDataRegisters()
	case 0x2:
		via.ddrb = data
		via.dumpDataDirectionRegisters()
	case 0x3:
		via.ddra = data
		via.dumpDataDirectionRegisters()
	}
}
