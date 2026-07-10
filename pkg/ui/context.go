package ui

// Context 是一个跨层数据通道（相当于 React.createContext）。
type Context[T any] struct {
	id  int
	def T
}

var ctxIDSeq int

// CreateContext 创建一个带默认值的 Context。
func CreateContext[T any](def T) *Context[T] {
	ctxIDSeq++
	return &Context[T]{id: ctxIDSeq, def: def}
}

// Provider 提供 value 给子树中的消费者。
func (c *Context[T]) Provider(value T, children ...*Node) *Node {
	n := &Node{typ: typeProvider, ctxID: c.id, ctxValue: value}
	for _, ch := range children {
		if ch != nil && ch.typ != typeAttr {
			n.kids = append(n.kids, ch)
		}
	}
	return n
}

// UseContext 读取最近 Provider 的值，并订阅其变化。
func UseContext[T any](c *Context[T]) T {
	f := currentFiber
	for p := f.parent; p != nil; p = p.parent {
		if p.typ == typeProvider && p.ctxID == c.id {
			p.subscribers = append(p.subscribers, f)
			return p.ctxValue.(T)
		}
	}
	return c.def
}
