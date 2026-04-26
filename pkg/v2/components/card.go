package components

import (
	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// Card is a content container with header, content and footer sections.
type Card struct {
	View
	headerEl  *View
	titleEl   *Text
	descEl    *Text
	contentEl *View
	footerEl  *View
}

// NewCard creates a card component.
func NewCard() *Card {
	theme := core.GetTheme()
	c := &Card{}
	c.Init(c)
	c.SetFlexDirection(yoga.FlexDirectionColumn)
	c.SetBackgroundColor(theme.CardColor)
	c.SetBorderRadius(theme.BorderRadius)
	c.SetBorderColor(theme.BorderColor)

	c.headerEl = NewView()
	c.headerEl.SetFlexDirection(yoga.FlexDirectionColumn)
	c.headerEl.SetGap(yoga.GutterAll, 4)
	c.headerEl.SetPadding(yoga.EdgeAll, 24)
	c.headerEl.SetPadding(yoga.EdgeBottom, 0)

	c.titleEl = NewText("").SetFontSize(18).SetColor(theme.CardForegroundColor)
	c.descEl = NewText("").SetFontSize(14).SetColor(theme.MutedForegroundColor)
	c.headerEl.Add(c.titleEl, c.descEl)

	c.contentEl = NewView()
	c.contentEl.SetFlexDirection(yoga.FlexDirectionColumn)
	c.contentEl.SetGap(yoga.GutterAll, 8)
	c.contentEl.SetPadding(yoga.EdgeAll, 24)

	c.footerEl = NewView()
	c.footerEl.SetFlexDirection(yoga.FlexDirectionRow)
	c.footerEl.SetGap(yoga.GutterAll, 8)
	c.footerEl.SetJustifyContent(yoga.JustifyFlexEnd)
	c.footerEl.SetPadding(yoga.EdgeAll, 24)
	c.footerEl.SetPadding(yoga.EdgeTop, 0)

	c.headerEl.SetVisible(false)
	c.contentEl.SetVisible(false)
	c.footerEl.SetVisible(false)

	c.Add(c.headerEl, c.contentEl, c.footerEl)
	return c
}

// ElementType returns type identifier.
func (c *Card) ElementType() string { return "Card" }

// SetTitle sets the card title.
func (c *Card) SetTitle(title string) *Card {
	c.titleEl.SetContent(title)
	c.titleEl.SetVisible(title != "")
	c.headerEl.SetVisible(title != "" || c.descEl.content != "")
	c.Mark(core.FlagNeedLayout | core.FlagNeedDraw)
	return c
}

// SetDescription sets the card description.
func (c *Card) SetDescription(desc string) *Card {
	c.descEl.SetContent(desc)
	c.descEl.SetVisible(desc != "")
	c.headerEl.SetVisible(c.titleEl.content != "" || desc != "")
	c.Mark(core.FlagNeedLayout | core.FlagNeedDraw)
	return c
}

// AddContent adds elements to the card content area.
func (c *Card) AddContent(children ...core.Element) *Card {
	for _, child := range children {
		c.contentEl.AppendChild(child)
	}
	c.contentEl.SetVisible(len(c.contentEl.GetChildren()) > 0)
	c.Mark(core.FlagNeedLayout | core.FlagNeedDraw)
	return c
}

// SetFooter adds elements to the card footer.
func (c *Card) SetFooter(children ...core.Element) *Card {
	c.footerEl.ClearChildren()
	for _, child := range children {
		c.footerEl.AppendChild(child)
	}
	c.footerEl.SetVisible(len(children) > 0)
	c.Mark(core.FlagNeedLayout | core.FlagNeedDraw)
	return c
}
