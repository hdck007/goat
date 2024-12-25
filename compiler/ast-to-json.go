package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
)

// convertFragmentsToJSON converts a slice of Fragments to JSON-friendly format
func converAttributesToJSON(attributes []Attribute) []interface{} {
	result := make([]interface{}, 0, len(attributes))

	for _, attr := range attributes {
		result = append(result, map[string]interface{}{
			"type":  attr.Type,
			"name":  attr.Name,
			"value": convertGoExprToJSON(attr.Value),
		})
	}

	return result
}

func convertFragmentsToJSON(fragments []Fragment) []interface{} {
	result := make([]interface{}, 0, len(fragments))

	for _, fragment := range fragments {
		switch f := fragment.(type) {
		case *Element:
			result = append(result, map[string]interface{}{
				"type":       f.Type,
				"name":       f.Name,
				"attributes": converAttributesToJSON(f.Attributes),
				"events":     f.Events,
				"children":   convertFragmentsToJSON(f.Children),
			})
		case *Text:
			result = append(result, map[string]interface{}{
				"type":  "Text",
				"value": f.Value,
			})
		case *Expression:
			result = append(result, map[string]interface{}{
				"type":       "Expression",
				"expression": convertGoExprToJSON(f.Expression),
			})
		}
	}

	return result
}

// convertGoExprToJSON converts Go AST expressions to JSON-friendly format
func convertGoExprToJSON(expr interface{}) interface{} {
	switch e := expr.(type) {
	case *ast.Ident:
		return map[string]interface{}{
			"type":  "Identifier",
			"value": e.Name,
		}
	case *ast.BasicLit:
		return map[string]interface{}{
			"type":  "Literal",
			"value": e.Value,
		}
	case *ast.BinaryExpr:
		return map[string]interface{}{
			"type":     "BinaryExpression",
			"operator": e.Op.String(),
			"left":     convertGoExprToJSON(e.X),
			"right":    convertGoExprToJSON(e.Y),
		}
	case *ast.CallExpr:
		return map[string]interface{}{
			"type":     "CallExpression",
			"function": convertGoExprToJSON(e.Fun),
			"args":     convertExprSliceToJSON(e.Args),
		}
	case *ast.SelectorExpr:
		return map[string]interface{}{
			"type":     "MemberExpression",
			"object":   convertGoExprToJSON(e.X),
			"property": convertGoExprToJSON(e.Sel),
		}
	default:
		return map[string]interface{}{
			"type": fmt.Sprintf("%T", expr),
		}
	}
}

// convertExprSliceToJSON converts a slice of Go AST expressions to JSON-friendly format
func convertExprSliceToJSON(exprs []ast.Expr) []interface{} {
	result := make([]interface{}, 0, len(exprs))
	for _, expr := range exprs {
		result = append(result, convertGoExprToJSON(expr))
	}
	return result
}

// ToJSON converts the AST to a JSON string
func (a *AST) ToJSON() (string, error) {
	astMap := map[string]interface{}{
		"html": convertFragmentsToJSON(a.HTML),
	}

	if a.Script != nil {
		if file, ok := a.Script.(*ast.File); ok {
			astMap["script"] = map[string]interface{}{
				"type": "Program",
				"body": convertGoASTToJSON(file),
			}
		}
	}

	jsonBytes, err := json.MarshalIndent(astMap, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal AST to JSON: %v", err)
	}

	return string(jsonBytes), nil
}

// convertGoASTToJSON converts Go AST program to JSON-friendly format
func convertGoASTToJSON(file *ast.File) []interface{} {
	var statements []interface{}

	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.ValueSpec:
					for i, name := range s.Names {
						var value interface{}
						if i < len(s.Values) {
							value = convertGoExprToJSON(s.Values[i])
						}
						statements = append(statements, map[string]interface{}{
							"type":         "VariableDeclaration",
							"name":         name.Name,
							"initialValue": value,
						})
					}
				}
			}
		case *ast.FuncDecl:
			statements = append(statements, map[string]interface{}{
				"type": "FunctionDeclaration",
				"name": d.Name.Name,
				"body": convertGoExprToJSON(d.Body),
			})
		}
	}

	return statements
}
