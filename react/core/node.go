package core

import (
	"github.com/sjm1327605995/tenon/react/yoga"

	"gioui.org/layout"
)

type Gio interface {
	Layout(gtx layout.Context) layout.Dimensions
}
type Node interface {
	Yoga() *yoga.Node
	Children() []Node
	Update(ctx layout.Context)
	Gio() Gio
}
type emptyGio struct {
}

func (e *emptyGio) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Dimensions{
		Size: gtx.Constraints.Max,
	}
}

type emptyNode struct {
	yogaNode *yoga.Node
	gio      *emptyGio
}

func newEmptyNode() *emptyNode {
	node := yoga.NewNode()
	return &emptyNode{yogaNode: node, gio: &emptyGio{}}
}

func (e *emptyNode) Gio() Gio {
	return e.gio
}
func (e *emptyNode) Update(ctx layout.Context) {

}
func (e *emptyNode) Yoga() *yoga.Node {
	return e.yogaNode
}

func (e *emptyNode) Children() []Node {
	return nil
}
