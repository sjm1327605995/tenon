// Package declarative 提供 React/SwiftUI 风格的声明式 UI API。
//
//	app := declarative.VStack(
//	    declarative.Text("Hello").FontSize(24),
//	    declarative.Button("Click").Style(ButtonPrimary),
//	).Padding(16)
package declarative

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/render"
	"github.com/sjm1327605995/tenon/pkg/engine"
	"github.com/sjm1327605995/tenon/pkg/widgets"
	"github.com/sjm1327605995/tenon/yoga"
)

// ==================== 类型别名 ====================

type Widget = engine.Widget
type BuildContext = engine.BuildContext
type State = engine.State
type RouteParams = engine.RouteParams

// ==================== 颜色常量 ====================

var (
	Black       = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	White       = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	Red         = color.RGBA{R: 239, G: 68, B: 68, A: 255}
	Green       = color.RGBA{R: 34, G: 197, B: 94, A: 255}
	Blue        = color.RGBA{R: 59, G: 130, B: 246, A: 255}
	Gray        = color.RGBA{R: 115, G: 115, B: 115, A: 255}
	Transparent = color.RGBA{}
)

// ==================== builder: 统一的声明式 Element ====================

// builder 是所有声明式包装类型的统一 Element。
// 它持有底层 Element，Update 时解包声明式类型再 diff。
type builder struct {
	engine.BaseElement
	inner engine.Element // 底层真实的 Element
}

func newBuilder(w engine.Widget) *builder {
	b := &builder{}
	b.BaseElement.Init(b, w)
	return b
}

func (b *builder) Mount(parent engine.Element, slot int) {
	b.BaseElement.Mount(parent, slot)
	real := b.unwrap(b.GetWidget())
	b.inner = real.CreateElement()
	b.inner.Mount(b, 0)
}

func (b *builder) Update(newWidget engine.Widget) {
	b.BaseElement.Update(newWidget)
	real := b.unwrap(newWidget)
	if b.inner != nil && engine.CanUpdate(b.inner.GetWidget(), real) {
		b.inner.Update(real)
	} else {
		if b.inner != nil {
			b.inner.Unmount()
		}
		b.inner = real.CreateElement()
		b.inner.Mount(b, 0)
	}
}

func (b *builder) Unmount() {
	if b.inner != nil {
		b.inner.Unmount()
		b.inner = nil
	}
	b.BaseElement.Unmount()
}

func (b *builder) GetChildren() []engine.Element {
	if b.inner == nil {
		return nil
	}
	return []engine.Element{b.inner}
}

func (b *builder) FindRenderObject() render.RenderObject {
	if b.inner != nil {
		return b.inner.FindRenderObject()
	}
	return nil
}

// unwrap 将声明式包装类型解包为底层 Widget。
type unwrapper interface {
	unwrap() engine.Widget
}

func (b *builder) unwrap(w engine.Widget) engine.Widget {
	if u, ok := w.(unwrapper); ok {
		return u.unwrap()
	}
	return w
}

// ==================== VStack / HStack ====================

func VStack(children ...Widget) FlexBox {
	return FlexBox{children: children, direction: yoga.FlexDirectionColumn}
}

func HStack(children ...Widget) FlexBox {
	return FlexBox{children: children, direction: yoga.FlexDirectionRow}
}

func Spacer() FlexBox {
	return FlexBox{flexGrow: 1}
}

type FlexBox struct {
	children   []Widget
	direction  yoga.FlexDirection
	gap        float32
	padding    engine.EdgeInsets
	justify    yoga.Justify
	alignItems yoga.Align
	flexGrow   float32
}

func (f FlexBox) Gap(v float32) FlexBox           { f.gap = v; return f }
func (f FlexBox) Gapf(v float32) FlexBox          { return f.Gap(v) }
func (f FlexBox) Padding(v float32) FlexBox       { f.padding = engine.EdgeInsetsAll(v); return f }
func (f FlexBox) Paddingf(insets engine.EdgeInsets) FlexBox { f.padding = insets; return f }
func (f FlexBox) Justify(v yoga.Justify) FlexBox  { f.justify = v; return f }
func (f FlexBox) JustifyContent(v yoga.Justify) FlexBox { f.justify = v; return f }
func (f FlexBox) Align(v yoga.Align) FlexBox      { f.alignItems = v; return f }
func (f FlexBox) AlignItems(v yoga.Align) FlexBox { f.alignItems = v; return f }
func (f FlexBox) Grow(v float32) FlexBox          { f.flexGrow = v; return f }
func (f FlexBox) CreateElement() engine.Element       { return newBuilder(f) }
func (f FlexBox) GetKey() engine.Key                  { return engine.NilKey{} }
func (f FlexBox) unwrap() engine.Widget {
	if f.direction == yoga.FlexDirectionRow {
		r := widgets.Row(f.children...)
		if f.gap > 0 {
			r = r.Gapf(f.gap)
		}
		r = r.Paddingf(f.padding)
		if f.justify != 0 {
			r = r.JustifyContent(f.justify)
		}
		if f.alignItems != 0 {
			r = r.AlignItems(f.alignItems)
		}
		return r
	}
	c := widgets.Column(f.children...)
	if f.gap > 0 {
		c = c.Gapf(f.gap)
	}
	c = c.Paddingf(f.padding)
	if f.justify != 0 {
		c = c.JustifyContent(f.justify)
	}
	if f.alignItems != 0 {
		c = c.AlignItems(f.alignItems)
	}
	return c
}

// ==================== Text ====================

func Text(content string) TextWidget {
	return TextWidget{content: content}
}

type TextWidget struct {
	content  string
	fontSize float32
	c        color.Color
	maxLines int
}

func (t TextWidget) FontSize(v float32) TextWidget { t.fontSize = v; return t }
func (t TextWidget) Color(c color.Color) TextWidget { t.c = c; return t }
func (t TextWidget) MaxLines(n int) TextWidget      { t.maxLines = n; return t }
func (t TextWidget) CreateElement() engine.Element       { return newBuilder(t) }
func (t TextWidget) GetKey() engine.Key                  { return engine.NilKey{} }
func (t TextWidget) unwrap() engine.Widget {
	w := widgets.Text(t.content)
	if t.fontSize > 0 {
		w = w.FontSize(t.fontSize)
	}
	if t.c != nil {
		w = w.Color(t.c)
	}
	if t.maxLines > 0 {
		w = w.MaxLines(t.maxLines)
	}
	return w
}

// ==================== Button ====================

type ButtonStyle = widgets.ButtonVariant

const (
	ButtonPrimary     = widgets.ButtonDefault
	ButtonSecondary   = widgets.ButtonSecondary
	ButtonOutline     = widgets.ButtonOutline
	ButtonGhost       = widgets.ButtonGhost
	ButtonDestructive = widgets.ButtonDestructive
	ButtonLink        = widgets.ButtonLink
)

func Button(label string) ButtonWidget {
	return ButtonWidget{label: label}
}

type ButtonWidget struct {
	label    string
	style    ButtonStyle
	onClick  func()
	disabled bool
	loading  bool
	height   float32
}

func (b ButtonWidget) Style(v ButtonStyle) ButtonWidget { b.style = v; return b }
func (b ButtonWidget) Variantf(v ButtonStyle) ButtonWidget { return b.Style(v) }
func (b ButtonWidget) OnClick(fn func()) ButtonWidget   { b.onClick = fn; return b }
func (b ButtonWidget) OnTap(fn func()) ButtonWidget     { return b.OnClick(fn) }
func (b ButtonWidget) Disabled(v bool) ButtonWidget     { b.disabled = v; return b }
func (b ButtonWidget) SetDisabled(v bool) ButtonWidget  { return b.Disabled(v) }
func (b ButtonWidget) Loading(v bool) ButtonWidget      { b.loading = v; return b }
func (b ButtonWidget) SetLoading(v bool) ButtonWidget   { return b.Loading(v) }
func (b ButtonWidget) H(v float32) ButtonWidget         { b.height = v; return b }
func (b ButtonWidget) CreateElement() engine.Element        { return newBuilder(b) }
func (b ButtonWidget) GetKey() engine.Key                   { return engine.NilKey{} }
func (b ButtonWidget) unwrap() engine.Widget {
	w := widgets.Button(b.label).
		Variantf(b.style).
		OnTap(b.onClick).
		SetDisabled(b.disabled).
		SetLoading(b.loading)
	if b.height > 0 {
		w = w.H(b.height)
	}
	return w
}

// ==================== Input ====================

func Input(placeholder string) InputWidget {
	return InputWidget{placeholder: placeholder}
}

type InputWidget struct {
	placeholder string
	value       string
	onChange    func(string)
	onSubmit    func(string)
	multiline   bool
	width       float32
	height      float32
	fontSize    float32
	bg          *render.Color
	borderColor *render.Color
	borderW     float32
	radius      float32
	padding     engine.EdgeInsets
}

func (i InputWidget) Value(v string) InputWidget           { i.value = v; return i }
func (i InputWidget) OnChange(fn func(string)) InputWidget { i.onChange = fn; return i }
func (i InputWidget) OnSubmit(fn func(string)) InputWidget { i.onSubmit = fn; return i }
func (i InputWidget) Multiline() InputWidget                { i.multiline = true; return i }
func (i InputWidget) Width(v float32) InputWidget           { i.width = v; return i }
func (i InputWidget) Height(v float32) InputWidget          { i.height = v; return i }
func (i InputWidget) FontSize(v float32) InputWidget        { i.fontSize = v; return i }
func (i InputWidget) Background(c color.Color) InputWidget  { i.bg = render.NewColorFrom(c); return i }
func (i InputWidget) Border(c color.Color, w float32) InputWidget { i.borderColor = render.NewColorFrom(c); i.borderW = w; return i }
func (i InputWidget) CornerRadius(v float32) InputWidget    { i.radius = v; return i }
func (i InputWidget) Padding(v float32) InputWidget         { i.padding = engine.EdgeInsetsAll(v); return i }
func (i InputWidget) CreateElement() engine.Element             { return newBuilder(i) }
func (i InputWidget) GetKey() engine.Key                        { return engine.NilKey{} }
func (i InputWidget) unwrap() engine.Widget {
	if i.multiline {
		w := widgets.Textarea(i.value).Placeholder(i.placeholder)
		if i.width > 0 {
			w = w.W(i.width)
		}
		if i.height > 0 {
			w = w.H(i.height)
		}
		if i.fontSize > 0 {
			w = w.FontSize(i.fontSize)
		}
		if i.bg != nil {
			w = w.Background(*i.bg)
		}
		if i.borderColor != nil && i.borderW > 0 {
			w = w.Border(*i.borderColor, i.borderW)
		}
		if i.radius > 0 {
			w = w.Radius(i.radius)
		}
		return w
	}
	w := widgets.TextField(i.value).Placeholder(i.placeholder)
	if i.width > 0 {
		w = w.W(i.width)
	}
	if i.height > 0 {
		w = w.H(i.height)
	}
	if i.fontSize > 0 {
		w = w.FontSize(i.fontSize)
	}
	if i.bg != nil {
		w = w.Bg(*i.bg)
	}
	if i.borderColor != nil && i.borderW > 0 {
		w = w.Border(*i.borderColor, i.borderW)
	}
	if i.radius > 0 {
		w = w.Radius(i.radius)
	}
	if i.padding != (engine.EdgeInsets{}) {
		w = w.Pad(i.padding)
	}
	if i.onChange != nil {
		w = w.OnChange(i.onChange)
	}
	if i.onSubmit != nil {
		w = w.OnSubmit(i.onSubmit)
	}
	return w
}

// ==================== Container ====================

func Container(child Widget) ContainerBox {
	return ContainerBox{child: child}
}

type ContainerBox struct {
	child         Widget
	bg            *render.Color
	bgImage       *ebiten.Image
	bgSlice       render.BorderSlice
	radius        float32
	border        *render.Color
	borderW       float32
	padding       engine.EdgeInsets
	margin        engine.EdgeInsets
	width         float32
	height        float32
	onClick       func()
	shadowColor   *render.Color
	shadowBlur    float32
	shadowOffsetX float32
	shadowOffsetY float32
}

func (c ContainerBox) Background(cl color.Color) ContainerBox      { c.bg = render.NewColorFrom(cl); return c }
func (c ContainerBox) BackgroundImage(img *ebiten.Image, slice render.BorderSlice) ContainerBox { c.bgImage = img; c.bgSlice = slice; return c }
func (c ContainerBox) CornerRadius(v float32) ContainerBox         { c.radius = v; return c }
func (c ContainerBox) Border(cl color.Color, w float32) ContainerBox { c.border = render.NewColorFrom(cl); c.borderW = w; return c }
func (c ContainerBox) Padding(v float32) ContainerBox              { c.padding = engine.EdgeInsetsAll(v); return c }
func (c ContainerBox) Pad(insets engine.EdgeInsets) ContainerBox       { c.padding = insets; return c }
func (c ContainerBox) Margin(v float32) ContainerBox               { c.margin = engine.EdgeInsetsAll(v); return c }
func (c ContainerBox) Marginf(insets engine.EdgeInsets) ContainerBox   { c.margin = insets; return c }
func (c ContainerBox) Width(v float32) ContainerBox                { c.width = v; return c }
func (c ContainerBox) W(v float32) ContainerBox                    { return c.Width(v) }
func (c ContainerBox) Height(v float32) ContainerBox               { c.height = v; return c }
func (c ContainerBox) H(v float32) ContainerBox                    { return c.Height(v) }
func (c ContainerBox) OnClick(fn func()) ContainerBox              { c.onClick = fn; return c }
func (c ContainerBox) OnTap(fn func()) ContainerBox                { return c.OnClick(fn) }
func (c ContainerBox) Shadow(cl color.Color, blur, offsetX, offsetY float32) ContainerBox {
	c.shadowColor = render.NewColorFrom(cl)
	c.shadowBlur = blur
	c.shadowOffsetX = offsetX
	c.shadowOffsetY = offsetY
	return c
}
func (c ContainerBox) CreateElement() engine.Element                   { return newBuilder(c) }
func (c ContainerBox) GetKey() engine.Key                              { return engine.NilKey{} }
func (c ContainerBox) unwrap() engine.Widget {
	w := widgets.Container(c.child)
	if c.bg != nil {
		w = w.Background(*c.bg)
	}
	if c.radius > 0 {
		w = w.Radius(c.radius)
	}
	if c.border != nil && c.borderW > 0 {
		w = w.Border(*c.border, c.borderW)
	}
	if c.bgImage != nil {
		w = w.BackgroundImage(c.bgImage, c.bgSlice)
	}
	w = w.Padding(c.padding).Margin(c.margin)
	if c.width > 0 {
		w = w.W(c.width)
	}
	if c.height > 0 {
		w = w.H(c.height)
	}
	if c.onClick != nil {
		w = w.OnTap(c.onClick)
	}
	if c.shadowColor != nil {
		w = w.Shadow(*c.shadowColor, c.shadowBlur, c.shadowOffsetX, c.shadowOffsetY)
	}
	return w
}

// ==================== Card ====================

func Card(child Widget) ContainerBox {
	return Container(child).Background(White).CornerRadius(8).Padding(16)
}

// ==================== Animated ====================

func Animated(child Widget) AnimatedBox {
	return AnimatedBox{child: child, duration: 300 * time.Millisecond}
}

type AnimatedBox struct {
	child    Widget
	duration time.Duration
	curve    engine.Curve
}

func (a AnimatedBox) Duration(d time.Duration) AnimatedBox { a.duration = d; return a }
func (a AnimatedBox) EaseInOut() AnimatedBox               { a.curve = engine.EaseInOutCurve{}; return a }
func (a AnimatedBox) CreateElement() engine.Element             { return newBuilder(a) }
func (a AnimatedBox) GetKey() engine.Key                        { return engine.NilKey{} }
func (a AnimatedBox) unwrap() engine.Widget {
	ac := widgets.NewAnimatedContainer().WithChild(a.child).WithDuration(a.duration)
	if a.curve != nil {
		ac = ac.WithCurve(a.curve)
	}
	return ac
}

// ==================== Navigator ====================

func Navigator(routes map[string]engine.RouteBuilder, initial string) NavigatorBox {
	return NavigatorBox{routes: routes, initial: initial}
}

type NavigatorBox struct {
	routes     map[string]engine.RouteBuilder
	initial    string
	transition engine.PageTransition
}

func (n NavigatorBox) Transition(t engine.PageTransition) NavigatorBox { n.transition = t; return n }
func (n NavigatorBox) CreateElement() engine.Element                   { return newBuilder(n) }
func (n NavigatorBox) GetKey() engine.Key                              { return engine.NilKey{} }
func (n NavigatorBox) unwrap() engine.Widget {
	w := engine.Navigator(n.routes, n.initial)
	if n.transition != engine.TransitionNone {
		w = w.WithTransition(n.transition)
	}
	return w
}

// ==================== 国际化 ====================

func Localization(locale string, translations map[string]map[string]string, child Widget) Widget {
	return engine.NewLocalization(locale, translations, child)
}

// ==================== 便捷函数 ====================

func L(ctx BuildContext, key string) string          { return engine.L(ctx, key) }
func GetNavigator(ctx BuildContext) engine.NavigatorState { return engine.GetNavigator(ctx) }
func Push(ctx BuildContext, name string, params ...RouteParams) { engine.NavPush(ctx, name, params...) }
func Pop(ctx BuildContext)                           { engine.NavPop(ctx) }
func SetTheme(t *engine.Theme)                           { engine.SetTheme(t) }
func Run(buildFunc engine.BuildFunc, width, height int) {
	wrappedBuildFunc := func() engine.Widget {
		return engine.ThemeProvider(engine.GetTheme(), buildFunc())
	}
	engine := engine.NewEngine(wrappedBuildFunc, width, height)
	engine.Run()
}


// ==================== Image ====================

func Image(src *ebiten.Image) ImageBox {
	return ImageBox{src: src}
}

type ImageBox struct {
	src          *ebiten.Image
	fit          render.ObjectFit
	width        float32
	height       float32
	borderRadius float32
	tint         color.Color
}

func (i ImageBox) Fit(f render.ObjectFit) ImageBox     { i.fit = f; return i }
func (i ImageBox) Width(v float32) ImageBox             { i.width = v; return i }
func (i ImageBox) W(v float32) ImageBox                 { return i.Width(v) }
func (i ImageBox) Height(v float32) ImageBox            { i.height = v; return i }
func (i ImageBox) H(v float32) ImageBox                 { return i.Height(v) }
func (i ImageBox) CornerRadius(v float32) ImageBox      { i.borderRadius = v; return i }
func (i ImageBox) Radius(v float32) ImageBox            { return i.CornerRadius(v) }
func (i ImageBox) Tint(c color.Color) ImageBox          { i.tint = c; return i }
func (i ImageBox) CreateElement() engine.Element            { return newBuilder(i) }
func (i ImageBox) GetKey() engine.Key                       { return engine.NilKey{} }
func (i ImageBox) unwrap() engine.Widget {
	w := widgets.Image(i.src)
	if i.fit != 0 {
		w = w.Fit(i.fit)
	}
	if i.width > 0 {
		w = w.W(i.width)
	}
	if i.height > 0 {
		w = w.H(i.height)
	}
	if i.borderRadius > 0 {
		w = w.Radius(i.borderRadius)
	}
	if i.tint != nil {
		w = w.Tint(i.tint)
	}
	return w
}

// ==================== Icon ====================

func Icon(name string) IconBox {
	return IconBox{name: name}
}

type IconBox struct {
	name     string
	size     float32
	iconColor color.Color
}

func (i IconBox) Size(v float32) IconBox       { i.size = v; return i }
func (i IconBox) Color(c color.Color) IconBox  { i.iconColor = c; return i }
func (i IconBox) CreateElement() engine.Element    { return newBuilder(i) }
func (i IconBox) GetKey() engine.Key               { return engine.NilKey{} }
func (i IconBox) unwrap() engine.Widget {
	w := widgets.Icon(i.name)
	if i.size > 0 {
		w = w.Size(i.size)
	}
	if i.iconColor != nil {
		w = w.Color(i.iconColor)
	}
	return w
}

// ==================== Scroll ====================

func Scroll(child Widget) ScrollBox {
	return ScrollBox{child: child}
}

type ScrollBox struct {
	child     Widget
	width     float32
	height    float32
	maxHeight float32
}

func (s ScrollBox) Width(v float32) ScrollBox     { s.width = v; return s }
func (s ScrollBox) Height(v float32) ScrollBox    { s.height = v; return s }
func (s ScrollBox) MaxHeight(v float32) ScrollBox { s.maxHeight = v; return s }
func (s ScrollBox) CreateElement() engine.Element     { return newBuilder(s) }
func (s ScrollBox) GetKey() engine.Key                { return engine.NilKey{} }
func (s ScrollBox) unwrap() engine.Widget {
	w := widgets.Scroll(s.child)
	if s.width > 0 {
		w = w.W(s.width)
	}
	if s.height > 0 {
		w = w.H(s.height)
	}
	if s.maxHeight > 0 {
		w = w.MaxHeight(s.maxHeight)
	}
	return w
}

// ==================== Stack / Positioned ====================

func Stack(children ...Widget) StackBox {
	return StackBox{children: children}
}

type StackBox struct {
	children []Widget
	width    float32
	height   float32
	bg       *render.Color
	radius   float32
}

func (s StackBox) Width(v float32) StackBox        { s.width = v; return s }
func (s StackBox) W(v float32) StackBox            { return s.Width(v) }
func (s StackBox) Height(v float32) StackBox       { s.height = v; return s }
func (s StackBox) H(v float32) StackBox            { return s.Height(v) }
func (s StackBox) Background(c color.Color) StackBox { s.bg = render.NewColorFrom(c); return s }
func (s StackBox) BackgroundColor(c render.Color) StackBox { s.bg = &c; return s }
func (s StackBox) CornerRadius(v float32) StackBox  { s.radius = v; return s }
func (s StackBox) Radius(v float32) StackBox        { return s.CornerRadius(v) }
func (s StackBox) CreateElement() engine.Element        { return newBuilder(s) }
func (s StackBox) GetKey() engine.Key                   { return engine.NilKey{} }
func (s StackBox) unwrap() engine.Widget {
	realChildren := make([]engine.Widget, len(s.children))
	for i, c := range s.children {
		if u, ok := c.(unwrapper); ok {
			realChildren[i] = u.unwrap()
		} else {
			realChildren[i] = c
		}
	}
	w := widgets.Stack(realChildren...)
	if s.width > 0 {
		w = w.W(s.width)
	}
	if s.height > 0 {
		w = w.H(s.height)
	}
	if s.bg != nil {
		w = w.Background(*s.bg)
	}
	if s.radius > 0 {
		w = w.Radius(s.radius)
	}
	return w
}

func Positioned(child Widget) PositionedBox {
	return PositionedBox{child: child, left: -1, top: -1, right: -1, bottom: -1}
}

type PositionedBox struct {
	child  Widget
	left   float32
	top    float32
	right  float32
	bottom float32
	width  float32
	height float32
	center bool
}

func (p PositionedBox) Left(v float32) PositionedBox   { p.left = v; return p }
func (p PositionedBox) L(v float32) PositionedBox      { return p.Left(v) }
func (p PositionedBox) Top(v float32) PositionedBox    { p.top = v; return p }
func (p PositionedBox) T(v float32) PositionedBox      { return p.Top(v) }
func (p PositionedBox) Right(v float32) PositionedBox  { p.right = v; return p }
func (p PositionedBox) R(v float32) PositionedBox      { return p.Right(v) }
func (p PositionedBox) Bottom(v float32) PositionedBox { p.bottom = v; return p }
func (p PositionedBox) B(v float32) PositionedBox      { return p.Bottom(v) }
func (p PositionedBox) Width(v float32) PositionedBox  { p.width = v; return p }
func (p PositionedBox) W(v float32) PositionedBox      { return p.Width(v) }
func (p PositionedBox) Height(v float32) PositionedBox { p.height = v; return p }
func (p PositionedBox) H(v float32) PositionedBox      { return p.Height(v) }
func (p PositionedBox) Center() PositionedBox           { p.center = true; return p }
func (p PositionedBox) CreateElement() engine.Element       { return newBuilder(p) }
func (p PositionedBox) GetKey() engine.Key                  { return engine.NilKey{} }
func (p PositionedBox) unwrap() engine.Widget {
	w := widgets.Positioned(p.child)
	if p.center {
		w = w.Center()
	} else {
		if p.left >= 0 {
			w = w.L(p.left)
		}
		if p.top >= 0 {
			w = w.T(p.top)
		}
		if p.right >= 0 {
			w = w.R(p.right)
		}
		if p.bottom >= 0 {
			w = w.B(p.bottom)
		}
	}
	if p.width > 0 {
		w = w.W(p.width)
	}
	if p.height > 0 {
		w = w.H(p.height)
	}
	return w
}

// ==================== Select ====================

func Select(options []widgets.SelectOption) SelectBox {
	return SelectBox{options: options}
}

type SelectBox struct {
	options     []widgets.SelectOption
	value       string
	placeholder string
	onChange    func(string)
	disabled    bool
	width       float32
}

func (s SelectBox) Value(v string) SelectBox         { s.value = v; return s }
func (s SelectBox) WithValue(v string) SelectBox     { return s.Value(v) }
func (s SelectBox) Placeholder(v string) SelectBox  { s.placeholder = v; return s }
func (s SelectBox) WithPlaceholder(v string) SelectBox { return s.Placeholder(v) }
func (s SelectBox) OnChange(fn func(string)) SelectBox { s.onChange = fn; return s }
func (s SelectBox) WithOnChange(fn func(string)) SelectBox { return s.OnChange(fn) }
func (s SelectBox) Disabled(v bool) SelectBox        { s.disabled = v; return s }
func (s SelectBox) WithDisabled(v bool) SelectBox    { return s.Disabled(v) }
func (s SelectBox) Width(v float32) SelectBox        { s.width = v; return s }
func (s SelectBox) WithWidth(v float32) SelectBox    { return s.Width(v) }
func (s SelectBox) CreateElement() engine.Element        { return newBuilder(s) }
func (s SelectBox) GetKey() engine.Key                   { return engine.NilKey{} }
func (s SelectBox) unwrap() engine.Widget {
	w := widgets.Select(s.options)
	w = w.SelectedValue(s.value)
	w = w.PlaceholderText(s.placeholder)
	w = w.OnChanged(s.onChange)
	w = w.SetDisabled(s.disabled)
	if s.width > 0 {
		w = w.Width_(s.width)
	}
	return w
}

// ==================== NinePatch ====================

func NinePatch(src *ebiten.Image, slice render.BorderSlice) NinePatchBox {
	return NinePatchBox{src: src, slice: slice}
}

type NinePatchBox struct {
	src    *ebiten.Image
	slice  render.BorderSlice
	width  float32
	height float32
}

func (n NinePatchBox) Width(v float32) NinePatchBox  { n.width = v; return n }
func (n NinePatchBox) Height(v float32) NinePatchBox { n.height = v; return n }
func (n NinePatchBox) CreateElement() engine.Element      { return newBuilder(n) }
func (n NinePatchBox) GetKey() engine.Key                 { return engine.NilKey{} }
func (n NinePatchBox) unwrap() engine.Widget {
	w := widgets.NinePatch(n.src, n.slice)
	if n.width > 0 {
		w = w.W(n.width)
	}
	if n.height > 0 {
		w = w.H(n.height)
	}
	return w
}

// ==================== Sprite ====================

func Sprite(frames []*ebiten.Image) SpriteBox {
	return SpriteBox{frames: frames}
}

type SpriteBox struct {
	frames        []*ebiten.Image
	frameDuration time.Duration
	loop          bool
	autoPlay      bool
	width         float32
	height        float32
}

func (s SpriteBox) FrameDuration(d time.Duration) SpriteBox { s.frameDuration = d; return s }
func (s SpriteBox) Loop(v bool) SpriteBox                   { s.loop = v; return s }
func (s SpriteBox) AutoPlay(v bool) SpriteBox               { s.autoPlay = v; return s }
func (s SpriteBox) Width(v float32) SpriteBox               { s.width = v; return s }
func (s SpriteBox) Height(v float32) SpriteBox              { s.height = v; return s }
func (s SpriteBox) CreateElement() engine.Element               { return newBuilder(s) }
func (s SpriteBox) GetKey() engine.Key                          { return engine.NilKey{} }
func (s SpriteBox) unwrap() engine.Widget {
	w := widgets.Sprite(s.frames)
	if s.frameDuration > 0 {
		w = w.FrameDuration(s.frameDuration)
	}
	w = w.Loop(s.loop)
	w = w.AutoPlay(s.autoPlay)
	if s.width > 0 {
		w = w.W(s.width)
	}
	if s.height > 0 {
		w = w.H(s.height)
	}
	return w
}

// ==================== SpriteButton ====================

func SpriteButton(normal *ebiten.Image) SpriteButtonBox {
	return SpriteButtonBox{normal: normal}
}

type SpriteButtonBox struct {
	normal   *ebiten.Image
	hover    *ebiten.Image
	pressed  *ebiten.Image
	disabled *ebiten.Image
	slice    render.BorderSlice
	label    string
	onClick  func()
	width    float32
	height   float32
	isDisabled bool
	loading    bool
}

func (b SpriteButtonBox) Hover(img *ebiten.Image) SpriteButtonBox       { b.hover = img; return b }
func (b SpriteButtonBox) Pressed(img *ebiten.Image) SpriteButtonBox     { b.pressed = img; return b }
func (b SpriteButtonBox) DisabledImg(img *ebiten.Image) SpriteButtonBox { b.disabled = img; return b }
func (b SpriteButtonBox) Slice(s render.BorderSlice) SpriteButtonBox    { b.slice = s; return b }
func (b SpriteButtonBox) Label(text string) SpriteButtonBox              { b.label = text; return b }
func (b SpriteButtonBox) OnClick(fn func()) SpriteButtonBox             { b.onClick = fn; return b }
func (b SpriteButtonBox) Width(v float32) SpriteButtonBox               { b.width = v; return b }
func (b SpriteButtonBox) Height(v float32) SpriteButtonBox              { b.height = v; return b }
func (b SpriteButtonBox) Disabled(v bool) SpriteButtonBox               { b.isDisabled = v; return b }
func (b SpriteButtonBox) Loading(v bool) SpriteButtonBox                { b.loading = v; return b }
func (b SpriteButtonBox) CreateElement() engine.Element                      { return newBuilder(b) }
func (b SpriteButtonBox) GetKey() engine.Key                                 { return engine.NilKey{} }
func (b SpriteButtonBox) unwrap() engine.Widget {
	w := widgets.SpriteButton(b.normal).
		Hover(b.hover).
		Pressed(b.pressed).
		Disabled(b.disabled).
		Slice(b.slice).
		Label(b.label).
		OnTap(b.onClick).
		W(b.width).
		H(b.height).
		Disable(b.isDisabled).
		Loading(b.loading)
	return w
}

// ==================== GameProgressBar ====================

func GameProgressBar(progress float32, track, fill *ebiten.Image) GameProgressBarBox {
	return GameProgressBarBox{progress: progress, track: track, fill: fill}
}

type GameProgressBarBox struct {
	progress float32
	track    *ebiten.Image
	fill     *ebiten.Image
	slice    render.BorderSlice
	width    float32
	height   float32
}

func (p GameProgressBarBox) Slice(s render.BorderSlice) GameProgressBarBox { p.slice = s; return p }
func (p GameProgressBarBox) Width(v float32) GameProgressBarBox             { p.width = v; return p }
func (p GameProgressBarBox) Height(v float32) GameProgressBarBox            { p.height = v; return p }
func (p GameProgressBarBox) CreateElement() engine.Element                      { return newBuilder(p) }
func (p GameProgressBarBox) GetKey() engine.Key                                 { return engine.NilKey{} }
func (p GameProgressBarBox) unwrap() engine.Widget {
	w := widgets.GameProgressBar(p.progress, p.track, p.fill).
		Slice(p.slice)
	if p.width > 0 {
		w = w.W(p.width)
	}
	if p.height > 0 {
		w = w.H(p.height)
	}
	return w
}

// ==================== Grid ====================

func Grid(cols int, gap float32, children ...Widget) Widget {
	realChildren := make([]engine.Widget, len(children))
	for i, c := range children {
		if u, ok := c.(unwrapper); ok {
			realChildren[i] = u.unwrap()
		} else {
			realChildren[i] = c
		}
	}
	return widgets.Grid(cols, gap, realChildren...)
}

// ==================== Draggable / DropTarget ====================

func Draggable(child Widget, data any) Widget {
	return widgets.Draggable(child, data)
}

func DropTarget(child Widget, onAccept func(any)) Widget {
	return widgets.DropTarget(child, onAccept)
}

// ==================== FlipCard ====================

func FlipCard(front, back Widget) FlipCardBox {
	return FlipCardBox{front: front, back: back, duration: 300 * time.Millisecond}
}

type FlipCardBox struct {
	front    Widget
	back     Widget
	duration time.Duration
}

func (f FlipCardBox) Duration(d time.Duration) FlipCardBox { f.duration = d; return f }
func (f FlipCardBox) CreateElement() engine.Element             { return newBuilder(f) }
func (f FlipCardBox) GetKey() engine.Key                        { return engine.NilKey{} }
func (f FlipCardBox) unwrap() engine.Widget {
	w := widgets.FlipCard(f.front, f.back).WithDuration(f.duration)
	return w
}

// ==================== Transform ====================

func Transform(child Widget) TransformBox {
	return TransformBox{child: child, scaleX: 1, scaleY: 1, originX: 0.5, originY: 0.5, alpha: 1}
}

type TransformBox struct {
	child       Widget
	rotation    float32
	rotateX     float32
	rotateY     float32
	perspective float32
	scaleX      float32
	scaleY      float32
	skewX       float32
	skewY       float32
	originX     float32
	originY     float32
	alpha       float32
}

func (t TransformBox) Rotate(deg float32) TransformBox     { t.rotation = deg; return t }
func (t TransformBox) RotateX(deg float32) TransformBox     { t.rotateX = deg; return t }
func (t TransformBox) RotateY(deg float32) TransformBox     { t.rotateY = deg; return t }
func (t TransformBox) Perspective(dist float32) TransformBox { t.perspective = dist; return t }
func (t TransformBox) Scale(x, y float32) TransformBox     { t.scaleX = x; t.scaleY = y; return t }
func (t TransformBox) ScaleUniform(s float32) TransformBox { t.scaleX = s; t.scaleY = s; return t }
func (t TransformBox) Skew(x, y float32) TransformBox      { t.skewX = x; t.skewY = y; return t }
func (t TransformBox) Anchor(ox, oy float32) TransformBox  { t.originX = ox; t.originY = oy; return t }
func (t TransformBox) Opacity(a float32) TransformBox      { t.alpha = a; return t }
func (t TransformBox) CreateElement() engine.Element            { return newBuilder(t) }
func (t TransformBox) GetKey() engine.Key                       { return engine.NilKey{} }
func (t TransformBox) unwrap() engine.Widget {
	w := widgets.Transform(t.child).
		Rotate(t.rotation).
		RotateX(t.rotateX).
		RotateY(t.rotateY).
		Perspective(t.perspective).
		Scale(t.scaleX, t.scaleY).
		Skew(t.skewX, t.skewY).
		Anchor(t.originX, t.originY).
		Opacity(t.alpha)
	return w
}


// ==================== Card3D ====================

func Card3D(front, back *ebiten.Image) Card3DBox {
	return Card3DBox{front: front, back: back}
}

type Card3DBox struct {
	front       *ebiten.Image
	back        *ebiten.Image
	rotateX     float32
	rotateY     float32
	perspective float32
	width       float32
	height      float32
}

func (c Card3DBox) RotateX(deg float32) Card3DBox     { c.rotateX = deg; return c }
func (c Card3DBox) RotateY(deg float32) Card3DBox     { c.rotateY = deg; return c }
func (c Card3DBox) Perspective(dist float32) Card3DBox { c.perspective = dist; return c }
func (c Card3DBox) Width(v float32) Card3DBox          { c.width = v; return c }
func (c Card3DBox) Height(v float32) Card3DBox         { c.height = v; return c }
func (c Card3DBox) W(v float32) Card3DBox              { return c.Width(v) }
func (c Card3DBox) H(v float32) Card3DBox              { return c.Height(v) }
func (c Card3DBox) CreateElement() engine.Element          { return newBuilder(c) }
func (c Card3DBox) GetKey() engine.Key                     { return engine.NilKey{} }
func (c Card3DBox) unwrap() engine.Widget {
	w := widgets.Card3D(c.front, c.back).
		RotateX(c.rotateX).
		RotateY(c.rotateY).
		Perspective(c.perspective).
		W(c.width).
		H(c.height)
	return w
}
