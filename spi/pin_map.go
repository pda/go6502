package spi

// PinMap associates SPI lines with parallel port pin numbers (0..7).
type PinMap struct {
	Sclk uint
	Mosi uint
	Miso uint
	Ss   uint
}

func (p PinMap) PinMask() byte {
	return 1<<p.Sclk | 1<<p.Mosi | 1<<p.Miso | 1<<p.Ss
}
