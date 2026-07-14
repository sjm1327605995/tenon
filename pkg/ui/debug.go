package ui

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

// 调试帧捕获：把引擎自己渲染的帧（离屏目标）存成 PNG，用于无头验证渲染效果。
// 只捕获应用自身的像素（绝不涉及桌面/其它窗口），安全。
var (
	capturePath  string
	captureAfter int
	captureExit  bool
	frameCount   int
)

// Capture 请求在第 afterFrames 帧把渲染结果保存为 PNG 到 path；
// exit 为真时保存后退出进程（便于脚本化/无头截图）。
// 也可用环境变量：TENON_CAPTURE=out.png（可选 TENON_CAPTURE_FRAMES，默认 90）。
func Capture(path string, afterFrames int, exit bool) {
	capturePath, captureAfter, captureExit = path, afterFrames, exit
	frameCount = 0
}

// maybeCapture 在每次 Draw 末尾调用；到达目标帧则保存并（可选）退出。
func (g *game) maybeCapture(screen *ebiten.Image) {
	if capturePath == "" {
		return
	}
	frameCount++
	if frameCount < captureAfter {
		return
	}
	path := capturePath
	capturePath = ""
	if err := saveFramePNG(screen, path); err != nil {
		fmt.Fprintln(os.Stderr, "capture failed:", err)
	} else {
		fmt.Fprintln(os.Stderr, "captured frame ->", path)
	}
	if captureExit {
		os.Exit(0)
	}
}

func saveFramePNG(img *ebiten.Image, path string) error {
	b := img.Bounds()
	rgba := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	img.ReadPixels(rgba.Pix)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, rgba)
}
