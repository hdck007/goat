package main

import (
	"syscall/js"

	"github.com/hdck007/goat/goat"
)

func start() {

	TEXT654 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		element := goat.CreateVirtualElements(
			"",
			"text",
			nil,
			"1",
			goat.VElement{},
		)
		return element
	}, map[string]any{})
	text558 := TEXT654(map[string]any{})
	H1103 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		childElements := []goat.VElement{text558}
		element := goat.CreateVirtualElements(
			"",
			"h1",
			nil,
			"",
			childElements...,
		)
		return element
	}, map[string]any{})
	h1277 := H1103(map[string]any{})
	DIV391 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		childElements := []goat.VElement{h1277}
		element := goat.CreateVirtualElements(
			"",
			"div",
			nil,
			"",
			childElements...,
		)
		return element
	}, map[string]any{})
	div400 := DIV391(map[string]any{})
	body := js.Global().Get("document").Call("getElementById", "root")
	div400.Mount(body)

	select {}
}
