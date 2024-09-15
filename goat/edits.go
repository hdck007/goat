package goat

type EditAttribute struct {
	editType  string
	path      []int
	attribute string
	hole      string
}

type EditChild struct {
	editType string
	path     []int
	index    int
	hole     string
}

type Edits interface {
	isChildEdit() bool
	isAttributeEdit() bool
	getChildEditValue() *EditChild
	getAttributeEditValue() *EditAttribute
}

func (edit *EditAttribute) isChildEdit() bool {
	return false
}

func (edit *EditChild) isChildEdit() bool {
	return true
}

func (edit *EditAttribute) isAttributeEdit() bool {
	return true
}

func (edit *EditChild) isAttributeEdit() bool {
	return false
}

func (edit *EditAttribute) getAttributeEditValue() *EditAttribute {
	return edit
}

func (edit *EditChild) getAttributeEditValue() *EditAttribute {
	return nil
}

func (edit *EditChild) getChildEditValue() *EditChild {
	return edit
}

func (edit *EditAttribute) getChildEditValue() *EditChild {
	return nil
}
