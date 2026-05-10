package engine

import (
	"reflect"
	"time"

	"github.com/sjm1327605995/tenon/pkg/render"
)

// ==================== Route & Page ====================

// RouteBuilder 是页面构建函数，接收 BuildContext 返回 Widget。
type RouteBuilder func(ctx BuildContext, params RouteParams) Widget

// RouteParams 是路由参数，用于在页面间传递数据。
type RouteParams map[string]any

// Get 获取字符串参数。
func (p RouteParams) Get(key string) string {
	if p == nil {
		return ""
	}
	if v, ok := p[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetInt 获取整数参数。
func (p RouteParams) GetInt(key string) int {
	if p == nil {
		return 0
	}
	if v, ok := p[key]; ok {
		if i, ok := v.(int); ok {
			return i
		}
	}
	return 0
}

// Page 表示导航栈中的一个页面。
type Page struct {
	Name    string
	Params  RouteParams
	Builder RouteBuilder
}

// ==================== Transition ====================

// PageTransition 定义页面转场动画类型。
type PageTransition int

const (
	// TransitionNone 无动画，直接切换。
	TransitionNone PageTransition = iota
	// TransitionSlide 从右向左滑入（push），从左向右滑出（pop）。
	TransitionSlide
	// TransitionFade 淡入淡出。
	TransitionFade
	// TransitionSlideUp 从下向上滑入。
	TransitionSlideUp
)

// Duration 返回各转场类型的默认时长。
func (t PageTransition) Duration() time.Duration {
	switch t {
	case TransitionNone:
		return 0
	default:
		return 300 * time.Millisecond
	}
}

// ==================== NavigatorState (公共接口) ====================

// NavigatorState 对外暴露的导航操作接口。
type NavigatorState interface {
	Push(name string, params ...RouteParams)
	Pop()
	PopToRoot()
	PushReplacement(name string, params ...RouteParams)
	CurrentPage() string
	PageCount() int
}

// ==================== NavigatorContext ====================

var navigatorContextType = reflect.TypeOf(NavigatorContext{})

// NavigatorContext 是一个 InheritedWidget，向下传递导航操作能力。
type NavigatorContext struct {
	BaseWidget
	Navigator NavigatorState
	Child     Widget
}

// NewNavigatorContext 创建 NavigatorContext。
func NewNavigatorContext(nav NavigatorState, child Widget) NavigatorContext {
	return NavigatorContext{Navigator: nav, Child: child}
}

func (n NavigatorContext) CreateElement() Element {
	return NewInheritedElement(n)
}

func (n NavigatorContext) UpdateShouldNotify(oldWidget InheritedWidget) bool {
	old, ok := oldWidget.(NavigatorContext)
	if !ok {
		return true
	}
	return n.Navigator != old.Navigator
}

func (n NavigatorContext) BuildChild(ctx BuildContext) Widget {
	return n.Child
}

// GetNavigator 从 BuildContext 获取导航器状态。
func GetNavigator(ctx BuildContext) NavigatorState {
	if ctx == nil {
		return nil
	}
	iw, ok := ctx.DependOnInheritedWidgetOfExactType(navigatorContextType)
	if !ok || iw == nil {
		return nil
	}
	nc, ok := iw.(NavigatorContext)
	if !ok {
		return nil
	}
	return nc.Navigator
}

// ==================== Navigator Widget ====================

// NavigatorWidget 是路由导航器，管理页面栈和转场动画。
type NavigatorWidget struct {
	BaseWidget
	routes     map[string]RouteBuilder
	initial    string
	params     RouteParams
	transition PageTransition
}

// Navigator 创建路由导航器 Widget。
func Navigator(routes map[string]RouteBuilder, initial string) NavigatorWidget {
	return NavigatorWidget{
		routes:     routes,
		initial:    initial,
		transition: TransitionSlide,
	}
}

// NavigatorWithParams 创建带初始参数的路由导航器。
func NavigatorWithParams(routes map[string]RouteBuilder, initial string, params RouteParams) NavigatorWidget {
	return NavigatorWidget{
		routes:     routes,
		initial:    initial,
		params:     params,
		transition: TransitionSlide,
	}
}

// WithTransition 设置转场动画类型。
func (n NavigatorWidget) WithTransition(t PageTransition) NavigatorWidget {
	n.transition = t
	return n
}

func (n NavigatorWidget) CreateElement() Element {
	return NewStatefulElement(n)
}

func (n NavigatorWidget) CreateState() State {
	s := &navigatorState{}
	s.Init(s)
	return s
}

// ==================== Navigator State ====================

type navigatorState struct {
	BaseStateOf[NavigatorWidget]
	pageStack []Page
	transType PageTransition
	animating bool
	animCtrl  *AnimationController
	progress  float32
}

func (s *navigatorState) InitState() {
	w := s.Widget()
	s.transType = w.transition
	if w.initial != "" {
		s.pageStack = []Page{
			{Name: w.initial, Params: w.params, Builder: w.routes[w.initial]},
		}
	}
}

func (s *navigatorState) Dispose() {
	s.stopAnim()
}

func (s *navigatorState) DidUpdateWidget(oldWidget Widget) {}

// ---- NavigatorState 接口 ----

func (s *navigatorState) Push(name string, params ...RouteParams) {
	w := s.Widget()
	builder, ok := w.routes[name]
	if !ok {
		return
	}
	var p RouteParams
	if len(params) > 0 {
		p = params[0]
	}
	s.pageStack = append(s.pageStack, Page{Name: name, Params: p, Builder: builder})
	s.startTransition(true)
}

func (s *navigatorState) Pop() {
	if len(s.pageStack) <= 1 {
		return
	}
	s.pageStack = s.pageStack[:len(s.pageStack)-1]
	s.startTransition(false)
}

func (s *navigatorState) PopToRoot() {
	if len(s.pageStack) <= 1 {
		return
	}
	s.pageStack = s.pageStack[:1]
	s.startTransition(false)
}

func (s *navigatorState) PushReplacement(name string, params ...RouteParams) {
	w := s.Widget()
	builder, ok := w.routes[name]
	if !ok {
		return
	}
	var p RouteParams
	if len(params) > 0 {
		p = params[0]
	}
	page := Page{Name: name, Params: p, Builder: builder}
	if len(s.pageStack) > 0 {
		s.pageStack[len(s.pageStack)-1] = page
	} else {
		s.pageStack = append(s.pageStack, page)
	}
	s.startTransition(true)
}

func (s *navigatorState) CurrentPage() string {
	if len(s.pageStack) == 0 {
		return ""
	}
	return s.pageStack[len(s.pageStack)-1].Name
}

func (s *navigatorState) PageCount() int {
	return len(s.pageStack)
}

// ---- 转场动画 ----

func (s *navigatorState) startTransition(forward bool) {
	if s.transType == TransitionNone {
		s.animating = false
		s.SetState(nil)
		return
	}
	s.stopAnim()
	s.animating = true
	s.animCtrl = &AnimationController{
		Duration:   s.transType.Duration(),
		LowerBound: 0,
		UpperBound: 1,
	}
	if forward {
		s.animCtrl.Value = 0
		s.animCtrl.Forward()
	} else {
		s.animCtrl.Value = 1
		s.animCtrl.Reverse()
	}
	s.animCtrl.AddListener(func() {
		s.progress = float32(s.animCtrl.Value)
		if s.animCtrl.Status == AnimationCompleted || s.animCtrl.Status == AnimationDismissed {
			s.animating = false
		}
		s.SetState(nil)
	})
	if defaultEngine != nil {
		defaultEngine.RegisterAnimation(s.animCtrl)
	}
}

func (s *navigatorState) stopAnim() {
	if s.animCtrl != nil {
		s.animCtrl.Stop()
		if defaultEngine != nil {
			defaultEngine.UnregisterAnimation(s.animCtrl)
		}
		s.animCtrl = nil
	}
}

// ==================== navTransitionWidget ====================

// navTransitionWidget 是导航器转场容器，始终使用 RenderStack 作为根 render object。
// 这样可以避免转场时 widget 类型变化导致的旧页面 unmount 和闪烁。
type navTransitionWidget struct {
	BaseWidget
	children []Widget
}

// navTransitionElement 在 mount/update 时自动将子页面拉伸填满父容器。
type navTransitionElement struct {
	MultiChildRenderObjectElement
}

func (e *navTransitionElement) Mount(parent Element, slot int) {
	e.MultiChildRenderObjectElement.Mount(parent, slot)
	e.stretchChildren()
}

func (e *navTransitionElement) Update(newWidget Widget) {
	e.MultiChildRenderObjectElement.Update(newWidget)
	e.stretchChildren()
}

func (e *navTransitionElement) stretchChildren() {
	for _, child := range e.Children {
		if ro := child.FindRenderObject(); ro != nil {
			if y := ro.GetYoga(); y != nil {
				y.StyleSetWidthPercent(100)
				y.StyleSetHeightPercent(100)
			}
		}
	}
}

func (n navTransitionWidget) CreateElement() Element {
	e := &navTransitionElement{}
	e.MultiChildRenderObjectElement.BaseElement.Init(e, n)
	return e
}

func (n navTransitionWidget) CreateRenderObject(element Element) render.RenderObject {
	return render.NewRenderStack()
}

func (n navTransitionWidget) UpdateRenderObject(ro render.RenderObject, oldWidget Widget) {
	// Stack 不需要更新
}

func (n navTransitionWidget) GetChildrenWidgets() []Widget {
	return n.children
}

// ---- Build ----

func (s *navigatorState) Build(ctx BuildContext) Widget {
	if len(s.pageStack) == 0 {
		return buildEmptyPage()
	}

	w := s.Widget()
	current := s.pageStack[len(s.pageStack)-1]

	// 构建当前页面
	currentBuilder := NewBuilder(func(innerCtx BuildContext) Widget {
		var pageWidget Widget
		if current.Builder != nil {
			pageWidget = current.Builder(innerCtx, current.Params)
		} else if builder, ok := w.routes[current.Name]; ok {
			pageWidget = builder(innerCtx, current.Params)
		}
		return pageWidget
	})
	currentPage := NewNavigatorContext(s, currentBuilder)

	// 非动画状态：只显示当前页面
	if !s.animating || len(s.pageStack) < 2 {
		return navTransitionWidget{children: []Widget{Opacity(currentPage, 1)}}
	}

	// 动画状态：同时显示旧页面和新页面
	prev := s.pageStack[len(s.pageStack)-2]
	prevBuilder := NewBuilder(func(innerCtx BuildContext) Widget {
		var pageWidget Widget
		if prev.Builder != nil {
			pageWidget = prev.Builder(innerCtx, prev.Params)
		} else if builder, ok := w.routes[prev.Name]; ok {
			pageWidget = builder(innerCtx, prev.Params)
		}
		return pageWidget
	})
	prevPage := NewNavigatorContext(s, prevBuilder)

	switch s.transType {
	case TransitionFade:
		return navTransitionWidget{children: []Widget{
			Opacity(prevPage, 1-s.progress),
			Opacity(currentPage, s.progress),
		}}
	case TransitionSlide, TransitionSlideUp:
		return navTransitionWidget{children: []Widget{
			SlideOffset(prevPage, -s.progress, 0),
			SlideOffset(currentPage, 1-s.progress, 0),
		}}
	default:
		return navTransitionWidget{children: []Widget{Opacity(currentPage, 1)}}
	}
}

func buildEmptyPage() Widget {
	return navTransitionWidget{children: nil}
}

// ==================== 便捷导航函数 ====================

// NavPush 从 BuildContext 推入新页面。
func NavPush(ctx BuildContext, name string, params ...RouteParams) {
	if nav := GetNavigator(ctx); nav != nil {
		nav.Push(name, params...)
	}
}

// NavPop 从 BuildContext 弹出页面。
func NavPop(ctx BuildContext) {
	if nav := GetNavigator(ctx); nav != nil {
		nav.Pop()
	}
}

// NavPushReplacement 从 BuildContext 替换当前页面。
func NavPushReplacement(ctx BuildContext, name string, params ...RouteParams) {
	if nav := GetNavigator(ctx); nav != nil {
		nav.PushReplacement(name, params...)
	}
}

// NavPopToRoot 从 BuildContext 弹出到栈底。
func NavPopToRoot(ctx BuildContext) {
	if nav := GetNavigator(ctx); nav != nil {
		nav.PopToRoot()
	}
}
