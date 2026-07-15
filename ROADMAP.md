# Roadmap

Status of Tenon as a GUI toolkit — what's implemented and what's next. Honest, not a wish list.

## Done

**Core engine (`pkg/ui`)**
- Three-tree model (Node → Fiber → renderNode), reconciliation (function-pointer + key identity), keyed lists, Fragment, Portal.
- Hooks: `UseState`, `UseReducer`, `UseEffect`, `UseMemo`, `UseCallback`, `UseRef`, `UseContext`; `Memo` (shallow-prop bailout); stable setters.
- `ErrorBoundary` — catches panics thrown during a subtree's render and shows a fallback (with retry) instead of crashing the app; propagates normally when no boundary is present.
- Local re-render (only the setState'd component); incremental layout (paint-only changes don't recompute; resize recomputes only size-dependent subtrees; idle frames do no layout).
- On-demand repaint: frames repaint only on real visual change (re-render / animation / FLIP / selection / IME / focus / caret blink); otherwise a cached frame is blitted. Idle drops to **0 repaints & 0 layouts/s** (verified) while the loop still runs at refresh rate. Perf HUD via `ui.ShowStats` / F12 (repaint & layout per-second, frame paint cost).
- Text-shaping cache: per-node memoization of `wrapForWidth` / `layoutRuns` (keyed by text/font/width), so repaints and re-layouts of unchanged text skip re-shaping — cache hit `~2ns`, `0 allocs` vs `~7.6µs`, `22 allocs` uncached (benchmarked).
- Refresh-rate adaptive: animation is wall-clock `dt`-based (speed constant across 60/120/144Hz); `ui.FrameSync` (default on) ties logic TPS to the display refresh via `SyncWithFPS` so high-refresh screens get more animation steps; long-press key-repeat is time-based (TPS-independent).
- `ui.Post` for thread-safe updates from background goroutines.

**Layout**
- Yoga flexbox, absolute positioning, `WidthPct`/`HeightPct`/`Fill` (window-adaptive), scroll (`ScrollView`), overflow clip.
- `VirtualList` for large lists — `UseScroll` feeds scroll offset/viewport back to the component so only the visible window (+overscan) renders; 100k rows stay smooth.

**Rendering**
- `ebiten/vector` rounded rects / borders, `text/v2` text, images; supersampling AA; HiDPI (device-scale) rendering.
- SVG icons: `Icon`/`IconFill` render `pkg/svg` paths (stroke or fill), color inherited like text; a small built-in lucide set (`IconCheck`, `IconChevronDown`, …). Rounded-rect clipping (a rounded container clips its children to the corners via an offscreen mask). Linear gradients (`LinearGradient(from, to, angle)` background, follows the corner radius). `Img` object-fit (`Fit(FitContain/FitCover/FitFill)`).
- Paint goes through a `painter` backend interface (draw primitives + clip + layer), so the render walk is backend-agnostic — an ebiten backend for the window, a recording backend for headless golden tests (and room for a future Skia backend).

**Animation**
- `UseTween` + easings; `UseTransition` (enter/exit); FLIP layout animation; transforms (scale/rotate/translate, hit-test aware); per-node and group opacity.

**Input & text**
- Click (bubbling), hover, drag, wheel scroll; keyboard: Tab focus nav, Enter/Space activate, Esc stack; press state. Modal focus trapping (`Portal(TrapFocus(), …)`) — Tab stays inside the top modal; wired into shadcn Dialog/Sheet. Roving arrow-key navigation (`ArrowNav(NavVertical/NavHorizontal)`) inside menus/lists/tabs — wired into shadcn Tabs (←→) and DropdownMenu (↑↓).
- Controlled `Input` with caret, multi-line (`Multiline`), selection (Shift+arrows/drag/Ctrl+A, double-click word, triple-click all) and cut/copy/paste (pluggable clipboard); IME composition (`exp/textinput`) with underlined preedit at the caret. Grapheme-aware caret/backspace/delete and word-wise nav/delete (Ctrl+←→/Backspace) via `rivo/uniseg` — emoji, combining marks, ZWJ sequences move & delete as one unit.
- Text wrapping via Unicode line-breaking (UAX#14, `rivo/uniseg`) — hyphen breaks, CJK per-char, closing punctuation never at line start, non-breaking spaces; style inheritance, font weights/italic, rich-text spans (`RichText`), embedded CJK font, anchored overlays (`UseMeasure`).

**Components**
- Base kit (Checkbox/Switch/Radio/Slider/Progress/Spinner/Badge/Avatar/Divider/Card/Tabs).
- `pkg/shadcn`: ~41 shadcn/ui-style components on a theme + interaction foundation.

**Tooling**
- Live preview / hot reload (`pkg/hotreload`): edit a plain-Go `View() *ui.Node` file and the running window updates in-process — no rebuild, no restart (yaegi-interpreted; interpreted code uses non-generic `pkg/ui` + all `pkg/shadcn`, host owns state).
- Debug frame capture: `ui.Capture(path, afterFrames, exit)` or env `TENON_CAPTURE=out.png` saves the engine's own rendered frame to PNG (only the app's pixels — safe), for visually verifying rendering headlessly.
- 40+ engine tests incl. headless **golden paint tests** (`Harness.Paint()` → `[]PaintOp` via the recording backend) + wrap/measure benchmarks; per-package READMEs + godoc; runnable examples.
- `ui.Mount` headless test harness (mount + drive click/hover/press/drag/type, query the render tree) — `pkg/shadcn` uses it for real behavior tests.

## In progress / next (priority order)

1. **BiDi text** — right-to-left / mixed-direction (Arabic, Hebrew) via `x/text/unicode/bidi`. Deferred deliberately: it's an all-or-nothing change touching visual reordering, caret/selection, and hit-testing, so it needs its own careful pass rather than being bolted onto the LTR path.
2. **Accessibility** — ~~focus trapping in modals~~ **done** (`TrapFocus()`); ~~arrow-key navigation inside menus/lists~~ **done** (`ArrowNav`, roving focus). Still: an accessibility tree for screen readers (needs AccessKit/platform APIs).
3. **Performance at scale** — ~~list virtualization~~ **done** (`VirtualList` + `UseScroll` renders only the visible window; 100k rows stay smooth). Still: sub-tree-scoped `resolveInherited`.
4. **Rendering extras** — ~~SVG icons~~, ~~rounded-rect clipping~~, ~~linear gradients~~, ~~`Img` object-fit~~ **all done**. (Remaining polish: radial gradients, image filters/blur — lower priority.)
5. **Native integration** — OS clipboard binding, native file/context menus; (multi-window is bounded by Ebiten).

**Recently done:** the text layer is now feature-complete — font weights/italic, rich-text spans (`RichText`), and IME composition (`exp/textinput`, underlined preedit) join wrapping, style inheritance, and multi-line selection. A headless test-mount helper (`ui.Mount`) lets `pkg/shadcn` and app code assert real click/input/hover behavior.

## Non-goals (for now)

- VDOM-style full diffing (we use Flutter/React fiber identity + keys).
- Vue-style Proxy reactivity (Go has no Proxy; hooks cover it).

Contributions toward any of the "next" items are very welcome — see [CONTRIBUTING.md](CONTRIBUTING.md).
