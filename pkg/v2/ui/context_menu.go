package ui

import (
	"image/color"
)

// MenuItem 定义菜单项。
type MenuItem struct {
	Label    string
	Shortcut string       // 快捷键描述，如 "Ctrl+C"
	OnClick  func()
	Disabled bool
	Divider  bool         // 是否为分割线
	Children []MenuItem   // 子菜单
}

// ContextMenuWidget 右键菜单。
type ContextMenuWidget struct {
	BaseWidget
	Child    Widget
	Items    []MenuItem
	BgColor  color.Color
	TextColor color.Color
	HoverBg  color.Color
	Width    float32
}

// ContextMenu 创建右键菜单 Widget。
func ContextMenu(child Widget, items []MenuItem) ContextMenuWidget {
	return ContextMenuWidget{
		Child:     child,
		Items:     items,
		BgColor:   color.RGBA{R: 255, G: 255, B: 255, A: 255},
		TextColor: color.RGBA{R: 10, G: 10, B: 10, A: 255},
		HoverBg:   color.RGBA{R: 245, G: 245, B: 245, A: 255},
		Width:     200,
	}
}

func (c ContextMenuWidget) WithBgColor(cl color.Color) ContextMenuWidget { c.BgColor = cl; return c }
func (c ContextMenuWidget) WithTextColor(cl color.Color) ContextMenuWidget { c.TextColor = cl; return c }
func (c ContextMenuWidget) WithWidth(v float32) ContextMenuWidget { c.Width = v; return c }

func (c ContextMenuWidget) CreateElement() Element {
	return NewStatefulElement(c)
}

func (c ContextMenuWidget) CreateState() State {
	s := &contextMenuState{}
	s.Init(s)
	return s
}

type contextMenuState struct {
	BaseState
	visible bool
	x, y    float32
}

func (s *contextMenuState) InitState() {}

func (s *contextMenuState) Dispose() {}

func (s *contextMenuState) DidUpdateWidget(oldWidget Widget) {}

func (s *contextMenuState) Build(ctx BuildContext) Widget {
	w := s.GetWidget().(ContextMenuWidget)

	// 包裹子 Widget，监听右键点击
	return NewBuilder(func(innerCtx BuildContext) Widget {
		return w.Child
	})
}

// Show 显示菜单。
func (s *contextMenuState) Show(x, y float32) {
	s.visible = true
	s.x = x
	s.y = y
	s.SetState(nil)
}

// Hide 隐藏菜单。
func (s *contextMenuState) Hide() {
	s.visible = false
	s.SetState(nil)
}

// ==================== MenuBar ====================

// MenuBarWidget 顶部菜单栏。
type MenuBarWidget struct {
	BaseWidget
	Items     []MenuBarItem
	BgColor   color.Color
	TextColor color.Color
	Height    float32
}

// MenuBarItem 菜单栏项。
type MenuBarItem struct {
	Label string
	Items []MenuItem
}

// MenuBar 创建菜单栏 Widget。
func MenuBar(items []MenuBarItem) MenuBarWidget {
	return MenuBarWidget{
		Items:     items,
		BgColor:   color.RGBA{R: 245, G: 245, B: 245, A: 255},
		TextColor: color.RGBA{R: 10, G: 10, B: 10, A: 255},
		Height:    32,
	}
}

func (m MenuBarWidget) WithBgColor(c color.Color) MenuBarWidget { m.BgColor = c; return m }
func (m MenuBarWidget) WithHeight(v float32) MenuBarWidget { m.Height = v; return m }

func (m MenuBarWidget) CreateElement() Element {
	return NewStatefulElement(m)
}

func (m MenuBarWidget) CreateState() State {
	s := &menuBarState{}
	s.Init(s)
	return s
}

type menuBarState struct {
	BaseState
	activeIndex int
	dropVisible bool
}

func (s *menuBarState) InitState() {}
func (s *menuBarState) Dispose() {}
func (s *menuBarState) DidUpdateWidget(oldWidget Widget) {}

func (s *menuBarState) Build(ctx BuildContext) Widget {
	w := s.GetWidget().(MenuBarWidget)
	theme := ThemeOf(ctx)

	// 构建菜单栏
	items := make([]Widget, len(w.Items))
	for i, item := range w.Items {
		items[i] = NewBuilder(func(ctx BuildContext) Widget {
			return buildMenuBarItem(item, i == s.activeIndex && s.dropVisible, theme, func() {
				s.activeIndex = i
				s.dropVisible = !s.dropVisible
				s.SetState(nil)
			})
		})
	}

	return buildMenuBarContainer(items, w, theme)
}

func buildMenuBarItem(item MenuBarItem, active bool, theme *Theme, onClick func()) Widget {
	// 由具体实现构建
	return nil
}

func buildMenuBarContainer(items []Widget, w MenuBarWidget, theme *Theme) Widget {
	// 由具体实现构建
	return nil
}

// ==================== 辅助函数 ====================

// buildMenuItems 构建菜单项列表。
func buildMenuItems(items []MenuItem, theme *Theme, onClose func()) []Widget {
	result := make([]Widget, 0, len(items))
	for _, item := range items {
		if item.Divider {
			result = append(result, buildMenuDivider(theme))
			continue
		}
		result = append(result, buildMenuItemWidget(item, theme, onClose))
	}
	return result
}

func buildMenuItemWidget(item MenuItem, theme *Theme, onClose func()) Widget {
	// 由具体实现构建
	return nil
}

func buildMenuDivider(theme *Theme) Widget {
	return nil
}

// buildMenuContainer 构建菜单容器。
func buildMenuContainer(items []MenuItem, x, y, width float32, theme *Theme, onClose func()) Widget {
	menuItems := buildMenuItems(items, theme, onClose)
	if len(menuItems) == 0 {
		return nil
	}

	return NewBuilder(func(ctx BuildContext) Widget {
		// 垂直排列菜单项
		return nil
	})
}
