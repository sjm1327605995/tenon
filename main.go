package main

import (
	"github.com/millken/yoga"
	"github.com/sjm1327605995/tenon/react"
	"github.com/sjm1327605995/tenon/react/api/styles"
	"github.com/sjm1327605995/tenon/react/core"
	"github.com/sjm1327605995/tenon/react/elements"
	"image/color"
)

type Hello struct {
}

func NewHello() *Hello {
	return &Hello{}
}

func (h *Hello) Render() *core.VNode {
	return elements.NewView().Style(
		styles.NewStyle().
			HeightPercent(100).WidthPercent(100).
			JustifyContent(yoga.JustifyCenter).AlignItem(yoga.AlignCenter)).
		Child(
			elements.NewView().Style(
				styles.NewStyle().Width(200).Height(150).
					Border(yoga.EdgeTop, 2).Border(yoga.EdgeBottom, 4).
					Border(yoga.EdgeLeft, 6).
					Border(yoga.EdgeRight, 8).
					CornerRadius(styles.CornerRadius{
						TopLeft:     10,
						TopRight:    20,
						BottomRight: 30,
						BottomLeft:  40,
					}).
					BackgroundColor(color.NRGBA{R: 97, G: 218, B: 251, A: 255}),
			),
			elements.NewImage().
				Style(styles.NewStyle().Width(100)).
				Source("react.svg"),
			elements.NewText().
				Style(styles.NewStyle()).
				Content("hello world"),
		).Render()
}

func main() {
	dom := react.NewReactDOM()
	// This will block until the window is closed.
	dom.Render(NewHello())
}
