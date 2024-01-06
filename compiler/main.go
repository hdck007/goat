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

type DynamicMapping map[string][]string

func createCodeForTextNodes(text string, parentHtml string, dynamicMap *DynamicMapping) (*bytes.Buffer, []string) {
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
				currentTextNodeCode, _ := createCodeForTextNodes(node, parentHtml, dynamicMap)
				code.WriteString(currentTextNodeCode.String())
			}
			if i != len(textNodes)-1 {
				dynamicName := strings.TrimSuffix(strings.TrimPrefix(dynamic, "{"), "}")
				(*dynamicMap)[dynamicName] = append((*dynamicMap)[dynamicName], variableName)
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

func createCodeForContainerNodes(node *html.Node, parentHtml string, dynamicMap *DynamicMapping) (*bytes.Buffer, string) {

	if node.Type == 1 {
		return nil, ""
	}

	codeTemplate := template.New("code")
	functionCodeTemplate, functionTemplateErr := template.New("function").Parse(`
		var {{.VariableName}} js.Func
		cb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			{{.FunctionName}}()
			{{.UpdateBindings}}
			return nil
		})
		{{.VariableName}}.Call("addEventListener", "{{.Event}}", cb)
	`)

	if functionTemplateErr != nil {
		log.Fatalf("Parse error: %s", functionTemplateErr)
	}

	variableName := strings.ToLower(node.Data) + strconv.Itoa(rand.Intn(1000))
	htmlReferenceVariableName := strings.ToLower(node.Data) + strconv.Itoa(rand.Intn(1000))
	functionCode := bytes.NewBuffer([]byte{})

	props := ""
	if len(node.Attr) > 0 {
		props += "map[string]any{" + "\n"
		for _, attr := range node.Attr {
			if strings.Contains(attr.Key, "@") {
				eventName := strings.ReplaceAll(attr.Key, "@", "")
				functionName := attr.Val
				functionVariableName := strings.ToLower(functionName) + strconv.Itoa(rand.Intn(1000))

				// for dependentVariableName := range (*dynamicMap)[functionName] {

				// }

				functionCodeTemplate.Execute(functionCode, map[string]string{
					"VariableName":   functionVariableName,
					"FunctionName":   functionName,
					"Event":          eventName,
					"UpdateBindings": "",
				})
			} else {
				props += "\"" + attr.Key + "\": \"" + attr.Val + "\",\n"
			}
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
				code, _ := createCodeForContainerNodes(c, htmlReferenceVariableName, dynamicMap)
				if code != nil {
					childCode += code.String()
				}
			} else {
				code, _ := createCodeForTextNodes(c.Data, htmlReferenceVariableName, dynamicMap)
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

func analyzeNode(node *html.Node, dynamicMap *DynamicMapping) {
	if node.Type == html.ElementNode && node.Data == "script" {
		return
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		analyzeNode(c, dynamicMap)
	}
}

func generateCode(fileName string) {
	contents, err := os.ReadFile(fileName)
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

	dynamicMap := make(DynamicMapping)

	// for _, node := range n {
	// 	if node.Type == html.ElementNode && node.Data == "script" {
	// 		continue
	// 	}
	// 	analyseNode(node, &dynamicMap)
	// }

	for _, node := range n {
		if node.Type == html.ElementNode && node.Data == "script" {
			scriptCode += returnCodeForScriptNode(node)
			continue
		}
		code, variableName := createCodeForContainerNodes(node, "body", &dynamicMap)
		rootVariableName = variableName
		if code != nil {
			componentsCode += code.String()
		}
	}

	createBoilerPlate(componentsCode, rootVariableName, scriptCode)
}

func main() {
	generateCode("example.goat")
}
