package sd

import "fmt"

// states
const (
	sCmd = iota // expect command
	sArg        // expect argument
	sChk        // expect checksum
)

const (
	r1_ready = 0x00
	r1_idle  = 0x01
)

// sdState is the state of SD protocol (layer above SPI protocol).
type sdState struct {
	state     uint8
	acmd      bool // next command is an application-specific command
	cmd       uint8
	arg       uint32
	argByte   uint8
	misoQueue []byte // data waiting to be sent from card.
	prevCmd   uint8
	prevAcmd  uint8
}

func newSdState() (s *sdState) {
	return &sdState{
		misoQueue: make([]byte, 0, 1024),
	}
}

func (s *sdState) consumeByte(b byte) {
	switch s.state {
	case sCmd:
		if b>>6 == 1 {
			s.cmd = b & (0xFF >> 2)
			s.state = sArg
			s.argByte = 0
		}
	case sArg:
		s.arg |= uint32(b) << ((3 - s.argByte) * 8)
		if s.argByte == 3 {
			s.state = sChk
		} else {
			s.argByte++
		}
	case sChk:
		if s.acmd {
			s.handleAcmd()
		} else {
			s.handleCmd()
		}

	default:
		panic("Unhandled state")
	}
}

func (s *sdState) handleCmd() {
	fmt.Printf("SD CMD%d arg: 0x%08X\n", s.cmd, s.arg)
	switch s.cmd {
	case 0: // GO_IDLE_STATE
		fmt.Println("SD CMD0 response: r1_idle")
		s.queueMisoBytes(0xFF, 0xFF, r1_idle) // busy then idle
		s.state = sCmd
	case 55: // APP_CMD
		fmt.Println("SD CMD55 response: r1_idle")
		s.queueMisoBytes(r1_idle) // busy then idle
		s.acmd = true
		s.state = sCmd
	default:
		panic(fmt.Sprintf("Unhandled CMD%d", s.cmd))
	}
	s.prevCmd = s.cmd
}

func (s *sdState) handleAcmd() {
	fmt.Printf("SD ACMD%d arg: 0x%08X\n", s.cmd, s.arg)
	switch s.cmd {
	case 41: // SD_SEND_OP_COND
		if s.prevAcmd == 41 {
			// on second attempt, busy, busy, then ready.
			fmt.Println("SD ACMD41 response: r1_ready")
			s.queueMisoBytes(0xFF, 0xFF, r1_ready)
		} else {
			// on first attempt, busy, busy, then idle (not yet ready).
			fmt.Println("SD ACMD41 response: r1_idle")
			s.queueMisoBytes(0xFF, 0xFF, r1_idle)
		}
		s.state = sCmd
	default:
		panic(fmt.Sprintf("Unhandled ACMD%d", s.cmd))
	}
	s.prevAcmd = s.cmd
	s.acmd = false
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
