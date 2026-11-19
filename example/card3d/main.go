// Command card3d 演示伪 3D 卡片：Perspective + RotateX/RotateY/TranslateZ，
// 对应 CSS 的 transform: perspective(p) rotateX() rotateY() translateZ()。
//
// 拖动卡片改变倾角、松手回正；hover 时整张卡朝观察者浮起。注意卡片里的文字和照片
// 也跟着一起透视变形 —— 内容仍是平面，只是被投影到一个四边形上，所以叫「伪 3D」。
// 照片是纯图片元素：投影层里的整块图自动走精确单应贴图（四角真实投影，无斜切）。
//
// 这些变换只作用于绘制、不进入布局，因此无论卡片怎么倾，它在栈里占的位置和下方文字的
// 排版都不变。
//
//	go run ./example/card3d
package main

import (
	"image"
	"image/color"
	"math"

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
		// 照片是纯图片元素：作为投影层的直接子元素，自动按四角真实投影精确贴图。
		ui.Img(ui.SrcImage("hero", heroImage()), ui.Style(ui.Width(300), ui.Height(210))),
		ui.Box([]ui.StyleOpt{
			ui.Height(190), ui.ItemsCenter, ui.JustifyCenter, ui.Gap(6),
			ui.LinearGradient(ui.Hex("#0f172a00"), ui.Hex("#0f172a"), 90),
		},
			ui.Text("HERO", ui.FontSize(26), ui.Bold, ui.TextColor(ui.White)),
			ui.Text("整块图自动精确单应贴图", ui.FontSize(12), ui.TextColor(ui.Hex("#94a3b8"))),
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

// heroImage 程序化生成一张照片：天空渐变 + 太阳 + 远山剪影。
// 不用读盘资源，演示里拖动画质差异（斜切 vs 精确单应）时图案得可控。
func heroImage() image.Image {
	const w, h = 600, 420
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	top, bot := color.RGBA{59, 130, 246, 255}, color.RGBA{139, 92, 246, 255}
	for y := 0; y < h; y++ {
		t := float32(y) / (h - 1)
		c := color.RGBA{
			R: uint8(float32(top.R)*(1-t) + float32(bot.R)*t),
			G: uint8(float32(top.G)*(1-t) + float32(bot.G)*t),
			B: uint8(float32(top.B)*(1-t) + float32(bot.B)*t),
			A: 255,
		}
		for x := 0; x < w; x++ {
			img.SetRGBA(x, y, c)
		}
	}
	// 太阳：柔边圆
	const cx, cy, r = 430, 110, 56
	for y := cy - r - 8; y <= cy+r+8; y++ {
		for x := cx - r - 8; x <= cx+r+8; x++ {
			if x < 0 || y < 0 || x >= w || y >= h {
				continue
			}
			dx, dy := float32(x-cx), float32(y-cy)
			d := float32(math.Sqrt(float64(dx*dx + dy*dy)))
			if d > r+6 {
				continue
			}
			a := float32(1)
			if d > r {
				a = 1 - (d-r)/6
			}
			img.SetRGBA(x, y, blend(img.RGBAAt(x, y), color.RGBA{254, 240, 138, 255}, a))
		}
	}
	// 远山：两条正弦剪影
	for x := 0; x < w; x++ {
		y1 := h - 90 - int(34*math.Sin(float64(x)/w*6.28))
		y2 := h - 40 - int(50*math.Sin(float64(x)/w*4.1+2))
		for y := y1; y < h; y++ {
			img.SetRGBA(x, y, color.RGBA{30, 41, 59, 255})
		}
		for y := y2; y < h; y++ {
			img.SetRGBA(x, y, color.RGBA{15, 23, 42, 255})
		}
	}
	return img
}

func blend(dst, src color.RGBA, a float32) color.RGBA {
	return color.RGBA{
		R: uint8(float32(dst.R)*(1-a) + float32(src.R)*a),
		G: uint8(float32(dst.G)*(1-a) + float32(src.G)*a),
		B: uint8(float32(dst.B)*(1-a) + float32(src.B)*a),
		A: 255,
	}
}
