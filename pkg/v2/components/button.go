package components

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/v2/core"
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
	borderRadius float32
	loading      bool
	disabled     bool
	labelEl      *Text
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
		borderRadius: theme.ButtonBorderRadius,
	}
	b.Init(b)
	b.BaseElement.SetPadding(yoga.EdgeHorizontal, 16)
	b.BaseElement.SetPadding(yoga.EdgeVertical, 10)
	b.SetJustifyContent(yoga.JustifyCenter)
	b.SetAlignItems(yoga.AlignCenter)
	b.labelEl = NewText(label).SetColor(b.textColor)
	b.AppendChild(b.labelEl)
	return b
}

// ElementType returns type identifier.
func (b *Button) ElementType() string { return "Button" }

// Draw renders button background.
func (b *Button) Draw(screen *ebiten.Image) {
	if !b.IsVisible() {
		return
	}
	bounds := b.GetBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	// Background color based on state
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

	br := core.BorderRadius{
		TopLeft: b.borderRadius, TopRight: b.borderRadius,
		BottomRight: b.borderRadius, BottomLeft: b.borderRadius,
	}

	if bg != nil {
		drawRoundedRectFill(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, br, bg)
	}
	if borderClr != nil && !b.disabled {
		drawRoundedRectStroke(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, br, 1, borderClr)
	}

	// Loading spinner
	if b.loading {
		b.drawLoading(screen, bounds)
	}
}

// Update syncs hover state with engine's hoverTarget.
func (b *Button) Update() error {
	if b.disabled {
		return nil
	}
	newState := ButtonNormal
	if b.state == ButtonPressed && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		newState = ButtonPressed
	} else if b.GetEngine() != nil && b.isHovered() {
		newState = ButtonHover
	}
	if newState != b.state {
		b.state = newState
		b.Mark(core.FlagNeedDraw)
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
		b.Mark(core.FlagNeedDraw)
		return true
	case core.EventMouseUp:
		bounds := b.GetBounds()
		if e.X >= bounds.X && e.X < bounds.X+bounds.Width &&
			e.Y >= bounds.Y && e.Y < bounds.Y+bounds.Height {
			b.state = ButtonHover
		} else {
			b.state = ButtonNormal
		}
		b.Mark(core.FlagNeedDraw)
		return true
	case core.EventClick:
		if b.onClick != nil {
			b.onClick()
		}
		return true
	}
	return false
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
	b.Mark(core.FlagNeedDraw)
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
	b.Mark(core.FlagNeedDraw)
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
	b.Mark(core.FlagNeedDraw)
	return b
}

func (b *Button) SetBorderRadius(r float32) *Button {
	b.borderRadius = r
	b.Mark(core.FlagNeedDraw)
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

func (b *Button) DebugProps() map[string]interface{} {
	props := make(map[string]interface{})
	variantNames := []string{"default", "secondary", "outline", "ghost", "destructive", "link"}
	if int(b.variant) < len(variantNames) {
		props["variant"] = variantNames[b.variant]
	}
	stateNames := []string{"normal", "hover", "pressed"}
	if int(b.state) < len(stateNames) {
		props["state"] = stateNames[b.state]
	}
	if b.labelEl != nil {
		props["label"] = b.labelEl.content
	}
	if b.normalColor != nil {
		props["normalColor"] = colorToCSS(b.normalColor)
	}
	if b.textColor != nil {
		props["textColor"] = colorToCSS(b.textColor)
	}
	if b.borderColor != nil {
		props["borderColor"] = colorToCSS(b.borderColor)
	}
	props["borderRadius"] = b.borderRadius
	props["disabled"] = b.disabled
	props["loading"] = b.loading
	return props
}
