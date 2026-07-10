# 调研报告：React / Flutter / Vue Native 架构对比

> 父文档：[ARCHITECTURE.md](./ARCHITECTURE.md)

---

## 1. React Fiber

### 核心机制

| 维度 | 说明 |
|------|------|
| **声明层** | JSX → React Element（immutable 描述对象，轻量） |
| **Reconciler** | Fiber 树 = mutable 的工作单元，持有 state、props、hooks、DOM 引用 |
| **Diff** | 双缓冲（current / workInProgress），单链表 DFS，O(n) diff |
| **Commit** | 副作用（增删改 DOM）同步一次性执行 |
| **状态触发** | `setState` → schedule → re-render → diff → commit |
| **渲染粒度** | Virtual DOM 节点与实际 DOM 节点**不强制一一对应**（Portal/Fragment） |

### Fiber 节点结构

```ts
type Fiber = {
  tag: WorkTag           // 组件类型标记
  key: null | string
  elementType: any       // element.type
  type: any              // 解析后的函数/类
  stateNode: any         // DOM 实例或组件实例
  return: Fiber | null   // 父节点
  child: Fiber | null    // 第一个子节点
  sibling: Fiber | null  // 下一个兄弟节点
  pendingProps: any      // 新 props
  memoizedProps: any     // 当前 props
  memoizedState: any     // 当前 state（hooks 链表）
  alternate: Fiber | null // 双缓冲另一棵树
  flags: Flags           // 副作用标记
  // ...
}
```

### 对 Go 的启示

- 声明层用轻量 immutable 对象是可取的，Go 的 GC 适合管理短生命周期描述对象。
- React 的 VDOM full-diff 在 Go 中性能开销大，且没有内置 key 优化机制。
- Hooks 依赖数组在编译期无检查，Go 更难模拟。

---

## 2. Flutter 三棵树

### 核心机制

| 树 | 职责 | 可变性 | 生命周期 |
|----|------|--------|----------|
| **Widget** | 配置描述（Blueprint） | Immutable | 每次 setState 重建 |
| **Element** | 身份 + 状态 + 生命周期 | Mutable | 跨重建复用 |
| **RenderObject** | Layout + Paint + Hit-test | Mutable | 由 Element 管理 |

### 核心 `canUpdate` 规则

```dart
static bool canUpdate(Widget oldWidget, Widget newWidget) {
  return oldWidget.runtimeType == newWidget.runtimeType
      && oldWidget.key == newWidget.key;
}
```

- **相同类型 + 相同 Key** → 复用 Element 和 RenderObject，只更新属性
- **不同** → 卸载整个子树，重建

### 为什么需要三层？

> Widget 非常轻量，每次 `setState` 都重建 Widget 子树。Element 作为"持久身份层"决定是否复用。RenderObject 做真正昂贵的 layout 和 paint。

### 对 Go 的启示

- Flutter 的"组件粒度"（Widget 可以任意嵌套）与"渲染粒度"（RenderObject 只包含实际要画的节点）**明确分离**。
- `canUpdate` 规则极其简单，避免了 VDOM 全量 diff，**非常适合 Go 实现**。
- Element 层作为"可变身份层"，完美解决了"声明式重建"与"状态保持"的矛盾。

---

## 3. Vue Native / Lynx

| 方案 | 状态 | 核心机制 |
|------|------|----------|
| Vue Native | **已废弃** | Vue 2 + React Native 渲染器 |
| Lynx (Vue + Lynx) | 活跃（字节跳动） | 双线程架构（主线程 UI + 后台线程 JS），编译时拆分 template/script |

### 关键问题：Proxy 依赖

Vue 3 的响应式系统重度依赖 JavaScript `Proxy`：

```js
// Vue 3 自动追踪依赖
const count = ref(0)
// 在 render 中读取 count.value，Vue 自动建立订阅关系
```

**Go 没有等价物。**

- Vue 的 `ref/reactive` 在运行时自动追踪依赖。
- Lynx 编译时将 `<template>` 拆到主线程、`<script>` 拆到后台线程。

### 结论

Vue 的响应式模式在 Go 中**不可复制**，必须放弃。详见 [02-core-constraints.md](./02-core-constraints.md)。

---

## 4. 三种架构对比总结

| 特性 | React | Flutter | Vue/Lynx |
|------|-------|---------|----------|
| 声明层 | JSX / Element | Widget (immutable) | Template / SFC |
| 中间层 | Fiber (mutable) | Element (mutable) | Proxy 响应式系统 |
| 渲染层 | DOM / Native | RenderObject | DOM / Native |
| Diff 策略 | VDOM 全树 diff | `canUpdate` 类型+Key 匹配 | 细粒度依赖追踪 |
| 状态模型 | 显式 setState | 显式 setState | 自动追踪 (Proxy) |
| 跨平台 | React Native | 自绘 (Skia/Impeller) | Lynx 双线程 |
| Go 可行性 | 中（diff 开销大） | **高**（规则简单） | **低**（无 Proxy） |

### 选型结论

**采用 Flutter 的三层树模型 + React 的函数组件声明风格。**

理由：
1. `canUpdate` 复用规则比 VDOM diff 更简单、更快、更易在 Go 中实现。
2. 三层分离天然解决"声明式重建 vs 状态保持"问题。
3. 函数组件语法可以通过 Go 的嵌套函数调用来模拟。
