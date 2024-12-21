package goat

type Placeholder struct {
	key string
}

func (h *Placeholder) isString() bool {
	return false
}

func (h *Placeholder) isVElement() bool {
	return false
}

func (h *Placeholder) isPlaceholder() bool {
	return true
}

func (h *Placeholder) getStringValue() string {
	return ""
}

func (h *Placeholder) getVElement() *VElement {
	return nil
}

func (h *Placeholder) getPlaceholderValue() *Placeholder {
	return h
}

type Props map[string]any

func Get(keyValue string) *Placeholder {
	currentElement := Placeholder{
		key: keyValue,
	}
	return &currentElement
}
