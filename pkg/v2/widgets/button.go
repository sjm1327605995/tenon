package widgets

import (
	"image/color"

	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
)

// ButtonVariant defines button visual variants.
type ButtonVariant int

const (
	ButtonDefault ButtonVariant = iota
	ButtonSecondary
	ButtonOutline
	ButtonGhost
	ButtonDestructive
	ButtonLink
)

// ButtonWidget is a button component.
type ButtonWidget struct {
	ui.BaseWidget
	label    string
	variant  ButtonVariant
	onClick  func()
	disabled bool
	loading  bool
}

// Button creates a button Widget.
func Button(label string) ButtonWidget {
	return ButtonWidget{label: label, variant: ButtonDefault}
}

func (b ButtonWidget) Variantf(v ButtonVariant) ButtonWidget {
	b.variant = v
	return b
}

func (b ButtonWidget) OnTap(fn func()) ButtonWidget {
	b.onClick = fn
	return b
}

func (b ButtonWidget) SetDisabled(v bool) ButtonWidget {
	b.disabled = v
	return b
}

func (b ButtonWidget) SetLoading(v bool) ButtonWidget {
	b.loading = v
	return b
}

func (b ButtonWidget) CreateElement() ui.Element {
	e := &ButtonElement{}
	e.SingleChildRenderObjectElement.RenderObjectElement.BaseElement.Init(e, b)
	return e
}

// ButtonElement is the Element corresponding to ButtonWidget.
type ButtonElement struct {
	ui.SingleChildRenderObjectElement
}

func (e *ButtonElement) CreateRenderObject() render.RenderObject {
	t := ui.GetTheme()
	r := render.NewRenderButton()
	r.SetBorderRadiusDetailed(render.UniformBorderRadius(t.ButtonBorderRadius))

	e.applyButtonVariant(r, e.GetWidget().(ButtonWidget).variant, t)

	r.DisabledColor = t.ButtonDisabledColor
	r.Loading = e.GetWidget().(ButtonWidget).loading
	if e.GetWidget().(ButtonWidget).disabled {
		r.SetState(render.ButtonStateDisabled)
	}
	r.SetOnClick(e.GetWidget().(ButtonWidget).onClick)

	if e.GetWidget().(ButtonWidget).label == "" {
		r.StyleSetWidthPercent(100)
	}

	r.SetOnMouseEnter(func() {
		w := e.GetWidget().(ButtonWidget)
		if !w.disabled && !w.loading {
			r.SetState(render.ButtonStateHover)
		}
	})
	r.SetOnMouseLeave(func() {
		w := e.GetWidget().(ButtonWidget)
		if !w.disabled && !w.loading {
			r.SetState(render.ButtonStateNormal)
		}
	})
	r.SetOnMouseDown(func() {
		w := e.GetWidget().(ButtonWidget)
		if !w.disabled && !w.loading {
			r.SetState(render.ButtonStatePressed)
		}
	})
	r.SetOnMouseUp(func() {
		w := e.GetWidget().(ButtonWidget)
		if !w.disabled && !w.loading {
			r.SetState(render.ButtonStateHover)
		}
	})

	return r
}

func (e *ButtonElement) applyButtonVariant(r *render.RenderButton, variant ButtonVariant, t *ui.Theme) {
	switch variant {
	case ButtonDefault:
		r.NormalColor = t.ButtonNormalColor
		r.HoverColor = t.ButtonHoverColor
		r.PressedColor = t.ButtonPressedColor
		r.TextColor = t.ButtonTextColor
	case ButtonSecondary:
		r.NormalColor = t.SecondaryColor
		r.HoverColor = render.Darken(t.SecondaryColor, 20)
		r.PressedColor = render.Darken(t.SecondaryColor, 40)
		r.TextColor = t.SecondaryForegroundColor
	case ButtonOutline:
		r.IsOutline = true
		r.NormalColor = nil
		r.HoverColor = t.AccentColor
		r.PressedColor = t.AccentColor
		r.BorderColor = t.BorderColor
		r.BorderWidth = 1
		r.TextColor = t.TextColor
	case ButtonGhost:
		r.IsGhost = true
		r.NormalColor = nil
		r.HoverColor = t.AccentColor
		r.PressedColor = t.AccentColor
		r.TextColor = t.TextColor
	case ButtonDestructive:
		r.NormalColor = t.DestructiveColor
		r.HoverColor = render.Lighten(t.DestructiveColor, 20)
		r.PressedColor = render.Lighten(t.DestructiveColor, 40)
		r.TextColor = t.DestructiveForegroundColor
	case ButtonLink:
		r.IsLink = true
		r.NormalColor = nil
		r.HoverColor = nil
		r.PressedColor = nil
		r.TextColor = t.PrimaryColor
		r.HoverTextColor = t.TextMutedColor
	}
}

func (e *ButtonElement) UpdateRenderObject(oldWidget ui.Widget) {
	w := e.GetWidget().(ButtonWidget)
	r := e.GetRenderObject().(*render.RenderButton)
	t := ui.GetTheme()

	old := oldWidget.(ButtonWidget)

	// Reconfigure colors when variant changes
	if old.variant != w.variant {
		r.IsOutline = false
		r.IsGhost = false
		r.IsLink = false
		r.NormalColor = nil
		r.HoverColor = nil
		r.PressedColor = nil
		r.BorderColor = nil
		r.BorderWidth = 0
		r.HoverTextColor = nil

		e.applyButtonVariant(r, w.variant, t)
		r.MarkNeedsPaint()
	}

	r.SetOnClick(w.onClick)
	if old.loading != w.loading {
		r.Loading = w.loading
		r.MarkNeedsPaint()
	}
	if old.disabled != w.disabled {
		if w.disabled {
			r.SetState(render.ButtonStateDisabled)
		} else {
			r.SetState(render.ButtonStateNormal)
		}
	}
}

func (e *ButtonElement) UpdateChild(oldWidget ui.Widget) {
	w := e.GetWidget().(ButtonWidget)
	textColor := buttonTextColor(w.variant, ui.GetTheme())
	textWidget := Text(w.label).FontSize(ui.GetTheme().FontSizeBase).Color(textColor)
	e.Child = ui.UpdateChild(e, e.Child, textWidget)
}

func (e *ButtonElement) Mount(parent ui.Element, slot int) {
	e.RenderObject = e.CreateRenderObject()
	e.SingleChildRenderObjectElement.Mount(parent, slot)
	w := e.GetWidget().(ButtonWidget)
	textColor := buttonTextColor(w.variant, ui.GetTheme())
	textWidget := Text(w.label).FontSize(ui.GetTheme().FontSizeBase).Color(textColor)
	e.Child = ui.UpdateChild(e, nil, textWidget)
}

// buttonTextColor returns the text color for a given button variant.
func buttonTextColor(variant ButtonVariant, t *ui.Theme) color.Color {
	switch variant {
	case ButtonLink:
		return t.PrimaryColor
	case ButtonSecondary:
		return t.SecondaryForegroundColor
	case ButtonDestructive:
		return t.DestructiveForegroundColor
	case ButtonGhost, ButtonOutline:
		return t.TextColor
	default:
		return t.ButtonTextColor
	}
}
