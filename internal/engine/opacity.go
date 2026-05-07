package engine

import "github.com/sjm1327605995/tenon/internal/render"

// OpacityWidget 对子 Widget 应用透明度。
// 轻量级版本，用于 Navigator 转场等 ui 包内部场景。
// 完整变换请使用 widgets.TransformWidget。
type OpacityWidget struct {
	BaseWidget
	Child Widget
	Alpha float32
}

// Opacity 创建透明度包装器。
func Opacity(child Widget, alpha float32) OpacityWidget {
	if alpha < 0 {
		alpha = 0
	}
	if alpha > 1 {
		alpha = 1
	}
	return OpacityWidget{Child: child, Alpha: alpha}
}

func (o OpacityWidget) CreateElement() Element {
	e := &opacityElement{}
	e.SingleChildComponentElement.ComponentElement.BaseElement.Init(e, o)
	return e
}

type opacityElement struct {
	SingleChildComponentElement
}

func (e *opacityElement) Mount(parent Element, slot int) {
	e.SingleChildComponentElement.Mount(parent, slot)
	w, ok := e.GetWidget().(OpacityWidget)
	if !ok {
		return
	}
	if w.Child != nil {
		e.Child = UpdateChild(e, nil, w.Child)
	}
	e.applyAlpha(w.Alpha)
}

func (e *opacityElement) PerformRebuild(oldWidget Widget) {
	w, ok := e.GetWidget().(OpacityWidget)
	if !ok {
		return
	}
	e.Child = UpdateChild(e, e.Child, w.Child)
	e.applyAlpha(w.Alpha)
}

func (e *opacityElement) applyAlpha(alpha float32) {
	if ro := e.Child.FindRenderObject(); ro != nil {
		t := ro.GetTransform()
		t.Alpha = alpha
		ro.SetTransform(t)
	}
}

// ==================== SlideOffsetWidget ====================

// SlideOffsetWidget 对子 Widget 应用偏移量（用于滑动转场）。
// 偏移量是百分比相对于子元素尺寸，1.0 = 一个元素宽/高。
type SlideOffsetWidget struct {
	BaseWidget
	Child    Widget
	OffsetX  float32 // 水平偏移比例，正值向右
	OffsetY  float32 // 垂直偏移比例，正值向下
}

// SlideOffset 创建滑动偏移包装器。
func SlideOffset(child Widget, offsetX, offsetY float32) SlideOffsetWidget {
	return SlideOffsetWidget{Child: child, OffsetX: offsetX, OffsetY: offsetY}
}

func (s SlideOffsetWidget) CreateElement() Element {
	e := &slideOffsetElement{}
	e.SingleChildComponentElement.ComponentElement.BaseElement.Init(e, s)
	return e
}

type slideOffsetElement struct {
	SingleChildComponentElement
}

func (e *slideOffsetElement) Mount(parent Element, slot int) {
	e.SingleChildComponentElement.Mount(parent, slot)
	w, ok := e.GetWidget().(SlideOffsetWidget)
	if !ok {
		return
	}
	if w.Child != nil {
		e.Child = UpdateChild(e, nil, w.Child)
	}
	e.applyOffset(w.OffsetX, w.OffsetY)
}

func (e *slideOffsetElement) PerformRebuild(oldWidget Widget) {
	w, ok := e.GetWidget().(SlideOffsetWidget)
	if !ok {
		return
	}
	e.Child = UpdateChild(e, e.Child, w.Child)
	e.applyOffset(w.OffsetX, w.OffsetY)
}

func (e *slideOffsetElement) applyOffset(offsetX, offsetY float32) {
	if ro := e.Child.FindRenderObject(); ro != nil {
		bounds := ro.GetBounds()
		t := ro.GetTransform()
		t.TranslateX = bounds.Width * offsetX
		t.TranslateY = bounds.Height * offsetY
		ro.SetTransform(t)
	}
}

func (e *slideOffsetElement) FindRenderObject() render.RenderObject {
	if e.Child != nil {
		return e.Child.FindRenderObject()
	}
	return nil
}
