package widgets

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sjm1327605995/tenon/pkg/render"
	"github.com/sjm1327605995/tenon/pkg/engine"
)

// SpriteWidget 是一个精灵图 / 帧动画 Widget。
// 通过 SpriteSheet 可以切出多帧，或者传入单张图片作为静态精灵。
type SpriteWidget struct {
	engine.BaseWidget

	// frames 精灵帧序列；单张图片时传 []*ebiten.Image{img}
	frames []*ebiten.Image

	// frameDuration 每帧显示时长（默认 100ms）
	frameDuration time.Duration

	// loop 是否循环播放
	loop bool

	// autoPlay 是否自动播放
	autoPlay bool

	// width / height 显式尺寸；若为 0 则使用第一帧图片尺寸
	width  float32
	height float32
}

// Sprite 创建 SpriteWidget；传入图片切片，可配合 SpriteSheet 切图使用。
func Sprite(frames []*ebiten.Image) SpriteWidget {
	return SpriteWidget{
		frames:        frames,
		frameDuration: 100 * time.Millisecond,
		loop:          true,
		autoPlay:      true,
	}
}

func (s SpriteWidget) FrameDuration(d time.Duration) SpriteWidget {
	s.frameDuration = d
	return s
}

func (s SpriteWidget) Loop(v bool) SpriteWidget {
	s.loop = v
	return s
}

func (s SpriteWidget) AutoPlay(v bool) SpriteWidget {
	s.autoPlay = v
	return s
}

func (s SpriteWidget) W(v float32) SpriteWidget {
	s.width = v
	return s
}

func (s SpriteWidget) H(v float32) SpriteWidget {
	s.height = v
	return s
}

func (s SpriteWidget) CreateElement() engine.Element {
	return engine.NewRenderObjectElement(s)
}

func (s SpriteWidget) CreateRenderObject(e engine.Element) render.RenderObject {
	r := render.NewRenderSprite()
	r.Frames = s.frames
	r.FrameDuration = s.frameDuration
	r.Loop = s.loop
	r.AutoPlay = s.autoPlay
	r.SetWidth(s.width)
	r.SetHeight(s.height)
	if s.autoPlay && len(s.frames) > 1 {
		r.Play()
	}
	return r
}

func (s SpriteWidget) UpdateRenderObject(ro render.RenderObject, old engine.Widget) {
	r := ro.(*render.RenderSprite)
	r.Frames = s.frames
	r.FrameDuration = s.frameDuration
	r.Loop = s.loop
	r.AutoPlay = s.autoPlay
	r.SetWidth(s.width)
	r.SetHeight(s.height)
}
