/*
	Package memory provides ROM & RAM for go6502; 16-bit address, 8-bit data.
*/
package memory

// Memory is a general interface for reading and writing bytes to and from
// 16-bit addresses.
type Memory interface {
	Shutdown()
	Read(uint16) byte
	Write(uint16, byte)
	Size() int
}
