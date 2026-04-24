package antdesign

import (
	"image/color"

	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/components"
	"github.com/sjm1327605995/tenon/yoga"
)

// AntInput 是 Ant Design 风格的输入框 Widget。
// 内部复用 *components.TextInput 实例以保留光标/选区状态。
type AntInput struct {
	tenon.BaseWidget

	input       *components.TextInput
	placeholder string
	value       string // 缓存当前文本（TextInput 无 getter）
	onChange    func(string)
	onSubmit    func(string)
	width       float32
	prefix      string
	suffix      string
	password    bool
	search      bool
	onSearch    func(string)
	disabled    bool
}

// NewAntInput 创建 Ant Design 输入框。
func NewAntInput() *AntInput {
	ai := &AntInput{}
	ai.Init(ai)
	return ai
}

// Render 返回包裹在 AntD 风格容器中的 TextInput。
func (ai *AntInput) Render() tenon.Component {
	theme := NewAntTheme()

	// 初始化或复用 TextInput（保留输入状态）
	if ai.input == nil {
		ai.input = components.NewTextInput()
		ai.input.SetBackgroundColor(color.RGBA{A: 0})  // 透明，容器负责背景
		ai.input.SetBorderColor(color.RGBA{A: 0})
		ai.input.SetFocusBorderColor(color.RGBA{A: 0})
		ai.input.SetPadding(0)
	}
	ai.input.SetPlaceholder(ai.placeholder)
	ai.input.SetOnChange(func(v string) {
		ai.value = v
		if ai.onChange != nil {
			ai.onChange(v)
		}
	})
	ai.input.SetOnSubmit(func(v string) {
		if ai.onSubmit != nil {
			ai.onSubmit(v)
		}
		if ai.search && ai.onSearch != nil {
			ai.onSearch(v)
		}
	})
	if ai.width > 0 {
		ai.input.SetWidth(ai.width - ai.calcExtraWidth())
	}

	// 容器 View：负责边框、背景、圆角
	root := components.NewView().
		SetFlexDirection(yoga.FlexDirectionRow).
		SetAlignItems(yoga.AlignCenter).
		SetBackgroundColor(theme.SurfaceColor).
		SetBorderRadius(theme.InputBorderRadius).
		SetBorder(yoga.EdgeAll, 1).
		SetBorderColor(theme.InputBorderColor).
		SetPadding(yoga.EdgeHorizontal, 11).
		SetPadding(yoga.EdgeVertical, 4)
	if ai.width > 0 {
		root.SetWidth(ai.width)
	}

	// 前缀
	if ai.prefix != "" {
		root.Add(components.NewText(ai.prefix).
			SetColor(theme.TextMutedColor).
			SetFontSize(theme.FontSizeBase).
			SetMargin(yoga.EdgeRight, 8))
	}

	root.AddChild(ai.input)

	// 后缀 / 搜索按钮
	if ai.suffix != "" {
		root.Add(components.NewText(ai.suffix).
			SetColor(theme.TextMutedColor).
			SetFontSize(theme.FontSizeBase).
			SetMargin(yoga.EdgeLeft, 8))
	}
	if ai.search {
		root.Add(components.NewButton("🔍").
			SetWidth(28).SetHeight(28).
			SetMargin(yoga.EdgeLeft, 8).
			SetOnClick(func() {
				if ai.onSearch != nil {
					ai.onSearch(ai.value)
				}
			}))
	}

	return root
}

// calcExtraWidth 计算前缀/后缀/按钮占用的宽度。
func (ai *AntInput) calcExtraWidth() float32 {
	extra := float32(22 + 8 + 8) // padding horizontal + 左右间隙
	if ai.prefix != "" {
		extra += float32(len(ai.prefix))*8 + 8
	}
	if ai.suffix != "" {
		extra += float32(len(ai.suffix))*8 + 8
	}
	if ai.search {
		extra += 36 // 按钮 28 + margin 8
	}
	return extra
}

// ==================== 链式 API ====================

func (ai *AntInput) SetPlaceholder(v string) *AntInput {
	ai.placeholder = v
	return ai
}
func (ai *AntInput) SetOnChange(fn func(string)) *AntInput {
	ai.onChange = fn
	return ai
}
func (ai *AntInput) SetOnSubmit(fn func(string)) *AntInput {
	ai.onSubmit = fn
	return ai
}
func (ai *AntInput) SetWidth(w float32) *AntInput {
	ai.width = w
	return ai
}
func (ai *AntInput) SetPrefix(v string) *AntInput {
	ai.prefix = v
	return ai
}
func (ai *AntInput) SetSuffix(v string) *AntInput {
	ai.suffix = v
	return ai
}
func (ai *AntInput) SetPassword(v bool) *AntInput {
	ai.password = v
	return ai
}
func (ai *AntInput) SetSearch(v bool) *AntInput {
	ai.search = v
	return ai
}
func (ai *AntInput) SetOnSearch(fn func(string)) *AntInput {
	ai.onSearch = fn
	return ai
}
func (ai *AntInput) SetDisabled(v bool) *AntInput {
	ai.disabled = v
	return ai
}
