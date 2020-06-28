package main

import "fmt"
import "tparse"
import "flag"
import "os"

func main() {
	inputFile := flag.String("in", "", "The file to parse")
	outputFile := flag.String("out", "out.tnp", "The file to store the parse in")

	flag.Parse()

	fd, err := os.Create(*outputFile)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fd.WriteString(fmt.Sprint(tparse.ParseFile(*inputFile)))

	fd.Close()
}
