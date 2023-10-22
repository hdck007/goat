package main

import (
	"reflect"
	"syscall/js"
)

type Block struct {
	props Props
	edits *[]Edits
	patch func(Block)
	mount func(js.Value)
}

func blockElement(fn func(proxy *Props, originalProp Props) VElement, props Props) func(prop Props) Block {

	proxy := make(Props)

	vnode := fn(&proxy, props)

	edits := make([]Edits, 0)

	path := make([]int, 0)

	root := proxy.render(vnode, &edits, path)

	return func(prop Props) Block {

		elements := make([]js.Value, len(edits))

		mount := func(parent js.Value) {
			el := root.Call("cloneNode", true)

			parent.Set("textContent", "")
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
						value.(*Block).mount(thisEl)
						continue
					}
					textNode := js.Global().Get("document").Call("createTextNode", value)
					thisEl.Call("insertBefore", textNode, thisEl.Get("childNodes").Index(edit.index))
				}
			}
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
						value.(*Block).patch((*newBlock.edits)[editIndex].hole.(Block))
						continue
					}
					thisEl.Get("childNodes").Call("item", edit.index).Set("textContent", newValue)
				}
			}

		}

		return Block{
			props: prop,
			edits: &edits,
			patch: patch,
			mount: mount,
		}
	}

}
