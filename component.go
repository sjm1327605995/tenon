package tenon

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/internal/engine"
	"github.com/sjm1327605995/tenon/internal/render"
)

// Component 是 go-tui 风格的组件接口。
// 用户自定义组件只需实现 Render 方法，直接返回 Widget 树即可。
type Component interface {
	Render(app *App) Widget
}

// KeyMapper 是组件的可选接口，用于声明键盘绑定。
// 每次重渲染时都会重新收集，可基于状态返回不同的绑定。
type KeyMapper interface {
	KeyMap() KeyMap
}

// Initializer 是组件的可选接口，用于在组件首次挂载时执行初始化。
// 返回的函数在组件卸载时调用（cleanup）。
type Initializer interface {
	Init() func()
}

// AppBinder 是组件的可选接口，用于将组件内的 State 绑定到 App。
type AppBinder interface {
	BindApp(app *App)
}

// KeyMap 是一组键盘绑定。
type KeyMap []KeyBinding

// KeyBinding 描述一个按键及其处理函数。
type KeyBinding struct {
	Key       ebiten.Key
	Modifiers []engine.ShortcutKey
	Handler   func(KeyEvent)
}

// KeyEvent 在按键触发时传递给处理函数。
type KeyEvent struct {
	App *App
}

// On 创建一个 KeyBinding。
func On(key ebiten.Key, handler func(KeyEvent)) KeyBinding {
	return KeyBinding{Key: key, Handler: handler}
}

// ==================== 内核适配器：Component → Widget → Element ====================

// componentWidget 直接实现 Widget 接口，作为 Component 的轻量包装。
type componentWidget struct {
	BaseWidget
	comp Component
	app  *App
}

func newComponentWidget(comp Component, app *App) *componentWidget {
	return &componentWidget{comp: comp, app: app}
}

func (c *componentWidget) CreateElement() Element {
	e := &componentElement{widget: c}
	e.Init(e, c)
	return e
}

// componentElement 直接管理 Component 的生命周期与子树，不经过 StatefulElement/State。
type componentElement struct {
	engine.BaseElement
	child    Element
	cleanup  func()
	bindings []shortcutReg
	widget   *componentWidget
}

type shortcutReg struct {
	key       ebiten.Key
	modifiers []engine.ShortcutKey
}

func (e *componentElement) Mount(parent Element, slot int) {
	e.BaseElement.Mount(parent, slot)

	if binder, ok := e.widget.comp.(AppBinder); ok {
		binder.BindApp(e.widget.app)
	}
	if init, ok := e.widget.comp.(Initializer); ok {
		e.cleanup = init.Init()
	}
	if km, ok := e.widget.comp.(KeyMapper); ok {
		e.registerKeyMap(e.widget.app, km.KeyMap())
	}

	w := e.widget.comp.Render(e.widget.app)
	e.child = engine.UpdateChild(e, nil, w)
}

func (e *componentElement) Update(newWidget Widget) {
	e.BaseElement.Update(newWidget)
	e.widget = newWidget.(*componentWidget)

	// 支持动态 KeyMap：注销旧的，注册新的
	e.unregisterKeyMap()
	if km, ok := e.widget.comp.(KeyMapper); ok {
		e.registerKeyMap(e.widget.app, km.KeyMap())
	}

	w := e.widget.comp.Render(e.widget.app)
	e.child = engine.UpdateChild(e, e.child, w)
}

func (e *componentElement) Unmount() {
	e.unregisterKeyMap()
	if e.cleanup != nil {
		e.cleanup()
		e.cleanup = nil
	}
	if e.child != nil {
		e.child.Unmount()
		e.child = nil
	}
	e.BaseElement.Unmount()
}

func (e *componentElement) GetChildren() []Element {
	if e.child == nil {
		return nil
	}
	return []Element{e.child}
}

func (e *componentElement) FindRenderObject() render.RenderObject {
	if e.child != nil {
		return e.child.FindRenderObject()
	}
	return nil
}

func (e *componentElement) registerKeyMap(app *App, km KeyMap) {
	if app == nil || app.engine == nil {
		return
	}
	sm := app.engine.GetShortcutManager()
	for _, kb := range km {
		kb := kb
		sc := engine.Shortcut{
			Key:       kb.Key,
			Modifiers: kb.Modifiers,
			Handler: func() {
				if kb.Handler != nil {
					kb.Handler(KeyEvent{App: app})
				}
			},
			Global: true,
		}
		sm.Register(sc)
		e.bindings = append(e.bindings, shortcutReg{key: kb.Key, modifiers: kb.Modifiers})
	}
}

func (e *componentElement) unregisterKeyMap() {
	if len(e.bindings) == 0 {
		return
	}
	if e.widget.app == nil || e.widget.app.engine == nil {
		e.bindings = nil
		return
	}
	sm := e.widget.app.engine.GetShortcutManager()
	for _, b := range e.bindings {
		sm.Unregister(b.key, b.modifiers...)
	}
	e.bindings = nil
}
