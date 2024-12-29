package main

import (
	"fmt"
	"go/ast"
	"math/rand"
	"strings"
	"text/template"
)

type CompilerVisitor struct {
	scriptBuffer strings.Builder
	htmlBuffer   strings.Builder
	variables    []string
}

func Generate(ast AST) string {

	compilerVisitor := CompilerVisitor{
		variables:    []string{},
		htmlBuffer:   strings.Builder{},
		scriptBuffer: strings.Builder{},
	}

	compilerVisitor.visitFragment(ast.HTML[0])

	if ast.Script != nil {
		compilerVisitor.scriptBuffer.WriteString(ast.Script.Value)
	}

	compilerVisitor.scriptBuffer.WriteString("derivedProps := goat.Props {")

	for _, variable := range compilerVisitor.variables {
		compilerVisitor.scriptBuffer.WriteString(fmt.Sprintf("%q: %s,", variable, variable))
	}

	compilerVisitor.scriptBuffer.WriteString("}")

	compilerVisitor.scriptBuffer.WriteString(`
		for k, v := range props {
			derivedProps[k] = v
		}
	`)

	randomComponentName := fmt.Sprintf("componentName%d", rand.Int())

	template := crateTemplate("component-template", fmt.Sprintf(`
		func App(props goat.Props) goat.Block {
			%s := goat.Block{
				Patch: func(goat.Block) {},
				Mount: func(js.Value) {},
			}

			{{.Script}}

			%s = goat.BlockElement(func(p goat.Props) goat.Vnode {
				return {{.Html}}
			}, derivedProps)()

			return %s
		}
	`, randomComponentName, randomComponentName, randomComponentName))

	resultBuffer := strings.Builder{}
	err := template.Execute(&resultBuffer, map[string]string{
		"Script": compilerVisitor.scriptBuffer.String(),
		"Html":   compilerVisitor.htmlBuffer.String(),
	})

	if err != nil {
		panic("Failed to write code")
	}

	return resultBuffer.String()
}

func (v *CompilerVisitor) visitHtmlElement(node *Element) {
	v.htmlBuffer.WriteString(fmt.Sprintf("goat.CreateVirtualElements(%q, goat.Props{", node.Name))

	for _, attribute := range node.Attributes {
		v.visitAttribute(attribute)
	}
	v.htmlBuffer.WriteString("}")

	for _, child := range node.Children {
		v.htmlBuffer.WriteString(",")
		v.visitFragment(child)
	}

	v.htmlBuffer.WriteString(",)")

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
		v.htmlBuffer.WriteString(fmt.Sprintf("goat.Get(%q)", e.Name))
		v.variables = append(v.variables, e.Name)
	case *ast.BasicLit:
		v.htmlBuffer.WriteString(fmt.Sprintf("%q,", e.Value))
	}
}

func crateTemplate(name, t string) *template.Template {
	return template.Must(template.New(name).Parse(t))
}
