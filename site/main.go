package main

import (
	"syscall/js"

	"github.com/hdck007/goat/goat"
)

func main() {

	name := "John"

	body := js.Global().Get("document").Call("getElementById", "root")

	DIV544 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		element := goat.CreateVirtualElements(
			"",
			"div",
			nil,
			"",
		)
		return element
	}, map[string]any{})
	div379 := DIV544(map[string]any{})
	div368 := div379.Mount(body)

	H1758 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		element := goat.CreateVirtualElements(
			"",
			"h1",
			map[string]any{
				"class": "text-red",
			},
			"",
		)
		return element
	}, map[string]any{})
	h1917 := H1758(map[string]any{})
	h1275 := h1917.Mount(div368)

	TEXT71 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		element := goat.CreateVirtualElements(
			"",
			"text",
			nil,
			"    Hello ",
		)
		return element
	}, map[string]any{})
	text191 := TEXT71(map[string]any{})
	text191.Mount(h1275)

	TEXT688 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		element := goat.CreateVirtualElements(
			"",
			"text",
			nil,
			prop[proxy.Get("name").Key],
		)
		return element
	}, map[string]any{"name": name})
	text116 := TEXT688(map[string]any{})
	text116.Mount(h1275)

	TEXT201 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		element := goat.CreateVirtualElements(
			"",
			"text",
			nil,
			", How are you?    ",
		)
		return element
	}, map[string]any{})
	text305 := TEXT201(map[string]any{})
	text305.Mount(h1275)

	select {}
}
