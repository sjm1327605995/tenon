package core

import (
	"fmt"
	"image/color"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/yoga"
)

type LayoutBounds struct {
	X, Y, Width, Height float32
}

type Element struct {
	Visible         bool
	Yoga            *yoga.Node
	BackgroundColor color.Color
	BorderColor     color.Color
	ShadowColor     color.Color
	ShadowBlur      float32
	ShadowOffsetX   float32
	ShadowOffsetY   float32
	PointerEvents   PointerEvents
	BorderRadius    BorderRadius
}

type BorderRadius struct {
	TopLeft     float32
	TopRight    float32
	BottomRight float32
	BottomLeft  float32
}

type PointerEvents int

const (
	PointerEventsAuto PointerEvents = iota
	PointerEventsNone
)

type ErrorInfo struct {
	ComponentStack string
}

// ============================================================================
// 【用户开放接口】Lifecycle - 用户自定义组件时可重写的生命周期钩子
// ============================================================================
type Lifecycle interface {
	// ComponentDidMount 组件挂载后执行（用户重写）
	// 适合数据请求、订阅、DOM 操作等初始化逻辑
	ComponentDidMount()

	// ComponentDidUpdate 组件更新后执行（用户重写）
	// prevProps: 更新前的 props
	// prevState: 更新前的 state  
	// snapshot: GetSnapshotBeforeUpdate 返回的值
	ComponentDidUpdate(prevProps interface{}, prevState interface{}, snapshot interface{})

	// ComponentWillUnmount 组件卸载前执行（用户重写）
	// 适合清理定时器、取消订阅等资源释放操作
	ComponentWillUnmount()

	// ShouldComponentUpdate 性能优化钩子（用户重写）
	// 返回 false 可阻止不必要的重新渲染
	ShouldComponentUpdate(nextProps interface{}, nextState interface{}) bool

	// GetSnapshotBeforeUpdate 在渲染提交到屏幕前调用（用户重写）
	// 用于捕获 DOM 信息（如滚动位置），返回值将传递给 ComponentDidUpdate
	GetSnapshotBeforeUpdate(prevProps interface{}, prevState interface{}) interface{}

	// GetDerivedStateFromProps 静态方法，在每次渲染前调用（用户重写）
	// 用于根据 props 同步更新 state
	GetDerivedStateFromProps(props interface{}, state interface{}) interface{}

	// ComponentDidCatch 错误边界钩子（用户重写）
	// 用于捕获子组件渲染错误，记录日志或显示备用 UI
	ComponentDidCatch(err error, info ErrorInfo)

	// GetDerivedStateFromError 静态方法，渲染错误时调用（用户重写）
	// 用于更新 state 显示备用 UI
	GetDerivedStateFromError(err error) interface{}
}

// ============================================================================
// 【用户开放接口】Stateful - 用户组件的状态管理接口
// ============================================================================
type Stateful interface {
	GetProps() interface{}
	GetState() interface{}
	SetProps(props interface{})
	SetState(state interface{})
}

// ============================================================================
// 【框架内部接口】InternalComponent - 框架内部使用的核心接口
// 注意：这些方法主要由框架调用，用户通常不需要直接调用或重写
// ============================================================================
type InternalComponent interface {
	// 渲染相关
	Draw(screen *ebiten.Image)
	DrawOverlay(screen *ebiten.Image)
	Update() error

	// 输入处理
	HandleInput() bool

	// 子组件管理
	GetChildren() []Component
	AddChild(child Component) error

	// 布局相关
	GetLayoutBounds() LayoutBounds
	GetElement() *Element
	Render() *Element
	UpdateLayoutBounds(parentX, parentY float32)

	// 脏检查机制
	ID() string
	MarkDirty()
	IsDirty() bool
	ClearDirty()
}

// ============================================================================
// 【综合接口】Component - 完整的组件接口（用户使用）
// ============================================================================
type Component interface {
	InternalComponent
	Lifecycle
	Stateful
}

// ============================================================================
// BaseComponent - 基础组件实现（用户组合使用）
// 提供所有接口的默认实现，用户只需重写需要的方法
// ============================================================================
type BaseComponent struct {
	// 核心属性
	children []Component
	element  *Element
	yogaNode *yoga.Node
	bounds   LayoutBounds

	// 内部状态
	id    string
	dirty bool
	mu    sync.RWMutex

	// Hooks 支持
	hooksRenderFunc func()

	// Self 引用（用于链式方法返回正确类型）
	self Component

	// 生命周期状态
	isMounted        bool
	prevProps        interface{}
	prevState        interface{}
	errorState       interface{}
	hasError         bool
}

// NewBaseComponent 创建一个新的 BaseComponent（用户调用）
func NewBaseComponent() BaseComponent {
	node := yoga.NewNode()
	return BaseComponent{
		children: make([]Component, 0),
		element: &Element{
			Visible:       true,
			Yoga:          node,
			PointerEvents: PointerEventsAuto,
		},
		yogaNode:  node,
		id:        generateComponentID(),
		dirty:     false,
		isMounted: false,
	}
}

// Init 初始化 self 引用（用户在构造函数中调用）
// 使用示例:
// func NewMyComponent() *MyComponent {
//     c := &MyComponent{BaseComponent: core.NewBaseComponent()}
//     c.Init(c)
//     return c
// }
func (b *BaseComponent) Init(self Component) {
	b.self = self
}

// Self 返回组件自身引用（框架内部使用）
func (b *BaseComponent) Self() Component {
	return b.self
}

// ----------------------------------------------------------------------------
// 【框架内部方法】InternalComponent 接口实现
// ----------------------------------------------------------------------------

func (b *BaseComponent) GetChildren() []Component {
	return b.children
}

func (b *BaseComponent) AddChild(child Component) error {
	childYoga := child.GetElement().Yoga
	if childYoga != nil && b.yogaNode != nil {
		if childYoga.GetOwner() != nil {
			childYoga.GetOwner().RemoveChild(childYoga)
		}
		b.yogaNode.InsertChild(childYoga, b.yogaNode.GetChildCount())
	}
	b.children = append(b.children, child)
	return nil
}

func (b *BaseComponent) GetLayoutBounds() LayoutBounds {
	return b.bounds
}

func (b *BaseComponent) GetElement() *Element {
	return b.element
}

func (b *BaseComponent) Render() *Element {
	return b.element
}

func (b *BaseComponent) UpdateLayoutBounds(parentX, parentY float32) {
	if b.yogaNode == nil {
		return
	}

	relX := b.yogaNode.LayoutLeft()
	relY := b.yogaNode.LayoutTop()

	b.bounds = LayoutBounds{
		X:      parentX + relX,
		Y:      parentY + relY,
		Width:  b.yogaNode.LayoutWidth(),
		Height: b.yogaNode.LayoutHeight(),
	}
}

func (b *BaseComponent) Update() error {
	for _, child := range b.children {
		if err := child.Update(); err != nil {
			return err
		}
	}
	return nil
}

func (b *BaseComponent) Draw(screen *ebiten.Image) {
}

func (b *BaseComponent) DrawOverlay(screen *ebiten.Image) {
}

func (b *BaseComponent) HandleInput() bool {
	return false
}

// CalculateLayout 计算布局（框架内部调用）
func CalculateLayout(root Component, width, height float32) {
	element := root.GetElement()
	if element == nil || element.Yoga == nil {
		return
	}

	element.Yoga.StyleSetWidth(width)
	element.Yoga.StyleSetHeight(height)
	element.Yoga.CalculateLayout(width, height, yoga.DirectionLTR)
	updateAllBounds(root, 0, 0)
}

func updateAllBounds(component Component, parentX, parentY float32) {
	if component == nil {
		return
	}

	component.UpdateLayoutBounds(parentX, parentY)
	bounds := component.GetLayoutBounds()

	for _, child := range component.GetChildren() {
		updateAllBounds(child, bounds.X, bounds.Y)
	}
}

var (
	componentCounter int64
	componentMutex   sync.Mutex
)

func generateComponentID() string {
	componentMutex.Lock()
	defer componentMutex.Unlock()

	componentCounter++
	return fmt.Sprintf("component_%d_%d", time.Now().UnixNano(), componentCounter)
}

func (b *BaseComponent) ID() string {
	return b.id
}

func (b *BaseComponent) MarkDirty() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.dirty = true
}

func (b *BaseComponent) IsDirty() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.dirty
}

func (b *BaseComponent) ClearDirty() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.dirty = false
}

// ----------------------------------------------------------------------------
// 【用户开放方法】样式设置方法（链式调用）
// ----------------------------------------------------------------------------

func (b *BaseComponent) ApplyStyle(style *Style) {
	if style != nil && b.yogaNode != nil {
		style.Apply(b.yogaNode)
	}
}

func (b *BaseComponent) ApplyVisualStyle(style *VisualStyle) {
	if style != nil && b.element != nil {
		style.Apply(b.element)
	}
}

func (b *BaseComponent) SetHooksRenderFunc(fn func()) {
	b.hooksRenderFunc = fn
}

func (b *BaseComponent) GetHooksRenderFunc() func() {
	return b.hooksRenderFunc
}

func (b *BaseComponent) SetWidth(width float32) Component {
	b.element.Yoga.StyleSetWidth(width)
	return b.self
}

func (b *BaseComponent) SetHeight(height float32) Component {
	b.element.Yoga.StyleSetHeight(height)
	return b.self
}

func (b *BaseComponent) SetMinWidth(width float32) Component {
	b.element.Yoga.StyleSetMinWidth(width)
	return b.self
}

func (b *BaseComponent) SetMinHeight(height float32) Component {
	b.element.Yoga.StyleSetMinHeight(height)
	return b.self
}

func (b *BaseComponent) SetMaxWidth(width float32) Component {
	b.element.Yoga.StyleSetMaxWidth(width)
	return b.self
}

func (b *BaseComponent) SetMaxHeight(height float32) Component {
	b.element.Yoga.StyleSetMaxHeight(height)
	return b.self
}

func (b *BaseComponent) SetFlexDirection(dir yoga.FlexDirection) Component {
	b.element.Yoga.StyleSetFlexDirection(dir)
	return b.self
}

func (b *BaseComponent) SetJustifyContent(justify yoga.Justify) Component {
	b.element.Yoga.StyleSetJustifyContent(justify)
	return b.self
}

func (b *BaseComponent) SetAlignItems(align yoga.Align) Component {
	b.element.Yoga.StyleSetAlignItems(align)
	return b.self
}

func (b *BaseComponent) SetAlignSelf(align yoga.Align) Component {
	b.element.Yoga.StyleSetAlignSelf(align)
	return b.self
}

func (b *BaseComponent) SetFlexGrow(grow float32) Component {
	b.element.Yoga.StyleSetFlexGrow(grow)
	return b.self
}

func (b *BaseComponent) SetFlexShrink(shrink float32) Component {
	b.element.Yoga.StyleSetFlexShrink(shrink)
	return b.self
}

func (b *BaseComponent) SetFlexBasis(basis float32) Component {
	b.element.Yoga.StyleSetFlexBasis(basis)
	return b.self
}

func (b *BaseComponent) SetFlexWrap(wrap yoga.Wrap) Component {
	b.element.Yoga.StyleSetFlexWrap(wrap)
	return b.self
}

func (b *BaseComponent) SetPadding(edge yoga.Edge, value float32) Component {
	b.element.Yoga.StyleSetPadding(edge, value)
	return b.self
}

func (b *BaseComponent) SetMargin(edge yoga.Edge, value float32) Component {
	b.element.Yoga.StyleSetMargin(edge, value)
	return b.self
}

func (b *BaseComponent) SetBorder(edge yoga.Edge, value float32) Component {
	b.element.Yoga.StyleSetBorder(edge, value)
	return b.self
}

func (b *BaseComponent) SetPosition(edge yoga.Edge, value float32) Component {
	b.element.Yoga.StyleSetPosition(edge, value)
	return b.self
}

func (b *BaseComponent) SetPositionType(positionType yoga.PositionType) Component {
	b.element.Yoga.StyleSetPositionType(positionType)
	return b.self
}

func (b *BaseComponent) SetAspectRatio(ratio float32) Component {
	b.element.Yoga.StyleSetAspectRatio(ratio)
	return b.self
}

func (b *BaseComponent) SetDisplay(display yoga.Display) Component {
	b.element.Yoga.StyleSetDisplay(display)
	return b.self
}

func (b *BaseComponent) SetOverflow(overflow yoga.Overflow) Component {
	b.element.Yoga.StyleSetOverflow(overflow)
	return b.self
}

func (b *BaseComponent) SetBackgroundColor(clr color.Color) Component {
	b.element.BackgroundColor = clr
	return b.self
}

func (b *BaseComponent) SetBorderColor(clr color.Color) Component {
	b.element.BorderColor = clr
	return b.self
}

func (b *BaseComponent) SetBorderRadius(radius float32) Component {
	b.element.BorderRadius = BorderRadius{
		TopLeft:     radius,
		TopRight:    radius,
		BottomRight: radius,
		BottomLeft:  radius,
	}
	return b.self
}

func (b *BaseComponent) SetBorderRadius4(topLeft, topRight, bottomRight, bottomLeft float32) Component {
	b.element.BorderRadius = BorderRadius{
		TopLeft:     topLeft,
		TopRight:    topRight,
		BottomRight: bottomRight,
		BottomLeft:  bottomLeft,
	}
	return b.self
}

func (b *BaseComponent) SetShadow(color color.Color, blur, offsetX, offsetY float32) Component {
	b.element.ShadowColor = color
	b.element.ShadowBlur = blur
	b.element.ShadowOffsetX = offsetX
	b.element.ShadowOffsetY = offsetY
	return b.self
}

func (b *BaseComponent) SetVisible(visible bool) Component {
	b.element.Visible = visible
	return b.self
}

func (b *BaseComponent) SetPointerEvents(pe PointerEvents) Component {
	b.element.PointerEvents = pe
	return b.self
}

func (b *BaseComponent) Add(children ...Component) Component {
	for _, child := range children {
		_ = b.AddChild(child)
	}
	return b.self
}

func (b *BaseComponent) SetStyle(style *Style) Component {
	b.ApplyStyle(style)
	return b.self
}

func (b *BaseComponent) SetVisualStyle(style *VisualStyle) Component {
	b.ApplyVisualStyle(style)
	return b.self
}

// ----------------------------------------------------------------------------
// 【用户开放方法】Lifecycle 接口默认实现
// ----------------------------------------------------------------------------

func (b *BaseComponent) ComponentDidMount() {
	b.isMounted = true
}

func (b *BaseComponent) ComponentDidUpdate(prevProps interface{}, prevState interface{}, snapshot interface{}) {
}

func (b *BaseComponent) ComponentWillUnmount() {
	b.isMounted = false
}

func (b *BaseComponent) ShouldComponentUpdate(nextProps interface{}, nextState interface{}) bool {
	return b.isMounted
}

func (b *BaseComponent) GetSnapshotBeforeUpdate(prevProps interface{}, prevState interface{}) interface{} {
	return nil
}

func (b *BaseComponent) GetDerivedStateFromProps(props interface{}, state interface{}) interface{} {
	return state
}

func (b *BaseComponent) ComponentDidCatch(err error, info ErrorInfo) {
}

func (b *BaseComponent) GetDerivedStateFromError(err error) interface{} {
	return nil
}

// ----------------------------------------------------------------------------
// 【用户开放方法】Stateful 接口实现
// ----------------------------------------------------------------------------

func (b *BaseComponent) IsMounted() bool {
	return b.isMounted
}

func (b *BaseComponent) SetProps(props interface{}) {
	b.prevProps = props
}

func (b *BaseComponent) SetState(state interface{}) {
	b.prevState = state
	b.MarkDirty()
}

func (b *BaseComponent) GetProps() interface{} {
	return b.prevProps
}

func (b *BaseComponent) GetState() interface{} {
	return b.prevState
}
