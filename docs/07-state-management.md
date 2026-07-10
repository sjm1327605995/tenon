# 状态管理：`SetState` 机制

> 父文档：[ARCHITECTURE.md](./ARCHITECTURE.md)
> 前置阅读：[02-core-constraints.md](./02-core-constraints.md)

---

## 核心问题：Go 没有 Proxy

JavaScript 的 Vue/Svelte 能做自动响应式，是因为有 `Proxy`：

```js
// Vue 3: 自动追踪
const count = ref(0)
// 读取 count.value 时自动注册依赖，修改时自动触发更新
```

Go 没有等价物。任何状态变更必须是**显式的**。

---

## 推荐语法：`SetState` 包裹闭包

借鉴 Flutter 的 `setState(() { ... })`，在 Go 中实现为方法包裹：

```go
func (c *Component) SetState(fn func()) func() {
    return func() {
        fn()
        if c.element != nil {
            c.element.MarkNeedsBuild()
        }
    }
}
```

### 使用示例

```go
type Counter struct {
    ui.Component
    count int
}

func (c *Counter) Build(ctx ui.BuildContext) ui.Widget {
    return ui.Row(
        ui.Gap(16),
        ui.AlignItems(ui.AlignCenter),

        ui.Button(
            ui.OnTap(c.SetState(func() { c.count-- })),
            ui.Text("-"),
        ),

        ui.Text(fmt.Sprintf("%d", c.count), ui.FontSize(24)),

        ui.Button(
            ui.OnTap(c.SetState(func() { c.count++ })),
            ui.Text("+"),
        ),
    )
}
```

---

## 为什么不用 `c.count++; c.SetState()`

之前的设计需要在每个回调末尾手动调用 `SetState()`：

```go
// 不够友好：啰嗦，容易忘
OnTap(func() {
    c.count++
    c.SetState()
})
```

`SetState(func(){ ... })` 的优势：

1. **一眼看懂**：状态变更和重建触发在一个闭包里，边界明确
2. **不可遗漏**：把状态修改和重建绑定在一起，无法忘记
3. **可组合**：可以在闭包前/后加逻辑（乐观更新、日志等）
4. **零黑魔法**：不需要反射或运行时追踪

---

## 批量状态更新

一次修改多个字段：

```go
OnTap(c.SetState(func() {
    c.count++
    c.loading = true
    c.error = nil
}))
```

---

## 与现有 Signal 系统的兼容

当前 tenon 已有 `state.Signal`。可以保留 Signal，但**不自动追踪**：

```go
type Counter struct {
    ui.Component
    count *state.Signal[int]
}

func (c *Counter) Build(ctx ui.BuildContext) ui.Widget {
    val := c.count.Get()
    return ui.Button(
        ui.OnTap(c.SetState(func() {
            c.count.Set(val + 1)
        })),
        ui.Text(fmt.Sprintf("%d", val)),
    )
}
```

> 未来可扩展：在 `BuildContext` 中自动收集 `Signal.Get()` 调用，隐式建立订阅关系（类似 SolidJS），但这需要运行时 hook，增加复杂度。建议第一步保持显式 `SetState()`。

---

## 异步状态更新

```go
func (c *SearchPage) search() {
    c.SetState(func() {
        c.loading = true
        c.results = nil
    })()  // 立即执行，触发重建

    go func() {
        results := fetchResults(c.query)
        c.SetState(func() {
            c.loading = false
            c.results = results
        })()
    }()
}
```

---

## 与 React / Flutter 的对比

| 框架 | 语法 | 特点 |
|------|------|------|
| React | `setCount(c => c + 1)` | Hook dispatch，异步批处理 |
| Flutter | `setState(() { count++; })` | 同步包裹，立即标记脏 |
| **本架构** | `c.SetState(func() { count++ })` | 同步包裹，返回闭包供事件绑定 |
