# Reconcile Diff 算法

> 父文档：[ARCHITECTURE.md](./ARCHITECTURE.md)
> 前置阅读：[03-three-tree-model.md](./03-three-tree-model.md)

---

## 核心规则：`canUpdate`

借鉴 Flutter 的极简复用规则：

```go
func canUpdate(oldW, newW Widget) bool {
    if oldW == nil || newW == nil {
        return false
    }
    return oldW.WidgetType() == newW.WidgetType() && oldW.Key() == newW.Key()
}
```

- **相同 WidgetType + 相同 Key** → 复用 Element，只更新属性
- **不同** → 卸载旧子树，创建新 Element

> Key 的作用：在列表中保持元素身份稳定，避免状态丢失。

---

## 单节点 Reconcile

```go
func ReconcileChild(parent Element, oldChild Element, newWidget Widget) Element {
    // 新 Widget 为 nil，卸载旧节点
    if newWidget == nil {
        if oldChild != nil {
            oldChild.Unmount()
        }
        return nil
    }

    // 尝试复用旧节点
    if oldChild != nil && canUpdate(oldChild.Widget(), newWidget) {
        oldChild.Update(newWidget)
        return oldChild
    }

    // 无法复用：卸载旧节点，创建新节点
    if oldChild != nil {
        oldChild.Unmount()
    }
    newChild := createElement(newWidget)
    newChild.Mount(parent)
    return newChild
}
```

---

## 子节点列表 Reconcile（带 Key）

列表场景需要更复杂的 diff，类似 React 的 Key-based reconciliation：

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

        // 优先按 Key 匹配
        if k := w.Key(); k != "" {
            matched = oldKeyed[k]
            delete(oldKeyed, k)
        } else if unkeyedIdx < len(oldUnkeyed) {
            // 无 Key：按顺序匹配同类型
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

### 复杂度

- 时间：O(n)，单次遍历
- 空间：O(k)，k 为有 Key 的旧节点数

---

## Element 创建分发

```go
func createElement(w Widget) Element {
    if cw, ok := w.(ComponentWidget); ok {
        return &ComponentElement{widget: cw}
    }
    if rw, ok := w.(RenderWidget); ok {
        return &RenderObjectElement{widget: rw}
    }
    panic("unknown widget type")
}
```

---

## ComponentElement 的重建流程

```go
func (e *ComponentElement) rebuild() {
    // 1. 调用用户 Build() 生成新 Widget 描述
    newWidget := e.widget.Build(e)

    // 2. 对单根子节点做 Reconcile
    e.child = ReconcileChild(e, e.child, newWidget)
}
```

---

## RenderObjectElement 的更新流程

```go
func (e *RenderObjectElement) Update(newWidget Widget) {
    oldWidget := e.widget
    e.widget = newWidget.(RenderWidget)

    // 1. 更新 RenderNode 的属性
    e.renderNode.SyncYogaProps(e.widget)

    // 2. 对子节点列表做 Reconcile
    oldChildren := e.children
    newWidgets := e.widget.GetChildren()
    e.children = reconcileChildren(e, oldChildren, newWidgets)

    // 3. 同步 Yoga 子节点树（增删子节点）
    syncYogaChildren(e.renderNode, e.children)
}
```

---

## 与 React / Flutter 的对比

| 特性 | React | Flutter | 本架构 |
|------|-------|---------|--------|
| Diff 单元 | Fiber 节点 | Element | Element |
| 复用条件 | Key + type | `canUpdate` (type + Key) | `canUpdate` (WidgetType + Key) |
| 列表 Diff | 完整 Key diff | 完整 Key diff | 简化 Key diff |
| 副作用收集 | flags 链表 | Element 直接操作 | RenderNode 直接操作 |
| 时间复杂度 | O(n) | O(n) | O(n) |
