package ui

import (
	"testing"

	"github.com/sjm1327605995/tenon/yoga"
)

type bidiTestFace struct{}

func (bidiTestFace) Measure(s string, _ float64) float32 {
	return float32(len([]rune(s))) * 10
}
func (bidiTestFace) Metrics() (float32, bool, bool) { return 10, false, false }

func TestBidiDisplayAndLogicalVisualMapping(t *testing.T) {
	face := bidiTestFace{}
	if got := bidiDisplay("abc אבג def", face, 13); got != "abc גבא def" {
		t.Fatalf("mixed display=%q, want %q", got, "abc גבא def")
	}
	if got := bidiDisplay("אבג", face, 13); got != "גבא" {
		t.Fatalf("RTL display=%q, want גבא", got)
	}

	const text = "אבג"
	for _, test := range []struct {
		offset int
		wantX  float32
	}{
		{offset: 0, wantX: 30},
		{offset: 2, wantX: 20},
		{offset: 4, wantX: 10},
		{offset: 6, wantX: 0},
	} {
		if got := bidiCaretX(text, test.offset, face, 13, 0); got != test.wantX {
			t.Fatalf("caret offset %d x=%v, want %v", test.offset, got, test.wantX)
		}
		if got := bidiOffset(text, test.wantX, face, 13); got != test.offset {
			t.Fatalf("x %v offset=%d, want %d", test.wantX, got, test.offset)
		}
	}
}

func TestBidiExtendedScriptsAndBracketPairs(t *testing.T) {
	face := bidiTestFace{}
	for _, test := range []struct{ input, want string }{
		{input: "abc أبجد def", want: "abc دجبأ def"},
		{input: "abc 123 אבג", want: "abc 123 גבא"},
		{input: "abc (אב) def", want: "abc (בא) def"},
	} {
		if got := bidiDisplay(test.input, face, 13); got != test.want {
			t.Fatalf("display %q=%q, want %q", test.input, got, test.want)
		}
	}
}

func TestBidiSelectionUsesVisualRanges(t *testing.T) {
	face := bidiTestFace{}
	ranges := bidiRangeX("abc אבג", 3, 6, face, 13, 0)
	want := [][2]float32{{30, 40}, {60, 70}}
	if len(ranges) != len(want) {
		t.Fatalf("selection ranges=%v, want %v", ranges, want)
	}
	for i := range want {
		if ranges[i] != want[i] {
			t.Fatalf("selection ranges=%v, want %v", ranges, want)
		}
	}
}

func TestRichTextBidiPreservesStyledRuns(t *testing.T) {
	face := bidiTestFace{}
	runs := []textRun{
		{text: "abc ", face: face, lineH: 13, ascent: 10},
		{text: "אבג", face: face, lineH: 13, ascent: 10},
		{text: " def", face: face, lineH: 13, ascent: 10},
	}
	lines, maxWidth, _ := layoutRuns(runs, 0)
	if len(lines) != 1 || len(lines[0].segs) != 4 {
		t.Fatalf("rich layout=%+v", lines)
	}
	segments := lines[0].segs
	for i, want := range []struct {
		text string
		run  int
		x    float32
	}{{"abc ", 0, 0}, {"גבא", 1, 40}, {" ", 2, 70}, {"def", 2, 80}} {
		if segments[i].text != want.text || segments[i].run != want.run || segments[i].x != want.x {
			t.Fatalf("segment %d=%+v, want %+v", i, segments[i], want)
		}
	}
	if lines[0].width != 110 || maxWidth != 110 {
		t.Fatalf("rich width=%v max=%v, want 110", lines[0].width, maxWidth)
	}
}

func TestRichTextBidiSplitsStylesInVisualOrder(t *testing.T) {
	face := bidiTestFace{}
	runs := []textRun{
		{text: "א", face: face, lineH: 13, ascent: 10},
		{text: "בג", face: face, lineH: 13, ascent: 10},
	}
	lines, _, _ := layoutRuns(runs, 0)
	segments := lines[0].segs
	if len(segments) != 2 || segments[0].text != "גב" || segments[0].run != 1 || segments[0].x != 0 ||
		segments[1].text != "א" || segments[1].run != 0 || segments[1].x != 20 {
		t.Fatalf("styled RTL segments=%+v", segments)
	}
}

func TestBidiPlainTextAndInputPainting(t *testing.T) {
	face := bidiTestFace{}
	text := &renderNode{kind: rnText, text: "abc אבג def", face: face, lineH: 13, bounds: Rect{W: 200, H: 20}}
	painter := &recordPainter{}
	paintSelf(painter, text)
	if len(painter.ops) != 1 || painter.ops[0].Text != "abc גבא def" {
		t.Fatalf("plain text paint=%+v", painter.ops)
	}

	input := &renderNode{
		kind: rnInput, value: "abc אבג", face: face, lineH: 13,
		bounds: Rect{W: 200, H: 30}, yn: newBoxRenderNode().yn,
		selAnchor: 3, caretPos: 6,
	}
	input.yn.StyleSetPadding(yoga.EdgeAll, 0)
	input.yn.CalculateLayout(200, 30, yoga.DirectionLTR)
	painter = &recordPainter{}
	paintInput(painter, input)
	var selection []PaintOp
	for _, op := range painter.ops {
		if op.Kind == "rect" && op.Color.A == 90 {
			selection = append(selection, op)
		}
	}
	if len(selection) != 2 || selection[0].Rect.X != 30 || selection[1].Rect.X != 60 {
		t.Fatalf("input selection paint=%+v", selection)
	}
}

func TestBidiLogicalVisualRoundTripWithinRuns(t *testing.T) {
	face := bidiTestFace{}
	texts := []string{
		"abc אבג def",
		"abc أبجد def",
		"abc (אב) 123 أبجد def",
		"אבג",
		"أبجد",
		"English 中文 עברית العربية",
	}
	for _, text := range texts {
		for _, run := range bidiRuns(text, face, 13, 0) {
			// A shared visual boundary can be a dual caret; verify interior offsets.
			for offset := run.start + 1; offset < run.end; offset++ {
				if !isRuneStart(text[offset]) {
					continue
				}
				x := bidiCaretX(text, offset, face, 13, 0)
				if got := bidiOffset(text, x, face, 13); got != offset {
					t.Fatalf("round trip %q offset=%d x=%v produced %d", text, offset, x, got)
				}
			}
		}
	}
}

func TestBidiHasRTLDetection(t *testing.T) {
	for _, test := range []struct {
		text string
		want bool
	}{
		{text: "Plain English 中文 123", want: false},
		{text: "abc אבג", want: true},
		{text: "abc أبجد", want: true},
		{text: "abc \u202Edef", want: true},
	} {
		if got := bidiHasRTL(test.text); got != test.want {
			t.Fatalf("bidiHasRTL(%q)=%t, want %t", test.text, got, test.want)
		}
	}
}
