package antdesign

import (
	"github.com/sjm1327605995/tenon"
	"github.com/sjm1327605995/tenon/pkg/components"
	"github.com/sjm1327605995/tenon/pkg/fonts"
	"github.com/sjm1327605995/tenon/yoga"
)

// AntColumn 定义表格列。
type AntColumn struct {
	Title     string
	Key       string
	Width     float32 // 0 = auto
	Render    func(row map[string]any) tenon.Component // 自定义渲染
}

// AntTable 是简化版表格组件。
type AntTable struct {
	tenon.BaseWidget
	columns []AntColumn
	data    []map[string]any
	stripe  bool
	small   bool
}

// NewAntTable 创建表格。
func NewAntTable(columns []AntColumn) *AntTable {
	t := &AntTable{columns: columns}
	t.Init(t)
	return t
}

// Render 返回表格 UI。
func (t *AntTable) Render() tenon.Component {
	theme := NewAntTheme()
	padding := float32(16)
	if t.small {
		padding = 8
	}

	// 根容器
	root := components.NewView().
		SetFlexDirection(yoga.FlexDirectionColumn).
		SetBackgroundColor(theme.SurfaceColor).
		SetBorderRadius(theme.BorderRadius).
		SetBorder(yoga.EdgeAll, 1).
		SetBorderColor(theme.BorderColor).
		SetOverflow(yoga.OverflowHidden)

	// 表头
	header := components.NewView().
		SetFlexDirection(yoga.FlexDirectionRow).
		SetBackgroundColor(theme.TableHeaderBg).
		SetBorder(yoga.EdgeBottom, 1).
		SetBorderColor(theme.BorderColor)

	for _, col := range t.columns {
		th := components.NewView().
			SetPadding(yoga.EdgeAll, padding).
			Add(components.NewText(col.Title).
				SetFontSize(theme.FontSizeSM).
				SetColor(theme.TableHeaderColor).
				SetFontWeight(fonts.FontWeightBold))
		if col.Width > 0 {
			th.SetWidth(col.Width)
		} else {
			th.SetFlexGrow(1)
		}
		header.AddChild(th)
	}
	root.AddChild(header)

	// 数据行
	for i, row := range t.data {
		tr := components.NewView().
			SetFlexDirection(yoga.FlexDirectionRow).
			SetBorder(yoga.EdgeBottom, 1).
			SetBorderColor(theme.BorderColor)

		// 斑马纹
		if t.stripe && i%2 == 1 {
			tr.SetBackgroundColor(theme.TableStripeBg)
		}

		for _, col := range t.columns {
			td := components.NewView().
				SetPadding(yoga.EdgeAll, padding)
			if col.Width > 0 {
				td.SetWidth(col.Width)
			} else {
				td.SetFlexGrow(1)
			}

			var cell tenon.Component
			if col.Render != nil {
				cell = col.Render(row)
			} else {
				v, _ := row[col.Key].(string)
				cell = components.NewText(v).
					SetFontSize(theme.FontSizeBase).
					SetColor(theme.TextColor)
			}
			td.AddChild(cell)
			tr.AddChild(td)
		}
		root.AddChild(tr)
	}

	return root
}

// ==================== 链式 API ====================

func (t *AntTable) SetData(data []map[string]any) *AntTable {
	t.data = data
	return t
}
func (t *AntTable) SetStripe(v bool) *AntTable {
	t.stripe = v
	return t
}
func (t *AntTable) SetSmall(v bool) *AntTable {
	t.small = v
	return t
}
