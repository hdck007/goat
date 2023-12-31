### What the compiler must do?

Starting with a very simple example

/count.goat
```jsx
<div>
  <h1>{count}</h1>
</div>

```

To

/count-component.go
```go
package main

import (
	"syscall/js"

	"github.com/hdck007/goat/goat"
)

func main() {

	TEXT46 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		element := goat.CreateVirtualElements(
			"",
			"text",
			nil,
			"1",
			goat.VElement{},
		)
		return element
	}, map[string]any{})
	text609 := TEXT46(map[string]any{})

	H1919 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		childElements := []goat.VElement{text609}
		element := goat.CreateVirtualElements(
			"",
			"h1",
			nil,
			"",
			childElements...,
		)
		return element
	}, map[string]any{})
	h1268 := H1919(map[string]any{})

	DIV571 := goat.BlockElement(func(proxy *goat.Props, prop goat.Props) goat.VElement {
		childElements := []goat.VElement{h1268}
		element := goat.CreateVirtualElements(
			"",
			"div",
			nil,
			"",
			childElements...,
		)
		return element
	}, map[string]any{})
	div971 := DIV571(map[string]any{})

	body := js.Global().Get("document").Call("getElementById", "root")
	div971.Mount(body)

	select {}
}

```
This is easy, just create a tree and create equivalent code
I have updated the code for this as well in the compiler

But one caveat and my second dilemma in this compiler is that the child should be a `velement` but can be block elements, so the vdom needs to adjust this. 
But these things kind of shape the dx and stuff so I will have to think about it rather than just writing code for compatibility.
