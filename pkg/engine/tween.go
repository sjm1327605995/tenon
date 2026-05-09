package engine

import (
	"image/color"

	"github.com/sjm1327605995/tenon/pkg/render"
)

// OffsetTween 对二维偏移量进行线性插值。
type OffsetTween struct {
	Begin, End render.Offset
}

func (t *OffsetTween) Evaluate(progress float64) render.Offset {
	p := float32(progress)
	return render.Offset{
		X: t.Begin.X + (t.End.X-t.Begin.X)*p,
		Y: t.Begin.Y + (t.End.Y-t.Begin.Y)*p,
	}
}

// SizeTween 对二维尺寸进行线性插值。
type SizeTween struct {
	Begin, End render.Size
}

func (t *SizeTween) Evaluate(progress float64) render.Size {
	p := float32(progress)
	return render.Size{
		Width:  t.Begin.Width + (t.End.Width-t.Begin.Width)*p,
		Height: t.Begin.Height + (t.End.Height-t.Begin.Height)*p,
	}
}

// ColorTween 对 RGBA 颜色进行线性插值。
type ColorTween struct {
	Begin, End color.Color
}

func (t *ColorTween) Evaluate(progress float64) color.Color {
	r1, g1, b1, a1 := t.Begin.RGBA()
	r2, g2, b2, a2 := t.End.RGBA()
	p := float32(progress)
	return color.RGBA{
		R: lerpUint8(uint8(r1>>8), uint8(r2>>8), p),
		G: lerpUint8(uint8(g1>>8), uint8(g2>>8), p),
		B: lerpUint8(uint8(b1>>8), uint8(b2>>8), p),
		A: lerpUint8(uint8(a1>>8), uint8(a2>>8), p),
	}
}

func lerpUint8(a, b uint8, p float32) uint8 {
	return uint8(float32(a) + (float32(b)-float32(a))*p)
}

// TransformTween 对 2D 仿射变换参数进行线性插值。
type TransformTween struct {
	Begin, End render.Transform
}

func (t *TransformTween) Evaluate(progress float64) render.Transform {
	p := float32(progress)
	return render.Transform{
		Rotation: t.Begin.Rotation + (t.End.Rotation-t.Begin.Rotation)*p,
		ScaleX:   t.Begin.ScaleX + (t.End.ScaleX-t.Begin.ScaleX)*p,
		ScaleY:   t.Begin.ScaleY + (t.End.ScaleY-t.Begin.ScaleY)*p,
		SkewX:    t.Begin.SkewX + (t.End.SkewX-t.Begin.SkewX)*p,
		SkewY:    t.Begin.SkewY + (t.End.SkewY-t.Begin.SkewY)*p,
		OriginX:  t.Begin.OriginX + (t.End.OriginX-t.Begin.OriginX)*p,
		OriginY:  t.Begin.OriginY + (t.End.OriginY-t.Begin.OriginY)*p,
		Alpha:    t.Begin.Alpha + (t.End.Alpha-t.Begin.Alpha)*p,
	}
}

// Float32Tween 是 float32 专用插值器，避免泛型开销。
type Float32Tween struct {
	Begin, End float32
}

func (t *Float32Tween) Evaluate(progress float64) float32 {
	return t.Begin + (t.End-t.Begin)*float32(progress)
}
