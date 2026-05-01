package ui

import (
	"image/color"
	"sort"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/yoga"
)

// defaultEngine 是当前运行的全局 Engine 实例，供 Rebuild() 等快捷函数使用。
var defaultEngine *Engine

// DefaultEngine 返回当前全局 Engine 实例。
func DefaultEngine() *Engine {
	return defaultEngine
}

// Engine is the core engine of Tenon v2, managing Widget, Element, and RenderObject trees.
type Engine struct {
	buildFunc BuildFunc

	rootWidget       Widget
	rootElement      Element
	rootRenderObject render.RenderObject

	owner *render.PipelineOwner

	screenWidth  int
	screenHeight int

	needsBuild    bool
	dirtyElements []Element
	building      bool

	activeAnimations []*AnimationController
	lastTickTime     time.Time

	lastMouseX      float32
	lastMouseY      float32
	mouseDown       bool
	mouseDownTarget render.RenderObject
	hoverTarget     render.RenderObject
	focusTarget     render.RenderObject
}

// NewEngine creates a new Engine.
func NewEngine(buildFunc BuildFunc, width, height int) *Engine {
	e := &Engine{
		buildFunc:    buildFunc,
		screenWidth:  width,
		screenHeight: height,
		owner:        render.NewPipelineOwner(),
	}
	defaultEngine = e
	return e
}

// Rebuild marks the engine to rebuild on the next frame.
func (e *Engine) Rebuild() {
	e.needsBuild = true
}

// RebuildDefault marks the default engine to rebuild on the next frame.
func RebuildDefault() {
	if defaultEngine != nil {
		defaultEngine.Rebuild()
	}
}

// scheduleBuildFor 将指定 Element 加入 dirty 队列，供局部 rebuild 使用。
func (e *Engine) scheduleBuildFor(element Element) {
	if e.building {
		// 如果在 build 过程中触发，直接加入 dirty 队列，下一轮处理
		e.dirtyElements = append(e.dirtyElements, element)
		return
	}
	e.dirtyElements = append(e.dirtyElements, element)
}

func (e *Engine) GetRootRenderObject() render.RenderObject {
	return e.rootRenderObject
}

func (e *Engine) GetRootElement() Element {
	return e.rootElement
}

// Mount mounts the root widget and triggers the first build.
func (e *Engine) Mount() {
	e.needsBuild = true
	e.flushBuild()
}

// ==================== ebiten.Game interface ====================

func (e *Engine) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth != e.screenWidth || outsideHeight != e.screenHeight {
		e.screenWidth = outsideWidth
		e.screenHeight = outsideHeight
		if e.rootRenderObject != nil {
			e.rootRenderObject.MarkNeedsLayout()
		}
	}
	return outsideWidth, outsideHeight
}

func (e *Engine) RegisterAnimation(a *AnimationController) {
	for _, existing := range e.activeAnimations {
		if existing == a {
			return
		}
	}
	e.activeAnimations = append(e.activeAnimations, a)
}

func (e *Engine) UnregisterAnimation(a *AnimationController) {
	for i, existing := range e.activeAnimations {
		if existing == a {
			e.activeAnimations = append(e.activeAnimations[:i], e.activeAnimations[i+1:]...)
			return
		}
	}
}

func (e *Engine) Update() error {
	now := time.Now()
	if e.lastTickTime.IsZero() {
		e.lastTickTime = now
	}
	dt := now.Sub(e.lastTickTime)
	e.lastTickTime = now

	// Tick active animations
	if len(e.activeAnimations) > 0 {
		for _, a := range e.activeAnimations {
			a.Tick(dt)
		}
	}

	if e.needsBuild || len(e.dirtyElements) > 0 {
		e.flushBuild()
	}

	if e.rootRenderObject != nil && e.rootRenderObject.GetYoga() != nil {
		rootYoga := e.rootRenderObject.GetYoga()
		rootYoga.StyleSetWidth(float32(e.screenWidth))
		rootYoga.StyleSetHeight(float32(e.screenHeight))
		rootYoga.CalculateLayout(float32(e.screenWidth), float32(e.screenHeight), yoga.DirectionLTR)
	}

	e.owner.FlushLayout()

	if e.rootRenderObject != nil {
		e.syncBounds(e.rootRenderObject)
	}

	e.owner.FlushPaint()

	e.handleMouseInput()

	e.handleKeyboardInput()

	e.tickBlink()

	return nil
}

func (e *Engine) handleMouseInput() {
	mx, my := ebiten.CursorPosition()
	e.lastMouseX = float32(mx)
	e.lastMouseY = float32(my)

	// Wheel
	wheelX, wheelY := ebiten.Wheel()
	if wheelY != 0 || wheelX != 0 {
		if e.rootRenderObject != nil {
			target := e.hitTest(e.rootRenderObject, e.lastMouseX, e.lastMouseY)
			if target != nil {
				if scroll := e.findScrollParent(target); scroll != nil {
					if s, ok := scroll.(interface {
						ScrollBy(dx, dy float32)
					}); ok {
						s.ScrollBy(-float32(wheelX)*30, -float32(wheelY)*30)
					}
				}
			}
		}
	}

	// Hover
	if e.rootRenderObject != nil {
		target := e.hitTest(e.rootRenderObject, e.lastMouseX, e.lastMouseY)
		if target != e.hoverTarget {
			if e.hoverTarget != nil && e.hoverTarget.GetOnMouseLeave() != nil {
				e.hoverTarget.GetOnMouseLeave()()
			}
			if target != nil && target.GetOnMouseEnter() != nil {
				target.GetOnMouseEnter()()
			}
			e.hoverTarget = target
		}
	}

	isPressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	if isPressed && !e.mouseDown {
		// Mouse down
		e.mouseDown = true
		if e.rootRenderObject != nil {
			target := e.hitTest(e.rootRenderObject, e.lastMouseX, e.lastMouseY)
			e.mouseDownTarget = target
			if target != nil && target.GetOnMouseDown() != nil {
				target.GetOnMouseDown()()
			}
		}
	} else if isPressed && e.mouseDown {
		// Drag
		if e.mouseDownTarget != nil {
			e.mouseDownTarget.HandleDrag(e.lastMouseX, e.lastMouseY)
		}
	} else if !isPressed && e.mouseDown {
		// Mouse up / click
		e.mouseDown = false
		if e.mouseDownTarget != nil && e.rootRenderObject != nil {
			target := e.hitTest(e.rootRenderObject, e.lastMouseX, e.lastMouseY)
			if target == e.mouseDownTarget && target.GetOnClick() != nil {
				target.GetOnClick()()
			}
			if e.mouseDownTarget.GetOnMouseUp() != nil {
				e.mouseDownTarget.GetOnMouseUp()()
			}
			e.updateFocus(target)
		}
		e.mouseDownTarget = nil
	}
}

func (e *Engine) updateFocus(target render.RenderObject) {
	focuser := e.findFocuser(target)
	if focuser != nil {
		if e.focusTarget != focuser {
			if e.focusTarget != nil {
				if b, ok := e.focusTarget.(interface{ Blur() }); ok {
					b.Blur()
				}
			}
			if f, ok := focuser.(interface{ Focus() }); ok {
				f.Focus()
			}
			e.focusTarget = focuser
		}
	} else {
		if e.focusTarget != nil {
			if b, ok := e.focusTarget.(interface{ Blur() }); ok {
				b.Blur()
			}
			e.focusTarget = nil
		}
	}
}

func (e *Engine) findFocuser(ro render.RenderObject) render.RenderObject {
	for ro != nil {
		if _, ok := ro.(interface {
			Focus()
			Blur()
			IsFocused() bool
		}); ok {
			return ro
		}
		ro = ro.GetParent()
	}
	return nil
}

func (e *Engine) handleKeyboardInput() {
	if e.focusTarget == nil {
		return
	}
	f, ok := e.focusTarget.(interface {
		HandleInput(chars []rune, backspace, enter, left, right, home, end, selectAll bool)
	})
	if !ok {
		return
	}

	chars := ebiten.AppendInputChars(nil)

	backspace := inpututil.IsKeyJustPressed(ebiten.KeyBackspace)
	if inpututil.KeyPressDuration(ebiten.KeyBackspace) > 30 && inpututil.KeyPressDuration(ebiten.KeyBackspace)%5 == 0 {
		backspace = true
	}

	enter := inpututil.IsKeyJustPressed(ebiten.KeyEnter)
	left := inpututil.IsKeyJustPressed(ebiten.KeyLeft)
	right := inpututil.IsKeyJustPressed(ebiten.KeyRight)
	if inpututil.KeyPressDuration(ebiten.KeyLeft) > 30 && inpututil.KeyPressDuration(ebiten.KeyLeft)%5 == 0 {
		left = true
	}
	if inpututil.KeyPressDuration(ebiten.KeyRight) > 30 && inpututil.KeyPressDuration(ebiten.KeyRight)%5 == 0 {
		right = true
	}

	home := inpututil.IsKeyJustPressed(ebiten.KeyHome)
	end := inpututil.IsKeyJustPressed(ebiten.KeyEnd)
	selectAll := inpututil.IsKeyJustPressed(ebiten.KeyA) && ebiten.IsKeyPressed(ebiten.KeyControl)

	if len(chars) > 0 || backspace || enter || left || right || home || end || selectAll {
		f.HandleInput(chars, backspace, enter, left, right, home, end, selectAll)
	}
}

func (e *Engine) tickBlink() {
	if e.focusTarget == nil {
		return
	}
	if b, ok := e.focusTarget.(interface{ TickBlink() }); ok {
		b.TickBlink()
	}
}

func (e *Engine) findScrollParent(ro render.RenderObject) render.RenderObject {
	for ro != nil {
		if _, ok := ro.(interface{ GetScrollOffset() render.Offset }); ok {
			return ro
		}
		ro = ro.GetParent()
	}
	return nil
}

func (e *Engine) hitTest(ro render.RenderObject, x, y float32) render.RenderObject {
	if ro == nil || !ro.IsVisible() {
		return nil
	}
	if !ro.HitTest(x, y) {
		return nil
	}

	scrollX, scrollY := float32(0), float32(0)
	if scroller, ok := ro.(interface{ GetScrollOffset() render.Offset }); ok {
		so := scroller.GetScrollOffset()
		scrollX, scrollY = so.X, so.Y
	}

	bounds := ro.GetBounds()
	for _, child := range ro.GetChildren() {
		if result := e.hitTest(child, x-bounds.X+scrollX, y-bounds.Y+scrollY); result != nil {
			return result
		}
	}
	return ro
}

func (e *Engine) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	if e.rootRenderObject != nil {
		e.drawRenderObject(screen, e.rootRenderObject, render.Offset{})
	}
}

// ==================== Widget diff & build ====================

func (e *Engine) flushBuild() {
	if e.building {
		return
	}
	e.building = true
	defer func() { e.building = false }()

	// 循环处理，因为 rebuild 可能产生新的 dirty elements
	// 限制循环次数防止无限递归
	for loop := 0; loop < 10; loop++ {
		hasWork := len(e.dirtyElements) > 0 || e.needsBuild
		if !hasWork {
			break
		}

		// 1. 处理局部 dirty elements（StatefulElement 的 setState）
		dirtyList := e.dirtyElements
		e.dirtyElements = nil
		for _, el := range dirtyList {
			if se, ok := el.(*StatefulElement); ok {
				se.rebuild()
			}
		}

		// 2. 处理全局 rebuild（buildFunc + Rebuild）
		if e.needsBuild {
			e.needsBuild = false
			e.flushGlobalBuild()
		}
	}
}

func (e *Engine) flushGlobalBuild() {
	newRootWidget := e.buildFunc()
	if newRootWidget == nil {
		return
	}

	if e.rootElement == nil {
		// First mount
		e.rootWidget = newRootWidget
		e.rootElement = newRootWidget.CreateElement()
		e.rootElement.Mount(nil, 0)
		if e.rootElement != nil {
			e.rootRenderObject = e.rootElement.FindRenderObject()
			if e.rootRenderObject != nil {
				e.rootRenderObject.Attach(e.owner)
				e.rootRenderObject.MarkNeedsLayout()
			}
		}
	} else {
		// Diff update
		if CanUpdate(e.rootWidget, newRootWidget) {
			e.rootWidget = newRootWidget
			e.rootElement.Update(newRootWidget)
		} else {
			// Type changed, remount
			e.rootElement.Unmount()
			e.rootWidget = newRootWidget
			e.rootElement = newRootWidget.CreateElement()
			e.rootElement.Mount(nil, 0)
			if e.rootElement != nil {
				e.rootRenderObject = e.rootElement.FindRenderObject()
				if e.rootRenderObject != nil {
					e.rootRenderObject.Attach(e.owner)
					e.rootRenderObject.MarkNeedsLayout()
				}
			}
		}
	}
}

// ==================== Sync bounds ====================

func (e *Engine) syncBounds(ro render.RenderObject) {
	if ro == nil {
		return
	}
	y := ro.GetYoga()
	if y == nil {
		return
	}
	bounds := render.Bounds{
		X:      y.LayoutLeft(),
		Y:      y.LayoutTop(),
		Width:  y.LayoutWidth(),
		Height: y.LayoutHeight(),
	}
	ro.SetBounds(bounds)
	for _, child := range ro.GetChildren() {
		e.syncBounds(child)
	}
}

// ==================== Draw ====================

func (e *Engine) drawRenderObject(screen *ebiten.Image, ro render.RenderObject, parentOffset render.Offset) {
	if ro == nil || !ro.IsVisible() {
		return
	}
	bounds := ro.GetBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	ro.Paint(screen, parentOffset)

	absX := parentOffset.X + bounds.X
	absY := parentOffset.Y + bounds.Y

	scrollX, scrollY := float32(0), float32(0)
	if scroller, ok := ro.(interface{ GetScrollOffset() render.Offset }); ok {
		so := scroller.GetScrollOffset()
		scrollX, scrollY = so.X, so.Y
	}

	childScreen := screen
	childOffset := render.Offset{X: absX - scrollX, Y: absY - scrollY}
	if clipper, ok := ro.(interface{ ClipChildren() bool }); ok && clipper.ClipChildren() {
		sub := render.SubImage(screen,
			int(absX), int(absY),
			int(bounds.Width), int(bounds.Height))
		if sub != nil {
			childScreen = sub
			childOffset = render.Offset{X: -scrollX, Y: -scrollY}
		}
	}

	children := ro.GetChildren()
	if len(children) > 1 {
		children = append([]render.RenderObject(nil), children...)
		sort.SliceStable(children, func(i, j int) bool {
			return children[i].GetZIndex() < children[j].GetZIndex()
		})
	}
	for _, child := range children {
		e.drawRenderObject(childScreen, child, childOffset)
	}
}

// ==================== Run ====================

// Run starts the ebiten game loop.
func (e *Engine) Run() {
	e.Mount()
	ebiten.SetWindowSize(e.screenWidth, e.screenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(e); err != nil {
		panic(err)
	}
}

// ==================== EdgeInsets ====================

type EdgeInsets struct {
	Top, Right, Bottom, Left float32
}

func (e EdgeInsets) All() float32 {
	return e.Top + e.Right + e.Bottom + e.Left
}

func EdgeInsetsAll(v float32) EdgeInsets {
	return EdgeInsets{Top: v, Right: v, Bottom: v, Left: v}
}

func EdgeInsetsSymmetric(horizontal, vertical float32) EdgeInsets {
	return EdgeInsets{Top: vertical, Right: horizontal, Bottom: vertical, Left: horizontal}
}

func EdgeInsetsOnly(top, right, bottom, left float32) EdgeInsets {
	return EdgeInsets{Top: top, Right: right, Bottom: bottom, Left: left}
}

// Yoga type aliases
type (
	FlexDirection = yoga.FlexDirection
	Justify       = yoga.Justify
	Align         = yoga.Align
	Wrap          = yoga.Wrap
	PositionType  = yoga.PositionType
	Display       = yoga.Display
	Overflow      = yoga.Overflow
	Edge          = yoga.Edge
	Gutter        = yoga.Gutter
	Direction     = yoga.Direction
)

const (
	FlexDirectionColumn    = yoga.FlexDirectionColumn
	FlexDirectionRow       = yoga.FlexDirectionRow
	JustifyFlexStart       = yoga.JustifyFlexStart
	JustifyCenter          = yoga.JustifyCenter
	JustifyFlexEnd         = yoga.JustifyFlexEnd
	JustifySpaceBetween    = yoga.JustifySpaceBetween
	JustifySpaceAround     = yoga.JustifySpaceAround
	JustifySpaceEvenly     = yoga.JustifySpaceEvenly
	AlignAuto              = yoga.AlignAuto
	AlignFlexStart         = yoga.AlignFlexStart
	AlignCenter            = yoga.AlignCenter
	AlignFlexEnd           = yoga.AlignFlexEnd
	AlignStretch           = yoga.AlignStretch
	AlignBaseline          = yoga.AlignBaseline
	WrapNoWrap             = yoga.WrapNoWrap
	WrapWrap               = yoga.WrapWrap
	PositionTypeRelative   = yoga.PositionTypeRelative
	PositionTypeAbsolute   = yoga.PositionTypeAbsolute
	DisplayFlex            = yoga.DisplayFlex
	DisplayNone            = yoga.DisplayNone
	OverflowVisible        = yoga.OverflowVisible
	OverflowHidden         = yoga.OverflowHidden
	OverflowScroll         = yoga.OverflowScroll
	EdgeLeft               = yoga.EdgeLeft
	EdgeTop                = yoga.EdgeTop
	EdgeRight              = yoga.EdgeRight
	EdgeBottom             = yoga.EdgeBottom
	EdgeStart              = yoga.EdgeStart
	EdgeEnd                = yoga.EdgeEnd
	EdgeHorizontal         = yoga.EdgeHorizontal
	EdgeVertical           = yoga.EdgeVertical
	EdgeAll                = yoga.EdgeAll
	GutterColumn           = yoga.GutterColumn
	GutterRow              = yoga.GutterRow
	GutterAll              = yoga.GutterAll
	DirectionLTR           = yoga.DirectionLTR
	DirectionRTL           = yoga.DirectionRTL
)
