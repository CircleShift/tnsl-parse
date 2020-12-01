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
			if t.Data != "~" {
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

		makeParent(working, Node{})
		working = &(working.Sub[0])
	}

	return out, tok
}
