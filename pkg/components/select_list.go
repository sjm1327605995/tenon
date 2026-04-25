package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// SelectItem 定义选择列表中的一个选项。
type SelectItem struct {
	Key      string
	Label    string
	SubLabel string // 副标题，如卡片的 ATK/DEF
	Selected bool
	Disabled bool
}

// SelectList 是选择列表组件，支持单选/多选。
type SelectList struct {
	core.BaseHost
	items       []SelectItem
	multiSelect bool
	onChange    func(selected []string)
	itemHosts   map[string]*View
	checkHosts  map[string]*View
	labelHosts  map[string]*Text
}

// NewSelectList 创建一个选择列表。
func NewSelectList() *SelectList {
	sl := &SelectList{
		itemHosts:  make(map[string]*View),
		checkHosts: make(map[string]*View),
		labelHosts: make(map[string]*Text),
	}
	sl.Init(sl)
	sl.GetElement().Yoga.StyleSetFlexDirection(yoga.FlexDirectionColumn)
	sl.GetElement().Yoga.StyleSetGap(yoga.GutterRow, 2)
	sl.GetElement().Yoga.StyleSetMaxHeight(300)
	sl.GetElement().Yoga.StyleSetOverflow(yoga.OverflowScroll)
	return sl
}

// SetMultiSelect 设置是否多选。
func (sl *SelectList) SetMultiSelect(multi bool) *SelectList {
	sl.multiSelect = multi
	return sl
}

// SetItems 设置选项数据。
func (sl *SelectList) SetItems(items []SelectItem) *SelectList {
	sl.items = items
	sl.rebuild()
	return sl
}

// GetSelected 获取已选中的 key 列表。
func (sl *SelectList) GetSelected() []string {
	var selected []string
	for _, item := range sl.items {
		if item.Selected {
			selected = append(selected, item.Key)
		}
	}
	return selected
}

// SetOnChange 设置选择变化回调。
func (sl *SelectList) SetOnChange(fn func(selected []string)) *SelectList {
	sl.onChange = fn
	return sl
}

func (sl *SelectList) rebuild() {
	sl.ClearChildren()
	sl.itemHosts = make(map[string]*View)
	sl.checkHosts = make(map[string]*View)
	sl.labelHosts = make(map[string]*Text)

	for _, item := range sl.items {
		row := NewView()
		row.Init(row)
		row.GetElement().Yoga.StyleSetFlexDirection(yoga.FlexDirectionRow)
		row.GetElement().Yoga.StyleSetAlignItems(yoga.AlignCenter)
		row.GetElement().Yoga.StyleSetPadding(yoga.EdgeVertical, 8)
		row.GetElement().Yoga.StyleSetPadding(yoga.EdgeHorizontal, 12)
		row.GetElement().Yoga.StyleSetBorderRadius(4)
		row.GetElement().PointerEvents = core.PointerEventsAuto

		// 选中指示器
		check := NewView()
		check.Init(check)
		check.GetElement().Yoga.StyleSetWidth(16)
		check.GetElement().Yoga.StyleSetHeight(16)
		check.GetElement().Yoga.StyleSetMargin(yoga.EdgeRight, 8)
		check.GetElement().Yoga.StyleSetBorderRadius(3)
		check.GetElement().PointerEvents = core.PointerEventsNone
		if item.Selected {
			check.GetElement().BackgroundColor = core.GetTheme().PrimaryColor
		} else {
			check.GetElement().BackgroundColor = core.GetTheme().SurfaceColor
			check.GetElement().Yoga.StyleSetBorder(yoga.EdgeAll, 1)
			check.GetElement().BorderColor = core.GetTheme().BorderColor
		}
		row.AddChild(check)
		sl.checkHosts[item.Key] = check

		// 标签
		label := NewText(item.Label)
		label.SetFontSize(core.GetTheme().FontSizeBase)
		if item.Disabled {
			label.SetColor(core.GetTheme().TextMutedColor)
		} else if item.Selected {
			label.SetColor(core.GetTheme().PrimaryColor)
		} else {
			label.SetColor(core.GetTheme().TextColor)
		}
		label.GetElement().PointerEvents = core.PointerEventsNone
		row.AddChild(label)
		sl.labelHosts[item.Key] = label

		// 副标签
		if item.SubLabel != "" {
			subLabel := NewText(item.SubLabel)
			subLabel.SetFontSize(core.GetTheme().FontSizeSM)
			subLabel.SetColor(core.GetTheme().TextMutedColor)
			subLabel.GetElement().Yoga.StyleSetMargin(yoga.EdgeLeft, 8)
			subLabel.GetElement().PointerEvents = core.PointerEventsNone
			row.AddChild(subLabel)
		}

		// 悬停效果
		row.SetOnHover(func(hovered bool) {
			if !item.Disabled && !item.Selected {
				if hovered {
					row.GetElement().BackgroundColor = core.GetTheme().MenuItemHoverBg
				} else {
					row.GetElement().BackgroundColor = nil
				}
			}
		})

		key := item.Key
		row.SetOnClick(func() {
			sl.toggleItem(key)
		})

		sl.itemHosts[item.Key] = row
		sl.AddChild(row)
	}
}

func (sl *SelectList) toggleItem(key string) {
	for i := range sl.items {
		if sl.items[i].Key != key {
			continue
		}
		if sl.items[i].Disabled {
			return
		}
		if sl.multiSelect {
			sl.items[i].Selected = !sl.items[i].Selected
		} else {
			// 单选：取消其他选中
			for j := range sl.items {
				if j != i {
					sl.items[j].Selected = false
				}
			}
			sl.items[i].Selected = true
		}
		break
	}
	sl.refreshSelection()
	if sl.onChange != nil {
		sl.onChange(sl.GetSelected())
	}
}

func (sl *SelectList) refreshSelection() {
	for _, item := range sl.items {
		check, ok := sl.checkHosts[item.Key]
		if ok {
			if item.Selected {
				check.GetElement().BackgroundColor = core.GetTheme().PrimaryColor
				check.GetElement().BorderColor = nil
				check.GetElement().Yoga.StyleSetBorder(yoga.EdgeAll, 0)
			} else {
				check.GetElement().BackgroundColor = core.GetTheme().SurfaceColor
				check.GetElement().BorderColor = core.GetTheme().BorderColor
				check.GetElement().Yoga.StyleSetBorder(yoga.EdgeAll, 1)
			}
		}
		label, ok := sl.labelHosts[item.Key]
		if ok {
			if item.Disabled {
				label.SetColor(core.GetTheme().TextMutedColor)
			} else if item.Selected {
				label.SetColor(core.GetTheme().PrimaryColor)
			} else {
				label.SetColor(core.GetTheme().TextColor)
			}
		}
	}
}

// ==================== 链式 API ====================

func (sl *SelectList) SetWidth(width float32) *SelectList {
	sl.GetElement().Yoga.StyleSetWidth(width)
	return sl
}
func (sl *SelectList) SetHeight(height float32) *SelectList {
	sl.GetElement().Yoga.StyleSetHeight(height)
	return sl
}
func (sl *SelectList) SetMargin(edge yoga.Edge, value float32) *SelectList {
	sl.GetElement().Yoga.StyleSetMargin(edge, value)
	return sl
}

// SyncFrom 同步选择列表属性。
func (sl *SelectList) SyncFrom(other core.Host) {
	if o, ok := other.(*SelectList); ok {
		sl.items = o.items
		sl.multiSelect = o.multiSelect
		sl.onChange = o.onChange
		sl.rebuild()
	}
}
