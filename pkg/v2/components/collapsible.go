package components

import (
	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// Collapsible is a container that can be expanded or collapsed.
type Collapsible struct {
	core.BaseElement
	trigger     *Button
	contentWrap *View
	open        bool
	onChange    func(open bool)
}

// NewCollapsible creates a collapsible container.
func NewCollapsible(triggerText string, content core.Element) *Collapsible {
	c := &Collapsible{open: false}
	c.Init(c)
	c.SetFlexDirection(yoga.FlexDirectionColumn)
	c.SetWidthPercent(100)

	c.trigger = NewButton(triggerText).SetVariant(ButtonGhost)
	c.trigger.SetOnClick(func() { c.toggle() })
	c.trigger.SetWidthPercent(100)
	c.trigger.SetJustifyContent(yoga.JustifyFlexStart)

	c.contentWrap = NewView()
	c.contentWrap.SetFlexDirection(yoga.FlexDirectionColumn)
	c.contentWrap.SetWidthPercent(100)
	if content != nil {
		c.contentWrap.Add(content)
	}
	c.contentWrap.SetVisible(false)

	c.Add(c.trigger, c.contentWrap)
	return c
}

// ElementType returns type identifier.
func (c *Collapsible) ElementType() string { return "Collapsible" }

func (c *Collapsible) toggle() {
	c.open = !c.open
	c.contentWrap.SetVisible(c.open)
	c.Mark(core.FlagNeedLayout)
	if c.onChange != nil {
		c.onChange(c.open)
	}
}

// SetOpen sets the open state.
func (c *Collapsible) SetOpen(open bool) *Collapsible {
	c.open = open
	c.contentWrap.SetVisible(open)
	c.Mark(core.FlagNeedLayout)
	return c
}

// IsOpen returns the current open state.
func (c *Collapsible) IsOpen() bool { return c.open }

// SetOnChange sets the change callback.
func (c *Collapsible) SetOnChange(fn func(open bool)) *Collapsible {
	c.onChange = fn
	return c
}

// SetTrigger updates the trigger button text.
func (c *Collapsible) SetTrigger(text string) *Collapsible {
	c.trigger.SetText(text)
	return c
}
