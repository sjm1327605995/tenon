package antdesign

import (
	"image/color"

	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/components"
	"github.com/sjm1327605995/tenon/yoga"
)

// AntButtonType 定义按钮类型。
type AntButtonType string

const (
	AntButtonPrimary AntButtonType = "primary"
	AntButtonDefault AntButtonType = "default"
	AntButtonDashed  AntButtonType = "dashed"
	AntButtonText    AntButtonType = "text"
	AntButtonLink    AntButtonType = "link"
)

// AntButtonSize 定义按钮尺寸。
type AntButtonSize string

const (
	AntButtonSmall  AntButtonSize = "small"
	AntButtonMiddle AntButtonSize = "middle"
	AntButtonLarge  AntButtonSize = "large"
)

// AntButton 是 Ant Design 风格的按钮组件。
type AntButton struct {
	tenon.BaseWidget
	label    string
	btnType  AntButtonType
	size     AntButtonSize
	danger   bool
	disabled bool
	loading  bool
	onClick  func()

	// 缓存 Host 以在重渲染时复用（保留 Yoga 节点和交互状态）
	btn *components.Button
}

// NewAntButton 创建 Ant Design 按钮。
func NewAntButton(label string) *AntButton {
	b := &AntButton{
		label:   label,
		btnType: AntButtonDefault,
		size:    AntButtonMiddle,
	}
	b.Init(b)
	return b
}

// Render 返回按钮 UI。
func (b *AntButton) Render() tenon.Component {
	theme := NewAntTheme()

	// 复用旧实例：如果存在且类型一致，只同步属性
	if b.btn == nil {
		b.btn = components.NewButton(b.label)
		b.btn.SetOnClick(func() {
			if b.disabled || b.loading {
				return
			}
			if b.onClick != nil {
				b.onClick()
			}
		})
	} else {
		b.btn.SetText(b.label)
	}

	b.applyStyle(theme)
	return b.btn
}

// applyStyle 根据当前配置应用 AntD 样式。
func (b *AntButton) applyStyle(theme *AntTheme) {
	if b.btn == nil {
		return
	}

	// 尺寸
	var h, padding float32
	switch b.size {
	case AntButtonSmall:
		h, padding = 24, 7
	case AntButtonLarge:
		h, padding = 40, 15
	default: // middle
		h, padding = 32, 15
	}
	b.btn.SetHeight(h)
	b.btn.SetPadding(yoga.EdgeHorizontal, padding)
	b.btn.SetPadding(yoga.EdgeVertical, 0)

	// 边框圆角统一 AntD 风格
	b.btn.SetBorderRadius(6)

	if b.disabled {
		b.btn.SetDisabled(true)
		b.btn.SetBackgroundColors(theme.DisabledBgColor, theme.DisabledBgColor, theme.DisabledBgColor)
		b.btn.SetTextColor(theme.DisabledTextColor)
		return
	}
	b.btn.SetDisabled(false)

	if b.loading {
		b.btn.SetText("Loading...") // 简化 loading 显示
		b.btn.SetDisabled(true)
	}

	// 根据类型和 danger 设置颜色
	normal, hover, pressed, textClr := b.resolveColors(theme)
	b.btn.SetBackgroundColors(normal, hover, pressed)
	b.btn.SetTextColor(textClr)

	// 虚线按钮加边框
	if b.btnType == AntButtonDashed {
		b.btn.GetElement().BorderColor = theme.BorderColor
		b.btn.GetElement().Yoga.StyleSetBorder(yoga.EdgeAll, 1)
	} else if b.btnType == AntButtonDefault {
		b.btn.GetElement().BorderColor = theme.BorderColor
		b.btn.GetElement().Yoga.StyleSetBorder(yoga.EdgeAll, 1)
	} else {
		b.btn.GetElement().Yoga.StyleSetBorder(yoga.EdgeAll, 0)
	}
}

func (b *AntButton) resolveColors(theme *AntTheme) (normal, hover, pressed, text color.Color) {
	if b.danger {
		switch b.btnType {
		case AntButtonPrimary:
			return theme.DangerNormalColor, theme.DangerHoverColor, theme.DangerPressedColor, theme.DangerTextColor
		case AntButtonDefault, AntButtonDashed:
			return theme.DangerBgColor, theme.DangerHoverColor, theme.DangerPressedColor, theme.DangerNormalColor
		case AntButtonText, AntButtonLink:
			return color.RGBA{A: 0}, theme.DangerBgColor, theme.DangerBgColor, theme.DangerNormalColor
		}
	}

	switch b.btnType {
	case AntButtonPrimary:
		return theme.PrimaryColor, theme.PrimaryHoverColor, theme.ButtonPressedColor, theme.ButtonTextColor
	case AntButtonDefault, AntButtonDashed:
		return theme.SurfaceColor, theme.BackgroundColor, theme.BackgroundColor, theme.TextColor
	case AntButtonText:
		return color.RGBA{A: 0}, theme.BackgroundColor, theme.BackgroundColor, theme.TextButtonColor
	case AntButtonLink:
		return color.RGBA{A: 0}, color.RGBA{A: 0}, color.RGBA{A: 0}, theme.LinkColor
	}
	return theme.SurfaceColor, theme.BackgroundColor, theme.BackgroundColor, theme.TextColor
}

// ==================== 链式 API ====================

func (b *AntButton) SetType(t AntButtonType) *AntButton {
	b.btnType = t
	return b
}
func (b *AntButton) SetSize(s AntButtonSize) *AntButton {
	b.size = s
	return b
}
func (b *AntButton) SetDanger(v bool) *AntButton {
	b.danger = v
	return b
}
func (b *AntButton) SetDisabled(v bool) *AntButton {
	b.disabled = v
	return b
}
func (b *AntButton) SetLoading(v bool) *AntButton {
	b.loading = v
	return b
}
func (b *AntButton) SetOnClick(fn func()) *AntButton {
	b.onClick = fn
	return b
}
func (b *AntButton) SetWidth(w float32) *AntButton {
	if b.btn != nil {
		b.btn.SetWidth(w)
	}
	return b
}
func (b *AntButton) SetMargin(edge yoga.Edge, v float32) *AntButton {
	if b.btn != nil {
		b.btn.SetMargin(edge, v)
	}
	return b
}
