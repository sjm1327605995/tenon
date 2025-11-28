package common

import (
	"github.com/sjm1327605995/tenon/react/common/component"
)

// Renderer 是第三方渲染器需要实现的接口
type Renderer interface {
	DrawView(view *component.View)
	Run() error
	SetElement(element Element)
}
