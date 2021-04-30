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

func getClosing(start string) string {
	switch start {
	case "{":
		return "}"
	case "[":
		return "]"
	case "(":
		return ")"
	}

	return ""
}

// Parse a list of values
func parseValueList(tokens *[]Token, tok, max int) (Node, int) {
	out := Node{}
	out.Data = Token{Type: 10, Data: "value"}
	var tmp Node

	c := getClosing((*tokens)[tok].Data)

	tok++

	for ; tok < max; tok++ {
		t := (*tokens)[tok]

		switch t.Data {
		case c:
			return out, tok
		case ",":
			tok++
		default:
			errOut("Error: unexpected token when parsing a list of values", t)
		}

		tmp, tok = parseValue(tokens, tok, max)
		out.Sub = append(out.Sub, tmp)
	}

	return out, tok
}

func parseTypeList(tokens *[]Token, tok, max int) (Node, int) {
	out := Node{}
	out.Data = Token{Type: 10, Data: "type"}
	var tmp Node

	c := getClosing((*tokens)[tok].Data)

	tok++

	for ; tok < max; tok++ {
		t := (*tokens)[tok]

		switch t.Data {
		case c:
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

func parseStatementList(tokens *[]Token, tok, max int) (Node, int) {
	out := Node{}
	out.Data = Token{Type: 10, Data: "statement"}
	var tmp Node

	c := getClosing((*tokens)[tok].Data)

	tok++

	for ; tok < max; tok++ {
		t := (*tokens)[tok]

		switch t.Data {
		case c:
			return out, tok
		case ",":
			tok++
		default:
			errOut("Error: unexpected token when parsing a list of statements", t)
		}

		tmp, tok = parseStatement(tokens, tok, max)
		out.Sub = append(out.Sub, tmp)
	}

	return out, tok
}
