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

A [go][golang]-based emulator and debugger for the [6502][6502]-based
[pda6502 homebrew computer][pda6502].

[![GoDoc](https://godoc.org/github.com/pda/go6502?status.png)](https://godoc.org/github.com/pda/go6502)

Background
----------

I've been designing and building a [6502-based homebrew computer][pda6502].

It's powered by an 8-bit 6502 (WDC 65C02), varitions of which powered the
venerable Commodore 64, Apple II, Vic 20, Nintendo and lots more.

74HC-series logic chips map the 64K address space to 32K RAM, 8K ROM, a VIA
6522 I/O controller, and room for expansion.

The first output device (beyond flashing LEDs) is a [128x32 pixel OLED][oled],
connected to one of the VIA 6522 parallel ports, with bit-banged serial comms.


go6502
------

go6502 emulates the 6502, address bus, RAM, ROM, 6522 and OLED display well
enough to run the current pda6502 code and get the same display output.

It has a flexible address bus, which paves the way to emulating other
6502-based systems.

go6502 features a stepping debugger with breakpoints on instruction type,
register values and memory location. This makes it far easier to get code
working correctly before writing it to an actual EEPROM.


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


Example usage
-------------

Various invocations from my shell history; some run from this project, some from pda6502.

```sh
time go run go6502.go --debug --debug-commands='bi nop;run;q' --via-ssd1306 && open output.png
make && go6502 --debug --debug-symbol-file=build/debug --debug-commands="bi nop;c;q" --via-ssd1306 --sd-card=sd.bin
go build -o g6 go6502.go && gtimeout -s INT 0.2 ./g6 --via-ssd1306 --speedometer
time go run go6502.go --via-ssd1306 --debug
time go run go6502.go --debug --debug-commands='bi nop;r;q' --via-ssd1306 && open output.png
make && go6502 --debug --debug-symbol-file=build/debug --via-ssd1306 --sd-card=sd.bin
time go run go6502.go --via-ssd1306 --debug --debug-commands='bi nop;r;q' && open output.png
go run go6502.go --debug --debug-symbol-file=$HOME/code/pda6502/build/debug --via-ssd1306 --sd-card=$HOME/code/pda6502/sd.bin
make && go install github.com/pda/go6502 && gtimeout -s INT 0.1 go6502 --via-ssd1306 --sd-card="sd.bin" ; open ssd1306.png
go run go6502.go --via-ssd1306 --debug --debug-commands='bi nop;run'
go build go6502.go && gtimeout -s INT 1 ./go6502 --via-ssd1306 --speedometer
go build go6502.go && gtimeout -s INT 0.1 ./go6502 --via-ssd1306 --speedometer
make && go6502 -via-ssd1306 -sd-card=sd4gb.fat32 -debug -debug-symbol-file=build/debug -debug-commands="ba Halt; c; q" && hd -s 0x6000 -n 512 core
make && go6502 --debug --debug-symbol-file=build/debug --via-ssd1306 --sd-card=sd.bin --debug-commands="bi nop; c; q"
make && go6502 --debug --debug-symbol-file=build/debug --debug-commands="bi nop;c;q" --via-ssd1306 --sd-card=sd.bin && open ssd1306.png
go6502 -via-ssd1306 -sd-card=sd4gb.fat32
go6502 --debug --debug-symbol-file=build/debug --via-ssd1306 --sd-card=sd.bin
go run go6502.go --via-ssd1306 --debug
go run go6502.go --debug --debug-commands='bi nop;run;q' --via-ssd1306
```


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
continue (alias: c) Run continuously until breakpoint.
exit (alias: quit, q) Shut down the emulator.
help (alias: h, ?) This help.
read <address> - Read and display 8-bit integer at address.
read16 <address> - Read and display 16-bit integer at address.
step (alias: s) Run only the current instruction.
(blank) Repeat the previous command.

Hex input formats: 0x1234 $1234
Commands expecting uint16 treat . as current address (PC).
$E000> step
CPU pc:0xE001 ac:0x00 x:0x00 y:0x00 sp:0xFF sr:--------
Instruction[LDX op:A2 addr:5 bytes:2 cycles:2] op8:0xFF op16:0x0000
$E001> break-instruction NOP
$E001> break-address $E003
$E001> continue
Breakpoint for PC address = $E003
CPU pc:0xE003 ac:0x00 x:0xFF y:0x00 sp:0xFF sr:n-------
Instruction[TXS op:9A addr:6 bytes:1 cycles:2] op8:0x00 op16:0x0000
$E003> continue
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
[oled]: https://www.adafruit.com/products/661
