package components

import (
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// ScrollView is a scrollable container with vertical scrolling support.
type ScrollView struct {
	core.BaseElement
	content        *View
	scrollX        float32
	scrollY        float32
	maxScrollX     float32
	maxScrollY     float32
	scrollbarWidth float32
	scrollbarColor color.Color
	trackColor     color.Color
	backgroundColor color.Color
	dragging       bool
	lastMouseX     float32
	lastMouseY     float32
	velocityY      float32
	flinging       bool
	lastDragY      float32
	lastDragTime   time.Time
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

	// Inner content view
	sv.content = NewView()
	sv.content.SetWidthPercent(100)
	sv.BaseElement.AppendChild(sv.content)

	return sv
}

// ElementType returns type identifier.
func (sv *ScrollView) ElementType() string { return "ScrollView" }

// Content returns the inner content container.
func (sv *ScrollView) Content() *View { return sv.content }

// Draw renders the background, scrollbar track and thumb.
func (sv *ScrollView) Draw(screen *ebiten.Image) {
	if !sv.IsVisible() {
		return
	}
	bounds := sv.GetBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

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
}

func (sv *ScrollView) drawVerticalScrollbar(screen *ebiten.Image, bounds core.LayoutBounds) {
	trackX := bounds.X + bounds.Width - sv.scrollbarWidth - 2
	trackY := bounds.Y + 2
	trackH := bounds.Height - 4

	// Track
	vector.FillRect(screen, trackX, trackY, sv.scrollbarWidth, trackH, sv.trackColor, false)

	// Thumb
	contentH := sv.content.GetBounds().Height
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

// HandleEvent processes wheel and drag events.
func (sv *ScrollView) HandleEvent(e *core.Event) bool {
	bounds := sv.GetBounds()

	switch e.Type {
	case core.EventScroll:
		// Cap wheel delta to prevent touchpad/large deltas from jumping too far.
		deltaY := e.DeltaY
		if deltaY > 1 {
			deltaY = 1
		} else if deltaY < -1 {
			deltaY = -1
		}
		sv.scrollY -= deltaY * 20
		// Prevent negative scroll immediately; upper bound clamped in Update().
		if sv.scrollY < 0 {
			sv.scrollY = 0
		}
		sv.applyScrollToContent()
		sv.Mark(core.FlagNeedDraw)
		core.LogDebug("[ScrollView] EventScroll", "rawDeltaY", e.DeltaY, "cappedDeltaY", deltaY, "scrollY", sv.scrollY, "maxScrollY", sv.maxScrollY)
		return true

	case core.EventMouseDown:
		sv.flinging = false
		sv.velocityY = 0
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
		if e.X >= bounds.X && e.X <= bounds.X+bounds.Width &&
			e.Y >= bounds.Y && e.Y <= bounds.Y+bounds.Height {
			sv.dragging = true
			sv.lastMouseX = e.X
			sv.lastMouseY = e.Y
			sv.lastDragY = e.Y
			sv.lastDragTime = time.Now()
			return true
		}

	case core.EventMouseUp:
		if sv.dragging {
			sv.dragging = false
			if math.Abs(float64(sv.velocityY)) > 50 {
				sv.flinging = true
			}
		}

	case core.EventMouseMove:
		if sv.dragging {
			dy := e.Y - sv.lastMouseY
			sv.scrollY += dy
			// Prevent negative scroll immediately; upper bound clamped in Update().
			if sv.scrollY < 0 {
				sv.scrollY = 0
			}
			sv.applyScrollToContent()
			now := time.Now()
			dt := float32(now.Sub(sv.lastDragTime).Seconds())
			if dt > 0 {
				sv.velocityY = (e.Y - sv.lastDragY) / dt
			}
			sv.lastDragY = e.Y
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
	contentBounds := sv.content.GetBounds()
	targetY := svBounds.Y - sv.scrollY
	dy := targetY - contentBounds.Y
	if dy == 0 {
		return
	}
	sv.applyOffsetRecursive(sv.content, 0, dy)
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

// Update recalculates max scroll and handles fling.
func (sv *ScrollView) Update() error {
	contentH := sv.content.GetBounds().Height
	viewportH := sv.GetBounds().Height
	if contentH > viewportH {
		sv.maxScrollY = contentH - viewportH
		// Reserve right margin for scrollbar so content doesn't cover it.
		if y := sv.content.GetYoga(); y != nil {
			y.StyleSetMargin(yoga.EdgeRight, sv.scrollbarWidth+4)
		}
	} else {
		sv.maxScrollY = 0
		sv.scrollY = 0
		// No scrollbar needed, remove right margin.
		if y := sv.content.GetYoga(); y != nil {
			y.StyleSetMargin(yoga.EdgeRight, 0)
		}
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

	sv.clampScroll()
	sv.applyScrollToContent()
	return nil
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
