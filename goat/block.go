package goat

import (
	"reflect"
	"syscall/js"
)

type Block struct {
	props Props
	edits *[]Edits
	Patch func(Block)
	Mount func(js.Value) js.Value
}

func BlockElement(fn func(proxy *Props, originalProp Props) VElement, props Props) func(prop Props) Block {

	proxy := make(Props)

	vnode := fn(&proxy, props)

	edits := make([]Edits, 0)

	path := make([]int, 0)

	root := proxy.Render(vnode, &edits, path)

	return func(prop Props) Block {

		elements := make([]js.Value, len(edits))

		mount := func(parent js.Value) js.Value {
			el := root.Call("cloneNode", true)

			// I may not completely understand this yet
			// parent.Set("textContent", "")
			parent.Call("appendChild", el)

			for editIndex, edit := range edits {
				thisEl := el
				for _, path := range edit.path {
					elementChild := thisEl.Get("childNodes")
					thisEl = elementChild.Index(path)
				}

				elements[editIndex] = thisEl

				value := props[edit.hole.(string)]

				if edit.editType == "attribute" {
					thisEl.Call("setAttribute", edit.attribute, value)
				} else if edit.editType == "child" {
					if reflect.TypeOf(value).String() == "main.Block" {
						value.(*Block).Mount(thisEl)
						continue
					}
					textNode := js.Global().Get("document").Call("createTextNode", value)
					thisEl.Call("insertBefore", textNode, thisEl.Get("childNodes").Index(edit.index))
				}
			}

			return el
		}

		patch := func(newBlock Block) {
			for editIndex, edit := range edits {
				value := props[edit.hole.(string)]
				newValue := newBlock.props[edit.hole.(string)]

				if value == newValue {
					return
				}

				thisEl := elements[editIndex]

				if edit.editType == "attribute" {
					thisEl.Call("setAttribute", edit.attribute, newValue)
				} else if edit.editType == "child" {
					if reflect.TypeOf(value).String() == "main.Block" {
						value.(*Block).Patch((*newBlock.edits)[editIndex].hole.(Block))
						continue
					}
					thisEl.Get("childNodes").Call("item", edit.index).Set("textContent", newValue)
				}
			}

		}

		return Block{
			props: prop,
			edits: &edits,
			Patch: patch,
			Mount: mount,
		}
	}

}
