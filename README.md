# tnsl-parse

The tokenizer for the [TNSL Language](https://github.com/CoreChg/tnsl-lang).  Written in Go for the moment.

The goal is to get this part written, and then write a backend for any arch (x86, arm, risc-v, etc.).
After that, work can begin on the real parser/compiler which will be written in native TNSL

## interpreter

This project was originally supposed to form a go based compiler for the language, but in the interest of time, it seems more efficient to build an interpreter instead so we can work on the TNSL based compiler sooner.

To build the parser:

    ./gobuild.sh parse

To build the interpreter:

    ./gobuild.sh exec

Binaries will be dumped in the "build" folder.

## Status:

    Parser: broke
    Interpreter: not started