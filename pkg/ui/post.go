package ui

import "sync"

var (
	postMu  sync.Mutex
	postFns []func()
)

// backendWake 唤醒后端的事件循环，由当前后端登记；可从任意 goroutine 调用。
// 界面按需重绘（静止时不出帧）后，后台任务排队的更新必须靠它把循环叫醒，
// 否则要等到下一次用户输入才会生效。
var backendWake func()

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
	if backendWake != nil {
		backendWake() // 界面静止时循环是睡着的，必须叫醒它来执行本次排队
	}
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
