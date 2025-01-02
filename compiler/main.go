package main

import (
	"fmt"
	"os"
)

func main() {

	// create a new folder dist if it doesn't exist if exists delete all the files in it
	if _, err := os.Stat("dist"); os.IsNotExist(err) {
		os.Mkdir("dist", 0755)
	} else {
		os.RemoveAll("dist/*")
	}

	ast, err := ParseFile(os.Args[1])
	if err != nil {
		fmt.Println("Error occurred while parsing: ", err)
		return
	}

	// astJson, err := ast.ToJSON()
	// if err != nil {
	// 	fmt.Print(err)
	// }

	Generate(*ast, "dist/App.go", "App", os.Args[1])

	mainFunctionCode := `
		package main

		import (
			"syscall/js"
			"github.com/hdck007/goat/goat"
		)

		func main() {
			done := make(chan struct{}, 0)
			component := App(goat.Props{})
			root := js.Global().Get("document").Call("getElementById", "root")
			component.Mount(root)

			<-done
		}
	`

	os.WriteFile("dist/main.go", []byte(mainFunctionCode), 0644)

	// // run GOOS=js GOARCH=wasm go build ./dist -o main.wasm
	// cmd := exec.Command("GOOS=js", "GOARCH=wasm", "go", "build", "-o", "main.wasm", "./dist")
	// stdout, err := cmd.Output()

	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }

	// fmt.Println(string(stdout))

}
