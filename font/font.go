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

//go:embed Inter-Bold.ttf
var interBold []byte

var (
	once         sync.Once
	parsedRegular *opentype.Font
	parsedBold    *opentype.Font
	faceCache    sync.Map // key: "size:bold" -> font.Face
)

func initFonts() {
	var err error
	parsedRegular, err = opentype.Parse(interRegular)
	if err != nil {
		panic(fmt.Sprintf("failed to parse Inter-Regular.ttf: %v", err))
	}
	parsedBold, err = opentype.Parse(interBold)
	if err != nil {
		panic(fmt.Sprintf("failed to parse Inter-Bold.ttf: %v", err))
	}
}

// Face returns a font.Face for the given size in pixels and weight.
func Face(size float32, bold bool) font.Face {
	once.Do(initFonts)

	key := fmt.Sprintf("%.1f:%v", size, bold)
	if cached, ok := faceCache.Load(key); ok {
		return cached.(font.Face)
	}

	src := parsedRegular
	if bold {
		src = parsedBold
	}
	face, err := opentype.NewFace(src, &opentype.FaceOptions{
		Size:    float64(size),
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to create font face: %v", err))
	}

	faceCache.Store(key, face)
	return face
}
