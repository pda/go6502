package go6502

import (
	"flag"
	"os"
)

type Options struct {
	viaDumpBinary bool
	viaDumpAscii  bool
	LogFile       string
}

func ParseOptions() *Options {
	options := &Options{}

	// Logging
	flag.StringVar(&options.LogFile, "log-file", os.DevNull, "Emulator debug log")

	// VIA
	flag.BoolVar(&options.viaDumpBinary, "via-dump-binary", false, "VIA6522 dumps binary output")
	flag.BoolVar(&options.viaDumpAscii, "via-dump-ascii", false, "VIA6522 dumps ASCII output")

	flag.Parse()
	return options
}
