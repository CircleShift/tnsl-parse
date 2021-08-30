#!/bin/bash

SRCDIR=$(pwd)

export GOPATH="$GOPATH:$SRCDIR"
export GO111MODULE="off"

go build -o build/${1} src/${1}.go
