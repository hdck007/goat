package goat

type Context struct {
	CreateBlock func(props Props) Block
	Block       *Block
	Props       *Props
}

func (c *Context) CreateState(initialValue any, key string) (any, func(any)) {
	var state any = initialValue
	return state, func(newValue any) {
		state = newValue
		if c.Block != nil {
			(*c.Props)[key] = newValue
			(*c.Block).Patch(c.CreateBlock(*c.Props))
		}
	}
}
