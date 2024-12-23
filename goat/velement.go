package goat

type TextOrElement struct {
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
	isPlaceholder() bool
	getStringValue() string
	getVElement() *VElement
	getPlaceholderValue() *Placeholder
}

func (element *TextOrElement) isString() bool {
	return element.Element == nil
}

func (element *TextOrElement) isVElement() bool {
	return element.Element != nil
}

func (element *TextOrElement) isPlaceholder() bool {
	return false
}

func (element *TextOrElement) getStringValue() string {
	if element.isString() {
		return element.StringValue
	}
	return ""
}

func (element *TextOrElement) getVElement() *VElement {
	return element.Element
}

func (element *TextOrElement) getPlaceholderValue() *Placeholder {
	return nil
}

func CreateVirtualElements(
	elementType string,
	props Props,
	children ...Vnode,
) Vnode {
	if elementType == "text" {
		return &TextOrElement{
			Element:     nil,
			StringValue: children[0].getStringValue(),
		}
	}

	return &TextOrElement{
		Element: &VElement{
			elementType: elementType,
			props:       props,
			children:    children,
		},
		StringValue: "",
	}
}
