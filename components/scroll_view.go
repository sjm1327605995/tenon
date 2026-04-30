package components

import (
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/internal/core"
	"github.com/sjm1327605995/tenon/internal/native"
	"github.com/sjm1327605995/tenon/yoga"
)

// ScrollView is a scrollable container with vertical and horizontal scrolling support.
type ScrollView struct {
	core.BaseElement
	contentView *native.View
	scrollX         float32
	scrollY         float32
	maxScrollX      float32
	maxScrollY      float32
	scrollbarWidth  float32
	scrollbarColor  color.Color
	trackColor      color.Color
	backgroundColor color.Color
	dragging        bool
	lastMouseX      float32
	lastMouseY      float32
	velocityY       float32
	velocityX       float32
	flinging        bool
	flingingX       bool
	lastDragY       float32
	lastDragX       float32
	lastDragTime    time.Time
	panStartX       float32
	panStartY       float32
	panning         bool
}

// NewScrollView creates a scrollable container.
func NewScrollView() *ScrollView {
	theme := core.GetTheme()
	sv := &ScrollView{
		scrollbarWidth: 6,
		scrollbarColor: theme.ScrollbarColor,
		trackColor:     theme.ScrollbarTrackColor,
	}
	sv.Init(sv)
	sv.SetFlag(core.FlagFocusable | core.FlagClipChildren)
	sv.SetFlexDirection(yoga.FlexDirectionColumn)
	sv.SetOverflow(yoga.OverflowHidden)

	sv.contentView = native.NewView()
	sv.contentView.SetWidthPercent(100)
	sv.contentView.SetFlexShrink(0)
	sv.BaseElement.AppendChild(sv.contentView)

	return sv
}

// ElementType returns type identifier.
func (sv *ScrollView) ElementType() string { return "ScrollView" }

// Content returns the inner Content container.
func (sv *ScrollView) Content() *native.View { return sv.contentView }

// Draw renders the background, scrollbar track and thumb.
func (sv *ScrollView) Draw(screen *ebiten.Image) {
	bounds := sv.GetBounds()

	// Background - align to integer pixels to avoid gaps with SubImage clipping
	if sv.backgroundColor != nil {
		x := float32(int(bounds.X))
		y := float32(int(bounds.Y))
		w := bounds.X + bounds.Width - x
		h := bounds.Y + bounds.Height - y
		vector.FillRect(screen, x, y, w, h, sv.backgroundColor, false)
	}

	// Vertical scrollbar
	if sv.maxScrollY > 0 {
		sv.drawVerticalScrollbar(screen, bounds)
	}
	// Horizontal scrollbar
	if sv.maxScrollX > 0 {
		sv.drawHorizontalScrollbar(screen, bounds)
	}
}

func (sv *ScrollView) drawVerticalScrollbar(screen *ebiten.Image, bounds core.LayoutBounds) {
	trackX := bounds.X + bounds.Width - sv.scrollbarWidth - 2
	trackY := bounds.Y + 2
	trackH := bounds.Height - 4

	// Track
	vector.FillRect(screen, trackX, trackY, sv.scrollbarWidth, trackH, sv.trackColor, false)

	// Thumb
	contentH := sv.contentView.GetBounds().Height
	viewportH := bounds.Height
	if contentH > viewportH {
		ratio := viewportH / contentH
		thumbH := ratio * trackH
		if thumbH < 10 {
			thumbH = 10
		}
		scrollRatio := sv.scrollY / sv.maxScrollY
		thumbY := trackY + scrollRatio*(trackH-thumbH)
		vector.FillRect(screen, trackX, thumbY, sv.scrollbarWidth, thumbH, sv.scrollbarColor, false)
	}
}

// recalcMaxScroll computes max scroll based on current bounds.
func (sv *ScrollView) recalcMaxScroll() {
	contentH := sv.contentView.GetBounds().Height
	viewportH := sv.GetBounds().Height
	if contentH > viewportH {
		sv.maxScrollY = contentH - viewportH
	} else {
		sv.maxScrollY = 0
		sv.scrollY = 0
	}
	contentW := sv.contentView.GetBounds().Width
	viewportW := sv.GetBounds().Width
	if contentW > viewportW {
		sv.maxScrollX = contentW - viewportW
	} else {
		sv.maxScrollX = 0
		sv.scrollX = 0
	}
}

// HandleEvent processes wheel and drag events.
func (sv *ScrollView) HandleEvent(e *core.Event) bool {
	bounds := sv.GetBounds()

	switch e.Type {
	case core.EventScroll:
		// Ensure max scroll is computed before handling (bounds may already be valid from layout).
		sv.recalcMaxScroll()

		// Cap wheel delta to prevent touchpad/large deltas from jumping too far.
		deltaY := e.DeltaY
		if deltaY > 1 {
			deltaY = 1
		} else if deltaY < -1 {
			deltaY = -1
		}
		deltaX := e.DeltaX
		if deltaX > 1 {
			deltaX = 1
		} else if deltaX < -1 {
			deltaX = -1
		}
		if sv.maxScrollY > 0 && deltaY != 0 {
			sv.scrollY -= deltaY * 20
			if sv.scrollY < 0 {
				sv.scrollY = 0
			}
			sv.applyScrollToContent()
			sv.Mark(core.FlagNeedDraw)
		}
		if sv.maxScrollX > 0 && deltaX != 0 {
			sv.scrollX -= deltaX * 20
			if sv.scrollX < 0 {
				sv.scrollX = 0
			}
			sv.applyScrollToContent()
			sv.Mark(core.FlagNeedDraw)
		}
		core.LogDebug("[ScrollView] EventScroll", "rawDeltaY", e.DeltaY, "cappedDeltaY", deltaY, "scrollY", sv.scrollY, "maxScrollY", sv.maxScrollY, "scrollX", sv.scrollX, "maxScrollX", sv.maxScrollX)
		return true

	case core.EventMouseDown:
		sv.flinging = false
		sv.velocityY = 0
		sv.velocityX = 0
		// Only consume if clicking on scrollbar (so child elements in Content area still get clicks)
		if sv.maxScrollY > 0 {
			trackX := bounds.X + bounds.Width - sv.scrollbarWidth - 2
			trackY := bounds.Y + 2
			trackH := bounds.Height - 4
			if e.X >= trackX && e.X <= trackX+sv.scrollbarWidth+4 &&
				e.Y >= trackY && e.Y <= trackY+trackH {
				sv.dragging = true
				sv.lastMouseX = e.X
				sv.lastMouseY = e.Y
				sv.lastDragY = e.Y
				sv.lastDragTime = time.Now()
				return true
			}
		}
		if sv.maxScrollX > 0 {
			trackY := bounds.Y + bounds.Height - sv.scrollbarWidth - 2
			trackX := bounds.X + 2
			trackW := bounds.Width - 4
			if e.Y >= trackY && e.Y <= trackY+sv.scrollbarWidth+4 &&
				e.X >= trackX && e.X <= trackX+trackW {
				sv.dragging = true
				sv.lastMouseX = e.X
				sv.lastMouseY = e.Y
				sv.lastDragX = e.X
				sv.lastDragTime = time.Now()
				return true
			}
		}
		// Content area: do NOT consume, let children handle their own clicks.
		// Pan-scrolling is handled in Update() instead.

	case core.EventMouseUp:
		if sv.dragging {
			sv.dragging = false
			if math.Abs(float64(sv.velocityY)) > 50 {
				sv.flinging = true
			}
			if math.Abs(float64(sv.velocityX)) > 50 {
				sv.flingingX = true
			}
		}

	case core.EventMouseMove:
		if sv.dragging {
			dy := e.Y - sv.lastMouseY
			dx := e.X - sv.lastMouseX
			if sv.maxScrollY > 0 {
				sv.scrollY += dy
				if sv.scrollY < 0 {
					sv.scrollY = 0
				}
				sv.applyScrollToContent()
			}
			if sv.maxScrollX > 0 {
				sv.scrollX += dx
				if sv.scrollX < 0 {
					sv.scrollX = 0
				}
				sv.applyScrollToContent()
			}
			now := time.Now()
			dt := float32(now.Sub(sv.lastDragTime).Seconds())
			if dt > 0 {
				if sv.maxScrollY > 0 {
					sv.velocityY = (e.Y - sv.lastDragY) / dt
				}
				if sv.maxScrollX > 0 {
					sv.velocityX = (e.X - sv.lastDragX) / dt
				}
			}
			sv.lastDragY = e.Y
			sv.lastDragX = e.X
			sv.lastDragTime = now
			sv.lastMouseX = e.X
			sv.lastMouseY = e.Y
			sv.Mark(core.FlagNeedLayout | core.FlagNeedDraw)
			return true
		}
	}

	return false
}

func (sv *ScrollView) clampScroll() {
	if sv.scrollY < 0 {
		sv.scrollY = 0
	}
	if sv.scrollY > sv.maxScrollY {
		sv.scrollY = sv.maxScrollY
	}
	if sv.scrollX < 0 {
		sv.scrollX = 0
	}
	if sv.scrollX > sv.maxScrollX {
		sv.scrollX = sv.maxScrollX
	}
}

func (sv *ScrollView) applyScrollToContent() {
	svBounds := sv.GetBounds()
	contentBounds := sv.contentView.GetBounds()
	targetY := svBounds.Y - sv.scrollY
	targetX := svBounds.X - sv.scrollX
	dy := targetY - contentBounds.Y
	dx := targetX - contentBounds.X
	if dx == 0 && dy == 0 {
		return
	}
	sv.applyOffsetRecursive(sv.contentView, dx, dy)
}

func (sv *ScrollView) applyOffsetRecursive(el core.Element, dx, dy float32) {
	if el == nil {
		return
	}
	b := el.GetBounds()
	b.X += dx
	b.Y += dy
	el.SetBounds(b)
	for _, child := range el.GetChildren() {
		sv.applyOffsetRecursive(child, dx, dy)
	}
}

// Update recalculates max scroll, handles fling, and pan-gesture scrolling.
func (sv *ScrollView) Update() error {
	contentH := sv.contentView.GetBounds().Height
	viewportH := sv.GetBounds().Height
	contentW := sv.contentView.GetBounds().Width
	viewportW := sv.GetBounds().Width

	if contentH > viewportH {
		sv.maxScrollY = contentH - viewportH
		if y := sv.contentView.GetYoga(); y != nil {
			y.StyleSetMargin(yoga.EdgeRight, sv.scrollbarWidth+4)
		}
	} else {
		sv.maxScrollY = 0
		sv.scrollY = 0
		if y := sv.contentView.GetYoga(); y != nil {
			y.StyleSetMargin(yoga.EdgeRight, 0)
		}
	}

	if contentW > viewportW {
		sv.maxScrollX = contentW - viewportW
		if y := sv.contentView.GetYoga(); y != nil {
			y.StyleSetMargin(yoga.EdgeBottom, sv.scrollbarWidth+4)
		}
	} else {
		sv.maxScrollX = 0
		sv.scrollX = 0
		if y := sv.contentView.GetYoga(); y != nil {
			y.StyleSetMargin(yoga.EdgeBottom, 0)
		}
	}

	// Pan gesture: if mouse is down inside scrollview bounds, track movement for scrolling.
	if core.IsMouseButtonPressed(core.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		bounds := sv.GetBounds()
		inBounds := float32(mx) >= bounds.X && float32(mx) < bounds.X+bounds.Width &&
			float32(my) >= bounds.Y && float32(my) < bounds.Y+bounds.Height
		if inBounds && !sv.dragging {
			if !sv.panning {
				sv.panning = true
				sv.panStartX = float32(mx)
				sv.panStartY = float32(my)
				sv.lastMouseX = float32(mx)
				sv.lastMouseY = float32(my)
			} else {
				dx := float32(mx) - sv.lastMouseX
				dy := float32(my) - sv.lastMouseY
				moved := math.Abs(float64(dx)) > 0.5 || math.Abs(float64(dy)) > 0.5
				if moved {
					if sv.maxScrollY > 0 {
						sv.scrollY -= dy
					}
					if sv.maxScrollX > 0 {
						sv.scrollX -= dx
					}
					sv.applyScrollToContent()
					sv.Mark(core.FlagNeedDraw)
				}
				sv.lastMouseX = float32(mx)
				sv.lastMouseY = float32(my)
			}
		}
	} else {
		sv.panning = false
	}

	if sv.flinging {
		sv.scrollY += sv.velocityY * (1.0 / 60.0)
		sv.velocityY *= 0.9
		if math.Abs(float64(sv.velocityY)) < 10 {
			sv.flinging = false
			sv.velocityY = 0
		}
		sv.Mark(core.FlagNeedDraw)
	}

	if sv.flingingX {
		sv.scrollX += sv.velocityX * (1.0 / 60.0)
		sv.velocityX *= 0.9
		if math.Abs(float64(sv.velocityX)) < 10 {
			sv.flingingX = false
			sv.velocityX = 0
		}
		sv.Mark(core.FlagNeedDraw)
	}

	sv.clampScroll()
	sv.applyScrollToContent()
	return nil
}

func (sv *ScrollView) drawHorizontalScrollbar(screen *ebiten.Image, bounds core.LayoutBounds) {
	trackY := bounds.Y + bounds.Height - sv.scrollbarWidth - 2
	trackX := bounds.X + 2
	trackW := bounds.Width - 4

	// Track
	vector.FillRect(screen, trackX, trackY, trackW, sv.scrollbarWidth, sv.trackColor, false)

	// Thumb
	contentW := sv.contentView.GetBounds().Width
	viewportW := bounds.Width
	if contentW > viewportW {
		ratio := viewportW / contentW
		thumbW := ratio * trackW
		if thumbW < 10 {
			thumbW = 10
		}
		scrollRatio := sv.scrollX / sv.maxScrollX
		thumbX := trackX + scrollRatio*(trackW-thumbW)
		vector.FillRect(screen, thumbX, trackY, thumbW, sv.scrollbarWidth, sv.scrollbarColor, false)
	}
}

// SetScrollbarWidth sets the scrollbar width.
func (sv *ScrollView) SetScrollbarWidth(w float32) *ScrollView {
	sv.scrollbarWidth = w
	sv.Mark(core.FlagNeedDraw)
	return sv
}

// SetScrollbarColor sets the thumb color.
func (sv *ScrollView) SetScrollbarColor(clr color.Color) *ScrollView {
	sv.scrollbarColor = clr
	sv.Mark(core.FlagNeedDraw)
	return sv
}

// SetBackgroundColor sets the scrollview background color.
func (sv *ScrollView) SetBackgroundColor(clr color.Color) *ScrollView {
	sv.backgroundColor = clr
	sv.Mark(core.FlagNeedDraw)
	return sv
}

// SetTrackColor sets the track background color.
func (sv *ScrollView) SetTrackColor(clr color.Color) *ScrollView {
	sv.trackColor = clr
	sv.Mark(core.FlagNeedDraw)
	return sv
}

// SyncFrom 同步新 ScrollView 的属性到当前 Element（声明式重建）。
func (sv *ScrollView) SyncFrom(src core.Element) {
	other, ok := src.(*ScrollView)
	if !ok {
		return
	}
	sb := &core.SyncBuilder{}
	core.SyncField(sb, &sv.scrollX, other.scrollX)
	core.SyncField(sb, &sv.scrollY, other.scrollY)
	core.SyncField(sb, &sv.scrollbarWidth, other.scrollbarWidth)
	core.SyncColor(sb, &sv.scrollbarColor, other.scrollbarColor)
	core.SyncColor(sb, &sv.trackColor, other.trackColor)
	core.SyncColor(sb, &sv.backgroundColor, other.backgroundColor)
	sb.MarkDraw(sv)
}

