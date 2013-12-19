package go6502

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
// 0x03: DDRB; Data Direction Register A
//
// External interface relevant to peripheral ports:
// PORTA: 8-bit independently bidirectional data to peripheral.
// PORTB: 8-bit independently bidirectional data to peripheral.
// DATA: 8-bit bidirectional data to microprocessor.
// RS: 4-bit register select.
// CA: 2-bit control lines for PORTA.
// CB: 2-bit control lines for PORTB.


