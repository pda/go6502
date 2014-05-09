/*
Emulates a 128x32 pixel monochrome OLED display with SPI interface.
Exposes the display as a dynamically generated PNG available from an HTTP URL.

Physical hardware example: https://www.adafruit.com/products/661
*/
package ssd1306

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const (
	// Filename where SSD1306 will write its display upon exit.
	DumpFilename = "ssd1306.png"

	// HttpUrl where screen data will be available.
	HttpUrl = "http://localhost:1234/ssd1306.png"
)

// Ssd1306 implements ParallelPeripheral interface for Via6522.

type Ssd1306 struct {
	lastClock   bool
	inputBuffer byte
	inputIndex  uint8
	img         *image.Gray
	imgPixel    uint32
}

func NewSsd1306() *Ssd1306 {
	s := Ssd1306{}
	s.inputIndex = 7 // MSB-first, decrementing index.
	s.img = image.NewGray(image.Rect(0, 0, 128, 32))
	s.serveHttp()
	return &s
}

// TODO: configurable lines
const (
	mosiMask  = 1 << 0
	clockMask = 1 << 1
	dcMask    = 1 << 2
	resetMask = 1 << 3
)

// Notify expects a byte representing the updated status of the parallel port
// that the display is connected to.
// Four bits are considered: MOSI, CLK, D/C, RST.
// The other four bits are ignored.
func (s *Ssd1306) Notify(data byte) {

	mosi := data&mosiMask > 0
	clock := data&clockMask > 0

	if clock && !s.lastClock {
		// rising clock
		s.lastClock = clock
		if mosi {
			s.inputBuffer |= (1 << s.inputIndex)
		}
		if s.inputIndex == 0 {
			//fmt.Printf("Ssd1306: 0x%02X 0b%08b\n", s.inputBuffer, s.inputBuffer)
			s.inputIndex = 7
			s.inputBuffer = 0x00
		} else {
			s.inputIndex--
		}

		if data&dcMask > 0 {
			//fmt.Printf("dat:%08b ", data)

			x := (s.imgPixel / 8) % 128
			y := (7 - s.imgPixel%8) + 8*(s.imgPixel/1024)

			//fmt.Printf("x:% 3d,y:% 3d ", x, y)
			if mosi {
				s.img.Set(int(x), int(y), color.White)
			} else {
				s.img.Set(int(x), int(y), color.Black)
			}

			s.imgPixel++
			s.imgPixel %= (128 * 64)
		}
	}

	if !clock && s.lastClock {
		// falling clock
		s.lastClock = clock
	}

}

func (s *Ssd1306) Close() {
	fmt.Println("Writing SSD1306 screen to", DumpFilename)
	writer, err := os.Create(DumpFilename)
	if err != nil {
		panic(err)
	}
	_ = png.Encode(writer, s.img)
}

func (s *Ssd1306) httpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "image/png")

	if refreshString := r.URL.Query().Get("refresh"); len(refreshString) > 0 {
		refresh, err := strconv.ParseFloat(refreshString, 64)
		if err == nil {
			w.Header().Add("Refresh", fmt.Sprintf("%0.2f", refresh))
		}
	}

	png.Encode(w, s.img)
}

func (s *Ssd1306) serveHttp() {
	url, err := url.Parse(HttpUrl)
	if err != nil {
		panic(err)
	}

	srv := &http.Server{
		Addr:    url.Host,
		Handler: http.HandlerFunc(s.httpHandler),
	}
	fmt.Printf("Ssd1306 output at %s\n", url)
	go srv.ListenAndServe()
}
