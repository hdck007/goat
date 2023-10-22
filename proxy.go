package main

type Hole struct {
	key string
}

type Props map[string]any

func (proxy *Props) Get(key string) Hole {
	currentElement := Hole{
		key: key,
	}
	(*proxy)[key] = currentElement
	return currentElement
}
