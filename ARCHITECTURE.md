# Architecture

Tenon uses React's model, adapted to Go (no JSX, no Proxy). Three trees:

```
Node            immutable description you return from a component (a "React element")
  │  reconcile (function-pointer + key identity, like Flutter canUpdate)
Fiber           persistent identity: holds hooks, child fibers; the mutable layer
  │  host/text fibers own a
renderNode      a yoga.Node + paint data; the only thing that enters layout/paint
```

- **Node** (`node.go`) — cheap value objects built each render: host elements (`Div`/`Button`/…), text, components (`Use`/`Memo`), attributes. Immutable.
- **Fiber** (`fiber.go`) — reused across renders. A component fiber stores its hook slots; reconciliation reuses a fiber when the node's identity matches (same function pointer / tag + key), otherwise unmounts and remounts.
- **renderNode** (`render.go`) — wraps a retained `yoga.Node` and the paint state (bg, radius, text face, transform…). Component/Fragment/Portal fibers have none; the yoga tree is built from the host/text render nodes only.

## Hooks

Rendering is single-threaded (Gio's frame loop), so the hook dispatcher is just a package-level "current fiber" + slot index (`hooks.go`). `UseState` stores state on the fiber by call order; the returned setter closes over the fiber and, when called, marks it dirty. Rules of Hooks (stable call order) apply.

## Frame loop (`run.go`)

Each `Update`:
1. drain `ui.Post` queue (cross-goroutine updates) → input → tick animations.
2. drain the dirty-fiber queue: re-run each dirty component, reconcile its subtree.
3. `layout()` — **incremental**: `relink` only rebuilds a host's yoga children when they actually changed (so paint-only changes keep yoga's cache valid); `CalculateLayout` runs only when yoga is dirty or the window resized; `computeBounds`/`syncMeasures` run only when bounds actually changed (or on scroll).
4. `Draw` paints the renderNode tree, then portal overlays; the scene is rendered at `DeviceScaleFactor × SuperSample` and downscaled for antialiased, HiDPI-crisp edges.

## Layers on top

Styling (`Style(...StyleOpt)`), animation (`UseTween`/`UseTransition`/FLIP), input (transform-aware hit-testing, focus/Esc stacks, controlled inputs), theming (`ThemeProvider`/`UseTheme`), and the component libraries (`pkg/ui` kit, `pkg/shadcn`) are all built on these three trees. See `pkg/ui/README.md` for the user-facing API.

> Historical design notes (a pre-rewrite Flutter/SwiftUI-style exploration) live in `docs/archive/` and `docs/` — they do **not** describe the current code.
