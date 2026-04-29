package components

import (
	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// Carousel is a simplified content slider.
type Carousel struct {
	core.BaseElement
	items       []core.Element
	currentIdx  int
	contentArea *View
}

// NewCarousel creates a carousel.
func NewCarousel(items []core.Element) *Carousel {
	c := &Carousel{
		items: items,
	}
	c.Init(c)
	c.SetFlexDirection(yoga.FlexDirectionColumn)
	c.SetGap(yoga.GutterAll, 8)

	c.contentArea = NewView()
	c.contentArea.SetFlexDirection(yoga.FlexDirectionRow)
	c.contentArea.SetOverflow(yoga.OverflowHidden)
	c.contentArea.SetWidthPercent(100)
	c.Add(c.contentArea)

	// Controls
	controls := NewView()
	controls.SetFlexDirection(yoga.FlexDirectionRow)
	controls.SetGap(yoga.GutterAll, 8)
	controls.SetJustifyContent(yoga.JustifyCenter)

	prev := NewButton("←").SetVariant(ButtonOutline)
	prev.SetOnClick(func() { c.prev() })
	next := NewButton("→").SetVariant(ButtonOutline)
	next.SetOnClick(func() { c.next() })
	controls.Add(prev, next)
	c.Add(controls)

	c.refresh()
	return c
}

// ElementType returns type identifier.
func (c *Carousel) ElementType() string { return "Carousel" }

func (c *Carousel) refresh() {
	c.contentArea.ClearChildren()
	if c.currentIdx >= 0 && c.currentIdx < len(c.items) {
		c.contentArea.Add(c.items[c.currentIdx])
	}
	c.Mark(core.FlagNeedLayout)
}

func (c *Carousel) prev() {
	if c.currentIdx > 0 {
		c.currentIdx--
	} else {
		c.currentIdx = len(c.items) - 1
	}
	c.refresh()
}

func (c *Carousel) next() {
	if c.currentIdx < len(c.items)-1 {
		c.currentIdx++
	} else {
		c.currentIdx = 0
	}
	c.refresh()
}

// SetIndex sets the current item index.
func (c *Carousel) SetIndex(idx int) *Carousel {
	if idx >= 0 && idx < len(c.items) {
		c.currentIdx = idx
		c.refresh()
	}
	return c
}
