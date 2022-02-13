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

import "tparse"

// TArtifact represents the path to a specific named object in the node tree.
type TArtifact struct {
	Path       []string
	Name       string
}

// TType represents the type of a variable (including pre and post unary ops)
type TType struct {
	Pre  []string
	T    TArtifact
	Post string
}

// TVariable represents a single variable in the program
type TVariable struct {
	Type TType
	Data interface{}
}

type VarMap map[string]*TVariable

// TModule represents a collection of files and sub-modules in a program
type TModule struct {
	Name       string
	Artifacts  []tparse.Node
	Defs       VarMap
	Sub        []TModule
}

