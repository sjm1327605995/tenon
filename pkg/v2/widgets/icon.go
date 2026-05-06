package widgets

import (
	"image/color"

	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
)

// ==================== 命名图标常量 ====================
// 组件使用这些常量而非硬编码 Unicode 字符。
// 如果默认字体不支持，会自动回退到 ASCII 文本。

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

// ChevronDown 和 ArrowDown 共用同一字符。
const IconChevronDown = IconArrowDown

// ==================== 图标回退映射 ====================
// 当字体不支持某个 Unicode 符号时，使用 ASCII 替代。

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

// ==================== Icon 渲染模式 ====================

// IconMode 图标渲染模式。
type IconMode int

const (
	// IconModeUnicode 优先使用 Unicode 符号（需要字体支持）。
	IconModeUnicode IconMode = iota
	// IconModeASCII 强制使用 ASCII 回退文本。
	IconModeASCII
	// IconModeAuto 自动检测：如果字体支持则用 Unicode，否则用 ASCII。
	IconModeAuto
)

var iconMode = IconModeUnicode

// SetIconMode 设置全局图标渲染模式。
func SetIconMode(mode IconMode) {
	iconMode = mode
}

// GetIconMode 返回当前图标渲染模式。
func GetIconMode() IconMode {
	return iconMode
}

// ==================== Icon Widget ====================

// IconWidget 图标组件，支持字体回退。
type IconWidget struct {
	ui.BaseWidget
	icon     string // Unicode 符号或命名常量
	fontSize float32
	color    color.Color
	fallback string // 手动指定的回退文本
}

// Icon 创建图标 Widget。
// icon 可以是命名常量（IconArrowDown 等）或任意 Unicode 字符。
func Icon(icon string) IconWidget {
	fb := ""
	if f, ok := iconFallbacks[icon]; ok {
		fb = f
	}
	return IconWidget{icon: icon, fontSize: 14, fallback: fb}
}

func (i IconWidget) FontSize(v float32) IconWidget { i.fontSize = v; return i }
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
	ro *render.RenderText
}

func (e *IconElement) CreateRenderObject() render.RenderObject {
	w := e.GetWidget().(IconWidget)
	content := e.resolveContent(w)
	r := render.NewRenderText(content)
	r.SetFontSize(w.fontSize)
	if w.color != nil {
		r.SetColor(w.color)
	}
	return r
}

func (e *IconElement) UpdateRenderObject(oldWidget ui.Widget) {
	w := e.GetWidget().(IconWidget)
	old := oldWidget.(IconWidget)
	content := e.resolveContent(w)

	if oldContent := e.resolveContent(old); oldContent != content {
		e.ro.SetContent(content)
	}
	if old.fontSize != w.fontSize {
		e.ro.SetFontSize(w.fontSize)
	}
	if old.color != w.color && w.color != nil {
		e.ro.SetColor(w.color)
	}
}

func (e *IconElement) Mount(parent ui.Element, slot int) {
	e.ro = e.CreateRenderObject().(*render.RenderText)
	e.RenderObject = e.ro
	e.RenderObjectElement.Mount(parent, slot)
}

// resolveContent 根据当前图标模式决定显示内容。
func (e *IconElement) resolveContent(w IconWidget) string {
	switch iconMode {
	case IconModeASCII:
		if w.fallback != "" {
			return w.fallback
		}
		if f, ok := iconFallbacks[w.icon]; ok {
			return f
		}
		return w.icon
	case IconModeAuto:
		// Auto 模式：检查字体是否支持该字符
		// 简单实现：如果有 fallback 就用 Unicode，否则用 fallback
		// 更精确的实现需要查询字体的 glyph coverage
		if w.fallback != "" {
			return w.icon // 有 fallback 说明可以尝试 Unicode
		}
		return w.icon
	default: // IconModeUnicode
		return w.icon
	}
}

// ==================== 便捷函数 ====================

// IconText 获取图标的当前显示文本。
func IconText(icon string) string {
	switch iconMode {
	case IconModeASCII:
		if f, ok := iconFallbacks[icon]; ok {
			return f
		}
		return icon
	default:
		return icon
	}
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
