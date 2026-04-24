package components

import (
	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// RadioGroup manages a set of Radio buttons with single-selection.
type RadioGroup struct {
	core.BaseElement
	options    []string
	radios     []*Radio
	selected   int
	onChange   func(index int, label string)
	gap        float32
}

// NewRadioGroup creates a radio group with the given options.
func NewRadioGroup(options []string) *RadioGroup {
	rg := &RadioGroup{
		options:  options,
		selected: -1,
		gap:      12,
	}
	rg.Init(rg)
	rg.SetFlexDirection(yoga.FlexDirectionColumn)
	rg.SetGap(yoga.GutterAll, rg.gap)
	rg.buildRadios()
	return rg
}

// ElementType returns type identifier.
func (rg *RadioGroup) ElementType() string { return "RadioGroup" }

func (rg *RadioGroup) buildRadios() {
	rg.ClearChildren()
	rg.radios = make([]*Radio, 0, len(rg.options))
	for i, opt := range rg.options {
		idx := i
		radio := NewRadio(opt)
		radio.SetSelected(idx == rg.selected)
		radio.SetOnChange(func(selected bool) {
			if selected {
				rg.setSelected(idx)
			}
		})
		rg.radios = append(rg.radios, radio)
		rg.AppendChild(radio)
	}
}

func (rg *RadioGroup) setSelected(index int) {
	if index == rg.selected {
		return
	}
	old := rg.selected
	rg.selected = index
	for i, r := range rg.radios {
		r.SetSelected(i == rg.selected)
	}
	if rg.onChange != nil {
		label := ""
		if rg.selected >= 0 && rg.selected < len(rg.options) {
			label = rg.options[rg.selected]
		}
		rg.onChange(rg.selected, label)
	}
	// If unselecting, also trigger change with old index for completeness
	if old >= 0 && old < len(rg.radios) && index < 0 {
		rg.radios[old].SetSelected(false)
	}
}

// SetSelected sets the selected index.
func (rg *RadioGroup) SetSelected(index int) *RadioGroup {
	rg.setSelected(index)
	return rg
}

// GetSelected returns the current selected index.
func (rg *RadioGroup) GetSelected() int {
	return rg.selected
}

// GetSelectedLabel returns the label of the selected option.
func (rg *RadioGroup) GetSelectedLabel() string {
	if rg.selected >= 0 && rg.selected < len(rg.options) {
		return rg.options[rg.selected]
	}
	return ""
}

// SetOnChange sets the selection change callback.
func (rg *RadioGroup) SetOnChange(fn func(index int, label string)) *RadioGroup {
	rg.onChange = fn
	return rg
}

// SetOptions rebuilds the group with new options.
func (rg *RadioGroup) SetOptions(options []string) *RadioGroup {
	rg.options = options
	rg.selected = -1
	rg.buildRadios()
	rg.Mark(core.FlagNeedLayout | core.FlagNeedDraw)
	return rg
}

// SetGap sets the spacing between radio items.
func (rg *RadioGroup) SetGap(gap float32) *RadioGroup {
	rg.gap = gap
	rg.SetGap(yoga.GutterAll, gap)
	return rg
}

// SetDirection sets the layout direction (Row or Column).
func (rg *RadioGroup) SetDirection(dir yoga.FlexDirection) *RadioGroup {
	rg.SetFlexDirection(dir)
	return rg
}
