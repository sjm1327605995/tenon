# Tenon

> A fine-grained reactive UI framework for Go, powered by Ebiten and Yoga.

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.25-blue)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

**English** | [简体中文](README.zh-CN.md)

---

## ⚠️ Early Stage Notice

**This project is in early development.** Component APIs and rendering elements may undergo significant changes. The codebase is not yet stable. Use in production with caution, and expect breaking changes in future releases.

---

## What is Tenon?

Tenon is a cross-platform UI framework for Go that brings modern reactive programming patterns to desktop application development. Unlike traditional immediate-mode or VDOM-based frameworks, Tenon adopts a **fine-grained reactive architecture**.

- **Widget** describes structure — lightweight, stateful, `Render()` produces the element tree
- **Element** is the native component — persistent, reusable, holds its own Yoga node
- **State** is a signal — changes propagate directly to subscribers, bypassing Widget rebuild
- **Engine** orchestrates everything — layout (Yoga), rendering (Ebiten), events

## Features

- **Fine-Grained Reactivity** — State changes drive direct element updates, zero diff cost for property changes
- **50+ Built-in Components** — View, Text, Button, Input, Modal, Toast, ScrollView, Table, Sidebar, Drawer, Tabs, and more
- **Flexbox Layout** — Full Yoga layout engine with gap, align, justify, wrap support
- **Theme System** — Built-in light / dark / Shadcn themes, fully customizable
- **Animation Engine** — Tween animations with multiple easing functions
- **Style System** — Register global or type-scoped styles via tags and classes
- **Conditional Rendering** — `Switch` builder for clean multi-branch UI logic
- **Auto Dependency Tracking** — State accessed during `Render()` is automatically subscribed
- **Cross-Platform** — Windows, macOS, Linux via Ebiten
- **Debug Tooling** — Built-in remote debugger for inspecting the element tree

## Quick Start

```go
package main

import (
    "image/color"
    
    "github.com/sjm1327605995/tenon"
    "github.com/sjm1327605995/tenon/pkg/fonts"
    "github.com/sjm1327605995/tenon/pkg/v2/components"
    "github.com/sjm1327605995/tenon/pkg/v2/core"
    "github.com/sjm1327605995/tenon/yoga"
)

type App struct {
    core.BaseWidget
    count *core.State[int]
}

func NewApp() *App {
    a := &App{count: core.NewState(0)}
    a.Init(a)
    return a
}

func (a *App) Render() core.Element {
    return components.NewView().
        SetWidthPercent(100).
        SetHeightPercent(100).
        SetPadding(yoga.EdgeAll, 24).
        SetBackgroundColor(core.GetTheme().BackgroundColor).
        Add(
            components.NewText("Hello, Tenon!").
                SetFontSize(24).
                SetColor(core.GetTheme().TextColor),
            components.NewButton("Clicked: 0").
                SetOnClick(func() {
                    a.count.Set(a.count.Get() + 1)
                }),
        )
}

func main() {
    fonts.InitDefaultFont()
    core.SetTheme(core.DefaultShadcnLightTheme())
    
    tenon.Run(NewApp(), 800, 600)
}
```

## Architecture

```
┌─────────────────────────────────────────┐
│  Widget (structure description)         │
│  Render() → Element tree                │
│  Called on mount & structural changes   │
├─────────────────────────────────────────┤
│  Element (persistent node)              │
���  *View, *Text, *Button...               │
│  Holds Yoga node, handles draw + events │
│  Chainable setters apply immediately    │
├─────────────────────────────────────────┤
│  State (fine-grained signal)            │
│  Set() notifies subscribers             │
│  Auto-tracked during Render()           │
├─────────────────────────────────────────┤
│  Engine                                 │
│  Build queue → Yoga layout → Draw loop  │
│  Event routing → Animation tick         │
└─────────────────────────────────────────┘
```

### Two Update Paths

**Path A — Property Update (high frequency, zero diff)**
```
user action → count.Set(42) → State notifies subscribers
→ text.SetContent("42") → MarkNeedDraw → next frame redraw
```
No `Render()` call. No diff. No Yoga recalculation. Only the subscribed element updates.

**Path B — Structure Update (low frequency, sibling shallow diff)**
```
page change → RequestBuild() → Render() new tree
→ sibling type comparison → reuse / replace / move Elements
→ mark Yoga dirty → next frame relayout
```
Only direct children are compared by type. Properties are not diffed — they are handled by Path A.

## Component Gallery

Tenon ships with a comprehensive component library covering the full spectrum of desktop UI needs:

**Layout** — `View`, `ScrollView`, `SplitView`, `Resizable`, `AspectRatio`, `Sidebar`, `Sheet`, `Drawer`

**Display** — `Text`, `Image`, `SVGIcon`, `Badge`, `Divider`, `Kbd`, `Skeleton`, `Table`, `Carousel`

**Form Controls** — `Button` (5 variants), `TextInput`, `TextArea`, `Select`, `Checkbox`, `Radio`, `RadioGroup`, `Switch`, `Slider`, `InputOTP`

**Feedback** — `Modal`, `Alert`, `AlertDialog`, `Toast`, `Tooltip`, `Popover`, `ProgressBar`, `LoadingSpinner`

**Navigation** — `Tab`, `Menu`, `Menubar`, `Breadcrumb`, `Pagination`, `Command`, `NavigationMenu`

**Overlay** — `Dropdown`, `ContextMenu`, `HoverCard`, `FloatingButton`

**Data** — `Accordion`, `Collapsible`, `Calendar`, `ListView`

## Theming

Tenon includes three built-in themes and full customization support:

```go
// Use a built-in theme
core.SetTheme(core.DefaultShadcnLightTheme())
core.SetTheme(core.DefaultShadcnDarkTheme())
core.SetTheme(core.DefaultAntTheme())

// Or build your own
theme := &core.Theme{
    PrimaryColor:      color.RGBA{59, 130, 246, 255},
    BackgroundColor:   color.RGBA{255, 255, 255, 255},
    TextColor:         color.RGBA{15, 23, 42, 255},
    BorderRadius:      8,
    // ... all theme fields
}
core.SetTheme(theme)
```

All components read from the current theme by default and can be individually overridden via chainable setters.

## Animation

```go
// Tween animation
tween := core.NewTween(300*time.Millisecond, core.EaseOutCubic).
    OnUpdate(func(p float32) {
        el.SetOpacity(p)
    }).
    OnComplete(func() {
        // animation done
    })
tween.Start()

// State-driven transitions
core.NewTween(200*time.Millisecond, core.EaseInOut).
    OnUpdate(func(p float32) {
        x := core.LerpFloat32(0, 200, p)
        el.SetPosition(x, 0)
    }).Start()
```

## Styling

Register reusable styles globally or per element type:

```go
// Global style by class
tenon.RegisterStyle("card", func(e tenon.Element) {
    if v, ok := e.(*components.View); ok {
        v.SetBackgroundColor(core.GetTheme().CardColor).
          SetBorderRadius(core.GetTheme().BorderRadius).
          SetShadow(color.RGBA{A: 20}, 12, 0, 2)
    }
})

// Usage in Render()
components.NewView().SetClass("card").Add(...)
```

## Examples

Explore the [`example/`](example/) directory:

| Example | Description |
|---------|-------------|
| [`gallery`](example/gallery) | Complete component showcase with 50+ UI components |
| [`v2-demo`](example/v2-demo) | Full component gallery with navigation |
| [`shadcn-gallery`](example/shadcn-gallery) | Shadcn/UI-inspired styling showcase |
| [`card`](example/card) | Card layout demonstration |

```bash
cd example/gallery
go run main.go
```

## Documentation

- **[Architecture Deep Dive](ARCHITECTURE.md)** — Detailed design decisions, update mechanisms, and internal structure
- **[Contributing Guide](CONTRIBUTING.md)** — Development setup, code conventions, and PR workflow
- **[API Reference](pkg/v2/core/)** — Core interfaces: `Widget`, `Element`, `State`, `Engine`, `Theme`, `Animation`

## Project Structure

```
tenon/
├── tenon.go              # Public API entry: Run(), types, helpers
├── pkg/
│   ├── v2/
│   │   ├── core/         # Engine, Widget, Element, State, Theme, Animation
│   │   └── components/   # 50+ built-in UI components
│   └── fonts/            # Font loading and glyph management
├── yoga/                 # Yoga Flexbox layout engine (Go port)
├── example/              # Demo applications
├── ARCHITECTURE.md       # Architecture documentation
├── CONTRIBUTING.md       # Contribution guidelines
└── LICENSE               # MIT License
```

## Why Tenon?

| | Immediate Mode (imgui) | VDOM (React-like) | **Tenon** |
|---|:---:|:---:|:---:|
| State → UI | Manual every frame | Diff + patch | Direct signal subscription |
| Widget rebuild per state change | N/A | Full tree re-render | Auto-tracked, minimal rebuild |
| Element persistence | Recreated each frame | Reconciled | Persistent, reusable |
| Performance model | CPU-bound draw calls | Diff overhead | Zero-diff property updates |
| Go ergonomics | Procedural | Hooks everywhere | Idiomatic Go, struct-based |

Tenon sits at a unique intersection: it offers the **performance characteristics of fine-grained reactivity** (à la Solid.js, Svelte) with the **ergonomics of Go structs and interfaces**, rendering via a high-performance 2D engine.

## License

MIT License — see [LICENSE](LICENSE) for details.
