# 完整渲染管线

> 父文档：[ARCHITECTURE.md](./ARCHITECTURE.md)
> 前置阅读：[03-three-tree-model.md](./03-three-tree-model.md), [06-yoga-ebiten-integration.md](./06-yoga-ebiten-integration.md)

---

## 管线总览

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐     ┌─────────────┐
│   SetState  │────▶│   Rebuild    │────▶│  Reconcile  │────▶│ Yoga Layout │
│  (用户触发)  │     │  Build()调用  │     │  diff更新   │     │  计算约束    │
└─────────────┘     └──────────────┘     └─────────────┘     └─────────────┘
                                                                    │
                                                                    ▼
┌─────────────┐     ┌──────────────┐     ┌─────────────┐     ┌─────────────┐
│  下一帧循环  │◀────│  MarkDirty   │◀────│  Commit更新  │◀────│  SyncBounds │
│             │     │  标记脏节点   │     │  应用副作用   │     │  同步布局结果│
└─────────────┘     └──────────────┘     └─────────────┘     └─────────────┘

┌──────────────────────────────────────────────────────────────────────────┐
│                              Ebiten Draw()                                │
│  遍历 RenderNode Tree ──▶ Paint() ──▶ Ebiten DrawRect/DrawText/DrawImage │
└──────────────────────────────────────────────────────────────────────────┘
```

---

## 阶段一：Rebuild（构建）

触发条件：`SetState()` 或首次挂载。

```
ComponentElement.BuildIfNeeded()
  └── widget.Build(ctx) → newWidget
      └── 用户代码执行，生成 Widget 描述树
```

特点：
- 只执行脏组件及其祖先路径上的 `Build()`
- Widget 描述对象是轻量 immutable 值对象
- 不触及 Yoga Node 或 Ebiten 资源

---

## 阶段二：Reconcile（协调）

```
ReconcileChild(oldChild, newWidget)
  ├── canUpdate ? Update() : Unmount() + Mount()
  └── RenderObjectElement.Update()
      ├── SyncYogaProps()           // 更新 Yoga 样式
      └── reconcileChildren()       // 递归处理子节点
```

副作用：
- 新增 Element → 创建 Yoga Node
- 删除 Element → 释放 Yoga Node
- 更新 Element → 修改 Yoga 样式属性

---

## 阶段三：Yoga Layout（布局）

```
Engine.layout()
  └── rootRenderNode.YogaNode().CalculateLayout(width, height, DirectionLTR)
      └── Yoga 内部算法：
          1. 自顶向下传递约束（Constraints）
          2. 自底向上计算尺寸（MeasureFunc 回调）
          3. 自顶向下确定位置（Position）
```

同步结果：
```go
func syncBoundsFromYoga(node RenderNode) {
    yn := node.YogaNode()
    bounds := Rect{
        X: yn.LayoutLeft(),
        Y: yn.LayoutTop(),
        W: yn.LayoutWidth(),
        H: yn.LayoutHeight(),
    }
    node.SetBounds(bounds)
    for _, child := range node.Children() {
        syncBoundsFromYoga(child)
    }
}
```

---

## 阶段四：Ebiten Draw（绘制）

```
Engine.Draw(screen)
  └── paintTree(screen, rootRenderNode)
      └── for each node:
          ├── node.Paint(screen, bounds)
          │   ├── Box: vector.DrawFilledRect
          │   ├── Text: text/v2 Draw
          │   ├── Image: screen.DrawImage
          │   └── 其他: 自定义绘制
          └── paintTree(screen, child) // DFS 递归
```

绘制顺序：先父后子（DFS），子节点自然覆盖父节点（ painters algorithm ）。

---

## 输入事件管线

```
Ebiten Update()
  ├── 处理鼠标/键盘/触摸输入
  ├── HitTest(rootRenderNode, mousePos)
  │   └── DFS 遍历 RenderNode，检查 bounds 包含
  ├── 构建事件路径（冒泡链）
  └── 分发事件到对应 Widget 的 OnTap/OnPress 处理器
```

---

## 与当前 tenon 的管线对比

| 阶段 | 当前 tenon | 新架构 |
|------|-----------|--------|
| 状态触发 | Signal callback | `SetState()` |
| 构建 | 无（命令式直接改） | `Build()` 生成 Widget 树 |
| Diff | 无 | `Reconcile()` O(n) |
| 布局 | 每帧新建 Yoga Node | 保留 Yoga Node，脏时重新计算 |
| 绘制 | 遍历 Widget Tree | 遍历 RenderNode Tree（更浅） |
| 事件 | Widget 直接处理 | Hit-test on RenderNode bounds |
