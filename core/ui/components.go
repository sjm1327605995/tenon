package ui

type IComponent interface {
	Element
	Hooks()
}
type Component struct {
}

func (c *Component) Render() Element {
	//TODO implement me
	panic("implement me")
}

func (c *Component) Hooks() {
	//TODO implement me
	panic("implement me")
}


