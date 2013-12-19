go6502
======

```
 | | | | | | | | | | | | | | | | | | | |
.----------------------------------------.
|                   GO                   |
|                  6502                  |
|                  1213                  |
`----------------------------------------'
 | | | | | | | | | | | | | | | | | | | |
```

A [go][golang]-based emulator for the [6502][6502]-based
[pda6502 homebrew computer][pda6502].


Background
----------

This started as a golang port of my [very incomplete Ruby 6502/6510/c64
emulator][c64.rb].

Since then, I've started working on an actual [6502-based homebrew
computer][pda6502], including designing the address decoder, RAM/ROM/IO memory
layout etc.

go6502 has become the emulator for that system, but has a flexible address
bus which could be repurposed to other 6502-based systems.


Running it
----------

* Drop an 8 KB `kernal.rom` into the `rom/` directory.
    * ([pda6502][pda6502] can help; see `memory.conf` and `Makefile`)
* `GOPATH=$(pwd) go run main.go`
    * (You can `brew install go` if you use [Homebrew][homebrew] on a Mac)


License
-------

Copyright 2013 Paul Annesley, released under MIT license.


[6502]: http://en.wikipedia.org/wiki/MOS_Technology_6502
[golang]: http://golang.org/
[c64.rb]: https://github.com/pda/c64.rb
[pda6502]: https://github.com/pda/pda6502
[homebrew]: http://brew.sh/
