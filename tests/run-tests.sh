PARSECMD=../build/parse
PARSEFILE=" "

parse () {
    $PARSECMD -in $1-test.tnsl -out $1-test.tnt
}

parse block
parse comment
parse literal
parse parameter
parse statement