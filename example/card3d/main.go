// Command card3d 演示伪 3D 卡片：Perspective + RotateX/RotateY/TranslateZ，
// 对应 CSS 的 transform: perspective(p) rotateX() rotateY() translateZ()。
//
// 拖动卡片改变倾角、松手回正；hover 时整张卡朝观察者浮起。注意卡片里的文字也跟着一起
// 透视变形 —— 内容仍是平面，只是被投影到一个四边形上，所以叫「伪 3D」。
//
// 这些变换只作用于绘制、不进入布局，因此无论卡片怎么倾，它在栈里占的位置和下方文字的
// 排版都不变。
//
//	go run ./example/card3d
package main

import (
	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

func main() {
	ui.WindowSize(720, 760)
	ui.Run(ui.Use(App, struct{}{}))
}

func App(_ struct{}) *ui.Node {
	return ui.Box([]ui.StyleOpt{
		ui.Fill, ui.Bg(ui.Hex("#0f172a")), ui.ItemsCenter, ui.JustifyCenter, ui.Gap(28),
	},
		ui.Use(card, struct{}{}),
		ui.Text("拖动卡片倾斜它 · 下方文字的排版始终不受 3D 影响",
			ui.FontSize(14), ui.TextColor(ui.Hex("#94a3b8"))),
	)
}

func card(_ struct{}) *ui.Node {
	rx, setRX := ui.UseState[float32](0)
	ry, setRY := ui.UseState[float32](0)
	hover, setHover := ui.UseState(false)

	// 平滑过渡：松手回正、hover 时朝观察者浮起。
	ax := ui.UseTween(rx, 180, ui.EaseOut)
	ay := ui.UseTween(ry, 180, ui.EaseOut)
	lift := float32(0)
	if hover {
		lift = 60
	}
	tz := ui.UseTween(lift, 220, ui.EaseOut)

	return ui.Box([]ui.StyleOpt{
		ui.Width(300), ui.Height(400), ui.Radius(18), ui.Clip, ui.Column,
		ui.Bg(ui.Hex("#1e293b")), ui.Border(1, ui.Hex("#334155")),
		// —— 伪 3D：透视 + 绕 X/Y 轴旋转 + Z 位移（锚定卡片中心）——
		ui.Perspective(700), ui.RotateX(ax), ui.RotateY(ay), ui.TranslateZ(tz),
	},
		ui.Box([]ui.StyleOpt{
			ui.Height(210), ui.ItemsCenter, ui.JustifyCenter,
			ui.LinearGradient(ui.Hex("#3b82f6"), ui.Hex("#8b5cf6"), 135),
		},
			ui.Text("HERO", ui.FontSize(30), ui.Bold, ui.TextColor(ui.White)),
		),
		ui.Box([]ui.StyleOpt{ui.Padding(20), ui.Gap(8), ui.Grow(1)},
			ui.Text("伪 3D 卡片", ui.FontSize(20), ui.Semibold, ui.TextColor(ui.White)),
			ui.Text("perspective + rotateX/rotateY。拖动我，松手回正。",
				ui.FontSize(13), ui.TextColor(ui.Hex("#94a3b8"))),
		),

		// 交互：拖动改变倾角，离开时回正。附加为子节点即绑定到本卡片。
		ui.OnHover(func(on bool) {
			setHover(on)
			if !on {
				setRX(0)
				setRY(0)
			}
		}),
		ui.OnDrag(func(dx, dy float32) {
			setRX(clamp(rx-dy*0.4, -35, 35)) // 上下拖 -> 绕 X 轴
			setRY(clamp(ry+dx*0.4, -35, 35)) // 左右拖 -> 绕 Y 轴
		}),
	)
}

func clamp(v, lo, hi float32) float32 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
