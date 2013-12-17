go6502
======

A [MOS 6502][6502] emulator in [go][golang].

I have no idea what I'm doing.


Overview
--------

This started as a port of my [very incomplete Ruby 6502/6510/c64
emulator][c64.rb].

Since then, I've started working on a [hardware 6502-based computer][pda6502],
including designing the address decoder, RAM/ROM/IO layout etc. I'll be
repurposing this repository to test layouts and code for that system. Maybe.


License
-------

Copyright 2013 Paul Annesley, released under MIT license.


[6502]: http://en.wikipedia.org/wiki/MOS_Technology_6502
[golang]: http://golang.org/
[c64.rb]: https://github.com/pda/c64.rb
[pda6502]: https://github.com/pda/pda6502
