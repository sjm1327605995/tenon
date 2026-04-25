package components

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/core"
	"github.com/sjm1327605995/tenon/yoga"
)

// LPDisplay 是 LP 数值显示组件，支持变化动画。
type LPDisplay struct {
	core.BaseHost
	currentLP   int
	targetLP    int
	changeColor color.Color
	changeTime  float64
	onComplete  func()
	label       *Text
	changeLabel *Text
	animating   bool
}

// NewLPDisplay 创建一个 LP 显示组件。
func NewLPDisplay() *LPDisplay {
	ld := &LPDisplay{
		currentLP: 8000,
		changeColor: color.RGBA{R: 255, G: 0, B: 0, A: 255},
	}
	ld.Init(ld)
	ld.GetElement().Yoga.StyleSetFlexDirection(yoga.FlexDirectionColumn)
	ld.GetElement().Yoga.StyleSetAlignItems(yoga.AlignCenter)

	// LP 数值
	lpText := NewText("8000")
	lpText.SetFontSize(36)
	lpText.SetColor(color.RGBA{R: 0, G: 200, B: 255, A: 255})
	lpText.GetElement().PointerEvents = core.PointerEventsNone
	ld.label = lpText
	ld.AddChild(lpText)

	// LP 标签
	lpLabel := NewText("LP")
	lpLabel.SetFontSize(14)
	lpLabel.SetColor(core.GetTheme().TextMutedColor)
	lpLabel.GetElement().PointerEvents = core.PointerEventsNone
	ld.AddChild(lpLabel)

	// 变化数值（动画时显示）
	changeText := NewText("")
	changeText.SetFontSize(18)
	changeText.GetElement().Visible = false
	changeText.GetElement().PointerEvents = core.PointerEventsNone
	ld.changeLabel = changeText
	ld.AddChild(changeText)

	return ld
}

// SetLP 设置当前 LP 值（无动画）。
func (ld *LPDisplay) SetLP(lp int) *LPDisplay {
	ld.currentLP = lp
	ld.targetLP = lp
	ld.label.SetContent(formatLP(lp))
	return ld
}

// ChangeLP 动画改变 LP 值。
func (ld *LPDisplay) ChangeLP(target int) *LPDisplay {
	ld.targetLP = target
	ld.changeTime = 0
	ld.animating = true

	delta := target - ld.currentLP
	if delta < 0 {
		ld.changeColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
		ld.changeLabel.SetContent("-" + formatLP(-delta))
	} else {
		ld.changeColor = color.RGBA{R: 0, G: 255, B: 0, A: 255}
		ld.changeLabel.SetContent("+" + formatLP(delta))
	}
	ld.changeLabel.SetColor(ld.changeColor)
	ld.changeLabel.GetElement().Visible = true
	return ld
}

// SetOnComplete 设置动画完成回调。
func (ld *LPDisplay) SetOnComplete(fn func()) *LPDisplay {
	ld.onComplete = fn
	return ld
}

// Update 每帧更新 LP 动画。
func (ld *LPDisplay) Update() error {
	if !ld.animating {
		return nil
	}

	ld.changeTime += 0.016 // ~60fps
	duration := 1.5 // 1.5秒动画
	progress := float64(ld.changeTime) / duration
	if progress > 1 {
		progress = 1
	}

	// 缓动函数 easeOutCubic
	ease := 1 - math.Pow(1-progress, 3)
	current := float64(ld.currentLP) + ease*float64(ld.targetLP-ld.currentLP)
	ld.label.SetContent(formatLP(int(current)))

	// 变化数值淡出
	alpha := uint8(255 * (1 - progress))
	if alpha < 0 {
		alpha = 0
	}
	if c, ok := ld.changeColor.(color.RGBA); ok {
		c.A = alpha
		ld.changeLabel.SetColor(c)
	}

	if progress >= 1 {
		ld.currentLP = ld.targetLP
		ld.animating = false
		ld.changeLabel.GetElement().Visible = false
		if ld.onComplete != nil {
			ld.onComplete()
		}
	}
	return nil
}

func formatLP(lp int) string {
	if lp < 0 {
		return "0"
	}
	return itoa(lp)
}

// ==================== 链式 API ====================

func (ld *LPDisplay) SetWidth(width float32) *LPDisplay {
	ld.GetElement().Yoga.StyleSetWidth(width)
	return ld
}
func (ld *LPDisplay) SetMargin(edge yoga.Edge, value float32) *LPDisplay {
	ld.GetElement().Yoga.StyleSetMargin(edge, value)
	return ld
}

// SyncFrom 同步 LP 显示属性。
func (ld *LPDisplay) SyncFrom(other core.Host) {
	if o, ok := other.(*LPDisplay); ok {
		ld.currentLP = o.currentLP
		ld.targetLP = o.targetLP
		ld.animating = o.animating
		ld.onComplete = o.onComplete
		ld.label.SetContent(formatLP(ld.currentLP))
	}
}
