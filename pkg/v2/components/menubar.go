package components

import (
	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// MenubarItem represents a top-level menu item.
type MenubarItem struct {
	Label   string
	SubMenu []string
}

// Menubar is a horizontal menu bar.
type Menubar struct {
	View
	items []MenubarItem
}

// NewMenubar creates a menubar.
func NewMenubar(items []MenubarItem) *Menubar {
	m := &Menubar{items: items}
	m.Init(m)
	m.SetFlexDirection(yoga.FlexDirectionRow)
	m.SetGap(yoga.GutterAll, 4)
	m.SetPadding(yoga.EdgeHorizontal, 16)
	m.SetPadding(yoga.EdgeVertical, 8)
	m.SetBackgroundColor(core.GetTheme().BackgroundColor)
	m.SetBorder(yoga.EdgeBottom, 1)
	m.SetBorderColor(core.GetTheme().BorderColor)

	for _, item := range items {
		btn := NewButton(item.Label).SetVariant(ButtonGhost)
		m.Add(btn)
	}
	return m
}

// ElementType returns type identifier.
func (m *Menubar) ElementType() string { return "Menubar" }
