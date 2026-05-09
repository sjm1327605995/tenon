package engine

import "github.com/sjm1327605995/tenon/pkg/render"

// RenderObjectFactory 负责根据 Widget 配置创建对应的 RenderObject。
// element 参数供闭包获取最新 Widget 状态（对应 Flutter 的 BuildContext）。
type RenderObjectFactory interface {
	CreateRenderObject(element Element) render.RenderObject
}

// RenderObjectUpdater 负责在 Widget 变化时将新属性同步到 RenderObject。
type RenderObjectUpdater interface {
	UpdateRenderObject(ro render.RenderObject, oldWidget Widget)
}

// SingleChildProvider 为 SingleChildRenderObjectElement 提供单个子 Widget。
type SingleChildProvider interface {
	GetChildWidget() Widget
}

// MultiChildProvider 为 MultiChildRenderObjectElement 提供多个子 Widget。
type MultiChildProvider interface {
	GetChildrenWidgets() []Widget
}
