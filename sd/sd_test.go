package sd

import (
	"fmt"
	"testing"

	"github.com/pda/go6502/spi"
)

func TestSdPinMask(t *testing.T) {
	sd, _ := NewSdCardPeripheral(spi.PinMap{Sclk: 4, Mosi: 5, Miso: 6, Ss: 7})
	if sd.PinMask() != 0xF0 {
		t.Error(fmt.Sprintf("0b%08b != 0b%08b", sd.PinMask(), 0xF0))
	}
}
