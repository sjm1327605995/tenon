package components

import (
	"github.com/sjm1327605995/tenon/internal/core"
	"github.com/sjm1327605995/tenon/internal/native"
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
	panel       *native.View
	listView    *ListView
	onChange    func(index int, value string)
}

// NewDropdown creates a Dropdown.
func NewDropdown() *Dropdown {
	dd := &Dropdown{selectedIdx: -1}
	dd.Init(dd)
	dd.SetFlexDirection(yoga.FlexDirectionColumn)
	dd.SetMinWidth(120)
	dd.setupUI()
	return dd
}

// ElementType returns type identifier.
func (dd *Dropdown) ElementType() string { return "Dropdown" }

func (dd *Dropdown) setupUI() {
	dd.trigger = NewButton("").SetOnClick(func() { dd.toggle() })
	dd.BaseElement.AppendChild(dd.trigger)

	th := core.GetTheme()
	dd.panel = native.NewView()
	dd.panel.SetFlexDirection(yoga.FlexDirectionColumn)
	dd.panel.SetBackgroundColor(th.SurfaceColor)
	dd.panel.SetBorder(yoga.EdgeAll, 1)
	dd.panel.SetBorderColor(th.BorderColor)
	dd.panel.SetShadow(th.ShadowColor, 8, 0, 4)
	dd.panel.SetVisible(false)
	dd.panel.SetDisplay(yoga.DisplayNone)
	dd.panel.SetWidthPercent(100)
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
	dd.Mark(core.FlagNeedLayout)
}

func (dd *Dropdown) close() {
	if !dd.isOpen {
		return
	}
	dd.isOpen = false
	dd.panel.SetVisible(false)
	dd.panel.SetDisplay(yoga.DisplayNone)
	dd.Mark(core.FlagNeedLayout)
}

// SetItems replaces all options.
func (dd *Dropdown) SetItems(items []string) *Dropdown {
	dd.items = append([]string{}, items...)
	dd.listView.Clear()
	for i, item := range dd.items {
		txt := native.NewText(item)
		dd.listView.AddItem(txt)
		_ = i
	}
	if dd.selectedIdx >= len(dd.items) {
		dd.selectedIdx = -1
		dd.trigger.SetText("")
	}
	dd.Mark(core.FlagNeedLayout)
	return dd
}

// AddItem appends an option.
func (dd *Dropdown) AddItem(item string) *Dropdown {
	dd.items = append(dd.items, item)
	txt := native.NewText(item)
	dd.listView.AddItem(txt)
	dd.Mark(core.FlagNeedLayout)
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
// 注意：isOpen 是命令式状态，不由声明式重建控制，避免重建时意外关闭。
func (dd *Dropdown) SyncFrom(src core.Element) {
	other, ok := src.(*Dropdown)
	if !ok {
		return
	}
	if dd.selectedIdx != other.selectedIdx {
		dd.selectedIdx = other.selectedIdx
		if dd.selectedIdx >= 0 && dd.selectedIdx < len(dd.items) {
			dd.trigger.SetText(dd.items[dd.selectedIdx])
		} else {
			dd.trigger.SetText("")
		}
		// native.Text child auto-marks dirty; no need to mark composite
	}
}
