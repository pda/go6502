package c64

import(
  "fmt"
  "strings"
)

// status register bits
const(
  sCarry = iota
  sZero
  sInterrupt
  sDecimal
  sBreak
  _
  sOverflow
  sNegative
)

func srName(bit uint8) (s string) {
  switch bit {
  case 0: s = "carry"
  case 1: s = "zero"
  case 2: s = "interrupt"
  case 3: s = "decimal"
  case 4: s = "break"
  case 5: s = "_"
  case 6: s = "overflow"
  case 7: s = "negative"
  }
  return
}

type Cpu struct {
  pc address
  ac byte
  x byte
  y byte
  sp byte
  sr byte
  Bus *Bus
}

type Iop struct {
  in *Instruction
  op8 uint8
  op16 address
}

func (iop *Iop) String() string {
  return fmt.Sprintf("%v op8:0x%02X op16:0x%04X", iop.in, iop.op8, iop.op16)
}

func (c *Cpu) Reset() {
  c.pc = c.Bus.Read16(0xFFFC)
  c.ac = 0x00
  c.x = 0x00
  c.y = 0x00
  c.sp = 0xFF // address relative to second page of memory (0x0100 ~ 0x01FF)
  c.sr = 0x00
}

func (c *Cpu) Step() {
  op := c.Bus.Read(c.pc)
  c.pc++
  in := findInstruction(op)
  iop := c.readOperand(in)
  fmt.Println(iop)
  c.Execute(iop)
  fmt.Println(c)
}

func (c *Cpu) readOperand(in *Instruction) *Iop {
  // read instruction with operand
  iop := &Iop{in: in}
  switch in.bytes {
  case 1: // no operand
  case 2: iop.op8 = c.Bus.Read(c.pc)
  case 3: iop.op16 = c.Bus.Read16(c.pc)
  default: panic(fmt.Sprintf("unhandled instruction length: %d", in.bytes))
  }
  c.pc += address(in.bytes - 1)
  return iop
}

func (c *Cpu) String() string {
  return fmt.Sprintf(
    "CPU pc:0x%04X ac:0x%02X x:0x%02X y:0x%02X sp:0x%02X sr:%s",
    c.pc, c.ac, c.x, c.y, c.sp,
    c.statusString(),
  )
}

func (c *Cpu) StackHead(offset int8) address {
  return address(0x0100) + address(c.sp) + address(offset)
}

func (c *Cpu) resolveOperand(iop *Iop) uint8 {
  switch iop.in.addressing {
  case immediate: return iop.op8
  default: return c.Bus.Read(c.memoryAddress(iop))
  }
}

func (c *Cpu) memoryAddress(iop *Iop) address {
  switch iop.in.addressing {
  case absolute: return iop.op16
  case absoluteX: return iop.op16 + address(c.x)
  case absoluteY: return iop.op16 + address(c.y)

  // Indexed indirect addressing.
  // In indexed indirect addressing (referred to as [Indirect, X]), the second
  // byte of the instruction is added to the contents of the X index register,
  // discarding the carry.  The result of this addition points to a memory
  // location on page zero whose contents is the low order eight bits of the
  // effective address.  The next memory location in page zero contains the
  // high order eight bits of the effective address.  Both memory locations
  // specitying the high and low order bytes of the effective address must be
  // in page zero.
  //
  // TODO: fix this implementation.
  case indirectX: return c.Bus.Read16(address(iop.op8) + address(c.x))

  // Indirect indexed addressing.
  // In indirect indexed addressing (referred to as [Indirect, Y]), the second
  // byte of the instruction points to a memory location in page zero. The
  // contents of this memory location is added to the contents of the Y index
  // register, the result being the low oeder eight bits of the effective
  // address. The carry from this addition is added to the contents of the next
  // page zero memory location, the result being the high order eight bits of
  // the effective address.
  //
  // TODO: fix this implementation.
  case indirectY: return c.Bus.Read16(address(iop.op8) + address(c.y))

  case zeropage: return address(iop.op8)
  default: panic("unhandled addressing")
  }
}

func (c *Cpu) getStatus(bit uint8) bool {
  return c.getStatusInt(bit) == 1
}

func (c *Cpu) getStatusInt(bit uint8) uint8 {
  return (c.sr >> bit) & 1
}

func (c *Cpu) setStatus(bit uint8, state bool) {
  fmt.Printf("sr %s = %v\n", srName(bit), state)
  if state {
    c.sr |= 1 << bit
  } else {
    c.sr &^= 1 << bit
  }
}

func (c *Cpu) updateStatus(value uint8) {
  c.setStatus(sZero, value == 0)
  c.setStatus(sNegative, (value >> 7) == 1)
}

func (c *Cpu) statusString() string {
  // this is horrible, I think. Should be much simpler?
  var chars = "nv_bdizc"
  var out [8]string
  for i := 0; i < 8; i++ {
    if c.getStatus(uint8(7 - i)) {
      out[i] = string(chars[i])
    } else {
      out[i] = "-"
    }
  }
  return strings.Join(out[0:], "")
}

func (c *Cpu) branch(iop *Iop) {
  relative := int8(iop.op8) // signed
  if relative >= 0 {
    c.pc += address(relative)
  } else {
    c.pc -= address(-relative)
  }
}

func (c *Cpu) Execute(iop *Iop) {
  switch iop.in.id {
  case AND: c.AND(iop)
  case BCS: c.BCS(iop)
  case BEQ: c.BEQ(iop)
  case BMI: c.BMI(iop)
  case BNE: c.BNE(iop)
  case CLC: c.CLC(iop)
  case CLD: c.CLD(iop)
  case CMP: c.CMP(iop)
  case DEX: c.DEX(iop)
  case DEY: c.DEY(iop)
  case INC: c.INC(iop)
  case INX: c.INX(iop)
  case INY: c.INY(iop)
  case JMP: c.JMP(iop)
  case JSR: c.JSR(iop)
  case LDA: c.LDA(iop)
  case LDX: c.LDX(iop)
  case LDY: c.LDY(iop)
  case ORA: c.ORA(iop)
  case ROL: c.ROL(iop)
  case RTS: c.RTS(iop)
  case SEI: c.SEI(iop)
  case STA: c.STA(iop)
  case STX: c.STX(iop)
  case STY: c.STY(iop)
  case TAX: c.TAX(iop)
  case TAY: c.TAY(iop)
  case TXA: c.TXA(iop)
  case TXS: c.TXS(iop)
  case TYA: c.TYA(iop)
  default: panic(fmt.Sprintf("unhandled instruction: %v", iop.in.name()))
  }
}

// bitwise AND with accumulator
func (c *Cpu) AND(iop *Iop) {
  c.ac &= c.resolveOperand(iop)
  c.updateStatus(c.ac)
}

// branch on carry (when carry set)
func (c *Cpu) BCS(iop *Iop) {
  if c.getStatus(sCarry) {
    c.branch(iop)
  }
}

// branch on equal (zero set)
// (branch on z = 1)
func (c *Cpu) BEQ(iop *Iop) {
  if c.getStatus(sZero) {
    c.branch(iop)
  }
}

// branch on result minus (status negative)
func (c *Cpu) BMI(iop *Iop) {
  if c.getStatus(sNegative) {
    c.branch(iop)
  }
}

// branch on not-equal (zero clear)
func (c *Cpu) BNE(iop *Iop) {
  if !c.getStatus(sZero) {
    c.branch(iop)
  }
}

// clear carry
func (c *Cpu) CLC(iop *Iop) {
  c.setStatus(sCarry, false)
}

// clear decimal
func (c *Cpu) CLD(iop *Iop) {
  c.setStatus(sDecimal, false)
}

// compare (with accumulator)
func (c *Cpu) CMP(iop *Iop) {
  value := c.resolveOperand(iop)
  c.setStatus(sCarry, c.ac >= value)
  c.updateStatus(c.ac - value)
}

// decrement x
func (c *Cpu) DEX(iop *Iop) {
  c.x--
  c.updateStatus(c.x)
}

// decrement y
func (c *Cpu) DEY(iop *Iop) {
  c.y--
  c.updateStatus(c.y)
}

// increment value in memory
func (c *Cpu) INC(iop *Iop) {
  address := c.memoryAddress(iop)
  value := c.Bus.Read(address) + 1
  c.Bus.Write(address, value)
  c.updateStatus(value)
}

// increment x
func (c *Cpu) INX(iop *Iop) {
  c.x++
  c.updateStatus(c.x)
}

// increment y
func (c *Cpu) INY(iop *Iop) {
  c.y++
  c.updateStatus(c.y)
}

// jump
func (c *Cpu) JMP(iop *Iop) {
  c.pc = c.memoryAddress(iop)
}

// jump to subroutine
func (c *Cpu) JSR(iop *Iop) {
  c.Bus.Write16(c.StackHead(-1), c.pc - 1)
  c.sp -= 2
  c.pc = iop.op16
}

// load accumulator
func (c *Cpu) LDA(iop *Iop) {
  c.ac = c.resolveOperand(iop)
  c.updateStatus(c.ac)
}

// load Y
func (c *Cpu) LDY(iop *Iop) {
  c.y = c.resolveOperand(iop)
  c.updateStatus(c.y)
}

// load X
func (c *Cpu) LDX(iop *Iop) {
  c.x = c.resolveOperand(iop)
  c.updateStatus(c.x)
}

// bitwise OR memory with accumulator
func (c *Cpu) ORA(iop *Iop) {
  c.ac |= c.resolveOperand(iop)
  c.updateStatus(c.ac)
}

// bitwise rotate left
// SR carry into bit 0, original bit 7 into SR carry.
func (c *Cpu) ROL(iop *Iop) {
  switch iop.in.addressing {
  case accumulator:
    carry := c.getStatusInt(sCarry)
    c.setStatus(sCarry, (c.ac >> 7) == 1)
    c.ac = (c.ac << 1) | carry
    c.updateStatus(c.ac)
  default: panic("ROL non-accumulator addressing not implemented")
  }
}

// return from subroutine
func (c *Cpu) RTS(iop *Iop) {
  c.pc = c.Bus.Read16(c.StackHead(1))
  c.sp += 2
  c.pc += 1
}

// set interrupt disable
func (c *Cpu) SEI(iop *Iop) {
  c.setStatus(sInterrupt, false)
}

// store from accumulator
func (c *Cpu) STA(iop *Iop) {
  c.Bus.Write(c.memoryAddress(iop), c.ac)
}

// store from X
func (c *Cpu) STX(iop *Iop) {
  c.Bus.Write(c.memoryAddress(iop), c.x)
}

// store from Y
func (c *Cpu) STY(iop *Iop) {
  c.Bus.Write(c.memoryAddress(iop), c.y)
}

// transfer accumulator for index Y
func (c *Cpu) TAX(iop *Iop) {
  c.x = c.ac
  c.updateStatus(c.x)
}

// transfer accumulator for index Y
func (c *Cpu) TAY(iop *Iop) {
  c.y = c.ac
  c.updateStatus(c.y)
}

// Copies the current contents of the X register into the accumulator and sets
// the zero and negative flags as appropriate.
func (c *Cpu) TXA(iop *Iop) {
  c.ac = c.x
  c.updateStatus(c.ac)
}

// transfer X to stack pointer
func (c *Cpu) TXS(iop *Iop) {
  c.sp = c.x
  c.updateStatus(c.sp)
}

// Transfer Y to Accumulator
func (c *Cpu) TYA(iop *Iop) {
  c.ac = c.y
  c.updateStatus(c.ac)
}
