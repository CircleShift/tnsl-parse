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

func parseBlock(tokens *[]Token, tok, max int) (Node, int) {
	out := Node{}
	out.Data = Token{Type: 10, Data: "block"}
	var tmp Node

	tok++

	for ;tok < max; tok++{
		t := (*tokens)[tok]

		switch t.Type {
		case DELIMIT:
			if t.Data == "(" {

			} else if t.Data == "(" {

			} else {
				goto BREAK
			}
		case DEFWORD:
		case KEYWORD:
		case LINESEP:
			goto BREAK
		}
	}

	BREAK:

	for ;tok < max; {
		t := (*tokens)[tok]

		switch t.Data {
		case ";/", ";;", ";:":
			return out, tok
		case ";":
			tmp, tok = parseStatement(tokens, tok, max)
		case "/;":
			REBLOCK:
			
			tmp, tok = parseBlock(tokens, tok, max)

			if (*tokens)[tok].Data == ";;" {
				out.Sub = append(out.Sub, tmp)
				goto REBLOCK
			} else if (*tokens)[tok].Data == ";/" {
				tok++
			}
		default:
			errOut("Error: unexpected token when parsing a code block", t)
		}

		out.Sub = append(out.Sub, tmp)
	}

	return out, tok
}

func parseStatement(tokens *[]Token, tok, max int) (Node, int) {
	out := Node{}
	out.Data = Token{Type: 11, Data: ";"}
	var tmp Node

	tok++

	for ; tok < max; tok++ {
		t := (*tokens)[tok]

		switch t.Type {
		case LINESEP, DELIMIT:
			return out, tok
		case INLNSEP:
			tok++
		default:
			errOut("Error: unexpected token when parsing a list of types", t)
		}

		tmp, tok = parseType(tokens, tok, max, true)
		out.Sub = append(out.Sub, tmp)
	}

	return out, tok
}

func parseDef(tokens *[]Token, tok, max int) (Node, int) {
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

		tmp, tok = parseType(tokens, tok, max, true)
		out.Sub = append(out.Sub, tmp)
	}

	return out, tok
}
