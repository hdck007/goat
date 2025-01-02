package main

import (
	"syscall/js"

	"github.com/hdck007/goat/goat"
)

func main() {
	done := make(chan struct{}, 0)
	component := App(goat.Props{})
	root := js.Global().Get("document").Call("getElementById", "root")
	component.Mount(root)

	<-done
}
