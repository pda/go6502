package go6502

import (
	"flag"
)

type Options struct {
}

func ParseOptions() *Options {
	options := &Options{}
	flag.Parse()
	return options
}
