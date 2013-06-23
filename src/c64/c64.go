package c64

// Address
type address uint16;

// Ram

type Ram [0x10000]byte
func (r *Ram) String() string {
  return "(RAM 64K)"
}
