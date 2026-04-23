package antdesign

import (
	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/components"
	"github.com/sjm1327605995/tenon/pkg/fonts"
	"github.com/sjm1327605995/tenon/yoga"
)

// AntAvatarShape defines avatar shape.
type AntAvatarShape string

const (
	AntAvatarCircle AntAvatarShape = "circle"
	AntAvatarSquare AntAvatarShape = "square"
)

// AntAvatar is an avatar component.
type AntAvatar struct {
	tenon.BaseWidget
	size     float32
	shape    AntAvatarShape
	src      string
	icon     string
	alt      string
	gap      int
	children string // text content if no src/icon
}

// NewAntAvatar creates an AntAvatar.
func NewAntAvatar() *AntAvatar {
	a := &AntAvatar{size: 40, shape: AntAvatarCircle}
	a.Init(a)
	return a
}

// Render returns the avatar UI.
func (a *AntAvatar) Render() tenon.Component {
	theme := NewAntTheme()

	root := components.NewView().
		SetWidth(a.size).
		SetHeight(a.size).
		SetBackgroundColor(theme.BorderColor).
		SetOverflow(yoga.OverflowHidden)

	if a.shape == AntAvatarCircle {
		root.SetBorderRadius(a.size / 2)
	} else {
		root.SetBorderRadius(4)
	}

	if a.src != "" {
		root.AddChild(components.NewImage().
			SetSource(a.src).
			SetWidth(a.size).
			SetHeight(a.size))
	} else if a.icon != "" {
		root.Add(components.NewText(a.icon).
			SetFontSize(a.size / 2).
			SetColor(theme.SurfaceColor))
	} else if a.children != "" {
		root.Add(components.NewText(a.children).
			SetFontSize(a.size / 2).
			SetColor(theme.SurfaceColor).
			SetFontWeight(fonts.FontWeightBold))
	}

	return root
}

func (a *AntAvatar) SetSize(v float32) *AntAvatar         { a.size = v; return a }
func (a *AntAvatar) SetShape(s AntAvatarShape) *AntAvatar { a.shape = s; return a }
func (a *AntAvatar) SetSrc(s string) *AntAvatar           { a.src = s; return a }
func (a *AntAvatar) SetIcon(i string) *AntAvatar          { a.icon = i; return a }
func (a *AntAvatar) SetAlt(s string) *AntAvatar           { a.alt = s; return a }
func (a *AntAvatar) SetGap(g int) *AntAvatar              { a.gap = g; return a }
func (a *AntAvatar) SetChildren(c string) *AntAvatar      { a.children = c; return a }

// AntAvatarGroup is a group of avatars.
type AntAvatarGroup struct {
	tenon.BaseWidget
	maxCount int
	size     float32
	children []*AntAvatar
}

// NewAntAvatarGroup creates an AntAvatarGroup.
func NewAntAvatarGroup() *AntAvatarGroup {
	g := &AntAvatarGroup{size: 40}
	g.Init(g)
	return g
}

// Render returns the avatar group UI.
func (g *AntAvatarGroup) Render() tenon.Component {
	root := components.NewView().
		SetFlexDirection(yoga.FlexDirectionRow).
		SetAlignItems(yoga.AlignCenter)

	show := g.children
	extra := 0
	if g.maxCount > 0 && len(g.children) > g.maxCount {
		show = g.children[:g.maxCount]
		extra = len(g.children) - g.maxCount
	}

	for _, avatar := range show {
		avatar.SetSize(g.size)
		root.AddChild(avatar)
	}

	if extra > 0 {
		extraAvatar := NewAntAvatar().
			SetSize(g.size).
			SetChildren("+" + string(rune('0'+extra))) // simplified
		root.AddChild(extraAvatar)
	}

	return root
}

func (g *AntAvatarGroup) Add(avatars ...*AntAvatar) *AntAvatarGroup {
	g.children = append(g.children, avatars...)
	return g
}

func (g *AntAvatarGroup) SetMaxCount(n int) *AntAvatarGroup { g.maxCount = n; return g }
func (g *AntAvatarGroup) SetSize(v float32) *AntAvatarGroup { g.size = v; return g }
