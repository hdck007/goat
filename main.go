package main

import (
	"reflect"
	"syscall/js"
)

type Props map[string]any

type VElement struct {
	elementType string
	props       Props
	children    []VElement
	text        any
	key         string
}

type Edits struct {
	editType    string
	path        []int
	attribute   string
	index       int
	hole        any
	elementType string
}

type Block struct {
	props Props
	edits *[]Edits
	patch func(Block)
	mount func(js.Value)
}

type ReturnType struct {
	edits []Edits
	props Props
	mount func(parent js.Value)
	patch func(block Block)
}

func createVirtualElements(
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

func (proxy *Props) render(
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
					hole:      v.(Hole).key,
				})
				continue
			}
			el.Call("setAttribute", k, v)
		}
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
		el.Call("appendChild", proxy.render(child, edits, newPath))
	}

	return el
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

func main() {
	value := 0

	Div := blockElement(func(proxy *Props, prop Props) VElement {
		childElements := []VElement{
			createVirtualElements(
				"",
				"text",
				nil,
				"The current value is ",
				VElement{},
			),
			createVirtualElements(
				"number",
				"text",
				nil,
				prop[proxy.Get("number").key],
				VElement{},
			),
		}

		element := createVirtualElements(
			"root",
			"div",
			nil,
			"",
			childElements...,
		)
		return element
	}, map[string]any{
		"number": value,
	})

	div := Div(map[string]any{
		"number": 1,
	})

	body := js.Global().Get("document").Call("getElementById", "root")
	var cb js.Func
	cb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		value++
		div.patch(Div(map[string]any{
			"number": value,
		}))
		return nil
	})
	js.Global().Get("document").Call("getElementById", "increment").Call("addEventListener", "click", cb)
	div.mount(body)

	select {}
}
