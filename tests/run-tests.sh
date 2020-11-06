PARSECMD=../build/parse
PARSEFILE=" "

parse () {
    $PARSECMD -in $PARSEFILE-test.tnsl -out $PARSEFILE-test.tnt
}

PARSEFILE=block
parse

PARSEFILE=comment
parse

PARSEFILE=literal
parse

PARSEFILE=parameter
parse