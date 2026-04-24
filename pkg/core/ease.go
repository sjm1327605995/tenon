package core

import (
	"image/color"
	"math"
)

// EasingFunction 是缓动函数类型，输入 t 的范围是 [0, 1]，输出也应在 [0, 1]。
type EasingFunction func(t float32) float32

// LerpFloat32 在 a 和 b 之间进行线性插值。
func LerpFloat32(a, b, t float32) float32 {
	return a + (b-a)*t
}

// LerpColor 在两个颜色之间进行线性插值。
func LerpColor(a, b color.Color, t float32) color.Color {
	ar, ag, ab, aa := a.RGBA()
	br, bg, bb, ba := b.RGBA()
	return color.RGBA{
		R: uint8(LerpFloat32(float32(ar>>8), float32(br>>8), t)),
		G: uint8(LerpFloat32(float32(ag>>8), float32(bg>>8), t)),
		B: uint8(LerpFloat32(float32(ab>>8), float32(bb>>8), t)),
		A: uint8(LerpFloat32(float32(aa>>8), float32(ba>>8), t)),
	}
}

// EaseLinear 线性缓动。
func EaseLinear(t float32) float32 {
	return t
}

// EaseInOutQuad 二次方缓入缓出。
func EaseInOutQuad(t float32) float32 {
	if t < 0.5 {
		return 2 * t * t
	}
	return -1 + (4-2*t)*t
}

// EaseOutQuad 二次方缓出。
func EaseOutQuad(t float32) float32 {
	return 1 - (1-t)*(1-t)
}

// EaseOutBounce 弹跳缓出。
func EaseOutBounce(t float32) float32 {
	tf := float64(t)
	const n1 = 7.5625
	const d1 = 2.75
	if tf < 1/d1 {
		return float32(n1 * tf * tf)
	} else if tf < 2/d1 {
		tf -= 1.5 / d1
		return float32(n1*tf*tf + 0.75)
	} else if tf < 2.5/d1 {
		tf -= 2.25 / d1
		return float32(n1*tf*tf + 0.9375)
	} else {
		tf -= 2.625 / d1
		return float32(n1*tf*tf + 0.984375)
	}
}

// EaseOutElastic 弹性缓出（Spring 效果）。
func EaseOutElastic(t float32) float32 {
	if t == 0 {
		return 0
	}
	if t == 1 {
		return 1
	}
	const c4 = (2 * math.Pi) / 3
	tf := float64(t)
	return float32(math.Pow(2, -10*tf)*math.Sin((tf*10-0.75)*c4) + 1)
}
