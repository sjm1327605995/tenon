# Yoga + Ebiten 声明式 GUI 架构设计

> 目标：在 Go 中实现类似 React 的声明式组件 API，以 Yoga 为布局引擎、Ebiten 为渲染后端。
> 
> 核心约束：**Go 没有对象代理（Proxy）**，无法做 Vue 式的细粒度自动依赖追踪。

---

## 一、调研总结：React / Flutter / Vue Native

### 1.1 React Fiber

| 维度 | 说明 |
|------|------|
| **声明层** | JSX → React Element（immutable 描述对象，轻量） |
| ** reconciler** | Fiber 树 = mutable 的工作单元，持有 state、props、hooks、DOM 引用 |
| **Diff** | 双缓冲（current / workInProgress），单链表 DFS，O(n) diff |
| **Commit** | 副作用（增删改 DOM）同步一次性执行 |
| **状态触发** | `setState` → schedule → re-render → diff → commit |
| **渲染粒度** | Virtual DOM 节点与实际 DOM 节点**不强制一一对应**（Portal/Fragment） |

**对 Go 的启示**：
- 声明层用轻量 immutable 对象是可取的，Go 的 GC 适合管理短生命周期描述对象。
- 但 React 的 VDOM full-diff 在 Go 中性能开销大，且没有内置 key 优化机制。
- Hooks 依赖数组在编译期无检查，Go 更难模拟。

### 1.2 Flutter 三棵树

| 树 | 职责 | 可变性 | 生命周期 |
|----|------|--------|----------|
| **Widget** | 配置描述（Blueprint） | Immutable | 每次 setState 重建 |
| **Element** | 身份 + 状态 + 生命周期 | Mutable | 跨重建复用 |
| **RenderObject** | Layout + Paint + Hit-test | Mutable | 由 Element 管理 |

**核心机制 `canUpdate`**：
```dart
static bool canUpdate(Widget oldWidget, Widget newWidget) {
  return oldWidget.runtimeType == newWidget.runtimeType
      && oldWidget.key == newWidget.key;
}
```
- **相同类型 + 相同 Key** → 复用 Element 和 RenderObject，只更新属性
- **不同** → 卸载整个子树，重建

**对 Go 的启示**：
- Flutter 的"组件粒度"（Widget 可以任意嵌套）与"渲染粒度"（RenderObject 只包含实际要画的节点）**明确分离**。
- `canUpdate` 规则极其简单，避免了 VDOM 全量 diff，非常适合 Go 实现。
- Element 层作为"可变身份层"，完美解决了"声明式重建"与"状态保持"的矛盾。

### 1.3 Vue Native / Lynx

| 方案 | 状态 |
|------|------|
| Vue Native | 已废弃 |
| Lynx (Vue + Lynx) | 字节跳动新方案，双线程（主线程 UI + 后台线程 JS） |

**关键问题**：Vue 3 的响应式系统重度依赖 `Proxy`，Go 没有等价物。
- Vue 的 `ref/reactive` 在运行时自动追踪依赖。
- Lynx 编译时将 `<template>` 拆到主线程、`<script>` 拆到后台线程。

**结论**：Vue 的响应式模式在 Go 中**不可复制**，必须放弃。

---

## 二、Go 的核心约束与破局点

### 2.1 约束清单

| 约束 | 影响 |
|------|------|
| **无 Proxy** | 无法自动追踪状态依赖 → 无法 Vue 式响应式 |
| **无泛型变参** | 函数式 children API 需要 `any` 或生成代码 |
| **无 JSX** | 声明式 API 只能依赖函数调用嵌套 |
| **有 GC** | 每帧创建短生命周期的 Widget 描述对象是可行的 |
| **有 interface** | 非常适合建模多态的 Widget / Element / RenderNode |

### 2.2 破局思路

**采用 Flutter 的三层树模型 + React 的函数组件声明风格。**

- **声明式**：用户写函数返回 Widget 描述（像 React 函数组件）
- **高性能**：用 Flutter 的 `canUpdate` 复用规则，避免 VDOM 全量 diff
- **无代理**：状态更新显式调用 `SetState()`（像 React `useState` + `setState`）

---

## 三、推荐架构：三层树 + 显式状态

### 3.1 架构总览

```
┌─────────────────────────────────────────────────────────────┐
│  User Code (声明式)                                          │
│  func (c *Counter) Build(ctx BuildContext) Widget {         │
│      return Column(Gap(16),                                 │
│          Text(fmt.Sprintf("Count: %d", c.count)),           │
│          Button(OnPress(c.inc), Text("+")),                 │
│      )                                                      │
│  }                                                          │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼ Build()
┌─────────────────────────────────────────────────────────────┐
│  Widget Tree (Immutable 描述层)                              │
│  每次状态变化重建，极轻量（只有配置数据）                      │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼ Reconcile (canUpdate)
┌─────────────────────────────────────────────────────────────┐
│  Element Tree (Mutable 身份 + 状态层)                        │
│  跨重建复用，持有组件实例、Yoga Node 引用、state              │
│  ComponentElement: 管理子 Element                            │
│  RenderObjectElement: 连接 Yoga/Ebiten 渲染节点              │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼ Yoga CalculateLayout
┌─────────────────────────────────────────────────────────────┐
│  RenderNode Tree (Mutable 布局 + 绘制层)                     │
│  每个节点 = Yoga.Node + Ebiten 绘制命令                       │
│  Layout:  Yoga 约束计算 (O(n))                               │
│  Paint:   Ebiten DrawRect / DrawText / DrawImage            │
└─────────────────────────────────────────────────────────────┘
```

### 3.2 为什么需要三层？

**组件粒度 ≠ 渲染粒度。**

用户写的组件可以任意嵌套：
```go
// 用户组件，包含逻辑和子组件组合
func (c *Pagination) Build() Widget {
    return Row(
        PrevButton(c.onPrev),
        // ... 10 个 PageButton ...
        NextButton(c.onNext),
    )
}
```

但 Yoga 和 Ebiten 只需要知道：
- 一个 Row（FlexDirectionRow）
- 12 个子节点（PrevButton 对应一个 Text RenderNode，PageButton 对应 Text RenderNode...）

中间的 `Pagination`、`PageButton` 等用户组件**不创建 RenderNode**，只创建 Element 来组织子树。

---

## 四、核心接口设计

### 4.1 Widget（声明层）

```go
// Widget 是所有 UI 描述的根接口，immutable
// 每次重建都是新实例，但框架通过 Type+Key 判断是否复用 Element
type Widget interface {
    // WidgetType 用于 canUpdate 判断（可用 reflect.Type 或字符串）
    WidgetType() string
    // Key 可选，用于列表等场景稳定身份
    Key() string
}

// ComponentWidget = 用户自定义组件，有 Build() 方法
type ComponentWidget interface {
    Widget
    Build(ctx BuildContext) Widget
}

// RenderWidget = 会生成实际渲染节点的组件（Text, Box, Image, Row, Column...）
type RenderWidget interface {
    Widget
    CreateRenderNode() RenderNode
}
```

### 4.2 Element（身份层）

```go
// Element 是 mutable 的，跨 Widget 重建复用
type Element interface {
    // 挂载时调用，创建子 Element 或 RenderNode
    Mount(parent Element)
    // 更新时调用，传入新 Widget，决定复用或重建
    Update(newWidget Widget)
    // 卸载时释放资源（Yoga Node、Ebiten 资源）
    Unmount()
    
    Parent() Element
    Children() []Element
    Widget() Widget
}

// ComponentElement: 用户组件的 Element，持有组件实例
type ComponentElement struct {
    widget   ComponentWidget
    children []Element
    parent   Element
    // 组件实例（用户自定义 struct）
    state    Component
}

// RenderObjectElement: 连接 RenderWidget 和 RenderNode
type RenderObjectElement struct {
    widget     RenderWidget
    renderNode RenderNode
    children   []Element  // 只有容器类 RenderWidget 才有子 Element
    parent     Element
}
```

### 4.3 RenderNode（渲染层）

```go
// RenderNode = Yoga Node + Ebiten 绘制
type RenderNode interface {
    // Yoga 相关
    YogaNode() *yoga.Node
    
    // 从 Widget 配置同步到 Yoga Node（样式设置）
    SyncYogaProps(w RenderWidget)
    
    // 绘制到 Ebiten screen
    Paint(screen *ebiten.Image, bounds geometry.Rect)
    
    // 布局回调：Yoga MeasureFunc 调用，计算内容尺寸
    Measure(width, height float32) geometry.Size
    
    Parent() RenderNode
    Children() []RenderNode
}

// 具体实现示例：文本节点
type TextRenderNode struct {
    yogaNode *yoga.Node
    text     string
    style    TextStyle
    // Ebiten 预渲染文本...
}
```

### 4.4 BuildContext

```go
// BuildContext 在 Build() 时传入，提供上下文能力
type BuildContext interface {
    // 获取祖先节点（用于 Theme、依赖注入等）
    Ancestor(widgetType string) Element
    // 标记当前组件需要重建
    MarkNeedsBuild()
    // 获取 RenderNode（仅 RenderObjectElement 有）
    RenderNode() RenderNode
}
```

---

## 五、用户层 API 设计（类似 React）

### 5.1 函数式 Widget 声明

Go 没有 JSX，用函数调用嵌套模拟：

```go
// 用户组件 = struct + Build 方法
type Counter struct {
    count int
}

func (c *Counter) Build(ctx BuildContext) Widget {
    return Column(
        Align(AlignCenter),
        Gap(16),
        Text(fmt.Sprintf("Count: %d", c.count), FontSize(24)),
        Button(
            OnPress(c.inc),
            Background(ColorBlue),
            Text("Increment"),
        ),
    )
}

func (c *Counter) inc() {
    c.count++
    // 显式标记重建，因为没有 Proxy 自动追踪
    c.SetState()
}
```

### 5.2 无状态组件（纯函数）

```go
func Greeting(name string) Widget {
    return Text(fmt.Sprintf("Hello, %s!", name))
}

// 使用
Column(Text("Title"), Greeting("Alice"))
```

### 5.3 列表与 Key

```go
func (c *TodoList) Build(ctx BuildContext) Widget {
    var items []Widget
    for _, todo := range c.todos {
        // Key 用于稳定 Element 身份，避免复用混乱
        items = append(items, TodoItem(
            Key(todo.ID),
            todo,
        ))
    }
    return Column(Gap(8), items...)
}
```

### 5.4 与现有 tenon signal 系统的兼容

当前 tenon 已有 `state.Signal`。可以保留 Signal，但**不自动追踪**——用户显式绑定：

```go
type Counter struct {
    count *state.Signal[int]
}

func (c *Counter) Build(ctx BuildContext) Widget {
    // Signal 读值，但不自动订阅
    val := c.count.Get()
    return Button(OnPress(func() {
        c.count.Set(val + 1)
        c.SetState() // 仍然需要显式触发重建
    }), Text(fmt.Sprintf("%d", val)))
}
```

> 未来可扩展：在 `BuildContext` 中自动收集 `Signal.Get()` 调用，隐式建立订阅关系（类似 SolidJS 的追踪），但这需要运行时 hook，增加复杂度。建议第一步保持显式 `SetState()`。

---

## 六、Reconciliation 算法（简化版）

### 6.1 canUpdate 规则

```go
func canUpdate(oldW, newW Widget) bool {
    if oldW == nil || newW == nil {
        return false
    }
    return oldW.WidgetType() == newW.WidgetType() && oldW.Key() == newW.Key()
}
```

### 6.2 单个子节点 reconcile

```go
func reconcileChild(parent Element, oldChild Element, newWidget Widget) Element {
    if oldChild != nil && canUpdate(oldChild.Widget(), newWidget) {
        oldChild.Update(newWidget)
        return oldChild
    }
    if oldChild != nil {
        oldChild.Unmount()
    }
    newChild := createElement(newWidget)
    newChild.Mount(parent)
    return newChild
}
```

### 6.3 子节点列表 reconcile（带 Key）

类似 React 的 Key-based diff，但简化：

```go
func reconcileChildren(parent Element, oldChildren []Element, newWidgets []Widget) []Element {
    // 1. 按 Key 建立旧节点 map
    oldKeyed := make(map[string]Element)
    var oldUnkeyed []Element
    for _, c := range oldChildren {
        if k := c.Widget().Key(); k != "" {
            oldKeyed[k] = c
        } else {
            oldUnkeyed = append(oldUnkeyed, c)
        }
    }
    
    // 2. 遍历新 widgets，匹配或创建
    var newChildren []Element
    unkeyedIdx := 0
    for _, w := range newWidgets {
        var matched Element
        if k := w.Key(); k != "" {
            matched = oldKeyed[k]
            delete(oldKeyed, k)
        } else if unkeyedIdx < len(oldUnkeyed) {
            if canUpdate(oldUnkeyed[unkeyedIdx].Widget(), w) {
                matched = oldUnkeyed[unkeyedIdx]
            }
            unkeyedIdx++
        }
        
        if matched != nil {
            matched.Update(w)
        } else {
            matched = createElement(w)
            matched.Mount(parent)
        }
        newChildren = append(newChildren, matched)
    }
    
    // 3. 卸载未复用的旧节点
    for _, c := range oldUnkeyed[unkeyedIdx:] {
        c.Unmount()
    }
    for _, c := range oldKeyed {
        c.Unmount()
    }
    
    return newChildren
}
```

---

## 七、Yoga + Ebiten 集成

### 7.1 渲染管线

```
Ebiten Update():
  1. 处理输入事件（鼠标、键盘）
  2. 如果有 SetState 触发脏标记：
     a. Rebuild: 从脏节点开始调用 Build() → 生成新 Widget 子树
     b. Reconcile: diff 旧 Widget / 新 Widget → 更新 Element 树
     c. Yoga Layout: 自顶向下传递约束 → yoga.CalculateLayout → 自底向上确定 bounds
  3. 如果需要，触发动画更新

Ebiten Draw(screen):
  1. 遍历 RenderNode 树（DFS）
  2. 对每个 RenderNode 调用 Paint(screen, node.Bounds())
     - Box: ebitenvector.DrawFilledRect
     - Text: text/v2 Draw
     - Image: screen.DrawImage
  3. 处理 clip、opacity、transform
```

### 7.2 Yoga Node 生命周期

- **创建**：`RenderObjectElement.Mount()` 时，调用 `widget.CreateRenderNode()` → `yoga.NewNode()`
- **配置**：`RenderNode.SyncYogaProps()` 从 Widget 读取 flex 属性设置到 Yoga Node
- **布局**：`yoga.CalculateLayout(rootWidth, rootHeight, DirectionLTR)`
- **读取**：`node.GetLayoutLeft/Top/Width/Height` → 设置到 `RenderNode.Bounds`
- **销毁**：`RenderObjectElement.Unmount()` → `yoga.Node.Free()`

> **关键**：Yoga Node 应该** retained**（跨帧复用），而不是像当前 tenon 那样每帧新建。这是性能核心。

### 7.3 容器节点（Row / Column）

```go
type RowRenderNode struct {
    yogaNode *yoga.Node
    children []RenderNode
}

func (r *RowRenderNode) SyncYogaProps(w RenderWidget) {
    rw := w.(*RowWidget)
    r.yogaNode.SetFlexDirection(yoga.FlexDirectionRow)
    r.yogaNode.SetGap(yoga.GutterAll, rw.gap)
    r.yogaNode.SetPadding(yoga.EdgeAll, rw.padding)
    // ...
}

func (r *RowRenderNode) Paint(screen *ebiten.Image, bounds geometry.Rect) {
    // Row 本身通常不绘制背景（由 Box 处理）
    // 子节点的 Paint 在 RenderTree 遍历中自动调用
}
```

### 7.4 文本测量（Yoga MeasureFunc）

```go
func (t *TextRenderNode) SyncYogaProps(w RenderWidget) {
    t.yogaNode.SetMeasureFunc(func(node *yoga.Node, width float32, widthMode yoga.MeasureMode, height float32, heightMode yoga.MeasureMode) yoga.Size {
        // 用 Ebiten text/v2 测量文本尺寸
        sz := measureText(t.text, t.style, width, widthMode)
        return yoga.Size{Width: sz.Width, Height: sz.Height}
    })
}
```

---

## 八、与当前 tenon 架构的对比

| 维度 | 当前 tenon | 新架构 |
|------|-----------|--------|
| **组件模型** | 单层 Widget（命令式） | 三层树（Widget/Element/RenderNode） |
| **API 风格** | 命令式：`w.SetText("x")` | 声明式：`Build() Widget` |
| **状态更新** | Signal push + 直接 Mark | 显式 `SetState()` + Rebuild |
| **Yoga 使用** | 即时模式（每帧新建 Node） | 保留模式（Node 随 Element 复用） |
| **渲染粒度** | 每个 Widget 都 Layout+Draw | 用户组件不创建 RenderNode，只有叶子/容器创建 |
| **Ebiten 集成** | `App` 直接驱动 | `Game` 实现分 Update(Build+Layout) / Draw |

---

## 九、实现路径建议

### Phase 1：核心骨架（验证可行性）
1. 定义 `Widget`, `Element`, `RenderNode`, `BuildContext` 接口
2. 实现 `ComponentElement` 和 `RenderObjectElement`
3. 实现基础 Reconcile（单节点 + 列表）
4. 绑定 Ebiten `Game`：`Update` 做 Build+Reconcile+Layout，`Draw` 遍历 RenderNode 画矩形

### Phase 2：基础 Widget 库
1. `Box`（背景色、圆角、边框）
2. `Text`（Ebiten text/v2）
3. `Row`, `Column`（Yoga Flex）
4. `Image`
5. `Button`（点击区域 + 事件）

### Phase 3：交互与状态
1. 事件系统（Hit-test 遍历 RenderNode bounds）
2. `GestureDetector`（Tap, Pan, Scroll）
3. 与现有 `state.Signal` 集成（可选自动追踪）

### Phase 4：高级特性
1. `Key` 优化列表 reconcile
2. `Fragment` / `Portal`
3. 动画集成（当前 tenon 的 animation 包）
4. 主题 / 依赖注入（InheritedWidget 模式）

---

## 十、核心代码骨架

```go
// ============================================================
// 1. Widget 层
// ============================================================

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
}

// ============================================================
// 2. Element 层
// ============================================================

type Element interface {
    Mount(parent Element)
    Update(newWidget Widget)
    Unmount()
    Parent() Element
    Children() []Element
    Widget() Widget
    MarkNeedsBuild()
}

type ComponentElement struct {
    parent   Element
    widget   ComponentWidget
    child    Element   // Component 只有一个 child（Build 返回单根）
    state    Component // 用户组件实例
    dirty    bool
}

func (e *ComponentElement) Update(newWidget Widget) {
    e.widget = newWidget.(ComponentWidget)
    e.rebuild()
}

func (e *ComponentElement) rebuild() {
    newWidget := e.widget.Build(e)
    e.child = reconcileChild(e, e.child, newWidget)
}

type RenderObjectElement struct {
    parent     Element
    widget     RenderWidget
    renderNode RenderNode
    children   []Element
}

func (e *RenderObjectElement) Mount(parent Element) {
    e.parent = parent
    e.renderNode = e.widget.CreateRenderNode()
    e.renderNode.SyncYogaProps(e.widget)
    // 挂载子节点...
}

// ============================================================
// 3. RenderNode 层
// ============================================================

type RenderNode interface {
    YogaNode() *yoga.Node
    SyncYogaProps(w RenderWidget)
    Paint(screen *ebiten.Image, bounds geometry.Rect)
    Parent() RenderNode
    Children() []RenderNode
}

// ============================================================
// 4. Ebiten Game 驱动
// ============================================================

type Engine struct {
    rootWidget    Widget
    rootElement   Element
    rootRenderNode RenderNode
    needRebuild   bool
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
        X: node.YogaNode().GetLayoutLeft(),
        Y: node.YogaNode().GetLayoutTop(),
        W: node.YogaNode().GetLayoutWidth(),
        H: node.YogaNode().GetLayoutHeight(),
    }
    node.Paint(screen, bounds)
    for _, child := range node.Children() {
        g.paintTree(screen, child)
    }
}
```

---

## 十一、结论

**不要试图在 Go 中复制 Vue 的 Proxy 响应式。** 那是死胡同。

**最佳路径是：Flutter 的三层树架构 + React 的函数组件声明语法。**

- Widget 层给用户一个干净、声明式的 API（函数嵌套模拟 JSX）。
- Element 层解决"重建 Widget 但保留状态和 Yoga Node"的核心矛盾。
- RenderNode 层把 Yoga 和 Ebiten 粘合在一起，只做布局+绘制。
- 状态更新用显式 `SetState()`，简单、可预测、零魔法。
