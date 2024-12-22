package main

import (
	"syscall/js"

	"github.com/hdck007/goat/goat"
)

func Component(props goat.Props) goat.Block {
	// Script block variables
	name := "Hello"
	padre := "Adiosa"

	return goat.BlockElement(func(p goat.Props) goat.Vnode {
		return goat.CreateVirtualElements("div", nil,
			// Button element
			goat.CreateVirtualElements("button",
				goat.Props{
					"class": "bg-none py-4 rounded-mg border border-black m-4",
				},
				goat.CreateVirtualElements("text", nil,
					&goat.TextOrElement{
						StringValue: "Click",
						Element:     nil,
					},
				),
			),
			// H1 element with interpolated values
			goat.CreateVirtualElements("h1",
				goat.Props{
					"class": "text-red",
				},
				goat.CreateVirtualElements("text", nil,
					&goat.TextOrElement{
						StringValue: "Hello ",
						Element:     nil,
					},
				),
				goat.Get("name"),
				goat.CreateVirtualElements("text", nil,
					&goat.TextOrElement{
						StringValue: ", How are you? mello ",
						Element:     nil,
					},
				),
				goat.Get("padre"),
				goat.CreateVirtualElements("text", nil,
					&goat.TextOrElement{
						StringValue: " ami imorto",
						Element:     nil,
					},
				),
			),
		)
	}, goat.Props{
		"name":  name,
		"padre": padre,
	})()
}

func start() {
	component := Component(goat.Props{})
	root := js.Global().Get("document").Call("getElementById", "root")
	component.Mount(root)
}
