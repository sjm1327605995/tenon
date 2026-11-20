# tenon/pkg/ui

A declarative, React-style GUI toolkit for Go, built on **[yoga](../../yoga)** for flexbox layout and **[Gio](https://gioui.org)** for rendering (GPU-accelerated vector paths and text).

- **HTML-like elements** — `Div`, `Span`, `Button`, `Input`, `Img`, `Text`, `ScrollView`.
- **Object components** — embed `Component`, implement `Render() NodeSource`, and read lifecycle-bound capabilities through `Context()`.
- **Reactive data** — `Value[T]` / `Bind[T]` connect text and inputs to persistent data; see the [binding guide](../value/README.md).
- **Function components and hooks** — `UseState`, `UseEffect`, `UseReducer`, `UseMemo`, `UseCallback`, `UseRef`, `UseContext`.
- **Automatic, local re-render** — a state setter re-renders only its own component (fiber), no manual invalidation.
- **Typed props** — custom components take a strongly-typed props struct via `Use[P]` / `Memo[P]`.

```
Node (immutable description)  ──reconcile──▶  Fiber (identity + hooks)  ──holds──▶  renderNode (yoga.Node + paint)
```

---

## Quick start

For an application using an embedded component and reactive bindings, start with the [project quick start](../../README.md#quick-start). The following counter demonstrates the function component API:

```go
package main

import (
	"fmt"
	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

func Counter() *ui.Node {
	count, setCount := ui.UseState(0)
	return ui.HStack(16,
		ui.Button(ui.Text("-")).OnClick(func() { setCount(count - 1) }),
		ui.Text(fmt.Sprintf("%d", count)).FontSize(32),
		ui.Button(ui.Text("+")).OnClick(func() { setCount(count + 1) }),
	).Padding(24)
}

func main() { ui.Run(ui.View(Counter)) }
```

`ui.Run(root *Node)` opens the window and drives the frame loop. Mount no-props functions with `View(...)`, or typed props with `Use(fn, props)`.

---

## Components

An object component embeds `Component`, implements `Render() NodeSource`, and mounts through `Render(instance)`. Its `Context()` provides theme, viewport, effects, and UI scheduling. The context is available after mounting and canceled on unmount; see the [complete example](../../example/component).

A function component is a plain function `func(P) *ui.Node` taking a typed props struct. Mount it with `Use` (or `Memo`):

```go
type CardProps struct{ Title string; OnPick func() }

func Card(p CardProps) *ui.Node { /* hooks + return a Node */ }

ui.Use(Card, CardProps{Title: "Hello", OnPick: pick})   // re-renders when parent does
ui.Memo(Card, CardProps{Title: "Hello", OnPick: pick})  // bails out if props are shallow-equal
```

- `Memo` compares props field by field. Non-function fields use deep equality; non-nil function props conservatively count as changed because Go cannot reliably distinguish closures with the same code pointer but different captured values.
- `Keyed(key, node)` gives a node a stable identity for list reconciliation; `Key(k)` is the attribute form.

## Elements & attributes

| Element | Notes |
|---|---|
| `Div`, `Span`, `Button` | boxes; `Button` is focusable/clickable |
| `Text(s, ...StyleOpt)` | text node; wraps when width-constrained; inherits color/size |
| `Input(...)` | controlled text field: `Value`, `OnChange`, `Placeholder` |
| `Img(Src(path))` | image (PNG/JPEG), cached |
| `ScrollView(...)` | clips overflow, mouse-wheel scroll, scrollbar |
| `Fragment(...)` | groups children without a box |
| `Portal(...)` | renders to a top-level overlay (modals, tooltips) |

Attributes (passed alongside children, any order): `Class`, `Id`, `Key`, `OnClick`, `OnHover(func(bool))`, `OnDrag(func(dx,dy float32))`, `Value`, `OnChange(func(string))`, `Placeholder`, `Src`, and `Style(...)`.

Helpers: `If(cond, node)` conditional; nil children are ignored.

## Compact syntax

Use functions for the layout tree and chain modifiers for common properties:

```go
func Profile() *ui.Node {
    name, setName := ui.UseState("")
    return ui.VStack(12,
        ui.Text("Profile").FontSize(22),
        ui.Input().Bind(name, setName).Placeholder("Name"),
        ui.Button(ui.Text("Save")).OnClick(func() { fmt.Println(name) }),
    ).Width(320).Padding(24)
}

func main() { ui.Run(ui.View(Profile)) }
```

`View` mounts a no-props function with normal Hooks; typed props still use `Use`/`Memo`. `For(items, func(item T) *ui.Node { ... })` returns a Fragment that fits directly among other children. Use stable `.Key(item.ID)` on stateful list rows, with Hooks inside their `Use` component.

Modifiers return new nodes and preserve the source. They work with hosts, text, icons, and single-root components without adding layout containers. Common methods: `Width`, `Height`, `Padding`, `Gap`, `Bg`, `Radius`, `FontSize`, `TextColor`, `OnClick`, `Bind`, `Placeholder`, `OnSubmit`, `Disabled`, `Key`. Other options remain available through `.Style(...StyleOpt)` or `.With(...attributes)`. Multiple-root fragments need an explicit container for shared layout or events. Fresh component modifiers invalidate Memo conservatively; reusing a stable modified Memo node preserves its bailout.

For themed widgets, `shadcn.TextField(binding)` accepts a `value.Bind[string]`; `shadcn.TextField(value, setter)` remains available for controlled string state. `shadcn.TextButton(label, callback)` covers text buttons. Full props constructors remain available for variants, sizes and disabled appearance. Run `go run ./example/simple`; see the [API ergonomics research](../../docs/API_ERGONOMICS.md) for the comparison and design rationale.

Prefer struct literals? The optional [`pkg/decl`](../decl/README.md) provides `Column{Children: ...}`, `Input{...}`, and other declarations. Pages return `decl.View` and start with `decl.Run(Profile)`; conversion is automatic. Existing `*ui.Node` values can be mixed directly into declaration children.

## Styling

`Style(...StyleOpt)` carries layout and appearance. Options:

- **Size**: `Width`, `Height`, `MinWidth/Height`, `MaxWidth/Height`, `WidthPct`/`HeightPct` (% of parent/viewport), `Fill` (100%×100%, window-adaptive)
- **Spacing**: `Padding(v)`, `PaddingXY(h,v)`, `Margin`, `MarginXY`, `Gap`
- **Flex**: `Row`, `Column`, `Grow`, `Shrink`, `ItemsStart/Center/End`, `JustifyStart/Center/End/Between`
- **Appearance**: `Bg(Color)`, `Radius`, `Border(w, Color)`, `Opacity`, `Clip`
- **Position**: `Absolute`, `Top/Right/Bottom/Left`
- **Transform** (around center): `Scale`, `Rotate(deg)`, `TranslateXY`
- **Text** (inherited by descendants): `TextColor(Color)`, `FontSize`, `FontWeight(int)` / `Bold` / `Semibold` / `Medium`, `Italic`. Only one face ships (OPPOSans Medium), so bold and italic are **synthesized** — bold by stroking the glyph outline, italic by shearing it. That makes weight effectively binary: `Semibold` (600) and above look the same, and 400/500 look the same.
- **Animation**: `Animated` (FLIP — slides to new position when its layout moves)

Colors: `Hex("#rrggbb"|"#rrggbbaa")`, `Color{R,G,B,A}`, `c.Alpha(f)`, plus `White/Black/Red/Green/Blue/Gray/...`.

## Hooks

```go
count, setCount := ui.UseState(0)                          // local state; setter is stable
total, dispatch := ui.UseReducer(reducer, initial)         // reducer state; dispatch is stable
ui.UseEffect(func() ui.Cleanup { ... ; return cleanup }, deps...) // after commit; cleanup on dep-change/unmount
v := ui.UseMemo(func() T { ... }, deps...)                 // memoized value
cb := ui.UseCallback(fn, deps...)                          // memoized callback
ref := ui.UseRef(initial)                                  // stable *T across renders
```

Setters/dispatch return the same function value across renders. Go function values are not comparable, though: `Memo` conservatively treats non-nil function props as changed, and effect dependency comparison does the same. Don't rely on callbacks to suppress memoized renders or effect reruns.

## Context

```go
var ThemeCtx = ui.CreateContext(lightTheme)

func App(_ struct{}) *ui.Node {
	return ThemeCtx.Provider(currentTheme, ui.Use(Page, PageProps{}))
}

func Page(_ struct{}) *ui.Node {
	theme := ui.UseContext(ThemeCtx) // re-renders on provider value change, even across Memo
	...
}
```

## Animation

```go
x := ui.UseTween(target, 200, ui.EaseInOut) // eases toward target; only this component re-renders while animating

mounted, p := ui.UseTransition(open, 200)   // keeps a node alive through its exit animation
return ui.If(mounted, ui.Div(ui.Style(ui.Opacity(p), ui.Scale(0.9+0.1*p)), ...))
```

Easings: `Linear`, `EaseIn`, `EaseOut`, `EaseInOut`. Layout (FLIP) animation is opt-in per element via the `Animated` style.

## Input

[`value.Value[T]`](../value) unifies static, derived and custom observable data; `value.Bind[T]` adds `Set`. `ui.Text(name)` accepts a `Value[string]`, and `ui.Input(ui.BindInput(name))` uses a `Bind[string]` for two-way updates. The themed shorthand is `shadcn.TextField(name)`. Nodes release subscriptions when the source changes or they unmount. Use `ui.UseValue(source)` when layout or branches depend on a value.

Click, hover, drag, wheel, and keyboard are all **transform-aware** (scaled/rotated/translated elements are hit at their visual position).

Keyboard: **Tab/Shift-Tab** cycles focus (inputs + clickable elements), **Enter/Space** activates the focused control, **Esc** blurs; the focused element shows a ring. Clicking button text or an icon focuses its control. Keyboard and programmatic focus automatically scroll the target into view, including nested scroll areas and Portal overlays. Tab into an input places the caret at the end without selecting its existing value.

`ui.Disabled(true)` prevents focus, editing, activation, and dragging on an element and its rendered descendants. Disabled controls still occupy space and block clicks from reaching underlying or enclosing controls. Disabling a focused or pressed control clears its interaction state; external value or Signal updates still render. A surrounding Tooltip can still respond to hover. Apply the desired disabled appearance with `Style`; `pkg/shadcn` controls already apply reduced opacity.

Text inputs support **selection** (Shift+arrows/Home/End, drag, Ctrl+A) and **cut/copy/paste** (Ctrl+X/C/V) with a rendered highlight. The clipboard is in-app by default; plug the OS clipboard via `ui.SetClipboardProvider(get, set)`.

## Component kit

Ready-made controls composed from the primitives (all in package `ui`):

| Component | Signature |
|---|---|
| `Checkbox` | `Checkbox(checked bool, onChange func(bool))` |
| `Radio` | `Radio(selected bool, onChange func())` |
| `Switch` | `Switch(on bool, onChange func(bool))` — animated thumb |
| `Slider` | `Slider(value, min, max float32, onChange func(float32))` — draggable |
| `ProgressBar` | `ProgressBar(value float32)` — 0..1 |
| `Spinner` | `Spinner(size float32, c Color)` — continuous spin (`UseElapsed`) |
| `Badge` | `Badge(text string, c Color)` |
| `Avatar` | `Avatar(initials string, size float32)` |
| `Divider` | `Divider()` |
| `Card` | `Card(children ...*Node)` — override with a leading `Style(...)` |
| `Tabs` | `Tabs(TabsProps{Tabs, Active, OnChange})` — tab bar |

Form controls are **controlled** (value in, `onChange` out) — hold the state with `UseState` in the parent.

## Building style libraries (extensions)

The base package ships the foundation for building restyled component libraries (e.g. `pkg/shadcn`) that mix freely with base nodes:

- **Theme tokens** — `Theme` struct (shadcn-named: `Primary`, `Muted`, `Border`, `Ring`, `Radius`, …), `LightTheme`/`DarkTheme`, `ThemeProvider(theme, children...)`, `UseTheme()`.
- **Interaction state** — `hovered, pressed, attrs := UseInteraction()`; spread `attrs` onto the element. Low-level: `OnPress(func(bool))`, `OnHover(func(bool))`, `Attrs(...)` to bundle attributes.
- **Style composition** — `Styles(...StyleOpt)` bundles a reusable variant/size; `StyleIf(cond, ...StyleOpt)` applies state styles.
- **Color** — `Mix(a, b, t)` for hover/active shades.

A styled component is just a component that restyles a primitive:

```go
func button(p ButtonProps) *ui.Node {
	th := ui.UseTheme()
	hovered, pressed, ia := ui.UseInteraction()
	style := []ui.StyleOpt{ui.Row, ui.ItemsCenter, ui.Radius(th.Radius), ui.Bg(th.Primary), ui.TextColor(th.PrimaryForeground)}
	style = append(style, ui.StyleIf(hovered, ui.Bg(ui.Mix(th.Primary, th.Background, 0.12))))
	style = append(style, ui.StyleIf(pressed, ui.Scale(0.97)))
	return ui.Button(ui.Style(style...), ui.OnClick(p.OnClick), ia /* children... */)
}
```

`pkg/shadcn` seeds this pattern with `Button` (6 variants, 4 sizes); see `example/accordion`.

## Examples

Run with `go run ./example/<name>`:

| Example | Shows |
|---|---|
| `accordion` | shadcn-style docs site with interactive previews and light/dark themes |
| `card3d` | shared-camera pseudo-3D scene and projected images |
| `component` | object components, lifecycle context, and value binding |
| `router` | stack navigation with push, pop, and replace |
| `simple` | compact fluent syntax, controlled input, and lists |
| `struct` | struct declarations and automatic component conversion |

## Testing

`ui.Mount` gives a **headless harness** — the real reconcile → layout → hit-test → event path, without opening a window — so you can assert behavior, not just that a constructor returned non-nil:

```go
h := ui.Mount(ui.View(Counter), 300, 100) // or ui.MountDefault(node)
defer h.Close()

h.Root().ByText("+").Click()          // fires the nearest onClick, bubbling like a real click
if !h.Root().ByText("1").Exists() {   // re-query after state changes the tree
    t.Fatal("counter did not increment")
}
```

`Mount` returns a `*Harness`; `Root()` (and `Overlays()` for `Portal` content) return a `*Query` you traverse and drive:

- **Find** — `ByText`, `AllByText`, `ByKind`, `ByPlaceholder`, `Clickables`, or `Find`/`FindAll(pred)`. A miss returns an empty `Query`; read/action methods on it are safe no-ops, so chains never panic (guard with `Exists()`).
- **Read** — `Text`, `AllText`, `Value`, `Placeholder`, `Kind`, `Bounds`, `Focusable`, `Clickable`, `IsFocused`, `Children`.
- **Drive a node** (each settles the tree) — `Click`, `Hover(bool)`, `Press(bool)`, `Drag(dx,dy)`, `ScrollBy(dy)`, `Focus().Type("…")`, `Backspace(n)`, `Clear`, `SetValue`.
- **Drive the app** — `Harness.ClickAt(x,y)`; `MoveTo(x,y)` drives real hit-tested hover enter/leave events, including across Portal boundaries. Keyboard: `Tab`/`ShiftTab` (returns the newly `Focused()` node), `Enter` (activate focused), `Escape` (topmost `UseEscape` or clear focus); `Resize`; and `Step(dtMs)` to advance `UseTween`/`UseTransition`/`UseElapsed`. `Overlays()` returns `Portal` content (dialogs, popovers, tooltips) for querying and driving. Components that read `Viewport()` during rendering automatically rerender on window resize.

`Query` handles are snapshots; host nodes are reused across re-renders, but after an action that adds/removes/replaces nodes, re-query from `Root()`. See `harness_test.go` (engine) and `pkg/shadcn/behavior_test.go` (a downstream consumer testing real clicks/toggles).

## Notes & limits

- **Incremental layout**: yoga child links are only rebuilt when a node's children actually change, so paint-only updates (color/hover/opacity/transform) keep yoga's cache valid and `CalculateLayout` is a no-op. On window resize, only size-dependent subtrees recompute; fixed-size subtrees are reused. Idle frames run no layout at all.
- **Crisp edges**: Gio rasterizes vector paths on the GPU with built-in antialiasing, so rounded corners, circles, and borders are smooth at any scale — there is no supersampling knob to tune. Author everything in logical pixels; the engine scales layout, fonts, and pointer deltas by the display's density.
- Rendering runs on Gio's single-threaded frame loop. Call state setters from callbacks/effects (the render goroutine). **From other goroutines** (network callbacks, timers) wrap updates in `ui.Post(func(){ ... })` — it queues the closure to run on the render goroutine before the next frame, so `setState` inside it is safe.
- Per-node `Opacity` on a container becomes a **group** opacity (composited via an offscreen layer); transforms also use a layer.
- Not yet implemented: horizontal scroll; loading additional font families at runtime (one CJK face is embedded; see the note on synthesized bold/italic above); drag-and-drop drop targets; an accessibility tree; multi-window.
