PARSECMD=../build/parse
PARSEFILE=" "

parse () {
    $PARSECMD -in $PARSEFILE-test.tnsl -out $PARSEFILE-test.tnp
}

PARSEFILE=block
parse

PARSEFILE=comment
parse

PARSEFILE=literal
parse