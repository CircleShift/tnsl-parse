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
		- io.open_file
		- io.File API for file objects
*/

// tells if the stub supports a function
func tnslResolve(callPath TPath) bool {
	l := len(callPath.Module)
	if l < 2 || l > 3 || callPath.Module[0] != "tnsl" || callPath.Module[1] != "io" {
		return false
	}
	if l > 2 && callPath.Module[2] != "File" {
		return false
	}

	if l > 2 {
		if callPath.Artifact == "write" || callPath.Artifact == "read" || callPath.Artifact == "close" {
			return true;
		}
	} else {
		if callPath.Artifact == "print" || callPath.Artifact == "println" || callPath.Artifact == "open_file" {
			return true;
		}
	}

	return false
}

// evaluate a function call.
// in is the variable in (if any)
// out is the variable out (if any)
// function is the name of the function
func tnslEval(in, out *TVariable, function string) {
	switch function {
	case "print":
		tprint(*in)
	case "println":
		tprintln(*in)
	case "open_file":
		topen_file(*in, out)
	}
}

// evaluate a call on a file object
func tnslFileEval(file, in, out *TVariable, function string) {
	switch function {
	case "close":
		tfile_close(file)
	case "read":
		tfile_read(file, out)
	case "write":
		tfile_write(file, in)
	}
}

// Generic IO funcs

func tprint(in TVariable) {
	fmt.Printf("%v", in.Data)
}

func tprintln(in TVariable) {
	fmt.Printf("%v\n", in.Data)
}

func topen_file(in TVariable, out *TVariable) {
	if in.Type != "string" {
		panic("Tried to open a file, but did not use a string type for the file name.")
	}
	fd, err := os.Create(in.Data.(string))
	if err != nil {
		panic(fmt.Sprintf("Failed to open file %v as requested by the program. Aborting.\n%v", in.Data, err))
	}
	out.Type = "tnsl.io.File"
	out.Data = fd
}


// File API

// tnsl.io.File.close
func tfile_close(file *TVariable) {
	if file.Type == "tnsl.io.File" {
		(file.Data).(*os.File).Close()
	}
}

// tnsl.io.File.read
func tfile_read(file, out *TVariable) {
	b := []byte{1}
	(file.Data).(*os.File).Read(b)
	if out.Data == "uint8" || out.Data == "int8" {
		out.Data = b[0]
	}
}

// tnsl.io.File.write
func tfile_write(file, in *TVariable) {
	b := []byte{0}
	if in.Data == "uint8" || in.Data == "int8" {
		b[0] = (in.Data).(byte)
	} else {
		(file.Data).(*os.File).Close()
		panic(fmt.Sprintf("Failed to write to file, attempted to use unsupported type (%v)\n", in.Type))
	}
	(file.Data).(*os.File).Write(b)
}
