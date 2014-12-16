//go:generate stringer -type state
//go:generate stringer -type response

package sd

import "fmt"

type state uint8

// states
const (
	sCommand  state = iota // expect command
	sArgument              // expect argument
	sChecksum              // expect checksum
	sData                  // sending data until misoQueue empty.
)

type response uint8

const (
	r1_ready response = 0x00
	r1_idle  response = 0x01
)

const (
	// blockSize isn't strictly constant, but...
	blockSize = 512
)

// sdCard is the state of SD protocol (layer above SPI protocol).
type sdCard struct {
	state     state
	acmd      bool // next command is an application-specific command
	cmd       uint8
	arg       uint32
	argByte   uint8
	misoQueue []byte // data waiting to be sent from card.
	prevCmd   uint8
	prevAcmd  uint8
	data      []byte
}

func newSdCard() (sd *sdCard) {
	return &sdCard{
		misoQueue: make([]byte, 0, 1024),
	}
}

func (sd *sdCard) enter(state state) {
	fmt.Printf("SD state %s -> %s\n", sd.state, state)
	sd.state = state
}

func (sd *sdCard) consumeByte(b byte) {
	switch sd.state {
	case sCommand:
		if b>>6 == 1 {
			sd.cmd = b & (0xFF >> 2)
			sd.enter(sArgument)
			sd.arg = 0x00000000
			sd.argByte = 0
		}
	case sArgument:
		sd.arg |= uint32(b) << ((3 - sd.argByte) * 8)
		if sd.argByte == 3 {
			sd.enter(sChecksum)
		} else {
			sd.argByte++
		}
	case sChecksum:
		if sd.acmd {
			sd.handleAcmd()
		} else {
			sd.handleCmd()
		}
	case sData:
		// ignore; data it being sent.

	default:
		panic(fmt.Errorf("Unhandled state: %d", sd.state))
	}
}

func (sd *sdCard) handleCmd() {
	fmt.Printf("SD CMD%d arg: 0x%08X\n", sd.cmd, sd.arg)
	switch sd.cmd {
	case 0: // GO_IDLE_STATE
		fmt.Println("SD CMD0 response: r1_idle")
		sd.queueMisoBytes(0xFF, 0xFF, byte(r1_idle)) // busy then idle
		sd.enter(sCommand)
	case 17: // READ_SINGLE_BLOCK
		fmt.Println("SD CMD17 response: r1_ready, data start block, data")
		sd.queueMisoBytes(0xFF, 0xFF, byte(r1_ready)) // busy then ready
		sd.queueMisoBytes(0xFF, 0xFF, 0xFF, 0xFF)     // time before data block
		sd.queueMisoBytes(0xFE)                       // data start block
		sd.queueMisoBytes(sd.readBlock(sd.arg)...)
		sd.enter(sData)
	case 55: // APP_CMD
		fmt.Println("SD CMD55 response: r1_idle")
		sd.queueMisoBytes(byte(r1_idle)) // busy then idle
		sd.acmd = true
		sd.enter(sCommand)
	default:
		panic(fmt.Sprintf("Unhandled CMD%d", sd.cmd))
	}
	sd.prevCmd = sd.cmd
}

func (sd *sdCard) handleAcmd() {
	fmt.Printf("SD ACMD%d arg: 0x%08X\n", sd.cmd, sd.arg)
	switch sd.cmd {
	case 41: // SD_SEND_OP_COND
		if sd.prevAcmd == 41 {
			// on second attempt, busy, busy, then ready.
			fmt.Println("SD ACMD41 response: r1_ready")
			sd.queueMisoBytes(0xFF, 0xFF, byte(r1_ready))
		} else {
			// on first attempt, busy, busy, then idle (not yet ready).
			fmt.Println("SD ACMD41 response: r1_idle")
			sd.queueMisoBytes(0xFF, 0xFF, byte(r1_idle))
		}
		sd.enter(sCommand)
	default:
		panic(fmt.Sprintf("Unhandled ACMD%d", sd.cmd))
	}
	sd.prevAcmd = sd.cmd
	sd.acmd = false
}

func (sd *sdCard) queueMisoBytes(bytes ...byte) {
	sd.misoQueue = append(sd.misoQueue, bytes...)
}

func (sd *sdCard) shiftMiso() (b byte) {
	if len(sd.misoQueue) > 0 {
		b = sd.misoQueue[0]
		sd.misoQueue = sd.misoQueue[1:len(sd.misoQueue)]
		if len(sd.misoQueue) == 0 && sd.state == sData {
			// transition from sData to sCommand when all data sent.
			sd.enter(sCommand)
		}
	} else {
		b = 0x00 // default to low for empty buffer.
	}
	return
}

func (sd *sdCard) readBlock(start uint32) []byte {
	// TODO: bounds checking
	// TODO: zero-fill remainder of last page in sd.data?
	return sd.data[start : start+blockSize]
}
