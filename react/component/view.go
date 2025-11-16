package component

import (
	"github.com/sjm1327605995/tenon/react/core"
	"github.com/sjm1327605995/tenon/react/yoga"
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
)

type View struct {
	Base[View]
	background    color.NRGBA
	setBackground bool
	borderColor   color.NRGBA
	cornerRadius  unit.Dp
	borderWidth   unit.Dp
}

func NewView() *View {
	view := &View{
		borderColor: color.NRGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 0xff,
		},
	}
	view.Base = NewBase[View](view)
	return view
}

func (v *View) BorderWidth(borderWidth unit.Dp) *View {
	v.borderWidth = borderWidth
	v.Node.StyleSetBorder(yoga.EdgeAll, float32(borderWidth))
	return v
}

func (v *View) BorderRadius(cornerRadius unit.Dp) *View {
	v.cornerRadius = cornerRadius
	return v
}

func (v *View) BorderColor(color color.NRGBA) *View {
	v.borderColor = color
	return v
}
func (v *View) Background(color color.NRGBA) *View {
	v.setBackground = true
	v.background = color
	return v
}

func (v *View) Layout(gtx layout.Context) layout.Dimensions {
	if v.borderWidth > 0 {
		return widget.Border{
			Color:        v.borderColor,
			CornerRadius: v.cornerRadius,
			Width:        v.borderWidth,
		}.Layout(gtx, v.layout)
	}
	return v.layout(gtx)
}
func (v *View) layout(gtx layout.Context) layout.Dimensions {

	w := v.Node.StyleGetWidth()
	h := v.Node.StyleGetHeight()
	//border := v.Node.StyleGetBorder(yoga.EdgeAll)
	w -= float32(v.borderWidth)
	h -= float32(v.borderWidth)
	size := image.Pt(int(w), int(h))
	defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
	if v.setBackground {
		paint.ColorOp{Color: v.background}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
	}
	children := v.Node.GetChildren()

	for i := range children {
		offsetX, offsetY := children[i].LayoutLeft(), children[i].LayoutTop()
		n := children[i].GetContext().(core.Node)
		off := op.Offset(image.Pt(int(offsetX+float32(v.borderWidth)), int(offsetY+float32(v.borderWidth)))).Push(gtx.Ops)
		n.Layout(gtx)
		off.Pop()

	}
	return layout.Dimensions{Size: image.Pt(int(w), int(h))}
}
