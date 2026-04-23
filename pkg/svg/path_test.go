package svg

import (
	"fmt"
	"testing"
)

func TestParsePathBoundsSearch(t *testing.T) {
	// AntIconSearch 路径
	d := "M909.6 854.5L649.9 594.8C690.2 542.7 714 479.2 714 410.5 714 234.3 571.7 92 395.5 92S77 234.3 77 410.5 219.3 729 395.5 729c68.7 0 132.2-23.8 184.3-64.1l259.7 259.7c7.8 7.8 20.5 7.8 28.3 0l42.8-42.8c7.8-7.8 7.8-20.5 0-28.3zM395.5 643c-128.2 0-232.5-104.3-232.5-232.5S267.3 178 395.5 178 628 282.3 628 410.5 523.7 643 395.5 643z"

	minX, minY, maxX, maxY, err := ParsePathBounds(d)
	if err != nil {
		t.Fatalf("ParsePathBounds failed: %v", err)
	}
	fmt.Printf("Bounds: minX=%f minY=%f maxX=%f maxY=%f\n", minX, minY, maxX, maxY)

	viewW := maxX - minX
	viewH := maxY - minY
	scale := float32(14) / max(viewW, viewH)
	fmt.Printf("viewW=%f viewH=%f scale=%f\n", viewW, viewH, scale)

	path, err := ParsePathScaledAndShifted(d, scale, -minX*scale, -minY*scale)
	if err != nil {
		t.Fatalf("ParsePathScaledAndShifted failed: %v", err)
	}
	if path == nil {
		t.Fatal("path is nil")
	}

	// 验证顶点大致在 0~14 范围内
	// Ebiten vector.Path 没有公开访问顶点的方法，我们只能检查不报错
	fmt.Printf("Path parsed successfully\n")
}

func TestParsePathBoundsSimple(t *testing.T) {
	d := "M0 0 L100 0 L100 100 L0 100 Z"
	minX, minY, maxX, maxY, err := ParsePathBounds(d)
	if err != nil {
		t.Fatalf("ParsePathBounds failed: %v", err)
	}
	fmt.Printf("Simple Bounds: minX=%f minY=%f maxX=%f maxY=%f\n", minX, minY, maxX, maxY)

	scale := float32(1) / max(maxX-minX, maxY-minY)
	path, err := ParsePathScaledAndShifted(d, scale, 0, 0)
	if err != nil {
		t.Fatalf("ParsePathScaledAndShifted failed: %v", err)
	}
	if path == nil {
		t.Fatal("path is nil")
	}
	fmt.Printf("Simple path parsed successfully\n")
}
