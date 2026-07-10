# 组件库实现方案

> 调研范围：Material UI (MUI)、Ant Design、Arco Design、Element Plus、Flutter Widget Catalog、Chakra UI
>
> 父文档：[ARCHITECTURE.md](./ARCHITECTURE.md)

---

## 一、调研结论

### 1.1 主流组件库分类对比

| 组件库 | 分类方式 |
|--------|---------|
| **Ant Design** | General / Layout / Navigation / Data Entry / Data Display / Feedback |
| **Arco Design** | General / Layout / Navigation / Data Entry / Data Display / Feedback / Other |
| **Element Plus** | Basic / Configuration / Form / Data / Navigation / Feedback / Others |
| **MUI** | Inputs / Data Display / Feedback / Surfaces / Navigation / Layout |
| **Flutter** | Basics / Layout / Input / Interaction / Scrolling / Text / Assets |
| **Chakra UI** | Layout / Forms / Data Display / Feedback / Overlay / Disclosure |

### 1.2 共性功能模式

所有成熟组件库都遵循以下设计模式：

1. **原子化设计**：基础组件（Text/Box/Button）+ 复合组件（Card/Table/Form）
2. **Variants 变体**：Button 有 Primary/Secondary/Ghost/Dashed/Link
3. **Size 尺寸**：Small / Default / Large 三级尺寸
4. **State 状态**：Normal / Hover / Active / Focus / Disabled / Loading
5. **受控与非受控**：表单组件支持 `value` + `onChange`（受控）或仅 `defaultValue`
6. **Design Token 驱动**：颜色、圆角、间距、字体大小由主题系统统一管控

---

## 二、本框架组件分类体系

基于 Yoga Flexbox 布局特性和 Ebiten 桌面渲染能力，采用 **7 大类 + 辅助函数**：

```
components/
├── 01-primitives/           # 基础原子组件（产生实际像素）
│   ├── Text, RichText
│   ├── Image, Icon
│   └── Box
├── 02-layout/               # 布局组件（组织空间）
│   ├── Row, Column
│   ├── Stack
│   ├── ScrollView
│   ├── Spacer
│   └── Divider
├── 03-form/                 # 表单/输入组件
│   ├── Button
│   ├── Input, TextArea
│   ├── Checkbox, Radio, Switch
│   ├── Select, Dropdown
│   └── Slider, Range
├── 04-data-display/         # 数据展示
│   ├── List, ListView
│   ├── Table, DataTable
│   ├── Card
│   ├── Badge, Tag
│   ├── Avatar
│   ├── Tooltip
│   └── Progress, ProgressBar
├── 05-feedback/             # 反馈组件
│   ├── Dialog, Modal
│   ├── Toast, Snackbar, Message
│   ├── Loading, Spinner, Skeleton
│   └── Drawer, Sidebar
├── 06-navigation/           # 导航组件
│   ├── Tabs, TabBar
│   ├── Navbar, AppBar
│   ├── Breadcrumb
│   └── Pagination
└── 07-helpers/              # 辅助组件
    ├── If, ForEach
    ├── GestureDetector
    └── ThemeProvider
```

---

## 三、通用设计规范

### 3.1 所有组件遵循的 Options 模式

```go
// 通用属性（所有组件都支持）
func Width(v float32) Option
func Height(v float32) Option
func MinWidth(v float32) Option
func MaxWidth(v float32) Option
func Padding(v float32) Option          // 统一内边距
func PaddingHV(h, v float32) Option     // 水平/垂直分别设置
func Margin(v float32) Option
func Background(c Color) Option
func Opacity(v float32) Option
func BorderRadius(v float32) Option
func Border(width float32, c Color) Option
func FlexGrow(v float32) Option
func FlexShrink(v float32) Option
func Key(k string) Option               // 稳定身份
```

### 3.2 状态驱动样式（参考 MUI `sx` + Chakra `_hover`）

```go
Button(
    Background(ColorBlue),
    Hover(Background(ColorDarkBlue)),
    Active(Transform(Scale(0.98))),
    Disabled(Opacity(0.5)),
)
```

### 3.3 Size 枚举（参考 Ant Design）

```go
type Size int
const (
    SizeSmall Size = iota
    SizeDefault
    SizeLarge
)

// 使用
Button(Size(SizeSmall), Text("Small"))
Button(Size(SizeLarge), Text("Large"))
```

### 3.4 Variant 枚举（参考 MUI + Chakra）

```go
type Variant int
const (
    VariantSolid Variant = iota      // 实心填充（Primary）
    VariantOutlined                   // 描边
    VariantGhost                      // 透明背景，hover 显色
    VariantLink                       // 纯文字链接
    VariantDashed                     // 虚线描边
)
```

---

## 四、逐类组件设计

### 4.1 基础原子组件（Primitives）

#### Text / RichText

参考：Flutter Text、MUI Typography、Chakra Text

```go
// 基础文本
Text("Hello", FontSize(16), Color(ColorGray))

// 富文本
RichText(
    Span("Hello ", FontSize(18)),
    Span("World", Bold, Color(ColorRed)),
)
```

Options：
- `FontSize(float32)`
- `FontWeight(FontWeight)` — Normal, Bold, Medium, Light
- `Color(Color)`
- `LineHeight(float32)`
- `MaxLines(int)`
- `TextAlign(TextAlign)` — Left, Center, Right
- `TextDecoration(TextDecoration)` — None, Underline, LineThrough

#### Image / Icon

参考：Flutter Image、MUI Icon、Ant Design Icon

```go
Image("avatar.png", Width(48), Height(48), BorderRadius(24))
Icon(Heart, Color(ColorRed), Size(24))
```

Options：
- `Width/Height`
- `ObjectFit(ObjectFit)` — Cover, Contain, Fill
- `BorderRadius`

#### Box

参考：MUI Box、Chakra Box、Flutter Container

```go
Box(
    Width(200), Height(100),
    Background(ColorWhite),
    BorderRadius(8),
    BoxShadow(Offset(0,4), Blur(8), ColorRGBA(0,0,0,0.1)),
    Text("Content"),
)
```

Box 是通用容器，支持所有视觉属性 + 作为布局容器（可嵌套 Row/Column 等）。

---

### 4.2 布局组件（Layout）

#### Row / Column

参考：Flutter Row/Column、CSS Flexbox、MUI Stack

```go
Row(
    Gap(16),
    AlignItems(AlignCenter),
    JustifyContent(JustifyBetween),
    Text("Left"),
    Text("Right"),
)

Column(
    Gap(8),
    Text("A"),
    Text("B"),
)
```

Options：
- `Gap(float32)` — 子元素间距（Yoga gap）
- `AlignItems(Align)` — FlexStart, Center, FlexEnd, Stretch
- `JustifyContent(Justify)` — FlexStart, Center, FlexEnd, SpaceBetween, SpaceAround, SpaceEvenly
- `Wrap(bool)` — 是否换行（FlexWrap）

#### Stack

参考：Flutter Stack、MUI Stack、CSS Position

```go
Stack(
    Width(300), Height(200),
    Image("bg.jpg", ObjectFitCover),
    Box(
        Position(Fill),
        Background(LinearGradient(...)),
    ),
    Text("Overlay", Position(BottomLeft)),
)
```

#### ScrollView

参考：Flutter SingleChildScrollView、CSS overflow

```go
ScrollView(
    Direction(Vertical),
    ShowScrollbar(true),
    Column(
        Gap(16),
        ForEach(items, renderItem)...,
    ),
)
```

#### Spacer / Divider

参考：Flutter Spacer/Divider、Ant Design Divider、MUI Divider

```go
Row(
    Text("Left"),
    Spacer(),           // 占据剩余空间（FlexGrow(1) 语法糖）
    Text("Right"),
)

Divider(Color(ColorGray), Thickness(1))
Divider(Vertical, Height(20), Thickness(1))  // 垂直分隔线
```

---

### 4.3 表单组件（Form）

#### Button

参考：MUI Button、Ant Design Button、Chakra Button、Flutter ElevatedButton

```go
Button(
    Variant(VariantSolid),
    Size(SizeDefault),
    ColorTheme(ColorBlue),          // 主题色
    OnTap(c.SetState(func() { ... })),
    Text("Submit"),
)

Button(
    Variant(VariantOutlined),
    Size(SizeSmall),
    Disabled(c.loading),
    Loading(c.loading),             // 自动显示 Spinner
    Icon(Send),
    Text("Send"),
)
```

Options：
- `Variant(Variant)`
- `Size(Size)`
- `ColorTheme(Color)` — 主色调
- `Disabled(bool)`
- `Loading(bool)` — 显示加载状态
- `OnTap(func())`
- `OnLongPress(func())`
- `Icon(IconType)` — 前缀图标
- `IconRight(IconType)` — 后缀图标
- `Hover/Active/Disabled(Option)` — 状态样式

#### Input / TextArea

参考：MUI TextField、Ant Design Input、Flutter TextField、Chakra Input

```go
Input(
    Placeholder("Enter your name"),
    Value(c.name),
    OnChange(c.SetState(func(v string) { c.name = v })),
    Prefix(Icon(User)),
    Suffix(Button(OnTap(c.clear), Icon(Close))),
    Focus(Border(2, ColorBlue)),
)

TextArea(
    Rows(4),
    MaxLength(200),
    ShowCount(true),
)
```

Options：
- `Value(string)` / `DefaultValue(string)`
- `Placeholder(string)`
- `Type(InputType)` — Text, Password, Email, Number, URL
- `Disabled(bool)`
- `ReadOnly(bool)`
- `MaxLength(int)`
- `Prefix(Widget)` / `Suffix(Widget)`
- `OnChange(func(string))`
- `OnSubmit(func())`
- `OnFocus/OnBlur(func())`
- `Focus(Option)` — 聚焦状态样式

#### Checkbox / Radio / Switch

参考：MUI Checkbox/Radio/Switch、Ant Design、Flutter Checkbox

```go
Checkbox(
    Checked(c.agreed),
    OnChange(c.SetState(func(v bool) { c.agreed = v })),
    Text("I agree to terms"),
)

Radio(
    Checked(c.selected == "a"),
    OnChange(c.SetState(func() { c.selected = "a" })),
    Text("Option A"),
)

Switch(
    Checked(c.enabled),
    OnChange(c.SetState(func(v bool) { c.enabled = v })),
)
```

#### Select / Dropdown

参考：MUI Select、Ant Design Select、Flutter DropdownButton

```go
Select(
    Value(c.city),
    Options(
        Option("bj", "Beijing"),
        Option("sh", "Shanghai"),
    ),
    OnChange(c.SetState(func(v string) { c.city = v })),
)
```

#### Slider / Range

参考：MUI Slider、Ant Design Slider、Flutter Slider

```go
Slider(
    Min(0), Max(100),
    Value(c.volume),
    OnChange(c.SetState(func(v float32) { c.volume = v })),
)
```

---

### 4.4 数据展示（Data Display）

#### List / ListView

参考：Flutter ListView、Ant Design List、MUI List

```go
ListView(
    Data(c.items),
    KeyBy(func(item Item) string { return item.ID }),
    Render(func(item Item) Widget {
        return ListItem(
            Title(item.Name),
            Description(item.Desc),
        )
    }),
)
```

#### Table / DataTable

参考：Ant Design Table、MUI DataGrid、Flutter DataTable

```go
DataTable(
    Columns(
        Column("Name", Width(200)),
        Column("Age", Width(100)),
        Column("Address"),
    ),
    Data(c.users),
)
```

#### Card

参考：MUI Card、Ant Design Card、Flutter Card

```go
Card(
    Width(300),
    BorderRadius(12),
    BoxShadow(...),
    Cover(Image("cover.jpg")),
    Header(Text("Title", FontSize(18), Bold)),
    Body(Text("Card content...")),
    Footer(Button(Text("Action"))),
)
```

#### Badge / Tag / Avatar

参考：Ant Design Badge/Tag/Avatar、MUI Badge/Chip/Avatar、Chakra Badge/Tag/Avatar

```go
Badge(
    Count(5),
    Color(ColorRed),
    Icon(Bell),
)

Tag("New", Color(ColorGreen))
Tag("Deprecated", Variant(VariantOutlined), Color(ColorRed))

Avatar("user.png", Size(40))
Avatar("Alice", FallbackText(true), Size(40))  // 无图片时显示首字母
```

#### Tooltip

参考：MUI Tooltip、Ant Design Tooltip、Flutter Tooltip

```go
Tooltip(
    Text("Hover me"),
    Content(Text("This is a tooltip")),
    Placement(Top),
)
```

#### Progress / ProgressBar

参考：Ant Design Progress、MUI LinearProgress、Flutter LinearProgressIndicator

```go
ProgressBar(
    Value(0.6),                     // 0~1
    Color(ColorBlue),
    ShowLabel(true),                // 显示 "60%"
)

ProgressCircle(
    Value(0.75),
    Size(80),
)
```

---

### 4.5 反馈组件（Feedback）

#### Dialog / Modal / AlertDialog

参考：MUI Dialog、Ant Design Modal、Flutter AlertDialog、Chakra Modal

```go
// 声明式弹窗（推荐）
If(c.showDialog,
    Dialog(
        Title(Text("Confirm")),
        Content(Text("Are you sure?")),
        Actions(
            Button(OnTap(c.hideDialog), Text("Cancel")),
            Button(OnTap(c.confirm), Variant(VariantSolid), Text("OK")),
        ),
    ),
)

// 命令式弹窗（全局 API）
Dialog.Confirm("Delete?", func(ok bool) {
    if ok { c.delete() }
})
```

#### Toast / Snackbar / Message

参考：Ant Design message、MUI Snackbar、Flutter SnackBar、Element Plus ElMessage

```go
// 命令式全局调用（最常用）
Toast.Success("Operation successful")
Toast.Error("Something went wrong")
Toast.Info("Please wait...", Duration(3000))

// 带 Action 的 Snackbar
Snackbar(
    Text("Item deleted"),
    Action(Button(Text("Undo"), OnTap(c.undo))),
)
```

#### Loading / Spinner / Skeleton

参考：Ant Design Spin/Skeleton、MUI CircularProgress/Skeleton、Chakra Skeleton

```go
Spinner(Size(40), Color(ColorBlue))     // 旋转指示器
Skeleton(Width(200), Height(20))         // 骨架屏占位
Skeleton.Circle(Size(48))                // 圆形骨架

// 包裹模式
Loading(c.loading,
    Spinner(),
    Content(),
)
```

#### Drawer / Sidebar

参考：Ant Design Drawer、MUI Drawer、Flutter Drawer

```go
If(c.showDrawer,
    Drawer(
        Placement(Right),
        Width(360),
        Title(Text("Settings")),
        Content(settingsForm),
        Mask(OnTap(c.closeDrawer)),
    ),
)
```

---

### 4.6 导航组件（Navigation）

#### Tabs / TabBar

参考：MUI Tabs、Ant Design Tabs、Flutter TabBar

```go
Tabs(
    ActiveKey(c.activeTab),
    OnChange(c.SetState(func(k string) { c.activeTab = k })),
    Tab("home", "Home", HomePage()),
    Tab("profile", "Profile", ProfilePage()),
    Tab("settings", "Settings", SettingsPage()),
)
```

#### Navbar / AppBar

参考：MUI AppBar、Ant Design PageHeader、Flutter AppBar

```go
AppBar(
    Height(56),
    Background(ColorWhite),
    Shadow(Offset(0,1), Blur(2)),
    Leading(Button(Icon(Menu))),
    Title(Text("My App")),
    Actions(
        Button(Icon(Search)),
        Button(Icon(Notifications)),
    ),
)
```

#### Breadcrumb

参考：Ant Design Breadcrumb、MUI Breadcrumbs

```go
Breadcrumb(
    Item("Home", OnTap(c.goHome)),
    Item("Products"),
    Item("Detail"),
)
```

#### Pagination

参考：Ant Design Pagination、MUI Pagination

```go
Pagination(
    Current(c.page),
    Total(c.total),
    PageSize(10),
    OnChange(c.SetState(func(p int) { c.page = p })),
)
```

---

### 4.7 辅助组件（Helpers）

#### If / ForEach

```go
If(c.loading, Spinner())

Column(
    ForEach(c.items, func(item Item, i int) Widget {
        return Row(Key(item.ID), Text(item.Name))
    })...,
)
```

#### GestureDetector

参考：Flutter GestureDetector

```go
GestureDetector(
    OnTap(func() { ... }),
    OnPan(func(g PanGesture) { ... }),
    OnLongPress(func() { ... }),
    Box(Text("Interact with me")),
)
```

#### ThemeProvider

参考：MUI ThemeProvider、Chakra ChakraProvider

```go
ThemeProvider(c.theme,
    App(...),
)

// 子组件消费
Box(
    Background(ThemeColor("primary")),
    Text("Themed", Color(ThemeColor("text.primary"))),
)
```

---

## 五、组件实现优先级

### P0（Phase 1 骨架必备）

| 组件 | 理由 |
|------|------|
| Text | 最基本的渲染原语 |
| Box | 通用容器，所有视觉属性的载体 |
| Row / Column | Yoga Flex 布局的核心 |
| Button | 验证交互事件通路 |
| Stack | 验证绝对定位 |
| If / ForEach | 声明式逻辑的基础 |

### P1（Phase 2 基础库）

| 组件 | 理由 |
|------|------|
| Image | 常见渲染需求 |
| Icon | 按钮、导航必备 |
| ScrollView | 内容溢出处理 |
| Input | 表单基础 |
| Checkbox / Switch | 布尔输入 |
| Divider / Spacer | 布局微调 |
| Card | 最常见的复合容器 |
| Badge / Avatar | 列表/导航常见 |
| ProgressBar | 反馈基础 |

### P2（Phase 3 交互完善）

| 组件 | 理由 |
|------|------|
| TextArea | 多行输入 |
| Radio / Select | 单选/下拉 |
| Slider | 范围输入 |
| ListView | 数据列表 |
| Tooltip | 信息提示 |
| Dialog | 模态弹窗 |
| Toast | 全局反馈 |
| Loading / Spinner | 加载状态 |
| Tabs | 内容切换 |
| Drawer | 侧栏导航 |

### P3（Phase 4 高级特性）

| 组件 | 理由 |
|------|------|
| Table / DataTable | 复杂数据展示 |
| Breadcrumb | 深层导航 |
| Pagination | 数据分页 |
| Navbar / AppBar | 页面框架 |
| RichText | 富文本排版 |
| GestureDetector | 复杂手势 |
| ThemeProvider | 主题切换 |
| Skeleton | 骨架屏 |
| Snackbar | 底部提示 |
| ProgressCircle | 环形进度 |
