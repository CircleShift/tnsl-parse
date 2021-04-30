#!/bin/bash

SRCDIR=$(pwd)

export GOPATH="$GOPATH:$SRCDIR"

go env -w GO111MODULE=off
go build -o build/parse src/main.go
