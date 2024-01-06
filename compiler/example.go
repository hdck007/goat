package main

import (
	"syscall/js"

	"github.com/hdck007/goat/goat"
)

func start() {

	name := "Hello"

	updateName := func() interface{} {
		name = "Mello"

		return nil
	}

	body := js.Global().Get("document").Call("getElementById", "root")

	DIV233 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		element := goat.CreateVirtualElements(
			"",
			"div",
			nil,
			"",
		)
		return element
	}, map[string]any{})
	div775 := DIV233(map[string]any{})
	div860 := div775.Mount(body)

	BUTTON661 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		element := goat.CreateVirtualElements(
			"",
			"button",
			map[string]any{},
			"",
		)
		return element
	}, map[string]any{})
	button871 := BUTTON661(map[string]any{})
	button411 := button871.Mount(div860)

	TEXT146 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		element := goat.CreateVirtualElements(
			"",
			"text",
			nil,
			"Click",
		)
		return element
	}, map[string]any{})
	text42 := TEXT146(map[string]any{})
	text42.Mount(button411)

	H1206 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
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
	h1122 := H1206(map[string]any{})
	h120 := h1122.Mount(div860)

	TEXT911 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		element := goat.CreateVirtualElements(
			"",
			"text",
			nil,
			"    Hello ",
		)
		return element
	}, map[string]any{})
	text164 := TEXT911(map[string]any{})
	text164.Mount(h120)

	TEXT382 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		element := goat.CreateVirtualElements(
			"",
			"text",
			nil,
			prop[proxy.Get("name").Key],
		)
		return element
	}, map[string]any{"name": name})
	text993 := TEXT382(map[string]any{})
	text993.Mount(h120)

	TEXT309 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		element := goat.CreateVirtualElements(
			"",
			"text",
			nil,
			", How are you?    ",
		)
		return element
	}, map[string]any{})
	text839 := TEXT309(map[string]any{})
	text839.Mount(h120)

	select {}
}
