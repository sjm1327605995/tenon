// Package styles provides the style system implementation for the React framework.
package styles

import (
	"image/color"
)

// BackgroundColor represents the background color style of an element.
// It implements the IExtendedStyle interface and is used to set the background color of an element.
type BackgroundColor struct {
	// Color is the NRGBA color value of the background, containing red, green, blue, and alpha channels.
	Color color.NRGBA
}

// ExtendedStyle implements the marker method for the IExtendedStyle interface.
func (BackgroundColor) ExtendedStyle() {
}

// BorderColor represents the border color style of an element.
// It implements the IExtendedStyle interface and is used to set the border color of an element.
type BorderColor struct {
	// Color is the NRGBA color value of the border, containing red, green, blue, and alpha channels.
	Color color.NRGBA
}

// ExtendedStyle 实现IExtendedStyle接口的标记方法。
func (BorderColor) ExtendedStyle() {
}

type CornerRadius struct {
	TopLeft, TopRight, BottomRight, BottomLeft float32
}

func (CornerRadius) ExtendedStyle() {
}
