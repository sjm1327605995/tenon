package ui

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strconv"

	"gioui.org/gpu/headless"
	"gioui.org/op"
	gpaint "gioui.org/op/paint"
)

// Screenshot 无头渲染一棵 UI 树并返回像素（不需要窗口）。用于像素级黄金测试与调试：
// 与真实运行走同一条 gio 绘制路径（gioPainter），因此能复现真机上的渲染问题。
// 需要可用的 GPU/驱动；不可用时返回 error。
func Screenshot(root *Node, w, h int) (*image.RGBA, error) {
	previousGame, previousScale := activeGame, uiScale
	defer func() { activeGame, uiScale = previousGame, previousScale }()
	hn := Mount(root, w, h)
	defer hn.Close()
	return renderGioImage(hn.g)
}

func renderGioImage(g *game) (*image.RGBA, error) {
	if g == nil || g.w <= 0 || g.h <= 0 {
		return nil, fmt.Errorf("ui: screenshot requires a positive viewport")
	}
	win, err := headless.NewWindow(g.w, g.h)
	if err != nil {
		return nil, err
	}
	defer win.Release()

	var ops op.Ops
	gpaint.ColorOp{Color: nrgba(Color{247, 248, 250, 255})}.Add(&ops)
	gpaint.PaintOp{}.Add(&ops)
	p := newGioPainter(&ops, g.w, g.h)
	if g.rootRN != nil {
		paint(p, g.rootRN)
	}
	for _, pf := range g.portals {
		if pf.overlayRoot != nil {
			paint(p, pf.overlayRoot)
		}
	}
	if err := win.Frame(&ops); err != nil {
		return nil, err
	}
	img := image.NewRGBA(image.Rect(0, 0, g.w, g.h))
	if err := win.Screenshot(img); err != nil {
		return nil, err
	}
	return img, nil
}

type captureSpec struct {
	path      string
	remaining int
	exit      bool
}

var pendingCapture *captureSpec

// Capture saves the application's next rendered frame after afterFrames
// additional frames. An empty path cancels a pending capture. If exit is true,
// the window closes after the capture attempt.
func Capture(path string, afterFrames int, exit bool) {
	if path == "" {
		pendingCapture = nil
		return
	}
	if afterFrames < 0 {
		afterFrames = 0
	}
	pendingCapture = &captureSpec{path: path, remaining: afterFrames, exit: exit}
}

func init() {
	path := os.Getenv("TENON_CAPTURE")
	if path == "" {
		return
	}
	frames := 0
	if value := os.Getenv("TENON_CAPTURE_FRAMES"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
			frames = parsed
		}
	}
	Capture(path, frames, true)
}

func captureGioFrame(g *game) bool {
	spec := pendingCapture
	if spec == nil {
		return false
	}
	if spec.remaining > 0 {
		spec.remaining--
		return false
	}
	pendingCapture = nil
	if err := writeGioCapture(spec.path, g); err != nil {
		fmt.Fprintf(os.Stderr, "ui: capture %q: %v\n", spec.path, err)
	} else {
		fmt.Fprintf(os.Stdout, "ui: captured %s\n", spec.path)
	}
	return spec.exit
}

func writeGioCapture(path string, g *game) error {
	img, err := renderGioImage(g)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
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
