#!/bin/bash

SRCDIR=$(pwd)

# Go options
export GOPATH="$GOPATH:$SRCDIR"
export GO111MODULE="off"

# Build windows
win () {
	export GOOS=windows
	go build -o build/${1}.exe src/${1}.go
}

# Build linux
linux () {
	export GOOS=linux
	go build -o build/${1} src/${1}.go
}

# Build mac
mac () {
	export GOOS=darwin
	go build -o build/${1} src/${1}.go
}

# Build all
all () {
	win $1
	mac $1
	linux $1
}

# Help text
print_help () {
	echo ""
	echo "Usage: gobuild.sh [os] [program] <arch>"
	echo ""
	echo "   os: (mac, linux, win, all)"
	echo " prog: (tint, parse)"
	echo " arch: any supported go arch for the target os"
	echo ""
}

# Check if given os is valid
is_os () {
	if [[($1 == "mac") || ($1 == "linux") || ($1 == "win") || ($1 == "all")]]; then
		return 0
	fi
	return 1
}

if [[ -n $3 ]]; then
	export GOARCH=$3
fi

if $(is_os $1); then
	if [[ -z $2 ]]; then
		$1 tint
		$1 parse
	else
		$1 $2
	fi
else
	print_help
fi
