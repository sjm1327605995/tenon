// Package router 是 Tenon 的栈式导航器——参照 React Navigation 的「导航栈」模型，
// 而非 Web 的 URL 路由（桌面单窗口没有 URL）。一个 Router 维护一段路由栈，按名字渲染
// 栈顶屏；屏内通过 UseNavigate 入栈/出栈/替换，通过 UseRoute 读取当前路由与参数。
//
// 整个实现只建立在 pkg/ui 的 hooks（UseState + UseContext）之上——核心引擎无需改动。
// 每个屏由 ui.Use(screen, params) 挂成独立子组件，因此切屏时旧屏卸载、新屏挂载，
// 各屏自己的 UseState 等 hooks 天然隔离。
//
//	router.Router(router.Props{
//	    Initial: "list",
//	    Screens: map[string]router.Screen{
//	        "list":   ListScreen,
//	        "detail": DetailScreen,
//	    },
//	})
//
//	// 屏内：
//	nav := router.UseNavigate()
//	nav.Push("detail", router.Params{"id": "42"}) // 入栈
//	nav.Pop()                                      // 返回
//	nav.Replace("home", nil)                       // 原地替换
//	r := router.UseRoute()                         // r.Name, r.Params["id"]
package router

import ui "github.com/sjm1327605995/tenon/pkg/ui"

// Params 是一次导航携带的参数。
type Params map[string]string

// Route 是栈中的一项：目标屏名字与其参数。
type Route struct {
	Name   string
	Params Params
}

// Screen 按参数渲染一个屏。
type Screen = func(Params) *ui.Node

// Props 配置一个 Router。
type Props struct {
	Initial string            // 初始屏名字
	Params  Params            // 初始屏参数（可选）
	Screens map[string]Screen // 名字 -> 屏
}

// Navigator 是屏内的导航接口（用 UseNavigate 获取）。其方法改变路由栈并触发重渲染。
type Navigator struct {
	stack []Route
	set   func([]Route)
}

// Current 返回栈顶路由；Depth 是栈深度；CanPop 表示能否返回（深度 > 1）。
func (n *Navigator) Current() Route { return n.stack[len(n.stack)-1] }
func (n *Navigator) Depth() int     { return len(n.stack) }
func (n *Navigator) CanPop() bool   { return len(n.stack) > 1 }

// Push 入栈一个新屏。
func (n *Navigator) Push(name string, params Params) {
	n.set(append(n.clone(), Route{Name: name, Params: params}))
}

// Replace 用一个新屏替换栈顶（不改变深度）。
func (n *Navigator) Replace(name string, params Params) {
	s := n.clone()
	s[len(s)-1] = Route{Name: name, Params: params}
	n.set(s)
}

// Pop 返回上一屏（仅当 CanPop 时生效）。
func (n *Navigator) Pop() {
	if n.CanPop() {
		s := n.clone()
		n.set(s[:len(s)-1])
	}
}

// PopToRoot 一路返回到栈底屏。
func (n *Navigator) PopToRoot() {
	if n.CanPop() {
		n.set([]Route{n.stack[0]})
	}
}

func (n *Navigator) clone() []Route {
	c := make([]Route, len(n.stack))
	copy(c, n.stack)
	return c
}

// navCtx 把当前 Router 的 Navigator 传给子树中的屏。默认 nil（在 Router 之外调用
// UseNavigate 会得到 nil）。
var navCtx = ui.CreateContext[*Navigator](nil)

// Router 渲染栈顶屏，并向子树提供导航能力。Router 可嵌套（内层 Provider 会遮蔽外层）。
func Router(p Props) *ui.Node { return ui.Use(routerImpl, p) }

func routerImpl(p Props) *ui.Node {
	stack, setStack := ui.UseState([]Route{{Name: p.Initial, Params: p.Params}})
	nav := &Navigator{stack: stack, set: setStack}
	cur := nav.Current()

	screen := p.Screens[cur.Name]
	if screen == nil {
		return navCtx.Provider(nav, ui.Text(`router: 未注册的路由 "`+cur.Name+`"`))
	}
	// 用 ui.Use 把屏挂成独立子组件：切屏时 screen 函数指针变化 -> 旧屏卸载、新屏挂载，
	// 各屏 hooks 隔离；同屏换参数则复用同一 fiber、以新 props 重渲染。
	return navCtx.Provider(nav, ui.Use(screen, cur.Params))
}

// UseNavigate 返回当前 Router 的导航器。仅可在 Router 子树内调用。
func UseNavigate() *Navigator { return ui.UseContext(navCtx) }

// UseRoute 返回当前路由（名字 + 参数）。
func UseRoute() Route {
	if n := ui.UseContext(navCtx); n != nil {
		return n.Current()
	}
	return Route{}
}
