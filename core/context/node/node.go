package node

import (
	"github.com/dhconnelly/rtreego"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/millken/yoga"
)

type Rect struct {
	W float32
	H float32
}

type INode interface {
	Measure()
	OnLayout()
	OnDraw(r *ebiten.Image, rtree *rtreego.Rtree)
	Yoga() *yoga.Node
	SetPositon(x, y float32)
	GetPositon() (x, y float32)
	OnEvent()
}
type EventHandler interface {
	Order() int
	Handle()
}
type OnHoverEvent struct {
	f func()
}

func (o *OnHoverEvent) Order() int {
	return 1
}
func (o *OnHoverEvent) Handle() {
	o.f()
}

type OnClickEvent struct {
	f func()
}

func (o *OnClickEvent) Order() int {
	return 2
}
func (o *OnClickEvent) Handle() {
	o.f()
}

type OnMouseOverEvent struct {
	f func()
}

func (o *OnMouseOverEvent) Order() int {
	return 3
}
func (o *OnMouseOverEvent) Handle() {
	o.f()
}

type Node struct {
	*yoga.Node
	IsRtreeNode int //0 不是 1 待添加 2 添加完毕
	rtreeRect   rtreego.Rect
	children    []INode
	X           float32
	Y           float32
	Event       []EventHandler
}

func (n *Node) OnEvent() {
	if len(n.Event) > 0 {
		for i := range n.Event {
			n.Event[i].Handle()
		}
	}
}
func (n *Node) SetRtreeRect(rtree *rtreego.Rtree) {

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

func (n *Node) Measure() {
	for i := range n.children {
		n.children[i].Measure()
	}
}
func (n *Node) OnDraw(r *ebiten.Image, rtree *rtreego.Rtree) {
	for i := range n.children {
		yogaNode := n.children[i].Yoga()
		x, y := yogaNode.LayoutLeft(), yogaNode.LayoutTop()
		n.children[i].SetPositon(n.X+x, n.Y+y)
		n.children[i].OnDraw(r, rtree)
	}
	if n.IsRtreeNode == 1 {
		n.SetDirtiedFunc(func(node *yoga.Node) {
			n.SetRtreeRect(rtree)
		})
		n.SetRtreeRect(rtree)
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
