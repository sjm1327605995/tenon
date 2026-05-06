package widgets

import (
	"image/color"

	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
)

// ==================== 命名图标常量 ====================
// 使用 Unicode 字符作为内部标识符（兼容旧代码），
// 实际渲染通过 iconPaths 映射到 Lucide-style SVG path 数据。

const (
	IconArrowUp         = "▲"
	IconArrowDown       = "▼"
	IconArrowLeft       = "←"
	IconArrowRight      = "→"
	IconChevronRight    = "▶"
	IconCheckboxEmpty   = "☐"
	IconCheckboxChecked = "☑"
	IconInfo            = "ⓘ"
	IconClose           = "✕"
	IconCheck           = "✓"
	IconMinus           = "−"
	IconPlus            = "+"
	IconSearch          = "🔍"
	IconStar            = "★"
	IconStarEmpty       = "☆"
)

// Chevron 系列与 Arrow 系列共用同一字符作为标识符。
const (
	IconChevronDown = IconArrowDown
	IconChevronUp   = IconArrowUp
	IconChevronLeft = IconArrowLeft
)

// iconPaths 将命名常量映射到 SVG path 数据（viewBox="0 0 24 24"）。
var iconPaths = map[string]string{
	IconArrowUp:         "m6 15 6-6 6 6",
	IconArrowDown:       "m6 9 6 6 6-6",
	IconArrowLeft:       "m15 18-6-6 6-6",
	IconArrowRight:      "m9 18 6-6-6-6",
	IconChevronRight:    "m9 18 6-6-6-6",
	IconCheckboxEmpty:   "M3 3h18v18H3V3z",
	IconCheckboxChecked: "M3 3h18v18H3zM9 12l2 2 5-5",
	IconInfo:            "M12 2C16.97 2 21 6.03 21 11C21 15.97 16.97 20 12 20C7.03 20 3 15.97 3 11C3 6.03 7.03 2 12 2zM12 16v-4M11 8h2",
	IconClose:           "M18 6 6 18M6 6l12 12",
	IconCheck:           "M20 6 9 17l-5-5",
	IconMinus:           "M5 12h14",
	IconPlus:            "M5 12h14M12 5v14",
	IconSearch:          "m21 21-4.3-4.3M11 3C15.42 3 19 6.58 19 11C19 15.42 15.42 19 11 19C6.58 19 3 15.42 3 11C3 6.58 6.58 3 11 3z",
	IconStar:            "M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z",
	IconStarEmpty:       "M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z",
}

// iconFallbacks 提供 ASCII 文本回退（用于 Button 等需要文本的场景）。
var iconFallbacks = map[string]string{
	IconArrowUp:         "^",
	IconArrowDown:       "v",
	IconArrowLeft:       "<",
	IconArrowRight:      ">",
	IconChevronRight:    ">",
	IconCheckboxEmpty:   "[ ]",
	IconCheckboxChecked: "[x]",
	IconInfo:            "(i)",
	IconClose:           "x",
	IconCheck:           "*",
	IconMinus:           "-",
	IconSearch:          "[Q]",
	IconStar:            "*",
	IconStarEmpty:       ".",
}

// ==================== Icon Widget ====================

// IconWidget 图标组件，基于 SVG path 矢量渲染。
type IconWidget struct {
	ui.BaseWidget
	path     string // SVG path 数据
	size     float32
	color    color.Color
	fallback string // 手动指定的回退文本
}

// Icon 创建图标 Widget。
// icon 可以是命名常量（IconArrowDown 等）或任意 SVG path 数据。
func Icon(icon string) IconWidget {
	path := icon
	if p, ok := iconPaths[icon]; ok {
		path = p
	}
	fb := ""
	if f, ok := iconFallbacks[icon]; ok {
		fb = f
	}
	return IconWidget{path: path, size: 16, fallback: fb}
}

// Size 设置图标尺寸（默认 16）。
func (i IconWidget) Size(v float32) IconWidget { i.size = v; return i }

// FontSize 兼容旧 API，等价于 Size。
func (i IconWidget) FontSize(v float32) IconWidget { i.size = v; return i }

func (i IconWidget) Color(c color.Color) IconWidget { i.color = c; return i }

// Fallback 手动指定回退文本（覆盖默认回退）。
func (i IconWidget) Fallback(text string) IconWidget { i.fallback = text; return i }

func (i IconWidget) CreateElement() ui.Element {
	e := &IconElement{}
	e.RenderObjectElement.BaseElement.Init(e, i)
	return e
}

func (i IconWidget) GetKey() ui.Key { return ui.NilKey{} }

// IconElement 是 IconWidget 对应的 Element。
type IconElement struct {
	ui.RenderObjectElement
	ro *render.RenderSvgIcon
}

func (e *IconElement) CreateRenderObject() render.RenderObject {
	w := e.GetWidget().(IconWidget)
	r := render.NewRenderSvgIcon()
	r.SetPathData(w.path)
	r.SetIconSize(w.size)
	if w.color != nil {
		r.SetIconColor(w.color)
	}
	return r
}

func (e *IconElement) UpdateRenderObject(oldWidget ui.Widget) {
	w := e.GetWidget().(IconWidget)
	old := oldWidget.(IconWidget)

	if old.path != w.path {
		e.ro.SetPathData(w.path)
	}
	if old.size != w.size {
		e.ro.SetIconSize(w.size)
	}
	if old.color != w.color && w.color != nil {
		e.ro.SetIconColor(w.color)
	}
}

func (e *IconElement) Mount(parent ui.Element, slot int) {
	e.ro = e.CreateRenderObject().(*render.RenderSvgIcon)
	e.RenderObject = e.ro
	e.RenderObjectElement.Mount(parent, slot)
}

// ==================== 便捷函数 ====================

// IconText 获取图标的当前显示文本（用于 Button 等需要文本的场景）。
func IconText(icon string) string {
	if f, ok := iconFallbacks[icon]; ok {
		return f
	}
	return icon
}

// HasIconFallback 检查图标是否有 ASCII 回退。
func HasIconFallback(icon string) bool {
	_, ok := iconFallbacks[icon]
	return ok
}

// RegisterIconFallback 注册自定义图标的 ASCII 回退。
func RegisterIconFallback(icon, fallback string) {
	iconFallbacks[icon] = fallback
}
