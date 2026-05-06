package render

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// BorderRadius 描述四个角的圆角半径。
type BorderRadius struct {
	TopLeft, TopRight, BottomRight, BottomLeft float32
}

// IsZero 判断是否所有角半径都为 0。
func (br BorderRadius) IsZero() bool {
	return br.TopLeft == 0 && br.TopRight == 0 && br.BottomRight == 0 && br.BottomLeft == 0
}

// Uniform 返回四角相同的 BorderRadius。
func UniformBorderRadius(r float32) BorderRadius {
	return BorderRadius{TopLeft: r, TopRight: r, BottomRight: r, BottomLeft: r}
}

// clampRadius 将每个角的半径限制为不超过半宽和半高。
func clampRadius(br BorderRadius, w, h float32) BorderRadius {
	maxR := w / 2
	if h/2 < maxR {
		maxR = h / 2
	}
	if br.TopLeft > maxR {
		br.TopLeft = maxR
	}
	if br.TopRight > maxR {
		br.TopRight = maxR
	}
	if br.BottomRight > maxR {
		br.BottomRight = maxR
	}
	if br.BottomLeft > maxR {
		br.BottomLeft = maxR
	}
	return br
}

func toColorScale(c color.Color) ebiten.ColorScale {
	var cs ebiten.ColorScale
	cs.ScaleWithColor(c)
	return cs
}

// BuildRoundedRectPath 构建支持四角独立半径的圆角矩形路径。
func BuildRoundedRectPath(path *vector.Path, x, y, w, h float32, br BorderRadius) {
	br = clampRadius(br, w, h)
	tl, tr, brR, bl := br.TopLeft, br.TopRight, br.BottomRight, br.BottomLeft

	path.MoveTo(x+tl, y)
	path.LineTo(x+w-tr, y)
	if tr > 0 {
		path.Arc(x+w-tr, y+tr, tr, -float32(math.Pi/2), 0, vector.Clockwise)
	}
	path.LineTo(x+w, y+h-brR)
	if brR > 0 {
		path.Arc(x+w-brR, y+h-brR, brR, 0, float32(math.Pi/2), vector.Clockwise)
	}
	path.LineTo(x+bl, y+h)
	if bl > 0 {
		path.Arc(x+bl, y+h-bl, bl, float32(math.Pi/2), float32(math.Pi), vector.Clockwise)
	}
	path.LineTo(x, y+tl)
	if tl > 0 {
		path.Arc(x+tl, y+tl, tl, float32(math.Pi), float32(math.Pi*1.5), vector.Clockwise)
	}
	path.Close()
}

// DrawRoundedRectFill 绘制填充圆角矩形（支持四角独立半径）。
func DrawRoundedRectFill(screen *ebiten.Image, x, y, w, h float32, br BorderRadius, c color.Color) {
	if c == nil || w <= 0 || h <= 0 {
		return
	}
	path := &vector.Path{}
	BuildRoundedRectPath(path, x, y, w, h, br)
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(c)
	op.AntiAlias = true
	vector.FillPath(screen, path, &vector.FillOptions{}, op)
}

// DrawRoundedRectStroke 绘制描边圆角矩形（支持四角独立半径）。
func DrawRoundedRectStroke(screen *ebiten.Image, x, y, w, h float32, br BorderRadius, strokeWidth float32, c color.Color) {
	if c == nil || strokeWidth <= 0 || w <= 0 || h <= 0 {
		return
	}
	path := &vector.Path{}
	BuildRoundedRectPath(path, x, y, w, h, br)
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(c)
	op.AntiAlias = true
	vector.StrokePath(screen, path, &vector.StrokeOptions{
		Width: float32(strokeWidth),
	}, op)
}

// DrawRect 绘制填充矩形。
func DrawRect(screen *ebiten.Image, x, y, w, h float32, c color.Color) {
	if c == nil {
		return
	}
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), float32(h), c, false)
}

// DrawBorder 绘制直角边框（四边分别绘制）。
func DrawBorder(screen *ebiten.Image, x, y, w, h, width float32, c color.Color) {
	if c == nil || width <= 0 {
		return
	}
	// top
	vector.DrawFilledRect(screen, x, y, w, width, c, false)
	// right
	vector.DrawFilledRect(screen, x+w-width, y, width, h, c, false)
	// bottom
	vector.DrawFilledRect(screen, x, y+h-width, w, width, c, false)
	// left
	vector.DrawFilledRect(screen, x, y, width, h, c, false)
}

// DrawShadow 绘制矩形阴影（纯色扩展模拟）。
func DrawShadow(screen *ebiten.Image, x, y, w, h float32, br BorderRadius, c color.Color, blur, offsetX, offsetY float32) {
	if c == nil {
		return
	}
	sx := x + offsetX
	sy := y + offsetY
	if !br.IsZero() {
		DrawRoundedRectFill(screen, sx, sy, w, h, br, c)
	} else {
		DrawRect(screen, sx, sy, w, h, c)
	}
}

// DrawFilledCircle 绘制填充圆形。
func DrawFilledCircle(screen *ebiten.Image, cx, cy, radius float32, c color.Color) {
	if c == nil || radius <= 0 {
		return
	}
	vector.DrawFilledCircle(screen, cx, cy, radius, c, true)
}

// StrokeCircle 绘制描边圆形。
func StrokeCircle(screen *ebiten.Image, cx, cy, radius, strokeWidth float32, c color.Color) {
	if c == nil || strokeWidth <= 0 || radius <= 0 {
		return
	}
	path := &vector.Path{}
	path.MoveTo(cx+radius, cy)
	path.Arc(cx, cy, radius, 0, float32(math.Pi*2), vector.Clockwise)
	vector.StrokePath(screen, path, &vector.StrokeOptions{
		Width: float32(strokeWidth),
	}, &vector.DrawPathOptions{
		ColorScale: toColorScale(c),
		AntiAlias:  true,
	})
}

// DrawFilledCirclePath 用 vector.Path 绘制填充圆形（支持透明色）。
func DrawFilledCirclePath(screen *ebiten.Image, cx, cy, radius float32, c color.Color) {
	if c == nil || radius <= 0 {
		return
	}
	path := &vector.Path{}
	path.MoveTo(cx+radius, cy)
	path.Arc(cx, cy, radius, 0, float32(math.Pi*2), vector.Clockwise)
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(c)
	op.AntiAlias = true
	vector.FillPath(screen, path, &vector.FillOptions{}, op)
}

// DrawStrokedCirclePath 用 vector.Path 绘制描边圆形（支持透明色）。
func DrawStrokedCirclePath(screen *ebiten.Image, cx, cy, radius, strokeWidth float32, c color.Color) {
	if c == nil || radius <= 0 {
		return
	}
	path := &vector.Path{}
	path.MoveTo(cx+radius, cy)
	path.Arc(cx, cy, radius, 0, float32(math.Pi*2), vector.Clockwise)
	strokeOp := &vector.StrokeOptions{Width: strokeWidth, MiterLimit: 10}
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(c)
	op.AntiAlias = true
	vector.StrokePath(screen, path, strokeOp, op)
}

// Transform 定义 2D 仿射变换参数。
type Transform struct {
	Rotation  float32
	ScaleX    float32
	ScaleY    float32
	SkewX     float32
	SkewY     float32
	TranslateX float32 // 平移 X（像素）
	TranslateY float32 // 平移 Y（像素）
	OriginX   float32 // 0~1，相对于元素宽高的比例
	OriginY   float32 // 0~1，相对于元素宽高的比例
	Alpha     float32
}

// DefaultTransform 返回无变换的默认值。
func DefaultTransform() Transform {
	return Transform{ScaleX: 1, ScaleY: 1, OriginX: 0.5, OriginY: 0.5, Alpha: 1}
}

// IsIdentity 检查是否接近无变换状态。
func (t Transform) IsIdentity() bool {
	return t.Rotation == 0 && t.ScaleX == 1 && t.ScaleY == 1 &&
		t.SkewX == 0 && t.SkewY == 0 && t.TranslateX == 0 && t.TranslateY == 0 && t.Alpha == 1
}

// BuildTransformGeoM 以元素中心（由 OriginX/OriginY 比例决定）为锚点，构建变换矩阵。
// 变换顺序：Translate(原点) → Skew → Rotate → Scale → Translate(-原点) → Translate(offset)。
func BuildTransformGeoM(bounds Bounds, t Transform) ebiten.GeoM {
	var g ebiten.GeoM
	ox := float64(bounds.X + bounds.Width*t.OriginX)
	oy := float64(bounds.Y + bounds.Height*t.OriginY)

	g.Translate(-ox, -oy)
	g.Skew(float64(t.SkewX), float64(t.SkewY))
	g.Rotate(float64(t.Rotation) * math.Pi / 180)
	g.Scale(float64(t.ScaleX), float64(t.ScaleY))
	g.Translate(ox, oy)
	g.Translate(float64(t.TranslateX), float64(t.TranslateY))
	return g
}

// ApplyColorScaleAlpha 在已有 ColorScale 基础上应用透明度。
func ApplyColorScaleAlpha(cs *ebiten.ColorScale, alpha float32) {
	if alpha < 0 {
		alpha = 0
	}
	if alpha > 1 {
		alpha = 1
	}
	cs.Scale(float32(alpha), float32(alpha), float32(alpha), float32(alpha))
}

// SubImage 创建用于 ClipChildren 的子图。
func SubImage(screen *ebiten.Image, x, y, w, h int) *ebiten.Image {
	if w <= 0 || h <= 0 {
		return nil
	}
	sub := screen.SubImage(image.Rect(x, y, x+w, y+h))
	if img, ok := sub.(*ebiten.Image); ok {
		return img
	}
	return nil
}

// PaintBackground 绘制背景（自动判断圆角/矩形）。
func PaintBackground(screen *ebiten.Image, x, y, w, h float32, br BorderRadius, bg color.Color) {
	if bg == nil {
		return
	}
	if !br.IsZero() {
		DrawRoundedRectFill(screen, x, y, w, h, br, bg)
	} else {
		DrawRect(screen, x, y, w, h, bg)
	}
}

// PaintBorder 绘制边框（自动判断圆角/矩形）。
func PaintBorder(screen *ebiten.Image, x, y, w, h float32, br BorderRadius, borderWidth float32, borderColor color.Color) {
	if borderColor == nil || borderWidth <= 0 {
		return
	}
	if !br.IsZero() {
		DrawRoundedRectStroke(screen, x, y, w, h, br, borderWidth, borderColor)
	} else {
		DrawBorder(screen, x, y, w, h, borderWidth, borderColor)
	}
}

// PaintBoxShadow 绘制矩形阴影。
func PaintBoxShadow(screen *ebiten.Image, x, y, w, h float32, br BorderRadius, shadow color.Color, blur, offsetX, offsetY float32) {
	if shadow == nil {
		return
	}
	DrawShadow(screen, x, y, w, h, br, shadow, blur, offsetX, offsetY)
}

// Darken 对颜色做暗化处理（各分量减 delta，最低到 0）。
func Darken(c color.Color, delta uint8) color.Color {
	r, g, b, a := c.RGBA()
	d := func(v uint8) uint8 {
		if v >= delta {
			return v - delta
		}
		return 0
	}
	return color.RGBA{R: d(uint8(r >> 8)), G: d(uint8(g >> 8)), B: d(uint8(b >> 8)), A: uint8(a >> 8)}
}

// Lighten 对颜色做亮化处理（各分量加 delta，最高到 255）。
func Lighten(c color.Color, delta uint8) color.Color {
	r, g, b, a := c.RGBA()
	l := func(v uint8) uint8 {
		if v >= 255-delta {
			return 255
		}
		return v + delta
	}
	return color.RGBA{R: l(uint8(r >> 8)), G: l(uint8(g >> 8)), B: l(uint8(b >> 8)), A: uint8(a >> 8)}
}
