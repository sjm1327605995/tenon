package ui

import "testing"

type recSink struct {
	moves, lines, cubes, quads, closes int
	lastX, lastY                       float32
}

func (s *recSink) moveTo(x, y float32) { s.moves++; s.lastX, s.lastY = x, y }
func (s *recSink) lineTo(x, y float32) { s.lines++; s.lastX, s.lastY = x, y }
func (s *recSink) cubeTo(x1, y1, x2, y2, x, y float32) {
	s.cubes++
	s.lastX, s.lastY = x, y
}
func (s *recSink) quadTo(x1, y1, x, y float32) { s.quads++; s.lastX, s.lastY = x, y }
func (s *recSink) closePath()                  { s.closes++ }

func TestSVGParseBasic(t *testing.T) {
	s := &recSink{}
	parseSVGInto("M0 0 L10 0 L10 10 Z", 1, s)
	if s.moves != 1 || s.lines != 2 || s.closes != 1 {
		t.Fatalf("矩形路径: moves=%d lines=%d closes=%d want 1/2/1", s.moves, s.lines, s.closes)
	}
}

func TestSVGParseScaleAndCurves(t *testing.T) {
	s := &recSink{}
	// 绝对三次贝塞尔 + 相对折线，scale=2
	parseSVGInto("M1 1 C2 2 3 3 4 4 l1 0", 2, s)
	if s.moves != 1 || s.cubes != 1 || s.lines != 1 {
		t.Fatalf("moves=%d cubes=%d lines=%d want 1/1/1", s.moves, s.cubes, s.lines)
	}
	if s.lastX != (4+1)*2 || s.lastY != 4*2 { // 终点 (4,4)+相对(1,0)=(5,4)，再*scale
		t.Fatalf("终点=%v,%v want 10,8", s.lastX, s.lastY)
	}
}

func TestSVGParseArcSamples(t *testing.T) {
	s := &recSink{}
	// 半圆弧应被采样成多段折线（而非退化成一条线）
	parseSVGInto("M0 0 A5 5 0 0 1 10 0", 1, s)
	if s.moves != 1 || s.lines < 4 {
		t.Fatalf("弧采样: moves=%d lines=%d want 1 / >=4", s.moves, s.lines)
	}
}
