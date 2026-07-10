# 文字排版与输入专项设计

> 父文档：[COMPONENT_LIBRARY_DESIGN.md](./COMPONENT_LIBRARY_DESIGN.md)
>
> 技术栈：Ebiten text/v2 + Yoga MeasureFunc + Ebiten inpututil

---

## 一、文字排版调研

### 1.1 Ebiten text/v2 核心 API

```go
package textv2 "github.com/hajimehoshi/ebiten/v2/text/v2"

// 字体源（可共享，线程安全缓存）
type GoTextFaceSource struct{}
func NewGoTextFaceSource(source io.Reader) (*GoTextFaceSource, error)

// 字体面（包含字号、粗细、颜色等选项）
type GoTextFace struct {
    Source *GoTextFaceSource
    Size   float64
    // ... Weight, Style, Direction 等
}

// 测量文本尺寸（关键：用于 Yoga MeasureFunc）
func Measure(text string, face Face, lineSpacingInPixels float64) (width, height float64)

// 单行进阶距离
func Advance(text string, face Face) float64

// 绘制文本
func Draw(dst *ebiten.Image, text string, face Face, options *DrawOptions)

// 布局选项

type LayoutOptions struct {
    LineSpacing    float64  // 行间距（像素）
    PrimaryAlign   Align    // 水平对齐：AlignStart/AlignCenter/AlignEnd
    SecondaryAlign Align    // 垂直对齐
}

type DrawOptions struct {
    LayoutOptions LayoutOptions
    GeoM          ebiten.GeoM      // 变换
    ColorScale    ebiten.ColorScale // 颜色缩放
}
```

### 1.2 text/v2 的关键特性

| 特性 | 说明 |
|------|------|
| **HarfBuzz shaping** | `GoTextFace` 内部使用 HarfBuzz，支持复杂文字（阿拉伯语、印地语、中文等） |
| **自动换行** | `Measure()` 在水平方向会根据可用宽度自动换行 |
| **多行** | 支持 `\n` 手动换行，`LineSpacing` 控制行距 |
| **对齐** | `PrimaryAlign` 控制水平对齐，`SecondaryAlign` 控制垂直对齐 |
| **缓存** | 字形自动 LRU 缓存，每帧重复绘制性能良好 |
| **颜色** | 通过 `DrawOptions.ColorScale` 控制，默认白色 |

### 1.3 与 Yoga MeasureFunc 的配合

Yoga 在布局时会调用 `MeasureFunc` 询问节点的自然尺寸：

```go
type MeasureFunc func(
    node *yoga.Node,
    width float32, widthMode yoga.MeasureMode,
    height float32, heightMode yoga.MeasureMode,
) yoga.Size
```

| MeasureMode | 含义 | Text 处理 |
|-------------|------|-----------|
| `MeasureModeUndefined` | 无约束 | 返回文本自然尺寸 |
| `MeasureModeExactly` | 固定尺寸 | 不需要测量，直接返回给定值 |
| `MeasureModeAtMost` | 最大约束 | 在 maxWidth 内换行，返回实际尺寸 |

```go
func textMeasureFunc(text string, face textv2.Face, lineSpacing float64) yoga.MeasureFunc {
    return func(node *yoga.Node, width float32, widthMode yoga.MeasureMode,
        height float32, heightMode yoga.MeasureMode) yoga.Size {
        
        var measureWidth float64
        if widthMode == yoga.MeasureModeAtMost {
            measureWidth = float64(width)
        } else if widthMode == yoga.MeasureModeExactly {
            return yoga.Size{Width: width, Height: height}
        }
        // Undefined: measureWidth = 0，Measure 会返回单行最大宽度
        
        w, h := textv2.Measure(text, face, lineSpacing)
        if widthMode == yoga.MeasureModeAtMost && w > measureWidth {
            // 重新测量，限制宽度（自动换行）
            // text/v2 的 Measure 目前不支持直接传入 maxWidth
            // 需要手动实现自动换行...（见下方方案）
        }
        
        return yoga.Size{Width: float32(w), Height: float32(h)}
    }
}
```

> **关键问题**：`text/v2.Measure()` 没有直接的 `maxWidth` 参数，需要我们自己实现自动换行逻辑。

---

## 二、文字排版实现方案

### 2.1 自动换行（Word Wrap）

text/v2 目前没有内建的 `maxWidth` 换行，需要手动实现：

```go
// wrapText 将文本按 maxWidth 自动换行
func wrapText(text string, face textv2.Face, maxWidth float64) string {
    if maxWidth <= 0 {
        return text // 无约束，不处理
    }
    
    var result strings.Builder
    lines := strings.Split(text, "\n")
    
    for li, line := range lines {
        if li > 0 {
            result.WriteByte('\n')
        }
        
        runes := []rune(line)
        if len(runes) == 0 {
            continue
        }
        
        // 尝试整行测量
        w, _ := textv2.Measure(line, face, 0)
        if w <= maxWidth {
            result.WriteString(line)
            continue
        }
        
        // 需要换行：逐字累积，超过宽度时换行
        start := 0
        for start < len(runes) {
            end := start + 1
            // 二分查找当前行能放下的最大字符数
            lo, hi := start+1, len(runes)+1
            for lo < hi {
                mid := (lo + hi) / 2
                substr := string(runes[start:mid])
                w, _ := textv2.Measure(substr, face, 0)
                if w <= maxWidth {
                    lo = mid + 1
                    end = mid
                } else {
                    hi = mid
                }
            }
            
            if start > 0 {
                result.WriteByte('\n')
            }
            result.WriteString(string(runes[start:end]))
            start = end
        }
    }
    
    return result.String()
}
```

> 优化：实际实现中可使用更高效的断行策略（如按单词边界 CJK 字符边界断行）。

### 2.2 Text 组件设计

```go
// Text 组件选项
type TextStyle struct {
    Content      string
    FontSize     float64
    FontWeight   FontWeight      // Normal, Bold, Medium
    Color        Color
    LineHeight   float64         // 倍数，如 1.5
    MaxLines     int             // 0 = 无限制
    TextAlign    TextAlign       // Left, Center, Right
    Overflow     TextOverflow    // Clip, Ellipsis
}
```

**Text RenderNode 实现**：

```go
type TextRenderNode struct {
    yogaNode    *yoga.Node
    bounds      Rect
    
    // 文本属性
    text        string
    face        textv2.Face
    lineSpacing float64
    color       Color
    align       textv2.Align
    maxWidth    float64        // 用于换行
    wrappedText string         // 换行后的文本
}

func (n *TextRenderNode) SyncYogaProps(w RenderWidget) {
    tw := w.(*TextWidget)
    
    // 创建/更新字体面
    n.face = GetFontFace(FontOptions{
        Size:   tw.FontSize,
        Weight: tw.FontWeight,
    })
    n.text = tw.Content
    n.color = tw.Color
    n.lineSpacing = tw.LineHeight * tw.FontSize
    
    // 对齐转换
    switch tw.TextAlign {
    case TextAlignCenter:
        n.align = textv2.AlignCenter
    case TextAlignRight:
        n.align = textv2.AlignEnd
    default:
        n.align = textv2.AlignStart
    }
    
    // 设置 MeasureFunc
    n.yogaNode.SetMeasureFunc(n.measure)
}

func (n *TextRenderNode) measure(node *yoga.Node,
    width float32, widthMode yoga.MeasureMode,
    height float32, heightMode yoga.MeasureMode) yoga.Size {
    
    // 根据 Yoga 约束决定是否换行
    maxW := float64(0)
    if widthMode == yoga.MeasureModeAtMost || widthMode == yoga.MeasureModeExactly {
        maxW = float64(width)
    }
    
    if maxW > 0 && maxW != n.maxWidth {
        n.maxWidth = maxW
        n.wrappedText = wrapText(n.text, n.face, maxW)
    } else if maxW <= 0 {
        n.wrappedText = n.text
    }
    
    w, h := textv2.Measure(n.wrappedText, n.face, n.lineSpacing)
    
    return yoga.Size{
        Width:  float32(w),
        Height: float32(h),
    }
}

func (n *TextRenderNode) Paint(screen *ebiten.Image, bounds Rect) {
    op := &textv2.DrawOptions{}
    op.GeoM.Translate(float64(bounds.X), float64(bounds.Y))
    op.ColorScale.Scale(
        float32(n.color.R)/255,
        float32(n.color.G)/255,
        float32(n.color.B)/255,
        float32(n.color.A)/255,
    )
    op.LayoutOptions.LineSpacing = n.lineSpacing
    op.LayoutOptions.PrimaryAlign = n.align
    
    textv2.Draw(screen, n.wrappedText, n.face, op)
}
```

### 2.3 字体管理

```go
package font

import (
    textv2 "github.com/hajimehoshi/ebiten/v2/text/v2"
)

var defaultSource *textv2.GoTextFaceSource

func Init(fontPath string) error {
    f, err := os.Open(fontPath)
    if err != nil { return err }
    defer f.Close()
    
    defaultSource, err = textv2.NewGoTextFaceSource(f)
    return err
}

type FontOptions struct {
    Size   float64
    Weight FontWeight
}

func GetFace(opts FontOptions) textv2.Face {
    if defaultSource == nil {
        panic("font not initialized")
    }
    return &textv2.GoTextFace{
        Source: defaultSource,
        Size:   opts.Size,
    }
}
```

> **注意**：需要选择包含中文的字体（如 Noto Sans CJK、Source Han Sans、或系统默认字体）。

### 2.4 富文本（RichText）

text/v2 不直接支持富文本（同一段文字不同样式），需要手动分段绘制：

```go
type Span struct {
    Text  string
    Face  textv2.Face
    Color Color
}

// 手动计算每段位置并绘制
func (n *RichTextRenderNode) Paint(screen *ebiten.Image, bounds Rect) {
    x, y := float64(bounds.X), float64(bounds.Y)
    lineHeight := n.calculateLineHeight()
    
    for _, span := range n.spans {
        // 分段绘制，手动管理 x/y 偏移
        op := &textv2.DrawOptions{}
        op.GeoM.Translate(x, y)
        op.ColorScale.Scale(...)
        textv2.Draw(screen, span.Text, span.Face, op)
        
        advance, _ := textv2.Measure(span.Text, span.Face, 0)
        x += advance
    }
}
```

> 简化方案：Phase 1 先不支持 RichText，Phase 2+ 再实现。

---

## 三、键盘输入调研（inpututil）

### 3.1 Ebiten 输入 API 总览

| API | 用途 | 适用场景 |
|-----|------|----------|
| `ebiten.AppendInputChars(runes)` | 获取当前帧输入的可打印字符 | **文本输入（核心）** |
| `inpututil.IsKeyJustPressed(key)` | 按键是否刚按下 | 特殊键（Backspace, Enter, Arrow） |
| `inpututil.IsKeyJustReleased(key)` | 按键是否刚释放 | 较少使用 |
| `ebiten.IsKeyPressed(key)` | 按键是否持续按住 | 方向键长按移动光标 |
| `inpututil.KeyPressDuration(key)` | 按键持续帧数 | 长按退格连续删除 |

### 3.2 Input 组件输入处理方案

```go
func (input *InputComponent) handleKeyboard() {
    // 1. 处理字符输入（支持中文输入法）
    var runes []rune
    runes = ebiten.AppendInputChars(runes)
    for _, r := range runes {
        input.insertRune(r)
    }
    
    // 2. 处理退格（Backspace）
    if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
        input.deleteBackward()
    }
    // 长按退格连续删除
    if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
        duration := inpututil.KeyPressDuration(ebiten.KeyBackspace)
        if duration > 30 && duration%5 == 0 { // 按住 0.5s 后开始连续删除
            input.deleteBackward()
        }
    }
    
    // 3. 处理删除键（Delete）
    if inpututil.IsKeyJustPressed(ebiten.KeyDelete) {
        input.deleteForward()
    }
    
    // 4. 处理方向键（光标移动）
    if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
        input.moveCursor(-1)
    }
    if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
        input.moveCursor(1)
    }
    // 长按方向键
    if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
        if d := inpututil.KeyPressDuration(ebiten.KeyArrowLeft); d > 30 && d%3 == 0 {
            input.moveCursor(-1)
        }
    }
    
    // 5. 处理 Home / End
    if inpututil.IsKeyJustPressed(ebiten.KeyHome) {
        input.cursor = 0
    }
    if inpututil.IsKeyJustPressed(ebiten.KeyEnd) {
        input.cursor = len([]rune(input.value))
    }
    
    // 6. 处理回车（Enter）
    if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
        if input.onSubmit != nil {
            input.onSubmit()
        }
    }
    
    // 7. 处理全选（Ctrl+A）
    if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyA) {
        input.selectAll()
    }
    
    // 8. 处理复制粘贴（Ctrl+C/V/X）
    if ebiten.IsKeyPressed(ebiten.KeyControl) {
        if inpututil.IsKeyJustPressed(ebiten.KeyC) {
            input.copy()
        }
        if inpututil.IsKeyJustPressed(ebiten.KeyV) {
            input.paste()
        }
        if inpututil.IsKeyJustPressed(ebiten.KeyX) {
            input.cut()
        }
    }
}
```

### 3.3 光标与选择

```go
type InputState struct {
    value       string        // 当前值
    cursor      int           // 光标位置（rune 索引）
    selectionStart int        // 选择起始
    selectionEnd   int        // 选择结束（-1 表示无选择）
    focused     bool
}

func (s *InputState) insertRune(r rune) {
    runes := []rune(s.value)
    if s.hasSelection() {
        s.deleteSelection()
    }
    // 在光标位置插入
    newRunes := append(runes[:s.cursor], append([]rune{r}, runes[s.cursor:]...)...)
    s.value = string(newRunes)
    s.cursor++
    s.notifyChange()
}

func (s *InputState) deleteBackward() {
    if s.hasSelection() {
        s.deleteSelection()
        return
    }
    if s.cursor <= 0 { return }
    runes := []rune(s.value)
    s.value = string(append(runes[:s.cursor-1], runes[s.cursor:]...))
    s.cursor--
    s.notifyChange()
}

func (s *InputState) moveCursor(delta int) {
    runes := []rune(s.value)
    s.cursor += delta
    if s.cursor < 0 { s.cursor = 0 }
    if s.cursor > len(runes) { s.cursor = len(runes) }
    s.clearSelection()
}
```

### 3.4 光标绘制

```go
func (n *InputRenderNode) Paint(screen *ebiten.Image, bounds Rect) {
    // 1. 绘制背景、边框...
    
    // 2. 绘制文本（带 clip，滚动显示）
    clipBounds := bounds.Inset(padding)
    // ...
    
    // 3. 绘制光标（如果 focused）
    if n.focused && n.cursorVisible {
        cursorX := n.measureCursorX()
        vector.DrawFilledRect(screen,
            float32(bounds.X+cursorX), float32(bounds.Y+paddingTop),
            1, float32(bounds.H-paddingTop-paddingBottom),
            color.Black, true)
    }
    
    // 4. 绘制选区高亮
    if n.hasSelection() {
        selX1 := n.measureCursorXAt(n.selectionStart)
        selX2 := n.measureCursorXAt(n.selectionEnd)
        vector.DrawFilledRect(screen,
            float32(bounds.X+selX1), float32(bounds.Y+paddingTop),
            float32(selX2-selX1), float32(bounds.H-paddingTop-paddingBottom),
            Color{R: 0, G: 100, B: 255, A: 100}, true)
    }
}
```

### 3.5 光标闪烁

```go
func (n *InputRenderNode) Update() {
    if n.focused {
        n.cursorTick++
        if n.cursorTick%30 == 0 { // 约 0.5s 闪烁一次（60fps）
            n.cursorVisible = !n.cursorVisible
        }
    } else {
        n.cursorVisible = false
    }
}
```

### 3.6 剪贴板支持

Ebiten 2.6+ 提供剪贴板 API：

```go
import "github.com/hajimehoshi/ebiten/v2"

func (input *InputComponent) copy() {
    if input.hasSelection() {
        ebiten.SetClipboardText(input.selectedText())
    }
}

func (input *InputComponent) paste() {
    text := ebiten.ClipboardText()
    for _, r := range text {
        input.insertRune(r)
    }
}
```

---

## 四、Text + Input API 设计（用户可见层）

### 4.1 Text 组件

```go
Text("Hello World",
    FontSize(16),           // float64
    FontWeightBold,         // Normal | Bold | Medium | Light
    Color(ColorHex("#333")),
    LineHeight(1.5),        // 倍数
    MaxLines(2),            // 超过截断
    TextAlignCenter,        // Left | Center | Right
    OverflowEllipsis,       // Clip | Ellipsis
)
```

### 4.2 Input 组件

```go
Input(
    Placeholder("Enter text..."),
    Value(c.text),
    OnChange(c.SetState(func(v string) { c.text = v })),
    
    // 样式
    FontSize(14),
    Color(ColorBlack),
    Background(ColorWhite),
    BorderRadius(4),
    Border(1, ColorGray),
    Padding(8, 12),
    
    // Focus 状态
    Focus(Border(1, ColorBlue)),
    
    // 限制
    MaxLength(100),
    Type(InputTypeText),     // Text | Password | Email | Number | URL
    
    // 装饰
    Prefix(Icon(Search)),
    Suffix(If(c.text != "", Button(OnTap(c.clear), Icon(Close)))),
)
```

---

## 五、实现优先级

| 阶段 | 内容 |
|------|------|
| **Phase 1** | Text 组件（单行 + 自动换行 + 对齐）、字体初始化 |
| **Phase 2** | Input 组件（字符输入 + 退格 + 光标 + Focus） |
| **Phase 3** | Input 高级（选择、复制粘贴、方向键、Home/End、长按连续输入） |
| **Phase 4** | RichText（多 Span、不同字号/颜色混排） |
| **Phase 5** | 输入法优化（中文候选词、Composition 事件） |
