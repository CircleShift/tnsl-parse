#!/bin/bash

SRCDIR=$(pwd)

GOPATH="$GOPATH:$SRCDIR"
GO111MODULE=off

go build -o build/${1} src/${1}.go
