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

// Token represents a token in a program
type Token struct {
	Type int
	Data string
	Line int
	Char int
}

// Node represents a node in an AST
type Node struct {
	Data Token
	Sub  []Node
}

func makeParent(parent *Node, child Node) {
	parent.Sub = append(parent.Sub, child)
}
