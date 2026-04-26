package components

import (
	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

type Popover struct {
	core.BaseElement
	trigger *Button
	content *View
	isOpen  bool
}

func NewPopover(triggerLabel string) *Popover {
	p := &Popover{}
	p.Init(p)
	p.SetFlexDirection(yoga.FlexDirectionColumn)

	p.trigger = NewButton(triggerLabel).SetVariant(ButtonOutline)
	p.trigger.SetOnClick(func() { p.toggle() })
	p.AppendChild(p.trigger)

	p.content = NewView()
	p.content.SetFlexDirection(yoga.FlexDirectionColumn)
	p.content.SetBackgroundColor(core.GetTheme().CardColor)
	p.content.SetBorderColor(core.GetTheme().BorderColor)
	p.content.SetBorderRadius(core.GetTheme().BorderRadius)
	p.content.SetPadding(yoga.EdgeAll, 16)
	p.content.SetShadow(core.GetTheme().ShadowColor, 8, 0, 4)
	p.content.SetPositionType(yoga.PositionTypeAbsolute)
	p.content.SetVisible(false)
	p.content.SetMinWidth(200)
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
	p.content.SetPosition(yoga.EdgeLeft, tb.X)
	p.content.SetPosition(yoga.EdgeTop, tb.Y+tb.Height)
	p.content.SetVisible(true)
	if eng := p.GetEngine(); eng != nil {
		eng.AddOverlay(p.content)
	}
	return p
}

func (p *Popover) Close() *Popover {
	p.isOpen = false
	p.content.SetVisible(false)
	if eng := p.GetEngine(); eng != nil {
		eng.RemoveOverlay(p.content)
	}
	return p
}

func (p *Popover) SetTrigger(text string) *Popover {
	p.trigger.SetText(text)
	return p
}

func (p *Popover) AddContent(children ...core.Element) *Popover {
	for _, child := range children {
		p.content.AppendChild(child)
	}
	return p
}

func (p *Popover) HandleEvent(e *core.Event) bool {
	if e.Type == core.EventClick && p.isOpen {
		target := e.Target
		for target != nil {
			if target == p || target == p.content {
				return false
			}
			target = target.GetParent()
		}
		p.Close()
		return false
	}
	return false
}
