#!/bin/bash

PARSECMD=../build/parse
PARSEFILE=" "

parse () {
	echo "ATTEMPTING TO PARSE $1-test.tnsl"
    $PARSECMD $2 -in $1-test.tnsl -out $1-test.tnt
	if [ $? -eq 0 ]; then
		echo "SUCCESS!"
	fi
}

parse block "$1"
parse comment "$1"
parse literal "$1"
parse parameter "$1"
parse statement "$1"
