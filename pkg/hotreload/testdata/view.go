package view

// 本文件是 hotreload 的测试夹具：一个被 yaegi 解释的无状态视图。
// 约束（yaegi）：可调用 pkg/ui、pkg/shadcn 的任意非泛型 API 与所有 shadcn 组件，
// 但自身不得使用泛型（即不能用 ui.UseState / ui.Use）。
// 放在 testdata/ 下，go build 会忽略它，仅在运行时被解释。

import (
	"github.com/sjm1327605995/tenon/pkg/shadcn"
	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

// View 返回被解释渲染的树，包含一段文本与一个（编译好的）shadcn 按钮。
func View() *ui.Node {
	return ui.Div(
		ui.Style(ui.Column, ui.Gap(16), ui.Padding(24)),
		ui.Text("Hot Reload 🔥", ui.FontSize(24), ui.Bold),
		shadcn.Button(shadcn.ButtonProps{}, ui.Text("Primary")),
	)
}
