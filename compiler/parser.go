package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"regexp"
	"strings"
)

type AST struct {
	HTML   []Fragment
	Script interface{} // Will hold parsed Go code
}

type Fragment interface{}

type Element struct {
	Type       string
	Name       string
	Attributes []Attribute
	Children   []Fragment
}

type Attribute struct {
	Type  string
	Name  string
	Value interface{} // Will hold parsed Go expression
}

type Expression struct {
	Type       string
	Expression interface{} // Will hold parsed Go expression
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
	fragments, err := p.parseFragments(func() bool {
		return p.pos < len(p.content)
	})
	if err != nil {
		return nil, err
	}

	ast.HTML = fragments
	return ast, nil
}

func (p *Parser) parseFragments(condition func() bool) ([]Fragment, error) {
	var fragments []Fragment

	for condition() {
		fragment, err := p.parseFragment()
		if err != nil {
			return nil, err
		}
		if fragment != nil {
			fragments = append(fragments, fragment)
		}
	}

	return fragments, nil
}

func (p *Parser) parseFragment() (Fragment, error) {
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
	// if expr, err := p.parseExpression(); err != nil {
	// 	return nil, err
	// } else if expr != nil {
	// 	return expr, nil
	// }
	if text := p.parseText(); text != nil {
		return text, nil
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

		// Parse Go code using go/parser
		expr, err := parser.ParseFile(p.fset, "", "package main\n"+code, parser.AllErrors)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Go code: %v", err)
		}

		p.pos = endIndex
		p.eat("</script>")
		return expr, nil
	}
	return nil, nil
}

func (p *Parser) parseElement() (*Element, error) {
	if p.match("<") {
		p.eat("<")
		tagName := p.readWhileMatching("[a-z0-9]")
		attributes, err := p.parseAttributeList()
		if err != nil {
			return nil, err
		}

		p.eat(">")
		endTag := fmt.Sprintf("</%s>", tagName)

		children, err := p.parseFragments(func() bool {
			return !p.match(endTag)
		})
		if err != nil {
			return nil, err
		}

		p.eat(endTag)

		return &Element{
			Type:       "Element",
			Name:       tagName,
			Attributes: attributes,
			Children:   children,
		}, nil
	}
	return nil, nil
}

func (p *Parser) parseAttributeList() ([]Attribute, error) {
	var attributes []Attribute
	p.skipWhitespace()

	for !p.match(">") {
		attr, err := p.parseAttribute()
		if err != nil {
			return nil, err
		}
		attributes = append(attributes, *attr)
		p.skipWhitespace()
	}

	return attributes, nil
}

func (p *Parser) parseAttribute() (*Attribute, error) {
	name := p.readWhileMatching("[^=]")
	fmt.Print(name)
	p.eat(`={`)

	// Parse Go expression
	exprStr := p.readWhileMatching("[^}]")
	fmt.Println(exprStr)
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

func (p *Parser) parseExpression() (*Expression, error) {
	if p.match("{") {
		p.eat("{")

		// Parse Go expression
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
		fmt.Println(string(p.content[p.pos]), pattern)
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
