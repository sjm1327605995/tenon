package render

import (
	"github.com/sjm1327605995/tenon/yoga"
)

// RenderFlex 是 Flex 布局容器（Row / Column）。
// 内嵌 RenderBox，主要通过 Yoga 的 flex 能力进行布局。
type RenderFlex struct {
	RenderBox
}

func NewRenderFlex() *RenderFlex {
	r := &RenderFlex{}
	r.RenderBox.Init(r)
	r.yoga = yoga.NewNode()
	r.yoga.StyleSetFlexDirection(yoga.FlexDirectionColumn)
	return r
}

func (r *RenderFlex) PerformLayout() {
	// Yoga 布局由 Engine 在全局统一计算（从根节点调用 CalculateLayout）。
	// RenderFlex 作为布局容器，其布局逻辑已内嵌在 Yoga 中。
	// 此处在需要时可添加额外的自定义布局逻辑。
}
