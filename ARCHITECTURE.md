# Tenon v2 架构设计与开发规范

## 一、项目结构

```
tenon/
├── tenon.go                 # 统一入口，导出 v2 类型和 Run()
├── pkg/
│   ├── v2/
│   │   ├── core/            # 引擎、Element、Widget、事件、动画、主题
│   │   │   ├── engine.go      # 引擎：帧循环、布局、绘制、事件分发
│   │   │   ├── element.go     # Element 接口 + BaseElement 默认实现
│   │   │   ├── widget.go      # Widget 接口 + BaseWidget 默认实现
│   │   │   ├── event.go       # 事件类型与分发
│   │   │   ├── animation.go   # Tween 动画
│   │   │   ├── transform.go   # 2D 仿射变换工具
│   │   │   ├── theme.go       # 主题系统
│   │   │   └── styles.go      # 全局样式注册
│   │   └── components/      # 通用 UI 组件
│   │       ├── view.go          # 基础容器（背景、边框、阴影、圆角）
│   │       ├── text.go          # 文本
│   │       ├── shadow_text.go   # 阴影文字
│   │       ├── image.go         # 图片
│   │       ├── button.go        # 按钮
│   │       ├── scroll_view.go   # 滚动容器
│   │       ├── list_view.go     # 列表
│   │       ├── dropdown.go      # 下拉选择
│   │       ├── modal.go         # 模态对话框
│   │       ├── window.go        # 浮动窗口
│   │       ├── menu.go          # 菜单
│   │       ├── checkbox.go      # 复选框
│   │       ├── radio.go         # 单选
│   │       ├── switch.go        # 开关
│   │       ├── slider.go        # 滑块
│   │       ├── text_input.go    # 输入框
│   │       ├── progress_bar.go  # 进度条
│   │       ├── badge.go         # 徽标
│   │       ├── tooltip.go       # 提示
│   │       ├── divider.go       # 分隔线
│   │       ├── loading_spinner.go
│   │       ├── svg_icon.go
│   │       └── draw_utils.go    # 绘制工具（圆角矩形、圆形）
│   ├── fonts/               # 字体管理器（加载、回退、缓存）
│   └── svg/                 # SVG 路径解析与渲染
├── yoga/                    # Yoga Flexbox 引擎 Go 绑定（独立模块）
├── example/
│   ├── v2-demo/             # 框架功能演示
│   └── card/                # 游戏卡牌扩展示例（非框架部分）
└── docs/                    # 文档
```

**清理原则**：
- 框架核心只保留通用 UI 组件
- 游戏/业务专用组件（如 Card）放到 `example/` 作为扩展示例
- 空目录、日志文件、编译产物定期清理

---

## 二、架构总览：三层模型

Tenon 采用 **声明式 UI + 增量更新** 架构，分为三层：

```
┌─────────────────────────────────────────┐
│  Layer 1: Widget 层（声明式配置）        │
│  - 轻量级、可频繁重建                     │
│  - 只描述「应该长什么样」                 │
│  - 状态变化 → RequestBuild() → 重建      │
├─────────────────────────────────────────┤
│  Layer 2: Element 层（持久化渲染节点）    │
│  - 重量级、长期存活                       │
│  - 持有 Yoga 节点、Bounds、脏标记         │
│  - 负责 Draw、HandleEvent、Update        │
├─────────────────────────────────────────┤
│  Layer 3: Yoga 层（布局计算引擎）         │
│  - 纯布局计算，无渲染逻辑                 │
│  - 输入样式 → 输出 left/top/width/height │
│  - 通过 MeasureFunc 与 Element 双向交互  │
└─────────────────────────────────────────┘
```

---

## 三、三棵树的关系

Tenon 运行时会同时维护 **三棵逻辑树**：

### 树 1：Widget 树（声明树）

```go
type Widget interface {
    Build() Element      // 描述 Element 应该长什么样
    RequestBuild()       // 请求下一帧重建
    OnMount(engine *Engine)
    OnUnmount()
}
```

- **生命周期短**：状态一变就可以丢弃重建
- **无 Yoga 节点**：不直接参与布局
- **示例**：`type MyScreen struct { BaseWidget; count int }`

### 树 2：Element 树（渲染树）

```go
type Element interface {
    Draw(screen *ebiten.Image)
    HandleEvent(e *Event) bool
    Update() error
    GetYoga() *yoga.Node
    GetBounds() LayoutBounds
    // ... tree ops, flags, engine, context
}
```

- **生命周期长**：首次 Mount 后持久化，通过 `patchElement` 增量更新
- **持有 Yoga 节点**：每个 Element 对应一个 Yoga Node（`BaseElement.yoga`）
- **持有屏幕坐标**：`LayoutBounds{X,Y,Width,Height}`
- **脏标记驱动**：`FlagNeedMeasure | FlagNeedLayout | FlagNeedDraw`

### 树 3：Yoga 节点树（布局树）

```go
type Node struct { /* C++ Yoga 节点的 Go 包装 */ }
```

- **纯计算树**：只负责根据 Flex 属性计算每个节点的位置尺寸
- **挂载在 Element 上**：`Element.GetYoga()` 返回对应的 Yoga Node
- **MeasureFunc 回调**：Yoga 在需要测量内容尺寸时调用 Element 的 `MeasureFunc`

### 三棵树的映射关系

```
Widget 树（声明式）          Element 树（持久化）           Yoga 节点树（计算）
                              
MyScreen (Widget)  ──Build()──►  root View (Element)  ──1:1──►  Yoga Node
    │                              │    │                         │
    │ Build()                      │    ├─ Button (Element)  ────► Yoga Node
    ▼                              │    ├─ Text (Element)  ─────► Yoga Node
  Button (Widget)  ──Build()────►  │    └─ ScrollView (Element) ► Yoga Node
                                   │         └─ content View      ► Yoga Node
                                   │              ├─ Text          ► Yoga Node
                                   │              └─ Text          ► Yoga Node
```

**关键规则**：
1. Widget 通过 `Build()` 产出 Element 树，但 Widget 本身不直接挂载到树上
2. Element 树首次构建后持久化，后续 Widget 重建时通过 `patchElement` 做同级对比复用
3. Yoga 节点随 Element 的 `AppendChild`/`RemoveChild` 同步增删
4. Yoga 计算完成后，引擎通过 `updateBounds` 把 Yoga 的 `LayoutLeft/Top/Width/Height` 写回 Element 的 `LayoutBounds`

---

## 四、核心数据流

### 4.1 首次挂载（Mount）

```
Widget.OnMount(engine)
    ↓
Widget.Build()  →  生成新 Element 树
    ↓
递归 onElementMounted(el)
    ├── 注入 Engine
    ├── 应用全局样式（styles.go）
    └── 调用 el.OnMount(engine)
    ↓
Yoga CalculateLayout → 计算所有节点位置尺寸
    ↓
updateBounds(root, 0, 0) → 把 Yoga 结果写回 Element.Bounds
```

### 4.2 状态变化重建（Rebuild）

```
Widget 状态变化
    ↓
widget.RequestBuild()
    ↓
下一帧 flushBuildQueue
    ├── Widget.Build() 生成新 Element 树
    ├── patchElement(oldRoot, newRoot) 同级对比复用
    │   ├── 同类型：复用旧节点，CopyStyleFrom，递归 patchChildren
    │   └── 不同类型：替换节点，RemoveChild + AppendChild
    └── 标记 FlagNeedLayout | FlagNeedDraw
    ↓
flushDirtyElements
    ├── 测量（FlagNeedMeasure）
    ├── 布局（FlagNeedLayout → calculateLayout）
    └── 清除脏标记
```

### 4.3 帧循环（Update + Draw）

```
每帧 Update():
    1. flushBuildQueue      ← Widget 重建
    2. handleEvents          ← 鼠标/键盘/滚轮事件
    3. flushDirtyElements    ← 测量 → 布局 → 清标记
    4. updateAnimations      ← Tween 更新
    5. updateElements        ← 每帧 Element.Update()

每帧 Draw():
    screen.Fill(背景色)
    drawElement(root)
        ├── el.Draw(screen)        ← 组件自绘制
        ├── FlagClipChildren ?     ← SubImage 裁剪
        └── 递归 drawElement(子元素)
```

### 4.4 事件分发

```
鼠标/键盘输入
    ↓
hitTest(root, x, y) → 后序遍历，子元素优先
    ↓
dispatchEvent(target, event)
    ├── target.HandleEvent(event)
    ├── 若返回 true → 停止冒泡
    └── 若返回 false → parent.HandleEvent(event)
        └── 一直冒泡到根
```

---

## 五、脏标记系统

`ElementFlags` 是 `uint64` 位图：

```
低 32 位：持久状态（不会被 ClearDirty 清除）
    FlagVisible      → 是否可见
    FlagFocusable    → 是否可聚焦
    FlagClipChildren → 是否裁剪子元素

高 32 位：脏标记（flushDirtyElements 后清除）
    FlagNeedMeasure  → 需要重新测量（文字排版）
    FlagNeedLayout   → 需要重新布局（样式变化）
    FlagNeedDraw     → 需要重绘（视觉属性变化）
```

**使用规范**：
- 修改布局属性（width/height/margin/flex）→ `Mark(FlagNeedLayout)`
- 修改视觉属性（color/text/radius）→ `Mark(FlagNeedDraw)`
- 修改内容（文字内容、图片源）→ `Mark(FlagNeedMeasure | FlagNeedLayout | FlagNeedDraw)`
- 不要直接修改 `bounds`，它由 Yoga 计算后通过 `updateBounds` 写入

---

## 六、组件开发规范

### 6.1 目录与包结构

- 所有通用组件放在 `pkg/v2/components/`
- 每个组件一个文件，命名：`snake_case.go`
- 测试文件：`xxx_test.go`
- 游戏/业务专用组件不要放在框架中，在外部项目 import tenon 后自行扩展

### 6.2 组件结构模板

所有组件遵循以下模板：

```go
package components

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/sjm1327605995/tenon/pkg/v2/core"
)

// ComponentName 一句话描述组件用途。
type ComponentName struct {
    core.BaseElement
    // 状态字段（视觉状态、数据）
    // 不要放 Yoga 节点（BaseElement 已有）
    // 不要放 Engine 引用（BaseElement 已有）
}

// NewComponentName 创建组件，必须在构造函数中调用 Init。
func NewComponentName() *ComponentName {
    c := &ComponentName{}
    c.Init(c)                    // ← 必须！绑定 self + 创建 Yoga 节点
    // 设置默认 Yoga 样式
    // 设置 MeasureFunc（如果需要）
    return c
}

// ElementType 返回类型标识符，用于 patchElement 同级对比复用。
func (c *ComponentName) ElementType() string { return "ComponentName" }

// Draw 渲染组件。使用 screen 坐标（从 GetBounds 读取）。
func (c *ComponentName) Draw(screen *ebiten.Image) {
    if !c.IsVisible() { return }
    bounds := c.GetBounds()
    if bounds.Width <= 0 || bounds.Height <= 0 { return }
    
    // 1. 获取 Transform
    tr := c.GetTransform()
    
    // 2. 绘制逻辑
    //    - 使用 vector 包绘制几何图形
    //    - 使用 ebiten.DrawImage 绘制图片
    //    - 使用 text/v2 绘制文字
    
    // 3. 若需要 Transform，用 core.BuildTransformGeoM(bounds, tr)
    //    或绘制到临时 image 再变换（View 的做法）
}

// HandleEvent 处理事件。返回 true 表示消费事件，停止冒泡。
func (c *ComponentName) HandleEvent(e *core.Event) bool {
    // switch e.Type { ... }
    return false
}

// Update 每帧调用，处理动画、hover 检测等。
func (c *ComponentName) Update() error {
    return nil
}

// ==================== Chain API ====================
// 
// 所有修改状态的 setter 必须：
// 1. 检查值是否真的变化（避免无效重绘）
// 2. 标记正确的脏标记
// 3. 返回 *ComponentName 支持链式调用

func (c *ComponentName) SetSomeProperty(v SomeType) *ComponentName {
    if c.someField == v { return c }
    c.someField = v
    c.Mark(core.FlagNeedDraw)   // 或 FlagNeedLayout / FlagNeedMeasure
    return c
}
```

### 6.3 关键原则

#### 原则 1：构造函数必须调用 `Init(self)`

```go
func NewXxx() *Xxx {
    x := &Xxx{}
    x.Init(x)  // ← 必须
    return x
}
```

`Init` 会：
- 绑定 `self`（让 BaseElement 知道外部类型是谁）
- 创建 Yoga 节点
- 设置默认 `FlagVisible`

#### 原则 2：所有子元素在构造函数中创建

**禁止**在 setter 中懒创建子元素。Yoga `RemoveAllChildren` 在空节点上会 panic。

```go
// ✅ 正确：所有 child 在 NewXxx() 中创建
func NewModal() *Modal {
    m := &Modal{}
    m.Init(m)
    m.panel = NewView()   // ← 构造函数中创建
    m.titleEl = NewText("")  // ← 构造函数中创建
    m.panel.Add(m.titleEl)
    m.Add(m.panel)
    return m
}

// ❌ 错误：在 SetTitle 中创建
func (m *Modal) SetTitle(t string) {
    if m.titleEl == nil {  // ← 不要这样做
        m.titleEl = NewText(t)
    }
}
```

#### 原则 3：链式 API 返回具体类型，但 Element 接口方法返回 `Element`

```go
// Element 接口的方法必须返回 Element
func (b *BaseElement) SetWidth(v float32) core.Element { ... }

// 组件特有的方法返回具体类型
func (v *View) SetBackgroundColor(c color.Color) *View { ... }
```

注意：如果组件内嵌 `BaseElement`，不能重写 `SetWidth` 等接口方法来返回具体类型，因为 Go 的接口实现要求方法签名完全一致。

#### 原则 4：绘制使用 Transform

如果组件使用 `ebiten.DrawImage` 或 `text.Draw`：

```go
op := &ebiten.DrawImageOptions{}
op.GeoM.Concat(core.BuildTransformGeoM(bounds, c.GetTransform()))
core.ApplyColorScaleAlpha(&op.ColorScale, c.GetTransform().Alpha)
screen.DrawImage(img, op)
```

如果组件使用 `vector` 包绘制几何图形，参考 `View.Draw` 的做法：当 Transform 非单位矩阵时，先绘制到临时 `ebiten.Image`，再用 `DrawImage` 应用变换。

#### 原则 5：事件消费要明确

```go
func (c *ComponentName) HandleEvent(e *core.Event) bool {
    switch e.Type {
    case core.EventClick:
        if c.onClick != nil {
            c.onClick()
        }
        return true  // ← 消费了事件，停止冒泡
    case core.EventMouseMove:
        // 只更新 hover 状态，不阻止其他组件接收
        return false
    }
    return false
}
```

#### 原则 6：不要直接操作 Yoga 节点做树操作

Element 的 `AppendChild`/`RemoveChild`/`ClearChildren` 已经同步了 Yoga 节点树：

```go
// ✅ 正确：通过 Element API
parent.AppendChild(child)
parent.RemoveChild(child)

// ❌ 错误：直接操作 Yoga（除非你知道在做什么）
parent.GetYoga().InsertChild(child.GetYoga(), 0)
```

### 6.4 测试规范

- 每个组件至少一个测试文件
- 测试覆盖：链式 API、状态变化、事件交互
- 使用 `tenon.Run` 或 `core.NewEngine` 创建引擎实例进行集成测试

---

## 七、Widget vs Element 选择指南

| 场景 | 使用 Widget | 使用纯 Element |
|------|------------|---------------|
| 需要管理内部状态（计数器、表单数据） | ✅ | ❌ |
| 结构随状态变化（条件渲染、循环子元素） | ✅ | ❌ |
| 结构固定，只有属性变化（Button、Text） | ❌ | ✅ |
| 需要跨帧持久化复杂状态 | ✅ | ❌ |
| 简单容器，无内部状态 | ❌ | ✅ |

**简单规则**：如果组件的结构会随状态变化，用 Widget；如果结构固定只有属性变化，用纯 Element。

---

## 八、扩展：在游戏客户端中使用 Tenon

游戏客户端 import tenon 后，基于 `core.Element` 扩展专用组件：

```go
package main  // 你的游戏项目

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/sjm1327605995/tenon/pkg/v2/core"
)

// Card 是游戏专用组件，不在 Tenon 框架中
type Card struct {
    core.BaseElement
    frontImage *ebiten.Image
    backImage  *ebiten.Image
    // ...
}

func NewCard() *Card {
    c := &Card{}
    c.Init(c)
    return c
}

func (c *Card) ElementType() string { return "Card" }
func (c *Card) Draw(screen *ebiten.Image) { /* ... */ }
```

参考 `example/card/card.go` 获取完整实现。
