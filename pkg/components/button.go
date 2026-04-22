package components

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

type ButtonState int

const (
	ButtonStateNormal ButtonState = iota
	ButtonStateHover
	ButtonStatePressed
)

type Button struct {
	core.BaseComponent
	text         *Text
	state        ButtonState
	onClick      func()
	hoverColor   color.Color
	pressedColor color.Color
	normalColor  color.Color
	disabled     bool
}

func NewButton(label string) *Button {
	b := &Button{
		BaseComponent: core.NewBaseComponent(),
		state:         ButtonStateNormal,
		hoverColor:    color.RGBA{R: 70, G: 130, B: 180, A: 255},
		pressedColor:  color.RGBA{R: 30, G: 144, B: 255, A: 255},
		normalColor:   color.RGBA{R: 0, G: 123, B: 255, A: 255},
		disabled:      false,
	}
	b.Init(b)

	b.SetPadding(yoga.EdgeAll, 12)
	b.SetBorderRadius(8)
	b.SetBackgroundColor(b.normalColor)
	b.SetJustifyContent(yoga.JustifyCenter)
	b.SetAlignItems(yoga.AlignCenter)

	b.text = NewText(label)
	b.text.SetColor(color.White)
	b.text.SetFontSize(16)
	b.AddChild(b.text)

	return b
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

func (b *Button) Update() error {
	if b.disabled {
		b.state = ButtonStateNormal
		return nil
	}

	x, y := ebiten.CursorPosition()
	bounds := b.GetLayoutBounds()

	isInside := float32(x) >= bounds.X && float32(x) <= bounds.X+bounds.Width &&
		float32(y) >= bounds.Y && float32(y) <= bounds.Y+bounds.Height

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

	return b.BaseComponent.Update()
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

func (b *Button) GetButtonState() ButtonState {
	return b.state
}

func (b *Button) Draw(screen *ebiten.Image) {
	element := b.Render()
	if element == nil || !element.Visible {
		return
	}

	bounds := b.GetLayoutBounds()
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return
	}

	if element.BackgroundColor != nil {
		vector.FillRect(screen, bounds.X, bounds.Y, bounds.Width, bounds.Height, element.BackgroundColor, false)
	}

	if element.BorderColor != nil {
		yogaNode := element.Yoga
		borderTop := yogaNode.StyleGetBorder(yoga.EdgeTop)
		if borderTop > 0 {
			vector.FillRect(screen, bounds.X, bounds.Y, bounds.Width, borderTop, element.BorderColor, false)
		}
	}

	for _, child := range b.GetChildren() {
		child.Draw(screen)
	}
}

func (b *Button) DrawOverlay(screen *ebiten.Image) {
	for _, child := range b.GetChildren() {
		child.DrawOverlay(screen)
	}
}
