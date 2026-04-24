# Tenon - 细粒度响应式 UI 框架

Tenon 是一个基于 Go 的跨平台 UI 框架，采用**细粒度响应式**架构：

- **Widget 负责描述结构**，`Build()` 产出声明式的 Element 树
- **Element 就是 Native 组件**（`*View`、`*Text`、`*Button`…），持久化、可复用、直接持有 Yoga 节点
- **State 即信号**，属性变更直接驱动 Element 局部更新，**不触发 Widget 重新 Build**
- **同级浅 Diff** 只用于结构协调（条件渲染、列表增删），不做属性级全量对比

渲染后端基于 [Ebiten](https://ebitengine.org/)，布局引擎基于 [Yoga](https://www.yogalayout.dev/)。

---

## 目录

- [核心架构](#核心架构)
- [三层模型](#三层模型)
- [更新机制](#更新机制)
- [快速开始](#快速开始)
- [Widget 开发](#widget-开发)
- [Element 开发](#element-开发)
- [项目结构](#项目结构)

---

## 核心架构

```
┌─────────────────────────────────────────────────────────────┐
│                      Widget（结构描述）                       │
│              轻量对象，Build() 产出 Element 树                │
│              只在结构变化时重新 Build                         │
├─────────────────────────────────────────────────────────────┤
│                      Element（持久节点）                      │
│         *View, *Text, *Button... 直接就是 Element            │
│         持有 Yoga 节点，负责绘制、事件、属性更新               │
├─────────────────────────────────────────────────────────────┤
│                      State（细粒度信号）                      │
│         属性变化直接通知订阅者，不走 Build / Diff             │
├─────────────────────────────────────────────────────────────┤
│                      Engine（框架核心）                       │
│         Build Diff（结构协调）→ Yoga 布局 → 绘制 → 事件      │
└─────────────────────────────────────────────────────────────┘
```

### 关键设计决策

| 问题 | 决策 | 原因 |
|------|------|------|
| Element 是什么？ | `*View`、`*Button` 等 Native 组件本身就是 Element | 消除包装层，减少抽象损耗 |
| 属性如何更新？ | State 信号直接驱动 Element Setter | 零 Diff 开销，精准局部刷新 |
| 结构如何更新？ | Widget `Build()` 产出新树，框架做**同级浅对比** | 条件渲染、列表增删需要协调新旧子树 |
| 是否需要 VDOM？ | 否。只做同级类型对比，不递归属性 Diff | Go 没有 JSX，全量 Diff 收益低 |
| Yoga 角色 | Element 内嵌 Yoga 节点，1:1 绑定 | Element 树即布局树 |

---

## 三层模型

### 1. Widget —— 结构描述层

Widget 是用户编写的轻量对象，**只关心子树长什么样**，不关心怎么绘制、怎么更新。

```go
type Counter struct {
    core.Widget
    count *core.State[int]
}

func NewCounter() *Counter {
    c := &Counter{}
    c.count = core.NewState(0)
    return c
}

func (c *Counter) Build() core.Element {
    text := components.NewText("").
        SetFontSize(18)
    
    // 状态变化直接驱动 Element 更新，不走 Build
    c.count.Subscribe(func(v int) {
        text.SetContent(fmt.Sprintf("Count: %d", v))
    })
    
    btn := components.NewButton("+1").
        SetOnClick(func() {
            c.count.Set(c.count.Get() + 1)
        })
    
    return components.NewView().
        SetPadding(yoga.EdgeAll, 20).
        Add(text, btn)
}
```

**Widget 的规则：**
- `Build()` 只在**首次挂载**和**结构变化**时调用
- 不要在 `Build()` 里做异步操作、副作用、手动订阅清理
- 属性变化（文字、颜色、状态）全部通过 `State` + `Element.Setter` 处理

---

### 2. Element —— 持久渲染层

`*components.View`、`*components.Text`、`*components.Button`… 这些就是 Element。它们：

- 内嵌 `*yoga.Node`，参与 Flexbox 布局
- 自己管理绘制状态（`Draw`）和交互（`HandleEvent`）
- 提供链式 Setter，调用后立即生效，自行决定是否需要重绘
- 组成一棵持久的树，`Parent` / `Children` 由框架维护

```go
type View struct {
    core.BaseElement
    backgroundColor color.Color
}

func (v *View) Draw(screen *ebiten.Image) {
    // 按自己的 bounds 绘制背景
}

func (v *View) SetBackgroundColor(c color.Color) *View {
    v.backgroundColor = c
    v.MarkNeedDraw()  // 标记需要重绘，下一帧刷新
    return v
}
```

**Element 的规则：**
- Element 是**长期存活**的，除非被框架从树上移除
- Setter 是**命令式**的：调用即生效，不需要等 Build
- `BaseElement` 提供子节点管理（`AppendChild` / `RemoveChild`），但日常增删建议让框架通过 Build Diff 处理

---

### 3. State —— 细粒度信号层

`State[T]` 是可观察的状态容器。`Set()` 时只通知订阅者，**不触发 Widget Build**。

```go
// 创建状态
count := core.NewState(0)

// 订阅：变化时直接操作 Element
count.Subscribe(func(v int) {
    text.SetContent(fmt.Sprintf("%d", v))
})

// 修改：只触发订阅回调
count.Set(42)
```

**State 的生命周期：**
- State 通常由 Widget 持有
- State 的订阅可以绑定到 Element，也可以绑定到任意回调
- Widget 卸载时，框架负责清理其关联的 State 订阅

---

## 更新机制

Tenon 有两条完全独立的更新路径：

### 路径 A：属性更新（高频，零 Diff）

```
用户操作 / 定时器 / 网络回调
    ↓
count.Set(42)
    ↓
State 通知所有订阅者
    ↓
text.SetContent("42")   ← Element 自己更新，自己标记重绘
    ↓
下一帧 Engine 只重绘该 Element
```

**特点：**
- 不调用 `Build()`
- 不做 Diff
- 不重建 Yoga 节点
- 只影响订阅该 State 的 Element

---

### 路径 B：结构更新（低频，同级浅 Diff）

当 Widget 主动请求结构更新（如切换页面、条件渲染变化）：

```
Widget.RequestBuild()  或  父级结构变化
    ↓
框架调用 Widget.Build()，产出新的 Element 子树
    ↓
同级浅对比：新子节点 vs 旧子节点（只比直接子节点，不递归）
    ↓
  ├─ 类型相同 → 复用旧 Element，同步 Yoga 样式
  ├─ 类型不同 → 卸载旧 Element，挂载新 Element
  ├─ 位置变化 → 移动 Element
  └─ 新增 → 挂载新 Element
    ↓
标记受影响的 Yoga 节点 dirty，下一帧重新布局
```

**Diff 规则（极简）：**

```go
// 只对比同一父节点的直接子节点
for i, newChild := range newChildren {
    if i < len(oldChildren) && oldChildren[i].Type() == newChild.Type() {
        // 复用：保留 Yoga 节点，只同步样式
        syncStyle(oldChildren[i].Yoga, newChild.Yoga)
    } else {
        // 替换或新增
        parent.InsertChild(newChild, i)
    }
}
// 删除多余的旧子节点
for i := len(newChildren); i < len(oldChildren); i++ {
    parent.RemoveChild(oldChildren[i])
}
```

**特点：**
- 只做**一层**，不递归到孙子节点
- 只对比**类型**（`"View"`、`"Text"`…），不对比属性
- 属性由路径 A（State）处理，Diff 不管

---

## 快速开始

```go
package main

import (
    "image/color"
    
    "github.com/sjm1327605995/tenon"
    "github.com/sjm1327605995/tenon/pkg/components"
    "github.com/sjm1327605995/tenon/yoga"
)

type App struct {
    core.Widget
}

func (a *App) Build() core.Element {
    return components.NewView().
        SetWidthPercent(100).
        SetHeightPercent(100).
        SetBackgroundColor(color.White).
        Add(
            components.NewText("Hello, Tenon!").
                SetFontSize(24),
        )
}

func main() {
    tenon.Run(&App{}, 800, 600)
}
```

---

## Widget 开发

### 有状态 Widget

```go
type Counter struct {
    core.Widget
    count *core.State[int]
}

func NewCounter() *Counter {
    return &Counter{count: core.NewState(0)}
}

func (c *Counter) Build() core.Element {
    text := components.NewText("").
        SetFontSize(18)
    
    // 属性绑定：State → Element
    c.count.Subscribe(func(v int) {
        text.SetContent(fmt.Sprintf("Count: %d", v))
    })
    
    btn := components.NewButton("+1").
        SetOnClick(func() {
            c.count.Set(c.count.Get() + 1)
        })
    
    return components.NewView().
        SetFlexDirection(yoga.FlexDirectionRow).
        SetGap(12).
        Add(text, btn)
}
```

### 条件渲染

条件变化时，Widget 需要告诉框架结构变了：

```go
type Toggle struct {
    core.Widget
    show *core.State[bool]
}

func (t *Toggle) Build() core.Element {
    root := components.NewView()
    
    // State 变化时请求结构重建
    t.show.Subscribe(func(v bool) {
        t.RequestBuild()  // 通知框架：我的结构可能变了
    })
    
    if t.show.Get() {
        root.Add(components.NewText("Visible"))
    } else {
        root.Add(components.NewText("Hidden"))
    }
    
    return root
}
```

> 注意：条件分支里产生的 Element 类型不同（或子节点数量不同）时，框架的同级 Diff 会自动复用/替换。

### 列表渲染

```go
type ItemList struct {
    core.Widget
    items *core.State[[]string]
}

func (l *ItemList) Build() core.Element {
    root := components.NewView().
        SetFlexDirection(yoga.FlexDirectionColumn)
    
    for _, item := range l.items.Get() {
        root.Add(components.NewText(item))
    }
    
    return root
}
```

列表数据变化时调用 `RequestBuild()`，框架的同级 Diff 会尽可能复用已有的 `*Text` 节点。

---

## Element 开发

内置组件（`View`、`Text`、`Button`…）都是 Element。你也可以创建自定义 Element。

```go
type MyBox struct {
    core.BaseElement
    color color.Color
}

func NewMyBox() *MyBox {
    b := &MyBox{color: color.White}
    b.Init()  // 初始化 BaseElement（创建 Yoga 节点等）
    return b
}

func (b *MyBox) Draw(screen *ebiten.Image) {
    bounds := b.Bounds()
    vector.FillRect(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, b.color, false)
}

func (b *MyBox) SetColor(c color.Color) *MyBox {
    b.color = c
    b.MarkNeedDraw()
    return b
}

// 类型标识，用于 Diff 时判断是否可复用
func (b *MyBox) ElementType() string { return "MyBox" }
```

---

## 项目结构

```
tenon/
├── pkg/
│   ├── core/                 # 框架核心
│   │   ├── widget.go         # Widget 接口、BaseWidget
│   │   ├── element.go        # Element 接口、BaseElement
│   │   ├── state.go          # State 信号系统
│   │   ├── engine.go         # Engine（Build/Diff/Layout/Draw/Event）
│   │   └── event.go          # 事件类型
│   ├── components/           # 内置 Element（View, Text, Button, Image...）
│   └── fonts/                # 字体管理
├── yoga/                     # Yoga 布局引擎（Go 移植）
├── internal/
│   └── renderer/             # Ebiten 渲染适配
├── example/
│   └── demo/                 # 示例应用
└── tenon.go                  # 公共 API 入口
```

---

## 与旧架构的区别

| 旧架构 | 新架构 |
|--------|--------|
| `Widget.Render()` | `Widget.Build()` |
| Host / Element 分离 | **Element 就是 Native 组件** |
| `UseState` 触发 `invalidate()` 重渲染 | `State` 直接通知订阅者，**不触发 Build** |
| `syncHost` / `replaceHostTree` 全量属性同步 | **同级浅 Diff 只做结构协调**，属性由 State 直连 |
| 组合组件（AntButton）是 Widget | 组合组件是 **Element 的包装器/工厂**，不是 Widget |

---

## License

MIT License
