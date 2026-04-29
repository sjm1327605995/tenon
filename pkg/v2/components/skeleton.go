package components

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/v2/core"
)

// Skeleton is a loading placeholder that mimics the shape of content.
type Skeleton struct {
	core.BaseElement
}

// NewSkeleton creates a skeleton placeholder.
func NewSkeleton() *Skeleton {
	s := &Skeleton{}
	s.Init(s)
	s.SetWidth(100)
	s.SetHeight(20)
	return s
}

// ElementType returns type identifier.
func (s *Skeleton) ElementType() string { return "Skeleton" }

// Draw renders the skeleton placeholder.
func (s *Skeleton) Draw(screen *ebiten.Image) {
	bounds := s.GetBounds()
	theme := core.GetTheme()
	br := core.BorderRadius{TopLeft: theme.BorderRadius / 2, TopRight: theme.BorderRadius / 2, BottomRight: theme.BorderRadius / 2, BottomLeft: theme.BorderRadius / 2}
	drawRoundedRectFill(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, br, theme.MutedColor)
}
