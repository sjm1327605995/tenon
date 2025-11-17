package component

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/sjm1327605995/tenon/react/core"
	"github.com/sjm1327605995/tenon/react/yoga"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
)

type View struct {
	Base[View]
	background    color.NRGBA
	setBackground bool
	borderColor   color.NRGBA
	radius        Radius
	borderWidth   unit.Dp
}
type Radius struct {
	SE, SW, NW, NE unit.Dp
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

func (v *View) Radius(cornerRadius ...unit.Dp) *View {
	if len(cornerRadius) == 0 {
		return v
	}
	if len(cornerRadius) == 1 {
		v.radius.NE = cornerRadius[3]
		v.radius.NW = cornerRadius[2]
		v.radius.SW = cornerRadius[1]
		v.radius.SE = cornerRadius[0]
		return v
	}
	switch len(cornerRadius) {
	case 4:
		v.radius.NE = cornerRadius[3]
		fallthrough
	case 3:
		v.radius.NW = cornerRadius[2]
		fallthrough
	case 2:
		v.radius.SW = cornerRadius[1]
		fallthrough
	case 1:
		v.radius.SE = cornerRadius[0]
	}
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
func (v *View) Update(ctx layout.Context) {
	w := int(v.Node.StyleGetWidth())
	h := int(v.Node.StyleGetHeight())
	viewGio := &ViewGio{
		W:        w,
		H:        h,
		radiusSE: ctx.Dp(v.radius.SE),
		radiusSW: ctx.Dp(v.radius.SW),
		radiusNW: ctx.Dp(v.radius.NW),
		radiusNE: ctx.Dp(v.radius.NE),
	}
	if v.setBackground {
		viewGio.Layouts = append(viewGio.Layouts, func(gtx layout.Context) {
			paint.ColorOp{Color: v.background}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
		})
	}
	children := v.Node.GetChildren()

	for i := range children {
		offsetX, offsetY := children[i].LayoutLeft(), children[i].LayoutTop()
		n := children[i].GetContext().(core.Node)
		viewGio.Layouts = append(viewGio.Layouts, func(gtx layout.Context) {
			off := op.Offset(image.Pt(int(offsetX+float32(v.borderWidth)), int(offsetY+float32(v.borderWidth)))).Push(gtx.Ops)
			n.Gio().Layout(gtx)
			off.Pop()
		})

	}
	v.gio = viewGio
}

type ViewGio struct {
	W        int
	H        int
	radiusSE int
	radiusSW int
	radiusNW int
	radiusNE int
	Layouts  []func(gtx layout.Context)
}

func (v *ViewGio) Layout(gtx layout.Context) layout.Dimensions {
	//if v.borderWidth > 0 {
	//	return widget.Border{
	//		Color:        v.borderColor,
	//		CornerRadius: v.cornerRadius,
	//		Width:        v.borderWidth,
	//	}.Layout(gtx, v.layout)
	//}
	return v.layout(gtx)
}
func (v *ViewGio) layout(gtx layout.Context) layout.Dimensions {

	//border := v.Node.StyleGetBorder(yoga.EdgeAll)
	//w -= float32(v.borderWidth)
	//h -= float32(v.borderWidth)

	size := image.Pt(v.W, v.H)

	defer clip.RRect{Rect: image.Rectangle{
		Max: size,
	},
		SE: v.radiusSE,
		SW: v.radiusSW,
		NW: v.radiusNW,
		NE: v.radiusNE,
	}.Push(gtx.Ops).Pop()
	for i := range v.Layouts {
		v.Layouts[i](gtx)
	}
	return layout.Dimensions{Size: size}
}
