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
		Patch: func(elements []interface{}) {},
	}
	blockGenerator := func(props goat.Props) goat.ArrayBlock {
		return goat.ArrayBlock{
			Patch: func(elements []interface{}) {},
			Mount: func(js.Value) js.Value {
				return js.Null()
			},
		}
	}

	blockGenerator = func(currentProps goat.Props) goat.ArrayBlock {
		return goat.ArrayBlockElement("elements",
			func(element interface{}, index int) goat.Block {
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

	block8604571377814938500 := &goat.Block{
		Patch: func(goat.Block) {},
		Mount: func(js.Value) js.Value {
			return js.Undefined()
		},
	}
	derivedProps := &goat.Props{}

	block8604571377814938500generator := func(props goat.Props) goat.Block {
		return goat.Block{
			Patch: func(goat.Block) {},
			Mount: func(js.Value) js.Value {
				return js.Null()
			},
		}
	}

	context := &goat.Context{
		CreateBlock: block8604571377814938500generator,
		Block:       block8604571377814938500,
		Props:       &goat.Props{},
	}

	elements, setElements := context.CreateState([]string{}, "elements")
	inputValue, setInputValue := context.CreateState("", "inputValue")

	remove := func(slice []string, index int) []string {
		return append(slice[:index], slice[index+1:]...)
	}

	handleChange := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		inputValue = args[0].Get("target").Get("value").String()
		setInputValue(inputValue)
		return nil
	})

	handleClick := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if inputValue == "" {
			return nil
		}

		switch convertedElement := elements.(type) {
		case []string:
			elements = append(convertedElement, inputValue.(string))
		case []interface{}:
			elements = append(convertedElement, inputValue)
		}

		setElements(elements)
		setInputValue("")
		return nil
	})

	handleDelete := func(index int) {
		elements = remove(elements.([]string), index)
		setElements(elements)
	}

	component_1602885801659547044 := TodoList(goat.Props{
		"elements":     elements,
		"handleDelete": handleDelete,
	})

	derivedProps = &goat.Props{"inputValue": inputValue, "component_1602885801659547044": component_1602885801659547044, "elements": elements}
	for k, v := range props {
		(*derivedProps)[k] = v
	}

	block8604571377814938500generator = func(currentProps goat.Props) goat.Block {
		return goat.BlockElement(func(p goat.Props) goat.Vnode {
			return goat.CreateVirtualElements("div", goat.Props{}, goat.CreateVirtualElements("div", goat.Props{"class": "\"p-4 text-white\""}, goat.CreateVirtualElements("h1", goat.Props{"class": "\"text-3xl font-bold mb-2\""}, goat.CreateVirtualElements("text", nil,
				&goat.TextOrElement{
					StringValue: "\n        Todo List App\n    ",
					Element:     nil,
				},
			)), goat.CreateVirtualElements("p", goat.Props{"class": "\"mb-4\""}, goat.CreateVirtualElements("text", nil,
				&goat.TextOrElement{
					StringValue: "\n        A simple todo list app built using Goat VDOM and Go.\n    ",
					Element:     nil,
				},
			)), goat.CreateVirtualElements("input", goat.Props{"class": "\"border text-black border-white p-2 mr-4\"", "value": goat.Get("inputValue"), "@input": handleChange}), goat.CreateVirtualElements("button", goat.Props{"class": "\"bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded\"", "@click": handleClick}, goat.CreateVirtualElements("text", nil,
				&goat.TextOrElement{
					StringValue: "\n        Add\n    ",
					Element:     nil,
				},
			)), goat.Get("component_1602885801659547044")))
		}, currentProps)()
	}

	context.CreateBlock = block8604571377814938500generator
	block := block8604571377814938500generator(*derivedProps)
	block8604571377814938500 = &block

	context.Block = block8604571377814938500
	context.Props = derivedProps

	return *block8604571377814938500
}

func main() {
	done := make(chan struct{}, 0)
	component := App(goat.Props{})
	root := js.Global().Get("document").Call("getElementById", "root")
	component.Mount(root)

	<-done
}
