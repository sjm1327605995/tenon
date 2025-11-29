package react

import (
	"github.com/sjm1327605995/tenon/react/common"
	"github.com/sjm1327605995/tenon/react/common/renderer"
	"github.com/sjm1327605995/tenon/react/core"
	"github.com/sjm1327605995/tenon/react/style"
	"github.com/sjm1327605995/tenon/react/style/unit"
)

type ReactDOM struct {
	root     common.Component
	width    unit.Unit
	height   unit.Unit
	renderer common.Renderer
}

func (h *ReactDOM) Render(children ...common.Component) error {
	element := core.NewView().
		Style(style.NewStyle().Width(800).Height(600)).
		Child(children...)
	h.root = element
	h.root.Render()
	h.renderer.SetElement(element)
	return h.renderer.Run()
}
func NewReactDOM() *ReactDOM {
	return &ReactDOM{
		width:    unit.Pt(800),
		height:   unit.Pt(600),
		renderer: renderer.NewGio(),
	}
}
