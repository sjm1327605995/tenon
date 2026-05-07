package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/fonts"
)

// Counter 是一个 go-tui 风格的组件。
type Counter struct {
	count *tenon.State[int]
}

func NewCounter() *Counter {
	return &Counter{count: tenon.NewState(0)}
}

// KeyMap 声明组件级键盘绑定。
func (c *Counter) KeyMap() tenon.KeyMap {
	return tenon.KeyMap{
		tenon.On(ebiten.KeyArrowUp, func(ke tenon.KeyEvent) {
			c.count.Update(func(v int) int { return v + 1 })
		}),
		tenon.On(ebiten.KeyArrowDown, func(ke tenon.KeyEvent) {
			c.count.Update(func(v int) int { return v - 1 })
		}),
		tenon.On(ebiten.KeyQ, func(ke tenon.KeyEvent) {
			ke.App.Stop()
		}),
	}
}

// Render 返回组件的 UI 描述。
func (c *Counter) Render(app *tenon.App) tenon.Widget {
	// Fragment 让 Component 可以像 React 函数组件一样返回多个子节点
	return tenon.Fragment(
		tenon.Column(
			tenon.Text(fmt.Sprintf("Count: %d", c.count.Get())).FontSize(32),
			tenon.Text("↑ / ↓ to change, Q to quit").FontSize(14).Color(tenon.GetTheme().TextMutedColor),
		).Gapf(12).Paddingf(tenon.EdgeInsetsAll(48)).AlignItems(tenon.AlignCenter),
	)
}

func main() {
	if err := fonts.InitDefaultFont(); err != nil {
		panic(err)
	}

	app, err := tenon.NewApp(
		tenon.WithSize(640, 480),
		tenon.WithTitle("Counter Demo"),
		tenon.WithRootComponent(NewCounter()),
	)
	if err != nil {
		panic(err)
	}
	if err := app.Run(); err != nil {
		panic(err)
	}
}
