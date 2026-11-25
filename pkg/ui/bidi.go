package ui

import "golang.org/x/text/unicode/bidi"

func bidiHasRTL(s string) bool {
	for _, r := range s {
		properties, _ := bidi.LookupRune(r)
		switch properties.Class() {
		case bidi.R, bidi.AL, bidi.RLE, bidi.RLO, bidi.RLI:
			return true
		}
	}
	return false
}

type bidiVisualRun struct {
	logical    string
	text       string
	start, end int
	rtl        bool
	x, width   float32
}

func bidiRuns(s string, face fontFace, lineH float64, tx float32) []bidiVisualRun {
	if s == "" {
		return nil
	}
	if !bidiHasRTL(s) {
		return []bidiVisualRun{{logical: s, text: s, end: len(s), x: tx, width: measureW(s, face, lineH)}}
	}
	var paragraph bidi.Paragraph
	if _, err := paragraph.SetString(s); err != nil {
		return []bidiVisualRun{{logical: s, text: s, end: len(s), x: tx, width: measureW(s, face, lineH)}}
	}
	ordering, err := paragraph.Order()
	if err != nil || ordering.NumRuns() == 0 {
		return []bidiVisualRun{{logical: s, text: s, end: len(s), x: tx, width: measureW(s, face, lineH)}}
	}
	allLTR := true
	for i := 0; i < ordering.NumRuns(); i++ {
		run := ordering.Run(i)
		if run.Direction() == bidi.RightToLeft {
			allLTR = false
			break
		}
	}
	if allLTR {
		return []bidiVisualRun{{logical: s, text: s, end: len(s), x: tx, width: measureW(s, face, lineH)}}
	}
	byteOffsets := make([]int, 0, len([]rune(s))+1)
	for i := range s {
		byteOffsets = append(byteOffsets, i)
	}
	byteOffsets = append(byteOffsets, len(s))

	runs := make([]bidiVisualRun, 0, ordering.NumRuns())
	x := tx
	for i := 0; i < ordering.NumRuns(); i++ {
		run := ordering.Run(i)
		startRune, endRune := run.Pos()
		start, end := byteOffsets[startRune], byteOffsets[endRune+1]
		text := run.String()
		rtl := run.Direction() == bidi.RightToLeft
		if rtl {
			text = string(bidi.AppendReverse(nil, run.Bytes()))
		}
		width := measureW(text, face, lineH)
		runs = append(runs, bidiVisualRun{logical: s[start:end], text: text, start: start, end: end, rtl: rtl, x: x, width: width})
		x += width
	}
	return runs
}

func (rn *renderNode) displayText(s string) string {
	if rn.bidiValid && rn.bidiText == s {
		return rn.bidiVisual
	}
	rn.bidiText = s
	rn.bidiVisual = bidiDisplay(s, rn.face, rn.lineH)
	rn.bidiValid = true
	return rn.bidiVisual
}

func bidiDisplay(s string, face fontFace, lineH float64) string {
	runs := bidiRuns(s, face, lineH, 0)
	if len(runs) == 1 && !runs[0].rtl {
		return s
	}
	display := ""
	for _, run := range runs {
		display += run.text
	}
	return display
}

func bidiXInRun(run bidiVisualRun, off int, face fontFace, lineH float64) float32 {
	off = clampi(off, run.start, run.end)
	prefix := measureW(runText(run, run.start, off), face, lineH)
	if run.rtl {
		return run.x + run.width - prefix
	}
	return run.x + prefix
}

func runText(run bidiVisualRun, start, end int) string {
	return runTextAt(run, start, end)
}

func runTextAt(run bidiVisualRun, start, end int) string {
	return run.logical[start-run.start : end-run.start]
}

func bidiCaretX(s string, off int, face fontFace, lineH float64, tx float32) float32 {
	runs := bidiRuns(s, face, lineH, tx)
	for _, run := range runs {
		if off >= run.start && off <= run.end {
			return bidiXInRun(run, off, face, lineH)
		}
	}
	return tx
}

func bidiOffset(s string, relX float32, face fontFace, lineH float64) int {
	runs := bidiRuns(s, face, lineH, 0)
	if len(runs) == 0 {
		return 0
	}
	first := runs[0]
	if relX <= first.x {
		if first.rtl {
			return first.end
		}
		return first.start
	}
	last := runs[len(runs)-1]
	if relX >= last.x+last.width {
		if last.rtl {
			return last.start
		}
		return last.end
	}
	for _, run := range runs {
		if relX < run.x || relX > run.x+run.width {
			continue
		}
		x := relX - run.x
		if run.rtl {
			x = run.width - x
		}
		return run.start + bidiPrefixOffset(runTextAt(run, run.start, run.end), x, face, lineH)
	}
	return len(s)
}

func bidiPrefixOffset(s string, x float32, face fontFace, lineH float64) int {
	if x <= 0 || s == "" {
		return 0
	}
	previousIndex, previousWidth := 0, float32(0)
	for i := 1; i <= len(s); i++ {
		if i < len(s) && !isRuneStart(s[i]) {
			continue
		}
		width := measureW(s[:i], face, lineH)
		if width >= x {
			if width-x < x-previousWidth {
				return i
			}
			return previousIndex
		}
		previousIndex, previousWidth = i, width
	}
	return len(s)
}

func bidiRangeX(s string, lo, hi int, face fontFace, lineH float64, tx float32) [][2]float32 {
	var ranges [][2]float32
	for _, run := range bidiRuns(s, face, lineH, tx) {
		a, b := max(lo, run.start), min(hi, run.end)
		if a >= b {
			continue
		}
		x0 := bidiXInRun(run, a, face, lineH)
		x1 := bidiXInRun(run, b, face, lineH)
		if x0 > x1 {
			x0, x1 = x1, x0
		}
		ranges = append(ranges, [2]float32{x0, x1})
	}
	return ranges
}
