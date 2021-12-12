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


func modDef(n tparse.Node, m *TModule) {
	t := getType(n.Sub[0])
	s, vs := modDefVars(n.Sub[1], t)
	for i := 0; i < len(s); i++ {
		m.Defs[s[i]] = vs[i]
	}
}

// Generate a variable list for a module
// For sub = 0, give the vlist
// May be horribly broken.  Definitely doesn't support composite types.
func modDefVars(n tparse.Node, t TType) ([]string, []TVariable) {
	s := []string{}
	v := []TVariable{}
	for i := 0; i < len(n.Sub); i++ {
		if n.Sub[i].Data.Type == tparse.DEFWORD {
			s = append(s, n.Sub[i].Data.Data)
			v = append(v, TVariable{t, nil})
		} else if n.Sub[i].Data.Data == "=" && n.Sub[i].Sub[0].Data.Type == tparse.DEFWORD {
			s = append(s, n.Sub[i].Sub[0].Data.Data)
			v = append(v, TVariable{t, getLiteral(n.Sub[i].Sub[1], t)})
		} else {
			errOut("Unexpected thing in definition. Expected '=' or DEFWORD.", n.Sub[i].Data)
		}
	}
	return s, v
}

func modDefEnum(n tparse.Node, m *TModule) {
	name := n.Sub[0].Data.Data
	t := getType(n.Sub[1].Sub[0])
	s, vs := modDefVars(n.Sub[2], t)
	out := TVariable{tEnum, make(VarMap)}
	for i := 0; i < len(s); i++ {
		out.Data.(VarMap)[s[i]] = vs[i]
	}
	m.Defs[name] = out
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
		} else if froot.Sub[n].Data.Data == "define" {
			modDef(froot.Sub[n], m)
		} else if froot.Sub[n].Data.Data == "enum"{
			modDefEnum(froot.Sub[n], m)
		}else if froot.Sub[n].Data.Data == "struct" || froot.Sub[n].Data.Data == "raw"{
		} else {
			m.Artifacts = append(m.Artifacts, froot.Sub[n])
		}
	}
}

// Build a module from a module block node
func buildModule(module tparse.Node) TModule {
	out := TModule{}
	out.Defs = make(VarMap)
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
	out.Defs = make(VarMap)

	importFile(file, &out)

	return out
}
