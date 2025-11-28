package common

type Node interface {
	Component
	Element
}

type Component interface {
	Render() Node
}
