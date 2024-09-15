package main

import (
	"syscall/js"

	"github.com/hdck007/goat/goat"
)

func main() {

	name := "Hello"

	updateName := func(value string) interface{} {
		name = value

		return nil
	}

	body := js.Global().Get("document").Call("getElementById", "root")

	h1712 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.Vnode {
		children := []goat.Vnode{
			&goat.UnionNode{
				Element:     nil,
				StringValue: "Click me to change ",
			},
			proxy.Get("name"),
		}

		element := goat.CreateVirtualElements(
			"h1",
			map[string]any{
				"class": "bg-red-500 cursor-pointer",
			},
			children...,
		)
		return element
	}, map[string]any{
		"name": name,
	})

	h1238 := h1712(map[string]any{
		"name": name,
	})

	h1736 := h1238.Mount(body)

	var updatename628 js.Func
	updatename628 = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if name == "Hello" {
			updateName("Mello")
			h1238.Patch(h1712(map[string]any{
				"name": name,
			}))
		} else {
			updateName("Hello")
			h1238.Patch(h1712(map[string]any{
				"name": name,
			}))
		}
		return nil
	})
	h1736.Call("addEventListener", "click", updatename628)

	select {}
}
