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

func createCodeForTextNodes(text string, listOfVariables *[]string, parentChildArrayName string) (*bytes.Buffer, []string) {
	codeTemplate := template.New("code")
	variableName := "text" + strconv.Itoa(rand.Intn(1000))
	variableNames := []string{}

	finalCode, err := codeTemplate.Parse(`
		{{.VariableName}} := goat.CreateVirtualElements(
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
				currentTextNodeCode, childVariableNames := createCodeForTextNodes(node, listOfVariables, parentChildArrayName)
				variableNames = append(variableNames, childVariableNames...)
				code.WriteString(currentTextNodeCode.String())
			}
			if i != len(textNodes)-1 {
				dynamicName := strings.TrimSuffix(strings.TrimPrefix(dynamic, "{"), "}")
				variableNames = append(variableNames, dynamicName)
				finalCode.Execute(code, map[string]string{
					"VariableName":         variableName,
					"data":                 "proxy.Get(\"" + strings.ReplaceAll(dynamicName, "\n", "") + "\")",
					"ParentChildArrayName": parentChildArrayName,
				})
			}
		}

		return code, variableNames
	}

	finalCode.Execute(code, map[string]string{
		"VariableName":         variableName,
		"data":                 "&goat.UnionNode{Element: nil, StringValue:" + "\"" + strings.ReplaceAll(text, "\n", "") + "\"" + "}",
		"ParentChildArrayName": parentChildArrayName,
	})

	return code, variableNames
}

func createCodeForChildContainers(node *html.Node, parentChildrenArrayName string, listOfVariables *[]string, parentHtmlReference string, parentComponentName string, parentBlockName string) (*bytes.Buffer, string, string, []string) {

	if node.Type == 1 {
		return nil, "", "", []string{}
	}

	dependentVariables := []string{}

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
				code, _, childFunctionCode, childVariableNames := createCodeForChildContainers(c, nodeArrayName, listOfVariables, parentHtmlReference, parentComponentName, parentBlockName)
				dependentVariables = append(dependentVariables, childVariableNames...)
				functionCode.WriteString(childFunctionCode)
				if code != nil {
					childCode += code.String()
				}
			} else {
				code, childVariableNames := createCodeForTextNodes(c.Data, listOfVariables, nodeArrayName)
				dependentVariables = append(dependentVariables, childVariableNames...)
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

				dependencyMap := createDependencyMap(*listOfVariables)

				bindings := parentComponentName + ".Patch(" + parentBlockName + "(map[string]any{" + dependencyMap + "}))" + "\n"

				functionCodeTemplate.Execute(functionCode, map[string]string{
					"VariableName":   functionVariableName,
					"FunctionName":   functionName,
					"Event":          eventName,
					"UpdateBindings": bindings,
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
		{{.NodeArrayName}} := []goat.Vnode{}
		{{.ChildCode}}
		{{.VariableName}} := goat.CreateVirtualElements(
			"{{.ComponentType}}",
			{{.Props}},
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

	return code, variableName, functionCode.String(), dependentVariables

}

func createDependencyMap(variables []string) string {

	codeTemplate := template.New("code")
	dependencyMap := bytes.NewBuffer([]byte{})

	codeTemplate, err := codeTemplate.Parse(`
		"{{.VariableName}}": {{.VariableName}},
	`)

	if err != nil {
		log.Fatalf("Parse error: %s", err)
	}

	for _, variable := range variables {
		codeTemplate.Execute(dependencyMap, map[string]string{
			"VariableName": variable,
		})
	}

	return dependencyMap.String()
}

func createCodeForContainerNodes(node *html.Node, parentHtml string, listOfVariables *[]string) (*bytes.Buffer, string, string) {

	if node.Type == 1 {
		return nil, "", ""
	}

	dependentVariableNames := []string{}

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
				code, _, childFunctionCode, childVariableNames := createCodeForChildContainers(c, "children", listOfVariables, htmlReferenceVariableName, variableName, componentName)
				dependentVariableNames = append(dependentVariableNames, childVariableNames...)
				functionCode.WriteString(childFunctionCode)
				if code != nil {
					childCode += code.String()
				}
			} else {
				code, _ := createCodeForTextNodes(c.Data, listOfVariables, "children")
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

				dependencyMap := createDependencyMap(*listOfVariables)

				bindings := variableName + ".Patch(" + componentName + "(map[string]any{" + dependencyMap + "}))" + "\n"

				functionCodeTemplate.Execute(functionCode, map[string]string{
					"VariableName":   functionVariableName,
					"FunctionName":   functionName,
					"Event":          eventName,
					"UpdateBindings": bindings,
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

	dependencyMap := createDependencyMap(dependentVariableNames)

	finalCode, err := codeTemplate.Parse(`
	{{.ComponentName}} := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.Vnode {
		children := []goat.Vnode{}
		{{.ChildCode}}
		element := goat.CreateVirtualElements(
			"{{.ComponentType}}",
			{{.Props}},
			children...,
		)
		return element
	}, map[string]any{
		{{.DependencyMap}}
	})
	{{.VariableName}} := {{.ComponentName}}(map[string]any{
		{{.DependencyMap}}
	})
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
		"DependencyMap": dependencyMap,
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

func analyzeNode(node *html.Node, dynamicMap *[]string) {
	if node.Type == html.ElementNode && node.Data == "script" {
		return
	}
	r, _ := regexp.Compile("{[^{}]*}")

	matches := r.FindAllStringSubmatch(node.Data, -1)
	for _, v := range matches {
		matchingString := v[0]
		nameOfVariable := strings.ReplaceAll(strings.ReplaceAll(matchingString, "{", ""), "}", "")

		println(nameOfVariable)

		(*dynamicMap) = append((*dynamicMap), nameOfVariable)
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

	dynamicMap := []string{}

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
	generateCode("example.goat")
}
