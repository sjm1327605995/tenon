// Package decl provides struct-based UI declarations that build ordinary ui.Node
// trees. Declarations are values; state and component lifecycle stay in pkg/ui.
package decl

import (
	"fmt"
	"reflect"

	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

// View is a UI declaration. Declaration values and *ui.Node can be mixed in the
// same tree. The runtime converts them automatically; callers do not need Build.
type View = ui.NodeSource

// Component is embedded in object components; implement Render() View and use
// the inherited Context() accessor for state, providers, and lifecycle.
type Component = ui.Component

type IComponent = ui.IComponent
type Context = ui.RenderContext

// Render mounts an IComponent and dispatches its Render method through the interface.
func Render(component IComponent) *ui.Node { return ui.Render(component) }

// State declares typed state scoped to a component context.
func State[T any](ctx *Context, initial T) (T, func(T)) { return ui.StateIn(ctx, initial) }

// Consume reads a typed provider using a component context.
func Consume[T any](ctx *Context, key *ui.Context[T]) T { return ui.ContextValueIn(ctx, key) }

// Func mounts a no-props declaration function with normal Hooks and stable
// function identity, converting its return value automatically.
func Func(render func() View) *ui.Node { return ui.ViewSource(render) }

// Use mounts a declaration function with typed props.
func Use[P any](render func(P) View, props P) *ui.Node { return ui.UseSource(render, props) }

// Memo mounts a declaration function with typed props and shallow-prop bailout.
func Memo[P any](render func(P) View, props P) *ui.Node { return ui.MemoSource(render, props) }

// Run starts a declaration-based application. The runtime handles conversion
// after each render, including renders caused by Hooks updates.
func Run(root any) { ui.Run(rootNode(root)) }

// Mount drives a declaration-based page in the normal headless test harness.
func Mount(root any, w, h int) *ui.Harness {
	return ui.Mount(rootNode(root), w, h)
}

func rootNode(root any) *ui.Node {
	switch value := root.(type) {
	case nil:
		return nil
	case IComponent:
		return Render(value)
	case func() View:
		return Func(value)
	case View:
		return build(value)
	default:
		panic(fmt.Sprintf("decl: root must be an IComponent, View, or func() View; got %T", root))
	}
}

func build(v View) *ui.Node {
	if v == nil {
		return nil
	}
	if node, ok := v.(*ui.Node); ok {
		return node
	}
	if value := reflect.ValueOf(v); value.Kind() == reflect.Ptr && value.IsNil() {
		return nil
	}
	if component, ok := v.(IComponent); ok {
		return Render(component)
	}
	return v.Build()
}

func children(views []View) []*ui.Node {
	nodes := make([]*ui.Node, 0, len(views))
	for _, v := range views {
		if n := build(v); n != nil {
			nodes = append(nodes, n)
		}
	}
	return nodes
}

func finish(n *ui.Node, key string, style []ui.StyleOpt) *ui.Node {
	if len(style) > 0 {
		n = n.Style(style...)
	}
	if key != "" {
		n = n.Key(key)
	}
	return n
}

// Box describes a container. Zero dimensions use automatic sizing; zero colors
// inherit the default. Style is applied last and can express explicit zero sizes
// and advanced options without adding another container.
type Box struct {
	Key                 string
	Width, Height       float32
	WidthPct, HeightPct float32
	MinWidth, MinHeight float32
	MaxWidth, MaxHeight float32
	Padding, Gap        float32
	Grow, Shrink        float32
	Radius              float32
	Background          ui.Color
	Disabled            bool
	OnClick             func()
	Style               []ui.StyleOpt
	Attrs               []*ui.Node
	Children            []View
}

func (b Box) styles(defaults ...ui.StyleOpt) []ui.StyleOpt {
	opts := append([]ui.StyleOpt(nil), defaults...)
	for _, p := range []struct {
		value float32
		opt   func(float32) ui.StyleOpt
	}{
		{b.Width, ui.Width}, {b.Height, ui.Height},
		{b.WidthPct, ui.WidthPct}, {b.HeightPct, ui.HeightPct},
		{b.MinWidth, ui.MinWidth}, {b.MinHeight, ui.MinHeight},
		{b.MaxWidth, ui.MaxWidth}, {b.MaxHeight, ui.MaxHeight},
		{b.Padding, ui.Padding}, {b.Gap, ui.Gap},
		{b.Grow, ui.Grow}, {b.Shrink, ui.Shrink}, {b.Radius, ui.Radius},
	} {
		if p.value != 0 {
			opts = append(opts, p.opt(p.value))
		}
	}
	if b.Background != (ui.Color{}) {
		opts = append(opts, ui.Bg(b.Background))
	}
	return append(opts, b.Style...)
}

func (b Box) container(makeNode func(...*ui.Node) *ui.Node, defaults ...ui.StyleOpt) *ui.Node {
	args := []*ui.Node{ui.Style(b.styles(defaults...)...), ui.Disabled(b.Disabled), ui.OnClick(b.OnClick)}
	args = append(args, b.Attrs...)
	args = append(args, children(b.Children)...)
	return finish(makeNode(args...), b.Key, nil)
}

// Build creates a plain container using the existing ui.Div layout behavior.
func (b Box) Build() *ui.Node { return b.container(ui.Div) }

// Column describes a vertical container with the same fields as Box.
type Column Box

// Build creates a vertical container without any adapter host nodes.
func (c Column) Build() *ui.Node { return Box(c).container(ui.Div, ui.Column) }

// Row describes a horizontal container. Children are vertically centered by
// default, matching ui.HStack; Style can override alignment.
type Row Box

// Build creates a horizontal container.
func (r Row) Build() *ui.Node { return Box(r).container(ui.Div, ui.Row, ui.ItemsCenter) }

// ScrollView describes a vertical scroll area with the same fields as Box.
type ScrollView Box

// Build creates a scroll area using the existing scroll and clipping behavior.
func (s ScrollView) Build() *ui.Node { return Box(s).container(ui.ScrollView, ui.Column) }

// Fragment groups declarations without adding a layout container.
type Fragment struct {
	Key      string
	Children []View
}

// Build creates a normal ui.Fragment, preserving keys on its children.
func (f Fragment) Build() *ui.Node { return finish(ui.Fragment(children(f.Children)...), f.Key, nil) }

// ThemeProvider supplies a theme to its declared or existing-node children.
type ThemeProvider struct {
	Theme    ui.Theme
	Children []View
}

// Build creates the existing provider; component Hooks and context remain on
// their original fibers.
func (p ThemeProvider) Build() *ui.Node { return ui.ThemeProvider(p.Theme, children(p.Children)...) }

var (
	_ View = Box{}
	_ View = Column{}
	_ View = Row{}
	_ View = ScrollView{}
	_ View = Fragment{}
	_ View = ThemeProvider{}
	_ View = (*ui.Node)(nil)
)
