package shadcn

import (
	"github.com/sjm1327605995/tenon/pkg/engine"
	"github.com/sjm1327605995/tenon/pkg/widgets"
)

// Label creates a form label with default text styling.
func Label(text string) widgets.TextWidget {
	theme := engine.GetTheme()
	return widgets.Text(text).
		FontSize(theme.FontSizeBase).
		Color(theme.TextColor)
}
