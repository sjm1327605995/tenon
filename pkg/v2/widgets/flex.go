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

func (r RowWidget) Gapf(v float32) RowWidget { return r.Gap(v) }
func (r RowWidget) Gap(v float32) RowWidget {
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

func (r RowWidget) Paddingf(insets ui.EdgeInsets) RowWidget { return r.Padding(insets) }
func (r RowWidget) Padding(insets ui.EdgeInsets) RowWidget {
	r.padding = insets
	return r
}

func (r RowWidget) CreateElement() ui.Element {
	return ui.NewMultiChildRenderObjectElement(r)
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

func (c ColumnWidget) Gapf(v float32) ColumnWidget { return c.Gap(v) }
func (c ColumnWidget) Gap(v float32) ColumnWidget {
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

func (c ColumnWidget) Paddingf(insets ui.EdgeInsets) ColumnWidget { return c.Padding(insets) }
func (c ColumnWidget) Padding(insets ui.EdgeInsets) ColumnWidget {
	c.padding = insets
	return c
}

func (c ColumnWidget) CreateElement() ui.Element {
	return ui.NewMultiChildRenderObjectElement(c)
}

// CreateRenderObject implements RenderObjectFactory.
func (r RowWidget) CreateRenderObject(element ui.Element) render.RenderObject {
	ro := render.NewRenderFlex()
	ro.StyleSetFlexDirection(ui.FlexDirectionRow)
	applyFlexProps(ro, 0, ui.JustifyFlexStart, ui.AlignFlexStart, ui.WrapNoWrap, ui.EdgeInsets{}, r.gap, r.justify, r.alignItems, r.wrapMode, r.padding)
	return ro
}

// UpdateRenderObject implements RenderObjectUpdater.
func (r RowWidget) UpdateRenderObject(ro render.RenderObject, oldWidget ui.Widget) {
	old := oldWidget.(RowWidget)
	applyFlexProps(ro.(*render.RenderFlex), old.gap, old.justify, old.alignItems, old.wrapMode, old.padding, r.gap, r.justify, r.alignItems, r.wrapMode, r.padding)
}

// GetChildrenWidgets implements MultiChildProvider.
func (r RowWidget) GetChildrenWidgets() []ui.Widget {
	return r.children
}

// CreateRenderObject implements RenderObjectFactory.
func (c ColumnWidget) CreateRenderObject(element ui.Element) render.RenderObject {
	ro := render.NewRenderFlex()
	ro.StyleSetFlexDirection(ui.FlexDirectionColumn)
	applyFlexProps(ro, 0, ui.JustifyFlexStart, ui.AlignFlexStart, ui.WrapNoWrap, ui.EdgeInsets{}, c.gap, c.justify, c.alignItems, c.wrapMode, c.padding)
	return ro
}

// UpdateRenderObject implements RenderObjectUpdater.
func (c ColumnWidget) UpdateRenderObject(ro render.RenderObject, oldWidget ui.Widget) {
	old := oldWidget.(ColumnWidget)
	applyFlexProps(ro.(*render.RenderFlex), old.gap, old.justify, old.alignItems, old.wrapMode, old.padding, c.gap, c.justify, c.alignItems, c.wrapMode, c.padding)
}

// GetChildrenWidgets implements MultiChildProvider.
func (c ColumnWidget) GetChildrenWidgets() []ui.Widget {
	return c.children
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
