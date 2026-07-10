# Yoga + Ebiten 集成

> 父文档：[ARCHITECTURE.md](./ARCHITECTURE.md)
> 前置阅读：[03-three-tree-model.md](./03-three-tree-model.md), [05-reconciliation.md](./05-reconciliation.md)

---

## Yoga Node 生命周期

与当前 tenon"即时模式"（每帧新建 Yoga Node）不同，新架构采用**保留模式**：

| 阶段 | 操作 | 触发时机 |
|------|------|----------|
| **创建** | `yoga.NewNode()` | `RenderObjectElement.Mount()` |
| **配置** | `StyleSetXxx()` | `RenderNode.SyncYogaProps()` |
| **插入子节点** | `InsertChild()` | Reconcile 新增子节点时 |
| **移除子节点** | `RemoveChild()` | Reconcile 卸载子节点时 |
| **布局** | `CalculateLayout()` | Engine.Update() 中统一调用 |
| **读取结果** | `LayoutLeft/Top/Width/Height` | 布局后同步到 RenderNode.Bounds |
| **销毁** | `Reset()` | `RenderObjectElement.Unmount()` |

> **关键**：Yoga Node 随 Element 复用，不是每帧新建。这是性能核心。

---

## 渲染管线

```
Ebiten Update():
  1. 处理输入事件（鼠标、键盘）
  2. 如果有 SetState 触发脏标记：
     a. Rebuild: 从脏节点开始调用 Build() → 生成新 Widget 子树
     b. Reconcile: diff 旧 Widget / 新 Widget → 更新 Element 树
     c. Yoga Layout: 自顶向下传递约束 → yoga.CalculateLayout → 自底向上确定 bounds
  3. 触发动画更新（如有）

Ebiten Draw(screen):
  1. 遍历 RenderNode 树（DFS）
  2. 对每个 RenderNode 调用 Paint(screen, node.Bounds())
     - Box: ebitenvector.DrawFilledRect
     - Text: text/v2 Draw
     - Image: screen.DrawImage
  3. 处理 clip、opacity、transform
```

---

## RenderNode 与 Yoga 的绑定

```go
type RenderNode interface {
    YogaNode() *yoga.Node
    SyncYogaProps(w RenderWidget)
    Paint(screen *ebiten.Image, bounds geometry.Rect)
    Parent() RenderNode
    Children() []RenderNode
}
```

### 容器节点（Row / Column）

```go
type RowRenderNode struct {
    yogaNode *yoga.Node
    children []RenderNode
    bounds   geometry.Rect
}

func (r *RowRenderNode) SyncYogaProps(w RenderWidget) {
    rw := w.(*RowWidget)
    r.yogaNode.StyleSetFlexDirection(yoga.FlexDirectionRow)
    r.yogaNode.StyleSetGap(yoga.GutterAll, rw.gap)
    r.yogaNode.StyleSetPadding(yoga.EdgeAll, rw.padding)
    // ...
}

func (r *RowRenderNode) Paint(screen *ebiten.Image, bounds geometry.Rect) {
    // Row 本身通常不绘制背景（由 Box 处理）
    // 子节点的 Paint 在 RenderTree 遍历中自动调用
}
```

### 文本测量（Yoga MeasureFunc）

```go
func (t *TextRenderNode) SyncYogaProps(w RenderWidget) {
    t.yogaNode.SetMeasureFunc(func(
        node *yoga.Node,
        width float32, widthMode yoga.MeasureMode,
        height float32, heightMode yoga.MeasureMode,
    ) yoga.Size {
        // 用 Ebiten text/v2 测量文本尺寸
        sz := measureText(t.text, t.style, width, widthMode)
        return yoga.Size{Width: sz.Width, Height: sz.Height}
    })
}
```

---

## Ebiten Game 驱动

```go
type Engine struct {
    rootWidget     Widget
    rootElement    Element
    rootRenderNode RenderNode
    needRebuild    bool
    width, height  int
}

func (g *Engine) Update() error {
    if g.needRebuild {
        g.reconcile()
        g.layout()
        g.needRebuild = false
    }
    return nil
}

func (g *Engine) Draw(screen *ebiten.Image) {
    g.paintTree(screen, g.rootRenderNode)
}

func (g *Engine) Layout(outsideWidth, outsideHeight int) (int, int) {
    g.width, g.height = outsideWidth, outsideHeight
    return outsideWidth, outsideHeight
}

func (g *Engine) reconcile() {
    if g.rootElement == nil {
        g.rootElement = createElement(g.rootWidget)
        g.rootElement.Mount(nil)
        g.rootRenderNode = findRootRenderNode(g.rootElement)
    } else {
        g.rootElement.Update(g.rootWidget)
    }
}

func (g *Engine) layout() {
    // Yoga 自顶向下约束，自底向上尺寸
    g.rootRenderNode.YogaNode().CalculateLayout(
        float32(g.width), float32(g.height), yoga.DirectionLTR,
    )
    // 遍历 RenderNode 树，把 Yoga 计算结果写入 Bounds
    syncBoundsFromYoga(g.rootRenderNode)
}

func (g *Engine) paintTree(screen *ebiten.Image, node RenderNode) {
    bounds := geometry.Rect{
        X: node.YogaNode().LayoutLeft(),
        Y: node.YogaNode().LayoutTop(),
        W: node.YogaNode().LayoutWidth(),
        H: node.YogaNode().LayoutHeight(),
    }
    node.Paint(screen, bounds)
    for _, child := range node.Children() {
        g.paintTree(screen, child)
    }
}
```

---

## 与当前 tenon 的差异

| 维度 | 当前 tenon | 新架构 |
|------|-----------|--------|
| Yoga Node 生命周期 | 即时模式（每帧新建） | 保留模式（随 Element 复用） |
| 布局触发 | 每帧 `Layout()` | 仅在脏时 `CalculateLayout()` |
| Ebiten 驱动 | `App` 直接驱动 Widget | `Engine` 驱动三层树 |
| 绘制遍历 | Widget Tree | RenderNode Tree（更浅） |
