package components

import (
	"github.com/sjm1327605995/tenon/internal/core"
	"github.com/sjm1327605995/tenon/internal/native"
	"github.com/sjm1327605995/tenon/yoga"
)

// Card is a Content container with header, Content and footer sections.
type Card struct {
	native.View
	headerEl  *native.View
	titleEl   *native.Text
	descEl    *native.Text
	contentEl *native.View
	footerEl  *native.View
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

	c.headerEl = native.NewView()
	c.headerEl.SetFlexDirection(yoga.FlexDirectionColumn)
	c.headerEl.SetGap(yoga.GutterAll, 4)
	c.headerEl.SetPadding(yoga.EdgeAll, 24)
	c.headerEl.SetPadding(yoga.EdgeBottom, 0)

	c.titleEl = native.NewText("").SetFontSize(18).SetColor(theme.CardForegroundColor)
	c.descEl = native.NewText("").SetFontSize(14).SetColor(theme.MutedForegroundColor)
	c.headerEl.Add(c.titleEl, c.descEl)

	c.contentEl = native.NewView()
	c.contentEl.SetFlexDirection(yoga.FlexDirectionColumn)
	c.contentEl.SetGap(yoga.GutterAll, 8)
	c.contentEl.SetPadding(yoga.EdgeAll, 24)

	c.footerEl = native.NewView()
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
	c.headerEl.SetVisible(title != "" || c.descEl.Content != "")
	c.Mark(core.FlagNeedLayout)
	return c
}

// SetDescription sets the card description.
func (c *Card) SetDescription(desc string) *Card {
	c.descEl.SetContent(desc)
	c.descEl.SetVisible(desc != "")
	c.headerEl.SetVisible(c.titleEl.Content != "" || desc != "")
	c.Mark(core.FlagNeedLayout)
	return c
}

// AddContent adds elements to the card Content area.
func (c *Card) AddContent(children ...core.Element) *Card {
	for _, child := range children {
		c.contentEl.AppendChild(child)
	}
	c.contentEl.SetVisible(len(c.contentEl.GetChildren()) > 0)
	c.Mark(core.FlagNeedLayout)
	return c
}

// SetFooter adds elements to the card footer.
func (c *Card) SetFooter(children ...core.Element) *Card {
	c.footerEl.ClearChildren()
	for _, child := range children {
		c.footerEl.AppendChild(child)
	}
	c.footerEl.SetVisible(len(children) > 0)
	c.Mark(core.FlagNeedLayout)
	return c
}
