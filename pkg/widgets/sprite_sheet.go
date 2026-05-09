package widgets

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// SpriteSheet 从单张大图中按固定网格切出多帧。
type SpriteSheet struct {
	// Source 原始精灵图
	Source *ebiten.Image

	// FrameWidth / FrameHeight 每帧的宽度和高度（像素）
	FrameWidth  int
	FrameHeight int

	// Columns / Rows 精灵图的列数和行数；若为 0 则自动根据图片尺寸计算
	Columns int
	Rows    int

	// Margin 帧与帧之间的间隔像素（默认 0）
	Margin int

	// Spacing 切图后每帧之间的额外间隔（默认 0）
	Spacing int
}

// NewSpriteSheet 创建一个 SpriteSheet。
// frameW/frameH 为每帧尺寸；cols/rows 为 0 时自动计算。
func NewSpriteSheet(source *ebiten.Image, frameW, frameH int) *SpriteSheet {
	return &SpriteSheet{
		Source:      source,
		FrameWidth:  frameW,
		FrameHeight: frameH,
	}
}

// Frames 切出所有帧并返回切片，可直接传入 Sprite()。
func (s *SpriteSheet) Frames() []*ebiten.Image {
	if s.Source == nil || s.FrameWidth <= 0 || s.FrameHeight <= 0 {
		return nil
	}

	srcW := s.Source.Bounds().Dx()
	srcH := s.Source.Bounds().Dy()

	cols := s.Columns
	rows := s.Rows
	if cols <= 0 {
		cols = srcW / (s.FrameWidth + s.Spacing)
	}
	if rows <= 0 {
		rows = srcH / (s.FrameHeight + s.Spacing)
	}
	if cols <= 0 || rows <= 0 {
		return nil
	}

	frames := make([]*ebiten.Image, 0, cols*rows)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			sx := s.Margin + x*(s.FrameWidth+s.Spacing)
			sy := s.Margin + y*(s.FrameHeight+s.Spacing)
			if sx+s.FrameWidth > srcW || sy+s.FrameHeight > srcH {
				continue
			}
			frame := s.Source.SubImage(image.Rect(sx, sy, sx+s.FrameWidth, sy+s.FrameHeight)).(*ebiten.Image)
			frames = append(frames, frame)
		}
	}
	return frames
}

// Frame 获取指定索引的单帧；越界返回 nil。
func (s *SpriteSheet) Frame(index int) *ebiten.Image {
	frames := s.Frames()
	if index < 0 || index >= len(frames) {
		return nil
	}
	return frames[index]
}

// SpriteFromSheet 从 SpriteSheet 直接创建 SpriteWidget，等价于 Sprite(sheet.Frames())。
func SpriteFromSheet(sheet *SpriteSheet) SpriteWidget {
	return Sprite(sheet.Frames())
}
