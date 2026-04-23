package components

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// MenuItemData 定义一个菜单项的配置。
type MenuItemData struct {
	Key   string
	Label string
}

// MenuItem 是单个菜单项宿主组件。
// 由 View + hoverBg(View) + indicator(View) + label(Text) 组合而成，零自定义绘制代码。
type MenuItem struct {
	View
	selected  bool
	hovered   bool
	onClick   func()
	hoverBg   *View
	indicator *View
	label     *Text
}

// NewMenuItem 创建一个菜单项。
func NewMenuItem(label string) *MenuItem {
	mi := &MenuItem{}
	mi.Init(mi)
	mi.GetElement().Yoga.StyleSetHeight(40)
	mi.GetElement().Yoga.StyleSetWidthPercent(100)
	mi.GetElement().Yoga.StyleSetAlignItems(yoga.AlignCenter)
	mi.GetElement().Yoga.StyleSetFlexDirection(yoga.FlexDirectionRow)

	theme := core.GetTheme()

	// Hover 背景：absolute 定位，左右上下各留 4px 边距，圆角
	mi.hoverBg = NewView().
		SetPositionType(yoga.PositionTypeAbsolute).
		SetPosition(yoga.EdgeLeft, 4).
		SetPosition(yoga.EdgeTop, 4).
		SetPosition(yoga.EdgeRight, 4).
		SetPosition(yoga.EdgeBottom, 4).
		SetBorderRadius(4).
		SetPointerEvents(core.PointerEventsNone).
		SetVisible(false)
	mi.AddChild(mi.hoverBg)

	// 选中指示器：左侧 3px 竖条
	mi.indicator = NewView().
		SetPositionType(yoga.PositionTypeAbsolute).
		SetPosition(yoga.EdgeLeft, 0).
		SetPosition(yoga.EdgeTop, 0).
		SetPosition(yoga.EdgeBottom, 0).
		SetWidth(3).
		SetBackgroundColor(theme.PrimaryColor).
		SetPointerEvents(core.PointerEventsNone).
		SetVisible(false)
	mi.AddChild(mi.indicator)

	// 文本标签
	mi.label = NewText(label)
	mi.label.SetMargin(yoga.EdgeLeft, 24)
	mi.label.GetElement().PointerEvents = core.PointerEventsNone
	mi.AddChild(mi.label)

	mi.refreshStyle()
	return mi
}

func (mi *MenuItem) refreshStyle() {
	theme := core.GetTheme()
	if mi.selected {
		mi.SetBackgroundColor(theme.MenuItemSelectedBg)
		mi.label.SetColor(theme.MenuItemSelectedText)
		mi.indicator.SetVisible(true)
		mi.hoverBg.SetVisible(false)
	} else {
		mi.SetBackgroundColor(nil)
		mi.label.SetColor(theme.TextColor)
		mi.indicator.SetVisible(false)
	}
}

// SetSelected 设置选中状态。
func (mi *MenuItem) SetSelected(selected bool) *MenuItem {
	mi.selected = selected
	mi.hovered = false
	mi.refreshStyle()
	return mi
}

// SetOnClick 设置点击回调。
func (mi *MenuItem) SetOnClick(fn func()) *MenuItem {
	mi.onClick = fn
	return mi
}

// Update 每帧检测鼠标是否悬停。
func (mi *MenuItem) Update() error {
	mx, my := ebiten.CursorPosition()
	bounds := mi.GetLayoutBounds()
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

// HandleEvent 处理点击事件。
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

// SyncFrom 同步菜单项属性（颜色不同步，以便主题切换后使用新主题色）。
func (mi *MenuItem) SyncFrom(other core.Host) {
	if o, ok := other.(*MenuItem); ok {
		mi.selected = o.selected
		mi.onClick = o.onClick
		if mi.label != nil && o.label != nil {
			mi.label.Content = o.label.Content
			mi.label.cachedLayout = nil
		}
		mi.refreshStyle()
	}
}

// ==================== Menu 容器 ====================

// Menu 是垂直导航菜单容器。
type Menu struct {
	core.BaseHost
	items       []MenuItemData
	selectedKey string
	onSelect    func(key string)
	menuWidth   float32
}

// NewMenu 创建一个菜单容器。
func NewMenu() *Menu {
	theme := core.GetTheme()
	m := &Menu{
		menuWidth: 200,
	}
	m.Init(m)
	m.GetElement().Yoga.StyleSetWidth(m.menuWidth)
	m.GetElement().Yoga.StyleSetFlexDirection(yoga.FlexDirectionColumn)
	m.GetElement().BackgroundColor = theme.MenuBg
	return m
}

// SetItems 设置菜单项数据并重建子组件。
func (m *Menu) SetItems(items []MenuItemData) *Menu {
	m.items = items
	m.rebuild()
	return m
}

// SetSelectedKey 设置当前选中的菜单项。
func (m *Menu) SetSelectedKey(key string) *Menu {
	m.selectedKey = key
	m.refreshSelection()
	return m
}

// SetOnSelect 设置选中回调。
func (m *Menu) SetOnSelect(fn func(key string)) *Menu {
	m.onSelect = fn
	return m
}

func (m *Menu) rebuild() {
	m.ClearChildren()
	for _, item := range m.items {
		if item.Key == "" {
			// 分组标题
			title := NewText(item.Label)
			title.SetFontSize(core.GetTheme().FontSizeSM)
			title.SetColor(core.GetTheme().TextMutedColor)
			title.GetElement().Yoga.StyleSetMargin(yoga.EdgeTop, 12)
			title.GetElement().Yoga.StyleSetMargin(yoga.EdgeBottom, 4)
			title.GetElement().Yoga.StyleSetMargin(yoga.EdgeLeft, 24)
			m.AddChild(title)
			continue
		}
		mi := NewMenuItem(item.Label)
		mi.SetSelected(item.Key == m.selectedKey)
		key := item.Key // 闭包捕获
		mi.SetOnClick(func() {
			if m.onSelect != nil {
				m.onSelect(key)
			}
		})
		m.AddChild(mi)
	}
}

func (m *Menu) refreshSelection() {
	for _, child := range m.GetChildren() {
		if mi, ok := child.(*MenuItem); ok {
			// 通过 SyncFrom 无法直接找到对应项，这里简单遍历
			for _, item := range m.items {
				if item.Key == "" {
					continue
				}
				if mi.label != nil && mi.label.Content == item.Label {
					mi.SetSelected(item.Key == m.selectedKey)
					break
				}
			}
		}
	}
}

// SyncFrom 同步菜单属性。
func (m *Menu) SyncFrom(other core.Host) {
	if o, ok := other.(*Menu); ok {
		m.items = o.items
		m.selectedKey = o.selectedKey
		m.onSelect = o.onSelect
		m.menuWidth = o.menuWidth
	}
}
