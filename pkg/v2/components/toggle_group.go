package components

import (
	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// ToggleGroup is a set of toggles where only one can be active (single) or multiple can be active.
type ToggleGroup struct {
	core.BaseElement
	toggles  []*Toggle
	single   bool
	selected map[int]bool
	onChange func(selected []int)
}

// NewToggleGroup creates a toggle group.
func NewToggleGroup(labels []string) *ToggleGroup {
	g := &ToggleGroup{
		single:   true,
		selected: make(map[int]bool),
	}
	g.Init(g)
	g.SetFlexDirection(yoga.FlexDirectionRow)
	g.SetGap(yoga.GutterAll, 4)

	for i, label := range labels {
		idx := i
		t := NewToggle(label)
		t.SetOnChange(func(pressed bool) {
			g.handleToggleChange(idx, pressed)
		})
		g.toggles = append(g.toggles, t)
		g.AppendChild(t)
	}
	return g
}

// ElementType returns type identifier.
func (g *ToggleGroup) ElementType() string { return "ToggleGroup" }

func (g *ToggleGroup) handleToggleChange(idx int, pressed bool) {
	if g.single {
		for i := range g.selected {
			g.selected[i] = false
			g.toggles[i].SetPressed(false)
		}
		if pressed {
			g.selected[idx] = true
			g.toggles[idx].SetPressed(true)
		}
	} else {
		g.selected[idx] = pressed
	}
	if g.onChange != nil {
		var res []int
		for i, v := range g.selected {
			if v {
				res = append(res, i)
			}
		}
		g.onChange(res)
	}
}

// SetSingle controls whether only one toggle can be active.
func (g *ToggleGroup) SetSingle(single bool) *ToggleGroup {
	g.single = single
	return g
}

// Select sets the selected toggle(s) by index.
func (g *ToggleGroup) Select(index int) *ToggleGroup {
	if index >= 0 && index < len(g.toggles) {
		g.handleToggleChange(index, true)
	}
	return g
}

// SetOnChange sets the change callback.
func (g *ToggleGroup) SetOnChange(fn func(selected []int)) *ToggleGroup {
	g.onChange = fn
	return g
}
