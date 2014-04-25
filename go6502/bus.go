package go6502

import (
	"fmt"
)

type busEntry struct {
	mem   Memory
	name  string
	start address
	end   address
}

type Bus struct {
	entries []busEntry
}

func CreateBus() (*Bus, error) {
	return &Bus{entries: make([]busEntry, 0)}, nil
}

func (b *Bus) Attach(mem Memory, name string, offset address) error {
	om := OffsetMemory{offset: offset, Memory: mem}
	end := offset + address(mem.Size()-1)
	entry := busEntry{mem: om, name: name, start: offset, end: end}
	b.entries = append(b.entries, entry)
	return nil
}

func (b *Bus) backendFor(a address) (Memory, error) {
	for _, be := range b.entries {
		if a >= be.start && a <= be.end {
			return be.mem, nil
		}
	}
	return nil, fmt.Errorf("No backend for address 0x%04X", a)
}

func (b *Bus) Read(a address) byte {
	mem, err := b.backendFor(a)
	if err != nil {
		panic(err)
	}
	value := mem.Read(a)
	return value
}

func (b *Bus) Read16(a address) address {
	lo := address(b.Read(a))
	hi := address(b.Read(a + 1))
	return hi<<8 | lo
}

func (b *Bus) String() string {
	return fmt.Sprintf("Address bus (TODO: describe)")
}

func (b *Bus) Write(a address, value byte) {
	mem, err := b.backendFor(a)
	if err != nil {
		panic(err)
	}
	mem.Write(a, value)
}

func (b *Bus) Write16(a address, value address) {
	b.Write(a, byte(value))
	b.Write(a+1, byte(value>>8))
}
