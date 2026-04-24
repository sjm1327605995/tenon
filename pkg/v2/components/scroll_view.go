package components

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// ScrollView is a scrollable container.
type ScrollView struct {
	core.BaseElement
	content       *View
	scrollX       float32
	scrollY       float32
	maxScrollX    float32
	maxScrollY    float32
	scrollbarW    float32
	scrollColor   color.Color
	trackColor    color.Color
	showScrollbar bool
}

// NewScrollView creates a scrollable container.
func NewScrollView() *ScrollView {
	sv := &ScrollView{
		scrollbarW:    6,
		scrollColor:   color.RGBA{R: 150, G: 150, B: 150, A: 200},
		trackColor:    color.RGBA{R: 230, G: 230, B: 230, A: 100},
		showScrollbar: true,
	}
	sv.Init(sv)
	sv.BaseElement.SetOverflow(yoga.OverflowHidden)

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

// Draw renders the scrollview with optional scrollbar.
func (sv *ScrollView) Draw(screen *ebiten.Image) {
	if !sv.IsVisible() {
		return
	}
	bounds := sv.GetBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	// Draw track
	if sv.showScrollbar && sv.maxScrollY > 0 {
		trackX := bounds.X + bounds.Width - sv.scrollbarW
		vector.FillRect(screen, trackX, bounds.Y, sv.scrollbarW, bounds.Height, sv.trackColor, false)

		// Draw thumb
		viewHeight := bounds.Height
		contentHeight := viewHeight + sv.maxScrollY
		thumbRatio := viewHeight / contentHeight
		thumbHeight := math.Max(float64(thumbRatio*viewHeight), 20)
		scrollRatio := float64(0)
		if sv.maxScrollY > 0 {
			scrollRatio = float64(sv.scrollY) / float64(sv.maxScrollY)
		}
		thumbY := float64(bounds.Y) + scrollRatio*(float64(viewHeight)-thumbHeight)

		vector.FillRect(screen, trackX, float32(thumbY), sv.scrollbarW, float32(thumbHeight), sv.scrollColor, false)
	}
}

// HandleEvent processes scroll events.
func (sv *ScrollView) HandleEvent(e *core.Event) bool {
	if e.Type == core.EventScroll {
		// Vertical scroll
		sv.scrollY += e.DeltaY * 20
		sv.clampScroll()
		sv.updateContentPosition()
		sv.Mark(core.FlagNeedDraw)
		return true
	}
	return false
}

// Update recalculates max scroll after layout.
func (sv *ScrollView) Update() error {
	bounds := sv.GetBounds()
	contentBounds := sv.content.GetBounds()

	// Calculate max scroll based on content size vs viewport size
	sv.maxScrollY = float32(math.Max(0, float64(contentBounds.Height-bounds.Height)))
	sv.maxScrollX = float32(math.Max(0, float64(contentBounds.Width-bounds.Width)))

	sv.clampScroll()
	sv.updateContentPosition()
	return nil
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

func (sv *ScrollView) updateContentPosition() {
	// Apply scroll offset to content via yoga position
	if y := sv.content.GetYoga(); y != nil {
		y.StyleSetPosition(yoga.EdgeTop, -sv.scrollY)
	}
}

// Chain API

func (sv *ScrollView) SetScrollbarWidth(w float32) *ScrollView {
	sv.scrollbarW = w
	sv.Mark(core.FlagNeedDraw)
	return sv
}

func (sv *ScrollView) SetShowScrollbar(show bool) *ScrollView {
	sv.showScrollbar = show
	sv.Mark(core.FlagNeedDraw)
	return sv
}
