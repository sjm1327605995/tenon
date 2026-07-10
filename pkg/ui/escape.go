package ui

type escEntry struct{ fn *func() }

// UseEscape 在 active 为真时注册一个 Esc 处理器；按下 Esc 时只触发最近注册的那个（栈顶），
// 用于让最上层浮层优先响应关闭。active 为假或组件卸载时自动注销。
func UseEscape(active bool, fn func()) {
	ref := UseRef(fn)
	*ref = fn // 始终保留最新回调
	UseEffect(func() Cleanup {
		if !active || activeGame == nil {
			return nil
		}
		e := &escEntry{fn: ref}
		activeGame.escStack = append(activeGame.escStack, e)
		return func() {
			if activeGame != nil {
				activeGame.removeEsc(e)
			}
		}
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
