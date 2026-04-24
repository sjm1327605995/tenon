package antdesign

import (
	"image/color"
)

// lighten 将颜色变亮，factor 范围 [0, 1]，越大越亮。
func lighten(c color.Color, factor float32) color.Color {
	r, g, b, a := c.RGBA()
	fr := float32(r>>8) + (255-float32(r>>8))*factor
	fg := float32(g>>8) + (255-float32(g>>8))*factor
	fb := float32(b>>8) + (255-float32(b>>8))*factor
	if fr > 255 {
		fr = 255
	}
	if fg > 255 {
		fg = 255
	}
	if fb > 255 {
		fb = 255
	}
	return color.RGBA{R: uint8(fr), G: uint8(fg), B: uint8(fb), A: uint8(a >> 8)}
}

// darken 将颜色变暗，factor 范围 [0, 1]，越大越暗。
func darken(c color.Color, factor float32) color.Color {
	r, g, b, a := c.RGBA()
	fr := float32(r>>8) * (1 - factor)
	fg := float32(g>>8) * (1 - factor)
	fb := float32(b>>8) * (1 - factor)
	if fr < 0 {
		fr = 0
	}
	if fg < 0 {
		fg = 0
	}
	if fb < 0 {
		fb = 0
	}
	return color.RGBA{R: uint8(fr), G: uint8(fg), B: uint8(fb), A: uint8(a >> 8)}
}

// withAlpha 修改颜色的透明度。
func withAlpha(c color.Color, alpha uint8) color.Color {
	r, g, b, _ := c.RGBA()
	return color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: alpha}
}
