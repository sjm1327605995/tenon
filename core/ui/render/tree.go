package render

import (
	"image"
	"slices"

	"gioui.org/layout"
	"gioui.org/op"
)

type Node struct {
	Max      image.Point
	Offset   image.Point
	Render   Render
	Children []*Node
}

func NewNode(render Render) *Node {
	return &Node{Render: render}
}
func (n *Node) InsertChild(child *Node) {
	n.Children = append(n.Children, child)
}
func (n *Node) Remove(child *Node) {
	n.Children = slices.DeleteFunc(n.Children, func(node *Node) bool {
		return node == child
	})
}
func (n *Node) Paint(ctx layout.Context) {
	defer op.Offset(n.Offset).Push(ctx.Ops).Pop()
	ctx.Constraints.Max = n.Max
	n.Render.Layout(ctx)

	for _, child := range n.Children {
		child.Paint(ctx)
	}
}
