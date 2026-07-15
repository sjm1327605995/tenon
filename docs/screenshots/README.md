# Screenshots

Drop screenshots/GIFs here and reference them from the root README.

`accordion.png` (the README hero) is the `example/accordion` docs-page demo. To
regenerate it headlessly, capture the engine's own framebuffer to a PNG:

```bash
# Windows (Git Bash)
TENON_CAPTURE=docs/screenshots/accordion.png TENON_CAPTURE_FRAMES=45 go run ./example/accordion
```

`TENON_CAPTURE_FRAMES` waits a few frames so the scroll-reveal fade-in settles
before the frame is saved. This captures only the app's own rendered frame — not
the OS desktop.
