# 三层树模型：Widget / Element / RenderNode

> 父文档：[ARCHITECTURE.md](./ARCHITECTURE.md)

---

## 为什么需要三层？

**组件粒度 ≠ 渲染粒度。**

用户写的组件可以任意嵌套：

```go
func (c *Pagination) Build() Widget {
    return Row(
        PrevButton(c.onPrev),
        ForEach(c.pages, func(p Page) Widget {
            return PageButton(p.number, p.active)
        }),
        NextButton(c.onNext),
    )
}
```

在这个例子中，`Pagination`、`PageButton`、`PrevButton` 都是用户组件，但 Yoga 和 Ebiten 只需要知道：
- 一个 `Row`（`FlexDirectionRow`）
- N 个文本/按钮叶子节点

中间的**用户组件不创建 RenderNode**，只创建轻量的 Element 来组织子树。

---

## 三层职责

### Widget Tree（Immutable 描述层）

- **职责**：描述 UI 应该长什么样
- **可变性**：Immutable（每次状态变化重建）
- **重量**：极轻量（只有配置数据，无 Yoga Node，无 Ebiten 资源）
- **生命周期**：每帧都可能重建

```go
type Widget interface {
    WidgetType() string
    Key() string
}

type ComponentWidget interface {
    Widget
    Build(ctx BuildContext) Widget
}

type RenderWidget interface {
    Widget
    CreateRenderNode() RenderNode
    UpdateRenderNode(node RenderNode)
    GetChildren() []Widget
}
```

### Element Tree（Mutable 身份 + 状态层）

- **职责**：决定复用或重建，持有状态，连接 Widget 和 RenderNode
- **可变性**：Mutable（跨重建复用）
- **重量**：中等（持有组件实例、子 Element 引用）
- **生命周期**：随挂载创建，随卸载销毁

```go
type Element interface {
    Mount(parent Element)
    Update(newWidget Widget)
    Unmount()
    Parent() Element
    Children() []Element
    Widget() Widget
    MarkNeedsBuild()
    BuildIfNeeded() bool
}
```

**两种 Element**：

| 类型 | 对应 Widget | 职责 |
|------|------------|------|
| `ComponentElement` | `ComponentWidget` | 管理子 Element，不直接渲染 |
| `RenderObjectElement` | `RenderWidget` | 连接 Yoga Node 和 Ebiten 绘制 |

### RenderNode Tree（Mutable 布局 + 绘制层）

- **职责**：Yoga 布局计算 + Ebiten 像素绘制
- **可变性**：Mutable
- **重量**：重量级（持有 Yoga Node、Ebiten 纹理/字体资源）
- **生命周期**：由 RenderObjectElement 管理，尽量复用

```go
type RenderNode interface {
    YogaNode() *yoga.Node
    SyncYogaProps(w RenderWidget)
    Paint(screen *ebiten.Image, bounds geometry.Rect)
    Parent() RenderNode
    Children() []RenderNode
}
```

---

## 三棵树的关系

```
Widget Tree (声明)          Element Tree (身份)         RenderNode Tree (像素)
─────────────────           ─────────────────           ───────────────────
Pagination(Build)           ComponentElement            (无 RenderNode)
    │                           │
    ▼                           ▼ Build()
Row(Gap(16),...)              RenderObjectElement         RowRenderNode(Yoga: Row)
    │                           │                               │
    ├── Text("Prev")            ├── RenderObjectElement         ├── TextRenderNode
    │                               (Text)                          (Yoga: Leaf)
    ├── PageButton(1)           ├── ComponentElement            │
    │   │                           │                               │
    │   ▼                           ▼ Build()                       │
    │   Text("1")                 RenderObjectElement             ├── TextRenderNode
    │                               (Text)                          (Yoga: Leaf)
    ├── PageButton(2)           ├── ComponentElement            │
    │   ▼                           ▼                               │
    │   Text("2")                 RenderObjectElement             ├── TextRenderNode
    │                               (Text)                          (Yoga: Leaf)
    └── Text("Next")            └── RenderObjectElement         └── TextRenderNode
                                    (Text)                          (Yoga: Leaf)
```

**关键观察**：
- `Pagination`、`PageButton` 等用户组件只有 Element，没有 RenderNode
- 只有 `Row`、`Text` 等渲染原语才创建 RenderNode
- Element 树深度 ≥ RenderNode 树深度

---

## 核心流程

### 首次挂载（Mount）

```
Engine.Mount(rootWidget)
  └── createElement(rootWidget) → ComponentElement
      └── ComponentElement.Mount()
          └── widget.Build(ctx) → RowWidget
              └── createElement(RowWidget) → RenderObjectElement
                  └── RenderObjectElement.Mount()
                      ├── renderNode = RowWidget.CreateRenderNode() // Yoga Node
                      ├── for each child Widget:
                      │     createElement(child) → Element
                      │     Element.Mount(this)
                      │     renderNode.AddChild(childRenderNode)
```

### 状态更新（Update）

```
c.SetState(func() { c.count++ })
  └── ComponentElement.MarkNeedsBuild()
      └── Engine.Update()
          └── ComponentElement.BuildIfNeeded()
              └── widget.Build(ctx) → newWidget
                  └── ReconcileChild(oldChild, newWidget)
                      ├── canUpdate(old, new) ?
                      │   ├── true → oldChild.Update(newWidget) → 复用
                      │   └── false → oldChild.Unmount() + createElement(new) → 重建
```

---

## 与 Flutter 的映射

| 本架构 | Flutter |
|--------|---------|
| Widget | Widget |
| ComponentElement | ComponentElement |
| RenderObjectElement | RenderObjectElement |
| RenderNode | RenderObject |
| BuildContext | BuildContext |
| `canUpdate` | `Widget.canUpdate` |
| `SetState()` | `setState()` |

---

## 与 React 的映射

| 本架构 | React |
|--------|-------|
| Widget | React Element |
| ComponentElement | Fiber (Class/Function Component) |
| RenderNode | DOM Node |
| `canUpdate` + Reconcile | Reconciliation |
| `SetState()` | `setState()` / `useState` dispatcher |
| Build() | render() / 函数组件 |
