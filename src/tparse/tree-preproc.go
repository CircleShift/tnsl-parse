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

func parsePreBlock (tokens *[]Token, tok, max int) (Node, int) {
	out := Node{}
	out.Data = Token{Type: 11, Data: (*tokens)[tok].Data}

	tok++

	for ; tok < max; tok++ {
		t := (*tokens)[tok]

		if t.Data == ":/" || t.Data == ":;" {
			break
		}

		tmp := Node{Data: t}
		out.Sub = append(out.Sub, tmp)
	}

	return out, tok
}

func parsePre (tokens *[]Token, tok, max int) (Node, int) {
	out := Node{}
	out.Data = Token{Type: 11, Data: (*tokens)[tok].Data}

	tok++

	tmp := Node{Data: (*tokens)[tok]}
	out.Sub = append(out.Sub, tmp)

	tok++

	return out, tok
}