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

// TODO: re-validate this code.  I forgot if it works or not.
func parseBlock(tokens *[]Token, tok, max int) (Node, int) {
	out, tmp, def, name := Node{}, Node{}, Node{}, false
	out.Data = Token{Type: 10, Data: "block"}
	out.IsBlock = true
	def.Data = Token{Type: 10, Data: "bdef"}

	for ;tok < max; tok++{
		t := (*tokens)[tok]
		tmp = Node{}
		switch t.Type {
		case DELIMIT:
			if t.Data == "(" {
				tmp, tok = parseParamList(tokens, tok + 1, max)
				def.Sub = append(def.Sub, tmp)
			} else if t.Data == "[" {
				tmp, tok = parseTypeList(tokens, tok + 1, max)
				def.Sub = append(def.Sub, tmp)
			} else {
				goto BREAK
			}
		case DEFWORD:
			if name {
				errOut("Unexpected defword, the block has already been given a name or keyword.", t)
			}
			tmp.Data = t
			def.Sub = append(def.Sub, tmp)
			name = true
		case KEYWORD:
			if name {
				errOut("Unexpected keyword, the block has already been given a name or keyword.", t)
			}
			switch t.Data {
			case "else":
				if (*tokens)[tok+1].Data == "if" {
					tmp.Data = Token{KEYWORD, "elif", t.Line, t.Char}
					def.Sub = append(def.Sub, tmp)
					tok++
					continue
				}
			case "if", "match", "case", "loop":
				name = true
				fallthrough
			case "export", "inline", "raw", "override":
				tmp.Data = t
				def.Sub = append(def.Sub, tmp)
			default:
				errOut("Unexpected keyword in block definition.", t)
			}
		case LINESEP:
			goto BREAK
		}
	}

	BREAK:

	out.Sub = append(out.Sub, def)

	for ;tok < max; {
		t := (*tokens)[tok]

		switch t.Data {
		case ";/", ";;", ";:":
			return out, tok
		case ";":

			tmp, tok = parseStatement(tokens, tok + 1, max)
		case "/;":
			REBLOCK:
			
			tmp, tok = parseBlock(tokens, tok + 1, max)

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

// This should work once isTypeThenValue properly functions
func parseStatement(tokens *[]Token, tok, max int) (Node, int) {
	out := Node{}
	out.Data = Token{Type: 11, Data: ";"}
	var tmp Node

	// Check for keyword, definition, then if none of those apply, assume it's a value.
	if (*tokens)[tok].Type == KEYWORD {
		tmp, tok = keywordStatement(tokens, tok, max)
		out.Sub = append(out.Sub, tmp)
	} else {
		// do check for definition
		if isTypeThenValue(tokens, tok, max) {
			// if not, parse a value
			tmp, tok = parseDef(tokens, tok, max)
		} else {
			// if not, parse a value
			tmp, tok = parseValue(tokens, tok, max)
		}
		out.Sub = append(out.Sub, tmp)
	}

	return out, tok
}

// Works?  Please test.
func keywordStatement(tokens *[]Token, tok, max int) (Node, int) {
	out := Node{}
	out.Data = (*tokens)[tok]
	out.IsBlock = false
	var tmp Node
	tmp.IsBlock = false

	if tok + 1 < max {
		tok++
	} else {
		return out, max
	}

	switch out.Data.Data {
	case "raw":
		// Something, something... code.
		if (*tokens)[tok].Data != "struct" {
			errOut("Unexpected use of raw operator in a statement", out.Data)
		}
		tok++
		fallthrough
	case "struct":
		// Check for defword, (), and {} then dip
		if (*tokens)[tok].Type != DEFWORD {
			errOut("Expected defword after struct keyword.", (*tokens)[tok])
		}
		tmp.Data = (*tokens)[tok]
		out.Sub = append(out.Sub, tmp)
		tok++
		if (*tokens)[tok].Data == "(" {
			tmp, tok = parseValueList(tokens, tok, max)
			out.Sub = append(out.Sub, tmp)
		}

		if (*tokens)[tok].Data != "{" {
			errOut("Could not find struct member list", (*tokens)[tok])
		}

		tmp, tok = parseParamList(tokens, tok, max)

	case "goto", "label":
		if (*tokens)[tok].Type != DEFWORD {
			errOut("Expected defword after goto or label keyword.", out.Data)
		}
		tmp.Data = (*tokens)[tok]
		tok++
		// Check for a defword and dip
	case "continue", "break":
		if (*tokens)[tok].Type != LITERAL {
			return out, tok
		}
		tmp.Data = (*tokens)[tok]
		tok++
		// Check for a numerical value and dip
	case "return":
		if (*tokens)[tok].Type != DELIMIT || (*tokens)[tok].Data == "{" || (*tokens)[tok].Data == "(" {
			tmp, tok = parseValue(tokens, tok, max)
		}
	case "alloc", "salloc":
		// Parse value list
		tmp, tok = parseValueList(tokens, tok, max)
	case "delete":
		// Parse value list
		tmp, tok = parseValueList(tokens, tok, max)
	}

	out.Sub = append(out.Sub, tmp)

	return out, tok
}

// Should work, but none of this is tested.
func parseDef(tokens *[]Token, tok, max int) (Node, int) {
	out := Node{Data: Token{11, "vdef", 0, 0}}
	var tmp Node

	tmp, tok = parseType(tokens, tok, max, false)
	out.Sub = append(out.Sub, tmp)
	tmp, tok = parseValueList(tokens, tok, max)
	out.Sub = append(out.Sub, tmp)

	return out, tok
}
