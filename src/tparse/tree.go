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

import "fmt"

// ID 9 = ast thing

func errOut(message string, token Token) {
	fmt.Println(message)
	fmt.Println(token)
	panic("AST Error")
}

// Parse a list of values
func parseValueList(tokens *[]Token, tok, max int) (Node, int) {
	out := Node{}
	out.Data = Token{Type: 10, Data: "list"}
	var tmp Node

	tok++

	for ; tok < max; tok++ {
		t := (*tokens)[tok]

		switch t.Data {
		case ")", "]", "}":
			return out, tok
		case ",":
			tok++
		default:
			errOut("Error: unexpected token when parsing a list of types", t)
		}

		tmp, tok = parseValue(tokens, tok, max)
		out.Sub = append(out.Sub, tmp)
	}

	return out, tok
}

// Parses a list of things
func parseDefList(tokens *[]Token, tok, max int) (Node, int) {
	out := Node{}
	out.Data = Token{Type: 9, Data: "list"}

	currentType := Node{}
	currentType.Data = Token{Data: "undefined"}

	for ; tok < max; tok++ {
		t0 := (*tokens)[tok]
		t1 := (*tokens)[tok+1]

		switch t1.Data {
		case ")", "]", "}", ",":
		default:
			currentType, tok = parseType(tokens, tok, max, true)
			t0 = (*tokens)[tok]
			t1 = (*tokens)[tok+1]
		}

		switch t0.Type {

		case DEFWORD:
			var tmp Node
			if currentType.Data.Data == "undefined" {
				errOut("Error: expected type before first parameter", t0)
			} else if currentType.Data.Data == "type" {
				tmp, tok = parseType(tokens, tok, max, true)
			} else {
				tmp = Node{Data: t0}
			}

			typ := currentType
			makeParent(&typ, tmp)
			makeParent(&out, typ)

		default:
			errOut("Error: unexpected token when parsing list, expected user-defined variable or ", t0)
		}

		switch t1.Data {
		case ")", "]", "}":
			return out, tok
		case ",":
		default:
			errOut("Error: unexpected token when parsing list, expected ',' or end of list", t1)
		}

		tok++
	}

	return out, tok
}

func parseTypeList(tokens *[]Token, tok, max int) (Node, int) {
	out := Node{}
	out.Data = Token{Type: 9, Data: "list"}
	var tmp Node

	tok++

	for ; tok < max; tok++ {
		t := (*tokens)[tok]

		switch t.Data {
		case ")", "]", "}":
			return out, tok
		case ",":
			tok++
		default:
			errOut("Error: unexpected token when parsing a list of types", t)
		}

		tmp, tok = parseType(tokens, tok, max, true)
		out.Sub = append(out.Sub, tmp)
	}

	return out, tok
}

func parseVoidType(tokens *[]Token, tok, max int) (Node, int) {
	out := Node{}

	for ; tok < max; tok++ {
		//t := (*tokens)[tok]
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
				errOut("Error: parameter types cannot be static", t)
			} else if t.Data != "const" && t.Data != "volatile" && t.Data != "static" {
				errOut("Error: unexpected keyword when parsing type", t)
			}
			working.Data = t

		case DELIMIT:
			if t.Data == "{" {
				if (*tokens)[tok+1].Data == "}" {
					working.Data = Token{9, "array", t.Line, t.Char}
					tok++
				} else {
					errOut("Error: start of list when parsing type (did you mean {} ?)", t)
				}
			} else {
				errOut("Error: start of list when parsing type", t)
			}

		default:
			errOut("Error: unexpected token when parsing type", t)
		}

		working.Sub = append(working.Sub, Node{Parent: working})
		working = &(working.Sub[0])
	}

	return out, tok
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

// MakeTree creates an AST out of a set of tokens
func MakeTree(tokens *[]Token, file string) Node {
	out := Node{}
	out.Data = Token{9, file, 0, 0}
	out.Parent = &out

	tmp := Node{}
	working := &tmp

	for _, t := range *tokens {
		switch t.Type {
		case LINESEP:

		case DELIMIT:

		}
		tmp = Node{Data: t}

		working.Sub = append(working.Sub, tmp)
	}

	return out
}
