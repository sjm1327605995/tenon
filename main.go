package main

import (
	"github.com/sjm1327605995/tenon/react"
	"github.com/sjm1327605995/tenon/react/common"
	"github.com/sjm1327605995/tenon/react/core"
	"github.com/sjm1327605995/tenon/react/style"
	"github.com/sjm1327605995/tenon/react/style/unit"
)

type Hello struct {
}

func NewHello() *Hello {
	return &Hello{}
}

func (h *Hello) Render() common.Node {
	return core.NewView().
		Style(style.Height(unit.Percent(50)), style.Width(unit.Percent(100)))

}

func main() {
	err := react.NewReactDOM().
		Render(NewHello())
	if err != nil {
		panic(err)
	}
}
