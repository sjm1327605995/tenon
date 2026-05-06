# Tenon v2 生产级进阶方案

> 目标：将 Tenon v2 从 "可写 Demo" 推进到 "可支撑生产级应用"

---

## 一、现状与差距分析

### 已具备的能力
- ✅ 声明式 Widget API（SwiftUI-like 链式调用）
- ✅ 三树分离架构（Widget/Element/RenderObject）
- ✅ Yoga Flexbox 布局 + Ebiten GPU 渲染
- ✅ 完整的鼠标/键盘/滚轮事件系统
- ✅ shadcn 风格组件库（20+ 组件）
- ✅ 主题系统（Light/Dark）

### 距离生产级的核心差距

| 维度 | 现状 | 生产级要求 | 差距等级 |
|------|------|-----------|---------|
| 状态管理 | 纯手动 `Rebuild()` | 组件级状态封装 + 自动变更通知 | 🔴 致命 |
| 动画系统 | 完全缺失 | Tween/Transition + 60fps 插值 | 🔴 致命 |
| 跨层数据传递 | 全局 `GetTheme()` | BuildContext + InheritedWidget | 🟡 严重 |
| 部分重绘 | 每帧全树遍历 | 脏标记选择性绘制 | 🟡 严重 |
| 测试体系 | 无单元测试 | Widget 测试 + Golden 测试 + 交互测试 | 🟡 严重 |
| 调试工具 | 无 | Inspector + 性能 Overlay | 🟢 重要 |
| 原生集成 | Ebiten 基础输入 | 文件对话框/托盘/剪贴板/通知 | 🟢 重要 |
| 无障碍 | 无 | 屏幕阅读器/键盘导航 | 🟢 重要 |
| 国际化 | 无 | 文本替换/RTL/时区数字 | 🟢 重要 |

---

## 二、总体策略

采用 **"底层能力先行，上层组件跟进"** 的阶梯策略：

```
Phase 1: 状态管理 + BuildContext（地基）
    ↓
Phase 2: 动画系统 + 部分重绘（体验）
    ↓
Phase 3: 测试体系 + 调试工具（工程质量）
    ↓
Phase 4: 原生集成 + 无障碍 + 国际化（平台化）
```

**为什么状态管理必须先做？**
- 动画系统依赖状态驱动（`AnimationController` 需要 `setState` 触发动画帧）
- 部分重绘的脏标记依赖 BuildContext 的依赖追踪
- 测试体系需要 `StatefulWidget` 来测试组件内部状态变化
- 后续所有上层能力都假设存在 `BuildContext`

---

## 三、分阶段详细方案

---

### Phase 1: 状态管理与 BuildContext（预计 2-3 周）

#### 1.1 StatefulWidget / StatefulElement

**目标**：让组件能够封装自己的状态，`setState` 自动触发 rebuild，告别手动 `Rebuild()`。

**设计**：

```go
// StatefulWidget 是有状态的 Widget。
type StatefulWidget interface {
    Widget
    CreateState() State
}

// State 是有状态组件的状态持有者。
type State[T StatefulWidget] interface {
    // Build 描述该状态下的 UI。
    Build(BuildContext) Widget
    
    // SetState 标记该组件需要 rebuild。
    // fn 是可选的状态变更函数，在标记 dirty 之前执行。
    SetState(fn func())
    
    // GetWidget 返回当前关联的 Widget（可能已更新）。
    GetWidget() T
    
    // GetContext 返回 BuildContext。
    GetContext() BuildContext
    
    // 生命周期钩子（可选实现）
    InitState()
    DidUpdateWidget(oldWidget T)
    Dispose()
}

// StatefulElement 管理 StatefulWidget 的生命周期和 State 实例。
type StatefulElement struct {
    ComponentElement
    state State
}
```

**关键行为**：
- `StatefulElement.Mount()` 时调用 `widget.CreateState()`，然后 `state.InitState()`
- `StatefulElement.Update()` 时调用 `state.DidUpdateWidget(oldWidget)`，然后 `state.Build(ctx)`
- `SetState()` 将 StatefulElement 标记为 dirty，下一帧 `flushBuild` 时调用 `state.Build()` 生成新 widget 树
- `StatefulElement.Unmount()` 时调用 `state.Dispose()`

**与现有架构的关系**：
- `StatefulElement` 是 `ComponentElement` 的子类（无 RenderObject）
- `Build()` 返回的 widget 树走现有的 `UpdateChild` diff 逻辑
- 外部状态（如 Gallery 中的 `g.counter`）仍然可以存在，StatefulWidget 用于组件内部状态

#### 1.2 BuildContext

**目标**：提供跨层数据查找能力，替代全局 `GetTheme()`，支持依赖追踪。

**设计**：

```go
// BuildContext 是 Widget 构建时的上下文。
type BuildContext interface {
    // 获取当前 Widget
    GetWidget() Widget
    
    // 向上查找最近的指定类型的祖先 Widget
    FindAncestorWidgetOfExactType[T Widget]() T
    
    // 获取指定类型的 InheritedWidget，并注册依赖关系
    GetInheritedWidgetOfExactType[T InheritedWidget]() T
    
    // 获取当前 RenderObject 的尺寸（仅在 layout 后有效）
    GetSize() Size
}
```

#### 1.3 InheritedWidget

**目标**：实现跨层数据传递 + 自动依赖通知。

**设计**：

```go
// InheritedWidget 是一种特殊 Widget，用于在树中向下传递数据。
// 当数据变化时，自动 rebuild 所有依赖它的子节点。
type InheritedWidget interface {
    Widget
    // UpdateShouldNotify 判断新数据是否发生了变化。
    UpdateShouldNotify(oldWidget InheritedWidget) bool
}

// InheritedElement 管理 InheritedWidget 的依赖关系。
type InheritedElement struct {
    ComponentElement
    dependents map[Element]struct{} // 依赖该数据的子元素
}
```

**关键行为**：
- 子元素在 `Build()` 中调用 `ctx.GetInheritedWidgetOfExactType[Theme>]()` 时，自动注册到 `InheritedElement.dependents`
- 当 `InheritedWidget` 更新且 `UpdateShouldNotify` 返回 true 时，遍历 `dependents` 并逐个标记 dirty
- `Theme` 包装为 `InheritedTheme`，全局 `GetTheme()` 逐步迁移到 `Theme.Of(ctx)`

#### 1.4 Gallery 改造示例

**改造前**：
```go
type gallery struct {
    counter int
}

func (g *gallery) Build() Widget {
    return Button(fmt.Sprintf("Clicked %d", g.counter)).OnTap(func() {
        g.counter++
        Rebuild()
    })
}
```

**改造后**：
```go
type CounterButton struct {
    BaseWidget
}

func (c CounterButton) CreateState() State { return &counterState{} }

type counterState struct {
    count int
}

func (s *counterState) Build(ctx BuildContext) Widget {
    return Button(fmt.Sprintf("Clicked %d", s.count)).OnTap(func() {
        s.SetState(func() {
            s.count++
        })
    })
}
```

#### 1.5 验收标准
- [ ] `StatefulWidget` + `State` + `StatefulElement` 完整实现
- [ ] `BuildContext` 接口及实现
- [ ] `InheritedWidget` + `InheritedElement` 依赖追踪实现
- [ ] `Theme` 改造为 `InheritedTheme`，提供 `Theme.Of(ctx)`
- [ ] Gallery 中至少 50% 的手动状态迁移到 `StatefulWidget`
- [ ] `go test ./pkg/v2/ui/...` 包含 StatefulElement 的 mount/update/unmount 测试

---

### Phase 2: 动画系统与部分重绘（预计 2-3 周）

#### 2.1 AnimationController + Tween

**目标**：提供声明式动画能力，支持 60fps 插值。

**设计**：

```go
// AnimationController 控制动画的时间线。
type AnimationController struct {
    Duration   time.Duration
    LowerBound float64 // 默认 0
    UpperBound float64 // 默认 1
    Value      float64 // 当前值
    Status     AnimationStatus // Dismissed/Forward/Reverse/Completed
}

func (a *AnimationController) Forward()
func (a *AnimationController) Reverse()
func (a *AnimationController) Stop()
func (a *AnimationController) SetStateCallback(fn func()) // 绑定到 StatefulState.SetState

// Tween 定义值区间。
type Tween[T any] struct {
    Begin T
    End   T
}

func (t *Tween[T]) Evaluate(progress float64) T

// CurvedAnimation 提供缓动曲线。
type CurvedAnimation struct {
    Parent   *AnimationController
    Curve    Curve // EaseIn/EaseOut/EaseInOut/Linear/etc.
}
```

**关键行为**：
- `AnimationController` 持有 `Ticker`（基于 Ebiten 的帧回调），每帧更新 `Value`
- 值变化时调用绑定的 `SetStateCallback`（即 `StatefulState.SetState`），触发 widget rebuild
- RenderObject 在 `Paint()` 中读取当前动画值（如 opacity、offset、scale）

**使用示例**：
```go
type FadeInState struct {
    controller *AnimationController
    opacity    *Tween[float64]
}

func (s *FadeInState) InitState() {
    s.controller = &AnimationController{Duration: 300 * time.Millisecond}
    s.opacity = &Tween[float64]{Begin: 0, End: 1}
    s.controller.SetStateCallback(func() { s.SetState(nil) })
    s.controller.Forward()
}

func (s *FadeInState) Build(ctx BuildContext) Widget {
    opacity := s.opacity.Evaluate(s.controller.Value)
    return Container(Text("Hello")).Alpha(opacity)
}
```

#### 2.2 AnimatedWidget / ImplicitlyAnimatedWidget

**目标**：提供更高层的动画封装，减少样板代码。

**设计**：

```go
// AnimatedContainer 在属性变化时自动过渡动画。
func AnimatedContainer(
    child Widget,
    duration time.Duration,
    backgroundColor *render.Color,
    width float32,
    // ... 其他属性
) Widget

// 内部使用 AnimationController + Tween，
// UpdateRenderObject 检测到属性变化时自动启动动画。
```

#### 2.3 脏标记与按需绘制

**目标**：避免每帧全树遍历构建、布局和绘制。

**现状问题**：
- `Engine.Update()` 每帧都调用 `flushBuild()`、`CalculateLayout()`、`syncBounds()` 和全树 `Draw()`，即使没有任何变化。
- `FlushPaint()` 只清除脏标记，实际绘制仍遍历整棵 RenderObject 树。

**Yoga 已覆盖的布局脏标记**：

Yoga 内部已实现完整的脏标记机制：
- 节点样式变化时自动调用 `markDirtyAndPropagate()`，向上标记自己和所有祖先
- `CalculateLayout()` 内部只遍历 dirty 节点，clean 子树直接复用上次结果

因此 **Tenon 不需要自建一套 Yoga 级别的脏区域系统**。当前每帧调用 `CalculateLayout()` 的开销已因 Yoga 内部优化而大幅降低。

**Tenon 层需要补充的优化**：

1. **Build 层**：已有的 `dirtyElements` + `needsBuild` 机制已支持局部 rebuild，但 build 完成后仍触发全量 layout/paint。build 完成后应只标记受影响的 RenderObject 为 `needsPaint`，而非全部。
2. **syncBounds 层**：`syncBounds()` 每帧全树遍历赋值。可配合 Yoga 的 `IsDirty()` 或 Tenon 自有的 `needsLayout` 标记，跳过未变化的分支。
3. **Paint 层**：Ebiten 的 `screen` 是单张全屏纹理，GPU 绘制全屏代价很低，"像素级部分重绘"收益有限。真正的优化是跳过 **未标记 `needsPaint` 的子树**（减少 draw call 和状态切换），而非裁剪绘制区域。

**优化策略**：

```go
func (e *Engine) Update() error {
    // ...

    if e.needsBuild || len(e.dirtyElements) > 0 {
        e.flushBuild() // 局部 rebuild，只影响 dirty 分支
    }

    // Yoga 内部已做 dirty 优化，CalculateLayout 开销可控。
    // 若根节点 clean，可进一步跳过 syncBounds。
    if e.rootYoga.IsDirty() || e.hasLayoutDirty() {
        e.rootYoga.CalculateLayout(...)
        e.syncBounds() // 可增量：只同步 Yoga dirty 节点
        e.flushLayout()
    }

    // 绘制：全树遍历，但跳过未标记 needsPaint 且子树也无 dirty 的分支
    e.Draw(screen)

    // ...
}
```

**分层绘制（动画优化）**：
- 对于高频动画的 RenderObject（如 loading spinner、粒子效果），单独绘制到离屏 `ebiten.Image` 缓存
- 主绘制循环只贴图，避免每帧重复构建复杂路径
- 动画结束或属性变化时失效缓存

#### 2.4 验收标准
- [ ] `AnimationController` + `Tween` + `Curve` 完整实现
- [ ] 至少 3 种内置缓动曲线（Linear、EaseInOut、Spring）
- [ ] `AnimatedContainer` 实现（背景色/尺寸/圆角过渡）
- [ ] ProgressBar loading spinner 改为旋转动画
- [ ] `syncBounds` 改为增量同步（跳过 Yoga clean 节点）；`Draw()` 跳过未标记 `needsPaint` 的子树
- [ ] Gallery 添加动画 Demo 页面（FadeIn/Slide/Scale）
- [ ] 帧率稳定在 60fps（通过 Ebiten 的 `ebiten.ActualFPS()` 验证）

---

### Phase 3: 测试体系与调试工具（预计 2 周）

#### 3.1 Widget 测试框架

**目标**：无需 GUI 即可测试组件布局和行为。

**设计**：

```go
// TestWidget 在测试环境中构建并布局 Widget。
func TestWidget(t *testing.T, w Widget, width, height float32) *TestEnvironment

type TestEnvironment struct {
    RootRenderObject render.RenderObject
}

func (e *TestEnvironment) FindText(text string) *render.RenderText
func (e *TestEnvironment) TapAt(x, y float32) // 模拟点击
func (e *TestEnvironment) AssertBounds(ro render.RenderObject, expected Bounds)
```

**使用示例**：
```go
func TestButtonTap(t *testing.T) {
    var clicked bool
    env := TestWidget(t, Button("OK").OnTap(func() {
        clicked = true
    }), 400, 300)
    
    btn := env.FindText("OK")
    env.TapAt(btn.GetBounds().X+5, btn.GetBounds().Y+5)
    
    if !clicked {
        t.Fatal("button not clicked")
    }
}
```

#### 3.2 Golden 测试

**目标**：捕获组件的视觉输出并与基准对比。

**设计**：

```go
func TestButtonSnapshot(t *testing.T) {
    env := TestWidget(t, Button("Submit"), 200, 60)
    env.Snapshot("button_default.png")
    // 与基准图片对比，像素差异 > threshold 则失败
}
```

**实现**：使用 Ebiten 的 `screen.Dump()` 或 `ebiten.NewImage` 离屏渲染，保存为 PNG。

#### 3.3 Widget Inspector

**目标**：运行时查看 Widget/Element/RenderObject 树结构和属性。

**设计**：

```go
// 在 Gallery 中按 F12 打开 Inspector Overlay
inspector := NewInspector(engine)
inspector.Show() // 绘制半透明覆盖层，显示：
// - 鼠标悬停节点的 bounds（红色边框）
// - 节点类型、widget 属性、renderObject 属性
// - Yoga 布局结果（width/height/x/y）
// - 脏标记状态（needsLayout/needsPaint）
```

#### 3.4 性能 Overlay

**目标**：实时显示性能指标。

**设计**：

```go
// 在 Engine.Draw() 最后叠加绘制：
// FPS: 60
// Frame Time: 16.2ms
// Layout Objects: 45
// Paint Objects: 45
// Build Time: 0.3ms
// Layout Time: 1.2ms
// Paint Time: 2.1ms
```

#### 3.5 验收标准
- [ ] Widget 测试框架支持 mount/update/tap/assertBounds
- [ ] 核心组件（Button、Container、Row、Column、Text）100% 覆盖 Widget 测试
- [ ] Golden 测试基础设施（截图 + 像素对比）
- [ ] Inspector 支持查看树结构、bounds、属性
- [ ] 性能 Overlay 显示 FPS / Frame Time / Layout&Paint 耗时

---

### Phase 4: 原生集成与平台化（预计 3-4 周）

#### 4.1 原生平台集成（通过 CGO / syscall）

**目标**：弥补 Ebiten 不提供的高级平台能力。

**能力清单**：

| 能力 | 平台 | 实现方式 |
|------|------|---------|
| 文件对话框（打开/保存） | Win/Mac/Linux | `github.com/sqweek/dialog` 或 CGO |
| 系统托盘 | Win/Mac/Linux | `github.com/getlantern/systray` |
| 剪贴板读写 | 全平台 | `golang.design/x/clipboard` |
| 本地通知 | Win/Mac/Linux | 平台特定 CGO / syscall |
| 窗口控制（最小化/置顶/透明） | 全平台 | Ebiten API 扩展 |
| 拖拽文件进窗口 | 全平台 | Ebiten `DropFile` 事件 |
| 深度链接（URL Scheme） | Win/Mac | 注册表/Info.plist |

**封装方式**：

```go
package native

func ShowOpenFileDialog(opts FileDialogOptions) (string, error)
func ShowSaveFileDialog(opts FileDialogOptions) (string, error)
func SetSystemTray(icon []byte, menu *TrayMenu)
func ShowNotification(title, body string)
func ReadClipboard() string
func WriteClipboard(text string)
```

#### 4.2 键盘导航与无障碍（Accessibility）

**目标**：支持键盘操作和屏幕阅读器。

**设计**：

```go
// FocusManager 管理焦点树。
type FocusManager struct {
    focused Element
    focusNodes []FocusNode
}

// FocusNode 附加到可聚焦的 Widget。
type FocusNode struct {
    CanFocus bool
    CanRequestFocus bool
    OnFocus func()
    OnUnfocus func()
}

// 快捷键系统
type Shortcut struct {
    Key ebiten.Key
    Modifiers []ebiten.Key // Ctrl/Shift/Alt
    Handler func()
}
```

**无障碍基础**：
- 为每个 RenderObject 添加 `SemanticLabel` 和 `Role`（button/text/input 等）
- 聚焦时通过平台 API 通知屏幕阅读器（Windows UI Automation / macOS Accessibility / Linux AT-SPI）
- Tab 键在可聚焦节点间切换

#### 4.3 国际化（i18n）

**目标**：支持多语言和 RTL 布局。

**设计**：

```go
// Localization 作为 InheritedWidget 注入。
type Localization struct {
    Locale string // "zh-CN", "en-US"
    Translations map[string]string
}

func L(ctx BuildContext, key string) string {
    loc := ctx.GetInheritedWidgetOfExactType[Localization]()
    return loc.Translations[key]
}

// RTL 支持：在 Localization 中标记 Direction，
// Row 的 flex-direction 自动反转，Text 对齐自动适配。
```

#### 4.4 路由系统（Navigator）

**目标**：支持多页面/模态框/侧边栏导航。

**设计**：

```go
// Navigator 管理页面栈。
type Navigator struct {
    pages []Page
}

type Page struct {
    Name string
    Builder func(ctx BuildContext) Widget
}

func (n *Navigator) Push(page Page)
func (n *Navigator) Pop()
func (n *Navigator) PushReplacement(page Page)

// 路由声明
Navigator(
    routes: map[string]RouteBuilder{
        "/":       HomePage,
        "/settings": SettingsPage,
        "/detail/:id": DetailPage,
    },
)
```

#### 4.5 验收标准
- [ ] 文件对话框在 Win/Mac/Linux 正常工作
- [ ] 剪贴板读写跨平台正常
- [ ] Tab 键在 Button/Input/Checkbox 之间正常切换焦点
- [ ] 至少支持中英文切换 + RTL 布局适配
- [ ] Navigator 支持 Push/Pop/Replace + 页面转场动画
- [ ] 系统托盘 + 本地通知在至少一个平台验证

---

## 四、工程保障

### 代码规范
- 所有新增核心代码必须通过 Widget 测试
- `go vet` 零警告
- 新增组件必须同步更新 Gallery Demo

### 版本策略
- 每个 Phase 完成后打 tag（`v2.1.0`、`v2.2.0`、`v2.3.0`、`v2.4.0`）
- Phase 1-2 之间不做破坏性 API 变更
- Phase 1 引入的 `StatefulWidget` 与现有纯函数 API 共存，逐步迁移

### 性能基准
- Gallery 启动时间 < 500ms
- 100 个组件的列表滚动 FPS > 55
- 内存占用增长稳定（无泄漏）

---

## 五、风险与应对

| 风险 | 影响 | 应对 |
|------|------|------|
| StatefulWidget 引入后代码复杂度上升 | 🔴 高 | 先做一个最小可行版本（只支持 `setState` + `InitState/Dispose`），暂不实现 `didChangeDependencies` 等高级生命周期 |
| Animation 系统与 Ebiten 帧循环耦合 | 🟡 中 | AnimationController 基于 Ebiten 的 `Update()` 回调，不独立开线程，避免同步问题 |
| 部分重绘优化收益不明显 | 🟡 中 | 先做 profiling 确认瓶颈。如果 GPU 绘制不是瓶颈，优先优化 Yoga layout 调用频率 |
| 原生平台集成维护成本高 | 🟢 低 | 先只支持 Windows，再逐步扩展 Mac/Linux。优先使用成熟第三方库 |

---

## 六、总结

Tenon v2 的**架构底子是好的**（Flutter 范式 + Yoga + Ebiten），当前最大的瓶颈是**状态管理缺失**和**动画缺失**，这两点补齐后，框架会从 "能写 Demo" 跃升到 "能写产品"。

**推荐执行顺序**：
1. **立即启动 Phase 1**（状态管理 + BuildContext），这是所有上层能力的地基
2. Phase 1 完成后发布 `v2.1.0`，将 Gallery 全部迁移到 `StatefulWidget`
3. 并行启动 Phase 2 动画系统和 Phase 3 测试框架（两者无强依赖）
4. Phase 4 原生集成按平台逐步推进，不急迫

预估总工期：**9-12 周**（1 个全职开发者），关键里程碑在 Phase 1 结束（第 3 周）。
