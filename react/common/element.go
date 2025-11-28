package common

import (
	"github.com/sjm1327605995/tenon/react/yoga"
)

// Element 是所有可渲染元素的基础接口，为第三方渲染器提供必要的方法
type Element interface {
	// Yoga 返回元素的Yoga节点，用于获取布局信息（如位置、大小等）
	Yoga() *yoga.Node
	// GetChildrenCount 返回子元素的数量
	GetChildrenCount() int
	// GetChildAt 返回指定索引的子元素
	GetChildAt(index int) Element
	Rendering(renderer Renderer)
	GetChildren() []Element
}
