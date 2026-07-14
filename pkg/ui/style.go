package ui

import (
	"math"

	"github.com/sjm1327605995/tenon/yoga"
)

// StyleProps 汇聚布局与外观属性，映射到 yoga 与绘制。
// 尺寸类字段用 NaN 表示“未设置/auto”，以区分显式的 0。
type StyleProps struct {
	width, height          float32
	widthPct, heightPct    float32 // 百分比尺寸（相对父容器/视口），NaN 表示未设置
	minW, minH, maxW, maxH float32

	padT, padR, padB, padL float32
	marT, marR, marB, marL float32

	gap    float32
	grow   float32
	shrink float32

	dir     yoga.FlexDirection
	align   yoga.Align
	justify yoga.Justify
	hasDir  bool
	hasAl   bool
	hasJu   bool

	bg          Color
	gradFrom    Color
	gradTo      Color
	gradAngle   float32
	hasGradient bool
	radius      float32
	borderW     float32
	borderColor Color
	clip        bool

	// 定位
	absolute               bool
	posT, posR, posB, posL float32

	opacity float32

	// 变换（围绕自身中心）
	scale          float32
	rotate         float32 // 角度
	transX, transY float32

	// 阴影（box-shadow）
	shadowColor              Color
	shadowX, shadowY         float32
	shadowBlur, shadowSpread float32
	hasShadow                bool

	animateLayout bool

	// 文本
	color       Color
	fontSize    float32
	hasColor    bool
	hasFontSize bool
	weight      int // 400 常规 / 700 粗体
	hasWeight   bool
	italic      bool
	hasItalic   bool
}

// StyleOpt 是作用于 StyleProps 的选项。
type StyleOpt func(*StyleProps)

func newStyleProps() StyleProps {
	n := float32(math.NaN())
	return StyleProps{
		width: n, height: n, widthPct: n, heightPct: n,
		minW: n, minH: n, maxW: n, maxH: n,
		posT: n, posR: n, posB: n, posL: n,
		opacity: 1, scale: 1,
	}
}

func isNaN(f float32) bool { return f != f }

// ---- 尺寸 ----

func Width(v float32) StyleOpt  { return func(s *StyleProps) { s.width = v } }
func Height(v float32) StyleOpt { return func(s *StyleProps) { s.height = v } }

// WidthPct/HeightPct 按父容器（根为视口）的百分比设定尺寸，随窗口自适应。
func WidthPct(v float32) StyleOpt  { return func(s *StyleProps) { s.widthPct = v } }
func HeightPct(v float32) StyleOpt { return func(s *StyleProps) { s.heightPct = v } }

// Fill 让元素填满父容器（根为视口），窗口缩放时自适应。
func Fill(s *StyleProps)           { s.widthPct, s.heightPct = 100, 100 }
func MinWidth(v float32) StyleOpt  { return func(s *StyleProps) { s.minW = v } }
func MinHeight(v float32) StyleOpt { return func(s *StyleProps) { s.minH = v } }
func MaxWidth(v float32) StyleOpt  { return func(s *StyleProps) { s.maxW = v } }
func MaxHeight(v float32) StyleOpt { return func(s *StyleProps) { s.maxH = v } }

// ---- 间距 ----

func Padding(v float32) StyleOpt {
	return func(s *StyleProps) { s.padT, s.padR, s.padB, s.padL = v, v, v, v }
}
func PaddingXY(h, v float32) StyleOpt {
	return func(s *StyleProps) { s.padL, s.padR, s.padT, s.padB = h, h, v, v }
}
func Margin(v float32) StyleOpt {
	return func(s *StyleProps) { s.marT, s.marR, s.marB, s.marL = v, v, v, v }
}
func MarginXY(h, v float32) StyleOpt {
	return func(s *StyleProps) { s.marL, s.marR, s.marT, s.marB = h, h, v, v }
}
func Gap(v float32) StyleOpt { return func(s *StyleProps) { s.gap = v } }

// ---- flex ----

func Grow(v float32) StyleOpt   { return func(s *StyleProps) { s.grow = v } }
func Shrink(v float32) StyleOpt { return func(s *StyleProps) { s.shrink = v } }

func Row(s *StyleProps)    { s.dir, s.hasDir = yoga.FlexDirectionRow, true }
func Column(s *StyleProps) { s.dir, s.hasDir = yoga.FlexDirectionColumn, true }

func ItemsStart(s *StyleProps)  { s.align, s.hasAl = yoga.AlignFlexStart, true }
func ItemsCenter(s *StyleProps) { s.align, s.hasAl = yoga.AlignCenter, true }
func ItemsEnd(s *StyleProps)    { s.align, s.hasAl = yoga.AlignFlexEnd, true }

func JustifyStart(s *StyleProps)   { s.justify, s.hasJu = yoga.JustifyFlexStart, true }
func JustifyCenter(s *StyleProps)  { s.justify, s.hasJu = yoga.JustifyCenter, true }
func JustifyEnd(s *StyleProps)     { s.justify, s.hasJu = yoga.JustifyFlexEnd, true }
func JustifyBetween(s *StyleProps) { s.justify, s.hasJu = yoga.JustifySpaceBetween, true }

// ---- 外观 ----

func Bg(c Color) StyleOpt { return func(s *StyleProps) { s.bg = c } }

// LinearGradient 用线性渐变作为背景填充：颜色从 from 到 to，
// angleDeg 为方向角（0=左→右，90=上→下，45=左上→右下）。会遵循圆角。
func LinearGradient(from, to Color, angleDeg float32) StyleOpt {
	return func(s *StyleProps) {
		s.gradFrom, s.gradTo, s.gradAngle, s.hasGradient = from, to, angleDeg, true
	}
}
func Radius(v float32) StyleOpt { return func(s *StyleProps) { s.radius = v } }
func Border(w float32, c Color) StyleOpt {
	return func(s *StyleProps) { s.borderW, s.borderColor = w, c }
}

// Clip 裁剪超出自身边界的子内容（overflow: hidden）。
func Clip(s *StyleProps) { s.clip = true }

// Opacity 设置不透明度（0..1）。叶子节点作用于自身；有子节点时作为整组透明度。
func Opacity(v float32) StyleOpt { return func(s *StyleProps) { s.opacity = v } }

// ---- 变换（围绕元素中心）----

func Scale(v float32) StyleOpt    { return func(s *StyleProps) { s.scale = v } }
func Rotate(deg float32) StyleOpt { return func(s *StyleProps) { s.rotate = deg } }
func TranslateXY(x, y float32) StyleOpt {
	return func(s *StyleProps) { s.transX, s.transY = x, y }
}

// Shadow 设置投影（box-shadow）：颜色、水平/垂直偏移、模糊半径、扩散。offY 正值向下。
// 叶子与容器均可用；柔和边缘由分层近似实现。
func Shadow(c Color, offX, offY, blur, spread float32) StyleOpt {
	return func(s *StyleProps) {
		s.shadowColor = c
		s.shadowX, s.shadowY = offX, offY
		s.shadowBlur, s.shadowSpread = blur, spread
		s.hasShadow = c.A > 0
	}
}

// Animated 开启布局动画：当该元素的布局位置变化时，从旧位置平滑滑到新位置（FLIP）。
func Animated(s *StyleProps) { s.animateLayout = true }

// ---- 定位 ----

// Absolute 使元素脱离流，按 Top/Left/Right/Bottom 相对父容器定位。
func Absolute(s *StyleProps)    { s.absolute = true }
func Top(v float32) StyleOpt    { return func(s *StyleProps) { s.posT = v } }
func Right(v float32) StyleOpt  { return func(s *StyleProps) { s.posR = v } }
func Bottom(v float32) StyleOpt { return func(s *StyleProps) { s.posB = v } }
func Left(v float32) StyleOpt   { return func(s *StyleProps) { s.posL = v } }

// ---- 文本 ----

func TextColor(c Color) StyleOpt { return func(s *StyleProps) { s.color, s.hasColor = c, true } }
func FontSize(v float32) StyleOpt {
	return func(s *StyleProps) { s.fontSize, s.hasFontSize = v, true }
}

// FontWeight 设置字重（400 常规、500 中等、600 半粗、700 粗体…）；>=600 时启用（合成）粗体。
func FontWeight(w int) StyleOpt { return func(s *StyleProps) { s.weight, s.hasWeight = w, true } }

// Bold / Medium / Semibold 是常用字重的便捷别名。
func Bold(s *StyleProps)     { s.weight, s.hasWeight = 700, true }
func Semibold(s *StyleProps) { s.weight, s.hasWeight = 600, true }
func Medium(s *StyleProps)   { s.weight, s.hasWeight = 500, true }

// Italic 启用（合成）斜体。
func Italic(s *StyleProps) { s.italic, s.hasItalic = true, true }

// ---- 组合 ----

// Styles 把多个样式选项合成一个（便于把变体/尺寸定义为可复用的单个 StyleOpt）。
func Styles(opts ...StyleOpt) StyleOpt {
	return func(s *StyleProps) {
		for _, o := range opts {
			o(s)
		}
	}
}

// StyleIf 条件应用一组样式（cond 为真才生效），用于 hover/pressed/disabled 等状态样式。
func StyleIf(cond bool, opts ...StyleOpt) StyleOpt {
	return func(s *StyleProps) {
		if cond {
			for _, o := range opts {
				o(s)
			}
		}
	}
}
