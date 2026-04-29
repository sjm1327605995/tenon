package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/v2/core"
)

// ListView is a vertical scrollable list with selectable items.
type ListView struct {
	core.BaseElement
	scrollView   *ScrollView
	items        []*listItemWrapper
	selectedIdx  int
	onSelect     func(index int)
	bgNormal     color.Color
	bgSelected   color.Color
	bgHover      color.Color
}

// listItemWrapper wraps a user-provided element and handles click selection.
type listItemWrapper struct {
	core.BaseElement
	content    core.Element
	index      int
	listView   *ListView
	isHovered  bool
}

// NewListView creates a ListView.
func NewListView() *ListView {
	lv := &ListView{
		selectedIdx: -1,
		bgNormal:    nil,
		bgSelected:  color.RGBA{200, 220, 255, 255},
		bgHover:     color.RGBA{230, 230, 230, 255},
	}
	lv.Init(lv)
	lv.scrollView = NewScrollView()
	lv.scrollView.SetWidthPercent(100)
	lv.scrollView.SetHeightPercent(100)
	lv.BaseElement.AppendChild(lv.scrollView)
	return lv
}

// ElementType returns type identifier.
func (lv *ListView) ElementType() string { return "ListView" }

// AddItem appends an element as a list row.
func (lv *ListView) AddItem(item core.Element) *ListView {
	wrapper := newListItemWrapper(lv, item, len(lv.items))
	lv.items = append(lv.items, wrapper)
	lv.scrollView.Content().AppendChild(wrapper)
	return lv
}

// RemoveItem removes the item at the given index.
func (lv *ListView) RemoveItem(index int) *ListView {
	if index < 0 || index >= len(lv.items) {
		return lv
	}
	wrapper := lv.items[index]
	lv.scrollView.Content().RemoveChild(wrapper)
	lv.items = append(lv.items[:index], lv.items[index+1:]...)
	// Re-index remaining items
	for i := index; i < len(lv.items); i++ {
		lv.items[i].index = i
	}
	if lv.selectedIdx == index {
		lv.selectedIdx = -1
	} else if lv.selectedIdx > index {
		lv.selectedIdx--
	}
	lv.Mark(core.FlagNeedDraw)
	return lv
}

// Clear removes all items.
func (lv *ListView) Clear() *ListView {
	lv.scrollView.Content().ClearChildren()
	lv.items = lv.items[:0]
	lv.selectedIdx = -1
	lv.Mark(core.FlagNeedDraw)
	return lv
}

// Select sets the selected index (-1 for none).
func (lv *ListView) Select(index int) *ListView {
	lv.selectItem(index)
	return lv
}

func (lv *ListView) selectItem(index int) {
	if lv.selectedIdx == index {
		return
	}
	lv.selectedIdx = index
	lv.Mark(core.FlagNeedDraw)
	if lv.onSelect != nil && index >= 0 && index < len(lv.items) {
		lv.onSelect(index)
	}
}

// SelectedIndex returns the current selected index.
func (lv *ListView) SelectedIndex() int { return lv.selectedIdx }

// SelectedItem returns the currently selected wrapper, or nil.
func (lv *ListView) SelectedItem() core.Element {
	if lv.selectedIdx < 0 || lv.selectedIdx >= len(lv.items) {
		return nil
	}
	return lv.items[lv.selectedIdx].content
}

// OnSelect sets the selection callback.
func (lv *ListView) OnSelect(fn func(index int)) *ListView {
	lv.onSelect = fn
	return lv
}

// SetBackgroundColors sets normal/selected/hover row backgrounds.
func (lv *ListView) SetBackgroundColors(normal, selected, hover color.Color) *ListView {
	lv.bgNormal = normal
	lv.bgSelected = selected
	lv.bgHover = hover
	lv.Mark(core.FlagNeedDraw)
	return lv
}

// SyncFrom 同步新 ListView 的属性到当前 Element（声明式重建）。
func (lv *ListView) SyncFrom(src core.Element) {
	other, ok := src.(*ListView)
	if !ok {
		return
	}
	sb := &SyncBuilder{}
	syncField(sb, &lv.selectedIdx, other.selectedIdx)
	syncColor(sb, &lv.bgNormal, other.bgNormal)
	syncColor(sb, &lv.bgSelected, other.bgSelected)
	syncColor(sb, &lv.bgHover, other.bgHover)
	sb.MarkDraw(lv)
}

// ScrollView returns the internal scroll view for advanced customization.
func (lv *ListView) ScrollView() *ScrollView { return lv.scrollView }

// ==================== listItemWrapper ====================

func newListItemWrapper(lv *ListView, content core.Element, index int) *listItemWrapper {
	w := &listItemWrapper{
		content:  content,
		index:    index,
		listView: lv,
	}
	w.Init(w)
	w.SetWidthPercent(100)
	w.Add(content)
	return w
}

func (w *listItemWrapper) ElementType() string { return "ListItem" }

func (w *listItemWrapper) Draw(screen *ebiten.Image) {
	bounds := w.GetBounds()

	var bg color.Color
	if w.listView.selectedIdx == w.index && w.listView.bgSelected != nil {
		bg = w.listView.bgSelected
	} else if w.isHovered && w.listView.bgHover != nil {
		bg = w.listView.bgHover
	} else if w.listView.bgNormal != nil {
		bg = w.listView.bgNormal
	}

	if bg != nil {
		vector.FillRect(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, bg, false)
	}
}

func (w *listItemWrapper) HandleEvent(e *core.Event) bool {
	switch e.Type {
	case core.EventMouseMove:
		if !w.isHovered {
			w.isHovered = true
			w.Mark(core.FlagNeedDraw)
		}
		return false // let it bubble for scrollview etc.
	case core.EventClick:
		w.listView.selectItem(w.index)
		return true
	}
	return false
}

func (w *listItemWrapper) Update() error {
	// Check if mouse is still over this item; clear hover if not.
	mx, my := ebiten.CursorPosition()
	b := w.GetBounds()
	inside := float32(mx) >= b.X && float32(mx) < b.X+b.Width &&
		float32(my) >= b.Y && float32(my) < b.Y+b.Height
	if w.isHovered && !inside {
		w.isHovered = false
		w.Mark(core.FlagNeedDraw)
	}
	return nil
}
