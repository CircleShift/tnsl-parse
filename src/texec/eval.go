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

import "tparse"
import "fmt"

/*
	So here's what I care to support at present:
	Type checking, basic types, writing to stdout or a file
	Variable and state contexts
	Reading from files
	Raw structs
	Appending to arrays
	Calling functions and methods
	libtnsl stub

	This subset should theoretically be enough to write a compiler.
*/

//################
//# Helper Funcs #
//################

// Error helper

func errOut(path TArtifact, place tparse.Token) {
	fmt.Println("Error in eval:")
	fmt.Println(path)
	fmt.Println(place)
	panic("EVAL ERROR")
}

// Names of artifacts, finding artifacts

func getDefNames(def tparse.Node) []string {
	out := []string{}
	for i := 0; i < len(def.Sub); i++ {
		if def.Sub[i].Data.Data == "vlist" && def.Sub[i].Data.Type == 10 {
			vl := def.Sub[i]
			for i := 0; i < len(vl.Sub); i++ {
				if vl.Sub[i].Data.Type == tparse.DEFWORD {
					out = append(out, vl.Sub[i].Data.Data)
				} else if vl.Sub[i].Data.Data == "=" && vl.Sub[i].Sub[0].Data.Type == tparse.DEFWORD {
					out = append(out, vl.Sub[i].Sub[0].Data.Data)
				}
			}
		}
	}
	return out
}

func getBlockName(block tparse.Node) []string {
	out := []string{}
	for i := 0; i < len(block.Sub[0].Sub); i++ {
		if block.Sub[0].Sub[i].Data.Type == tparse.DEFWORD {
			out = append(out, block.Sub[0].Sub[i].Data.Data)
		}
	}
	return out
}

func getTypeName(t tparse.Node) []string {
	out := []string{}
	for i := 0; i < len(t.Sub); i++ {
		if t.Sub[i].Data.Type == tparse.DEFWORD {
			out = append(out, t.Sub[i].Data.Data)
		}
	}
	return out
}

// Get the list of names defined by the block or variable definition
func getNames(root tparse.Node) []string {
	switch root.Data.Data {
	case "block":
		return getBlockName(root)
	case "define":
		return getDefNames(root)
	case "raw", "switch", "enum":
		return getTypeName(root)
	}

	return []string{}
}

// Find an artifact from a path and the root node
func getArtifact(a TArtifact, root *TModule) *tparse.Node {
	mod := root
	for i := 0; i < len(a.Path); i++ {
		for j := 0; j < len(mod.Sub); j++ {
			if mod.Sub[j].Name == a.Path[i] {
				mod = &(mod.Sub[j])
				break
			}
		}
	}

	for i := 0; i < len(mod.Artifacts); i++ {
		n := getNames(mod.Artifacts[i])
		for i := 0; i < len(n); i++ {
			if n[i] == a.Name {
				return &(mod.Artifacts[i])
			}
		}
	}

	return nil
}

// Type related stuff

// Checking type equality
func equateType(a, b TType) bool {
	if len(a.Pre) != len(b.Pre) || len(a.Post) != len(b.Post) {
		return false
	} else if len(a.T.Path) != len(b.T.Path) {
		return false
	}

	for i := 0; i < len(a.Pre); i++ {
		if a.Pre[i] != b.Pre[i] {
			return false
		}
	}

	for i := 0; i < len(a.T.Path); i++ {
		if a.T.Path[i] != b.T.Path[i] {
			return false
		}
	}

	if a.T.Name != b.T.Name {
		return false
	}

	for i := 0; i < len(a.Post); i++ {
		if a.Post[i] != b.Post[i] {
			return false
		}
	}

	return true;
}

// Generate a TType from a 'type' node
func getType(t tparse.Node) TType {
	out := TType{}

	return out
}

// Value generation

func getStringLiteral(v tparse.Node) []byte {

}

func getCharLiteral(v tparse.Node) byte {

}

func getIntLiteral(v tparse.Node) int {

}

// Get a literal value from nodes.  Must specify type of literal to generate.
func getLiteral(v tparse.Node, t TType) interface{} {

}

//#################
//# Runtime funcs #
//#################

