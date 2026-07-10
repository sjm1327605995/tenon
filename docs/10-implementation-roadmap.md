# 分阶段实现路径

> 父文档：[ARCHITECTURE.md](./ARCHITECTURE.md)

---

## Phase 1：核心骨架（4 周）

目标：能运行一个最小示例（Counter），验证三层树 + Ebiten 驱动可行性。

### 1.1 UI 核心 (`pkg/ui/`)

- [ ] `Widget` / `ComponentWidget` / `RenderWidget` 接口
- [ ] `Element` / `ComponentElement` / `RenderObjectElement`
- [ ] `BuildContext` 接口
- [ ] `Component` 基类（`SetState()`）
- [ ] `ReconcileChild()` + `reconcileChildren()`
- [ ] `createElement()` 工厂

### 1.2 渲染核心 (`pkg/ui/`)

- [ ] `RenderNode` 接口
- [ ] `BaseRenderNode`（Yoga Node 托管、Bounds 管理）
- [ ] `syncBoundsFromYoga()`

### 1.3 Ebiten 驱动 (`pkg/driver/`)

- [ ] `Engine`（实现 `ebiten.Game`）
- [ ] `EbitenCanvas`（DrawRect, DrawText, DrawImage）
- [ ] `Update()` 触发 Rebuild → Reconcile → Layout
- [ ] `Draw()` 遍历 RenderNode 绘制

### 1.4 基础组件 (`pkg/components/`)

- [ ] `Box`（背景色、padding、圆角）
- [ ] `Text`（Ebiten text/v2）
- [ ] `Row` / `Column`（Yoga Flex）
- [ ] `Button`（点击区域 + OnTap）

### 1.5 示例 (`example/`)

- [ ] Counter（+ / - 按钮，数字显示）

**验收标准**：`go run example/main.go` 弹出窗口，点击按钮数字变化。

---

## Phase 2：基础组件库（3 周）

- [ ] `Image`（加载 + 绘制）
- [ ] `Stack`（层叠 + 定位）
- [ ] `ScrollView`（滚动区域 + 滚动条）
- [ ] `Input`（文本输入 + 焦点）
- [ ] `Checkbox` / `Switch`
- [ ] `Divider`
- [ ] `If` / `ForEach` 辅助组件
- [ ] `Icon`（SVG/矢量图标支持）

**验收标准**：能用纯组件库拼出一个 Settings 页面。

---

## Phase 3：交互与动画（3 周）

- [ ] 完整事件系统（Hit-test → 冒泡 → 捕获）
- [ ] `GestureDetector`（Tap, Pan, LongPress, Scroll）
- [ ] 手势冲突解决
- [ ] 入场动画 `Animate()`
- [ ] 布局动画 `LayoutAnimation()`
- [ ] 过渡效果（Fade, Slide, Scale）

**验收标准**：实现一个可拖拽卡片，带入场动画。

---

## Phase 4：高级特性（4 周）

- [ ] `Key` 优化列表 reconcile（Virtual List / ListView 回收）
- [ ] `Fragment` / `Portal` 支持
- [ ] 主题系统 `ThemeProvider` + `Theme()`
- [ ] 依赖注入 / `InheritedWidget` 模式
- [ ] 与现有 `state.Signal` 集成（可选自动订阅）
- [ ] 多窗口 / 弹出层 / Dialog 支持
- [ ] 无障碍（Accessibility）基础

**验收标准**：实现一个深色模式切换 + 大型列表（1000 项）不卡顿。

---

## 与现有 tenon 代码的关系

| 现有代码 | 处理方式 |
|----------|----------|
| `pkg/yoga/` | **保留**，作为布局引擎核心 |
| `state/` (Signal) | **可选集成**，Phase 4 中对接 |
| `animation/` | **复用**，Phase 3 中对接 |
| `event/` | **参考**，Phase 3 中改造为 RenderNode Hit-test 模型 |
| `theme/` | **参考**，Phase 4 中重构为 ThemeProvider 模式 |
| `core/*` (Button, Dialog...) | **重写**，基于新组件 API 重新实现 |
| `widget/` (旧 base, lifecycle) | **废弃**，由新 `pkg/ui/` 替代 |
| `primitives/` (旧 Row, Column) | **废弃**，由新 `pkg/components/` 替代 |
| `app/app.go` | **废弃**，由新 `pkg/driver/ebiten.go` 替代 |

---

## 风险与应对

| 风险 | 影响 | 应对 |
|------|------|------|
| Ebiten text/v2 中文渲染 | 高 | Phase 1 即验证， fallback 到 image/font |
| Yoga  retained Node 内存泄漏 | 中 | 每个 Element.Unmount() 必须调用 yoga.Reset() |
| `...any` API 缺乏类型安全 | 低 | 补充 `go vet` 插件或文档约束 |
| 大量 Widget 重建的 GC 压力 | 中 | Phase 1 用 pprof 验证，必要时引入对象池 |
| 事件系统与 Yoga Hit-test 冲突 | 中 | Phase 3 预留 GestureArena 设计 |
