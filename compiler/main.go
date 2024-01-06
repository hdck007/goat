package main

import (
	"bytes"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func createCodeForTextNodes(text string, parentHtml string) (*bytes.Buffer, []string) {
	codeTemplate := template.New("code")
	variableName := "text" + strconv.Itoa(rand.Intn(1000))
	binding := variableName + ".Mount(" + parentHtml + ")"
	variableNames := []string{}

	finalCode, err := codeTemplate.Parse(`
	{{.ComponentName}} := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		element := goat.CreateVirtualElements(
			"",
			"text",
			nil,
			{{.data}},
		)
		return element
	}, {{.Props}})
	{{.VariableName}} := {{.ComponentName}}(map[string]any{})
	{{.Binding}}
	`)

	if err != nil {
		log.Fatalf("Parse error: %s", err)
	}

	code := bytes.NewBuffer([]byte{})

	r, _ := regexp.Compile("{.*}")

	match := r.MatchString(text)

	if match {
		dynamic := r.FindString(text)
		textNodes := strings.Split(text, dynamic)

		for i, node := range textNodes {

			if node != "" {
				currentTextNodeCode, _ := createCodeForTextNodes(node, parentHtml)
				code.WriteString(currentTextNodeCode.String())
			}
			if i != len(textNodes)-1 {
				dynamicName := strings.TrimSuffix(strings.TrimPrefix(dynamic, "{"), "}")
				finalCode.Execute(code, map[string]string{
					"ComponentName": "TEXT" + strconv.Itoa(rand.Intn(1000)),
					"ComponentType": "text",
					"VariableName":  variableName,
					"data":          "prop[proxy.Get(\"" + strings.ReplaceAll(dynamicName, "\n", "") + "\").Key]",
					"Binding":       binding,
					"Props":         "map[string]any{" + "\"" + strings.ReplaceAll(dynamicName, "\n", "") + "\"" + ": " + strings.ReplaceAll(dynamicName, "\n", "") + "}",
				})
			}
		}

		return code, variableNames
	}

	finalCode.Execute(code, map[string]string{
		"ComponentName": "TEXT" + strconv.Itoa(rand.Intn(1000)),
		"ComponentType": "text",
		"VariableName":  variableName,
		"data":          "\"" + strings.ReplaceAll(text, "\n", "") + "\"",
		"Binding":       binding,
		"Props":         "map[string]any{}",
	})

	return code, variableNames
}

func createCodeForContainerNodes(node *html.Node, parentHtml string) (*bytes.Buffer, string) {

	if node.Type == 1 {
		return nil, ""
	}

	codeTemplate := template.New("code")
	variableName := strings.ToLower(node.Data) + strconv.Itoa(rand.Intn(1000))
	htmlReferenceVariableName := strings.ToLower(node.Data) + strconv.Itoa(rand.Intn(1000))

	props := ""
	if len(node.Attr) > 0 {
		props += "map[string]any{" + "\n"
		for _, attr := range node.Attr {
			props += "\"" + attr.Key + "\": \"" + attr.Val + "\",\n"
		}
		props += "}"
	} else {
		props = "nil"
	}

	binding := htmlReferenceVariableName + " := " + variableName + ".Mount(" + parentHtml + ")"

	childCode := ""
	childBindings := ""

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if strings.TrimSpace(c.Data) != "" {
			if c.Type != html.TextNode {
				code, _ := createCodeForContainerNodes(c, htmlReferenceVariableName)
				if code != nil {
					childCode += code.String()
				}
			} else {
				code, _ := createCodeForTextNodes(c.Data, htmlReferenceVariableName)
				if code != nil {
					childCode += code.String()
				}
			}
		}
	}

	finalCode, err := codeTemplate.Parse(`
	{{.ComponentName}} := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		element := goat.CreateVirtualElements(
			"",
			"{{.ComponentType}}",
			{{.Props}},
			"",
		)
		return element
	}, map[string]any{})
	{{.VariableName}} := {{.ComponentName}}(map[string]any{})
	{{.Binding}}
	{{.ChildCode}}
	`)

	if err != nil {
		log.Fatalf("Parse error: %s", err)
	}

	code := bytes.NewBuffer([]byte{})

	finalCode.Execute(code, map[string]string{
		"ComponentName": strings.ToUpper(node.Data) + strconv.Itoa(rand.Intn(1000)),
		"ComponentType": node.Data,
		"VariableName":  variableName,
		"ChildCode":     childCode,
		"ChildBindings": childBindings,
		"Binding":       binding,
		"Props":         props,
	})

	return code, variableName
}

func returnCodeForScriptNode(node *html.Node) string {
	code := ""

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if strings.TrimSpace(c.Data) != "" {
			code += c.Data + "\n"
		}
	}

	return code
}

func createBoilerPlate(componentsCode string, rootVariableName string, scriptCode string) {

	boilerplate := template.New("boilerplate")

	boilerplate, err := boilerplate.Parse(`
	package main
	import (
		"syscall/js"
	
		"github.com/hdck007/goat/goat"
	)
	
	func start() {
		{{.Script}}
		body := js.Global().Get("document").Call("getElementById", "root")
		{{.ComponentsCode}}
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
		"Script":           scriptCode,
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

	scriptCode := ""

	rootVariableName := ""

	for _, node := range n {
		if node.Type == html.ElementNode && node.Data == "script" {
			scriptCode += returnCodeForScriptNode(node)
			continue
		}
		code, variableName := createCodeForContainerNodes(node, "body")
		rootVariableName = variableName
		if code != nil {
			componentsCode += code.String()
		}
	}

	createBoilerPlate(componentsCode, rootVariableName, scriptCode)
}
