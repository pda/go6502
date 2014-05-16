/*
	Package SD emulates an SD/MMC card.
*/
package sd

import (
	"fmt"
	"io/ioutil"
)

// TODO: runtime configurable lines
const (
	mosiMask  = 1 << 4
	clockMask = 1 << 6
	csMask    = 1 << 7
)

type spiReader struct {
	clock  bool // the previous clock state
	buffer byte // the byte being built from bits
	index  uint8
}

type SdCard struct {
	data []byte
	size int
	spiReader
}

// SdFromFile creates a new SdCard based on the contents of a file.
func SdFromFile(path string) (sd *SdCard, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	sd = &SdCard{size: len(data), data: data}
	sd.spiReader.index = 7
	return
}

func (sd *SdCard) Shutdown() {
}

func (sd *SdCard) Notify(data byte) {

	cs := data&csMask > 0
	if cs { // high = inactive
		return
	}

	mosi := data&mosiMask > 0
	clock := data&clockMask > 0

	if clock && !sd.clock { // rising clock
		sd.clock = clock
		if mosi {
			sd.buffer |= (1 << sd.index)
		}
		if sd.index == 0 {
			fmt.Printf("SD: 0x%02X 0b%08b\n", sd.buffer, sd.buffer)
			sd.index = 7
			sd.buffer = 0x00
		} else {
			sd.index--
		}
	}

	if !clock && sd.clock {
		// falling clock
		sd.clock = clock
	}
}
