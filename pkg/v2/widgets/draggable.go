package widgets

import (
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
)

// DraggableWidget 让子 Widget 可以被拖拽。
// 拖拽时携带 Data，Engine 会在 mouseDown 时读取并触发拖放协议。
type DraggableWidget struct {
	ui.BaseWidget
	Child ui.Widget
	Data  any // 拖拽时携带的数据（nil 表示不可拖拽）
}

// Draggable 创建可拖拽包装器。
func Draggable(child ui.Widget, data any) DraggableWidget {
	return DraggableWidget{Child: child, Data: data}
}

func (d DraggableWidget) CreateElement() ui.Element {
	e := &DraggableElement{widget: d}
	e.SingleChildComponentElement.ComponentElement.BaseElement.Init(e, d)
	return e
}

// DraggableElement 负责将 Data 同步到子节点的 RenderObject。
type DraggableElement struct {
	ui.SingleChildComponentElement
	widget DraggableWidget
}

func (e *DraggableElement) Update(newWidget ui.Widget) {
	e.widget = newWidget.(DraggableWidget)
	e.SingleChildComponentElement.Update(newWidget)
}

func (e *DraggableElement) UpdateChild(oldWidget ui.Widget) {
	e.Child = ui.UpdateChild(e, e.Child, e.widget.Child)
	if ro := e.Child.FindRenderObject(); ro != nil {
		ro.SetDragData(e.widget.Data)
	}
}

func (e *DraggableElement) Mount(parent ui.Element, slot int) {
	e.SingleChildComponentElement.Mount(parent, slot)
	if e.widget.Child != nil {
		e.Child = ui.UpdateChild(e, nil, e.widget.Child)
	}
	if ro := e.Child.FindRenderObject(); ro != nil {
		ro.SetDragData(e.widget.Data)
	}
}

// DropTargetWidget 让子 Widget 可以接受拖放。
type DropTargetWidget struct {
	ui.BaseWidget
	Child        ui.Widget
	OnWillAccept func(data any) bool // 可选：返回 false 拒绝
	OnAccept     func(data any)      // 必需：drop 成功回调
	OnReject     func(data any)      // 可选：drop 被拒绝回调
}

// DropTarget 创建放置目标包装器。
func DropTarget(child ui.Widget, onAccept func(data any)) DropTargetWidget {
	return DropTargetWidget{Child: child, OnAccept: onAccept}
}

func (d DropTargetWidget) CreateElement() ui.Element {
	e := &DropTargetElement{widget: d}
	e.SingleChildComponentElement.ComponentElement.BaseElement.Init(e, d)
	return e
}

// DropTargetElement 负责将拖放回调同步到子节点的 RenderObject。
type DropTargetElement struct {
	ui.SingleChildComponentElement
	widget DropTargetWidget
}

func (e *DropTargetElement) Update(newWidget ui.Widget) {
	e.widget = newWidget.(DropTargetWidget)
	e.SingleChildComponentElement.Update(newWidget)
}

func (e *DropTargetElement) UpdateChild(oldWidget ui.Widget) {
	e.Child = ui.UpdateChild(e, e.Child, e.widget.Child)
	if ro := e.Child.FindRenderObject(); ro != nil {
		ro.SetOnDragEnter(func(data any) {
			if e.widget.OnWillAccept != nil && !e.widget.OnWillAccept(data) {
				return
			}
			// visual feedback 可在此触发
		})
		ro.SetOnDragLeave(func() {
			// 取消 visual feedback
		})
		ro.SetOnDrop(func(data any) {
			if e.widget.OnWillAccept != nil && !e.widget.OnWillAccept(data) {
				if e.widget.OnReject != nil {
					e.widget.OnReject(data)
				}
				return
			}
			e.widget.OnAccept(data)
		})
	}
}

func (e *DropTargetElement) Mount(parent ui.Element, slot int) {
	e.SingleChildComponentElement.Mount(parent, slot)
	if e.widget.Child != nil {
		e.Child = ui.UpdateChild(e, nil, e.widget.Child)
	}
	if ro := e.Child.FindRenderObject(); ro != nil {
		ro.SetOnDragEnter(func(data any) {
			if e.widget.OnWillAccept != nil && !e.widget.OnWillAccept(data) {
				return
			}
		})
		ro.SetOnDragLeave(func() {})
		ro.SetOnDrop(func(data any) {
			if e.widget.OnWillAccept != nil && !e.widget.OnWillAccept(data) {
				if e.widget.OnReject != nil {
					e.widget.OnReject(data)
				}
				return
			}
			e.widget.OnAccept(data)
		})
	}
}
