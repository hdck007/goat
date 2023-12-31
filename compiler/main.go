package main

import (
	"bytes"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"text/template"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func createCodeForTextNodes(node *html.Node) (*bytes.Buffer, string) {
	codeTemplate := template.New("code")

	finalCode, err := codeTemplate.Parse(`
	{{.ComponentName}} := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		element := goat.CreateVirtualElements(
			"",
			"text",
			nil,
			"{{.data}}",
			goat.VElement{},
		)
		return element
	}, map[string]any{})
	{{.VariableName}} := {{.ComponentName}}(map[string]any{})`)

	if err != nil {
		log.Fatalf("Parse error: %s", err)
	}

	code := bytes.NewBuffer([]byte{})

	variableName := "text" + strconv.Itoa(rand.Intn(1000))

	finalCode.Execute(code, map[string]string{
		"ComponentName": "TEXT" + strconv.Itoa(rand.Intn(1000)),
		"ComponentType": "text",
		"VariableName":  variableName,
		"data":          node.Data,
	})

	return code, variableName
}

func createCodeForContainerNodes(node *html.Node) (*bytes.Buffer, string) {
	codeTemplate := template.New("code")

	childcode := ""
	childvariables := "{"

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if strings.TrimSpace(c.Data) != "" {
			if c.Type != html.TextNode {
				code, variableName := createCodeForContainerNodes(c)
				if code != nil {
					childcode += code.String()
					childvariables += variableName + ", "
				}
			} else {
				code, variableName := createCodeForTextNodes(c)
				if code != nil {
					childcode += code.String()
					childvariables += variableName + ", "
				}
			}
		}
	}
	childvariables = strings.TrimSuffix(childvariables, ", ") + "}"

	finalCode, err := codeTemplate.Parse(`
	{{.ChildCode}}
	{{.ComponentName}} := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		childElements := []goat{{.ChildVariables}}
		element := goat.CreateVirtualElements(
			"",
			"{{.ComponentType}}",
			nil,
			"",
			childElements...,
		)
		return element
	}, map[string]any{})
	{{.VariableName}} := {{.ComponentName}}(map[string]any{})`)

	if err != nil {
		log.Fatalf("Parse error: %s", err)
	}

	code := bytes.NewBuffer([]byte{})

	variableName := strings.ToLower(node.Data) + strconv.Itoa(rand.Intn(1000))

	finalCode.Execute(code, map[string]string{
		"ComponentName":  strings.ToUpper(node.Data) + strconv.Itoa(rand.Intn(1000)),
		"ComponentType":  node.Data,
		"VariableName":   variableName,
		"ChildCode":      childcode,
		"ChildVariables": childvariables,
	})

	return code, variableName
}

func createBoilerPlate(componentsCode string, rootVariableName string) {

	boilerplate := template.New("boilerplate")

	boilerplate, err := boilerplate.Parse(`
	package main
	import (
		"syscall/js"
	
		"github.com/hdck007/goat/goat"
	)
	
	func start() {
		{{.ComponentsCode}}
		body := js.Global().Get("document").Call("getElementById", "root")
		{{.RootVariableName}}.Mount(body)

		select {}
	}
	`)

	if err != nil {
		log.Fatalf("Parse error: %s", err)
	}

	code := bytes.NewBuffer([]byte{})

	boilerplate.Execute(code, map[string]string{
		"ComponentsCode":   componentsCode,
		"RootVariableName": rootVariableName,
	})

	os.WriteFile("example.go", code.Bytes(), 0644)

}

func main() {
	contents, err := os.ReadFile("example.goat")
	if err != nil {
		return
	}
	n, err := html.ParseFragment(bytes.NewReader(contents), &html.Node{
		Type:     html.ElementNode,
		Data:     "body",
		DataAtom: atom.Body,
	})
	if err != nil {
		log.Fatalf("Parse error: %s", err)
	}

	componentsCode := ""

	rootVariableName := ""

	for _, node := range n {
		code, variableName := createCodeForContainerNodes(node)
		rootVariableName = variableName
		if code != nil {
			componentsCode += code.String()
		}
	}

	createBoilerPlate(componentsCode, rootVariableName)
}
