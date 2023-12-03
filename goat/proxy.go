package goat

type Hole struct {
	Key string
}

type Props map[string]any

func (proxy *Props) Get(key string) Hole {
	currentElement := Hole{
		Key: key,
	}
	(*proxy)[key] = currentElement
	return currentElement
}
