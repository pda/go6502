package go6502

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type Options struct {
	viaDumpBinary bool
	viaDumpAscii  bool
	ViaSsd1306    bool
	LogFile       string
	Debug         bool
	DebugCmds     commandList
}

func ParseOptions() *Options {
	opt := &Options{}

	// Debug
	flag.BoolVar(&opt.Debug, "debug", false, "Run debugger")
	flag.Var(&opt.DebugCmds, "debug-commands", "Debugger commands to run, semicolon separated.")

	// Logging
	flag.StringVar(&opt.LogFile, "log-file", os.DevNull, "Emulator debug log")

	// VIA
	flag.BoolVar(&opt.viaDumpBinary, "via-dump-binary", false, "VIA6522 dumps binary output")
	flag.BoolVar(&opt.viaDumpAscii, "via-dump-ascii", false, "VIA6522 dumps ASCII output")
	flag.BoolVar(&opt.ViaSsd1306, "via-ssd1306", false, "SSD1306 OLED display on VIA6522 port B")

	flag.Parse()
	return opt
}

type commandList []string

func (cl *commandList) Set(value string) error {
	list := strings.Split(value, ";")
	for i, value := range list {
		list[i] = strings.TrimSpace(value)
	}
	*cl = list
	return nil
}

func (cl *commandList) String() string {
	return fmt.Sprint(*cl)
}
