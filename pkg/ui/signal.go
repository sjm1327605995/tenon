package ui

import "reflect"

// Signal 是细粒度可观察值：Set 时不触发组件重渲染，直接把新值推给订阅者。
// 引擎用它实现叶子级数据绑定（Bind/BindInput）：订阅回调直接写 renderNode，
// 一次 O(1) 更新，不经过 reconcile；只有文本宽度变化才重排。
//
// 与 UseState / Store 的分工：
//   - UseState：组件内部状态，变化重渲染本组件；
//   - Store：树外共享的业务状态，变化重渲染订阅组件（粒度=组件）；
//   - Signal：跨组件共享、且变化频率高（每帧/每秒多次）的叶子值，直写 renderNode（粒度=节点）。
//
// 线程约定与 setState 一致：Set/Subscribe 须在 UI 线程（事件回调、render、ui.Post 内）调用。
type Signal[T any] struct {
	v    T
	subs map[uint64]func(T)
	seq  uint64
}

// NewSignal 创建信号，initial 为初始值。
func NewSignal[T any](initial T) *Signal[T] { return &Signal[T]{v: initial} }

// Get 返回当前值。render 期读取即参与初始显示；要响应变化请用 Bind/BindInput 或 Subscribe。
func (s *Signal[T]) Get() T { return s.v }

// Set 写入新值并通知订阅者；与旧值深度相等时不通知（同 UseState 的 bailout 语义）。
func (s *Signal[T]) Set(v T) {
	if reflect.DeepEqual(s.v, v) {
		return
	}
	s.v = v
	for _, fn := range s.subs {
		fn(v)
	}
}

// Subscribe 注册变更回调，返回退订函数。注册时不会立即回调当前值。
// 引擎内部绑定（Bind/BindInput）经此订阅；自定义集成（如把信号画进自定义组件）
// 也可直接使用，回调里需要更新界面时请自行标脏（g.needsLayout 不可达时可走 Post + Store）。
func (s *Signal[T]) Subscribe(fn func(T)) (unsubscribe func()) {
	s.seq++
	id := s.seq
	if s.subs == nil {
		s.subs = map[uint64]func(T){}
	}
	s.subs[id] = fn
	return func() { delete(s.subs, id) }
}

// key 返回信号的稳定标识，绑定去重（同一 renderNode 对同一信号只订阅一次）。
func (s *Signal[T]) key() uintptr { return reflect.ValueOf(s).Pointer() }

// sigSub 是一条活跃的信号订阅：unsub 退订；fire 是当前生效的更新闭包——
// 它随每次 render 重新 attach 而刷新（formatter/effect 每轮 render 都是新闭包，
// 不能靠函数身份识别新旧）。
type sigSub struct {
	unsub func()
	fire  func(any)
}

// ---- 文本绑定（Bind）----

// textBinding 是 Signal → 文本节点的类型擦除绑定（*Binding[T] 实现）。
type textBinding interface {
	attach(f *Fiber, rn *renderNode)
}

// Binding 把一个信号绑定到文本节点的内容：信号变化时直写节点文字，
// 所在组件不重渲染。由 Bind 构造，混在 ui.Text 的参数里传入。
type Binding[T any] struct {
	sig    *Signal[T]
	format func(T) string
}

// Bind 把 sig 绑定到 ui.Text 的文本内容。format 把信号值格式化成显示字符串。
// s 可传 ""（占位，挂载时立即用 format(sig.Get()) 回填）。
//
//	ui.Text("", ui.Bind(fps, func(v float64) string { return fmt.Sprintf("%.0f", v) }))
//
// 注意：只支持普通 Text；RichText 的 run 在 render 期快照，绑定不生效。
func Bind[T any](sig *Signal[T], format func(T) string) *Binding[T] {
	return &Binding[T]{sig: sig, format: format}
}

func (b *Binding[T]) attach(f *Fiber, rn *renderNode) {
	key := b.sig.key()
	cell := rn.sigSubs[key]
	if cell == nil {
		if rn.sigSubs == nil {
			rn.sigSubs = map[uintptr]*sigSub{}
		}
		cell = &sigSub{}
		rn.sigSubs[key] = cell
		cell.unsub = b.sig.Subscribe(func(v T) { cell.fire(v) })
		f.addCleanup(cell.unsub)
	}
	// 刷新生效闭包（捕获本轮 render 的 formatter），并立即对齐当前值。
	cell.fire = func(a any) { rn.applyBoundText(b.format(a.(T))) }
	cell.fire(b.sig.Get())
}

// applyBoundText 由 Signal 绑定回调直写文本：不经过组件重渲染。
// 测量函数在布局时动态读 rn.text，故只需标脏节点并请求一次重排。
func (rn *renderNode) applyBoundText(s string) {
	if rn.text == s {
		return
	}
	rn.text = s
	rn.yn.MarkDirty()
	if g := activeGame; g != nil {
		g.needsLayout = true
	}
}

// ---- 输入框双向绑定（BindInput）----

// valueBinding 是 Signal[string] → 输入框值的类型擦除绑定（*inputBinding 实现）。
type valueBinding interface {
	attach(f *Fiber, rn *renderNode)
	current() string
}

// inputBinding 把信号双向接到输入框：输入时写信号（onChange），外部 Set 时直写
// rn.value（不走组件重渲染）。打字路径由引擎直接改 rn.value 再回调 onChange，
// 信号回写时 rn.value 已是最新，bailout 不会自激循环。
type inputBinding struct {
	sig *Signal[string]
}

func (vb *inputBinding) current() string { return vb.sig.Get() }

func (vb *inputBinding) attach(f *Fiber, rn *renderNode) {
	key := vb.sig.key()
	cell := rn.inputSubs[key]
	if cell == nil {
		if rn.inputSubs == nil {
			rn.inputSubs = map[uintptr]*sigSub{}
		}
		cell = &sigSub{}
		rn.inputSubs[key] = cell
		cell.unsub = vb.sig.Subscribe(func(v string) { cell.fire(v) })
		f.addCleanup(cell.unsub)
	}
	cell.fire = func(a any) {
		v := a.(string)
		if rn.value == v {
			return
		}
		rn.value = v
		if rn.caretPos > len(v) {
			rn.caretPos = len(v)
		}
		if rn.selAnchor > len(v) {
			rn.selAnchor = len(v)
		}
		rn.yn.MarkDirty()
		if g := activeGame; g != nil {
			g.needsLayout = true
		}
	}
	cell.fire(vb.sig.Get())
}

// BindInput 把信号双向绑定到输入框的值：显示值由信号直写（外部 Set 不走重渲染），
// 输入经 onChange 回写信号。已内含 OnChange 语义，不要再叠加 ui.OnChange（后者会覆盖回写）。
// 低层用法：ui.Input(ui.Style(...), ui.BindInput(sig))。
func BindInput(sig *Signal[string]) *Node {
	vb := &inputBinding{sig: sig}
	return &Node{typ: typeAttr, applyAttr: func(hp *hostProps) {
		hp.bindValue = vb
		hp.hasValue = true
		hp.onChange = func(v string) { sig.Set(v) }
	}}
}

func attachTextBinds(f *Fiber, binds []textBinding, rn *renderNode) {
	for _, b := range binds {
		b.attach(f, rn)
	}
}

// 约束：Binding 必须实现 textBinding；inputBinding 必须实现 valueBinding。
var (
	_ textBinding  = (*Binding[int])(nil)
	_ valueBinding = (*inputBinding)(nil)
)
