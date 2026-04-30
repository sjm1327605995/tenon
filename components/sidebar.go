package components

import (
	"github.com/sjm1327605995/tenon/internal/core"
	"github.com/sjm1327605995/tenon/internal/native"
	"github.com/sjm1327605995/tenon/yoga"
)

// SidebarNavItem represents a navigation item in the sidebar.
type SidebarNavItem struct {
	Label    string
	Icon     string
	Page     string
	Children []SidebarNavItem
}

// Sidebar is a collapsible navigation sidebar.
type Sidebar struct {
	native.View
	items    []SidebarNavItem
	onSelect func(page string)
	selected string
	Content  *native.View
}

// NewSidebar creates a sidebar.
func NewSidebar(items []SidebarNavItem) *Sidebar {
	theme := core.GetTheme()
	s := &Sidebar{
		items: items,
	}
	s.Init(s)
	s.SetWidth(260)
	s.SetHeightPercent(100)
	s.SetFlexDirection(yoga.FlexDirectionColumn)
	s.SetPadding(yoga.EdgeAll, 16)
	s.SetGap(yoga.GutterAll, 8)
	s.SetBackgroundColor(theme.BackgroundColor)
	s.SetBorder(yoga.EdgeRight, 1)
	s.SetBorderColor(theme.BorderColor)

	s.Content = native.NewView()
	s.Content.SetFlexDirection(yoga.FlexDirectionColumn)
	s.Content.SetGap(yoga.GutterAll, 4)
	s.buildItems()
	s.Add(s.Content)
	return s
}

// ElementType returns type identifier.
func (s *Sidebar) ElementType() string { return "Sidebar" }

func (s *Sidebar) buildItems() {
	s.Content.ClearChildren()
	for _, group := range s.items {
		if group.Label != "" && group.Page == "" {
			label := native.NewText(group.Label).SetFontSize(12).SetColor(core.GetTheme().MutedForegroundColor)
			label.SetPadding(yoga.EdgeVertical, 8)
			s.Content.Add(label)
		}

		for _, item := range group.Children {
			it := item
			btn := NewButton(it.Label).SetVariant(ButtonGhost)
			btn.SetJustifyContent(yoga.JustifyFlexStart)
			btn.SetWidthPercent(100)
			if s.selected == it.Page {
				btn.SetVariant(ButtonSecondary)
			}
			btn.SetOnClick(func() {
				s.selected = it.Page
				if s.onSelect != nil {
					s.onSelect(it.Page)
				}
				s.buildItems()
			})
			s.Content.Add(btn)
		}
	}
	s.Mark(core.FlagNeedLayout)
}

// SetOnSelect sets the navigation selection callback.
func (s *Sidebar) SetOnSelect(fn func(page string)) *Sidebar {
	s.onSelect = fn
	return s
}

// SetSelected sets the selected page.
func (s *Sidebar) SetSelected(page string) *Sidebar {
	s.selected = page
	s.buildItems()
	return s
}
