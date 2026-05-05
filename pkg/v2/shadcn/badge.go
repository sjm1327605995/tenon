package shadcn

import (
	"fmt"

	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
	"github.com/sjm1327605995/tenon/pkg/v2/widgets"
)

// BadgeVariant defines the visual style of a badge.
type BadgeVariant int

const (
	BadgeDefault BadgeVariant = iota
	BadgeSecondary
	BadgeOutline
	BadgeDestructive
)

// Badge creates a small status indicator.
func Badge(content string, variant BadgeVariant) ui.Widget {
	return badgeWidget{content: content, variant: variant, maxCount: 99}
}

// DotBadge creates a small dot indicator without text.
func DotBadge() ui.Widget {
	return badgeWidget{dotMode: true, variant: BadgeDefault}
}

// CountBadge creates a badge from a numeric count with overflow.
func CountBadge(count int) ui.Widget {
	return badgeWidget{content: formatCount(count, 99), variant: BadgeDefault, maxCount: 99}
}

// internal widget type — not exported, keeps the public API as pure functions.
type badgeWidget struct {
	ui.BaseWidget
	content  string
	variant  BadgeVariant
	dotMode  bool
	maxCount int
}

func (b badgeWidget) CreateElement() ui.Element {
	e := &badgeElement{}
	e.SingleChildRenderObjectElement.RenderObjectElement.BaseElement.Init(e, b)
	return e
}

type badgeElement struct {
	ui.SingleChildRenderObjectElement
	ro *render.RenderBox
}

func (e *badgeElement) CreateRenderObject() render.RenderObject {
	r := render.NewRenderBox()
	applyBadgeProps(r, badgeWidget{}, e.GetWidget().(badgeWidget), true)
	return r
}

func (e *badgeElement) UpdateRenderObject(oldWidget ui.Widget) {
	r := e.ro
	old := oldWidget.(badgeWidget)
	w := e.GetWidget().(badgeWidget)
	applyBadgeProps(r, old, w, false)
}

func (e *badgeElement) Mount(parent ui.Element, slot int) {
	e.ro = e.CreateRenderObject().(*render.RenderBox)
	e.RenderObject = e.ro
	e.SingleChildRenderObjectElement.Mount(parent, slot)
	w := e.GetWidget().(badgeWidget)
	if !w.dotMode && w.content != "" {
		child := widgets.Text(w.content).FontSize(12).Color(getBadgeTextColor(w.variant))
		e.Child = ui.UpdateChild(e, nil, child)
	}
}

func (e *badgeElement) UpdateChild(oldWidget ui.Widget) {
	w := e.GetWidget().(badgeWidget)
	if w.dotMode {
		e.Child = ui.UpdateChild(e, e.Child, nil)
		return
	}
	child := widgets.Text(w.content).FontSize(12).Color(getBadgeTextColor(w.variant))
	e.Child = ui.UpdateChild(e, e.Child, child)
}

func applyBadgeProps(r *render.RenderBox, old, w badgeWidget, force bool) {
	if force || old.dotMode != w.dotMode {
		if w.dotMode {
			r.StyleSetWidth(8)
			r.StyleSetHeight(8)
			r.SetBorderRadius(4)
			r.StyleSetAlignSelf(ui.AlignCenter)
		} else {
			r.StyleSetWidthAuto()
			r.StyleSetHeightAuto()
			r.SetBorderRadius(999)
			r.StyleSetPadding(ui.EdgeHorizontal, 10)
			r.StyleSetPadding(ui.EdgeVertical, 4)
		}
	} else if !w.dotMode && (old.dotMode || old.content != w.content) {
		// content changed but not dot mode: ensure padding is set
		r.StyleSetPadding(ui.EdgeHorizontal, 10)
		r.StyleSetPadding(ui.EdgeVertical, 4)
	}

	r.StyleSetFlexDirection(ui.FlexDirectionRow)
	r.StyleSetJustifyContent(ui.JustifyCenter)
	r.StyleSetAlignItems(ui.AlignCenter)

	bg, border := getBadgeColors(w.variant)
	oldBg, oldBorder := getBadgeColors(old.variant)
	if force || !render.ColorPtrEquals(oldBg, bg) {
		r.SetBackgroundColor(bg)
	}
	if force || !render.ColorPtrEquals(oldBorder, border) {
		if border != nil {
			r.SetBorderColor(border)
			r.SetBorderWidth(1)
		} else {
			r.SetBorderWidth(0)
		}
	}
}

func getBadgeColors(v BadgeVariant) (*render.Color, *render.Color) {
	theme := ui.GetTheme()
	switch v {
	case BadgeDefault:
		return newColor(theme.PrimaryColor), nil
	case BadgeSecondary:
		return newColor(theme.SecondaryColor), nil
	case BadgeOutline:
		return nil, newColor(theme.BorderColor)
	case BadgeDestructive:
		return newColor(theme.DestructiveColor), nil
	}
	return newColor(theme.PrimaryColor), nil
}

func getBadgeTextColor(v BadgeVariant) *render.Color {
	theme := ui.GetTheme()
	switch v {
	case BadgeDefault:
		return newColor(theme.PrimaryForegroundColor)
	case BadgeSecondary:
		return newColor(theme.SecondaryForegroundColor)
	case BadgeOutline:
		return newColor(theme.TextColor)
	case BadgeDestructive:
		return newColor(theme.DestructiveForegroundColor)
	}
	return newColor(theme.PrimaryForegroundColor)
}

func formatCount(count, max int) string {
	if count > max {
		return fmt.Sprintf("%d+", max)
	}
	return fmt.Sprintf("%d", count)
}
