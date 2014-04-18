package go6502

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
)

type Rom struct {
	name string
	size int // bytes
	data []byte
}

func (rom *Rom) Read(a address) byte {
	return rom.data[a]
}

func RomFromFile(path string) (*Rom, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return &Rom{name: path, size: len(data), data: data}, nil
}

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

func (r *Rom) Write(_ address, _ byte) {
	panic(fmt.Sprintf("%v is read-only", r))
}
