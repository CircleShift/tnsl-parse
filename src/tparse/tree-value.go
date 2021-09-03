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

var UNARY = map[string]int {
	"~": 0,
	"`": 0,
	"++": 2,
	"--": 2,
	"!": 6,

}

var ORDER = map[string]int{
	// Get
	".": 1,

	"is": 2,

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

	// Assignement
	"=": 8,
}

// Works? Please test. 
func parseUnaryOps(tokens *[]Token, tok, max int) (Node) {
	out := Node{Data: Token{Type: 10, Data: "value"}, IsBlock: false}
	val := false

	// Pre-value op scan
	for ; tok < max && !val; tok++ {
		t := (*tokens)[tok]
		switch t.Type {
		case DELIMIT:
			var tmp Node
			switch t.Data {
			case "{", "(": // Array or struct evaluation, parenthetical value
				tmp, tok = parseValueList(tokens, tok, max)
				out.Sub = append(out.Sub, tmp)
				val = true
			default:
				errOut("Unexpected delimiter when parsing value", t)
			}
		case LITERAL, DEFWORD:
			out.Sub = append(out.Sub, Node{Data: t, IsBlock: false})
			val = true
		case AUGMENT:
			_, prs := UNARY[t.Data]
			if !prs {
				errOut("Parser bug!  Operator failed to load into AST.", t)
			} else {
				out.Sub = append(out.Sub, Node{Data: t, IsBlock: false})
			}
		default:
			errOut("Unexpected token in value declaration", t)
		}
	}

	// Sanity check: make sure there's actually a value here
	if !val {
		errOut("Expected to find value, but there wasn't one", (*tokens)[max])
	}

	// Post-value op scan
	for ; tok < max; tok++ {
		t := (*tokens)[tok]
		switch t.Type {
		case DELIMIT:
			var tmp Node
			switch t.Data {
			case "(": // Function call
				tmp, tok = parseValueList(tokens, tok, max)
			case "[": // Typecasting
				tmp, tok = parseTypeList(tokens, tok, max)
			case "{": // Array indexing
				tmp = Node{Data: Token{Type: 10, Data: "index"}}
				var tmp2 Node
				tmp2, tok = parseValue(tokens, tok + 1, max)
				tmp.Sub = append(tmp.Sub, tmp2)
			default:
				errOut("Unexpected delimiter when parsing value", t)
			}
			out.Sub = append(out.Sub, tmp)
		case AUGMENT:
			_, prs := UNARY[t.Data]
			if !prs {
				errOut("Parser bug!  Operator failed to load into AST.", t)
			} else {
				out.Sub = append(out.Sub, Node{Data: t, IsBlock: false})
			}
		default:
			errOut("Unexpected token in value declaration", t)
		}
	}

	return out
}

// Works? Please test.
func parseBinaryOp(tokens *[]Token, tok, max int) (Node) {
	out := Node{IsBlock: false}
	first := tok
	var high, highOrder, bincount int = first, 8, 0
	var curl, brak, parn int = 0, 0, 0

	// Find first high-order op
	for ; tok < max; tok++ {
		t := (*tokens)[tok]
		if t.Type == DELIMIT {
			switch t.Data {
			case "{":
				curl++
			case "[":
				brak++
			case "(":
				parn++
			
			case "}":
				curl--
			case "]":
				brak--
			case ")":
				parn--
			}

			if curl < 0 || brak < 0 || parn < 0 {
				if curl > 0 || brak > 0 || parn > 0 {
					errOut("Un-matched closing delimiter when parsing a type.", t)
				}
			}
		} else if t.Type == AUGMENT {
			order, prs := ORDER[t.Data]
			if !prs || curl > 0 || brak > 0 || parn > 0 {
				continue
			} else if order > highOrder {
				high, highOrder = tok, order
			}
			// TODO: Add in case for the "is" operator
			bincount++
		}
	}

	out.Data = (*tokens)[high]

	if bincount == 0 {
		// No binops means we have a value to parse.  Parse all unary ops around it.
		return parseUnaryOps(tokens, first, max)
	} else {
		// Recursive split to lower order operations
		out.Sub = append(out.Sub, parseBinaryOp(tokens, first, high))
		out.Sub = append(out.Sub, parseBinaryOp(tokens, high + 1, max))
	}

	return out
}

// Works? Please test.
func parseValue(tokens *[]Token, tok, max int) (Node, int) {
	first := tok
	var curl, brak, parn int = 0, 0, 0

	for ; tok < max; tok++ {
		t := (*tokens)[tok]
		switch t.Type {
		case LINESEP:
			if curl > 0 || brak > 0 || parn > 0 {
				errOut("Encountered end of statement before all delimiter pairs were closed while looking for the end of a value.", t)
			}
			goto PARSEBIN
		case INLNSEP:
			if curl > 0 || brak > 0 || parn > 0 {
				continue
			}
			goto PARSEBIN
		case DELIMIT:
			switch t.Data {
			case "{":
				curl++
			case "[":
				brak++
			case "(":
				parn++
			
			case "}":
				curl--
			case "]":
				brak--
			case ")":
				parn--
			}

			// TODO: Support blocks as values

			if curl < 0 || brak < 0 || parn < 0 {
				if curl > 0 || brak > 0 || parn > 0 {
					errOut("Un-matched closing delimiter when parsing a value.", t)
				} else if curl + brak + parn == -1 {
					goto PARSEBIN
				} else {
					errOut("Strange bracket values detected when parsing value.  Possibly a parser bug.", t)
				}
			}
		}
	}

	PARSEBIN:

	return parseBinaryOp(tokens, first, tok), tok
}

// Works? Please test.
func parseTypeParams(tokens *[]Token, tok, max int) (Node, int) {
	out := Node{Data: (*tokens)[tok], IsBlock: false}
	tok++

	for ; tok < max; tok++ {
		t := (*tokens)[tok]
		tmp := Node{IsBlock: false}
		switch t.Type {
		case DELIMIT:
			if tok < max-1 {
				if t.Data == "(" {
					tmp, tok = parseValueList(tokens, tok, max)
				} else if t.Data == "[" {
					tmp, tok = parseTypeList(tokens, tok, max)
				} else if t.Data == ")" || t.Data == "]" || t.Data == "}" {
					// End of type
					tok--
					goto VOIDDONE
				} else {
					errOut("Error: unexpected delimeter when parsing type", t)
				}
			} else if tok >= max-1 {
				errOut("Error: unexpected end of file when parsing type", t)
			}

		default:
			// End of type
			tok--
			goto VOIDDONE
		}

		out.Sub = append(out.Sub, tmp)
	}

	VOIDDONE:

	return out, tok
}

// TODO: make sure this actually works
func parseType(tokens *[]Token, tok, max int, param bool) (Node, int) {
	out := Node{Data: Token{Type: 10, Data: "type"}}

	for ; tok < max; tok++ {
		t := (*tokens)[tok]
		var tmp Node
		switch t.Type {
		case AUGMENT:
			if t.Data != "~" && t.Data != "`" {
				errOut("Error: unexpected augment token when parsing type", t)
			}
			tmp.Data = t

		case KEYTYPE:
			if t.Data == "void" {
				tmp, tok = parseTypeParams(tokens, tok, max)
			} else {
				tmp.Data = t
			}
			out.Sub = append(out.Sub, tmp)
			return out, tok
		case DEFWORD:
			if (*tokens)[tok+1].Data == "(" {
				tmp, tok = parseTypeParams(tokens, tok, max)
			}

			return out, tok

		case KEYWORD:
			if param && t.Data == "static" {
				// Nonstandard keyword in parameter definition
				errOut("Error: parameter or value types cannot be static", t)
			} else if t.Data != "const" && t.Data != "volatile" && t.Data != "static" {
				// Nonstandard keyword in variable definition
				errOut("Error: unexpected keyword when parsing type", t)
			}
			tmp.Data = t

		case DELIMIT:
			if tok < max-1 {
				if t.Data == "{" {
					// What happens when an array type is defined
					tmp.Data = Token{AUGMENT, "{}", t.Line, t.Char}
					if (*tokens)[tok+1].Data == "}" {
						// Length variable array, add no sub-nodes and increment
						tok++
					} else {
						// Constant length array.  Parse value for length and increment
						var tmp2 Node
						tmp2, tok = parseValue(tokens, tok + 1, max)
						tmp.Sub = append(tmp.Sub, tmp2)
					}
				} else {
					errOut("Error: unexpected delimeter when parsing type", t)
				}
			} else {
				// End of file with open delimiter after type parsing has begun
				errOut("Error: unexpected end of file when parsing type", t)
			}

		default:
			errOut("Error: unexpected token when parsing type", t)
		}

		out.Sub = append(out.Sub, tmp)
	}
	
	errOut("End of token list when trying to parse type", (*tokens)[max - 1])

	return out, tok
}