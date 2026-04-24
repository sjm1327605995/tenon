package antdesign

import (
	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/components"
	"github.com/sjm1327605995/tenon/yoga"
)

// AntLayout is the top-level layout container.
type AntLayout struct {
	tenon.BaseWidget
	children []tenon.Component
}

// NewAntLayout creates an AntLayout.
func NewAntLayout() *AntLayout {
	l := &AntLayout{}
	l.Init(l)
	return l
}

// Render returns the layout UI.
func (l *AntLayout) Render() tenon.Component {
	root := components.NewView().
		SetFlexDirection(yoga.FlexDirectionColumn).
		SetWidthPercent(100).
		SetHeightPercent(100)
	for _, child := range l.children {
		root.AddChild(child)
	}
	return root
}

func (l *AntLayout) Add(children ...tenon.Component) *AntLayout {
	l.children = append(l.children, children...)
	return l
}

// AntHeader is the top navigation area.
type AntHeader struct {
	tenon.BaseWidget
	height   float32
	children []tenon.Component
}

// NewAntHeader creates an AntHeader.
func NewAntHeader() *AntHeader {
	h := &AntHeader{height: 64}
	h.Init(h)
	return h
}

// Render returns the header UI.
func (h *AntHeader) Render() tenon.Component {
	theme := NewAntTheme()
	header := components.NewView().
		SetFlexDirection(yoga.FlexDirectionRow).
		SetAlignItems(yoga.AlignCenter).
		SetHeight(h.height).
		SetPadding(yoga.EdgeHorizontal, 24).
		SetBackgroundColor(theme.SurfaceColor)
	for _, child := range h.children {
		header.AddChild(child)
	}
	return header
}

func (h *AntHeader) Add(children ...tenon.Component) *AntHeader {
	h.children = append(h.children, children...)
	return h
}

func (h *AntHeader) SetHeight(v float32) *AntHeader { h.height = v; return h }

// AntSider is the sidebar area.
type AntSider struct {
	tenon.BaseWidget
	width          float32
	collapsed      bool
	collapsedWidth float32
	children       []tenon.Component
}

// NewAntSider creates an AntSider.
func NewAntSider() *AntSider {
	s := &AntSider{width: 200, collapsedWidth: 80}
	s.Init(s)
	return s
}

// Render returns the sider UI.
func (s *AntSider) Render() tenon.Component {
	theme := NewAntTheme()
	w := s.width
	if s.collapsed {
		w = s.collapsedWidth
	}
	sider := components.NewView().
		SetFlexDirection(yoga.FlexDirectionColumn).
		SetWidth(w).
		SetHeightPercent(100).
		SetBackgroundColor(theme.SurfaceColor)
	for _, child := range s.children {
		sider.AddChild(child)
	}
	return sider
}

func (s *AntSider) Add(children ...tenon.Component) *AntSider {
	s.children = append(s.children, children...)
	return s
}

func (s *AntSider) SetWidth(v float32) *AntSider          { s.width = v; return s }
func (s *AntSider) SetCollapsed(v bool) *AntSider         { s.collapsed = v; return s }
func (s *AntSider) SetCollapsedWidth(v float32) *AntSider { s.collapsedWidth = v; return s }

// AntContent is the main content area.
type AntContent struct {
	tenon.BaseWidget
	children []tenon.Component
}

// NewAntContent creates an AntContent.
func NewAntContent() *AntContent {
	c := &AntContent{}
	c.Init(c)
	return c
}

// Render returns the content UI.
func (c *AntContent) Render() tenon.Component {
	theme := NewAntTheme()
	content := components.NewView().
		SetFlexDirection(yoga.FlexDirectionColumn).
		SetFlexGrow(1).
		SetPadding(yoga.EdgeAll, 24).
		SetBackgroundColor(theme.BackgroundColor)
	for _, child := range c.children {
		content.AddChild(child)
	}
	return content
}

func (c *AntContent) Add(children ...tenon.Component) *AntContent {
	c.children = append(c.children, children...)
	return c
}

// AntFooter is the bottom area.
type AntFooter struct {
	tenon.BaseWidget
	height   float32
	children []tenon.Component
}

// NewAntFooter creates an AntFooter.
func NewAntFooter() *AntFooter {
	f := &AntFooter{height: 64}
	f.Init(f)
	return f
}

// Render returns the footer UI.
func (f *AntFooter) Render() tenon.Component {
	theme := NewAntTheme()
	footer := components.NewView().
		SetFlexDirection(yoga.FlexDirectionRow).
		SetAlignItems(yoga.AlignCenter).
		SetHeight(f.height).
		SetPadding(yoga.EdgeHorizontal, 24).
		SetBackgroundColor(theme.SurfaceColor)
	for _, child := range f.children {
		footer.AddChild(child)
	}
	return footer
}

func (f *AntFooter) Add(children ...tenon.Component) *AntFooter {
	f.children = append(f.children, children...)
	return f
}

func (f *AntFooter) SetHeight(v float32) *AntFooter { f.height = v; return f }
