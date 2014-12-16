package debugger

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type debugSymbol struct {
	address uint16
	name    string
}

type debugSymbols []debugSymbol

func (d debugSymbol) String() string {
	return fmt.Sprintf("{%s => $%04X}", d.name, d.address)
}

// addressesFor returns the addresses labelled with the given name.
func (symbols debugSymbols) addressesFor(name string) (result []uint16) {
	for _, s := range symbols {
		if strings.EqualFold(name, s.name) {
			result = append(result, s.address)
		}
	}
	return
}

// labelsFor returns label name(s) for the given address.
func (symbols debugSymbols) labelsFor(addr uint16) (result []string) {
	for _, l := range symbols {
		if l.address == addr {
			result = append(result, l.name)
		}
	}
	return
}

// uniqueLabels is label names which resolve to a single address.
func (symbols debugSymbols) uniqueLabels() (result []string) {
	counter := make(map[string]int)
	for _, l := range symbols {
		counter[l.name]++
	}
	for l, count := range counter {
		if count == 1 {
			result = append(result, l)
		}
	}
	return
}

func readDebugSymbols(debugFile string) (symbols debugSymbols, err error) {
	file, err := os.Open(debugFile)
	if err != nil {
		return
	}

	symbols = make([]debugSymbol, 128)
	t := &tokenizer{state: sBegin}

	handleLine := func() {
		// old format: "label", new format: "lab"
		if t.line.data["type"][0:3] == "lab" {
			val, ok := t.line.data["val"] // new format
			if !ok {
				val = t.line.data["value"] // old format
			}
			addr, err := strconv.ParseUint(val, 0, 16)
			if err != nil {
				panic(err)
			}
			symbols = append(symbols, debugSymbol{address: uint16(addr), name: t.line.name})
		}
	}

	s := bufio.NewScanner(file)
	s.Split(t.splitter)
	for s.Scan() {
		bytes := s.Bytes()
		switch t.state {
		case sBegin:
			if s.Text() == "sym" {
				t.line = debugLine{prefix: "sym", data: make(map[string]string)}
				t.enter(sTab)
			} else {
				t.enter(sReject)
			}
		case sReject:
			if bytes[0] == '\n' {
				t.enter(sBegin)
			}
		case sTab:
			if bytes[0] == '\t' {
				t.enter(sNameOrMap)
			} else {
				panic("Expected TAB after line type")
			}
		case sNameOrMap:
			if bytes[0] == '"' {
				// name (old debug format)
				text := s.Text()
				t.line.name = text[1 : len(text)-1] // strip quotes
				t.enter(sMap)
			} else {
				// map key (new debug format)
				t.line.key = s.Text()
				t.enter(sMapEquals)
			}
		case sMap:
			if bytes[0] == ',' {
				t.enter(sMapKey)
			} else if bytes[0] == '\n' {
				t.enter(sBegin)
				handleLine()
			}
		case sMapKey:
			t.line.key = s.Text()
			t.enter(sMapEquals)
		case sMapEquals:
			if bytes[0] != '=' {
				panic("Expected '=' in sMapEquals state")
			} else {
				t.enter(sMapValue)
			}
		case sMapValue:
			t.enter(sMap)
			if t.line.key == "name" {
				text := s.Text()
				t.line.name = text[1 : len(text)-1] // strip quotes
			} else {
				t.line.data[t.line.key] = s.Text()
			}
		}
	}
	if err = s.Err(); err != nil {
		return
	}

	return
}

// Tokenizer states.
const (
	sBegin     = iota // initial state
	sReject           // line is being rejected
	sTab              // expect tab
	sNameOrMap        // expecting name in old format, map in new format.
	sMap              // expecting ,key=value,key=value
	sMapKey
	sMapEquals
	sMapValue
)

type debugLine struct {
	prefix string
	name   string
	key    string
	data   map[string]string
}

type tokenizer struct {
	state int
	line  debugLine
}

func (t *tokenizer) enter(state int) {
	t.state = state
}

func (t *tokenizer) splitter(data []byte, atEOF bool) (advance int, token []byte, err error) {

	// return separators as separate tokens.
	if len(data) > 0 {
		switch data[0] {
		case '\t', '\n', ',', '=':
			return 1, data[0:1], nil
		}
	}

	// return tokens up to but not including the separator.
	for i, b := range data {
		switch b {
		case '\t', '\n', ',', '=':
			return i, data[0:i], nil
		}
	}

	return 0, nil, nil
}
