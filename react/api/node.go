package api

type Node interface {
	Component
	Element
}

type Component interface {
	Render() Node
}
