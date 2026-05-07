package tenon

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
)

// App 是 go-tui 风格的应用对象，包装 tenon 引擎。
type App struct {
	engine  *ui.Engine
	stopped bool
	width   int
	height  int
	title   string
	root    Widget
}

// AppOption 是 App 的配置选项。
type AppOption func(*App) error

// WithSize 设置窗口初始大小。
func WithSize(width, height int) AppOption {
	return func(a *App) error {
		a.width = width
		a.height = height
		return nil
	}
}

// WithTitle 设置窗口标题。
func WithTitle(title string) AppOption {
	return func(a *App) error {
		a.title = title
		return nil
	}
}

// WithRootComponent 设置根组件。
func WithRootComponent(comp Component) AppOption {
	return func(a *App) error {
		a.root = newComponentWidget(comp, a)
		return nil
	}
}

// WithRootWidget 设置根 Widget。
func WithRootWidget(widget Widget) AppOption {
	return func(a *App) error {
		a.root = widget
		return nil
	}
}

// NewApp 创建一个新的 App。
func NewApp(opts ...AppOption) (*App, error) {
	app := &App{
		width:  800,
		height: 600,
	}
	for _, opt := range opts {
		if err := opt(app); err != nil {
			return nil, err
		}
	}
	return app, nil
}

// Run 启动应用事件循环。阻塞直到应用退出。
func (a *App) Run() error {
	if a.root == nil {
		panic("tenon.App: no root component or widget set")
	}

	wrappedBuild := func() ui.Widget {
		return ui.ThemeProvider(ui.GetTheme(), a.root)
	}

	a.engine = ui.NewEngine(wrappedBuild, a.width, a.height)
	a.engine.Mount()

	if a.title != "" {
		ebiten.SetWindowTitle(a.title)
	}
	ebiten.SetWindowSize(a.width, a.height)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := &appGame{app: a, engine: a.engine}
	return ebiten.RunGame(game)
}

// Stop 请求停止应用事件循环。
func (a *App) Stop() {
	a.stopped = true
}

// Close 关闭应用。等价于 Stop。
func (a *App) Close() error {
	a.Stop()
	return nil
}

// Rebuild 标记应用需要在下一帧重渲染。
func (a *App) Rebuild() {
	if a.engine != nil {
		a.engine.Rebuild()
	}
}

// appGame 将 Engine 包装为 ebiten.Game，以支持 Stop。
type appGame struct {
	app    *App
	engine *ui.Engine
}

func (g *appGame) Update() error {
	if g.app.stopped {
		return ebiten.Termination
	}
	return g.engine.Update()
}

func (g *appGame) Draw(screen *ebiten.Image) {
	g.engine.Draw(screen)
}

func (g *appGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.engine.Layout(outsideWidth, outsideHeight)
}
