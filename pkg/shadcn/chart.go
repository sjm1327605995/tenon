package shadcn

import (
	"fmt"
	"math"
	"strings"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

type BarChartProps struct {
	Data          []float32
	Labels        []string
	Width, Height float32
}

// BarChart 是简单柱状图。
func BarChart(p BarChartProps) *ui.Node { return ui.Use(barChart, p) }

func barChart(p BarChartProps) *ui.Node {
	th := ui.UseTheme()
	maxV := float32(1)
	for _, v := range p.Data {
		if v > maxV {
			maxV = v
		}
	}
	bars := []*ui.Node{ui.Style(ui.Row, ui.Gap(10), ui.Width(p.Width), ui.Height(p.Height))}
	for i, v := range p.Data {
		h := v / maxV * (p.Height - 26)
		label := ""
		if i < len(p.Labels) {
			label = p.Labels[i]
		}
		bars = append(bars, ui.Div(
			ui.Style(ui.Column, ui.Grow(1), ui.JustifyEnd, ui.ItemsCenter, ui.Gap(4)),
			ui.Div(ui.Style(ui.Width(26), ui.Height(h), ui.Radius(4), ui.Bg(th.Primary))),
			ui.Text(label, ui.FontSize(11), ui.TextColor(th.MutedForeground)),
		))
	}
	return ui.Div(bars...)
}

// ---- 折线 / 面积图 ----

type LineChartProps struct {
	Data          []float32
	Width, Height float32
	Area          bool // 折线下方填充
}

// LineChart 折线图；LineChartProps.Area=true 时为面积图。
func LineChart(p LineChartProps) *ui.Node { return ui.Use(lineChart, p) }

func lineChart(p LineChartProps) *ui.Node {
	th := ui.UseTheme()
	w, h := p.Width, p.Height
	if w <= 0 {
		w = 320
	}
	if h <= 0 {
		h = 160
	}
	if len(p.Data) < 2 {
		return ui.Div(ui.Style(ui.Width(w), ui.Height(h)))
	}
	pad := float32(6)
	pts := points(p.Data, w, h, pad)

	kids := []*ui.Node{ui.Style(ui.Width(w), ui.Height(h))}
	if p.Area {
		// 闭合到底边形成面积
		area := "M" + pts[0] + " L" + strings.Join(pts, " L") +
			fmt.Sprintf(" L%.1f %.1f L%.1f %.1f Z", w-pad, h-pad, pad, h-pad)
		kids = append(kids, ui.Div(ui.Style(ui.Absolute, ui.Left(0), ui.Top(0), ui.Opacity(0.18)),
			ui.Vector(area, w, h, 0, ui.TextColor(th.Primary))))
	}
	line := "M" + pts[0] + " L" + strings.Join(pts[1:], " L")
	kids = append(kids, ui.Div(ui.Style(ui.Absolute, ui.Left(0), ui.Top(0)),
		ui.Vector(line, w, h, 2, ui.TextColor(th.Primary))))
	return ui.Div(kids...)
}

// points 把数据映射到 [pad, w-pad]×[pad, h-pad] 的坐标串。
func points(data []float32, w, h, pad float32) []string {
	minV, maxV := data[0], data[0]
	for _, v := range data {
		minV = float32(math.Min(float64(minV), float64(v)))
		maxV = float32(math.Max(float64(maxV), float64(v)))
	}
	span := maxV - minV
	if span == 0 {
		span = 1
	}
	iw, ih := w-2*pad, h-2*pad
	pts := make([]string, len(data))
	for i, v := range data {
		x := pad + iw*float32(i)/float32(len(data)-1)
		y := pad + ih*(1-(v-minV)/span)
		pts[i] = fmt.Sprintf("%.1f %.1f", x, y)
	}
	return pts
}

// ---- 饼图 / 环形图 ----

type PieSlice struct {
	Value float32
	Color ui.Color
	Label string
}

type PieChartProps struct {
	Slices []PieSlice
	Size   float32
	Donut  bool // 环形（暂以饼图近似）
}

// PieChart 饼图：各扇区按值占比，颜色可自定义。
func PieChart(p PieChartProps) *ui.Node { return ui.Use(pieChart, p) }

func pieChart(p PieChartProps) *ui.Node {
	sz := p.Size
	if sz <= 0 {
		sz = 180
	}
	var total float32
	for _, s := range p.Slices {
		total += s.Value
	}
	if total <= 0 {
		return ui.Div(ui.Style(ui.Width(sz), ui.Height(sz)))
	}
	cx, cy, r := sz/2, sz/2, sz/2-2
	kids := []*ui.Node{ui.Style(ui.Width(sz), ui.Height(sz))}
	ang := -math.Pi / 2 // 从顶部开始
	for _, s := range p.Slices {
		frac := float64(s.Value / total)
		next := ang + frac*2*math.Pi
		x1, y1 := cx+r*float32(math.Cos(ang)), cy+r*float32(math.Sin(ang))
		x2, y2 := cx+r*float32(math.Cos(next)), cy+r*float32(math.Sin(next))
		large := 0
		if frac > 0.5 {
			large = 1
		}
		d := fmt.Sprintf("M%.1f %.1f L%.1f %.1f A%.1f %.1f 0 %d 1 %.1f %.1f Z",
			cx, cy, x1, y1, r, r, large, x2, y2)
		kids = append(kids, ui.Div(ui.Style(ui.Absolute, ui.Left(0), ui.Top(0)),
			ui.Vector(d, sz, sz, 0, ui.TextColor(s.Color))))
		ang = next
	}
	return ui.Div(kids...)
}
