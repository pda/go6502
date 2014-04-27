package go6502

import (
	"io/ioutil"
)

// Ram (32K)

type Ram [0x8000]byte

func (r *Ram) String() string {
	return "(RAM 32K)"
}

func (mem *Ram) Read(a Address) byte {
	return mem[a]
}

func (mem *Ram) Write(a Address, value byte) {
	mem[a] = value
}

func (mem *Ram) Size() int {
	return 0x8000 // 32K
}

func (mem *Ram) Dump(path string) {
	err := ioutil.WriteFile(path, mem[:], 0640)
	if err != nil {
		panic(err)
	}
}
