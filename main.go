package main

import (
	"image/color"

	"github.com/sjm1327605995/tenon/react"
	"github.com/sjm1327605995/tenon/react/common"
	"github.com/sjm1327605995/tenon/react/core"
	"github.com/sjm1327605995/tenon/react/style"
	"github.com/sjm1327605995/tenon/react/yoga"
)

type Hello struct {
}

func NewHello() *Hello {
	return &Hello{}
}

func (h *Hello) Render() common.Node {
	return core.NewView().
		Style(style.NewStyle().
			BackgroundColor(color.NRGBA{R: 255, A: 255}).
			Direction(yoga.DirectionInherit).WidthPercent(50).HeightPercent(50))

}

func main() {
	err := react.NewReactDOM().
		Render(NewHello())
	if err != nil {
		panic(err)
	}
}
