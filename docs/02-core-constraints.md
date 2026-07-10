# Go 核心约束与破局思路

> 父文档：[ARCHITECTURE.md](./ARCHITECTURE.md)
> 相关调研：[01-research.md](./01-research.md)

---

## 约束清单

| 约束 | 影响 |
|------|------|
| **无 Proxy** | 无法自动追踪状态依赖 → Vue 式响应式不可行 |
| **无泛型变参** | 函数式 children API 需要 `any` 或代码生成 |
| **无 JSX** | 声明式 API 只能依赖函数调用嵌套 |
| **无命名参数** | 属性设置需用 Options Pattern 或 Struct Literal |
| **有 GC** | 每帧创建短生命周期的 Widget 描述对象是可行的 |
| **有 interface** | 非常适合建模多态的 Widget / Element / RenderNode |
| **有 goroutine** | 可用于后台任务，但 UI 操作必须单线程 |

---

## 约束 1：无 Proxy → 放弃 Vue 式响应式

### 为什么不可行

Vue 3 的核心是运行时依赖追踪：

```js
// Vue: 读取 count 时自动注册依赖
const count = ref(0)
function render() {
  return h('div', count.value)  // 自动建立订阅
}
```

Go 没有 Proxy/Reflect API，无法在属性访问时插入钩子。任何模拟方案都需要：
- 代码生成（改造编译器/ast）
- 显式 `Get()/Set()` 调用（失去透明性）
- 全局读写拦截（不可能）

**结论**：Vue 的响应式模式在 Go 中是死胡同。

### 替代方案

采用 React 的**显式状态更新**模型：

```go
// 用户显式标记重建
func (c *Counter) inc() {
    c.count++
    c.SetState()  // 或包裹在 SetState 闭包里
}
```

---

## 约束 2：无 JSX → 函数嵌套模拟

### React JSX

```jsx
<Column gap={16}>
  <Text>Hello</Text>
  <Button onTap={handleTap}>Click</Button>
</Column>
```

### Go 等价写法

```go
Column(
    Gap(16),
    Text("Hello"),
    Button(OnTap(handleTap), Text("Click")),
)
```

实现技巧：
- 容器接收 `...any`，运行时区分 `Option` 和 `Widget`
- 属性用函数选项（Functional Options）
- Children 直接传入 `Widget` 返回值

详见 [04-declarative-api.md](./04-declarative-api.md)。

---

## 约束 3：无命名参数 → Options Pattern

Go 不支持：

```go
// 不支持
Box(background: ColorRed, padding: 16)
```

采用 Functional Options：

```go
Box(
    Background(ColorRed),
    Padding(16),
)
```

或 Struct Literal（适合大量属性）：

```go
Box(BoxProps{
    Background: ColorRed,
    Padding:    16,
    Children:   []Widget{Text("Hello")},
})
```

框架内部统一为 Options 模式，对外暴露函数式 API。

---

## 约束 4：泛型变参限制

Go 1.18+ 泛型不支持变参中的类型约束：

```go
// 不支持：T 不能同时是 Option 和 Widget
func Row[T Option | Widget](args ...T) Widget
```

解决方案：

```go
// 顶层 API 用 any，运行时类型断言
func Row(args ...any) Widget {
    cfg := &Config{}
    for _, a := range args {
        switch v := a.(type) {
        case Option:
            v.apply(cfg)
        case Widget:
            cfg.Children = append(cfg.Children, v)
        }
    }
    return &rowWidget{cfg: cfg}
}
```

牺牲类型安全换取 API 优雅度。可在 `go vet` 或 linter 中补充检查。

---

## 破局总结

| 约束 | 破局方案 |
|------|----------|
| 无 Proxy | 显式 `SetState()`，类似 React |
| 无 JSX | 函数嵌套 + `...any` 变参 |
| 无命名参数 | Functional Options Pattern |
| 泛型变参限制 | 运行时类型断言（顶层 API） |
| GC 开销 | Widget 描述对象极轻量，依赖 GC 无压力 |

**最终选择**：Flutter 三层树架构 + React 函数组件声明风格 + 显式 SetState。
