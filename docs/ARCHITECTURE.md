# Yoga + Ebiten 声明式 GUI 架构总纲

> 在 Go 中实现类似 React 的声明式组件 API，以 Yoga 为布局引擎、Ebiten 为渲染后端。

---

## 文档导航

本文档为父文档，各子主题详见对应子文档：

| 序号 | 文档 | 内容 |
|------|------|------|
| 1 | [01-research.md](./01-research.md) | React / Flutter / Vue Native 架构调研对比 |
| 2 | [02-core-constraints.md](./02-core-constraints.md) | Go 语言核心约束与破局思路 |
| 3 | [03-three-tree-model.md](./03-three-tree-model.md) | 三层树模型（Widget/Element/RenderNode） |
| 4 | [04-declarative-api.md](./04-declarative-api.md) | 用户层声明式 API 设计 |
| 5 | [05-reconciliation.md](./05-reconciliation.md) | Reconcile  diff 算法 |
| 6 | [06-yoga-ebiten-integration.md](./06-yoga-ebiten-integration.md) | Yoga 布局与 Ebiten 渲染集成 |
| 7 | [07-state-management.md](./07-state-management.md) | 状态管理（`SetState` 机制） |
| 8 | [08-component-library.md](./08-component-library.md) | 组件库设计规范 |
| 9 | [09-rendering-pipeline.md](./09-rendering-pipeline.md) | 完整渲染管线 |
| 10 | [10-implementation-roadmap.md](./10-implementation-roadmap.md) | 分阶段实现路径 |

---

## 架构全景图

```
┌─────────────────────────────────────────────────────────────┐
│  User Code (声明式) 详见 04-declarative-api.md               │
│  func (c *Counter) Build(ctx BuildContext) Widget {         │
│      return Column(Gap(16),                                 │
│          Text(fmt.Sprintf("Count: %d", c.count)),           │
│          Button(OnTap(c.SetState(func() { c.count++ })),    │
│              Text("+"),                                     │
│          ),                                                 │
│      )                                                      │
│  }                                                          │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼ Build()
┌─────────────────────────────────────────────────────────────┐
│  Widget Tree (Immutable 描述层)  详见 03-three-tree-model.md │
│  每次状态变化重建，极轻量（只有配置数据）                      │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼ Reconcile (canUpdate)
┌─────────────────────────────────────────────────────────────┐
│  Element Tree (Mutable 身份 + 状态层)  详见 03-three-tree-model.md │
│  跨重建复用，持有组件实例、Yoga Node 引用、state              │
│  ComponentElement: 管理子 Element                            │
│  RenderObjectElement: 连接 Yoga/Ebiten 渲染节点              │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼ Yoga CalculateLayout
┌─────────────────────────────────────────────────────────────┐
│  RenderNode Tree (Mutable 布局 + 绘制层)                     │
│  详见 06-yoga-ebiten-integration.md + 09-rendering-pipeline.md │
│  每个节点 = Yoga.Node + Ebiten 绘制命令                       │
│  Layout:  Yoga 约束计算 (O(n))                               │
│  Paint:   Ebiten DrawRect / DrawText / DrawImage            │
└─────────────────────────────────────────────────────────────┘
```

---

## 核心设计决策

### 1. 采用 Flutter 三层树 + React 声明语法

- **不要试图在 Go 中复制 Vue 的 Proxy 响应式**（详见 [02-core-constraints.md](./02-core-constraints.md)）
- Widget 层给用户一个干净、声明式的 API（函数嵌套模拟 JSX）
- Element 层解决"重建 Widget 但保留状态和 Yoga Node"的核心矛盾
- RenderNode 层把 Yoga 和 Ebiten 粘合在一起，只做布局+绘制

### 2. 状态更新用显式 `SetState()`

```go
// 包裹状态变更，框架在闭包执行后自动触发重建
OnTap(c.SetState(func() {
    c.count++
}))
```

详见 [07-state-management.md](./07-state-management.md)。

### 3. 组件粒度 ≠ 渲染粒度

用户组件可以任意嵌套，但只有"渲染 Widget"（Text、Box、Row 等）才创建 RenderNode 和 Yoga Node。

详见 [03-three-tree-model.md#为什么需要三层](./03-three-tree-model.md)。

---

## 与当前 tenon 架构的对比

| 维度 | 当前 tenon | 新架构 |
|------|-----------|--------|
| **组件模型** | 单层 Widget（命令式） | 三层树（Widget/Element/RenderNode） |
| **API 风格** | 命令式：`w.SetText("x")` | 声明式：`Build() Widget` |
| **状态更新** | Signal push + 直接 Mark | 显式 `SetState()` + Rebuild |
| **Yoga 使用** | 即时模式（每帧新建 Node） | 保留模式（Node 随 Element 复用） |
| **渲染粒度** | 每个 Widget 都 Layout+Draw | 用户组件不创建 RenderNode，只有叶子/容器创建 |
| **Ebiten 集成** | `App` 直接驱动 | `Game` 实现分 Update(Build+Layout) / Draw |

详见 [10-implementation-roadmap.md#与现有代码的迁移关系](./10-implementation-roadmap.md)。
