package ui

type escEntry struct{ fn *func() }

// UseEscape 在 active 为真时注册一个 Esc 处理器；按下 Esc 时只触发最近注册的那个（栈顶），
// 用于让最上层浮层优先响应关闭。active 为假或组件卸载时自动注销。
func UseEscape(active bool, fn func()) {
	f := currentFiber
	ref := UseRef(fn)
	*ref = fn // 始终保留最新回调
	UseEffect(func() Cleanup {
		if !active || f.game == nil {
			return nil
		}
		g := f.game
		e := &escEntry{fn: ref}
		g.escStack = append(g.escStack, e)
		return func() { g.removeEsc(e) }
	}, active)
}

func (g *game) removeEsc(e *escEntry) {
	for i, x := range g.escStack {
		if x == e {
			g.escStack = append(g.escStack[:i], g.escStack[i+1:]...)
			return
		}
	}
}
