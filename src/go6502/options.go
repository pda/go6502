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
	opt := &Options{}

	// Logging
	flag.StringVar(&opt.LogFile, "log-file", os.DevNull, "Emulator debug log")

	// VIA
	flag.BoolVar(&opt.viaDumpBinary, "via-dump-binary", false, "VIA6522 dumps binary output")
	flag.BoolVar(&opt.viaDumpAscii, "via-dump-ascii", false, "VIA6522 dumps ASCII output")

	flag.Parse()
	return opt
}
