package goat

import (
	"reflect"
	"syscall/js"
)

func (proxy *Props) Render(
	element Vnode,
	edits *[]Edits,
	path []int,
) js.Value {

	document := js.Global().Get("document")

	if element.isString() {
		return document.Call("createTextNode", element.getStringValue())
	}

	isVElement := element.isVElement()

	if !isVElement {
		return js.Value{}
	}

	virtualElement := element.getVElement()

	el := document.Call("createElement", virtualElement.elementType)

	if virtualElement.props != nil {
		for k, v := range virtualElement.props {
			if reflect.TypeOf(v).String() == "goat.Hole" {
				*edits = append(*edits, &EditAttribute{
					editType:  "attribute",
					path:      path,
					attribute: k,
					hole:      v.(*Hole).key,
				})
				continue
			}

			el.Call("setAttribute", k, v)
		}
	}

	if virtualElement.children == nil {
		return el
	}

	for childIndex, child := range virtualElement.children {
		if child.isHole() {
			*edits = append(*edits, &EditChild{
				editType: "child",
				path:     path,
				index:    childIndex,
				hole:     child.getHoleValue().key,
			})
			continue
		}
		newPath := make([]int, len(path))
		copy(newPath, path)
		newPath = append(newPath, childIndex)
		theChild := proxy.Render(child, edits, newPath)
		el.Call("appendChild", theChild)
	}

	return el
}
