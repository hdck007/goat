### What the compiler must do?

Starting with a very simple example

/count.goat
```jsx
<div>
  <h1>1</h1>
</div>

```

To

/count-component.go
```go
package main

import (
	"syscall/js"

	"github.com/hdck007/goat/goat"
)

func main() {

	TEXT46 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		element := goat.CreateVirtualElements(
			"",
			"text",
			nil,
			"1",
			goat.VElement{},
		)
		return element
	}, map[string]any{})
	text609 := TEXT46(map[string]any{})

	H1919 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		childElements := []goat.VElement{text609}
		element := goat.CreateVirtualElements(
			"",
			"h1",
			nil,
			"",
			childElements...,
		)
		return element
	}, map[string]any{})
	h1268 := H1919(map[string]any{})

	DIV571 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		childElements := []goat.VElement{h1268}
		element := goat.CreateVirtualElements(
			"",
			"div",
			nil,
			"",
			childElements...,
		)
		return element
	}, map[string]any{})
	div971 := DIV571(map[string]any{})

	body := js.Global().Get("document").Call("getElementById", "root")
	div971.Mount(body)

	select {}
}

```
This is easy, just create a tree and create an equivalent code
I have updated the code for this as well in the compiler

But one caveat and my second dilemma in this compiler is that the child should be a `velement` but can be block elements, so the vdom needs to adjust this.
But these things kind of shape the dx and stuff so I will have to think about it rather than just writing code for compatibility.

Well, there were a few changes and finally got it up and running for the example. The final generated code looked like this

```go
package main

import (
	"syscall/js"

	"github.com/hdck007/goat/goat"
)

func main() {
	body := js.Global().Get("document").Call("getElementById", "root")

	DIV983 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		element := goat.CreateVirtualElements(
			"",
			"div",
			nil,
			"",
		)
		return element
	}, map[string]any{})
	div359 := DIV983(map[string]any{})
	div937 := div359.Mount(body)

	H1628 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		element := goat.CreateVirtualElements(
			"",
			"h1",
			nil,
			"",
		)
		return element
	}, map[string]any{})
	h1193 := H1628(map[string]any{})
	h1972 := h1193.Mount(div937)

	TEXT66 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		element := goat.CreateVirtualElements(
			"",
			"text",
			nil,
			"1",
		)
		return element
	}, map[string]any{})
	text445 := TEXT66(map[string]any{})
	text445.Mount(h1972)

	select {}
}
```

But heading on for now....
### How the compiler will work for the user?
I mean is the compiler a package like react or something or a CLI that takes in a .goat file and converts it into a wasm file
Well since the VDOM is different from the compiler I think the latter will be better.

So for now it will be a cli that runs like
```
goat build index.goat --config=some.json
```

And it would create a dist directory with a main.wasm file.

Adding a dump here
```go
package compiler

import (
	"fmt"
	"regexp"
	"strings"
)

type NodeType int

const (
	ElementNode NodeType = iota
	TextNode
	ExpressionNode
	ScriptNode
)

type Node struct {
	Type       NodeType
	Tag        string
	Attrs      map[string]string
	Children   []Node
	TextValue  string
	Expression string
}

// Lexer splits the template into tokens
func tokenize(template string) []string {
	// Add spaces around brackets for easier splitting
	template = regexp.MustCompile(`([<>{}])`).ReplaceAllString(template, " $1 ")
	return strings.Fields(template)
}

// Parser converts tokens into an AST
func parse(tokens []string) []Node {
	var nodes []Node
	var currentNode *Node
	var scriptContent strings.Builder
	inScript := false

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		switch {
		case token == "<script>":
			inScript = true
			continue
		case token == "</script>":
			inScript = false
			nodes = append(nodes, Node{
				Type:      ScriptNode,
				TextValue: scriptContent.String(),
			})
			continue
		case inScript:
			scriptContent.WriteString(token + " ")
			continue
		case token == "<":
			if i+1 < len(tokens) && !strings.HasPrefix(tokens[i+1], "/") {
				currentNode = &Node{
					Type:  ElementNode,
					Tag:   tokens[i+1],
					Attrs: make(map[string]string),
				}
				i++ // Skip tag name
			}
		case strings.HasPrefix(token, "class="):
			if currentNode != nil {
				currentNode.Attrs["class"] = strings.Trim(token[6:], "\"")
			}
		case token == ">":
			if currentNode != nil {
				nodes = append(nodes, *currentNode)
				currentNode = nil
			}
		case token == "{":
			if i+1 < len(tokens) && tokens[i+1] != "}" {
				nodes = append(nodes, Node{
					Type:       ExpressionNode,
					Expression: tokens[i+1],
				})
				i++ // Skip expression
			}
		case !strings.HasPrefix(token, "<") && !strings.HasPrefix(token, ">") && !strings.HasPrefix(token, "{") && !strings.HasPrefix(token, "}"):
			if currentNode == nil {
				nodes = append(nodes, Node{
					Type:      TextNode,
					TextValue: token,
				})
			}
		}
	}

	return nodes
}

// CodeGenerator converts AST to Goat component code
func generateCode(nodes []Node) string {
	var code strings.Builder
	var variables []string

	// Generate package and imports
	code.WriteString("package main\n\n")
	code.WriteString("import (\n")
	code.WriteString("\t\"syscall/js\"\n")
	code.WriteString("\t\"github.com/hdck007/goat/goat\"\n")
	code.WriteString(")\n\n")

	// Start component function
	code.WriteString("func Component(props goat.Props) goat.Block {\n")

	// Process script node first to get variables
	for _, node := range nodes {
		if node.Type == ScriptNode {
			lines := strings.Split(node.TextValue, "\n")
			for _, line := range strings.Fields(node.TextValue) {
				if strings.Contains(line, ":=") {
					varName := strings.TrimSpace(strings.Split(line, ":=")[0])
					variables = append(variables, varName)
					code.WriteString("\t" + line + "\n")
				}
			}
			// Add computed values
			code.WriteString("\tdoubledCount := count * 2\n")
		}
	}

	// Start BlockElement
	code.WriteString("\treturn goat.BlockElement(func(p goat.Props) goat.Vnode {\n")
	code.WriteString("\t\treturn ")

	// Generate element structure
	var generateNode func(node Node) string
	generateNode = func(node Node) string {
		switch node.Type {
		case ElementNode:
			result := fmt.Sprintf("goat.CreateVirtualElements(\"%s\", ", node.Tag)
			if len(node.Attrs) > 0 {
				result += "goat.Props{"
				for k, v := range node.Attrs {
					result += fmt.Sprintf("\"%s\": \"%s\", ", k, v)
				}
				result += "}, "
			} else {
				result += "nil, "
			}
			return result
		case TextNode:
			return fmt.Sprintf("goat.CreateVirtualElements(\"text\", nil, &goat.TextOrElement{StringValue: \"%s\", Element: nil})", node.TextValue)
		case ExpressionNode:
			return fmt.Sprintf("goat.Get(\"%s\")", node.Expression)
		}
		return ""
	}

	// Generate the virtual DOM structure
	for _, node := range nodes {
		if node.Type != ScriptNode {
			code.WriteString(generateNode(node))
		}
	}

	// Close parent element
	code.WriteString(")\n")

	// Add props
	code.WriteString("\t}, goat.Props{\n")
	for _, v := range variables {
		code.WriteString(fmt.Sprintf("\t\t\"%s\": %s,\n", v, v))
	}
	code.WriteString("\t})()") // Close BlockElement
	code.WriteString("\n}\n\n")

	// Add main function
	code.WriteString("func main() {\n")
	code.WriteString("\tcomponent := Component(goat.Props{})\n")
	code.WriteString("\troot := js.Global().Get(\"document\").Call(\"getElementById\", \"root\")\n")
	code.WriteString("\tcomponent.Mount(root)\n")
	code.WriteString("}")

	return code.String()
}

// CompileTemplate is the main entry point that converts a template to Goat code
func CompileTemplate(template string) string {
	tokens := tokenize(template)
	ast := parse(tokens)
	return generateCode(ast)
}
```
