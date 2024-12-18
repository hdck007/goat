package main

import (
	"syscall/js"

	"github.com/hdck007/goat/goat"
)

func main() {

	x, y := 100, 100

	name := "Hello"

	padre := "Adiosa"

	updateName := func() interface{} {
		name = "Mello"

		return nil
	}

	updatePostion := func(event js.Value) interface{} {
		x = event.Get("clientX").Int()
		y = event.Get("clientY").Int()
		return nil
	}

	body := js.Global().Get("document").Call("getElementById", "root")

	DIV917 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.Vnode {
		children := []goat.Vnode{}

		children761 := []goat.Vnode{}

		text326 := goat.CreateVirtualElements(
			"text",
			nil,
			&goat.UnionNode{Element: nil, StringValue: "Click"},
		)
		children761 = append(children761, text326)

		button243 := goat.CreateVirtualElements(
			"button",
			map[string]any{},
			children761...,
		)
		children = append(children, button243)

		children27 := []goat.Vnode{}

		text254 := goat.CreateVirtualElements(
			"text",
			nil,
			&goat.UnionNode{Element: nil, StringValue: "    Hello "},
		)
		children27 = append(children27, text254)

		text205 := goat.CreateVirtualElements(
			"text",
			nil,
			proxy.Get("name"),
		)
		children27 = append(children27, text205)

		text13 := goat.CreateVirtualElements(
			"text",
			nil,
			&goat.UnionNode{Element: nil, StringValue: ", How are you?    mello "},
		)
		children27 = append(children27, text13)

		text689 := goat.CreateVirtualElements(
			"text",
			nil,
			proxy.Get("padre"),
		)
		children27 = append(children27, text689)

		text849 := goat.CreateVirtualElements(
			"text",
			nil,
			&goat.UnionNode{Element: nil, StringValue: " ami imorto    "},
		)
		children27 = append(children27, text849)

		h1147 := goat.CreateVirtualElements(
			"h1",
			map[string]any{
				"class":    "text-red",
				"left":     proxy.Get("x"),
				"top":      proxy.Get("y"),
				"position": "absolute",
				// "style": "postion: absolute; left: " + strconv.Itoa(newX) + "px; top: " + strconv.Itoa(newY) + "px;",
			},
			children27...,
		)
		children = append(children, h1147)

		element := goat.CreateVirtualElements(
			"div",
			nil,
			children...,
		)
		return element
	}, map[string]any{
		"x":     x,
		"y":     y,
		"name":  name,
		"padre": padre,
	})
	div335 := DIV917(map[string]any{
		"x":     x,
		"y":     y,
		"name":  name,
		"padre": padre,
	})
	div327 := div335.Mount(body)

	var updatename88 js.Func
	updatename88 = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		updateName()
		div335.Patch(DIV917(map[string]any{
			"name":  name,
			"padre": padre,
			"x":     x,
			"y":     y,
		}))

		return nil
	})

	updatePostionRandom := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		updatePostion(js.Global().Get("event"))
		// println(x, y)
		div335.Patch(DIV917(map[string]any{
			"name":  name,
			"padre": padre,
			"x":     x,
			"y":     y,
		}))
		return nil
	})

	div327.Call("addEventListener", "click", updatename88)
	body.Call("addEventListener", "mousemove", updatePostionRandom)

	select {}
}
