package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/internal/core"
	"github.com/sjm1327605995/tenon/internal/native"
	"github.com/sjm1327605995/tenon/yoga"
)

// AlertVariant defines the visual style of an alert.
type AlertVariant int

const (
	AlertDefault AlertVariant = iota
	AlertDestructive
)

// Alert is a callout box for important messages.
type Alert struct {
	native.View
	variant     AlertVariant
	titleEl     *native.Text
	descEl      *native.Text
	iconText    *native.Text
	borderColor color.Color
	bgColor     color.Color
}

// NewAlert creates an alert.
func NewAlert(title string) *Alert {
	theme := core.GetTheme()
	a := &Alert{
		variant: AlertDefault,
	}
	a.Init(a)
	a.SetFlexDirection(yoga.FlexDirectionRow)
	a.SetGap(yoga.GutterAll, 12)
	a.SetPadding(yoga.EdgeAll, 16)
	a.SetBorderRadius(theme.BorderRadius)
	a.SetWidthPercent(100)

	a.borderColor = theme.BorderColor
	a.bgColor = theme.CardColor
	//TODO emoji
	a.iconText = native.NewText("ⓘ").SetColor(theme.PrimaryColor)
	a.iconText.SetFontSize(16)

	Content := native.NewView().SetFlexDirection(yoga.FlexDirectionColumn).SetGap(yoga.GutterAll, 4)
	a.titleEl = native.NewText(title).SetFontSize(14).SetColor(theme.TextColor)
	a.descEl = native.NewText("").SetFontSize(14).SetColor(theme.MutedForegroundColor)
	Content.Add(a.titleEl, a.descEl)

	a.Add(a.iconText, Content)
	return a
}

// ElementType returns type identifier.
func (a *Alert) ElementType() string { return "Alert" }

// Draw renders the alert background and border.
func (a *Alert) Draw(screen *ebiten.Image) {
	bounds := a.GetBounds()
	br := core.BorderRadius{TopLeft: core.GetTheme().BorderRadius, TopRight: core.GetTheme().BorderRadius, BottomRight: core.GetTheme().BorderRadius, BottomLeft: core.GetTheme().BorderRadius}
	if a.bgColor != nil {
		native.DrawRoundedRectFill(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, br, a.bgColor)
	}
	if a.borderColor != nil {
		native.DrawRoundedRectStroke(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, br, 1, a.borderColor)
	}
	// Left accent bar for destructive variant
	if a.variant == AlertDestructive {
		vector.FillRect(screen, bounds.X, bounds.Y+4, 3, bounds.Height-8, core.GetTheme().DestructiveColor, false)
	}
}

// SetDescription sets the alert description.
func (a *Alert) SetDescription(desc string) *Alert {
	a.descEl.SetContent(desc)
	a.descEl.SetVisible(desc != "")
	a.Mark(core.FlagNeedLayout | core.FlagNeedDraw)
	return a
}

// SetVariant sets the alert variant.
func (a *Alert) SetVariant(v AlertVariant) *Alert {
	a.variant = v
	if v == AlertDestructive {
		a.iconText.SetColor(core.GetTheme().DestructiveColor)
		a.titleEl.SetColor(core.GetTheme().DestructiveColor)
	}
	a.Mark(core.FlagNeedDraw)
	return a
}
