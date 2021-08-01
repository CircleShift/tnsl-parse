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

// MakeTree creates an AST out of a set of tokens
func MakeTree(tokens *[]Token, file string) Node {
	out := Node{}
	out.Data = Token{9, file, 0, 0}

	tmp := Node{}

	max := len(*tokens)

	for tok := 0; tok < max; tok++ {
		t := (*tokens)[tok]
		switch t.Data {
		case "/;":
			tmp, tok = parseBlock(tokens, tok, max)
		case ";":
			tmp, tok = parseStatement(tokens, tok, max)
		case "/:":
			tmp = Node{}
		case ":":
			tmp = Node{}
		default:
			errOut("Unexpected token in file root", t)
		}

		out.Sub = append(out.Sub, tmp)
	}

	return out
}
