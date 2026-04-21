package components

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sjm1327605995/tenon/yoga"
)

type ButtonState int

const (
	ButtonStateNormal ButtonState = iota
	ButtonStateHover
	ButtonStatePressed
)

type Button struct {
	*View
	text         *Text
	state        ButtonState
	onClick      func()
	hoverColor   color.Color
	pressedColor color.Color
	normalColor  color.Color
	disabled     bool
}

func NewButton(label string) *Button {
	button := &Button{
		View:         NewView(),
		state:        ButtonStateNormal,
		hoverColor:   color.RGBA{R: 70, G: 130, B: 180, A: 255},
		pressedColor: color.RGBA{R: 30, G: 144, B: 255, A: 255},
		normalColor:  color.RGBA{R: 0, G: 123, B: 255, A: 255},
		disabled:     false,
	}

	button.SetPadding(yoga.EdgeAll, 12)
	button.SetBorderRadius(8)
	button.SetBackgroundColor(button.normalColor)
	button.SetJustifyContent(yoga.JustifyCenter)
	button.SetAlignItems(yoga.AlignCenter)

	button.text = NewText(label)
	button.text.SetColor(color.White)
	button.text.SetFontSize(16)
	button.AddChild(button.text)

	return button
}

func (b *Button) SetOnClick(callback func()) *Button {
	b.onClick = callback
	return b
}

func (b *Button) SetDisabled(disabled bool) *Button {
	b.disabled = disabled
	if disabled {
		b.SetBackgroundColor(color.RGBA{R: 108, G: 117, B: 125, A: 255})
	} else {
		b.SetBackgroundColor(b.normalColor)
	}
	return b
}

// Button 链式方法重写
func (b *Button) SetWidth(width float32) *Button {
	b.View.SetWidth(width)
	return b
}

func (b *Button) SetHeight(height float32) *Button {
	b.View.SetHeight(height)
	return b
}

func (b *Button) SetBackgroundColor(clr color.Color) *Button {
	b.View.SetBackgroundColor(clr)
	return b
}

func (b *Button) SetPadding(edge yoga.Edge, value float32) *Button {
	b.View.SetPadding(edge, value)
	return b
}

func (b *Button) SetBorderRadius(radius float32) *Button {
	b.View.SetBorderRadius(radius)
	return b
}

func (b *Button) SetJustifyContent(justify yoga.Justify) *Button {
	b.View.SetJustifyContent(justify)
	return b
}

func (b *Button) SetAlignItems(align yoga.Align) *Button {
	b.View.SetAlignItems(align)
	return b
}

func (b *Button) SetFlexDirection(dir yoga.FlexDirection) *Button {
	b.View.SetFlexDirection(dir)
	return b
}

func (b *Button) SetMargin(edge yoga.Edge, value float32) *Button {
	b.View.SetMargin(edge, value)
	return b
}

func (b *Button) SetBorder(edge yoga.Edge, value float32) *Button {
	b.View.SetBorder(edge, value)
	return b
}

func (b *Button) SetBorderColor(clr color.Color) *Button {
	b.View.SetBorderColor(clr)
	return b
}

func (b *Button) Update() error {
	if b.disabled {
		b.state = ButtonStateNormal
		return nil
	}

	x, y := ebiten.CursorPosition()
	bounds := b.GetLayoutBounds()

	isInside := float32(x) >= bounds.X && float32(x) <= bounds.X+bounds.Width &&
		float32(y) >= bounds.Y && float32(y) <= bounds.Y+bounds.Height

	// 调试信息
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		fmt.Printf("鼠标点击: (%d, %d), 按钮位置: (%.1f, %.1f, %.1f, %.1f), 是否在按钮内: %v\n",
			x, y, bounds.X, bounds.Y, bounds.Width, bounds.Height, isInside)
	}

	if isInside {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			b.state = ButtonStatePressed
			fmt.Println("按钮按下状态")
		} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			b.state = ButtonStatePressed
		} else {
			b.state = ButtonStateHover
		}

		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) && b.state == ButtonStatePressed {
			fmt.Println("按钮点击事件触发")
			if b.onClick != nil {
				b.onClick()
			}
			b.state = ButtonStateHover
		}
	} else {
		b.state = ButtonStateNormal
	}

	switch b.state {
	case ButtonStateHover:
		b.SetBackgroundColor(b.hoverColor)
	case ButtonStatePressed:
		b.SetBackgroundColor(b.pressedColor)
	default:
		b.SetBackgroundColor(b.normalColor)
	}

	return b.View.Update()
}

func (b *Button) HandleInput() bool {
	return true
}

func (b *Button) SetText(text string) *Button {
	b.text.SetContent(text)
	return b
}

func (b *Button) SetTextColor(clr color.Color) *Button {
	b.text.SetColor(clr)
	return b
}

func (b *Button) SetBackgroundColors(normal, hover, pressed color.Color) *Button {
	b.normalColor = normal
	b.hoverColor = hover
	b.pressedColor = pressed
	b.SetBackgroundColor(b.normalColor)
	return b
}

func (b *Button) GetState() ButtonState {
	return b.state
}
