package render

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/yoga"
)

// ObjectFit 定义图片在目标矩形内的填充方式。
type ObjectFit int

const (
	ObjectFitCover ObjectFit = iota
	ObjectFitContain
	ObjectFitFill
	ObjectFitNone
	ObjectFitScaleDown
)

// RenderImage 负责图片的 Yoga 测量和 Ebiten 绘制。
type RenderImage struct {
	BaseRenderObject

	Source       *ebiten.Image
	ObjectFit    ObjectFit
	BorderRadius BorderRadius
	TintColor    color.Color
}

func NewRenderImage() *RenderImage {
	r := &RenderImage{
		ObjectFit: ObjectFitCover,
	}
	r.BaseRenderObject.Init(r)
	r.yoga = yoga.NewNode()
	r.yoga.SetMeasureFunc(r.measure)
	return r
}

func (r *RenderImage) measure(
	node *yoga.Node,
	width float32,
	widthMode yoga.MeasureMode,
	height float32,
	heightMode yoga.MeasureMode,
) yoga.Size {
	if r.Source == nil {
		return yoga.Size{Width: 0, Height: 0}
	}
	bounds := r.Source.Bounds()
	imgW := float32(bounds.Dx())
	imgH := float32(bounds.Dy())

	// 宽高都有约束，直接返回约束尺寸
	if widthMode != yoga.MeasureModeUndefined && heightMode != yoga.MeasureModeUndefined {
		return yoga.Size{Width: width, Height: height}
	}

	// 只有宽度约束，按图片比例计算高度
	if widthMode != yoga.MeasureModeUndefined && heightMode == yoga.MeasureModeUndefined {
		if imgW > 0 {
			return yoga.Size{Width: width, Height: width * imgH / imgW}
		}
		return yoga.Size{Width: width, Height: imgH}
	}

	// 只有高度约束，按图片比例计算宽度
	if widthMode == yoga.MeasureModeUndefined && heightMode != yoga.MeasureModeUndefined {
		if imgH > 0 {
			return yoga.Size{Width: height * imgW / imgH, Height: height}
		}
		return yoga.Size{Width: imgW, Height: height}
	}

	// 无约束，返回原始尺寸
	return yoga.Size{Width: imgW, Height: imgH}
}

func (r *RenderImage) Paint(screen *ebiten.Image, offset Offset) {
	if r.Source == nil {
		return
	}
	bounds := r.bounds
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	x := float64(offset.X + bounds.X)
	y := float64(offset.Y + bounds.Y)
	w := float64(bounds.Width)
	h := float64(bounds.Height)

	srcBounds := r.Source.Bounds()
	srcW := float64(srcBounds.Dx())
	srcH := float64(srcBounds.Dy())

	geoM := buildObjectFitGeoM(srcW, srcH, w, h, r.ObjectFit)

	// 有圆角时，走 mask 裁剪流程
	if !r.BorderRadius.IsZero() {
		r.paintWithClip(screen, x, y, w, h, geoM)
		return
	}

	// 无圆角时，直接绘制到 SubImage 防止溢出
	sub := SubImage(screen, int(x), int(y), int(w), int(h))
	if sub == nil {
		return
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM = geoM
	if r.TintColor != nil {
		op.ColorScale.ScaleWithColor(r.TintColor)
	}
	sub.DrawImage(r.Source, op)
}

func (r *RenderImage) paintWithClip(screen *ebiten.Image, x, y, w, h float64, geoM ebiten.GeoM) {
	iw, ih := int(w), int(h)
	if iw <= 0 || ih <= 0 {
		return
	}

	// 1. mask：白色圆角矩形（alpha=1），其余透明
	mask := ebiten.NewImage(iw, ih)
	defer mask.Dispose()

	path := &vector.Path{}
	BuildRoundedRectPath(path, 0, 0, float32(w), float32(h), r.BorderRadius)
	vector.FillPath(mask, path, &vector.FillOptions{}, &vector.DrawPathOptions{
		ColorScale: toColorScale(color.White),
		AntiAlias:  true,
	})

	// 2. 用 blend 把原图按 mask 的 alpha 直接裁剪到 mask 上
	//    source = 原图, destination = mask
	//    out = source * destAlpha + dest * 0
	op := &ebiten.DrawImageOptions{}
	op.GeoM = geoM
	if r.TintColor != nil {
		op.ColorScale.ScaleWithColor(r.TintColor)
	}
	op.Blend = ebiten.Blend{
		BlendFactorSourceRGB:        ebiten.BlendFactorDestinationAlpha,
		BlendFactorSourceAlpha:      ebiten.BlendFactorDestinationAlpha,
		BlendFactorDestinationRGB:   ebiten.BlendFactorZero,
		BlendFactorDestinationAlpha: ebiten.BlendFactorZero,
		BlendOperationRGB:           ebiten.BlendOperationAdd,
		BlendOperationAlpha:         ebiten.BlendOperationAdd,
	}
	mask.DrawImage(r.Source, op)

	// 3. 把裁剪后的结果画到 screen
	screenOp := &ebiten.DrawImageOptions{}
	screenOp.GeoM.Translate(x, y)
	screen.DrawImage(mask, screenOp)
}

// buildObjectFitGeoM 计算图片按 ObjectFit 缩放/位移后的 GeoM。
// 返回的 GeoM 是相对于目标矩形左上角的变换。
func buildObjectFitGeoM(imgW, imgH, targetW, targetH float64, fit ObjectFit) ebiten.GeoM {
	if imgW <= 0 || imgH <= 0 {
		return ebiten.GeoM{}
	}

	var sx, sy, dx, dy float64

	switch fit {
	case ObjectFitFill:
		sx = targetW / imgW
		sy = targetH / imgH
	case ObjectFitContain:
		scale := math.Min(targetW/imgW, targetH/imgH)
		sx, sy = scale, scale
		dx = (targetW - imgW*scale) / 2
		dy = (targetH - imgH*scale) / 2
	case ObjectFitCover:
		scale := math.Max(targetW/imgW, targetH/imgH)
		sx, sy = scale, scale
		dx = (targetW - imgW*scale) / 2
		dy = (targetH - imgH*scale) / 2
	case ObjectFitNone:
		dx = (targetW - imgW) / 2
		dy = (targetH - imgH) / 2
	case ObjectFitScaleDown:
		scale := math.Min(targetW/imgW, targetH/imgH)
		if scale > 1 {
			scale = 1
		}
		sx, sy = scale, scale
		dx = (targetW - imgW*scale) / 2
		dy = (targetH - imgH*scale) / 2
	}

	var g ebiten.GeoM
	g.Scale(sx, sy)
	g.Translate(dx, dy)
	return g
}

func (r *RenderImage) SetSource(src *ebiten.Image) {
	if r.Source == src {
		return
	}
	r.Source = src
	r.yoga.MarkDirty()
	r.MarkNeedsLayout()
}

func (r *RenderImage) SetObjectFit(fit ObjectFit) {
	if r.ObjectFit == fit {
		return
	}
	r.ObjectFit = fit
	r.MarkNeedsLayout()
}

func (r *RenderImage) SetBorderRadius(radius float32) {
	br := UniformBorderRadius(radius)
	if r.BorderRadius == br {
		return
	}
	r.BorderRadius = br
	r.MarkNeedsPaint()
}

func (r *RenderImage) SetTintColor(c color.Color) {
	if ColorEquals(r.TintColor, c) {
		return
	}
	r.TintColor = c
	r.MarkNeedsPaint()
}

func (r *RenderImage) SetWidth(v float32) {
	if v > 0 {
		r.yoga.StyleSetWidth(v)
	} else {
		r.yoga.StyleSetWidthAuto()
	}
	r.MarkNeedsLayout()
}

func (r *RenderImage) SetHeight(v float32) {
	if v > 0 {
		r.yoga.StyleSetHeight(v)
	} else {
		r.yoga.StyleSetHeightAuto()
	}
	r.MarkNeedsLayout()
}
