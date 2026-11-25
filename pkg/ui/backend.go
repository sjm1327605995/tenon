package ui

import (
	"image"
	"reflect"
	"time"
)

// 渲染后端句柄：这些接口是「UI 引擎」与「具体渲染后端」之间的中立边界。
// 引擎侧（renderNode 字段、Painter 接口、布局/折行/富文本代码）只依赖这些接口，不出现任何
// Gio 的具体类型；后端定义自己的实现类型，并在其 Painter 内部类型断言回去。
//
// 与 Gio 生态的边界（约定，改动前先想清楚）：
//
//   - Tenon 的引擎是 Gio 的**对等层**，不是 Gio widgets 的使用者：我们有自己的
//     fiber/hooks 与 yoga 布局，因此 Gio 的 layout / widget / material 用不上（范式冲突）。
//   - Gio 的**低层能力必须走它的惯例，不要自己造**：op/clip/paint、text.Shaper、
//     io/key 的 IME 协议、io/clipboard、io/pointer 的指针与光标语义。
//     有现成的就照 Gio 自己的实现（如 widget.Editor）当参考。
//   - Gio 的**单位与语义原样透传，禁止自创换算**。反例：曾把 Gio 的像素滚动量除以 24
//     伪造成「档位」再让引擎乘回 24*uiScale，白白多乘一次 uiScale，高 DPI 上滚动速度失真。
//     后端与引擎统一用物理像素。

// FontFace 是后端拥有的字体句柄：布局用它测量文本，Painter 用它绘制字形。
type FontFace interface {
	// Measure 返回字符串 s 在给定行高下的推进宽度（物理像素）。
	Measure(s string, lineH float64) float32
	// Metrics 返回基线以上高度（ascent，物理像素）与是否需要伪粗/伪斜合成。
	Metrics() (ascent float32, fauxBold, fauxItalic bool)
}

// Bitmap 是后端拥有的位图句柄（图片元素与解码后的图像）。
type Bitmap interface {
	// Size 返回像素宽高。
	Size() (w, h int)
}

// VecPath 是后端拥有的矢量路径句柄（图标）。仅作为不透明句柄在引擎与 Painter 间传递。
type VecPath interface{ IsVecPath() }

type fontFace = FontFace
type bitmap = Bitmap
type vecPath = VecPath

// Backend 定义一个完整渲染后端。实现通常在独立模块中提供，并在 init 中调用 SetBackend；
// 不导入任何可选后端时，pkg/ui 默认使用 Gio。所有尺寸均为物理像素。
type Backend interface {
	NewFont(px float32, weight int, italic bool) FontFace
	NewBitmap(img image.Image) Bitmap
	NewVecPath(svgPath string, scale float32) VecPath
	Run(root *Node, cfg WindowConfig)
}

// BackendCapabilities reports complete feature support, not merely the presence
// of a primitive that may degrade or ignore the feature.
type BackendCapabilities struct {
	SVGPaths        bool
	RoundedClipping bool
	LayerTransforms bool
	ExactProjection bool
	Accessibility   bool
}

// BackendCapabilityProvider is an optional backend extension.
type BackendCapabilityProvider interface {
	Capabilities() BackendCapabilities
}

// SetBackend 设置后续 Run 使用的渲染后端。须在启动窗口前调用，不是并发切换 API。
func SetBackend(backend Backend) {
	if backend == nil {
		panic("ui: cannot use a nil render backend")
	}
	activeBackend = backend
	backendNewFont = backend.NewFont
	backendNewBitmap = backend.NewBitmap
	backendNewVecPath = backend.NewVecPath
	backendRun = backend.Run
}

// CurrentBackendCapabilities returns the active backend's declared capabilities.
// Backends without BackendCapabilityProvider conservatively report no advanced features.
func CurrentBackendCapabilities() BackendCapabilities {
	if provider, ok := activeBackend.(BackendCapabilityProvider); ok {
		return provider.Capabilities()
	}
	return BackendCapabilities{}
}

// 后端构造钩子：由当前激活的渲染后端登记。引擎在需要新建句柄/启动窗口时经此调用。
var (
	activeBackend     Backend
	backendNewFont    func(px float32, weight int, italic bool) fontFace
	backendNewBitmap  func(img image.Image) bitmap
	backendNewVecPath func(svgPath string, scale float32) vecPath
	backendRun        func(root *Node, cfg windowConfig)
)

// BackendRuntime 向后端暴露引擎帧驱动能力，而不公开 Fiber/renderNode 的内部结构。
// 后端在自己的窗口循环中依次调用 Resize、Frame 和 Paint，退出时调用 Close。
// Close 是终态：关闭后所有帧驱动方法保持惰性，复用界面必须创建新实例。
// 每个实例保存独立输入源和设备缩放，Frame/Paint 可交错驱动多个实例，但不支持并发调用。
type BackendRuntime struct {
	g               *game
	input           InputSource
	scale           float32
	previousInput   InputSource
	previousScale   float32
	previousRuntime *BackendRuntime
	closed          bool
}

var activeBackendRuntime *BackendRuntime

// NewBackendRuntime 创建一棵由外部后端驱动的 UI 运行时。input 不得为 nil。
func NewBackendRuntime(root *Node, cfg WindowConfig, input InputSource) *BackendRuntime {
	if input == nil {
		panic("ui: backend runtime requires an input source")
	}
	value := reflect.ValueOf(input)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		if value.IsNil() {
			panic("ui: backend runtime requires an input source")
		}
	}
	enableYogaLayoutDebug()
	g := &game{root: root, w: cfg.w, h: cfg.h}
	previousScale := uiScale
	activeGame = g
	uiScale = 1
	previous := setGlobalInput(input)
	runtime := &BackendRuntime{
		g: g, input: input, scale: 1,
		previousInput: previous, previousScale: previousScale,
		previousRuntime: activeBackendRuntime,
	}
	activeBackendRuntime = runtime
	return runtime
}

func (r *BackendRuntime) activate() func() {
	previousGame, previousInput, previousScale := activeGame, input, uiScale
	previousRuntime := activeBackendRuntime
	activeGame, input, uiScale = r.g, r.input, r.scale
	activeBackendRuntime = r
	return func() {
		activeGame, input, uiScale = previousGame, previousInput, previousScale
		activeBackendRuntime = previousRuntime
	}
}

// Resize 更新窗口的物理像素尺寸和设备缩放。scale 小于 1 时按 1 处理。
func (r *BackendRuntime) Resize(w, h int, scale float32) {
	if r == nil || r.closed {
		return
	}
	if scale < 1 {
		scale = 1
	}
	if w > 0 && (w != r.g.w || h != r.g.h || scale != r.scale) {
		r.scale = scale
		if activeGame == r.g {
			uiScale = scale
		}
		r.g.w, r.g.h = w, h
		r.g.needsLayout = true
	}
}

// Frame 执行输入、重渲染、布局和动画，返回 FLIP 布局动画是否仍在移动。
func (r *BackendRuntime) Frame(dtMs float32) bool {
	if r == nil || r.closed {
		return false
	}
	defer r.activate()()
	return gioFrame(r.g, dtMs)
}

// PaintVersion 返回当前视觉版本，包含输入光标闪烁相位。
// 外部后端可据此复用原生绘制节点，仅在版本变化时重建。
func (r *BackendRuntime) PaintVersion(now time.Time) uint64 {
	if r == nil || r.closed {
		return 0
	}
	return r.g.framePaintVersion(now)
}

// Paint 按主树、Portal 浮层的顺序提交绘制指令。
func (r *BackendRuntime) Paint(p Painter) {
	if r == nil || r.closed {
		return
	}
	defer r.activate()()
	if r.g.rootRN != nil {
		paint(p, r.g.rootRN)
	}
	for _, pf := range r.g.portals {
		if pf.overlayRoot != nil {
			paint(p, pf.overlayRoot)
		}
	}
}

// WantsNextFrame 报告引擎是否还有需要立即处理的更新。
func (r *BackendRuntime) WantsNextFrame(layoutMoving bool) bool {
	return r != nil && !r.closed && wantsNextFrame(r.g, layoutMoving)
}

// Size 返回当前窗口的物理像素尺寸。
func (r *BackendRuntime) Size() (w, h int) {
	if r == nil || r.closed {
		return 0, 0
	}
	return r.g.w, r.g.h
}

// PointerInteractive 报告物理像素坐标是否落在交互 UI 上，供游戏后端避免穿透点击。
func (r *BackendRuntime) PointerInteractive(x, y float32) bool {
	if r == nil || r.closed {
		return false
	}
	for n := r.g.hitTop(x, y); n != nil; n = n.parent {
		if n.isDisabled() || n.kind == rnInput || n.kind == rnScroll {
			return true
		}
		if n.onClick != nil || n.onContextMenu != nil || n.onHover != nil || n.onPress != nil || n.onDrag != nil {
			return true
		}
	}
	return false
}

// Close 幂等卸载组件、释放订阅，并恢复最近的有效后端或此前输入源。
func (r *BackendRuntime) Close() {
	if r == nil || r.closed {
		return
	}
	r.closed = true
	wasActive := activeBackendRuntime == r
	wasActiveGame := activeGame == r.g
	r.g.dispose()
	if !wasActive {
		return
	}

	previous := r.previousRuntime
	for previous != nil && previous.closed {
		previous = previous.previousRuntime
	}
	if previous != nil {
		input, uiScale = previous.input, previous.scale
		activeBackendRuntime = previous
		if wasActiveGame {
			activeGame = previous.g
		}
		return
	}

	fallback := r
	for fallback.previousRuntime != nil {
		fallback = fallback.previousRuntime
	}
	input, uiScale = fallback.previousInput, fallback.previousScale
	activeBackendRuntime = nil
	if wasActiveGame {
		activeGame = nil
	}
}
