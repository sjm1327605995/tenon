# Runtime lifecycle

Tenon has two deliberately different runtime contracts. `BackendRuntime` drives a retained UI tree; `Scene3DRuntime` coordinates native 3D components. Their similarly named `Close` methods are not interchangeable.

## UI: `BackendRuntime`

A backend creates an instance with `NewBackendRuntime`, then calls `Resize`, `Frame` and `Paint` from its window loop. Each instance owns its component tree, input source, device scale, focus, subscriptions and effects.

`Close` is **terminal and idempotent**:

- The first call unmounts the tree, runs effect cleanup and releases Tenon-owned subscriptions.
- Later calls do nothing. `Frame`, `Paint`, `Resize` and other operations remain inert or return zero values after close.
- There is no `Activate` or restart operation. Create a new `BackendRuntime` to mount another UI tree.
- Closing an inactive runtime does not dispose the active runtime. Closing the active one restores the most recent live runtime, or the input and scale that preceded the outermost instance.
- `InputSource` is supplied by the backend and is not closed by Tenon; its owner remains responsible for any native resources.

Instances may be driven in interleaved order, but an individual instance is not safe for concurrent calls. Schedule background state changes with `ui.Post` rather than mutating Hooks or Stores from a worker.

`Harness.Close` follows the same terminal ownership model for headless UI tests.

## 3D: `Scene3DRuntime`

A scene runtime owns an ordered component list. `Update` visits components in declaration order; `Close` releases them once in reverse order. Repeated `Close` calls without an intervening activation do nothing.

Unlike `BackendRuntime`, a 3D runtime retains its component list after `Close`. The backend can remount that same scene with this sequence:

1. Call `Activate` to allow updates again.
2. Mount the retained components into the native scene through the backend.
3. Resume `Update` calls.

`Activate` only changes the runtime state; it does not mount native resources itself. Components must therefore support their backend's `Mount → Update → Close → Mount` cycle. `Components` returns a copy of the list so callers cannot mutate runtime ownership through the returned slice.

## Construction failure and independent roots

`Scene3DPluginRegistry.Build` validates and adapts declarations before returning a runtime. If any declaration fails, it closes already-constructed components in reverse order and returns an error. Retry by calling `Build` again after correcting the declaration or registration; do not call `Activate` on a nil runtime.

Registries, scene runtimes and native scene roots are independent. Closing one scene must not remove another scene's native objects. Native G3N resource cleanup and window behavior still require execution in a CGO-enabled environment with G3N's native libraries; backend-neutral tests alone are not native verification.
