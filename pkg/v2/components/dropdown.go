package components

import (
	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// Dropdown is a select control that shows a trigger button and
// expands a scrollable list of options below it.
type Dropdown struct {
	core.BaseElement
	items       []string
	selectedIdx int
	isOpen      bool
	trigger     *Button
	panel       *View
	listView    *ListView
	onChange    func(index int, value string)
}

// NewDropdown creates a Dropdown.
func NewDropdown() *Dropdown {
	dd := &Dropdown{selectedIdx: -1}
	dd.Init(dd)
	dd.SetFlexDirection(yoga.FlexDirectionColumn)
	dd.setupUI()
	return dd
}

// ElementType returns type identifier.
func (dd *Dropdown) ElementType() string { return "Dropdown" }

func (dd *Dropdown) setupUI() {
	dd.trigger = NewButton("").SetOnClick(func() { dd.toggle() })
	dd.BaseElement.AppendChild(dd.trigger)

	th := core.GetTheme()
	dd.panel = NewView()
	dd.panel.SetFlexDirection(yoga.FlexDirectionColumn)
	dd.panel.SetBackgroundColor(th.SurfaceColor)
	dd.panel.SetBorder(yoga.EdgeAll, 1)
	dd.panel.SetBorderColor(th.BorderColor)
	dd.panel.SetShadow(th.ShadowColor, 8, 0, 4)
	dd.panel.SetVisible(false)
	dd.panel.SetDisplay(yoga.DisplayNone)
	dd.BaseElement.AppendChild(dd.panel)

	dd.listView = NewListView()
	dd.listView.OnSelect(func(idx int) {
		dd.selectedIdx = idx
		if idx >= 0 && idx < len(dd.items) {
			dd.trigger.SetText(dd.items[idx])
		}
		dd.close()
		if dd.onChange != nil {
			dd.onChange(idx, dd.items[idx])
		}
	})
	dd.listView.ScrollView().SetMaxHeight(200)
	dd.panel.Add(dd.listView)
}

func (dd *Dropdown) toggle() {
	if dd.isOpen {
		dd.close()
	} else {
		dd.open()
	}
}

func (dd *Dropdown) open() {
	if dd.isOpen {
		return
	}
	dd.isOpen = true
	dd.panel.SetVisible(true)
	dd.panel.SetDisplay(yoga.DisplayFlex)
	dd.panel.SetMinWidth(dd.trigger.GetBounds().Width)
	dd.Mark(core.FlagNeedLayout | core.FlagNeedDraw)
}

func (dd *Dropdown) close() {
	if !dd.isOpen {
		return
	}
	dd.isOpen = false
	dd.panel.SetVisible(false)
	dd.panel.SetDisplay(yoga.DisplayNone)
	dd.Mark(core.FlagNeedLayout | core.FlagNeedDraw)
}

// SetItems replaces all options.
func (dd *Dropdown) SetItems(items []string) *Dropdown {
	dd.items = append([]string{}, items...)
	dd.listView.Clear()
	for i, item := range dd.items {
		txt := NewText(item)
		dd.listView.AddItem(txt)
		_ = i
	}
	if dd.selectedIdx >= len(dd.items) {
		dd.selectedIdx = -1
		dd.trigger.SetText("")
	}
	dd.Mark(core.FlagNeedLayout | core.FlagNeedDraw)
	return dd
}

// AddItem appends an option.
func (dd *Dropdown) AddItem(item string) *Dropdown {
	dd.items = append(dd.items, item)
	txt := NewText(item)
	dd.listView.AddItem(txt)
	dd.Mark(core.FlagNeedLayout | core.FlagNeedDraw)
	return dd
}

// Select sets the selected index.
func (dd *Dropdown) Select(index int) *Dropdown {
	if index < -1 || index >= len(dd.items) {
		return dd
	}
	dd.selectedIdx = index
	if index >= 0 {
		dd.trigger.SetText(dd.items[index])
	} else {
		dd.trigger.SetText("")
	}
	dd.listView.Select(index)
	return dd
}

// SelectedIndex returns the current selection.
func (dd *Dropdown) SelectedIndex() int { return dd.selectedIdx }

// SelectedValue returns the currently selected text.
func (dd *Dropdown) SelectedValue() string {
	if dd.selectedIdx < 0 || dd.selectedIdx >= len(dd.items) {
		return ""
	}
	return dd.items[dd.selectedIdx]
}

// OnChange sets the change callback.
func (dd *Dropdown) OnChange(fn func(index int, value string)) *Dropdown {
	dd.onChange = fn
	return dd
}

// SetPlaceholder sets the trigger text when nothing is selected.
func (dd *Dropdown) SetPlaceholder(text string) *Dropdown {
	if dd.selectedIdx < 0 {
		dd.trigger.SetText(text)
	}
	return dd
}

// IsOpen returns whether the dropdown panel is visible.
func (dd *Dropdown) IsOpen() bool { return dd.isOpen }

// HandleEvent closes the dropdown on outside clicks.
func (dd *Dropdown) HandleEvent(e *core.Event) bool {
	if e.Type == core.EventClick && dd.isOpen {
		target := e.Target
		for target != nil {
			if target == dd || target == dd.panel {
				return false
			}
			target = target.GetParent()
		}
		dd.close()
		return false
	}
	return false
}

// SyncFrom 同步新 Dropdown 的属性到当前 Element（声明式重建）。
func (dd *Dropdown) SyncFrom(src core.Element) {
	other, ok := src.(*Dropdown)
	if !ok {
		return
	}
	needDraw := false
	layout := false
	if dd.selectedIdx != other.selectedIdx {
		dd.selectedIdx = other.selectedIdx
		if dd.selectedIdx >= 0 && dd.selectedIdx < len(dd.items) {
			dd.trigger.SetText(dd.items[dd.selectedIdx])
		} else {
			dd.trigger.SetText("")
		}
		needDraw = true
	}
	if dd.isOpen != other.isOpen {
		dd.isOpen = other.isOpen
		if dd.isOpen {
			dd.panel.SetVisible(true)
			dd.panel.SetDisplay(yoga.DisplayFlex)
		} else {
			dd.panel.SetVisible(false)
			dd.panel.SetDisplay(yoga.DisplayNone)
		}
		needDraw = true
		layout = true
	}
	if needDraw {
		dd.Mark(core.FlagNeedDraw)
	}
	if layout {
		dd.Mark(core.FlagNeedLayout)
	}
}
