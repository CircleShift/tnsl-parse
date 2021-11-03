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
func isMain(artifact tparse.Node) bool {
	return false
}

// Check if a block is control flow
func isCF(artifact tparse.Node) bool {

}

// Get the control flow's name
func cfType(artifact tparse.Node) string {

}

// Get type as string from nodes
func evalType(artifact tparse.Node) string {
	return ""
}
 
// Returns generated value and general "type" of value (string, number, character)
func evalLiteral(artifact tparse.Node) (interface{}, string) {

}

// Evaluates a definition and sets up a TVariable in the context's var map
func evalDef(artifact tparse.Node, ctx *TContext) {
	vars := len(ctx.VarMap) - 1

	t := evalType(artifact.Sub[0])

	for i := 0; i < len(artifact.Sub[1].Sub); i++ {
		if artifact.Sub[1].Sub[i].Data.Data == "=" {
			artifact.Sub[1].Sub[i].Sub[0]
		}
	}
}

// Evaluates a value statement
func evalValue(artifact tparse.Node, ctx *TContext) {
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