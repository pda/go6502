/*
	Package SD emulates an SD/MMC card.
*/
package sd

import (
	"fmt"
	"io/ioutil"
)

type SdCard struct {
	data  []byte
	size  int
	state *sdState
	spi   *spi
	PinMap
}

// PinMap associates SD card lines with parallel port pin numbers (0..7).
type PinMap struct {
	Sclk uint
	Mosi uint
	Miso uint
	Ss   uint
}

func (p PinMap) PinMask() byte {
	return 1<<p.Sclk | 1<<p.Mosi | 1<<p.Miso | 1<<p.Ss
}

// SdFromFile creates a new SdCard based on the contents of a file.
func NewSdCard(pm PinMap) (sd *SdCard, err error) {
	sd = &SdCard{
		PinMap: pm,
		state:  newSdState(),
		spi:    newSpi(pm),
	}

	// two busy bytes, then ready.
	sd.state.queueMiso(0x00, 0x00, 0xFF)

	return
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
			miso := sd.state.shiftMiso()

			sd.state.consumeByte(mosi)
			sd.spi.QueueMiso(miso)

			fmt.Printf("SD MOSI $%02X %08b <-> $%02X %08b MISO\n",
				mosi, mosi, miso, miso)
		}
	}
}
