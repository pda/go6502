/*
	Package SD emulates an SD/MMC card.
*/
package sd

import (
	"fmt"
	"io/ioutil"
)

type spiState struct {
	clock      bool   // the most recent clock state
	index      uint8  // the bit index of the current byte (mod 8).
	misoBuffer byte   // current byte being sent one bit at a time via Read().
	misoQueue  []byte // data waiting to be sent via Read().
	readByte   byte   // the state of the pins as read by the VIA controller.
	mosiBuffer byte   // the byte being built from bits
}

type SdCard struct {
	data []byte
	size int
	spiState
	PinMap

	maskSclk uint8
	maskMosi uint8
	maskMiso uint8
	maskSs   uint8
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
	sd = &SdCard{PinMap: pm}

	sd.maskSclk = 1 << pm.Sclk
	sd.maskMosi = 1 << pm.Mosi
	sd.maskMiso = 1 << pm.Miso
	sd.maskSs = 1 << pm.Ss

	sd.spiState.index = 7
	sd.misoBuffer = 0xFF
	sd.spiState.misoQueue = make([]byte, 0, 1024)
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
	return sd.readByte
}

// Write takes an updated parallel port state.
func (sd *SdCard) Write(data byte) {

	cs := data&sd.maskSs > 0
	if cs { // high = inactive
		return
	}

	mosi := data&sd.maskMosi > 0
	clock := data&sd.maskSclk > 0

	rising := !sd.clock && clock
	falling := sd.clock && !clock
	sd.clock = clock

	// sclk:rise -> miso -> sclk:fall -> mosi -> ...

	if rising {
		if sd.misoBuffer&(1<<sd.index) > 0 {
			sd.readByte = 0x00 | sd.maskMiso
		} else {
			sd.readByte = 0x00
		}
	}

	if falling {
		if mosi {
			sd.mosiBuffer |= (1 << sd.index)
		}

		// after eigth bit
		if sd.index == 0 {
			sd.handleMosiByte()
			sd.handleMisoByte()

			sd.index = 7
			sd.mosiBuffer = 0x00
		} else {
			sd.index--
		}
	}
}

func (sd *SdCard) handleMisoByte() {
	if len(sd.misoQueue) > 0 {
		fmt.Printf("%08b ($%02X) misoQueue -> misoBuffer\n", sd.misoQueue[0], sd.misoQueue[0])
		sd.misoBuffer = sd.misoQueue[0]
		sd.misoQueue = sd.misoQueue[1:len(sd.misoQueue)]
	} else {
		sd.misoBuffer = 0x00 // default to low for empty buffer.
	}
}

func (sd *SdCard) handleMosiByte() {
	data := sd.mosiBuffer
	fmt.Printf("SD MOSI: 0x%02X 0b%08b\n", sd.mosiBuffer, sd.mosiBuffer)
	switch data {
	case 0x40:
		fmt.Printf("Got 0x40; queueing response bytes.\n")
		sd.misoQueue = append(sd.misoQueue, 0xAA, 0xAB, 0xAC, 0xAD)
	}
}
