package components

import (
	"github.com/sjm1327605995/tenon/internal/core"
	"github.com/sjm1327605995/tenon/internal/native"
	"github.com/sjm1327605995/tenon/yoga"
)

// Table is a data table component.
type Table struct {
	native.View
	headerRow *native.View
	bodyEl    *native.View
}

// NewTable creates a table.
func NewTable() *Table {
	theme := core.GetTheme()
	t := &Table{}
	t.Init(t)
	t.SetFlexDirection(yoga.FlexDirectionColumn)
	t.SetWidthPercent(100)
	t.SetBorderColor(theme.BorderColor)

	t.headerRow = native.NewView()
	t.headerRow.SetFlexDirection(yoga.FlexDirectionRow)
	t.headerRow.SetWidthPercent(100)
	t.headerRow.SetBackgroundColor(theme.MutedColor)
	t.headerRow.SetBorder(yoga.EdgeBottom, 1)
	t.headerRow.SetBorderColor(theme.BorderColor)

	t.bodyEl = native.NewView()
	t.bodyEl.SetFlexDirection(yoga.FlexDirectionColumn)
	t.bodyEl.SetWidthPercent(100)

	t.Add(t.headerRow, t.bodyEl)
	return t
}

// ElementType returns type identifier.
func (t *Table) ElementType() string { return "Table" }

// SetHeaders sets the column headers.
func (t *Table) SetHeaders(headers []string) *Table {
	t.headerRow.ClearChildren()
	for _, h := range headers {
		cell := native.NewText(h).SetFontSize(12).SetColor(core.GetTheme().TextColor)
		cell.SetPadding(yoga.EdgeAll, 12)
		cell.SetFlexGrow(1)
		t.headerRow.Add(cell)
	}
	t.Mark(core.FlagNeedLayout)
	return t
}

// AddRow adds a data row.
func (t *Table) AddRow(cells []string) *Table {
	row := native.NewView()
	row.SetFlexDirection(yoga.FlexDirectionRow)
	row.SetWidthPercent(100)
	row.SetBorder(yoga.EdgeBottom, 1)
	row.SetBorderColor(core.GetTheme().BorderColor)
	for _, c := range cells {
		cell := native.NewText(c).SetFontSize(14).SetColor(core.GetTheme().TextColor)
		cell.SetPadding(yoga.EdgeAll, 12)
		cell.SetFlexGrow(1)
		row.Add(cell)
	}
	t.bodyEl.Add(row)
	t.Mark(core.FlagNeedLayout)
	return t
}

// ClearRows removes all data rows.
func (t *Table) ClearRows() *Table {
	t.bodyEl.ClearChildren()
	t.Mark(core.FlagNeedLayout)
	return t
}
