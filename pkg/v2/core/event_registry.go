package core

import (
	"reflect"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

// EventCallback 是事件监听回调函数签名。
// 返回 true 表示事件已被消费，不再继续传递。
type EventCallback func(e *Event) bool

// eventListenerKey 用于在注册表中唯一标识一个监听器。
type eventListenerKey struct {
	target    Element
	eventType EventType
}

// EventRegistry 管理组件到引擎的事件监听注册。
type EventRegistry struct {
	mu        sync.RWMutex
	listeners map[eventListenerKey][]EventCallback
	// capture 用于在捕获阶段触发的监听器（从根到目标）
	captureListeners map[eventListenerKey][]EventCallback
}

// NewEventRegistry 创建事件注册表。
func NewEventRegistry() *EventRegistry {
	return &EventRegistry{
		listeners:        make(map[eventListenerKey][]EventCallback),
		captureListeners: make(map[eventListenerKey][]EventCallback),
	}
}

// AddEventListener 为指定组件注册事件监听器。
// target: 目标组件
// eventType: 监听的事件类型
// callback: 回调函数，返回 true 消费事件
func (r *EventRegistry) AddEventListener(target Element, eventType EventType, callback EventCallback) {
	if target == nil || callback == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	key := eventListenerKey{target: target, eventType: eventType}
	r.listeners[key] = append(r.listeners[key], callback)
}

// RemoveEventListener 移除指定组件的某个事件监听器。
// 使用回调函数指针进行精确匹配移除。
func (r *EventRegistry) RemoveEventListener(target Element, eventType EventType, callback EventCallback) {
	if target == nil || callback == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	key := eventListenerKey{target: target, eventType: eventType}
	callbacks := r.listeners[key]
	for i, cb := range callbacks {
		if compareFunc(cb, callback) {
			r.listeners[key] = append(callbacks[:i], callbacks[i+1:]...)
			return
		}
	}
}

// RemoveAllListeners 移除指定组件的所有事件监听器。
func (r *EventRegistry) RemoveAllListeners(target Element) {
	if target == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for key := range r.listeners {
		if key.target == target {
			delete(r.listeners, key)
		}
	}
	for key := range r.captureListeners {
		if key.target == target {
			delete(r.captureListeners, key)
		}
	}
}

// HasListener 检查指定组件是否有指定类型的监听器。
func (r *EventRegistry) HasListener(target Element, eventType EventType) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	key := eventListenerKey{target: target, eventType: eventType}
	return len(r.listeners[key]) > 0
}

// Dispatch 分发事件到目标组件及其所有祖先的监听器。
// 先执行捕获阶段（从根到目标），再执行冒泡阶段（从目标到根）。
// 如果任一阶段返回 true，则事件被消费，停止传递。
func (r *EventRegistry) Dispatch(event *Event) bool {
	if event == nil || event.Target == nil {
		return false
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	// 构建从根到目标的路径
	path := buildEventPath(event.Target)

	// 捕获阶段：从根到目标
	for _, el := range path {
		key := eventListenerKey{target: el, eventType: event.Type}
		for _, cb := range r.captureListeners[key] {
			if cb(event) {
				return true
			}
		}
	}

	// 冒泡阶段：从目标到根
	for i := len(path) - 1; i >= 0; i-- {
		el := path[i]
		key := eventListenerKey{target: el, eventType: event.Type}
		for _, cb := range r.listeners[key] {
			if cb(event) {
				return true
			}
		}
	}

	return false
}

// DispatchToTarget 只分发事件到指定目标（不冒泡）。
func (r *EventRegistry) DispatchToTarget(target Element, event *Event) bool {
	if target == nil || event == nil {
		return false
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	key := eventListenerKey{target: target, eventType: event.Type}
	for _, cb := range r.listeners[key] {
		if cb(event) {
			return true
		}
	}
	return false
}

// buildEventPath 构建从根节点到目标节点的路径。
func buildEventPath(target Element) []Element {
	var path []Element
	for el := target; el != nil; el = el.GetParent() {
		path = append(path, el)
	}
	// 反转，使根节点在前
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

// compareFunc 比较两个函数指针是否相同。
// Go 中函数不能直接比较，使用 reflect 比较函数值。
func compareFunc(a, b EventCallback) bool {
	return reflect.ValueOf(a).Pointer() == reflect.ValueOf(b).Pointer()
}

// ==================== Engine 便捷方法 ====================

// AddEventListener 为组件注册事件监听器。
func (e *Engine) AddEventListener(target Element, eventType EventType, callback EventCallback) {
	if e.eventRegistry == nil {
		e.eventRegistry = NewEventRegistry()
	}
	e.eventRegistry.AddEventListener(target, eventType, callback)
}

// RemoveEventListener 移除组件的某个事件监听器。
func (e *Engine) RemoveEventListener(target Element, eventType EventType, callback EventCallback) {
	if e.eventRegistry == nil {
		return
	}
	e.eventRegistry.RemoveEventListener(target, eventType, callback)
}

// RemoveAllEventListeners 移除组件的所有事件监听器。
func (e *Engine) RemoveAllEventListeners(target Element) {
	if e.eventRegistry == nil {
		return
	}
	e.eventRegistry.RemoveAllListeners(target)
}

// dispatchRegistryEvent 通过注册表分发事件。
// 如果注册表消费了事件，返回 true；否则继续调用 HandleEvent。
func (e *Engine) dispatchRegistryEvent(event *Event) bool {
	if e.eventRegistry == nil {
		return false
	}
	return e.eventRegistry.Dispatch(event)
}

// ==================== Element 便捷方法 ====================

// OnClick 注册点击事件监听器。
func (b *BaseElement) OnClick(callback EventCallback) Element {
	if b.engine != nil {
		b.engine.AddEventListener(b.self, EventClick, callback)
	} else {
		// 延迟注册：等 OnMount 后再注册
		b.delayedListeners = append(b.delayedListeners, delayedListener{
			eventType: EventClick,
			callback:  callback,
		})
	}
	return b.self
}

// OnMouseDown 注册鼠标按下事件监听器。
func (b *BaseElement) OnMouseDown(callback EventCallback) Element {
	if b.engine != nil {
		b.engine.AddEventListener(b.self, EventMouseDown, callback)
	} else {
		b.delayedListeners = append(b.delayedListeners, delayedListener{
			eventType: EventMouseDown,
			callback:  callback,
		})
	}
	return b.self
}

// OnMouseUp 注册鼠标松开事件监听器。
func (b *BaseElement) OnMouseUp(callback EventCallback) Element {
	if b.engine != nil {
		b.engine.AddEventListener(b.self, EventMouseUp, callback)
	} else {
		b.delayedListeners = append(b.delayedListeners, delayedListener{
			eventType: EventMouseUp,
			callback:  callback,
		})
	}
	return b.self
}

// OnMouseMove 注册鼠标移动事件监听器。
func (b *BaseElement) OnMouseMove(callback EventCallback) Element {
	if b.engine != nil {
		b.engine.AddEventListener(b.self, EventMouseMove, callback)
	} else {
		b.delayedListeners = append(b.delayedListeners, delayedListener{
			eventType: EventMouseMove,
			callback:  callback,
		})
	}
	return b.self
}

// OnScroll 注册滚动事件监听器。
func (b *BaseElement) OnScroll(callback EventCallback) Element {
	if b.engine != nil {
		b.engine.AddEventListener(b.self, EventScroll, callback)
	} else {
		b.delayedListeners = append(b.delayedListeners, delayedListener{
			eventType: EventScroll,
			callback:  callback,
		})
	}
	return b.self
}

// OnMouseEnter 注册鼠标进入事件监听器。
func (b *BaseElement) OnMouseEnter(callback EventCallback) Element {
	if b.engine != nil {
		b.engine.AddEventListener(b.self, EventMouseEnter, callback)
	} else {
		b.delayedListeners = append(b.delayedListeners, delayedListener{
			eventType: EventMouseEnter,
			callback:  callback,
		})
	}
	return b.self
}

// OnMouseLeave 注册鼠标离开事件监听器。
func (b *BaseElement) OnMouseLeave(callback EventCallback) Element {
	if b.engine != nil {
		b.engine.AddEventListener(b.self, EventMouseLeave, callback)
	} else {
		b.delayedListeners = append(b.delayedListeners, delayedListener{
			eventType: EventMouseLeave,
			callback:  callback,
		})
	}
	return b.self
}

// OnFocusIn 注册获得焦点事件监听器。
func (b *BaseElement) OnFocusIn(callback EventCallback) Element {
	if b.engine != nil {
		b.engine.AddEventListener(b.self, EventFocusIn, callback)
	} else {
		b.delayedListeners = append(b.delayedListeners, delayedListener{
			eventType: EventFocusIn,
			callback:  callback,
		})
	}
	return b.self
}

// OnFocusOut 注册失去焦点事件监听器。
func (b *BaseElement) OnFocusOut(callback EventCallback) Element {
	if b.engine != nil {
		b.engine.AddEventListener(b.self, EventFocusOut, callback)
	} else {
		b.delayedListeners = append(b.delayedListeners, delayedListener{
			eventType: EventFocusOut,
			callback:  callback,
		})
	}
	return b.self
}

// OnKeyDown 注册键盘按下事件监听器。
func (b *BaseElement) OnKeyDown(callback EventCallback) Element {
	if b.engine != nil {
		b.engine.AddEventListener(b.self, EventKeyDown, callback)
	} else {
		b.delayedListeners = append(b.delayedListeners, delayedListener{
			eventType: EventKeyDown,
			callback:  callback,
		})
	}
	return b.self
}

// OnKeyUp 注册键盘松开事件监听器。
func (b *BaseElement) OnKeyUp(callback EventCallback) Element {
	if b.engine != nil {
		b.engine.AddEventListener(b.self, EventKeyUp, callback)
	} else {
		b.delayedListeners = append(b.delayedListeners, delayedListener{
			eventType: EventKeyUp,
			callback:  callback,
		})
	}
	return b.self
}

// RemoveOnClick 移除点击事件监听器。
func (b *BaseElement) RemoveOnClick(callback EventCallback) Element {
	if b.engine != nil {
		b.engine.RemoveEventListener(b.self, EventClick, callback)
	}
	return b.self
}

// delayedListener 用于在 engine 尚未注入时延迟注册的事件监听器。
type delayedListener struct {
	eventType EventType
	callback  EventCallback
}

// flushDelayedListeners 在 OnMount 时刷新延迟注册的监听器。
func (b *BaseElement) flushDelayedListeners() {
	if b.engine == nil || len(b.delayedListeners) == 0 {
		return
	}
	for _, dl := range b.delayedListeners {
		b.engine.AddEventListener(b.self, dl.eventType, dl.callback)
	}
	b.delayedListeners = b.delayedListeners[:0]
}

// ==================== 兼容旧 HandleEvent 的桥接 ====================

// legacyEventHandler 用于兼容仍使用 HandleEvent 方法的组件。
// 它会将事件包装为回调注册到注册表中。
type legacyEventHandler struct {
	el Element
}

func (h *legacyEventHandler) handle(event *Event) bool {
	return h.el.HandleEvent(event)
}

// RegisterLegacyHandler 为仍使用 HandleEvent 的组件注册桥接。
// 在 Engine 的 onElementMounted 中自动调用。
func (e *Engine) RegisterLegacyHandler(el Element) {
	if el == nil {
		return
	}
	// 检查组件是否重写了 HandleEvent（不是 BaseElement 的默认实现）
	if _, ok := el.(interface{ HandleEvent(*Event) bool }); ok {
		// 为所有事件类型注册桥接
		for et := EventMouseMove; et <= EventResize; et++ {
			e.AddEventListener(el, et, (&legacyEventHandler{el: el}).handle)
		}
	}
}

// IsKeyPressed 是 ebiten.IsKeyPressed 的包装，便于组件使用。
func IsKeyPressed(key ebiten.Key) bool {
	return ebiten.IsKeyPressed(key)
}
