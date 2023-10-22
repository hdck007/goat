# Goat

A Golang-based virtual dom.

Live demo https://goat-two.vercel.app/

## Steps to run locally
1. Clone the repository `git clone https://github.com/hdck007/goat.git`
2. navigate inside the repository `cd goat`
3. build wasm `GOOS=js GOARCH=wasm go build -o main.wasm -buildvcs=false`. Yes, it's a big command feel free to use npm scripts, or as in my case I have created an alias
4. Serve the HTML using some HTTP server (I use the vscode live server extension currently)

> Note: I will soon develop a dev server for this

## Motivation
https://github.com/aidenybai/million

## Architecture:

Currently working on it.....

## TODOs
1. Figuring out how to bind event listeners using Golang
2. Prop serialization logic, (this would help me convert `style: map[string]string|number` to inline styles for dom elements)
3. JSX parser
4. A development server
5. Some benchmarking --> At least starting with some basics

> These todos are the first things that came to my mind at the moment of writing them and the list doesn't indicate the order or priority in which it will be implemented
