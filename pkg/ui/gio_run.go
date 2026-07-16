package ui

import (
	"image"
	"time"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/op"
	"gioui.org/op/clip"
	gpaint "gioui.org/op/paint"
	"gioui.org/unit"
)

// ---- gio 后端：窗口与渲染循环 ----
//
// 每帧：排空输入事件 -> 应用输入法编辑 -> 驱动引擎（reconcile/布局）-> 绘制 -> 声明输入区
// -> 同步输入法状态 -> 按需安排下一帧。引擎的 game 编排（reconcile / flushDirty / layout /
// paint）与后端无关，这里只做转接。

// gio 是唯一渲染后端：包加载即登记构造钩子、输入源与运行循环。
func init() {
	backendNewFont = gioNewFont
	backendNewBitmap = func(img image.Image) bitmap { return &gioImage{src: img} }
	backendNewVecPath = func(d string, scale float32) vecPath {
		if d == "" {
			return nil
		}
		return &gioPath{d: d, scale: scale}
	}
	backendRun = gioRun
	input = gioIn
}

func gioRun(root *Node, w, h int, title string, sync bool) {
	g := &game{root: root, w: w, h: h}
	activeGame = g

	win := new(app.Window)
	win.Option(app.Title(title), app.Size(unit.Dp(float32(w)), unit.Dp(float32(h))))
	backendWake = win.Invalidate // 后台任务（Post）用它唤醒睡着的循环；可跨 goroutine 调用
	defer func() { backendWake = nil }()

	var ops op.Ops
	last := time.Time{}
	for {
		switch e := win.Event().(type) {
		case app.DestroyEvent:
			return
		case app.FrameEvent:
			scale := e.Metric.PxPerDp
			if scale < 1 {
				scale = 1
			}
			if e.Size.X > 0 && (e.Size.X != g.w || e.Size.Y != g.h || scale != uiScale) {
				uiScale = scale
				g.w, g.h = e.Size.X, e.Size.Y
				g.needsLayout = true
			}

			// 排空本帧输入事件到 gioIn，再驱动一帧（handleInput 读取 gioIn）。
			gioIn.resetFrame()
			gioIn.process(e.Source)
			// 输入法的编辑按 Range 替换，先落到聚焦输入框上，再让引擎跑这一帧。
			gioIME.applyEdits(g, gioIn)

			// dt 用浮点毫秒：整数 Milliseconds() 会把高刷新率下的帧间隔截断（144Hz 的
			// 6.9ms 变 6），帧间隔不足 1ms 时更会截成 0，而 tickAnims/tickLoops 遇 dt<=0
			// 直接返回 —— 动画会卡住不动。
			var dt float32
			now := e.Now
			if !last.IsZero() {
				dt = float32(now.Sub(last).Seconds() * 1000)
			}
			last = now
			layoutMoving := gioFrame(g, dt)

			ops.Reset()
			gpaint.ColorOp{Color: nrgba(Color{247, 248, 250, 255})}.Add(&ops)
			gpaint.PaintOp{}.Add(&ops)
			p := newGioPainter(&ops, g.w, g.h)
			if g.rootRN != nil {
				paint(p, g.rootRN)
			}
			for _, pf := range g.portals {
				if pf.overlayRoot != nil {
					paint(p, pf.overlayRoot)
				}
			}
			// 声明整窗为输入命中区（引擎自管内部焦点，这里整窗恒接收）。
			area := clip.Rect{Max: e.Size}.Push(&ops)
			event.Op(&ops, gioTag)
			key.InputHintOp{Tag: gioTag, Hint: key.HintAny}.Add(&ops)
			area.Pop()
			// 仅在尚未获得焦点时请求一次；每帧无脑请求会和 gio 的焦点管理打架。
			if !gioIn.focused {
				e.Source.Execute(key.FocusCmd{Tag: gioTag})
			}
			// 把聚焦输入框的光标位置与上下文文本发布给输入法（否则无法组字，只能打英文）。
			if gioIn.snippetReq != nil {
				gioIME.handleSnippetReq(e.Source, g, *gioIn.snippetReq)
			}
			gioIME.sync(e.Source, g)

			e.Frame(&ops)
			scheduleNextFrame(e, g, layoutMoving)
		}
	}
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
func scheduleNextFrame(e app.FrameEvent, g *game, layoutMoving bool) {
	if g.needsLayout || len(g.dirty) > 0 || len(g.anims) > 0 || len(g.loops) > 0 ||
		layoutMoving || g.inputSelecting || g.imeComposing {
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
	if g.rootFiber == nil {
		g.rootFiber = reconcile(nil, nil, g.root)
		g.needsLayout = true
	} else {
		drainPosts()
		g.handleInput() // 读取 gioIn（指针/键盘/滚轮/编辑），分发命中/焦点/拖拽/文本编辑
		g.tickAnims(dt)
		g.tickLoops(dt)
		for guard := 0; len(g.dirty) > 0 && guard < 100; guard++ {
			g.flushDirty()
		}
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
