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

import "fmt"

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
// out is the variable out (if any)
// in is the variable in (if any)
// callPath is the function being called.
func tnslEval(out, in *TVaraiable, callPath TPath) {

}