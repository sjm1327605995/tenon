package shadcn

import (
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
	"github.com/sjm1327605995/tenon/pkg/v2/widgets"
)

// Label creates a form label with default text styling.
func Label(text string) widgets.TextWidget {
	theme := ui.GetTheme()
	return widgets.Text(text).
		FontSize(theme.FontSizeBase).
		Color(theme.TextColor)
}
