package main

type Hole struct {
	key string
}

type Proxy map[string]Hole

func (proxy *Props) Get(key string) Hole {
	currentElement := Hole{
		key: key,
	}
	(*proxy)[key] = currentElement

	return currentElement
}
