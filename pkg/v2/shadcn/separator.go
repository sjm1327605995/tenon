package shadcn

import (
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
	"github.com/sjm1327605995/tenon/pkg/v2/widgets"
)

// SeparatorOrientation defines the orientation of a separator.
type SeparatorOrientation int

const (
	SeparatorHorizontal SeparatorOrientation = iota
	SeparatorVertical
)

// ShadcnSeparator visually separates content.
func ShadcnSeparator(orientation SeparatorOrientation) ui.Widget {
	t := ui.GetTheme()
	if orientation == SeparatorHorizontal {
		return widgets.Container(nil).Background(colorToRender(t.DividerColor)).WPct(100).H(1)
	}
	return widgets.Container(nil).Background(colorToRender(t.DividerColor)).W(1).HPct(100)
}
