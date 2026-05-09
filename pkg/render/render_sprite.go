package render

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/yoga"
)

// RenderSprite 负责绘制精灵图 / 帧动画。
// 支持多帧循环播放或单帧静止显示。
type RenderSprite struct {
	BaseRenderObject

	// Frames 是精灵的帧序列，每帧都是一个 *ebiten.Image
	Frames []*ebiten.Image

	// FrameDuration 是每帧显示时长（默认 100ms）
	FrameDuration time.Duration

	// Loop 为 true 时循环播放，false 时停在最后一帧
	Loop bool

	// AutoPlay 为 true 且在 Mount 时帧数>1，则自动开始播放
	AutoPlay bool

	elapsed      time.Duration
	currentFrame int
	playing      bool
}

func NewRenderSprite() *RenderSprite {
	r := &RenderSprite{
		FrameDuration: 100 * time.Millisecond,
		Loop:          true,
		AutoPlay:      true,
	}
	r.BaseRenderObject.Init(r)
	r.yoga = yoga.NewNode()
	return r
}

func (r *RenderSprite) SetWidth(v float32) {
	if v > 0 {
		r.yoga.StyleSetWidth(v)
	} else {
		r.yoga.StyleSetWidthAuto()
	}
	r.MarkNeedsLayout()
}

func (r *RenderSprite) SetHeight(v float32) {
	if v > 0 {
		r.yoga.StyleSetHeight(v)
	} else {
		r.yoga.StyleSetHeightAuto()
	}
	r.MarkNeedsLayout()
}

// Tick 推进动画时间，返回是否需要重绘。
func (r *RenderSprite) Tick(dt time.Duration) bool {
	if !r.playing || len(r.Frames) <= 1 {
		return false
	}
	r.elapsed += dt
	total := time.Duration(len(r.Frames)) * r.FrameDuration
	if total == 0 {
		return false
	}
	oldFrame := r.currentFrame
	if r.Loop {
		r.currentFrame = int((r.elapsed % total) / r.FrameDuration)
	} else {
		r.currentFrame = int(r.elapsed / r.FrameDuration)
		if r.currentFrame >= len(r.Frames) {
			r.currentFrame = len(r.Frames) - 1
			r.playing = false
		}
	}
	return oldFrame != r.currentFrame
}

// Play 开始播放动画
func (r *RenderSprite) Play() {
	if len(r.Frames) <= 1 {
		return
	}
	r.playing = true
	r.MarkNeedsPaint()
}

// Stop 停止播放动画
func (r *RenderSprite) Stop() {
	r.playing = false
}

// Rewind 重置到第一帧
func (r *RenderSprite) Rewind() {
	r.elapsed = 0
	r.currentFrame = 0
	r.MarkNeedsPaint()
}

func (r *RenderSprite) Paint(screen *ebiten.Image, offset Offset) {
	bounds := r.GetBounds()
	if len(r.Frames) == 0 {
		return
	}
	frame := r.Frames[r.currentFrame]
	if frame == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	// 获取实际绘制尺寸：使用 bounds 尺寸；若为 0 则使用图片原始尺寸
	drawW := bounds.Width
	drawH := bounds.Height
	if drawW == 0 {
		drawW = float32(frame.Bounds().Dx())
	}
	if drawH == 0 {
		drawH = float32(frame.Bounds().Dy())
	}
	sw := float64(drawW) / float64(frame.Bounds().Dx())
	sh := float64(drawH) / float64(frame.Bounds().Dy())
	if sw != 1 || sh != 1 {
		op.GeoM.Scale(sw, sh)
	}
	op.GeoM.Translate(float64(bounds.X+offset.X), float64(bounds.Y+offset.Y))
	screen.DrawImage(frame, op)
}
