# Tenon 架构设计（Flutter 三棵树 + SwiftUI-like API）

## 一、设计哲学

- **参考 Flutter 三棵树**：严格分离 Widget（配置）、Element（状态/生命周期）、RenderObject（布局/绘制）。
- **SwiftUI-like 用户 API**：声明式函数、链式修饰符、不可变 Widget。
- **显式重建**：先不考虑自动数据绑定，用户通过 `engine.Rebuild()` 显式触发 UI 更新。
- **RenderObject 自治**：所有布局/视觉属性修改由 RenderObject 内部标记脏区，触发重布局/重绘制。

---

## 二、三棵树职责

| 树 | 职责 | 生命周期 | 对应包 |
|---|---|---|---|
| **Widget Tree** | 不可变的 UI 配置描述 | 极短，每次 rebuild 重建 | `pkg/v2/widgets/` |
| **Element Tree** | 可变状态持有者，连接 Widget ↔ RenderObject | 长，按 Key 复用 | `pkg/v2/ui/element.go` |
| **RenderObject Tree** | 布局计算 + 绘制原语 | 随 Element 同步创建/销毁 | `pkg/v2/render/` |

### Widget → Element → RenderObject 的映射

```
TextWidget ──CreateElement──► TextElement ──createRenderObject──► RenderText
     │                              │                                 │
     │  Update()                    │  updateRenderObject()           │ MarkNeedsLayout/Paint
     ▼                              ▼                                 ▼
TextWidget' ──复用 Element──► TextElement' ──同步属性──► RenderText'
```

---

## 三、目录结构

```
tenon/
├── tenon.go                   # 统一入口，导出 v2 类型和 Run()
├── pkg/
│   ├── v2/ui/                 # 三棵树核心
│   │   ├── widget.go          # Widget 接口、Key、BaseWidget
│   │   ├── element.go         # Element 接口、ComponentElement、RenderObjectElement
│   │   ├── engine.go          # Engine、PipelineOwner、帧循环
│   │   ├── element.go         # BuildFunc、updateChild、diff 逻辑（已合并至 element.go）
│   │   └── theme.go           # Theme
│   ├── v2/render/             # RenderObject 实现
│   │   ├── render_object.go   # RenderObject 基类、脏标记系统
│   │   ├── render_box.go      # RenderBox（基础矩形，内嵌 yoga.Node）
│   │   ├── render_text.go     # RenderText（文字测量与绘制）
│   │   ├── render_flex.go     # RenderFlex（Row/Column 布局）
│   │   └── render_painter.go  # 通用绘制辅助
│   ├── v2/widgets/            # 内置 Widget（SwiftUI-like API）
│   │   ├── text.go            # TextWidget
│   │   ├── flex.go            # RowWidget, ColumnWidget
│   │   ├── container.go       # ContainerWidget
│   │   └── button.go          # ButtonWidget
│   ├── fonts/                 # 字体管理器
│   └── svg/                   # SVG 工具
├── yoga/                      # Flexbox 布局引擎（独立模块，无改动）
├── example/
│   ├── gallery/               # 组件画廊示例
│   └── v2-demo/               # v2 架构示例
└── _debug_snapshot/           # 无 GUI 测试（test_table.go）
```

---

## 四、核心流程

### 4.1 启动流程

```
用户调用 tenon.Run(buildFunc, width, height)
  → ui.NewEngine(buildFunc, width, height)
  → engine.Mount()
      → buildFunc() 产出 root Widget
      → rootWidget.CreateElement() 产出 root Element
      → rootElement.Mount(nil, 0)
          → Element.Mount 递归创建子 Element
          → RenderObjectElement.Mount 创建 RenderObject
          → 子 RenderObject 挂载到父 RenderObject 树
      → rootRenderObject.Attach(PipelineOwner)
  → ebiten.RunGame(engine)
```

### 4.2 帧循环

```
Engine.Update()
  ├── 1. flushBuild()          # Widget diff（如 needsBuild）
  ├── 2. Yoga 全局布局计算    # 根节点 CalculateLayout
  ├── 3. PipelineOwner.FlushLayout()  # RenderObject.PerformLayout
  ├── 4. syncBounds()          # Yoga 结果 → RenderObject.Bounds
  ├── 5. PipelineOwner.FlushPaint()   # 清绘制脏标记
  └── 6. handleMouseInput()    # 鼠标事件、点击回调

Engine.Draw()
  └── 遍历 RenderObject 树，调用 Paint(screen, offset)
```

### 4.3 Widget Diff（updateChild / UpdateChildren）

```
CanUpdate(oldWidget, newWidget)
  ├── 同 runtimeType + 同 Key → element.Update(newWidget) 复用
  └── 不同 → unmount 旧 Element，mount 新 Element

Element.Update(newWidget)
  ├── 更新自身 widget 引用
  ├── RenderObjectElement: updateRenderObject(oldWidget) 同步属性
  └── ComponentElement: 递归 diff 子树
```

---

## 五、用户 API 风格

```go
func MyApp() tenon.Widget {
    return tenon.Column(
        tenon.Text("Hello").FontSize(24).Color(tenon.GetTheme().TextColor),
        tenon.Row(
            tenon.Button("Click").OnTap(func() {
                counter++
                engine.Rebuild()
            }),
        ).Gapf(8),
    ).Gapf(16).Paddingf(tenon.EdgeInsetsAll(24))
}

func main() {
    tenon.Run(MyApp, 900, 600)
}
```

- **不可变性**：所有链式方法返回新 Widget 实例（值拷贝）。
- **无数据绑定**：状态变化后显式调用 `engine.Rebuild()`。
- **声明式**：UI 由纯函数描述，每次重建重新执行函数产出新 Widget 树。

---

## 六、开发规范

### 6.1 新增 Widget 的步骤

1. 在 `pkg/v2/widgets/` 中定义 `XxxWidget` 结构体（内嵌 `ui.BaseWidget`）。
2. 提供构造函数（如 `Xxx(...) XxxWidget`）和链式修饰方法。
3. 实现 `CreateElement() ui.Element`。
4. 定义 `XxxElement`，根据需求继承：
   - `RenderObjectElement`：单孩子 + RenderObject（如 Container）
   - `SingleChildRenderObjectElement`：单孩子 + RenderObject（如 Padding）
   - `MultiChildRenderObjectElement`：多孩子 + RenderObject（如 Row/Column）
   - `SingleChildComponentElement`：单孩子，无 RenderObject（如条件包装器）
5. 覆盖 `CreateRenderObject()` 和 `UpdateRenderObject()`。
6. 在 `Mount` 中确保 `e.RenderObject = e.CreateRenderObject()` 在调用父类 `Mount` 之前执行（Go 嵌入方法不触发动态分发）。

### 6.2 新增 RenderObject 的步骤

1. 在 `pkg/v2/render/` 中定义 `RenderXxx`，内嵌 `BaseRenderObject` 或 `RenderBox`。
2. 在 `Init(self)` 中创建 `yoga.NewNode()`（如果需要 Flexbox 布局）。
3. 覆盖 `PerformLayout()`（如有自定义布局逻辑）和 `Paint()`。
4. 属性 setter 内部调用 `MarkNeedsLayout()` / `MarkNeedsPaint()`。

---

## 七、与旧架构的区别

| 维度 | 旧架构 | 新架构 |
|---|---|---|
| 树模型 | Element + Yoga 强耦合 | Widget + Element + RenderObject 三棵树 |
| 用户 API | 命令式链式 Setter（`SetWidth`, `Add`） | 声明式不可变 Widget + 链式修饰符 |
| 状态更新 | `State[T]` 自动追踪 + `RequestBuild()` | 无自动绑定，显式 `Rebuild()` |
| 组件定义 | 直接创建 Element，命令式操作 | Widget 描述配置，Element 管理生命周期 |
| 布局职责 | BaseElement 内嵌 Yoga 节点 | RenderObject 内嵌 Yoga 节点，Element 不碰布局 |
