package components

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/pkg/hooks"
)

type HooksHost struct {
	core.BaseComponent
}

func NewHooksHost() *HooksHost {
	h := &HooksHost{
		BaseComponent: core.NewBaseComponent(),
	}
	h.Init(h)
	return h
}

func (h *HooksHost) UseHooks(renderFunc func()) {
	h.SetHooksRenderFunc(renderFunc)
}

func (h *HooksHost) Update() error {
	if h.GetHooksRenderFunc() != nil && h.IsDirty() {
		hooks.WithHooks(h, h.GetHooksRenderFunc())
		h.ClearDirty()
	}
	return h.BaseComponent.Update()
}

func (h *HooksHost) Draw(screen *ebiten.Image) {
}

func (h *HooksHost) DrawOverlay(screen *ebiten.Image) {
}

func (h *HooksHost) HandleInput() bool {
	return false
}
