package main

import (
	"syscall/js"

	"github.com/hdck007/goat/goat"
)

func start() {

	count := "Hello"

	body := js.Global().Get("document").Call("getElementById", "root")

	DIV52 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		element := goat.CreateVirtualElements(
			"",
			"div",
			nil,
			"",
		)
		return element
	}, map[string]any{})
	div233 := DIV52(map[string]any{})
	div433 := div233.Mount(body)

	H1589 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
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
	h1844 := H1589(map[string]any{})
	h1308 := h1844.Mount(div433)

	TEXT587 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		element := goat.CreateVirtualElements(
			"",
			"text",
			nil,
			"    ",
		)
		return element
	}, map[string]any{})
	text206 := TEXT587(map[string]any{})
	text206.Mount(h1308)

	TEXT669 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		element := goat.CreateVirtualElements(
			"",
			"text",
			nil,
			prop[proxy.Get("count").Key],
		)
		return element
	}, map[string]any{"count": count})
	text23 := TEXT669(map[string]any{})
	text23.Mount(h1308)

	TEXT532 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		element := goat.CreateVirtualElements(
			"",
			"text",
			nil,
			" World!    ",
		)
		return element
	}, map[string]any{})
	text561 := TEXT532(map[string]any{})
	text561.Mount(h1308)

	select {}
}
