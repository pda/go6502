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

Set up Go:

* `brew install go` / `aptitude install golang` / whatever.
* Spend a few days making sense of, fighting against, and eventually
  tollerating the golang directory structure as described at
  http://golang.org/doc/code.html
* Set your `$GOPATH` the way Go wants you to, e.g. `$HOME/code/go`.
* Put Go's bin dir in your path, e.g. `$HOME/code/go/bin`

Get and run go6502:

* Drop an 8 KB `kernal.rom` into `$PWD/rom/`, where ever that may be.
    * ([pda6502][pda6502] can help; see `memory.conf` and `Makefile`)
* `go get github.com/pda/go6502`
* `go6502`
* `go6502 --help`
* `go6502 --debug`


Debugger / Monitor
------------------

Given there's almost no I/O, you'll probably want a debugger / monitor session.

```
$ go6502 --debug
CPU pc:0xE000 ac:0x00 x:0x00 y:0x00 sp:0xFF sr:--------
Instruction[SEI op:78 addr:6 bytes:1 cycles:2] op8:0x00 op16:0x0000
$E000> help

pda6502 debuger
---------------
break-address <addr> (alias: ba) e.g. ba 0x1000
break-instruction <mnemonic> (alias: bi) e.g. bi NOP
break-register <x|y|a> <value> (alias: br) e.g. br x 128
exit (alias: quit, q) Shut down the emulator.
help (alias: h, ?) This help.
read <address> - Read and display 8-bit integer at address.
read16 <address> - Read and display 16-bit integer at address.
run (alias: r) Run continuously until breakpoint.
step (alias: s) Run only the current instruction.
(blank) Repeat the previous command.

Hex input formats: 0x1234 $1234
Commands expecting uint16 treat . as current address (PC).
$E000> step
CPU pc:0xE001 ac:0x00 x:0x00 y:0x00 sp:0xFF sr:--------
Instruction[LDX op:A2 addr:5 bytes:2 cycles:2] op8:0xFF op16:0x0000
$E001> break-instruction NOP
$E001> break-address $E003
$E001> run
Breakpoint for PC address = $E003
CPU pc:0xE003 ac:0x00 x:0xFF y:0x00 sp:0xFF sr:n-------
Instruction[TXS op:9A addr:6 bytes:1 cycles:2] op8:0x00 op16:0x0000
$E003> run
Breakpoint for instruction NOP
CPU pc:0xE0FC ac:0xFF x:0x00 y:0x00 sp:0xFF sr:-----izc
Instruction[NOP op:EA addr:6 bytes:1 cycles:2] op8:0x00 op16:0x0000
$E0FC> read $FFFC
$FFFC => $00 0b00000000 0 '\x00'
$E0FC> read16 $FFFC
$FFFC,FFFD => $E000 0b1110000000000000 57344
$E0FC> quit
```


License
-------

Copyright 2013–2014 Paul Annesley, released under MIT license.


[6502]: http://en.wikipedia.org/wiki/MOS_Technology_6502
[golang]: http://golang.org/
[c64.rb]: https://github.com/pda/c64.rb
[pda6502]: https://github.com/pda/pda6502
[homebrew]: http://brew.sh/
