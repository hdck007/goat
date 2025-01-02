package main

import (
	"syscall/js"

	"github.com/hdck007/goat/goat"
)

func App(props goat.Props) goat.Block {
	block530748580410125483 := &goat.Block{
		Patch: func(goat.Block) {},
		Mount: func(js.Value) js.Value {
			return js.Undefined()
		},
	}
	derivedProps := &goat.Props{}

	block530748580410125483generator := func(props goat.Props) goat.Block {
		return goat.Block{
			Patch: func(goat.Block) {},
			Mount: func(js.Value) js.Value {
				return js.Null()
			},
		}
	}

	context := &goat.Context{
		CreateBlock: block530748580410125483generator,
		Block:       block530748580410125483,
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
		elements = append(elements.([]string), inputValue.(string))
		setElements(elements)
		setInputValue("")
		return nil
	})

	handleDelete := func(index int) {
		elements = remove(elements.([]string), index)
		setElements(elements)
	}

	component_4680478051752220265 := TodoElement(goat.Props{
		"handleDelete": handleDelete})

	derivedProps = &goat.Props{"inputValue": inputValue, "handleChange": handleChange, "handleClick": handleClick, "component_4680478051752220265": component_4680478051752220265}
	for k, v := range props {
		(*derivedProps)[k] = v
	}
	propsToPass := goat.Props{"inputValue": inputValue, "handleChange": handleChange, "handleClick": handleClick, "component_4680478051752220265": component_4680478051752220265}
	for k, v := range props {
		propsToPass[k] = v
	}

	block530748580410125483generator = func(currentProps goat.Props) goat.Block {
		return goat.BlockElement(func(p goat.Props) goat.Vnode {
			return goat.CreateVirtualElements("div", goat.Props{},
				goat.CreateVirtualElements("div", goat.Props{"class": "\"p-4 text-white\""},
					goat.CreateVirtualElements("h1", goat.Props{"class": "\"text-3xl font-bold mb-2\""},
						goat.CreateVirtualElements("text", nil,
							&goat.TextOrElement{
								StringValue: "\n        Todo List App\n    ",
								Element:     nil,
							},
						),
					),
					goat.CreateVirtualElements("p", goat.Props{"class": "\"mb-4\""}, goat.CreateVirtualElements("text", nil,
						&goat.TextOrElement{
							StringValue: "\n        A simple todo list app built using Goat VDOM and Go.\n    ",
							Element:     nil,
						},
					)), goat.CreateVirtualElements("input", goat.Props{"class": "\"border text-black border-white p-2 mr-4\"", "value": goat.Get("inputValue"), "handleChange": handleChange}), goat.CreateVirtualElements("button", goat.Props{"class": "\"bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded\"", "handleClick": handleClick}, goat.CreateVirtualElements("text", nil,
						&goat.TextOrElement{
							StringValue: "\n        Add\n    ",
							Element:     nil,
						},
					)), goat.Get("component_4680478051752220265")))
		}, currentProps)()
	}

	context.CreateBlock = block530748580410125483generator
	block := block530748580410125483generator(propsToPass)

	context.Block = &block
	context.Props = derivedProps

	return block
}
