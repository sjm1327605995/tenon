package components

import (
	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// HoverCard is a floating card shown on hover.
type HoverCard struct {
	core.BaseElement
	trigger *View
	content *View
	visible bool
}

// NewHoverCard creates a hover card.
func NewHoverCard(trigger core.Element) *HoverCard {
	h := &HoverCard{}
	h.Init(h)
	h.SetFlexDirection(yoga.FlexDirectionColumn)

	if t, ok := trigger.(*View); ok {
		h.trigger = t
	} else {
		h.trigger = NewView()
		h.trigger.Add(trigger)
	}
	h.Add(h.trigger)

	h.content = NewView()
	h.content.SetFlexDirection(yoga.FlexDirectionColumn)
	h.content.SetBackgroundColor(core.GetTheme().CardColor)
	h.content.SetBorderColor(core.GetTheme().BorderColor)
	h.content.SetBorderRadius(core.GetTheme().BorderRadius)
	h.content.SetPadding(yoga.EdgeAll, 16)
	h.content.SetPositionType(yoga.PositionTypeAbsolute)
	h.content.SetVisible(false)
	h.content.SetShadow(core.GetTheme().ShadowColor, 8, 0, 4)
	h.Add(h.content)
	return h
}

// ElementType returns type identifier.
func (h *HoverCard) ElementType() string { return "HoverCard" }

// Update handles hover detection.
func (h *HoverCard) Update() error {
	eng := h.GetEngine()
	if eng == nil {
		return nil
	}
	hover := eng.GetHoverTarget()
	shouldShow := false
	for target := hover; target != nil; target = target.GetParent() {
		if target == h.trigger || target == h.content {
			shouldShow = true
			break
		}
	}
	if shouldShow != h.visible {
		h.visible = shouldShow
		if shouldShow {
			tb := h.trigger.GetBounds()
			h.content.SetPosition(yoga.EdgeLeft, tb.X)
			h.content.SetPosition(yoga.EdgeTop, tb.Y+tb.Height)
			h.content.SetVisible(true)
			eng.AddOverlay(h.content)
		} else {
			h.content.SetVisible(false)
			eng.RemoveOverlay(h.content)
		}
	}
	return nil
}

// AddContent adds elements to the hover card.
func (h *HoverCard) AddContent(children ...core.Element) *HoverCard {
	for _, child := range children {
		h.content.AppendChild(child)
	}
	return h
}
