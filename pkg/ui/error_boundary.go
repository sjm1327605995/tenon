package ui

import "fmt"

type ebProps struct {
	fallback func(err any, retry func()) *Node
	children []*Node
}

// ErrorBoundary 捕获其子树在 render 期抛出的 panic，改为渲染 fallback，
// 避免一个组件的 bug 让整个应用崩溃（类似 React 的 Error Boundary）。
// fallback 收到 panic 值与一个 retry 回调（清除错误、重新渲染子树，用于可恢复的错误）。
// 无 ErrorBoundary 兜底时 panic 照常向上抛出（保留原始栈，便于开发期调试）。
//
//	ui.ErrorBoundary(
//	    func(err any, retry func()) *ui.Node {
//	        return ui.Div(ui.Text("出错了: "+ui.ErrText(err)),
//	            ui.Button(ui.OnClick(retry), ui.Text("重试")))
//	    },
//	    RiskyComponent(...),
//	)
func ErrorBoundary(fallback func(err any, retry func()) *Node, children ...*Node) *Node {
	return Use(errorBoundaryC, ebProps{fallback: fallback, children: children})
}

func errorBoundaryC(p ebProps) *Node {
	f := currentFiber
	f.errBoundary = true
	// 同步计数（必须在子树渲染前就 >0，否则初次挂载时的 panic 无法被捕获）。
	counted := UseRef(false)
	if !*counted {
		*counted = true
		boundaryCount++ // 卸载时在 unmount 里递减
	}
	if f.caughtErr != nil {
		err := f.caughtErr
		retry := func() {
			f.caughtErr = nil
			if activeGame != nil {
				activeGame.markDirty(f)
			}
		}
		if p.fallback != nil {
			return p.fallback(err, retry)
		}
		return nil
	}
	return Fragment(p.children...)
}

// ErrText 返回 panic 值的可读字符串，便于在 fallback 中显示。
func ErrText(err any) string { return fmt.Sprint(err) }
