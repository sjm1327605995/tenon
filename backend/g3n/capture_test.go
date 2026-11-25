package g3n

import "testing"

func TestNativeCaptureFromEnv(t *testing.T) {
	t.Setenv("TENON_CAPTURE", "frame.png")
	t.Setenv("TENON_CAPTURE_FRAMES", "2")
	capture := nativeCaptureFromEnv()
	if capture == nil || capture.path != "frame.png" || capture.remaining != 2 {
		t.Fatalf("capture config = %+v", capture)
	}
	if capture.capture(nil, 0, 0) || capture.remaining != 1 {
		t.Fatal("capture countdown did not advance")
	}
}

func TestNativeCaptureDisabled(t *testing.T) {
	t.Setenv("TENON_CAPTURE", "")
	if capture := nativeCaptureFromEnv(); capture != nil {
		t.Fatal("empty capture path produced a capture config")
	}
	var capture *nativeCapture
	if capture.capture(nil, 0, 0) {
		t.Fatal("nil capture requested a frame")
	}
}
