# tnsl-parse

The tokenizer for the [TNSL Language](https://github.com/CircleShift/tnsl-lang).  Written in Go for the moment.

The goal is to get this part written, and then write a backend for any arch (x86, arm, risc-v, etc.).
After that, work can begin on the real parser/compiler which will be written in native TNSL

## interpreter

This project was originally supposed to form a go based compiler for the language, but in the interest of time, it seems more efficient to build an interpreter instead so we can work on the TNSL based compiler sooner.

To build the parser:

    ./gobuild.sh parse

To build the interpreter:

    ./gobuild.sh exec

Binaries will be dumped in the "build" folder.

## Status

    Parser: sorta works sometimes (subtext: AST generater is at least a little broken, but works for at least some use cases)
    Interpreter: broken

## Usage

Once you have built the parser, it can be invoked in the build folder with `./parse`.  The cli options are as follows:

- `-writelevel <0, 1, or 2>` tells the parser at what stage it should stop parsing and output a file.
	- `0` will output the list of tokens generated by the tokenizer.
	- `1` will output the Abstract Syntax Tree as generated by the parser.
	- `2` will output a "virtual program".  This is what the interpreter acts on.
	- The default value is `1`

- `-in <file>` tells the parser what file to parse. This is the only manditory option.

- `-out <file>` tells the parser where to write the data.  The default is `out.tnt`.

### Other notes

With some of the code I've written, I'm kinda supprised that this even compiles.
