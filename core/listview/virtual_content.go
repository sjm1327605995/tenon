package listview

import (
	"github.com/sjm1327605995/tenon/event"
	"github.com/sjm1327605995/tenon/geometry"
	"github.com/sjm1327605995/tenon/widget"
)

// virtualContent is an internal widget that represents the entire scrollable
// content area. It reports the total content height to the parent ScrollView
// but only renders visible items.
//
// This is the SwiftUI LazyVStack pattern: the content widget lies about its
// actual children (only visible ones), but truthfully reports its total height
// so the scroll view can calculate correct scrollbar position and size.
type virtualContent struct {
	widget.WidgetBase
	list *Widget // back-reference to parent ListView
}

// Layout returns the total content size: full viewport width and the sum of
// all item heights (real or estimated). This tells the ScrollView how much
// scrollable area exists.
func (vc *virtualContent) Layout(_ widget.Context, c geometry.Constraints) geometry.Size {
	if vc.list == nil {
		return geometry.Size{}
	}

	totalHeight := vc.list.heights.totalHeight()

	// Width fills the available space.
	width := c.MaxWidth
	if width >= geometry.Infinity {
		width = c.MinWidth
	}
	if width < c.MinWidth {
		width = c.MinWidth
	}

	return geometry.Sz(width, totalHeight)
}

// Draw renders only the visible items within the current viewport.
func (vc *virtualContent) Draw(ctx widget.Context, canvas widget.Canvas) {
	if vc.list == nil {
		return
	}

	lv := vc.list
	itemCount := lv.cfg.ResolvedItemCount()

	// Handle empty state.
	if itemCount == 0 {
		lv.painter.PaintEmptyState(canvas, lv.viewportBounds())
		return
	}

	scrollY := lv.currentScrollY()
	viewportH := lv.viewportHeight

	start, end := lv.heights.visibleRange(scrollY, viewportH, lv.cfg.overscan)
	selectedIdx := lv.cfg.ResolvedSelectedIndex()

	// Update the widget cache for the visible range.
	lv.cache.update(start, end, lv.cfg.itemContent, selectedIdx, lv.hoveredIndex)

	// Wire parent chain on item widgets so dirty propagation
	// (SetNeedsRedraw  - propagateDirtyUpward) can reach the root WidgetBase
	// boundary. Flutter adoptChild pattern.
	for i := 0; i < end-start; i++ {
		if w := lv.cache.widgetAt(i); w != nil {
			if setter, ok := w.(interface{ SetParent(widget.Widget) }); ok {
				setter.SetParent(vc)
			}
		}
	}

	// Content width excludes scrollbar inset so items don't render under it.
	contentWidth := lv.viewportWidth - lv.scroll.ScrollbarInset()

	// Layout and draw each visible item.
	for i := start; i < end; i++ {
		offset := i - start
		w := lv.cache.widgetAt(offset)
		if w == nil {
			continue
		}

		y := lv.heights.offsetAt(i)

		// Layout the item widget with fixed width, flexible height.
		itemConstraints := geometry.Constraints{
			MinWidth:  contentWidth,
			MaxWidth:  contentWidth,
			MinHeight: 0,
			MaxHeight: geometry.Infinity,
		}
		actualSize := w.Layout(ctx, itemConstraints)

		// Update measured height in lazy mode.
		lv.heights.setMeasured(i, actualSize.Height)

		// Compute item bounds using actual measured height.
		itemBounds := geometry.NewRect(0, y, contentWidth, actualSize.Height)

		// Set widget bounds.
		if setter, ok := w.(interface{ SetBounds(geometry.Rect) }); ok {
			setter.SetBounds(itemBounds)
		}

		// Paint item background (hover).
		ips := ItemPaintState{
			Bounds:   itemBounds,
			Index:    i,
			Selected: i == selectedIdx,
			Focused:  i == selectedIdx && lv.IsFocused(),
			Hovered:  i == lv.hoveredIndex,
			Disabled: lv.cfg.ResolvedDisabled(),
		}

		lv.painter.PaintItemBackground(canvas, ips)

		// Paint selection highlight.
		if ips.Selected {
			lv.painter.PaintSelection(canvas, ips)
		}

		widget.StampScreenOrigin(w, canvas)
		widget.DrawChild(w, ctx, canvas)

		// Draw divider between items (not after the last visible item).
		if lv.cfg.divider && i < end-1 {
			divY := y + actualSize.Height
			lv.painter.PaintDivider(canvas, DividerState{
				Bounds:    geometry.NewRect(0, divY, contentWidth, dividerHeight),
				ItemIndex: i,
			})
		}
	}

	// Check end-reached callback.
	lv.checkEndReached(end, itemCount)

	// Clear dirty  - individual items track their own dirty state.
	// Without this, virtualContent (bounds=full content height) stays
	// permanently dirty, causing huge dirty regions in the overlay.
	vc.ClearRedraw()
}

// Event delegates events back to the parent list for item interaction.
func (vc *virtualContent) Event(ctx widget.Context, e event.Event) bool {
	if vc.list == nil {
		return false
	}
	return handleContentEvent(vc.list, ctx, e)
}

// Children returns the cached item widgets for dirty-region collection.
// Their ScreenBounds (set during the previous Draw) allow the dirty.Collector to
// report item-level dirty rects clipped to the viewport.
func (vc *virtualContent) Children() []widget.Widget {
	if vc.list == nil {
		return nil
	}
	widgets := vc.list.cache.widgets
	if len(widgets) == 0 {
		return nil
	}
	children := make([]widget.Widget, 0, len(widgets))
	for _, w := range widgets {
		if w != nil {
			children = append(children, w)
		}
	}
	return children
}
