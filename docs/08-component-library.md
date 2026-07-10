# 组件库设计规范

> 父文档：[ARCHITECTURE.md](./ARCHITECTURE.md)
> 前置阅读：[04-declarative-api.md](./04-declarative-api.md)

---

## 组件分类

| 分类 | 职责 | 示例 |
|------|------|------|
| **布局组件** | 组织子节点位置 | Row, Column, Stack, ScrollView |
| **渲染组件** | 产生实际像素 | Text, Image, Box |
| **交互组件** | 响应用户输入 | Button, Input, Checkbox, Slider |
| **容器组件** | 视觉包装 | Card, Badge, Divider |
| **辅助组件** | 条件/列表/动画 | If, ForEach, Animate |

---

## 布局组件

### Row

```go
func Row(args ...any) Widget
```

属性：
- `Gap(float32)` — 子元素间距
- `AlignItems(Align)` — 交叉轴对齐
- `JustifyContent(Justify)` — 主轴分布
- `FlexGrow(float32)` — 占据剩余空间

### Column

```go
func Column(args ...any) Widget
```

属性同 Row，默认 `FlexDirectionColumn`。

### Stack

```go
func Stack(args ...any) Widget
```

属性：
- `Position(Position)` — 子元素定位方式（Fill, TopLeft, BottomRight 等）

### ScrollView

```go
func ScrollView(args ...any) Widget
```

属性：
- `Direction(ScrollDirection)` — Vertical / Horizontal
- `ShowScrollIndicator(bool)`
- `AlwaysBounce(bool)`

---

## 渲染组件

### Text

```go
func Text(content string, args ...any) Widget
```

属性：
- `FontSize(float32)`
- `FontWeight(FontWeight)` — Normal, Bold, Medium
- `Color(Color)`
- `LineHeight(float32)`
- `MaxLines(int)`
- `TextAlign(TextAlign)` — Left, Center, Right

### Image

```go
func Image(src string, args ...any) Widget
```

属性：
- `Width(float32)`, `Height(float32)`
- `ObjectFit(ObjectFit)` — Cover, Contain, Fill
- `BorderRadius(float32)`

### Box

```go
func Box(args ...any) Widget
```

通用容器，支持所有视觉属性。

---

## 交互组件

### Button

```go
func Button(args ...any) Widget
```

交互状态：
- `Hover(...Option)` — 鼠标悬停样式
- `Active(...Option)` — 按下样式
- `Disabled(...Option)` — 禁用样式

事件：
- `OnTap(func())`
- `OnLongPress(func())`

### Input

```go
func Input(placeholder string, args ...any) Widget
```

状态：
- `Focus(...Option)` — 聚焦样式
- `Value(string)`
- `Placeholder(string)`

事件：
- `OnChange(func(string))`
- `OnSubmit(func())`

---

## 辅助组件

### If（条件渲染）

```go
func If(cond bool, w Widget) Widget
```

使用：

```go
If(c.isLoading, Spinner())
```

### ForEach（列表渲染）

```go
func ForEach[T any](items []T, fn func(item T, index int) Widget) []Widget
```

使用：

```go
Column(
    ForEach(c.todos, func(todo Todo, i int) Widget {
        return TodoItem(Key(todo.ID), todo)
    })...,
)
```

> 注意：`ForEach` 返回 `[]Widget`，需要用 `...` 展开传入 Column。

---

## 视觉属性清单

所有组件支持以下通用视觉属性：

| 属性 | 类型 | 说明 |
|------|------|------|
| `Width` | float32 / Percent | 宽度 |
| `Height` | float32 / Percent | 高度 |
| `MinWidth/MinHeight` | float32 | 最小尺寸 |
| `MaxWidth/MaxHeight` | float32 | 最大尺寸 |
| `Padding` | Insets / float32 | 内边距 |
| `Margin` | Insets / float32 | 外边距 |
| `Background` | Color / Gradient | 背景 |
| `BorderWidth` | float32 | 边框宽度 |
| `BorderColor` | Color | 边框颜色 |
| `BorderRadius` | float32 | 圆角 |
| `Opacity` | float32 | 不透明度 |
| `Transform` | Transform | 变换 |
| `BoxShadow` | Shadow | 阴影 |
| `Clip` | Clip | 裁剪 |
| `Overflow` | Overflow | 溢出处理 |
