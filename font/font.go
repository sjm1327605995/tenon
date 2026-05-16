package font

import (
	_ "embed"
	"fmt"
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed Inter-Regular.ttf
var interRegular []byte

var (
	once      sync.Once
	parsed    *opentype.Font
	faceCache sync.Map // key: float32(size) -> font.Face
)

func initFont() {
	var err error
	parsed, err = opentype.Parse(interRegular)
	if err != nil {
		panic(fmt.Sprintf("failed to parse Inter-Regular.ttf: %v", err))
	}
}

// Face returns a font.Face for the given size in pixels.
func Face(size float32) font.Face {
	once.Do(initFont)

	if cached, ok := faceCache.Load(size); ok {
		return cached.(font.Face)
	}

	face, err := opentype.NewFace(parsed, &opentype.FaceOptions{
		Size:    float64(size),
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to create font face: %v", err))
	}

	faceCache.Store(size, face)
	return face
}

