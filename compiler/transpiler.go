package main

import (
	"fmt"
	"go/ast"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type CompilerVisitor struct {
	scriptBuffer    strings.Builder
	htmlBuffer      strings.Builder
	componentBuffer strings.Builder
	variables       []string
}

func getFileRelativeToCurrentFile(currentFilePath string, importFilePath string) string {
	directoryForCurrenFile := filepath.Dir(currentFilePath)
	return filepath.Join(directoryForCurrenFile, importFilePath)
}

func Generate(ast AST, outputFileName string, functionName string, currentFileName string) {

	if ast.Imports == nil {
		ast.Imports = &Imports{
			Imports: []Import{},
		}
	}

	compilerVisitor := CompilerVisitor{
		variables:    []string{},
		htmlBuffer:   strings.Builder{},
		scriptBuffer: strings.Builder{},
	}

	golangImports := make([]string, 0)

	for _, importStatement := range ast.Imports.Imports {

		if importStatement.Type == "golangImport" {
			golangImports = append(golangImports, importStatement.Value.(string))
			continue
		}

		fileName := importStatement.File
		componentName := importStatement.Name

		relativeFileName := getFileRelativeToCurrentFile(currentFileName, fileName)

		fileNameWithoutExtension := strings.TrimSuffix(fileName, ".goat")

		ast, error := ParseFile(relativeFileName)
		if error != nil {
			fmt.Println("Error occurred while parsing: ", error)
		}

		Generate(*ast, fmt.Sprintf("dist/%s.go", fileNameWithoutExtension), componentName, relativeFileName)
	}

	compilerVisitor.visitFragment(ast.HTML[0])

	if ast.Script != nil {
		compilerVisitor.scriptBuffer.WriteString(ast.Script.Value)
	}

	compilerVisitor.scriptBuffer.WriteString(compilerVisitor.componentBuffer.String())
	compilerVisitor.scriptBuffer.WriteString("derivedProps = &goat.Props{")

	for _, variable := range compilerVisitor.variables {
		compilerVisitor.scriptBuffer.WriteString(fmt.Sprintf("%q: %s,", variable, variable))
	}

	compilerVisitor.scriptBuffer.WriteString("}")
	compilerVisitor.scriptBuffer.WriteString(`
		for k, v := range props {
			(*derivedProps)[k] = v
		}
	`)

	compilerVisitor.scriptBuffer.WriteString("propsToPass := goat.Props{")

	for _, variable := range compilerVisitor.variables {
		compilerVisitor.scriptBuffer.WriteString(fmt.Sprintf("%q: %s,", variable, variable))
	}

	compilerVisitor.scriptBuffer.WriteString("}")
	compilerVisitor.scriptBuffer.WriteString(`
		for k, v := range props {
			propsToPass[k] = v
		}
	`)

	randomComponentName := fmt.Sprintf("block%d", rand.Int())

	importString := strings.Builder{}

	if len(golangImports) > 0 {
		importString.WriteString("import (\n")
	}

	for _, golangImport := range golangImports {
		importString.WriteString(fmt.Sprintf("  %q\n", golangImport))
	}

	if len(golangImports) > 0 {
		importString.WriteString(")\n")
	}

	template := crateTemplate("component-template", fmt.Sprintf(`
		package main

		{{.Imports}}

		func %s(props goat.Props) goat.Block {
			%s := &goat.Block{
				Patch: func(goat.Block) {},
				Mount: func(js.Value) js.Value {
					return js.Undefined()
				},
			}
			derivedProps := &goat.Props{}

			%s := func(props goat.Props) goat.Block {
				return goat.Block{
					Patch: func(goat.Block) {},
					Mount: func(js.Value) js.Value {
						return js.Null()
					},
				}
			}

			context := &goat.Context{
				CreateBlock: %s,
				Block:       %s,
				Props:       &goat.Props{},
			}

			{{.Script}}

			%s = func(currentProps goat.Props) goat.Block {
				return goat.BlockElement(func(p goat.Props) goat.Vnode {
					return {{.Html}}
				}, currentProps)()
			}

			context.CreateBlock = %s
			block := %s(propsToPass)

			context.Block = &block
			context.Props = derivedProps

			return block
		}
	`,
		functionName,
		randomComponentName,
		randomComponentName+"generator",
		randomComponentName+"generator",
		randomComponentName,
		randomComponentName+"generator",
		randomComponentName+"generator",
		randomComponentName+"generator",
	))

	resultBuffer := strings.Builder{}
	err := template.Execute(&resultBuffer, map[string]string{
		"Script":     compilerVisitor.scriptBuffer.String(),
		"Html":       compilerVisitor.htmlBuffer.String(),
		"Components": compilerVisitor.componentBuffer.String(),
		"Imports":    importString.String(),
	})

	if err != nil {
		panic("Failed to write code")
	}

	os.WriteFile(outputFileName, []byte(resultBuffer.String()), 0644)
}

func (v *CompilerVisitor) visitHtmlElement(node *Element) {

	if node.Type == "Component" {
		v.visitComponent(node)
		return
	}

	v.htmlBuffer.WriteString(fmt.Sprintf("goat.CreateVirtualElements(%q, goat.Props{", node.Name))

	for _, attribute := range node.Attributes {
		v.visitAttribute(attribute)
	}

	for _, event := range node.Events {
		v.visitEvent(event)
	}

	v.htmlBuffer.WriteString("}")

	for _, child := range node.Children {
		v.htmlBuffer.WriteString(",")
		v.visitFragment(child)
	}

	v.htmlBuffer.WriteString(",)")
}

func (v *CompilerVisitor) visitComponent(node *Element) {
	randomComponentName := fmt.Sprintf("component_%d", rand.Int())

	v.variables = append(v.variables, randomComponentName)

	v.htmlBuffer.WriteString(fmt.Sprintf("goat.Get(%q)", randomComponentName))

	v.componentBuffer.WriteString(fmt.Sprintf(`
		%s := %s(goat.Props{
	`, randomComponentName, node.Name))

	for _, attribute := range node.Attributes {
		v.componentBuffer.WriteString(fmt.Sprintf("%q: %s,", attribute.Name, attribute.Name))
	}

	v.componentBuffer.WriteString("},)\n\n")
}

func (v *CompilerVisitor) visitTextElement(node *Text) {
	v.htmlBuffer.WriteString(fmt.Sprintf(`goat.CreateVirtualElements("text", nil,
		&goat.TextOrElement{
			StringValue: %q,
			Element:     nil,
		},
	)`, node.Value))
}

func (v *CompilerVisitor) visitExpressionElement(node *Expression) {

	v.htmlBuffer.WriteString(fmt.Sprintf(`goat.CreateVirtualElements("text", nil,
		goat.Get(%q),
	)`, node.Expression.(*ast.Ident).Name))

	v.variables = append(v.variables, node.Expression.(*ast.Ident).Name)
}

func (v *CompilerVisitor) visitFragment(node Fragment) {

	switch node := node.(type) {
	case *Element:
		{
			v.visitHtmlElement(node)
			break
		}
	case *Text:
		{
			v.visitTextElement(node)
			break
		}
	default:
		{
			v.visitExpressionElement(node.(*Expression))
			break
		}
	}
}

func (v *CompilerVisitor) visitAttribute(attr Attribute) {
	v.htmlBuffer.WriteString(fmt.Sprintf("%q:", attr.Name))
	switch e := attr.Value.(type) {
	case *ast.Ident:
		v.htmlBuffer.WriteString(fmt.Sprintf("goat.Get(%q),", e.Name))
		v.variables = append(v.variables, e.Name)
	case *ast.BasicLit:
		v.htmlBuffer.WriteString(fmt.Sprintf("%q,", e.Value))
	}
}

func (v *CompilerVisitor) visitEvent(event Event) {
	v.htmlBuffer.WriteString(fmt.Sprintf("%q: %s", event.Name, event.Name))
	v.variables = append(v.variables, event.Name)
}

func crateTemplate(name, t string) *template.Template {
	return template.Must(template.New(name).Parse(t))
}
