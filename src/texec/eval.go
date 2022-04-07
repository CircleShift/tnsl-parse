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

func errOutCTX(msg string, ctx *VarMap) {
	fmt.Println("==== BEGIN ERROR ====")
	fmt.Println(msg)
	fmt.Println(cart)
	fmt.Println(*ctx)
	fmt.Println("====  END  ERROR ====")
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

// Attempt to get a module from a path starting at the given module
// Returns nil if the module was not found.
func getModuleRelative(mod *TModule, a TArtifact) *TModule {
	if mod == nil {
		return nil
	}

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

// Attempt to get a module from the root module using a specified path
// Returns nil if the module was not found.
func getModule(a TArtifact) *TModule {
	return getModuleRelative(prog, a)
}

// Get a module ion the current path.
// Returns nil if the index is out of range.
func getModuleInPath(p int) *TModule {
	m := len(cart.Path)

	if p < 0 || p > m {
		return nil
	}

	return getModule( TArtifact{ cart.Path[:p] , "" } )
}

// Find an artifact from a name and the module to search
// Returns nil if the node is not found in the module
func getNode(mod *TModule, n string) *tparse.Node {

	for i := 0; i < len(mod.Artifacts); i++ {
		chk := getNames(mod.Artifacts[i])
		for j := 0; j < len(chk); j++ {
			if chk[j] == n {
				return &(mod.Artifacts[i])
			}
		}
	}

	return nil
}

func getDef(mod *TModule, n string) *TVariable {
	ret, prs := (*mod).Defs[n]
	if prs {
		return ret
	}
	return nil
}

// This is a horrible way to search through nodes with this structure.  O(n^3).
// This could (and should) be made better by using a dictionary like structure for sub-modules and artifacts.
// By sacrificing this type of tree it could probably get down to O(n^2) or even O(n) if you were good enough.
// Most probably, the following is not how it will be implemented in the final version of the compiler.
// If this wasn't a bootstrap/hoby project, I would probably fire myself for the following code.

// Yes, I am aware that the following code is bad.
// No, I don't care.

func searchNode(s TArtifact) *tparse.Node {

	// i-- because we are doing a reverse lookup
	for i := len(cart.Path); i >= 0; i-- { // O(n)
		tst := getModuleInPath(i) // O(n^2) (O(n^3) total here)
		tst = getModuleRelative(tst, s) // O(n^2) (O(n^3) total here)

		if tst == nil {
			continue
		}

		ret := getNode(tst, s.Name) // O(n^2) (O(n^3) total here)

		if ret != nil {
			return ret
		}
	} // Block total complexity 3*O(n^2) * O(n) = 3*O(n^3)

	return nil
}

func searchDef(s TArtifact) *TVariable {

	// i-- because of reverse lookup
	for i := len(cart.Path); i >= 0; i-- {
		tst := getModuleInPath(i)
		tst = getModuleRelative(tst, s)

		if tst == nil {
			continue
		}

		ret := getDef(tst, s.Name)

		if ret != nil {
			return ret
		}
	}
	return nil
}

// End block of complexity horror

// Type related stuff

// Checking type equality
// Assumes a is an unknown and b is also an unknown.
func equateTypePS(a, b TType, psa, psb int) bool {
	if len(a.T.Path) != len(b.T.Path) || len(a.Pre) - psa != len(b.Pre) - psb {
		return false
	}

	for i := 0; i < len(a.Pre) - psa; i++ {
		if a.Pre[psa + i] != b.Pre[psb + i] {
			return false
		}
	}

	for i := 0; i < len(a.T.Path);i++ {
		if a.T.Path[i] != b.T.Path[i] {
			return false
		}
	}

	if a.T.Name != b.T.Name || a.Post != b.Post {
		return false
	}

	return true
}

func equateTypePSB(a, b TType, ps int) bool {
	return equateTypePS(a, b, ps, ps)
}

// Checking type equality
// Assumes a is an unknown type and b is a known good type.
func equateTypePSO(a, b TType, ps int) bool {
	return equateTypePS(a, b, ps, 0)
}

func equateType(a, b TType) bool {
	return equateTypePS(a, b, 0, 0)
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

func stripType(t TType, s int) TType {
	return TType{t.Pre[s:], t.T, t.Post}
}

// Value generation

func getStringLiteral(v tparse.Node) []interface{} {
	str, err := strconv.Unquote(v.Data.Data)

	if err != nil {
		errOut(fmt.Sprintf("Failed to parse string literal %v", v.Data))
	}

	dat := []byte(str)
	out := []interface{}{}

	for i := 0; i < len(dat); i++ {
		out = append(out, dat)
	}

	return out
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

func getFloatLiteral(v tparse.Node) float64 {
	i, err := strconv.ParseFloat(v.Data.Data, 64)

	if err != nil {
		errOut(fmt.Sprintf("Failed to parse float literal. %v", v.Data))
	}

	return float64(i)
}

func getLiteralComposite(v tparse.Node) []interface{} {
	out := []interface{}{}

	for i := 0; i < len(v.Sub); i++ {
		if v.Sub[i].Data.Data[0] == '"' {
			out = append(out, getStringLiteral(v.Sub[i]))
		} else if v.Sub[i].Data.Data[0] == '\'' {
			out = append(out, getCharLiteral(v.Sub[i]))
		} else if v.Sub[i].Data.Data == "comp" {
			out = append(out, getLiteralComposite(v.Sub[i]))
		} else if v.Sub[i].Data.Data[0] == '0' {
			out = append(out, getIntLiteral(v.Sub[i]))
		} else {
			out = append(out, getFloatLiteral(v.Sub[i]))
		}
	}

	return out
}

func getBoolLiteral(v tparse.Node) bool {
	return v.Data.Data == "true"
}

func getLiteral(v tparse.Node, t TType) interface{} {
	if equateType(t, tFloat) {
		return getFloatLiteral(v)
	} else if equateType(t, tCharp) {
		return getCharLiteral(v)
	} else if equateType(t, tString) {
		return getStringLiteral(v)
	} else if equateType(t, tBool) {
		return getBoolLiteral(v)
	} else if equateType(t, tInt) {
		return getIntLiteral(v)
	}

	return getLiteralComposite(v)
}

func getLiteralType(v tparse.Node) TType {
	if v.Data.Data[0] == '"' {
		return tString
	} else if v.Data.Data[0] == '\'' {
		return tCharp
	} else if v.Data.Data == "comp" {
		return tStruct
	} else if v.Data.Data == "true" || v.Data.Data == "false" {
		return tBool
	} else if v.Data.Data[0] == '0' {
		return tInt
	} else {
		return tFloat
	}
}

// Convert Value to Struct from Array (cvsa)
// USE ONLY IN THE CASE OF tStruct!
func cvsa(sct TArtifact, dat []interface{}) VarMap {
	sv := searchDef(sct)
	
	old_c := cart
	cart = sct
	
	vars := sv.Data.([]TVariable)
	if len(vars) != len(dat) {
		return nil
	}

	out := make(VarMap)

	for i:=0;i<len(vars);i++ {
		tmp := TVariable{vars[i].Type, nil}
		if isStruct(vars[i].Type, 0) {
			tmp.Data = cvsa(vars[i].Type.T, dat[i].([]interface{}))
		} else {
			tmp.Data = dat[i]
		}
		out[vars[i].Data.(string)] = &tmp
	}

	cart = old_c

	return out
}

// Copy aray to aray (cata)
func cata(st TArtifact, dat []interface{}) []interface{} {
	out := []interface{}{}

	for i := 0; i < len(dat); i++ {
		switch v := dat[i].(type) {
		case []interface{}:
			out = append(out, cata(st, v))
		case VarMap:
			out = append(out, csts(st, v))
		default:
			out = append(out, v)
		}
	}

	return out
}

// Copy struct to struct
// Makes a deep copy of a struct.
func csts(st TArtifact, dat VarMap) VarMap {
	sv := searchDef(st)
	
	old_c := cart
	cart = st
	
	vars := sv.Data.([]TVariable)

	out := make(VarMap)

	for i := 0; i < len(vars); i++ {
		var dts interface{} = nil

		switch v := dat[vars[i].Data.(string)].Data.(type) {
		case []interface{}:
			dts = cata(vars[i].Type.T, v)
		case VarMap:
			dts = csts(vars[i].Type.T, v)
		default:
			dts = v
		}
		
		out[vars[i].Data.(string)] = &TVariable{vars[i].Type, dts}
	}

	cart = old_c

	return out
}

func convertValPS(to TType, sk int, dat interface{}) interface{} {
	var numcv float64
	switch v := dat.(type) {
	case []interface{}:
		if isArray(to, sk) {
			return cata(to.T, v)
		} else if isStruct(to, sk) {
			fmt.Println(to)
			return cvsa(to.T, v)
		}
	case VarMap:
		return csts(to.T, v)
	case int:
		numcv = float64(v)
		goto NCV
	case byte:
		numcv = float64(v)
		goto NCV
	case float64:
		numcv = v
		goto NCV
	case bool:
		numcv = 0
		if v {
			numcv = 1
		}
		goto NCV
	}

	errOut(fmt.Sprintf("Unable to convert between two types.\nTO: %v\nSK: %d\nDT: %v", to, sk, dat))
	return nil

	NCV:
	if equateTypePSO(to, tInt, sk) {
		return int(numcv)
	} else if equateTypePSO(to, tFloat, sk) {
		return float64(numcv)
	} else if equateTypePSO(to, tByte, sk) {
		return byte(numcv)
	} else if equateTypePSO(to, tBool, sk) {
		return numcv != 0
	}

	errOut(fmt.Sprintf("Unable to convert between two types.\nTO: %v\nSK: %d\nDT: %v", to, sk, dat))
	return nil
}

func convertVal(dat TVariable, to TType) *TVariable {
	return &TVariable{to, convertValPS(to, 0, dat.Data)}
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
	val, prs := (*ctx)[a.Name]
	if !prs || len(a.Path) != 0 {
		// Try searching the modules for it
		val = searchDef(a)
		if val == nil {
			errOutCTX(fmt.Sprintf("Could not resolve %s in the current context.", a.Name), ctx)
		}
	}
	return val
}

//#################
//# Runtime funcs #
//#################

// Value statement parsing

func isStruct(t TType, skp int) bool {
	ch := false

	ch = ch || isPointer(t, skp)
	ch = ch || isArray(t, skp)
	ch = ch || equateTypePSO(t, tFile, skp)
	ch = ch || equateTypePSO(t, tInt, skp)
	ch = ch || equateTypePSO(t, tByte, skp)
	ch = ch || equateTypePSO(t, tFloat, skp)
	ch = ch || equateTypePSO(t, tCharp, skp)
	ch = ch || equateTypePSO(t, tBool, skp)
	ch = ch || equateTypePSO(t, tNull, skp)

	return !ch
}

func isPointer(t TType, skp int) bool {
	for ;skp < len(t.Pre) && t.Pre[skp] == "const"; skp++ {}

	if len(t.Pre) >= skp {
		return false
	}

	return t.Pre[skp] == "~"
}

func isArray(t TType, skp int) bool {
	if len(t.Pre) <= skp {
		return false
	}

	return t.Pre[skp] == "{}"
}

func evalDotChain(v tparse.Node, ctx *VarMap, wk *TVariable) *TVariable {
	var wrvm *VarMap
	wrvm = ctx

	if isStruct((*wk).Type, 0) {
		wrvm = (*wk).Data.(*VarMap)
	}

	// Check if current name relates to a variable in context or working var
	dat, prs := (*wrvm)[v.Sub[0].Data.Data]
	if prs {
		return evalDotChain(v.Sub[1], ctx, dat)
	}

	//


	return &null
}

func setVal(v tparse.Node, ctx *VarMap, val *TVariable) *TVariable {
	return &null
}

func evalCall() {

}

func evalIndex(v tparse.Node, sk, i int) *TVariable {
	
	return &null
}

// Parse a value node
func evalValue(v tparse.Node, ctx *VarMap) *TVariable {

	// STRUCT/ARRAY DEF
	if v.Data.Data == "comp" {
		out := []interface{}{}

		for i := 0; i < len(v.Sub); i++ {
			tmp := evalValue(v.Sub[i], ctx)
			out = append(out, (*tmp).Data)
		}

		return &TVariable{tStruct, out}
	}

	switch v.Data.Type {
	case tparse.LITERAL:
		t := getLiteralType(v)
		return &TVariable{t, getLiteral(v, t)}
	case tparse.DEFWORD:
		if len(v.Sub) > 0 {
			if v.Sub[0].Data.Data == "index" {
				//return evalIndex()
			} else if v.Sub[0].Data.Data == "call" {
				params := []TVariable{}

				for i := 0; i < len(v.Sub[0].Sub); i++ {
					params = append(params, *evalValue(v.Sub[0].Sub[i], ctx))
				}
				
				out := evalBlock(*searchNode(TArtifact{[]string{}, v.Data.Data}), params)

				return &out
			}
		}

		return resolveArtifact(TArtifact{[]string{}, v.Data.Data}, ctx)

	case tparse.AUGMENT:
		// Special case for =
		if v.Data.Data == "=" {
			return setVal(v.Sub[0], ctx, evalValue(v.Sub[1], ctx))

		} else if v.Data.Data == "." {
			return evalDotChain(v, ctx, &null)

		} else if v.Data.Data == "!" {

			a := evalValue(v.Sub[0], ctx)
			return &TVariable{tBool, !(a.Data.(bool))}
		}

		// General case setup
		
		a, b := evalValue(v.Sub[0], ctx), evalValue(v.Sub[1], ctx)
		var out TVariable
		out.Type = tInt

		// General math and bool cases
		switch v.Data.Data {
		case "+":
			out.Data = a.Data.(int) + b.Data.(int)
		case "-":
			out.Data = a.Data.(int) - b.Data.(int)
		case "*":
			out.Data = a.Data.(int) * b.Data.(int)
		case "/":
			out.Data = a.Data.(int) / b.Data.(int)
		case "%":
			out.Data = a.Data.(int) % b.Data.(int)
		case "&&":
			out.Type = tBool
			out.Data = a.Data.(bool) && b.Data.(bool)
		case "||":
			out.Type = tBool
			out.Data = a.Data.(bool) || b.Data.(bool)
		case "==":
			out.Type = tBool
			if equateType(a.Type, b.Type) {
				out.Data = a.Data == b.Data
			} else {
				out.Data = a.Data.(int) == b.Data.(int)
			}
		case "!=":
			out.Type = tBool
			if equateType(a.Type, b.Type) {
				out.Data = a.Data != b.Data
			} else {
				out.Data = a.Data.(int) != b.Data.(int)
			}
		}

		return &out
	}

	return &null
}

// Eval a definition
func evalDef(v tparse.Node, ctx *VarMap) {
	t := getType(v.Sub[0])
	
	for i := 0; i < len(v.Sub[1].Sub); i++ {
		if v.Sub[1].Sub[i].Data.Data == "=" {
			(*ctx)[v.Sub[1].Sub[i].Sub[0].Data.Data] = convertVal(*evalValue(v.Sub[1].Sub[i].Sub[1], ctx), t)
		} else {
			(*ctx)[v.Sub[1].Sub[i].Data.Data] = &TVariable{t, nil}
		}
		
	}
}

// Eval a control flow
func evalCF(v tparse.Node, ctx *VarMap) (bool, TVariable) {
	//scopeVars := []string{}
	return false, null
}

func evalParams(pd tparse.Node, params *[]TVariable, ctx *VarMap) {
	if len(pd.Sub) == 0 {
		return
	}
	cvt := getType(pd.Sub[0])
	pi := 0
	for i := 1; i < len(pd.Sub); i++ {
		if pd.Sub[i].Data.Type == 10 && pd.Sub[i].Data.Data == "type" {
			cvt = getType(pd.Sub[i])
		} else if pd.Sub[i].Data.Type == tparse.DEFWORD {
			(*ctx)[pd.Sub[i].Data.Data] = convertVal((*params)[pi], cvt)
			pi++
		}
	}
}

func evalBlock(b tparse.Node, params []TVariable) TVariable {
	ctx := make(VarMap)

	var rty TType = tNull

	if b.Sub[0].Data.Data == "bdef" {
		for i := 0; i < len(b.Sub[0].Sub); i++ {
			if b.Sub[0].Sub[i].Data.Data == "[]" {
				rty = getType(b.Sub[0].Sub[i])
			} else if b.Sub[0].Sub[i].Data.Data == "()" {
				evalParams(b.Sub[0].Sub[i], &params, &ctx)
			}
		}
	}

	for i := 0; i < len(b.Sub); i++ {
		switch b.Sub[i].Data.Data {
		case "define":
			evalDef(b.Sub[i], &ctx)
		case "value":
			fmt.Println("--- Eval Value ---")
			fmt.Println(b.Sub[i].Sub[0])
			fmt.Println(*evalValue(b.Sub[i].Sub[0], &ctx))
			fmt.Println("--- End Value ---")
		case "block":
			ret, val := evalCF(b.Sub[i].Sub[0], &ctx)
			if ret {
				return *convertVal(val, rty)
			}
		case "return":
			fmt.Println("--- Block return ---")
			fmt.Println(b.Sub[i].Sub[0].Sub[0])
			ret := *evalValue(b.Sub[i].Sub[0].Sub[0], &ctx)
			fmt.Println(ret)
			fmt.Println("--- Return end ---")
			return *convertVal(ret, rty)
		}
	}

	return null
}

func EvalTNSL(root *TModule, args string) TVariable {
	prog = root
	cart = TArtifact { []string{}, "main" }

	sarg := strings.Split(args, " ")

	saif := []interface{}{}

	for i := 0; i < len(sarg); i++ {
		tmp := []interface{}{}
		dat := []byte(sarg[i])
		for j := 0; j < len(dat); j++ {
			tmp = append(tmp, dat[j])
		}

		saif = append(saif, tmp)
	}

	targ := TVariable {
		TType {
			[]string{"{}", "{}"},
			TArtifact { []string{}, "charp" },
			"" },
		saif }

	mainNode := getNode(prog, "main")

	fmt.Println(mainNode)

	return evalBlock(*mainNode, []TVariable{targ})
}
