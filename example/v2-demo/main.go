package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/fonts"
	"github.com/sjm1327605995/tenon/pkg/v2/components"
	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// Counter demonstrates state-driven updates without rebuilding.
type Counter struct {
	core.BaseWidget
	count *core.State[int]
}

func NewCounter() *Counter {
	c := &Counter{count: core.NewState(0)}
	c.Init(c)
	return c
}

func (c *Counter) Build() core.Element {
	text := components.NewText("")
	btn := components.NewButton("+1")

	c.count.Subscribe(func(v int) {
		text.SetContent(fmt.Sprintf("Count: %d", v))
	})

	btn.SetOnClick(func() {
		c.count.Set(c.count.Get() + 1)
	})

	card := components.NewView().
		SetClass("card").
		SetPadding(yoga.EdgeAll, 24).
		SetFlexDirection(yoga.FlexDirectionRow).
		SetGap(yoga.GutterAll, 12).
		SetAlignItems(yoga.AlignCenter).
		Add(text, btn)
	return card
}

// ToggleDemo demonstrates conditional rendering with If.
type ToggleDemo struct {
	core.BaseWidget
	show *core.State[bool]
}

func NewToggleDemo() *ToggleDemo {
	t := &ToggleDemo{show: core.NewState(false)}
	t.Init(t)
	return t
}

func (t *ToggleDemo) Build() core.Element {
	toggleBtn := components.NewButton("Toggle")
	toggleBtn.SetOnClick(func() {
		t.show.Set(!t.show.Get())
	})

	card := components.NewView().
		SetClass("card").
		SetPadding(yoga.EdgeAll, 24).
		SetFlexDirection(yoga.FlexDirectionColumn).
		SetGap(yoga.GutterAll, 12).
		Add(
			toggleBtn,
			core.NewIf(t.show).
				Then(func() core.Element {
					return components.NewText("Hello, World!").SetColor(color.RGBA{R: 0, G: 150, B: 0, A: 255})
				}).
				Else(func() core.Element {
					return components.NewText("Hidden...").SetColor(color.RGBA{R: 150, G: 0, B: 0, A: 255})
				}),
		)
	return card
}

// App is the root widget.
type App struct {
	core.BaseWidget
	page *core.State[string]
}

func NewApp() *App {
	a := &App{page: core.NewState("counter")}
	a.Init(a)
	return a
}

func (a *App) Build() core.Element {
	counterBtn := components.NewButton("Counter")
	counterBtn.SetOnClick(func() { a.page.Set("counter") })

	toggleBtn := components.NewButton("Toggle")
	toggleBtn.SetOnClick(func() { a.page.Set("toggle") })

	// Navigation bar
	nav := components.NewView().
		SetFlexDirection(yoga.FlexDirectionRow).
		SetGap(yoga.GutterAll, 8).
		SetPadding(yoga.EdgeBottom, 16).
		Add(counterBtn, toggleBtn)

	// Page content via Switch
	content := core.NewSwitch(a.page).
		Case("counter", NewCounter().Build).
		Case("toggle", NewToggleDemo().Build).
		Default(func() core.Element {
			return components.NewText("Select a page")
		})

	// Root view - directly use *View to set background color
	root := components.NewView()
	root.SetWidthPercent(100)
	root.SetHeightPercent(100)
	root.SetPadding(yoga.EdgeAll, 24)
	root.SetBackgroundColor(color.RGBA{R: 245, G: 245, B: 245, A: 255})
	root.Add(nav, content)

	return root
}

func main() {
	if err := fonts.InitDefaultFont(); err != nil {
		panic("failed to init font: " + err.Error())
	}
	if err := fonts.LoadFontFromFile(fonts.FontFamilySans, "font/OPPOSans-Medium.ttf"); err == nil {
		fonts.SetDefaultFontFamily(fonts.FontFamilySans)
	}

	core.RegisterStyle("card", func(e core.Element) {
		if v, ok := e.(*components.View); ok {
			v.SetBackgroundColor(color.White).
				SetBorderRadius(12).
				SetShadow(color.RGBA{A: 32}, 8, 0, 2)
		}
	})

	app := NewApp()
	engine := core.NewEngine(app, 900, 600)
	engine.Mount()

	fmt.Println("Tenon v2 demo started at 900x600")

	if err := ebiten.RunGame(engine); err != nil {
		panic(err)
	}
}
