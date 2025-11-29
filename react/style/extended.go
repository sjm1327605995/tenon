package style

import "image/color"

// BackgroundColor 背景色样式结构体
type BackgroundColor struct {
	Color color.NRGBA
}

func (BackgroundColor) ExtendedStyle() {
}

// BorderColor 边框色样式结构体
type BorderColor struct {
	Color color.NRGBA
}

func (BorderColor) ExtendedStyle() {
}
