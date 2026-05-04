package main

import (
	"fmt"
	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/fonts"
)

type galleryApp struct{}

func (g *galleryApp) Build() tenon.Widget {
	t := tenon.GetTheme()
	muted := t.TextMutedColor
	name, email := "", ""
	text := ""

	return tenon.Column(
		// Header
		tenon.Text("TextInput Debug").FontSize(28).Color(t.TextColor),
		tenon.Text("Debug page for text input").FontSize(14).Color(muted),
		tenon.Separator(tenon.SeparatorHorizontal),

		// TextField
		sectionTitle("TextField"),
		tenon.Row(
			tenon.Container(tenon.Text("Name:").FontSize(14)).W(60),
			tenon.TextField(name).
				W(200).
				Placeholder("Enter your name").
				OnChange(func(v string) { fmt.Printf("[APP] name=%q\n", v) }),
		).Gapf(8).AlignItems(tenon.AlignCenter),
		tenon.Row(
			tenon.Container(tenon.Text("Email:").FontSize(14)).W(60),
			tenon.TextField(email).
				W(200).
				Placeholder("Enter your email").
				OnChange(func(v string) { fmt.Printf("[APP] email=%q\n", v) }),
		).Gapf(8).AlignItems(tenon.AlignCenter),

		tenon.Separator(tenon.SeparatorHorizontal),

		// Textarea
		sectionTitle("Textarea"),
		tenon.Textarea(text, func(v string) {
			fmt.Printf("[APP] textarea=%q\n", v)
		}),

		tenon.Separator(tenon.SeparatorHorizontal),

		// EditableText
		sectionTitle("EditableText"),
		tenon.EditableText(text).
			Size(16).
			Placeholder("Click to edit").
			OnChange(func(v string) { fmt.Printf("[APP] editable=%q\n", v) }),
	).Gapf(16).Paddingf(tenon.EdgeInsetsAll(24))
}

func sectionTitle(text string) tenon.Widget {
	return tenon.Label(text)
}

func main() {
	if err := fonts.InitDefaultFont(); err != nil {
		panic("failed to init font: " + err.Error())
	}
	if err := fonts.ReloadFontFromFile(fonts.FontFamilyDefault, "font/OPPOSans-Medium.ttf"); err != nil {
		panic("failed to load CJK font: " + err.Error())
	}
	if err := fonts.PreloadCommonSizes(fonts.FontFamilyDefault, []float32{12, 14, 16, 18, 20, 24, 32, 48}); err != nil {
		panic("failed to preload CJK font sizes: " + err.Error())
	}
	app := &galleryApp{}
	tenon.SetTheme(tenon.DefaultLightTheme())
	tenon.Run(app.Build, 900, 600)
}
