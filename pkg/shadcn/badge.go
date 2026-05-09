package shadcn

import (
	"fmt"

	"github.com/sjm1327605995/tenon/pkg/render"
	"github.com/sjm1327605995/tenon/pkg/engine"
	"github.com/sjm1327605995/tenon/pkg/widgets"
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
func Badge(content string, variant BadgeVariant) engine.Widget {
	return badgeWidget{content: content, variant: variant, maxCount: 99}
}

// DotBadge creates a small dot indicator without text.
func DotBadge() engine.Widget {
	return badgeWidget{dotMode: true, variant: BadgeDefault}
}

// CountBadge creates a badge from a numeric count with overflow.
func CountBadge(count int) engine.Widget {
	return badgeWidget{content: formatCount(count, 99), variant: BadgeDefault, maxCount: 99}
}

// internal widget type — not exported, keeps the public API as pure functions.
type badgeWidget struct {
	engine.BaseWidget
	content  string
	variant  BadgeVariant
	dotMode  bool
	maxCount int
}

func (b badgeWidget) CreateElement() engine.Element {
	return engine.NewSingleChildRenderObjectElement(b)
}

// CreateRenderObject implements RenderObjectFactory.
func (b badgeWidget) CreateRenderObject(element engine.Element) render.RenderObject {
	r := render.NewRenderBox()
	applyBadgeProps(r, badgeWidget{}, b, true)
	return r
}

// UpdateRenderObject implements RenderObjectUpdater.
func (b badgeWidget) UpdateRenderObject(ro render.RenderObject, oldWidget engine.Widget) {
	old := oldWidget.(badgeWidget)
	r := ro.(*render.RenderBox)
	applyBadgeProps(r, old, b, false)
}

// GetChildWidget implements SingleChildProvider.
func (b badgeWidget) GetChildWidget() engine.Widget {
	if b.dotMode || b.content == "" {
		return nil
	}
	return widgets.Text(b.content).FontSize(12).Color(getBadgeTextColor(b.variant))
}

func applyBadgeProps(r *render.RenderBox, old, w badgeWidget, force bool) {
	if force || old.dotMode != w.dotMode {
		if w.dotMode {
			r.StyleSetWidth(8)
			r.StyleSetHeight(8)
			r.SetBorderRadius(4)
			r.StyleSetAlignSelf(engine.AlignCenter)
		} else {
			r.StyleSetWidthAuto()
			r.StyleSetHeightAuto()
			r.SetBorderRadius(999)
			r.StyleSetPadding(engine.EdgeHorizontal, 10)
			r.StyleSetPadding(engine.EdgeVertical, 4)
		}
	} else if !w.dotMode && (old.dotMode || old.content != w.content) {
		// content changed but not dot mode: ensure padding is set
		r.StyleSetPadding(engine.EdgeHorizontal, 10)
		r.StyleSetPadding(engine.EdgeVertical, 4)
	}

	r.StyleSetFlexDirection(engine.FlexDirectionRow)
	r.StyleSetJustifyContent(engine.JustifyCenter)
	r.StyleSetAlignItems(engine.AlignCenter)

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
	theme := engine.GetTheme()
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
	theme := engine.GetTheme()
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
