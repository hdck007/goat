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

func BlockElement(fn func(proxy *Props, originalProp Props) Vnode, props Props) func(prop Props) Block {

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
			parent.Set("textContent", "")
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
					value := props[edit.hole]
					thisEl.Call("setAttribute", edit.attribute, value)
				} else {

					// edit for composition
					// what if the child is a block itself

					edit := editUnionObject.getChildEditValue()
					value := props[edit.hole]

					// if reflect.TypeOf(value).String() == "goat.Block" {
					// 	value.(*Block).Mount(thisEl)
					// 	continue
					// }
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
					value := props[edit.hole]
					newValue := newBlock.props[edit.hole]

					if value == newValue {
						return
					}

					thisEl := elements[editIndex]

					thisEl.Call("setAttribute", edit.attribute, newValue)

					return
				}

				edit := editUnionObject.getChildEditValue()

				value := props[edit.hole]
				newValue := newBlock.props[edit.hole]

				if value == newValue {
					return
				}

				// update the existing edit value
				props[edit.hole] = newValue

				thisEl := elements[editIndex]

				if reflect.TypeOf(value).String() == "goat.Block" {
					// edit for composition
					// value.(*Block).Patch((*newBlock.edits)[editIndex].getChildEditValue().hole.(Block))
					continue
				}
				thisEl.Get("childNodes").Call("item", edit.index).Set("textContent", newValue)

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
