package api

import (
	"github.com/sjm1327605995/tenon/react/api/styles"
)

// Element 是所有可渲染元素的基础接口，为第三方渲染器提供必要的方法
type Element interface {
	styles.StyleElement
	// GetChildrenCount 返回子元素的数量
	GetChildrenCount() int
	// GetChildAt 返回指定索引的子元素
	GetChildAt(index int) Element
	Rendering(renderer Renderer)
	GetChildren() []Element
	// SetExtendedStyle 设置扩展样式，由每个元素自行实现具体逻辑
	SetExtendedStyle(style styles.IExtendedStyle)
	SetStyle(style *styles.Style)
}
