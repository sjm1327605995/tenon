package node

import (
	"github.com/dhconnelly/rtreego"
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
	OnHover()
	OnClick()
}

type Node struct {
	*yoga.Node
	IsRtreeNode int //0 不是 1 待添加 2 添加完毕
	rtreeRect   rtreego.Rect
	children    []INode
	X           float32
	Y           float32
	Hover       func()
	Click       func()
}

func (n *Node) SetRtreeRect() {

	if n.IsRtreeNode == 2 {
		rtree.Delete(n)
	}
	n.rtreeRect, _ = rtreego.NewRect(rtreego.Point{float64(n.X), float64(n.Y)},
		[]float64{float64(n.Node.StyleGetWidth()), float64(n.Node.StyleGetHeight())})
	rtree.Insert(n)
	n.IsRtreeNode = 2
}
func (n *Node) Bounds() rtreego.Rect {
	return n.rtreeRect
}
func (n *Node) OnLayout() {
	yoga.CalculateLayout(n.Node, yoga.Undefined, yoga.Undefined, yoga.DirectionLTR)
}
func (n *Node) OnHover() {
	if n.Hover != nil {
		n.Hover()
	}
}
func (n *Node) OnClick() {
	if n.Click != nil {
		n.Click()
	}
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
	if n.IsRtreeNode == 1 {
		n.SetDirtiedFunc(func(node *yoga.Node) {
			n.SetRtreeRect()
		})
		n.SetRtreeRect()
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
