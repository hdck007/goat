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
		fmt.Println("Error occurred while parsing: ", error)
		return
	}

	generatedCode := Generate(*ast)

	os.WriteFile("output.goat", []byte(generatedCode), 0644)
}
