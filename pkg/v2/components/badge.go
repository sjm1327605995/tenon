package components

import (
	"fmt"
	"image/color"

	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// BadgeVariant defines the visual style of a badge.
type BadgeVariant int

const (
	BadgeDefault BadgeVariant = iota
	BadgeSecondary
	BadgeOutline
	BadgeDestructive
)

// Badge is a small status indicator, typically attached to another element.
type Badge struct {
	View
	textEl   *Text
	dotMode  bool
	maxCount int
	variant  BadgeVariant
}

func NewBadge(content string) *Badge {
	theme := core.GetTheme()
	b := &Badge{
		maxCount: 99,
		variant:  BadgeDefault,
	}
	b.Init(b)
	b.SetFlexDirection(yoga.FlexDirectionRow)
	b.SetJustifyContent(yoga.JustifyCenter)
	b.SetAlignItems(yoga.AlignCenter)
	b.applyVariantStyles(theme)

	if content == "" {
		b.dotMode = true
		b.SetWidth(8)
		b.SetHeight(8)
		b.SetBorderRadius(4)
	} else {
		b.SetPadding(yoga.EdgeHorizontal, 10)
		b.SetPadding(yoga.EdgeVertical, 4)
		b.textEl = NewText(content).SetColor(theme.PrimaryForegroundColor)
		b.textEl.SetFontSize(12)
		b.AppendChild(b.textEl)
	}
	return b
}

// ElementType returns type identifier.
func (b *Badge) ElementType() string { return "Badge" }

// Update auto-computes pill-shaped borderRadius.
func (b *Badge) Update() error {
	if !b.dotMode {
		bounds := b.GetBounds()
		if bounds.Height > 0 {
			halfH := bounds.Height / 2
			if b.borderRadius.TopLeft != halfH {
				b.SetBorderRadius(halfH)
			}
		}
	}
	return nil
}

// SetVariant changes the badge style.
func (b *Badge) SetVariant(v BadgeVariant) *Badge {
	b.variant = v
	b.applyVariantStyles(core.GetTheme())
	return b
}

func (b *Badge) applyVariantStyles(theme *core.Theme) {
	switch b.variant {
	case BadgeDefault:
		b.SetBackgroundColor(theme.PrimaryColor)
		if b.textEl != nil {
			b.textEl.SetColor(theme.PrimaryForegroundColor)
		}
	case BadgeSecondary:
		b.SetBackgroundColor(theme.SecondaryColor)
		if b.textEl != nil {
			b.textEl.SetColor(theme.SecondaryForegroundColor)
		}
	case BadgeOutline:
		b.SetBackgroundColor(nil)
		b.SetBorderColor(theme.BorderColor)
		if b.textEl != nil {
			b.textEl.SetColor(theme.TextColor)
		}
	case BadgeDestructive:
		b.SetBackgroundColor(theme.DestructiveColor)
		if b.textEl != nil {
			b.textEl.SetColor(theme.DestructiveForegroundColor)
		}
	}
}

// SetCount sets numeric content with optional "99+" overflow.
func (b *Badge) SetCount(count int) *Badge {
	if b.dotMode {
		return b
	}
	var text string
	if count > b.maxCount {
		text = fmt.Sprintf("%d+", b.maxCount)
	} else {
		text = fmt.Sprintf("%d", count)
	}
	if b.textEl != nil {
		b.textEl.SetContent(text)
	}
	return b
}

// SetMaxCount sets the overflow threshold (default 99).
func (b *Badge) SetMaxCount(max int) *Badge {
	b.maxCount = max
	return b
}

// SetDotMode switches to a small dot without text.
func (b *Badge) SetDotMode(dot bool) *Badge {
	b.dotMode = dot
	if dot {
		b.SetWidth(8)
		b.SetHeight(8)
		b.SetBorderRadius(4)
		if b.textEl != nil {
			b.textEl.SetVisible(false)
		}
	} else {
		if b.textEl != nil {
			b.textEl.SetVisible(true)
		}
	}
	b.Mark(core.FlagNeedDraw)
	return b
}

// SetTextColor sets the badge text color.
func (b *Badge) SetTextColor(clr color.Color) *Badge {
	if b.textEl != nil {
		b.textEl.SetColor(clr)
	}
	return b
}

// SyncFrom 同步新 Badge 的属性到当前 Element（声明式重建）。
func (b *Badge) SyncFrom(src core.Element) {
	other, ok := src.(*Badge)
	if !ok {
		return
	}
	b.View.SyncFrom(&other.View)
	sb := &SyncBuilder{}
	syncField(sb, &b.dotMode, other.dotMode)
	syncField(sb, &b.maxCount, other.maxCount)
	if b.variant != other.variant {
		b.variant = other.variant
		b.applyVariantStyles(core.GetTheme())
		sb.NeedDraw = true
	}
	sb.MarkDraw(b)
}
