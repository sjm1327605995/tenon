package node

import (
	"github.com/millken/yoga"
)

type Rect struct {
	W float32
	H float32
}

type INode interface {
	Measure()
	OnLayout()
	OnDraw(r Renderer)
	Yoga() *yoga.Node
	SetPositon(x, y float32)
	GetPositon() (x, y float32)
}

type Node struct {
	*yoga.Node
	children []INode
	X        float32
	Y        float32
}

func (n *Node) OnLayout() {
	yoga.CalculateLayout(n.Node, yoga.Undefined, yoga.Undefined, yoga.DirectionLTR)
}
func (n *Node) Measure() {
	for i := range n.children {
		n.children[i].Measure()
	}
}
func (n *Node) OnDraw(r Renderer) {
	for i := range n.children {
		yogaNode := n.children[i].Yoga()
		x, y := yogaNode.LayoutLeft(), yogaNode.LayoutTop()
		n.children[i].SetPositon(n.X+x, n.Y+y)
		n.children[i].OnDraw(r)
	}
}
func (n *Node) SetPositon(x, y float32) {
	n.X = x
	n.Y = y
}
func (n *Node) GetPositon() (x, y float32) {
	return n.X, n.Y
}
func (n *Node) Yoga() *yoga.Node {
	return n.Node
}
