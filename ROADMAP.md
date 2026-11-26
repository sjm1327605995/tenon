# Roadmap

Status of Tenon as a GUI toolkit — what's implemented and what's next. Honest, not a wish list.

## Done

**Core engine (`pkg/ui`)**
- Three-tree model (Node → Fiber → renderNode), reconciliation (function-pointer + key identity), keyed lists, Fragment, Portal.
- Hooks: `UseState`, `UseReducer`, `UseEffect`, `UseMemo`, `UseCallback`, `UseRef`, `UseContext`; `Memo` (shallow-prop bailout); stable setters.
- `ErrorBoundary` — catches panics thrown during a subtree's render and shows a fallback (with retry) instead of crashing the app; propagates normally when no boundary is present.
- Local re-render (only the setState'd component); incremental layout (paint-only changes don't recompute; resize recomputes only size-dependent subtrees; idle frames do no layout).
- On-demand repaint: frames repaint only on real visual change (re-render / animation / FLIP / selection / IME / focus / caret blink); static frames replay a cached Gio macro, and hit-testing is cached by cursor position and layout version. Idle drops to **0 repaints & 0 layouts/s** (verified) while the loop still runs at refresh rate. Perf HUD via `ui.ShowStats` / F12 (repaint & layout per-second, frame paint cost).
- Text-shaping cache: per-node memoization of `wrapForWidth` / `layoutRuns` (keyed by text/font/width), so repaints and re-layouts of unchanged text skip re-shaping — cache hit `~2ns`, `0 allocs` vs `~7.6µs`, `22 allocs` uncached (benchmarked).
- Refresh-rate adaptive: animation is wall-clock `dt`-based (speed constant across 60/120/144Hz), and long-press key-repeat is time-based. Gio schedules frames on demand rather than running a fixed-TPS loop.
- `ui.Post` for thread-safe updates from background goroutines.

**Layout**
- Yoga 3.2.1 feature baseline: flexbox, `display: contents`, border/content box sizing, absolute positioning, config-aware layout caches, `WidthPct`/`HeightPct`/`Fill` (window-adaptive), scroll (`ScrollView`), overflow clip.
- `VirtualList` for large lists — `UseScroll` feeds scroll offset/viewport back to the component so only the visible window (+overscan) renders; 100k rows stay smooth.

**Rendering**
- Gio-rendered rounded rects / borders / text / images — GPU vector rasterization with built-in AA; HiDPI (device-scale) rendering.
- SVG icons: `Icon`/`IconFill` render `pkg/svg` paths (stroke or fill), color inherited like text; a small built-in lucide set (`IconCheck`, `IconChevronDown`, …). Rounded-rect clipping (a rounded container clips its children to the corners via an offscreen mask). Linear gradients (`LinearGradient(from, to, angle)` background, follows the corner radius). `Img` object-fit (`Fit(FitContain/FitCover/FitFill)`).
- Paint goes through a `painter` backend interface (draw primitives + clip + layer), so the render walk is backend-agnostic — a Gio backend for the window, a recording backend for headless golden tests. The boundary rule lives at the top of `pkg/ui/backend.go`.

**Animation**
- `UseTween` + easings; `UseTransition` (enter/exit); FLIP layout animation; transforms (scale/rotate/translate, hit-test aware); per-node and group opacity.

**Input & text**
- Click (bubbling), hover, drag, wheel scroll; keyboard: Tab focus nav, Enter/Space activate, Esc stack; press state. Modal focus trapping (`Portal(TrapFocus(), …)`) — Tab stays inside the top modal; wired into shadcn Dialog/Sheet. Roving arrow-key navigation (`ArrowNav(NavVertical/NavHorizontal)`) inside menus/lists/tabs — wired into shadcn Tabs (←→) and DropdownMenu (↑↓).
- Controlled `Input` with caret, multi-line (`Multiline`), selection (Shift+arrows/drag/Ctrl+A, double-click word, triple-click all) and cut/copy/paste (pluggable clipboard); IME composition (`exp/textinput`) with underlined preedit at the caret. Grapheme-aware caret/backspace/delete and word-wise nav/delete (Ctrl+←→/Backspace) via `rivo/uniseg` — emoji, combining marks, ZWJ sequences move & delete as one unit.
- Text wrapping via Unicode line-breaking (UAX#14, `rivo/uniseg`) — hyphen breaks, CJK per-char, closing punctuation never at line start, non-breaking spaces; style inheritance, synthesized font weights/italic (one embedded CJK face — weight is effectively binary), rich-text spans (`RichText`), anchored overlays (`UseMeasure`).

**Components**
- Base kit (Checkbox/Switch/Radio/Slider/Progress/Spinner/Badge/Avatar/Divider/Card/Tabs).
- `pkg/shadcn`: ~41 shadcn/ui-style components on a theme + interaction foundation.

**Tooling**
- Live preview / hot reload (`pkg/hotreload`): edit a plain-Go `View() *ui.Node` file and the running window updates in-process — no rebuild, no restart (yaegi-interpreted; interpreted code uses non-generic `pkg/ui` + all `pkg/shadcn`, host owns state).
- Debug frame capture: `ui.Capture(path, afterFrames, exit)` or env `TENON_CAPTURE=out.png` saves the engine's own rendered frame to PNG (only the app's pixels — safe), for visually verifying rendering headlessly. `TENON_DEBUG_LAYOUT=1` enables Yoga layout/caching/measurement summaries through the event API.
- 40+ engine tests incl. headless **golden paint tests** (`Harness.Paint()` → `[]PaintOp` via the recording backend) + wrap/measure benchmarks; per-package READMEs + godoc; runnable examples.
- `ui.Mount` headless test harness (mount + drive click/hover/press/drag/type, query the render tree) — `pkg/shadcn` uses it for real behavior tests.

## Stabilization roadmap

The phases are dependency-ordered rather than date-based. **Done** requires code, regression tests and the listed acceptance checks on the deliverable itself; a configured but never executed CI job does not count as platform verification.

### Phase 0 — Core correctness (done)

**Scope:** component identity, Hooks, external state and deterministic cleanup.

- Preserve nil interface values in State, Reducer, Memo and object-component state.
- Follow Store source changes without leaking the old subscription.
- Reconcile duplicate keys without orphaning Fibers, effects or subscriptions.
- Bind dirty queues, effects, animation, focus and viewport updates to their owning tree.
- Keep the root module build, vet, tests and 35 interaction-bounds checks green.

**Acceptance:** `go build ./...`, `go vet ./...`, `go test ./...` and `go run _debug_snapshot/test_v3.go` all pass. The regression suite covers every fixed bug.

### Phase 1 — Runtime and backend lifecycle (done)

**Completed:** Harness and external backends share mount/update/close conformance tests; multiple `BackendRuntime` instances isolate input and device scale; out-of-order close skips disposed runtimes. Terminal `BackendRuntime.Close` now has an explicit inert-after-close regression test, while `Scene3DRuntime` retains components for `Activate` and remounting. `docs/RUNTIME_LIFECYCLE.md` documents the distinct ownership and recovery contracts. The backend-neutral 3D runtime supports `Mount → Update → Close → Mount`, reverse cleanup and failed-mount retry. CI jobs are configured for Linux, Windows, macOS and the separate G3N module.

**Configured:** `TENON_SMOKE_FRAMES=3` runs a native G3N plugin window for three frames and exits. Ubuntu CI has Xvfb lifecycle, UI interaction and projected-image smokes configured, and `REQUIRE_WINDOW=1` adds them to the native release check. `example/lifecycle` also performs an intentional failed mount, retry and two close/remount cycles during the native window run, asserting the final mount-attempt count. Remote execution remains required before claiming verification.

**Native verification:** G3N v0.2 build, full tests and race passed locally in WSL with CGO/GCC. G3N `TENON_CAPTURE` now reads final framebuffer pixels, and the same shared UI parity tree has been captured and visually compared with Gio, including text styles, icons, inputs and buttons. Real WSLg OpenGL runs passed for the versioned plugin and intentional failed-mount/two-cycle close-remount examples. These runs exposed and now regression-test optional Mount, pointer dispatch, light recreation and shared-camera ownership. A separate native/race test mounts two scene roots on one Game and verifies independent close, continued updates and exactly-once cleanup.

**Remaining:** execute the configured native CI jobs and verify other supported operating systems; do not treat configuration as a passed check.

**Acceptance:** lifecycle tests pass for headless UI, external UI and native 3D; no input, scale, subscription or native resource leaks remain after close/remount. The required local lifecycle matrix is complete; native G3N event-driven button, Tab-focus, select-all/delete and Chinese-input smoke also passes. Remote platform evidence remains a Phase 7 gate.

### Phase 2 — Yoga 3.2.1 compatibility (curated baseline done)

**Completed:** default `LayoutResults` now leaves `computedFlexBasis` undefined, and cached-layout event counters are no longer incremented twice. Upstream-derived fixtures cover flex growth/shrink/basis, percentage positioning/padding, absolute positioning, aspect ratios, min/max constraints, baseline callbacks, display none/contents, child-list ownership/reordering, dirtied propagation, measure-cache reuse, pixel-grid rounding, per-node scale factors and cache event statistics. This is a tested compatibility baseline, not a claim that every upstream test has been ported; future regressions should add their corresponding upstream fixture.

**Curated coverage:**

1. Flex grow/shrink, flex basis, min/max constraints and fractional factors.
2. Percentage dimensions/insets/padding, RTL, absolute/static containing blocks and aspect ratio.
3. Measure functions, baseline, display contents/none, tree mutation and cache invalidation.
4. Pixel-grid rounding at different scale factors and layout-event ordering.
5. Record intentional differences from upstream instead of changing expected values to match Tenon.

**Acceptance:** each curated upstream fixture has a stable Go regression test and a documented source; all Yoga and UI tests pass; intentional incompatibilities are listed explicitly.

### Phase 3 — Public API and plugin contracts (done)

**Completed:** the recommended object-component and function-component paths are documented with Hook ordering, unique-key, nil and UI-thread rules. `Plugin3DVersion` and `RegisterPluginVersion` enforce exact positive-version matching without changing the unversioned API. `Scene3DValidator` runs before native adaptation and cleans up earlier components on failure; `SupportsPlugin` provides side-effect-free plugin queries. `BackendCapabilityProvider` and `CurrentBackendCapabilities` expose SVG, rounded-clipping, transform and projection support without breaking existing backends. The G3N plugin example uses a validated versioned adapter. The backend-neutral `example/plugin` is also built and executed from an independent module by `make release-check`, covering invalid-payload rejection and two update/close cycles with an intervening `Activate`. Native mounting remains covered by Phase 1 rather than being assumed here.

**Scope:**

- Recommend one default component style and mark chain nodes, Hooks, object components and struct declarations as complementary rather than interchangeable defaults.
- Document Hook ordering, key rules, nil handling, UI-thread ownership and background `Post` behavior.
- Define plugin naming, payload validation, construction/adaptation errors, lifecycle errors and optional capability negotiation.
- Add plugin API versioning before the ecosystem depends on `Plugin3D(name, payload)`; preserve per-scene registries and native escape hatches.
- Specify backend capability levels so unsupported G3N effects are explicit rather than silently assumed to match Gio.

**Acceptance:** an external module can build and remount a versioned plugin without importing core internals; invalid payloads fail before native mounting; migration and compatibility rules are documented.

### Phase 4 — Rendering, text and platform completion (done)

**Completed:** `ui.Capture(path, frames, exit)` and `TENON_CAPTURE` now save the current Gio scene through the existing painter, keep static applications producing frames until the capture point, and optionally close the window. A real Windows `example/accordion` run produced a valid 1200×860 PNG. WSLg BiDi, real two-window, interactive-catalog, gradient/filter and tilted 3D-card captures also passed locally. The card example now uses a polished destination-card layout with stronger typography and a subtle initial tilt, while captures exercise actual projection rather than only a flat card. Public Harness regressions assert the third content region stays inside the 400px card and initial visual bounds differ from layout bounds. Linux CI has Xvfb accordion, gradient/filter, tilted 3D-card and projected-image screenshot smokes with artifact uploads configured; Windows and macOS CI also have BiDi, multiwindow, catalog, gradient/filter, tilted 3D-card and projected-image window smokes with per-platform artifacts configured. `docs/PLATFORM_SUPPORT.md` separates implemented capabilities, configured CI and executed verification. Gio now emits basic text, button, editor, checkbox, switch and radio semantics through its native semantic operations; backend-neutral snapshots include role, label, value, enabled/focused/selected state and hierarchy, with tests through Gio's input router and real shadcn components. Tabs, Toggle and ToggleGroup now expose native radio/checkbox selection semantics. Gio `SemanticAt` hit-testing is covered, including transformed controls and topmost Portal overlays. Enabled interactive nodes expose ClickGesture through stable per-node tags and frame filters, while disabled controls and static text do not. Interactive controls no longer expose duplicate child-text nodes that would steal semantic hits. Additional composite-widget roles and complete screen-reader workflows remain open. Plain text and single/multiline inputs now use `x/text/unicode/bidi` for visual runs, UTF-8 logical/visual mapping, caret hit-testing and segmented selection painting. RichText maps visual runs back to styled textRun segments while preserving baseline alignment. Arabic, Hebrew, number and bracket fixtures pass, and `example/bidi` has been executed and captured with real fonts. Run-interior logical/visual round trips are tested across mixed scripts; shared visual boundaries retain their dual-caret semantics. Full platform BiDi workflows remain open. `RadialGradient(from,to)`, `ImageBlur` and `ImageGrayscale` render through cached backend-neutral bitmaps with pixel/cache tests and real `example/gradient` captures. `RunWindows` now opens independent Gio windows with per-window input, IME, clipboard, wake registration and lifecycle; capture-driven close shuts down the whole group. A real two-window run passed locally, and desktop CI has multiwindow smoke configured. Frame submission remains serialized and other backends are not covered.

**Scope:**

- Close documented G3N gaps for SVG, rounded clipping, transforms and projection, or keep them as declared backend limitations.
- Add BiDi visual ordering together with caret, selection and hit-testing as one coordinated change.
- Add a screen-reader accessibility tree through platform APIs; retain existing focus trapping and roving navigation.
- Evaluate multi-window using the isolated runtime model; add native file/context menus only after ownership semantics are clear.
- Add real-window screenshot smoke tests on Linux, Windows and macOS.

**Optional follow-ups:** rounded clipping and complete 3D layer projection remain declared backend limitations rather than core blockers; pure projected images now share exact homography with Gio. Native file/context menus are optional platform enhancements; the shadcn ContextMenu remains the primary in-app menu. Richer composite semantics require capabilities beyond Gio v0.10.1's five existing classes, and real screen-reader workflows remain unverified. Other-backend multi-window support is also optional follow-up work.

**Acceptance:** platform support matrices match actual tests; text/input behavior is covered for CJK, emoji and mixed directions; no silent cross-window state sharing occurs. Local Gio/G3N rendering and native interaction acceptance is complete. Remote CI execution and broader screen-reader/platform workflows remain Phase 7 verification or optional enhancements.

### Phase 5 — Components and developer experience (done)

**Completed:** `example/catalog` provides a real, scrollable application for all major documented components, with connected state, callbacks and application-level interaction tests. Its polished layout now uses a centered maximum width, title hierarchy and intrinsic control widths instead of stretching buttons and badges into full-row blocks; regression tests cover both interaction and these dimensions. The README quick start now exists as `example/quickstart`, uses the recommended object-component path and has been run and captured as a valid 480×320 profile window. Its behavior test drives input, binding updates and the clear button. An independent temporary Go module successfully imported and built it through public APIs; the same external-consumer build is configured in CI. Harness now drains `ui.Post` like the real frame loop, and hot reload has an end-to-end test for file changes, retaining the last good view on interpreter errors and recovering after a fix. Switch and RadioGroup now have controlled interaction tests that also verify selection semantics; disabled Switch is covered. Slider, Toggle/ToggleGroup, Carousel, Resizable and Toast now test controlled values, visual transforms, panel bounds and expiry. `Query.VisualBounds()` exposes transform-aware rectangles for animation tests, and empty root trees now clear stale Portals correctly. Every documented shadcn component fixture is also mounted and closed through the real Fiber/layout lifecycle, rather than only checking constructor output. DropdownMenu now has a Tab/Arrow/Enter selection test, and Command has filtering/selection/close coverage. Menubar, Pagination, DatePicker, Calendar, NavigationMenu and Breadcrumb also test selection, month/page navigation, controlled updates and overlay close. Sheet, Collapsible and Tooltip now cover controlled close, natural-height expansion/collapse and hover visibility. Collapsible measures content independently of the collapsed zero-height container, fixing content that could never expand. `Focusable(false)` keeps pointer-only overlay wrappers out of keyboard focus while preserving click dismissal.

**Scope:** complete shadcn behavior coverage for keyboard, focus, disabled state and overlays; keep examples and bilingual API docs synchronized; test hot reload state ownership and interpreter limitations; provide a small application composed from the recommended API rather than feature-specific demos only.

**Acceptance:** every documented component has a runnable example and behavior test; quick-start code is built in CI; historical documents remain outside the current API path.

### Phase 6 — Performance and resource validation (done)

**Completed:** `docs/PERFORMANCE.md` records reproducible UI and Yoga benchmark environments, timings and allocations. `BenchmarkMountClose` covers the full Fiber/Yoga lifecycle with Store, Signal and Effect cleanup; its regression test verifies 50 cycles leave no subscriptions and exactly 50 cleanups. The local hover bitmap remains scoped rather than expanded without profiles. A per-renderNode, single-entry BiDi display cache removes repeated visual-run allocation and is documented with before/after benchmark evidence. A character-class fast path also avoids full paragraph construction for pure-LTR text while preserving the full algorithm for RTL content. `TestLifecycleMemoryTrend` adds a 200-cycle complex-tree heap check; the documented run shows no retained-heap growth after GC, without claiming it proves every leak is impossible. `BenchmarkFullFrameComplexTree` adds a 500-row combined frame/layout/paint-command baseline with measured time and allocation data. `Harness.PaintTo` allows discard/custom painters, and allocation profiling identified and removed test-recorder overhead from this baseline rather than triggering an unjustified production rewrite. `BenchmarkStartupComplexTree` and `BenchmarkIdleFrameComplexTree` now record first-mount/layout/paint/close cost and settled zero-allocation idle stepping for the same 500-row workload, with native-window and GPU exclusions documented. `BenchmarkVirtualList100kStartup` and `BenchmarkVirtualList100kScroll` now substantiate the 100k-row virtualization claim with a checked 17-row mounted window and measured startup/scroll allocations, while explicitly excluding OS input and GPU submission. Image-filter benchmarks motivated direct common-format pixel conversion and blur buffer reuse: 512×512 blur/grayscale allocations fell from roughly 262k per operation to 4/2, with standard-color equivalence tests. JPEG-style YCbCr grayscale subsequently fell from 262,146 to 2 allocations; Gray was measured and kept on its already allocation-free generic path.

**Scope:**

- Publish reproducible benchmark environments for startup, idle frames, reconciliation, text shaping, layout and complex virtualized lists.
- Track allocations and memory after repeated mount/close cycles, not only microbenchmark latency.
- Profile real component trees before changing keyed maps, slices or event registries.
- Keep the current local child bitmap for hover state; expand bitmap/index structures only with measured time and memory benefits.

**Acceptance:** baseline and candidate results include hardware, OS, Go version, workload and allocation counts; optimizations do not regress correctness or idle behavior. The complete documented UI/Yoga benchmark set, 50-cycle subscription cleanup and 200-cycle memory-trend test passed locally. Native backend profiling and long-running soak tests remain Phase 7 release evidence rather than unverified claims here.

### Phase 7 — Release gates (in progress)

**Configured:** CI builds the quick-start application and runs the backend-neutral versioned-plugin lifecycle from independent modules. Remote run 45042687522 passed Ubuntu build-test/race at `f30269f`. Current run 45792815811 has passing Gio desktop jobs, but `build-test-g3n` fails at `go build`; its compiler log requires authenticated inspection, so Phase 7 remains open. `make release-check` now covers formatting, root checks, interaction bounds, independent quick-start/plugin consumers, optional race/G3N and real-window capture. It ran successfully on the Windows host with valid BiDi, multiwindow, interactive-catalog, gradient/filter and tilted 3D-card captures, while race and G3N were explicitly skipped because host CGO is disabled. WSL subsequently passed G3N build/tests/race and native plugin/lifecycle windows. Root race passed after excluding six named headless GPU pixel tests; the unfiltered run fails those tests on empty frames, so full root race remains unverified even though real WSLg rendering works. UI, shadcn and hotreload also passed five repeated test runs. `CHANGELOG.md` and `docs/RELEASING.md` define version synchronization and stable-release evidence; skips never count as release proof. `docs/public-api.txt` and `make api-check` now detect public signature drift, with `make api-update` reserved for reviewed intentional changes.

**Before a stable 1.0:** freeze the supported API, synchronize root and backend module versions, add changelog/migration notes, pass race-enabled root tests and native G3N tests, verify packages from an external consumer module, and complete the relevant real-window smoke tests. Pre-1.0 tags may be used earlier with explicit instability notes.

**Acceptance:** a clean checkout on each supported OS can build, test and run the documented examples without repository-local replace directives or undocumented native setup. The published pseudo-version G3N consumer exposed root/backend API version skew (`ui.MouseButton`, `ui.InputKey`, `ui.Layer`, `ui.FontFace` and related types missing). The workspace now requires the matching root pseudo-version and keeps local replacement only in ignored `go.work`; a new push and external consumer run are required before claiming this gate.

### Immediate execution queue

1. Execute the configured G3N CI jobs and verify native lifecycle behavior on the other supported operating systems; local Linux/WSL build, race and OpenGL lifecycle checks now pass.
2. Execute the configured Linux, Windows and macOS workflows, inspect their real-window artifacts, and keep executed results distinct from configuration. Resolve full root race coverage in an environment where headless GPU pixel tests are reliable.
3. Keep optional backend limitations explicit: rounded clipping and exact 3D projection remain unimplemented, while shared SVG, text styles, 2D transforms and native pointer/keyboard interactions are complete. Do not expand native menus ahead of core release gates.
4. For a stable release, require race/G3N/native-window checks and external consumers built from synchronized published module versions rather than local replace directives.
5. Keep Phase 2 open for upstream regressions and new fixtures even though its curated baseline is complete.

## In progress / next (priority order)

1. **BiDi text** — ~~plain text, inputs and style-aware RichText visual ordering, caret/selection and hit-testing~~ **done** via `x/text/unicode/bidi`, with logical-byte/visual-x and styled-run mapping tests. Still: broader language/font matrices and end-to-end platform workflows.
2. **Accessibility** — ~~focus trapping in modals~~ **done** (`TrapFocus()`); ~~arrow-key navigation inside menus/lists~~ **done** (`ArrowNav`, roving focus); ~~basic Gio semantic tree~~ **done** for text, buttons and editors. Still: additional composite-widget roles, richer properties and end-to-end screen-reader workflows across platforms.
3. **Performance at scale** — ~~list virtualization~~ **done** (`VirtualList` + `UseScroll` renders only the visible window; 100k rows stay smooth); ~~sub-tree-scoped `resolveInherited`~~ **done** (unchanged siblings reuse their resolved inheritance context); structure-gated `relink` and Portal collection skip full Fiber walks for content-only updates.
4. **Rendering extras** — ~~SVG icons~~, ~~rounded-rect clipping~~, ~~linear gradients~~, ~~`Img` object-fit~~ **all done**. ~~Radial gradients~~, ~~cached image box blur (`ImageBlur`)~~ and ~~luminance grayscale (`ImageGrayscale`)~~ are also done. Image filters share source/parameter cache keys, run after decoding on workers and reject stale asynchronous results. Remaining polish: additional compositing filters.
5. **Native integration** — ~~OS clipboard binding~~ and ~~Gio multi-window~~ **done** (`RunWindows` with independent per-window state and serialized frame submission). Optional: native file/context menus and multi-window support in other backends.

**Recently done:** migrated the renderer from Ebiten to **Gio** (Ebiten is gone from `go.mod`). Pseudo-3D (`Perspective`/`RotateX`/`RotateY`/`TranslateZ`) plus `Scene3D` — a shared camera so a table of elements agrees on one vanishing point — and `PlaneImage` for an exactly-projected floor texture. Hit-testing follows the 3D projection. Also: `SrcImage` (in-memory image source), `ImageBlur` (cached two-pass box blur), `ImageGrayscale` (backend-neutral luminance filtering), `OnSubmit` (Enter on a single-line `Input`), `ZIndex`, flex `Wrap`, an LRU budget on the image cache, and window title / min-max size / fullscreen. The public `Backend`/`BackendRuntime` contract now supports external backends; the independent `backend/g3n` module overlays Tenon UI on native G3N game scenes without adapter-specific asset downloads. Its `GameScene` provides declarative Camera/Mesh/Material/Light components, shared resources, value-bound transforms, and Spin animation. G3N advanced UI effects are not yet at Gio parity (see the backend guide).

## Non-goals (for now)

- VDOM-style full diffing (we use Flutter/React fiber identity + keys).
- Vue-style Proxy reactivity (Go has no Proxy; hooks cover it).

Contributions toward any of the "next" items are very welcome — see [CONTRIBUTING.md](CONTRIBUTING.md).
