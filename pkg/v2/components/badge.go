package components

import (
	"fmt"
	"image/color"

	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// Badge is a small status indicator, typically attached to another element.
type Badge struct {
	View
	textEl   *Text
	dotMode  bool
	maxCount int
}

// NewBadge creates a badge with numeric or text content.
// Pass "" for a dot-only badge.
func NewBadge(content string) *Badge {
	theme := core.GetTheme()
	b := &Badge{
		maxCount: 99,
	}
	b.Init(b)
	b.SetFlexDirection(yoga.FlexDirectionRow)
	b.SetJustifyContent(yoga.JustifyCenter)
	b.SetAlignItems(yoga.AlignCenter)
	b.SetBackgroundColor(theme.PrimaryColor)

	if content == "" {
		b.dotMode = true
		b.SetWidth(8)
		b.SetHeight(8)
		b.SetBorderRadius(4)
	} else {
		b.SetPadding(yoga.EdgeHorizontal, 6)
		b.SetPadding(yoga.EdgeVertical, 2)
		b.SetBorderRadius(999) // pill shape
		b.textEl = NewText(content).SetColor(color.White)
		b.textEl.SetFontSize(10)
		b.AppendChild(b.textEl)
	}
	return b
}

// ElementType returns type identifier.
func (b *Badge) ElementType() string { return "Badge" }

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
