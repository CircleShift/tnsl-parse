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

// TVaraiable represents a single variable in the program
type TVariable struct {
	Type string
	Data interface{}
}

// TPath represents a pointer to the current module and file
// that the thread is working in.
type TPath struct {
	Module     []string,
	Artifact   string
}

// TContext represents a single thread.
type TContext struct {
	CallStack []Node,
	CallEnv   []TPath,
	VarMap    []map[string]TVariable
}

// TModule represents a collection of files and sub-modules in a program
type TModule struct {
	Files   []Node,
	Globals []map[string]TVariable
	Sub     []TModule
}

// TWorld represents the full program
type TWorld struct {
	Modules  []TModule,
	MainPath TPath,
	MainFunc Node
}
