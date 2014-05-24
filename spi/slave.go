package spi

// Slave represents an 8-bit MSB-first mode-0 SPI slave device.
type Slave struct {

	// Done is true after a write() completed a byte transfer.
	Done bool

	// Mosi is the byte most recently transferred from master to slave.
	Mosi byte

	// Miso is the byte most recently transferred from slave to master.
	Miso byte

	PinMap

	clock      bool  // the most recent clock state
	index      uint8 // the bit index of the current byte.
	misoBuffer byte  // current byte being sent one bit at a time via Read().
	readByte   byte  // the state of the pins as read by the VIA controller.
	mosiBuffer byte  // the byte being built from bits

	maskSclk uint8
	maskMosi uint8
	maskMiso uint8
	maskSs   uint8
}

func NewSlave(pm PinMap) *Slave {
	return &Slave{
		index:    7,
		maskSclk: 1 << pm.Sclk,
		maskMosi: 1 << pm.Mosi,
		maskMiso: 1 << pm.Miso,
		maskSs:   1 << pm.Ss,
	}
}

// Read returns the current output (MISO) state for the parallel interface.
func (s *Slave) Read() byte {
	return s.readByte
}

// Write takes a byte of parallel data containing Sclk, Mosi, Miso, Ss.
// It may update the result of Read().
// spi.Done is updated to reflect whether the write completed a byte transfer,
// in which case spi.Mosi is set.
func (s *Slave) Write(data byte) bool {
	if data&s.maskSs != 0 {
		// do nothing unless SS is low (active)
		return false
	}

	s.Done = false

	mosi := data&s.maskMosi > 0
	clock := data&s.maskSclk > 0

	rising := !s.clock && clock
	falling := s.clock && !clock
	s.clock = clock

	// sclk:rise -> miso -> sclk:fall -> mosi -> ...

	if rising {
		if s.misoBuffer&(1<<s.index) > 0 {
			s.readByte = 0x00 | s.maskMiso
		} else {
			s.readByte = 0x00
		}
	}

	if falling {
		if mosi {
			s.mosiBuffer |= (1 << s.index)
		}

		// after eigth bit
		if s.index == 0 {
			s.index = 7
			s.Mosi = s.mosiBuffer
			s.Miso = s.misoBuffer
			s.Done = true
			s.mosiBuffer = 0x00
		} else {
			s.index--
		}
	}

	return true
}

// QueueMisoBits loads a byte into the MISO buffer, to be sent during the next
// eight clock cycles.
func (s *Slave) QueueMisoBits(b byte) {
	if s.index != 7 {
		panic("Cannot queue MISO; byte send in progress.")
	}
	s.misoBuffer = b
}
