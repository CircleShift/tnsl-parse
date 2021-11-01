/*
   Copyright 2020 Kyle Gunger

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package tparse

// LINESEP represents a line seperator
const LINESEP = 0

// INLNSEP represents an inline seperator
const INLNSEP = 1

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

// PREWORDS represents all the pre-processor directives
var PREWORDS = []string{
	"include",
	"define",
	"extern",
	"size",
	"align",
	"address",
	"rootfile",
	"if",
	"else",
	"abi",
	//"mark",
	"using",
}

func checkPreWord(s string) int {
	for _, str := range PREWORDS {
		if str == s {
			return PREWORD
		}
	}

	return -1
}

// RESWORD represents all the reserved words and what type of tokens they are.
var RESWORD = map[string]int{
	"bool":  KEYTYPE,
	"char": KEYTYPE,
	"charp": KEYTYPE,

	"int":    KEYTYPE,
	"int8":   KEYTYPE,
	"int16":  KEYTYPE,
	"int32":  KEYTYPE,
	"int64":  KEYTYPE,
	"uint":   KEYTYPE,
	"uint8":  KEYTYPE,
	"uint16": KEYTYPE,
	"uint32": KEYTYPE,
	"uint64": KEYTYPE,

	"float":   KEYTYPE,
	"float32": KEYTYPE,
	"float64": KEYTYPE,

	"void": KEYTYPE,
	"type": KEYTYPE,

	"struct":    KEYWORD,
	"interface": KEYWORD,
	"enum":      KEYWORD,
	"is":        AUGMENT,
	"extends":   KEYWORD,

	"loop":     KEYWORD,
	"continue": KEYWORD,
	"break":    KEYWORD,
	"return":   KEYWORD,

	"match":   KEYWORD,
	"case":    KEYWORD,
	"default": KEYWORD,

	"label": KEYWORD,
	"goto":  KEYWORD,

	"if":   KEYWORD,
	"else": KEYWORD,

	"const":    KEYWORD,
	"static":   KEYWORD,
	"volatile": KEYWORD,

	"method":   KEYWORD,
	"override": KEYWORD,
	"self":     LITERAL,
	"super":    LITERAL,
	"operator": KEYWORD,

	"raw":    KEYWORD,
	"asm":    KEYWORD,
	"inline": KEYWORD,

	"true":  LITERAL,
	"false": LITERAL,

	"alloc":   KEYWORD,
	"salloc":  KEYWORD,
	"realloc": KEYWORD,
	"delete":  KEYWORD,

	"module": KEYWORD,
	"export": KEYWORD,
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
	// Array/set mark open
	'{': DELIMIT,
	// Array/set mark close
	'}': DELIMIT,

	// Start of pre-proc directive
	':': LINESEP,
	// Statement seperator
	';': LINESEP,
	// Comment line
	'#': LINESEP,

	// Seperate arguments or enclosed statements
	',': INLNSEP,

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
	'`': AUGMENT,
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
	// Preproc block
	"/:": DELIMIT,
	":/": DELIMIT,

	// Redef blocks
	";;": DELIMIT,
	"::": DELIMIT,
	";#": DELIMIT,
	":#": DELIMIT,
	"#;": DELIMIT,
	"#:": DELIMIT,

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
	"+=": AUGMENT,
	"-=": AUGMENT,
	"*=": AUGMENT,
	"/=": AUGMENT,
	"%=": AUGMENT,
	"~=": AUGMENT,
	"`=": AUGMENT,

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

	// Increment and De-increment
	"++": AUGMENT,
	"--": AUGMENT,

	"len": AUGMENT,
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

func checkToken(s string, pre bool) int {
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

	if pre {
		o = checkPreWord(s)
	}

	if o > -1 {
		return o
	}

	o = checkRuneGroup(s)

	if o > -1 {
		return o
	}

	return DEFWORD
}
