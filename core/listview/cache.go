package listview

import (
	"github.com/sjm1327605995/tenon/cdk"
	"github.com/sjm1327605995/tenon/widget"
)

// widgetCache caches the currently visible item widgets between frames.
//
// When the visible range has not changed and data has not been invalidated,
// the cache returns the same widgets without calling the builder again.
// This avoids unnecessary allocations during static frames (no scroll).
//
// Each item widget is automatically wrapped in a [primitives.RepaintBoundary]
// to enable per-item pixel caching. When the item subtree is clean (no dirty
// widgets), the cached pixels are blitted directly instead of re-rendering.
// Item background, selection, and dividers are painted OUTSIDE the boundary
// by the painter on the main canvas.
type widgetCache struct {
	startIndex    int
	endIndex      int
	selectedIndex int
	hoveredIndex  int
	widgets       []widget.Widget
	valid         bool
}

// rebuildAffected rebuilds only items whose selection or hover state changed.
// Android RecyclerView pattern: only affected ViewHolders are rebound.
func (wc *widgetCache) rebuildAffected(start int, content cdk.Content[ItemContext], selectedIndex, hoveredIndex int) {
	affectedIndices := make(map[int]bool)
	if wc.selectedIndex != selectedIndex {
		affectedIndices[wc.selectedIndex] = true
		affectedIndices[selectedIndex] = true
	}
	if wc.hoveredIndex != hoveredIndex {
		affectedIndices[wc.hoveredIndex] = true
		affectedIndices[hoveredIndex] = true
	}

	for idx := range affectedIndices {
		offset := idx - start
		if offset < 0 || offset >= len(wc.widgets) {
			continue
		}
		w := content.Render(ItemContext{
			Index:    idx,
			Selected: idx == selectedIndex,
			Focused:  idx == selectedIndex,
			Hovered:  idx == hoveredIndex,
		})
		if w != nil {
			if setter, ok := w.(interface{ SetRepaintBoundary(bool) }); ok {
				setter.SetRepaintBoundary(true)
			}
		}
		wc.widgets[offset] = w
	}
}

// fullRebuild recreates all items in the range (scroll or first build).
func (wc *widgetCache) fullRebuild(start, _, count int, content cdk.Content[ItemContext], selectedIndex, hoveredIndex int) {
	if cap(wc.widgets) >= count {
		wc.widgets = wc.widgets[:count]
	} else {
		wc.widgets = make([]widget.Widget, count)
	}

	if content == nil { //nolint:nestif // cache miss path with lazy initialization and RepaintBoundary wrapping
		for i := range wc.widgets {
			wc.widgets[i] = nil
		}
	} else {
		for i := range count {
			idx := start + i
			w := content.Render(ItemContext{
				Index:    idx,
				Selected: idx == selectedIndex,
				Focused:  idx == selectedIndex,
				Hovered:  idx == hoveredIndex,
			})
			if w != nil {
				if setter, ok := w.(interface{ SetRepaintBoundary(bool) }); ok {
					setter.SetRepaintBoundary(true)
				}
			}
			wc.widgets[i] = w
		}
	}
}

// update ensures the cache contains widgets for the range [start, end).
// If the range matches and the cache is valid, this is a no-op.
// Otherwise, it calls the content's Render method for each index in the range
// and wraps each widget in a RepaintBoundary.
func (wc *widgetCache) update(start, end int, content cdk.Content[ItemContext], selectedIndex, hoveredIndex int) {
	count := end - start
	if count <= 0 {
		wc.clear()
		return
	}

	// Fast path: same range, only selection/hover changed  - rebuild only affected items.
	// Android RecyclerView pattern: notifyItemChanged(pos) rebinds single ViewHolder.
	if wc.valid && wc.startIndex == start && wc.endIndex == end && content != nil {
		if wc.selectedIndex != selectedIndex || wc.hoveredIndex != hoveredIndex {
			wc.rebuildAffected(start, content, selectedIndex, hoveredIndex)
			wc.selectedIndex = selectedIndex
			wc.hoveredIndex = hoveredIndex
			return
		}
		return // nothing changed
	}

	// Full rebuild: range changed or first build.
	wc.fullRebuild(start, end, count, content, selectedIndex, hoveredIndex)
	wc.startIndex = start
	wc.endIndex = end
	wc.selectedIndex = selectedIndex
	wc.hoveredIndex = hoveredIndex
	wc.valid = true
}

// widgetAt returns the cached widget at the given offset from startIndex.
// This returns the raw (unwrapped) widget for compatibility with existing code
// that needs direct widget access (e.g., for Layout and bounds setting).
func (wc *widgetCache) widgetAt(offset int) widget.Widget {
	if offset < 0 || offset >= len(wc.widgets) {
		return nil
	}
	return wc.widgets[offset]
}

// invalidate marks the cache as needing a rebuild.
func (wc *widgetCache) invalidate() {
	wc.valid = false
}

// clear resets the cache entirely and unmounts boundaries to free pixel caches.
func (wc *widgetCache) clear() {
	for i := range wc.widgets {
		wc.widgets[i] = nil
	}
	wc.widgets = wc.widgets[:0]
	wc.startIndex = 0
	wc.endIndex = 0
	wc.valid = false
}
