package ui

import (
	"time"

	"gioui.org/op"
	gpaint "gioui.org/op/paint"
)

type gioFrameCache struct {
	version uint64
	call    op.CallOp
	valid   bool
	builds  uint64
}

func (g *game) framePaintVersion(now time.Time) uint64 {
	version := g.visualVersion
	if focused := g.focusedFiber; focused != nil && focused.rnode != nil && focused.rnode.kind == rnInput {
		phase := uint64(now.UnixMilli() / int64(caretBlinkPeriod/time.Millisecond))
		version = version<<1 | phase&1
	}
	return version
}

func (c *gioFrameCache) Add(g *game, ops *op.Ops, version uint64) {
	if c.valid && c.version == version {
		c.call.Add(ops)
		return
	}

	cachedOps := new(op.Ops)
	macro := op.Record(cachedOps)
	gpaint.ColorOp{Color: nrgba(Color{247, 248, 250, 255})}.Add(cachedOps)
	gpaint.PaintOp{}.Add(cachedOps)
	painter := newGioPainter(cachedOps, g.w, g.h)
	g.inputTags = painter.accessibilityTags
	if g.rootRN != nil {
		paint(painter, g.rootRN)
	}
	for _, portal := range g.portals {
		if portal.overlayRoot != nil {
			paint(painter, portal.overlayRoot)
		}
	}
	c.call = macro.Stop()
	c.version = version
	c.valid = true
	c.builds++
	c.call.Add(ops)
}
