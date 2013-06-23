package c64

import(
  "fmt"
  "io/ioutil"
  "encoding/hex"
)

type Rom struct {
  name string
  size int // bytes
  data []byte
}

func RomFromFile(path string) *Rom {
  data, _ := ioutil.ReadFile(path)
  return &Rom{name: path, size: len(data), data: data}
}

func (r *Rom) String() string {
  return fmt.Sprintf("ROM[%dk:%s:%s..%s]",
  r.size / 1024,
  r.name,
  hex.EncodeToString(r.data[0:4]),
  hex.EncodeToString(r.data[len(r.data) - 4:]))
}
