# Backend and platform support

This matrix distinguishes implementation, configured CI and executed verification. A CI job is not considered passed until it has run on the target platform.

## Verification status

| Backend | Platform | Build and tests | Real-window verification |
| --- | --- | --- | --- |
| Gio | Windows | Passing locally; Windows CI configured | Accordion, BiDi, real two-window, interactive-catalog and gradient/filter captures passed locally; BiDi, multiwindow, catalog, gradient/filter and tilted 3D-card smoke jobs configured in CI |
| Gio | macOS | macOS CI configured | BiDi, multiwindow, catalog and gradient/filter smoke jobs configured; not yet executed remotely |
| Gio | Linux | Ubuntu build-test and race passed remotely at `f30269f`; WSL root race passed locally after excluding six named headless GPU pixel tests, so full race is not claimed locally | BiDi, real two-window, interactive-catalog, gradient/filter and tilted 3D-card captures passed locally in WSLg; current screenshot jobs are configured but not yet executed remotely |
| G3N | Windows | Requires CGO, MinGW-w64 and G3N's native runtime | Not verified in the current environment (`CGO_ENABLED=0`) |
| G3N | macOS | Requires CGO, OpenAL and Vorbis | Not verified; no dedicated CI job yet |
| G3N | Linux | Build, full tests and race passed locally in WSL with CGO/GCC; Ubuntu CI configured with G3N's documented native dependencies | Native plugin, failed-mount/two-cycle close-remount and SVG UI parity windows passed locally in WSLg; Xvfb smokes configured in CI but not yet executed remotely |

## Remote CI evidence

[CI run 45042687522](https://github.com/sjm1327605995/tenon/actions/runs/45042687522) completed successfully for commit `f30269f`: Ubuntu `build-test`, including `go test -race ./...`, passed. That run predates the current desktop and G3N jobs. Those jobs and current local commits require a push and new workflow run; no current remote success is inferred from the older run.

## Rendering capabilities

`ui.CurrentBackendCapabilities()` is the programmatic source for optional advanced effects. Undeclared third-party capabilities are false.

| Capability | Gio | G3N | Notes |
| --- | --- | --- | --- |
| Text, rectangles and images | Yes | Yes | G3N shares Tenon's embedded font |
| `SVGPaths` | Yes | Yes | Both backends share SVG parsing; G3N uses cached antialiased CPU rasterization |
| `RoundedClipping` | Yes | No | G3N currently degrades or ignores advanced clipping |
| `LayerTransforms` | Yes | Yes | G3N supports 2D scale, shared-center rotation, translation and group opacity with geometry tests |
| `ExactProjection` | Yes | Partial | Pure projected images share homography and cached warps with Gio; complete 3D layer projection is not equivalent |
| `Accessibility` | Basic | No | Gio v0.10.1 exposes five classes—button, editor, checkbox, switch and radio—all mapped by Tenon, while static text uses labels rather than a class; shadcn Tabs, Toggle and ToggleGroup use interactive roles automatically, while richer composite classes and full screen-reader workflows remain incomplete |
| Backend-neutral 3D declarations | N/A | Yes | Mesh, camera, light, material and transform props |
| Versioned 3D plugins | N/A | Public API tested | Native plugin execution requires the G3N environment |

## Runtime and platform features

| Feature | Status |
| --- | --- |
| UI mount/update/close and repeated close | Tested through Harness and `BackendRuntime` conformance tests; terminal close and inert-after-close behavior are covered by the [lifecycle contract](RUNTIME_LIFECYCLE.md) |
| Multiple external runtime input and scale isolation | Tested headlessly; this is not OS multi-window support |
| 3D close/remount and failed-mount retry | Backend-neutral tests and native Linux G3N tests/race passed, including independent scene roots and shared-camera ownership; other native platforms remain unverified |
| Gio frame capture | Implemented through `ui.Capture` and `TENON_CAPTURE`; locally verified |
| G3N frame capture | `TENON_CAPTURE` uses final-framebuffer `ReadPixels`; locally verified against the same shared UI parity tree as Gio, with CI screenshot artifacts configured |
| Pointer feedback | Gio and G3N share `BackendRuntime.PointerCursor`: text inputs use IBeam, clickable targets use Hand, and disabled/empty targets use Arrow |
| BiDi text | Plain text, inputs and style-aware RichText support mixed-direction visual runs, caret mapping and segmented selection; extended Arabic/Hebrew platform workflows remain pending |
| Screen-reader accessibility tree | Basic Gio text, button, editor, checkbox, switch and radio semantics tested through Gio's input router, including Tabs, Toggle and ToggleGroup; semantic coordinate hit-testing, transforms, topmost Portal overlays and ClickGesture activation are covered through Gio's input router, while full composite-widget and platform workflows remain pending |
| Multiple OS windows | `RunWindows` supports multiple Gio windows with independent input, IME, clipboard and lifecycle; a real two-window run passed locally. Frame submission is serialized; other backends are not covered |
| Native file/context menus | Optional enhancement, not a core priority; shadcn ContextMenu is the primary in-app implementation. Native view events exist, but no OS-menu implementation is currently delivered |

## Reproduce verification

Run the root checks:

```sh
go build ./...
go vet ./...
go test ./...
go run _debug_snapshot/test_v3.go
```

On the verified WSL configuration, the unfiltered race run produces empty frames in `TestFauxBoldAndItalicChangePixels`, `TestProjectedOutlineIsTrapezoidNotParallelogram`, `TestProjectedContentFillsClipQuad`, `TestGio3DCardIsProjected`, `TestGio3DWarpsTextToo` and `TestScene3DSharesVanishingPoint`. The remaining root race suite passes with only those six names excluded. Forcing OpenGL, software rendering or Lavapipe did not fix the headless path; the exclusion is diagnostic coverage, not release evidence. A real WSLg window capture succeeds independently.

Capture a real Gio application frame and exit automatically:

```sh
TENON_CAPTURE=out.png TENON_CAPTURE_FRAMES=45 go run ./example/accordion
```

The capture contains only Tenon's rendered pixels, not the OS desktop. On headless Linux, prefix the command with `xvfb-run -a` after installing Xvfb and its native Gio dependencies.

Validate G3N in an environment with CGO and its native libraries:

```sh
cd backend/g3n
go build ./...
go test ./...
```

A real G3N window smoke is still required before claiming native platform support.
