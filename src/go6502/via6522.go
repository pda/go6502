package go6502

import (
	"fmt"
	"log"
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
// 0x0C: PCR; Peripheral Control Register.
//            0: CA1 control, 1..3: CA2 control
//            4: CB1 control, 5..7: CB2 control.
//
// External interface relevant to peripheral ports:
// PORTA: 8-bit independently bidirectional data to peripheral.
// PORTB: 8-bit independently bidirectional data to peripheral.
// DATA: 8-bit bidirectional data to microprocessor.
// RS: 4-bit register select.
// CA: 2-bit control lines for PORTA.
// CB: 2-bit control lines for PORTB.

// Write handshake control (PORT A as example, PORT B is same for writes):
//   CA2 (output) indicates data has been written to ORA and is ready.
//   CA1 (input) indicates data has been taken.
// Default modes assuming PCR == 0x00:
//   CA2: Input-negative active edge (one of eight options).
//   CA1: negative active edge (one of two options).

const (
	VIA_ORB = 0x0
	VIA_IRB = 0x0
	VIA_ORA = 0x1
	VIA_IRA = 0x1

	VIA_DDRB = 0x2
	VIA_DDRA = 0x3

	// bit-offset into PCR for port A & B
	VIA_PCR_OFFSET_A = 0
	VIA_PCR_OFFSET_B = 4
)

/**
 * Memory interface implementation.
 */

type Via6522 struct {
	// Note: It may be a mistake to consider ORx and IRx separate registers.
	//       If so... fix it?
	ora    byte // output register port A
	orb    byte // output register port B
	ira    byte // input register port A
	irb    byte // input register port B
	ddra   byte // data direction port A
	ddrb   byte // data direction port B
	pcr    byte // peripheral control register
	logger *log.Logger
}

func NewVia6522(l *log.Logger) *Via6522 {
	return &Via6522{logger: l}
}

// CA1 or CB1 1-bit mode for the given port offset (VIA_PCR_OFFSET_x)
func (via *Via6522) control1Mode(portOffset uint8) byte {
	return (via.pcr >> portOffset) & 1
}

// CA2 or CB2 3-bit mode for the given port offset (VIA_PCR_OFFSET_x)
func (via *Via6522) control2Mode(portOffset uint8) byte {
	return (via.pcr >> (portOffset + 1)) & 0x7
}

func (via *Via6522) dumpDataDirectionRegisters() {
	via.logger.Printf("VIA DDRA:%08b DDRB:%08b\n", via.ddra, via.ddrb)
}

func (via *Via6522) dumpDataRegisters() {
	via.logger.Printf("VIA ORA:0x%02X ORB:0x%02X IRA:0x%02X IRB:0x%02X\n", via.ora, via.orb, via.ira, via.irb)
}

func (via *Via6522) handleDataWrite(portOffset uint8) {
	mode := via.control2Mode(portOffset)
	switch mode {
	default:
		panic(fmt.Sprintf("VIA: Unhanded PCR mode 0x%X for write (PCR offset %d)", mode, portOffset))
	case 0x5:
		// pulse output
		fmt.Printf("%c", via.ora) // TODO: something useful
	}
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
	case 0xC:
		return via.pcr
	}
}

// From the datasheet:
// Reset clears all internal registers
// (except T1 and T2 counters and latches, and the SR.)
func (via *Via6522) Reset() {
	via.ora = 0
	via.orb = 0
	via.ira = 0
	via.irb = 0
	via.ddra = 0
	via.ddrb = 0
	via.pcr = 0
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
		via.handleDataWrite(VIA_PCR_OFFSET_B)
		via.dumpDataRegisters()
	case 0x1:
		via.ora = data
		via.handleDataWrite(VIA_PCR_OFFSET_A)
		via.dumpDataRegisters()
	case 0x2:
		via.ddrb = data
		via.dumpDataDirectionRegisters()
	case 0x3:
		via.ddra = data
		via.dumpDataDirectionRegisters()
	case 0xC:
		via.pcr = data
	}
}
