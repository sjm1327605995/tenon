# Tenon v2 — AI 开发约束与踩坑记录

> 本文档记录框架中**非显而易见的设计约束**和**已踩过的坑**，供 AI 编码时参考。
> 与 `ARCHITECTURE.md` 的区别：ARCHITECTURE 描述「应该怎样做」，本文档描述「千万不要这样做以及为什么」。

---

## 一、架构红线：组件禁止直接操纵子树

### ❌ 禁止的模式

```go
// 在组件的 setter / 事件回调中直接重建子元素
func (p *Pagination) SetPage(page int) {
    p.current = page
    p.ClearChildren()           // ← 禁止
    p.buildButtons()            // ← 禁止：动态创建子元素
    p.Mark(core.FlagNeedLayout) // ← 不应由组件直接触发
}

func (p *Pagination) buildButtons() {
    p.ClearChildren()           // ← 禁止
    for i := 1; i <= p.totalPages; i++ {
        btn := NewButton(strconv.Itoa(i))  // ← 禁止：在 setter 中创建子元素
        p.AppendChild(btn)      // ← 禁止
    }
}
```

### ✅ 正确的模式：固定子元素 + 属性更新

```go
// 1. 构造函数中一次性创建所有子元素
func NewPagination(totalPages int) *Pagination {
    p := &Pagination{totalPages: totalPages, current: 1}
    p.Init(p)
    // ... 创建 prev/next/page buttons，保存到字段
    p.initButtons()   // 只创建一次
    p.updateButtons() // 设置初始属性
    return p
}

// 2. setter 只更新状态 + 修改已有元素的属性
func (p *Pagination) SetPage(page int) {
    p.current = page
    p.updateButtons() // 只改 text/variant/display，不创建/销毁节点
    // 不需要手动 Mark，SetText/SetDisplay 内部会自动 Mark
}

func (p *Pagination) updateButtons() {
    p.prevBtn.SetDisabled(p.current <= 1)
    for i, btn := range p.pageButtons {
        if pageIsVisible(i) {
            btn.SetDisplay(yoga.DisplayFlex)
            btn.SetText(fmt.Sprintf("%d", pageNum))
            btn.SetVariant(selectedVariant)
        } else {
            btn.SetDisplay(yoga.DisplayNone)
        }
    }
    p.nextBtn.SetDisabled(p.current >= p.totalPages)
}
```

### 为什么禁止

框架的设计契约是：
- **Widget 层**（声明式）通过 `Render() → BuildEngine.patchElement()` 驱动结构变化
- **Element 层**（命令式）只负责绘制和事件，不主动改变树结构
- **只有 native 元素**（如单个 Button、Text）在属性变化时触发 framework signals

如果组件直接 `ClearChildren()` + `AppendChild()`，会绕过 `BuildEngine` 的 diff/patch 逻辑，导致：
1. 新旧 Element 无法复用，性能下降
2. Yoga 节点树与 Element 树不同步
3. 子元素的 dirty flags 和 engine 注入时机混乱（见下一节 zero-bounds bug）

---

## 二、Dirty Bus 陷阱：Pre-mount 构造的 Zero-Bounds Bug

### 现象

组件在构造函数中创建的子元素，调用 `Mark(FlagNeedLayout)` 后，mount 到引擎并执行 `SetPage()`/`SetDate()` 后，**新创建的子元素 bounds 为 `{0,0,0,0}`**。

### 根因分析

```
构造时（pre-mount）                    Mount 后
───────────────────                    ────────
Mark(FlagNeedLayout)                    hadDirty=true（旧 flag 还在）
  → engine == nil                       Mark(FlagNeedLayout)
  → 无法进入 dirtyBus                     → hadDirty=true（旧 flag 导致）
  → dirty flag 留在元素上                 → 不会重新 post 到 dirtyBus
                                         → dirtyBus.Len() == 0
                                         → flushDirtyElements 跳过
                                         → 新元素永远得不到布局
```

### 修复方案

在 `Engine.onElementMounted()` 中，递归挂载子元素时，**如果元素已携带 dirty flags，重新 post 到 dirtyBus**：

```go
func (e *Engine) onElementMounted(el Element) {
    // ... 现有 mount 逻辑 ...
    
    // FIX: 重新投递 pre-mount 期间累积的 dirty flags
    if el.GetFlags()&FlagDirtyMask != 0 {
        e.dirtyBus.Post(el)
    }
    
    for _, child := range el.GetChildren() {
        if child.GetEngine() == nil {
            e.onElementMounted(child)
        }
    }
}
```

### AI 编码约束

1. **不要假设 pre-mount 的 `Mark()` 会生效**。构造组件时如果必须触发布局，确保在 mount 后有补救机制，或干脆等 mount 后再做首次布局。
2. **如果修改了 `onElementMounted`，必须同步更新此文档**。
3. **新增组件时**，如果构造函数中创建了子元素并调用了子元素的 `Mark()`，请用 `test_repro.go` 模式验证 mount 后的 bounds 是否正确。

---

## 三、调试方案：无 GUI 测试 vs HTML 预览

### 已有基础设施

项目同时支持两种调试方式，**根据场景选择**：

| 场景 | 推荐方式 | 命令/工具 |
|------|---------|----------|
| 回归测试、验证布局数值 | **无 GUI 测试** | `go run _debug_snapshot/test_repro.go` |
| 开发新组件、肉眼确认视觉效果 | **HTML 预览** | `go run _debug_snapshot/main.go` |
| CI 集成 | 无 GUI 测试 | `go test ./...` |
| 实时调试（运行时） | WebSocket 调试器 | `debug.NewDebugger(engine, port)` |

### 方案 A：无 GUI 测试（主力，推荐用于 AI 编码）

```go
// _debug_snapshot/test_repro.go 模式
func main() {
    // 1. 初始化字体和主题
    fonts.InitDefaultFont()
    core.SetTheme(core.DefaultShadcnLightTheme())

    // 2. 构建测试场景
    root := buildTestScene()

    // 3. 创建引擎并挂载
    app := &testWidget{root: root}
    e := core.NewEngine(app, 900, 700)
    e.Mount()

    // 4. 触发状态变化
    pg.SetPage(5)
    e.Update()  // ← 手动触发一帧更新

    // 5. 断言 bounds
    assertBounds(pg.GetChildren()[0], 24, 60, 48, 38)
}
```

**优点**：
- 快速（秒级运行）
- 可自动化断言 bounds、Yoga 树结构
- 不需要 Ebiten 窗口或 GPU
- CI 友好

**缺点**：
- 无法肉眼确认视觉问题（颜色、间距细节）
- 无法捕获渲染层面的 bug

### 方案 B：HTML 预览（辅助，用于视觉确认）

```go
// _debug_snapshot/main.go 模式
d := debug.NewDebugger(engine, 8765)
html := d.GenerateHTML(root)
os.WriteFile("preview.html", []byte(html), 0644)
```

生成的 HTML 用绝对定位 div 模拟 Yoga 布局结果，可在浏览器中打开查看。

**优点**：
- 可视化好，一眼看出错位
- 可对比修改前后的 HTML diff

**缺点**：
- Yoga flex 与 CSS flex 有细微差异（如 `MeasureFunc` 回调的文字测量）
- HTML 只是近似，不是像素级精确
- 无法断言，不适合回归测试

### AI 编码时的推荐流程

```
1. 修改代码
   ↓
2. go test ./pkg/v2/...          ← 快速验证不破坏现有功能
   ↓
3. go run _debug_snapshot/test_repro.go  ← 验证具体 bug 是否修复
   ↓
4.（可选）go run _debug_snapshot/main.go
   → 打开生成的 preview_*.html 肉眼确认
```

---

## 四、常见 API 误用

### `Mark()` 的脏标记语义

```go
// 修改布局属性（width/height/margin/flex/gap/display）
→ Mark(core.FlagNeedLayout)

// 修改视觉属性（color/background/border/radius/alpha）
→ Mark(core.FlagNeedDraw)

// 修改内容（文字内容、图片源）
→ Mark(core.FlagNeedMeasure | core.FlagNeedLayout | core.FlagNeedDraw)
```

注意：`BaseElement.SetDisplay()`、`Text.SetContent()` 等标准 API **内部已自动调用 `Mark`**，外部不需要重复调用。

### 构造函数中必须 `Init(self)`

```go
func NewXxx() *Xxx {
    x := &Xxx{}
    x.Init(x)  // ← 忘记这行会导致 yoga 节点未创建、self 未绑定
    return x
}
```

### Yoga `RemoveAllChildren` 在空节点上会 panic

不要在 setter 中懒创建子元素，所有子元素必须在构造函数中创建（参见 6.2 原则 2）。

---

## 五、验证清单（修改核心引擎或组件后）

- [ ] `go test ./pkg/v2/...` 通过
- [ ] `go run _debug_snapshot/test_repro.go` 无 zero-bounds 输出
- [ ] 如果是新组件，`test_repro.go` 中包含状态变化后的 bounds 断言
- [ ] 组件 setter 中没有 `ClearChildren()` 或动态 `NewXxx()` 创建子元素
- [ ] 如果修改了 `onElementMounted`/`flushDirtyElements`/`Mark`，检查 pre-mount 场景
