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

package main

import "fmt"
import "tparse"
import "flag"
import "os"

func main() {
	inputFile := flag.String("in", "", "The file to parse")
	outputFile := flag.String("out", "out.tnt", "The file to store the node tree")

	flag.Parse()

	fd, err := os.Create(*outputFile)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	tokens := tparse.TokenizeFile(*inputFile)
	tree := tparse.MakeTree(&tokens, *inputFile)

	fd.WriteString(fmt.Sprint(tree) + "\n")

	fd.Close()
}
