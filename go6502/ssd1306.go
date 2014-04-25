package go6502

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"os"
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

func (s *Ssd1306) serveHttp() {
	httpHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "image/png")
		w.Header().Add("Refresh", "0.1")
		png.Encode(w, s.img)
	}
	address := "localhost:1234"
	srv := &http.Server{
		Addr:    address,
		Handler: http.HandlerFunc(httpHandler),
	}
	fmt.Printf("Ssd1306 output at http://%s/screen.png\n", address)
	go srv.ListenAndServe()
}

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
	fmt.Println("Writing SSD1306 data to PNG")
	writer, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}
	_ = png.Encode(writer, s.img)
}
