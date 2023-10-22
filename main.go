package main

import (
	"syscall/js"
)

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
