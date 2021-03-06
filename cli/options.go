/*
	Package cli provides command line support for go6502.

	It parses CLI flags and exposes the resulting options.
*/
package cli

import (
	"flag"
	"fmt"
	"strings"
)

// Options stores the value of command line options after they're parsed.
type Options struct {
	Debug           bool
	DebugCmds       commandList
	DebugSymbolFile string
	Ili9340         bool
	SdCard          string
	Speedometer     bool
	ViaDumpAscii    bool
	ViaDumpBinary   bool
	ViaSsd1306      bool
}

// ParseFlags uses the flag stdlib package to parse CLI options.
func ParseFlags() *Options {
	opt := &Options{}

	flag.BoolVar(&opt.Debug, "debug", false, "Run debugger")
	flag.Var(&opt.DebugCmds, "debug-commands", "Debugger commands to run, semicolon separated.")
	flag.StringVar(&opt.DebugSymbolFile, "debug-symbol-file", "", "ld65 debug file to load.")
	flag.StringVar(&opt.SdCard, "sd-card", "", "Load file as SD card")
	flag.BoolVar(&opt.Speedometer, "speedometer", false, "Measure effective clock speed")
	flag.BoolVar(&opt.ViaDumpBinary, "via-dump-binary", false, "6522 dumps binary output")
	flag.BoolVar(&opt.ViaDumpAscii, "via-dump-ascii", false, "6522 dumps ASCII output")
	flag.BoolVar(&opt.ViaSsd1306, "via-ssd1306", false, "SSD1306 OLED display on 6522")
	flag.BoolVar(&opt.Ili9340, "ili9340", false, "ILI9340 TFT display on 6522")

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
