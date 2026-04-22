package hooks

import (
	"reflect"
)

// StateHook 实现 useState Hook
type StateHook struct {
	value     *HookValue
	component Component
}

// NewStateHook 创建新的状态 Hook
func NewStateHook(initialValue interface{}, component Component) *StateHook {
	hook := &StateHook{
		value:     NewHookValue(initialValue),
		component: component,
	}

	// 订阅值变化，触发组件重新渲染
	hook.value.Subscribe(func(newValue interface{}) {
		if hook.component != nil {
			hook.component.MarkDirty()
		}
	})

	return hook
}

// Type 返回 Hook 类型
func (sh *StateHook) Type() HookType {
	return HookTypeState
}

// Update 更新状态值
func (sh *StateHook) Update(value interface{}) {
	sh.value.Set(value)
}

// Cleanup 清理 Hook
func (sh *StateHook) Cleanup() {
	// 状态 Hook 不需要特殊清理
}

// GetValue 获取当前值
func (sh *StateHook) GetValue() interface{} {
	return sh.value.Get()
}

// SetValue 设置新值
func (sh *StateHook) SetValue(newValue interface{}) {
	sh.value.Set(newValue)
}

// useState 类似 React 的 useState Hook
func useState(initialValue interface{}) (interface{}, func(interface{})) {
	hooks := getCurrentHooks()

	var hook *StateHook
	if hooks.index < len(hooks.hooks) {
		// 重用现有的 Hook
		existingHook := nextHook()
		var ok bool
		hook, ok = existingHook.(*StateHook)
		if !ok {
			panic("Hook type mismatch: expected StateHook")
		}
	} else {
		// 创建新的 Hook
		hook = NewStateHook(initialValue, currentComponent)
		addHook(hook)
	}

	// 返回当前值和设置函数
	return hook.GetValue(), func(newValue interface{}) {
		hook.SetValue(newValue)
	}
}

// UseState 导出的 useState Hook（泛型版本）
func UseState[T any](initialValue T) (T, func(T)) {
	currentValue, setter := useState(initialValue)

	// 类型安全的包装器
	typedSetter := func(newValue T) {
		setter(newValue)
	}

	// 确保返回正确的类型
	if typedValue, ok := currentValue.(T); ok {
		return typedValue, typedSetter
	}

	// 如果类型不匹配，返回初始值
	return initialValue, typedSetter
}

// UseStateInt 整数类型的便捷函数
func UseStateInt(initialValue int) (int, func(int)) {
	return UseState(initialValue)
}

// UseStateString 字符串类型的便捷函数
func UseStateString(initialValue string) (string, func(string)) {
	return UseState(initialValue)
}

// UseStateBool 布尔类型的便捷函数
func UseStateBool(initialValue bool) (bool, func(bool)) {
	return UseState(initialValue)
}

// UseStateFloat 浮点数类型的便捷函数
func UseStateFloat(initialValue float64) (float64, func(float64)) {
	return UseState(initialValue)
}

// StateConfig 状态配置选项
type StateConfig struct {
	// 比较函数，用于判断值是否变化
	Equals func(prev, next interface{}) bool
	// 延迟更新，批量处理状态变更
	Debounce bool
}

// useStateWithConfig 带配置的 useState
func useStateWithConfig(initialValue interface{}, config StateConfig) (interface{}, func(interface{})) {
	currentValue, setter := useState(initialValue)

	// 应用配置的包装器
	configuredSetter := func(newValue interface{}) {
		if config.Equals != nil {
			// 使用自定义比较函数
			if !config.Equals(currentValue, newValue) {
				setter(newValue)
			}
		} else {
			// 默认深度比较
			if !reflect.DeepEqual(currentValue, newValue) {
				setter(newValue)
			}
		}
	}

	return currentValue, configuredSetter
}

// UseStateWithConfig 导出的带配置版本
func UseStateWithConfig[T any](initialValue T, config StateConfig) (T, func(T)) {
	currentValue, setter := useStateWithConfig(initialValue, config)

	typedSetter := func(newValue T) {
		setter(newValue)
	}

	if typedValue, ok := currentValue.(T); ok {
		return typedValue, typedSetter
	}

	return initialValue, typedSetter
}

// StateHookInfo 状态 Hook 的调试信息
type StateHookInfo struct {
	Value      interface{}
	ComponentID string
}

// GetStateHooksInfo 获取所有状态 Hooks 的调试信息
func GetStateHooksInfo() []StateHookInfo {
	var infos []StateHookInfo

	for _, hooks := range components {
		for _, hook := range hooks.hooks {
			if stateHook, ok := hook.(*StateHook); ok {
				infos = append(infos, StateHookInfo{
					Value:       stateHook.GetValue(),
					ComponentID: hooks.component.ID(),
				})
			}
		}
	}

	return infos
}