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

	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
	"github.com/sjm1327605995/tenon/pkg/v2/widgets"
	"github.com/sjm1327605995/tenon/yoga"
)

// ==================== 类型别名 ====================

type Widget = ui.Widget
type BuildContext = ui.BuildContext
type State = ui.State
type RouteParams = ui.RouteParams

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
	ui.BaseElement
	inner ui.Element // 底层真实的 Element
}

func newBuilder(w ui.Widget) *builder {
	b := &builder{}
	b.BaseElement.Init(b, w)
	return b
}

func (b *builder) Mount(parent ui.Element, slot int) {
	b.BaseElement.Mount(parent, slot)
	real := b.unwrap(b.GetWidget())
	b.inner = real.CreateElement()
	b.inner.Mount(b, 0)
}

func (b *builder) Update(newWidget ui.Widget) {
	b.BaseElement.Update(newWidget)
	real := b.unwrap(newWidget)
	if b.inner != nil && ui.CanUpdate(b.inner.GetWidget(), real) {
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

func (b *builder) GetChildren() []ui.Element {
	if b.inner == nil {
		return nil
	}
	return []ui.Element{b.inner}
}

func (b *builder) FindRenderObject() render.RenderObject {
	if b.inner != nil {
		return b.inner.FindRenderObject()
	}
	return nil
}

// unwrap 将声明式包装类型解包为底层 Widget。
type unwrapper interface {
	unwrap() ui.Widget
}

func (b *builder) unwrap(w ui.Widget) ui.Widget {
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
	padding    ui.EdgeInsets
	justify    yoga.Justify
	alignItems yoga.Align
	flexGrow   float32
}

func (f FlexBox) Gap(v float32) FlexBox           { f.gap = v; return f }
func (f FlexBox) Padding(v float32) FlexBox       { f.padding = ui.EdgeInsetsAll(v); return f }
func (f FlexBox) Justify(v yoga.Justify) FlexBox  { f.justify = v; return f }
func (f FlexBox) Align(v yoga.Align) FlexBox      { f.alignItems = v; return f }
func (f FlexBox) Grow(v float32) FlexBox          { f.flexGrow = v; return f }
func (f FlexBox) CreateElement() ui.Element       { return newBuilder(f) }
func (f FlexBox) GetKey() ui.Key                  { return ui.NilKey{} }
func (f FlexBox) unwrap() ui.Widget {
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
func (t TextWidget) CreateElement() ui.Element       { return newBuilder(t) }
func (t TextWidget) GetKey() ui.Key                  { return ui.NilKey{} }
func (t TextWidget) unwrap() ui.Widget {
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
}

func (b ButtonWidget) Style(v ButtonStyle) ButtonWidget { b.style = v; return b }
func (b ButtonWidget) OnClick(fn func()) ButtonWidget   { b.onClick = fn; return b }
func (b ButtonWidget) Disabled(v bool) ButtonWidget     { b.disabled = v; return b }
func (b ButtonWidget) Loading(v bool) ButtonWidget      { b.loading = v; return b }
func (b ButtonWidget) CreateElement() ui.Element        { return newBuilder(b) }
func (b ButtonWidget) GetKey() ui.Key                   { return ui.NilKey{} }
func (b ButtonWidget) unwrap() ui.Widget {
	return widgets.Button(b.label).
		Variantf(b.style).
		OnTap(b.onClick).
		SetDisabled(b.disabled).
		SetLoading(b.loading)
}

// ==================== Input ====================

func Input(placeholder string) InputWidget {
	return InputWidget{placeholder: placeholder}
}

type InputWidget struct {
	placeholder string
	onChange    func(string)
	onSubmit    func(string)
	multiline   bool
}

func (i InputWidget) OnChange(fn func(string)) InputWidget { i.onChange = fn; return i }
func (i InputWidget) OnSubmit(fn func(string)) InputWidget { i.onSubmit = fn; return i }
func (i InputWidget) Multiline() InputWidget                { i.multiline = true; return i }
func (i InputWidget) CreateElement() ui.Element             { return newBuilder(i) }
func (i InputWidget) GetKey() ui.Key                        { return ui.NilKey{} }
func (i InputWidget) unwrap() ui.Widget {
	if i.multiline {
		return widgets.Textarea("").Placeholder(i.placeholder)
	}
	w := widgets.TextField("").Placeholder(i.placeholder)
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
	child   Widget
	bg      *render.Color
	radius  float32
	border  *render.Color
	borderW float32
	padding ui.EdgeInsets
	margin  ui.EdgeInsets
	width   float32
	height  float32
	onClick func()
}

func (c ContainerBox) Background(cl color.Color) ContainerBox      { c.bg = render.NewColorFrom(cl); return c }
func (c ContainerBox) CornerRadius(v float32) ContainerBox         { c.radius = v; return c }
func (c ContainerBox) Border(cl color.Color, w float32) ContainerBox { c.border = render.NewColorFrom(cl); c.borderW = w; return c }
func (c ContainerBox) Padding(v float32) ContainerBox              { c.padding = ui.EdgeInsetsAll(v); return c }
func (c ContainerBox) Margin(v float32) ContainerBox               { c.margin = ui.EdgeInsetsAll(v); return c }
func (c ContainerBox) Width(v float32) ContainerBox                { c.width = v; return c }
func (c ContainerBox) Height(v float32) ContainerBox               { c.height = v; return c }
func (c ContainerBox) OnClick(fn func()) ContainerBox              { c.onClick = fn; return c }
func (c ContainerBox) CreateElement() ui.Element                   { return newBuilder(c) }
func (c ContainerBox) GetKey() ui.Key                              { return ui.NilKey{} }
func (c ContainerBox) unwrap() ui.Widget {
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
	w = w.Pad(c.padding).Marginf(c.margin)
	if c.width > 0 {
		w = w.W(c.width)
	}
	if c.height > 0 {
		w = w.H(c.height)
	}
	if c.onClick != nil {
		w = w.OnTap(c.onClick)
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
	curve    ui.Curve
}

func (a AnimatedBox) Duration(d time.Duration) AnimatedBox { a.duration = d; return a }
func (a AnimatedBox) EaseInOut() AnimatedBox               { a.curve = ui.EaseInOutCurve{}; return a }
func (a AnimatedBox) CreateElement() ui.Element             { return newBuilder(a) }
func (a AnimatedBox) GetKey() ui.Key                        { return ui.NilKey{} }
func (a AnimatedBox) unwrap() ui.Widget {
	ac := widgets.NewAnimatedContainer().WithChild(a.child).WithDuration(a.duration)
	if a.curve != nil {
		ac = ac.WithCurve(a.curve)
	}
	return ac
}

// ==================== Navigator ====================

func Navigator(routes map[string]ui.RouteBuilder, initial string) NavigatorBox {
	return NavigatorBox{routes: routes, initial: initial}
}

type NavigatorBox struct {
	routes     map[string]ui.RouteBuilder
	initial    string
	transition ui.PageTransition
}

func (n NavigatorBox) Transition(t ui.PageTransition) NavigatorBox { n.transition = t; return n }
func (n NavigatorBox) CreateElement() ui.Element                   { return newBuilder(n) }
func (n NavigatorBox) GetKey() ui.Key                              { return ui.NilKey{} }
func (n NavigatorBox) unwrap() ui.Widget {
	w := ui.Navigator(n.routes, n.initial)
	if n.transition != ui.TransitionNone {
		w = w.WithTransition(n.transition)
	}
	return w
}

// ==================== 国际化 ====================

func Localization(locale string, translations map[string]map[string]string, child Widget) Widget {
	return ui.NewLocalization(locale, translations, child)
}

// ==================== 便捷函数 ====================

func L(ctx BuildContext, key string) string          { return ui.L(ctx, key) }
func GetNavigator(ctx BuildContext) ui.NavigatorState { return ui.GetNavigator(ctx) }
func Push(ctx BuildContext, name string, params ...RouteParams) { ui.NavPush(ctx, name, params...) }
func Pop(ctx BuildContext)                           { ui.NavPop(ctx) }
func SetTheme(t *ui.Theme)                           { ui.SetTheme(t) }
func Run(buildFunc ui.BuildFunc, width, height int)  { tenon.Run(buildFunc, width, height) }
