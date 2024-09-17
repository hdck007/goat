package goat

import (
	"strings"
)

type UnionNode struct {
	Element     *VElement
	StringValue string
}

type VElement struct {
	elementType string
	props       Props
	children    []Vnode
}

type Vnode interface {
	isString() bool
	isVElement() bool
	isHole() bool
	getStringValue() string
	getVElement() *VElement
	getHoleValue() *Hole
}

func (element *UnionNode) isString() bool {
	trimmedString := strings.Trim(element.StringValue, " ")
	return len(trimmedString) > 0
}

func (element *UnionNode) isVElement() bool {
	return element.Element != nil
}

func (element *UnionNode) isHole() bool {
	return false
}

func (element *UnionNode) getStringValue() string {
	if element.isString() {
		return element.StringValue
	}
	return ""
}

func (element *UnionNode) getVElement() *VElement {
	return element.Element
}

func (element *UnionNode) getHoleValue() *Hole {
	return nil
}

func CreateVirtualElements(
	elementType string,
	props Props,
	children ...Vnode,
) Vnode {

	if children[0].isHole() {
		return children[0].getHoleValue()
	}

	if elementType == "text" {
		return &UnionNode{
			Element:     nil,
			StringValue: children[0].getStringValue(),
		}
	}

	return &UnionNode{
		Element: &VElement{
			elementType: elementType,
			props:       props,
			children:    children,
		},
		StringValue: "",
	}
}
