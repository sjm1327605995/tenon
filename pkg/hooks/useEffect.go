package hooks

import (
	"reflect"
	"sync"
	"time"
)

// EffectHook 实现 useEffect Hook
type EffectHook struct {
	effectFunc  func() func()
	cleanupFunc func()
	deps        []interface{}
	component   Component
	firstRun    bool
}

// NewEffectHook 创建新的副作用 Hook
func NewEffectHook(effectFunc func() func(), deps []interface{}, component Component) *EffectHook {
	return &EffectHook{
		effectFunc: effectFunc,
		deps:       deps,
		component:  component,
		firstRun:   true,
	}
}

// Type 返回 Hook 类型
func (eh *EffectHook) Type() HookType {
	return HookTypeEffect
}

// Update 更新副作用
func (eh *EffectHook) Update(value interface{}) {
	// Effect Hook 的更新逻辑在 shouldRun 中处理
}

// Cleanup 清理副作用
func (eh *EffectHook) Cleanup() {
	if eh.cleanupFunc != nil {
		eh.cleanupFunc()
		eh.cleanupFunc = nil
	}
}

// shouldRun 检查是否需要运行副作用
func (eh *EffectHook) shouldRun(newDeps []interface{}) bool {
	// 首次运行总是执行
	if eh.firstRun {
		eh.firstRun = false
		return true
	}

	// 如果没有依赖项，每次渲染都执行
	if len(eh.deps) == 0 && len(newDeps) == 0 {
		return true
	}

	// 依赖项数量变化，需要执行
	if len(eh.deps) != len(newDeps) {
		return true
	}

	// 检查依赖项是否变化
	for i := range eh.deps {
		if !reflect.DeepEqual(eh.deps[i], newDeps[i]) {
			return true
		}
	}

	return false
}

// run 运行副作用
func (eh *EffectHook) run() {
	eh.Cleanup() // 先清理之前的副作用

	if eh.effectFunc != nil {
		eh.cleanupFunc = eh.effectFunc()
	}
}

// useEffect 类似 React 的 useEffect Hook
func useEffect(effectFunc func() func(), deps []interface{}) {
	hooks := getCurrentHooks()

	var hook *EffectHook
	if hooks.index < len(hooks.hooks) {
		// 重用现有的 Hook
		existingHook := nextHook()
		var ok bool
		hook, ok = existingHook.(*EffectHook)
		if !ok {
			panic("Hook type mismatch: expected EffectHook")
		}

		// 检查是否需要运行副作用
		if hook.shouldRun(deps) {
			hook.deps = deps
			hook.run()
		}
	} else {
		// 创建新的 Hook 并立即运行
		hook = NewEffectHook(effectFunc, deps, currentComponent)
		addHook(hook)
		hook.run()
	}
}

// UseEffect 导出的 useEffect Hook
func UseEffect(effectFunc func() func(), deps ...interface{}) {
	useEffect(effectFunc, deps)
}

// UseLayoutEffect 类似 useEffect，但在布局更新后同步执行
type UseLayoutEffect func() func()

// EffectConfig 副作用配置选项
type EffectConfig struct {
	// 延迟执行，避免频繁触发
	Debounce time.Duration
	// 最大执行频率限制
	Throttle time.Duration
	// 是否在组件卸载时自动清理
	AutoCleanup bool
}

// useEffectWithConfig 带配置的 useEffect
func useEffectWithConfig(effectFunc func() func(), deps []interface{}, config EffectConfig) {
	useEffect(func() func() {
		// 应用配置的包装器
		var cleanup func()
		
		// 延迟执行逻辑
		if config.Debounce > 0 {
			timer := time.NewTimer(config.Debounce)
			go func() {
				<-timer.C
				cleanup = effectFunc()
			}()
			
			return func() {
				timer.Stop()
				if cleanup != nil {
					cleanup()
				}
			}
		}
		
		// 节流逻辑
		if config.Throttle > 0 {
			var lastRun time.Time
			var mu sync.Mutex
			
			return func() func() {
				mu.Lock()
				defer mu.Unlock()
				
				now := time.Now()
				if now.Sub(lastRun) < config.Throttle {
					return nil
				}
				
				lastRun = now
				return effectFunc()
			}()
		}
		
		// 默认直接执行
		return effectFunc()
	}, deps)
}

// UseEffectWithConfig 导出的带配置版本
func UseEffectWithConfig(effectFunc func() func(), config EffectConfig, deps ...interface{}) {
	useEffectWithConfig(effectFunc, deps, config)
}

// EffectHookInfo 副作用 Hook 的调试信息
type EffectHookInfo struct {
	ComponentID  string
	HasCleanup   bool
	Dependencies []interface{}
}

// GetEffectHooksInfo 获取所有副作用 Hooks 的调试信息
func GetEffectHooksInfo() []EffectHookInfo {
	var infos []EffectHookInfo

	for _, hooks := range components {
		for _, hook := range hooks.hooks {
			if effectHook, ok := hook.(*EffectHook); ok {
				infos = append(infos, EffectHookInfo{
					ComponentID:  hooks.component.ID(),
					HasCleanup:   effectHook.cleanupFunc != nil,
					Dependencies: effectHook.deps,
				})
			}
		}
	}

	return infos
}

// 便捷的副作用函数

// UseInterval 定时器副作用
func UseInterval(callback func(), interval time.Duration, deps ...interface{}) {
	UseEffect(func() func() {
		ticker := time.NewTicker(interval)
		go func() {
			for range ticker.C {
				callback()
			}
		}()
		
		return func() {
			ticker.Stop()
		}
	}, deps...)
}

// UseTimeout 超时副作用
func UseTimeout(callback func(), timeout time.Duration, deps ...interface{}) {
	UseEffect(func() func() {
		timer := time.NewTimer(timeout)
		go func() {
			<-timer.C
			callback()
		}()
		
		return func() {
			timer.Stop()
		}
	}, deps...)
}

// UseEventListener 事件监听器副作用
func UseEventListener(eventName string, handler func(interface{}), deps ...interface{}) {
	UseEffect(func() func() {
		// 这里需要集成到事件系统中
		// 暂时返回空清理函数
		return func() {}
	}, deps...)
}

// UseAsyncEffect 异步副作用
func UseAsyncEffect(asyncFunc func() func(), deps ...interface{}) {
	UseEffect(func() func() {
		done := make(chan bool)
		var cleanup func()
		
		go func() {
			cleanup = asyncFunc()
			done <- true
		}()
		
		return func() {
			<-done
			if cleanup != nil {
				cleanup()
			}
		}
	}, deps...)
}