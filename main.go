package main

import (
	"github.com/sjm1327605995/tenon/react"
	"github.com/sjm1327605995/tenon/react/api"
	"github.com/sjm1327605995/tenon/react/api/styles"
	"github.com/sjm1327605995/tenon/react/elements"
	"github.com/sjm1327605995/tenon/react/yoga"
	"image/color"
)

type Hello struct {
}

func NewHello() *Hello {
	return &Hello{}
}

func (h *Hello) Render() api.Node {
	return elements.NewView().Style(
		styles.NewStyle().
			BackgroundColor(color.NRGBA{G: 255, A: 255}).HeightPercent(100).WidthPercent(100).
			JustifyContent(yoga.JustifyCenter).AlignItem(yoga.AlignCenter)).
		Child(elements.NewImage().
			Style(styles.NewStyle().WidthPercent(50).HeightPercent(50)).
			Source("react.svg"), elements.NewText().Content("hello"))

}

func main() {
	dom := react.NewReactDOM()
	// 渲染组件
	err := dom.Render(NewHello())
	if err != nil {
		panic(err)
	}
}
