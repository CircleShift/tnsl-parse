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

// ID 9 = ast root
// ID 10 = ast token

func errOut(message string, token Token) {
	fmt.Println(message)
	fmt.Println(token)
	panic("AST Error")
}

func errOutV(message string, tok, max int, token Token) {
	fmt.Println(message)
	fmt.Println(token)
	fmt.Println(tok)
	fmt.Println(max)
	panic("AST Error")
}
// MakeTree creates an AST out of a set of tokens
func MakeTree(tokens *[]Token, file string) Node {
	out := Node{}
	out.Data = Token{9, file, 0, 0}

	tmp := Node{}

	max := len(*tokens)

	for tok := 0; tok < max; {
		t := (*tokens)[tok]
		switch t.Data {
		case "/;", ";;", ":;":
			REBLOCK:
			
			tmp, tok = parseBlock(tokens, tok + 1, max)

			if (*tokens)[tok].Data == ";;" {
				out.Sub = append(out.Sub, tmp)
				goto REBLOCK
			} else if (*tokens)[tok].Data == ";/" {
				tok++
			}
		case ";":
			tmp, tok = parseStatement(tokens, tok + 1, max)
		case "/:", ";:":
			tmp, tok = parsePreBlock(tokens, tok + 1, max)
		case ":":
			tmp, tok = parsePre(tokens, tok + 1, max)
		default:
			errOut("Unexpected token in file root", t)
		}

		out.Sub = append(out.Sub, tmp)
	}

	return out
}

func findClosing(tokens *[]Token, tok int) int {
	t := (*tokens)[tok]
	var match string
	
	switch (t.Data) {
	case "(":
		match = ")"
	case "[":
		match = "]"
	case "{":
		match = "}"
	default:
		errOut("[Internal] Attempted to find closing for a non-delim token.", t)
	}

	var curl, brak, parn int = 0, 0, 0

	for tok++; tok < len(*tokens); tok++ {
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
				if (curl < 0 && match == "}") || (brak < 0 && match == "]") || (parn < 0 && match == ")") {
					return tok
				} else {
					errOut("Un-matched closing delimiter when searching for a closing delim.", t)
				}
			}
		}
	}

	return -1
}
