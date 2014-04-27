package memory

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
)

// A Rom provides read-only memory, with data generally pre-loaded from a file.
type Rom struct {
	name string
	size int // bytes
	data []byte
}

// Read a byte from the given address.
func (rom *Rom) Read(a uint16) byte {
	return rom.data[a]
}

// Create a new ROM, loading the contents from a file.
// The size of the ROM is determined by the size of the file.
func RomFromFile(path string) (*Rom, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return &Rom{name: path, size: len(data), data: data}, nil
}

// Size of the Rom in bytes.
func (r *Rom) Size() int {
	return r.size
}

func (r *Rom) String() string {
	return fmt.Sprintf("ROM[%dk:%s:%s..%s]",
		r.Size()/1024,
		r.name,
		hex.EncodeToString(r.data[0:2]),
		hex.EncodeToString(r.data[len(r.data)-2:]))
}

// Rom meets the go6502.Memory interface, but Write is not supported, and will
// cause an error.
func (r *Rom) Write(_ uint16, _ byte) {
	panic(fmt.Sprintf("%v is read-only", r))
}
