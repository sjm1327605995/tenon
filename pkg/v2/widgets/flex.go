package widgets

import (
	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
	"github.com/sjm1327605995/tenon/yoga"
)

// RowWidget 水平 Flex 容器。
type RowWidget struct {
	ui.BaseWidget
	children     []ui.Widget
	gap          float32
	justify      yoga.Justify
	alignItems   yoga.Align
	alignContent yoga.Align
	wrapMode     yoga.Wrap
	padding      ui.EdgeInsets
}

// Row 创建水平排列容器。
func Row(children ...ui.Widget) RowWidget {
	return RowWidget{
		children:   children,
		justify:    ui.JustifyFlexStart,
		alignItems: ui.AlignFlexStart,
		wrapMode:   ui.WrapNoWrap,
	}
}

func (r RowWidget) Gapf(v float32) RowWidget {
	r.gap = v
	return r
}

func (r RowWidget) JustifyContent(v yoga.Justify) RowWidget {
	r.justify = v
	return r
}

func (r RowWidget) AlignItems(v yoga.Align) RowWidget {
	r.alignItems = v
	return r
}

func (r RowWidget) Paddingf(insets ui.EdgeInsets) RowWidget {
	r.padding = insets
	return r
}

func (r RowWidget) CreateElement() ui.Element {
	e := &RowElement{}
	e.MultiChildRenderObjectElement.RenderObjectElement.BaseElement.Init(e, r)
	return e
}

// ColumnWidget 垂直 Flex 容器。
type ColumnWidget struct {
	ui.BaseWidget
	children     []ui.Widget
	gap          float32
	justify      yoga.Justify
	alignItems   yoga.Align
	alignContent yoga.Align
	wrapMode     yoga.Wrap
	padding      ui.EdgeInsets
}

// Column 创建垂直排列容器。
func Column(children ...ui.Widget) ColumnWidget {
	return ColumnWidget{
		children:   children,
		justify:    ui.JustifyFlexStart,
		alignItems: ui.AlignFlexStart,
		wrapMode:   ui.WrapNoWrap,
	}
}

func (c ColumnWidget) Gapf(v float32) ColumnWidget {
	c.gap = v
	return c
}

func (c ColumnWidget) JustifyContent(v yoga.Justify) ColumnWidget {
	c.justify = v
	return c
}

func (c ColumnWidget) AlignItems(v yoga.Align) ColumnWidget {
	c.alignItems = v
	return c
}

func (c ColumnWidget) Paddingf(insets ui.EdgeInsets) ColumnWidget {
	c.padding = insets
	return c
}

func (c ColumnWidget) CreateElement() ui.Element {
	e := &ColumnElement{}
	e.MultiChildRenderObjectElement.RenderObjectElement.BaseElement.Init(e, c)
	return e
}

// ==================== RowElement ====================

type RowElement struct {
	ui.MultiChildRenderObjectElement
}

func (e *RowElement) CreateRenderObject() render.RenderObject {
	r := render.NewRenderFlex()
	r.StyleSetFlexDirection(ui.FlexDirectionRow)
	w := e.GetWidget().(RowWidget)
	applyFlexProps(r, 0, ui.JustifyFlexStart, ui.AlignFlexStart, ui.WrapNoWrap, ui.EdgeInsets{}, w.gap, w.justify, w.alignItems, w.wrapMode, w.padding)
	return r
}

func (e *RowElement) UpdateRenderObject(oldWidget ui.Widget) {
	r := e.GetRenderObject().(*render.RenderFlex)
	old := oldWidget.(RowWidget)
	w := e.GetWidget().(RowWidget)
	applyFlexProps(r, old.gap, old.justify, old.alignItems, old.wrapMode, old.padding, w.gap, w.justify, w.alignItems, w.wrapMode, w.padding)
}

func (e *RowElement) UpdateChildren(oldWidget ui.Widget) {
	w := e.GetWidget().(RowWidget)
	newWidgets := make([]ui.Widget, len(w.children))
	for i, child := range w.children {
		newWidgets[i] = child
	}
	e.Children = ui.UpdateChildren(e, e.Children, newWidgets)
}

func (e *RowElement) Mount(parent ui.Element, slot int) {
	e.RenderObject = e.CreateRenderObject()
	e.MultiChildRenderObjectElement.Mount(parent, slot)
	w := e.GetWidget().(RowWidget)
	newWidgets := make([]ui.Widget, len(w.children))
	for i, child := range w.children {
		newWidgets[i] = child
	}
	e.Children = ui.UpdateChildren(e, nil, newWidgets)
}

// ==================== ColumnElement ====================

type ColumnElement struct {
	ui.MultiChildRenderObjectElement
}

func (e *ColumnElement) CreateRenderObject() render.RenderObject {
	r := render.NewRenderFlex()
	r.StyleSetFlexDirection(ui.FlexDirectionColumn)
	w := e.GetWidget().(ColumnWidget)
	applyFlexProps(r, 0, ui.JustifyFlexStart, ui.AlignFlexStart, ui.WrapNoWrap, ui.EdgeInsets{}, w.gap, w.justify, w.alignItems, w.wrapMode, w.padding)
	return r
}

func (e *ColumnElement) UpdateRenderObject(oldWidget ui.Widget) {
	r := e.GetRenderObject().(*render.RenderFlex)
	old := oldWidget.(ColumnWidget)
	w := e.GetWidget().(ColumnWidget)
	applyFlexProps(r, old.gap, old.justify, old.alignItems, old.wrapMode, old.padding, w.gap, w.justify, w.alignItems, w.wrapMode, w.padding)
}

func (e *ColumnElement) UpdateChildren(oldWidget ui.Widget) {
	w := e.GetWidget().(ColumnWidget)
	newWidgets := make([]ui.Widget, len(w.children))
	for i, child := range w.children {
		newWidgets[i] = child
	}
	e.Children = ui.UpdateChildren(e, e.Children, newWidgets)
}

func (e *ColumnElement) Mount(parent ui.Element, slot int) {
	e.RenderObject = e.CreateRenderObject()
	e.MultiChildRenderObjectElement.Mount(parent, slot)
	w := e.GetWidget().(ColumnWidget)
	newWidgets := make([]ui.Widget, len(w.children))
	for i, child := range w.children {
		newWidgets[i] = child
	}
	e.Children = ui.UpdateChildren(e, nil, newWidgets)
}

// ==================== 辅助函数 ====================

func applyFlexProps(
	r *render.RenderFlex,
	oldGap float32, oldJustify yoga.Justify, oldAlignItems yoga.Align, oldWrapMode yoga.Wrap, oldPadding ui.EdgeInsets,
	newGap float32, newJustify yoga.Justify, newAlignItems yoga.Align, newWrapMode yoga.Wrap, newPadding ui.EdgeInsets,
) {
	if oldGap != newGap {
		r.StyleSetGap(ui.GutterAll, newGap)
	}
	if oldJustify != newJustify {
		r.StyleSetJustifyContent(newJustify)
	}
	if oldAlignItems != newAlignItems {
		r.StyleSetAlignItems(newAlignItems)
	}
	if oldWrapMode != newWrapMode {
		r.StyleSetFlexWrap(newWrapMode)
	}
	if oldPadding != newPadding {
		r.StyleSetPadding(ui.EdgeTop, newPadding.Top)
		r.StyleSetPadding(ui.EdgeRight, newPadding.Right)
		r.StyleSetPadding(ui.EdgeBottom, newPadding.Bottom)
		r.StyleSetPadding(ui.EdgeLeft, newPadding.Left)
	}
}
