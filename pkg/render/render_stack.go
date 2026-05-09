package render

import (
	"github.com/sjm1327605995/tenon/yoga"
)

// RenderStack 是层叠布局容器，支持绝对定位子元素和 Z-Index。
// 对应 Flutter 的 RenderStack 和 CSS 的 position:absolute/relative。
type RenderStack struct {
	RenderBox
}

func NewRenderStack() *RenderStack {
	r := &RenderStack{}
	r.RenderBox.Init(r)
	r.yoga = yoga.NewNode()
	return r
}

func (r *RenderStack) PerformLayout() {
	// Yoga 已经处理了 absolute 子节点的布局（通过 left/top/right/bottom）
	// RenderStack 不需要额外的布局逻辑
}
