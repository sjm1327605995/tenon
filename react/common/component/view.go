package component

import "github.com/sjm1327605995/tenon/react/yoga"

type View struct {
	yoga *yoga.Node
}

func (v *View) Yoga() *yoga.Node {
	return v.yoga
}
func NewView() *View {
	return &View{
		yoga: yoga.NewNode(),
	}
}
