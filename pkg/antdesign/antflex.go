package antdesign

import (
	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/components"
	"github.com/sjm1327605995/tenon/yoga"
)

// AntFlex is a flex layout container.
type AntFlex struct {
	tenon.BaseWidget
	vertical bool
	justify  yoga.Justify
	align    yoga.Align
	gap      float32
	wrap     bool
	children []tenon.Component
}

// NewAntFlex creates an AntFlex.
func NewAntFlex() *AntFlex {
	f := &AntFlex{
		justify: yoga.JustifyFlexStart,
		align:   yoga.AlignFlexStart,
	}
	f.Init(f)
	return f
}

// Render returns the flex UI.
func (f *AntFlex) Render() tenon.Component {
	dir := yoga.FlexDirectionRow
	if f.vertical {
		dir = yoga.FlexDirectionColumn
	}
	root := components.NewView().
		SetFlexDirection(dir).
		SetJustifyContent(f.justify).
		SetAlignItems(f.align)
	if f.wrap {
		root.SetFlexWrap(yoga.WrapWrap)
	}

	for i, child := range f.children {
		if i > 0 && f.gap > 0 {
			if f.vertical {
				root.Add(components.NewView().SetHeight(f.gap))
			} else {
				root.Add(components.NewView().SetWidth(f.gap))
			}
		}
		root.AddChild(child)
	}
	return root
}

func (f *AntFlex) Add(children ...tenon.Component) *AntFlex {
	f.children = append(f.children, children...)
	return f
}

func (f *AntFlex) SetVertical(v bool) *AntFlex        { f.vertical = v; return f }
func (f *AntFlex) SetJustify(j yoga.Justify) *AntFlex { f.justify = j; return f }
func (f *AntFlex) SetAlign(a yoga.Align) *AntFlex     { f.align = a; return f }
func (f *AntFlex) SetGap(g float32) *AntFlex          { f.gap = g; return f }
func (f *AntFlex) SetWrap(w bool) *AntFlex            { f.wrap = w; return f }
