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

// nand: output to bits 0,1 then read NAND result on bit 7.

type nand struct {
	value byte
}

func (nand *nand) PinMask() byte {
	return 0x83 // 0b10000011
}

func (nand *nand) Read() byte {
	return nand.value
}

func (nand *nand) Shutdown() {
}

func (nand *nand) Write(in byte) {
	if (in & 0x3) == 0x3 {
		nand.value = 0
	} else {
		nand.value = (1 << 7)
	}
}

func (nand *nand) String() string {
	return "NAND gate test peripheral"
}
