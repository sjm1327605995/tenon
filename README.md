# Tenon - React-like UI Framework for Go

Tenon is a React-like UI framework for Go, combining React's component-based approach with Gio's high-performance rendering, using Yoga layout engine for flexible styling, providing modern UI development experience for Go developers.

[📖 中文版本](README.zh-CN.md) | [🏠 Homepage](https://github.com/sjm1327605995/tenon)

## 📋 Core Features

- **React-like Component System**: Supports functional components and Hooks
- **Declarative UI**: Build views using chained API
- **Yoga Layout Engine**: Supports Flexbox and Grid layouts
- **State Management**: Built-in `useState` Hook for component state management
- **Router System**: Multi-page application support
- **Gio Rendering**: High-performance rendering based on Gio library
- **Event Handling**: Support for click and other user interactions

## 🏗️ Architecture Design

```mermaid
flowchart TD
    subgraph Application Layer
        A[Example Applications example]
        A1[Data Binding Demo data_binding_demo.go]
        A2[Router Demo router_demo.go]
    end

    subgraph Framework Core Layer
        B[Core UI Library core/ui]
        B1[Component System]
        B2[State Management]
        B3[Hooks System]
        B4[Router System]
        B5[Element Tree Management]
        B6[View Builder]
    end

    subgraph Rendering Layer
        C[Gio Integration core/gio]
        C1[Window Management]
        C2[Event Handling]
        C3[Render Bridge]
        D[Render Backend core/ui/render]
        D1[RenderObject Tree]
        D2[Style Rendering]
        D3[Click Event]
        D4[Text Rendering]
    end

    subgraph Dependency Layer
        E[Yoga Layout Engine yoga]
        E1[Style System]
        E2[Layout Calculation]
        F[Gio Library gio]
        F1[Layout Engine]
        F2[Painting System]
        F3[Event System]
    end

    %% Application Layer Connections
    A --> A1
    A --> A2
    A1 --> B
    A2 --> B

    %% Core Layer Internal Connections
    B --> B1
    B --> B2
    B --> B3
    B --> B4
    B --> B5
    B --> B6
    B2 --> B3
    B3 --> B1
    B4 --> B1

    %% Core Layer to Rendering Layer
    B5 --> C
    B5 --> D
    B6 --> D

    %% Rendering Layer Internal Connections
    C --> C1
    C --> C2
    C --> C3
    D --> D1
    D --> D2
    D --> D3
    D --> D4
    D1 --> D2
    D1 --> D3
    D1 --> D4

    %% Rendering Layer to Dependency Layer
    C --> F
    D --> F
    B5 --> E
    B6 --> E

    %% Dependency Layer Internal Connections
    E --> E1
    E --> E2
    F --> F1
    F --> F2
    F --> F3
```

## 🚀 Quick Start

### Installation

```bash
go get github.com/sjm1327605995/tenon
```

### Running Examples

```bash
# Run data binding example
go run example/data_binding_demo.go

# Run router example
go run example/router_demo.go
```

## 📚 Core Features

### 1. Component System

```go
// Define component props
type CounterProps struct {
    InitialCount int
}

// Create functional component
func Counter(props CounterProps) ui.UI {
    // Use useState Hook
    count, setCount := ui.UseState(props.InitialCount)
    
    return ui.View(
        ui.Text().Content(fmt.Sprintf("Count: %d", count)),
        ui.View(
            ui.Text().Content("+"),
        ).Background(color.NRGBA{G: 255, A: 255}).OnClick(func() {
            setCount(count + 1)
        }),
    )
}
```

### 2. Hooks System

- `useState`: Manage component state
- `useNavigate`: Implement programmatic navigation

### 3. Router System

```go
// Define routes
routes := []ui.RouteProps{
    {Path: "/", Component: Counter, Props: CounterProps{InitialCount: 0}},
    {Path: "/about", Component: AboutPage, Props: AboutPageProps{}},
}

// Create router
router := ui.NewRouter(routes)
```

### 4. View Construction

```go
// Build view using chained API
view := ui.View(
    ui.Text().Content("Hello, Tenon!").FontSize(24),
    ui.View(
        ui.Text().Content("Subview 1"),
        ui.Text().Content("Subview 2"),
    ).FlexDirection(yoga.FlexDirectionRow),
).
    Width(ui.Percent(100)).
    Height(ui.Percent(100)).
    Background(color.NRGBA{B: 255, A: 128})
```

## 📁 Project Structure

```
tenon/
├── core/
│   ├── gio/          # Gio integration
│   │   └── app.go    # Application startup and window management
│   └── ui/           # Core UI library
│       ├── binding.go # Components and Hooks system
│       ├── element.go # Element tree management
│       ├── ui.go      # View construction API
│       └── render/    # Render backend
│           ├── click_able.go  # Click event handling
│           ├── render.go       # Render base class
│           ├── text.go         # Text rendering
│           └── tree.go         # RenderObject tree
├── example/          # Example applications
│   └── data_binding_demo.go # Data binding example
└── yoga/             # Yoga layout engine
    ├── enum.go       # Enum definitions
    └── style.go      # Style system
```

## 📖 Usage Guide

### Creating Components

1. Define component props structure
2. Create functional component that receives props and returns UI
3. Use Hooks to manage state inside components
4. Build views using chained API

### Using Hooks

```go
// State management
count, setCount := ui.UseState(0)

// Navigation
navigate := ui.UseNavigate()
navigate("/about")
```

### Configuring Routes

1. Define route configuration, mapping paths to components
2. Create router manager
3. Use `useNavigate` in components for programmatic navigation

## 🤝 Contribution Guide

1. Fork this repository
2. Create feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Open Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## 📞 Contact

For questions or suggestions, please submit Issues or Pull Requests.

---

**Tenon** - Making Go UI development simpler and more efficient! 🎉