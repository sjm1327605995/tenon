package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// TabItem defines a single tab label and its content builder.
type TabItem struct {
	Label   string
	Content func() core.Element
}

// Tab is a tab container with a tab bar and content area.
type Tab struct {
	core.BaseElement
	items       []TabItem
	activeIndex int
	tabBar      *View
	contentArea *View
	tabButtons  []*Button
	underline   *View
	onChange    func(index int, label string)
	activeColor color.Color
	inactiveColor color.Color
	indicatorColor color.Color
	barBgColor  color.Color
}

// NewTab creates a tab component.
func NewTab(items []TabItem) *Tab {
	theme := core.GetTheme()
	t := &Tab{
		items:          items,
		activeIndex:    0,
		activeColor:    theme.PrimaryColor,
		inactiveColor:  theme.TextMutedColor,
		indicatorColor: theme.PrimaryColor,
		barBgColor:     theme.SurfaceColor,
	}
	t.Init(t)
	t.SetFlexDirection(yoga.FlexDirectionColumn)

	// Tab bar
	t.tabBar = NewView()
	t.tabBar.SetFlexDirection(yoga.FlexDirectionRow)
	t.tabBar.SetAlignItems(yoga.AlignCenter)
	t.tabBar.SetBackgroundColor(t.barBgColor)
	t.tabBar.SetBorder(yoga.EdgeBottom, 1)
	t.tabBar.SetBorderColor(theme.BorderColor)

	// Content area
	t.contentArea = NewView()
	t.contentArea.SetFlexGrow(1)
	t.contentArea.SetWidthPercent(100)

	// Underline indicator
	t.underline = NewView()
	t.underline.SetBackgroundColor(t.indicatorColor)
	t.underline.SetHeight(2)
	t.underline.SetPositionType(yoga.PositionTypeAbsolute)
	t.underline.SetPosition(yoga.EdgeBottom, 0)

	t.buildTabs()
	t.tabBar.Add(t.underline)
	t.Add(t.tabBar, t.contentArea)
	return t
}

// ElementType returns type identifier.
func (t *Tab) ElementType() string { return "Tab" }

// Draw renders the tab bar underline and background.
func (t *Tab) Draw(screen *ebiten.Image) {
	bounds := t.GetBounds()
	// Background of the whole tab container
	if t.barBgColor != nil {
		vector.FillRect(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, t.barBgColor, false)
	}
}

// Update repositions the underline indicator.
func (t *Tab) Update() error {
	if t.activeIndex >= 0 && t.activeIndex < len(t.tabButtons) {
		btn := t.tabButtons[t.activeIndex]
		bb := btn.GetBounds()
		tbBounds := t.tabBar.GetBounds()
		if bb.Width > 0 {
			t.underline.SetWidth(bb.Width)
			t.underline.SetPosition(yoga.EdgeLeft, bb.X-tbBounds.X)
		}
	}
	return nil
}

func (t *Tab) buildTabs() {
	t.tabBar.ClearChildren()
	t.tabButtons = make([]*Button, 0, len(t.items))

	for i, item := range t.items {
		idx := i
		btn := NewButton(item.Label).
			SetColors(t.barBgColor, t.barBgColor, t.barBgColor)
		if i == t.activeIndex {
			btn.labelEl.SetColor(t.activeColor)
		} else {
			btn.labelEl.SetColor(t.inactiveColor)
		}
		btn.SetOnClick(func() {
			t.setActive(idx)
		})
		// btn.SetBorderRadius(0) // Button does not have SetBorderRadius
		btn.SetPadding(yoga.EdgeHorizontal, 16)
		btn.SetPadding(yoga.EdgeVertical, 12)
		btn.SetBorder(yoga.EdgeBottom, 0)

		t.tabButtons = append(t.tabButtons, btn)
		t.tabBar.AppendChild(btn)
	}

	// Re-add underline (it must be last child to render on top)
	t.tabBar.AppendChild(t.underline)
	t.refreshContent()
}

func (t *Tab) setActive(index int) {
	if index < 0 || index >= len(t.items) || index == t.activeIndex {
		return
	}
	t.activeIndex = index
	for i, btn := range t.tabButtons {
		if i == t.activeIndex {
			btn.labelEl.SetColor(t.activeColor)
		} else {
			btn.labelEl.SetColor(t.inactiveColor)
		}
	}
	t.refreshContent()
	t.Mark(core.FlagNeedLayout | core.FlagNeedDraw)
	if t.onChange != nil {
		t.onChange(t.activeIndex, t.items[t.activeIndex].Label)
	}
}

func (t *Tab) refreshContent() {
	t.contentArea.ClearChildren()
	if t.activeIndex >= 0 && t.activeIndex < len(t.items) {
		if builder := t.items[t.activeIndex].Content; builder != nil {
			el := builder()
			if el != nil {
				t.contentArea.AppendChild(el)
			}
		}
	}
}

// SetActive sets the active tab index.
func (t *Tab) SetActive(index int) *Tab {
	t.setActive(index)
	return t
}

// GetActive returns the current active index.
func (t *Tab) GetActive() int {
	return t.activeIndex
}

// SetOnChange sets the tab switch callback.
func (t *Tab) SetOnChange(fn func(index int, label string)) *Tab {
	t.onChange = fn
	return t
}

// SetItems rebuilds tabs with new items.
func (t *Tab) SetItems(items []TabItem) *Tab {
	t.items = items
	t.activeIndex = 0
	if len(items) == 0 {
		t.activeIndex = -1
	}
	t.buildTabs()
	t.Mark(core.FlagNeedLayout | core.FlagNeedDraw)
	return t
}

// SetActiveColor sets the active tab text and indicator color.
func (t *Tab) SetActiveColor(clr color.Color) *Tab {
	t.activeColor = clr
	t.indicatorColor = clr
	t.underline.SetBackgroundColor(clr)
	t.setActive(t.activeIndex)
	return t
}

// SetInactiveColor sets the inactive tab text color.
func (t *Tab) SetInactiveColor(clr color.Color) *Tab {
	t.inactiveColor = clr
	t.setActive(t.activeIndex)
	return t
}

// SetBarBgColor sets the tab bar background.
func (t *Tab) SetBarBgColor(clr color.Color) *Tab {
	t.barBgColor = clr
	t.tabBar.SetBackgroundColor(clr)
	t.Mark(core.FlagNeedDraw)
	return t
}
