package main

import (
	"fmt"
	"image/color"

	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/fonts"
	"github.com/sjm1327605995/tenon/pkg/v2/components"
	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// ComponentGallery showcases all v2 native components.
type ComponentGallery struct {
	core.BaseWidget
	page *core.State[string]
}

func NewComponentGallery() *ComponentGallery {
	g := &ComponentGallery{page: core.NewState("scroll")}
	g.Init(g)
	return g
}

func (g *ComponentGallery) Build() core.Element {
	// Navigation bar
	controlsBtn := components.NewButton("Controls")
	controlsBtn.SetOnClick(func() { g.page.Set("controls") })

	scrollBtn := components.NewButton("Scroll")
	scrollBtn.SetOnClick(func() { g.page.Set("scroll") })

	inputBtn := components.NewButton("Input")
	inputBtn.SetOnClick(func() { g.page.Set("input") })

	nav := components.NewView().
		SetFlexDirection(yoga.FlexDirectionRow).
		SetGap(yoga.GutterAll, 8).
		SetPadding(yoga.EdgeBottom, 16).
		Add(controlsBtn, scrollBtn, inputBtn)

	// Page content
	content := core.NewSwitch(g.page).
		Case("controls", g.buildControlsPage).
		Case("scroll", g.buildScrollPage).
		Case("input", g.buildInputPage).
		Default(func() core.Element {
			return components.NewText("Select a page")
		})

	// Root
	root := components.NewView()
	root.SetWidthPercent(100)
	root.SetHeightPercent(100)
	root.SetPadding(yoga.EdgeAll, 24)
	root.SetBackgroundColor(color.RGBA{R: 100, G: 100, B: 100, A: 255})
	root.Add(nav, content)

	return root
}

func (g *ComponentGallery) buildControlsPage() core.Element {
	// Title
	title := components.NewText("Controls Demo").SetFontSize(24)

	// ProgressBar
	progressLabel := components.NewText("Progress: 30%")
	pb := components.NewProgressBar().SetProgress(0.3)
	progressBtn := components.NewButton("+10%")
	progressBtn.SetOnClick(func() {
		p := pb.GetProgress() + 0.1
		if p > 1 {
			p = 0
		}
		pb.SetProgress(p)
		progressLabel.SetContent(fmt.Sprintf("Progress: %d%%", int(p*100)))
	})

	// Switch
	switchLabel := components.NewText("Switch: Off")
	sw := components.NewSwitch()
	sw.SetOnChange(func(checked bool) {
		if checked {
			switchLabel.SetContent("Switch: On")
		} else {
			switchLabel.SetContent("Switch: Off")
		}
	})

	// Checkbox
	checkLabel := components.NewText("Checkbox: Unchecked")
	cb := components.NewCheckbox("Accept terms")
	cb.SetOnChange(func(checked bool) {
		if checked {
			checkLabel.SetContent("Checkbox: Checked")
		} else {
			checkLabel.SetContent("Checkbox: Unchecked")
		}
	})

	// Radio
	radioLabel := components.NewText("Radio: Unselected")
	rb := components.NewRadio("Select me")
	rb.SetOnChange(func(selected bool) {
		if selected {
			radioLabel.SetContent("Radio: Selected")
		}
	})

	// Slider
	sliderLabel := components.NewText("Slider: 0.00")
	sl := components.NewSlider(0, 100)
	sl.SetOnChange(func(v float32) {
		sliderLabel.SetContent(fmt.Sprintf("Slider: %.2f", v))
	})

	// Divider
	div := components.NewDivider()

	// Layout
	card := components.NewView().
		SetClass("card").
		SetFlexDirection(yoga.FlexDirectionColumn).
		SetGap(yoga.GutterAll, 12).
		SetPadding(yoga.EdgeAll, 24).
		Add(
			title,
			div,
			progressLabel, pb, progressBtn,
			switchLabel, sw,
			checkLabel, cb,
			radioLabel, rb,
			sliderLabel, sl,
		)

	return card
}

func (g *ComponentGallery) buildScrollPage() core.Element {
	title := components.NewText("ScrollView Demo").SetFontSize(24)

	scrollView := components.NewScrollView()
	scrollView.SetHeight(300)
	scrollView.SetBackgroundColor(color.RGBA{R: 0, G: 255, B: 255, A: 255})
	content := scrollView.Content()
	content.SetFlexDirection(yoga.FlexDirectionColumn)
	content.SetGap(yoga.GutterAll, 8)
	content.SetPadding(yoga.EdgeAll, 16)

	for i := 0; i < 50; i++ {
		row := components.NewView().
			SetFlexDirection(yoga.FlexDirectionRow).
			SetAlignItems(yoga.AlignCenter).
			SetGap(yoga.GutterAll, 8)
		row.SetHeight(40)
		row.Add(
			components.NewText(fmt.Sprintf("Item %d", i+1)),
			components.NewView().SetFlexGrow(1),
			func() *components.Button {
				btn := components.NewButton("Action")
				btn.SetOnClick(func() {
					fmt.Println("Action button clicked!")
				})
				return btn
			}(),
		)
		content.Add(row)
	}

	card := components.NewView().
		SetClass("card").
		SetFlexDirection(yoga.FlexDirectionColumn).
		SetGap(yoga.GutterAll, 12).
		SetPadding(yoga.EdgeAll, 24).
		Add(title, scrollView)

	return card
}

func (g *ComponentGallery) buildInputPage() core.Element {
	title := components.NewText("TextInput Demo").SetFontSize(24)

	input := components.NewTextInput()
	input.SetPlaceholder("Type something...")
	input.SetOnChange(func(text string) {
		core.LogDebug("Input", "text", text)
	})
	input.SetOnSubmit(func(text string) {
		core.LogDebug("Submit", "text", text)
	})

	card := components.NewView().
		SetClass("card").
		SetFlexDirection(yoga.FlexDirectionColumn).
		SetGap(yoga.GutterAll, 12).
		SetPadding(yoga.EdgeAll, 24).
		Add(title, input)

	return card
}

func main() {
	if err := fonts.InitDefaultFont(); err != nil {
		panic("failed to init font: " + err.Error())
	}
	if err := fonts.LoadFontFromFile(fonts.FontFamilySans, "font/OPPOSans-Medium.ttf"); err == nil {
		fonts.SetDefaultFontFamily(fonts.FontFamilySans)
	}

	tenon.RegisterStyle("card", func(e tenon.Element) {
		if v, ok := e.(*components.View); ok {
			v.SetBackgroundColor(color.White).
				SetBorderRadius(12).
				SetShadow(color.RGBA{A: 32}, 8, 0, 2)
		}
	})

	app := NewComponentGallery()
	tenon.Run(app, 900, 600)
}
