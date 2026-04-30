package core

import "image/color"

// SyncBuilder 辅助声明式重建时同步组件属性，消除 needDraw 样板代码。
type SyncBuilder struct {
	NeedDraw bool
}

// MarkDraw 在 NeedDraw 为 true 时标记元素需要重绘。
func (s *SyncBuilder) MarkDraw(el Element) {
	if s != nil && s.NeedDraw {
		el.Mark(FlagNeedDraw)
	}
}

// SyncField 同步可比较字段，值变化时设置 NeedDraw。
func SyncField[T comparable](s *SyncBuilder, dst *T, src T) {
	if *dst != src {
		*dst = src
		if s != nil {
			s.NeedDraw = true
		}
	}
}

// SyncColor 同步 color.Color，值变化时设置 NeedDraw。
func SyncColor(s *SyncBuilder, dst *color.Color, src color.Color) {
	if !ColorsEqual(*dst, src) {
		*dst = src
		if s != nil {
			s.NeedDraw = true
		}
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

// HasRadius 检查 BorderRadius 是否非零。
func HasRadius(r BorderRadius) bool {
	return r.TopLeft > 0 || r.TopRight > 0 || r.BottomRight > 0 || r.BottomLeft > 0
}
