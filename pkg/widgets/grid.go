package widgets

import (
	"github.com/sjm1327605995/tenon/pkg/engine"
)

// Grid 创建等宽网格布局。
// cols 为列数，gap 为行列间距，children 按从左到右、从上到下的顺序填充。
func Grid(cols int, gap float32, children ...engine.Widget) engine.Widget {
	if cols <= 0 {
		cols = 1
	}
	rows := make([]engine.Widget, 0, (len(children)+cols-1)/cols)
	for i := 0; i < len(children); i += cols {
		rowChildren := make([]engine.Widget, cols)
		for j := 0; j < cols; j++ {
			if idx := i + j; idx < len(children) {
				rowChildren[j] = Container(children[idx]).Grow(1)
			} else {
				rowChildren[j] = Container(nil).Grow(1)
			}
		}
		rows = append(rows, Row(rowChildren...).Gapf(gap).AlignItems(engine.AlignStretch))
	}
	return Column(rows...).Gapf(gap).AlignItems(engine.AlignStretch)
}
