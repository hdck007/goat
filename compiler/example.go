// package main

// import (
// 	"syscall/js"

// 	"github.com/hdck007/goat/goat"
// )

// func start() {

// 	name := "Hello"

// 	updateName := func() interface{} {
// 		name = "Mello"

// 		return nil
// 	}

// 	body := js.Global().Get("document").Call("getElementById", "root")

// 	DIV712 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
// 		children := []goat.VElement{}

// 		children330 := []goat.VElement{}

// 		text76 := goat.CreateVirtualElements(
// 			"",
// 			"text",
// 			nil,
// 			"Click",
// 		)
// 		children330 = append(children330, text76)

// 		button752 := goat.CreateVirtualElements(
// 			"",
// 			"button",
// 			map[string]any{},
// 			"",
// 			children330...,
// 		)
// 		children = append(children, button752)

// 		children972 := []goat.VElement{}

// 		text738 := goat.CreateVirtualElements(
// 			"",
// 			"text",
// 			nil,
// 			"    Hello ",
// 		)
// 		children972 = append(children972, text738)

// 		text167 := goat.CreateVirtualElements(
// 			"",
// 			"text",
// 			nil,
// 			proxy.Get("name"),
// 		)
// 		children972 = append(children972, text167)

// 		text779 := goat.CreateVirtualElements(
// 			"",
// 			"text",
// 			nil,
// 			", How are you?    ",
// 		)
// 		children972 = append(children972, text779)

// 		h1364 := goat.CreateVirtualElements(
// 			"",
// 			"h1",
// 			map[string]any{
// 				"class": "text-red",
// 			},
// 			"",
// 			children972...,
// 		)
// 		children = append(children, h1364)

// 		element := goat.CreateVirtualElements(
// 			"",
// 			"div",
// 			nil,
// 			"",
// 			children...,
// 		)
// 		return element
// 	}, map[string]any{
// 		"name": name,
// 	})
// 	div238 := DIV712(map[string]any{})
// 	div736 := div238.Mount(body)

// 	var updatename628 js.Func
// 	updatename628 = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
// 		updateName()
// 		div238.Patch(DIV712(map[string]any{
// 			"name": name,
// 		}))

// 		return nil
// 	})
// 	div736.Call("addEventListener", "click", updatename628)

// 	select {}
// }
