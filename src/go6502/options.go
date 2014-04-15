package go6502

import (
	"flag"
)

type Options struct {
	viaDumpBinary bool
	viaDumpAscii  bool
}

func ParseOptions() *Options {
	options := &Options{}
	flag.BoolVar(&options.viaDumpBinary, "via-dump-binary", false, "VIA6522 dumps binary output")
	flag.BoolVar(&options.viaDumpAscii, "via-dump-ascii", false, "VIA6522 dumps ASCII output")
	flag.Parse()
	return options
}
