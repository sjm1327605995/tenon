package antdesign

import (
	"image/color"

	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/components"
	"github.com/sjm1327605995/tenon/pkg/fonts"
	"github.com/sjm1327605995/tenon/yoga"
)

// AntTextType defines text type.
type AntTextType string

const (
	AntTextPrimary   AntTextType = "primary"
	AntTextSecondary AntTextType = "secondary"
	AntTextSuccess   AntTextType = "success"
	AntTextWarning   AntTextType = "warning"
	AntTextDanger    AntTextType = "danger"
)

// AntText is a text component with variants.
type AntText struct {
	tenon.BaseWidget
	content    string
	textType   AntTextType
	customClr  color.Color
	code       bool
	mark       bool
	keyboard   bool
	underline  bool
	deleteLine bool
	strong     bool
	italic     bool
}

// NewAntText creates an AntText.
func NewAntText(content string) *AntText {
	t := &AntText{content: content}
	t.Init(t)
	return t
}

// Render returns the text UI.
func (t *AntText) Render() tenon.Component {
	theme := NewAntTheme()
	clr := t.resolveColor(theme)

	text := components.NewText(t.content).SetColor(clr)

	if t.strong {
		text.SetFontWeight(fonts.FontWeightBold)
	}
	if t.italic {
		text.SetFontStyle(fonts.FontStyleItalic)
	}

	if t.code {
		text.SetFontFamily(fonts.FontFamilyMono)
		return t.wrapWithBg(text, theme.BackgroundColor)
	}
	if t.mark {
		return t.wrapWithBg(text, theme.WarningBgColor)
	}
	if t.keyboard {
		return t.wrapKeyboard(text, theme)
	}
	if t.underline {
		return t.wrapUnderline(text, clr)
	}
	if t.deleteLine {
		return t.wrapDelete(text, clr)
	}
	return text
}

func (t *AntText) resolveColor(theme *AntTheme) color.Color {
	if t.customClr != nil {
		return t.customClr
	}
	switch t.textType {
	case AntTextPrimary:
		return theme.PrimaryColor
	case AntTextSecondary:
		return theme.TextMutedColor
	case AntTextSuccess:
		return theme.SuccessColor
	case AntTextWarning:
		return theme.WarningColor
	case AntTextDanger:
		return theme.ErrorColor
	default:
		return theme.TextColor
	}
}

func (t *AntText) wrapWithBg(text *components.Text, bg color.Color) tenon.Component {
	wrapper := components.NewView().
		SetBackgroundColor(bg).
		SetPadding(yoga.EdgeHorizontal, 4).
		SetPadding(yoga.EdgeVertical, 1).
		SetBorderRadius(4)
	wrapper.AddChild(text)
	return wrapper
}

func (t *AntText) wrapKeyboard(text *components.Text, theme *AntTheme) tenon.Component {
	wrapper := components.NewView().
		SetBackgroundColor(theme.BackgroundColor).
		SetBorder(yoga.EdgeAll, 1).
		SetBorderColor(theme.BorderColor).
		SetBorderRadius(4).
		SetPadding(yoga.EdgeHorizontal, 4).
		SetPadding(yoga.EdgeVertical, 1)
	wrapper.AddChild(text)
	return wrapper
}

func (t *AntText) wrapUnderline(text *components.Text, clr color.Color) tenon.Component {
	wrapper := components.NewView().
		SetFlexDirection(yoga.FlexDirectionColumn)
	wrapper.AddChild(text)
	line := components.NewView().
		SetHeight(1).
		SetBackgroundColor(clr)
	wrapper.AddChild(line)
	return wrapper
}

func (t *AntText) wrapDelete(text *components.Text, clr color.Color) tenon.Component {
	wrapper := components.NewView().
		SetFlexDirection(yoga.FlexDirectionColumn).
		SetJustifyContent(yoga.JustifyCenter)
	wrapper.AddChild(text)
	line := components.NewView().
		SetHeight(1).
		SetBackgroundColor(clr)
	wrapper.AddChild(line)
	return wrapper
}

func (t *AntText) SetType(tp AntTextType) *AntText { t.textType = tp; return t }
func (t *AntText) SetCode(v bool) *AntText         { t.code = v; return t }
func (t *AntText) SetMark(v bool) *AntText         { t.mark = v; return t }
func (t *AntText) SetKeyboard(v bool) *AntText     { t.keyboard = v; return t }
func (t *AntText) SetUnderline(v bool) *AntText    { t.underline = v; return t }
func (t *AntText) SetDelete(v bool) *AntText       { t.deleteLine = v; return t }
func (t *AntText) SetStrong(v bool) *AntText       { t.strong = v; return t }
func (t *AntText) SetItalic(v bool) *AntText       { t.italic = v; return t }
func (t *AntText) SetColor(c color.Color) *AntText { t.customClr = c; return t }
