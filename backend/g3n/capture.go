package g3n

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strconv"

	"github.com/g3n/engine/gls"
)

type nativeCapture struct {
	path      string
	remaining int
}

func nativeCaptureFromEnv() *nativeCapture {
	path := os.Getenv("TENON_CAPTURE")
	if path == "" {
		return nil
	}
	frames := 0
	if value := os.Getenv("TENON_CAPTURE_FRAMES"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
			frames = parsed
		}
	}
	return &nativeCapture{path: path, remaining: frames}
}

func (c *nativeCapture) capture(gs *gls.GLS, width, height int) bool {
	if c == nil {
		return false
	}
	if c.remaining > 0 {
		c.remaining--
		return false
	}
	if err := writeNativeCapture(c.path, gs, width, height); err != nil {
		fmt.Fprintf(os.Stderr, "g3n: capture %q: %v\n", c.path, err)
	} else {
		fmt.Fprintf(os.Stdout, "g3n: captured %s\n", c.path)
	}
	return true
}

func writeNativeCapture(path string, gs *gls.GLS, width, height int) error {
	if width <= 0 || height <= 0 {
		return fmt.Errorf("invalid framebuffer size %dx%d", width, height)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	img := readNativeImage(gs, width, height)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	if err := png.Encode(file, img); err != nil {
		file.Close()
		return err
	}
	return file.Close()
}

func readNativeImage(gs *gls.GLS, width, height int) *image.NRGBA {
	pixels := gs.ReadPixels(0, 0, width, height, gls.RGBA, gls.UNSIGNED_BYTE)
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	stride := width * 4
	for y := 0; y < height; y++ {
		copy(img.Pix[y*stride:(y+1)*stride], pixels[(height-y-1)*stride:(height-y)*stride])
	}
	return img
}
