package components

import (
	"image/color"

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
	b := &Button{
		state:        ButtonNormal,
		normalColor:  color.RGBA{R: 240, G: 240, B: 240, A: 255},
		hoverColor:   color.RGBA{R: 220, G: 220, B: 220, A: 255},
		pressedColor: color.RGBA{R: 200, G: 200, B: 200, A: 255},
		textColor:    color.Black,
	}
	b.Init(b)
	b.BaseElement.SetPadding(yoga.EdgeAll, 8)
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
}

// Update handles hover state per frame.
func (b *Button) Update() error {
	if b.disabled {
		return nil
	}
	mx, my := ebiten.CursorPosition()
	bounds := b.GetBounds()
	hovered := float32(mx) >= bounds.X && float32(mx) < bounds.X+bounds.Width &&
		float32(my) >= bounds.Y && float32(my) < bounds.Y+bounds.Height

	newState := ButtonNormal
	if hovered {
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
	if e.Type == core.EventClick {
		b.state = ButtonPressed
		b.Mark(core.FlagNeedDraw)
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
