package react

import "context"

// React 的生命周期简介
//React 的生命周期是指组件从创建到销毁的整个过程，分为 挂载阶段、更新阶段 和 卸载阶段。在每个阶段，React 提供了特定的生命周期方法，允许开发者在适当的时机执行代码。
//挂载阶段
//挂载阶段是组件被创建并插入到 DOM 的过程，主要方法包括：
//constructor()：用于初始化组件的状态和绑定方法。
//getDerivedStateFromProps()：静态方法，用于根据 props 更新 state。
//render()：返回组件的 UI 结构，不应包含副作用。
//componentDidMount()：组件挂载完成后调用，常用于发起异步请求或操作 DOM。
//更新阶段
//更新阶段发生在组件的状态或属性变化时，主要方法包括：
//getDerivedStateFromProps()：与挂载阶段相同，用于更新 state。
//shouldComponentUpdate()：决定组件是否需要重新渲染，可用于性能优化。
//render()：重新渲染组件的 UI。
//getSnapshotBeforeUpdate()：在 DOM 更新前调用，用于捕获更新前的信息。
//componentDidUpdate()：组件更新后调用，适合处理 DOM 操作或依赖更新后的数据。
//卸载阶段
//卸载阶段是组件从 DOM 中移除的过程，主要方法包括：
//componentWillUnmount()：在组件卸载前调用，用于清理定时器、取消订阅等操作。
//错误处理
//React 提供了以下方法处理组件中的错误：
//getDerivedStateFromError()：捕获错误并更新 state。
//componentDidCatch()：记录错误信息，适合用于日志记录。
//函数组件中的生命周期
//在 React 16.8 之后，函数组件通过 Hooks 模拟生命周期：
//useEffect：结合了 componentDidMount、componentDidUpdate 和 componentWillUnmount 的功能。
//useState：管理组件状态。
//useRef：保持引用不变。

// Props 定义组件属性的类型
type Props map[string]interface{}

// State 定义组件状态的类型
type State map[string]interface{}

// Element 表示虚拟DOM元素
type Element interface {
	Type() string
	Props() Props
	Children() []Element
}

// Snapshot 表示在DOM更新前捕获的信息
type Snapshot interface{}

// PlatformRenderer 定义平台渲染器接口
type PlatformRenderer interface {
	CreateView(viewType string, props Props) interface{}              // 创建原生视图
	UpdateView(view interface{}, props Props)                         // 更新原生视图
	MountView(view interface{}, parent interface{})                    // 挂载视图到父容器
	UnmountView(view interface{})                                     // 卸载视图
	AddViewListener(view interface{}, event string, handler func(map[string]interface{})) // 添加视图事件监听器
	RemoveViewListener(view interface{}, event string)                 // 移除视图事件监听器
}

// NativeRenderer 全局原生渲染器实例
var NativeRenderer PlatformRenderer

// SetPlatformRenderer 设置全局平台渲染器
func SetPlatformRenderer(renderer PlatformRenderer) {
	NativeRenderer = renderer
}

// NativeElement 表示原生元素
type NativeElement struct {
	typeName string
	props    Props
	children []Element
}

// Type 返回原生元素类型
func (e *NativeElement) Type() string {
	return e.typeName
}

// Props 返回原生元素属性
func (e *NativeElement) Props() Props {
	return e.props
}

// Children 返回原生元素子节点
func (e *NativeElement) Children() []Element {
	return e.children
}

// CreateNativeElement 创建原生元素
func CreateNativeElement(typeName string, props Props, children []Element) *NativeElement {
	return &NativeElement{
		typeName: typeName,
		props:    props,
		children: children,
	}
}

// NativeComponentFactory 原生组件工厂函数类型
type NativeComponentFactory func() NativeComponent

// nativeComponentRegistry 原生组件注册表
var nativeComponentRegistry = make(map[string]NativeComponentFactory)

// RegisterNativeComponent 注册原生组件
func RegisterNativeComponent(name string, factory NativeComponentFactory) {
	nativeComponentRegistry[name] = factory
}

// CreateNativeComponent 创建原生组件实例
func CreateNativeComponent(name string, props Props) (NativeComponent, bool) {
	if factory, exists := nativeComponentRegistry[name]; exists {
		component := factory()
		component.Initialize(props)
		return component, true
	}
	return nil, false
}

// NativeRendererManager 原生渲染管理器，负责协调组件树到原生视图树的转换
type NativeRendererManager struct {
	componentMap map[Component]NativeComponent // 组件映射关系
}

// NewNativeRendererManager 创建新的原生渲染管理器
func NewNativeRendererManager() *NativeRendererManager {
	return &NativeRendererManager{
		componentMap: make(map[Component]NativeComponent),
	}
}

// RenderComponent 渲染组件到原生视图
func (m *NativeRendererManager) RenderComponent(component Component, parentView interface{}) NativeComponent {
	// 简单实现，实际需要更复杂的逻辑来处理组件树
	// 这里仅作为示例框架
	element := component.Render()
	if nativeElement, ok := element.(*NativeElement); ok {
		if nativeComp, exists := CreateNativeComponent(nativeElement.Type(), nativeElement.Props()); exists {
			// 创建原生视图
			nativeComp.CreateNativeView()
			// 挂载到父容器
			nativeComp.MountNativeView(parentView)
			// 完成挂载
			nativeComp.DidMount()
			// 保存映射关系
			m.componentMap[component] = nativeComp
			return nativeComp
		}
	}
	return nil
}

// Component 定义React组件的基础接口
type Component interface {
	// 挂载阶段生命周期方法
	Constructor(props Props)
	GetDerivedStateFromProps(props Props, state State) State
	Render() Element
	ComponentDidMount()

	// 更新阶段生命周期方法
	ShouldComponentUpdate(nextProps Props, nextState State) bool
	GetSnapshotBeforeUpdate(prevProps Props, prevState State) Snapshot
	ComponentDidUpdate(prevProps Props, prevState State, snapshot Snapshot)

	// 卸载阶段生命周期方法
	ComponentWillUnmount()

	// 错误处理生命周期方法
	GetDerivedStateFromError(error error) State
	ComponentDidCatch(error error, info ErrorInfo)

	// 状态和属性管理
	SetState(newState State, callback func())
	ForceUpdate(callback func())

	// 获取组件信息
	GetProps() Props
	GetState() State
	GetContext() context.Context
	GetRefs() map[string]interface{}
}

// NativeComponent 原生组件接口
type NativeComponent interface {
	// 基础方法
	Initialize(props Props)          // 初始化组件
	GetProps() Props                 // 获取组件属性
	GetState() State                 // 获取组件状态
	UpdateProps(props Props)         // 更新组件属性
	SetState(state State)            // 更新组件状态
	
	// 视图操作方法
	CreateNativeView() interface{}   // 创建原生视图
	MountNativeView(parent interface{}) // 挂载原生视图到父容器
	UnmountNativeView()              // 卸载原生视图
	GetNativeView() interface{}      // 获取原生视图引用
	SetParentView(parentView interface{}) // 设置父视图引用
	GetParentView() interface{}      // 获取父视图引用
	
	// 渲染方法
	RenderToNative()                 // 渲染到原生视图
	
	// 生命周期方法
	DidMount()                       // 挂载完成后调用
	DidUpdate()                      // 更新完成后调用
	WillUnmount()                    // 卸载前调用
	
	// 事件处理方法
	AddNativeListener(event string, handler func(map[string]interface{})) // 添加原生事件监听器
	RemoveNativeListener(event string) // 移除原生事件监听器
	UpdateNativeProps(props Props)    // 更新原生属性
}

// ErrorInfo 包含错误发生时的组件栈信息
type ErrorInfo struct {
	ComponentStack string
}

// ErrorBoundary 错误边界组件，用于捕获子组件树中的错误
type ErrorBoundary struct {
	BaseComponent
	child    Component
	fallback Component
}

// NewErrorBoundary 创建一个新的错误边界组件
func NewErrorBoundary(child Component, fallback Component) *ErrorBoundary {
	boundary := &ErrorBoundary{
		child:    child,
		fallback: fallback,
	}
	boundary.Constructor(nil)
	return boundary
}

// ComponentType 定义组件类型
type ComponentType interface{}

// FunctionalComponent 定义函数组件类型
type FunctionalComponent func(props Props) Element

// HookType 定义Hook类型
type HookType int

const (
	HookTypeState HookType = iota
	HookTypeEffect
	HookTypeRef
	// 其他Hook类型
)

// Hook 定义Hook基础结构
type Hook struct {
	Type          HookType
	MemoizedState interface{}
	Next          *Hook
}

// Effect 定义Effect结构
type Effect struct {
	Create       func() func()
	Destroy      func()
	Dependencies []interface{}
}

// RefObject 定义Ref对象结构
type RefObject struct {
	Current interface{}
}

// hooksContext 保存当前函数组件的hooks状态
var hooksContext = &struct {
	currentComponent *FunctionalComponentAdapter
	hookIndex        int
}{}

// FunctionalComponentAdapter 将函数组件适配为标准Component接口
type FunctionalComponentAdapter struct {
	BaseComponent
	function  FunctionalComponent
	firstHook *Hook
	effects   []*Effect
}

// NewFunctionalComponentAdapter 创建函数组件适配器
func NewFunctionalComponentAdapter(fn FunctionalComponent, props Props) *FunctionalComponentAdapter {
	adapter := &FunctionalComponentAdapter{
		function: fn,
		effects:  make([]*Effect, 0),
	}
	adapter.Constructor(props)
	return adapter
}

// Render 调用函数组件进行渲染
func (a *FunctionalComponentAdapter) Render() Element {
	// 设置hooks上下文
	prevComponent := hooksContext.currentComponent
	prevIndex := hooksContext.hookIndex

	hooksContext.currentComponent = a
	hooksContext.hookIndex = 0

	// 调用函数组件
	result := a.function(a.props)

	// 恢复上下文
	hooksContext.currentComponent = prevComponent
	hooksContext.hookIndex = prevIndex

	return result
}

// ComponentDidMount 函数组件挂载后处理effects
func (a *FunctionalComponentAdapter) ComponentDidMount() {
	a.BaseComponent.ComponentDidMount()
	a.processEffects()
}

// ComponentDidUpdate 函数组件更新后处理effects
func (a *FunctionalComponentAdapter) ComponentDidUpdate(prevProps Props, prevState State, snapshot Snapshot) {
	a.BaseComponent.ComponentDidUpdate(prevProps, prevState, snapshot)
	a.processEffects()
}

// ComponentWillUnmount 清理所有effects
func (a *FunctionalComponentAdapter) ComponentWillUnmount() {
	a.BaseComponent.ComponentWillUnmount()
	for _, effect := range a.effects {
		if effect.Destroy != nil {
			effect.Destroy()
		}
	}
	a.effects = a.effects[:0]
}

// processEffects 处理所有effects
func (a *FunctionalComponentAdapter) processEffects() {
	for _, effect := range a.effects {
		if effect.Destroy != nil {
			// 安全地调用destroy函数
			safeCall(effect.Destroy)
		}
		// 安全地调用create函数
		destroy := safeCallWithReturn(effect.Create)
		effect.Destroy = destroy
	}
}

// safeCall 安全地调用函数，捕获并处理错误
func safeCall(fn func()) {
	if fn == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			// 处理panic，可以记录日志或向上传播
			// 在实际应用中，这里可以将错误冒泡到ErrorBoundary
		}
	}()
	fn()
}

// safeCallWithReturn 安全地调用有返回值的函数，捕获并处理错误
func safeCallWithReturn(fn func() func()) func() {
	if fn == nil {
		return nil
	}
	defer func() {
		if r := recover(); r != nil {
			// 处理panic，可以记录日志或向上传播
		}
	}()
	return fn()
}

// useState 模拟React的useState Hook
func useState(initialState interface{}) (interface{}, func(interface{})) {
	component := hooksContext.currentComponent
	index := hooksContext.hookIndex

	// 查找或创建hook
	hook := component.getHookAt(index, HookTypeState)

	// 首次调用时初始化状态
	if hook.MemoizedState == nil {
		if initFn, ok := initialState.(func() interface{}); ok {
			hook.MemoizedState = initFn()
		} else {
			hook.MemoizedState = initialState
		}
	}

	// 创建setState函数
	setState := func(newState interface{}) {
		var value interface{}
		if updateFn, ok := newState.(func(interface{}) interface{}); ok {
			value = updateFn(hook.MemoizedState)
		} else {
			value = newState
		}

		hook.MemoizedState = value
		component.ForceUpdate(nil)
	}

	hooksContext.hookIndex++
	return hook.MemoizedState, setState
}

// useEffect 模拟React的useEffect Hook
func useEffect(create func() func(), deps []interface{}) {
	component := hooksContext.currentComponent
	index := hooksContext.hookIndex

	// 查找或创建hook
	hook := component.getHookAt(index, HookTypeEffect)

	// 检查依赖项是否变化
	shouldRun := hook.MemoizedState == nil
	if !shouldRun && deps != nil {
		oldDeps, ok := hook.MemoizedState.([]interface{})
		shouldRun = !ok || !depsEqual(oldDeps, deps)
	}

	if shouldRun || deps == nil {
		// 创建新的effect
		effect := &Effect{
			Create:       create,
			Dependencies: deps,
		}
		component.effects = append(component.effects, effect)
		hook.MemoizedState = deps
	}

	hooksContext.hookIndex++
}

// useRef 模拟React的useRef Hook
func useRef(initialValue interface{}) *RefObject {
	component := hooksContext.currentComponent
	index := hooksContext.hookIndex

	// 查找或创建hook
	hook := component.getHookAt(index, HookTypeRef)

	// 首次调用时初始化ref
	if hook.MemoizedState == nil {
		hook.MemoizedState = &RefObject{
			Current: initialValue,
		}
	}

	hooksContext.hookIndex++
	return hook.MemoizedState.(*RefObject)
}

// getHookAt 获取指定索引的hook，如果不存在则创建
func (a *FunctionalComponentAdapter) getHookAt(index int, hookType HookType) *Hook {
	if a.firstHook == nil {
		a.firstHook = &Hook{}
	}

	current := a.firstHook
	for i := 0; i < index; i++ {
		if current.Next == nil {
			current.Next = &Hook{}
		}
		current = current.Next
	}

	current.Type = hookType
	return current
}

// depsEqual 比较两个依赖数组是否相等
func depsEqual(oldDeps, newDeps []interface{}) bool {
	if len(oldDeps) != len(newDeps) {
		return false
	}

	for i := range oldDeps {
		// 简化实现，实际React使用Object.is比较
		if oldDeps[i] != newDeps[i] {
			return false
		}
	}

	return true
}

// UpdateQueue 定义更新队列结构
type UpdateQueue struct {
	pendingStates     []State
	pendingCallbacks  []func()
	isBatchingUpdates bool
}

// BaseComponent 提供Component接口的基础实现
type BaseComponent struct {
	props       Props
	state       State
	context     context.Context
	refs        map[string]interface{}
	mounted     bool
	updateQueue *UpdateQueue
	// 其他内部状态管理字段
}

// BaseNativeComponent 提供NativeComponent接口的基础实现
type BaseNativeComponent struct {
	props           Props               // 组件属性
	state           State               // 简单状态管理
	nativeView      interface{}         // 原生视图引用
	nativeProps     Props               // 原生属性缓存
	nativeListeners map[string]func(map[string]interface{}) // 原生事件监听器
	parentView      interface{}         // 父视图引用
	isMounted       bool                // 挂载状态
	children        []NativeComponent   // 子组件列表
}

// GetProps 获取组件属性
func (n *BaseNativeComponent) GetProps() Props {
	return n.props
}

// GetState 获取组件状态
func (n *BaseNativeComponent) GetState() State {
	return n.state
}

// SetParentView 设置父视图引用
func (n *BaseNativeComponent) SetParentView(parentView interface{}) {
	n.parentView = parentView
}

// GetParentView 获取父视图引用
func (n *BaseNativeComponent) GetParentView() interface{} {
	return n.parentView
}

// AddChild 添加子组件
func (n *BaseNativeComponent) AddChild(child NativeComponent) {
	if n.children == nil {
		n.children = make([]NativeComponent, 0)
	}
	child.SetParentView(n.nativeView)
	n.children = append(n.children, child)
}

// GetChildren 获取所有子组件
func (n *BaseNativeComponent) GetChildren() []NativeComponent {
	return n.children
}

// RemoveChild 移除子组件
func (n *BaseNativeComponent) RemoveChild(index int) {
	if index >= 0 && index < len(n.children) {
		// 先卸载子组件
		n.children[index].UnmountNativeView()
		n.children[index].WillUnmount()
		// 移除子组件
		n.children = append(n.children[:index], n.children[index+1:]...)
	}
}

// Initialize 初始化原生组件
func (n *BaseNativeComponent) Initialize(props Props) {
	if props == nil {
		n.props = make(Props)
	} else {
		n.props = props
	}
	n.state = make(State)
	n.nativeListeners = make(map[string]func(map[string]interface{}))
	n.isMounted = false
	// 初始化原生属性缓存
	n.nativeProps = make(Props)
}

// RenderToNative 渲染到原生视图
func (n *BaseNativeComponent) RenderToNative() {
	// 如果还没有创建原生视图，则创建
	if n.nativeView == nil {
		n.CreateNativeView()
	}
	// 更新原生属性
	n.UpdateNativeProps(n.props)
	
	// 递归渲染所有子组件
	for _, child := range n.children {
		child.RenderToNative()
	}
}

// CreateNativeView 创建原生视图 - 基础实现
func (n *BaseNativeComponent) CreateNativeView() interface{} {
	// 基础实现返回nil，子类应该覆盖此方法
	return nil
}

// DidMount 原生视图挂载完成后调用
func (n *BaseNativeComponent) DidMount() {
	n.isMounted = true
	// 子类可以覆盖此方法以实现挂载后的逻辑
}

// DidUpdate 原生视图更新完成后调用
func (n *BaseNativeComponent) DidUpdate() {
	// 子类可以覆盖此方法以实现更新后的逻辑
}

// WillUnmount 原生视图卸载前调用
func (n *BaseNativeComponent) WillUnmount() {
	// 清理所有事件监听器
	for event := range n.nativeListeners {
		n.RemoveNativeListener(event)
	}
	n.nativeListeners = nil
	n.isMounted = false
}

// GetNativeView 获取原生视图引用
func (n *BaseNativeComponent) GetNativeView() interface{} {
	return n.nativeView
}

// UnmountNativeView 卸载原生视图
func (n *BaseNativeComponent) UnmountNativeView() {
	// 先卸载所有子组件
	for _, child := range n.children {
		child.UnmountNativeView()
		child.WillUnmount()
	}
	
	// 清空子组件列表
	n.children = nil
	
	// 执行自身卸载逻辑
	if NativeRenderer != nil && n.nativeView != nil {
		NativeRenderer.UnmountView(n.nativeView)
	}
}

// UpdateNativeProps 更新原生属性
func (n *BaseNativeComponent) UpdateNativeProps(props Props) {
	// 计算属性差异
	diff := make(Props)
	for key, newValue := range props {
		if oldValue, exists := n.nativeProps[key]; !exists || oldValue != newValue {
			diff[key] = newValue
		}
	}
	
	// 更新属性缓存
	for key, value := range props {
		n.nativeProps[key] = value
	}
	
	// 如果有差异且存在原生渲染器，则通过渲染器更新视图
	if len(diff) > 0 && NativeRenderer != nil && n.nativeView != nil {
		NativeRenderer.UpdateView(n.nativeView, diff)
	}
}

// UpdateProps 更新组件属性
func (n *BaseNativeComponent) UpdateProps(props Props) {
	for key, value := range props {
		n.props[key] = value
	}
	// 属性更新后重新渲染
	if n.isMounted {
		n.RenderToNative()
		n.DidUpdate()
	}
}

// AddNativeListener 添加原生事件监听器
func (n *BaseNativeComponent) AddNativeListener(event string, handler func(map[string]interface{})) {
	// 保存事件处理器
	n.nativeListeners[event] = handler
	
	// 如果有原生渲染器和视图，使用渲染器添加事件监听
	if NativeRenderer != nil && n.nativeView != nil {
		// 包装处理器以确保在组件挂载状态下执行
		wrappedHandler := func(eventData map[string]interface{}) {
			if n.isMounted {
				handler(eventData)
			}
		}
		NativeRenderer.AddViewListener(n.nativeView, event, wrappedHandler)
	}
}

// RemoveNativeListener 移除原生事件监听器
func (n *BaseNativeComponent) RemoveNativeListener(event string) {
	// 从映射中删除处理器
	delete(n.nativeListeners, event)
	
	// 如果有原生渲染器和视图，使用渲染器移除事件监听
	if NativeRenderer != nil && n.nativeView != nil {
		NativeRenderer.RemoveViewListener(n.nativeView, event)
	}
}

// EmitEvent 触发原生事件（用于平台特定实现调用）
func (n *BaseNativeComponent) EmitEvent(event string, eventData map[string]interface{}) {
	// 安全地触发事件处理器
	if handler, exists := n.nativeListeners[event]; exists && n.isMounted {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					// 处理panic，避免整个应用崩溃
				}
			}()
			handler(eventData)
		}()
	}
}

// NativeText 文本原生组件
type NativeText struct {
	BaseNativeComponent
}

// NewNativeText 创建新的文本原生组件
func NewNativeText() *NativeText {
	return &NativeText{}
}

// CreateNativeView 创建文本原生视图
func (t *NativeText) CreateNativeView() interface{} {
	if NativeRenderer != nil {
		t.nativeView = NativeRenderer.CreateView("text", t.props)
	}
	return t.nativeView
}

// MountNativeView 挂载文本视图到父容器
func (t *NativeText) MountNativeView(parent interface{}) {
	t.parentView = parent
	if NativeRenderer != nil && t.nativeView != nil {
		NativeRenderer.MountView(t.nativeView, parent)
	}
}

// UnmountNativeView 卸载文本视图
func (t *NativeText) UnmountNativeView() {
	if NativeRenderer != nil && t.nativeView != nil {
		NativeRenderer.UnmountView(t.nativeView)
	}
	t.parentView = nil
	t.nativeView = nil
}

// NativeView 容器原生组件
type NativeView struct {
	BaseNativeComponent
	children []NativeComponent // 子原生组件
}

// NewNativeView 创建新的容器原生组件
func NewNativeView() *NativeView {
	return &NativeView{
		children: make([]NativeComponent, 0),
	}
}

// CreateNativeView 创建容器原生视图
func (v *NativeView) CreateNativeView() interface{} {
	if NativeRenderer != nil {
		v.nativeView = NativeRenderer.CreateView("view", v.props)
	}
	return v.nativeView
}

// MountNativeView 挂载容器视图到父容器
func (v *NativeView) MountNativeView(parent interface{}) {
	v.parentView = parent
	if NativeRenderer != nil && v.nativeView != nil {
		NativeRenderer.MountView(v.nativeView, parent)
	}
}

// UnmountNativeView 卸载容器视图
func (v *NativeView) UnmountNativeView() {
	// 先卸载所有子组件
	for _, child := range v.children {
		child.UnmountNativeView()
	}
	v.children = v.children[:0]
	
	if NativeRenderer != nil && v.nativeView != nil {
		NativeRenderer.UnmountView(v.nativeView)
	}
	v.parentView = nil
	v.nativeView = nil
}

// AddChild 添加子原生组件
func (v *NativeView) AddChild(child NativeComponent) {
	v.children = append(v.children, child)
	// 将子组件挂载到当前视图
	if child != nil && v.nativeView != nil {
		child.CreateNativeView()
		child.MountNativeView(v.nativeView)
		child.DidMount()
	}
}

// RemoveChild 移除子原生组件
func (v *NativeView) RemoveChild(child NativeComponent) {
	for i, c := range v.children {
		if c == child {
			// 从数组中移除
			v.children = append(v.children[:i], v.children[i+1:]...)
			// 卸载子组件
			child.UnmountNativeView()
			break
		}
	}
}

// InitializeNativeComponents 初始化并注册内置原生组件
func InitializeNativeComponents() {
	// 注册Text组件
	RegisterNativeComponent("Text", func() NativeComponent {
		return NewNativeText()
	})
	
	// 注册View组件
	RegisterNativeComponent("View", func() NativeComponent {
		return NewNativeView()
	})
	
	// 可以在这里注册更多内置原生组件
}

// Constructor 初始化组件的状态和属性
func (c *BaseComponent) Constructor(props Props) {
	if props == nil {
		c.props = make(Props)
	} else {
		c.props = props
	}
	c.state = make(State)
	c.refs = make(map[string]interface{})
	c.context = context.Background()
	c.mounted = false
	c.updateQueue = &UpdateQueue{
		pendingStates:     make([]State, 0),
		pendingCallbacks:  make([]func(), 0),
		isBatchingUpdates: false,
	}
}

// GetDerivedStateFromProps 根据props更新state
func (c *BaseComponent) GetDerivedStateFromProps(props Props, state State) State {
	// 基础实现返回原state，子类可以覆盖此方法
	return state
}

// Render 返回组件的UI结构，基础实现返回nil，子类必须实现
func (c *BaseComponent) Render() Element {
	return nil
}

// ComponentDidMount 组件挂载完成后调用
func (c *BaseComponent) ComponentDidMount() {
	c.mounted = true
	// 基础实现为空，子类可以覆盖此方法
}

// ShouldComponentUpdate 决定组件是否需要重新渲染
func (c *BaseComponent) ShouldComponentUpdate(nextProps Props, nextState State) bool {
	// 基础实现总是返回true，表示需要重新渲染
	return true
}

// GetSnapshotBeforeUpdate 在DOM更新前调用，捕获更新前的信息
func (c *BaseComponent) GetSnapshotBeforeUpdate(prevProps Props, prevState State) Snapshot {
	// 基础实现返回nil，子类可以覆盖此方法
	return nil
}

// ComponentDidUpdate 组件更新后调用
func (c *BaseComponent) ComponentDidUpdate(prevProps Props, prevState State, snapshot Snapshot) {
	// 基础实现为空，子类可以覆盖此方法
}

// ComponentWillUnmount 在组件卸载前调用
func (c *BaseComponent) ComponentWillUnmount() {
	c.mounted = false
	// 基础实现为空，子类可以覆盖此方法进行清理工作
}

// GetDerivedStateFromError 捕获错误并更新state
func (c *BaseComponent) GetDerivedStateFromError(err error) State {
	// 基础实现返回原state，子类可以覆盖此方法
	return c.state
}

// ComponentDidCatch 记录错误信息
func (c *BaseComponent) ComponentDidCatch(err error, info ErrorInfo) {
	// 基础实现为空，子类可以覆盖此方法进行错误处理
}

// Render 错误边界组件的渲染方法
func (b *ErrorBoundary) Render() Element {
	// 检查是否有错误状态
	hasError, exists := b.state["hasError"]
	if exists && hasError.(bool) {
		return b.fallback.Render()
	}
	return b.child.Render()
}

// GetDerivedStateFromError 错误边界特有的错误状态处理
func (b *ErrorBoundary) GetDerivedStateFromError(err error) State {
	return State{"hasError": true}
}

// ComponentDidCatch 错误边界特有的错误处理
func (b *ErrorBoundary) ComponentDidCatch(err error, info ErrorInfo) {
	// 记录错误信息，可以在这里添加日志记录
	// 在实际应用中，可以发送错误报告到监控服务
	// 简化实现，仅设置错误信息到状态
	b.state["error"] = err
	b.state["errorInfo"] = info
}

// SetState 更新组件状态
func (c *BaseComponent) SetState(newState State, callback func()) {
	if !c.mounted {
		// 组件未挂载时，直接更新state
		c.updateState(newState)
		if callback != nil {
			callback()
		}
		return
	}

	// 将更新添加到队列
	c.updateQueue.pendingStates = append(c.updateQueue.pendingStates, newState)
	if callback != nil {
		c.updateQueue.pendingCallbacks = append(c.updateQueue.pendingCallbacks, callback)
	}

	// 如果不在批量更新模式，则处理更新
	if !c.updateQueue.isBatchingUpdates {
		c.performUpdateIfNeeded()
	}
}

// ForceUpdate 强制组件重新渲染
func (c *BaseComponent) ForceUpdate(callback func()) {
	if c.mounted {
		// 强制重新渲染逻辑
		// 这里简化实现，实际应该触发重新渲染
		if callback != nil {
			callback()
		}
	}
}

// GetProps 获取组件属性
func (c *BaseComponent) GetProps() Props {
	return c.props
}

// GetState 获取组件状态
func (c *BaseComponent) GetState() State {
	return c.state
}

// GetContext 获取组件上下文
func (c *BaseComponent) GetContext() context.Context {
	return c.context
}

// GetRefs 获取组件引用
func (c *BaseComponent) GetRefs() map[string]interface{} {
	return c.refs
}

// updateState 内部方法，更新组件状态
func (c *BaseComponent) updateState(newState State) {
	if newState == nil {
		return
	}

	// 合并新状态到当前状态
	for key, value := range newState {
		c.state[key] = value
	}
}

// performUpdateIfNeeded 执行必要的更新
func (c *BaseComponent) performUpdateIfNeeded() {
	if len(c.updateQueue.pendingStates) == 0 {
		return
	}

	// 进入批量更新模式
	prevBatchingState := c.updateQueue.isBatchingUpdates
	c.updateQueue.isBatchingUpdates = true

	// 处理所有待处理的状态更新
	for _, newState := range c.updateQueue.pendingStates {
		c.updateState(newState)
	}

	// 清空队列
	callbacks := c.updateQueue.pendingCallbacks
	c.updateQueue.pendingStates = c.updateQueue.pendingStates[:0]
	c.updateQueue.pendingCallbacks = c.updateQueue.pendingCallbacks[:0]

	// 恢复批量更新状态
	c.updateQueue.isBatchingUpdates = prevBatchingState

	// 执行所有回调
	for _, callback := range callbacks {
		if callback != nil {
			callback()
		}
	}

	// 触发重新渲染逻辑
	// 这里简化实现，实际应该检查ShouldComponentUpdate等
}
