package tenon

import (
	"github.com/sjm1327605995/tenon/internal/renderer"
	"github.com/sjm1327605995/tenon/pkg/core"
)

// 导出公共类型
type Component = core.Component
type LayoutBounds = core.LayoutBounds
type Element = core.Element
type Style = core.Style
type VisualStyle = core.VisualStyle

var NewStyle = core.NewStyle
var NewVisualStyle = core.NewVisualStyle

// FunctionComponent 函数组件类型
// 函数接收 props 返回一个组件实例（通常是 View）
type FunctionComponent func(props map[string]interface{}) core.Component

// Run 启动应用
func Run(root core.Component, width, height int) {
	renderer.Run(root, width, height)
}
