package g3n

import (
	"image"
	"image/draw"
	"sync"

	"github.com/g3n/engine/text"
	"github.com/g3n/engine/texture"
	"github.com/sjm1327605995/tenon/pkg/ui"
)

type fontKey struct {
	px     float32
	weight int
	italic bool
}

type font struct {
	font *text.Font
	px   float32
}

func (f *font) Measure(s string, _ float64) float32 {
	w, _ := f.font.MeasureText(s)
	return float32(w)
}

func (f *font) Metrics() (float32, bool, bool) {
	return float32(f.font.Metrics().Ascent.Ceil()), false, false
}

type bitmap struct {
	texture *texture.Texture2D
	source  *image.NRGBA
	mu      sync.Mutex
	frames  map[uint8]*texture.Texture2D
}

// TextureFromImage creates a native texture using the same pixel path as UI images.
func TextureFromImage(img image.Image) *texture.Texture2D {
	if img == nil {
		return nil
	}
	return newBitmap(img).texture
}

func newBitmap(img image.Image) *bitmap {
	source := image.NewNRGBA(img.Bounds())
	draw.Draw(source, source.Bounds(), img, img.Bounds().Min, draw.Src)
	return &bitmap{texture: texture.NewTexture2DFromRGBA(bitmapRGBA(source, 255)), source: source}
}

func bitmapRGBA(source *image.NRGBA, opacity uint8) *image.RGBA {
	rgba := image.NewRGBA(source.Bounds())
	for i := 0; i < len(source.Pix); i += 4 {
		rgba.Pix[i] = source.Pix[i]
		rgba.Pix[i+1] = source.Pix[i+1]
		rgba.Pix[i+2] = source.Pix[i+2]
		rgba.Pix[i+3] = uint8((uint16(source.Pix[i+3])*uint16(opacity) + 127) / 255)
	}
	return rgba
}

func (b *bitmap) textureForOpacity(opacity uint8) *texture.Texture2D {
	if opacity == 255 {
		return b.texture
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if cached := b.frames[opacity]; cached != nil {
		return cached
	}
	if b.frames == nil {
		b.frames = make(map[uint8]*texture.Texture2D)
	}
	created := texture.NewTexture2DFromRGBA(bitmapRGBA(b.source, opacity))
	b.frames[opacity] = created
	return created
}

func (b *bitmap) Size() (int, int) {
	if b == nil || b.texture == nil {
		return 0, 0
	}
	return b.texture.Width(), b.texture.Height()
}

type vectorPath struct {
	svg   string
	scale float32
}

func (*vectorPath) IsVecPath() {}

var _ ui.FontFace = (*font)(nil)
var _ ui.Bitmap = (*bitmap)(nil)
var _ ui.VecPath = (*vectorPath)(nil)
