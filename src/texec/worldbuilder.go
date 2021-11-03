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
	"path"
)

/**
	worldbuilder.go - take in a file name and construct a root TModule based on it.
*/

func parseFile(p string) tparse.Node {
	tokens := tparse.TokenizeFile(p)
	return tparse.MakeTree(&(tokens), p)
}

func buildModule(module tparse.Node) TModule {
	out := TModule{}

	for n := 0 ; n < len(module.Sub) ; n++ {

		switch module.Sub[n].Data.Type {
		case 11:
			
		case 10:

		}
	}

	return out
}

// BuildRoot builds the root module, ready for eval
func BuildRoot(file tparse.Node) TModule {
	out := TModule{}

	out.Files = append(out.Files, file)

	for n := 0 ; n < len(file.Sub) ; n++ {

		switch file.Sub[n].Data.Type {
		case 11:

		case 10:
		}
	}

	return out
}
