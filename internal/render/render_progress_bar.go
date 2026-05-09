package render

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/yoga"
)

// RenderProgressBar is a linear progress indicator.
type RenderProgressBar struct {
	RenderBox
	Progress     float32 // 0.0 ~ 1.0
	TrackColor   color.Color
	FillColor    color.Color
	BarRadius    float32

	// TrackImage / FillImage 是游戏风格的图片背景；非 nil 时优先于颜色。
	TrackImage *ebiten.Image
	FillImage  *ebiten.Image
	Slice      BorderSlice // 9-slice 参数；全 0 时普通拉伸
}

func NewRenderProgressBar() *RenderProgressBar {
	r := &RenderProgressBar{}
	r.RenderBox.Init(r)
	r.yoga = yoga.NewNode()
	r.BarRadius = 4
	r.StyleSetHeight(8)
	r.StyleSetWidthPercent(100)
	return r
}

func (r *RenderProgressBar) SetProgress(v float32) {
	if v < 0 {
		v = 0
	}
	if v > 1 {
		v = 1
	}
	if r.Progress != v {
		r.Progress = v
		r.MarkNeedsPaint()
	}
}

func (r *RenderProgressBar) SetImages(track, fill *ebiten.Image, slice BorderSlice) {
	if r.TrackImage == track && r.FillImage == fill && r.Slice == slice {
		return
	}
	r.TrackImage = track
	r.FillImage = fill
	r.Slice = slice
	r.MarkNeedsPaint()
}

func (r *RenderProgressBar) Paint(screen *ebiten.Image, offset Offset) {
	bounds := r.bounds
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}
	x := offset.X + bounds.X
	y := offset.Y + bounds.Y
	w := bounds.Width
	h := bounds.Height

	br := UniformBorderRadius(r.BarRadius)
	br = clampRadius(br, w, h)

	// Track
	DrawRoundedRectFill(screen, x, y, w, h, br, r.TrackColor)

	// Fill
	fillW := w * r.Progress
	if fillW > w {
		fillW = w
	}
	if r.FillImage != nil && fillW > 0 {
		if r.Slice.Top > 0 || r.Slice.Right > 0 || r.Slice.Bottom > 0 || r.Slice.Left > 0 {
			ninePatch := &RenderNinePatch{Source: r.FillImage, Slice: r.Slice}
			ninePatch.Paint(screen, Offset{X: x, Y: y})
		} else {
			op := &ebiten.DrawImageOptions{}
			sw := float64(fillW) / float64(r.FillImage.Bounds().Dx())
			sh := float64(h) / float64(r.FillImage.Bounds().Dy())
			op.GeoM.Scale(sw, sh)
			op.GeoM.Translate(float64(x), float64(y))
			screen.DrawImage(r.FillImage, op)
		}
	} else if r.FillColor != nil && r.Progress > 0 {
		DrawRoundedRectFill(screen, x, y, fillW, h, br, r.FillColor)
	}
}
