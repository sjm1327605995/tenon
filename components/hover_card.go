package components

import (
	"github.com/sjm1327605995/tenon/internal/core"
	"github.com/sjm1327605995/tenon/internal/native"
	"github.com/sjm1327605995/tenon/yoga"
)

// HoverCard is a floating card shown on hover.
type HoverCard struct {
	core.BaseElement
	trigger *native.View
	Content *native.View
	visible bool
}

// NewHoverCard creates a hover card.
func NewHoverCard(trigger core.Element) *HoverCard {
	h := &HoverCard{}
	h.Init(h)
	h.SetFlexDirection(yoga.FlexDirectionColumn)

	if t, ok := trigger.(*native.View); ok {
		h.trigger = t
	} else {
		h.trigger = native.NewView()
		h.trigger.Add(trigger)
	}
	h.Add(h.trigger)

	h.Content = native.NewView()
	h.Content.SetFlexDirection(yoga.FlexDirectionColumn)
	h.Content.SetBackgroundColor(core.GetTheme().CardColor)
	h.Content.SetBorderColor(core.GetTheme().BorderColor)
	h.Content.SetBorderRadius(core.GetTheme().BorderRadius)
	h.Content.SetPadding(yoga.EdgeAll, 16)
	h.Content.SetPositionType(yoga.PositionTypeAbsolute)
	h.Content.SetVisible(false)
	h.Content.SetShadow(core.GetTheme().ShadowColor, 8, 0, 4)
	h.Add(h.Content)
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
		if target == h.trigger || target == h.Content {
			shouldShow = true
			break
		}
	}
	if shouldShow != h.visible {
		h.visible = shouldShow
		if shouldShow {
			tb := h.trigger.GetBounds()
			h.Content.SetPosition(yoga.EdgeLeft, tb.X)
			h.Content.SetPosition(yoga.EdgeTop, tb.Y+tb.Height)
			h.Content.SetVisible(true)
			eng.AddOverlay(h.Content)
		} else {
			h.Content.SetVisible(false)
			eng.RemoveOverlay(h.Content)
		}
	}
	return nil
}

// AddContent adds elements to the hover card.
func (h *HoverCard) AddContent(children ...core.Element) *HoverCard {
	for _, child := range children {
		h.Content.AppendChild(child)
	}
	return h
}
