# 声明式 API 设计

> 父文档：[ARCHITECTURE.md](./ARCHITECTURE.md)
> 前置阅读：[02-core-constraints.md](./02-core-constraints.md)

---

## 设计原则

1. **Children 永远在最后**，用变参接收，支持任意数量
2. **属性用函数选项**（Functional Options），模拟 JSX 的 named props
3. **CSS 效果原子化**，每个视觉属性都是一个 option 函数
4. **布局即样式**，Flex 属性直接挂在容器上
5. **条件/列表用辅助函数**，避免 Go 的 `if` 嵌套过深

---

## 核心 trick：`...any` 变参

Go 没有泛型变参联合类型，顶层 API 用 `...any`，运行时区分 `Option` 和 `Widget`：

```go
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

用户写法：

```go
Row(
    Gap(16),
    AlignItems(AlignCenter),
    Text("Hello"),   // Widget 自动识别为 Children
    Text("World"),
)
```

---

## 基础组件

### 文本

```go
// 最简单
Text("Hello World")

// 带样式
Text("Hello World",
    FontSize(24),
    FontWeightBold,
    Color(ColorHex("#1a1a1a")),
    LineHeight(1.5),
    MaxLines(2),
    TextAlignCenter,
)

// 富文本（内联样式）
RichText(
    Span("Hello ", FontSize(18)),
    Span("World", FontSize(18), Color(ColorRed), FontWeightBold),
)
```

### 盒子（View / Div）

```go
Box(
    // 尺寸
    Width(200),
    HeightAuto,
    MinHeight(100),

    // 背景
    Background(ColorWhite),
    Background(LinearGradient(
        ColorHex("#667eea"), ColorHex("#764ba2"),
        Direction(135),
    )),

    // 边框
    Border(1, ColorHex("#e5e7eb")),
    BorderLeft(4, ColorHex("#3b82f6")),
    BorderRadius(12),
    BorderTopLeftRadius(0),

    // 阴影
    BoxShadow(
        Offset(0, 4),
        Blur(12),
        Spread(0),
        ColorRGBA(0, 0, 0, 0.08),
    ),

    // 内部布局
    Flex(Column),
    Gap(8),
    AlignItems(AlignStart),
    JustifyContent(JustifyCenter),

    // 变换
    Transform(
        Translate(0, -2),
        Rotate(5),
        Scale(1.02),
        Origin(Center),
    ),

    // 其他
    Opacity(0.95),
    OverflowHidden,
    Clip(ClipRounded(12)),

    // Children
    Text("Card Title", FontSize(18), FontWeightBold),
    Text("Card description goes here.", FontSize(14), Color(ColorGray)),
)
```

### 图片

```go
Image("avatar.png",
    Width(48),
    Height(48),
    BorderRadius(24),        // 圆形
    ObjectFitCover,
    Border(2, ColorWhite),
    BoxShadow(Offset(0, 2), Blur(4), ColorRGBA(0,0,0,0.1)),
)
```

---

## 布局系统

### Row（水平排列）

```go
Row(
    // Flex 容器属性
    Gap(16),
    AlignItems(AlignCenter),
    JustifyContent(JustifyBetween),

    // 自身样式
    WidthPercent(100),
    Padding(16, 20),
    Background(ColorHex("#f8fafc")),

    // Children
    Avatar("user.png"),
    Column(
        Gap(4),
        FlexGrow(1),
        Text("User Name", FontSize(16), FontWeightBold),
        Text("Online", FontSize(12), Color(ColorGreen)),
    ),
    Icon(ChevronRight, Color(ColorGray)),
)
```

### Column（垂直排列）

```go
Column(
    Gap(24),
    Width(360),
    Padding(24),
    Background(ColorWhite),
    BorderRadius(16),
    BoxShadow(Offset(0, 8), Blur(24), ColorRGBA(0,0,0,0.12)),

    Text("Sign In", FontSize(28), FontWeightBold, TextAlignCenter),

    Column(
        Gap(12),
        Input("Email", Placeholder("you@example.com")),
        Input("Password", Placeholder("••••••••"), SecureTextEntry),
    ),

    Button(
        WidthPercent(100),
        Padding(Vertical(14)),
        Background(ColorHex("#3b82f6")),
        BorderRadius(8),
        Active(Background(ColorHex("#2563eb"))),
        Text("Continue", Color(ColorWhite), TextAlignCenter, FontWeightBold),
    ),

    Row(
        AlignItems(AlignCenter),
        Gap(12),
        Divider(ColorHex("#e5e7eb")),
        Text("OR", FontSize(12), Color(ColorGray)),
        Divider(ColorHex("#e5e7eb")),
    ),
)
```

### Stack（层叠布局）

```go
Stack(
    Width(300), Height(200),

    // 底层：背景图
    Image("cover.jpg", ObjectFitCover),

    // 中层：渐变遮罩
    Box(
        Position(Fill),
        Background(LinearGradient(
            ColorTransparent, ColorRGBA(0,0,0,0.7),
            DirectionBottom,
        )),
    ),

    // 顶层：文字内容
    Column(
        Position(BottomLeft),
        Padding(16),
        Gap(4),
        Text("Article Title", FontSize(20), Color(ColorWhite), FontWeightBold),
        Text("5 min read", FontSize(12), Color(ColorWhite).Alpha(0.8)),
    ),
)
```

---

## 交互组件

### Button

```go
Button(
    Padding(12, 24),
    Background(ColorHex("#3b82f6")),
    BorderRadius(8),
    Gap(8),
    AlignItems(AlignCenter),

    // 交互状态
    Hover(
        Background(ColorHex("#2563eb")),
        Transform(Translate(0, -1)),
        BoxShadow(Offset(0, 4), Blur(8), ColorRGBA(59,130,246,0.3)),
    ),
    Active(
        Background(ColorHex("#1d4ed8")),
        Transform(Scale(0.98)),
    ),
    Disabled(
        Background(ColorHex("#e5e7eb")),
        Opacity(0.6),
    ),

    // 事件
    OnTap(c.SetState(func() { c.handleSubmit() })),

    // Content
    Icon(Send, Color(ColorWhite), Size(16)),
    Text("Send Message", Color(ColorWhite), FontWeightBold),
)
```

### Input

```go
Input(
    Value(c.searchText),
    Placeholder("Search..."),
    WidthPercent(100),
    Padding(12, 16),
    Background(ColorHex("#f3f4f6")),
    BorderRadius(8),
    FontSize(16),

    Focus(
        Background(ColorWhite),
        Border(2, ColorHex("#3b82f6")),
        BoxShadow(Offset(0, 0), Blur(0), Spread(3), ColorHex("#bfdbfe")),
    ),

    OnChange(c.SetState(func(v string) {
        c.searchText = v
    })),

    Prefix(Icon(Search, Color(ColorGray))),
)
```

---

## 条件与列表

### 条件渲染

```go
func (c *Page) Build() Widget {
    return Column(
        Navbar(),
        If(c.isLoading, Spinner(Center, Size(40))),
        If(!c.isLoading && c.hasError, ErrorView(c.err)),
        If(!c.isLoading && !c.hasError, ContentView(c.data)),
    )
}
```

### 列表渲染

```go
func (c *TodoApp) Build() Widget {
    return Column(
        ForEach(c.todos, func(todo Todo, index int) Widget {
            return TodoItem(
                Key(todo.ID),
                todo,
            )
        }),
    )
}
```

---

## 与 JSX 的映射关系

| 概念 | JSX (React) | Go API |
|------|-------------|--------|
| 元素 | `<Text>Hello</Text>` | `Text("Hello")` |
| 属性 | `<Box bg="red" p={16}>` | `Box(Background(ColorRed), Padding(16))` |
| Children | `<Col><A/><B/></Col>` | `Column(A(), B())` |
| 条件 | `{isOpen && <Modal/>}` | `If(isOpen, Modal())` |
| 列表 | `{items.map(x => <Item key={x.id} />)}` | `ForEach(items, func(x Item) Widget { return Item(Key(x.ID)) })` |
| 事件 | `onClick={() => ...}` | `OnTap(c.SetState(func() { ... }))` |
| 伪类 | `:hover`, `:active` | `Hover(...)` / `Active(...)` options |
| 样式类 | `className="card active"` | 无 className，全部内联原子属性 |
