package core

import (
	"github.com/sjm1327605995/tenon/react/yoga"

	"gioui.org/layout"
)

type Node interface {
	Yoga() *yoga.Node

	Layout(gtx layout.Context) layout.Dimensions
}
type emptyNode struct {
	yogaNode *yoga.Node
}

func newEmptyNode() *emptyNode {
	node := yoga.NewNode()
	return &emptyNode{yogaNode: node}
}

func (e *emptyNode) Yoga() *yoga.Node {
	return e.yogaNode
}

func (e *emptyNode) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Dimensions{
		Size: gtx.Constraints.Max,
	}
}
