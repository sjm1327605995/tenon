package api

import (
	"github.com/sjm1327605995/tenon/react/components"
)

// Renderer 是第三方渲染器需要实现的接口
type Renderer interface {
	DrawView(view *components.View)
	SetElement(element Element)
	Run() error
}
