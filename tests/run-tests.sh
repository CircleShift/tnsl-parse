PARSECMD=../build/parse
PARSEFILE=" "

parse () {
    $PARSECMD -in $PARSEFILE.tnsl -out $PARSEFILE.tnp
}

PARSEFILE=block-test
parse
PARSEFILE=comment-test
parse
PARSEFILE=literal-test
parse