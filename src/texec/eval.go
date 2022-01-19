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
	"fmt"
	"strconv"
	"strings"
)

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

func errOut(msg string, place tparse.Token) {
	fmt.Println("Error in eval:")
	fmt.Println(msg)
	fmt.Println(place)
	panic("EVAL ERROR")
}

func errOutCTX(msg string, place tparse.Token, ctx TContext) {
	fmt.Println("Error in eval:")
	fmt.Println(msg)
	fmt.Println(place)
	fmt.Println(ctx)
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
		} else if block.Sub[0].Sub[i].Data.Data == "method" {
			out = append(out, block.Sub[0].Sub[i].Sub[0].Data.Data)
		}
	}
	fmt.Println(out)
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
	case "raw", "struct", "enum":
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
		for j := 0; j < len(n); j++ {
			if n[j] == a.Name {
				return &(mod.Artifacts[i])
			}
		}
	}

	return nil
}

// Type related stuff

// Checking type equality
// Assumes a is an unknown type and b is a known good type.
func equateTypePS(a, b TType, preskip int) bool {
	cc := 0
	for i := 0; i < len(a.Pre); i++ {
		if a.Pre[i] == "const" {
			cc++
		}
	}

	if len(a.T.Path) != len(b.T.Path) || len(a.Pre) - preskip - cc != len(b.Pre) {
		fmt.Println("[EVAL] Equate type died at len check.")
		return false
	}

	for i := preskip; i < len(a.Pre); i++ {
		if a.Pre[i] == "const" {
			preskip++
			continue
		} else if a.Pre[i] != b.Pre[i - preskip] {
			fmt.Println("[EVAL] Equate type died at pre check.")
			return false
		}
	}

	for i := 0; i < len(a.T.Path); i++ {
		if a.T.Path[i] != b.T.Path[i] {
			fmt.Println("[EVAL] Equate type died at path check.")
			return false
		}
	}

	if a.T.Name != b.T.Name {
		fmt.Println("[EVAL] Equate type died at name check.")
		return false
	}

	if (a.Post == "`" && b.Post != "`") || (b.Post == "`" && a.Post != "`") {
		fmt.Println("[EVAL] Equate type died at rel check.")
		return false
	}

	return true;
}

func equateType(a, b TType) bool {
	return equateTypePS(a, b, 0)
}

// Generate a TType from a 'type' node
func getType(t tparse.Node) TType {
	out := TType{}
	i := 0

	// Pre
	for ; i < len(t.Sub); i++ {
		if t.Sub[i].Data.Type == tparse.DEFWORD || t.Sub[i].Data.Type == tparse.KEYTYPE {
			break
		} else {
			out.Pre = append(out.Pre, t.Sub[i].Data.Data)
		}
	}

	// T
	for ; i < len(t.Sub); i++ {
		if t.Sub[i].Data.Type == tparse.KEYTYPE {
			out.T.Name = t.Sub[i].Data.Data
			i++
			break
		} else if t.Sub[i].Data.Type == tparse.DEFWORD {
			if i < len(t.Sub) - 1 {
				if t.Sub[i + 1].Data.Type == tparse.DEFWORD {
					out.T.Path = append(out.T.Path, t.Sub[i].Data.Data)
				} else {
					out.T.Name = t.Sub[i].Data.Data
					break
				}
			} else {
				out.T.Name = t.Sub[i].Data.Data
			}
		}
	}

	// Post
	if i < len(t.Sub) {
		if t.Sub[i].Data.Data == "`" {
			out.Post = "`"
		}
	}

	return out
}

// Value generation

func getStringLiteral(v tparse.Node) []byte {
	str, err := strconv.Unquote(v.Data.Data)

	if err != nil {
		errOut("Failed to parse string literal.", v.Data)
	}

	return []byte(str)
}

func getCharLiteral(v tparse.Node) byte {
	val, mb, _, err := strconv.UnquoteChar(v.Data.Data, byte('\''))

	if err != nil || mb == true{
		errOut("Failed to parse character as single byte.", v.Data)
	}

	return byte(val)
}

func getIntLiteral(v tparse.Node) int {
	i, err := strconv.ParseInt(v.Data.Data, 0, 64)

	if err != nil {
		errOut("Failed to parse integer literal.", v.Data)
	}

	return int(i)
}

func getLiteral(v tparse.Node, t TType) interface{} {

	if equateType(t, tInt) {
		return getIntLiteral(v)
	} else if equateType(t, tCharp) {
		return getCharLiteral(v)
	} else if equateType(t, tString) {
		return getStringLiteral(v)
	}

	return nil
}



//#####################
//# Finding Artifacts #
//#####################

func resolveModArtifact(a TArtifact) *TVariable {
	return nil
}

func resolveArtifactCall(a TArtifact, params []TVariable) TVariable {
	return TVariable{tNull, nil}
}

func resolveArtifact(a TArtifact, ctx *TContext, root *TModule) *TVariable {
	return nil
}

//#################
//# Runtime funcs #
//#################

// Value statement parsing

// Get a value from nodes.  Must specify type of value to generate.
func evalValue(v tparse.Node, ctx *TContext) TVariable {
	return TVariable{tNull, nil}
}

// Get a value from nodes.  Must specify type of value to generate.
func evalDef(v tparse.Node, ctx *TContext) {
	
}

// Get a value from nodes.  Must specify type of value to generate.
func evalCF(v tparse.Node, ctx *TContext) (bool, TVariable) {
	//scopeVars := []string{}
	return false, TVariable{tNull, nil}
}

func evalBlock(b tparse.Node, m TArtifact, params []TVariable) TVariable {
	ctx := TContext { m, make(VarMap) }

	for i := 0; i < len(b.Sub); i++ {
		switch b.Sub[i].Data.Data {
		case "define":
			evalDef(b.Sub[i], &ctx)
		case "value":
			evalValue(b.Sub[i], &ctx)
		case "block":
			ret, val := evalCF(b.Sub[i], &ctx)
			if ret {
				return val
			}
		case "return":
			return evalValue(b.Sub[i].Sub[0], &ctx)
		}
	}

	return TVariable{tNull, nil}
}

func EvalTNSL(root *TModule, args string) TVariable {
	sarg := strings.Split(args, " ")
	
	targ := TVariable {
		TType {
			[]string{"{}", "{}"},
			TArtifact { []string{}, "charp" },
			"" },
		sarg }

	mainArt := TArtifact { []string{}, "main" }

	mainNod := getArtifact(mainArt, root)
	
	fmt.Println(mainNod)

	return evalBlock(*mainNod, mainArt, []TVariable{targ})
}