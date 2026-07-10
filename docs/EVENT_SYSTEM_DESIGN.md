# 事件系统设计

> 父文档：[ARCHITECTURE.md](./ARCHITECTURE.md)
>
> 调研范围：Flutter GestureDetector/GestureArena、Web DOM Event Flow、React SyntheticEvent

---

## 一、设计目标

1. **支持所有交互形式**：Tap、Pan、Scroll、LongPress、DragAndDrop、Hover、键盘输入
2. **手势消歧**：同区域多个手势竞争时，通过规则自动决出胜者（如 Tap vs Pan vs LongPress）
3. **事件冒泡与捕获**：支持父子组件间的事件传递控制
4. **Overlay 隔离**：弹窗/下拉/Toast 等覆盖层自动阻断底层事件
5. **避免堆叠代码**：Recognizer 状态机 + 事件对象池 + 统一分发管线，不写 if-else 嵌套

---

## 二、分层架构

```
┌─────────────────────────────────────────────────────────────┐
│  Component Layer (用户层)                                    │
│  Button(OnTap), Input(OnChange), GestureDetector(OnPan)...  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│  Gesture Layer (手势识别)                                    │
│  TapRecognizer / PanRecognizer / LongPressRecognizer        │
│  ScrollRecognizer / DnDRecognizer / HoverRecognizer         │
│  ──▶ GestureArena (每帧裁决胜者)                             │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│  Event Flow Layer (事件流)                                   │
│  Capture Phase → Target Phase → Bubble Phase                │
│  EventPath = [root, ..., target]                            │
│  event.StopPropagation() / event.StopImmediatePropagation() │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│  Hit-Test Layer (命中测试)                                   │
│  每帧 PointerDown 时计算一次，缓存路径                        │
│  优先 Overlay 树 → 后主树                                   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│  Pointer Layer (原始输入)                                    │
│  Ebiten: CursorPosition / IsMouseButtonPressed              │
│  AppendInputChars / IsKeyJustPressed                        │
│  ──▶ 封装为统一 PointerEvent / KeyEvent                      │
└─────────────────────────────────────────────────────────────┘
```

---

## 三、Pointer 层：统一原始输入

Ebiten 提供多种输入 API，我们在框架内部统一为 `PointerEvent`：

```go
type PointerEvent struct {
    Type      PointerEventType  // Down, Move, Up, Cancel, Scroll
    PointerID int               // 鼠标=0, 触摸=1,2,3...
    X, Y      float32           // 屏幕坐标
    DX, DY    float32           // 相对上一帧位移
    Pressure  float32           // 压力（触摸支持时）
    Buttons   PointerButtons    // 左键/右键/中键
    Modifiers KeyModifiers      // Shift/Ctrl/Alt/Meta
    Timestamp int64             // 纳秒时间戳
}

type PointerEventType int
const (
    PointerDown PointerEventType = iota
    PointerMove
    PointerUp
    PointerCancel
    PointerScroll  // 滚轮事件
)
```

**Ebiten 输入 → PointerEvent 转换：**

```go
// 每帧 Update 中
func (e *Engine) pollPointers() []PointerEvent {
    var events []PointerEvent
    
    // 鼠标左键
    if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
        x, y := ebiten.CursorPosition()
        events = append(events, PointerEvent{Type: PointerDown, X: float32(x), Y: float32(y)})
    }
    if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
        x, y := ebiten.CursorPosition()
        events = append(events, PointerEvent{Type: PointerMove, X: float32(x), Y: float32(y)})
    }
    if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
        x, y := ebiten.CursorPosition()
        events = append(events, PointerEvent{Type: PointerUp, X: float32(x), Y: float32(y)})
    }
    
    // 滚轮
    _, dy := ebiten.Wheel()
    if dy != 0 {
        x, y := ebiten.CursorPosition()
        events = append(events, PointerEvent{Type: PointerScroll, X: float32(x), Y: float32(y), DY: float32(dy)})
    }
    
    return events
}
```

---

## 四、Hit-Test 层：命中测试

**核心规则**：
1. 从 `rootRenderNode` 递归 DFS
2. 检查 `event.X/Y` 是否在 `node.GetBounds()` 内
3. **优先 Overlay 树**：如果存在 Overlay，先测试 Overlay，命中则返回 Overlay 路径
4. 路径从 root 到 target 按顺序存储

```go
type EventPath struct {
    Nodes []ui.RenderNode  // 从 root 到 target
}

func (p *EventPath) Target() ui.RenderNode {
    if len(p.Nodes) == 0 { return nil }
    return p.Nodes[len(p.Nodes)-1]
}

func hitTest(root ui.RenderNode, x, y float32) *EventPath {
    var nodes []ui.RenderNode
    var dfs func(node ui.RenderNode)
    dfs = func(node ui.RenderNode) {
        bounds := node.GetBounds()
        if x < bounds.X || x > bounds.X+bounds.W ||
           y < bounds.Y || y > bounds.Y+bounds.H {
            return
        }
        nodes = append(nodes, node)
        // 子节点按 z-index 倒序遍历（后绘制的在上层）
        children := node.Children()
        for i := len(children) - 1; i >= 0; i-- {
            dfs(children[i])
        }
    }
    dfs(root)
    return &EventPath{Nodes: nodes}
}
```

**缓存策略**：同一个 PointerID 的 Down→Move→Up 序列复用同一份 EventPath，只有 Down 时重新计算。

---

## 五、事件流层：Capture → Target → Bubble

参考 Web DOM 三阶段：

```go
type EventPhase int
const (
    EventPhaseCapture EventPhase = iota
    EventPhaseTarget
    EventPhaseBubble
)

type UIEvent struct {
    Type       string           // "tap", "pan", "scroll", "key", "hover"...
    Target     ui.RenderNode    // 目标节点
    Path       *EventPath       // 完整路径
    Phase      EventPhase       // 当前阶段
    LocalX     float32          // 相对 target 的坐标
    LocalY     float32
    
    // 控制位
    propagationStopped        bool
    immediatePropagationStopped bool
    defaultPrevented          bool
}

func (e *UIEvent) StopPropagation() {
    e.propagationStopped = true
}

func (e *UIEvent) StopImmediatePropagation() {
    e.immediatePropagationStopped = true
}
```

**分发流程**：

```go
func dispatchEvent(event *UIEvent) {
    path := event.Path.Nodes
    target := event.Target
    
    // 1. Capture Phase
    event.Phase = EventPhaseCapture
    for i := 0; i < len(path)-1 && !event.propagationStopped; i++ {
        path[i].HandleEvent(event)
        if event.immediatePropagationStopped { return }
    }
    
    // 2. Target Phase
    if !event.propagationStopped {
        event.Phase = EventPhaseTarget
        target.HandleEvent(event)
        if event.immediatePropagationStopped { return }
    }
    
    // 3. Bubble Phase
    event.Phase = EventPhaseBubble
    for i := len(path) - 2; i >= 0 && !event.propagationStopped; i-- {
        path[i].HandleEvent(event)
        if event.immediatePropagationStopped { return }
    }
}
```

---

## 六、手势识别层：GestureRecognizer + GestureArena

### 6.1 设计原则

- **每个手势一个独立状态机**：避免一个大 switch 处理所有手势
- **Arena 竞争机制**：同 Pointer 的多个 Recognizer 竞争，只有一个胜出
- **延迟裁决**：Tap 需要等 200ms 确认不是 DoubleTap；Pan 需要移动超过阈值

### 6.2 GestureRecognizer 接口

```go
type GestureRecognizer interface {
    // 接收指针事件，更新内部状态
    HandlePointerEvent(event *PointerEvent)
    
    // Arena 裁决接口
    Accept()   // 胜出，触发回调
    Reject()   // 失败，重置状态
    
    // 查询状态
    IsActive() bool
    Name() string
}
```

### 6.3 具体 Recognizer

#### TapRecognizer

```go
type TapRecognizer struct {
    onTap       func()
    onTapDown   func()
    onTapUp     func()
    onTapCancel func()
    
    state       TapState
    downPos     Point
    downTime    int64
}

type TapState int
const (
    TapStateIdle TapState = iota
    TapStateDown
    TapStateAccepted
    TapStateRejected
)

func (r *TapRecognizer) HandlePointerEvent(e *PointerEvent) {
    switch e.Type {
    case PointerDown:
        r.state = TapStateDown
        r.downPos = Point{e.X, e.Y}
        r.downTime = e.Timestamp
        if r.onTapDown != nil { r.onTapDown() }
        
    case PointerMove:
        if r.state == TapStateDown {
            if distance(r.downPos, Point{e.X, e.Y}) > TapSlop {
                r.Reject()  // 移动太远，不是 Tap
            }
        }
        
    case PointerUp:
        if r.state == TapStateDown {
            r.Accept()  // 胜出，触发 onTap
        }
        
    case PointerCancel:
        r.Reject()
    }
}

func (r *TapRecognizer) Accept() {
    r.state = TapStateAccepted
    if r.onTap != nil { r.onTap() }
}

func (r *TapRecognizer) Reject() {
    if r.state == TapStateDown && r.onTapCancel != nil {
        r.onTapCancel()
    }
    r.state = TapStateRejected
}
```

#### PanRecognizer（拖拽）

```go
type PanRecognizer struct {
    onPanStart  func(details PanDetails)
    onPanUpdate func(details PanDetails)
    onPanEnd    func(details PanDetails)
    
    state       PanState
    startPos    Point
    currentPos  Point
}

func (r *PanRecognizer) HandlePointerEvent(e *PointerEvent) {
    switch e.Type {
    case PointerDown:
        r.state = PanStatePossible
        r.startPos = Point{e.X, e.Y}
        
    case PointerMove:
        switch r.state {
        case PanStatePossible:
            if distance(r.startPos, Point{e.X, e.Y}) > PanSlop {
                r.Accept()  // 移动超过阈值，胜出
                r.state = PanStateActive
                r.currentPos = Point{e.X, e.Y}
                if r.onPanStart != nil {
                    r.onPanStart(PanDetails{Start: r.startPos, Current: r.currentPos})
                }
            }
        case PanStateActive:
            r.currentPos = Point{e.X, e.Y}
            if r.onPanUpdate != nil {
                r.onPanUpdate(PanDetails{Start: r.startPos, Current: r.currentPos, Delta: Point{e.DX, e.DY}})
            }
        }
        
    case PointerUp:
        if r.state == PanStateActive {
            if r.onPanEnd != nil {
                r.onPanEnd(PanDetails{Start: r.startPos, Current: r.currentPos})
            }
            r.state = PanStateIdle
        } else {
            r.Reject()
        }
        
    case PointerCancel:
        r.Reject()
    }
}
```

#### LongPressRecognizer（长按）

```go
type LongPressRecognizer struct {
    onLongPress func()
    state       LongPressState
    downTime    int64
}

func (r *LongPressRecognizer) HandlePointerEvent(e *PointerEvent) {
    switch e.Type {
    case PointerDown:
        r.state = LongPressStatePossible
        r.downTime = e.Timestamp
        
    case PointerMove:
        if r.state == LongPressStatePossible && distance(...) > Slop {
            r.Reject()
        }
        
    case PointerUp:
        r.Reject()  // 提前抬起，不是长按
    }
}

// 每帧 Update 中检查
func (r *LongPressRecognizer) Tick(now int64) {
    if r.state == LongPressStatePossible && now-r.downTime > LongPressDuration {
        r.Accept()
        if r.onLongPress != nil { r.onLongPress() }
    }
}
```

#### ScrollRecognizer（滚动）

```go
type ScrollRecognizer struct {
    onScroll    func(details ScrollDetails)
    state       ScrollState
    velocity    Point  // 速度（用于惯性滚动）
}

func (r *ScrollRecognizer) HandlePointerEvent(e *PointerEvent) {
    if e.Type == PointerScroll {
        r.Accept()
        if r.onScroll != nil {
            r.onScroll(ScrollDetails{DeltaX: e.DX, DeltaY: e.DY})
        }
        return
    }
    // 也支持拖拽式滚动（Scrollbar / 触摸板）
    // ... 类似 PanRecognizer
}
```

#### DnDRecognizer（拖放）

```go
type DnDRecognizer struct {
    onDragStart func(details DnDDetails)
    onDragUpdate func(details DnDDetails)
    onDragEnd   func(details DnDDetails)
    onDrop      func(details DnDDetails)  // 放置到目标
    
    state       DnDState
    dragSource  ui.RenderNode
    dragData    any
}

func (r *DnDRecognizer) HandlePointerEvent(e *PointerEvent) {
    switch e.Type {
    case PointerDown:
        r.state = DnDStatePossible
        
    case PointerMove:
        if r.state == DnDStatePossible && distance(...) > DnDSlop {
            r.state = DnDStateDragging
            if r.onDragStart != nil { r.onDragStart(...) }
        }
        if r.state == DnDStateDragging {
            if r.onDragUpdate != nil { r.onDragUpdate(...) }
            // 实时 Hit-test 查找放置目标
            r.updateDropTarget(e.X, e.Y)
        }
        
    case PointerUp:
        if r.state == DnDStateDragging {
            if r.dropTarget != nil && r.onDrop != nil {
                r.onDrop(...)
            }
            if r.onDragEnd != nil { r.onDragEnd(...) }
        }
        r.state = DnDStateIdle
    }
}
```

### 6.4 GestureArena（手势竞技场）

```go
type GestureArena struct {
    members map[string]GestureRecognizer  // key = recognizer name
    winner  GestureRecognizer
    closed  bool
}

func (a *GestureArena) Add(r GestureRecognizer) {
    a.members[r.Name()] = r
}

func (a *GestureArena) Close() {
    a.closed = true
    a.sweep()
}

func (a *GestureArena) sweep() {
    if a.winner != nil { return }
    
    // 规则：只有一个成员时直接胜出
    if len(a.members) == 1 {
        for _, r := range a.members {
            a.declareWinner(r)
            return
        }
    }
    
    // 规则：有明确接受者时，其他全部拒绝
    for _, r := range a.members {
        if r.IsActive() {
            a.declareWinner(r)
            return
        }
    }
    
    // 默认：全部拒绝
    for _, r := range a.members {
        r.Reject()
    }
}

func (a *GestureArena) declareWinner(r GestureRecognizer) {
    a.winner = r
    r.Accept()
    for _, other := range a.members {
        if other != r {
            other.Reject()
        }
    }
}
```

---

## 七、组件事件层：用户 API

### 7.1 GestureDetector 组件

```go
func GestureDetector(args ...any) Widget {
    cfg := &Config{}
    parseOptions(cfg, args...)
    return wrap(&gestureDetectorWidget{cfg: cfg})
}

// 手势选项
func OnTap(fn func()) Option           { return optFunc(func(c *Config) { c.OnTap = fn }) }
func OnPanStart(fn func()) Option      { return optFunc(func(c *Config) { c.OnPanStart = fn }) }
func OnPanUpdate(fn func()) Option     { return optFunc(func(c *Config) { c.OnPanUpdate = fn }) }
func OnPanEnd(fn func()) Option        { return optFunc(func(c *Config) { c.OnPanEnd = fn }) }
func OnLongPress(fn func()) Option     { return optFunc(func(c *Config) { c.OnLongPress = fn }) }
func OnScroll(fn func()) Option        { return optFunc(func(c *Config) { c.OnScroll = fn }) }
```

### 7.2 Button 内部集成

Button 内部自动挂载 `TapRecognizer`，用户只需传 `OnTap`：

```go
func (n *buttonRenderNode) SyncYogaProps(w ui.RenderWidget) {
    // ... 样式设置
    
    // 注册 TapRecognizer
    tw := w.(*buttonWidget)
    if tw.cfg.OnTap != nil {
        n.recognizers = append(n.recognizers, &TapRecognizer{
            onTap: tw.cfg.OnTap,
        })
    }
}
```

### 7.3 Input 键盘事件

Input 不走 Pointer 流，走独立的 `TextInputConnection`：

```go
type TextInputConnection struct {
    node     *inputRenderNode
    value    string
    cursor   int
    focused  bool
}

// Engine 每帧将输入事件分发给活跃的 connection
func (e *Engine) dispatchTextInput() {
    if e.activeInput == nil { return }
    
    // 字符输入
    var runes []rune
    runes = ebiten.AppendInputChars(runes)
    for _, r := range runes {
        e.activeInput.InsertRune(r)
    }
    
    // 退格
    if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
        e.activeInput.DeleteBackward()
    }
    
    // ... 方向键、Home、End、Enter
}
```

---

## 八、Overlay / 弹窗 / 下拉事件隔离

### 8.1 架构

```
Engine
├── rootRenderNode      (主树)
└── overlayRoot         (覆盖层树，独立 RenderNode)
```

**Hit-test 规则**：
1. 先测试 `overlayRoot`，如果命中任何节点，只返回 Overlay 路径
2. 否则测试 `rootRenderNode`

### 8.2 弹窗点击外部关闭

```go
func (n *dialogRenderNode) HandleEvent(event *UIEvent) {
    if event.Type == "tap" && event.Phase == EventPhaseTarget {
        // 点击了 Dialog 背景（自身），而非内容
        if n.onDismiss != nil {
            n.onDismiss()
        }
    }
}
```

### 8.3 下拉菜单（Dropdown / Select）

- 下拉面板通过 Portal 插入到 `overlayRoot`
- 点击面板外部 → 自动关闭（Hit-test 命中主树 → 关闭 dropdown）
- 滚动主树时 → dropdown 跟随锚点更新位置

---

## 九、事件对象池

参考 React SyntheticEvent 池化，避免每帧创建大量事件对象：

```go
var eventPool = sync.Pool{
    New: func() any { return &UIEvent{} },
}

func acquireEvent() *UIEvent {
    return eventPool.Get().(*UIEvent)
}

func releaseEvent(e *UIEvent) {
    e.propagationStopped = false
    e.immediatePropagationStopped = false
    e.defaultPrevented = false
    eventPool.Put(e)
}
```

> 注意：异步访问事件时需要 `event.Persist()`（从池中移除，不回收）。

---

## 十、与当前代码的集成点

| 当前代码 | 修改 |
|---------|------|
| `RenderNode` | 添加 `HandleEvent(event *UIEvent)` 方法 |
| `BaseRenderNode` | 添加 `recognizers []GestureRecognizer` 字段 |
| `Engine.Update()` | 在 `buildIfNeeded()` 之后调用 `pollPointers()` + `dispatchEvents()` |
| `Engine.Draw()` | 不变 |
| `Button` | 内部自动创建 `TapRecognizer`，不再直接存 `OnTap` |
| `Input` | 通过 `TextInputConnection` 与 Engine 交互 |
| `GestureDetector` | 新增组件，用于包裹任意子组件添加手势 |
