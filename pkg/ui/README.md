# tenon/pkg/ui

A declarative, React-style GUI toolkit for Go, built on **[yoga](../../yoga)** for flexbox layout and **[Ebiten](https://ebiten.org)** (`vector` + `text/v2`) for rendering.

- **HTML-like elements** — `Div`, `Span`, `Button`, `Input`, `Img`, `Text`, `ScrollView`.
- **React-latest hooks** — `UseState`, `UseEffect`, `UseReducer`, `UseMemo`, `UseCallback`, `UseRef`, `UseContext`.
- **Automatic, local re-render** — a state setter re-renders only its own component (fiber), no manual invalidation.
- **Typed props** — custom components take a strongly-typed props struct via `Use[P]` / `Memo[P]`.

```
Node (immutable description)  ──reconcile──▶  Fiber (identity + hooks)  ──holds──▶  renderNode (yoga.Node + paint)
```

---

## Quick start

```go
package main

import (
	"fmt"
	ui "github.com/sjm1327605995/tenon/pkg/ui"
)

func Counter(_ struct{}) *ui.Node {
	count, setCount := ui.UseState(0)
	return ui.Div(
		ui.Style(ui.Row, ui.Gap(16), ui.Padding(24), ui.ItemsCenter),
		ui.Button(ui.OnClick(func() { setCount(count - 1) }), ui.Text("-")),
		ui.Text(fmt.Sprintf("%d", count), ui.FontSize(32)),
		ui.Button(ui.OnClick(func() { setCount(count + 1) }), ui.Text("+")),
	)
}

func main() { ui.Run(ui.Use(Counter, struct{}{})) }
```

`ui.Run(root *Node)` opens the window and drives the frame loop. `root` is usually a `Use(...)` component.

---

## Components

A component is a plain function `func(P) *ui.Node` taking a typed props struct. Mount it with `Use` (or `Memo`):

```go
type CardProps struct{ Title string; OnPick func() }

func Card(p CardProps) *ui.Node { /* hooks + return a Node */ }

ui.Use(Card, CardProps{Title: "Hello", OnPick: pick})   // re-renders when parent does
ui.Memo(Card, CardProps{Title: "Hello", OnPick: pick})  // bails out if props are shallow-equal
```

- `Memo` shallow-compares props: function fields compare by identity, so stable callbacks (hook setters, `UseCallback`) don't defeat the bail-out.
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

## Styling

`Style(...StyleOpt)` carries layout and appearance. Options:

- **Size**: `Width`, `Height`, `MinWidth/Height`, `MaxWidth/Height`, `WidthPct`/`HeightPct` (% of parent/viewport), `Fill` (100%×100%, window-adaptive)
- **Spacing**: `Padding(v)`, `PaddingXY(h,v)`, `Margin`, `MarginXY`, `Gap`
- **Flex**: `Row`, `Column`, `Grow`, `Shrink`, `ItemsStart/Center/End`, `JustifyStart/Center/End/Between`
- **Appearance**: `Bg(Color)`, `Radius`, `Border(w, Color)`, `Opacity`, `Clip`
- **Position**: `Absolute`, `Top/Right/Bottom/Left`
- **Transform** (around center): `Scale`, `Rotate(deg)`, `TranslateXY`
- **Text** (inherited by descendants): `TextColor(Color)`, `FontSize`
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

Setters/dispatch have stable identity across renders (like React), so they're safe as `Memo` props and effect deps.

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

Click, hover, drag, wheel, and keyboard are all **transform-aware** (scaled/rotated/translated elements are hit at their visual position).

Keyboard: **Tab/Shift-Tab** cycles focus (inputs + clickable elements), **Enter/Space** activates the focused control, **Esc** blurs; the focused element shows a ring.

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

`pkg/shadcn` seeds this pattern with `Button` (6 variants, 4 sizes); see `example/shadcn-demo`.

## Examples

Run with `go run ./example/<name>`:

| Example | Shows |
|---|---|
| `hooks-counter` | state + effect, local re-render |
| `hooks-app` | todo: Context theme, controlled `Input`, `UseReducer`, keyed + `Memo` list, `ScrollView` |
| `hooks-anim` | `UseTween` collapse (height + fade) |
| `hooks-hover` | `OnHover` + tween scale spring |
| `hooks-drag` | `OnDrag` + `TranslateXY`, transform-aware hit-testing |
| `hooks-modal` | `Portal` dialog with enter/exit transition |
| `hooks-reorder` | FLIP list-reorder animation |
| `hooks-text` | wrapping + style inheritance |
| `hooks-kit` | component kit: Checkbox/Switch/Radio/Slider/Progress/Badge/Avatar/Spinner/Tabs/Card |

## Notes & limits

- **Incremental layout**: yoga child links are only rebuilt when a node's children actually change, so paint-only updates (color/hover/opacity/transform) keep yoga's cache valid and `CalculateLayout` is a no-op. On window resize, only size-dependent subtrees recompute; fixed-size subtrees are reused. Idle frames run no layout at all.
- **Crisp edges**: the scene is rendered at `DeviceScaleFactor × SuperSample` (default 2×, capped at 2.5×) and downscaled by Ebiten, so rounded corners, circles, and borders are antialiased. Author everything in logical pixels; the engine scales layout, fonts, and pointer deltas. Lower `ui.SuperSample` (e.g. to 1) to trade sharpness for performance.
- Rendering runs on Ebiten's single Update/Draw loop. Call state setters from callbacks/effects (the render goroutine). **From other goroutines** (network callbacks, timers) wrap updates in `ui.Post(func(){ ... })` — it queues the closure to run on the render goroutine before the next frame, so `setState` inside it is safe. See `example/hooks-async`.
- Per-node `Opacity` on a container becomes a **group** opacity (composited via an offscreen layer); transforms also use a layer.
- Not yet implemented: `Img` object-fit, horizontal scroll, rich-text spans, multiple font families/weights (one embedded CJK face is the default).
