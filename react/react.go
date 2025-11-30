package react

import (
	"github.com/sjm1327605995/tenon/react/api"
	"github.com/sjm1327605995/tenon/react/api/renderer"
	"github.com/sjm1327605995/tenon/react/elements"
	"github.com/sjm1327605995/tenon/react/event"
	"github.com/sjm1327605995/tenon/react/styles"
)

type ReactDOM struct {
	root     api.Component
	renderer api.Renderer
	style    *styles.Style
	event    chan event.Event
}

func (h *ReactDOM) Render(children ...api.Component) error {
	element := elements.NewView().
		Style(h.style).
		Child(children...)
	h.root = element
	h.root.Render()
	h.renderer.SetElement(element)
	return h.renderer.Run()
}
func NewReactDOM() *ReactDOM {
	return &ReactDOM{
		renderer: renderer.NewGio(),
		style:    styles.NewStyle().Width(800).Height(600),
		event:    make(chan event.Event, 10),
	}
}
