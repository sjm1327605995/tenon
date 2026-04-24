package components

import (
	"testing"

	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

type testAppForMenu struct {
	core.BaseWidget
	page string
}

func pageTitle2(text string) core.Component {
	title := NewText(text).SetFontSize(core.GetTheme().FontSizeBase + 24).SetFontWeight(700)
	return NewView().SetMargin(yoga.EdgeBottom, 24).Add(title)
}

func demoBlock2(title string, demos ...core.Component) core.Component {
	card := NewView().
		SetPadding(yoga.EdgeAll, 24).
		SetBackgroundColor(core.GetTheme().SurfaceColor).
		SetBorderRadius(core.GetTheme().BorderRadius).
		SetBorder(yoga.EdgeAll, 1).
		SetBorderColor(core.GetTheme().BorderColor)

	card.Add(NewText(title).SetFontSize(core.GetTheme().FontSizeLG).SetMargin(yoga.EdgeBottom, 16))
	for _, d := range demos {
		card.AddChild(d)
	}
	return NewView().SetMargin(yoga.EdgeBottom, 16).Add(card)
}

func (a *testAppForMenu) pageA() core.Component {
	return NewView().SetFlexDirection(yoga.FlexDirectionColumn).Add(
		pageTitle2("Page A"),
		demoBlock2("Buttons",
			NewView().SetFlexDirection(yoga.FlexDirectionRow).Add(
				NewButton("Btn 1").SetWidth(100).SetHeight(36),
				NewButton("Btn 2").SetWidth(100).SetHeight(36).SetMargin(yoga.EdgeLeft, 8),
			),
		),
		demoBlock2("Cards",
			NewView().SetWidth(200).SetHeight(100).SetBackgroundColor(core.GetTheme().PrimaryColor),
		),
	)
}

func (a *testAppForMenu) pageB() core.Component {
	return NewView().SetFlexDirection(yoga.FlexDirectionColumn).Add(
		pageTitle2("Page B"),
		demoBlock2("Inputs",
			NewTextInput().SetWidth(300).SetPlaceholder("Type here"),
		),
		demoBlock2("Tags",
			NewView().SetFlexDirection(yoga.FlexDirectionRow).Add(
				NewText("Tag1").SetFontSize(core.GetTheme().FontSizeSM),
				NewText("Tag2").SetFontSize(core.GetTheme().FontSizeSM).SetMargin(yoga.EdgeLeft, 8),
			),
		),
	)
}

func (a *testAppForMenu) Render() core.Component {
	menu := NewMenu().SetItems([]MenuItemData{
		{Key: "a", Label: "Page A"},
		{Key: "b", Label: "Page B"},
	}).SetSelectedKey(a.page).SetOnSelect(func(key string) {
		a.SetState(func() { a.page = key })
	})

	content := NewView().SetFlexGrow(1).SetPadding(yoga.EdgeAll, 24)
	if a.page == "a" {
		content.AddChild(a.pageA())
	} else {
		content.AddChild(a.pageB())
	}

	scrollView := NewScrollView().SetWidthPercent(100).SetHeightPercent(100)
	scrollView.Content().AddChild(content)

	return NewView().SetFlexDirection(yoga.FlexDirectionRow).SetWidthPercent(100).SetHeightPercent(100).Add(menu, scrollView)
}

func TestMenuSwitchPanic(t *testing.T) {
	app := &testAppForMenu{page: "a"}
	app.Init(app)

	root := NewView()
	root.Init(root)
	root.AddChild(app)

	engine := core.NewEngine(root, 800, 600)
	engine.Mount()

	// 第一次切换
	app.SetState(func() { app.page = "b" })
	engine.Update()

	// 第二次切换
	app.SetState(func() { app.page = "a" })
	engine.Update()

	// 第三次切换
	app.SetState(func() { app.page = "b" })
	engine.Update()
}
