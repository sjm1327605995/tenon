package components

import (
	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// NavigationMenuItem represents a navigation item.
type NavigationMenuItem struct {
	Label    string
	Href     string
	Children []NavigationMenuItem
}

// NavigationMenu is a complex navigation component.
type NavigationMenu struct {
	View
	items []NavigationMenuItem
}

// NewNavigationMenu creates a navigation menu.
func NewNavigationMenu(items []NavigationMenuItem) *NavigationMenu {
	n := &NavigationMenu{items: items}
	n.Init(n)
	n.SetFlexDirection(yoga.FlexDirectionRow)
	n.SetGap(yoga.GutterAll, 4)
	n.SetPadding(yoga.EdgeHorizontal, 16)
	n.SetPadding(yoga.EdgeVertical, 12)
	n.SetBackgroundColor(core.GetTheme().BackgroundColor)

	for _, item := range items {
		btn := NewButton(item.Label).SetVariant(ButtonGhost)
		n.Add(btn)
	}
	return n
}

// ElementType returns type identifier.
func (n *NavigationMenu) ElementType() string { return "NavigationMenu" }
