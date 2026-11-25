package ui

import "testing"

func layoutBenchTree(rows int) *Node {
	children := make([]*Node, 0, rows)
	for i := 0; i < rows; i++ {
		children = append(children, Div(
			Style(Row, Gap(8)),
			Text("label"),
			Div(Style(Width(40), Height(20)), Text("value")),
		))
	}
	return Div(append([]*Node{Style(Column, Gap(4))}, children...)...)
}

func BenchmarkRelink(b *testing.B) {
	h := MountDefault(layoutBenchTree(200))
	defer h.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		relink(h.g.rootFiber)
	}
}

func BenchmarkLayoutCleanTree(b *testing.B) {
	h := MountDefault(layoutBenchTree(200))
	defer h.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.g.layout()
	}
}

func BenchmarkFullFrameComplexTree(b *testing.B) {
	h := MountDefault(layoutBenchTree(500))
	defer h.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Step(0)
		h.PaintTo(discardPainter{})
	}
}

func BenchmarkStartupComplexTree(b *testing.B) {
	root := layoutBenchTree(500)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h := MountDefault(root)
		h.PaintTo(discardPainter{})
		h.Close()
	}
}

func BenchmarkIdleFrameComplexTree(b *testing.B) {
	h := MountDefault(layoutBenchTree(500))
	defer h.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Step(0)
	}
}

func BenchmarkResolveInherited(b *testing.B) {
	h := MountDefault(layoutBenchTree(200))
	defer h.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resolveInherited(h.g.rootRN, inhText{})
	}
}

func BenchmarkResolveInheritedDirtySubtree(b *testing.B) {
	h := MountDefault(layoutBenchTree(200))
	defer h.Close()
	branch := h.g.rootRN.children[0]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		branch.markInheritedDirty()
		resolveInherited(h.g.rootRN, inhText{})
	}
}

type discardPainter struct{}

func (discardPainter) FillRect(float32, float32, float32, float32, float32, Color) {}
func (discardPainter) FillGradient(float32, float32, float32, float32, float32, Color, Color, float32) {
}
func (discardPainter) StrokeRect(float32, float32, float32, float32, float32, float32, Color) {
}
func (discardPainter) Shadow(float32, float32, float32, float32, float32, float32, float32, float32, float32, Color) {
}
func (discardPainter) Line(float32, float32, float32, float32, Color) {}
func (discardPainter) DrawText(string, FontFace, Color, float32, float32, bool, bool) {
}
func (discardPainter) DrawImage(Bitmap, Rect, float32)                                  {}
func (discardPainter) FillPath(VecPath, float32, float32, Color)                        {}
func (discardPainter) StrokePath(VecPath, float32, float32, float32, Color)             {}
func (discardPainter) PushClip(Rect, float32)                                           {}
func (discardPainter) PopClip()                                                         {}
func (discardPainter) BeginLayer()                                                      {}
func (discardPainter) EndLayer(LayerTransform)                                          {}
func (discardPainter) BeginProjected(LayerTransform)                                    {}
func (discardPainter) BeginProjectedElement(Rect)                                       {}
func (discardPainter) EndProjectedElement()                                             {}
func (discardPainter) ProjectedImage(Bitmap, [4]Point, float32, string, *Rect, float32) {}
func (discardPainter) EndProjected()                                                    {}
