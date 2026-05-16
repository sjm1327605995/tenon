# Tenon

A pure-Go UI framework built on [Ebiten](https://ebiten.org/) and a custom [Yoga](https://www.yogalayout.dev/) layout engine.

> **Note:** This project is a fork/simplification of the `gogpu/ui` architecture, replacing the WebGPU rendering stack with Ebiten v2 for broader compatibility and simpler deployment.

## Architecture

| Layer | Technology |
|-------|-----------|
| **Render Backend** | Ebiten v2 (`github.com/hajimehoshi/ebiten/v2`) |
| **Layout Engine** | Pure-Go Yoga (`yoga/` package) |
| **Widget System** | Retained-mode with pluggable Painters |
| **State Management** | Reactive signals (`state/` package) |
| **Animation** | Time-based interpolation (`animation/` + `transition/`) |
| **Accessibility** | ARIA roles + actions (`a11y/` package) |

## Features

### Core Widgets (24 packages)

| Widget | Status | Theme Painters |
|--------|--------|---------------|
| Button | ✅ | M3, Cupertino, Fluent |
| Checkbox | ✅ | M3, Cupertino, Fluent |
| Radio | ✅ | M3, Cupertino, Fluent |
| TextField | ✅ | M3, Cupertino, Fluent |
| Slider | ✅ | M3, Cupertino, Fluent |
| Progress / ProgressBar | ✅ | M3, Cupertino, Fluent |
| TabView | ✅ | M3, Cupertino, Fluent |
| Dialog | ✅ | M3, Cupertino, Fluent |
| Dropdown | ✅ | M3, Cupertino, Fluent |
| Popover | ✅ | M3, Cupertino, Fluent |
| ScrollView | ✅ | M3 |
| ListView | ✅ | M3 |
| GridView | ✅ | M3 |
| SplitView | ✅ | M3 |
| Docking | ✅ | M3 |
| Collapsible | ✅ | M3 |
| LineChart | ✅ | M3 |
| **Menu** (Bar + Context) | ✅ | M3 |
| **TreeView** | ✅ | M3 |
| **DataTable** | ✅ | M3 |
| **Toolbar** | ✅ | M3 |
| **TitleBar** | ✅ | M3 |
| **Stripe** | ✅ | M3 |

### Theme Systems

- **Material 3** — Full color scheme (HCT-based), type scale, shape scale
- **Cupertino** — iOS-style controls
- **Fluent** — Windows-style controls

### Layout Primitives

- `primitives.Row` / `primitives.Column` — Flexbox via Yoga
- `primitives.Box` — Single-child container with padding/background
- `primitives.Text` — Text rendering with font cache

## Quick Start

```bash
go get github.com/sjm1327605995/tenon
```

```go
package main

import (
    "log"

    "github.com/sjm1327605995/tenon/app"
    "github.com/sjm1327605995/tenon/core/button"
    "github.com/sjm1327605995/tenon/primitives"
    "github.com/sjm1327605995/tenon/theme/material3"
    "github.com/sjm1327605995/tenon/widget"
)

func main() {
    m3 := material3.New(widget.Hex(0x6750A4))
    bp := material3.ButtonPainter{Theme: m3}

    btn := button.New(
        button.Text("Hello Tenon"),
        button.VariantOpt(button.Filled),
        button.PainterOpt(bp),
        button.OnClick(func() {
            log.Println("Clicked!")
        }),
    )

    root := primitives.NewColumn(primitives.NewBox(btn))
    a := app.New(root)

    if err := a.Run(); err != nil {
        log.Fatal(err)
    }
}
```

## Examples

```bash
# Gallery demo with Material3 widgets
go run ./examples/gallery
```

## Project Status

This is an active work-in-progress. The following are implemented and working:

- ✅ 24 core widgets with pluggable painters
- ✅ 3 design system themes (Material 3, Cupertino, Fluent)
- ✅ Reactive signal-based state management
- ✅ Full mouse/keyboard/touch input dispatch
- ✅ Overlay system (modal dialogs, popovers, context menus)
- ✅ Animation and transition framework
- ✅ Accessibility foundation (ARIA roles, actions, tree)
- ✅ Ebiten v2 game loop integration

### Known Limitations

- Partial redraw / dirty-region optimization not yet implemented (full-screen clear each frame)
- RepaintBoundary scene caching inactive
- Limited unit test coverage
- No drag-and-drop (OS-level DnD)

## License

MIT
