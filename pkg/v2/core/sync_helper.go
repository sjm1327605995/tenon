package core

import "image/color"

// SyncField 同步一个可比较字段，若变化则设置 changed 标记。
func SyncField[T comparable](dst *T, src T, changed *bool) {
	if *dst != src {
		*dst = src
		*changed = true
	}
}

// SyncColor 同步一个颜色字段，若变化则设置 changed 标记。
func SyncColor(dst *color.Color, src color.Color, changed *bool) {
	if !ColorsEqual(*dst, src) {
		*dst = src
		*changed = true
	}
}

// ColorsEqual 比较两个 color.Color 是否相等。
func ColorsEqual(a, b color.Color) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	ra, ga, ba, aa := a.RGBA()
	rb, gb, bb, ab := b.RGBA()
	return ra == rb && ga == gb && ba == bb && aa == ab
}
