package main

import (
	"syscall/js"

	"github.com/hdck007/goat/goat"
)

func start() {

	TEXT924 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		element := goat.CreateVirtualElements(
			"",
			"text",
			nil,
			"1",
			goat.VElement{},
		)
		return element
	}, map[string]any{})
	text763 := TEXT924(map[string]any{})
	H1670 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		childElements := []goat{text763}
		element := goat.CreateVirtualElements(
			"",
			"h1",
			nil,
			"",
			childElements...,
		)
		return element
	}, map[string]any{})
	h138 := H1670(map[string]any{})
	DIV643 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		childElements := []goat{h138}
		element := goat.CreateVirtualElements(
			"",
			"div",
			nil,
			"",
			childElements...,
		)
		return element
	}, map[string]any{})
	div343 := DIV643(map[string]any{})
	body := js.Global().Get("document").Call("getElementById", "root")
	div343.Mount(body)

	select {}
}
