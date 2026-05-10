package widgets

import (
	"github.com/sjm1327605995/tenon/pkg/engine"
	"github.com/sjm1327605995/tenon/pkg/render"
)

// RenderComponent 是可嵌入的 RenderObjectWidget 基类，大幅简化自定义组件开发。
//
// 设计目标：用户只关心 "怎么渲染" 和 "子组件是什么"，无需手写 CreateElement / GetKey 等样板。
// 支持嵌套其他 tenon 组件（通过 Child() 返回值）。
//
// 使用示例：
//
//	type GlowButton struct {
//	    widgets.RenderComponent
//	    label   string
//	    onClick func()
//	}
//
//	func NewGlowButton(label string, onClick func()) engine.Widget {
//	    b := &GlowButton{label: label, onClick: onClick}
//	    b.Impl = b // 关键：将自身注册为渲染实现
//	    return b
//	}
//
//	// 必须实现：创建 RenderObject
//	func (b *GlowButton) RenderObject(element engine.Element) render.RenderObject {
//	    r := render.NewRenderBox()
//	    r.SetBackgroundColor(render.NewColorFrom(color.RGBA{R: 50, G: 150, B: 255, A: 255}))
//	    r.SetOnClick(b.onClick)
//	    return r
//	}
//
//	// 可选实现：返回子组件（可以嵌套任意 tenon 组件）
//	func (b *GlowButton) Child() engine.Widget {
//	    return tenon.HStack(
//	        tenon.Icon("star"),
//	        tenon.Text(b.label),
//	    ).Gap(4)
//	}
//
// 用户只需实现 RenderObject(element) 方法；Patch 和 Child 都是可选的。
type RenderComponent struct {
	engine.BaseWidget
	// Impl 是用户的实际实现，必须在构造时赋值。
	// 它只需要实现 RenderObject(element) render.RenderObject（必须），
	// Patch(ro, oldWidget) 和 Child() engine.Widget 是可选的。
	Impl any
}

// CreateElement 返回 SingleChildRenderObjectElement。
func (c RenderComponent) CreateElement() engine.Element {
	return engine.NewSingleChildRenderObjectElement(c)
}

// GetKey 返回 NilKey（无 key）。
func (c RenderComponent) GetKey() engine.Key {
	return engine.NilKey{}
}

// CreateRenderObject 通过 Impl 委托创建 RenderObject。
// 如果 Impl 未实现 RenderObject 方法，返回一个空的 RenderBox。
func (c RenderComponent) CreateRenderObject(element engine.Element) render.RenderObject {
	if maker, ok := c.Impl.(interface {
		RenderObject(element engine.Element) render.RenderObject
	}); ok {
		return maker.RenderObject(element)
	}
	return render.NewRenderBox()
}

// UpdateRenderObject 通过 Impl 委托更新 RenderObject。
func (c RenderComponent) UpdateRenderObject(ro render.RenderObject, oldWidget engine.Widget) {
	if patcher, ok := c.Impl.(interface {
		Patch(ro render.RenderObject, oldWidget engine.Widget)
	}); ok {
		patcher.Patch(ro, oldWidget)
	}
}

// GetChildWidget 通过 Impl 委托获取子 widget。
func (c RenderComponent) GetChildWidget() engine.Widget {
	if provider, ok := c.Impl.(interface {
		Child() engine.Widget
	}); ok {
		return provider.Child()
	}
	return nil
}
