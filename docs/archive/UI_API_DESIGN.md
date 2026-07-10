# 声明式 UI API 设计 —— Go 中的 "JSX 体验"

> 目标：在 Go 的语法限制下，提供尽可能接近 React/JSX 的声明式体验。
> 核心设计：**Options Pattern（函数式属性）+ 变参 Children**

---

## 一、设计原则

1. **Children 永远在最后**，用变参 `...Widget` 接收，支持任意数量
2. **属性用函数选项**（Functional Options），模拟 JSX 的 named props
3. **CSS 效果原子化**，每个视觉属性都是一个 option 函数
4. **布局即样式**，Flex 属性直接挂在容器上
5. **条件/列表用辅助函数**，避免 Go 的 `if` 嵌套过深

---

## 二、基础组件

### 2.1 文本

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

### 2.2 盒子（View / Div）

```go
// 纯容器
Box(
    Padding(16),
    Text("Content"),
)

// 视觉丰富的卡片
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
    BorderLeft(4, ColorHex("#3b82f6")), // 左边强调线
    BorderRadius(12),
    BorderTopLeftRadius(0),              // 可单独覆盖

    // 阴影
    BoxShadow(
        Offset(0, 4),
        Blur(12),
        Spread(0),
        ColorRGBA(0, 0, 0, 0.08),
    ),
    BoxShadowInset( // 内阴影
        Offset(0, 1),
        Blur(0),
        ColorRGBA(255, 255, 255, 0.5),
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

### 2.3 图片

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

## 三、布局系统

### 3.1 Row（水平排列）

```go
Row(
    // Flex 容器属性
    Gap(16),                          // 子元素间距
    AlignItems(AlignCenter),          // 交叉轴对齐
    JustifyContent(JustifyBetween),   // 主轴分布

    // 自身样式
    WidthPercent(100),
    Padding(16, 20),                  // H, V
    Background(ColorHex("#f8fafc")),

    // Children
    Avatar("user.png"),
    Column(
        Gap(4),
        FlexGrow(1),                  // 占据剩余空间
        Text("User Name", FontSize(16), FontWeightBold),
        Text("Online", FontSize(12), Color(ColorGreen)),
    ),
    Icon(ChevronRight, Color(ColorGray)),
)
```

### 3.2 Column（垂直排列）

```go
Column(
    Gap(24),
    Width(360),
    Padding(24),
    Background(ColorWhite),
    BorderRadius(16),
    BoxShadow(Offset(0, 8), Blur(24), ColorRGBA(0,0,0,0.12)),

    // Header
    Text("Sign In", FontSize(28), FontWeightBold, TextAlignCenter),

    // Form
    Column(
        Gap(12),
        Input("Email",
            Placeholder("you@example.com"),
            KeyboardType(Email),
        ),
        Input("Password",
            Placeholder("••••••••"),
            SecureTextEntry,
        ),
    ),

    // Actions
    Button(
        WidthPercent(100),
        Padding(Vertical(14)),
        Background(ColorHex("#3b82f6")),
        BorderRadius(8),
        Active(Background(ColorHex("#2563eb"))),  // :active 状态

        Text("Continue", Color(ColorWhite), TextAlignCenter, FontWeightBold),
    ),

    // Divider
    Row(
        AlignItems(AlignCenter),
        Gap(12),
        Divider(ColorHex("#e5e7eb")),   // 横线
        Text("OR", FontSize(12), Color(ColorGray)),
        Divider(ColorHex("#e5e7eb")),
    ),

    // Social Login
    Row(
        Gap(12),
        JustifyContent(JustifyCenter),
        SocialButton(Google),
        SocialButton(Apple),
        SocialButton(GitHub),
    ),
)
```

### 3.3 Stack（层叠布局）

```go
Stack(
    Width(300), Height(200),

    // 底层：背景图
    Image("cover.jpg", ObjectFitCover),

    // 中层：渐变遮罩
    Box(
        Position(Fill),           // 填满 Stack
        Background(LinearGradient(
            ColorTransparent, ColorRGBA(0,0,0,0.7),
            DirectionBottom,
        )),
    ),

    // 顶层：文字内容
    Column(
        Position(BottomLeft),     // Stack 内定位
        Padding(16),
        Gap(4),
        Text("Article Title", FontSize(20), Color(ColorWhite), FontWeightBold),
        Text("5 min read", FontSize(12), Color(ColorWhite).Alpha(0.8)),
    ),
)
```

### 3.4 ScrollView

```go
ScrollView(
    Direction(Vertical),
    ShowScrollIndicator(false),
    AlwaysBounce(false),

    Column(
        Gap(16),
        Padding(16),
        // 大量内容...
        ForEach(cards, func(card Card) Widget {
            return CardWidget(card)
        }),
    ),
)
```

---

## 四、交互组件

### 4.1 Button

```go
Button(
    // 样式
    Padding(12, 24),
    Background(ColorHex("#3b82f6")),
    BorderRadius(8),
    Gap(8),                       // 图标和文字间距
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
    OnTap(func() { c.handleSubmit() }),
    OnLongPress(func() { c.showMenu() }),

    // Content
    Icon(Send, Color(ColorWhite), Size(16)),
    Text("Send Message", Color(ColorWhite), FontWeightBold),
)
```

### 4.2 Input / TextField

```go
Input(
    Value(c.searchText),
    Placeholder("Search..."),

    // 样式
    WidthPercent(100),
    Padding(12, 16),
    Background(ColorHex("#f3f4f6")),
    BorderRadius(8),
    FontSize(16),

    // Focus 状态
    Focus(
        Background(ColorWhite),
        Border(2, ColorHex("#3b82f6")),
        BoxShadow(Offset(0, 0), Blur(0), Spread(3), ColorHex("#bfdbfe")), // ring
    ),

    // 事件
    OnChange(func(v string) {
        c.searchText = v
        c.SetState()
    }),
    OnSubmit(func() { c.doSearch() }),

    // 前缀/后缀
    Prefix(Icon(Search, Color(ColorGray))),
    Suffix(If(c.searchText != "", 
        Button(OnTap(c.clearSearch), Icon(X, Color(ColorGray))),
    )),
)
```

### 4.3 Checkbox / Switch

```go
Row(
    Gap(12),
    AlignItems(AlignCenter),
    Checkbox(
        Checked(c.agreed),
        OnChange(func(v bool) {
            c.agreed = v
            c.SetState()
        }),
        Size(20),
        Color(ColorHex("#3b82f6")),
    ),
    Text("I agree to the Terms of Service"),
)
```

---

## 五、状态与逻辑

### 5.1 有状态组件

```go
type Counter struct {
    count int
}

func (c *Counter) Build() Widget {
    return Box(
        Padding(24),
        Background(ColorHex("#f0f9ff")),
        BorderRadius(12),

        Row(
            Gap(16),
            AlignItems(AlignCenter),

            Button(
                Width(40), Height(40),
                Background(ColorHex("#e0f2fe")),
                BorderRadius(20),
                OnTap(func() {
                    c.count--
                    c.SetState()
                }),
                Text("-", FontSize(20), Color(ColorHex("#0284c7"))),
            ),

            Text(fmt.Sprintf("%d", c.count),
                FontSize(32),
                FontWeightBold,
                Width(60),
                TextAlignCenter,
                Color(If(c.count > 0, ColorHex("#16a34a"), ColorHex("#dc2626"))),
            ),

            Button(
                Width(40), Height(40),
                Background(ColorHex("#0284c7")),
                BorderRadius(20),
                OnTap(func() {
                    c.count++
                    c.SetState()
                }),
                Text("+", FontSize(20), Color(ColorWhite)),
            ),
        ),
    )
}
```

### 5.2 条件渲染

```go
func (c *Page) Build() Widget {
    return Column(
        // 顶部导航（始终显示）
        Navbar(),

        // 条件内容
        Switch(c.state,
            Case("loading", Spinner(Center, Size(40))),
            Case("error", ErrorView(c.err)),
            Case("empty", EmptyView("No data yet")),
            Default(ContentView(c.data)),
        ),

        // 或者用更 Go 风格
        c.buildContent(),
    )
}

func (c *Page) buildContent() Widget {
    switch c.state {
    case "loading":
        return Box(
            FlexGrow(1),
            AlignItems(AlignCenter),
            JustifyContent(JustifyCenter),
            Spinner(Size(40), Color(ColorPrimary)),
        )
    case "error":
        return Box(
            FlexGrow(1),
            AlignItems(AlignCenter),
            JustifyContent(JustifyCenter),
            Gap(12),
            Icon(AlertCircle, Color(ColorRed), Size(48)),
            Text(c.err.Error(), Color(ColorRed)),
            Button(OnTap(c.retry), Text("Retry")),
        )
    default:
        return ListView(...)
    }
}
```

### 5.3 列表渲染

```go
func (c *TodoApp) Build() Widget {
    return Column(
        ForEach(c.todos, func(todo Todo, index int) Widget {
            return TodoItem(
                Key(todo.ID),              // 列表必须带 Key
                todo,
                OnToggle(c.toggleTodo),
                OnDelete(c.deleteTodo),
            )
        }),
    )
}

// TodoItem 组件
func TodoItem(todo Todo, opts ...Option) Widget {
    return Row(
        Key(todo.ID),                      // 透传 Key
        Gap(12),
        AlignItems(AlignCenter),
        Padding(12, 16),
        Background(If(todo.Done, ColorHex("#f0fdf4"), ColorWhite)),
        BorderRadius(8),

        // 完成时文字带删除线
        Text(todo.Text,
            FlexGrow(1),
            FontSize(16),
            TextDecoration(If(todo.Done, LineThrough, None)),
            Color(If(todo.Done, ColorHex("#86efac"), ColorHex("#1f2937"))),
        ),

        Checkbox(Checked(todo.Done), OnChange(opts.OnToggle)),
        Button(OnTap(opts.OnDelete), Icon(Trash2, Color(ColorRed))),
    )
}
```

### 5.4 动画

```go
Box(
    // 入场动画
    Animate(
        Duration(300),
        Ease(OutCubic),
        From(Opacity(0), TranslateY(20)),
        To(Opacity(1), TranslateY(0)),
    ),

    // 布局动画（位置/尺寸变化时自动过渡）
    LayoutAnimation(
        Duration(250),
        Ease(InOutQuad),
    ),

    // 手势驱动动画
    Gesture(
        OnPan(func(gesture PanGesture) {
            c.offset = gesture.Translation
            c.SetState()
        }),
        OnPanEnd(func(gesture PanGesture) {
            if gesture.Velocity.X > 500 {
                c.dismiss()
            } else {
                c.offset = Point{0, 0} // spring back
                c.SetState()
            }
        }),
    ),

    Transform(Translate(c.offset.X, c.offset.Y)),

    Text("Swipe me"),
)
```

---

## 六、主题与上下文

```go
// 根组件注入主题
func (c *App) Build() Widget {
    return ThemeProvider(c.theme,
        Box(
            Background(Theme("bg.primary")),
            TextColor(Theme("text.primary")),
            FontFamily(Theme("font.sans")),

            Router(
                Route("/", HomePage),
                Route("/profile", ProfilePage),
            ),
        ),
    )
}

// 子组件消费主题
func Card() Widget {
    return Box(
        Background(Theme("surface")),           // 自动响应主题切换
        BorderRadius(Theme("radius.lg")),
        BoxShadow(Theme("shadow.md")),
        // ...
    )
}

// 暗黑模式切换
Button(
    OnTap(func() {
        c.isDark = !c.isDark
        c.SetState()
    }),
    Icon(If(c.isDark, Sun, Moon)),
)
```

---

## 七、完整页面示例

```go
type ChatApp struct {
    messages []Message
    input    string
}

func (c *ChatApp) Build() Widget {
    return Box(
        WidthPercent(100),
        HeightPercent(100),
        Background(ColorHex("#f3f4f6")),
        Flex(Column),

        // 顶部栏
        Box(
            Height(56),
            Padding(Horizontal(16)),
            Background(ColorWhite),
            BorderBottom(1, ColorHex("#e5e7eb")),
            Flex(Row),
            AlignItems(AlignCenter),
            JustifyContent(JustifyBetween),

            Row(
                Gap(12),
                AlignItems(AlignCenter),
                Avatar("bot.png", Size(36)),
                Text("AI Assistant", FontWeightBold),
            ),
            Icon(MoreVertical, Color(ColorGray)),
        ),

        // 消息列表
        ScrollView(
            FlexGrow(1),
            Reverse(true),                    // 底部对齐
            Padding(16),
            Gap(12),

            ForEach(c.messages, func(msg Message) Widget {
                isMe := msg.From == "user"
                return Row(
                    Key(msg.ID),
                    Gap(8),
                    AlignItems(AlignEnd),
                    If(isMe, FlexDirection(RowReverse)),

                    Avatar(msg.Avatar, Size(32)),

                    Box(
                        MaxWidth(280),
                        Padding(12, 16),
                        BorderRadius(16),
                        BorderTopRightRadius(If(isMe, 4, 16)),
                        BorderTopLeftRadius(If(isMe, 16, 4)),
                        Background(If(isMe, ColorHex("#3b82f6"), ColorWhite)),
                        BoxShadow(Offset(0, 1), Blur(2), ColorRGBA(0,0,0,0.05)),

                        Text(msg.Text,
                            Color(If(isMe, ColorWhite, ColorHex("#1f2937"))),
                            FontSize(15),
                            LineHeight(1.4),
                        ),
                        Text(msg.Time,
                            MarginTop(4),
                            FontSize(11),
                            Color(If(isMe, ColorRGBA(255,255,255,0.7), ColorGray)),
                        ),
                    ),
                )
            }),
        ),

        // 输入区
        Box(
            Padding(12, 16),
            Background(ColorWhite),
            BorderTop(1, ColorHex("#e5e7eb")),
            Flex(Row),
            Gap(8),
            AlignItems(AlignCenter),

            Input(
                FlexGrow(1),
                Value(c.input),
                Placeholder("Type a message..."),
                Background(ColorHex("#f3f4f6")),
                BorderRadius(20),
                Padding(10, 16),
                MaxHeight(120),
                OnChange(func(v string) {
                    c.input = v
                    c.SetState()
                }),
                OnSubmit(c.send),
            ),

            Button(
                Width(40), Height(40),
                BorderRadius(20),
                Background(If(c.input == "", ColorHex("#e5e7eb"), ColorHex("#3b82f6"))),
                Disabled(c.input == ""),
                OnTap(c.send),
                Icon(Send,
                    Color(If(c.input == "", ColorGray, ColorWhite)),
                    Size(18),
                ),
            ),
        ),
    )
}

func (c *ChatApp) send() {
    if c.input == "" { return }
    c.messages = append(c.messages, Message{
        ID:   uuid(),
        From: "user",
        Text: c.input,
        Time: now(),
    })
    c.input = ""
    c.SetState()
    // 触发 AI 回复...
}
```

---

## 八、API 风格总结

| 概念 | JSX (React) | Go API (本方案) |
|------|-------------|----------------|
| 元素 | `<Text>Hello</Text>` | `Text("Hello")` |
| 属性 | `<Box bg="red" p={16}>` | `Box(Background(ColorRed), Padding(16), ...)` |
| Children | `<Col><A/><B/></Col>` | `Column(A(), B())` |
| 条件 | `{isOpen && <Modal/>}` | `If(isOpen, Modal())` |
| 列表 | `{items.map(x => <Item key={x.id} />)}` | `ForEach(items, func(x Item) Widget { return Item(Key(x.ID)) })` |
| 事件 | `onClick={() => ...}` | `OnTap(func() { ... })` |
| 状态 | `const [n, setN] = useState(0)` | `c.n++; c.SetState()` |
| 样式类 | `className="card active"` | 无 className，全部内联原子属性 |
| 伪类 | `:hover`, `:active` | `Hover(...)` / `Active(...)` options |
