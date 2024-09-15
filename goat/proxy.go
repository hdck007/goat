package goat

type Hole struct {
	key string
}

func (h *Hole) isString() bool {
	return false
}

func (h *Hole) isVElement() bool {
	return false
}

func (h *Hole) isHole() bool {
	return true
}

func (h *Hole) getStringValue() string {
	return ""
}

func (h *Hole) getVElement() *VElement {
	return nil
}

func (h *Hole) getHoleValue() *Hole {
	return h
}

type Props map[string]any

func (proxy *Props) Get(keyValue string) *Hole {
	currentElement := Hole{
		key: keyValue,
	}
	(*proxy)[keyValue] = currentElement
	return &currentElement
}
