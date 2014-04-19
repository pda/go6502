package go6502

import "fmt"

// SpiDebugger implements ParallelPeripheral interface for Via6522.

type SpiDebugger struct {
	lastClock   bool
	inputBuffer byte
	inputIndex  uint8
}

func NewSpiDebugger() *SpiDebugger {
	s := SpiDebugger{}
	s.inputIndex = 7 // MSB-first, decrementing index.
	return &s
}

// TODO: configurable lines
const (
	mosiMask  = 1 << 0
	clockMask = 1 << 1
)

func (s *SpiDebugger) Notify(data byte) {

	mosi := data&mosiMask > 0
	clock := data&clockMask > 0

	if clock && !s.lastClock {
		// rising clock
		s.lastClock = clock
		if mosi {
			s.inputBuffer |= (1 << s.inputIndex)
		}
		if s.inputIndex == 0 {
			fmt.Printf("SpiDebugger: 0x%02X 0b%08b\n", s.inputBuffer, s.inputBuffer)
			s.inputIndex = 7
			s.inputBuffer = 0x00
		} else {
			s.inputIndex--
		}
	}

	if !clock && s.lastClock {
		// falling clock
		s.lastClock = clock
	}

}
