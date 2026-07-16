package ui

import (
	"image"

	"gioui.org/gpu/headless"
	"gioui.org/op"
	gpaint "gioui.org/op/paint"
)

// Screenshot 无头渲染一棵 UI 树并返回像素（不需要窗口）。用于像素级黄金测试与调试：
// 与真实运行走同一条 gio 绘制路径（gioPainter），因此能复现真机上的渲染问题。
// 需要可用的 GPU/驱动；不可用时返回 error。
func Screenshot(root *Node, w, h int) (*image.RGBA, error) {
	win, err := headless.NewWindow(w, h)
	if err != nil {
		return nil, err
	}
	defer win.Release()

	hn := Mount(root, w, h)
	g := hn.g

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
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	if err := win.Screenshot(img); err != nil {
		return nil, err
	}
	return img, nil
}
