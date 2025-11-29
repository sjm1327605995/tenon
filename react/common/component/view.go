package component

import (
	"image/color"

	"github.com/sjm1327605995/tenon/react/yoga"
)

type View struct {
	yoga        *yoga.Node
	Background  color.NRGBA
	BorderColor color.NRGBA
}

func (v *View) Yoga() *yoga.Node {
	return v.yoga
}
func NewView() *View {
	return &View{
		yoga: yoga.NewNode(),
	}
}
