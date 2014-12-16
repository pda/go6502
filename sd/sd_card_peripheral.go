package sd

import (
	"io/ioutil"

	"github.com/pda/go6502/spi"
)

type SdCardPeripheral struct {
	state *sdState
	spi   *spi.Slave
}

// SdFromFile creates a new SdCardPeripheral based on the contents of a file.
func NewSdCardPeripheral(pm spi.PinMap) (sd *SdCardPeripheral, err error) {
	sd = &SdCardPeripheral{
		state: newSdState(),
		spi:   spi.NewSlave(pm),
	}

	// two busy bytes, then ready.
	sd.state.queueMisoBytes(0x00, 0x00, 0xFF)

	return
}

// LoadFile is equivalent to inserting an SD card.
func (sd *SdCardPeripheral) LoadFile(path string) (err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	sd.state.data = data
	return
}

// via6522.ParallelPeripheral interface

func (sd *SdCardPeripheral) PinMask() byte {
	return sd.spi.PinMask()
}

func (sd *SdCardPeripheral) Read() byte {
	return sd.spi.Read()
}

func (sd *SdCardPeripheral) Shutdown() {
}

// Write takes an updated parallel port state.
func (sd *SdCardPeripheral) Write(data byte) {
	if sd.spi.Write(data) {
		if sd.spi.Done {
			mosi := sd.spi.Mosi
			//fmt.Printf("SD MOSI $%02X %08b <-> $%02X %08b MISO\n",
			//	mosi, mosi, sd.spi.Miso, sd.spi.Miso)

			// consume the byte read, queue miso bytes internally
			sd.state.consumeByte(mosi)
			// dequeues one miso byte, or a default byte if queue empty.
			sd.spi.QueueMisoBits(sd.state.shiftMiso())
		}
	}
}

func (sd *SdCardPeripheral) String() string {
	return "SD card"
}
