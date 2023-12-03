package goat

type VElement struct {
	elementType string
	props       Props
	children    []VElement
	text        any
	key         string
}

func CreateVirtualElements(
	key string,
	elementType string,
	props Props,
	text any,
	children ...VElement,
) VElement {
	if elementType == "text" {
		return VElement{
			elementType: "text",
			props:       nil,
			children:    nil,
			text:        text,
			key:         key,
		}
	}
	return VElement{
		elementType: elementType,
		props:       props,
		children:    children,
		text:        "",
		key:         key,
	}
}
