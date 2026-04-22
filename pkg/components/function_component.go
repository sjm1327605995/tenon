package components

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/pkg/hooks"
	"github.com/sjm1327605995/tenon/pkg/reconciler"
)

// FunctionComponentWrapper 包装函数组件
// 将 React 风格的函数组件适配到组件树中
type FunctionComponentWrapper struct {
	core.BaseComponent
	fn           func(props map[string]interface{}) core.Component
	props        map[string]interface{}
	rendered     core.Component
	fiber        *reconciler.Fiber
	isFirstRender bool
}

// NewFunctionComponent 创建函数组件包装器
func NewFunctionComponent(fn func(props map[string]interface{}) core.Component, props map[string]interface{}) *FunctionComponentWrapper {
	f := &FunctionComponentWrapper{
		BaseComponent: core.NewBaseComponent(),
		fn:            fn,
		props:         props,
		isFirstRender: true,
	}
	f.Init(f)
	return f
}

// SetProps 设置 props
func (f *FunctionComponentWrapper) SetProps(props map[string]interface{}) *FunctionComponentWrapper {
	f.props = props
	f.MarkDirty()
	return f
}

func (f *FunctionComponentWrapper) Update() error {
	// 每次 Update 都重新执行函数组件（render）
	f.render()
	return nil
}

func (f *FunctionComponentWrapper) render() {
	// 创建或复用 Fiber
	if f.fiber == nil {
		f.fiber = reconciler.CreateFiber(reconciler.FunctionComponent, "function", "")
		f.fiber.StateNode = f
		f.fiber.RenderFn = f.fn
		f.fiber.RenderProps = f.props
	}

	// 设置当前 Fiber 供 hooks 使用
	hooks.SetCurrentFiber(f.fiber)
	reconciler.SetWorkInProgress(f.fiber)
	reconciler.ResetWorkInProgressHook()

	// 执行函数组件
	newRendered := f.fn(f.props)

	// 重置 hooks 上下文
	hooks.SetCurrentFiber(nil)
	reconciler.ResetWorkInProgressHook()

	// 更新子组件
	if f.rendered != nil {
		f.RemoveChild(f.rendered)
	}
	f.rendered = newRendered
	if newRendered != nil {
		f.AddChild(newRendered)
	}
}

func (f *FunctionComponentWrapper) Draw(screen *ebiten.Image) {
	// 函数组件自己不绘制，由子组件绘制
	for _, child := range f.GetChildren() {
		child.Draw(screen)
	}
}

func (f *FunctionComponentWrapper) DrawOverlay(screen *ebiten.Image) {
	for _, child := range f.GetChildren() {
		child.DrawOverlay(screen)
	}
}

func (f *FunctionComponentWrapper) HandleInput() bool {
	children := f.GetChildren()
	for i := len(children) - 1; i >= 0; i-- {
		if children[i].HandleInput() {
			return true
		}
	}
	return false
}
