package components

import (
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// RadioOption 定义单选组中的一个选项。
type RadioOption struct {
	Label    string
	Value    string
	Disabled bool
}

// RadioGroup 是单选按钮组组件，管理多个 Radio 的选中状态。
type RadioGroup struct {
	core.BaseHost
	options    []RadioOption
	value      string
	onChange   func(value string)
	radios     []*Radio
	innerView  *View
}

// NewRadioGroup 创建一个单选按钮组。
func NewRadioGroup() *RadioGroup {
	rg := &RadioGroup{}
	rg.Init(rg)
	rg.GetElement().Yoga.StyleSetFlexDirection(yoga.FlexDirectionColumn)
	rg.GetElement().Yoga.StyleSetGap(yoga.GutterRow, 8)
	return rg
}

// SetOptions 设置选项数据并重建子组件。
func (rg *RadioGroup) SetOptions(options []RadioOption) *RadioGroup {
	rg.options = options
	rg.rebuild()
	return rg
}

// SetValue 设置当前选中的值。
func (rg *RadioGroup) SetValue(value string) *RadioGroup {
	if rg.value != value {
		rg.value = value
		rg.refreshSelection()
	}
	return rg
}

// SetOnChange 设置选中变化回调。
func (rg *RadioGroup) SetOnChange(fn func(value string)) *RadioGroup {
	rg.onChange = fn
	return rg
}

func (rg *RadioGroup) rebuild() {
	rg.ClearChildren()
	rg.radios = nil

	for _, opt := range rg.options {
		radio := NewRadio(opt.Label)
		radio.SetSelected(opt.Value == rg.value)
		if opt.Disabled {
			radio.SetOnChange(nil)
		} else {
			val := opt.Value
			radio.SetOnChange(func(selected bool) {
				if selected {
					rg.SetValue(val)
					if rg.onChange != nil {
						rg.onChange(val)
					}
				}
			})
		}
		rg.radios = append(rg.radios, radio)
		rg.AddChild(radio)
	}
}

func (rg *RadioGroup) refreshSelection() {
	for i, opt := range rg.options {
		if i < len(rg.radios) {
			rg.radios[i].SetSelected(opt.Value == rg.value)
		}
	}
}

// SyncFrom 同步单选组属性。
func (rg *RadioGroup) SyncFrom(other core.Host) {
	if o, ok := other.(*RadioGroup); ok {
		rg.options = o.options
		rg.value = o.value
		rg.onChange = o.onChange
		rg.rebuild()
	}
}
