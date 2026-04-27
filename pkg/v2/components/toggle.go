package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// ToggleVariant defines the visual style of a toggle.
type ToggleVariant int

const (
	ToggleDefault ToggleVariant = iota
	ToggleOutline
)

// Toggle is a button that can be on or off.
type Toggle struct {
	core.BaseElement
	pressed      bool
	hovered      bool
	onChange     func(pressed bool)
	variant      ToggleVariant
	normalColor  color.Color
	hoverColor   color.Color
	activeColor  color.Color
	textColor    color.Color
	borderColor  color.Color
	borderRadius float32
	labelEl      *Text
}

// NewToggle creates a toggle button.
func NewToggle(label string) *Toggle {
	theme := core.GetTheme()
	t := &Toggle{
		variant:      ToggleDefault,
		normalColor:  nil,
		hoverColor:   theme.AccentColor,
		activeColor:  theme.TextColor,
		textColor:    theme.TextColor,
		borderRadius: theme.BorderRadius / 2,
	}
	t.Init(t)
	t.SetPadding(yoga.EdgeHorizontal, 12)
	t.SetPadding(yoga.EdgeVertical, 8)
	t.SetJustifyContent(yoga.JustifyCenter)
	t.SetAlignItems(yoga.AlignCenter)

	t.labelEl = NewText(label).SetColor(t.textColor)
	t.labelEl.SetFontSize(14)
	t.AppendChild(t.labelEl)
	return t
}

// ElementType returns type identifier.
func (t *Toggle) ElementType() string { return "Toggle" }

// Draw renders the toggle background.
func (t *Toggle) Draw(screen *ebiten.Image) {
	if !t.IsVisible() {
		return
	}
	bounds := t.GetBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}
	var bg color.Color
	if t.pressed {
		bg = t.activeColor
	} else if t.hovered {
		bg = t.hoverColor
	} else {
		bg = t.normalColor
	}
	br := core.BorderRadius{TopLeft: t.borderRadius, TopRight: t.borderRadius, BottomRight: t.borderRadius, BottomLeft: t.borderRadius}
	if bg != nil {
		drawRoundedRectFill(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, br, bg)
	}
	if t.variant == ToggleOutline || t.pressed {
		var borderClr color.Color
		if t.variant == ToggleOutline {
			borderClr = t.borderColor
			if borderClr == nil {
				borderClr = core.GetTheme().BorderColor
			}
		} else if t.pressed {
			borderClr = t.activeColor
		}
		if borderClr != nil {
			drawRoundedRectStroke(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, br, 1, borderClr)
		}
	}
}

// Update syncs hover state.
func (t *Toggle) Update() error {
	newHovered := false
	if t.GetEngine() != nil {
		hover := t.GetEngine().GetHoverTarget()
		for el := hover; el != nil; el = el.GetParent() {
			if el == t {
				newHovered = true
				break
			}
		}
	}
	if newHovered != t.hovered {
		t.hovered = newHovered
		t.Mark(core.FlagNeedDraw)
	}
	return nil
}

// HandleEvent processes click events.
func (t *Toggle) HandleEvent(e *core.Event) bool {
	switch e.Type {
	case core.EventMouseDown:
		return true
	case core.EventMouseUp:
		return true
	case core.EventClick:
		t.pressed = !t.pressed
		if t.onChange != nil {
			t.onChange(t.pressed)
		}
		t.updateLabelColor()
		t.Mark(core.FlagNeedDraw)
		return true
	}
	return false
}

func (t *Toggle) updateLabelColor() {
	if t.pressed {
		t.labelEl.SetColor(core.GetTheme().BackgroundColor)
	} else {
		t.labelEl.SetColor(t.textColor)
	}
}

// SetPressed sets the pressed state.
func (t *Toggle) SetPressed(v bool) *Toggle {
	t.pressed = v
	t.updateLabelColor()
	t.Mark(core.FlagNeedDraw)
	return t
}

// IsPressed returns the current pressed state.
func (t *Toggle) IsPressed() bool { return t.pressed }

// SetOnChange sets the change callback.
func (t *Toggle) SetOnChange(fn func(pressed bool)) *Toggle {
	t.onChange = fn
	return t
}

// SetVariant sets the toggle variant.
func (t *Toggle) SetVariant(v ToggleVariant) *Toggle {
	t.variant = v
	t.Mark(core.FlagNeedDraw)
	return t
}

// SetColors sets custom colors.
func (t *Toggle) SetColors(normal, active, text color.Color) *Toggle {
	t.normalColor = normal
	t.activeColor = active
	t.textColor = text
	t.updateLabelColor()
	t.Mark(core.FlagNeedDraw)
	return t
}
