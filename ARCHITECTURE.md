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
3. `layout()` — **incremental**: mount/unmount and keyed reordering invalidate the structure cache; content-only updates skip the whole `relink` walk, while structural changes rebuild only the host yoga children that actually changed. The Portal list is cached under the same structural tracking. Inherited text styles resolve only in subtrees whose context, styles, or device scale changed; `CalculateLayout` runs only when yoga is dirty or the window resized; `computeBounds`/`syncMeasures` run only when bounds actually changed (or on scroll).
4. `Draw` paints the renderNode tree, then portal overlays. Static frames replay a cached Gio macro; hit-testing is cached by cursor position and layout version. Gio rasterizes antialiased vector paths and text at the display's device scale for HiDPI-crisp edges.

## Layers on top

Styling (`Style(...StyleOpt)`), animation (`UseTween`/`UseTransition`/FLIP), input (transform-aware hit-testing, focus/Esc stacks, controlled inputs), theming (`ThemeProvider`/`UseTheme`), and the component libraries (`pkg/ui` kit, `pkg/shadcn`) are all built on these three trees. See `pkg/ui/README.md` for the user-facing API.

## Rendering backends

`Backend`, `Painter`, `FontFace`, `Bitmap`, `VecPath`, `InputSource`, and `BackendRuntime` form the public backend contract (`pkg/ui/backend.go`, `pkg/ui/painter.go`, `pkg/ui/input.go`). The optional `BackendCapabilityProvider` reports support for SVG paths, rounded clipping, layer transforms, exact projection and basic accessibility semantics through `CurrentBackendCapabilities`; undeclared capabilities remain false. Gio registers itself as the default backend and retains frame caching, IME, vector paths, and exact projection. An optional backend calls `SetBackend` before `ui.Run`, owns the native window loop, drives reconciliation with `BackendRuntime.Frame`, and paints both the main tree and Portals with `BackendRuntime.Paint`. Gio's `RunWindows` creates independent window runtimes with separate input tags, IME, clipboard and wake registration; frame submission is serialized because the component renderer is single-threaded.

The `backend/g3n` directory is a separate Go module for game-level 3D. Its native `Game` exposes G3N's scene, camera, window, and renderer. Game objects live under the G3N scene root; Tenon panels live under a separate UI root, so UI updates do not clear game meshes. Applications declare scenes with backend-neutral `GameScene3D`, `Mesh3D`, `PerspectiveCamera3D`, `Light3D`, material and transform Props from `pkg/ui`; `g3n.Scene3D` is the adapter. `Scene3DPluginRegistry` resolves those declarations through backend adapters, while `Plugin3D(name, payload)` adds third-party node types without changing the core enum. `Plugin3DVersion` and `RegisterPluginVersion` require exact positive-version matches, `Scene3DValidator` rejects invalid payloads before adaptation, and `SupportsPlugin` queries registered capabilities without executing plugin code. `Scene3DRuntime` provides the same update and reverse-close lifecycle to built-in and plugin components. Close keeps declarations for remounting; `Activate` resets the runtime before another mount, so all backends share `Mount → Update → Close → Mount` semantics. Each registry is scene-scoped, not a global singleton. The adapter owns per-scene geometry, image/texture and material caches, value-bound transforms, pointer dispatch, and deterministic cleanup. Cache entries and each mounted mesh hold explicit references; mesh close cannot invalidate a sibling, and `GeometryShape` remains application-owned. `BackendRuntime.Close` is idempotent and restores the previous input source. `GameScene` and native components remain available as an advanced escape hatch. Each scene has its own root, so closing one cannot remove native objects or another scene. The adapter needs no extra asset downloads or linker flags beyond G3N's own platform requirements.

> Historical design notes (a pre-rewrite Flutter/SwiftUI-style exploration) live in `docs/archive/` and `docs/` — they do **not** describe the current code.
