package goat

import (
	"reflect"
	"syscall/js"
)

func Render(
	element Vnode,
	edits *[]Edits,
	path []int,
	listeners *[]EventListener,
) js.Value {

	document := js.Global().Get("document")

	if element.isString() {
		return document.Call("createTextNode", element.getStringValue())
	}

	virtualElement := element.getVElement()

	el := document.Call("createElement", virtualElement.elementType)

	if virtualElement.props != nil {
		for k, v := range virtualElement.props {

			if k[0] == '@' {
				*listeners = append(*listeners, EventListener{
					eventType:  k[1:],
					executable: v,
					path:       path,
				})
				continue
			}

			if reflect.TypeOf(v).String() == "*goat.Placeholder" {
				*edits = append(*edits, &EditAttribute{
					editType:  "attribute",
					path:      path,
					attribute: k,
					key:       v.(*Placeholder).key,
				})
				continue
			}

			if k == "value" {
				el.Set("value", v)
			}

			el.Call("setAttribute", k, v)
		}
	}

	if virtualElement.children == nil {
		return el
	}

	for childIndex, child := range virtualElement.children {
		if child.isPlaceholder() {
			*edits = append(*edits, &EditChild{
				editType: "child",
				path:     path,
				index:    childIndex,
				key:      child.getPlaceholderValue().key,
			})
			continue
		}
		newPath := make([]int, len(path))
		copy(newPath, path)
		newPath = append(newPath, childIndex)
		theChild := Render(child, edits, newPath, listeners)
		el.Call("appendChild", theChild)
	}

	return el
}
