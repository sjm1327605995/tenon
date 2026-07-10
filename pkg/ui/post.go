package ui

import "sync"

var (
	postMu  sync.Mutex
	postFns []func()
)

// Post 从任意 goroutine 安全地把一个函数排队到 UI（渲染）线程，在下一帧开始时执行。
// 渲染是单线程的，setState 只能在渲染线程调用；后台任务（网络请求、定时器等）完成后
// 要更新界面时，用 Post 包裹状态更新即可：
//
//	go func() {
//	    data := fetch()
//	    ui.Post(func() { setData(data) }) // 安全
//	}()
func Post(fn func()) {
	if fn == nil {
		return
	}
	postMu.Lock()
	postFns = append(postFns, fn)
	postMu.Unlock()
}

// drainPosts 在渲染线程执行所有已排队的函数（其中的 setState 因此是安全的）。
func drainPosts() {
	postMu.Lock()
	fns := postFns
	postFns = nil
	postMu.Unlock()
	for _, fn := range fns {
		fn()
	}
}
