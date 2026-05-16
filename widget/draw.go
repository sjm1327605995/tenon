package widget

// DrawStats holds statistics from a draw tree traversal.
type DrawStats struct {
	TotalWidgets   int
	DrawnWidgets   int
	SkippedWidgets int
	DirtyWidgets   int
	CleanWidgets   int
}

// DrawStatsProvider is an optional interface implemented by Context
// implementations that support draw statistics collection.
type DrawStatsProvider interface {
	DrawStats() *DrawStats
}

// DrawTree performs a draw traversal of the widget tree rooted at w.
func DrawTree(w Widget, ctx Context, canvas Canvas) DrawStats {
	var stats DrawStats

	type drawStatsSetter interface {
		SetDrawStats(*DrawStats)
	}
	if setter, ok := ctx.(drawStatsSetter); ok {
		setter.SetDrawStats(&stats)
		defer setter.SetDrawStats(nil)
	}

	drawTreeRecursive(w, ctx, canvas, &stats)
	return stats
}

// DrawChild draws a child widget.
func DrawChild(child Widget, ctx Context, canvas Canvas) {
	if child == nil {
		return
	}
	child.Draw(ctx, canvas)
}

func drawTreeRecursive(w Widget, ctx Context, canvas Canvas, stats *DrawStats) {
	if w == nil {
		return
	}

	stats.TotalWidgets++

	type redrawChecker interface {
		NeedsRedraw() bool
	}
	if rc, ok := w.(redrawChecker); ok {
		if rc.NeedsRedraw() {
			stats.DirtyWidgets++
		} else {
			stats.CleanWidgets++
		}
	} else {
		stats.DirtyWidgets++
	}

	StampScreenOrigin(w, canvas)
	w.Draw(ctx, canvas)
	stats.DrawnWidgets++
}

// CollectDrawStats walks the widget tree and counts dirty/clean widgets
// WITHOUT drawing anything.
func CollectDrawStats(w Widget) DrawStats {
	var stats DrawStats
	collectStatsRecursive(w, &stats)
	return stats
}

func collectStatsRecursive(w Widget, stats *DrawStats) {
	if w == nil {
		return
	}

	stats.TotalWidgets++

	type visibilityChecker interface {
		IsVisible() bool
	}
	if vc, ok := w.(visibilityChecker); ok && !vc.IsVisible() {
		stats.SkippedWidgets++
		return
	}

	type redrawChecker interface {
		NeedsRedraw() bool
	}
	if rc, ok := w.(redrawChecker); ok {
		if rc.NeedsRedraw() {
			stats.DirtyWidgets++
		} else {
			stats.CleanWidgets++
		}
	} else {
		stats.DirtyWidgets++
	}

	for _, child := range w.Children() {
		collectStatsRecursive(child, stats)
	}
}
