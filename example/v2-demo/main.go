package main

import (
	"fmt"
	"image/color"

	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/fonts"
	"github.com/sjm1327605995/tenon/components"
	"github.com/sjm1327605995/tenon/internal/core"
	"github.com/sjm1327605995/tenon/widgets"
	"github.com/sjm1327605995/tenon/yoga"
)

// ComponentGallery showcases all v2 native components with Shadcn/UI-inspired styling.
type ComponentGallery struct {
	core.BaseWidget
	page *core.State[string]
}

func NewComponentGallery() *ComponentGallery {
	g := &ComponentGallery{page: core.NewState("buttons")}
	g.Init(g)
	return g
}

func (g *ComponentGallery) Render() core.Element {
	// Navigation bar
	pages := []string{"buttons", "inputs", "controls", "scroll", "cards"}
	var navButtons []core.Element
	for _, p := range pages {
		btn := components.NewButton(p)
		btn.SetOnClick(func(page string) func() {
			return func() { g.page.Set(page) }
		}(p))
		if g.page.Get() == p {
			btn.SetVariant(components.ButtonSecondary)
		}
		navButtons = append(navButtons, btn)
	}

	nav := components.NewView().
		SetFlexDirection(yoga.FlexDirectionRow).
		SetGap(yoga.GutterAll, 8).
		SetPadding(yoga.EdgeBottom, 16).
		SetAlignItems(yoga.AlignCenter).
		Add(navButtons...)

	// Page content
	content := widgets.NewSwitch(g.page).
		Case("buttons", g.buildButtonsPage).
		Case("inputs", g.buildInputsPage).
		Case("controls", g.buildControlsPage).
		Case("scroll", g.buildScrollPage).
		Case("cards", g.buildCardsPage).
		Default(func() core.Element {
			return components.NewText("Select a page")
		})

	// Root
	root := components.NewView()
	root.SetWidthPercent(100)
	root.SetHeightPercent(100)
	root.SetPadding(yoga.EdgeAll, 24)
	root.SetBackgroundColor(core.GetTheme().BackgroundColor)
	root.Add(nav, content)

	return root
}

func (g *ComponentGallery) buildButtonsPage() core.Element {
	title := components.NewText("Button Variants").SetFontSize(20).SetColor(core.GetTheme().TextColor)
	subtitle := components.NewText("Shadcn/UI-inspired button styles").SetFontSize(14).SetColor(core.GetTheme().MutedForegroundColor)

	// Default
	defaultBtn := components.NewButton("Default")
	defaultBtn.SetOnClick(func() { fmt.Println("Default clicked") })

	// Secondary
	secondaryBtn := components.NewButton("Secondary").SetVariant(components.ButtonSecondary)

	// Outline
	outlineBtn := components.NewButton("Outline").SetVariant(components.ButtonOutline)

	// Ghost
	ghostBtn := components.NewButton("Ghost").SetVariant(components.ButtonGhost)

	// Destructive
	destructiveBtn := components.NewButton("Destructive").SetVariant(components.ButtonDestructive)

	// Disabled
	disabledBtn := components.NewButton("Disabled").SetDisabled(true)

	row1 := components.NewView().
		SetFlexDirection(yoga.FlexDirectionRow).
		SetGap(yoga.GutterAll, 8).
		SetAlignItems(yoga.AlignCenter).
		Add(defaultBtn, secondaryBtn, outlineBtn, ghostBtn, destructiveBtn, disabledBtn)

	// Sizes
	smallBtn := components.NewButton("Small").SetVariant(components.ButtonSecondary)
	smallBtn.BaseElement.SetPadding(yoga.EdgeHorizontal, 12)
	smallBtn.BaseElement.SetPadding(yoga.EdgeVertical, 6)

	largeBtn := components.NewButton("Large").SetVariant(components.ButtonSecondary)
	largeBtn.BaseElement.SetPadding(yoga.EdgeHorizontal, 24)
	largeBtn.BaseElement.SetPadding(yoga.EdgeVertical, 14)

	row2 := components.NewView().
		SetFlexDirection(yoga.FlexDirectionRow).
		SetGap(yoga.GutterAll, 8).
		SetAlignItems(yoga.AlignCenter).
		Add(smallBtn, largeBtn)

	card := components.NewView().
		SetClass("card").
		SetFlexDirection(yoga.FlexDirectionColumn).
		SetGap(yoga.GutterAll, 16).
		SetPadding(yoga.EdgeAll, 24).
		Add(title, subtitle, row1, row2)

	return card
}

func (g *ComponentGallery) buildInputsPage() core.Element {
	title := components.NewText("Input & Form Elements").SetFontSize(20).SetColor(core.GetTheme().TextColor)

	input := components.NewTextInput()
	input.SetPlaceholder("Type something...")
	input.SetOnChange(func(text string) {
		core.LogDebug("Input", "text", text)
	})

	input2 := components.NewTextInput()
	input2.SetPlaceholder("Another input field...")

	card := components.NewView().
		SetClass("card").
		SetFlexDirection(yoga.FlexDirectionColumn).
		SetGap(yoga.GutterAll, 12).
		SetPadding(yoga.EdgeAll, 24).
		Add(title, input, input2)

	return card
}

func (g *ComponentGallery) buildControlsPage() core.Element {
	title := components.NewText("Controls Demo").SetFontSize(20).SetColor(core.GetTheme().TextColor)

	// ProgressBar
	progressLabel := components.NewText("Progress: 30%").SetColor(core.GetTheme().TextColor)
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
	switchLabel := components.NewText("Switch: Off").SetColor(core.GetTheme().TextColor)
	sw := components.NewSwitch()
	sw.SetOnChange(func(checked bool) {
		if checked {
			switchLabel.SetContent("Switch: On")
		} else {
			switchLabel.SetContent("Switch: Off")
		}
	})

	// Checkbox
	checkLabel := components.NewText("Checkbox: Unchecked").SetColor(core.GetTheme().TextColor)
	cb := components.NewCheckbox("Accept terms")
	cb.SetOnChange(func(checked bool) {
		if checked {
			checkLabel.SetContent("Checkbox: Checked")
		} else {
			checkLabel.SetContent("Checkbox: Unchecked")
		}
	})

	// Radio
	radioLabel := components.NewText("Radio: Unselected").SetColor(core.GetTheme().TextColor)
	rb := components.NewRadio("Select me")
	rb.SetOnChange(func(selected bool) {
		if selected {
			radioLabel.SetContent("Radio: Selected")
		}
	})

	// Slider
	sliderLabel := components.NewText("Slider: 0.00").SetColor(core.GetTheme().TextColor)
	sl := components.NewSlider(0, 100)
	sl.SetOnChange(func(v float32) {
		sliderLabel.SetContent(fmt.Sprintf("Slider: %.2f", v))
	})

	// Badge variants
	badgeDefault := components.NewBadge("default").SetVariant(components.BadgeDefault)
	badgeSecondary := components.NewBadge("secondary").SetVariant(components.BadgeSecondary)
	badgeOutline := components.NewBadge("outline").SetVariant(components.BadgeOutline)
	badgeDestructive := components.NewBadge("destructive").SetVariant(components.BadgeDestructive)
	badgeRow := components.NewView().
		SetFlexDirection(yoga.FlexDirectionRow).
		SetGap(yoga.GutterAll, 8).
		SetAlignItems(yoga.AlignCenter).
		Add(badgeDefault, badgeSecondary, badgeOutline, badgeDestructive)

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
			components.NewText("Badges:").SetColor(core.GetTheme().TextColor),
			badgeRow,
		)

	return card
}

func (g *ComponentGallery) buildScrollPage() core.Element {
	title := components.NewText("ScrollView Demo").SetFontSize(20).SetColor(core.GetTheme().TextColor)

	scrollView := components.NewScrollView()
	scrollView.SetHeight(300)
	scrollView.SetBackgroundColor(core.GetTheme().MutedColor)
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
			components.NewText(fmt.Sprintf("Item %d", i+1)).SetColor(core.GetTheme().TextColor),
			components.NewView().SetFlexGrow(1),
			func() *components.Button {
				btn := components.NewButton("Action").SetVariant(components.ButtonGhost)
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

func (g *ComponentGallery) buildCardsPage() core.Element {
	title := components.NewText("Card & Modal Demo").SetFontSize(20).SetColor(core.GetTheme().TextColor)

	modalBtn := components.NewButton("Open Modal")
	modal := components.NewModal()
	modal.SetTitle("Shadcn-style Modal")
	modal.Panel().Add(
		components.NewText("This is a modal dialog with modern styling.").SetColor(core.GetTheme().TextColor),
		components.NewView().SetHeight(16),
		components.NewButton("Close").SetOnClick(func() {
			modal.Close()
		}),
	)
	modalBtn.SetOnClick(func() {
		modal.Open()
	})

	// Toast demo
	toastBtn := components.NewButton("Show Toast").SetVariant(components.ButtonOutline)
	toastManager := components.NewToastManager()
	toastManager.SetPosition(yoga.EdgeRight, 24)
	toastManager.SetPosition(yoga.EdgeTop, 24)
	toastBtn.SetOnClick(func() {
		toastManager.Show("Action completed successfully!", 3000)
	})

	card := components.NewView().
		SetClass("card").
		SetFlexDirection(yoga.FlexDirectionColumn).
		SetGap(yoga.GutterAll, 12).
		SetPadding(yoga.EdgeAll, 24).
		Add(title, modalBtn, toastBtn, modal, toastManager)

	return card
}

func main() {
	if err := fonts.InitDefaultFont(); err != nil {
		panic("failed to init font: " + err.Error())
	}
	if err := fonts.LoadFontFromFile(fonts.FontFamilySans, "font/OPPOSans-Medium.ttf"); err == nil {
		fonts.SetDefaultFontFamily(fonts.FontFamilySans)
	}

	// Use Shadcn/UI-inspired theme
	core.SetTheme(core.DefaultShadcnLightTheme())

	tenon.RegisterStyle("card", func(e tenon.Element) {
		if v, ok := e.(*components.View); ok {
			v.SetBackgroundColor(core.GetTheme().CardColor).
				SetBorderRadius(core.GetTheme().BorderRadius).
				SetBorderColor(core.GetTheme().BorderColor).
				SetShadow(color.RGBA{A: 20}, 12, 0, 2)
		}
	})

	app := NewComponentGallery()
	tenon.Run(app, 900, 600)
}
