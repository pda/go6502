package via6522

import (
	"fmt"
	"testing"
)

const (
	iorb = 0x0
	iora = 0x1
	ddrb = 0x2
	ddra = 0x3
)

func via() *Via6522 {
	return NewVia6522(Options{})
}

func TestViaReadAndWriteToDataDirectionRegisters(t *testing.T) {
	via := via()
	via.Write(ddrb, 0x12)
	via.Write(ddra, 0x34)
	b := via.Read(ddrb)
	a := via.Read(ddra)
	if b != 0x12 {
		t.Error(fmt.Errorf("DDRB read back $%02X instead of $%02X", a, 0x12))
	}
	if a != 0x34 {
		t.Error(fmt.Errorf("DDRA read back $%02X instead of $%02X", b, 0x34))
	}
}

func TestViaReadAndWriteInputOutputPortRegisters(t *testing.T) {
	via := via()
	via.Write(ddrb, 0xAA)
	via.Write(iorb, 0xDE)
	b := via.Read(iorb)
	expected := uint8(0xAA & 0xDE) // assumes zero values in input register
	if b != expected {
		t.Error(fmt.Errorf("$%02X != $%02X", b, expected))
	}
}

func TestViaReadAndWriteToNandPeripheral(t *testing.T) {
	via := via()
	via.AttachToPortB(&nand{})
	via.Write(ddrb, 0x03) // output to bits 0,1

	// read back the output pins, and the NAND result at bit 7.
	assertPortBWriteThenRead(t, via, 0x0, 0x0|1<<7)
	assertPortBWriteThenRead(t, via, 0x1, 0x1|1<<7)
	assertPortBWriteThenRead(t, via, 0x2, 0x2|1<<7)
	assertPortBWriteThenRead(t, via, 0x3, 0x3|0)

	// input pins that were written to are not read back.
	assertPortBWriteThenRead(t, via, 0xFF, 0x3|0)
	assertPortBWriteThenRead(t, via, 0xFE, 0x82)
}

func assertPortBWriteThenRead(t *testing.T, via *Via6522, write, expect byte) {
	via.Write(iorb, write)
	result := via.Read(iorb)
	if result != expect {
		t.Error(fmt.Errorf("wrote 0b%08b, expected 0b%08b, got 0b%08b", write, expect, result))
	}
}

func TestPinMaskBlocksWrites(t *testing.T) {
	nand := &nand{}
	via := via()
	via.AttachToPortB(nand)
	via.Write(ddrb, 0xFF) // all output
	via.Write(iorb, 0xFF) // write all bits
	if nand.value != nand.PinMask() {
		t.Error(fmt.Errorf("peripheral received 0b%08b despite 0b%08b pinmask", nand.value, nand.PinMask()))
	}
}

func TestDdrBlocksWrites(t *testing.T) {
	nand := &nand{}
	via := via()
	via.AttachToPortB(nand)
	via.Write(ddrb, 0x02)
	via.Write(iorb, 0xFF)
	if nand.value&^0x02 != 0 {
		t.Error(fmt.Errorf("peripheral received 0b%08b despite 0b%08b DDR", nand.value, 0xAA))
	}
}

func TestWriteToMultipleOverlappingPeripherals(t *testing.T) {
	one := &flipflop{pinmask: 0xF8} // 0b11111000
	two := &flipflop{pinmask: 0x1F} // 0b00011111
	via := via()
	via.AttachToPortA(one)
	via.AttachToPortA(two)
	via.Write(ddra, 0xFF) // all output
	via.Write(iora, 0xAA)

	expectedOne := 0xAA & one.PinMask()
	expectedTwo := 0xAA & two.PinMask()

	if one.value != expectedOne {
		t.Error(fmt.Errorf("one: wrote 0b%08b, expected 0b%08b, got 0b%08b", 0xAA, expectedOne, one.value))
	}

	if two.value != expectedTwo {
		t.Error(fmt.Errorf("two: wrote 0b%08b, expected 0b%08b, got 0b%08b", 0xAA, expectedTwo, two.value))
	}
}

func TestReadFromMultipleOverlappingPeripherals(t *testing.T) {
	one := &flipflop{pinmask: 0xF8, value: 0xAA} // pinmask: 11111000, value: 10101010
	two := &flipflop{pinmask: 0x1F, value: 0xF0} // pinmask: 00011111, value: 11110000
	via := via()
	via.AttachToPortA(one)
	via.AttachToPortA(two)
	via.Write(ddra, 0x00) // all input
	result := via.Read(iora)
	expected := byte(0xB8) // 0b10111000
	if result != expected {
		t.Error(fmt.Errorf("one: expected 0b%08b, got 0b%08b", expected, result))
	}
}

// ---------------------------------------
// Test ParallelPeripheral implementations

// nand: output to bits 0,1 then read NAND result on bit 7.

type nand struct {
	value byte
}

func (nand *nand) PinMask() byte {
	return 0x83 // 0b10000011
}

func (nand *nand) Read() byte {
	if (nand.value & 0x3) == 0x3 {
		return 0
	} else {
		return (1 << 7)
	}
}

func (nand *nand) Shutdown() {
}

func (nand *nand) Write(in byte) {
	nand.value = in
}

func (nand *nand) String() string {
	return "NAND gate test peripheral"
}

// flipflop: simple memory (subject to external DDR, pinmask etc)

type flipflop struct {
	value   byte
	pinmask byte
}

func (ff *flipflop) PinMask() byte {
	return ff.pinmask
}

func (ff *flipflop) Read() byte {
	return ff.value
}

func (ff *flipflop) Shutdown() {
}

func (ff *flipflop) Write(in byte) {
	ff.value = in
}

func (ff *flipflop) String() string {
	return "flipflop test peripheral"
}
