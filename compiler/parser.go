package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"strings"
)

type Script struct {
	Value string
}

type Imports struct {
	Imports []Import
}

type Import struct {
	Type  string
	File  string
	Name  string
	Value interface{}
}

type AST struct {
	HTML    []Fragment
	Script  *Script
	Imports *Imports
}

type Fragment interface{}

type Element struct {
	Type       string
	Name       string
	Events     []Event
	Attributes []Attribute
	Children   []Fragment
}

type Attribute struct {
	Type  string
	Name  string
	Value interface{} // Will hold parsed Go expression
}

type Event struct {
	Type  string
	Name  string
	Value interface{}
}

type ExpressionValue struct {
	Type  string
	Value interface{}
}

type Expression struct {
	Type       string
	Expression interface{}
}

type Text struct {
	Type  string
	Value string
}

type Parser struct {
	content string
	pos     int
	fset    *token.FileSet
}

func Parse(content string) (*AST, error) {
	p := &Parser{
		content: content,
		pos:     0,
		fset:    token.NewFileSet(),
	}

	ast := &AST{}
	fragments, script, imports, err := p.parseFragments(func() bool {
		return p.pos < len(p.content)
	})
	if err != nil {
		return nil, err
	}

	ast.HTML = fragments
	ast.Script = script
	ast.Imports = imports
	return ast, nil
}

func ParseFile(fileName string) (*AST, error) {
	b, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Print(err)
	}

	ast, error := Parse(string(b))

	if error != nil {
		fmt.Println("Error occurred while parsing: ", error)
	}

	return ast, nil
}

func (p *Parser) parseFragments(condition func() bool) ([]Fragment, *Script, *Imports, error) {
	var fragments []Fragment
	var script Script
	var imports Imports

	for condition() {
		fragment, err := p.parseFragment()
		if err != nil {
			return nil, nil, nil, err
		}
		if fragment != nil {

			switch fragmentValue := fragment.(type) {
			case Script:
				{
					script = Script{
						Value: fragmentValue.Value,
					}
					break
				}
			case *Element:
				{
					fragments = append(fragments, fragment)

					break
				}
			case Imports:
				{
					imports = fragmentValue
					break
				}
			default:
				{
					fragments = append(fragments, fragment)
					break
				}
			}

		}
	}

	return fragments, &script, &imports, nil
}

func (p *Parser) parseFragment() (Fragment, error) {
	if imports, err := p.parseImports(); err != nil {
		return nil, err
	} else if imports != nil {
		return imports, nil
	}
	if script, err := p.parseScript(); err != nil {
		return nil, err
	} else if script != nil {
		return script, nil
	}
	if element, err := p.parseElement(); err != nil {
		return nil, err
	} else if element != nil {
		return element, nil
	}
	if expr, err := p.parseExpression(); err != nil {
		return nil, err
	} else if expr != nil {
		return expr, nil
	}
	if text := p.parseText(); text != nil {
		return text, nil
	}
	return nil, nil
}

func (p *Parser) parseImports() (Fragment, error) {
	imports := []Import{}
	if p.match("<imports>") {
		p.eat("<imports>")
		startIndex := p.pos
		endIndex := strings.Index(p.content[p.pos:], "</imports>")
		if endIndex == -1 {
			return nil, fmt.Errorf("unclosed imports tag")
		}
		endIndex += p.pos
		code := p.content[startIndex:endIndex]
		p.pos = endIndex
		p.eat("</imports>")

		re := regexp.MustCompile(`import\s+(\w+)\s+from\s+"(.+)"`)

		matches := re.FindAllStringSubmatch(code, -1)

		for _, match := range matches {
			if len(match) == 3 {
				componentName := match[1]
				fileName := match[2]

				imports = append(imports, Import{
					File:  fileName,
					Name:  componentName,
					Type:  "import",
					Value: componentName,
				})
			}
		}

		goalngImportsRegex := regexp.MustCompile(`"([^"]+)"`)

		golangImportMatches := goalngImportsRegex.FindAllStringSubmatch(code, -1)

		for _, match := range golangImportMatches {
			if len(match) > 1 {
				importPath := match[1]
				if !strings.HasSuffix(importPath, ".goat") {
					imports = append(imports, Import{
						File:  "",
						Type:  "golangImport",
						Name:  importPath,
						Value: importPath,
					})
				}
			}
		}

		return Imports{
			Imports: imports,
		}, nil
	}
	return nil, nil
}

func (p *Parser) parseScript() (Fragment, error) {
	if p.match("<script>") {
		p.eat("<script>")
		startIndex := p.pos
		endIndex := strings.Index(p.content[p.pos:], "</script>")
		if endIndex == -1 {
			return nil, fmt.Errorf("unclosed script tag")
		}
		endIndex += p.pos
		code := p.content[startIndex:endIndex]
		p.pos = endIndex
		p.eat("</script>")
		return Script{
			Value: code,
		}, nil
	}
	return nil, nil
}

func (p *Parser) parseComponent() (*Element, error) {
	componentName := p.readWhileMatching("[A-Za-z0-9]")
	attributes, _, err := p.parseAttributeList()
	if err != nil {
		return nil, err
	}

	if p.match("/>") {
		p.eat("/>")
		return &Element{
			Type:       "Component",
			Name:       componentName,
			Events:     nil,
			Attributes: attributes,
			Children:   nil,
		}, nil
	}
	p.eat(">")
	endTag := fmt.Sprintf("</%s>", componentName)

	children, _, _, err := p.parseFragments(func() bool {
		return !p.match(endTag)
	})
	if err != nil {
		return nil, err
	}

	p.eat(endTag)

	return &Element{
		Type:       "Component",
		Name:       componentName,
		Attributes: attributes,
		Children:   children,
	}, nil
}

func (p *Parser) parseElement() (*Element, error) {
	if p.match("<") {
		p.eat("<")

		if p.startsWith("[A-Z]") {
			return p.parseComponent()
		}

		tagName := p.readWhileMatching("[a-z0-9]")
		attributes, events, err := p.parseAttributeList()
		if err != nil {
			return nil, err
		}

		if p.match("/>") {
			p.eat("/>")
			return &Element{
				Type:       "Element",
				Name:       tagName,
				Events:     events,
				Attributes: attributes,
				Children:   nil,
			}, nil
		}
		p.eat(">")
		endTag := fmt.Sprintf("</%s>", tagName)

		children, _, _, err := p.parseFragments(func() bool {
			return !p.match(endTag)
		})
		if err != nil {
			return nil, err
		}

		p.eat(endTag)

		return &Element{
			Type:       "Element",
			Name:       tagName,
			Events:     events,
			Attributes: attributes,
			Children:   children,
		}, nil
	}
	return nil, nil
}

func (p *Parser) parseAttributeList() ([]Attribute, []Event, error) {
	var attributes []Attribute
	var events []Event
	p.skipWhitespace()

	for !(p.match(">") || p.match("/>")) {

		if p.startsWith("@") {
			p.eat("@")
			event, err := p.parseEvent()
			if err != nil {
				return nil, nil, err
			}
			events = append(events, *event)
			p.skipWhitespace()
			continue
		}

		attr, err := p.parseAttribute()
		if err != nil {
			return nil, nil, err
		}
		attributes = append(attributes, *attr)
		p.skipWhitespace()
	}

	return attributes, events, nil
}

func (p *Parser) parseAttribute() (*Attribute, error) {
	name := p.readWhileMatching("[^=]")
	p.eat(`={`)

	exprStr := p.readWhileMatching("[^}]")
	expr, err := parser.ParseExpr(exprStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Go expression: %v", err)
	}

	p.eat(`}`)

	return &Attribute{
		Type:  "Attribute",
		Name:  name,
		Value: expr,
	}, nil
}

func (p *Parser) parseEvent() (*Event, error) {
	eventType := p.readWhileMatching("[^=]")
	p.eat(`={`)
	functionName := p.readWhileMatching("[^}]")
	p.eat(`}`)

	return &Event{
		Type:  eventType,
		Name:  functionName,
		Value: nil,
	}, nil
}

func (p *Parser) parseExpression() (*Expression, error) {
	if p.match("{") {
		p.eat("{")

		exprStr := p.readWhileMatching("[^}]")
		expr, err := parser.ParseExpr(exprStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Go expression: %v", err)
		}

		p.eat("}")

		return &Expression{
			Type:       "Expression",
			Expression: expr,
		}, nil
	}
	return nil, nil
}

func (p *Parser) parseText() *Text {
	text := p.readWhileMatching("[^<{]")
	if strings.TrimSpace(text) != "" {
		return &Text{
			Type:  "Text",
			Value: text,
		}
	}
	return nil
}

func (p *Parser) match(str string) bool {
	if p.pos+len(str) > len(p.content) {
		return false
	}
	return p.content[p.pos:p.pos+len(str)] == str
}

func (p *Parser) eat(str string) {
	if p.match(str) {
		p.pos += len(str)
	} else {
		panic(fmt.Sprintf("Parse error: expecting %q", str))
	}
}

func (p *Parser) readWhileMatching(pattern string) string {
	re := regexp.MustCompile(pattern)
	startPos := p.pos
	for p.pos < len(p.content) {
		if !re.MatchString(string(p.content[p.pos])) {
			break
		}
		p.pos++
	}
	return p.content[startPos:p.pos]
}

func (p *Parser) skipWhitespace() {
	p.readWhileMatching(`[\s\n]`)
}

func (p *Parser) startsWith(pattern string) bool {
	re := regexp.MustCompile(pattern)
	return re.MatchString(string(p.content[p.pos]))
}
