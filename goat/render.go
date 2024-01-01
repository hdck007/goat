package goat

import "syscall/js"

func (proxy *Props) Render(
	element VElement,
	edits *[]Edits,
	path []int,
) js.Value {

	document := js.Global().Get("document")

	if element.elementType == "text" {
		node := document.Call("createTextNode", element.text)
		return node
	}

	el := document.Call("createElement", element.elementType)

	if element.props != nil {
		for k, v := range element.props {
			_, ok := (*proxy)[k]
			if ok {
				*edits = append(*edits, Edits{
					editType:  "attribute",
					path:      path,
					attribute: k,
					index:     -1,
					hole:      v.(Hole).Key,
				})
				continue
			}
			el.Call("setAttribute", k, v)
		}
	}

	if element.children == nil {
		return el
	}

	for childIndex, child := range element.children {
		_, ok := (*proxy)[child.key]
		if ok {
			*edits = append(*edits, Edits{
				editType:    "child",
				path:        path,
				index:       childIndex,
				hole:        child.key,
				elementType: child.elementType,
			})
			continue
		}
		newPath := make([]int, len(path))
		copy(newPath, path)
		newPath = append(newPath, childIndex)
		el.Call("appendChild", proxy.Render(child, edits, newPath))
	}

	return el
}
