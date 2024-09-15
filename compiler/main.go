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

func createCodeForTextNodes(text string, dynamicMap *DynamicMapping, parentChildArrayName string) (*bytes.Buffer, []string) {
	codeTemplate := template.New("code")
	variableName := "text" + strconv.Itoa(rand.Intn(1000))
	variableNames := []string{}

	finalCode, err := codeTemplate.Parse(`
		{{.VariableName}} := goat.CreateVirtualElements(
			"",
			"text",
			nil,
			{{.data}},
		)
		{{.ParentChildArrayName}} = append({{.ParentChildArrayName}}, {{.VariableName}})
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
				currentTextNodeCode, _ := createCodeForTextNodes(node, dynamicMap, parentChildArrayName)
				code.WriteString(currentTextNodeCode.String())
			}
			if i != len(textNodes)-1 {
				dynamicName := strings.TrimSuffix(strings.TrimPrefix(dynamic, "{"), "}")
				finalCode.Execute(code, map[string]string{
					"VariableName":         variableName,
					"data":                 "prop[proxy.Get(\"" + strings.ReplaceAll(dynamicName, "\n", "") + "\").Key]",
					"ParentChildArrayName": parentChildArrayName,
				})
			}
		}

		return code, variableNames
	}

	finalCode.Execute(code, map[string]string{
		"VariableName":         variableName,
		"data":                 "\"" + strings.ReplaceAll(text, "\n", "") + "\"",
		"ParentChildArrayName": parentChildArrayName,
	})

	return code, variableNames
}

func createCodeForChildContainers(node *html.Node, parentChildrenArrayName string, dynamicMap *DynamicMapping, parentHtmlReference string, parentComponentName string, parentBlockName string) (*bytes.Buffer, string, string) {

	if node.Type == 1 {
		return nil, "", ""
	}

	codeTemplate := template.New("code")
	functionCodeTemplate, functionTemplateErr := template.New("function").Parse(`
		var {{.VariableName}} js.Func
		{{.VariableName}} = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			{{.FunctionName}}()
			{{.UpdateBindings}}
			return nil
		})
		{{.HTMLReference}}.Call("addEventListener", "{{.Event}}", {{.VariableName}})
	`)

	if functionTemplateErr != nil {
		log.Fatalf("Parse error: %s", functionTemplateErr)
	}

	variableName := strings.ToLower(node.Data) + strconv.Itoa(rand.Intn(1000))
	functionCode := bytes.NewBuffer([]byte{})
	nodeArrayName := "children" + strconv.Itoa(rand.Intn(1000))

	childCode := ""

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if strings.TrimSpace(c.Data) != "" {
			if c.Type != html.TextNode {
				code, _, childFunctionCode := createCodeForChildContainers(c, nodeArrayName, dynamicMap, parentHtmlReference, parentComponentName, parentBlockName)
				functionCode.WriteString(childFunctionCode)
				if code != nil {
					childCode += code.String()
				}
			} else {
				code, _ := createCodeForTextNodes(c.Data, dynamicMap, nodeArrayName)
				if code != nil {
					childCode += code.String()
				}
			}
		}
	}

	props := ""
	if len(node.Attr) > 0 {
		props += "map[string]any{" + "\n"
		for _, attr := range node.Attr {
			if strings.Contains(attr.Key, "@") {
				eventName := strings.ReplaceAll(attr.Key, "@", "")
				functionName := attr.Val
				functionVariableName := strings.ToLower(functionName) + strconv.Itoa(rand.Intn(1000))

				updateBindings := ""

				for _, dependentVariable := range (*dynamicMap)["name"] {
					println(dependentVariable)
					updateBindings += parentComponentName + ".Patch(" + parentBlockName + "(map[string]any{" + "\n" + "\"" + "name" + "\": " + "name" + ",\n" + "}))" + "\n"
				}

				functionCodeTemplate.Execute(functionCode, map[string]string{
					"VariableName":   functionVariableName,
					"FunctionName":   functionName,
					"Event":          eventName,
					"UpdateBindings": updateBindings,
					"HTMLReference":  parentHtmlReference,
				})
			} else {
				props += "\"" + attr.Key + "\": \"" + attr.Val + "\",\n"
			}
		}
		props += "}"
	} else {
		props = "nil"
	}

	finalCode, err := codeTemplate.Parse(`
		{{.NodeArrayName}} := []goat.VElement{}
		{{.ChildCode}}
		{{.VariableName}} := goat.CreateVirtualElements(
			"",
			"{{.ComponentType}}",
			{{.Props}},
			"",
			{{.NodeArrayName}}...,
		)
		{{.ParentArrayName}} = append({{.ParentArrayName}}, {{.VariableName}})
	`)

	if err != nil {
		log.Fatalf("Parse error: %s", err)
	}

	code := bytes.NewBuffer([]byte{})

	finalCode.Execute(code, map[string]string{
		"NodeArrayName":   nodeArrayName,
		"ChildCode":       childCode,
		"ComponentType":   node.Data,
		"Props":           props,
		"ParentArrayName": parentChildrenArrayName,
		"VariableName":    variableName,
	})

	return code, variableName, functionCode.String()

}

func createCodeForContainerNodes(node *html.Node, parentHtml string, dynamicMap *DynamicMapping) (*bytes.Buffer, string, string) {

	if node.Type == 1 {
		return nil, "", ""
	}

	codeTemplate := template.New("code")
	functionCodeTemplate, functionTemplateErr := template.New("function").Parse(`
		var {{.VariableName}} js.Func
		{{.VariableName}} = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			{{.FunctionName}}()
			{{.UpdateBindings}}
			return nil
		})
		{{.ComponentName}}.Call("addEventListener", "{{.Event}}", {{.VariableName}})
	`)

	if functionTemplateErr != nil {
		log.Fatalf("Parse error: %s", functionTemplateErr)
	}

	componentName := strings.ToUpper(node.Data) + strconv.Itoa(rand.Intn(1000))
	variableName := strings.ToLower(node.Data) + strconv.Itoa(rand.Intn(1000))
	htmlReferenceVariableName := strings.ToLower(node.Data) + strconv.Itoa(rand.Intn(1000))
	functionCode := bytes.NewBuffer([]byte{})

	binding := htmlReferenceVariableName + " := " + variableName + ".Mount(" + parentHtml + ")"

	childCode := ""
	childBindings := ""

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if strings.TrimSpace(c.Data) != "" {
			if c.Type != html.TextNode {
				code, _, childFunctionCode := createCodeForChildContainers(c, "children", dynamicMap, htmlReferenceVariableName, variableName, componentName)
				functionCode.WriteString(childFunctionCode)
				if code != nil {
					childCode += code.String()
				}
			} else {
				code, _ := createCodeForTextNodes(c.Data, dynamicMap, "children")
				if code != nil {
					childCode += code.String()
				}
			}
		}
	}

	props := ""
	if len(node.Attr) > 0 {
		props += "map[string]any{" + "\n"
		for _, attr := range node.Attr {
			if strings.Contains(attr.Key, "@") {
				eventName := strings.ReplaceAll(attr.Key, "@", "")
				functionName := attr.Val
				functionVariableName := strings.ToLower(functionName) + strconv.Itoa(rand.Intn(1000))

				updateBindings := ""

				for _, dependentVariable := range (*dynamicMap)["name"] {
					println(dependentVariable)
					updateBindings += dependentVariable + ".Patch(map[string]any{" + "\n" + "\"" + "name" + "\": " + "name" + ",\n" + "})" + "\n"
				}

				functionCodeTemplate.Execute(functionCode, map[string]string{
					"VariableName":   functionVariableName,
					"FunctionName":   functionName,
					"Event":          eventName,
					"UpdateBindings": updateBindings,
					"ComponentName":  htmlReferenceVariableName,
				})
			} else {
				props += "\"" + attr.Key + "\": \"" + attr.Val + "\",\n"
			}
		}
		props += "}"
	} else {
		props = "nil"
	}

	finalCode, err := codeTemplate.Parse(`
	{{.ComponentName}} := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		children := []goat.VElement{}
		{{.ChildCode}}
		element := goat.CreateVirtualElements(
			"",
			"{{.ComponentType}}",
			{{.Props}},
			"",
			children...,
		)
		return element
	}, map[string]any{})
	{{.VariableName}} := {{.ComponentName}}(map[string]any{})
	{{.Binding}}
	`)

	if err != nil {
		log.Fatalf("Parse error: %s", err)
	}

	code := bytes.NewBuffer([]byte{})

	finalCode.Execute(code, map[string]string{
		"ComponentName": componentName,
		"ComponentType": node.Data,
		"VariableName":  variableName,
		"ChildCode":     childCode,
		"ChildBindings": childBindings,
		"Binding":       binding,
		"Props":         props,
	})

	return code, variableName, functionCode.String()
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

func createBoilerPlate(componentsCode string, rootVariableName string, scriptCode string, functionCode string) {

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
		{{.FunctionCode}}
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
		"FunctionCode":     functionCode,
	})

	os.WriteFile("example.go", code.Bytes(), 0644)

}

func analyzeNode(node *html.Node, dynamicMap *DynamicMapping) {
	if node.Type == html.ElementNode && node.Data == "script" {
		return
	}
	r, _ := regexp.Compile("{.*}")

	match := r.MatchString(node.Data)

	if match {
		matchingString := r.FindString(node.Data)
		nameOfVariable := strings.ReplaceAll(strings.ReplaceAll(matchingString, "{", ""), "}", "")
		println(nameOfVariable)
		(*dynamicMap)[nameOfVariable] = append((*dynamicMap)[nameOfVariable], "hello")
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

	functionCode := ""

	dynamicMap := make(DynamicMapping)

	for _, node := range n {
		analyzeNode(node, &dynamicMap)
	}

	for _, node := range n {
		if node.Type == html.ElementNode && node.Data == "script" {
			scriptCode += returnCodeForScriptNode(node)
			continue
		}
		code, variableName, childFunctionCode := createCodeForContainerNodes(node, "body", &dynamicMap)
		rootVariableName = variableName
		functionCode += childFunctionCode
		if code != nil {
			componentsCode += code.String()
		}
	}

	createBoilerPlate(componentsCode, rootVariableName, scriptCode, functionCode)
}

func main() {
	// generateCode("example.goat")
	value, err := strconv.ParseBool("")

	if err == nil {
		println(err)
	}

	println(value)
}