# Tenon

> Go 语言的**声明式、React 风格** GUI 工具包 —— 函数组件 + Hooks，基于 [Yoga](https://www.yogalayout.dev/) flex 布局与 [Ebiten](https://ebiten.org) 渲染。

[![Go](https://img.shields.io/badge/go-%3E%3D1.24-blue)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

[English](README.md) | **简体中文**

---

## ⚠️ 状态

年轻但自洽。核心（`pkg/ui`）形态已稳定并有测试覆盖；1.0 之前 API 仍可能调整。适合工具、仪表盘、应用内/游戏内 UI 与原型。已实现与待办见 [ROADMAP.md](ROADMAP.md)（主要缺口：字重、输入法 preedit、无障碍、多窗口）。

## 是什么

把 React 的心智模型带到 Go 原生 GUI：

- **函数组件 + Hooks** —— `UseState/UseEffect/UseReducer/UseMemo/UseCallback/UseRef/UseContext`。无 class、无手动失效。
- **自动、局部重渲染** —— setter 只重渲染它所属的组件（fiber）。
- **HTML 风格元素** —— `Div/Span/Button/Input/Img/Text/ScrollView/Portal/Fragment`。
- **Yoga flex** 布局，**Ebiten `vector` + `text/v2`** 渲染（抗锯齿、高分屏自适应）。
- **开箱即用** —— 动画（补间/进出场/FLIP）、变换、拖拽/悬停/键盘，基础组件 kit，以及一套 **shadcn/ui 风格**组件库（~41 个）。

内部是类 React 的三棵树：不可变 `Node` 描述 → 持久 `Fiber`（身份 + hooks）→ `renderNode`（yoga 节点 + 绘制）。布局**增量**：纯绘制更新不重算布局。详见 [ARCHITECTURE.md](ARCHITECTURE.md)。

## 快速上手

```bash
go get github.com/sjm1327605995/tenon
```

```go
package main

import (
	"fmt"
	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

func Counter(_ struct{}) *ui.Node {
	count, setCount := ui.UseState(0)
	return ui.Div(
		ui.Style(ui.Row, ui.Gap(16), ui.Padding(24), ui.ItemsCenter),
		ui.Button(ui.OnClick(func() { setCount(count - 1) }), ui.Text("-")),
		ui.Text(fmt.Sprintf("%d", count), ui.FontSize(32)),
		ui.Button(ui.OnClick(func() { setCount(count + 1) }), ui.Text("+")),
	)
}

func main() { ui.Run(ui.Use(Counter, struct{}{})) }
```

## 包结构

| 包 | 说明 |
|---|---|
| [`pkg/ui`](pkg/ui) | 引擎 + 元素、Hooks、样式、动画、输入。从这里开始 —— 见其 [README](pkg/ui/README.md)。 |
| [`pkg/shadcn`](pkg/shadcn) | 基于 `pkg/ui` 的 shadcn/ui 风格组件库（Button/Card/Dialog/Select/Table/Toast…）。[README](pkg/shadcn/README.md)。 |
| [`yoga`](yoga) | 纯 Go 的 Yoga flex 布局引擎移植。 |
| [`pkg/font`](pkg/font) | 字体加载/测量（内置支持 CJK 的字体）。 |

## 示例

`go run ./example/<名称>`：

`hooks-counter` · `hooks-app`（待办）· `hooks-anim` · `hooks-hover` · `hooks-drag` · `hooks-modal` · `hooks-reorder`（FLIP）· `hooks-text` · `hooks-kit` · `hooks-async`（后台 goroutine + `ui.Post`）· `shadcn-gallery`（组件总览）。

`go run ./example/shadcn-gallery` 可查看全部组件（含明暗切换）。

<!-- 添加 docs/screenshots/gallery.png（见 docs/screenshots/README.md）后取消注释：
![组件总览](docs/screenshots/gallery.png)
-->


## 后台更新

渲染是单线程的。后台 goroutine 里更新界面请用 `ui.Post` 包裹：

```go
go func() {
	data := fetch()
	ui.Post(func() { setData(data) }) // 回到渲染线程，安全
}()
```

## 贡献

欢迎 Issue 与 PR。请 `gofmt`、保持 `go test ./...` 通过、并与周边代码风格一致。见 [CONTRIBUTING.md](CONTRIBUTING.md)。

## 许可

MIT —— 见 [LICENSE](LICENSE)。
