package components

import (
	"github.com/sjm1327605995/tenon/internal/core"
	"github.com/sjm1327605995/tenon/internal/native"
	"github.com/sjm1327605995/tenon/yoga"
)

type Popover struct {
	core.BaseElement
	trigger *Button
	Content *native.View
	isOpen  bool
}

func NewPopover(triggerLabel string) *Popover {
	p := &Popover{}
	p.Init(p)
	p.SetFlexDirection(yoga.FlexDirectionColumn)

	p.trigger = NewButton(triggerLabel).SetVariant(ButtonOutline)
	p.trigger.SetOnClick(func() { p.toggle() })
	p.AppendChild(p.trigger)

	p.Content = native.NewView()
	p.Content.SetFlexDirection(yoga.FlexDirectionColumn)
	p.Content.SetBackgroundColor(core.GetTheme().CardColor)
	p.Content.SetBorderColor(core.GetTheme().BorderColor)
	p.Content.SetBorderRadius(core.GetTheme().BorderRadius)
	p.Content.SetPadding(yoga.EdgeAll, 16)
	p.Content.SetShadow(core.GetTheme().ShadowColor, 8, 0, 4)
	p.Content.SetPositionType(yoga.PositionTypeAbsolute)
	p.Content.SetVisible(false)
	p.Content.SetMinWidth(200)
	return p
}

func (p *Popover) ElementType() string { return "Popover" }

func (p *Popover) toggle() {
	if p.isOpen {
		p.Close()
	} else {
		p.Open()
	}
}

func (p *Popover) Open() *Popover {
	p.isOpen = true
	tb := p.trigger.GetBounds()
	p.Content.SetPosition(yoga.EdgeLeft, tb.X)
	p.Content.SetPosition(yoga.EdgeTop, tb.Y+tb.Height)
	p.Content.SetVisible(true)
	if eng := p.GetEngine(); eng != nil {
		eng.AddOverlay(p.Content)
	}
	return p
}

func (p *Popover) Close() *Popover {
	p.isOpen = false
	p.Content.SetVisible(false)
	if eng := p.GetEngine(); eng != nil {
		eng.RemoveOverlay(p.Content)
	}
	return p
}

func (p *Popover) SetTrigger(text string) *Popover {
	p.trigger.SetText(text)
	return p
}

func (p *Popover) AddContent(children ...core.Element) *Popover {
	for _, child := range children {
		p.Content.AppendChild(child)
	}
	return p
}

func (p *Popover) HandleEvent(e *core.Event) bool {
	if e.Type == core.EventClick && p.isOpen {
		target := e.Target
		for target != nil {
			if target == p || target == p.Content {
				return false
			}
			target = target.GetParent()
		}
		p.Close()
		return false
	}
	return false
}
