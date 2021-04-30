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

// Ops order in TNSL
// Cast/Paren > Address > Get > Inc/Dec > Math > Bitwise > Logic

var ORDER = map[string]int{
	// Address of
	"~": 0,
	// De-ref
	"`": 0,

	// Get
	".": 1,

	// Inc/Dec
	"++": 2,
	"--": 2,

	// Multiplication
	"*": 3,
	// Division
	"/": 3,

	// Addition
	"+": 4,
	// Subtraction
	"-": 4,
	
	// Mod
	"%": 5,

	// Bitwise and
	"&": 6,
	// Bitwise or
	"|": 6,
	// Bitwise xor
	"^": 6,

	// Bitwise l-shift
	"<<": 6,
	// Bitwise r-shift
	">>": 6,

	"!&": 6,
	"!|": 6,
	"!^": 6,

	// Not (prefix any bool or bitwise)
	"!": 6,

	// Boolean and
	"&&": 7,
	// Boolean or
	"||": 7,
	// Truthy equals
	"==": 7,

	// Greater than
	">": 7,
	// Less than
	"<": 7,

	"!&&": 7,
	"!||": 7,
	"!==": 7,

	"!>": 7,
	"!<": 7,
}

func parseValue(tokens *[]Token, tok, max int) (Node, int) {
	out := Node{}

	for ; tok < max; tok++ {
		t := (*tokens)[tok]
		switch t.Type {
		case LITERAL:
		case DEFWORD:
		case DELIMIT:
		}
	}

	return out, tok
}

func parseVoidType(tokens *[]Token, tok, max int) (Node, int) {
	out := Node{}
	working := &out

	for ; tok < max; tok++ {
		t := (*tokens)[tok]
		switch t.Type {
		case AUGMENT:
			if t.Data != "~" && t.Data != "`" {
				errOut("Error: unexpected augment token when parsing type", t)
			}
			working.Data = t

		case KEYTYPE:
			if t.Data == "void" {
				*working, tok = parseVoidType(tokens, tok, max)
			} else {
				working.Data = t
			}

			return out, tok

		case DEFWORD:

		case DELIMIT:
			if t.Data == "{" && tok < max-1 {
				if (*tokens)[tok+1].Data == "}" {
					working.Data = Token{AUGMENT, "{}", t.Line, t.Char}
					tok++
				} else {
					errOut("Error: start of list when parsing type (did you mean \"{}\"?)", t)
				}
			} else if tok >= max-1 {
				errOut("Error: unexpected end of file when parsing type", t)
			} else {
				errOut("Error: unexpected delimeter when parsing type", t)
			}

		default:
			errOut("Error: unexpected token when parsing type", t)
		}

		makeParent(working, Node{})
		working = &(working.Sub[0])
	}

	return out, tok
}

func parseType(tokens *[]Token, tok, max int, param bool) (Node, int) {
	out := Node{}
	working := &out

	for ; tok < max; tok++ {
		t := (*tokens)[tok]
		switch t.Type {
		case AUGMENT:
			if t.Data != "~" && t.Data != "`" {
				errOut("Error: unexpected augment token when parsing type", t)
			}
			working.Data = t

		case KEYTYPE:
			if t.Data == "void" {
				*working, tok = parseVoidType(tokens, tok, max)
			} else {
				working.Data = t
			}

			return out, tok

		case DEFWORD:
			if (*tokens)[tok+1].Data == "(" {

			}

		case KEYWORD:
			if param && t.Data == "static" {
				// Nonstandard keyword in parameter definition
				errOut("Error: parameter types cannot be static", t)
			} else if t.Data != "const" && t.Data != "volatile" && t.Data != "static" {
				// Nonstandard keyword in variable definition
				errOut("Error: unexpected keyword when parsing type", t)
			}
			working.Data = t

		case DELIMIT:
			if t.Data == "{" && tok < max-1 {
				// What happens when an array type is defined
				if (*tokens)[tok+1].Data == "}" {
					// Length variable array
					working.Data = Token{AUGMENT, "{}", t.Line, t.Char}
					tok++
				} else if (*tokens)[tok+1].Type == LITERAL {
					// Array with constant length
				} else {
					// Undefined behaviour
					errOut("Error: start of list when parsing type (did you mean \"{}\"?)", t)
				}
			} else if tok >= max-1 {
				// End of file with open delimiter after type parsing has begun
				errOut("Error: unexpected end of file when parsing type", t)
			} else {
				// Other delimiter than {} used in variable definition
				errOut("Error: unexpected delimeter when parsing type", t)
			}

		default:
			errOut("Error: unexpected token when parsing type", t)
		}

		makeParent(working, Node{})
		working = &(working.Sub[0])
	}

	return out, tok
}
