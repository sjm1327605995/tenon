package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/v2/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// Modal is a dialog overlay with a mask backdrop and a content panel.
type Modal struct {
	core.BaseElement
	panel       *View
	titleEl     *Text
	onClose     func()
	closeOnMask bool
	closeOnEsc  bool
	maskColor   color.Color
}

// NewModal creates a modal dialog.
func NewModal() *Modal {
	theme := core.GetTheme()
	m := &Modal{
		closeOnMask: true,
		closeOnEsc:  true,
		maskColor:   color.RGBA{R: 0, G: 0, B: 0, A: 180},
	}
	m.Init(m)
	m.SetVisible(false)
	m.SetPositionType(yoga.PositionTypeAbsolute)
	m.SetPosition(yoga.EdgeLeft, 0)
	m.SetPosition(yoga.EdgeTop, 0)
	m.SetWidthPercent(100)
	m.SetHeightPercent(100)
	m.SetFlexDirection(yoga.FlexDirectionColumn)
	m.SetJustifyContent(yoga.JustifyCenter)
	m.SetAlignItems(yoga.AlignCenter)

	m.panel = NewView()
	m.panel.SetWidth(400)
	m.panel.SetMinHeight(200)
	m.panel.SetBackgroundColor(theme.SurfaceColor)
	m.panel.SetBorderRadius(theme.BorderRadius)
	m.panel.SetPadding(yoga.EdgeAll, 24)
	m.panel.SetFlexDirection(yoga.FlexDirectionColumn)
	m.AppendChild(m.panel)

	// Pre-create title element so user content is always appended after it
	m.titleEl = NewText("")
	m.titleEl.SetMargin(yoga.EdgeBottom, 16)
	m.titleEl.SetVisible(false)
	m.panel.AppendChild(m.titleEl)

	return m
}

// ElementType returns type identifier.
func (m *Modal) ElementType() string { return "Modal" }

// Draw renders the backdrop mask.
func (m *Modal) Draw(screen *ebiten.Image) {
	if !m.IsVisible() {
		return
	}
	bounds := m.GetBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}
	vector.FillRect(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, m.maskColor, false)
}

// Update handles ESC key to close.
func (m *Modal) Update() error {
	if m.IsVisible() && m.closeOnEsc && inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		m.Close()
	}
	return nil
}

// HandleEvent processes mask click to close.
func (m *Modal) HandleEvent(e *core.Event) bool {
	if !m.IsVisible() {
		return false
	}
	switch e.Type {
	case core.EventClick:
		if m.closeOnMask {
			pb := m.panel.GetBounds()
			// Click outside panel bounds = mask click
			if e.X < pb.X || e.X >= pb.X+pb.Width || e.Y < pb.Y || e.Y >= pb.Y+pb.Height {
				m.Close()
				return true
			}
		}
		return true // consume event, prevent穿透
	}
	return false
}

// Open shows the modal.
func (m *Modal) Open() *Modal {
	m.SetVisible(true)
	return m
}

// Close hides the modal.
func (m *Modal) Close() {
	m.SetVisible(false)
	if m.onClose != nil {
		m.onClose()
	}
}

// Panel returns the inner content container for adding custom children.
func (m *Modal) Panel() *View { return m.panel }

// SetTitle sets the modal title text.
func (m *Modal) SetTitle(title string) *Modal {
	m.titleEl.SetContent(title)
	m.titleEl.SetVisible(title != "")
	return m
}

// SetOnClose sets the close callback.
func (m *Modal) SetOnClose(fn func()) *Modal {
	m.onClose = fn
	return m
}

// SetCloseOnMask controls whether clicking the mask closes the modal.
func (m *Modal) SetCloseOnMask(v bool) *Modal {
	m.closeOnMask = v
	return m
}

// SetCloseOnEsc controls whether pressing ESC closes the modal.
func (m *Modal) SetCloseOnEsc(v bool) *Modal {
	m.closeOnEsc = v
	return m
}

// SetMaskColor sets the backdrop color.
func (m *Modal) SetMaskColor(clr color.Color) *Modal {
	m.maskColor = clr
	m.Mark(core.FlagNeedDraw)
	return m
}

// SetPanelSize sets the content panel size.
func (m *Modal) SetPanelSize(width, height float32) *Modal {
	m.panel.SetWidth(width)
	m.panel.SetHeight(height)
	return m
}
