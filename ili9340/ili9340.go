/*
Emulates 240x320 TFT color display with SPI interface.
http://www.adafruit.com/products/1480
*/
package ili9340

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/pda/go6502/spi"
)

const (
	dcMask uint8 = 1 << 2
)

const (
	stateUnknown = iota
	stateRamWrite
)

const (
	cmdRamWrite = 0x2C
)

const (
	width        = 320
	height       = 240
	dumpFilename = "ili9340.png"
)

type Display struct {
	spi       *spi.Slave
	dataMode  bool
	state     uint
	img       *image.RGBA
	nextX     uint16
	nextY     uint16
	pixel     uint16 // accumulator for current pixel
	pixelByte uint8  // next byte offset into the pixel (0 or 1)
}

func NewDisplay(pm spi.PinMap) (display *Display, err error) {
	display = &Display{
		spi: spi.NewSlave(pm),
		img: image.NewRGBA(image.Rect(0, 0, width, height)),
	}
	return
}

func (d *Display) PinMask() byte {
	return d.spi.PinMask() | dcMask
}

func (d *Display) Read() byte {
	return d.spi.Read()
}

func (d *Display) Write(b byte) {
	if b&dcMask == 0 && d.dataMode {
		d.dataMode = false
	} else if b&dcMask != 0 && !d.dataMode {
		d.dataMode = true
	}

	d.spi.Write(b)
	if d.spi.Done {
		d.acceptByte(d.spi.Mosi)
	}

}

func (d *Display) String() string {
	return "ILI9340"
}

func (d *Display) Shutdown() {
	d.writeImage()
}

func (d *Display) writeImage() {
	fmt.Println("Writing ILI9340 screen to", dumpFilename)
	writer, err := os.Create(dumpFilename)
	if err != nil {
		panic(err)
	}
	_ = png.Encode(writer, d.img)
}

func (d *Display) acceptByte(b byte) {
	if d.dataMode {
		d.acceptData(b)
	} else {
		d.acceptCommand(b)
	}
}

func (d *Display) acceptCommand(b byte) {
	switch b {
	case cmdRamWrite:
		d.state = stateRamWrite
	default:
		if d.state != stateUnknown {
			d.state = stateUnknown
		}
	}
}

func (d *Display) acceptData(b byte) {
	switch d.state {
	case stateRamWrite:
		d.ramWrite(b)
	}
}

func (d *Display) ramWrite(b byte) {
	if d.pixelByte == 0 {
		d.pixel |= uint16(b) << 8
		d.pixelByte++
	} else {
		d.pixel |= uint16(b)
		d.pixelByte--
		d.pixelWrite(d.pixel)
	}
}

func (d *Display) pixelWrite(p16 uint16) {
	r := uint8((p16 & 0xF800) >> 8) // map high 5-bit to 8-bit color
	g := uint8((p16 & 0x07E0) >> 3) // map mid 6-bit to 8-bit color
	b := uint8((p16 & 0x001F) << 3) // map low 5-bit to 8-bit color

	d.img.SetRGBA(int(d.nextX), int(d.nextY), color.RGBA{r, g, b, 0xFF})

	// move to next pixel
	d.nextX = (d.nextX + 1) % width
	if d.nextX == 0 {
		d.nextY = (d.nextY + 1) % height
	}
}
