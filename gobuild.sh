#!/bin/bash

SRCDIR=$(pwd)

export GOPATH="$GOPATH:$SRCDIR"

go build -o build/parse src/main.go
