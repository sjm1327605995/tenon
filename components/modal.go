package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/internal/core"
	"github.com/sjm1327605995/tenon/internal/native"
	"github.com/sjm1327605995/tenon/yoga"
)

// Modal is a dialog overlay with a mask backdrop and a Content panel.
type Modal struct {
	core.BaseElement
	panel       *native.View
	titleEl     *native.Text
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
		maskColor:   color.RGBA{R: 0, G: 0, B: 0, A: 120},
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

	m.panel = native.NewView()
	m.panel.SetWidth(400)
	m.panel.SetMinHeight(200)
	m.panel.SetBackgroundColor(theme.CardColor)
	m.panel.SetBorderRadius(theme.BorderRadius)
	m.panel.SetShadow(theme.ShadowColor, 16, 0, 4)
	m.panel.SetPadding(yoga.EdgeAll, 24)
	m.panel.SetFlexDirection(yoga.FlexDirectionColumn)
	m.AppendChild(m.panel)

	// Pre-create title element so user Content is always appended after it
	m.titleEl = native.NewText("")
	m.titleEl.SetMargin(yoga.EdgeBottom, 16)
	m.titleEl.SetVisible(false)
	m.panel.AppendChild(m.titleEl)

	return m
}

// ElementType returns type identifier.
func (m *Modal) ElementType() string { return "Modal" }

// Draw renders the backdrop mask.
func (m *Modal) Draw(screen *ebiten.Image) {
	bounds := m.GetBounds()
	//存在painc情况
	vector.FillRect(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, m.maskColor, false)
}

// Update handles ESC key to close.
func (m *Modal) Update() error {
	if m.IsVisible() && m.closeOnEsc && core.IsKeyJustPressed(core.KeyEscape) {
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
	if eng := m.GetEngine(); eng != nil {
		eng.AddOverlay(m)
	}
	return m
}

// Close hides the modal.
func (m *Modal) Close() {
	m.SetVisible(false)
	if eng := m.GetEngine(); eng != nil {
		eng.RemoveOverlay(m)
	}
	if m.onClose != nil {
		m.onClose()
	}
}

// Panel returns the inner Content container for adding custom children.
func (m *Modal) Panel() *native.View { return m.panel }

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

// SyncFrom 同步新 Modal 的属性到当前 Element（声明式重建）。
func (m *Modal) SyncFrom(src core.Element) {
	other, ok := src.(*Modal)
	if !ok {
		return
	}
	sb := &core.SyncBuilder{}
	core.SyncField(sb, &m.closeOnMask, other.closeOnMask)
	core.SyncField(sb, &m.closeOnEsc, other.closeOnEsc)
	core.SyncColor(sb, &m.maskColor, other.maskColor)
	sb.MarkDraw(m)
}

// SetMaskColor sets the backdrop color.
func (m *Modal) SetMaskColor(clr color.Color) *Modal {
	m.maskColor = clr
	m.Mark(core.FlagNeedDraw)
	return m
}

// SetPanelSize sets the Content panel size.
func (m *Modal) SetPanelSize(width, height float32) *Modal {
	m.panel.SetWidth(width)
	m.panel.SetHeight(height)
	return m
}

