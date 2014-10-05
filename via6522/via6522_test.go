package via6522

import (
	"fmt"
	"testing"
)

func via() *Via6522 {
	return NewVia6522(Options{})
}

func TestViaReadAndWriteToDataDirectionRegisters(t *testing.T) {
	via := via()
	via.Write(0x0002, 0x12)
	via.Write(0x0003, 0x34)
	b := via.Read(0x0002)
	a := via.Read(0x0003)
	if b != 0x12 {
		t.Error(fmt.Errorf("DDRB read back $%02X instead of $%02X", a, 0x12))
	}
	if a != 0x34 {
		t.Error(fmt.Errorf("DDRA read back $%02X instead of $%02X", b, 0x34))
	}
}
