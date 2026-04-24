package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// TabItem 定义单个标签页的数据。
type TabItem struct {
	Key      string
	Label    string
	Content  core.Component
	Disabled bool
}

// Tabs 是原生标签页宿主组件，包含标签栏和内容区。
type Tabs struct {
	core.BaseHost
	items       []TabItem
	activeKey   string
	tabType     string // "line" | "card"
	onChange    func(key string)
	tabBar      *View
	contentArea *View
	tabLabels   map[string]*Text
	tabWrappers map[string]*View
}

// NewTabs 创建一个标签页组件。
func NewTabs() *Tabs {
	t := &Tabs{
		tabType:     "line",
		tabLabels:   make(map[string]*Text),
		tabWrappers: make(map[string]*View),
	}
	t.Init(t)
	t.SetFocusable(true)
	t.GetElement().Yoga.StyleSetFlexDirection(yoga.FlexDirectionColumn)
	t.GetElement().Yoga.StyleSetWidthPercent(100)

	t.tabBar = NewView()
	t.tabBar.Init(t.tabBar)
	t.tabBar.GetElement().Yoga.StyleSetFlexDirection(yoga.FlexDirectionRow)
	t.tabBar.GetElement().Yoga.StyleSetHeight(44)
	t.tabBar.GetElement().Yoga.StyleSetAlignItems(yoga.AlignCenter)
	t.tabBar.GetElement().Yoga.StyleSetBorder(yoga.EdgeBottom, 1)
	t.tabBar.GetElement().BorderColor = core.GetTheme().BorderColor
	t.AddChild(t.tabBar)

	t.contentArea = NewView()
	t.contentArea.Init(t.contentArea)
	t.contentArea.GetElement().Yoga.StyleSetFlexGrow(1)
	t.contentArea.GetElement().Yoga.StyleSetWidthPercent(100)
	t.AddChild(t.contentArea)

	return t
}

// SetItems 设置标签页数据并重建 UI。
func (t *Tabs) SetItems(items []TabItem) *Tabs {
	t.items = items
	if len(items) > 0 && t.activeKey == "" {
		t.activeKey = items[0].Key
	}
	t.rebuild()
	return t
}

// SetActiveKey 设置当前激活的标签页。
func (t *Tabs) SetActiveKey(key string) *Tabs {
	if t.activeKey != key {
		t.activeKey = key
		t.refreshTabs()
		t.refreshContent()
	}
	return t
}

// SetType 设置标签栏样式："line" 或 "card"。
func (t *Tabs) SetType(tabType string) *Tabs {
	t.tabType = tabType
	t.refreshTabs()
	return t
}

// SetOnChange 设置切换回调。
func (t *Tabs) SetOnChange(fn func(key string)) *Tabs {
	t.onChange = fn
	return t
}

// HandleEvent 处理点击事件，根据点击位置切换标签页。
func (t *Tabs) HandleEvent(e *core.Event) bool {
	if e.Type != core.EventClick {
		return false
	}
	for _, item := range t.items {
		if item.Disabled {
			continue
		}
		wrapper, ok := t.tabWrappers[item.Key]
		if !ok {
			continue
		}
		bounds := wrapper.GetLayoutBounds()
		if e.X >= bounds.X && e.X < bounds.X+bounds.Width &&
			e.Y >= bounds.Y && e.Y < bounds.Y+bounds.Height {
			t.SetActiveKey(item.Key)
			if t.onChange != nil {
				t.onChange(item.Key)
			}
			return true
		}
	}
	return false
}

func (t *Tabs) rebuild() {
	t.tabBar.ClearChildren()
	t.tabLabels = make(map[string]*Text)
	t.tabWrappers = make(map[string]*View)

	theme := core.GetTheme()

	for _, item := range t.items {
		isActive := item.Key == t.activeKey

		wrapper := NewView()
		wrapper.Init(wrapper)
		wrapper.GetElement().Yoga.StyleSetFlexDirection(yoga.FlexDirectionColumn)
		wrapper.GetElement().Yoga.StyleSetAlignItems(yoga.AlignCenter)
		wrapper.GetElement().Yoga.StyleSetPadding(yoga.EdgeHorizontal, 16)
		wrapper.GetElement().Yoga.StyleSetPadding(yoga.EdgeVertical, 10)
		wrapper.GetElement().PointerEvents = core.PointerEventsAuto

		if t.tabType == "card" {
			if isActive {
				wrapper.GetElement().BackgroundColor = theme.SurfaceColor
				wrapper.GetElement().Yoga.StyleSetBorder(yoga.EdgeTop, 1)
				wrapper.GetElement().Yoga.StyleSetBorder(yoga.EdgeLeft, 1)
				wrapper.GetElement().Yoga.StyleSetBorder(yoga.EdgeRight, 1)
				wrapper.GetElement().BorderColor = theme.BorderColor
				wrapper.GetElement().Yoga.StyleSetMargin(yoga.EdgeBottom, -1)
			} else {
				wrapper.GetElement().BackgroundColor = theme.BackgroundColor
			}
		}

		label := NewText(item.Label)
		label.SetFontSize(theme.FontSizeBase)
		if isActive {
			label.SetColor(theme.PrimaryColor)
		} else if item.Disabled {
			label.SetColor(theme.TextMutedColor)
		} else {
			label.SetColor(theme.TextColor)
		}
		label.GetElement().PointerEvents = core.PointerEventsNone
		wrapper.AddChild(label)
		t.tabLabels[item.Key] = label

		if isActive && t.tabType == "line" {
			indicator := NewView()
			indicator.Init(indicator)
			indicator.GetElement().Yoga.StyleSetWidthPercent(100)
			indicator.GetElement().Yoga.StyleSetHeight(2)
			indicator.GetElement().Yoga.StyleSetMargin(yoga.EdgeTop, 4)
			indicator.GetElement().BackgroundColor = theme.PrimaryColor
			indicator.GetElement().PointerEvents = core.PointerEventsNone
			wrapper.AddChild(indicator)
		}

		t.tabWrappers[item.Key] = wrapper
		t.tabBar.AddChild(wrapper)
	}

	t.refreshContent()
}

func (t *Tabs) refreshTabs() {
	theme := core.GetTheme()
	for _, item := range t.items {
		isActive := item.Key == t.activeKey
		label, ok := t.tabLabels[item.Key]
		if ok {
			if isActive {
				label.SetColor(theme.PrimaryColor)
			} else if item.Disabled {
				label.SetColor(theme.TextMutedColor)
			} else {
				label.SetColor(theme.TextColor)
			}
		}

		wrapper, ok := t.tabWrappers[item.Key]
		if !ok {
			continue
		}
		if t.tabType == "card" {
			if isActive {
				wrapper.GetElement().BackgroundColor = theme.SurfaceColor
				wrapper.GetElement().Yoga.StyleSetBorder(yoga.EdgeTop, 1)
				wrapper.GetElement().Yoga.StyleSetBorder(yoga.EdgeLeft, 1)
				wrapper.GetElement().Yoga.StyleSetBorder(yoga.EdgeRight, 1)
				wrapper.GetElement().BorderColor = theme.BorderColor
				wrapper.GetElement().Yoga.StyleSetMargin(yoga.EdgeBottom, -1)
			} else {
				wrapper.GetElement().BackgroundColor = theme.BackgroundColor
				wrapper.GetElement().Yoga.StyleSetBorder(yoga.EdgeTop, 0)
				wrapper.GetElement().Yoga.StyleSetBorder(yoga.EdgeLeft, 0)
				wrapper.GetElement().Yoga.StyleSetBorder(yoga.EdgeRight, 0)
			}
		}
	}
}

func (t *Tabs) refreshContent() {
	t.contentArea.ClearChildren()
	for _, item := range t.items {
		if item.Key == t.activeKey && item.Content != nil {
			t.contentArea.AddChild(item.Content)
			break
		}
	}
	if t.GetEngine() != nil {
		t.GetEngine().InvalidateAll()
	}
}

// Draw 绘制标签页背景和边框。
func (t *Tabs) Draw(screen *ebiten.Image) {
	el := t.GetElement()
	if el == nil || !el.Visible {
		return
	}
	bounds := t.GetLayoutBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	if t.tabType == "card" && el.BackgroundColor != nil {
		vector.FillRect(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, el.BackgroundColor, false)
	}
}

// ==================== 链式 API ====================

func (t *Tabs) SetWidth(width float32) *Tabs {
	t.GetElement().Yoga.StyleSetWidth(width)
	return t
}
func (t *Tabs) SetWidthPercent(percent float32) *Tabs {
	t.GetElement().Yoga.StyleSetWidthPercent(percent)
	return t
}
func (t *Tabs) SetHeight(height float32) *Tabs {
	t.GetElement().Yoga.StyleSetHeight(height)
	return t
}
func (t *Tabs) SetMargin(edge yoga.Edge, value float32) *Tabs {
	t.GetElement().Yoga.StyleSetMargin(edge, value)
	return t
}
func (t *Tabs) SetPadding(edge yoga.Edge, value float32) *Tabs {
	t.GetElement().Yoga.StyleSetPadding(edge, value)
	return t
}
func (t *Tabs) SetBackgroundColor(clr color.Color) *Tabs {
	t.GetElement().BackgroundColor = clr
	return t
}

// SyncFrom 同步标签页属性。
func (t *Tabs) SyncFrom(other core.Host) {
	if o, ok := other.(*Tabs); ok {
		t.items = o.items
		t.activeKey = o.activeKey
		t.tabType = o.tabType
		t.onChange = o.onChange
		t.rebuild()
	}
}
