package hooks

// UseMemo 类似 React 的 useMemo Hook
func UseMemo[T any](factory func() T, deps []interface{}) T {
	hook := getWorkInProgressHook()

	var prevDeps []interface{}
	var prevValue T
	if hook.MemoizedState != nil {
		state := hook.MemoizedState.(map[string]interface{})
		prevDeps = state["deps"].([]interface{})
		if v, ok := state["value"].(T); ok {
			prevValue = v
		}
	}

	if hasChanged(prevDeps, deps) {
		newValue := factory()
		hook.MemoizedState = map[string]interface{}{
			"deps":  deps,
			"value": newValue,
		}
		return newValue
	}

	return prevValue
}

// UseCallback 类似 React 的 useCallback Hook
func UseCallback[T any](callback T, deps []interface{}) T {
	return UseMemo(func() T {
		return callback
	}, deps)
}
