package ui

import (
	"github.com/sjm1327605995/tenon/yoga"
)

type LayoutResults struct {
	Width  float32
	Height float32
	Left   float32
	Top    float32
}

func newLayoutResults(node *yoga.Node) LayoutResults {
	return LayoutResults{
		Width:  node.LayoutWidth(),
		Height: node.LayoutHeight(),
		Left:   node.LayoutLeft(),
		Top:    node.LayoutTop(),
	}
}

func (l LayoutResults) Equal(other LayoutResults) bool {
	return l.Width == other.Width &&
		l.Height == other.Height &&
		l.Left == other.Left &&
		l.Top == other.Top
}

func getNodeLayout(node *yoga.Node) LayoutResults {
	if node == nil {
		return LayoutResults{}
	}
	return newLayoutResults(node)
}


