package spi

import (
	"fmt"
	"testing"
)

func TestPinMapPinMask(t *testing.T) {
	p := PinMap{Sclk: 0, Mosi: 1, Miso: 2, Ss: 7}
	expected := byte(0x87)
	mask := p.PinMask()
	if mask != expected {
		t.Error(fmt.Sprintf("expected PinMask() to be 0b%08b, got 0b%08b", expected, mask))
	}
}
