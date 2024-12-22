package main

import (
	"fmt"
	"os"
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
			// lines := strings.Split(node.TextValue, "\n")
			for _, line := range strings.Fields(node.TextValue) {
				if strings.Contains(line, ":=") {
					varName := strings.TrimSpace(strings.Split(line, ":=")[0])
					variables = append(variables, varName)
					code.WriteString("\t" + line + "\n")
				}
			}
			// Add computed values
			// code.WriteString("\tdoubledCount := count * 2\n")
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

func CompileTemplate(template string) string {
	tokens := tokenize(template)
	ast := parse(tokens)
	return generateCode(ast)
}

func main() {
	b, err := os.ReadFile("example.goat") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	str := string(b) // convert content to a 'string'

	compiledCode := CompileTemplate(str)

	os.WriteFile("example.go", []byte(compiledCode), 0644)
}
