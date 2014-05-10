package speedometer

import (
	"fmt"
	"time"

	"github.com/pda/go6502/cpu"
)

// Speedometer tracks how many instructions and cycles have executed in how
// much time, to calculate an effective MHz etc.
type Speedometer struct {
	cycles       uint64
	instructions uint64
	timeStart    time.Time
	cycleChan    chan uint8
}

// NewSpeedometer creates a Speedometer, and starts a goroutine to receive
// cycle counts from Speedometer.BeforeExecute().
func NewSpeedometer() *Speedometer {
	s := &Speedometer{
		timeStart: time.Now(),
		cycleChan: make(chan uint8),
	}
	go func() {
		for {
			s.cycles += uint64(<-s.cycleChan)
			s.instructions++
		}
	}()
	return s
}

// BeforeExecute meets go6502.Monitor interface.
func (s *Speedometer) BeforeExecute(in cpu.Instruction) {
	s.cycleChan <- in.Cycles
}

// Shutdown the Speedometer session, reporting stats to stdout.
func (s *Speedometer) Shutdown() {
	duration := time.Since(s.timeStart)
	us := float64(duration) / float64(time.Microsecond)

	fmt.Printf("Speedometer\n")
	fmt.Printf("----------------------------------\n")
	fmt.Printf("Instructions: % 20d\n", s.instructions)
	fmt.Printf("Cycles:       % 20d\n", s.cycles)
	fmt.Printf("Seconds:      % 20.2f\n", duration.Seconds())
	fmt.Printf("MHz:          % 20.2f\n", float64(s.cycles)/us)
	fmt.Printf("MIPS:         % 20.2f\n", float64(s.instructions)/us)
	fmt.Printf("----------------------------------\n")
}
