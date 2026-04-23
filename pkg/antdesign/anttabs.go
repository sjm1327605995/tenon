package antdesign

import (
	"image/color"

	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/components"
	"github.com/sjm1327605995/tenon/yoga"
)

// AntTabType defines tab bar style.
type AntTabType string

const (
	AntTabLine AntTabType = "line"
	AntTabCard AntTabType = "card"
)

// AntTabItem defines a single tab item.
type AntTabItem struct {
	Key      string
	Label    string
	Disabled bool
	Children tenon.Component
}

// AntTabs is a tab component.
type AntTabs struct {
	tenon.BaseWidget
	items     []AntTabItem
	activeKey string
	tabType   AntTabType
	size      string // small/default/large
	centered  bool
	onChange  func(key string)
}

// NewAntTabs creates an AntTabs.
func NewAntTabs() *AntTabs {
	t := &AntTabs{tabType: AntTabLine}
	t.Init(t)
	return t
}

// Render returns the tabs UI.
func (t *AntTabs) Render() tenon.Component {
	theme := NewAntTheme()

	root := components.NewView().
		SetFlexDirection(yoga.FlexDirectionColumn)

	// Tab bar
	bar := components.NewView().
		SetFlexDirection(yoga.FlexDirectionRow).
		SetBorder(yoga.EdgeBottom, 1).
		SetBorderColor(theme.BorderColor)
	if t.centered {
		bar.SetJustifyContent(yoga.JustifyCenter)
	}

	for _, item := range t.items {
		isActive := item.Key == t.activeKey
		tab := t.renderTab(item, isActive, theme)
		bar.AddChild(tab)
	}
	root.AddChild(bar)

	// Content area
	for _, item := range t.items {
		if item.Key == t.activeKey && item.Children != nil {
			root.AddChild(item.Children)
			break
		}
	}

	return root
}

func (t *AntTabs) renderTab(item AntTabItem, isActive bool, theme *AntTheme) tenon.Component {
	var bg, clr color.Color
	if t.tabType == AntTabCard {
		if isActive {
			bg = theme.SurfaceColor
			clr = theme.PrimaryColor
		} else {
			bg = theme.BackgroundColor
			clr = theme.TextColor
		}
	} else {
		bg = nil
		if isActive {
			clr = theme.PrimaryColor
		} else {
			clr = theme.TextColor
		}
	}

	label := components.NewText(item.Label).
		SetFontSize(theme.FontSizeBase).
		SetColor(clr)

	wrapper := components.NewView().
		SetFlexDirection(yoga.FlexDirectionColumn).
		SetPadding(yoga.EdgeHorizontal, 16).
		SetPadding(yoga.EdgeVertical, 12)

	if bg != nil {
		wrapper.SetBackgroundColor(bg)
	}

	wrapper.AddChild(label)

	if isActive && t.tabType == AntTabLine {
		activeBar := components.NewView().
			SetHeight(2).
			SetBackgroundColor(theme.PrimaryColor)
		wrapper.AddChild(activeBar)
	}

	if !item.Disabled {
		// Ideally add click handler; simplified here
	}

	return wrapper
}

func (t *AntTabs) SetItems(items []AntTabItem) *AntTabs { t.items = items; return t }
func (t *AntTabs) SetActiveKey(k string) *AntTabs       { t.activeKey = k; return t }
func (t *AntTabs) SetType(tp AntTabType) *AntTabs       { t.tabType = tp; return t }
func (t *AntTabs) SetSize(s string) *AntTabs            { t.size = s; return t }
func (t *AntTabs) SetCentered(v bool) *AntTabs          { t.centered = v; return t }
func (t *AntTabs) SetOnChange(fn func(string)) *AntTabs { t.onChange = fn; return t }

// AntTabPane is a convenience wrapper for tab content (kept for API compatibility).
type AntTabPane struct {
	tenon.BaseWidget
	key      string
	label    string
	children []tenon.Component
}

// NewAntTabPane creates an AntTabPane.
func NewAntTabPane(key, label string) *AntTabPane {
	p := &AntTabPane{key: key, label: label}
	p.Init(p)
	return p
}

func (p *AntTabPane) Add(children ...tenon.Component) *AntTabPane {
	p.children = append(p.children, children...)
	return p
}

func (p *AntTabPane) Render() tenon.Component {
	root := components.NewView().SetFlexDirection(yoga.FlexDirectionColumn)
	for _, child := range p.children {
		root.AddChild(child)
	}
	return root
}
