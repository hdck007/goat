package main

import (
	"syscall/js"

	"github.com/hdck007/goat/goat"
)

func TodoElement(props goat.Props) goat.Block {
	block1379884031591310177 := &goat.Block{
		Patch: func(goat.Block) {},
		Mount: func(js.Value) js.Value {
			return js.Undefined()
		},
	}
	derivedProps := &goat.Props{}

	block1379884031591310177generator := func(props goat.Props) goat.Block {
		return goat.Block{
			Patch: func(goat.Block) {},
			Mount: func(js.Value) js.Value {
				return js.Null()
			},
		}
	}

	context := &goat.Context{
		CreateBlock: block1379884031591310177generator,
		Block:       block1379884031591310177,
		Props:       &goat.Props{},
	}

	derivedProps = &goat.Props{}
	for k, v := range props {
		(*derivedProps)[k] = v
	}
	propsToPass := goat.Props{}
	for k, v := range props {
		propsToPass[k] = v
	}

	block1379884031591310177generator = func(currentProps goat.Props) goat.Block {
		return goat.BlockElement(func(p goat.Props) goat.Vnode {
			return goat.CreateVirtualElements("div", goat.Props{}, goat.CreateVirtualElements("text", nil,
				&goat.TextOrElement{
					StringValue: "\n    Hello\n",
					Element:     nil,
				},
			))
		}, currentProps)()
	}

	context.CreateBlock = block1379884031591310177generator
	block := block1379884031591310177generator(propsToPass)

	context.Block = &block
	context.Props = derivedProps

	return block
}
