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
