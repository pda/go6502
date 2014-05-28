package debugger

import (
	"bufio"
	"os"
	"strconv"
)

type debugSymbol struct {
	address uint16
	name    string
}

type debugSymbols []debugSymbol

func (symbols debugSymbols) symbolsFor(addr uint16) (result []string) {
	result = make([]string, 0)
	for _, l := range symbols {
		if l.address == addr {
			result = append(result, l.name)
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
		if t.line.data["type"] == "label" {
			addr, err := strconv.ParseUint(t.line.data["value"], 0, 16)
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
				t.enter(sName)
			} else {
				t.enter(sReject)
			}
		case sReject:
			if bytes[0] == '\n' {
				t.enter(sBegin)
			}
		case sName:
			if bytes[0] != '\t' {
				text := s.Text()
				t.line.name = text[1 : len(text)-1] // strip quotes
				t.enter(sMap)
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
			t.line.data[t.line.key] = s.Text()
		}
	}
	if err = s.Err(); err != nil {
		return
	}

	return
}

// Tokenizer states.
const (
	sBegin  = iota // initial state
	sReject        // line is being rejected
	sName          // expecting name
	sMap           // expecting ,key=value,key=value
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
	//fmt.Printf("cc65 debug tokenizer: %d => %d\n", t.state, state)
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
