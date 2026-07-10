package ui

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"time"
	"unicode/utf8"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/sjm1327605995/tenon/pkg/font"
	"github.com/sjm1327605995/tenon/yoga"
)

// Rect 是绝对坐标下的矩形。
type Rect struct{ X, Y, W, H float32 }

func (r Rect) contains(px, py float32) bool {
	return px >= r.X && px < r.X+r.W && py >= r.Y && py < r.Y+r.H
}

type rnKind int

const (
	rnBox rnKind = iota
	rnText
	rnInput
	rnImage
	rnScroll
)

// renderNode = yoga.Node + 绘制数据。是唯一进入 yoga 树的节点类型。
type renderNode struct {
	yn       *yoga.Node
	kind     rnKind
	bounds   Rect
	owner    *Fiber
	parent   *renderNode
	children []*renderNode

	// box / input
	bg          Color
	radius      float32
	borderW     float32
	borderColor Color
	onClick     func()

	// text / input
	text  string
	face  *text.GoTextFace
	color Color
	lineH float64

	// 文本样式继承
	explicitColor bool // 本节点显式设置了颜色
	ownColor      Color
	explicitSize  bool // 本节点显式设置了字号
	ownSize       float32
	effSize       float32 // 生效字号（逻辑，含继承）
	effScale      float32 // 上次取字体所用的 uiScale
	inhColor      Color   // box 向下传递的颜色
	hasInhColor   bool
	inhSize       float32 // box 向下传递的字号
	hasInhSize    bool

	// input
	value       string
	placeholder string
	multiline   bool
	onChange    func(string)
	caretPos    int
	selAnchor   int // 选区另一端；等于 caretPos 表示无选区
	onHover     func(bool)
	onPress     func(bool)
	onDrag      func(dx, dy float32)
	measure     *measureHook
	focusable   bool

	// image
	imgSrc string
	img    *ebiten.Image

	// scroll / clip
	clip     bool
	scroll   bool
	scrollY  float32
	contentH float32

	opacity        float32
	scale          float32
	rotate         float32
	transX, transY float32

	// 投影（box-shadow）
	shadowColor              Color
	shadowX, shadowY         float32
	shadowBlur, shadowSpread float32
	hasShadow                bool

	// 布局动画（FLIP）：位置变化时的残余偏移，逐帧衰减到 0
	animatedLayout     bool
	hasPrevPos         bool
	prevPosX, prevPosY float32
	offX, offY         float32
}

func (rn *renderNode) effTransX() float32 { return rn.transX + rn.offX }
func (rn *renderNode) effTransY() float32 { return rn.transY + rn.offY }

func (rn *renderNode) hasTransform() bool {
	return rn.scale != 1 || rn.rotate != 0 || rn.effTransX() != 0 || rn.effTransY() != 0
}

func (rn *renderNode) needsLayer() bool {
	return rn.hasTransform() || (rn.opacity < 1 && len(rn.children) > 0)
}

func (rn *renderNode) container() bool { return rn.kind == rnBox || rn.kind == rnScroll }

func newHostRenderNode(tag string) *renderNode {
	switch tag {
	case "input":
		return newInputRenderNode()
	case "img":
		return newImageRenderNode()
	case "scroll":
		return &renderNode{yn: yoga.NewNode(), kind: rnScroll, clip: true, scroll: true, opacity: 1, scale: 1}
	default:
		return newBoxRenderNode()
	}
}

func newBoxRenderNode() *renderNode {
	return &renderNode{yn: yoga.NewNode(), kind: rnBox, opacity: 1, scale: 1}
}

func newTextRenderNode(s string, st StyleProps) *renderNode {
	rn := &renderNode{yn: yoga.NewNode(), kind: rnText, text: s, opacity: 1, scale: 1}
	rn.applyTextStyle(st)
	rn.yn.SetMeasureFunc(func(_ *yoga.Node, w float32, wm yoga.MeasureMode, _ float32, _ yoga.MeasureMode) yoga.Size {
		if rn.face == nil || rn.text == "" {
			return yoga.Size{}
		}
		avail := float32(0)
		if wm == yoga.MeasureModeExactly || wm == yoga.MeasureModeAtMost {
			avail = w
		}
		lines, mw := wrapForWidth(rn.text, rn.face, rn.lineH, avail)
		return yoga.Size{Width: mw, Height: float32(len(lines)) * float32(rn.lineH)}
	})
	return rn
}

func newInputRenderNode() *renderNode {
	rn := &renderNode{yn: yoga.NewNode(), kind: rnInput, opacity: 1, scale: 1}
	rn.applyTextStyle(newStyleProps())
	rn.yn.SetMeasureFunc(func(_ *yoga.Node, w float32, wm yoga.MeasureMode, _ float32, _ yoga.MeasureMode) yoga.Size {
		if rn.face == nil {
			return yoga.Size{Width: 40, Height: float32(rn.lineH)}
		}
		if rn.multiline {
			avail := float32(0)
			if wm == yoga.MeasureModeExactly || wm == yoga.MeasureModeAtMost {
				avail = w
			}
			s := rn.value
			if s == "" {
				s = " "
			}
			lines, mw := wrapForWidth(s, rn.face, rn.lineH, avail)
			cw := mw
			if avail > 0 {
				cw = avail
			}
			return yoga.Size{Width: cw, Height: float32(len(lines)) * float32(rn.lineH)}
		}
		wd := float32(40)
		s := rn.value
		if s == "" {
			s = rn.placeholder
		}
		if s != "" {
			tw, _ := text.Measure(s, rn.face, rn.lineH)
			wd = float32(tw) + 8
		}
		return yoga.Size{Width: wd, Height: float32(rn.lineH)}
	})
	return rn
}

func newImageRenderNode() *renderNode {
	rn := &renderNode{yn: yoga.NewNode(), kind: rnImage, opacity: 1, scale: 1}
	rn.yn.SetMeasureFunc(func(_ *yoga.Node, _ float32, _ yoga.MeasureMode, _ float32, _ yoga.MeasureMode) yoga.Size {
		if rn.img == nil {
			return yoga.Size{}
		}
		b := rn.img.Bounds()
		return yoga.Size{Width: float32(b.Dx()), Height: float32(b.Dy())}
	})
	return rn
}

// applyTextStyle 只记录本节点显式设置的文本样式；生效值由 resolveInherited 决定。
func (rn *renderNode) applyTextStyle(st StyleProps) {
	rn.explicitColor = st.hasColor
	rn.ownColor = st.color
	rn.explicitSize = st.hasFontSize
	rn.ownSize = st.fontSize
	if rn.effSize == 0 { // 初始回退，保证在 resolve 前也有可用字体
		rn.setEffectiveText(Black, 16)
	}
}

// setEffectiveText 应用最终生效的颜色/字号（字号在物理像素下取字体，保证高分屏清晰）。
func (rn *renderNode) setEffectiveText(c Color, size float32) {
	if size <= 0 {
		size = 16
	}
	rn.color = c
	if rn.effSize != size || rn.effScale != uiScale || rn.face == nil {
		rn.effSize = size
		rn.effScale = uiScale
		px := size * uiScale
		rn.lineH = float64(px) * 1.3
		if ff, err := font.GetDefaultFontFace(px); err == nil {
			rn.face = ff.Face
		}
		rn.yn.MarkDirty()
	}
}

// resolveInherited 自顶向下解析文本继承：box 贡献/透传 color 与 font-size，
// 文本/输入节点若未显式设置则采用继承值。须在测量（CalculateLayout）之前调用。
func resolveInherited(rn *renderNode, inC Color, hasC bool, inS float32, hasS bool) {
	switch rn.kind {
	case rnText, rnInput:
		c := Black
		if rn.explicitColor {
			c = rn.ownColor
		} else if hasC {
			c = inC
		}
		s := float32(16)
		if rn.explicitSize {
			s = rn.ownSize
		} else if hasS {
			s = inS
		}
		rn.setEffectiveText(c, s)
	default:
		cC, cHasC := inC, hasC
		if rn.hasInhColor {
			cC, cHasC = rn.inhColor, true
		}
		cS, cHasS := inS, hasS
		if rn.hasInhSize {
			cS, cHasS = rn.inhSize, true
		}
		for _, ch := range rn.children {
			resolveInherited(ch, cC, cHasC, cS, cHasS)
		}
	}
}

func (rn *renderNode) setText(s string, st StyleProps) {
	changed := rn.text != s
	rn.text = s
	rn.applyTextStyle(st)
	if changed {
		rn.yn.MarkDirty()
	}
}

var imgCache = map[string]*ebiten.Image{}

func (rn *renderNode) loadImage() {
	if img, ok := imgCache[rn.imgSrc]; ok {
		rn.img = img
		return
	}
	f, err := os.Open(rn.imgSrc)
	if err != nil {
		return
	}
	defer f.Close()
	src, _, err := image.Decode(f)
	if err != nil {
		return
	}
	ei := ebiten.NewImageFromImage(src)
	imgCache[rn.imgSrc] = ei
	rn.img = ei
}

// applyHostProps 把某个 host 元素的属性写入其 renderNode。
func applyHostProps(rn *renderNode, hp hostProps) {
	syncYoga(rn, hp.style)
	rn.onClick = hp.onClick
	rn.onHover = hp.onHover
	rn.onPress = hp.onPress
	rn.onDrag = hp.onDrag
	rn.measure = hp.measure
	rn.focusable = rn.kind == rnInput || hp.onClick != nil
	switch rn.kind {
	case rnInput:
		rn.applyTextStyle(hp.style)
		rn.value = hp.value
		rn.placeholder = hp.placeholder
		rn.multiline = hp.multiline
		rn.onChange = hp.onChange
		if rn.caretPos > len(rn.value) {
			rn.caretPos = len(rn.value)
		}
		if rn.selAnchor > len(rn.value) {
			rn.selAnchor = len(rn.value)
		}
		rn.yn.MarkDirty()
	case rnImage:
		if hp.src != "" && rn.imgSrc != hp.src {
			rn.imgSrc = hp.src
			rn.loadImage()
			rn.yn.MarkDirty()
		}
	}
}

// syncYoga 把 StyleProps 写进 yoga 节点，并缓存绘制属性。所有尺寸按 uiScale 换算到物理像素。
func syncYoga(rn *renderNode, s StyleProps) {
	yn := rn.yn
	k := uiScale

	if !isNaN(s.widthPct) {
		yn.StyleSetWidthPercent(s.widthPct)
	} else if isNaN(s.width) {
		yn.StyleSetWidthAuto()
	} else {
		yn.StyleSetWidth(s.width * k)
	}
	if !isNaN(s.heightPct) {
		yn.StyleSetHeightPercent(s.heightPct)
	} else if isNaN(s.height) {
		yn.StyleSetHeightAuto()
	} else {
		yn.StyleSetHeight(s.height * k)
	}
	if !isNaN(s.minW) {
		yn.StyleSetMinWidth(s.minW * k)
	}
	if !isNaN(s.minH) {
		yn.StyleSetMinHeight(s.minH * k)
	}
	if !isNaN(s.maxW) {
		yn.StyleSetMaxWidth(s.maxW * k)
	}
	if !isNaN(s.maxH) {
		yn.StyleSetMaxHeight(s.maxH * k)
	}

	yn.StyleSetPadding(yoga.EdgeTop, s.padT*k)
	yn.StyleSetPadding(yoga.EdgeRight, s.padR*k)
	yn.StyleSetPadding(yoga.EdgeBottom, s.padB*k)
	yn.StyleSetPadding(yoga.EdgeLeft, s.padL*k)

	yn.StyleSetMargin(yoga.EdgeTop, s.marT*k)
	yn.StyleSetMargin(yoga.EdgeRight, s.marR*k)
	yn.StyleSetMargin(yoga.EdgeBottom, s.marB*k)
	yn.StyleSetMargin(yoga.EdgeLeft, s.marL*k)

	yn.StyleSetGap(yoga.GutterAll, s.gap*k)
	yn.StyleSetFlexGrow(s.grow)
	yn.StyleSetFlexShrink(s.shrink)

	if s.hasDir {
		yn.StyleSetFlexDirection(s.dir)
	}
	if s.hasAl {
		yn.StyleSetAlignItems(s.align)
	}
	if s.hasJu {
		yn.StyleSetJustifyContent(s.justify)
	}
	if s.borderW > 0 {
		yn.StyleSetBorder(yoga.EdgeAll, s.borderW*k)
	}

	if s.absolute {
		yn.StyleSetPositionType(yoga.PositionTypeAbsolute)
	} else {
		yn.StyleSetPositionType(yoga.PositionTypeRelative)
	}
	setPos(yn, yoga.EdgeTop, s.posT*k)
	setPos(yn, yoga.EdgeRight, s.posR*k)
	setPos(yn, yoga.EdgeBottom, s.posB*k)
	setPos(yn, yoga.EdgeLeft, s.posL*k)

	rn.bg = s.bg
	rn.radius = s.radius * k
	rn.borderW = s.borderW * k
	rn.borderColor = s.borderColor
	rn.opacity = s.opacity
	rn.scale = s.scale
	rn.rotate = s.rotate
	rn.transX, rn.transY = s.transX*k, s.transY*k

	rn.hasShadow = s.hasShadow
	rn.shadowColor = s.shadowColor
	rn.shadowX, rn.shadowY = s.shadowX*k, s.shadowY*k
	rn.shadowBlur, rn.shadowSpread = s.shadowBlur*k, s.shadowSpread*k
	if s.clip {
		rn.clip = true
	}

	// 容器可通过 TextColor/FontSize 为后代文本设定继承值
	rn.hasInhColor = s.hasColor
	rn.inhColor = s.color
	rn.hasInhSize = s.hasFontSize
	rn.inhSize = s.fontSize

	rn.animatedLayout = s.animateLayout
	if s.animateLayout && activeGame != nil {
		activeGame.hasLayoutAnim = true
	}
}

func setPos(yn *yoga.Node, edge yoga.Edge, v float32) {
	if !isNaN(v) {
		yn.StyleSetPosition(edge, v)
	}
}

// computeBounds 自顶向下把 yoga 相对坐标累加为绝对 bounds，并应用滚动偏移。
func computeBounds(rn *renderNode, ox, oy float32) {
	x := ox + rn.yn.LayoutLeft()
	y := oy + rn.yn.LayoutTop()
	rn.bounds = Rect{X: x, Y: y, W: rn.yn.LayoutWidth(), H: rn.yn.LayoutHeight()}

	cox, coy := x, y
	if rn.scroll {
		var maxBottom float32
		for _, c := range rn.children {
			if b := c.yn.LayoutTop() + c.yn.LayoutHeight(); b > maxBottom {
				maxBottom = b
			}
		}
		rn.contentH = maxBottom
		maxScroll := rn.contentH - rn.bounds.H
		if maxScroll < 0 {
			maxScroll = 0
		}
		if rn.scrollY > maxScroll {
			rn.scrollY = maxScroll
		}
		if rn.scrollY < 0 {
			rn.scrollY = 0
		}
		coy -= rn.scrollY
	}
	for _, c := range rn.children {
		computeBounds(c, cox, coy)
	}
}

// paint 是绘制入口；需要变换/整组透明度时先合成到离屏图层再变换绘制。
func paint(dst *ebiten.Image, rn *renderNode) {
	if rn.needsLayer() {
		paintLayer(dst, rn)
		return
	}
	paintNode(dst, rn)
}

// paintLayer 把子树合成到全屏尺寸的离屏图层，再围绕中心做 scale/rotate/translate
// 及整组透明度，最后绘制回 dst（SubImage 裁剪仍然有效）。
func paintLayer(dst *ebiten.Image, rn *renderNode) {
	if activeGame == nil {
		paintNode(dst, rn)
		return
	}
	layer := acquireLayer(activeGame.w, activeGame.h)
	o := rn.opacity
	rn.opacity = 1 // 组透明度在合成时统一施加，内部按不透明绘制
	paintNode(layer, rn)
	rn.opacity = o

	b := rn.bounds
	cx, cy := float64(b.X+b.W/2), float64(b.Y+b.H/2)
	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
	op.GeoM.Translate(-cx, -cy)
	op.GeoM.Scale(float64(rn.scale), float64(rn.scale))
	op.GeoM.Rotate(float64(rn.rotate) * math.Pi / 180)
	op.GeoM.Translate(cx+float64(rn.effTransX()), cy+float64(rn.effTransY()))
	op.ColorScale.ScaleAlpha(o)
	dst.DrawImage(layer, op)
	releaseLayer(layer)
}

func paintNode(dst *ebiten.Image, rn *renderNode) {
	b := rn.bounds
	o := rn.opacity
	if rn.hasShadow {
		fillShadow(dst, b.X, b.Y, b.W, b.H, rn.radius,
			rn.shadowX, rn.shadowY, rn.shadowBlur, rn.shadowSpread, rn.shadowColor.Alpha(o))
	}
	switch rn.kind {
	case rnBox, rnScroll:
		fillRoundRect(dst, b.X, b.Y, b.W, b.H, rn.radius, rn.bg.Alpha(o))
		if rn.borderW > 0 {
			strokeRoundRect(dst, b.X, b.Y, b.W, b.H, rn.radius, rn.borderW, rn.borderColor.Alpha(o))
		}
	case rnInput:
		paintInput(dst, rn)
	case rnText:
		lines, _ := wrapForWidth(rn.text, rn.face, rn.lineH, b.W)
		for i, ln := range lines {
			drawText(dst, ln, rn.face, rn.color.Alpha(o), b.X, b.Y+float32(i)*float32(rn.lineH))
		}
	case rnImage:
		if rn.img != nil {
			ib := rn.img.Bounds()
			op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
			op.GeoM.Scale(float64(b.W)/float64(ib.Dx()), float64(b.H)/float64(ib.Dy()))
			op.GeoM.Translate(float64(b.X), float64(b.Y))
			op.ColorScale.ScaleAlpha(o)
			dst.DrawImage(rn.img, op)
		}
	}

	// 裁剪：把子节点画进自身矩形的 SubImage（共享坐标系，自动裁掉越界部分）。
	childDst := dst
	if rn.clip {
		ir := image.Rect(int(b.X), int(b.Y), int(b.X+b.W), int(b.Y+b.H))
		if ir.Dx() > 0 && ir.Dy() > 0 {
			childDst = dst.SubImage(ir).(*ebiten.Image)
		}
	}
	for _, c := range rn.children {
		paint(childDst, c)
	}

	if rn.scroll && rn.contentH > b.H {
		drawScrollbar(dst, rn)
	}

	// 键盘焦点环
	if rn.focusable && isFocused(rn) {
		strokeRoundRect(dst, b.X-2, b.Y-2, b.W+4, b.H+4, rn.radius+2, 2, Hex("#60a5fa"))
	}
}

// 离屏图层池（按屏幕尺寸复用）。
var layerPool []*ebiten.Image

func acquireLayer(w, h int) *ebiten.Image {
	for i, im := range layerPool {
		if b := im.Bounds(); b.Dx() == w && b.Dy() == h {
			layerPool = append(layerPool[:i], layerPool[i+1:]...)
			im.Clear()
			return im
		}
	}
	return ebiten.NewImage(w, h)
}

func releaseLayer(im *ebiten.Image) { layerPool = append(layerPool, im) }

func drawScrollbar(dst *ebiten.Image, rn *renderNode) {
	b := rn.bounds
	track := b.H
	thumb := track * b.H / rn.contentH
	if thumb < 24 {
		thumb = 24
	}
	maxScroll := rn.contentH - b.H
	var t float32
	if maxScroll > 0 {
		t = rn.scrollY / maxScroll
	}
	ty := b.Y + t*(track-thumb)
	fillRoundRect(dst, b.X+b.W-6, ty, 4, thumb, 2, Color{0, 0, 0, 90})
}

func paintInput(dst *ebiten.Image, rn *renderNode) {
	b := rn.bounds
	fillRoundRect(dst, b.X, b.Y, b.W, b.H, rn.radius, rn.bg)
	if rn.borderW > 0 {
		strokeRoundRect(dst, b.X, b.Y, b.W, b.H, rn.radius, rn.borderW, rn.borderColor)
	}
	padL := rn.yn.LayoutPadding(yoga.EdgeLeft)
	padT := rn.yn.LayoutPadding(yoga.EdgeTop)
	tx := b.X + padL
	usePlaceholder := rn.value == "" && rn.placeholder != ""
	lineH := float32(rn.lineH)

	if rn.multiline {
		ty := b.Y + padT
		display := rn.value
		col := rn.color
		if usePlaceholder {
			display, col = rn.placeholder, Gray
		}
		lines, _ := wrapForWidth(display, rn.face, rn.lineH, b.W-padL*2)
		for i, ln := range lines {
			drawText(dst, ln, rn.face, col, tx, ty+float32(i)*lineH)
		}
		if isFocused(rn) && caretVisible() && !usePlaceholder {
			prefix := rn.value[:min(rn.caretPos, len(rn.value))]
			plines, _ := wrapForWidth(prefix, rn.face, rn.lineH, b.W-padL*2)
			row := len(plines) - 1
			cx := tx
			if row >= 0 {
				w, _ := text.Measure(plines[row], rn.face, rn.lineH)
				cx += float32(w)
			} else {
				row = 0
			}
			cy := ty + float32(row)*lineH
			ebitenDrawLine(dst, cx, cy, cx, cy+lineH, rn.color)
		}
		return
	}

	ty := b.Y + (b.H-lineH)/2

	// 选区高亮（单行）
	if rn.face != nil && rn.selAnchor != rn.caretPos {
		a, c := min(rn.selAnchor, rn.caretPos), max(rn.selAnchor, rn.caretPos)
		a = clampi(a, 0, len(rn.value))
		c = clampi(c, 0, len(rn.value))
		wa, _ := text.Measure(rn.value[:a], rn.face, rn.lineH)
		wc, _ := text.Measure(rn.value[:c], rn.face, rn.lineH)
		fillRoundRect(dst, tx+float32(wa), ty, float32(wc-wa), lineH, 0, Color{R: 59, G: 130, B: 246, A: 90})
	}

	if usePlaceholder {
		drawText(dst, rn.placeholder, rn.face, Gray, tx, ty)
	} else {
		drawText(dst, rn.value, rn.face, rn.color, tx, ty)
	}
	if isFocused(rn) && caretVisible() {
		cx := tx
		if rn.face != nil && rn.caretPos > 0 {
			w, _ := text.Measure(rn.value[:min(rn.caretPos, len(rn.value))], rn.face, rn.lineH)
			cx += float32(w)
		}
		ebitenDrawLine(dst, cx, ty, cx, ty+lineH, rn.color)
	}
}

func drawText(dst *ebiten.Image, s string, face *text.GoTextFace, c Color, x, y float32) {
	if face == nil || s == "" {
		return
	}
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(c)
	text.Draw(dst, s, face, op)
}

func caretVisible() bool { return (time.Now().UnixMilli()/500)%2 == 0 }

func clampi(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// caretFromX 把绝对屏幕 x 映射到单行输入内最近的字节偏移（在 rune 边界上）。
func (rn *renderNode) caretFromX(px float32) int {
	if rn.face == nil || rn.value == "" {
		return 0
	}
	padL := rn.yn.LayoutPadding(yoga.EdgeLeft)
	rel := px - (rn.bounds.X + padL)
	if rel <= 0 {
		return 0
	}
	prevIdx, prevW := 0, float32(0)
	for i := 1; i <= len(rn.value); i++ {
		if i < len(rn.value) && !utf8.RuneStart(rn.value[i]) {
			continue
		}
		w, _ := text.Measure(rn.value[:i], rn.face, rn.lineH)
		if float32(w) >= rel {
			if float32(w)-rel < rel-prevW {
				return i
			}
			return prevIdx
		}
		prevIdx, prevW = i, float32(w)
	}
	return len(rn.value)
}

func isFocused(rn *renderNode) bool {
	return activeGame != nil && activeGame.focusedFiber != nil && activeGame.focusedFiber.rnode == rn
}

// hitTest 返回命中点最深、且带 onClick 的处理器（兼容旧调用）。
func hitTest(rn *renderNode, px, py float32) func() {
	n := hitNode(rn, px, py)
	for c := n; c != nil; c = c.parent {
		if c.onClick != nil {
			return c.onClick
		}
	}
	return nil
}

// invTransform 把父空间坐标反变换到本节点的未变换（布局）空间。
func (rn *renderNode) invTransform(px, py float32) (float32, float32) {
	if !rn.hasTransform() {
		return px, py
	}
	b := rn.bounds
	cx, cy := b.X+b.W/2, b.Y+b.H/2
	x := px - (cx + rn.effTransX())
	y := py - (cy + rn.effTransY())
	if rn.rotate != 0 {
		rad := -float64(rn.rotate) * math.Pi / 180
		cos, sin := float32(math.Cos(rad)), float32(math.Sin(rad))
		x, y = x*cos-y*sin, x*sin+y*cos
	}
	if rn.scale != 0 {
		x /= rn.scale
		y /= rn.scale
	}
	return x + cx, y + cy
}

// collectFocusables 按树序收集可聚焦节点（输入框与可点击元素）。
func collectFocusables(rn *renderNode, out *[]*renderNode) {
	if rn.focusable {
		*out = append(*out, rn)
	}
	for _, c := range rn.children {
		collectFocusables(c, out)
	}
}

// nextFocus 返回 Tab 顺序中的下一个/上一个可聚焦节点（环形）。
func nextFocus(list []*renderNode, cur *renderNode, forward bool) *renderNode {
	if len(list) == 0 {
		return nil
	}
	idx := -1
	for i, r := range list {
		if r == cur {
			idx = i
			break
		}
	}
	if idx == -1 {
		if forward {
			return list[0]
		}
		return list[len(list)-1]
	}
	if forward {
		return list[(idx+1)%len(list)]
	}
	return list[(idx-1+len(list))%len(list)]
}

// hitNode 返回包含该点的最深 renderNode（沿途反变换查询点，命中测试跟随 transform）。
func hitNode(rn *renderNode, px, py float32) *renderNode {
	lx, ly := rn.invTransform(px, py)
	if !rn.bounds.contains(lx, ly) {
		return nil
	}
	for i := len(rn.children) - 1; i >= 0; i-- {
		if h := hitNode(rn.children[i], lx, ly); h != nil {
			return h
		}
	}
	return rn
}
