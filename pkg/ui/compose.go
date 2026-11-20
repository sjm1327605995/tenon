package ui

import "reflect"

// NodeSource is implemented by UI descriptions that can be converted to nodes.
// Runtimes and declaration adapters call Build; application code can pass a
// source-returning function to ViewSource/UseSource without calling it explicitly.
type NodeSource interface {
	Build() *Node
}

func nodeFromSource(source NodeSource) *Node {
	if source == nil {
		return nil
	}
	if node, ok := source.(*Node); ok {
		return node
	}
	if value := reflect.ValueOf(source); value.Kind() == reflect.Ptr && value.IsNil() {
		return nil
	}
	if component, ok := source.(IComponent); ok {
		return Render(component)
	}
	return source.Build()
}

// View 挂载不需要 props 的函数组件，组件内部可正常使用 Hooks。
// 有参数的组件继续使用 Use(fn, props)。
func View(fn func() *Node) *Node {
	return ViewSource(fn)
}

// ViewSource mounts a function returning a node description, preserving the
// original function's component identity and converting its result automatically.
func ViewSource[V NodeSource](fn func() V) *Node {
	return &Node{
		typ:        typeComponent,
		fnPtr:      reflect.ValueOf(fn).Pointer(),
		props:      struct{}{},
		render:     func(any) *Node { return nodeFromSource(fn()) },
		propsEqual: shallowEqual,
	}
}

// UseSource is Use for functions returning node descriptions rather than *Node.
func UseSource[P any, V NodeSource](fn func(P) V, props P) *Node {
	return componentSource(fn, props, false)
}

// MemoSource is Memo for functions returning node descriptions.
func MemoSource[P any, V NodeSource](fn func(P) V, props P) *Node {
	return componentSource(fn, props, true)
}

// For 把切片映射为 Fragment，可直接与其他子节点组合，无需手动拼接切片。
// 会重排的有状态列表应为回调结果设置稳定的 Key。
func For[T any](items []T, render func(T) *Node) *Node {
	kids := make([]*Node, 0, len(items))
	for _, item := range items {
		if n := render(item); n != nil {
			kids = append(kids, n)
		}
	}
	return Fragment(kids...)
}
