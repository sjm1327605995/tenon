# Tenon

> Go 语言的细粒度响应式 UI 框架，基于 Ebiten 和 Yoga 构建。

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.25-blue)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

[English](README.md) | **简体中文**

---

## Tenon 是什么？

Tenon 是一款面向 Go 语言的跨平台 UI 框架，将现代响应式编程模式引入桌面应用开发。不同于传统的即时模式或 VDOM 框架，Tenon 采用**细粒度响应式架构**——状态变更直接流向原生元素，无需 Diff 开销，无需不必要的重建。

- **Widget** 描述结构 —— 轻量、有状态，`Render()` 产出元素树
- **Element** 即原生组件 —— 持久化、可复用、持有自身的 Yoga 节点
- **State** 即信号 —— 变更直接通知订阅者，绕过 Widget 重建
- **Engine** 统筹全局 —— 布局（Yoga）、渲染（Ebiten）、事件

## 核心特性

- **细粒度响应式** —— 状态变更驱动元素直接更新，属性更新零 Diff 开销
- **50+ 内置组件** —— View、Text、Button、Input、Modal、Toast、ScrollView、Table、Sidebar、Drawer、Tabs 等
- **Flexbox 布局** —— 完整的 Yoga 布局引擎，支持 gap、align、justify、wrap
- **主题系统** —— 内置亮色 / 暗色 / Shadcn 主题，完全可定制
- **动画引擎** —— Tween 补间动画，多种缓动函数
- **样式系统** —— 通过标签和类注册全局或类型级样式
- **条件渲染** —— `Switch` 构建器实现清晰的多分支 UI 逻辑
- **自动依赖追踪** —— `Render()` 期间访问的 State 自动订阅
- **跨平台** —— 通过 Ebiten 支持 Windows、macOS、Linux
- **调试工具** —— 内置远程调试器，可检查元素树

## 快速开始

```go
package main

import (
    "image/color"
    
    "github.com/sjm1327605995/tenon"
    "github.com/sjm1327605995/tenon/pkg/fonts"
    "github.com/sjm1327605995/tenon/pkg/v2/components"
    "github.com/sjm1327605995/tenon/pkg/v2/core"
    "github.com/sjm1327605995/tenon/yoga"
)

type App struct {
    core.BaseWidget
    count *core.State[int]
}

func NewApp() *App {
    a := &App{count: core.NewState(0)}
    a.Init(a)
    return a
}

func (a *App) Render() core.Element {
    return components.NewView().
        SetWidthPercent(100).
        SetHeightPercent(100).
        SetPadding(yoga.EdgeAll, 24).
        SetBackgroundColor(core.GetTheme().BackgroundColor).
        Add(
            components.NewText("Hello, Tenon!").
                SetFontSize(24).
                SetColor(core.GetTheme().TextColor),
            components.NewButton("Clicked: 0").
                SetOnClick(func() {
                    a.count.Set(a.count.Get() + 1)
                }),
        )
}

func main() {
    fonts.InitDefaultFont()
    core.SetTheme(core.DefaultShadcnLightTheme())
    
    tenon.Run(NewApp(), 800, 600)
}
```

## 架构设计

```
┌─────────────────────────────────────────┐
│  Widget（结构描述层）                    │
│  Render() → Element 树                  │
│  仅在挂载时和结构变化时调用              │
├─────────────────────────────────────────┤
│  Element（持久节点层）                   │
│  *View, *Text, *Button...               │
│  持有 Yoga 节点，负责绘制和事件          │
│  链式 Setter 立即生效                    │
├─────────────────────────────────────────┤
│  State（细粒度信号层）                   │
│  Set() 通知订阅者                       │
│  Render() 期间自动追踪依赖               │
├─────────────────────────────────────────┤
│  Engine（引擎层）                        │
│  Build 队列 → Yoga 布局 → 绘制循环      │
│  事件路由 → 动画帧驱动                   │
└─────────────────────────────────────────┘
```

### 两条更新路径

**路径 A — 属性更新（高频，零 Diff）**
```
用户操作 → count.Set(42) → State 通知订阅者
→ text.SetContent("42") → MarkNeedDraw → 下一帧重绘
```
不调用 `Render()`，不做 Diff，不重建 Yoga 节点。仅更新订阅了该 State 的元素。

**路径 B — 结构更新（低频，同级浅对比）**
```
页面切换 → RequestBuild() → Render() 新树
→ 同级类型对比 → 复用 / 替换 / 移动 Element
→ 标记 Yoga dirty → 下一帧重新布局
```
仅比较直接子节点的类型，不递归。属性不在这里 Diff —— 它们由路径 A 处理。

## 组件总览

Tenon 提供覆盖桌面 UI 全场景的丰富组件库：

**布局** — `View`、`ScrollView`、`SplitView`、`Resizable`、`AspectRatio`、`Sidebar`、`Sheet`、`Drawer`

**展示** — `Text`、`Image`、`SVGIcon`、`Badge`、`Divider`、`Kbd`、`Skeleton`、`Table`、`Carousel`

**表单控件** — `Button`（5 种变体）、`TextInput`、`TextArea`、`Select`、`Checkbox`、`Radio`、`RadioGroup`、`Switch`、`Slider`、`InputOTP`

**反馈** — `Modal`、`Alert`、`AlertDialog`、`Toast`、`Tooltip`、`Popover`、`ProgressBar`、`LoadingSpinner`

**导航** — `Tab`、`Menu`、`Menubar`、`Breadcrumb`、`Pagination`、`Command`、`NavigationMenu`

**浮层** — `Dropdown`、`ContextMenu`、`HoverCard`、`FloatingButton`

**数据** — `Accordion`、`Collapsible`、`Calendar`、`ListView`

## 主题定制

Tenon 内置三套主题，并支持完全自定义：

```go
// 使用内置主题
core.SetTheme(core.DefaultShadcnLightTheme())
core.SetTheme(core.DefaultShadcnDarkTheme())
core.SetTheme(core.DefaultAntTheme())

// 或自定义主题
theme := &core.Theme{
    PrimaryColor:      color.RGBA{59, 130, 246, 255},
    BackgroundColor:   color.RGBA{255, 255, 255, 255},
    TextColor:         color.RGBA{15, 23, 42, 255},
    BorderRadius:      8,
    // ... 完整主题字段
}
core.SetTheme(theme)
```

所有组件默认从当前主题读取配色，可通过链式 API 单独覆盖。

## 动画系统

```go
// Tween 补间动画
tween := core.NewTween(300*time.Millisecond, core.EaseOutCubic).
    OnUpdate(func(p float32) {
        el.SetOpacity(p)
    }).
    OnComplete(func() {
        // 动画完成
    })
tween.Start()

// 状态驱动过渡
core.NewTween(200*time.Millisecond, core.EaseInOut).
    OnUpdate(func(p float32) {
        x := core.LerpFloat32(0, 200, p)
        el.SetPosition(x, 0)
    }).Start()
```

## 样式系统

注册全局或类型级可复用样式：

```go
// 全局样式（按类名）
tenon.RegisterStyle("card", func(e tenon.Element) {
    if v, ok := e.(*components.View); ok {
        v.SetBackgroundColor(core.GetTheme().CardColor).
          SetBorderRadius(core.GetTheme().BorderRadius).
          SetShadow(color.RGBA{A: 20}, 12, 0, 2)
    }
})

// 在 Render() 中使用
components.NewView().SetClass("card").Add(...)
```

## 示例

探索 [`example/`](example/) 目录：

| 示例 | 说明 |
|------|------|
| [`v2-demo`](example/v2-demo) | 完整组件库导航演示 |
| [`shadcn-gallery`](example/shadcn-gallery) | Shadcn/UI 风格样式展示 |
| [`card`](example/card) | 卡片布局演示 |

```bash
cd example/v2-demo
go run main.go
```

## 文档

- **[架构深度解析](ARCHITECTURE.md)** —— 设计决策、更新机制、内部结构
- **[贡献指南](CONTRIBUTING.md)** —— 开发环境搭建、代码规范、PR 流程
- **[API 参考](pkg/v2/core/)** —— 核心接口：`Widget`、`Element`、`State`、`Engine`、`Theme`、`Animation`

## 项目结构

```
tenon/
├── tenon.go              # 公共 API 入口：Run()、类型别名、工具函数
├── pkg/
│   ├── v2/
│   │   ├── core/         # 引擎、Widget、Element、State、Theme、Animation
│   │   └── components/   # 50+ 内置 UI 组件
│   └── fonts/            # 字体加载与字形管理
├── yoga/                 # Yoga Flexbox 布局引擎（Go 移植版）
├── example/              # 示例应用
├── ARCHITECTURE.md       # 架构文档
├── CONTRIBUTING.md       # 贡献指南
└── LICENSE               # MIT 许可证
```

## 为什么选择 Tenon？

| | 即时模式 (imgui) | VDOM (React-like) | **Tenon** |
|---|:---:|:---:|:---:|
| State → UI | 每帧手动处理 | Diff + Patch | 直接信号订阅 |
| 状态变更时 Widget 重建 | 不适用 | 整树重渲染 | 自动追踪，最小重建 |
| Element 持久化 | 每帧重新创建 | Reconcile 对比 | 持久化、可复用 |
| 性能模型 | CPU 密集型绘制调用 | Diff 开销 | 属性更新零 Diff |
| Go 语言体验 | 过程式 | 到处 Hooks | 地道 Go，基于结构体 |

Tenon 处于独特的交叉点：它提供**细粒度响应式的性能特征**（类似 Solid.js、Svelte），配合**Go 结构体和接口的编程体验**，通过**高性能 2D 引擎**渲染。

## 许可证

MIT License —— 详见 [LICENSE](LICENSE)。
