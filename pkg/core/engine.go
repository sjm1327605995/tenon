package core

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sjm1327605995/tenon/yoga"
)

// Engine 是 Tenon 核心，负责 Widget → Element 的构建、布局计算、绘制和事件。
type Engine struct {
	rootWidget   Widget
	rootElement  Element

	// 帧调度
	dirtyElements []Element
	buildQueue    []Widget

	// 窗口
	screenWidth  int
	screenHeight int

	// 输入状态
	lastMouseX   float32
	lastMouseY   float32
	mouseDownTarget Element

	// 时间
	lastFrameTime time.Time
}

// NewEngine 创建引擎。
func NewEngine(rootWidget Widget, width, height int) *Engine {
	return &Engine{
		rootWidget:    rootWidget,
		screenWidth:   width,
		screenHeight:  height,
		lastFrameTime: time.Now(),
	}
}

// Mount 执行首次挂载，调用 Widget.Build() 并构建 Element 树。
func (e *Engine) Mount() {
	if e.rootWidget == nil {
		return
	}
	e.rootWidget.OnMount(e)
	e.rootElement = e.rootWidget.Build()
	if e.rootElement != nil {
		e.onElementMounted(e.rootElement)
		e.calculateLayout()
	}
}

// onElementMounted 递归给新挂载的 Element 注入 Engine。
func (e *Engine) onElementMounted(el Element) {
	if el == nil {
		return
	}
	el.SetEngine(e)
	el.OnMount(e)
	for _, child := range el.GetChildren() {
		e.onElementMounted(child)
	}
}

// ==================== 帧循环 ====================

func (e *Engine) Update() error {
	now := time.Now()
	deltaTime := float32(now.Sub(e.lastFrameTime).Seconds())
	e.lastFrameTime = now

	// 1. 处理请求重建的 Widget（结构变化）
	e.flushBuildQueue()

	// 2. 批量刷新脏 Element（测量 → 布局 → 清除标记）
	e.flushDirtyElements()

	// 3. 事件处理
	e.handleEvents()

	// 4. 每帧更新 Element（动画等）
	if e.rootElement != nil {
		e.updateElements(e.rootElement)
	}

	_ = deltaTime
	return nil
}

func (e *Engine) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 245, G: 245, B: 245, A: 255})
	if e.rootElement != nil {
		e.drawElement(screen, e.rootElement)
	}
}

// ==================== 结构重建 ====================

func (e *Engine) scheduleBuild(w Widget) {
	e.buildQueue = append(e.buildQueue, w)
}

func (e *Engine) flushBuildQueue() {
	if len(e.buildQueue) == 0 {
		return
	}
	queue := e.buildQueue
	e.buildQueue = e.buildQueue[:0]

	for _, w := range queue {
		newRoot := w.Build()
		if newRoot == nil {
			continue
		}
		if e.rootElement == nil {
			e.rootElement = newRoot
			e.onElementMounted(newRoot)
		} else {
			// 同级浅 Diff：复用或替换根节点
			e.patchElement(e.rootElement, newRoot)
		}
	}

	// 重建后布局可能变化
	if e.hasLayoutDirty() {
		e.calculateLayout()
	}
}

// patchElement 对新旧 Element 做同级浅对比。
// 目前只处理根级别替换，完整 Diff 后续扩展。
func (e *Engine) patchElement(oldEl, newEl Element) {
	if oldEl.ElementType() == newEl.ElementType() {
		// 同类型：复用旧节点，同步 Yoga 样式
		if oldYoga, newYoga := oldEl.GetYoga(), newEl.GetYoga(); oldYoga != nil && newYoga != nil {
			oldYoga.CopyStyleFrom(newYoga)
		}
		// TODO: 同步子节点
		oldEl.Mark(FlagNeedLayout | FlagNeedDraw)
	} else {
		// 类型不同：替换根节点
		if oldEl.GetParent() != nil {
			parent := oldEl.GetParent()
			parent.RemoveChild(oldEl)
			parent.AppendChild(newEl)
		} else {
			e.rootElement = newEl
			e.onElementMounted(newEl)
		}
	}
}

// ==================== 脏标记刷新 ====================

func (e *Engine) markDirty(el Element) {
	if el == nil {
		return
	}
	// 避免重复加入
	for _, d := range e.dirtyElements {
		if d == el {
			return
		}
	}
	e.dirtyElements = append(e.dirtyElements, el)
}

func (e *Engine) flushDirtyElements() {
	if len(e.dirtyElements) == 0 {
		return
	}

	// 1. 先测量（文字排版等）
	for _, el := range e.dirtyElements {
		if el.GetFlags()&FlagNeedMeasure != 0 {
			if m, ok := el.(interface{ FlushMeasure() }); ok {
				m.FlushMeasure()
			}
		}
	}

	// 2. 再布局（只要有一个节点标记了 NeedLayout，就从根重算）
	if e.hasLayoutDirty() {
		e.calculateLayout()
	}

	// 3. 清除所有标记（只清脏标记，保留持久状态）
	for _, el := range e.dirtyElements {
		el.ClearDirty()
	}
	e.dirtyElements = e.dirtyElements[:0]
}

func (e *Engine) hasLayoutDirty() bool {
	for _, el := range e.dirtyElements {
		if el.GetFlags()&FlagNeedLayout != 0 {
			return true
		}
	}
	return false
}

// ==================== 布局 ====================

func (e *Engine) calculateLayout() {
	if e.rootElement == nil || e.rootElement.GetYoga() == nil {
		return
	}
	rootYoga := e.rootElement.GetYoga()
	rootYoga.StyleSetWidth(float32(e.screenWidth))
	rootYoga.StyleSetHeight(float32(e.screenHeight))
	rootYoga.CalculateLayout(float32(e.screenWidth), float32(e.screenHeight), yoga.DirectionLTR)

	if e.rootElement != nil {
		e.updateBounds(e.rootElement, 0, 0)
	}
}

func (e *Engine) updateBounds(el Element, parentX, parentY float32) {
	if el == nil || el.GetYoga() == nil {
		return
	}
	y := el.GetYoga()
	el.SetBounds(LayoutBounds{
		X:      parentX + y.LayoutLeft(),
		Y:      parentY + y.LayoutTop(),
		Width:  y.LayoutWidth(),
		Height: y.LayoutHeight(),
	})
	for _, child := range el.GetChildren() {
		e.updateBounds(child, el.GetBounds().X, el.GetBounds().Y)
	}
}

// ==================== 绘制 ====================

func (e *Engine) drawElement(screen *ebiten.Image, el Element) {
	if el == nil {
		return
	}
	bounds := el.GetBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	// TODO: clip
	el.Draw(screen)

	for _, child := range el.GetChildren() {
		e.drawElement(screen, child)
	}
}

// ==================== Element 每帧更新 ====================

func (e *Engine) updateElements(el Element) {
	if el == nil {
		return
	}
	_ = el.Update()
	for _, child := range el.GetChildren() {
		e.updateElements(child)
	}
}

// ==================== 事件（简化版）====================

func (e *Engine) handleEvents() {
	mx, my := ebiten.CursorPosition()
	fx, fy := float32(mx), float32(my)

	// Hover / Click 简化处理
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		target := e.hitTest(e.rootElement, fx, fy)
		if target != nil {
			target.HandleEvent(&Event{Type: EventClick, X: fx, Y: fy, Target: target})
		}
	}

	e.lastMouseX = fx
	e.lastMouseY = fy
}

func (e *Engine) hitTest(el Element, x, y float32) Element {
	if el == nil {
		return nil
	}
	// 后序遍历：子节点优先
	children := el.GetChildren()
	for i := len(children) - 1; i >= 0; i-- {
		if h := e.hitTest(children[i], x, y); h != nil {
			return h
		}
	}
	b := el.GetBounds()
	if x >= b.X && x < b.X+b.Width && y >= b.Y && y < b.Y+b.Height {
		return el
	}
	return nil
}
