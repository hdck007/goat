package main

import (
	"syscall/js"

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
    count := 1
    userName := "John"
    message := "Welcome back"

    // Create the greeting component instance
    greetingComponent := Greeting(goat.Props{
        "name": userName,
    })

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
    }, goat.Props{
        "count":             count,
        "greetingComponent": greetingComponent,
    })()
}

func main() {
    component := ParentComponent(goat.Props{})
    root := js.Global().Get("document").Call("getElementById", "root")
    component.Mount(root)
}
