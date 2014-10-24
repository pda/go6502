/*
Emulates 240x320 TFT color display with SPI interface.
http://www.adafruit.com/products/1480
*/
package ili9340

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
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
	stateColumnAddressSet
	statePageAddressSet
)

const (
	cmdRamWrite         = 0x2C
	cmdColumnAddressSet = 0x2A
	cmdPageAddressSet   = 0x2B
)

const (
	width        = 320
	height       = 240
	dumpFilename = "ili9340.png"
)

type Display struct {
	spi        *spi.Slave
	dataMode   bool
	state      uint
	paramIndex uint8
	paramData  uint32 // accumulator for current parameter
	img        *image.RGBA
	nextX      uint16
	nextY      uint16
	startCol   uint16
	endCol     uint16
	startRow   uint16
	endRow     uint16
}

func NewDisplay(pm spi.PinMap) (display *Display, err error) {
	img := createImage()
	display = &Display{
		spi:      spi.NewSlave(pm),
		img:      img,
		startCol: 0,
		endCol:   width - 1,
		startRow: 0,
		endRow:   height - 1,
	}
	return
}

func createImage() (img *image.RGBA) {
	img = image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.Black}, image.ZP, draw.Src)
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
	d.paramIndex = 0
	switch b {
	case cmdRamWrite:
		d.nextX = d.startCol
		d.nextY = d.startRow
		d.state = stateRamWrite
	case cmdColumnAddressSet:
		d.state = stateColumnAddressSet
	case cmdPageAddressSet:
		d.state = statePageAddressSet
	default:
		if d.state != stateUnknown {
			d.state = stateUnknown
		}
	}
}

func (d *Display) acceptData(b byte) {
	if d.paramIndex == 0 {
		d.paramData = 0
	}
	d.paramData |= (uint32(b) << ((3 - d.paramIndex) * 8))
	d.paramIndex++

	switch d.state {
	case stateRamWrite:
		d.acceptRamWrite(b)
	case stateColumnAddressSet:
		d.acceptColumnAddressByte(b)
	case statePageAddressSet:
		d.acceptPageAddressByte(b)
	}
}

func (d *Display) acceptRamWrite(b byte) {
	if d.paramIndex == 2 {
		d.pixelWrite(uint16(d.paramData >> 16))
		d.paramIndex = 0
	}
}

func (d *Display) pixelWrite(p16 uint16) {
	r := uint8((p16 & 0xF800) >> 8) // map high 5-bit to 8-bit color
	g := uint8((p16 & 0x07E0) >> 3) // map mid 6-bit to 8-bit color
	b := uint8((p16 & 0x001F) << 3) // map low 5-bit to 8-bit color
	c := color.RGBA{r, g, b, 0xFF}

	d.img.SetRGBA(int(d.nextX), int(d.nextY), c)

	if d.nextX == d.endCol {
		d.nextX = d.startCol
		if d.nextY == d.endRow {
			d.nextY = d.startRow
		} else {
			d.nextY++
		}
	} else {
		d.nextX++
	}
}

func (d *Display) acceptColumnAddressByte(b byte) {
	if d.paramIndex == 4 {
		d.startCol = uint16(d.paramData >> 16)
		d.endCol = uint16(d.paramData & 0xFFFF)
		fmt.Printf("ILI9340 column address range %d:%d\n", d.startCol, d.endCol)
	}
}

func (d *Display) acceptPageAddressByte(b byte) {
	if d.paramIndex == 4 {
		d.startRow = uint16(d.paramData >> 16)
		d.endRow = uint16(d.paramData & 0xFFFF)
		fmt.Printf("ILI9340 row address range %d:%d\n", d.startRow, d.endRow)
	}
}
