package main

import (
	"fmt"
	"os"
)

func main() {
	b, err := os.ReadFile("./data/example.goat") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	ast, error := Parse(string(b))

	if error != nil {
		fmt.Println("Error occured while parsing", error)
		return
	}

	astJson, printingError := ast.ToJSON()

	if printingError != nil {
		println("Error while printing")
	}

	os.WriteFile("ast.json", []byte(astJson), 0644)
}
