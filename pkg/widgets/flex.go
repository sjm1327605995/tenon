package widgets

import (
	"github.com/sjm1327605995/tenon/pkg/render"
	"github.com/sjm1327605995/tenon/pkg/engine"
	"github.com/sjm1327605995/tenon/yoga"
)

// RowWidget 水平 Flex 容器。
type RowWidget struct {
	engine.BaseWidget
	children     []engine.Widget
	gap          float32
	justify      yoga.Justify
	alignItems   yoga.Align
	alignContent yoga.Align
	wrapMode     yoga.Wrap
	padding      engine.EdgeInsets
}

// Row 创建水平排列容器。
func Row(children ...engine.Widget) RowWidget {
	return RowWidget{
		children:   children,
		justify:    engine.JustifyFlexStart,
		alignItems: engine.AlignFlexStart,
		wrapMode:   engine.WrapNoWrap,
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

func (r RowWidget) Paddingf(insets engine.EdgeInsets) RowWidget { return r.Padding(insets) }
func (r RowWidget) Padding(insets engine.EdgeInsets) RowWidget {
	r.padding = insets
	return r
}

func (r RowWidget) CreateElement() engine.Element {
	return engine.NewMultiChildRenderObjectElement(r)
}

// ColumnWidget 垂直 Flex 容器。
type ColumnWidget struct {
	engine.BaseWidget
	children     []engine.Widget
	gap          float32
	justify      yoga.Justify
	alignItems   yoga.Align
	alignContent yoga.Align
	wrapMode     yoga.Wrap
	padding      engine.EdgeInsets
}

// Column 创建垂直排列容器。
func Column(children ...engine.Widget) ColumnWidget {
	return ColumnWidget{
		children:   children,
		justify:    engine.JustifyFlexStart,
		alignItems: engine.AlignFlexStart,
		wrapMode:   engine.WrapNoWrap,
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

func (c ColumnWidget) Paddingf(insets engine.EdgeInsets) ColumnWidget { return c.Padding(insets) }
func (c ColumnWidget) Padding(insets engine.EdgeInsets) ColumnWidget {
	c.padding = insets
	return c
}

func (c ColumnWidget) CreateElement() engine.Element {
	return engine.NewMultiChildRenderObjectElement(c)
}

// CreateRenderObject implements RenderObjectFactory.
func (r RowWidget) CreateRenderObject(element engine.Element) render.RenderObject {
	ro := render.NewRenderFlex()
	ro.StyleSetFlexDirection(engine.FlexDirectionRow)
	applyFlexProps(ro, 0, engine.JustifyFlexStart, engine.AlignFlexStart, engine.WrapNoWrap, engine.EdgeInsets{}, r.gap, r.justify, r.alignItems, r.wrapMode, r.padding)
	return ro
}

// UpdateRenderObject implements RenderObjectUpdater.
func (r RowWidget) UpdateRenderObject(ro render.RenderObject, oldWidget engine.Widget) {
	old := oldWidget.(RowWidget)
	applyFlexProps(ro.(*render.RenderFlex), old.gap, old.justify, old.alignItems, old.wrapMode, old.padding, r.gap, r.justify, r.alignItems, r.wrapMode, r.padding)
}

// GetChildrenWidgets implements MultiChildProvider.
func (r RowWidget) GetChildrenWidgets() []engine.Widget {
	return r.children
}

// CreateRenderObject implements RenderObjectFactory.
func (c ColumnWidget) CreateRenderObject(element engine.Element) render.RenderObject {
	ro := render.NewRenderFlex()
	ro.StyleSetFlexDirection(engine.FlexDirectionColumn)
	applyFlexProps(ro, 0, engine.JustifyFlexStart, engine.AlignFlexStart, engine.WrapNoWrap, engine.EdgeInsets{}, c.gap, c.justify, c.alignItems, c.wrapMode, c.padding)
	return ro
}

// UpdateRenderObject implements RenderObjectUpdater.
func (c ColumnWidget) UpdateRenderObject(ro render.RenderObject, oldWidget engine.Widget) {
	old := oldWidget.(ColumnWidget)
	applyFlexProps(ro.(*render.RenderFlex), old.gap, old.justify, old.alignItems, old.wrapMode, old.padding, c.gap, c.justify, c.alignItems, c.wrapMode, c.padding)
}

// GetChildrenWidgets implements MultiChildProvider.
func (c ColumnWidget) GetChildrenWidgets() []engine.Widget {
	return c.children
}

// ==================== 辅助函数 ====================

func applyFlexProps(
	r *render.RenderFlex,
	oldGap float32, oldJustify yoga.Justify, oldAlignItems yoga.Align, oldWrapMode yoga.Wrap, oldPadding engine.EdgeInsets,
	newGap float32, newJustify yoga.Justify, newAlignItems yoga.Align, newWrapMode yoga.Wrap, newPadding engine.EdgeInsets,
) {
	if oldGap != newGap {
		r.StyleSetGap(engine.GutterAll, newGap)
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
		r.StyleSetPadding(engine.EdgeTop, newPadding.Top)
		r.StyleSetPadding(engine.EdgeRight, newPadding.Right)
		r.StyleSetPadding(engine.EdgeBottom, newPadding.Bottom)
		r.StyleSetPadding(engine.EdgeLeft, newPadding.Left)
	}
}
