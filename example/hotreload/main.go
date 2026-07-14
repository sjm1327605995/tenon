package main

import (
	"os"

	"github.com/sjm1327605995/tenon/pkg/hotreload"
	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

// 热更新预览：编辑 example/hotreload/view/view.go 并保存，运行中的窗口即时刷新（无需重启）。
// 用法：从仓库根目录运行 `go run ./example/hotreload`，然后改 view/view.go 保存看效果。
func main() {
	path := "example/hotreload/view/view.go"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	ui.ShowStats = true
	ui.Run(hotreload.LivePreview(path))
}
