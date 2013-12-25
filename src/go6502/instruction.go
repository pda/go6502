package go6502

import(
  "fmt"
)

// addressing modes
const(
  _ = iota
  absolute
  absoluteX
  absoluteY
  accumulator
  immediate
  implied
  indirect
  indirectX
  indirectY
  relative
  zeropage
  zeropageX
  zeropageY
)

// instruction (mnemonic) ids.
const(
  _ = iota
  ADC
  AND
  ASL
  BCC
  BCS
  BEQ
  BIT
  BMI
  BNE
  BPL
  BRK
  BVC
  BVS
  CLC
  CLD
  CLI
  CLV
  CMP
  CPX
  CPY
  DEC
  DEX
  DEY
  EOR
  INC
  INX
  INY
  JMP
  JSR
  LDA
  LDX
  LDY
  LSR
  NOP
  ORA
  PHA
  PHP
  PLA
  PLP
  ROL
  ROR
  RTI
  RTS
  SBC
  SEC
  SED
  SEI
  STA
  STX
  STY
  TAX
  TAY
  TSX
  TXA
  TXS
  TYA
)

type Instruction struct {
  id uint8  // the const identifier of the instruction type, e.g. ADC
  opcode byte
  addressing int
  bytes int
  cycles int
  flags int
}

func findInstruction(opcode byte) *Instruction {
  // TODO: singleton instructions; they're immutable.
  var i Instruction;
  switch opcode {
  default: panic(fmt.Sprintf("Unknown opcode: 0x%02X", opcode))
  case 0x69: i = Instruction{ADC, opcode, immediate, 2, 2, 0}
  case 0x65: i = Instruction{ADC, opcode, zeropage, 2, 3, 0}
  case 0x75: i = Instruction{ADC, opcode, zeropageX, 2, 4, 0}
  case 0x6D: i = Instruction{ADC, opcode, absolute, 3, 4, 0}
  case 0x7D: i = Instruction{ADC, opcode, absoluteX, 3, 4, 0}
  case 0x79: i = Instruction{ADC, opcode, absoluteY, 3, 4, 0}
  case 0x61: i = Instruction{ADC, opcode, indirectX, 2, 6, 0}
  case 0x71: i = Instruction{ADC, opcode, indirectY, 2, 5, 0}
  case 0x29: i = Instruction{AND, opcode, immediate, 2, 2, 0}
  case 0x25: i = Instruction{AND, opcode, zeropage, 2, 3, 0}
  case 0x35: i = Instruction{AND, opcode, zeropageX, 2, 4, 0}
  case 0x2D: i = Instruction{AND, opcode, absolute, 3, 4, 0}
  case 0x3D: i = Instruction{AND, opcode, absoluteX, 3, 4, 0}
  case 0x39: i = Instruction{AND, opcode, absoluteY, 3, 4, 0}
  case 0x21: i = Instruction{AND, opcode, indirectX, 2, 6, 0}
  case 0x31: i = Instruction{AND, opcode, indirectY, 2, 5, 0}
  case 0x0A: i = Instruction{ASL, opcode, accumulator, 1, 2, 0}
  case 0x06: i = Instruction{ASL, opcode, zeropage, 2, 5, 0}
  case 0x16: i = Instruction{ASL, opcode, zeropageX, 2, 6, 0}
  case 0x0E: i = Instruction{ASL, opcode, absolute, 3, 6, 0}
  case 0x1E: i = Instruction{ASL, opcode, absoluteX, 3, 7, 0}
  case 0x90: i = Instruction{BCC, opcode, relative, 2, 2, 0}
  case 0xB0: i = Instruction{BCS, opcode, relative, 2, 2, 0}
  case 0xF0: i = Instruction{BEQ, opcode, relative, 2, 2, 0}
  case 0x24: i = Instruction{BIT, opcode, zeropage, 2, 3, 0}
  case 0x2C: i = Instruction{BIT, opcode, absolute, 3, 4, 0}
  case 0x30: i = Instruction{BMI, opcode, relative, 2, 2, 0}
  case 0xD0: i = Instruction{BNE, opcode, relative, 2, 2, 0}
  case 0x10: i = Instruction{BPL, opcode, relative, 2, 2, 0}
  case 0x00: i = Instruction{BRK, opcode, implied, 1, 7, 0}
  case 0x50: i = Instruction{BVC, opcode, relative, 2, 2, 0}
  case 0x70: i = Instruction{BVS, opcode, relative, 2, 2, 0}
  case 0x18: i = Instruction{CLC, opcode, implied, 1, 2, 0}
  case 0xD8: i = Instruction{CLD, opcode, implied, 1, 2, 0}
  case 0x58: i = Instruction{CLI, opcode, implied, 1, 2, 0}
  case 0xB8: i = Instruction{CLV, opcode, implied, 1, 2, 0}
  case 0xC9: i = Instruction{CMP, opcode, immediate, 2, 2, 0}
  case 0xC5: i = Instruction{CMP, opcode, zeropage, 2, 3, 0}
  case 0xD5: i = Instruction{CMP, opcode, zeropageX, 2, 4, 0}
  case 0xCD: i = Instruction{CMP, opcode, absolute, 3, 4, 0}
  case 0xDD: i = Instruction{CMP, opcode, absoluteX, 3, 4, 0}
  case 0xD9: i = Instruction{CMP, opcode, absoluteY, 3, 4, 0}
  case 0xC1: i = Instruction{CMP, opcode, indirectX, 2, 6, 0}
  case 0xD1: i = Instruction{CMP, opcode, indirectY, 2, 5, 0}
  case 0xE0: i = Instruction{CPX, opcode, immediate, 2, 2, 0}
  case 0xE4: i = Instruction{CPX, opcode, zeropage, 2, 3, 0}
  case 0xEC: i = Instruction{CPX, opcode, absolute, 3, 4, 0}
  case 0xC0: i = Instruction{CPY, opcode, immediate, 2, 2, 0}
  case 0xC4: i = Instruction{CPY, opcode, zeropage, 2, 3, 0}
  case 0xCC: i = Instruction{CPY, opcode, absolute, 3, 4, 0}
  case 0xC6: i = Instruction{DEC, opcode, zeropage, 2, 5, 0}
  case 0xD6: i = Instruction{DEC, opcode, zeropageX, 2, 6, 0}
  case 0xCE: i = Instruction{DEC, opcode, absolute, 3, 3, 0}
  case 0xDE: i = Instruction{DEC, opcode, absoluteX, 3, 7, 0}
  case 0xCA: i = Instruction{DEX, opcode, implied, 1, 2, 0}
  case 0x88: i = Instruction{DEY, opcode, implied, 1, 2, 0}
  case 0x49: i = Instruction{EOR, opcode, immediate, 2, 2, 0}
  case 0x45: i = Instruction{EOR, opcode, zeropage, 2, 3, 0}
  case 0x55: i = Instruction{EOR, opcode, zeropageX, 2, 4, 0}
  case 0x4D: i = Instruction{EOR, opcode, absolute, 3, 4, 0}
  case 0x5D: i = Instruction{EOR, opcode, absoluteX, 3, 4, 0}
  case 0x59: i = Instruction{EOR, opcode, absoluteY, 3, 4, 0}
  case 0x41: i = Instruction{EOR, opcode, indirectX, 2, 6, 0}
  case 0x51: i = Instruction{EOR, opcode, indirectY, 2, 5, 0}
  case 0xE6: i = Instruction{INC, opcode, zeropage, 2, 5, 0}
  case 0xF6: i = Instruction{INC, opcode, zeropageX, 2, 6, 0}
  case 0xEE: i = Instruction{INC, opcode, absolute, 3, 6, 0}
  case 0xFE: i = Instruction{INC, opcode, absoluteX, 3, 7, 0}
  case 0xE8: i = Instruction{INX, opcode, implied, 1, 2, 0}
  case 0xC8: i = Instruction{INY, opcode, implied, 1, 2, 0}
  case 0x4C: i = Instruction{JMP, opcode, absolute, 3, 3, 0}
  case 0x6C: i = Instruction{JMP, opcode, indirect, 3, 5, 0}
  case 0x20: i = Instruction{JSR, opcode, absolute, 3, 6, 0}
  case 0xA9: i = Instruction{LDA, opcode, immediate, 2, 2, 0}
  case 0xA5: i = Instruction{LDA, opcode, zeropage, 2, 3, 0}
  case 0xB5: i = Instruction{LDA, opcode, zeropageX, 2, 4, 0}
  case 0xAD: i = Instruction{LDA, opcode, absolute, 3, 4, 0}
  case 0xBD: i = Instruction{LDA, opcode, absoluteX, 3, 4, 0}
  case 0xB9: i = Instruction{LDA, opcode, absoluteY, 3, 4, 0}
  case 0xA1: i = Instruction{LDA, opcode, indirectX, 2, 6, 0}
  case 0xB1: i = Instruction{LDA, opcode, indirectY, 2, 5, 0}
  case 0xA2: i = Instruction{LDX, opcode, immediate, 2, 2, 0}
  case 0xA6: i = Instruction{LDX, opcode, zeropage, 2, 3, 0}
  case 0xB6: i = Instruction{LDX, opcode, zeropageY, 2, 4, 0}
  case 0xAE: i = Instruction{LDX, opcode, absolute, 3, 4, 0}
  case 0xBE: i = Instruction{LDX, opcode, absoluteY, 3, 4, 0}
  case 0xA0: i = Instruction{LDY, opcode, immediate, 2, 2, 0}
  case 0xA4: i = Instruction{LDY, opcode, zeropage, 2, 3, 0}
  case 0xB4: i = Instruction{LDY, opcode, zeropageX, 2, 4, 0}
  case 0xAC: i = Instruction{LDY, opcode, absolute, 3, 4, 0}
  case 0xBC: i = Instruction{LDY, opcode, absoluteX, 3, 4, 0}
  case 0x4A: i = Instruction{LSR, opcode, accumulator, 1, 2, 0}
  case 0x46: i = Instruction{LSR, opcode, zeropage, 2, 5, 0}
  case 0x56: i = Instruction{LSR, opcode, zeropageX, 2, 6, 0}
  case 0x4E: i = Instruction{LSR, opcode, absolute, 3, 6, 0}
  case 0x5E: i = Instruction{LSR, opcode, absoluteX, 3, 7, 0}
  case 0xEA: i = Instruction{NOP, opcode, implied, 1, 2, 0}
  case 0x09: i = Instruction{ORA, opcode, immediate, 2, 2, 0}
  case 0x05: i = Instruction{ORA, opcode, zeropage, 2, 3, 0}
  case 0x15: i = Instruction{ORA, opcode, zeropageX, 2, 4, 0}
  case 0x0D: i = Instruction{ORA, opcode, absolute, 3, 4, 0}
  case 0x1D: i = Instruction{ORA, opcode, absoluteX, 3, 4, 0}
  case 0x19: i = Instruction{ORA, opcode, absoluteY, 3, 4, 0}
  case 0x01: i = Instruction{ORA, opcode, indirectX, 2, 6, 0}
  case 0x11: i = Instruction{ORA, opcode, indirectY, 2, 5, 0}
  case 0x48: i = Instruction{PHA, opcode, implied, 1, 3, 0}
  case 0x08: i = Instruction{PHP, opcode, implied, 1, 3, 0}
  case 0x68: i = Instruction{PLA, opcode, implied, 1, 4, 0}
  case 0x28: i = Instruction{PHP, opcode, implied, 1, 4, 0}
  case 0x2A: i = Instruction{ROL, opcode, accumulator, 1, 2, 0}
  case 0x26: i = Instruction{ROL, opcode, zeropage, 2, 5, 0}
  case 0x36: i = Instruction{ROL, opcode, zeropageX, 2, 6, 0}
  case 0x2E: i = Instruction{ROL, opcode, absolute, 3, 6, 0}
  case 0x3E: i = Instruction{ROL, opcode, absoluteX, 3, 7, 0}
  case 0x6A: i = Instruction{ROR, opcode, accumulator, 1, 2, 0}
  case 0x66: i = Instruction{ROR, opcode, zeropage, 2, 5, 0}
  case 0x76: i = Instruction{ROR, opcode, zeropageX, 2, 6, 0}
  case 0x6E: i = Instruction{ROR, opcode, absolute, 3, 6, 0}
  case 0x7E: i = Instruction{ROR, opcode, absoluteX, 3, 7, 0}
  case 0x40: i = Instruction{RTI, opcode, implied, 1, 6, 0}
  case 0x60: i = Instruction{RTS, opcode, implied, 1, 6, 0}
  case 0xE9: i = Instruction{SBC, opcode, immediate, 2, 2, 0}
  case 0xE5: i = Instruction{SBC, opcode, zeropage, 2, 3, 0}
  case 0xF5: i = Instruction{SBC, opcode, zeropageX, 2, 4, 0}
  case 0xED: i = Instruction{SBC, opcode, absolute, 3, 4, 0}
  case 0xFD: i = Instruction{SBC, opcode, absoluteX, 3, 4, 0}
  case 0xF9: i = Instruction{SBC, opcode, absoluteY, 3, 4, 0}
  case 0xE1: i = Instruction{SBC, opcode, indirectX, 2, 6, 0}
  case 0xF1: i = Instruction{SBC, opcode, indirectY, 2, 5, 0}
  case 0x38: i = Instruction{SEC, opcode, implied, 1, 2, 0}
  case 0xF8: i = Instruction{SED, opcode, implied, 1, 2, 0}
  case 0x78: i = Instruction{SEI, opcode, implied, 1, 2, 0}
  case 0x85: i = Instruction{STA, opcode, zeropage, 2, 3, 0}
  case 0x95: i = Instruction{STA, opcode, zeropageX, 2, 4, 0}
  case 0x8D: i = Instruction{STA, opcode, absolute, 3, 4, 0}
  case 0x9D: i = Instruction{STA, opcode, absoluteX, 3, 5, 0}
  case 0x99: i = Instruction{STA, opcode, absoluteY, 3, 5, 0}
  case 0x81: i = Instruction{STA, opcode, indirectX, 2, 6, 0}
  case 0x91: i = Instruction{STA, opcode, indirectY, 2, 6, 0}
  case 0x86: i = Instruction{STX, opcode, zeropage, 2, 3, 0}
  case 0x96: i = Instruction{STX, opcode, zeropageY, 2, 4, 0}
  case 0x8E: i = Instruction{STX, opcode, absolute, 3, 4, 0}
  case 0x84: i = Instruction{STY, opcode, zeropage, 2, 3, 0}
  case 0x94: i = Instruction{STY, opcode, zeropageX, 2, 4, 0}
  case 0x8C: i = Instruction{STY, opcode, absolute, 3, 4, 0}
  case 0xAA: i = Instruction{TAX, opcode, implied, 1, 2, 0}
  case 0xA8: i = Instruction{TAY, opcode, implied, 1, 2, 0}
  case 0xBA: i = Instruction{TSX, opcode, implied, 1, 2, 0}
  case 0x8A: i = Instruction{TXA, opcode, implied, 1, 2, 0}
  case 0x9A: i = Instruction{TXS, opcode, implied, 1, 2, 0}
  case 0x98: i = Instruction{TYA, opcode, implied, 1, 2, 0}
  }
  return &i;
}

func (i *Instruction) String() string {
  return fmt.Sprintf("Instruction[%s op:%02X addr:%d bytes:%d cycles:%d]",
    i.name(), i.opcode, i.addressing, i.bytes, i.cycles)
}

func (i *Instruction) name() (s string) {
  switch i.id {
  case ADC: s = "ADC"
  case AND: s = "AND"
  case ASL: s = "ASL"
  case BCC: s = "BCC"
  case BCS: s = "BCS"
  case BEQ: s = "BEQ"
  case BIT: s = "BIT"
  case BMI: s = "BMI"
  case BNE: s = "BNE"
  case BPL: s = "BPL"
  case BRK: s = "BRK"
  case BVC: s = "BVC"
  case BVS: s = "BVS"
  case CLC: s = "CLC"
  case CLD: s = "CLD"
  case CLI: s = "CLI"
  case CLV: s = "CLV"
  case CMP: s = "CMP"
  case CPX: s = "CPX"
  case CPY: s = "CPY"
  case DEC: s = "DEC"
  case DEX: s = "DEX"
  case DEY: s = "DEY"
  case EOR: s = "EOR"
  case INC: s = "INC"
  case INX: s = "INX"
  case INY: s = "INY"
  case JMP: s = "JMP"
  case JSR: s = "JSR"
  case LDA: s = "LDA"
  case LDX: s = "LDX"
  case LDY: s = "LDY"
  case LSR: s = "LSR"
  case NOP: s = "NOP"
  case ORA: s = "ORA"
  case PHA: s = "PHA"
  case PHP: s = "PHP"
  case PLA: s = "PLA"
  case PLP: s = "PLP"
  case ROL: s = "ROL"
  case ROR: s = "ROR"
  case RTI: s = "RTI"
  case RTS: s = "RTS"
  case SBC: s = "SBC"
  case SEC: s = "SEC"
  case SED: s = "SED"
  case SEI: s = "SEI"
  case STA: s = "STA"
  case STX: s = "STX"
  case STY: s = "STY"
  case TAX: s = "TAX"
  case TAY: s = "TAY"
  case TSX: s = "TSX"
  case TXA: s = "TXA"
  case TXS: s = "TXS"
  case TYA: s = "TYA"
  }
  return
}
