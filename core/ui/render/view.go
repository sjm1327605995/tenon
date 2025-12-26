package render

import (
	"image"
	"image/color"
	"math"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type ViewStyle struct {
	Top         float32
	Right       float32
	Bottom      float32
	Left        float32
	Background  color.NRGBA
	BorderColor color.NRGBA
	CornerRadii CornerRadius
}

type CornerRadius struct {
	TopLeft, TopRight, BottomRight, BottomLeft float32
}

func (v ViewStyle) ToRender() Render {
	return &ViewRender{ViewStyle: v}
}

type ViewRender struct {
	ViewStyle
}

func (v *ViewRender) DefaultSize() image.Point {
	return image.Point{}
}

func (v *ViewRender) HasDefault() bool {
	return false
}

func (v *ViewRender) Layout(ctx layout.Context) layout.Dimensions {

	width := float32(ctx.Constraints.Max.X)
	height := float32(ctx.Constraints.Max.Y)

	if v.Background.A > 0 {
		var p clip.Path
		p.Begin(ctx.Ops)
		drawInnerLoop(&p, width, height, v.Top, v.Right, v.Bottom, v.Left, v.CornerRadii, false)
		paint.FillShape(ctx.Ops, v.Background, clip.Outline{Path: p.End()}.Op())
	}

	if v.BorderColor.A > 0 && (v.Top > 0 || v.Right > 0 || v.Bottom > 0 || v.Left > 0) {
		var p clip.Path
		p.Begin(ctx.Ops)
		drawOuterLoop(&p, width, height, v.CornerRadii)
		drawInnerLoop(&p, width, height, v.Top, v.Right, v.Bottom, v.Left, v.CornerRadii, true)
		paint.FillShape(ctx.Ops, v.BorderColor, clip.Outline{Path: p.End()}.Op())
	}
	return layout.Dimensions{Size: ctx.Constraints.Max}
}

const q = 4 * (math.Sqrt2 - 1) / 3
const k = 1 - q

func drawOuterLoop(p *clip.Path, w, h float32, r CornerRadius) {
	tl, tr, br, bl := r.TopLeft, r.TopRight, r.BottomRight, r.BottomLeft
	p.MoveTo(f32.Pt(tl, 0))
	p.LineTo(f32.Pt(w-tr, 0))
	if tr > 0 {
		p.CubeTo(f32.Pt(w-tr*(1-k), 0), f32.Pt(w, tr*(1-k)), f32.Pt(w, tr))
	}
	p.LineTo(f32.Pt(w, h-br))
	if br > 0 {
		p.CubeTo(f32.Pt(w, h-br*(1-k)), f32.Pt(w-br*(1-k), h), f32.Pt(w-br, h))
	}
	p.LineTo(f32.Pt(bl, h))
	if bl > 0 {
		p.CubeTo(f32.Pt(bl*(1-k), h), f32.Pt(0, h-bl*(1-k)), f32.Pt(0, h-bl))
	}

	p.LineTo(f32.Pt(0, tl))
	if tl > 0 {
		p.CubeTo(f32.Pt(0, tl*(1-k)), f32.Pt(tl*(1-k), 0), f32.Pt(tl, 0))
	}
	p.Close()
}

func drawInnerLoop(p *clip.Path, w, h, top, right, bottom, left float32, r CornerRadius, reverse bool) {
	tlRx, tlRy := max(0, r.TopLeft-left), max(0, r.TopLeft-top)
	trRx, trRy := max(0, r.TopRight-right), max(0, r.TopRight-top)
	brRx, brRy := max(0, r.BottomRight-right), max(0, r.BottomRight-bottom)
	blRx, blRy := max(0, r.BottomLeft-left), max(0, r.BottomLeft-bottom)

	ptTlStart := f32.Pt(left, top+tlRy)
	ptTlEnd := f32.Pt(left+tlRx, top)
	ptTrStart := f32.Pt(w-right-trRx, top)
	ptTrEnd := f32.Pt(w-right, top+trRy)
	ptBrStart := f32.Pt(w-right, h-bottom-brRy)
	ptBrEnd := f32.Pt(w-right-brRx, h-bottom)
	ptBlStart := f32.Pt(left+blRx, h-bottom)
	ptBlEnd := f32.Pt(left, h-bottom-blRy)

	if !reverse {
		p.MoveTo(ptTlEnd)
		p.LineTo(ptTrStart)
		if trRx > 0 || trRy > 0 {
			p.CubeTo(f32.Pt(w-right-trRx+trRx*k, top), f32.Pt(w-right, top+trRy-trRy*k), ptTrEnd)
		} else {
			p.LineTo(ptTrEnd)
		}
		p.LineTo(ptBrStart)
		if brRx > 0 || brRy > 0 {
			p.CubeTo(f32.Pt(w-right, h-bottom-brRy+brRy*k), f32.Pt(w-right-brRx+brRx*k, h-bottom), ptBrEnd)
		} else {
			p.LineTo(ptBrEnd)
		}
		p.LineTo(ptBlStart)
		if blRx > 0 || blRy > 0 {
			p.CubeTo(f32.Pt(left+blRx-blRx*k, h-bottom), f32.Pt(left, h-bottom-blRy+blRy*k), ptBlEnd)
		} else {
			p.LineTo(ptBlEnd)
		}
		p.LineTo(ptTlStart)
		if tlRx > 0 || tlRy > 0 {
			p.CubeTo(f32.Pt(left, top+tlRy-tlRy*k), f32.Pt(left+tlRx-tlRx*k, top), ptTlEnd)
		} else {
			p.LineTo(ptTlEnd)
		}
		p.Close()
	} else {
		p.MoveTo(ptTlEnd)
		if tlRx > 0 || tlRy > 0 {
			p.CubeTo(f32.Pt(left+tlRx-tlRx*k, top), f32.Pt(left, top+tlRy-tlRy*k), ptTlStart)
		} else {
			p.LineTo(ptTlStart)
		}
		p.LineTo(ptBlEnd)
		if blRx > 0 || blRy > 0 {
			p.CubeTo(f32.Pt(left, h-bottom-blRy+blRy*k), f32.Pt(left+blRx-blRx*k, h-bottom), ptBlStart)
		} else {
			p.LineTo(ptBlStart)
		}
		p.LineTo(ptBrEnd)
		if brRx > 0 || brRy > 0 {
			p.CubeTo(f32.Pt(w-right-brRx+brRx*k, h-bottom), f32.Pt(w-right, h-bottom-brRy+brRy*k), ptBrStart)
		} else {
			p.LineTo(ptBrStart)
		}
		p.LineTo(ptTrEnd)
		if trRx > 0 || trRy > 0 {
			p.CubeTo(f32.Pt(w-right, top+trRy-trRy*k), f32.Pt(w-right-trRx+trRx*k, top), ptTrStart)
		} else {
			p.LineTo(ptTrStart)
		}
		p.LineTo(ptTlEnd)
		p.Close()
	}
}
