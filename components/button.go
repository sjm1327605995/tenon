package components

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/internal/core"
	"github.com/sjm1327605995/tenon/internal/native"
	"github.com/sjm1327605995/tenon/yoga"
)

// ButtonVariant represents the visual style of a button.
type ButtonVariant int

const (
	ButtonDefault ButtonVariant = iota
	ButtonSecondary
	ButtonOutline
	ButtonGhost
	ButtonDestructive
	ButtonLink
)

// ButtonState represents button interaction state.
type ButtonState int

const (
	ButtonNormal ButtonState = iota
	ButtonHover
	ButtonPressed
)

// Button is an interactive button element with Shadcn/UI-inspired styling.
// Visual properties (background, border, radius) are delegated to the internal
// bgEl (native.View) so the debugger preview can reflect them without custom Draw.
type Button struct {
	core.BaseElement
	state        ButtonState
	onClick      func()
	variant      ButtonVariant
	normalColor  color.Color
	hoverColor   color.Color
	pressedColor color.Color
	textColor    color.Color
	borderColor  color.Color
	BorderRadius float32
	loading      bool
	disabled     bool
	bgEl         *native.View
	labelEl      *native.Text
}

// NewButton creates a button with default (primary) variant.
func NewButton(label string) *Button {
	theme := core.GetTheme()
	b := &Button{
		state:        ButtonNormal,
		variant:      ButtonDefault,
		normalColor:  theme.ButtonNormalColor,
		hoverColor:   theme.ButtonHoverColor,
		pressedColor: theme.ButtonPressedColor,
		textColor:    theme.ButtonTextColor,
		BorderRadius: theme.ButtonBorderRadius,
	}
	b.Init(b)

	b.bgEl = native.NewView()
	b.bgEl.SetFlexGrow(1)
	b.bgEl.SetPadding(yoga.EdgeHorizontal, 16)
	b.bgEl.SetPadding(yoga.EdgeVertical, 10)
	b.bgEl.SetJustifyContent(yoga.JustifyCenter)
	b.bgEl.SetAlignItems(yoga.AlignCenter)
	b.bgEl.SetBorderRadius(theme.ButtonBorderRadius)
	b.bgEl.SetBackgroundColor(b.normalColor)

	b.labelEl = native.NewText(label).SetColor(b.textColor)
	b.bgEl.Add(b.labelEl)
	b.AppendChild(b.bgEl)
	return b
}

// ElementType returns type identifier.
func (b *Button) ElementType() string { return "Button" }

// SyncFrom 同步新 Button 的属性到当前 Element（声明式重建）。
func (b *Button) SyncFrom(src core.Element) {
	other, ok := src.(*Button)
	if !ok {
		return
	}
	sb := &core.SyncBuilder{}
	core.SyncField(sb, &b.loading, other.loading)
	core.SyncField(sb, &b.disabled, other.disabled)
	core.SyncColor(sb, &b.normalColor, other.normalColor)
	core.SyncColor(sb, &b.hoverColor, other.hoverColor)
	core.SyncColor(sb, &b.pressedColor, other.pressedColor)
	core.SyncColor(sb, &b.textColor, other.textColor)
	if b.bgEl != nil {
		b.applyBgColors()
	}
	sb.MarkDraw(b)
}

// Draw renders loading spinner only; background/border are handled by bgEl (native.View).
func (b *Button) Draw(screen *ebiten.Image) {
	if b.loading {
		b.drawLoading(screen, b.GetBounds())
	}
}

// Update syncs hover state with engine's hoverTarget and updates bgEl colors.
func (b *Button) Update() error {
	if b.disabled {
		return nil
	}
	newState := ButtonNormal
	if b.state == ButtonPressed && core.IsMouseButtonPressed(core.MouseButtonLeft) {
		newState = ButtonPressed
	} else if b.GetEngine() != nil && b.isHovered() {
		newState = ButtonHover
	}
	if newState != b.state {
		b.state = newState
		b.applyBgColors()
	}
	return nil
}

func (b *Button) isHovered() bool {
	hover := b.GetEngine().GetHoverTarget()
	for el := hover; el != nil; el = el.GetParent() {
		if el == b {
			return true
		}
	}
	return false
}

// HandleEvent processes click events.
func (b *Button) HandleEvent(e *core.Event) bool {
	if b.disabled {
		return false
	}
	switch e.Type {
	case core.EventMouseDown:
		b.state = ButtonPressed
		b.applyBgColors()
		return true
	case core.EventMouseUp:
		bounds := b.GetBounds()
		if e.X >= bounds.X && e.X < bounds.X+bounds.Width &&
			e.Y >= bounds.Y && e.Y < bounds.Y+bounds.Height {
			b.state = ButtonHover
		} else {
			b.state = ButtonNormal
		}
		b.applyBgColors()
		return true
	case core.EventClick:
		if b.onClick != nil {
			b.onClick()
		}
		return true
	}
	return false
}

// applyBgColors updates bgEl background/border based on current state.
func (b *Button) applyBgColors() {
	if b.bgEl == nil {
		return
	}
	var bg color.Color
	var borderClr color.Color
	switch b.state {
	case ButtonHover:
		bg = b.hoverColor
	case ButtonPressed:
		bg = b.pressedColor
	default:
		bg = b.normalColor
	}
	borderClr = b.borderColor
	if b.disabled {
		bg = core.GetTheme().ButtonDisabledColor
		if b.variant == ButtonOutline || b.variant == ButtonGhost || b.variant == ButtonLink {
			bg = nil
		}
	}
	b.bgEl.SetBackgroundColor(bg)
	b.bgEl.SetBorderColor(borderClr)
}

// setVariantColors updates colors based on the current variant and theme.
func (b *Button) setVariantColors() {
	theme := core.GetTheme()
	switch b.variant {
	case ButtonDefault:
		b.normalColor = theme.ButtonNormalColor
		b.hoverColor = theme.ButtonHoverColor
		b.pressedColor = theme.ButtonPressedColor
		b.textColor = theme.ButtonTextColor
		b.borderColor = nil
	case ButtonSecondary:
		b.normalColor = theme.SecondaryColor
		if h, ok := b.normalColor.(color.RGBA); ok {
			darken := func(v uint8) uint8 {
				if v >= 20 {
					return v - 20
				}
				return 0
			}
			b.hoverColor = color.RGBA{R: darken(h.R), G: darken(h.G), B: darken(h.B), A: h.A}
		} else {
			b.hoverColor = theme.AccentColor
		}
		b.pressedColor = theme.BorderColor
		b.textColor = theme.SecondaryForegroundColor
		b.borderColor = nil
	case ButtonOutline:
		b.normalColor = nil
		b.hoverColor = theme.AccentColor
		b.pressedColor = theme.SecondaryColor
		b.textColor = theme.TextColor
		b.borderColor = theme.BorderColor
	case ButtonGhost:
		b.normalColor = nil
		b.hoverColor = theme.AccentColor
		b.pressedColor = theme.SecondaryColor
		b.textColor = theme.TextColor
		b.borderColor = nil
	case ButtonDestructive:
		b.normalColor = theme.DestructiveColor
		b.hoverColor = theme.DestructiveColor
		b.pressedColor = theme.DestructiveColor
		if h, ok := b.normalColor.(color.RGBA); ok {
			b.hoverColor = color.RGBA{R: min(255, h.R+20), G: min(255, h.G+20), B: min(255, h.B+20), A: h.A}
		}
		b.textColor = theme.DestructiveForegroundColor
		b.borderColor = nil
	case ButtonLink:
		b.normalColor = nil
		b.hoverColor = nil
		b.pressedColor = nil
		b.textColor = theme.PrimaryColor
		b.borderColor = nil
	}
	if b.labelEl != nil {
		b.labelEl.SetColor(b.textColor)
	}
	b.applyBgColors()
}

// Chain API

func (b *Button) SetLabel(label string) *Button {
	b.labelEl.SetContent(label)
	return b
}

// SetText is an alias for SetLabel.
func (b *Button) SetText(text string) *Button {
	b.labelEl.SetContent(text)
	return b
}

func (b *Button) SetOnClick(fn func()) *Button {
	b.onClick = fn
	return b
}

func (b *Button) SetLoading(loading bool) *Button {
	b.loading = loading
	b.Mark(core.FlagNeedDraw)
	return b
}

func (b *Button) SetDisabled(disabled bool) *Button {
	b.disabled = disabled
	b.applyBgColors()
	return b
}

func (b *Button) SetVariant(v ButtonVariant) *Button {
	b.variant = v
	b.setVariantColors()
	return b
}

func (b *Button) SetColors(normal, hover, pressed color.Color) *Button {
	b.normalColor = normal
	b.hoverColor = hover
	b.pressedColor = pressed
	b.applyBgColors()
	return b
}

func (b *Button) SetBorderRadius(r float32) *Button {
	b.BorderRadius = r
	if b.bgEl != nil {
		b.bgEl.SetBorderRadius(r)
	}
	return b
}

func (b *Button) drawLoading(screen *ebiten.Image, bounds core.LayoutBounds) {
	cx := bounds.X + bounds.Width/2
	cy := bounds.Y + bounds.Height/2
	r := float32(6)
	stroke := float32(2)
	clr := b.textColor
	if clr == nil {
		clr = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}
	angle := float32(0) // static for simplicity, could animate
	var path vector.Path
	startAngle := angle
	endAngle := angle + 1.2
	path.MoveTo(cx+r*float32(math.Cos(float64(startAngle))), cy+r*float32(math.Sin(float64(startAngle))))
	path.Arc(cx, cy, r, startAngle, endAngle, vector.Clockwise)
	strokeOp := &vector.StrokeOptions{Width: stroke, MiterLimit: 10}
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(clr)
	vector.StrokePath(screen, &path, strokeOp, op)
}
