package tenon

import (
	"github.com/sjm1327605995/tenon/react/yoga"
	"github.com/sjm1327605995/tenon/style"
)

type ViewNode struct {
	node      *yoga.Node
	textStyle *style.TextStyle
}

func (v *ViewNode) Yoga() *yoga.Node {
	return v.node
}

func (v *ViewNode) Body(node ...INode) INode {
	return v
}
func View() *ViewNode {
	v := &ViewNode{
		node:      yoga.NewNode(),
		textStyle: &style.TextStyle{},
	}
	return v
}
func (v *ViewNode) Style(opts ...style.Option) *ViewNode {

	style.ApplyOptions(v.textStyle, opts...)

	return v
}
