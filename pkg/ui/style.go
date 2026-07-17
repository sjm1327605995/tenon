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
	wrap    yoga.Wrap
	content yoga.Align // align-content：多行时行与行之间怎么排（仅换行时有意义）
	hasDir  bool
	hasAl   bool
	hasJu   bool
	hasWrap bool
	hasCont bool

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

	// 变换（围绕自身中心，2D）
	scale          float32
	rotate         float32 // 角度（绕 Z 轴，2D 旋转）
	transX, transY float32

	// 伪 3D（仅影响绘制，不参与布局）：绕 X/Y 轴旋转 + Z 位移 + 透视，
	// 对应 CSS 的 transform: perspective(p) rotateX(rx) rotateY(ry) translateZ(tz)。
	rotateX, rotateY float32 // 角度
	transZ           float32 // Z 位移（正值朝观察者，需配合 Perspective 才有远近感）
	perspective      float32 // 透视距离（px，越小透视越强；0=正交无透视）
	scene3D          bool    // 作为共享相机的场景根（见 Scene3D）

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

// ---- 换行（flex-wrap）----
//
// 默认不换行：子元素挤在一行里被压缩。Wrap 让放不下的子元素折到下一行，
// 用于图片墙、标签云、卡组编辑器一类「装得下多少算多少」的布局。

// Wrap 允许子元素折行（flex-wrap: wrap）。
func Wrap(s *StyleProps) { s.wrap, s.hasWrap = yoga.WrapWrap, true }

// WrapReverse 折行且交叉轴方向翻转（flex-wrap: wrap-reverse）：新行往上/往左堆。
func WrapReverse(s *StyleProps) { s.wrap, s.hasWrap = yoga.WrapWrapReverse, true }

// NoWrap 显式禁止折行（flex-wrap: nowrap），即默认值。
func NoWrap(s *StyleProps) { s.wrap, s.hasWrap = yoga.WrapNoWrap, true }

// ---- 多行的对齐（align-content）----
//
// 只在换行后有多行时才有意义：它管的是「行与行」在交叉轴上怎么排，
// 而 ItemsCenter 那组（align-items）管的是「一行内的子元素」怎么对齐。

func ContentStart(s *StyleProps)   { s.content, s.hasCont = yoga.AlignFlexStart, true }
func ContentCenter(s *StyleProps)  { s.content, s.hasCont = yoga.AlignCenter, true }
func ContentEnd(s *StyleProps)     { s.content, s.hasCont = yoga.AlignFlexEnd, true }
func ContentStretch(s *StyleProps) { s.content, s.hasCont = yoga.AlignStretch, true }
func ContentBetween(s *StyleProps) { s.content, s.hasCont = yoga.AlignSpaceBetween, true }
func ContentAround(s *StyleProps)  { s.content, s.hasCont = yoga.AlignSpaceAround, true }

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

// ---- 伪 3D 变换（只影响绘制，不参与 yoga 布局，因此不会改变排版）----
//
// 与 CSS 的 transform: perspective(p) rotateX(rx) rotateY(ry) translateZ(tz) 对应。
// 这不是真 3D：内容仍是平面，只是把它按透视投影贴回屏幕（做倾斜卡片一类的效果）。

// RotateX 绕水平 X 轴旋转（度），产生上/下翻起的立体感。需配合 Perspective。
func RotateX(deg float32) StyleOpt { return func(s *StyleProps) { s.rotateX = deg } }

// RotateY 绕垂直 Y 轴旋转（度），产生左/右翻转的立体感。需配合 Perspective。
func RotateY(deg float32) StyleOpt { return func(s *StyleProps) { s.rotateY = deg } }

// TranslateZ 沿 Z 轴平移（px，正值朝观察者）；需配合 Perspective 才有远近缩放。
func TranslateZ(z float32) StyleOpt { return func(s *StyleProps) { s.transZ = z } }

// Perspective 设置本元素 3D 变换的透视距离（px，越小透视越强；0=正交无透视）。
// 透视锚定在元素中心，等价于 CSS transform 里的 perspective(px)。
func Perspective(px float32) StyleOpt { return func(s *StyleProps) { s.perspective = px } }

// Scene3D 把本元素变成一台共享相机：它自己的 Perspective/RotateX/RotateY 定义了一个
// 倾斜的平面与视角，而它的每个直接子元素各自投影到这个平面上、共用同一个灭点。
//
// 这解决的是「摆一桌卡」的问题。不加 Scene3D 时每个元素绕自己的中心做透视，
// 灭点各不相同，并排摆开就不像同一张桌子；而把整块桌子当一个元素来倾斜也不行 ——
// 投影是按元素尺寸近似的，元素越大越不准（见 gio_3d.go 的「已知边界」）。
//
// 用法：桌子加 Scene3D + Perspective + RotateX，卡牌作为直接子元素按布局摆放即可，
// 它们无需自己声明 3D。子元素自身的 RotateX/RotateY 会叠加到相机的角度上
// （欧拉角相加，不是严格的矩阵复合；对「桌上的卡再翻一下」这类用法足够）。
// TranslateZ 照常生效，用来让卡片浮离桌面。
//
// 边界：
//   - 只对直接子元素生效 —— 中间再套一层容器，那层会被当作一个整体投影，尺寸一大就失真。
//     卡牌请直接挂在场景下（决斗盘那种按区绝对定位正合适）。
//   - 场景自身的背景同样受「大元素失真」约束（见 gio_3d.go），陡角度下会被画成斜切的
//     平行四边形而不是梯形。子元素只要够小就不受影响。
//   - 场景自身的 Clip 在 3D 下不生效（裁剪矩形是未投影的，会切错）。
func Scene3D(s *StyleProps) { s.scene3D = true }

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
