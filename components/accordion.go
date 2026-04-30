package components

import (
	"github.com/sjm1327605995/tenon/internal/core"
	"github.com/sjm1327605995/tenon/internal/native"
	"github.com/sjm1327605995/tenon/yoga"
)

// AccordionItem defines a single accordion entry.
type AccordionItem struct {
	Title   string
	Content func() core.Element
}

// Accordion is a vertically stacked set of collapsible panels.
type Accordion struct {
	core.BaseElement
	items      []AccordionItem
	expanded   map[int]bool
	single     bool
	contentMap map[int]*native.View
}

// NewAccordion creates an accordion.
func NewAccordion(items []AccordionItem) *Accordion {
	a := &Accordion{
		items:      items,
		expanded:   make(map[int]bool),
		single:     false,
		contentMap: make(map[int]*native.View),
	}
	a.Init(a)
	a.SetFlexDirection(yoga.FlexDirectionColumn)
	a.SetWidthPercent(100)
	a.buildItems()
	return a
}

// ElementType returns type identifier.
func (a *Accordion) ElementType() string { return "Accordion" }

func (a *Accordion) buildItems() {
	for i, item := range a.items {
		idx := i
		// Header row
		header := native.NewView()
		header.SetFlexDirection(yoga.FlexDirectionRow)
		header.SetAlignItems(yoga.AlignCenter)
		header.SetJustifyContent(yoga.JustifySpaceBetween)
		header.SetPadding(yoga.EdgeVertical, 12)
		header.SetWidthPercent(100)

		title := native.NewText(item.Title).SetFontSize(14).SetColor(core.GetTheme().TextColor)
		chevron := native.NewText("▶").SetFontSize(12).SetColor(core.GetTheme().MutedForegroundColor)

		header.Add(title, native.NewView().SetFlexGrow(1), chevron)

		// Clickable trigger wrapper
		trigger := NewButton("")
		trigger.SetOnClick(func() { a.toggle(idx) })
		trigger.SetVariant(ButtonGhost)
		trigger.labelEl.SetVisible(false)
		trigger.SetWidthPercent(100)
		trigger.BaseElement.SetPadding(yoga.EdgeAll, 0)
		trigger.RemoveChild(trigger.labelEl)
		trigger.AppendChild(header)

		a.AppendChild(trigger)

		// Content wrapper
		contentWrap := native.NewView()
		contentWrap.SetFlexDirection(yoga.FlexDirectionColumn)
		contentWrap.SetWidthPercent(100)
		if item.Content != nil {
			if el := item.Content(); el != nil {
				contentWrap.Add(el)
			}
		}
		contentWrap.SetVisible(false)
		a.contentMap[i] = contentWrap
		a.AppendChild(contentWrap)

		// Separator between items
		if i < len(a.items)-1 {
			a.AppendChild(native.NewDivider())
		}
	}
}

func (a *Accordion) toggle(idx int) {
	if a.single {
		for i := range a.expanded {
			if i != idx {
				a.expanded[i] = false
				a.contentMap[i].SetVisible(false)
			}
		}
	}
	a.expanded[idx] = !a.expanded[idx]
	a.contentMap[idx].SetVisible(a.expanded[idx])
	a.Mark(core.FlagNeedLayout)
}

// Expand opens a specific item.
func (a *Accordion) Expand(idx int) *Accordion {
	if idx >= 0 && idx < len(a.items) {
		a.toggle(idx)
	}
	return a
}

// SetSingle sets whether only one item can be expanded at a time.
func (a *Accordion) SetSingle(single bool) *Accordion {
	a.single = single
	return a
}
