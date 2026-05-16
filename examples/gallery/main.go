package main

import (
	"log"

	"github.com/sjm1327605995/tenon/app"
	"github.com/sjm1327605995/tenon/core/button"
	"github.com/sjm1327605995/tenon/core/checkbox"
	"github.com/sjm1327605995/tenon/core/datatable"
	"github.com/sjm1327605995/tenon/core/progress"
	"github.com/sjm1327605995/tenon/core/progressbar"
	"github.com/sjm1327605995/tenon/core/slider"
	"github.com/sjm1327605995/tenon/core/tabview"
	"github.com/sjm1327605995/tenon/core/textfield"
	"github.com/sjm1327605995/tenon/core/treeview"
	"github.com/sjm1327605995/tenon/geometry"
	"github.com/sjm1327605995/tenon/primitives"
	"github.com/sjm1327605995/tenon/theme/material3"
	"github.com/sjm1327605995/tenon/widget"
)

func main() {
	m3 := material3.New(widget.Hex(0x6750A4))

	root := buildGallery(m3)

	a := app.New(root)

	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}

func buildGallery(m3 *material3.Theme) widget.Widget {
	bp := material3.ButtonPainter{Theme: m3}
	cp := material3.CheckboxPainter{Theme: m3}
	tp := material3.TextFieldPainter{Theme: m3}
	sp := material3.SliderPainter{Theme: m3}
	pp := material3.ProgressPainter{Theme: m3}
	pbp := material3.ProgressBarPainter{Theme: m3}
	tabp := material3.TabViewPainter{Theme: m3}
	tvp := material3.TreeViewPainter{Theme: m3}
	dtp := material3.DataTablePainter{Theme: m3}

	// Title
	title := primitives.NewText("Tenon Gallery")
	title.FontSize = 28
	title.Bold = true
	title.Color = m3.Colors.Primary

	subtitle := primitives.NewText("Material 3 Theme")
	subtitle.FontSize = 14
	subtitle.Color = m3.Colors.OnSurfaceVariant

	// Buttons
	btnFilled := button.New(button.Text("Filled"), button.VariantOpt(button.Filled), button.PainterOpt(bp), button.OnClick(func() {}))
	btnOutlined := button.New(button.Text("Outlined"), button.VariantOpt(button.Outlined), button.PainterOpt(bp), button.OnClick(func() {}))
	btnText := button.New(button.Text("Text"), button.VariantOpt(button.TextOnly), button.PainterOpt(bp), button.OnClick(func() {}))
	btnTonal := button.New(button.Text("Tonal"), button.VariantOpt(button.Tonal), button.PainterOpt(bp), button.OnClick(func() {}))

	buttonRow := primitives.NewRow(btnFilled, btnOutlined, btnText, btnTonal)
	buttonRow.Gap = 12

	// Checkbox
	chk := checkbox.New(checkbox.LabelOpt("Enable feature"), checkbox.PainterOpt(cp))

	// TextField
	tf := textfield.New(textfield.Placeholder("Enter text..."), textfield.PainterOpt(tp))

	// Slider
	sl := slider.New(slider.Min(0), slider.Max(100), slider.Value(50), slider.PainterOpt(sp))

	// Progress
	prog := progress.New(progress.Value(0.65), progress.PainterOpt(pp))
	progBar := progressbar.New(progressbar.Value(0.65), progressbar.PainterOpt(pbp))

	// TabView
	tabs := tabview.New(
		[]tabview.Tab{
			{Label: "General", Content: primitives.NewBox(primitives.NewText("General settings"))},
			{Label: "Appearance", Content: primitives.NewBox(primitives.NewText("Appearance settings"))},
			{Label: "About", Content: primitives.NewBox(primitives.NewText("Version 1.0"))},
		},
		tabview.PainterOpt(tabp),
	)

	// TreeView
	rootNode := &treeview.TreeNode{
		ID:    "root",
		Label: "Project",
		Expanded: true,
		Children: []*treeview.TreeNode{
			{ID: "src", Label: "src", Expanded: true, Children: []*treeview.TreeNode{
				{ID: "main", Label: "main.go"},
				{ID: "app", Label: "app.go"},
			}},
			{ID: "assets", Label: "assets", Children: []*treeview.TreeNode{
				{ID: "logo", Label: "logo.png"},
			}},
			{ID: "readme", Label: "README.md"},
		},
	}
	tree := treeview.New(
		treeview.Root(rootNode),
		treeview.ItemHeight(28),
		treeview.IndentWidth(20),
		treeview.ShowLines(true),
		treeview.SelectionModeOpt(treeview.SelectionSingle),
		treeview.PainterOpt(tvp),
	)

	// DataTable
	cols := []datatable.Column{
		{Key: "name", Title: "Name", Width: 150, Sortable: true},
		{Key: "type", Title: "Type", Width: 100, Sortable: true},
		{Key: "size", Title: "Size", Width: 80, Sortable: true, Align: widget.TextAlignRight},
	}
	rows := []map[string]string{
		{"name": "main.go", "type": "Go", "size": "2.4 KB"},
		{"name": "app.go", "type": "Go", "size": "5.1 KB"},
		{"name": "widget.go", "type": "Go", "size": "8.7 KB"},
		{"name": "logo.png", "type": "Image", "size": "12.3 KB"},
		{"name": "README.md", "type": "Markdown", "size": "1.8 KB"},
	}
	dt := datatable.New(
		datatable.Columns(cols),
		datatable.RowCount(len(rows)),
		datatable.CellValue(func(row int, col string) string {
			return rows[row][col]
		}),
		datatable.RowHeight(32),
		datatable.PainterOpt(dtp),
	)

	// Layout everything in a scrollable column
	col := primitives.NewColumn(
		primitives.NewBox(title),
		primitives.NewBox(subtitle),
		sectionLabel("Buttons"),
		buttonRow,
		sectionLabel("Checkbox & TextField"),
		chk,
		tf,
		sectionLabel("Slider & Progress"),
		sl,
		prog,
		progBar,
		sectionLabel("TabView"),
		tabs,
		sectionLabel("TreeView"),
		tree,
		sectionLabel("DataTable"),
		dt,
	)
	col.Gap = 12
	col.Padding = geometry.UniformInsets(24)

	// Wrap in a scroll view for vertical scrolling if content exceeds window
	return col
}

func sectionLabel(text string) *primitives.Text {
	t := primitives.NewText(text)
	t.FontSize = 12
	t.Bold = true
	t.Color = widget.Hex(0x6750A4)
	return t
}
