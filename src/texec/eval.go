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

// Don't want to deal with this rn

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

//#################
//# Runtime funcs #
//#################

