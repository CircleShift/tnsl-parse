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
	"fmt"
	"os"
)

/**
	libtnsl module stub.  Contains only parts of the io sub-module.
	Parts included:
		- io.print
		- io.println
		- io.readFile
		- io.writeFile
		- io.File API for file objects
	
	Types included:
		- tnsl.io.File
		- string ({}charp)
		- charp
		- int
		- float
		- null
*/

// I really hope this works.

// Generic in-built types
var (
	
	tFile = TType{Pre: []string{}, T: TArtifact{Path: []string{"tnsl", "io"}, Name: "File"}, Post: ""}
	tString = TType{Pre: []string{"{}"}, T: TArtifact{Path: []string{}, Name:"uint8"}, Post: ""}
	tInt = TType{Pre: []string{}, T: TArtifact{Path: []string{}, Name:"int"}, Post: ""}
	tUint = TType{Pre: []string{}, T: TArtifact{Path: []string{}, Name:"uint"}, Post: ""}
	tByte = TType{Pre: []string{}, T: TArtifact{Path: []string{}, Name:"uint8"}, Post: ""}
	tFloat = TType{Pre: []string{}, T: TArtifact{Path: []string{}, Name:"float"}, Post: ""}
	tNull = TType{Pre: []string{}, T: TArtifact{Path: []string{}, Name: "null"}, Post: ""}
	tBool = TType{Pre: []string{}, T: TArtifact{Path: []string{}, Name: "bool"}, Post: ""}

	// used only in module definintion
	tEnum = TType{Pre: []string{}, T: TArtifact{Path: []string{}, Name: "enum"}, Post: ""}
	tStruct = TType{Pre: []string{}, T: TArtifact{Path: []string{}, Name: "struct"}, Post: ""}

	// Special types for if chain checking
	tIF = TType{Pre: []string{}, T: TArtifact{Path: []string{}, Name: "if"}, Post: ""}
)

// tells if the stub supports a function
func tnslResolve(callPath TArtifact) int {
	l := len(callPath.Path)
	if l < 2 || l > 3 || callPath.Path[0] != "tnsl" || callPath.Path[1] != "io" {
		return -1
	} else if l > 2 && callPath.Path[2] != "File" {
		return -1
	}

	if l > 2 {
		if callPath.Name == "write" || callPath.Name == "read" || callPath.Name == "close" {
			return 1;
		}
	} else {
		if callPath.Name == "print" || callPath.Name == "println" || callPath.Name == "readFile" || callPath.Name == "writeFile" {
			return 0;
		}
	}

	return -1
}

// evaluate a function call.
// in is the variable in (if any)
// out is the variable out (if any)
// function is the name of the function
func tnslEval(in TVariable, function string) TVariable {
	switch function {
	case "print":
		tprint(in)
	case "println":
		tprintln(in)
	case "readFile":
		return topenReadFile(in)
	case "writeFile":
		return topenWriteFile(in)
	}
	return TVariable{tNull, nil}
}

// evaluate a call on a file object
func tnslFileEval(file, in TVariable, function string) TVariable {
	switch function {
	case "close":
		tfile_close(file)
	case "read":
		return tfile_read(file)
	case "write":
		tfile_write(file, in)
	}
	return TVariable{tNull, nil}
}

// Generic IO funcs

func tprint(in TVariable) {
	if equateType(in.Type, tString) {
		fmt.Print(datToString(in.Data))
	} else {
		fmt.Print(in.Data)
	}
}

func tprintln(in TVariable) {
	if equateType(in.Type, tString) {
		fmt.Println(datToString(in.Data))
	} else {
		fmt.Println(in.Data)
	}
}

func datToString(dat interface{}) string {
	out := []byte{}
	in := dat.([]interface{})
	for i := 0; i < len(in); i++ {
		out = append(out, in[i].(byte))
	}

	return string(out)
}

func topenWriteFile(in TVariable) TVariable {
	if !equateType(in.Type, tString) {
		panic("Tried to open a file (for writing), but did not use a string type for the file name.")
	}
	fd, err := os.Create(datToString(in.Data))
	if err != nil {
		panic(fmt.Sprintf("Failed to open file (for writing) %v as requested by the program. Aborting.\n%v", in.Data, err))
	}
	return  TVariable{tFile, fd}
}

func topenReadFile(in TVariable) TVariable {
	if !equateType(in.Type, tString) {
		panic("Tried to open a file (for reading), but did not use a string type for the file name.")
	}
	fd, err := os.Open(datToString(in.Data))
	if err != nil {
		panic(fmt.Sprintf("Failed to open file (for reading) %v as requested by the program. Aborting.\n%v", in.Data, err))
	}
	return  TVariable{tFile, fd}
}

// File API

// tnsl.io.File.close
func tfile_close(file TVariable) {
	if equateType(file.Type, tFile) {
		(file.Data).(*os.File).Close()
	}
}

// tnsl.io.File.read
func tfile_read(file TVariable) TVariable {
	b := []byte{1}
	_, err := (file.Data).(*os.File).Read(b)
	if err != nil {
		return TVariable{tInt, -1}
	}
	return TVariable{tInt, int(b[0])}
}

// tnsl.io.File.write
func tfile_write(file, in TVariable) {
	if equateType(file.Type, tFile) {
		if equateType(in.Type, tByte) {
			b := []byte{0}
			b[0] = (in.Data).(byte)
			(file.Data).(*os.File).Write(b)
		} else if equateType(in.Type, tString) {
			dat := (in.Data).([]interface{})
			wrt := []byte{}
			for i := 0; i < len(dat); i++ {
				wrt = append(wrt, dat[i].(byte))
			}
			(file.Data).(*os.File).Write(wrt)
		}
	} else {
		(file.Data).(*os.File).Close()
		panic(fmt.Sprintf("Failed to write to file, attempted to use unsupported type (%v)\n", in.Type))
	}
}
