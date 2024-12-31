package main

import (
	"syscall/js"

	"github.com/hdck007/goat/goat"
)

func TodoElement(props goat.Props) goat.Block {

	derivedProps := make(goat.Props, 0)

	for k, v := range props {
		derivedProps[k] = v
	}

	handleDeleteCallback := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		derivedProps["handleDelete"].(func(index int))(derivedProps["index"].(int))
		return nil
	})

	return goat.BlockElement(func(p goat.Props) goat.Vnode {
		return goat.CreateVirtualElements("div",
			goat.Props{"class": "flex items-center gap-2 p-4"},
			goat.CreateVirtualElements(
				"span",
				goat.Props{
					"class": "text-white",
				},
				goat.Get("element"),
			),
			goat.CreateVirtualElements(
				"button",
				goat.Props{
					"class":  "bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded",
					"@click": handleDeleteCallback,
				},
				goat.CreateVirtualElements("text", nil,
					&goat.TextOrElement{
						StringValue: "Delete",
						Element:     nil,
					},
				),
			),
		)
	}, derivedProps)()

}

func TodoList(props goat.Props) goat.ArrayBlock {

	block := goat.ArrayBlock{
		Mount: func(js.Value) js.Value {
			return js.Null()
		},
		Patch: func(elements []string) {},
	}
	blockGenerator := func(props goat.Props) goat.ArrayBlock {
		return goat.ArrayBlock{
			Patch: func(elements []string) {},
			Mount: func(js.Value) js.Value {
				return js.Null()
			},
		}
	}

	blockGenerator = func(currentProps goat.Props) goat.ArrayBlock {
		return goat.ArrayBlockElement("elements",
			func(element string, index int) goat.Block {
				return TodoElement(goat.Props{
					"element":      element,
					"handleDelete": props["handleDelete"],
					"index":        index,
				})
			},
			currentProps,
		)()
	}

	block = blockGenerator(goat.Props{
		"elements":     props["elements"],
		"handleDelete": props["handleDelete"],
	})

	return block

}

func remove(slice []string, index int) []string {
	return append(slice[:index], slice[index+1:]...)
}

func App(props goat.Props) goat.Block {

	block := goat.Block{
		Patch: func(goat.Block) {},
		Mount: func(js.Value) js.Value {
			return js.Null()
		},
	}
	blockGenerator := func(props goat.Props) goat.Block {
		return goat.Block{
			Patch: func(goat.Block) {},
			Mount: func(js.Value) js.Value {
				return js.Null()
			},
		}
	}

	elements := make([]string, 0)

	text := ""

	todoList := goat.ArrayBlock{}

	handleChange := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		text = args[0].Get("target").Get("value").String()
		block.Patch(blockGenerator(
			goat.Props{
				"value":    text,
				"todoList": todoList,
				"elements": elements,
			},
		))
		return nil
	})

	handleClick := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if text == "" {
			return nil
		}
		elements = append(elements, text)
		text = ""
		block.Patch(blockGenerator(
			goat.Props{
				"value":    text,
				"todoList": todoList,
				"elements": elements,
			},
		))
		return nil
	})

	handleDelete := func(index int) {
		elements = remove(elements, index)
		block.Patch(blockGenerator(
			goat.Props{
				"elements": elements,
				"text":     text,
				"todoList": todoList,
			},
		))
	}

	todoList = TodoList(goat.Props{
		"elements":     elements,
		"handleDelete": handleDelete,
	})

	blockGenerator = func(currentProps goat.Props) goat.Block {
		return goat.BlockElement(func(p goat.Props) goat.Vnode {
			return goat.CreateVirtualElements("div",
				goat.Props{"class": "p-4 text-white"},
				// Header and Description
				goat.CreateVirtualElements("h1",
					goat.Props{"class": "text-3xl font-bold mb-2"},
					goat.CreateVirtualElements("text", nil, &goat.TextOrElement{
						StringValue: "Todo List App",
					}),
				),
				goat.CreateVirtualElements("p",
					goat.Props{"class": "mb-4"},
					goat.CreateVirtualElements("text", nil, &goat.TextOrElement{
						StringValue: "A simple todo list app built using Goat VDOM and Go.",
					}),
				),
				// Input and Button
				goat.CreateVirtualElements("input",
					goat.Props{"class": "border text-black border-white p-2 mr-4", "@input": handleChange,
						"value": goat.Get("value")},
				),
				goat.CreateVirtualElements("button", goat.Props{
					"class":  "bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded",
					"@click": handleClick,
				},
					goat.CreateVirtualElements("text", nil, &goat.TextOrElement{
						StringValue: "Add",
						Element:     nil,
					}),
				),
				// Todo List
				goat.Get("todoList"),
			)
		}, currentProps)()
	}

	block = blockGenerator(goat.Props{
		"value":    text,
		"todoList": todoList,
		"elements": elements,
	})

	return block
}

func main() {
	done := make(chan struct{}, 0)
	component := App(goat.Props{})
	root := js.Global().Get("document").Call("getElementById", "root")
	component.Mount(root)

	<-done
}
