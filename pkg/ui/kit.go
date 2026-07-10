package ui

import "math"

// 基础组件套件：基于 Div/Text/Button/Input 等原语组合而成的常用控件。

// 默认配色（可按需在各组件参数里覆盖）。
var (
	kitAccent = Hex("#3b82f6")
	kitBorder = Hex("#cbd5e1")
	kitTrack  = Hex("#e2e8f0")
	kitMuted  = Hex("#94a3b8")
	kitLine   = Hex("#e5e7eb")
)

func pickColor(cond bool, a, b Color) Color {
	if cond {
		return a
	}
	return b
}

func clampf(v, lo, hi float32) float32 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func lerpColor(a, b Color, t float32) Color { return Mix(a, b, t) }

// Checkbox 是受控复选框。
func Checkbox(checked bool, onChange func(bool)) *Node {
	return Div(
		Style(Width(20), Height(20), Radius(5), ItemsCenter, JustifyCenter,
			Bg(pickColor(checked, kitAccent, White)),
			Border(2, pickColor(checked, kitAccent, kitBorder))),
		OnClick(func() {
			if onChange != nil {
				onChange(!checked)
			}
		}),
		If(checked, Text("✓", FontSize(14), TextColor(White))),
	)
}

// Radio 是受控单选项（点击触发 onChange，选中态由外部管理）。
func Radio(selected bool, onChange func()) *Node {
	return Div(
		Style(Width(20), Height(20), Radius(10), ItemsCenter, JustifyCenter,
			Border(2, pickColor(selected, kitAccent, kitBorder))),
		OnClick(func() {
			if onChange != nil {
				onChange()
			}
		}),
		If(selected, Div(Style(Width(10), Height(10), Radius(5), Bg(kitAccent)))),
	)
}

// Switch 是受控开关，滑块带过渡动画。
type switchProps struct {
	On       bool
	OnChange func(bool)
}

func switchImpl(p switchProps) *Node {
	target := float32(0)
	if p.On {
		target = 1
	}
	x := UseTween(target, 140, EaseOut) // 0..1
	return Div(
		Style(Width(44), Height(24), Radius(12), JustifyStart, ItemsCenter,
			Bg(lerpColor(kitBorder, kitAccent, x))),
		OnClick(func() {
			if p.OnChange != nil {
				p.OnChange(!p.On)
			}
		}),
		Div(Style(Width(18), Height(18), Radius(9), Bg(White),
			Absolute, Top(3), Left(3+x*20))),
	)
}

func Switch(on bool, onChange func(bool)) *Node {
	return Use(switchImpl, switchProps{On: on, OnChange: onChange})
}

// Slider 是受控滑块（拖动滑块调整值）。
type sliderProps struct {
	Value, Min, Max float32
	OnChange        func(float32)
}

func sliderImpl(p sliderProps) *Node {
	w := float32(200)
	span := p.Max - p.Min
	if span <= 0 {
		span = 1
	}
	frac := clampf((p.Value-p.Min)/span, 0, 1)
	fill := frac * w
	return Div(
		Style(Width(w), Height(20), JustifyStart, ItemsCenter),
		Div(Style(Width(w), Height(6), Radius(3), Bg(kitTrack), Absolute, Top(7))),
		Div(Style(Width(fill), Height(6), Radius(3), Bg(kitAccent), Absolute, Top(7))),
		Div(
			Style(Width(18), Height(18), Radius(9), Bg(White), Border(2, kitAccent),
				Absolute, Top(1), Left(fill-9)),
			OnDrag(func(dx, _ float32) {
				if p.OnChange != nil {
					p.OnChange(clampf(p.Value+dx/w*span, p.Min, p.Max))
				}
			}),
		),
	)
}

func Slider(value, min, max float32, onChange func(float32)) *Node {
	return Use(sliderImpl, sliderProps{Value: value, Min: min, Max: max, OnChange: onChange})
}

// ProgressBar 展示 0..1 的进度。
func ProgressBar(value float32) *Node {
	value = clampf(value, 0, 1)
	w := float32(220)
	return Div(
		Style(Width(w), Height(8), Radius(4), Bg(kitTrack)),
		Div(Style(Width(w*value), Height(8), Radius(4), Bg(kitAccent))),
	)
}

// Spinner 是持续旋转的加载指示器。
type spinnerProps struct {
	Size  float32
	Color Color
}

func spinnerImpl(p spinnerProps) *Node {
	deg := UseElapsed() * 300 // ~0.83 圈/秒
	sz := p.Size
	dot := sz / 5
	r := sz/2 - dot/2
	kids := []*Node{Style(Width(sz), Height(sz), Rotate(deg))}
	const n = 8
	for i := 0; i < n; i++ {
		a := float64(i) / float64(n) * 2 * math.Pi
		cx := sz/2 + float32(math.Cos(a))*r - dot/2
		cy := sz/2 + float32(math.Sin(a))*r - dot/2
		kids = append(kids, Div(Style(
			Width(dot), Height(dot), Radius(dot/2), Bg(p.Color),
			Opacity(float32(i+1)/float32(n)), Absolute, Left(cx), Top(cy))))
	}
	return Div(kids...)
}

func Spinner(size float32, c Color) *Node {
	return Use(spinnerImpl, spinnerProps{Size: size, Color: c})
}

// Badge 是小圆角标签。
func Badge(text string, c Color) *Node {
	return Div(
		Style(PaddingXY(10, 3), Radius(999), Bg(c), ItemsCenter, JustifyCenter),
		Text(text, FontSize(12), TextColor(White)),
	)
}

// Avatar 是显示首字母的圆形头像。
func Avatar(initials string, size float32) *Node {
	return Div(
		Style(Width(size), Height(size), Radius(size/2), Bg(kitAccent),
			ItemsCenter, JustifyCenter),
		Text(initials, FontSize(size*0.4), TextColor(White)),
	)
}

// Divider 是一条横向分隔线（在会拉伸交叉轴的容器中铺满宽度）。
func Divider() *Node {
	return Div(Style(Height(1), Bg(kitLine)))
}

// Card 是带边框内边距的卡片容器；传入 Style(...) 可覆盖默认样式。
func Card(children ...*Node) *Node {
	base := Style(Bg(White), Radius(12), Border(1, kitLine), Padding(16), Column, Gap(8))
	return Div(append([]*Node{base}, children...)...)
}

// Tabs 渲染标签栏（内容区由外部按 Active 自行渲染）。
type TabsProps struct {
	Tabs     []string
	Active   int
	OnChange func(int)
}

func Tabs(p TabsProps) *Node {
	kids := []*Node{Style(Row, Gap(4))}
	for i, label := range p.Tabs {
		active := i == p.Active
		idx := i
		kids = append(kids, Button(
			Style(PaddingXY(16, 8), Radius(8), ItemsCenter, JustifyCenter,
				Bg(pickColor(active, kitAccent, Transparent))),
			OnClick(func() {
				if p.OnChange != nil {
					p.OnChange(idx)
				}
			}),
			Text(label, FontSize(14), TextColor(pickColor(active, White, kitMuted))),
		))
	}
	return Div(kids...)
}
