package g3n

import (
	"fmt"

	"github.com/sjm1327605995/tenon/pkg/ui"
)

type projectedImageEntry struct {
	bitmap *bitmap
	box    ui.Rect
}

var projectedImages = map[string]projectedImageEntry{}

func projectedImageKey(key string, quad [4]ui.Point) string {
	return fmt.Sprintf("%s|%.2f,%.2f|%.2f,%.2f|%.2f,%.2f|%.2f,%.2f",
		key, quad[0].X, quad[0].Y, quad[1].X, quad[1].Y, quad[2].X, quad[2].Y, quad[3].X, quad[3].Y)
}

func (p *painter) projectedImage(img ui.Bitmap, quad [4]ui.Point, opacity float32, cacheKey string) {
	source, ok := img.(*bitmap)
	if !ok || source == nil {
		return
	}
	key := projectedImageKey(cacheKey, quad)
	entry, ok := projectedImages[key]
	if !ok {
		warped, box := ui.WarpImageToQuad(source.source, quad)
		if warped == nil {
			return
		}
		entry = projectedImageEntry{bitmap: newBitmap(warped), box: box}
		if len(projectedImages) >= 64 {
			clear(projectedImages)
		}
		projectedImages[key] = entry
	}
	p.DrawImage(entry.bitmap, entry.box, opacity)
}
