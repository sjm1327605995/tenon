package hooks

import (
	"sync"
)

// Context 上下文接口
type Context interface {
	ID() string
	Value() interface{}
	SetValue(interface{})
	Subscribe(func(interface{}))
	Unsubscribe(func(interface{}))
}

// BaseContext 基础上下文实现
type BaseContext struct {
	id         string
	value      interface{}
	listeners  []func(interface{})
	mu         sync.RWMutex
}

// NewContext 创建新的上下文
func NewContext(id string, initialValue interface{}) Context {
	return &BaseContext{
		id:        id,
		value:     initialValue,
		listeners: make([]func(interface{}), 0),
	}
}

// ID 返回上下文 ID
func (bc *BaseContext) ID() string {
	return bc.id
}

// Value 获取当前值
func (bc *BaseContext) Value() interface{} {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.value
}

// SetValue 设置新值并通知监听器
func (bc *BaseContext) SetValue(newValue interface{}) {
	bc.mu.Lock()
	bc.value = newValue
	bc.mu.Unlock()
	
	// 通知所有监听器
	for _, listener := range bc.listeners {
		listener(newValue)
	}
}

// Subscribe 订阅值变化
func (bc *BaseContext) Subscribe(listener func(interface{})) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	bc.listeners = append(bc.listeners, listener)
}

// Unsubscribe 取消订阅
func (bc *BaseContext) Unsubscribe(listener func(interface{})) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	
	for i, l := range bc.listeners {
		if &l == &listener {
			bc.listeners = append(bc.listeners[:i], bc.listeners[i+1:]...)
			break
		}
	}
}

// ContextHook 实现 useContext Hook
type ContextHook struct {
	context   Context
	component Component
	listener  func(interface{})
}

// NewContextHook 创建新的上下文 Hook
func NewContextHook(context Context, component Component) *ContextHook {
	hook := &ContextHook{
		context:   context,
		component: component,
	}
	
	// 创建监听器
	hook.listener = func(newValue interface{}) {
		if hook.component != nil {
			hook.component.MarkDirty()
		}
	}
	
	// 订阅上下文变化
	context.Subscribe(hook.listener)
	
	return hook
}

// Type 返回 Hook 类型
func (ch *ContextHook) Type() HookType {
	return HookTypeContext
}

// Update 更新上下文 Hook
func (ch *ContextHook) Update(value interface{}) {
	// 上下文 Hook 的更新逻辑在 useContext 中处理
}

// Cleanup 清理 Hook
func (ch *ContextHook) Cleanup() {
	if ch.context != nil && ch.listener != nil {
		ch.context.Unsubscribe(ch.listener)
	}
}

// GetValue 获取上下文值
func (ch *ContextHook) GetValue() interface{} {
	if ch.context != nil {
		return ch.context.Value()
	}
	return nil
}

// useContext 类似 React 的 useContext Hook
func useContext(context Context) interface{} {
	hooks := getCurrentHooks()

	var hook *ContextHook
	if hooks.index < len(hooks.hooks) {
		// 重用现有的 Hook
		existingHook := nextHook()
		var ok bool
		hook, ok = existingHook.(*ContextHook)
		if !ok {
			panic("Hook type mismatch: expected ContextHook")
		}
		
		// 检查上下文是否变化
		if hook.context != context {
			hook.Cleanup()
			hook.context = context
			hook.context.Subscribe(hook.listener)
		}
	} else {
		// 创建新的 Hook
		hook = NewContextHook(context, currentComponent)
		addHook(hook)
	}

	return hook.GetValue()
}

// UseContext 导出的 useContext Hook
func UseContext(context Context) interface{} {
	return useContext(context)
}

// UseContextTyped 类型安全的 useContext
func UseContextTyped[T any](context Context) T {
	value := UseContext(context)
	if typedValue, ok := value.(T); ok {
		return typedValue
	}
	
	// 如果类型不匹配，返回零值
	var zero T
	return zero
}

// ContextManager 上下文管理器
type ContextManager struct {
	contexts map[string]Context
	mu       sync.RWMutex
}

// NewContextManager 创建新的上下文管理器
func NewContextManager() *ContextManager {
	return &ContextManager{
		contexts: make(map[string]Context),
	}
}

// CreateContext 创建新的上下文
func (cm *ContextManager) CreateContext(id string, initialValue interface{}) Context {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	context := NewContext(id, initialValue)
	cm.contexts[id] = context
	return context
}

// GetContext 获取上下文
func (cm *ContextManager) GetContext(id string) (Context, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	
	context, exists := cm.contexts[id]
	return context, exists
}

// RemoveContext 移除上下文
func (cm *ContextManager) RemoveContext(id string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	delete(cm.contexts, id)
}

// ContextProvider 上下文提供者组件接口
type ContextProvider interface {
	Component
	GetContext() Context
}

// CreateContextProvider 创建上下文提供者
func CreateContextProvider(context Context, child Component) Component {
	// 这里需要集成到组件系统中
	// 暂时返回子组件
	return child
}

// ContextHookInfo 上下文 Hook 的调试信息
type ContextHookInfo struct {
	ComponentID string
	ContextID   string
	Value       interface{}
}

// GetContextHooksInfo 获取所有上下文 Hooks 的调试信息
func GetContextHooksInfo() []ContextHookInfo {
	var infos []ContextHookInfo

	for _, hooks := range components {
		for _, hook := range hooks.hooks {
			if contextHook, ok := hook.(*ContextHook); ok && contextHook.context != nil {
				infos = append(infos, ContextHookInfo{
					ComponentID: hooks.component.ID(),
					ContextID:   contextHook.context.ID(),
					Value:       contextHook.GetValue(),
				})
			}
		}
	}

	return infos
}

// 全局上下文管理器实例
var globalContextManager = NewContextManager()

// CreateGlobalContext 创建全局上下文
func CreateGlobalContext(id string, initialValue interface{}) Context {
	return globalContextManager.CreateContext(id, initialValue)
}

// GetGlobalContext 获取全局上下文
func GetGlobalContext(id string) (Context, bool) {
	return globalContextManager.GetContext(id)
}

// UseGlobalContext 使用全局上下文的便捷函数
func UseGlobalContext(id string) interface{} {
	if context, exists := GetGlobalContext(id); exists {
		return UseContext(context)
	}
	return nil
}

// UseGlobalContextTyped 类型安全的全局上下文
func UseGlobalContextTyped[T any](id string) T {
	value := UseGlobalContext(id)
	if typedValue, ok := value.(T); ok {
		return typedValue
	}
	
	var zero T
	return zero
}

// ContextConfig 上下文配置选项
type ContextConfig struct {
	// 是否持久化到本地存储
	Persistent bool
	// 持久化键名
	StorageKey string
}

// CreateContextWithConfig 带配置创建上下文
func CreateContextWithConfig(id string, initialValue interface{}, config ContextConfig) Context {
	context := NewContext(id, initialValue)
	
	// 应用配置
	if config.Persistent {
		// 这里需要集成到持久化系统中
		// 暂时返回基础上下文
	}
	
	return context
}