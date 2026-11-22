# Contributing to Tenon

Thank you for your interest in contributing to Tenon! This document provides guidelines and instructions for participating in the project.

## 🌟 How to Contribute

### Reporting Issues

If you find a bug or have a feature request, please submit an issue on GitHub. Before submitting, please:

1. Search existing issues to avoid duplicates
2. Provide a clear description and reproduction steps (for bugs)
3. Include your Go version and operating system information

### Submitting Pull Requests

1. **Fork** the repository and create your branch from `main`
2. **Write** clear, concise commit messages
3. **Test** your changes (`go test ./...`)
4. **Ensure** your code follows the existing style
5. **Update** documentation if your changes affect the API
6. **Submit** a Pull Request with a clear description

## Development Setup

### Prerequisites

- Go 1.24 or later
- Git
- A desktop session for running windowed examples

On Ubuntu, install Gio's native build dependencies before building or testing:

```bash
sudo apt-get update
sudo apt-get install -y \
  pkg-config libwayland-dev libx11-dev libx11-xcb-dev \
  libxkbcommon-dev libxkbcommon-x11-dev \
  libgles2-mesa-dev libegl1-mesa-dev libffi-dev \
  libxcursor-dev libxfixes-dev libvulkan-dev
```

This matches the dependency list in [CI](.github/workflows/ci.yml). The Linux test wrappers additionally require `xvfb`.

### Local Development

```bash
# Clone your fork
git clone https://github.com/your-username/tenon.git
cd tenon

# Install dependencies
go mod download

# Run tests
go test ./...

# Run the demo
go run ./example/accordion
```

## 📋 Code Style

- Follow standard Go conventions (`gofmt`, `go vet`)
- Add comments for exported types and functions
- Keep functions focused and reasonably sized
- Prefer explicit error handling over panics in library code

## Testing

- Add tests for new features and bug fixes
- Ensure all tests pass before submitting PRs
- Tests are located alongside source files (e.g., `pkg/ui/engine_test.go`)

```bash
# Run all tests
go build ./...
go vet ./...
go test ./...

# Check interaction bounds (no window required)
go run _debug_snapshot/test_v3.go

# Detect data races (requires cgo and a C compiler)
go test -race ./...

# Run specific package tests with verbose output
go test ./pkg/ui/ -v
go test ./yoga/ -v
```

`make lint` runs `go vet ./...`. On Linux without a display, install `xvfb` and run `./scripts/test.sh ./...`; `make test` uses the environment-detecting wrapper and preserves test failures in its exit status. On Windows, the Go commands above can be run directly in PowerShell.

The root `./...` pattern does not cross module boundaries. Validate the optional game backend separately with `cd backend/g3n && go build ./... && go test ./...`; the repository's local `go.work` resolves the root module without adapter-specific downloads.

## 🏗️ Project Structure

When adding new features, please respect the existing architecture (see [ARCHITECTURE.md](ARCHITECTURE.md)):

- `pkg/ui/` — Framework core: elements, hooks, reconciler, layout, rendering, input
- `pkg/shadcn/` — shadcn/ui-style component library built on `pkg/ui`
- `pkg/decl/`, `pkg/value/`, `pkg/router/` — declarations, data binding, and navigation
- `pkg/hotreload/` — Yaegi-based live preview
- `backend/g3n/` — Optional G3N game backend in a separate Go module
- `yoga/` — Yoga layout engine (Go port, minimize changes)
- `example/` — Runnable examples: `accordion`, `card3d`, `component`, `router`, `simple`, `struct`

## 📝 Documentation

- Update `README.md` and `README.zh-CN.md` if you add or modify public APIs
- Keep both English and Chinese documentation in sync
- Add code examples for new components or features

## 🎯 Component Development Guide

Components describe a UI tree; the engine (`pkg/ui`) owns layout and rendering. For application
components with persistent fields and lifecycle context, embed `ui.Component` and implement
`Render() ui.NodeSource`; see the [quick start](README.md#quick-start) and
[component example](example/component). Use `value.Bind[T]` for data that controls text and inputs.

Function components remain useful for reusable widgets. A public constructor wraps a function
with `ui.Use`, which gives it a fiber identity so hooks have somewhere to live:

```go
// SwitchProps 集中描述组件配置，字段可选；引擎据此比较输入并决定是否重渲染。
// 返回的 *ui.Node 可继续链式设置通用布局、样式和事件；复杂配置仍保留在 props 中。
type SwitchProps struct {
    Checked  bool
    Disabled bool
    OnChange func(bool)
}

// 公开构造函数：ui.Use 把函数与一个 fiber 绑定，hooks 的状态就存在那上面。
func Switch(p SwitchProps) *ui.Node { return ui.Use(switchC, p) }

// 内部实现：读 props、调 hooks、返回节点树。
func switchC(p SwitchProps) *ui.Node {
    th := ui.UseTheme()                       // 主题
    x := ui.UseTween(bool2f(p.Checked), 140, ui.EaseOut) // 动画

    attrs := []*ui.Node{ui.Style(
        ui.Width(32), ui.Height(18), ui.Radius(9999),
        ui.Bg(ui.Mix(th.Input, th.Primary, x)),
    )}
    if !p.Disabled {
        attrs = append(attrs, ui.OnClick(func() {
            if p.OnChange != nil {
                p.OnChange(!p.Checked)
            }
        }))
    }
    return ui.Div(attrs...)
}
```

Rules of Hooks apply: call them unconditionally, in the same order every render (no hooks inside
`if`/loops), or state will bind to the wrong slot.

Where things go:

- `pkg/ui` — the engine: elements, hooks, layout, style, input, painting. Backend-neutral except
  the `gio_*.go` files; see the boundary rule at the top of `pkg/ui/backend.go` before touching them.
- `pkg/shadcn` — the component library (shadcn/ui port). New general-purpose components go here.

Test behavior, not construction. `ui.Mount` gives a headless harness that drives the real
reconcile → layout → hit-test → event path, so assert what a user would observe:

```go
h := ui.Mount(Switch(SwitchProps{Checked: false, OnChange: set}), 200, 100)
defer h.Close()
h.ClickAt(16, 9)
// 断言 OnChange 真的被调用了，而不是断言构造函数返回了非 nil
```

## 📄 License

By contributing to Tenon, you agree that your contributions will be licensed under the MIT License.

---

If you have any questions, feel free to open an issue or reach out. Happy coding! 🎉
