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
	background color.NRGBA

	borderColor color.NRGBA
	radius      Radius
	borderWidth unit.Dp
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
	view.Base = NewBase(view)
	return view
}

func (v *View) BorderWidth(borderWidth unit.Dp) *View {
	v.borderWidth = borderWidth
	return v
}

func (v *View) Radius(cornerRadius ...unit.Dp) *View {
	switch len(cornerRadius) {
	case 0:
		return v
	case 1:
		// All four corners
		r := cornerRadius[0]
		v.radius.NE = r
		v.radius.NW = r
		v.radius.SW = r
		v.radius.SE = r
	case 2:
		// Top-left/right, Bottom-right/left
		r1 := cornerRadius[0]
		r2 := cornerRadius[1]
		v.radius.NW = r1
		v.radius.NE = r1
		v.radius.SE = r2
		v.radius.SW = r2
	case 3:
		// Top-left, Top-right/bottom-right, Bottom-left
		v.radius.NW = cornerRadius[0]
		v.radius.NE = cornerRadius[1]
		v.radius.SE = cornerRadius[1]
		v.radius.SW = cornerRadius[2]
	case 4:
		// Top-left, Top-right, Bottom-right, Bottom-left (clockwise from top-left)
		v.radius.NW = cornerRadius[0]
		v.radius.NE = cornerRadius[1]
		v.radius.SE = cornerRadius[2]
		v.radius.SW = cornerRadius[3]
	default:
		// More than 4 values, use first 4
		v.radius.NW = cornerRadius[0]
		v.radius.NE = cornerRadius[1]
		v.radius.SE = cornerRadius[2]
		v.radius.SW = cornerRadius[3]
	}
	return v
}

func (v *View) BorderColor(color color.NRGBA) *View {
	v.borderColor = color
	return v
}
func (v *View) Background(color color.NRGBA) *View {
	v.background = color
	return v
}
func (v *View) Update(ctx layout.Context) {
	w := int(v.Node.StyleGetWidth())
	h := int(v.Node.StyleGetHeight())
	viewGio := &ViewGio{
		Size:        image.Pt(w, h),
		radiusSE:    ctx.Dp(v.radius.SE),
		radiusSW:    ctx.Dp(v.radius.SW),
		radiusNW:    ctx.Dp(v.radius.NW),
		radiusNE:    ctx.Dp(v.radius.NE),
		BorderWidth: v.borderWidth,
		BorderColor: v.borderColor,
	}
	v.Node.StyleSetBorder(yoga.EdgeAll, float32(ctx.Dp(v.borderWidth)))
	children := v.Node.GetChildren()
	viewGio.BackgroundColor = v.background
	for i := range children {
		offsetX, offsetY := children[i].LayoutLeft(), children[i].LayoutTop()
		n := children[i].GetContext().(core.Node)
		viewGio.Layouts = append(viewGio.Layouts, func(gtx layout.Context) {
			off := op.Offset(image.Pt(int(offsetX), int(offsetY))).Push(gtx.Ops)
			n.Gio().Layout(gtx)
			off.Pop()
		})
	}
	v.gio = viewGio
}

type ViewGio struct {
	Size            image.Point
	radiusSE        int
	radiusSW        int
	radiusNW        int
	radiusNE        int
	Layouts         []func(gtx layout.Context)
	BorderWidth     unit.Dp
	BorderColor     color.NRGBA
	BackgroundColor color.NRGBA
}

func (v *ViewGio) Layout(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min = v.Size
	gtx.Constraints.Max = v.Size
	return v.layout(gtx)
}

func (v *ViewGio) layout(gtx layout.Context) layout.Dimensions {

	if v.Size.X == 0 && v.Size.Y == 0 {
		return layout.Dimensions{Size: v.Size}
	}
	width := gtx.Dp(v.BorderWidth)
	whalf := (width + 1) / 2
	if v.BackgroundColor.A > 0 {
		bodySize := v.Size
		if v.BorderWidth > 0 {
			bodySize.X -= whalf
			bodySize.Y -= whalf
		}
		paint.FillShape(gtx.Ops, v.BackgroundColor, clip.Outline{
			Path: clip.RRect{
				Rect: image.Rectangle{Min: image.Pt(whalf, whalf), Max: bodySize},
				SE:   v.radiusSE,
				SW:   v.radiusSW,
				NW:   v.radiusNW,
				NE:   v.radiusNE,
			}.Path(gtx.Ops),
		}.Op())
	}
	paint.FillShape(gtx.Ops, v.BorderColor,
		clip.Stroke{
			Path: clip.RRect{
				Rect: image.Rect(whalf, whalf, v.Size.X-whalf, v.Size.Y-whalf),
				SE:   v.radiusSE,
				SW:   v.radiusSW,
				NW:   v.radiusNW,
				NE:   v.radiusNE,
			}.Path(gtx.Ops),
			Width: float32(width),
		}.Op(),
	)
	for i := range v.Layouts {
		v.Layouts[i](gtx)
	}
	return layout.Dimensions{Size: v.Size}
}
