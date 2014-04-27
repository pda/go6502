package memory

import "io/ioutil"

// Ram (32 KiB)
type Ram [0x8000]byte

func (r *Ram) String() string {
	return "(RAM 32K)"
}

// Read a byte from a 16-bit address.
func (mem *Ram) Read(a uint16) byte {
	return mem[a]
}

// Write a byte to a 16-bit address.
func (mem *Ram) Write(a uint16, value byte) {
	mem[a] = value
}

// Size of the RAM in bytes.
func (mem *Ram) Size() int {
	return 0x8000 // 32K
}

// Dump writes the RAM contents to the specified file path.
func (mem *Ram) Dump(path string) {
	err := ioutil.WriteFile(path, mem[:], 0640)
	if err != nil {
		panic(err)
	}
}
