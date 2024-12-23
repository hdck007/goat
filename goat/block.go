package goat

import (
	"reflect"
	"syscall/js"
)

type Block struct {
	props Props
	edits []Edits
	Patch func(Block)
	Mount func(js.Value) js.Value
}

func BlockElement(fn func(originalProp Props) Vnode, props Props) func() Block {
	// get the vnode from the current block
	vnode := fn(props)
	edits := make([]Edits, 0)
	path := make([]int, 0)
	root := Render(vnode, &edits, path)

	return func() Block {

		elements := make([]js.Value, len(edits))

		mount := func(parent js.Value) js.Value {
			el := root.Call("cloneNode", true)
			parent.Call("appendChild", el)

			for editIndex, editUnionObject := range edits {
				thisEl := el

				if editUnionObject.isAttributeEdit() {
					for _, path := range editUnionObject.getAttributeEditValue().path {
						elementChild := thisEl.Get("childNodes")
						thisEl = elementChild.Index(path)
					}
				} else {
					for _, path := range editUnionObject.getChildEditValue().path {
						elementChild := thisEl.Get("childNodes")
						thisEl = elementChild.Index(path)
					}
				}

				elements[editIndex] = thisEl

				if editUnionObject.isAttributeEdit() {
					edit := editUnionObject.getAttributeEditValue()
					value := props[edit.key]
					thisEl.Call("setAttribute", edit.attribute, value)
				} else {
					edit := editUnionObject.getChildEditValue()
					value := props[edit.key]
					if reflect.TypeOf(value).String() == "goat.Block" {
						value.(Block).Mount(thisEl)
						continue
					}
					textNode := js.Global().Get("document").Call("createTextNode", value)
					thisEl.Call("insertBefore", textNode, thisEl.Get("childNodes").Index(edit.index))
				}
			}

			return el
		}

		patch := func(newBlock Block) {
			for editIndex, editUnionObject := range edits {
				if editUnionObject.isAttributeEdit() {
					edit := editUnionObject.getAttributeEditValue()
					value := props[edit.key]
					newValue := newBlock.props[edit.key]
					if value == newValue {
						return
					}
					thisEl := elements[editIndex]
					thisEl.Call("setAttribute", edit.attribute, newValue)
					return
				}

				edit := editUnionObject.getChildEditValue()
				value := props[edit.key]
				newValue := newBlock.props[edit.key]

				if reflect.TypeOf(value).String() == "goat.Block" {
					value.(Block).Patch(newBlock.props[(newBlock.edits)[editIndex].getChildEditValue().key].(Block))
					continue
				}

				if value == newValue {
					return
				}

				props[edit.key] = newValue
				thisEl := elements[editIndex]

				thisEl.Get("childNodes").Call("item", edit.index).Set("textContent", newValue)
			}
		}

		return Block{
			props: props,
			edits: edits,
			Patch: patch,
			Mount: mount,
		}
	}
}
