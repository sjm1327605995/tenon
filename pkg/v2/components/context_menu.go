package components

import (
	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// ContextMenuItem represents a single context menu entry.
type ContextMenuItem struct {
	Label  string
	Action func()
}

// ContextMenu is a right-click menu.
type ContextMenu struct {
	core.BaseElement
	items  []ContextMenuItem
	panel  *View
	isOpen bool
}

// NewContextMenu creates a context menu.
func NewContextMenu() *ContextMenu {
	c := &ContextMenu{}
	c.Init(c)
	c.SetPositionType(yoga.PositionTypeAbsolute)
	c.SetVisible(false)

	c.panel = NewView()
	c.panel.SetFlexDirection(yoga.FlexDirectionColumn)
	c.panel.SetBackgroundColor(core.GetTheme().CardColor)
	c.panel.SetBorderColor(core.GetTheme().BorderColor)
	c.panel.SetBorderRadius(core.GetTheme().BorderRadius)
	c.panel.SetPadding(yoga.EdgeVertical, 8)
	c.panel.SetShadow(core.GetTheme().ShadowColor, 8, 0, 4)
	c.Add(c.panel)
	return c
}

// ElementType returns type identifier.
func (c *ContextMenu) ElementType() string { return "ContextMenu" }

// SetItems sets the menu items.
func (c *ContextMenu) SetItems(items []ContextMenuItem) *ContextMenu {
	c.items = items
	c.panel.ClearChildren()
	for _, item := range items {
		it := item
		btn := NewButton(it.Label).SetVariant(ButtonGhost)
		btn.SetJustifyContent(yoga.JustifyFlexStart)
		btn.SetWidthPercent(100)
		btn.SetBorderRadius(0)
		btn.SetOnClick(func() {
			c.Close()
			if it.Action != nil {
				it.Action()
			}
		})
		c.panel.Add(btn)
	}
	return c
}

// OpenAt shows the context menu at the given position.
func (c *ContextMenu) OpenAt(x, y float32) *ContextMenu {
	c.isOpen = true
	c.SetVisible(true)
	c.SetPosition(yoga.EdgeLeft, x)
	c.SetPosition(yoga.EdgeTop, y)
	c.Mark(core.FlagNeedLayout)
	if eng := c.GetEngine(); eng != nil {
		eng.AddOverlay(c)
	}
	return c
}

// Close hides the menu.
func (c *ContextMenu) Close() *ContextMenu {
	c.isOpen = false
	c.SetVisible(false)
	if eng := c.GetEngine(); eng != nil {
		eng.RemoveOverlay(c)
	}
	return c
}

// HandleEvent closes on outside click.
func (c *ContextMenu) HandleEvent(e *core.Event) bool {
	if e.Type == core.EventClick && c.isOpen {
		target := e.Target
		for target != nil {
			if target == c {
				return false
			}
			target = target.GetParent()
		}
		c.Close()
		return false
	}
	return false
}
