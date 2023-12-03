package main

import (
	"syscall/js"

	"github.com/hdck007/goat/goat"
)

func main() {
	value := 0

	Div := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		childElements := []goat.VElement{
			goat.CreateVirtualElements(
				"",
				"text",
				nil,
				"The current value is ",
				goat.VElement{},
			),
			goat.CreateVirtualElements(
				"number",
				"text",
				nil,
				prop[proxy.Get("number").Key],
				goat.VElement{},
			),
		}

		element := goat.CreateVirtualElements(
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
		div.Patch(Div(map[string]any{
			"number": value,
		}))
		return nil
	})
	js.Global().Get("document").Call("getElementById", "increment").Call("addEventListener", "click", cb)
	div.Mount(body)

	select {}
}
