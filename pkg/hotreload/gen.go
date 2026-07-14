package hotreload

// 重新生成解释器可见的符号表（当 pkg/ui 或 pkg/shadcn 的导出 API 变化时运行 `go generate ./pkg/hotreload`）。
// yaegi extract 会跳过泛型函数（解释器不支持），非泛型 API 与所有 shadcn 组件均可导出。
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract github.com/sjm1327605995/tenon/pkg/ui
//go:generate go run github.com/traefik/yaegi/cmd/yaegi extract github.com/sjm1327605995/tenon/pkg/shadcn
