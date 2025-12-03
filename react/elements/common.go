package elements

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/sjm1327605995/tenon/react/api"
	"github.com/sjm1327605995/tenon/react/api/styles"
	"github.com/sjm1327605995/tenon/react/yoga"
)

var Metric unit.Metric

type ElementBase struct {
	Node     *yoga.Node
	Children []api.Element
}

func (e *ElementBase) Yoga() *yoga.Node {
	return e.Node
}

func (e *ElementBase) SetExtendedStyle(style styles.IExtendedStyle) {

}

func (e *ElementBase) GetChildrenCount() int {
	return len(e.Children)
}

func (e *ElementBase) GetChildAt(index int) api.Element {
	return e.Children[index]
}

func (e *ElementBase) Paint(ctx layout.Context) {
	//TODO implement me
	panic("implement me")
}

func (e *ElementBase) GetChildren() []api.Element {
	return e.Children
}

func (e *ElementBase) SetStyle(style *styles.Style) {
	//TODO implement me
	panic("implement me")
}
