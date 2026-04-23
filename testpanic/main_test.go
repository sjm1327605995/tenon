package main

import (
	"testing"

	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/components"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

type testApp struct {
	tenon.BaseWidget
	page string
}

func (a *testApp) Render() tenon.Component {
	menu := components.NewMenu().SetItems([]components.MenuItemData{
		{Key: "a", Label: "Page A"},
		{Key: "b", Label: "Page B"},
	}).SetSelectedKey(a.page).SetOnSelect(func(key string) {
		a.page = key
		a.Invalidate()
	})

	content := components.NewView().SetFlexGrow(1).SetPadding(yoga.EdgeAll, 24)
	if a.page == "a" {
		content.Add(components.NewText("Page A"))
	} else {
		content.Add(components.NewText("Page B"))
	}

	scrollView := components.NewScrollView().SetWidthPercent(100).SetHeightPercent(100)
	scrollView.Content().AddChild(content)

	return components.NewView().SetFlexDirection(yoga.FlexDirectionRow).SetWidthPercent(100).SetHeightPercent(100).Add(menu, scrollView)
}

func TestMenuSwitch(t *testing.T) {
	app := &testApp{page: "a"}
	app.Init(app)

	engine := core.NewEngine(app, 800, 600)
	engine.Mount()

	// 模拟点击菜单切换到 page b
	app.page = "b"
	app.Invalidate()
	engine.FlushUpdates()
}
