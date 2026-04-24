package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// FloatButtonState represents float button interaction state.
type FloatButtonState int

const (
	FloatButtonNormal FloatButtonState = iota
	FloatButtonHover
	FloatButtonPressed
)

// FloatButton is a fixed-position circular floating button.
type FloatButton struct {
	View
	state        FloatButtonState
	onClick      func()
	normalColor  color.Color
	hoverColor   color.Color
	pressedColor color.Color
	textColor    color.Color
	labelEl      *Text
}

// NewFloatButton creates a floating button.
func NewFloatButton(label string) *FloatButton {
	fb := &FloatButton{
		normalColor:  color.RGBA{R: 0, G: 123, B: 255, A: 255},
		hoverColor:   color.RGBA{R: 70, G: 130, B: 180, A: 255},
		pressedColor: color.RGBA{R: 30, G: 144, B: 255, A: 255},
		textColor:    color.White,
	}
	fb.Init(fb)
	fb.SetPositionType(yoga.PositionTypeAbsolute)
	fb.SetPosition(yoga.EdgeRight, 24)
	fb.SetPosition(yoga.EdgeBottom, 24)
	fb.SetWidth(56)
	fb.SetHeight(56)
	fb.SetBorderRadius(28)
	fb.SetBackgroundColor(fb.normalColor)
	fb.SetJustifyContent(yoga.JustifyCenter)
	fb.SetAlignItems(yoga.AlignCenter)

	fb.labelEl = NewText(label).SetColor(fb.textColor)
	fb.AppendChild(fb.labelEl)
	return fb
}

// ElementType returns type identifier.
func (fb *FloatButton) ElementType() string { return "FloatButton" }

// Draw renders the float button background via View.Draw.
func (fb *FloatButton) Draw(screen *ebiten.Image) {
	fb.View.Draw(screen)
}

// Update handles hover state per frame.
func (fb *FloatButton) Update() error {
	mx, my := ebiten.CursorPosition()
	bounds := fb.GetBounds()
	hovered := float32(mx) >= bounds.X && float32(mx) < bounds.X+bounds.Width &&
		float32(my) >= bounds.Y && float32(my) < bounds.Y+bounds.Height

	newState := FloatButtonNormal
	if hovered {
		newState = FloatButtonHover
	}
	if newState != fb.state {
		fb.state = newState
		var bg color.Color
		switch fb.state {
		case FloatButtonHover:
			bg = fb.hoverColor
		case FloatButtonPressed:
			bg = fb.pressedColor
		default:
			bg = fb.normalColor
		}
		fb.SetBackgroundColor(bg)
	}
	return nil
}

// HandleEvent processes click events.
func (fb *FloatButton) HandleEvent(e *core.Event) bool {
	if e.Type == core.EventClick {
		fb.state = FloatButtonPressed
		fb.SetBackgroundColor(fb.pressedColor)
		if fb.onClick != nil {
			fb.onClick()
		}
		return true
	}
	return false
}

// SetOnClick sets the click callback.
func (fb *FloatButton) SetOnClick(fn func()) *FloatButton {
	fb.onClick = fn
	return fb
}

// SetBackgroundColor sets all state background colors.
func (fb *FloatButton) SetBackgroundColor(clr color.Color) *FloatButton {
	fb.normalColor = clr
	fb.hoverColor = clr
	fb.pressedColor = clr
	fb.View.SetBackgroundColor(clr)
	return fb
}

// SetBackgroundColors sets normal, hover and pressed colors.
func (fb *FloatButton) SetBackgroundColors(normal, hover, pressed color.Color) *FloatButton {
	fb.normalColor = normal
	fb.hoverColor = hover
	fb.pressedColor = pressed
	fb.View.SetBackgroundColor(normal)
	return fb
}

// SetTextColor sets the label text color.
func (fb *FloatButton) SetTextColor(clr color.Color) *FloatButton {
	fb.textColor = clr
	fb.labelEl.SetColor(clr)
	return fb
}
