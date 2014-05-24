/*
	Package SD emulates an SD/MMC card.
*/
package sd

import (
	"fmt"
	"io/ioutil"

	"github.com/pda/go6502/spi"
)

type SdCard struct {
	data  []byte
	size  int
	state *sdState
	spi   *spi.Slave
}

// SdFromFile creates a new SdCard based on the contents of a file.
func NewSdCard(pm spi.PinMap) (sd *SdCard, err error) {
	sd = &SdCard{
		state: newSdState(),
		spi:   spi.NewSlave(pm),
	}

	// two busy bytes, then ready.
	sd.state.queueMisoBytes(0x00, 0x00, 0xFF)

	return
}

func (sd *SdCard) PinMask() byte {
	return sd.spi.PinMask()
}

// LoadFile is equivalent to inserting an SD card.
func (sd *SdCard) LoadFile(path string) (err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	sd.size = len(data)
	sd.data = data
	return
}

func (sd *SdCard) Shutdown() {
}

func (sd *SdCard) Read() byte {
	return sd.spi.Read()
}

// Write takes an updated parallel port state.
func (sd *SdCard) Write(data byte) {
	if sd.spi.Write(data) {
		if sd.spi.Done {
			mosi := sd.spi.Mosi
			fmt.Printf("SD MOSI $%02X %08b <-> $%02X %08b MISO\n",
				mosi, mosi, sd.spi.Miso, sd.spi.Miso)

			// consume the byte read, queue miso bytes internally
			sd.state.consumeByte(mosi)
			// dequeues one miso byte, or a default byte if queue empty.
			sd.spi.QueueMisoBits(sd.state.shiftMiso())
		}
	}
}
