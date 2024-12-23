package main

import (
	"syscall/js"
	"time"

	"github.com/hdck007/goat/goat"
)

// Child Greeting component
func Greeting(props goat.Props) goat.Block {
	return goat.BlockElement(func(p goat.Props) goat.Vnode {
		return goat.CreateVirtualElements("div",
			goat.Props{"class": "greeting"},
			goat.CreateVirtualElements("p", nil,
				goat.CreateVirtualElements("text", nil,
					&goat.TextOrElement{
						StringValue: "Hello, ",
						Element:     nil,
					},
				),
				goat.Get("name"),
				goat.CreateVirtualElements("text", nil,
					&goat.TextOrElement{
						StringValue: "!",
						Element:     nil,
					},
				),
			),
		)
	}, props)()
}

// Parent component with nested components
func ParentComponent(props goat.Props) goat.Block {
	userName := "John"
	message := "Welcome back"

	// Create the greeting component instance
	greetingComponent := Greeting(goat.Props{
		"name": userName,
	})

	derivedProps := goat.Props{
		"count":             props["count"],
		"greetingComponent": greetingComponent,
	}

	return goat.BlockElement(func(p goat.Props) goat.Vnode {
		return goat.CreateVirtualElements("div",
			goat.Props{"class": "container"},
			// Counter section
			goat.CreateVirtualElements("div",
				goat.Props{"class": "counter-section"},
				goat.CreateVirtualElements("h2", nil,
					goat.CreateVirtualElements("text", nil,
						&goat.TextOrElement{
							StringValue: "Count: ",
							Element:     nil,
						},
					),
					goat.Get("count"),
				),
			),
			// Message section
			goat.CreateVirtualElements("div",
				goat.Props{"class": "message"},
				goat.CreateVirtualElements("text", nil,
					&goat.TextOrElement{
						StringValue: message,
						Element:     nil,
					},
				),
			),
			// Nested Greeting component using placeholder
			goat.CreateVirtualElements("div",
				goat.Props{"class": "nested-component"},
				goat.Get("greetingComponent"),
			),
		)
	}, derivedProps)()
}

func main() {
	count := 0
	component := ParentComponent(goat.Props{
		"count": 1,
	})
	root := js.Global().Get("document").Call("getElementById", "root")
	component.Mount(root)

	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				count += 1
				println(count)
				newComponent := ParentComponent(goat.Props{
					"count": count,
				})
				component.Patch(newComponent)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	select {}
}
