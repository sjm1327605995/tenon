package styles

import (
	"github.com/sjm1327605995/tenon/react/components"
)

// init 包初始化函数，注册所有样式处理器
func init() {
	// 注册View组件的样式处理器
	registerViewStyleHandlers()
}

// registerViewStyleHandlers 注册View组件的样式处理器
func registerViewStyleHandlers() {
	// 注册BackgroundColor样式处理器
	RegisterGlobalHandler("View", "*styles.BackgroundColor", func(view interface{}, style interface{}) {
		if v, ok := view.(interface{ GetView() *components.View }); ok {
			bgColor := style.(BackgroundColor)
			v.GetView().Background = bgColor.Color
		}
	})

	// 注册BorderColor样式处理器
	RegisterGlobalHandler("View", "*styles.BorderColor", func(view interface{}, style interface{}) {
		if v, ok := view.(interface{ GetView() *components.View }); ok {
			borderColor := style.(BorderColor)
			v.GetView().BorderColor = borderColor.Color
		}
	})
}
