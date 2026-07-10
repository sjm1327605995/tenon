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

## 🛠️ Development Setup

### Prerequisites

- Go 1.21 or later
- Git

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
go run cmd/demo/main.go
```

## 📋 Code Style

- Follow standard Go conventions (`gofmt`, `go vet`)
- Add comments for exported types and functions
- Keep functions focused and reasonably sized
- Prefer explicit error handling over panics in library code

## 🧪 Testing

- Add tests for new features and bug fixes
- Ensure all tests pass before submitting PRs
- Tests are located alongside source files (e.g., `pkg/ui/engine_test.go`)

```bash
# Run all tests
go test ./...

# Run specific package tests with verbose output
go test ./pkg/ui/ -v
go test ./yoga/ -v
```

## 🏗️ Project Structure

When adding new features, please respect the existing architecture (see [ARCHITECTURE.md](ARCHITECTURE.md)):

- `pkg/ui/` — Framework core: elements, hooks, reconciler, layout, rendering, input
- `pkg/shadcn/` — shadcn/ui-style component library built on `pkg/ui`
- `pkg/font/` — Font management
- `yoga/` — Yoga layout engine (Go port, minimize changes)
- `example/` — Runnable example programs (`hooks-*`, `shadcn-*`)

## 📝 Documentation

- Update `README.md` and `README.zh-CN.md` if you add or modify public APIs
- Keep both English and Chinese documentation in sync
- Add code examples for new components or features

## 🎯 Component Development Guide

When adding a new component to `pkg/components/`:

1. Embed `core.BaseHost`
2. Call `Init(self)` in the constructor
3. Implement `Draw(screen *ebiten.Image)` for rendering
4. Implement `HandleEvent(e *core.Event) bool` for interaction
5. Provide fluent chained API methods
6. Set `focusable = true` if the component should receive keyboard focus

Example:

```go
type MyComponent struct {
    core.BaseHost
}

func NewMyComponent() *MyComponent {
    m := &MyComponent{}
    m.Init(m)
    m.SetFocusable(true) // if interactive
    return m
}

func (m *MyComponent) Draw(screen *ebiten.Image) {
    // rendering logic
}

func (m *MyComponent) HandleEvent(e *core.Event) bool {
    // event handling logic
    return false
}
```

## 📄 License

By contributing to Tenon, you agree that your contributions will be licensed under the MIT License.

---

If you have any questions, feel free to open an issue or reach out. Happy coding! 🎉
