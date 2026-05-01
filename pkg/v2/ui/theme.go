package ui

import (
	"image/color"
)

// Theme 定义 Tenon UI 的全局视觉主题。
// 默认采用 Shadcn/UI Neutral 风格。
type Theme struct {
	PrimaryColor           color.Color
	PrimaryHoverColor      color.Color
	PrimaryForegroundColor color.Color

	SecondaryColor           color.Color
	SecondaryForegroundColor color.Color

	AccentColor           color.Color
	AccentForegroundColor color.Color

	DestructiveColor           color.Color
	DestructiveForegroundColor color.Color

	MutedColor           color.Color
	MutedForegroundColor color.Color

	CardColor           color.Color
	CardForegroundColor color.Color
	PopoverColor        color.Color
	PopoverForegroundColor color.Color

	BorderColor      color.Color
	FocusBorderColor color.Color
	BorderRadius     float32
	RingColor        color.Color

	BackgroundColor color.Color
	SurfaceColor    color.Color

	TextColor        color.Color
	TextMutedColor   color.Color
	PlaceholderColor color.Color

	FontSizeSM   float32
	FontSizeBase float32
	FontSizeLG   float32
	FontSizeXL   float32

	ButtonNormalColor   color.Color
	ButtonHoverColor    color.Color
	ButtonPressedColor  color.Color
	ButtonTextColor     color.Color
	ButtonDisabledColor color.Color
	ButtonBorderRadius  float32

	InputBgColor          color.Color
	InputBorderColor      color.Color
	InputFocusBorderColor color.Color
	InputTextColor        color.Color
	InputPlaceholderColor color.Color
	InputBorderRadius     float32
	InputSelectionColor   color.Color

	SwitchOnColor    color.Color
	SwitchOffColor   color.Color
	SwitchThumbColor color.Color

	CheckboxFillColor   color.Color
	CheckboxBorderColor color.Color
	CheckboxCheckColor  color.Color

	RadioFillColor   color.Color
	RadioBorderColor color.Color
	RadioInnerColor  color.Color

	SliderTrackColor color.Color
	SliderFillColor  color.Color
	SliderThumbColor color.Color

	ProgressBarTrackColor color.Color
	ProgressBarFillColor  color.Color

	ScrollbarColor      color.Color
	ScrollbarTrackColor color.Color

	DividerColor color.Color

	ShadowColor color.Color

	MenuBg               color.Color
	MenuItemSelectedBg   color.Color
	MenuItemSelectedText color.Color
	MenuItemHoverBg      color.Color
}

var currentTheme *Theme

func SetTheme(t *Theme) {
	currentTheme = t
}

func GetTheme() *Theme {
	if currentTheme == nil {
		currentTheme = DefaultLightTheme()
	}
	return currentTheme
}

// DefaultLightTheme 返回 Shadcn/UI Neutral 风格的浅色主题。
func DefaultLightTheme() *Theme {
	return &Theme{
		PrimaryColor:           color.RGBA{R: 23, G: 23, B: 23, A: 255},    // #171717
		PrimaryHoverColor:      color.RGBA{R: 50, G: 50, B: 50, A: 255},    // #323232
		PrimaryForegroundColor: color.RGBA{R: 250, G: 250, B: 250, A: 255}, // #fafafa

		SecondaryColor:           color.RGBA{R: 229, G: 229, B: 229, A: 255}, // #e5e5e5
		SecondaryForegroundColor: color.RGBA{R: 23, G: 23, B: 23, A: 255},    // #171717

		AccentColor:           color.RGBA{R: 229, G: 229, B: 229, A: 255}, // #e5e5e5
		AccentForegroundColor: color.RGBA{R: 23, G: 23, B: 23, A: 255},    // #171717

		DestructiveColor:           color.RGBA{R: 239, G: 68, B: 68, A: 255}, // #ef4444
		DestructiveForegroundColor: color.RGBA{R: 250, G: 250, B: 250, A: 255},

		MutedColor:           color.RGBA{R: 229, G: 229, B: 229, A: 255}, // #e5e5e5
		MutedForegroundColor: color.RGBA{R: 115, G: 115, B: 115, A: 255}, // #737373

		CardColor:              color.RGBA{R: 255, G: 255, B: 255, A: 255}, // #ffffff
		CardForegroundColor:    color.RGBA{R: 10, G: 10, B: 10, A: 255},    // #0a0a0a
		PopoverColor:           color.RGBA{R: 255, G: 255, B: 255, A: 255}, // #ffffff
		PopoverForegroundColor: color.RGBA{R: 10, G: 10, B: 10, A: 255},    // #0a0a0a

		BorderColor:      color.RGBA{R: 229, G: 229, B: 229, A: 255}, // #e5e5e5
		FocusBorderColor: color.RGBA{R: 23, G: 23, B: 23, A: 255},    // #171717
		BorderRadius:     8,
		RingColor:        color.RGBA{R: 23, G: 23, B: 23, A: 255}, // #171717

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
		ButtonBorderRadius:  8,

		InputBgColor:          color.RGBA{R: 255, G: 255, B: 255, A: 255},
		InputBorderColor:      color.RGBA{R: 229, G: 229, B: 229, A: 255}, // #e5e5e5
		InputFocusBorderColor: color.RGBA{R: 23, G: 23, B: 23, A: 255},    // #171717
		InputTextColor:        color.RGBA{R: 10, G: 10, B: 10, A: 255},
		InputPlaceholderColor: color.RGBA{R: 161, G: 161, B: 161, A: 255},
		InputBorderRadius:     8,
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
		SliderThumbColor: color.RGBA{R: 23, G: 23, B: 23, A: 255},    // #171717

		ProgressBarTrackColor: color.RGBA{R: 229, G: 229, B: 229, A: 255}, // #e5e5e5
		ProgressBarFillColor:  color.RGBA{R: 23, G: 23, B: 23, A: 255},    // #171717

		ScrollbarColor:      color.RGBA{R: 0, G: 0, B: 0, A: 30},
		ScrollbarTrackColor: color.RGBA{R: 0, G: 0, B: 0, A: 10},

		DividerColor: color.RGBA{R: 229, G: 229, B: 229, A: 255}, // #e5e5e5

		ShadowColor: color.RGBA{A: 30},

		MenuBg:               color.White,
		MenuItemSelectedBg:   color.RGBA{R: 245, G: 245, B: 245, A: 255}, // #f5f5f5
		MenuItemSelectedText: color.RGBA{R: 23, G: 23, B: 23, A: 255},    // #171717
		MenuItemHoverBg:      color.RGBA{R: 245, G: 245, B: 245, A: 255}, // #f5f5f5
	}
}

// DefaultDarkTheme 返回 Shadcn/UI Neutral 风格的深色主题。
func DefaultDarkTheme() *Theme {
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
		BorderRadius:     8,
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
		ButtonBorderRadius:  8,

		InputBgColor:          color.RGBA{R: 10, G: 10, B: 10, A: 255},
		InputBorderColor:      color.RGBA{R: 255, G: 255, B: 255, A: 25},
		InputFocusBorderColor: color.RGBA{R: 229, G: 229, B: 229, A: 255},
		InputTextColor:        color.RGBA{R: 229, G: 229, B: 229, A: 255},
		InputPlaceholderColor: color.RGBA{R: 115, G: 115, B: 115, A: 255},
		InputBorderRadius:     8,
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
