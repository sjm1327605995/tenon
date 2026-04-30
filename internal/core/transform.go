package core

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// BuildTransformGeoM 构建从局部坐标系到屏幕坐标系的完整变换矩阵。
// 变换以元素 bounds 内的 (OriginX, OriginY) 比例为变换中心。
func BuildTransformGeoM(bounds LayoutBounds, t Transform) ebiten.GeoM {
	var g ebiten.GeoM

	ox := float64(bounds.Width) * float64(t.OriginX)
	oy := float64(bounds.Height) * float64(t.OriginY)

	// 1. 将变换中心移到局部原点
	g.Translate(-ox, -oy)

	// 2. 应用 Skew → Rotate → Scale
	if t.SkewX != 0 || t.SkewY != 0 {
		g.Skew(float64(t.SkewX)*math.Pi/180, float64(t.SkewY)*math.Pi/180)
	}
	if t.Rotation != 0 {
		g.Rotate(float64(t.Rotation) * math.Pi / 180)
	}
	if t.ScaleX != 1 || t.ScaleY != 1 {
		g.Scale(float64(t.ScaleX), float64(t.ScaleY))
	}

	// 3. 移回并叠加屏幕位置
	g.Translate(float64(bounds.X)+ox, float64(bounds.Y)+oy)

	return g
}

// ApplyColorScaleAlpha 将 Alpha 应用到 ebiten.ColorScale。
func ApplyColorScaleAlpha(scale *ebiten.ColorScale, alpha float32) {
	if alpha >= 0 && alpha < 1 {
		scale.ScaleAlpha(alpha)
	}
}
