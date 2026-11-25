# Changelog

All notable changes to Tenon are documented here. The project is pre-1.0; public APIs may still change.

## Unreleased

### Added

- Versioned 3D plugins with exact adapter matching, payload validation and capability queries.
- Optional backend capability reporting for advanced rendering and accessibility support.
- Repeatable 3D mount/update/close/remount lifecycle and failed-mount cleanup.
- A runtime lifecycle guide distinguishing terminal UI backend closure from reactivatable 3D scenes, with inert-after-close regression coverage.
- Gio frame capture through `ui.Capture` and `TENON_CAPTURE`, including automatic exit and CI artifacts.
- Basic accessibility snapshots and Gio semantics for text, buttons, editors, checkboxes, switches and radio buttons, including shadcn Tabs, Toggle and ToggleGroup, with coordinate hit-testing and ClickGesture activation.
- BiDi visual runs for plain text, inputs and styled RichText, including caret, selection and hit-test mapping.
- Runnable quick-start, BiDi and interactive component-catalog examples with real-window verification.
- Gio `RunWindows` with independent per-window input, IME, clipboard and lifecycle, plus a real two-window example.
- `RadialGradient` with cached backend-neutral bitmaps, rounded clipping and a real gradient example.
- `ImageBlur` with two-pass cached box blur for file, network and in-memory images, plus an original/blurred comparison in the gradient example.
- `ImageGrayscale` with clamped Rec.709 luminance blending, alpha preservation, combined filter cache keys and an original/blur/grayscale example.
- Shared backend-neutral `ParseSVGPath` and G3N antialiased SVG fill/stroke rendering with appearance-scoped bitmap caching and a native UI parity example, wired into native release checks and CI.
- Reproducible performance and platform-support documents and a consolidated release check script.
- A checked-in public API baseline and `make api-check` / `make api-update` for deliberate signature changes.
- `Harness.PaintTo` for custom/discard painters and allocation profiles that exclude test-recorder overhead.

### Fixed

- Hook nil interface handling, Store source switching, duplicate-key cleanup and runtime-owned scheduling.
- Multiple runtime input, scale, focus, animation and viewport isolation, including out-of-order backend close.
- Yoga default layout initialization, flex-basis precedence and duplicate cached-layout event counting.
- Harness Post draining now preserves unmount lifetime semantics and supports hot-reload recovery.
- Disabled composite controls remain present in the accessibility tree.
- Empty root trees clear stale Portals, and `Query.VisualBounds()` supports transform-aware animation assertions.
- Pointer-only overlay wrappers no longer steal keyboard focus from menu items; `Focusable(false)` preserves their click behavior.
- Collapsible content now measures its natural height independently of the collapsed container, allowing it to expand correctly.
- Shared scroll updates now clamp to the available content range, and the interactive catalog connects real control state and callbacks.
- Interactive controls no longer expose duplicate child-text semantics that steal Gio `SemanticAt` hits from the actual control role.
- RTL hit-testing now maps coordinates to logical offsets without incorrectly subtracting the visual run origin; mixed-script interior mappings are regression-tested.
- Image cache and projected-texture keys now include the blur radius, and late asynchronous results cannot overwrite a newer radius on the same node.
- Removing an image source now clears the existing bitmap and filter state; re-adding it uses the new attributes instead of retaining stale visuals. A delayed HTTP test also confirms that pending results cannot resurrect a cleared image.
- Updated G3N v0.2 window/scene API usage and root-module version requirements, moved local replacement out of `go.mod` into ignored `go.work`, made native Mount optional as documented, fixed pointer dispatch, and repaired native plugin payload adaptation.
- G3N light objects are now recreated on remount, while closing a declarative scene no longer detaches the shared game camera.
- G3N now honors synthetic bold and italic text styles with Gio's thresholds, including correct fixed-point baseline conversion and pixel regressions.
- Declared G3N 2D layer transforms after verifying scale, shared-center rotation, translation and group opacity with geometry tests.
- Unified Gio/G3N pointer feedback through `BackendRuntime.PointerCursor`, with native Hand/IBeam/Arrow roles and disabled-target behavior.
- G3N now draws continuous lines for text carets and underlines instead of placing a square every four pixels.
- G3N transformed panels now own their `GraphicMaterial` render callbacks, so clip and transform setup is not bypassed. Rounded stencil clipping remains unimplemented after an unsuccessful experiment was fully reverted.
- Added G3N `TENON_CAPTURE` framebuffer screenshots and a shared Gio/G3N UI parity tree, including native-event button, Tab-focus, select-all/delete and Chinese-input assertions in CI. Real captures exposed and fixed synthetic-text background rectangles by tinting glyph pixels and keeping panel backgrounds transparent.
- Reworked the 3D card into a polished destination card with a light canvas, rounded surface, featured badge, stronger typography and a subtler initial tilt. All content regions remain visible, and Harness tests guard content bounds and initial projection.
- Improved the component catalog with a centered maximum width, page title hierarchy and intrinsic button/badge sizing, with interaction and control-width regressions.

### Performance

- Added lifecycle allocation benchmarks, 50-cycle subscription cleanup verification and a 200-cycle complex-tree heap trend check.
- Added reproducible 500-row startup and static idle-frame benchmarks with explicit native-window and GPU measurement boundaries.
- Added checked 100k-row virtual-list startup and bidirectional-scroll benchmarks, verifying that only the 17-row visible/overscan window mounts.
- Kept bitmap usage scoped to measured hover-state workloads rather than replacing general keyed maps.
- Added a per-node BiDi display cache and a pure-LTR character-class fast path with before/after allocation and latency benchmarks.
- Removed per-pixel image-filter allocations for common NRGBA/RGBA and JPEG-style YCbCr images and reused the blur input as its second-pass output, reducing 512×512 blur/grayscale to 4/2 allocations per operation.

### Verification

- Added curated Yoga 3.2.1 upstream fixtures.
- Added external consumer, hot-reload, lifecycle, accessibility, BiDi and real Gio window tests.
- Added an external versioned-plugin consumer that executes payload validation and close/reactivate lifecycle checks in the release script; equivalent CI coverage is configured. The release script also rejects development `replace` directives in `backend/g3n/go.mod` with `GOWORK=off`.
- Regenerated hot-reload symbols and added interpreted execution coverage for SVG parsing, icons, image blur and grayscale.
- Added shared `WarpImageToQuad` and G3N cached projected-image warps; Gio/G3N four-color projection captures have mean RGBA differences below one channel value.
- Completed the Phase 6 performance baseline with a full documented UI/Yoga benchmark run, lifecycle cleanup and memory-trend verification.
- Completed the Phase 4 local rendering/text/platform baseline, including shared Gio/G3N UI parity and native button, Tab-focus, selection and Chinese-input smoke.
- Configured Linux, Windows, macOS and separate G3N CI jobs, including gradient/image-filter, tilted 3D-card and projected-image window captures and artifacts; see `docs/PLATFORM_SUPPORT.md` for executed versus configured checks. Remote Ubuntu build-test/race passed at `f30269f`, while current jobs still require a new run.
- Passed native Linux G3N build/tests/race and real WSLg plugin/lifecycle windows locally, including independent scene-root isolation, plus real Linux Gio BiDi, two-window, interactive-catalog, filter and tilted 3D-card captures. Root race passed after excluding six named WSL headless-GPU pixel tests; the unfiltered run still fails those tests on empty frames, so full root race is not claimed.
