package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// ButtonState 表示按钮的交互状态。
type ButtonState int

const (
	ButtonStateNormal ButtonState = iota
	ButtonStateHover
	ButtonStatePressed
)

// Button 是按钮宿主组件。
type Button struct {
	core.BaseHost
	text         *Text
	state        ButtonState
	onClick      func()
	hoverColor   color.Color
	pressedColor color.Color
	normalColor  color.Color
	disabled     bool
}

// NewButton 创建一个按钮。
func NewButton(label string) *Button {
	theme := core.GetTheme()
	b := &Button{
		state:        ButtonStateNormal,
		normalColor:  theme.ButtonNormalColor,
		hoverColor:   theme.ButtonHoverColor,
		pressedColor: theme.ButtonPressedColor,
	}
	b.Init(b)
	b.SetFocusable(true)
	b.SetPadding(yoga.EdgeAll, 12)
	b.SetBorderRadius(theme.ButtonBorderRadius)
	b.SetBackgroundColor(b.normalColor)
	b.SetJustifyContent(yoga.JustifyCenter)
	b.SetAlignItems(yoga.AlignCenter)

	b.text = NewText(label)
	b.text.SetColor(theme.ButtonTextColor)
	b.text.SetFontSize(theme.FontSizeBase + 2)
	b.AddChild(b.text)

	return b
}

// Draw 绘制按钮背景。
func (b *Button) Draw(screen *ebiten.Image) {
	el := b.GetElement()
	if el == nil || !el.Visible {
		return
	}
	bounds := b.GetLayoutBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}
	if el.BackgroundColor != nil {
		vector.FillRect(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, el.BackgroundColor, false)
	}
}

// HandleEvent 处理按钮交互事件。
func (b *Button) HandleEvent(e *core.Event) bool {
	if b.disabled {
		return false
	}

	switch e.Type {
	case core.EventMouseDown:
		b.state = ButtonStatePressed
		b.refreshColor()
	case core.EventMouseUp:
		b.state = ButtonStateHover
		b.refreshColor()
	case core.EventClick:
		if b.onClick != nil {
			b.onClick()
		}
		return true
	}
	return false
}

func (b *Button) refreshColor() {
	switch b.state {
	case ButtonStateHover:
		b.SetBackgroundColor(b.hoverColor)
	case ButtonStatePressed:
		b.SetBackgroundColor(b.pressedColor)
	default:
		b.SetBackgroundColor(b.normalColor)
	}
}

// ==================== 链式 API ====================

func (b *Button) SetOnClick(fn func()) *Button {
	b.onClick = fn
	return b
}
func (b *Button) SetDisabled(disabled bool) *Button {
	b.disabled = disabled
	if disabled {
		b.SetBackgroundColor(color.RGBA{R: 108, G: 117, B: 125, A: 255})
	} else {
		b.SetBackgroundColor(b.normalColor)
	}
	return b
}
func (b *Button) SetText(text string) *Button {
	b.text.SetContent(text)
	return b
}
func (b *Button) SetTextColor(clr color.Color) *Button {
	b.text.SetColor(clr)
	return b
}
func (b *Button) SetBackgroundColors(normal, hover, pressed color.Color) *Button {
	b.normalColor = normal
	b.hoverColor = hover
	b.pressedColor = pressed
	b.SetBackgroundColor(b.normalColor)
	return b
}
func (b *Button) SetWidth(width float32) *Button {
	b.GetElement().Yoga.StyleSetWidth(width)
	return b
}
func (b *Button) SetHeight(height float32) *Button {
	b.GetElement().Yoga.StyleSetHeight(height)
	return b
}
func (b *Button) SetMargin(edge yoga.Edge, value float32) *Button {
	b.GetElement().Yoga.StyleSetMargin(edge, value)
	return b
}
func (b *Button) SetPadding(edge yoga.Edge, value float32) *Button {
	b.GetElement().Yoga.StyleSetPadding(edge, value)
	return b
}
func (b *Button) SetJustifyContent(justify yoga.Justify) *Button {
	b.GetElement().Yoga.StyleSetJustifyContent(justify)
	return b
}
func (b *Button) SetAlignItems(align yoga.Align) *Button {
	b.GetElement().Yoga.StyleSetAlignItems(align)
	return b
}
func (b *Button) SetBackgroundColor(clr color.Color) *Button {
	b.GetElement().BackgroundColor = clr
	return b
}
func (b *Button) SetBorderRadius(radius float32) *Button {
	b.GetElement().BorderRadius = core.BorderRadius{
		TopLeft: radius, TopRight: radius,
		BottomRight: radius, BottomLeft: radius,
	}
	return b
}

// SyncFrom 同步按钮属性。
func (b *Button) SyncFrom(other core.Host) {
	if o, ok := other.(*Button); ok {
		b.normalColor = o.normalColor
		b.hoverColor = o.hoverColor
		b.pressedColor = o.pressedColor
		b.disabled = o.disabled
		b.onClick = o.onClick
		if b.text != nil && o.text != nil {
			b.text.Content = o.text.Content
			b.text.cachedLayout = nil
		}
	}
}
