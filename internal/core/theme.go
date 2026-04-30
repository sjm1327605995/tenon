package core

import "image/color"

// Theme 定义 Tenon UI 的全局视觉主题。
// 所有组件在初始化时默认从当前主题读取配色，用户仍可通过链式 API 单独覆盖。
type Theme struct {
	// 通用品牌色
	PrimaryColor           color.Color // 主题色：按钮、滑块、开关、选区
	PrimaryHoverColor      color.Color
	PrimaryForegroundColor color.Color // 主题色上的文字

	// 次要色
	SecondaryColor           color.Color
	SecondaryForegroundColor color.Color

	// 强调色
	AccentColor           color.Color
	AccentForegroundColor color.Color

	// 破坏色（错误、删除）
	DestructiveColor           color.Color
	DestructiveForegroundColor color.Color

	// 静音/禁用色
	MutedColor           color.Color
	MutedForegroundColor color.Color

	// 卡片与浮层
	CardColor              color.Color
	CardForegroundColor    color.Color
	PopoverColor           color.Color
	PopoverForegroundColor color.Color

	// 边框
	BorderColor      color.Color
	FocusBorderColor color.Color
	BorderRadius     float32
	RingColor        color.Color // 聚焦环

	// 背景与表面
	BackgroundColor color.Color // 页面/窗口背景
	SurfaceColor    color.Color // 卡片、面板、输入框背景（兼容旧版）

	// 文字
	TextColor        color.Color
	TextMutedColor   color.Color // 次要文本、占位符（兼容旧版）
	PlaceholderColor color.Color

	// 字体大小层级
	FontSizeSM   float32
	FontSizeBase float32
	FontSizeLG   float32
	FontSizeXL   float32

	// 按钮
	ButtonNormalColor   color.Color
	ButtonHoverColor    color.Color
	ButtonPressedColor  color.Color
	ButtonTextColor     color.Color
	ButtonDisabledColor color.Color
	ButtonBorderRadius  float32

	// 输入框
	InputBgColor          color.Color
	InputBorderColor      color.Color
	InputFocusBorderColor color.Color
	InputTextColor        color.Color
	InputPlaceholderColor color.Color
	InputBorderRadius     float32
	InputSelectionColor   color.Color

	// Switch
	SwitchOnColor    color.Color
	SwitchOffColor   color.Color
	SwitchThumbColor color.Color

	// Checkbox
	CheckboxFillColor   color.Color
	CheckboxBorderColor color.Color
	CheckboxCheckColor  color.Color

	// Radio
	RadioFillColor   color.Color
	RadioBorderColor color.Color
	RadioInnerColor  color.Color

	// Slider
	SliderTrackColor color.Color
	SliderFillColor  color.Color
	SliderThumbColor color.Color

	// ProgressBar
	ProgressBarTrackColor color.Color
	ProgressBarFillColor  color.Color

	// ScrollView
	ScrollbarColor      color.Color
	ScrollbarTrackColor color.Color

	// Divider
	DividerColor color.Color

	// 阴影
	ShadowColor color.Color

	// 菜单
	MenuBg               color.Color
	MenuItemSelectedBg   color.Color
	MenuItemSelectedText color.Color
	MenuItemHoverBg      color.Color
}

var defaultTheme *Theme
var themeChangeNotifier func()

// SetTheme 设置全局默认主题，并自动触发所有已挂载 Widget 的重渲染。
func SetTheme(t *Theme) {
	defaultTheme = t
	if themeChangeNotifier != nil {
		themeChangeNotifier()
	}
}

// setThemeChangeNotifier 由 Engine.Mount 调用，注册主题变更通知器。
func setThemeChangeNotifier(fn func()) {
	themeChangeNotifier = fn
}

// GetTheme 返回当前全局主题，若未设置则返回默认浅色主题。
func GetTheme() *Theme {
	if defaultTheme == nil {
		defaultTheme = DefaultLightTheme()
	}
	return defaultTheme
}

// DefaultLightTheme 返回默认浅色主题。
func DefaultLightTheme() *Theme {
	return &Theme{
		PrimaryColor:           color.RGBA{R: 0, G: 123, B: 255, A: 255},
		PrimaryHoverColor:      color.RGBA{R: 0, G: 105, B: 217, A: 255},
		PrimaryForegroundColor: color.White,

		SecondaryColor:           color.RGBA{R: 108, G: 117, B: 125, A: 255},
		SecondaryForegroundColor: color.White,

		AccentColor:           color.RGBA{R: 245, G: 245, B: 245, A: 255},
		AccentForegroundColor: color.RGBA{R: 33, G: 37, B: 41, A: 255},

		DestructiveColor:           color.RGBA{R: 220, G: 53, B: 69, A: 255},
		DestructiveForegroundColor: color.White,

		MutedColor:           color.RGBA{R: 245, G: 245, B: 245, A: 255},
		MutedForegroundColor: color.RGBA{R: 108, G: 117, B: 125, A: 255},

		CardColor:              color.White,
		CardForegroundColor:    color.RGBA{R: 33, G: 37, B: 41, A: 255},
		PopoverColor:           color.White,
		PopoverForegroundColor: color.RGBA{R: 33, G: 37, B: 41, A: 255},

		BorderColor:      color.RGBA{R: 222, G: 226, B: 230, A: 255},
		FocusBorderColor: color.RGBA{R: 0, G: 123, B: 255, A: 255},
		BorderRadius:     8,
		RingColor:        color.RGBA{R: 0, G: 123, B: 255, A: 255},

		BackgroundColor: color.RGBA{R: 245, G: 245, B: 245, A: 255},
		SurfaceColor:    color.White,

		TextColor:        color.RGBA{R: 33, G: 37, B: 41, A: 255},
		TextMutedColor:   color.RGBA{R: 108, G: 117, B: 125, A: 255},
		PlaceholderColor: color.RGBA{R: 170, G: 170, B: 170, A: 255},

		FontSizeSM:   12,
		FontSizeBase: 14,
		FontSizeLG:   18,
		FontSizeXL:   24,

		ButtonNormalColor:   color.RGBA{R: 0, G: 123, B: 255, A: 255},
		ButtonHoverColor:    color.RGBA{R: 70, G: 130, B: 180, A: 255},
		ButtonPressedColor:  color.RGBA{R: 30, G: 144, B: 255, A: 255},
		ButtonTextColor:     color.White,
		ButtonDisabledColor: color.RGBA{R: 108, G: 117, B: 125, A: 255},
		ButtonBorderRadius:  8,

		InputBgColor:          color.White,
		InputBorderColor:      color.RGBA{R: 200, G: 200, B: 200, A: 255},
		InputFocusBorderColor: color.RGBA{R: 0, G: 123, B: 255, A: 255},
		InputTextColor:        color.RGBA{R: 33, G: 37, B: 41, A: 255},
		InputPlaceholderColor: color.RGBA{R: 170, G: 170, B: 170, A: 255},
		InputBorderRadius:     4,
		InputSelectionColor:   color.RGBA{R: 0, G: 123, B: 255, A: 100},

		SwitchOnColor:    color.RGBA{R: 0, G: 123, B: 255, A: 255},
		SwitchOffColor:   color.RGBA{R: 204, G: 204, B: 204, A: 255},
		SwitchThumbColor: color.White,

		CheckboxFillColor:   color.RGBA{R: 0, G: 123, B: 255, A: 255},
		CheckboxBorderColor: color.RGBA{R: 108, G: 117, B: 125, A: 255},
		CheckboxCheckColor:  color.White,

		RadioFillColor:   color.RGBA{R: 0, G: 123, B: 255, A: 255},
		RadioBorderColor: color.RGBA{R: 108, G: 117, B: 125, A: 255},
		RadioInnerColor:  color.White,

		SliderTrackColor: color.RGBA{R: 224, G: 224, B: 224, A: 255},
		SliderFillColor:  color.RGBA{R: 0, G: 123, B: 255, A: 255},
		SliderThumbColor: color.White,

		ProgressBarTrackColor: color.RGBA{R: 224, G: 224, B: 224, A: 255},
		ProgressBarFillColor:  color.RGBA{R: 0, G: 123, B: 255, A: 255},

		ScrollbarColor:      color.RGBA{R: 150, G: 150, B: 150, A: 200},
		ScrollbarTrackColor: color.RGBA{R: 230, G: 230, B: 230, A: 100},

		DividerColor: color.RGBA{R: 222, G: 226, B: 230, A: 255},

		ShadowColor: color.RGBA{A: 80},

		MenuBg:               color.White,
		MenuItemSelectedBg:   color.RGBA{R: 230, G: 244, B: 255, A: 255},
		MenuItemSelectedText: color.RGBA{R: 0, G: 123, B: 255, A: 255},
		MenuItemHoverBg:      color.RGBA{R: 245, G: 245, B: 245, A: 255},
	}
}

// DefaultAntTheme 返回 Ant Design 风格主题。
func DefaultAntTheme() *Theme {
	return &Theme{
		PrimaryColor:           color.RGBA{R: 22, G: 119, B: 255, A: 255}, // #1677ff
		PrimaryHoverColor:      color.RGBA{R: 64, G: 150, B: 255, A: 255},
		PrimaryForegroundColor: color.White,

		SecondaryColor:           color.RGBA{R: 217, G: 217, B: 217, A: 255},
		SecondaryForegroundColor: color.RGBA{R: 0, G: 0, B: 0, A: 224},

		AccentColor:           color.RGBA{R: 245, G: 245, B: 245, A: 255},
		AccentForegroundColor: color.RGBA{R: 0, G: 0, B: 0, A: 224},

		DestructiveColor:           color.RGBA{R: 255, G: 77, B: 79, A: 255},
		DestructiveForegroundColor: color.White,

		MutedColor:           color.RGBA{R: 245, G: 245, B: 245, A: 255},
		MutedForegroundColor: color.RGBA{R: 0, G: 0, B: 0, A: 115},

		CardColor:              color.White,
		CardForegroundColor:    color.RGBA{R: 0, G: 0, B: 0, A: 224},
		PopoverColor:           color.White,
		PopoverForegroundColor: color.RGBA{R: 0, G: 0, B: 0, A: 224},

		BorderColor:      color.RGBA{R: 217, G: 217, B: 217, A: 255}, // #d9d9d9
		FocusBorderColor: color.RGBA{R: 22, G: 119, B: 255, A: 255},
		BorderRadius:     6,
		RingColor:        color.RGBA{R: 22, G: 119, B: 255, A: 255},

		BackgroundColor: color.RGBA{R: 240, G: 242, B: 245, A: 255}, // #f0f2f5
		SurfaceColor:    color.White,

		TextColor:        color.RGBA{R: 0, G: 0, B: 0, A: 224}, // rgba(0,0,0,0.88)
		TextMutedColor:   color.RGBA{R: 0, G: 0, B: 0, A: 115}, // rgba(0,0,0,0.45)
		PlaceholderColor: color.RGBA{R: 0, G: 0, B: 0, A: 64},  // rgba(0,0,0,0.25)

		FontSizeSM:   12,
		FontSizeBase: 14,
		FontSizeLG:   16,
		FontSizeXL:   24,

		ButtonNormalColor:   color.RGBA{R: 22, G: 119, B: 255, A: 255},
		ButtonHoverColor:    color.RGBA{R: 64, G: 150, B: 255, A: 255},
		ButtonPressedColor:  color.RGBA{R: 9, G: 88, B: 217, A: 255},
		ButtonTextColor:     color.White,
		ButtonDisabledColor: color.RGBA{R: 191, G: 191, B: 191, A: 255},
		ButtonBorderRadius:  6,

		InputBgColor:          color.White,
		InputBorderColor:      color.RGBA{R: 217, G: 217, B: 217, A: 255},
		InputFocusBorderColor: color.RGBA{R: 22, G: 119, B: 255, A: 255},
		InputTextColor:        color.RGBA{R: 0, G: 0, B: 0, A: 224},
		InputPlaceholderColor: color.RGBA{R: 0, G: 0, B: 0, A: 64},
		InputBorderRadius:     6,
		InputSelectionColor:   color.RGBA{R: 22, G: 119, B: 255, A: 100},

		SwitchOnColor:    color.RGBA{R: 22, G: 119, B: 255, A: 255},
		SwitchOffColor:   color.RGBA{R: 0, G: 0, B: 0, A: 25},
		SwitchThumbColor: color.White,

		CheckboxFillColor:   color.RGBA{R: 22, G: 119, B: 255, A: 255},
		CheckboxBorderColor: color.RGBA{R: 217, G: 217, B: 217, A: 255},
		CheckboxCheckColor:  color.White,

		RadioFillColor:   color.RGBA{R: 22, G: 119, B: 255, A: 255},
		RadioBorderColor: color.RGBA{R: 217, G: 217, B: 217, A: 255},
		RadioInnerColor:  color.White,

		SliderTrackColor: color.RGBA{R: 245, G: 245, B: 245, A: 255},
		SliderFillColor:  color.RGBA{R: 22, G: 119, B: 255, A: 255},
		SliderThumbColor: color.White,

		ProgressBarTrackColor: color.RGBA{R: 245, G: 245, B: 245, A: 255},
		ProgressBarFillColor:  color.RGBA{R: 22, G: 119, B: 255, A: 255},

		ScrollbarColor:      color.RGBA{R: 0, G: 0, B: 0, A: 30},
		ScrollbarTrackColor: color.RGBA{R: 0, G: 0, B: 0, A: 10},

		DividerColor: color.RGBA{R: 240, G: 240, B: 240, A: 255},

		ShadowColor: color.RGBA{A: 40},

		MenuBg:               color.White,
		MenuItemSelectedBg:   color.RGBA{R: 230, G: 244, B: 255, A: 255}, // #e6f4ff
		MenuItemSelectedText: color.RGBA{R: 22, G: 119, B: 255, A: 255},  // #1677ff
		MenuItemHoverBg:      color.RGBA{R: 245, G: 245, B: 245, A: 255}, // #f5f5f5
	}
}

// DefaultDarkTheme 返回默认深色主题。
func DefaultDarkTheme() *Theme {
	return &Theme{
		PrimaryColor:           color.RGBA{R: 64, G: 169, B: 255, A: 255},
		PrimaryHoverColor:      color.RGBA{R: 89, G: 182, B: 255, A: 255},
		PrimaryForegroundColor: color.RGBA{R: 18, G: 18, B: 18, A: 255},

		SecondaryColor:           color.RGBA{R: 66, G: 66, B: 66, A: 255},
		SecondaryForegroundColor: color.RGBA{R: 230, G: 230, B: 230, A: 255},

		AccentColor:           color.RGBA{R: 45, G: 45, B: 45, A: 255},
		AccentForegroundColor: color.RGBA{R: 230, G: 230, B: 230, A: 255},

		DestructiveColor:           color.RGBA{R: 255, G: 100, B: 100, A: 255},
		DestructiveForegroundColor: color.RGBA{R: 18, G: 18, B: 18, A: 255},

		MutedColor:           color.RGBA{R: 45, G: 45, B: 45, A: 255},
		MutedForegroundColor: color.RGBA{R: 150, G: 150, B: 150, A: 255},

		CardColor:              color.RGBA{R: 30, G: 30, B: 30, A: 255},
		CardForegroundColor:    color.RGBA{R: 230, G: 230, B: 230, A: 255},
		PopoverColor:           color.RGBA{R: 30, G: 30, B: 30, A: 255},
		PopoverForegroundColor: color.RGBA{R: 230, G: 230, B: 230, A: 255},

		BorderColor:      color.RGBA{R: 66, G: 66, B: 66, A: 255},
		FocusBorderColor: color.RGBA{R: 64, G: 169, B: 255, A: 255},
		BorderRadius:     8,
		RingColor:        color.RGBA{R: 64, G: 169, B: 255, A: 255},

		BackgroundColor: color.RGBA{R: 18, G: 18, B: 18, A: 255},
		SurfaceColor:    color.RGBA{R: 30, G: 30, B: 30, A: 255},

		TextColor:        color.RGBA{R: 230, G: 230, B: 230, A: 255},
		TextMutedColor:   color.RGBA{R: 150, G: 150, B: 150, A: 255},
		PlaceholderColor: color.RGBA{R: 100, G: 100, B: 100, A: 255},

		FontSizeSM:   12,
		FontSizeBase: 14,
		FontSizeLG:   18,
		FontSizeXL:   24,

		ButtonNormalColor:   color.RGBA{R: 64, G: 169, B: 255, A: 255},
		ButtonHoverColor:    color.RGBA{R: 89, G: 182, B: 255, A: 255},
		ButtonPressedColor:  color.RGBA{R: 41, G: 151, B: 255, A: 255},
		ButtonTextColor:     color.White,
		ButtonDisabledColor: color.RGBA{R: 80, G: 80, B: 80, A: 255},
		ButtonBorderRadius:  8,

		InputBgColor:          color.RGBA{R: 30, G: 30, B: 30, A: 255},
		InputBorderColor:      color.RGBA{R: 66, G: 66, B: 66, A: 255},
		InputFocusBorderColor: color.RGBA{R: 64, G: 169, B: 255, A: 255},
		InputTextColor:        color.RGBA{R: 230, G: 230, B: 230, A: 255},
		InputPlaceholderColor: color.RGBA{R: 100, G: 100, B: 100, A: 255},
		InputBorderRadius:     4,
		InputSelectionColor:   color.RGBA{R: 64, G: 169, B: 255, A: 100},

		SwitchOnColor:    color.RGBA{R: 64, G: 169, B: 255, A: 255},
		SwitchOffColor:   color.RGBA{R: 80, G: 80, B: 80, A: 255},
		SwitchThumbColor: color.RGBA{R: 200, G: 200, B: 200, A: 255},

		CheckboxFillColor:   color.RGBA{R: 64, G: 169, B: 255, A: 255},
		CheckboxBorderColor: color.RGBA{R: 150, G: 150, B: 150, A: 255},
		CheckboxCheckColor:  color.RGBA{R: 30, G: 30, B: 30, A: 255},

		RadioFillColor:   color.RGBA{R: 64, G: 169, B: 255, A: 255},
		RadioBorderColor: color.RGBA{R: 150, G: 150, B: 150, A: 255},
		RadioInnerColor:  color.RGBA{R: 30, G: 30, B: 30, A: 255},

		SliderTrackColor: color.RGBA{R: 66, G: 66, B: 66, A: 255},
		SliderFillColor:  color.RGBA{R: 64, G: 169, B: 255, A: 255},
		SliderThumbColor: color.RGBA{R: 200, G: 200, B: 200, A: 255},

		ProgressBarTrackColor: color.RGBA{R: 66, G: 66, B: 66, A: 255},
		ProgressBarFillColor:  color.RGBA{R: 64, G: 169, B: 255, A: 255},

		ScrollbarColor:      color.RGBA{R: 120, G: 120, B: 120, A: 200},
		ScrollbarTrackColor: color.RGBA{R: 50, G: 50, B: 50, A: 100},

		DividerColor: color.RGBA{R: 66, G: 66, B: 66, A: 255},

		ShadowColor: color.RGBA{A: 60},

		MenuBg:               color.RGBA{R: 0, G: 21, B: 41, A: 255},
		MenuItemSelectedBg:   color.RGBA{R: 0, G: 123, B: 255, A: 255},
		MenuItemSelectedText: color.White,
		MenuItemHoverBg:      color.RGBA{R: 255, G: 255, B: 255, A: 30},
	}
}

// DefaultShadcnLightTheme 返回 Shadcn/UI Neutral 风格的浅色主题。
// 参考 https://ui.shadcn.com/docs/theming
func DefaultShadcnLightTheme() *Theme {
	return &Theme{
		PrimaryColor:           color.RGBA{R: 23, G: 23, B: 23, A: 255},    // #171717
		PrimaryHoverColor:      color.RGBA{R: 50, G: 50, B: 50, A: 255},    // #323232
		PrimaryForegroundColor: color.RGBA{R: 250, G: 250, B: 250, A: 255}, // #fafafa

		SecondaryColor:           color.RGBA{R: 245, G: 245, B: 245, A: 255}, // #f5f5f5
		SecondaryForegroundColor: color.RGBA{R: 23, G: 23, B: 23, A: 255},    // #171717

		AccentColor:           color.RGBA{R: 245, G: 245, B: 245, A: 255}, // #f5f5f5
		AccentForegroundColor: color.RGBA{R: 23, G: 23, B: 23, A: 255},    // #171717

		DestructiveColor:           color.RGBA{R: 239, G: 68, B: 68, A: 255}, // #ef4444
		DestructiveForegroundColor: color.RGBA{R: 250, G: 250, B: 250, A: 255},

		MutedColor:           color.RGBA{R: 245, G: 245, B: 245, A: 255}, // #f5f5f5
		MutedForegroundColor: color.RGBA{R: 115, G: 115, B: 115, A: 255}, // #737373

		CardColor:              color.RGBA{R: 255, G: 255, B: 255, A: 255}, // #ffffff
		CardForegroundColor:    color.RGBA{R: 10, G: 10, B: 10, A: 255},    // #0a0a0a
		PopoverColor:           color.RGBA{R: 255, G: 255, B: 255, A: 255}, // #ffffff
		PopoverForegroundColor: color.RGBA{R: 10, G: 10, B: 10, A: 255},    // #0a0a0a

		BorderColor:      color.RGBA{R: 229, G: 229, B: 229, A: 255}, // #e5e5e5
		FocusBorderColor: color.RGBA{R: 23, G: 23, B: 23, A: 255},    // #171717
		BorderRadius:     10,                                         // 0.625rem
		RingColor:        color.RGBA{R: 161, G: 161, B: 161, A: 255}, // #a1a1a1

		BackgroundColor: color.RGBA{R: 255, G: 255, B: 255, A: 255}, // #ffffff
		SurfaceColor:    color.RGBA{R: 255, G: 255, B: 255, A: 255}, // #ffffff

		TextColor:        color.RGBA{R: 10, G: 10, B: 10, A: 255},    // #0a0a0a
		TextMutedColor:   color.RGBA{R: 115, G: 115, B: 115, A: 255}, // #737373
		PlaceholderColor: color.RGBA{R: 161, G: 161, B: 161, A: 255}, // #a1a1a1

		FontSizeSM:   12,
		FontSizeBase: 14,
		FontSizeLG:   16,
		FontSizeXL:   24,

		ButtonNormalColor:   color.RGBA{R: 23, G: 23, B: 23, A: 255},    // #171717
		ButtonHoverColor:    color.RGBA{R: 50, G: 50, B: 50, A: 255},    // #323232
		ButtonPressedColor:  color.RGBA{R: 10, G: 10, B: 10, A: 255},    // #0a0a0a
		ButtonTextColor:     color.RGBA{R: 250, G: 250, B: 250, A: 255}, // #fafafa
		ButtonDisabledColor: color.RGBA{R: 161, G: 161, B: 161, A: 255}, // #a1a1a1
		ButtonBorderRadius:  10,

		InputBgColor:          color.RGBA{R: 255, G: 255, B: 255, A: 255},
		InputBorderColor:      color.RGBA{R: 229, G: 229, B: 229, A: 255}, // #e5e5e5
		InputFocusBorderColor: color.RGBA{R: 23, G: 23, B: 23, A: 255},    // #171717
		InputTextColor:        color.RGBA{R: 10, G: 10, B: 10, A: 255},
		InputPlaceholderColor: color.RGBA{R: 161, G: 161, B: 161, A: 255},
		InputBorderRadius:     10,
		InputSelectionColor:   color.RGBA{R: 23, G: 23, B: 23, A: 100},

		SwitchOnColor:    color.RGBA{R: 23, G: 23, B: 23, A: 255},    // #171717
		SwitchOffColor:   color.RGBA{R: 229, G: 229, B: 229, A: 255}, // #e5e5e5
		SwitchThumbColor: color.White,

		CheckboxFillColor:   color.RGBA{R: 23, G: 23, B: 23, A: 255},    // #171717
		CheckboxBorderColor: color.RGBA{R: 229, G: 229, B: 229, A: 255}, // #e5e5e5
		CheckboxCheckColor:  color.RGBA{R: 250, G: 250, B: 250, A: 255}, // #fafafa

		RadioFillColor:   color.RGBA{R: 23, G: 23, B: 23, A: 255},    // #171717
		RadioBorderColor: color.RGBA{R: 229, G: 229, B: 229, A: 255}, // #e5e5e5
		RadioInnerColor:  color.RGBA{R: 250, G: 250, B: 250, A: 255}, // #fafafa

		SliderTrackColor: color.RGBA{R: 229, G: 229, B: 229, A: 255}, // #e5e5e5
		SliderFillColor:  color.RGBA{R: 23, G: 23, B: 23, A: 255},    // #171717
		SliderThumbColor: color.White,

		ProgressBarTrackColor: color.RGBA{R: 229, G: 229, B: 229, A: 255},
		ProgressBarFillColor:  color.RGBA{R: 23, G: 23, B: 23, A: 255},

		ScrollbarColor:      color.RGBA{R: 0, G: 0, B: 0, A: 30},
		ScrollbarTrackColor: color.RGBA{R: 0, G: 0, B: 0, A: 10},

		DividerColor: color.RGBA{R: 229, G: 229, B: 229, A: 255}, // #e5e5e5

		ShadowColor: color.RGBA{A: 30},

		MenuBg:               color.White,
		MenuItemSelectedBg:   color.RGBA{R: 245, G: 245, B: 245, A: 255},
		MenuItemSelectedText: color.RGBA{R: 23, G: 23, B: 23, A: 255},
		MenuItemHoverBg:      color.RGBA{R: 245, G: 245, B: 245, A: 255},
	}
}

// DefaultShadcnDarkTheme 返回 Shadcn/UI Neutral 风格的深色主题。
func DefaultShadcnDarkTheme() *Theme {
	return &Theme{
		PrimaryColor:           color.RGBA{R: 229, G: 229, B: 229, A: 255}, // #e5e5e5
		PrimaryHoverColor:      color.RGBA{R: 255, G: 255, B: 255, A: 255}, // #ffffff
		PrimaryForegroundColor: color.RGBA{R: 23, G: 23, B: 23, A: 255},    // #171717

		SecondaryColor:           color.RGBA{R: 38, G: 38, B: 38, A: 255}, // #262626
		SecondaryForegroundColor: color.RGBA{R: 229, G: 229, B: 229, A: 255},

		AccentColor:           color.RGBA{R: 38, G: 38, B: 38, A: 255}, // #262626
		AccentForegroundColor: color.RGBA{R: 229, G: 229, B: 229, A: 255},

		DestructiveColor:           color.RGBA{R: 248, G: 113, B: 113, A: 255}, // #f87171
		DestructiveForegroundColor: color.RGBA{R: 23, G: 23, B: 23, A: 255},

		MutedColor:           color.RGBA{R: 38, G: 38, B: 38, A: 255},    // #262626
		MutedForegroundColor: color.RGBA{R: 161, G: 161, B: 161, A: 255}, // #a1a1a1

		CardColor:              color.RGBA{R: 23, G: 23, B: 23, A: 255},    // #171717
		CardForegroundColor:    color.RGBA{R: 229, G: 229, B: 229, A: 255}, // #e5e5e5
		PopoverColor:           color.RGBA{R: 23, G: 23, B: 23, A: 255},    // #171717
		PopoverForegroundColor: color.RGBA{R: 229, G: 229, B: 229, A: 255}, // #e5e5e5

		BorderColor:      color.RGBA{R: 255, G: 255, B: 255, A: 25}, // rgba(255,255,255,0.1)
		FocusBorderColor: color.RGBA{R: 229, G: 229, B: 229, A: 255},
		BorderRadius:     10,
		RingColor:        color.RGBA{R: 115, G: 115, B: 115, A: 255}, // #737373

		BackgroundColor: color.RGBA{R: 10, G: 10, B: 10, A: 255}, // #0a0a0a
		SurfaceColor:    color.RGBA{R: 23, G: 23, B: 23, A: 255}, // #171717

		TextColor:        color.RGBA{R: 229, G: 229, B: 229, A: 255}, // #e5e5e5
		TextMutedColor:   color.RGBA{R: 161, G: 161, B: 161, A: 255}, // #a1a1a1
		PlaceholderColor: color.RGBA{R: 115, G: 115, B: 115, A: 255}, // #737373

		FontSizeSM:   12,
		FontSizeBase: 14,
		FontSizeLG:   16,
		FontSizeXL:   24,

		ButtonNormalColor:   color.RGBA{R: 229, G: 229, B: 229, A: 255}, // #e5e5e5
		ButtonHoverColor:    color.RGBA{R: 255, G: 255, B: 255, A: 255}, // #ffffff
		ButtonPressedColor:  color.RGBA{R: 200, G: 200, B: 200, A: 255}, // #c8c8c8
		ButtonTextColor:     color.RGBA{R: 23, G: 23, B: 23, A: 255},    // #171717
		ButtonDisabledColor: color.RGBA{R: 66, G: 66, B: 66, A: 255},    // #424242
		ButtonBorderRadius:  10,

		InputBgColor:          color.RGBA{R: 10, G: 10, B: 10, A: 255},
		InputBorderColor:      color.RGBA{R: 255, G: 255, B: 255, A: 25},
		InputFocusBorderColor: color.RGBA{R: 229, G: 229, B: 229, A: 255},
		InputTextColor:        color.RGBA{R: 229, G: 229, B: 229, A: 255},
		InputPlaceholderColor: color.RGBA{R: 115, G: 115, B: 115, A: 255},
		InputBorderRadius:     10,
		InputSelectionColor:   color.RGBA{R: 229, G: 229, B: 229, A: 100},

		SwitchOnColor:    color.RGBA{R: 229, G: 229, B: 229, A: 255},
		SwitchOffColor:   color.RGBA{R: 38, G: 38, B: 38, A: 255}, // #262626
		SwitchThumbColor: color.RGBA{R: 66, G: 66, B: 66, A: 255}, // #424242

		CheckboxFillColor:   color.RGBA{R: 229, G: 229, B: 229, A: 255},
		CheckboxBorderColor: color.RGBA{R: 66, G: 66, B: 66, A: 255},
		CheckboxCheckColor:  color.RGBA{R: 23, G: 23, B: 23, A: 255},

		RadioFillColor:   color.RGBA{R: 229, G: 229, B: 229, A: 255},
		RadioBorderColor: color.RGBA{R: 66, G: 66, B: 66, A: 255},
		RadioInnerColor:  color.RGBA{R: 23, G: 23, B: 23, A: 255},

		SliderTrackColor: color.RGBA{R: 38, G: 38, B: 38, A: 255},
		SliderFillColor:  color.RGBA{R: 229, G: 229, B: 229, A: 255},
		SliderThumbColor: color.RGBA{R: 66, G: 66, B: 66, A: 255},

		ProgressBarTrackColor: color.RGBA{R: 38, G: 38, B: 38, A: 255},
		ProgressBarFillColor:  color.RGBA{R: 229, G: 229, B: 229, A: 255},

		ScrollbarColor:      color.RGBA{R: 255, G: 255, B: 255, A: 30},
		ScrollbarTrackColor: color.RGBA{R: 255, G: 255, B: 255, A: 10},

		DividerColor: color.RGBA{R: 255, G: 255, B: 255, A: 25},

		ShadowColor: color.RGBA{A: 40},

		MenuBg:               color.RGBA{R: 10, G: 10, B: 10, A: 255},
		MenuItemSelectedBg:   color.RGBA{R: 38, G: 38, B: 38, A: 255},
		MenuItemSelectedText: color.RGBA{R: 229, G: 229, B: 229, A: 255},
		MenuItemHoverBg:      color.RGBA{R: 38, G: 38, B: 38, A: 255},
	}
}
