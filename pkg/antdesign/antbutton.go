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

// AntButtonShape 定义按钮形状。
type AntButtonShape string

const (
	AntButtonShapeDefault AntButtonShape = "default"
	AntButtonShapeCircle  AntButtonShape = "circle"
	AntButtonShapeRound   AntButtonShape = "round"
)

// AntButtonIconPosition 定义图标位置。
type AntButtonIconPosition string

const (
	AntButtonIconLeft  AntButtonIconPosition = "left"
	AntButtonIconRight AntButtonIconPosition = "right"
)

// AntButton 是 Ant Design 风格的按钮组件。
type AntButton struct {
	tenon.BaseWidget
	label         string
	btnType       AntButtonType
	size          AntButtonSize
	shape         AntButtonShape
	danger        bool
	ghost         bool
	disabled      bool
	loading       bool
	loadingIcon   tenon.Component
	block         bool
	icon          tenon.Component
	iconPosition  AntButtonIconPosition
	onClick       func()

	// 缓存 Host 以在重渲染时复用（保留 Yoga 节点和交互状态）
	btn *components.Button
}

// NewAntButton 创建 Ant Design 按钮。
func NewAntButton(label string) *AntButton {
	b := &AntButton{
		label:        label,
		btnType:      AntButtonDefault,
		size:         AntButtonMiddle,
		shape:        AntButtonShapeDefault,
		iconPosition: AntButtonIconLeft,
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
			if b.disabled {
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
	if b.label == "" {
		// 纯图标按钮：去掉水平 padding，让图标在按钮内完全居中
		b.btn.SetPadding(yoga.EdgeHorizontal, 0)
	} else {
		b.btn.SetPadding(yoga.EdgeHorizontal, padding)
	}
	b.btn.SetPadding(yoga.EdgeVertical, 0)

	// 形状
	switch b.shape {
	case AntButtonShapeCircle:
		b.btn.SetBorderRadius(h / 2)
		b.btn.SetWidth(h)
	case AntButtonShapeRound:
		b.btn.SetBorderRadius(h / 2)
	default:
		b.btn.SetBorderRadius(6)
	}

	// block：宽度 100%
	if b.block {
		b.btn.SetWidthPercent(100)
	}

	// loading 状态
	b.btn.SetLoading(b.loading)

	// 图标与 loading 处理
	b.applyIcons(theme)

	// disabled（非 loading 时）
	if b.disabled && !b.loading {
		b.btn.SetDisabled(true)
		b.btn.SetBackgroundColors(theme.DisabledBgColor, theme.DisabledBgColor, theme.DisabledBgColor)
		b.btn.SetBorderColors(theme.DisabledBorderColor, theme.DisabledBorderColor, theme.DisabledBorderColor)
		b.btn.SetTextColors(theme.DisabledTextColor, theme.DisabledTextColor, theme.DisabledTextColor)
		return
	}
	if !b.loading {
		b.btn.SetDisabled(false)
	}

	// 根据类型、danger、ghost 设置颜色
	bgNormal, bgHover, bgPressed,
		borderNormal, borderHover, borderPressed,
		textNormal, textHover, textPressed := b.resolveColors(theme)

	b.btn.SetBackgroundColors(bgNormal, bgHover, bgPressed)
	b.btn.SetBorderColors(borderNormal, borderHover, borderPressed)
	b.btn.SetTextColors(textNormal, textHover, textPressed)

	// 虚线按钮
	b.btn.SetDashed(b.btnType == AntButtonDashed)

	// 边框宽度
	if b.btnType == AntButtonDefault || b.btnType == AntButtonDashed {
		b.btn.GetElement().Yoga.StyleSetBorder(yoga.EdgeAll, 1)
	} else {
		b.btn.GetElement().Yoga.StyleSetBorder(yoga.EdgeAll, 0)
	}
}

// applyIcons 处理图标和 loading 图标的显示逻辑。
func (b *AntButton) applyIcons(theme *AntTheme) {
	textColor := b.getTextColor(theme)

	if !b.loading {
		// 非 loading 状态，显示普通图标
		if b.icon != nil {
			// 同步图标颜色与按钮文字一致
			if icon, ok := b.icon.(*components.SVGIcon); ok {
				icon.SetColor(textColor)
			}
			if b.iconPosition == AntButtonIconRight {
				b.btn.SetRightIcon(b.icon)
				b.btn.SetLeftIcon(nil)
			} else {
				b.btn.SetLeftIcon(b.icon)
				b.btn.SetRightIcon(nil)
			}
		} else {
			b.btn.SetLeftIcon(nil)
			b.btn.SetRightIcon(nil)
		}
		return
	}

	// loading 状态
	var spinner tenon.Component
	if b.loadingIcon != nil {
		spinner = b.loadingIcon
	} else {
		// 使用默认 loading spinner，颜色与按钮文字一致
		spinner = components.NewLoadingSpinner().SetColor(textColor)
	}

	// loading 时：如果有普通 icon，用 loading 替换它；如果没有，loading 放在左边
	if b.icon != nil && b.iconPosition == AntButtonIconRight {
		b.btn.SetLeftIcon(nil)
		b.btn.SetRightIcon(spinner)
	} else {
		b.btn.SetLeftIcon(spinner)
		b.btn.SetRightIcon(nil)
	}
}

// getTextColor 返回当前按钮正常状态下的文字颜色，用于 loading spinner 配色。
func (b *AntButton) getTextColor(theme *AntTheme) color.Color {
	if b.disabled {
		return theme.DisabledTextColor
	}
	if b.ghost {
		switch b.btnType {
		case AntButtonPrimary:
			if b.danger {
				return theme.DangerNormalColor
			}
			return theme.PrimaryColor
		case AntButtonDefault, AntButtonDashed:
			if b.danger {
				return theme.DangerNormalColor
			}
			return theme.PrimaryColor
		case AntButtonText, AntButtonLink:
			if b.danger {
				return theme.DangerNormalColor
			}
			return theme.PrimaryColor
		}
	}
	if b.danger {
		switch b.btnType {
		case AntButtonPrimary:
			return theme.DangerTextColor
		case AntButtonDefault, AntButtonDashed, AntButtonText, AntButtonLink:
			return theme.DangerNormalColor
		}
	}
	switch b.btnType {
	case AntButtonPrimary:
		return theme.ButtonTextColor
	case AntButtonDefault, AntButtonDashed:
		return theme.TextColor
	case AntButtonText:
		return theme.TextColor
	case AntButtonLink:
		return theme.LinkColor
	}
	return theme.TextColor
}

func (b *AntButton) resolveColors(theme *AntTheme) (
	bgNormal, bgHover, bgPressed,
	borderNormal, borderHover, borderPressed,
	textNormal, textHover, textPressed color.Color,
) {
	if b.ghost {
		return b.resolveGhostColors(theme)
	}
	if b.danger {
		return b.resolveDangerColors(theme)
	}

	switch b.btnType {
	case AntButtonPrimary:
		return theme.PrimaryColor, theme.PrimaryHoverColor, theme.ButtonPressedColor,
			theme.PrimaryColor, theme.PrimaryHoverColor, theme.ButtonPressedColor,
			theme.ButtonTextColor, theme.ButtonTextColor, theme.ButtonTextColor
	case AntButtonDefault, AntButtonDashed:
		return theme.SurfaceColor, theme.SurfaceColor, theme.SurfaceColor,
			theme.BorderColor, theme.PrimaryColor, theme.ButtonPressedColor,
			theme.TextColor, theme.PrimaryColor, theme.ButtonPressedColor
	case AntButtonText:
		return color.RGBA{A: 0}, theme.TextButtonHoverBg, theme.TextButtonActiveBg,
			color.RGBA{A: 0}, color.RGBA{A: 0}, color.RGBA{A: 0},
			theme.TextColor, theme.TextColor, theme.TextColor
	case AntButtonLink:
		return color.RGBA{A: 0}, color.RGBA{A: 0}, color.RGBA{A: 0},
			color.RGBA{A: 0}, color.RGBA{A: 0}, color.RGBA{A: 0},
			theme.LinkColor, theme.LinkHoverColor, theme.LinkPressedColor
	}
	return theme.SurfaceColor, theme.SurfaceColor, theme.SurfaceColor,
		theme.BorderColor, theme.PrimaryColor, theme.ButtonPressedColor,
		theme.TextColor, theme.PrimaryColor, theme.ButtonPressedColor
}

func (b *AntButton) resolveDangerColors(theme *AntTheme) (
	bgNormal, bgHover, bgPressed,
	borderNormal, borderHover, borderPressed,
	textNormal, textHover, textPressed color.Color,
) {
	switch b.btnType {
	case AntButtonPrimary:
		return theme.DangerNormalColor, theme.DangerHoverColor, theme.DangerPressedColor,
			theme.DangerNormalColor, theme.DangerHoverColor, theme.DangerPressedColor,
			theme.DangerTextColor, theme.DangerTextColor, theme.DangerTextColor
	case AntButtonDefault, AntButtonDashed:
		return color.RGBA{R: 255, G: 255, B: 255, A: 255}, theme.DangerBgColor, theme.DangerBgColor,
			theme.DangerNormalColor, theme.DangerHoverColor, theme.DangerPressedColor,
			theme.DangerNormalColor, theme.DangerHoverColor, theme.DangerPressedColor
	case AntButtonText, AntButtonLink:
		return color.RGBA{A: 0}, theme.DangerBgColor, theme.DangerBgColor,
			color.RGBA{A: 0}, color.RGBA{A: 0}, color.RGBA{A: 0},
			theme.DangerNormalColor, theme.DangerHoverColor, theme.DangerPressedColor
	}
	return theme.SurfaceColor, theme.SurfaceColor, theme.SurfaceColor,
		theme.BorderColor, theme.PrimaryColor, theme.ButtonPressedColor,
		theme.TextColor, theme.PrimaryColor, theme.ButtonPressedColor
}

func (b *AntButton) resolveGhostColors(theme *AntTheme) (
	bgNormal, bgHover, bgPressed,
	borderNormal, borderHover, borderPressed,
	textNormal, textHover, textPressed color.Color,
) {
	transparent := color.RGBA{A: 0}
	switch b.btnType {
	case AntButtonPrimary:
		if b.danger {
			return transparent, theme.DangerNormalColor, theme.DangerHoverColor,
				theme.DangerNormalColor, theme.DangerNormalColor, theme.DangerHoverColor,
				theme.DangerNormalColor, theme.ButtonTextColor, theme.ButtonTextColor
		}
		return transparent, theme.PrimaryColor, theme.PrimaryHoverColor,
			theme.PrimaryColor, theme.PrimaryColor, theme.PrimaryHoverColor,
			theme.PrimaryColor, theme.ButtonTextColor, theme.ButtonTextColor
	case AntButtonDefault, AntButtonDashed:
		if b.danger {
			return transparent, theme.DangerNormalColor, theme.DangerHoverColor,
				theme.DangerNormalColor, theme.DangerNormalColor, theme.DangerHoverColor,
				theme.DangerNormalColor, theme.DangerNormalColor, theme.DangerHoverColor
		}
		return transparent, theme.PrimaryColor, theme.PrimaryHoverColor,
			theme.PrimaryColor, theme.PrimaryColor, theme.PrimaryHoverColor,
			theme.PrimaryColor, theme.PrimaryColor, theme.PrimaryHoverColor
	case AntButtonText, AntButtonLink:
		if b.danger {
			return transparent, theme.DangerBgColor, theme.DangerBgColor,
				transparent, transparent, transparent,
				theme.DangerNormalColor, theme.DangerHoverColor, theme.DangerPressedColor
		}
		return transparent, theme.TextButtonHoverBg, theme.TextButtonActiveBg,
			transparent, transparent, transparent,
			theme.PrimaryColor, theme.PrimaryHoverColor, theme.ButtonPressedColor
	}
	return transparent, theme.PrimaryColor, theme.PrimaryHoverColor,
		theme.PrimaryColor, theme.PrimaryColor, theme.PrimaryHoverColor,
		theme.PrimaryColor, theme.PrimaryColor, theme.PrimaryHoverColor
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
func (b *AntButton) SetShape(s AntButtonShape) *AntButton {
	b.shape = s
	return b
}
func (b *AntButton) SetDanger(v bool) *AntButton {
	b.danger = v
	return b
}
func (b *AntButton) SetGhost(v bool) *AntButton {
	b.ghost = v
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
func (b *AntButton) SetLoadingIcon(icon tenon.Component) *AntButton {
	b.loadingIcon = icon
	return b
}
func (b *AntButton) SetBlock(v bool) *AntButton {
	b.block = v
	return b
}
func (b *AntButton) SetIcon(icon tenon.Component) *AntButton {
	b.icon = icon
	return b
}
func (b *AntButton) SetAntIcon(name AntIconName) *AntButton {
	b.icon = NewAntIcon(name)
	return b
}
func (b *AntButton) SetIconPosition(pos AntButtonIconPosition) *AntButton {
	b.iconPosition = pos
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
