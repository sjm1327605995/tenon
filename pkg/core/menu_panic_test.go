package core

import (
	"testing"

	"github.com/sjm1327605995/tenon/pkg/components"
	"github.com/sjm1327605995/tenon/yoga"
)

type testAppForMenu struct {
	BaseWidget
	page string
}

func (a *testAppForMenu) Render() Component {
	menu := components.NewMenu().SetItems([]components.MenuItemData{
		{Key: "a", Label: "Page A"},
		{Key: "b", Label: "Page B"},
	}).SetSelectedKey(a.page).SetOnSelect(func(key string) {
		a.SetState(func() { a.page = key })
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

func TestMenuSwitchPanic(t *testing.T) {
	app := &testAppForMenu{page: "a"}
	app.Init(app)

	root := newTestHost("root")
	root.AddChild(app)

	engine := NewEngine(root, 800, 600)
	engine.Mount()

	// 模拟点击菜单切换到 page b
	app.SetState(func() { app.page = "b" })
	engine.flushUpdates()
}
