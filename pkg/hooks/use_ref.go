package hooks

// Ref 引用对象
type Ref[T any] struct {
	Current T
}

// UseRef 类似 React 的 useRef Hook
func UseRef[T any](initialValue T) *Ref[T] {
	hook := getWorkInProgressHook()

	if hook.MemoizedState == nil {
		ref := &Ref[T]{
			Current: initialValue,
		}
		hook.MemoizedState = ref
		return ref
	}

	return hook.MemoizedState.(*Ref[T])
}

// UseImperativeHandle 自定义暴露给父组件的实例值
func UseImperativeHandle[T any](ref *Ref[T], create func() T, deps []interface{}) {
	hook := getWorkInProgressHook()

	var prevDeps []interface{}
	if hook.MemoizedState != nil {
		state := hook.MemoizedState.(map[string]interface{})
		prevDeps = state["deps"].([]interface{})
	}

	if hasChanged(prevDeps, deps) {
		value := create()
		ref.Current = value

		hook.MemoizedState = map[string]interface{}{
			"deps":  deps,
			"value": value,
		}
	}
}
