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

var UNARY_PRE = map[string]int {
	"~": 0,
	"++": 2,
	"--": 2,
	"!": 6,
	"len": 0,
	"-": 0,
}

var UNARY_POST = map[string]int {
	"`": 0,
	"++": 2,
	"--": 2,
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
	var out Node
	var vnode *Node = &out
	val, comp := false, false
	// Pre-value op scan
	for ; tok < max && !val; tok++ {
		t := (*tokens)[tok]
		switch t.Type {
		case DELIMIT:
			
			switch t.Data {
			case "{": // Array or struct evaluation, parenthetical value
				if vnode != &out {
					errOut("Composite values may not use unary operators.", out.Data)
				}
				(*vnode) = Node{Token{10, "comp", 0, 0}, []Node{Node{}}}
				(*vnode).Sub[0], tok = parseValueList(tokens, tok + 1, max)
				val = true
				comp = true
			default:
				errOut("Unexpected delimiter when parsing value", t)
			}
		case LITERAL, DEFWORD:
			(*vnode).Data = t
			val = true
		case AUGMENT:
			_, prs := UNARY_PRE[t.Data]
			if !prs {
				errOut("Parser bug!  Operator failed to load into AST.", t)
			} else {
				(*vnode) = Node{t, []Node{Node{}}}
				vnode = &((*vnode).Sub[0])
			}
		default:
			errOut("Unexpected token in value declaration", t)
		}
	}

	// Sanity check: make sure there's actually a value here
	if !val {
		errOut("Expected to find value, but there wasn't one", (*tokens)[max - 1])
	}

	// Post-value op scan
	for ; tok < max; tok++ {
		t := (*tokens)[tok]
		var tmp Node
		switch t.Type {
		case DELIMIT:
			switch t.Data {
			case "(": // Function call
				if comp {
					errOut("Composite values can not be called as functions.", t)
				}
				tmp, tok = parseValueList(tokens, tok + 1, max)
				tmp.Data.Data = "call"
			case "[": // Typecasting
				tmp, tok = parseTypeList(tokens, tok + 1, max)
				tmp.Data.Data = "cast"
			case "{": // Array indexing
				if comp {
					errOut("Inline composite values can not be indexed.", t)
				}
				tmp, tok = parseValueList(tokens, tok + 1, max)
				tmp.Data.Data = "index"
			default:
				errOut("Unexpected delimiter when parsing value", t)
			}
		case AUGMENT:
			_, prs := UNARY_POST[t.Data]
			if !prs {
				errOut("Parser bug!  Operator failed to load into AST.", t)
			}else if comp {
				errOut("Composite values are not allowed to use unary operators.", t)
			}
			tmp = Node{}
			tmp.Data = t
		default:
			errOut("Unexpected token in value declaration", t)
		}
		(*vnode).Sub = append((*vnode).Sub, tmp)
	}

	return out
}

// Works? Please test.
func parseBinaryOp(tokens *[]Token, tok, max int) (Node) {
	out := Node{}
	first := tok
	var high, highOrder, bincount int = first, 0, 0
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
			if t.Data == "-" {
				_, prs := ORDER[(*tokens)[tok - 1].Data]
				if prs || (*tokens)[tok - 1].Data == "return" {
					continue
				} else if order > highOrder {
					high, highOrder = tok, order
				}
			} else if prs == false || curl > 0 || brak > 0 || parn > 0 {
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
		// No binops means we have a pure value to parse.  Parse all unary ops around it.
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
	var curl, brak, parn, block int = 0, 0, 0, 0

	for ; tok < max; tok++ {
		t := (*tokens)[tok]
		switch t.Type {
		case LINESEP:
			if block > 0 {
				continue
			}
			if curl > 0 || brak > 0 || parn > 0 {
				errOut("Encountered end of statement before all delimiter pairs were closed while looking for the end of a value.", t)
			}
			goto PARSEBIN
		case INLNSEP:
			if curl > 0 || brak > 0 || parn > 0 || block > 0 {
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
			
			case "/;":
				_, prs := ORDER[(*tokens)[tok - 1].Data]
				if !prs && block == 0 {
					goto PARSEBIN
				}
				block++
			case ";/":
				
				if block > 0 {
					block--
				}
				fallthrough
			case ";;":
				if block > 1 {
					continue
				} else if block == 1 {
					errOut("Error: redefinition of a block from a block as a value is not permitted.", t)
				}
				if curl > 0 || brak > 0 || parn > 0 {
					errOut("Delimeter pairs not closed before end of value", t)
				}
				goto PARSEBIN
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
	out := Node{Data: (*tokens)[tok]}
	tok++

	for ; tok < max; tok++{
		t := (*tokens)[tok]
		tmp := Node{}
		switch t.Type {
		case DELIMIT:
			if tok < max {
				if t.Data == "(" {
					tmp, tok = parseTypeList(tokens, tok + 1, max)
					tmp.Data.Data = "()"
				} else if t.Data == "[" {
					tmp, tok = parseTypeList(tokens, tok + 1, max)
					tmp.Data.Data = "[]"
				} else if t.Data == ")" || t.Data == "]" || t.Data == "}" {
					// End of type
					//errOutV("Test", tok, max, t)
					goto DONE
				} else {
					errOut("Error: unexpected delimeter when parsing type", t)
				}
			} else {
				errOut("Error: unexpected end of file when parsing type", t)
			}

		default:
			// End of type
			goto DONE
		}

		out.Sub = append(out.Sub, tmp)
	}

	DONE:

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
			if t.Data != "~" {
				errOut("Error: unexpected augment token when parsing type", t)
			}
			tmp.Data = t

		case KEYTYPE:
			if t.Data == "void" {
				tmp, tok = parseTypeParams(tokens, tok, max)
			} else {
				tmp.Data = t
				tok++
			}

			out.Sub = append(out.Sub, tmp)
			
			if param && (*tokens)[tok].Data == "`" {
				tmp = Node{(*tokens)[tok], []Node{}}
				out.Sub = append(out.Sub, tmp)
				tok++
			}
			
			return out, tok
		case DEFWORD:
			if (*tokens)[tok+1].Data == "(" {
				tmp, tok = parseTypeParams(tokens, tok, max)
			} else if (*tokens)[tok+1].Data == "." {
				tmp.Data = t
				out.Sub = append(out.Sub, tmp)
				tok++
				continue
			} else {
				tmp.Data = t
				tok++
			}
			
			out.Sub = append(out.Sub, tmp)

			if param && (*tokens)[tok].Data == "`" {
				tmp = Node{(*tokens)[tok], []Node{}}
				out.Sub = append(out.Sub, tmp)
				tok++
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
						tmp2, tok = parseValueList(tokens, tok + 1, max)
						tmp.Sub = append(tmp.Sub, tmp2)
					}
				} else if t.Data == ")" || t.Data == "]" || t.Data == "}"{
					// End of type
					goto TYPEDONE
				} else {
					errOut("Error: unexpected delimeter when parsing type", t)
				}
			} else {
				// End of file with open delimiter after type parsing has begun
				errOut("Error: unexpected end of file when parsing type", t)
			}

		default:
			goto TYPEDONE
		}

		out.Sub = append(out.Sub, tmp)
	}

	TYPEDONE:

	return out, tok
}

// TODO: Check if this works. This implimentation is probably bad, but I don't care.
func isTypeThenValue(tokens *[]Token, tok, max int) (bool) {
	//TODO: check for a standard type and then a value
	var stage int = 0
	var curl, brak, parn int = 0, 0, 0

	for ; tok < max && stage < 2; tok++ {
		t := (*tokens)[tok]
		switch t.Type {
		case KEYTYPE:
			if curl > 0 || brak > 0 || parn > 0 {
				continue
			} else if stage > 0 {
				errOut("Encountered a keytype where we weren't expecting one (iTTV).", t)
			}
			stage++
		case DEFWORD:
			if curl > 0 || brak > 0 || parn > 0 {
				continue
			}
			stage++
		case LINESEP:
			if curl > 0 || brak > 0 || parn > 0 {
				errOut("Encountered end of statement before all delimiter pairs were closed (iTTV).", t)
			}
			return false
		case INLNSEP:
			if curl > 0 || brak > 0 || parn > 0 {
				continue
			}
			return false
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
			default:
				return false
			}

			if curl < 0 || brak < 0 || parn < 0 {
				if curl > 0 || brak > 0 || parn > 0 {
					errOut("Un-matched closing delimiter (iTTV).", t)
				} else if curl + brak + parn == -1 {
					return false
				} else {
					errOut("Strange bracket values detected when parsing value.  Possibly a parser bug. (iTTV)", t)
				}
			}
		case AUGMENT:
			switch t.Data {
			case ".":
				if (*tokens)[tok + 1].Type == DEFWORD {
					tok++
				} else {
					errOut("Expected defword after 'get' operator (iTTV).", t)
				}
			case "~":
			case "`":
			default:
				return false
			}
		}
	}

	return stage == 2
}