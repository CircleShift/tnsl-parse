package tparse

import ()

// LINESEP represents a line seperator
const LINESEP = 0

// ARGNSEP represents an inline seperator
const ARGNSEP = 1

// DELIMIT represents an opening or closing delimiter
const DELIMIT = 2

// AUGMENT represents an augmentation
const AUGMENT = 3

// LITERAL represents a literal value
const LITERAL = 4

// KEYTYPE represents a built in type
const KEYTYPE = 5

// PREWORD represents a reserved pre-processor directive
const PREWORD = 6

// KEYWORD represents a reserved word
const KEYWORD = 7

// DEFWORD represents a user-defined word such as a variable, method, or struct
const DEFWORD = 8

// RESWORD represents all the reserved words and what type of tokens they are.
var RESWORD = map[string]int{
	"import": PREWORD,

	"bool": KEYTYPE,
	"byte": KEYTYPE,
	"char": KEYTYPE,

	"int":   KEYTYPE,
	"float": KEYTYPE,

	"struct": KEYWORD,
	"type":   KEYWORD,

	"loop":     KEYWORD,
	"continue": KEYWORD,
	"break":    KEYWORD,

	"switch":  KEYWORD,
	"case":    KEYWORD,
	"default": KEYWORD,

	"label": KEYWORD,
	"goto":  KEYWORD,

	"if":   KEYWORD,
	"else": KEYWORD,

	"const":    KEYWORD,
	"static":   KEYWORD,
	"volatile": KEYWORD,

	"true":  LITERAL,
	"false": LITERAL,

	"null": LITERAL,
}

func checkResWord(s string) int {
	out, prs := RESWORD[s]
	if !prs {
		return -1
	}
	return out
}

// RESRUNE represents all the reserved runes
var RESRUNE = map[rune]int{
	// Starting condition open
	'(': DELIMIT,
	// Starting condition close
	')': DELIMIT,
	// Ending condition open
	'[': DELIMIT,
	// Ending condition close
	']': DELIMIT,
	// Array mark open
	'{': DELIMIT,
	// Array mark close
	'}': DELIMIT,

	// Start of pre-proc directive
	':': LINESEP,
	// Start of line
	';': LINESEP,
	// Comment line
	'#': LINESEP,

	// Seperate arguments
	',': ARGNSEP,

	// Assignment
	'=': AUGMENT,

	// Get
	'.': AUGMENT,

	// Bitwise and
	'&': AUGMENT,
	// Bitwise or
	'|': AUGMENT,
	// Bitwise xor
	'^': AUGMENT,

	// Greater than
	'>': AUGMENT,
	// Less than
	'<': AUGMENT,

	// Not (prefix any bool or bitwise)
	'!': AUGMENT,

	// Addition
	'+': AUGMENT,
	// Subtraction
	'-': AUGMENT,
	// Multiplication
	'*': AUGMENT,
	// Division
	'/': AUGMENT,
	// Mod
	'%': AUGMENT,

	// Address of
	'~': AUGMENT,
	// De-ref
	'_': AUGMENT,
}

func checkResRune(r rune) int {
	out, prs := RESRUNE[r]
	if !prs {
		return -1
	}
	return out
}

// RESRUNES Reserved sets of reserved runes which mean something
var RESRUNES = map[string]int{
	// Code block
	"/;": DELIMIT,
	";/": DELIMIT,
	// Comment block
	"/#": DELIMIT,
	"#/": DELIMIT,

	";;": DELIMIT,

	// Boolean equ
	"==": AUGMENT,
	// Boolean and
	"&&": AUGMENT,
	// Boolean or
	"||": AUGMENT,

	// Bitwise l-shift
	"<<": AUGMENT,
	// Bitwise r-shift
	">>": AUGMENT,

	// PREaugmented augmentors
	"&=": AUGMENT,
	"|=": AUGMENT,
	"^=": AUGMENT,
	"!=": AUGMENT,
	"+=": AUGMENT,
	"-=": AUGMENT,
	"*=": AUGMENT,
	"/=": AUGMENT,
	"%=": AUGMENT,
	"~=": AUGMENT,
	"_=": AUGMENT,

	// POSTaugmented augmentors
	"!&":  AUGMENT,
	"!|":  AUGMENT,
	"!^":  AUGMENT,
	"!==": AUGMENT,
	"!&&": AUGMENT,
	"!||": AUGMENT,
	"!>":  AUGMENT,
	"!<":  AUGMENT,
	">==": AUGMENT,
	"<==": AUGMENT,
}

func maxResRunes() int {
	max := 0

	for k := range RESRUNES {
		if len(k) > max {
			max = len(k)
		}
	}

	return max
}

func checkRuneGroup(s string) int {
	rs := StringAsRunes(s)

	if len(rs) == 1 {
		return checkResRune(rs[0])
	}

	out, prs := RESRUNES[s]
	if !prs {
		return -1
	}
	return out
}

func checkToken(s string) int {
	rs := StringAsRunes(s)

	if len(rs) == 0 {
		return -1
	}

	if len(rs) == 1 {
		o := checkResRune(rs[0])
		if o > -1 {
			return o
		}
	}

	o := checkResWord(s)

	if o > -1 {
		return o
	}

	o = checkRuneGroup(s)

	if o > -1 {
		return o
	}

	return DEFWORD
}
