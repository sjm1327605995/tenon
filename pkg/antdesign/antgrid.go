package antdesign

import (
	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/components"
	"github.com/sjm1327605995/tenon/yoga"
)

// AntRow is a row container for the 24-column grid system.
type AntRow struct {
	tenon.BaseWidget
	gutter   float32
	justify  yoga.Justify
	align    yoga.Align
	wrap     bool
	children []tenon.Component
}

// NewAntRow creates an AntRow.
func NewAntRow() *AntRow {
	r := &AntRow{
		justify: yoga.JustifyFlexStart,
		align:   yoga.AlignFlexStart,
		wrap:    true,
	}
	r.Init(r)
	return r
}

// Render returns the row UI.
func (r *AntRow) Render() tenon.Component {
	row := components.NewView().
		SetFlexDirection(yoga.FlexDirectionRow).
		SetJustifyContent(r.justify).
		SetAlignItems(r.align)
	if r.wrap {
		row.SetFlexWrap(yoga.WrapWrap)
	}

	for _, child := range r.children {
		row.AddChild(child)
	}
	return row
}

func (r *AntRow) Add(children ...tenon.Component) *AntRow {
	r.children = append(r.children, children...)
	return r
}

func (r *AntRow) SetGutter(g float32) *AntRow       { r.gutter = g; return r }
func (r *AntRow) SetJustify(j yoga.Justify) *AntRow { r.justify = j; return r }
func (r *AntRow) SetAlign(a yoga.Align) *AntRow     { r.align = a; return r }
func (r *AntRow) SetWrap(w bool) *AntRow            { r.wrap = w; return r }

// AntCol is a column in the 24-column grid system.
type AntCol struct {
	tenon.BaseWidget
	span     int // 0-24
	offset   int
	push     int
	pull     int
	order    int
	children []tenon.Component
}

// NewAntCol creates an AntCol.
func NewAntCol() *AntCol {
	c := &AntCol{span: 24}
	c.Init(c)
	return c
}

// Render returns the column UI.
func (c *AntCol) Render() tenon.Component {
	col := components.NewView()

	if c.span > 0 && c.span <= 24 {
		percent := float32(c.span) * 100.0 / 24.0
		col.SetWidthPercent(percent)
	}
	if c.offset > 0 && c.offset <= 24 {
		percent := float32(c.offset) * 100.0 / 24.0
		col.SetMargin(yoga.EdgeLeft, percent) // margin-left as offset
	}
	if c.order != 0 {
		// Yoga does not support order directly in this binding
	}

	for _, child := range c.children {
		col.AddChild(child)
	}
	return col
}

func (c *AntCol) Add(children ...tenon.Component) *AntCol {
	c.children = append(c.children, children...)
	return c
}

func (c *AntCol) SetSpan(s int) *AntCol   { c.span = s; return c }
func (c *AntCol) SetOffset(o int) *AntCol { c.offset = o; return c }
func (c *AntCol) SetPush(p int) *AntCol   { c.push = p; return c }
func (c *AntCol) SetPull(p int) *AntCol   { c.pull = p; return c }
func (c *AntCol) SetOrder(o int) *AntCol  { c.order = o; return c }
