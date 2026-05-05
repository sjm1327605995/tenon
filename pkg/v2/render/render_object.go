package render

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/yoga"
)

// Offset 表示二维偏移量。
type Offset struct {
	X, Y float32
}

// Size 表示二维尺寸。
type Size struct {
	Width, Height float32
}

// BoxConstraints 描述布局约束。
type BoxConstraints struct {
	MinWidth, MaxWidth   float32
	MinHeight, MaxHeight float32
}

// RenderObject 是渲染树中的节点，负责布局和绘制。
// 对应 Flutter 的 RenderObject。
type RenderObject interface {
	// 树关系
	GetParent() RenderObject
	SetParent(p RenderObject)
	GetChildren() []RenderObject
	AddChild(child RenderObject)
	RemoveChild(child RenderObject)
	InsertChild(child RenderObject, index int)

	// 生命周期
	Attach(owner *PipelineOwner)
	Detach()
	GetOwner() *PipelineOwner

	// 布局
	MarkNeedsLayout()
	PerformLayout()
	ClearNeedsLayout()
	GetSize() Size
	SetSize(s Size)

	// 绘制
	MarkNeedsPaint()
	ClearNeedsPaint()
	Paint(screen *ebiten.Image, offset Offset)

	// Yoga 节点（需要 Flexbox 布局的 RenderObject 实现）
	GetYoga() *yoga.Node

	// 边界与可见性
	GetBounds() Bounds
	SetBounds(b Bounds)
	IsVisible() bool
	SetVisible(v bool)

	// 变换
	GetTransform() Transform
	SetTransform(t Transform)

	// Z-Index（绘制顺序）
	GetZIndex() int
	SetZIndex(v int)

	// 事件命中测试
	HitTest(x, y float32) bool

	// 拖动事件
	HandleDrag(x, y float32)

	// 鼠标事件回调
	SetOnClick(fn func())
	GetOnClick() func()
	SetOnMouseEnter(fn func())
	GetOnMouseEnter() func()
	SetOnMouseLeave(fn func())
	GetOnMouseLeave() func()
	SetOnMouseDown(fn func())
	GetOnMouseDown() func()
	SetOnMouseUp(fn func())
	GetOnMouseUp() func()

	// 拖放协议
	GetDragData() any
	SetDragData(v any)
	GetOnDragEnter() func(data any)
	SetOnDragEnter(fn func(data any))
	GetOnDragLeave() func()
	SetOnDragLeave(fn func())
	GetOnDrop() func(data any)
	SetOnDrop(fn func(data any))
}

// Bounds 描述 RenderObject 在父坐标系中的位置和尺寸。
type Bounds struct {
	X, Y, Width, Height float32
}

func (b Bounds) Contains(x, y float32) bool {
	return x >= b.X && x < b.X+b.Width && y >= b.Y && y < b.Y+b.Height
}

// BaseRenderObject 提供 RenderObject 的默认实现。
type BaseRenderObject struct {
	self     RenderObject
	parent   RenderObject
	children []RenderObject
	owner    *PipelineOwner
	size     Size
	bounds   Bounds
	visible  bool

	// 脏标记
	needsLayout bool
	needsPaint  bool

	// 鼠标事件回调
	onClick       func()
	onMouseEnter  func()
	onMouseLeave  func()
	onMouseDown   func()
	onMouseUp     func()

	// 拖放协议
	dragData      any
	onDragEnter   func(data any)
	onDragLeave   func()
	onDrop        func(data any)

	// 变换
	transform Transform

	// Z-Index（绘制顺序，值越大越在上层）
	zIndex int

	// 焦点状态
	focused bool

	// Yoga 节点（可选，由子类决定是否创建）
	yoga *yoga.Node
}

func (b *BaseRenderObject) Init(self RenderObject) {
	b.self = self
	b.visible = true
	b.transform = DefaultTransform()
}

// === 树关系 ===

func (b *BaseRenderObject) GetParent() RenderObject  { return b.parent }
func (b *BaseRenderObject) SetParent(p RenderObject) { b.parent = p }
func (b *BaseRenderObject) GetChildren() []RenderObject {
	return b.children
}

func (b *BaseRenderObject) AddChild(child RenderObject) {
	if child == nil {
		return
	}
	b.InsertChild(child, len(b.children))
}

func (b *BaseRenderObject) InsertChild(child RenderObject, index int) {
	if child == nil {
		return
	}
	if child.GetParent() != nil {
		child.GetParent().RemoveChild(child)
	}
	child.SetParent(b.self)
	if index >= len(b.children) {
		b.children = append(b.children, child)
	} else {
		b.children = append(b.children, nil)
		copy(b.children[index+1:], b.children[index:])
		b.children[index] = child
	}
	if b.yoga != nil {
		if childYoga := child.GetYoga(); childYoga != nil {
			b.yoga.InsertChild(childYoga, uint32(index))
		}
	}
	if b.owner != nil {
		child.Attach(b.owner)
	}
	b.self.MarkNeedsLayout()
}

func (b *BaseRenderObject) RemoveChild(child RenderObject) {
	for i, c := range b.children {
		if c == child {
			b.children = append(b.children[:i], b.children[i+1:]...)
			child.SetParent(nil)
			if b.yoga != nil {
				if childYoga := child.GetYoga(); childYoga != nil {
					b.yoga.RemoveChild(childYoga)
				}
			}
			child.Detach()
			b.self.MarkNeedsLayout()
			return
		}
	}
}

// === 生命周期 ===

func (b *BaseRenderObject) Attach(owner *PipelineOwner) {
	b.owner = owner
	if b.needsLayout && b.owner != nil {
		b.owner.NeedsLayoutAdd(b.self)
	}
	for _, child := range b.children {
		if child.GetOwner() == nil {
			child.Attach(owner)
		}
	}
}

func (b *BaseRenderObject) Detach() {
	b.owner = nil
	for _, child := range b.children {
		child.Detach()
	}
}

func (b *BaseRenderObject) GetOwner() *PipelineOwner { return b.owner }

// === 布局 ===

func (b *BaseRenderObject) MarkNeedsLayout() {
	if b.needsLayout {
		return
	}
	b.needsLayout = true
	if b.owner != nil {
		b.owner.NeedsLayoutAdd(b.self)
	}
	// 通知父节点也需要重新布局
	if b.parent != nil {
		b.parent.MarkNeedsLayout()
	}
}

func (b *BaseRenderObject) PerformLayout() {
	// 子类必须覆盖
}

func (b *BaseRenderObject) GetSize() Size   { return b.size }
func (b *BaseRenderObject) SetSize(s Size)  { b.size = s }

// === 绘制 ===

func (b *BaseRenderObject) MarkNeedsPaint() {
	if b.needsPaint {
		return
	}
	b.needsPaint = true
	if b.owner != nil {
		b.owner.NeedsPaintAdd(b.self)
	}
}

func (b *BaseRenderObject) Paint(screen *ebiten.Image, offset Offset) {
	// 子类覆盖实现具体绘制
}

// === Yoga ===

func (b *BaseRenderObject) GetYoga() *yoga.Node { return b.yoga }

// === 边界与可见性 ===

func (b *BaseRenderObject) GetBounds() Bounds    { return b.bounds }
func (b *BaseRenderObject) SetBounds(bounds Bounds) { b.bounds = bounds }
func (b *BaseRenderObject) IsVisible() bool      { return b.visible }
func (b *BaseRenderObject) SetVisible(v bool) {
	b.visible = v
	b.self.MarkNeedsPaint()
}

// === 命中测试 ===

func (b *BaseRenderObject) HitTest(x, y float32) bool {
	if !b.visible {
		return false
	}
	return b.bounds.Contains(x, y)
}

func (b *BaseRenderObject) HandleDrag(x, y float32) {}

func (b *BaseRenderObject) SetOnClick(fn func())       { b.onClick = fn }
func (b *BaseRenderObject) GetOnClick() func()         { return b.onClick }
func (b *BaseRenderObject) SetOnMouseEnter(fn func())  { b.onMouseEnter = fn }
func (b *BaseRenderObject) GetOnMouseEnter() func()    { return b.onMouseEnter }
func (b *BaseRenderObject) SetOnMouseLeave(fn func())  { b.onMouseLeave = fn }
func (b *BaseRenderObject) GetOnMouseLeave() func()    { return b.onMouseLeave }
func (b *BaseRenderObject) SetOnMouseDown(fn func())   { b.onMouseDown = fn }
func (b *BaseRenderObject) GetOnMouseDown() func()     { return b.onMouseDown }
func (b *BaseRenderObject) SetOnMouseUp(fn func())     { b.onMouseUp = fn }
func (b *BaseRenderObject) GetOnMouseUp() func()       { return b.onMouseUp }

func (b *BaseRenderObject) GetDragData() any          { return b.dragData }
func (b *BaseRenderObject) SetDragData(v any)         { b.dragData = v }
func (b *BaseRenderObject) GetOnDragEnter() func(data any) { return b.onDragEnter }
func (b *BaseRenderObject) SetOnDragEnter(fn func(data any)) { b.onDragEnter = fn }
func (b *BaseRenderObject) GetOnDragLeave() func()    { return b.onDragLeave }
func (b *BaseRenderObject) SetOnDragLeave(fn func())  { b.onDragLeave = fn }
func (b *BaseRenderObject) GetOnDrop() func(data any) { return b.onDrop }
func (b *BaseRenderObject) SetOnDrop(fn func(data any)) { b.onDrop = fn }

func (b *BaseRenderObject) GetTransform() Transform {
	return b.transform
}

func (b *BaseRenderObject) SetTransform(t Transform) {
	b.transform = t
	b.MarkNeedsPaint()
}

func (b *BaseRenderObject) GetZIndex() int { return b.zIndex }
func (b *BaseRenderObject) SetZIndex(v int) {
	if b.zIndex == v {
		return
	}
	b.zIndex = v
	b.MarkNeedsPaint()
}

// === 辅助 ===

func (b *BaseRenderObject) ClearNeedsLayout() { b.needsLayout = false }
func (b *BaseRenderObject) ClearNeedsPaint()  { b.needsPaint = false }
func (b *BaseRenderObject) NeedsLayout() bool { return b.needsLayout }
func (b *BaseRenderObject) NeedsPaint() bool  { return b.needsPaint }

// PipelineOwner 管理脏 RenderObject 队列，负责调度和执行布局与绘制。
type PipelineOwner struct {
	needsLayoutList []RenderObject
	needsPaintList  []RenderObject
}

func NewPipelineOwner() *PipelineOwner {
	return &PipelineOwner{}
}

func (p *PipelineOwner) NeedsLayoutAdd(ro RenderObject) {
	for _, existing := range p.needsLayoutList {
		if existing == ro {
			return
		}
	}
	p.needsLayoutList = append(p.needsLayoutList, ro)
}

func (p *PipelineOwner) NeedsPaintAdd(ro RenderObject) {
	for _, existing := range p.needsPaintList {
		if existing == ro {
			return
		}
	}
	p.needsPaintList = append(p.needsPaintList, ro)
}

// FlushLayout 执行所有标记了 needsLayout 的 RenderObject 的布局。
func (p *PipelineOwner) FlushLayout() {
	if len(p.needsLayoutList) == 0 {
		return
	}
	// 复制列表后清空，因为布局过程中可能产生新的 dirty
	list := p.needsLayoutList
	p.needsLayoutList = nil
	for _, ro := range list {
		ro.PerformLayout()
		ro.ClearNeedsLayout()
	}
}

// HasPendingLayout 返回是否有待执行的布局任务。
func (p *PipelineOwner) HasPendingLayout() bool {
	return len(p.needsLayoutList) > 0
}

// FlushPaint 执行所有标记了 needsPaint 的 RenderObject 的绘制准备（此处仅清标记，实际绘制在 Engine.Draw 中遍历）。
func (p *PipelineOwner) FlushPaint() {
	list := p.needsPaintList
	p.needsPaintList = nil
	for _, ro := range list {
		ro.ClearNeedsPaint()
	}
}

// ColorEquals 比较两个 color.Color 是否相等。
func ColorEquals(a, b color.Color) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	ra, ga, ba, aa := a.RGBA()
	rb, gb, bb, ab := b.RGBA()
	return ra == rb && ga == gb && ba == bb && aa == ab
}
