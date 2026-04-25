package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// GridItem 定义网格中的一个单元格数据。
type GridItem struct {
	Key      string
	Content  core.Component
	Row      int
	Col      int
	RowSpan  int
	ColSpan  int
}

// Grid 是网格布局容器，支持行列定义和单元格放置。
type Grid struct {
	core.BaseHost
	cols       []float32 // 列宽定义，支持百分比和固定值
	rows       []float32 // 行高定义
	items      []GridItem
	gap        float32
	cellHosts  map[string]*View
}

// NewGrid 创建一个网格布局容器。
func NewGrid() *Grid {
	g := &Grid{
		gap:       8,
		cellHosts: make(map[string]*View),
	}
	g.Init(g)
	g.GetElement().Yoga.StyleSetFlexDirection(yoga.FlexDirectionColumn)
	return g
}

// SetColumns 设置列定义。正值为固定宽度，负值为百分比（-1 = 100%）。
func (g *Grid) SetColumns(cols ...float32) *Grid {
	g.cols = cols
	return g
}

// SetRows 设置行定义。
func (g *Grid) SetRows(rows ...float32) *Grid {
	g.rows = rows
	return g
}

// SetGap 设置单元格间距。
func (g *Grid) SetGap(gap float32) *Grid {
	g.gap = gap
	return g
}

// AddItem 添加单元格内容。
func (g *Grid) AddItem(item GridItem) *Grid {
	if item.RowSpan <= 0 {
		item.RowSpan = 1
	}
	if item.ColSpan <= 0 {
		item.ColSpan = 1
	}
	g.items = append(g.items, item)
	g.rebuild()
	return g
}

// ClearItems 清空所有单元格。
func (g *Grid) ClearItems() *Grid {
	g.items = nil
	g.cellHosts = make(map[string]*View)
	g.ClearChildren()
	return g
}

// GetCell 获取指定 key 的单元格容器。
func (g *Grid) GetCell(key string) *View {
	return g.cellHosts[key]
}

func (g *Grid) rebuild() {
	g.ClearChildren()
	g.cellHosts = make(map[string]*View)

	if len(g.cols) == 0 || len(g.rows) == 0 {
		return
	}

	// 创建行容器
	for r := 0; r < len(g.rows); r++ {
		rowView := NewView()
		rowView.Init(rowView)
		rowView.GetElement().Yoga.StyleSetFlexDirection(yoga.FlexDirectionRow)
		if g.rows[r] > 0 {
			rowView.GetElement().Yoga.StyleSetHeight(g.rows[r])
		} else {
			rowView.GetElement().Yoga.StyleSetFlexGrow(1)
		}
		if r > 0 {
			rowView.GetElement().Yoga.StyleSetMargin(yoga.EdgeTop, g.gap)
		}
		g.AddChild(rowView)

		// 为每列创建单元格
		for c := 0; c < len(g.cols); c++ {
			cell := NewView()
			cell.Init(cell)
			cell.GetElement().Yoga.StyleSetFlexDirection(yoga.FlexDirectionColumn)
			if g.cols[c] > 0 {
				cell.GetElement().Yoga.StyleSetWidth(g.cols[c])
			} else {
				cell.GetElement().Yoga.StyleSetFlexGrow(1)
			}
			if c > 0 {
				cell.GetElement().Yoga.StyleSetMargin(yoga.EdgeLeft, g.gap)
			}
			rowView.AddChild(cell)
		}
	}

	// 放置内容到对应单元格
	for _, item := range g.items {
		if item.Row < 0 || item.Row >= len(g.rows) ||
			item.Col < 0 || item.Col >= len(g.cols) {
			continue
		}
		rowIdx := item.Row
		colIdx := item.Col
		if rowIdx >= len(g.children) {
			continue
		}
		rowHost := g.children[rowIdx].(Host)
		if colIdx >= len(rowHost.GetChildren()) {
			continue
		}
		cell := rowHost.GetChildren()[colIdx].(Host)
		if item.Content != nil {
			cell.AddChild(item.Content)
		}
		if item.Key != "" {
			g.cellHosts[item.Key] = cell.(*View)
		}
	}
}

// Draw 绘制网格背景和边框。
func (g *Grid) Draw(screen *ebiten.Image) {
	el := g.GetElement()
	if el == nil || !el.Visible {
		return
	}
	bounds := g.GetLayoutBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	if el.BackgroundColor != nil {
		vector.FillRect(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, el.BackgroundColor, false)
	}
	if el.BorderColor != nil {
		vector.StrokeRect(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, 1, el.BorderColor, false)
	}
}

// ==================== 链式 API ====================

func (g *Grid) SetWidth(width float32) *Grid {
	g.GetElement().Yoga.StyleSetWidth(width)
	return g
}
func (g *Grid) SetHeight(height float32) *Grid {
	g.GetElement().Yoga.StyleSetHeight(height)
	return g
}
func (g *Grid) SetMargin(edge yoga.Edge, value float32) *Grid {
	g.GetElement().Yoga.StyleSetMargin(edge, value)
	return g
}
func (g *Grid) SetPadding(edge yoga.Edge, value float32) *Grid {
	g.GetElement().Yoga.StyleSetPadding(edge, value)
	return g
}
func (g *Grid) SetBackgroundColor(clr color.Color) *Grid {
	g.GetElement().BackgroundColor = clr
	return g
}

// SyncFrom 同步网格属性。
func (g *Grid) SyncFrom(other core.Host) {
	if o, ok := other.(*Grid); ok {
		g.cols = o.cols
		g.rows = o.rows
		g.items = o.items
		g.gap = o.gap
		g.rebuild()
	}
}
