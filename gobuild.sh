#!/bin/bash

SRCDIR=$(pwd)

GOPATH="$GOPATH:$SRCDIR"

go env -w GOPATH=$GOPATH

go build -o build/parse src/main.go
