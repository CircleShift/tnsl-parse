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

var (
	// Program to run
	prog *TModule
	// Current artifact
	cart TArtifact

	//Default null value
	null = TVariable{tNull, nil}
)

//################
//# Helper Funcs #
//################

// Error helper

func errOut(msg string) {
	fmt.Println("==== BEGIN ERROR ====")
	fmt.Println(msg)
	fmt.Println(cart)
	fmt.Println("==== END ERROR ====")
	panic(">>> PANIC FROM EVAL <<<")
}

func errOutCTX(msg string, ctx VarMap) {
	fmt.Println("==== BEGIN ERROR ====")
	fmt.Println(msg)
	fmt.Println(cart)
	fmt.Println(ctx)
	fmt.Println("==== END ERROR ====")
	panic(">>> PANIC FROM EVAL <<<")
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
		} else if block.Sub[0].Sub[i].Data.Type == tparse.KEYWORD {
			switch block.Sub[0].Sub[i].Data.Data {
			case "if", "elif", "else", "loop", "match", "case", "default":
				out = append(out, block.Sub[0].Sub[i].Data.Data)
			default:
			}
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

func getModule(a TArtifact) *TModule {
	mod := prog
	
	for i := 0; i < len(a.Path); i++ {
		for j := 0; j < len(mod.Sub); j++ {
			if mod.Sub[j].Name == a.Path[i] {
				mod = &(mod.Sub[j])
				break
			}
			if j + 1 == len(mod.Sub) {
				errOut(fmt.Sprintf("Failed to find module %v", a))
			}
		}
	}

	return mod
}

func getModuleRelative(mod *TModule, a TArtifact) *TModule {
	for i := 0; i < len(a.Path); i++ {
		for j := 0; j < len(mod.Sub); j++ {
			if mod.Sub[j].Name == a.Path[i] {
				mod = &(mod.Sub[j])
				break
			}
			if j + 1 == len(mod.Sub) {
				return nil
			}
		}
	}

	return mod
}

func getModuleInPath(a TArtifact) *TModule {
	mod := prog
	m := len(cart.Path)

	out := getModuleRelative(mod, a)

	for i := 0; i < m; i++ {
		for j := 0; j < len(mod.Sub); j++ {
			if mod.Sub[j].Name == cart.Path[i] {
				mod = &(mod.Sub[j])
				break
			}
		}
		tmp := getModuleRelative(mod, a)
		if tmp != nil {
			out = tmp
		}
	}

	if out == nil {
		errOut(fmt.Sprintf("Failed to find module %d in path %v", m, cart))
	}

	return out
}

// Find an artifact from a path and the root node
func getNode(a TArtifact) *tparse.Node {
	mod := getModule(a)

	for i := 0; i < len(mod.Artifacts); i++ {
		n := getNames(mod.Artifacts[i])
		for j := 0; j < len(n); j++ {
			if n[j] == a.Name {
				return &(mod.Artifacts[i])
			}
		}
	}

	errOut(fmt.Sprintf("Failed to find node %v", a))
	return nil
}

func getNodeRelative(s TArtifact) *tparse.Node {

	tmpmod := getModuleInPath(s)
	if tmpmod == nil {
		errOut(fmt.Sprintf("Failed to get module to resolve node: %v", s))
	}

	for i := 0; i < len(tmpmod.Artifacts); i++ {
		n := getNames(tmpmod.Artifacts[i])
		for j := 0; j < len(n); j++ {
			if n[j] == s.Name {
				return &(tmpmod.Artifacts[i])
			}
		}
	}

	errOut(fmt.Sprintf("Failed to find node (relative) %v", s))
	return nil
}

func getModDefRelative(s TArtifact) *TVariable {

	tmpmod := getModuleInPath(s)
	if tmpmod == nil {
		errOut(fmt.Sprintf("Failed to get module to resolve artifact: %v", s))
	}

	val, prs := tmpmod.Defs[s.Name]

	if prs {
		return val
	}

	errOut(fmt.Sprintf("Failed to resolve mod def artifact (relative) %v", s))
	return nil
}

// Returns a mod definition, requires a resolved artifact
func getModDef(a TArtifact) *TVariable {
	mod := prog
	
	for i := 0; i < len(a.Path); i++ {
		for j := 0; j < len(mod.Sub); j++ {
			if mod.Sub[j].Name == a.Path[i] {
				mod = &(mod.Sub[j])
				break
			}
		}
	}

	v, prs := mod.Defs[a.Name]

	if prs {
		return v
	}

	errOut(fmt.Sprintf("Failed to resolve mod def artifact %v", a))
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
		return false
	}

	for i := preskip; i < len(a.Pre); i++ {
		if a.Pre[i] == "const" {
			preskip++
			continue
		} else if a.Pre[i] != b.Pre[i - preskip] {
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

	if (a.Post == "`" && b.Post != "`") || (b.Post == "`" && a.Post != "`") {
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
		errOut(fmt.Sprintf("Failed to parse string literal %v", v.Data))
	}

	return []byte(str)
}

func getCharLiteral(v tparse.Node) byte {
	val, mb, _, err := strconv.UnquoteChar(v.Data.Data, byte('\''))

	if err != nil || mb == true{
		errOut(fmt.Sprintf("Failed to parse character as single byte. %v", v.Data))
	}

	return byte(val)
}

func getIntLiteral(v tparse.Node) int {
	i, err := strconv.ParseInt(v.Data.Data, 0, 64)

	if err != nil {
		errOut(fmt.Sprintf("Failed to parse integer literal. %v", v.Data))
	}

	return int(i)
}

func getLiteralComposite(v tparse.Node) []interface{} {
	out := []interface{}{}

	for i := 0; i < len(v.Sub); i++ {
		if v.Sub[i].Data.Data[0] == '"' {
			out = append(out, getStringLiteral(v.Sub[i]))
		} else if v.Sub[i].Data.Data[0] == '\'' {
			out = append(out, getStringLiteral(v.Sub[i]))
		} else if v.Sub[i].Data.Data == "comp" {
			out = append(out, getLiteralComposite(v.Sub[i]))
		} else {
			out = append(out, getIntLiteral(v.Sub[i]))
		}
	}

	return out
}

func getLiteral(v tparse.Node, t TType) interface{} {

	if equateType(t, tInt) {
		return getIntLiteral(v)
	} else if equateType(t, tCharp) {
		return getCharLiteral(v)
	} else if equateType(t, tString) {
		return getStringLiteral(v)
	}

	return getLiteralComposite(v)
}

func compositeToStruct(str TArtifact, cmp []interface{}) VarMap {
	sv := getModDefRelative(str)

	vars := sv.Data.([]TVariable)
	if len(vars) != len(cmp) {
		return nil
	}
	
	out := make(VarMap)

	for i:=0;i<len(vars);i++ {
		if equateType(vars[i].Type, tInt) || equateType(vars[i].Type, tCharp) || equateType(vars[i].Type, tString) {
			out[vars[i].Data.(string)] = &(TVariable{vars[i].Type, cmp[i]})
		}
		
	}

	return out
}

//#####################
//# Finding Artifacts #
//#####################

func resolveArtifactCall(a TArtifact, params []TVariable) TVariable {
	tres := tnslResolve(a)
	if tres == 0 {
		if len(params) > 0 {
			return tnslEval(params[0], a.Name)
		} else {
			errOut("Need at least one arg to call tnsl.io func")
		}
	} else if tres == 1 {
		if len(params) > 1 {
			return tnslFileEval(params[0], params[1], a.Name)
		} else {
			errOut("Not enough args recieved to call tnsl.io.File method.")
		}
	}

	return null
}

func resolveArtifact(a TArtifact, ctx *VarMap) *TVariable {
	if len(a.Path) == 0 {
		val, prs := (*ctx)[a.Name]
		if !prs {
			errOutCTX(fmt.Sprintf("Could not resolve %s in the current context.", a.Name), *ctx)
		}
		return val
	} else {

	}

	return nil
}

//#################
//# Runtime funcs #
//#################

// Value statement parsing

// Parse a value node
func evalValue(v tparse.Node, ctx *VarMap) TVariable {
	switch v.Data.Data {
	case "=":
	case "+":
	case "-":
	case "*":
	case "/":
	case "&&":
	case "||":
	case "==":
	case ".":
	}
	return null
}

// Generate a value for a definition
func evalDefVal(v tparse.Node, ctx *VarMap) {
	
}

// Eval a definition
func evalDef(v tparse.Node, ctx *VarMap) {
	
}

// Eval a control flow
func evalCF(v tparse.Node, ctx *VarMap) (bool, TVariable) {
	//scopeVars := []string{}
	return false, null
}

func evalBlock(b tparse.Node, params []TVariable) TVariable {
	ctx := make(VarMap)

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

	return null
}

func EvalTNSL(root *TModule, args string) TVariable {
	prog = root
	cart = TArtifact { []string{}, "main" }

	sarg := strings.Split(args, " ")
	
	targ := TVariable {
		TType {
			[]string{"{}", "{}"},
			TArtifact { []string{}, "charp" },
			"" },
		sarg }

	mainNod := getNode(cart)
	
	fmt.Println(mainNod)

	return evalBlock(*mainNod, []TVariable{targ})
}
