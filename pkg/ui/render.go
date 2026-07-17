package ui

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

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
	rnIcon
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
	hasGradient bool
	gradFrom    Color
	gradTo      Color
	gradAngle   float32
	radius      float32
	borderW     float32
	borderColor Color
	onClick     func()

	// text / input
	text       string
	runs       []textRun // 富文本：非空时按多段混排绘制/测量
	runsRev    int       // 富文本解析版本（resolveRuns 改动字体时自增），用于排版缓存失效
	wc         wrapCache // 折行缓存（纯文本）
	rc         richCache // 排版缓存（富文本）
	face       fontFace
	color      Color
	lineH      float64
	fauxBold   bool
	fauxItalic bool

	// 文本样式继承
	explicitColor  bool // 本节点显式设置了颜色
	ownColor       Color
	explicitSize   bool // 本节点显式设置了字号
	ownSize        float32
	explicitWeight bool
	ownWeight      int
	explicitItalic bool
	ownItalic      bool
	effSize        float32 // 生效字号（逻辑，含继承）
	effScale       float32 // 上次取字体所用的 uiScale
	effWeight      int
	effItalic      bool
	inhColor       Color // box 向下传递的颜色
	hasInhColor    bool
	inhSize        float32 // box 向下传递的字号
	hasInhSize     bool
	inhWeight      int
	hasInhWeight   bool
	inhItalic      bool
	hasInhItalic   bool

	// input
	value       string
	placeholder string
	multiline   bool
	onChange    func(string)
	onSubmit    func(string)
	caretPos    int
	selAnchor   int // 选区另一端；等于 caretPos 表示无选区

	// 输入法组字区间（字节偏移）。gio 的模型里组字文本已经写进 value，这里只标出范围
	// 用于画预编辑下划线；composeHi<=composeLo 表示当前没在组字。
	composeLo     int
	composeHi     int
	onHover       func(bool)
	onPress       func(bool)
	onDrag        func(dx, dy float32)
	onContextMenu func(x, y float32)
	measure       *measureHook
	scrollRef     *scrollHook // UseScroll：写回滚动状态
	focusable     bool
	navGroup      bool // ArrowNav：本节点是方向键导航组
	navOrient     NavOrient

	// image
	imgSrc    string
	planeImg  image.Image // PlaneImage：非空则作为场景地板精确预变形（见 plane_image.go）
	img       bitmap
	objectFit ObjectFit

	// icon（SVG）/ vector（原始像素路径）
	iconPath    string
	iconSize    float32
	iconStroke  float32 // >0 描边，0 填充
	iconRaw     bool    // Vector：路径为逻辑像素坐标
	iconW       float32
	iconH       float32
	iconCache   vecPath // 已按缩放解析好的路径（局部坐标，绘制时平移到 bounds）
	iconCacheK  string  // 缓存键：path
	iconCacheSz float32 // 缓存键：缩放比例

	// scroll / clip
	clip     bool
	scroll   bool
	scrollY  float32
	contentH float32

	opacity        float32
	scale          float32
	rotate         float32
	transX, transY float32

	// 伪 3D（仅绘制阶段，不入 yoga）
	rotateX, rotateY float32
	transZ           float32
	perspective      float32
	scene3D          bool
	zIndex           int

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

// has3D 表示本节点带伪 3D 变换（需要走投影绘制路径）。
func (rn *renderNode) has3D() bool {
	return rn.rotateX != 0 || rn.rotateY != 0 || rn.transZ != 0
}

func (rn *renderNode) hasTransform() bool {
	return rn.scale != 1 || rn.rotate != 0 || rn.effTransX() != 0 || rn.effTransY() != 0 || rn.has3D()
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

func newTextRenderNode(s string, st StyleProps, runs []textRun) *renderNode {
	rn := &renderNode{yn: yoga.NewNode(), kind: rnText, text: s, runs: runs, opacity: 1, scale: 1}
	rn.applyTextStyle(st)
	rn.yn.SetMeasureFunc(func(_ *yoga.Node, w float32, wm yoga.MeasureMode, _ float32, _ yoga.MeasureMode) yoga.Size {
		avail := float32(0)
		if wm == yoga.MeasureModeExactly || wm == yoga.MeasureModeAtMost {
			avail = w
		}
		if len(rn.runs) > 0 {
			_, mw, h := rn.richLayout(avail)
			return yoga.Size{Width: mw, Height: h}
		}
		if rn.face == nil || rn.text == "" {
			return yoga.Size{}
		}
		lines, mw := rn.wrapped(avail)
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
			wd = measureW(s, rn.face, rn.lineH) + 8
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
		iw, ih := rn.img.Size()
		return yoga.Size{Width: float32(iw), Height: float32(ih)}
	})
	return rn
}

func newIconRenderNode(d string, size, stroke float32, raw bool, w, h float32, st StyleProps) *renderNode {
	rn := &renderNode{yn: yoga.NewNode(), kind: rnIcon, iconPath: d, iconSize: size, iconStroke: stroke,
		iconRaw: raw, iconW: w, iconH: h, opacity: 1, scale: 1}
	rn.applyTextStyle(st)                     // 复用文本颜色继承（currentColor）
	rn.yn.StyleSetAlignSelf(yoga.AlignCenter) // 固定尺寸（不随 align-items:stretch 拉伸），并在交叉轴居中
	rn.yn.SetMeasureFunc(func(_ *yoga.Node, _ float32, _ yoga.MeasureMode, _ float32, _ yoga.MeasureMode) yoga.Size {
		if rn.iconRaw {
			return yoga.Size{Width: rn.iconW * uiScale, Height: rn.iconH * uiScale}
		}
		s := rn.iconSize * uiScale
		return yoga.Size{Width: s, Height: s}
	})
	return rn
}

func (rn *renderNode) setIcon(d string, size, stroke float32, raw bool, w, h float32, st StyleProps) {
	if rn.iconPath != d || rn.iconSize != size || rn.iconStroke != stroke || rn.iconRaw != raw || rn.iconW != w || rn.iconH != h {
		rn.iconPath, rn.iconSize, rn.iconStroke = d, size, stroke
		rn.iconRaw, rn.iconW, rn.iconH = raw, w, h
		rn.iconCache = nil
		rn.yn.MarkDirty()
	}
	rn.applyTextStyle(st)
}

// iconScale 是路径坐标 -> 物理像素的缩放：Vector 原始像素用 uiScale；图标按 size/viewBox。
func (rn *renderNode) iconScale() float32 {
	if rn.iconRaw {
		return uiScale
	}
	return rn.iconSize * uiScale / iconViewBox
}

// scaledIconPath 返回按当前缩放解析好的路径（缓存，局部坐标，绘制时再平移）。
// 无可绘制路径时返回 nil（注意返回真正的 nil 接口，供调用方 != nil 判断）。
func (rn *renderNode) scaledIconPath() vecPath {
	sc := rn.iconScale()
	if rn.iconCache != nil && rn.iconCacheK == rn.iconPath && rn.iconCacheSz == sc {
		return rn.iconCache
	}
	ip := backendNewVecPath(rn.iconPath, sc)
	rn.iconCache, rn.iconCacheK, rn.iconCacheSz = ip, rn.iconPath, sc
	return ip
}

// applyTextStyle 只记录本节点显式设置的文本样式；生效值由 resolveInherited 决定。
func (rn *renderNode) applyTextStyle(st StyleProps) {
	rn.explicitColor = st.hasColor
	rn.ownColor = st.color
	rn.explicitSize = st.hasFontSize
	rn.ownSize = st.fontSize
	rn.explicitWeight = st.hasWeight
	rn.ownWeight = st.weight
	rn.explicitItalic = st.hasItalic
	rn.ownItalic = st.italic
	if rn.effSize == 0 { // 初始回退，保证在 resolve 前也有可用字体
		rn.setEffectiveText(Black, 16, 400, false)
	}
}

// setEffectiveText 应用最终生效的颜色/字号/字重/斜体（在物理像素下取字体，保证高分屏清晰）。
func (rn *renderNode) setEffectiveText(c Color, size float32, weight int, italic bool) {
	if size <= 0 {
		size = 16
	}
	if weight <= 0 {
		weight = 400
	}
	rn.color = c
	if rn.effSize != size || rn.effScale != uiScale || rn.effWeight != weight || rn.effItalic != italic || rn.face == nil {
		rn.effSize, rn.effScale, rn.effWeight, rn.effItalic = size, uiScale, weight, italic
		px := size * uiScale
		rn.lineH = float64(px) * 1.3
		if f := backendNewFont(px, weight, italic); f != nil {
			rn.face = f
			_, rn.fauxBold, rn.fauxItalic = f.Metrics()
		}
		rn.yn.MarkDirty()
	}
}

// inhText 是文本继承上下文（自顶向下传递）。
type inhText struct {
	color              Color
	size               float32
	weight             int
	hasColor, hasSize  bool
	hasWeight, hasItal bool
	italic             bool
}

// resolveInherited 自顶向下解析文本继承（颜色/字号/字重/斜体）；文本/输入节点未显式设置时采用继承值。
// 须在测量（CalculateLayout）之前调用。
func resolveInherited(rn *renderNode, ctx inhText) {
	switch rn.kind {
	case rnIcon:
		c := Black
		if rn.explicitColor {
			c = rn.ownColor
		} else if ctx.hasColor {
			c = ctx.color
		}
		rn.color = c
	case rnText, rnInput:
		if rn.kind == rnText && len(rn.runs) > 0 {
			rn.resolveRuns(ctx)
			return
		}
		c := Black
		if rn.explicitColor {
			c = rn.ownColor
		} else if ctx.hasColor {
			c = ctx.color
		}
		s := float32(16)
		if rn.explicitSize {
			s = rn.ownSize
		} else if ctx.hasSize {
			s = ctx.size
		}
		w := 400
		if rn.explicitWeight {
			w = rn.ownWeight
		} else if ctx.hasWeight {
			w = ctx.weight
		}
		it := false
		if rn.explicitItalic {
			it = rn.ownItalic
		} else if ctx.hasItal {
			it = ctx.italic
		}
		rn.setEffectiveText(c, s, w, it)
	default:
		if rn.hasInhColor {
			ctx.color, ctx.hasColor = rn.inhColor, true
		}
		if rn.hasInhSize {
			ctx.size, ctx.hasSize = rn.inhSize, true
		}
		if rn.hasInhWeight {
			ctx.weight, ctx.hasWeight = rn.inhWeight, true
		}
		if rn.hasInhItalic {
			ctx.italic, ctx.hasItal = rn.inhItalic, true
		}
		for _, ch := range rn.children {
			resolveInherited(ch, ctx)
		}
	}
}

func (rn *renderNode) setText(s string, st StyleProps, runs []textRun) {
	changed := rn.text != s || !runsEqual(rn.runs, runs)
	rn.text = s
	if changed {
		rn.runs = runs // 未变化时保留旧 runs（含已解析的字体缓存）
	}
	rn.applyTextStyle(st)
	if changed {
		rn.yn.MarkDirty()
	}
}

// 图片缓存见 imgcache.go（按字节预算的 LRU）。
var (
	imgLoading = map[string]bool{}          // 正在后台加载的 src
	imgWaiters = map[string][]*renderNode{} // 等待同一 src 完成的节点
)

// imgHTTPClient 带超时，避免后台加载 goroutine 因网络挂起而永久泄漏。
var imgHTTPClient = &http.Client{Timeout: 15 * time.Second}

// setImage 安装一张已解码的图片（SrcImage）。没有 IO 也没有解码，故不像 loadImage 那样
// 绕后台 goroutine 与 Post，直接同步建位图 —— 当帧即可见，不会先闪一帧空白。
// 仍进同一份 LRU：同一 key 的多个节点共用一张位图。
func (rn *renderNode) setImage(img image.Image) {
	if b, ok := lookupImage(rn.imgSrc); ok {
		rn.img = b
		return
	}
	b := backendNewBitmap(img)
	storeImage(rn.imgSrc, b)
	rn.img = b
}

// loadImage 异步加载图片：命中缓存立即返回；否则在后台 goroutine 解码，完成后经
// Post 回到渲染线程安装缓存并触发重绘。imgCache/imgLoading/imgWaiters 仅在渲染线程
// （reconcile 与 drainPosts）访问，故无需加锁；后台 goroutine 只做 IO/解码。
func (rn *renderNode) loadImage() {
	if img, ok := lookupImage(rn.imgSrc); ok {
		rn.img = img
		return
	}
	rn.img = nil // 加载完成前不显示旧图
	src := rn.imgSrc
	imgWaiters[src] = append(imgWaiters[src], rn)
	if imgLoading[src] {
		return // 已在加载，完成时统一通知全部等待者
	}
	imgLoading[src] = true
	go func() {
		img, err := decodeImageSource(src)
		Post(func() {
			delete(imgLoading, src)
			waiters := imgWaiters[src]
			delete(imgWaiters, src)
			if err != nil || img == nil {
				return // 失败：保持未加载（不缓存失败），src 再次设置时会重试
			}
			ei := backendNewBitmap(img)
			storeImage(src, ei)
			for _, w := range waiters {
				if w.imgSrc == src { // 期间 src 可能已改，仅回填仍需要它的节点
					w.img = ei
					w.yn.MarkDirty()
				}
			}
			if activeGame != nil {
				activeGame.needsLayout = true
			}
		})
	}()
}

// decodeImageSource 从本地路径或 http(s) URL 读取并解码一张图片（在后台 goroutine 调用）。
func decodeImageSource(src string) (image.Image, error) {
	if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
		resp, err := imgHTTPClient.Get(src)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("image: GET %s -> %s", src, resp.Status)
		}
		img, _, err := image.Decode(resp.Body)
		return img, err
	}
	f, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}

// applyHostProps 把某个 host 元素的属性写入其 renderNode。
func applyHostProps(rn *renderNode, hp hostProps) {
	syncYoga(rn, hp.style)
	rn.onClick = hp.onClick
	rn.onHover = hp.onHover
	rn.onPress = hp.onPress
	rn.onDrag = hp.onDrag
	rn.onContextMenu = hp.onContextMenu
	rn.measure = hp.measure
	rn.scrollRef = hp.scrollRef
	rn.navGroup = hp.navGroup
	rn.navOrient = hp.navOrient
	rn.focusable = rn.kind == rnInput || hp.onClick != nil
	switch rn.kind {
	case rnInput:
		rn.applyTextStyle(hp.style)
		rn.value = hp.value
		rn.placeholder = hp.placeholder
		rn.multiline = hp.multiline
		rn.onChange = hp.onChange
		rn.onSubmit = hp.onSubmit
		if rn.caretPos > len(rn.value) {
			rn.caretPos = len(rn.value)
		}
		if rn.selAnchor > len(rn.value) {
			rn.selAnchor = len(rn.value)
		}
		rn.yn.MarkDirty()
	case rnImage:
		rn.objectFit = hp.objectFit
		rn.planeImg = hp.planeImg
		if hp.src != "" && rn.imgSrc != hp.src {
			rn.imgSrc = hp.src
			if hp.imgData != nil {
				rn.setImage(hp.imgData)
			} else {
				rn.loadImage()
			}
			rn.yn.MarkDirty()
		}
	}
}

// fitRect 按 object-fit 计算图片在框 b 内的绘制矩形；bool 表示是否需要裁剪到 b（cover）。
func fitRect(iw, ih float32, b Rect, fit ObjectFit) (Rect, bool) {
	if iw <= 0 || ih <= 0 || fit == FitFill {
		return b, false
	}
	scale := b.W / iw
	switch fit {
	case FitContain:
		if s := b.H / ih; s < scale {
			scale = s
		}
	case FitCover:
		if s := b.H / ih; s > scale {
			scale = s
		}
	}
	dw, dh := iw*scale, ih*scale
	return Rect{X: b.X + (b.W-dw)/2, Y: b.Y + (b.H-dh)/2, W: dw, H: dh}, fit == FitCover
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
	if s.hasWrap {
		yn.StyleSetFlexWrap(s.wrap)
	}
	if s.hasCont {
		yn.StyleSetAlignContent(s.content)
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
	rn.hasGradient = s.hasGradient
	rn.gradFrom, rn.gradTo, rn.gradAngle = s.gradFrom, s.gradTo, s.gradAngle
	rn.radius = s.radius * k
	rn.borderW = s.borderW * k
	rn.borderColor = s.borderColor
	rn.opacity = s.opacity
	rn.scale = s.scale
	rn.rotate = s.rotate
	rn.transX, rn.transY = s.transX*k, s.transY*k
	rn.rotateX, rn.rotateY = s.rotateX, s.rotateY
	rn.transZ = s.transZ * k           // Z 位移随 uiScale 换算到物理像素
	rn.perspective = s.perspective * k // 透视距离同为物理像素，与投影坐标同一量纲
	rn.scene3D = s.scene3D
	rn.zIndex = s.zIndex

	rn.hasShadow = s.hasShadow
	rn.shadowColor = s.shadowColor
	rn.shadowX, rn.shadowY = s.shadowX*k, s.shadowY*k
	rn.shadowBlur, rn.shadowSpread = s.shadowBlur*k, s.shadowSpread*k
	if s.clip {
		rn.clip = true
	}

	// 容器可通过 TextColor/FontSize/FontWeight/Italic 为后代文本设定继承值
	rn.hasInhColor = s.hasColor
	rn.inhColor = s.color
	rn.hasInhSize = s.hasFontSize
	rn.inhSize = s.fontSize
	rn.hasInhWeight = s.hasWeight
	rn.inhWeight = s.weight
	rn.hasInhItalic = s.hasItalic
	rn.inhItalic = s.italic

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
// camera3D 是一台共享相机：由 Scene3D 容器确立，为其直接子元素提供同一个投影原点与
// 视角，使它们灭点一致。见 Scene3D 的文档。
type camera3D struct {
	ox, oy           float32 // 投影原点 = 场景中心
	rotateX, rotateY float32 // 场景平面的倾角
	perspective      float32
}

// cameraOf 返回某个 Scene3D 场景根确立的相机。
func cameraOf(rn *renderNode) *camera3D {
	b := rn.bounds
	return &camera3D{
		ox: b.X + b.W/2, oy: b.Y + b.H/2,
		rotateX: rn.rotateX, rotateY: rn.rotateY, perspective: rn.perspective,
	}
}

// layerOf 构造某节点合成时施加的变换。绘制与命中测试必须共用它 —— 两边各算一套正是
// 「画得到却点不到」的根源（本仓已因此修过一次：e19e310）。
func layerOf(rn *renderNode, cam *camera3D) layerTransform {
	b := rn.bounds
	t := layerTransform{
		cx: b.X + b.W/2, cy: b.Y + b.H/2,
		w: b.W, h: b.H,
		scale: rn.scale, rotate: rn.rotate,
		tx: rn.effTransX(), ty: rn.effTransY(), opacity: rn.opacity,
		rotateX: rn.rotateX, rotateY: rn.rotateY,
		transZ: rn.transZ, perspective: rn.perspective,
	}
	if cam != nil {
		t.camX, t.camY, t.hasCam = cam.ox, cam.oy, true
		t.perspective = cam.perspective
		t.rotateX += cam.rotateX // 欧拉角相加，非严格矩阵复合（见 Scene3D 文档）
		t.rotateY += cam.rotateY
	}
	return t
}

// paintOrder 返回子节点的绘制顺序：按 zIndex 升序（大的后画=在上面），同值保持兄弟顺序。
//
// 绘制正序遍历它、命中测试逆序遍历它 —— 必须共用这一个函数，否则「画在上面的」和
// 「点得到的」会是两个元素（本仓已两次栽在绘制与命中各算一套上：e19e310、27f5ce5）。
func paintOrder(rn *renderNode) []*renderNode {
	sorted := false
	for _, c := range rn.children {
		if c.zIndex != 0 {
			sorted = true
			break
		}
	}
	if !sorted { // 绝大多数容器没人设 zIndex：直接用原切片，不排序不分配
		return rn.children
	}
	out := make([]*renderNode, len(rn.children))
	copy(out, rn.children)
	// 必须是稳定排序：同 zIndex 时保持兄弟顺序，与不设 zIndex 的行为一致
	sort.SliceStable(out, func(i, j int) bool { return out[i].zIndex < out[j].zIndex })
	return out
}

func paint(p painter, rn *renderNode) { paintIn(p, rn, nil) }

// paintIn 绘制 rn；cam 非 nil 表示 rn 是某个 Scene3D 的直接子元素，应透过该相机投影。
func paintIn(p painter, rn *renderNode, cam *camera3D) {
	if rn.scene3D {
		paintScene(p, rn, cam)
		return
	}
	// 地板：必须绕开图层+仿射那条路 —— 精确投影正是靠「不走仿射」得来的。
	if rn.planeImg != nil && cam != nil && paintPlaneImage(p, rn, cam) {
		return
	}
	if rn.needsLayer() || cam != nil {
		paintLayer(p, rn, cam)
		return
	}
	paintNode(p, rn, nil)
}

// paintScene 画一个共享相机场景：自身（背景/边框）按相机投影，然后每个直接子元素
// 各自成层、共用这台相机 —— 关键在于子元素不能被卷进场景自己的图层，
// 否则就退化成「把整块桌子当一个元素倾斜」，尺寸一大投影就失真。
func paintScene(p painter, rn *renderNode, outer *camera3D) {
	b := rn.bounds
	cam := cameraOf(rn)
	t := layerOf(rn, nil) // 场景自身不透过自己的相机：原点本就是自己的中心
	// 场景自身（桌面）：原点就是自己的中心，故与无相机时等价 —— 也就是说它照样受
	// 「大元素失真」的约束：640x420 的桌面在 50° 下内容斜切达 211px（自身高度的 50%），
	// 画出来是被斜切的平行四边形而非梯形。卡牌没这个问题（78x104 只斜 6px）。
	// 要让桌面精确，得按投影四边形直接填充而不是走仿射 —— 那需要一个新的绘制原语，
	// 尚未做。眼下：陡角度下别指望场景自身的背景当桌面，或把倾角放缓。
	o := rn.opacity
	rn.opacity = 1
	p.BeginLayer()
	paintSelf(p, rn)
	rn.opacity = o
	p.EndLayer(t)
	// 场景的 Clip 在 3D 下不生效：裁剪矩形是未投影的，会把卡切错（见 Scene3D 文档）。
	for _, c := range paintOrder(rn) {
		paintIn(p, c, cam)
	}
	if rn.scroll && rn.contentH > b.H {
		drawScrollbar(p, rn)
	}
	_ = outer // 场景不嵌套：内层场景自成相机，不继承外层
}

// paintLayer 把子树合成到独立图层，再围绕中心做 scale/rotate/translate 及整组透明度合回。
// cam 非 nil 时改走共享相机投影：角度叠加到相机上，投影原点用场景中心。
func paintLayer(p painter, rn *renderNode, cam *camera3D) {
	t := layerOf(rn, cam) // 先取，下面会临时改 rn.opacity
	o := rn.opacity
	rn.opacity = 1 // 组透明度在合成时统一施加，内部按不透明绘制
	p.BeginLayer()
	paintNode(p, rn, nil) // 子树在本图层内扁平化
	rn.opacity = o
	p.EndLayer(t)
}

func paintNode(p painter, rn *renderNode, cam *camera3D) {
	paintSelf(p, rn)
	b := rn.bounds

	// 裁剪：把子节点画进自身矩形（越界部分被裁掉；有圆角则裁到圆角）。
	if rn.clip {
		p.PushClip(b, rn.radius)
		for _, c := range paintOrder(rn) {
			paintIn(p, c, cam)
		}
		p.PopClip()
	} else {
		for _, c := range paintOrder(rn) {
			paintIn(p, c, cam)
		}
	}

	if rn.scroll && rn.contentH > b.H {
		drawScrollbar(p, rn)
	}

	// 键盘焦点环
	if rn.focusable && isFocused(rn) {
		p.StrokeRect(b.X-2, b.Y-2, b.W+4, b.H+4, rn.radius+2, 2, Hex("#60a5fa"))
	}
}

// paintSelf 只画元素自身（阴影/背景/边框/文字/图片/图标），不含子节点。
func paintSelf(p painter, rn *renderNode) {
	b := rn.bounds
	o := rn.opacity
	if rn.hasShadow {
		p.Shadow(b.X, b.Y, b.W, b.H, rn.radius,
			rn.shadowX, rn.shadowY, rn.shadowBlur, rn.shadowSpread, rn.shadowColor.Alpha(o))
	}
	switch rn.kind {
	case rnBox, rnScroll:
		if rn.hasGradient {
			p.FillGradient(b.X, b.Y, b.W, b.H, rn.radius, rn.gradFrom.Alpha(o), rn.gradTo.Alpha(o), rn.gradAngle)
		} else {
			p.FillRect(b.X, b.Y, b.W, b.H, rn.radius, rn.bg.Alpha(o))
		}
		if rn.borderW > 0 {
			p.StrokeRect(b.X, b.Y, b.W, b.H, rn.radius, rn.borderW, rn.borderColor.Alpha(o))
		}
	case rnInput:
		paintInput(p, rn)
	case rnText:
		if len(rn.runs) > 0 {
			paintRichText(p, rn, o)
			break
		}
		lines, _ := rn.wrapped(b.W)
		for i, ln := range lines {
			p.DrawText(ln, rn.face, rn.color.Alpha(o), b.X, b.Y+float32(i)*float32(rn.lineH), rn.fauxBold, rn.fauxItalic)
		}
	case rnImage:
		if rn.img != nil {
			iw, ih := rn.img.Size()
			dr, needClip := fitRect(float32(iw), float32(ih), b, rn.objectFit)
			if needClip { // cover：裁剪到框（遵循圆角）
				p.PushClip(b, rn.radius)
				p.DrawImage(rn.img, dr, o)
				p.PopClip()
			} else {
				p.DrawImage(rn.img, dr, o)
			}
		}
	case rnIcon:
		if path := rn.scaledIconPath(); path != nil {
			if rn.iconStroke > 0 {
				sw := rn.iconStroke * rn.iconScale() // 描边宽随缩放
				p.StrokePath(path, b.X, b.Y, sw, rn.color.Alpha(o))
			} else {
				p.FillPath(path, b.X, b.Y, rn.color.Alpha(o))
			}
		}
	}

}

func drawScrollbar(p painter, rn *renderNode) {
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
	p.FillRect(b.X+b.W-6, ty, 4, thumb, 2, Color{0, 0, 0, 90})
}

func paintInput(p painter, rn *renderNode) {
	b := rn.bounds
	p.FillRect(b.X, b.Y, b.W, b.H, rn.radius, rn.bg)
	if rn.borderW > 0 {
		p.StrokeRect(b.X, b.Y, b.W, b.H, rn.radius, rn.borderW, rn.borderColor)
	}
	padL := rn.yn.LayoutPadding(yoga.EdgeLeft)
	padT := rn.yn.LayoutPadding(yoga.EdgeTop)
	tx := b.X + padL
	lineH := float32(rn.lineH)

	// 组字（IME 预编辑）：文本本身已在 value 内，这里只取出区间画下划线。
	val := rn.value
	caret := clampi(rn.caretPos, 0, len(val))
	preLo, preHi := -1, -1
	if isFocused(rn) && rn.composeHi > rn.composeLo {
		preLo = clampi(rn.composeLo, 0, len(val))
		preHi = clampi(rn.composeHi, 0, len(val))
	}
	usePlaceholder := val == "" && rn.placeholder != ""
	selLo := clampi(min(rn.selAnchor, rn.caretPos), 0, len(rn.value))
	selHi := clampi(max(rn.selAnchor, rn.caretPos), 0, len(rn.value))
	hasSel := preLo < 0 && selLo != selHi

	if rn.multiline {
		ty := b.Y + padT
		avail := b.W - padL*2
		if usePlaceholder {
			for i, sp := range wrapSpans(rn.placeholder, rn.face, rn.lineH, avail) {
				p.DrawText(sp.text, rn.face, Gray, tx, ty+float32(i)*lineH, rn.fauxBold, rn.fauxItalic)
			}
			return
		}
		spans := wrapSpans(val, rn.face, rn.lineH, avail)
		if hasSel {
			paintSpanRange(p, spans, selLo, selHi, rn.face, rn.lineH, tx, ty, lineH, false)
		}
		if preLo >= 0 { // 预编辑下划线
			paintSpanRange(p, spans, preLo, preHi, rn.face, rn.lineH, tx, ty, lineH, true)
		}
		for i, sp := range spans {
			p.DrawText(sp.text, rn.face, rn.color, tx, ty+float32(i)*lineH, rn.fauxBold, rn.fauxItalic)
		}
		if isFocused(rn) && caretVisible() {
			for i, sp := range spans {
				if caret <= sp.end {
					cx := sp.xInSpan(caret, rn.face, rn.lineH, tx)
					cy := ty + float32(i)*lineH
					p.Line(cx, cy, cx, cy+lineH, rn.color)
					break
				}
			}
		}
		return
	}

	ty := b.Y + (b.H-lineH)/2

	if hasSel && rn.face != nil {
		wa := measureW(val[:selLo], rn.face, rn.lineH)
		wc := measureW(val[:selHi], rn.face, rn.lineH)
		p.FillRect(tx+wa, ty, wc-wa, lineH, 0, Color{R: 59, G: 130, B: 246, A: 90})
	}
	if preLo >= 0 && rn.face != nil { // 预编辑下划线
		wa := measureW(val[:preLo], rn.face, rn.lineH)
		wc := measureW(val[:preHi], rn.face, rn.lineH)
		p.Line(tx+wa, ty+lineH-1, tx+wc, ty+lineH-1, rn.color)
	}

	if usePlaceholder {
		p.DrawText(rn.placeholder, rn.face, Gray, tx, ty, rn.fauxBold, rn.fauxItalic)
	} else {
		p.DrawText(val, rn.face, rn.color, tx, ty, rn.fauxBold, rn.fauxItalic)
	}
	if isFocused(rn) && caretVisible() {
		cx := tx
		if rn.face != nil && caret > 0 {
			cx += measureW(val[:caret], rn.face, rn.lineH)
		}
		p.Line(cx, ty, cx, ty+lineH, rn.color)
	}
}

// paintSpanRange 逐行为字节区间 [lo,hi) 绘制高亮块或底部下划线（多行输入用）。
func paintSpanRange(p painter, spans []wrapSpan, lo, hi int, face fontFace, lineH float64, tx, ty, lh float32, underline bool) {
	for i, sp := range spans {
		a, c := max(lo, sp.start), min(hi, sp.end)
		if a >= c {
			continue
		}
		x0 := sp.xInSpan(a, face, lineH, tx)
		x1 := sp.xInSpan(c, face, lineH, tx)
		y := ty + float32(i)*lh
		if underline {
			p.Line(x0, y+lh-1, x1, y+lh-1, Color{R: 30, G: 30, B: 30, A: 255})
		} else {
			p.FillRect(x0, y, x1-x0, lh, 0, Color{R: 59, G: 130, B: 246, A: 90})
		}
	}
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
		w := measureW(rn.value[:i], rn.face, rn.lineH)
		if w >= rel {
			if w-rel < rel-prevW {
				return i
			}
			return prevIdx
		}
		prevIdx, prevW = i, w
	}
	return len(rn.value)
}

// caretFromPoint 把绝对屏幕点映射到多行输入内的字节偏移。
func (rn *renderNode) caretFromPoint(px, py float32) int {
	if rn.face == nil || rn.value == "" {
		return 0
	}
	padL := rn.yn.LayoutPadding(yoga.EdgeLeft)
	padT := rn.yn.LayoutPadding(yoga.EdgeTop)
	spans := wrapSpans(rn.value, rn.face, rn.lineH, rn.bounds.W-padL*2)
	row := int((py - (rn.bounds.Y + padT)) / float32(rn.lineH))
	if row < 0 {
		row = 0
	}
	if row >= len(spans) {
		row = len(spans) - 1
	}
	return spans[row].offsetInSpan(px-(rn.bounds.X+padL), rn.face, rn.lineH)
}

// caretRect 返回 caret 处的屏幕矩形（物理像素），用于放置 IME 候选窗。
func (rn *renderNode) caretRect(caret int) image.Rectangle {
	b := rn.bounds
	padL := rn.yn.LayoutPadding(yoga.EdgeLeft)
	lineH := float32(rn.lineH)
	caret = clampi(caret, 0, len(rn.value))
	var cx, cy float32
	if rn.multiline {
		padT := rn.yn.LayoutPadding(yoga.EdgeTop)
		tx := b.X + padL
		cx, cy = tx, b.Y+padT
		for i, sp := range wrapSpans(rn.value, rn.face, rn.lineH, b.W-padL*2) {
			if caret <= sp.end {
				cx = sp.xInSpan(caret, rn.face, rn.lineH, tx)
				cy = b.Y + padT + float32(i)*lineH
				break
			}
		}
	} else {
		cx = b.X + padL
		if rn.face != nil && caret > 0 {
			cx += measureW(rn.value[:caret], rn.face, rn.lineH)
		}
		cy = b.Y + (b.H-lineH)/2
	}
	return image.Rect(int(cx), int(cy), int(cx)+1, int(cy+lineH))
}

// caretAt 根据单行/多行选择合适的光标定位方式。
func (rn *renderNode) caretAt(px, py float32) int {
	if rn.multiline {
		return rn.caretFromPoint(px, py)
	}
	return rn.caretFromX(px)
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
func (rn *renderNode) invTransform(px, py float32, cam *camera3D) (float32, float32) {
	if !rn.hasTransform() && cam == nil {
		return px, py
	}
	if t := layerOf(rn, cam); t.is3D() {
		// 3D：绘制施加的就是 contentAffine，这里取它的逆 —— 同一个矩阵，画哪点哪。
		q := contentAffine(t).invert().transform(pt{px, py})
		return q.X, q.Y
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
// hitNode 自顶向下找命中的最深节点（子节点逆序遍历 = 上层优先）。
//
// 命中必须与绘制用同一套裁剪语义：绘制只在设了 clip 时裁掉后代，所以这里也只有 clip
// 容器才能拦下后代。曾经无论是否 clip 都用自身 bounds 提前 return，导致溢出到未裁剪父
// 容器之外的元素「画得出来却点不到」。
//
// 代价是非裁剪容器无法靠 bounds 剪枝，命中要走完整棵树；对典型规模的树只是若干次矩形
// 比较，且悬停只在光标移动时才重算，可以接受。
func hitNode(rn *renderNode, px, py float32) *renderNode { return hitNodeIn(rn, px, py, nil) }

// hitNodeIn 的相机走向必须与 paintIn 逐字对应，否则又会「画在一处、点在另一处」：
//
//   - 场景根：子元素各自透过相机独立投影、不跟随场景自身的变换（对应 paintScene），
//     所以子节点拿到的是原始屏幕点 + 相机，而不是场景反变换后的点。
//   - 其他节点：子树随本节点一起被变换、在层内是扁平的（对应 paintLayer 传 nil），
//     所以子节点拿反变换后的点、且不再带相机。
func hitNodeIn(rn *renderNode, px, py float32, cam *camera3D) *renderNode {
	lx, ly := rn.invTransform(px, py, cam)
	inside := rn.bounds.contains(lx, ly)
	if rn.clip && !inside {
		return nil // 裁剪容器：框外的后代不可见，也就不可命中（scroll 容器同样置了 clip）
	}
	cx, cy, childCam := lx, ly, (*camera3D)(nil)
	if rn.scene3D {
		cx, cy, childCam = px, py, cameraOf(rn)
	}
	// 逆着绘制顺序找：上层优先。与 paint 共用 paintOrder，顺序不可能分岔。
	order := paintOrder(rn)
	for i := len(order) - 1; i >= 0; i-- {
		if h := hitNodeIn(order[i], cx, cy, childCam); h != nil {
			return h
		}
	}
	if inside {
		return rn
	}
	return nil
}
