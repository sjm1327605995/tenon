package ui

import (
	"fmt"
	"testing"
	"time"

	"gioui.org/gpu/headless"
	"gioui.org/op"
)

// nowMs 是毫秒墙钟（基准分段用）。
func nowMs() float64 { return float64(time.Now().UnixNano()) / 1e6 }

// 帧耗时基准：区分「本引擎 paint 录制」与「gio 渲染（CPU 处理 + GPU 提交）」两段。
// headless 的 GPU 是软件光栅（llvmpipe 类），绝对值偏悲观，但两段各占多少、
// 随场景规模怎么涨，是真实可信的 —— 用于回答「瓶颈在引擎还是在 gio」。
//
//	go test ./pkg/ui -bench Frame -benchtime 60x
func benchFrame(b *testing.B, tree *Node, w, h int) {
	win, err := headless.NewWindow(w, h)
	if err != nil {
		b.Skipf("headless GPU 不可用: %v", err)
	}
	defer win.Release()

	hn := Mount(Use(func(struct{}) *Node { return tree }, struct{}{}), w, h)
	layoutAll(hn.g)

	var ops op.Ops
	// 预热：首帧含字形图集上传等一次性成本，不计入。
	p := newGioPainter(&ops, w, h)
	paint(p, hn.g.rootRN)
	if err := win.Frame(&ops); err != nil {
		b.Fatal(err)
	}

	var paintMs, frameMs float64
	for i := 0; i < b.N; i++ {
		ops.Reset()
		p := newGioPainter(&ops, w, h)

		t0 := nowMs()
		paint(p, hn.g.rootRN)
		t1 := nowMs()
		if err := win.Frame(&ops); err != nil {
			b.Fatal(err)
		}
		t2 := nowMs()

		paintMs += t1 - t0
		frameMs += t2 - t1
	}
	b.ReportMetric(paintMs/float64(b.N), "paint-ms/frame")
	b.ReportMetric(frameMs/float64(b.N), "gioframe-ms/frame")
	b.ReportMetric((paintMs+frameMs)/float64(b.N), "total-ms/frame")
}

// 文本密集场景：类似文档站 —— 大量文本行、带阴影的按钮、SVG 图标。
// 同时覆盖字形缓存、阴影贴图、SVG 路径缓存与 Gio 每帧三角化的成本。
func BenchmarkFrameTextHeavy(b *testing.B) {
	row := func(i int) *Node {
		return Div(Style(Row, Gap(12), ItemsCenter, PaddingXY(16, 10), Bg(Hex("#ffffff")), Radius(8)),
			Icon(IconChevronDown, 16, TextColor(Hex("#334155"))),
			Div(Style(Column, Gap(4), Grow(1)),
				Text(fmt.Sprintf("第 %d 项：逐元素投影与数据绑定的基准验证", i), FontSize(14), Semibold),
				Text("绘制、布局、命中三线分离；每一帧的轮廓三角化成本值得单独测量", FontSize(12), TextColor(Hex("#64748b")))),
			Div(Style(Width(64), Height(28), Radius(6), Bg(Hex("#3b82f6")),
				Shadow(Hex("#00000044"), 0, 2, 8, 0))))
	}
	var kids []*Node
	kids = append(kids, Style(Column, Gap(6), Padding(12)))
	for i := 0; i < 200; i++ {
		kids = append(kids, row(i))
	}
	benchFrame(b, Div(append([]*Node{Style(Width(560), Height(900), Column, Bg(Hex("#f1f5f9")))}, kids...)...), 600, 900)
}

// 3D 场景：Scene3D 牌桌 + 20 张带文字的卡。投影层逐元素录制、四角三角函数、
// 裁剪四边形路径，每帧的投影数学与 clip 构造成本在这里测量。
func BenchmarkFrameScene3D(b *testing.B) {
	mkCard := func(i int) *Node {
		return Div(Style(Width(78), Height(104), Radius(6), Bg(Hex("#f8fafc")),
			Border(1, Hex("#cbd5e1")), Column, ItemsCenter, PaddingXY(4, 8), Gap(4)),
			Icon(IconChevronDown, 18, TextColor(Hex("#475569"))),
			Text(fmt.Sprintf("K%d", i+1), FontSize(11), Semibold))
	}
	var cards []*Node
	for i := 0; i < 20; i++ {
		cards = append(cards, mkCard(i))
	}
	benchFrame(b, Div(Style(Width(900), Height(700), ItemsCenter, JustifyCenter),
		Div(append([]*Node{Style(Width(760), Height(460), Scene3D, Perspective(750), RotateX(48),
			Row, Gap(10), ItemsCenter, JustifyCenter, Bg(Hex("#0f172a")), Radius(12))}, cards...)...)), 900, 700)
}

// 归因变体：同样的 200 行，每次只保留一种装饰，定位 gioframe 暴增的来源。
func benchVariant(b *testing.B, decorate func(i int) *Node) {
	row := func(i int) *Node {
		return Div(Style(Row, Gap(12), ItemsCenter, PaddingXY(16, 10)),
			decorate(i),
			Div(Style(Column, Gap(4)),
				Text(fmt.Sprintf("第 %d 项：归因变体", i), FontSize(14)),
				Text("对照组文本", FontSize(12), TextColor(Hex("#64748b")))))
	}
	var kids []*Node
	for i := 0; i < 200; i++ {
		kids = append(kids, row(i))
	}
	benchFrame(b, Div(append([]*Node{Style(Width(560), Height(900), Column)}, kids...)...), 600, 900)
}

func BenchmarkFrameIconsOnly(b *testing.B) {
	benchVariant(b, func(int) *Node {
		return Icon(IconChevronDown, 16, TextColor(Hex("#334155")))
	})
}

func BenchmarkFrameShadowsOnly(b *testing.B) {
	benchVariant(b, func(int) *Node {
		return Div(Style(Width(64), Height(28), Radius(6), Bg(Hex("#3b82f6")),
			Shadow(Hex("#00000044"), 0, 2, 8, 0)))
	})
}

func BenchmarkFrameRRectOnly(b *testing.B) {
	benchVariant(b, func(int) *Node {
		return Div(Style(Width(64), Height(28), Radius(6), Bg(Hex("#3b82f6"))))
	})
}

// 静态文本场景的对照：同样行数但去掉阴影/图标，分离出阴影与图标路径的成本占比。
func BenchmarkFrameTextPlain(b *testing.B) {
	row := func(i int) *Node {
		return Div(Style(Row, Gap(12), ItemsCenter, PaddingXY(16, 10)),
			Div(Style(Column, Gap(4)),
				Text(fmt.Sprintf("第 %d 项：纯文本行对照组", i), FontSize(14)),
				Text("无阴影、无图标、无背景 —— 只剩字形轮廓三角化", FontSize(12), TextColor(Hex("#64748b")))))
	}
	var kids []*Node
	kids = append(kids, Style(Column, Gap(6), Padding(12)))
	for i := 0; i < 200; i++ {
		kids = append(kids, row(i))
	}
	benchFrame(b, Div(append([]*Node{Style(Width(560), Height(900), Column)}, kids...)...), 600, 900)
}
