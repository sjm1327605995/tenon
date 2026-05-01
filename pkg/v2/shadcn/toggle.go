package shadcn

import (
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
	"github.com/sjm1327605995/tenon/pkg/v2/widgets"
)

// ShadcnToggle renders a toggle button.
// pressed=true uses Default variant (filled), otherwise Outline.
func ShadcnToggle(label string, pressed bool, onChange func(bool)) ui.Widget {
	var variant widgets.ButtonVariant
	if pressed {
		variant = widgets.ButtonDefault
	} else {
		variant = widgets.ButtonOutline
	}
	return widgets.Button(label).
		Variantf(variant).
		OnTap(func() {
			if onChange != nil {
				onChange(!pressed)
			}
		})
}
