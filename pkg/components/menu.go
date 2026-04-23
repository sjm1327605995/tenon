package components

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// MenuItemData 定义一个菜单项的配置。
type MenuItemData struct {
	Key   string
	Label string
}

// MenuItem 是单个菜单项宿主组件。
type MenuItem struct {
	core.BaseHost
	selected  bool
	hovered   bool
	onClick   func()
	textComp  *Text
}

// NewMenuItem 创建一个菜单项。
func NewMenuItem(label string) *MenuItem {
	mi := &MenuItem{}
	mi.Init(mi)
	mi.GetElement().Yoga.StyleSetHeight(40)
	mi.GetElement().Yoga.StyleSetWidthPercent(100)
	mi.GetElement().Yoga.StyleSetAlignItems(yoga.AlignCenter) // 垂直居中

	mi.textComp = NewText(label)
	mi.textComp.SetMargin(yoga.EdgeLeft, 24)
	mi.AddChild(mi.textComp)
	mi.refreshStyle()
	return mi
}

func (mi *MenuItem) refreshStyle() {
	theme := core.GetTheme()
	if mi.selected {
		mi.GetElement().BackgroundColor = theme.MenuItemSelectedBg
		mi.textComp.SetColor(theme.MenuItemSelectedText)
	} else {
		mi.GetElement().BackgroundColor = nil
		mi.textComp.SetColor(theme.TextColor)
	}
}

// SetSelected 设置选中状态。
func (mi *MenuItem) SetSelected(selected bool) *MenuItem {
	mi.selected = selected
	mi.refreshStyle()
	return mi
}

// SetOnClick 设置点击回调。
func (mi *MenuItem) SetOnClick(fn func()) *MenuItem {
	mi.onClick = fn
	return mi
}

// Draw 绘制左边选中竖条。
func (mi *MenuItem) Draw(screen *ebiten.Image) {
	el := mi.GetElement()
	if el == nil || !el.Visible {
		return
	}
	bounds := mi.GetLayoutBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	// 绘制背景（已由 Element.BackgroundColor 处理，但 nil 时引擎不绘，这里确保透明背景被覆盖）
	if el.BackgroundColor != nil {
		vector.FillRect(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, el.BackgroundColor, false)
	}

	// Hover 时绘制圆角背景
	if mi.hovered && !mi.selected {
		// Ant Design 风格的 hover 背景色
		hoverColor := color.RGBA{R: 245, G: 245, B: 245, A: 255}
		// 左右各留 4px 边距，形成圆角背景效果
		drawRoundedRect(screen, bounds.X+4, bounds.Y+4, bounds.Width-8, bounds.Height-8, 4, hoverColor)
	}

	// 选中时绘制左侧蓝色竖条
	if mi.selected {
		theme := core.GetTheme()
		vector.FillRect(screen, bounds.X, bounds.Y, 3, bounds.Height, theme.PrimaryColor, false)
	}
}

// drawRoundedRect 绘制圆角矩形填充。
func drawRoundedRect(screen *ebiten.Image, x, y, w, h, r float32, c color.Color) {
	if w <= 0 || h <= 0 {
		return
	}
	path := vector.Path{}
	// 左上角开始，顺时针绘制
	path.MoveTo(x+r, y)
	path.LineTo(x+w-r, y)
	path.Arc(x+w-r, y+r, r, -math.Pi/2, 0, vector.Clockwise)
	path.LineTo(x+w, y+h-r)
	path.Arc(x+w-r, y+h-r, r, 0, math.Pi/2, vector.Clockwise)
	path.LineTo(x+r, y+h)
	path.Arc(x+r, y+h-r, r, math.Pi/2, math.Pi, vector.Clockwise)
	path.LineTo(x, y+r)
	path.Arc(x+r, y+r, r, math.Pi, math.Pi*1.5, vector.Clockwise)
	path.Close()
	op := &vector.DrawPathOptions{}
	op.ColorScale.ScaleWithColor(c)
	vector.FillPath(screen, &path, &vector.FillOptions{}, op)
}

// Update 每帧检测鼠标是否悬停。
func (mi *MenuItem) Update() error {
	mx, my := ebiten.CursorPosition()
	bounds := mi.GetLayoutBounds()
	mi.hovered = float32(mx) >= bounds.X && float32(mx) < bounds.X+bounds.Width &&
		float32(my) >= bounds.Y && float32(my) < bounds.Y+bounds.Height
	return nil
}

// HandleEvent 处理点击事件。
func (mi *MenuItem) HandleEvent(e *core.Event) bool {
	switch e.Type {
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
		if mi.textComp != nil && o.textComp != nil {
			mi.textComp.Content = o.textComp.Content
			mi.textComp.cachedLayout = nil
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
				if mi.textComp != nil && mi.textComp.Content == item.Label {
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
