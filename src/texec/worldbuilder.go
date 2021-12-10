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

import (
	"tparse"
)

/**
	worldbuilder.go - take in a file name and construct a root TModule based on it.
*/

// Note: This is good enough, I guess. Gonna mark this as the final version, only update on major errors.

// Supported features:
// Importing other files
// Sub-modules across files

// Semi-borked sub-folders:
// Because the builder doesn't preserve the paths you are taking, it will not figure out which folder each file is in properly.
// Technically, you could work around this by making all imports in EVERY FILE EVERYWHERE look as if they are pathed from the folder
// where the root file is, but this would be a headache.  I am just planning on fixing this in the full compiler.

// Returns generated value and general "type" of value (string, number)
func evalPreLiteral(n tparse.Node) string {
	r := []rune(n.Data.Data)
	l := len(r)
	if r[0] == '"' || r[0] == '\'' {
		return string(r[1:l - 1])
	}
	return ""
}

//Generate a variable for a module
// For sub = 0, give the vlist
func modDef(n tparse.Node, t TType, sub int) ([]string, []TVariable) {

}

// Parse a file and make an AST from it.
func parseFile(p string) tparse.Node {
	tokens := tparse.TokenizeFile(p)
	return tparse.MakeTree(&(tokens), p)
}

// Import a file and auto-import sub-modules and files
func importFile(f string, m *TModule) {
	froot := parseFile(f)
	for n := 0 ; n < len(froot.Sub) ; n++ {
		if froot.Sub[n].Data.Data == "block" {
			if froot.Sub[n].Sub[0].Sub[0].Data.Data == "module" {
				m.Sub = append(m.Sub, buildModule(froot.Sub[n]))
			} else {
				m.Artifacts = append(m.Artifacts, froot.Sub[n])
			}
		} else if froot.Sub[n].Data.Data == "include" {
			importFile(evalPreLiteral(froot.Sub[n].Sub[0]), m)
		} else {
			m.Artifacts = append(m.Artifacts, froot.Sub[n])
		}
	}
}

// Build a module from a module block node
func buildModule(module tparse.Node) TModule {
	out := TModule{}
	out.Name = module.Sub[0].Sub[0].Sub[0].Data.Data

	for n := 1 ; n < len(module.Sub) ; n++ {
		if module.Sub[n].Data.Data == "include" {
			importFile(evalPreLiteral(module.Sub[n].Sub[0]), &out)
		}
	}

	return out
}

// BuildRoot builds the root module, ready for eval
func BuildRoot(file string) TModule {
	out := TModule{}

	importFile(file, &out)

	return out
}
