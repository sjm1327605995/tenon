package render

import (
	"bytes"
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers/gio"
)

type Svg struct {
	canvas     *canvas.Canvas
	data       *bytes.Reader
	call       op.CallOp
	Record     bool
	style      ImageStyle
	dimensions layout.Dimensions
}

func (s *Svg) HasDefault() bool {
	return true
}

const ptPerMm = 72.0 / 25.4

func (s *Svg) DefaultSize() image.Point {
	return image.Pt(int(ptPerMm*s.canvas.W), int(ptPerMm*s.canvas.H))
}
func NewSvg(reader *bytes.Reader, style ImageStyle) *Svg {

	fc, err := canvas.ParseSVG(reader)
	if err != nil {
		panic(err)
	}
	return &Svg{
		data:   reader,
		canvas: fc,
		style:  style,
	}
}
func (s *Svg) Layout(gtx layout.Context) layout.Dimensions {
	if !s.Record {
		ops := gtx.Ops
		cache := new(op.Ops)
		gtx.Ops = cache
		macro := op.Record(gtx.Ops)
		gtx.Constraints.Min = image.Pt(0, 0)
		c := gio.NewContain(gtx, s.canvas.W, s.canvas.H)
		s.canvas.RenderTo(c)
		s.call = macro.Stop()
		gtx.Ops = ops
		s.Record = true
		s.dimensions = c.Dimensions()
	}
	s.call.Add(gtx.Ops)
	return s.dimensions
}
