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

package texec

import "strings"
import "tparse"

// Check if a block is the main function
func isMain(n tparse.Node) bool {
	if n.Data.Data == "block" {
		if n.Sub[0].Data.Data == "bdef" {
			for i := 0; i < len(n.Sub.Sub); i++ {
				if n.Sub.Sub[i].Data.Type == tparse.DEFWORD && n.Sub.Sub[i].Data.Data == "main" {
					return true
				}
			}
		}
	}
	return false
}

// Check if a block is control flow
func isCF(n tparse.Node) bool {
	if n.Data.Data == "block" {
		if n.Sub[0].Data.Data == "bdef" {
			for i := 0; i < len(n.Sub.Sub); i++ {
				if n.Sub.Sub[i].Data.Type == tparse.KEYWORD {
					if n.Sub.Sub[i].Data.Data == "if" || n.Sub.Sub[i].Data.Data == "elif" || n.Sub.Sub[i].Data.Data == "else" || n.Sub.Sub[i].Data.Data == "match" || n.Sub.Sub[i].Data.Data == "case" || n.Sub.Sub[i].Data.Data == "loop" {
						return true
					}
				}
			}
		}
	}
	return false
}

// Get the control flow's name
func cfType(n tparse.Node) string {

}

// Get type as string from nodes
func evalType(n tparse.Node) string {
	return ""
}
 
// Returns generated value and general "type" of value (string, number, character)
func evalLiteral(n tparse.Node) (interface{}, string) {

}

// Evaluates a definition and sets up a TVariable in the context's var map
func evalDef(n tparse.Node, ctx *TContext) {
	vars := len(ctx.VarMap) - 1

	t := evalType(n.Sub[0])

	for i := 0; i < len(n.Sub[1].Sub); i++ {
		if n.Sub[1].Sub[i].Data.Data == "=" {
			n.Sub[1].Sub[i].Sub[0]
		}
	}
}

// Evaluates a value statement
func evalValue(artifact tparse.Node, ctx *TContext) TVariable {
	vars := len(ctx.VarMap) - 1
}

// Evaluates control flow
func evalCF(artifact tparse.Node, ctx *TContext) {}

// Evaluate a block (Assume that all blocks have only one output for now)
func evalBlock(artifact tparse.Node, ctx *TContext) TVariable {
	
}


// EvalTNSL starts the evaluation on the root TModule's main function with the given flags passed to the program
func EvalTNSL(world *TModule, f string) {
	flags := strings.Split(f, " ")
}