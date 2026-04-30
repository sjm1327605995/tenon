package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/internal/core"
	"github.com/sjm1327605995/tenon/internal/native"
	"github.com/sjm1327605995/tenon/yoga"
)

// Window is a draggable floating window with a title bar and Content area.
type Window struct {
	native.View
	titleBar     *native.View
	titleEl      *native.Text
	closeBtn     *native.View
	contentPanel *native.View
	dragging     bool
	dragOffsetX  float32
	dragOffsetY  float32
	titleHeight  float32
	onClose      func()
}

// NewWindow creates a draggable window.
func NewWindow(title string, width, height float32) *Window {
	theme := core.GetTheme()
	w := &Window{
		titleHeight: 32,
	}
	w.Init(w)
	w.SetPositionType(yoga.PositionTypeAbsolute)
	w.SetWidth(width)
	w.SetHeight(height)
	w.SetBackgroundColor(theme.SurfaceColor)
	w.SetBorderRadius(theme.BorderRadius)
	w.SetOverflow(yoga.OverflowHidden)

	// Title bar
	w.titleBar = native.NewView()
	w.titleBar.SetHeight(w.titleHeight)
	w.titleBar.SetWidthPercent(100)
	w.titleBar.SetBackgroundColor(theme.SurfaceColor)
	w.titleBar.SetBorder(yoga.EdgeBottom, 1)
	w.titleBar.SetBorderColor(theme.BorderColor)
	w.titleBar.SetFlexDirection(yoga.FlexDirectionRow)
	w.titleBar.SetAlignItems(yoga.AlignCenter)
	w.titleBar.SetPadding(yoga.EdgeHorizontal, 12)
	w.AppendChild(w.titleBar)

	// Title text
	w.titleEl = native.NewText(title).SetColor(theme.TextMutedColor)
	w.titleBar.AppendChild(w.titleEl)

	// Spacer to push close button to right
	spacer := native.NewView()
	spacer.SetFlexGrow(1)
	w.titleBar.AppendChild(spacer)

	// Close button area
	w.closeBtn = native.NewView()
	w.closeBtn.SetWidth(16)
	w.closeBtn.SetHeight(16)
	w.titleBar.AppendChild(w.closeBtn)

	// Content panel
	w.contentPanel = native.NewView()
	w.contentPanel.SetWidthPercent(100)
	w.contentPanel.SetFlexGrow(1)
	w.contentPanel.SetPadding(yoga.EdgeAll, 12)
	w.AppendChild(w.contentPanel)

	return w
}

// ElementType returns type identifier.
func (w *Window) ElementType() string { return "Window" }

// Draw renders the window background and close button X.
func (w *Window) Draw(screen *ebiten.Image) {
	w.View.Draw(screen)

	// Draw close button X mark
	cb := w.closeBtn.GetBounds()
	if cb.Width <= 0 || cb.Height <= 0 {
		return
	}
	cx := cb.X + cb.Width/2
	cy := cb.Y + cb.Height/2
	sz := cb.Width * 0.3
	clr := color.RGBA{R: 150, G: 150, B: 150, A: 255}
	vector.StrokeLine(screen, cx-sz, cy-sz, cx+sz, cy+sz, 1.5, clr, false)
	vector.StrokeLine(screen, cx+sz, cy-sz, cx-sz, cy+sz, 1.5, clr, false)
}

// HandleEvent processes drag start, drag end and close click.
func (w *Window) HandleEvent(e *core.Event) bool {
	switch e.Type {
	case core.EventMouseDown:
		// Start drag if click is on title bar
		tb := w.titleBar.GetBounds()
		if e.X >= tb.X && e.X < tb.X+tb.Width &&
			e.Y >= tb.Y && e.Y < tb.Y+tb.Height {
			w.dragging = true
			w.dragOffsetX = e.X - w.GetBounds().X
			w.dragOffsetY = e.Y - w.GetBounds().Y
			return true
		}
	case core.EventMouseUp:
		if w.dragging {
			w.dragging = false
			return true
		}
	case core.EventClick:
		// Close button click
		cb := w.closeBtn.GetBounds()
		if e.X >= cb.X && e.X < cb.X+cb.Width &&
			e.Y >= cb.Y && e.Y < cb.Y+cb.Height {
			w.Close()
			return true
		}
		// Click on title bar or Content area consumes event to block穿透
		return true
	}
	return false
}

// Update handles dragging movement.
func (w *Window) Update() error {
	if w.dragging {
		mx, my := ebiten.CursorPosition()
		newX := float32(mx) - w.dragOffsetX
		newY := float32(my) - w.dragOffsetY
		w.moveTo(newX, newY)
	}
	return nil
}

func (w *Window) moveTo(x, y float32) {
	oldBounds := w.GetBounds()
	dx := x - oldBounds.X
	dy := y - oldBounds.Y
	if dx == 0 && dy == 0 {
		return
	}

	// Update window bounds immediately for zero-latency drag
	oldBounds.X = x
	oldBounds.Y = y
	w.SetBounds(oldBounds)

	// Sync yoga position so next layout calculation stays consistent
	w.GetYoga().StyleSetPosition(yoga.EdgeLeft, x)
	w.GetYoga().StyleSetPosition(yoga.EdgeTop, y)

	// Recursively shift all child bounds so they follow the window instantly
	w.shiftChildBounds(w, dx, dy)
}

func (w *Window) shiftChildBounds(parent core.Element, dx, dy float32) {
	for _, child := range parent.GetChildren() {
		b := child.GetBounds()
		b.X += dx
		b.Y += dy
		child.SetBounds(b)
		w.shiftChildBounds(child, dx, dy)
	}
}

// Content returns the Content panel for adding custom children.
func (w *Window) Content() *native.View { return w.contentPanel }

// SetTitle updates the window title.
func (w *Window) SetTitle(title string) *Window {
	w.titleEl.SetContent(title)
	return w
}

// SetOnClose sets the close callback.
func (w *Window) SetOnClose(fn func()) *Window {
	w.onClose = fn
	return w
}

// Close hides the window and triggers onClose.
func (w *Window) Close() {
	w.SetVisible(false)
	if w.onClose != nil {
		w.onClose()
	}
}

// Show shows the window.
func (w *Window) Show() *Window {
	w.SetVisible(true)
	return w
}
