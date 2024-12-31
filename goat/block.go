package goat

import (
	"reflect"
	"syscall/js"
)

type Block struct {
	props   Props
	edits   []Edits
	Patch   func(Block)
	Mount   func(js.Value) js.Value
	Unmount func()
}

type ArrayBlock struct {
	elementKey string
	props      Props
	children   []Block
	Patch      func(element []string)
	Mount      func(js.Value) js.Value
}

func BlockElement(fn func(originalProp Props) Vnode, props Props) func() Block {
	// get the vnode from the current block
	vnode := fn(props)
	edits := make([]Edits, 0)
	listeners := make([]EventListener, 0)
	path := make([]int, 0)
	root := Render(vnode, &edits, path, &listeners)

	return func() Block {

		elements := make([]js.Value, len(edits))
		element := js.Undefined()
		parentOfCurrent := js.Undefined()

		mount := func(parent js.Value) js.Value {
			el := root.Call("cloneNode", true)
			parent.Call("appendChild", el)
			element = el
			parentOfCurrent = parent

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

					if edit.attribute == "value" {
						thisEl.Set("value", value)
					}

					thisEl.Call("setAttribute", edit.attribute, value)
				} else {
					edit := editUnionObject.getChildEditValue()
					value := props[edit.key]
					if reflect.TypeOf(value).String() == "goat.Block" {
						value.(Block).Mount(thisEl)
						continue
					}

					if reflect.TypeOf(value).String() == "goat.ArrayBlock" {
						value.(ArrayBlock).Mount(thisEl)
						continue
					}
					textNode := js.Global().Get("document").Call("createTextNode", value)
					thisEl.Call("insertBefore", textNode, thisEl.Get("childNodes").Index(edit.index))
				}
			}

			for _, listener := range listeners {

				thisEl := el
				for _, path := range listener.path {
					elementChild := thisEl.Get("childNodes")
					thisEl = elementChild.Index(path)
				}

				thisEl.Call("addEventListener", listener.eventType, listener.executable)
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
						continue
					}
					thisEl := elements[editIndex]
					props[edit.key] = newValue

					if edit.attribute == "value" {
						thisEl.Set("value", newValue)
					}

					thisEl.Call("setAttribute", edit.attribute, newValue)
					continue
				}

				edit := editUnionObject.getChildEditValue()
				value := props[edit.key]
				newValue := newBlock.props[edit.key]

				if reflect.TypeOf(value).String() == "goat.Block" {
					value.(Block).Patch(newBlock.props[(newBlock.edits)[editIndex].getChildEditValue().key].(Block))
					continue
				}

				if reflect.TypeOf(value).String() == "goat.ArrayBlock" {
					value.(ArrayBlock).Patch(newBlock.props[newBlock.props[(newBlock.edits)[editIndex].getChildEditValue().key].(ArrayBlock).elementKey].([]string))
					continue
				}

				if value == newValue {
					continue
				}

				props[edit.key] = newValue
				thisEl := elements[editIndex]

				thisEl.Get("childNodes").Call("item", edit.index).Set("textContent", newValue)
			}
		}

		unmount := func() {

			parentOfCurrent.Call("removeChild", element)
		}

		return Block{
			props:   props,
			edits:   edits,
			Patch:   patch,
			Mount:   mount,
			Unmount: unmount,
		}
	}
}

func ArrayBlockElement(
	elementKey string,
	rendererFunc func(element string, index int) Block,
	props Props,
) func() ArrayBlock {
	children := make([]Block, 0)
	elements := props[elementKey].([]string)

	for index, element := range elements {
		children = append(children, rendererFunc(element, index))
	}

	element := js.Null()

	arrayBlock := ArrayBlock{
		elementKey: elementKey,
		Patch:      func(elements []string) {},
		Mount:      func(js.Value) js.Value { return js.Null() },
		children:   children,
		props:      props,
	}

	return func() ArrayBlock {

		mount := func(parent js.Value) js.Value {
			element = parent
			for _, child := range children {
				child.Mount(element)
			}
			return element
		}

		patch := func(elements []string) {
			oldChildren := arrayBlock.children
			newChildren := make([]Block, 0)

			for index, element := range elements {
				newChildren = append(newChildren, rendererFunc(element, index))
			}

			maxLength := len(oldChildren)
			if len(newChildren) > maxLength {
				maxLength = len(newChildren)
			}

			for i := 0; i < maxLength; i++ {
				if i < len(oldChildren) && i < len(newChildren) {
					oldChildren[i].Unmount()
					newChildren[i].Mount(element)
				} else if i >= len(oldChildren) {
					newChildren[i].Mount(element)
				} else if i >= len(newChildren) {
					oldChildren[i].Unmount()
				}
			}

			arrayBlock.children = newChildren

		}

		arrayBlock.Patch = patch
		arrayBlock.Mount = mount

		return arrayBlock

	}
}
