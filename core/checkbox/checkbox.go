package checkbox

import "github.com/sjm1327605995/tenon/widget"

// Label sets the checkbox's static display label.
//
// This is a convenience alias for [LabelOpt].
func Label(s string) Option {
	return LabelOpt(s)
}

// Background sets a custom background color override via option.
//
// This is a convenience alias for [BackgroundOpt].
func Background(color widget.Color) Option {
	return BackgroundOpt(color)
}
