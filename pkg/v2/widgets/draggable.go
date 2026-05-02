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
	e := &DraggableElement{}
	e.SingleChildComponentElement.ComponentElement.BaseElement.Init(e, d)
	return e
}

// DraggableElement 负责将 Data 同步到子节点的 RenderObject。
type DraggableElement struct {
	ui.SingleChildComponentElement
}

func (e *DraggableElement) PerformRebuild(oldWidget ui.Widget) {
	w := e.GetWidget().(DraggableWidget)
	e.Child = ui.UpdateChild(e, e.Child, w.Child)
	if ro := e.Child.FindRenderObject(); ro != nil {
		ro.SetDragData(w.Data)
	}
}

func (e *DraggableElement) Mount(parent ui.Element, slot int) {
	e.SingleChildComponentElement.Mount(parent, slot)
	w := e.GetWidget().(DraggableWidget)
	if w.Child != nil {
		e.Child = ui.UpdateChild(e, nil, w.Child)
	}
	if ro := e.Child.FindRenderObject(); ro != nil {
		ro.SetDragData(w.Data)
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
	e := &DropTargetElement{}
	e.SingleChildComponentElement.ComponentElement.BaseElement.Init(e, d)
	return e
}

// DropTargetElement 负责将拖放回调同步到子节点的 RenderObject。
type DropTargetElement struct {
	ui.SingleChildComponentElement
}

func (e *DropTargetElement) PerformRebuild(oldWidget ui.Widget) {
	w := e.GetWidget().(DropTargetWidget)
	e.Child = ui.UpdateChild(e, e.Child, w.Child)
	if ro := e.Child.FindRenderObject(); ro != nil {
		ro.SetOnDragEnter(func(data any) {
			if w.OnWillAccept != nil && !w.OnWillAccept(data) {
				return
			}
			// visual feedback 可在此触发
		})
		ro.SetOnDragLeave(func() {
			// 取消 visual feedback
		})
		ro.SetOnDrop(func(data any) {
			if w.OnWillAccept != nil && !w.OnWillAccept(data) {
				if w.OnReject != nil {
					w.OnReject(data)
				}
				return
			}
			w.OnAccept(data)
		})
	}
}

func (e *DropTargetElement) Mount(parent ui.Element, slot int) {
	e.SingleChildComponentElement.Mount(parent, slot)
	w := e.GetWidget().(DropTargetWidget)
	if w.Child != nil {
		e.Child = ui.UpdateChild(e, nil, w.Child)
	}
	if ro := e.Child.FindRenderObject(); ro != nil {
		ro.SetOnDragEnter(func(data any) {
			if w.OnWillAccept != nil && !w.OnWillAccept(data) {
				return
			}
		})
		ro.SetOnDragLeave(func() {})
		ro.SetOnDrop(func(data any) {
			if w.OnWillAccept != nil && !w.OnWillAccept(data) {
				if w.OnReject != nil {
					w.OnReject(data)
				}
				return
			}
			w.OnAccept(data)
		})
	}
}
