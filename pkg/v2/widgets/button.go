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
	return ui.NewSingleChildRenderObjectElement(b)
}

// CreateRenderObject implements RenderObjectFactory.
func (b ButtonWidget) CreateRenderObject(element ui.Element) render.RenderObject {
	t := ui.GetTheme()
	r := render.NewRenderButton()
	r.SetBorderRadiusDetailed(render.UniformBorderRadius(t.ButtonBorderRadius))

	applyButtonVariant(r, b.variant, t)

	r.DisabledColor = t.ButtonDisabledColor
	r.Loading = b.loading
	if b.disabled {
		r.SetState(render.ButtonStateDisabled)
	}
	r.SetOnClick(b.onClick)

	if b.label == "" {
		r.StyleSetWidthPercent(100)
	}

	r.SetOnMouseEnter(func() {
		w, ok := element.GetWidget().(ButtonWidget)
		if ok && !w.disabled && !w.loading {
			r.SetState(render.ButtonStateHover)
		}
	})
	r.SetOnMouseLeave(func() {
		w, ok := element.GetWidget().(ButtonWidget)
		if ok && !w.disabled && !w.loading {
			r.SetState(render.ButtonStateNormal)
		}
	})
	r.SetOnMouseDown(func() {
		w, ok := element.GetWidget().(ButtonWidget)
		if ok && !w.disabled && !w.loading {
			r.SetState(render.ButtonStatePressed)
		}
	})
	r.SetOnMouseUp(func() {
		w, ok := element.GetWidget().(ButtonWidget)
		if ok && !w.disabled && !w.loading {
			r.SetState(render.ButtonStateHover)
		}
	})

	return r
}

func applyButtonVariant(r *render.RenderButton, variant ButtonVariant, t *ui.Theme) {
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

// UpdateRenderObject implements RenderObjectUpdater.
func (b ButtonWidget) UpdateRenderObject(ro render.RenderObject, oldWidget ui.Widget) {
	r := ro.(*render.RenderButton)
	old := oldWidget.(ButtonWidget)
	t := ui.GetTheme()

	// Reconfigure colors when variant changes
	if old.variant != b.variant {
		r.IsOutline = false
		r.IsGhost = false
		r.IsLink = false
		r.NormalColor = nil
		r.HoverColor = nil
		r.PressedColor = nil
		r.BorderColor = nil
		r.BorderWidth = 0
		r.HoverTextColor = nil

		applyButtonVariant(r, b.variant, t)
		r.MarkNeedsPaint()
	}

	r.SetOnClick(b.onClick)
	if old.loading != b.loading {
		r.Loading = b.loading
		r.MarkNeedsPaint()
	}
	if old.disabled != b.disabled {
		if b.disabled {
			r.SetState(render.ButtonStateDisabled)
		} else {
			r.SetState(render.ButtonStateNormal)
		}
	}
}

// GetChildWidget implements SingleChildProvider.
func (b ButtonWidget) GetChildWidget() ui.Widget {
	textColor := buttonTextColor(b.variant, ui.GetTheme())
	return Text(b.label).FontSize(ui.GetTheme().FontSizeBase).Color(textColor)
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
