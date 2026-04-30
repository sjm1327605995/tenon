package components

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/internal/core"
	"github.com/sjm1327605995/tenon/internal/native"
	"github.com/sjm1327605995/tenon/yoga"
)

// MenuItemData defines a menu item configuration.
type MenuItemData struct {
	Key   string
	Label string
}

// MenuItem is a single menu item.
type MenuItem struct {
	native.View
	itemKey   string
	selected  bool
	hovered   bool
	onClick   func()
	hoverBg   *native.View
	indicator *native.View
	label     *native.Text
}

// NewMenuItem creates a menu item.
func NewMenuItem(label string) *MenuItem {
	mi := &MenuItem{}
	mi.Init(mi)
	mi.SetHeight(40)
	mi.SetWidthPercent(100)
	mi.SetAlignItems(yoga.AlignCenter)
	mi.SetFlexDirection(yoga.FlexDirectionRow)

	theme := core.GetTheme()

	mi.hoverBg = native.NewView()
	mi.hoverBg.SetPositionType(yoga.PositionTypeAbsolute)
	mi.hoverBg.SetPosition(yoga.EdgeLeft, 4)
	mi.hoverBg.SetPosition(yoga.EdgeTop, 4)
	mi.hoverBg.SetPosition(yoga.EdgeRight, 4)
	mi.hoverBg.SetPosition(yoga.EdgeBottom, 4)
	mi.hoverBg.SetBorderRadius(4)
	mi.hoverBg.SetVisible(false)
	mi.AppendChild(mi.hoverBg)

	mi.indicator = native.NewView()
	mi.indicator.SetPositionType(yoga.PositionTypeAbsolute)
	mi.indicator.SetPosition(yoga.EdgeLeft, 0)
	mi.indicator.SetPosition(yoga.EdgeTop, 0)
	mi.indicator.SetPosition(yoga.EdgeBottom, 0)
	mi.indicator.SetWidth(3)
	mi.indicator.SetBackgroundColor(theme.PrimaryColor)
	mi.indicator.SetVisible(false)
	mi.AppendChild(mi.indicator)

	mi.label = native.NewText(label)
	mi.label.SetMargin(yoga.EdgeLeft, 24)
	mi.AppendChild(mi.label)

	mi.refreshStyle()
	return mi
}

func (mi *MenuItem) refreshStyle() {
	theme := core.GetTheme()
	if mi.selected {
		mi.View.SetBackgroundColor(theme.MenuItemSelectedBg)
		mi.label.SetColor(theme.MenuItemSelectedText)
		mi.indicator.SetVisible(true)
		mi.hoverBg.SetVisible(false)
	} else {
		mi.View.SetBackgroundColor(nil)
		mi.label.SetColor(theme.TextColor)
		mi.indicator.SetVisible(false)
	}
}

// ElementType returns type identifier.
func (mi *MenuItem) ElementType() string { return "MenuItem" }

// SetItemKey sets the internal key for selection matching.
func (mi *MenuItem) SetItemKey(key string) *MenuItem {
	mi.itemKey = key
	return mi
}

// SetSelected sets the selected state.
func (mi *MenuItem) SetSelected(selected bool) *MenuItem {
	mi.selected = selected
	mi.hovered = false
	mi.refreshStyle()
	return mi
}

// SetOnClick sets the click callback.
func (mi *MenuItem) SetOnClick(fn func()) *MenuItem {
	mi.onClick = fn
	return mi
}

// Update detects hover state per frame.
func (mi *MenuItem) Update() error {
	mx, my := ebiten.CursorPosition()
	bounds := mi.GetBounds()
	hovered := float32(mx) >= bounds.X && float32(mx) < bounds.X+bounds.Width &&
		float32(my) >= bounds.Y && float32(my) < bounds.Y+bounds.Height

	if hovered != mi.hovered {
		mi.hovered = hovered
		if hovered && !mi.selected {
			mi.hoverBg.SetBackgroundColor(core.GetTheme().MenuItemHoverBg).SetVisible(true)
		} else {
			mi.hoverBg.SetVisible(false)
		}
	}
	return nil
}

// HandleEvent processes click events.
func (mi *MenuItem) HandleEvent(e *core.Event) bool {
	switch e.Type {
	case core.EventMouseDown:
		return true
	case core.EventMouseUp:
		return true
	case core.EventClick:
		if mi.onClick != nil {
			mi.onClick()
		}
		return true
	}
	return false
}

// Menu is a vertical navigation menu container.
type Menu struct {
	native.View
	items       []MenuItemData
	selectedKey string
	onSelect    func(key string)
	menuWidth   float32
}

// NewMenu creates a menu container.
func NewMenu() *Menu {
	theme := core.GetTheme()
	m := &Menu{
		menuWidth: 200,
	}
	m.Init(m)
	m.SetWidth(m.menuWidth)
	m.SetFlexDirection(yoga.FlexDirectionColumn)
	m.SetBackgroundColor(theme.MenuBg)
	return m
}

// ElementType returns type identifier.
func (m *Menu) ElementType() string { return "Menu" }

// SetItems sets menu item data and rebuilds children.
func (m *Menu) SetItems(items []MenuItemData) *Menu {
	m.items = items
	m.rebuild()
	return m
}

// SetSelectedKey sets the currently selected item key.
func (m *Menu) SetSelectedKey(key string) *Menu {
	m.selectedKey = key
	m.refreshSelection()
	return m
}

// SetOnSelect sets the selection callback.
func (m *Menu) SetOnSelect(fn func(key string)) *Menu {
	m.onSelect = fn
	return m
}

func (m *Menu) rebuild() {
	m.ClearChildren()
	for _, item := range m.items {
		if item.Key == "" {
			title := native.NewText(item.Label)
			title.SetFontSize(float64(core.GetTheme().FontSizeSM))
			title.SetColor(core.GetTheme().TextMutedColor)
			title.SetMargin(yoga.EdgeTop, 12)
			title.SetMargin(yoga.EdgeBottom, 4)
			title.SetMargin(yoga.EdgeLeft, 24)
			m.AppendChild(title)
			continue
		}
		mi := NewMenuItem(item.Label)
		mi.SetItemKey(item.Key)
		mi.SetSelected(item.Key == m.selectedKey)
		key := item.Key
		mi.SetOnClick(func() {
			if m.onSelect != nil {
				m.onSelect(key)
			}
		})
		m.AppendChild(mi)
	}
}

func (m *Menu) refreshSelection() {
	for _, child := range m.GetChildren() {
		if mi, ok := child.(*MenuItem); ok {
			mi.SetSelected(mi.itemKey == m.selectedKey)
		}
	}
}
