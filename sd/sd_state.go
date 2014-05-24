package sd

import "fmt"

// sdState is the state of SD protocol (layer above SPI protocol).
type sdState struct {
	misoQueue []byte // data waiting to be sent from card.
}

func newSdState() (s *sdState) {
	return &sdState{
		misoQueue: make([]byte, 0, 1024),
	}
}

func (s *sdState) consumeByte(b byte) {
	switch b {
	case 0x40:
		fmt.Printf("SD: Got 0x40; queueing response bytes.\n")
		s.queueMisoBytes(0xAA, 0xAB, 0xAC, 0xAD)
	}
}

func (s *sdState) queueMisoBytes(bytes ...byte) {
	s.misoQueue = append(s.misoQueue, bytes...)
}

func (s *sdState) shiftMiso() (b byte) {
	if len(s.misoQueue) > 0 {
		b = s.misoQueue[0]
		s.misoQueue = s.misoQueue[1:len(s.misoQueue)]
	} else {
		b = 0x00 // default to low for empty buffer.
	}
	return
}
