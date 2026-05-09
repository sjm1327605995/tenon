package render

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/svg"
	"github.com/sjm1327605995/tenon/yoga"
)

const lucideViewBox = 24

// RenderSvgIcon 使用 SVG path 数据绘制矢量图标（描边风格）。
type RenderSvgIcon struct {
	BaseRenderObject
	pathData    string
	iconColor   color.Color
	size        float32
	strokeWidth float32

	// 缓存解析后的 path，避免每帧重复解析
	parsedPath *vector.Path
	pathDirty  bool
}

// NewRenderSvgIcon 创建 SVG 图标渲染对象。
func NewRenderSvgIcon() *RenderSvgIcon {
	r := &RenderSvgIcon{
		size:        24,
		strokeWidth: 2,
		pathDirty:   true,
	}
	r.BaseRenderObject.Init(r)
	r.yoga = yoga.NewNode()
	r.yoga.SetMeasureFunc(r.measure)
	return r
}

func (r *RenderSvgIcon) measure(
	node *yoga.Node,
	width float32,
	widthMode yoga.MeasureMode,
	height float32,
	heightMode yoga.MeasureMode,
) yoga.Size {
	return yoga.Size{Width: r.size, Height: r.size}
}

// SetPathData 设置 SVG path 数据。
func (r *RenderSvgIcon) SetPathData(d string) {
	if r.pathData == d {
		return
	}
	r.pathData = d
	r.pathDirty = true
	r.MarkNeedsPaint()
}

// SetIconColor 设置图标描边颜色。
func (r *RenderSvgIcon) SetIconColor(c color.Color) {
	if r.iconColor == c {
		return
	}
	r.iconColor = c
	r.MarkNeedsPaint()
}

// SetIconSize 设置图标尺寸。
func (r *RenderSvgIcon) SetIconSize(v float32) {
	if r.size == v {
		return
	}
	r.size = v
	r.pathDirty = true
	r.MarkNeedsLayout()
	r.MarkNeedsPaint()
}

func (r *RenderSvgIcon) ensurePath() {
	if !r.pathDirty || r.pathData == "" {
		return
	}
	// 缓存未缩放的 path，在 Paint 时根据实际 bounds 动态缩放，
	// 避免 Yoga 布局压缩后 path 被截断。
	path, err := svg.ParsePathScaledAndShifted(r.pathData, 1, 0, 0)
	if err != nil {
		r.parsedPath = nil
	} else {
		r.parsedPath = path
	}
	r.pathDirty = false
}

// Paint 绘制 SVG 图标。
func (r *RenderSvgIcon) Paint(screen *ebiten.Image, offset Offset) {
	r.ensurePath()
	if r.parsedPath == nil {
		return
	}

	bounds := r.bounds
	// 使用 bounds 的较小边作为缩放基准，确保 path 不超出 bounds
	scale := bounds.Width / lucideViewBox
	if hScale := bounds.Height / lucideViewBox; hScale < scale {
		scale = hScale
	}

	x := offset.X + bounds.X
	y := offset.Y + bounds.Y

	drawPath := &vector.Path{}
	gm := ebiten.GeoM{}
	gm.Scale(float64(scale), float64(scale))
	gm.Translate(float64(x), float64(y))
	drawPath.AddPath(r.parsedPath, &vector.AddPathOptions{GeoM: gm})

	op := &vector.DrawPathOptions{}
	if r.iconColor != nil {
		op.ColorScale.ScaleWithColor(r.iconColor)
	}
	op.AntiAlias = true

	strokeOp := &vector.StrokeOptions{
		Width:    r.strokeWidth * scale,
		LineCap:  vector.LineCapRound,
		LineJoin: vector.LineJoinRound,
	}

	vector.StrokePath(screen, drawPath, strokeOp, op)
}

// HitTest 始终返回 false，让事件穿透到父节点或兄弟节点。
func (r *RenderSvgIcon) HitTest(x, y float32) bool {
	return false
}
