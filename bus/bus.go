/*
	Package bus provides a mappable 16-bit addressable 8-bit data bus for go6502.
	Different Memory backends can be attached at different base addresses.
*/
package bus

import (
	"fmt"

	"github.com/pda/go6502/memory"
)

type busEntry struct {
	mem   memory.Memory
	name  string
	start uint16
	end   uint16
}

// Bus is a 16-bit address, 8-bit data bus, which maps reads and writes
// at different locations to different backend Memory. For example the
// lower 32K could be RAM, the upper 8KB ROM, and some I/O in the middle.
type Bus struct {
	entries []busEntry
}

func (b *Bus) String() string {
	return fmt.Sprintf("Address bus (TODO: describe)")
}

func CreateBus() (*Bus, error) {
	return &Bus{entries: make([]busEntry, 0)}, nil
}

// Attach maps a bus address range to a backend Memory implementation,
// which could be RAM, ROM, I/O device etc.
func (b *Bus) Attach(mem memory.Memory, name string, offset uint16) error {
	om := OffsetMemory{Offset: offset, Memory: mem}
	end := offset + uint16(mem.Size()-1)
	entry := busEntry{mem: om, name: name, start: offset, end: end}
	b.entries = append(b.entries, entry)
	return nil
}

func (b *Bus) backendFor(a uint16) (memory.Memory, error) {
	for _, be := range b.entries {
		if a >= be.start && a <= be.end {
			return be.mem, nil
		}
	}
	return nil, fmt.Errorf("No backend for address 0x%04X", a)
}

// Shutdown tells the address bus a shutdown is occurring, and to pass the
// message on to subordinates.
func (b *Bus) Shutdown() {
	for _, be := range b.entries {
		be.mem.Shutdown()
	}
}

// Read returns the byte from memory mapped to the given address.
// e.g. if ROM is mapped to 0xC000, then Read(0xC0FF) returns the byte at
// 0x00FF in that RAM device.
func (b *Bus) Read(a uint16) byte {
	mem, err := b.backendFor(a)
	if err != nil {
		panic(err)
	}
	value := mem.Read(a)
	return value
}

// Read16 returns the 16-bit value stored in little-endian format with the
// low byte at address, and the high byte at address+1.
func (b *Bus) Read16(a uint16) uint16 {
	lo := uint16(b.Read(a))
	hi := uint16(b.Read(a + 1))
	return hi<<8 | lo
}

// Write the byte to the device mapped to the given address.
func (b *Bus) Write(a uint16, value byte) {
	mem, err := b.backendFor(a)
	if err != nil {
		panic(err)
	}
	mem.Write(a, value)
}

// Write16 writes the given 16-bit value to the specifie address, storing it
// little-endian, with high byte at address+1.
func (b *Bus) Write16(a uint16, value uint16) {
	b.Write(a, byte(value))
	b.Write(a+1, byte(value>>8))
}
