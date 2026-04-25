package components

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// ButtonState represents button interaction state.
type ButtonState int

const (
	ButtonNormal ButtonState = iota
	ButtonHover
	ButtonPressed
)

// Button is an interactive button element.
type Button struct {
	core.BaseElement
	state        ButtonState
	onClick      func()
	normalColor  color.Color
	hoverColor   color.Color
	pressedColor color.Color
	textColor    color.Color
	loading      bool
	disabled     bool
	labelEl      *Text
}

// NewButton creates a button.
func NewButton(label string) *Button {
	theme := core.GetTheme()
	b := &Button{
		state:        ButtonNormal,
		normalColor:  theme.ButtonNormalColor,
		hoverColor:   theme.ButtonHoverColor,
		pressedColor: theme.ButtonPressedColor,
		textColor:    theme.ButtonTextColor,
	}
	b.Init(b)
	b.BaseElement.SetPadding(yoga.EdgeAll, 8)
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
	switch b.state {
	case ButtonHover:
		bg = b.hoverColor
	case ButtonPressed:
		bg = b.pressedColor
	default:
		bg = b.normalColor
	}

	if b.disabled {
		bg = color.RGBA{R: 200, G: 200, B: 200, A: 255}
	}

	vector.FillRect(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, bg, false)

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
	} else if b.GetEngine() != nil && b.GetEngine().GetHoverTarget() == b {
		newState = ButtonHover
	}
	if newState != b.state {
		b.state = newState
		b.Mark(core.FlagNeedDraw)
	}
	return nil
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

func (b *Button) SetColors(normal, hover, pressed color.Color) *Button {
	b.normalColor = normal
	b.hoverColor = hover
	b.pressedColor = pressed
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
