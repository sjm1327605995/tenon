package ui

import "reflect"

// shallowEqual 逐字段浅比较两个 props：
// 函数字段按引用（指针）比较，其余字段用 DeepEqual。
// 这样带回调的 props 在回调引用稳定时可命中 Memo 短路（React 语义）。
func shallowEqual(a, b any) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	va, vb := reflect.ValueOf(a), reflect.ValueOf(b)
	if va.Type() != vb.Type() {
		return false
	}
	if va.Kind() == reflect.Struct {
		for i := 0; i < va.NumField(); i++ {
			if !fieldEqual(va.Field(i), vb.Field(i)) {
				return false
			}
		}
		return true
	}
	return fieldEqual(va, vb)
}

func fieldEqual(a, b reflect.Value) bool {
	if a.Kind() == reflect.Func {
		return a.Pointer() == b.Pointer()
	}
	if !a.CanInterface() {
		// 非导出字段无法安全比较：保守地视为不等（触发重渲染，安全）。
		return false
	}
	return reflect.DeepEqual(a.Interface(), b.Interface())
}
