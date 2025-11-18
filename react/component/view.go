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
	"math"
)

type View struct {
	Base[View]
	background    color.NRGBA
	setBackground bool
	borderColor   color.NRGBA
	radius        Radius
	border        Border
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

func (v *View) BorderWidth(borderWidth ...unit.Dp) *View {
	switch len(borderWidth) {
	case 0:
		return v
	case 1:
		// All four sides
		w := borderWidth[0]
		v.border.TopWidth = w
		v.border.RightWidth = w
		v.border.BottomWidth = w
		v.border.LeftWidth = w
	case 2:
		// Top/Bottom, Left/Right
		v.border.TopWidth = borderWidth[0]
		v.border.BottomWidth = borderWidth[0]
		v.border.LeftWidth = borderWidth[1]
		v.border.RightWidth = borderWidth[1]
	case 3:
		// Top, Left/Right, Bottom
		v.border.TopWidth = borderWidth[0]
		v.border.LeftWidth = borderWidth[1]
		v.border.RightWidth = borderWidth[1]
		v.border.BottomWidth = borderWidth[2]
	case 4:
		// Top, Right, Bottom, Left (clockwise from top)
		v.border.TopWidth = borderWidth[0]
		v.border.RightWidth = borderWidth[1]
		v.border.BottomWidth = borderWidth[2]
		v.border.LeftWidth = borderWidth[3]
	default:
		// More than 4 values, use first 4
		v.border.TopWidth = borderWidth[0]
		v.border.RightWidth = borderWidth[1]
		v.border.BottomWidth = borderWidth[2]
		v.border.LeftWidth = borderWidth[3]
	}
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
	v.setBackground = true
	v.background = color
	return v
}
func (v *View) Update(ctx layout.Context) {
	w := int(v.Node.StyleGetWidth())
	h := int(v.Node.StyleGetHeight())
	viewGio := &ViewGio{
		W:           w,
		H:           h,
		radiusSE:    ctx.Dp(v.radius.SE),
		radiusSW:    ctx.Dp(v.radius.SW),
		radiusNW:    ctx.Dp(v.radius.NW),
		radiusNE:    ctx.Dp(v.radius.NE),
		Border:      v.border,
		BorderColor: v.borderColor,
	}
	v.Node.StyleSetBorder(yoga.EdgeLeft, float32(ctx.Dp(v.border.LeftWidth)))
	v.Node.StyleSetBorder(yoga.EdgeTop, float32(ctx.Dp(v.border.TopWidth)))
	v.Node.StyleSetBorder(yoga.EdgeBottom, float32(ctx.Dp(v.border.BottomWidth)))
	v.Node.StyleSetBorder(yoga.EdgeRight, float32(ctx.Dp(v.border.RightWidth)))
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
			off := op.Offset(image.Pt(int(offsetX), int(offsetY))).Push(gtx.Ops)
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
	Border
	BorderColor color.NRGBA
}
type Border struct {
	LeftWidth   unit.Dp
	RightWidth  unit.Dp
	BottomWidth unit.Dp
	TopWidth    unit.Dp
}

func (v *ViewGio) Layout(gtx layout.Context) layout.Dimensions {
	return v.layout(gtx)
}
func (v *ViewGio) layout(gtx layout.Context) layout.Dimensions {
	size := image.Pt(v.W, v.H)

	defer clip.RRect{Rect: image.Rectangle{
		Max: size,
	},
		SE: v.radiusSE,
		SW: v.radiusSW,
		NW: v.radiusNW,
		NE: v.radiusNE,
	}.Push(gtx.Ops).Pop()

	// Draw inner border
	leftWidthPx := gtx.Dp(v.Border.LeftWidth)
	rightWidthPx := gtx.Dp(v.Border.RightWidth)
	topWidthPx := gtx.Dp(v.Border.TopWidth)
	bottomWidthPx := gtx.Dp(v.Border.BottomWidth)

	// Only draw border if any width is non-zero
	if leftWidthPx > 0 || rightWidthPx > 0 || topWidthPx > 0 || bottomWidthPx > 0 {
		// Calculate the minimum border width to use for all sides
		// This ensures consistent border width around the view
		minBorderWidth := leftWidthPx
		if rightWidthPx < minBorderWidth {
			minBorderWidth = rightWidthPx
		}
		if topWidthPx < minBorderWidth {
			minBorderWidth = topWidthPx
		}
		if bottomWidthPx < minBorderWidth {
			minBorderWidth = bottomWidthPx
		}

		// Ensure we have at least a 1px border
		if minBorderWidth < 1 {
			minBorderWidth = 1
		}

		// Create outer RRect with the view's full size
		outerRRect := clip.RRect{
			Rect: image.Rect(0, 0, v.W, v.H),
			SE:   v.radiusSE,
			SW:   v.radiusSW,
			NW:   v.radiusNW,
			NE:   v.radiusNE,
		}

		// Create inner RRect with border width subtracted from all sides
		innerLeft := minBorderWidth
		innerTop := minBorderWidth
		innerRight := v.W - minBorderWidth
		innerBottom := v.H - minBorderWidth

		// Ensure inner rectangle is valid
		if innerLeft >= innerRight {
			innerRight = innerLeft + 1
		}
		if innerTop >= innerBottom {
			innerBottom = innerTop + 1
		}

		// Calculate inner radii (outer radii minus border width)
		innerRadius := max(int(math.Max(float64(v.radiusSE-minBorderWidth), 0.0)), 0)

		innerRRect := clip.RRect{
			Rect: image.Rect(innerLeft, innerTop, innerRight, innerBottom),
			SE:   innerRadius,
			SW:   innerRadius,
			NW:   innerRadius,
			NE:   innerRadius,
		}

		// Draw outer RRect with border color
		outerClip := outerRRect.Push(gtx.Ops)
		paint.ColorOp{Color: v.BorderColor}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		outerClip.Pop()

		// Draw inner RRect with background color (this creates the border effect)
		// Note: We're assuming the background color is solid here
		// In a real implementation, we would need to get the actual background color
		innerClip := innerRRect.Push(gtx.Ops)
		bgColor := color.NRGBA{255, 255, 255, 255} // Default to white background
		paint.ColorOp{Color: bgColor}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		innerClip.Pop()
	}

	for i := range v.Layouts {
		v.Layouts[i](gtx)
	}
	return layout.Dimensions{Size: size}
}
