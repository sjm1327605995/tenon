package ui

import "reflect"

// shallowEqual 逐字段比较两个 props。非函数字段用 DeepEqual；
// 非 nil 函数字段保守地视为不等，因为 Go 无法区分入口地址相同、捕获值不同的闭包。
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
		return a.IsNil() && b.IsNil()
	}
	if !a.CanInterface() {
		// 非导出字段无法安全比较：保守地视为不等（触发重渲染，安全）。
		return false
	}
	return reflect.DeepEqual(a.Interface(), b.Interface())
}
