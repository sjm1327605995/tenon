package ui

import "image/color"

// Color 是 RGBA 颜色，实现 image/color.Color 接口，可直接传给 vector/text。
type Color struct{ R, G, B, A uint8 }

// RGBA 实现 color.Color（返回 alpha 预乘的 16bit 分量）。
func (c Color) RGBA() (r, g, b, a uint32) {
	return color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A}.RGBA()
}

// Alpha 返回按比例调整透明度后的新颜色。
func (c Color) Alpha(a float32) Color {
	return Color{R: c.R, G: c.G, B: c.B, A: uint8(float32(c.A) * a)}
}

func (c Color) transparent() bool { return c.A == 0 }

// Mix 在 a、b 间按 t（0..1）线性插值，供扩展库做 hover/active 等状态配色。
func Mix(a, b Color, t float32) Color {
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}
	li := func(x, y uint8) uint8 { return uint8(float32(x) + (float32(y)-float32(x))*t) }
	return Color{R: li(a.R, b.R), G: li(a.G, b.G), B: li(a.B, b.B), A: li(a.A, b.A)}
}

// Hex 解析 "#RRGGBB" 或 "#RRGGBBAA"，解析失败返回黑色。
func Hex(s string) Color {
	if len(s) > 0 && s[0] == '#' {
		s = s[1:]
	}
	switch len(s) {
	case 6:
		r, g, b := hx(s[0:2]), hx(s[2:4]), hx(s[4:6])
		return Color{r, g, b, 255}
	case 8:
		r, g, b, a := hx(s[0:2]), hx(s[2:4]), hx(s[4:6]), hx(s[6:8])
		return Color{r, g, b, a}
	}
	return Color{0, 0, 0, 255}
}

func hx(s string) uint8 {
	v := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		var d int
		switch {
		case c >= '0' && c <= '9':
			d = int(c - '0')
		case c >= 'a' && c <= 'f':
			d = int(c-'a') + 10
		case c >= 'A' && c <= 'F':
			d = int(c-'A') + 10
		}
		v = v*16 + d
	}
	return uint8(v)
}

// 常用颜色
var (
	Transparent = Color{0, 0, 0, 0}
	Black       = Color{0, 0, 0, 255}
	White       = Color{255, 255, 255, 255}
	Red         = Color{220, 38, 38, 255}
	Green       = Color{34, 197, 94, 255}
	Blue        = Color{59, 130, 246, 255}
	Gray        = Color{156, 163, 175, 255}
	LightGray   = Color{243, 244, 246, 255}
	DarkGray    = Color{75, 85, 99, 255}
)
