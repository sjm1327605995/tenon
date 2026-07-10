# Roadmap

Status of Tenon as a GUI toolkit — what's implemented and what's next. Honest, not a wish list.

## Done

**Core engine (`pkg/ui`)**
- Three-tree model (Node → Fiber → renderNode), reconciliation (function-pointer + key identity), keyed lists, Fragment, Portal.
- Hooks: `UseState`, `UseReducer`, `UseEffect`, `UseMemo`, `UseCallback`, `UseRef`, `UseContext`; `Memo` (shallow-prop bailout); stable setters.
- Local re-render (only the setState'd component); incremental layout (paint-only changes don't recompute; resize recomputes only size-dependent subtrees; idle frames do no layout).
- `ui.Post` for thread-safe updates from background goroutines.

**Layout**
- Yoga flexbox, absolute positioning, `WidthPct`/`HeightPct`/`Fill` (window-adaptive), scroll (`ScrollView`), overflow clip.

**Rendering**
- `ebiten/vector` rounded rects / borders, `text/v2` text, images; supersampling AA; HiDPI (device-scale) rendering.

**Animation**
- `UseTween` + easings; `UseTransition` (enter/exit); FLIP layout animation; transforms (scale/rotate/translate, hit-test aware); per-node and group opacity.

**Input & text**
- Click (bubbling), hover, drag, wheel scroll; keyboard: Tab focus nav, Enter/Space activate, Esc stack; press state.
- Controlled `Input` with caret, multi-line (`Multiline`), selection (Shift+arrows/drag/Ctrl+A) and cut/copy/paste (pluggable clipboard).
- Text wrapping (latin word-break + CJK), style inheritance, embedded CJK font, anchored overlays (`UseMeasure`).

**Components**
- Base kit (Checkbox/Switch/Radio/Slider/Progress/Spinner/Badge/Avatar/Divider/Card/Tabs).
- `pkg/shadcn`: ~41 shadcn/ui-style components on a theme + interaction foundation.

**Tooling**
- 28 engine tests; per-package READMEs + godoc; runnable examples.

## In progress / next (priority order)

1. **Font weights & richer text** — real bold/italic (load or synthesize), rich-text spans in one `Text`, IME preedit/composition, multi-line selection highlight. *(Biggest remaining gap for text-heavy apps.)*
2. **Component behavior tests** — expose a test-mount helper from `pkg/ui` so `pkg/shadcn` can assert real click/input/hover behavior, not just construction.
3. **Accessibility** — arrow-key navigation inside menus/lists, focus trapping in modals; an accessibility tree.
4. **Performance at scale** — benchmark suite; list virtualization for large data; sub-tree-scoped `resolveInherited`.
5. **Rendering extras** — box-shadow, gradients, `Img` object-fit, integrate `pkg/svg` for icons.
6. **Native integration** — OS clipboard binding, native file/context menus; (multi-window is bounded by Ebiten).

## Non-goals (for now)

- VDOM-style full diffing (we use Flutter/React fiber identity + keys).
- Vue-style Proxy reactivity (Go has no Proxy; hooks cover it).

Contributions toward any of the "next" items are very welcome — see [CONTRIBUTING.md](CONTRIBUTING.md).
