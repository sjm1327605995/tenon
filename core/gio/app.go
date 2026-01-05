package gio

import (
	"image"
	"os"

	"github.com/sjm1327605995/tenon/core/ui"
	"github.com/sjm1327605995/tenon/core/ui/render"
	"github.com/sjm1327605995/tenon/yoga"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"
)

// AppConfig 应用程序配置
type AppConfig struct {
	Title  string
	Width  unit.Dp
	Height unit.Dp
}

// 状态管理，保存所有Element的Clickable状态
var clickableMap = make(map[uintptr]*ui.Element)
var clickableCounter uintptr

// RunApp 启动应用程序，渲染用户提供的UI组件
func RunApp(config AppConfig, component ui.UI) {
	// 创建上下文
	ctx := ui.NewContext(component)
	var (
		element       *ui.Element // 当前根元素
		rootRender    *render.Node
		windowSize    = image.Point{X: -1, Y: -1} // 初始化为无效值，避免首次ConfigEvent误触发
		isFirstRender = true
	)

	// 创建Gio窗口
	go func() {
		w := new(app.Window)
		w.Option(
			app.Title(config.Title),
			app.Size(config.Width, config.Height),
		)

		var ops op.Ops
		for {
			switch e := w.Event().(type) {
			case app.DestroyEvent:
				os.Exit(0)
			case app.ConfigEvent:
				// 窗口配置变化（如大小变化），需要重新布局
				if e.Config.Size.X > 0 && e.Config.Size.Y > 0 && !windowSize.Eq(e.Config.Size) {
					windowSize = e.Config.Size

					// 获取根元素
					element = ctx.GetRootElement()
					if isFirstRender {
						// 首次渲染，需要挂载Element树，创建RenderObject树
						element.Mount()
						isFirstRender = false
					}

					// 执行布局计算，使用实际像素值
					element.Yoga.CalculateLayout(float32(windowSize.X), float32(windowSize.Y), yoga.DirectionInherit)

					// 后续更新，只需要更新RenderObject的布局信息
					element.UpdateRenderObject()

					// 获取根RenderObject
					rootRender = element.RenderObject()
				}
			case app.FrameEvent:
				ui.Metric = e.Metric
				// 检查窗口大小是否变化
				if !windowSize.Eq(e.Size) && e.Size.X > 0 && e.Size.Y > 0 {
					windowSize = e.Size
					// 获取根元素
					element = ctx.GetRootElement()
					// 执行布局计算，使用实际像素值
					element.Yoga.CalculateLayout(float32(windowSize.X), float32(windowSize.Y), yoga.DirectionInherit)

					if isFirstRender {
						// 首次渲染，需要挂载Element树，创建RenderObject树
						element.Mount()
						isFirstRender = false
					}
					// 后续更新，只需要更新RenderObject的布局信息
					element.UpdateRenderObject()

					// 获取根RenderObject
					rootRender = element.RenderObject()
				}

				// 检查是否需要更新组件树
				if ctx.NeedsUpdate() {
					// 获取更新后的根元素
					element = ctx.GetRootElement()

					// 重新挂载Element树，创建RenderObject树
					element.Unmount()
					element.Mount()

					// 执行布局计算
					element.Yoga.CalculateLayout(float32(windowSize.X), float32(windowSize.Y), yoga.DirectionInherit)

					// 更新RenderObject的布局信息
					element.UpdateRenderObject()

					// 获取更新后的根RenderObject
					rootRender = element.RenderObject()

					// 清除更新标志
					ctx.ClearNeedsUpdate()
				}

				ctx := app.NewContext(&ops, e)
				// 使用RenderObject树进行绘制
				if rootRender != nil {
					rootRender.Paint(ctx)
				}

				// 交换缓冲区显示绘制结果
				e.Frame(&ops)
			}
		}
	}()

	app.Main()
}
