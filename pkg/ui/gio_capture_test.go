package ui

import "testing"

func TestCaptureSchedulingAndCancellation(t *testing.T) {
	previous := pendingCapture
	t.Cleanup(func() { pendingCapture = previous })

	Capture("frame.png", 2, true)
	if pendingCapture == nil || pendingCapture.remaining != 2 || !pendingCapture.exit {
		t.Fatalf("capture config=%+v, want two remaining frames and exit", pendingCapture)
	}
	if captureGioFrame(nil) || pendingCapture.remaining != 1 {
		t.Fatal("first waiting frame did not decrement the capture countdown")
	}
	if captureGioFrame(nil) || pendingCapture.remaining != 0 {
		t.Fatal("second waiting frame did not reach the capture frame")
	}
	Capture("", 0, false)
	if pendingCapture != nil {
		t.Fatal("empty capture path did not cancel the request")
	}
	Capture("frame.png", -1, false)
	if pendingCapture == nil || pendingCapture.remaining != 0 || pendingCapture.exit {
		t.Fatalf("normalized capture config=%+v, want immediate capture without exit", pendingCapture)
	}
}
