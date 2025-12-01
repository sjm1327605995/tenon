package components

import (
	"github.com/sjm1327605995/tenon/react/yoga"
)

type Text struct {
	yoga    *yoga.Node // Yoga layout node for handling Flexbox layout
	Content string
}

// Yoga returns the View component's Yoga layout node.
// This method enables the View to implement the api.StyleElement interface, allowing styles to be applied.
func (t *Text) Yoga() *yoga.Node {
	return t.yoga
}

func NewText() *Text {
	return &Text{
		yoga: yoga.NewNode(),
	}
}
