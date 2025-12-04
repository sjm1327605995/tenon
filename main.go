package main

import (
	"image/color"

	"github.com/sjm1327605995/tenon/react"
	"github.com/sjm1327605995/tenon/react/api/styles"
	"github.com/sjm1327605995/tenon/react/core"
	"github.com/sjm1327605995/tenon/react/elements"
	"github.com/sjm1327605995/tenon/react/yoga"
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
					Border(yoga.EdgeTop, 10).Border(yoga.EdgeBottom, 20).
					Border(yoga.EdgeLeft, 20).
					Border(yoga.EdgeRight, 30).
					CornerRadius(styles.CornerRadius{
						TopLeft:     10,
						TopRight:    20,
						BottomRight: 30,
						BottomLeft:  40,
					}).
					BackgroundColor(color.NRGBA{R: 255, A: 255}),
			),
			elements.NewImage().
				Style(styles.NewStyle().WidthPercent(50).HeightPercent(50)).
				Source("react.svg"),
			elements.NewText().
				Style(styles.NewStyle().Width(100)).
				Content("你好世界"),
		).Render()
}

func main() {
	dom := react.NewReactDOM()
	// This will block until the window is closed.
	dom.Render(NewHello())
}
