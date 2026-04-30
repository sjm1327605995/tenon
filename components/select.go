package components

import (
	"github.com/sjm1327605995/tenon/internal/core"
	"github.com/sjm1327605995/tenon/internal/native"
	"github.com/sjm1327605995/tenon/yoga"
)

type Select struct {
	core.BaseElement
	trigger     *Button
	panel       *native.View
	options     []string
	selectedIdx int
	placeholder string
	isOpen      bool
	onChange    func(value string)
}

func NewSelect() *Select {
	s := &Select{selectedIdx: -1, placeholder: "Select..."}
	s.Init(s)
	s.SetFlexDirection(yoga.FlexDirectionColumn)

	s.trigger = NewButton(s.placeholder).SetVariant(ButtonOutline)
	s.trigger.SetOnClick(func() { s.toggle() })
	s.trigger.SetJustifyContent(yoga.JustifySpaceBetween)
	s.AppendChild(s.trigger)

	s.panel = native.NewView()
	s.panel.SetFlexDirection(yoga.FlexDirectionColumn)
	s.panel.SetBackgroundColor(core.GetTheme().CardColor)
	s.panel.SetBorderColor(core.GetTheme().BorderColor)
	s.panel.SetBorderRadius(core.GetTheme().BorderRadius)
	s.panel.SetShadow(core.GetTheme().ShadowColor, 8, 0, 4)
	s.panel.SetPositionType(yoga.PositionTypeAbsolute)
	s.panel.SetVisible(false)
	s.panel.SetMinWidth(150)
	s.panel.SetMaxHeight(200)
	return s
}

func (s *Select) ElementType() string { return "Select" }

func (s *Select) toggle() {
	if s.isOpen {
		s.close()
	} else {
		s.open()
	}
}

func (s *Select) open() {
	s.isOpen = true
	tb := s.trigger.GetBounds()
	s.panel.SetPosition(yoga.EdgeLeft, tb.X)
	s.panel.SetPosition(yoga.EdgeTop, tb.Y+tb.Height)
	s.panel.SetVisible(true)
	if eng := s.GetEngine(); eng != nil {
		eng.AddOverlay(s.panel)
	}
}

func (s *Select) close() {
	s.isOpen = false
	s.panel.SetVisible(false)
	if eng := s.GetEngine(); eng != nil {
		eng.RemoveOverlay(s.panel)
	}
}

func (s *Select) SetOptions(options []string) *Select {
	s.options = append([]string{}, options...)
	s.panel.ClearChildren()
	for i, opt := range options {
		idx := i
		btn := NewButton(opt).SetVariant(ButtonGhost)
		btn.SetJustifyContent(yoga.JustifyFlexStart)
		btn.SetWidthPercent(100)
		btn.SetOnClick(func() {
			s.selectedIdx = idx
			s.trigger.SetText(s.options[idx])
			s.close()
			if s.onChange != nil {
				s.onChange(s.options[idx])
			}
		})
		s.panel.Add(btn)
	}
	return s
}

func (s *Select) SetValue(value string) *Select {
	for i, opt := range s.options {
		if opt == value {
			s.selectedIdx = i
			s.trigger.SetText(value)
			return s
		}
	}
	return s
}

func (s *Select) SetPlaceholder(text string) *Select {
	s.placeholder = text
	if s.selectedIdx < 0 {
		s.trigger.SetText(text)
	}
	return s
}

func (s *Select) Value() string {
	if s.selectedIdx < 0 || s.selectedIdx >= len(s.options) {
		return ""
	}
	return s.options[s.selectedIdx]
}

func (s *Select) SetOnChange(fn func(value string)) *Select {
	s.onChange = fn
	return s
}

func (s *Select) HandleEvent(e *core.Event) bool {
	if e.Type == core.EventClick && s.isOpen {
		target := e.Target
		for target != nil {
			if target == s || target == s.panel {
				return false
			}
			target = target.GetParent()
		}
		s.close()
		return false
	}
	return false
}
