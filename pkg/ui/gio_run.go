package ui

import (
	"image"
	"time"

	"gioui.org/app"
	"gioui.org/op"
	gpaint "gioui.org/op/paint"
	"gioui.org/unit"
)

// ---- gio 后端：窗口与渲染循环 ----
//
// 每帧：排空输入事件 -> 应用输入法编辑 -> 驱动引擎（reconcile/布局）-> 绘制 -> 声明输入区
// -> 同步输入法状态 -> 按需安排下一帧。引擎的 game 编排（reconcile / flushDirty / layout /
// paint）与后端无关，这里只做转接。

// Gio 是默认渲染后端；可选后端模块初始化时可通过 SetBackend 替换。
type gioBackend struct{}

func (gioBackend) NewFont(px float32, weight int, italic bool) FontFace {
	return gioNewFont(px, weight, italic)
}

func (gioBackend) NewBitmap(img image.Image) Bitmap {
	return &gioImage{src: img, op: gpaint.NewImageOp(img)}
}

func (gioBackend) NewVecPath(d string, scale float32) VecPath {
	if d == "" {
		return nil
	}
	return &gioPath{d: d, scale: scale}
}

func (gioBackend) Capabilities() BackendCapabilities {
	return BackendCapabilities{
		SVGPaths:        true,
		RoundedClipping: true,
		LayerTransforms: true,
		ExactProjection: true,
		Accessibility:   true,
	}
}

func (gioBackend) Run(root *Node, cfg WindowConfig) { gioRun(root, cfg) }

func init() {
	SetBackend(gioBackend{})
	input = gioIn
}

// gioWindowOptions 把中立的 windowConfig 映射成 gio 的窗口选项。
// 注意 app.Size 会把模式重置为 Windowed，所以全屏必须放在它之后。
func gioWindowOptions(cfg windowConfig) []app.Option {
	dp := func(v int) unit.Dp { return unit.Dp(float32(v)) }
	opts := []app.Option{app.Title(cfg.title), app.Size(dp(cfg.w), dp(cfg.h))}
	if cfg.minW > 0 || cfg.minH > 0 {
		opts = append(opts, app.MinSize(dp(cfg.minW), dp(cfg.minH)))
	}
	if cfg.maxW > 0 || cfg.maxH > 0 {
		opts = append(opts, app.MaxSize(dp(cfg.maxW), dp(cfg.maxH)))
	}
	if cfg.fullscreen {
		opts = append(opts, app.Fullscreen.Option())
	}
	return opts
}

func gioRun(root *Node, cfg windowConfig) {
	newGioWindowRuntime(root, cfg).run()
}

// caretBlinkPeriod 是光标闪烁的半周期，须与 caretVisible() 的判据一致。
const caretBlinkPeriod = 500 * time.Millisecond

// scheduleNextFrame 决定还要不要下一帧 —— 界面静止时就不出帧，让 gio 睡到下次输入。
//
// 迁移期这里曾无条件 Invalidate，等于恒定 60~144fps 空转；而「按需重绘、静止时不烧电」
// 正是选 gio 而非 ebiten 的主要理由之一，所以要按引擎的实际状态来判断：
//   - 有动画/待处理重渲染/布局变更 -> 立刻要下一帧
//   - 只是光标在闪 -> 用 InvalidateCmd 定时唤醒到下个闪烁边界，而不是每帧都画
//   - 其它情况 -> 什么都不做；用户输入或 Post（backendWake）会把循环叫醒
//
// 注意不能用 g.hasLayoutAnim 判断：它只表示「树里存在 Animated 节点」且置位后从不复位，
// 拿它当条件会让循环永远醒着。要用 tickLayoutAnim 的返回值（是否真的还在移动）。
// wantsNextFrame 表示引擎还有事情要做、需要立刻再出一帧。
func hasRunningLoops(g *game) bool {
	for _, loop := range g.loops {
		if loop.running && !loop.fiber.unmounted {
			return true
		}
	}
	return false
}

func wantsNextFrame(g *game, layoutMoving bool) bool {
	return g.needsLayout || len(g.dirty) > 0 || len(g.anims) > 0 || hasRunningLoops(g) ||
		layoutMoving || g.inputSelecting || g.imeComposing
}

func scheduleNextFrame(e app.FrameEvent, g *game, layoutMoving bool) {
	if pendingCapture != nil || wantsNextFrame(g, layoutMoving) {
		e.Source.Execute(op.InvalidateCmd{})
		return
	}
	// 聚焦输入框时光标要闪：只在下一个闪烁边界唤醒。
	if f := g.focusedFiber; f != nil && f.rnode != nil && f.rnode.kind == rnInput {
		now := e.Now
		next := now.Truncate(caretBlinkPeriod).Add(caretBlinkPeriod)
		e.Source.Execute(op.InvalidateCmd{At: next})
	}
}

// gioFrame 驱动一帧的 reconcile + 布局（对应 game.Update 中与输入无关的部分）。
// 返回布局动画（FLIP）是否仍在移动，供 scheduleNextFrame 判断要不要继续出帧。
func gioFrame(g *game, dt float32) (layoutMoving bool) {
	drainPosts()
	if g.rootFiber == nil {
		g.rootFiber = reconcile(nil, nil, g.root)
		g.needsLayout = true
	} else {
		g.handleInput() // 读取 gioIn（指针/键盘/滚轮/编辑），分发命中/焦点/拖拽/文本编辑
		g.tickAnims(dt)
		g.tickLoops(dt)
	}
	for guard := 0; len(g.dirty) > 0 && guard < 100; guard++ {
		g.flushDirty()
	}
	if g.needsLayout {
		g.rootRN = rootRenderNode(g.rootFiber)
		g.layout()
		g.needsLayout = false
	}
	layoutMoving = g.tickLayoutAnim(dt)
	g.flushEffects()
	return layoutMoving
}
