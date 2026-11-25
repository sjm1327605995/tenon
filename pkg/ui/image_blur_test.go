package ui

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestBlurredImageMixesEdgeColors(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 3, 1))
	img.SetNRGBA(0, 0, color.NRGBA{A: 255})
	img.SetNRGBA(1, 0, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
	img.SetNRGBA(2, 0, color.NRGBA{A: 255})

	blurred := blurredImage(img, 1)
	for x, want := range []uint8{127, 85, 127} {
		got := color.NRGBAModel.Convert(blurred.At(x, 0)).(color.NRGBA)
		if got.R != want || got.G != want || got.B != want || got.A != 255 {
			t.Errorf("pixel %d = %#v, want gray %d with full alpha", x, got, want)
		}
	}
	if middle := color.NRGBAModel.Convert(img.At(1, 0)).(color.NRGBA); middle.R != 255 {
		t.Fatal("blur modified the source image")
	}
}

func TestImageBlurUsesDistinctCacheEntries(t *testing.T) {
	resetImageCache(t)
	img := image.NewNRGBA(image.Rect(0, 0, 3, 1))
	img.SetNRGBA(0, 0, color.NRGBA{A: 255})
	img.SetNRGBA(1, 0, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
	img.SetNRGBA(2, 0, color.NRGBA{A: 255})

	h := Mount(Div(
		Img(SrcImage("mem:blur-cache", img)),
		Img(SrcImage("mem:blur-cache", img), Style(ImageBlur(1))),
		Img(SrcImage("mem:blur-cache", img), Style(ImageBlur(3))),
	), 100, 100)
	images := h.Root().FindAll(func(q *Query) bool { return q.rn.kind == rnImage })
	if len(images) != 3 {
		t.Fatalf("found %d image nodes, want 3", len(images))
	}
	if images[0].rn.img == images[1].rn.img || images[1].rn.img == images[2].rn.img {
		t.Fatalf("original and blurred images share a bitmap; radii = %v/%v/%v", images[0].rn.imgBlur, images[1].rn.imgBlur, images[2].rn.imgBlur)
	}
	if n, _, _ := ImageCacheStats(); n != 3 {
		t.Fatalf("cache entries = %d, want original plus two blur radii", n)
	}
	middle := color.NRGBAModel.Convert(images[1].rn.img.(*gioImage).src.At(1, 0)).(color.NRGBA)
	if middle.R != 85 {
		t.Fatalf("blurred middle pixel = %d, want 85", middle.R)
	}
}

func TestImageBlurChangeReloadsBitmap(t *testing.T) {
	resetImageCache(t)
	img := image.NewNRGBA(image.Rect(0, 0, 3, 1))
	img.SetNRGBA(0, 0, color.NRGBA{A: 255})
	img.SetNRGBA(1, 0, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
	img.SetNRGBA(2, 0, color.NRGBA{A: 255})
	app := func(_ struct{}) *Node {
		blurred, setBlurred := UseState(false)
		radius := float32(0)
		if blurred {
			radius = 1
		}
		return Div(
			Img(SrcImage("mem:blur-change", img), Style(ImageBlur(radius))),
			Button(Text("toggle"), OnClick(func() { setBlurred(true) })),
		)
	}
	h := Mount(Use(app, struct{}{}), 100, 100)
	original := h.Root().ByKind("image").rn.img
	if original == nil {
		t.Fatal("original image was not installed")
	}

	bounds := h.Root().ByText("toggle").Bounds()
	h.ClickAt(bounds.X+bounds.W/2, bounds.Y+bounds.H/2)
	changed := h.Root().ByKind("image").rn.img
	if changed == nil || changed == original {
		t.Fatal("changing blur radius did not replace the bitmap")
	}
	middle := color.NRGBAModel.Convert(changed.(*gioImage).src.At(1, 0)).(color.NRGBA)
	if middle.R != 85 {
		t.Fatalf("middle pixel after enabling blur = %d, want 85", middle.R)
	}
	if _, ok := lookupImage("mem:blur-change"); !ok {
		t.Fatal("original cache entry was not retained")
	}
	if _, ok := lookupImage("mem:blur-change\x00blur=1"); !ok {
		t.Fatal("blurred cache entry is missing")
	}
}

func TestAsyncImageBlurDoesNotOverwriteNewRadius(t *testing.T) {
	resetImageCache(t)
	source := image.NewNRGBA(image.Rect(0, 0, 3, 1))
	source.SetNRGBA(0, 0, color.NRGBA{A: 255})
	source.SetNRGBA(1, 0, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
	source.SetNRGBA(2, 0, color.NRGBA{A: 255})
	var encoded bytes.Buffer
	if err := png.Encode(&encoded, source); err != nil {
		t.Fatal(err)
	}

	var requests atomic.Int32
	firstStarted := make(chan struct{})
	secondStarted := make(chan struct{})
	releaseFirst := make(chan struct{})
	var releaseOnce sync.Once
	release := func() { releaseOnce.Do(func() { close(releaseFirst) }) }
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if requests.Add(1) == 1 {
			close(firstStarted)
			<-releaseFirst
		} else {
			close(secondStarted)
		}
		_, _ = w.Write(encoded.Bytes())
	}))
	defer server.Close()
	defer release()

	rn := newImageRenderNode()
	rn.imgSrc, rn.imgBlur = server.URL, 1
	rn.loadImage()
	waitForImageRequest(t, firstStarted)
	oldKey := rn.imageKey()

	rn.imgBlur = 3
	rn.loadImage()
	newKey := rn.imageKey()
	waitForImageRequest(t, secondStarted)
	waitForImageLoad(t, func() bool { return !imgLoading[newKey] })
	newBitmap := rn.img
	if newBitmap == nil {
		t.Fatal("new blur radius did not install a bitmap")
	}

	release()
	waitForImageLoad(t, func() bool { return !imgLoading[oldKey] })
	if rn.img != newBitmap {
		t.Fatal("late result for the old radius overwrote the current bitmap")
	}
	if _, ok := lookupImage(oldKey); !ok {
		t.Fatal("old radius result was not stored under its own key")
	}
	if _, ok := lookupImage(newKey); !ok {
		t.Fatal("new radius result was not stored under its own key")
	}
}

func waitForImageRequest(t *testing.T, started <-chan struct{}) {
	t.Helper()
	select {
	case <-started:
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for image request")
	}
}

func waitForImageLoad(t *testing.T, done func() bool) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for !done() && time.Now().Before(deadline) {
		drainPosts()
		time.Sleep(5 * time.Millisecond)
	}
	drainPosts()
	if !done() {
		t.Fatal("timed out waiting for image load")
	}
}

func TestAsyncImageClearIgnoresLateResult(t *testing.T) {
	resetImageCache(t)
	var encoded bytes.Buffer
	if err := png.Encode(&encoded, solidImage(2, 2, color.RGBA{R: 255, A: 255})); err != nil {
		t.Fatal(err)
	}
	started := make(chan struct{})
	releaseLoad := make(chan struct{})
	var releaseOnce sync.Once
	release := func() { releaseOnce.Do(func() { close(releaseLoad) }) }
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		close(started)
		<-releaseLoad
		_, _ = w.Write(encoded.Bytes())
	}))
	defer server.Close()
	defer release()

	rn := newImageRenderNode()
	rn.imgSrc = server.URL
	rn.loadImage()
	waitForImageRequest(t, started)
	key := rn.imageKey()

	applyHostProps(&Fiber{}, rn, &hostProps{style: newStyleProps()})
	if rn.img != nil || rn.imgSrc != "" {
		t.Fatal("clearing the source did not unload the pending image")
	}
	release()
	waitForImageLoad(t, func() bool { return !imgLoading[key] })
	if rn.img != nil || rn.imgSrc != "" {
		t.Fatal("late asynchronous result resurrected a cleared image")
	}
	if imgLoading[""] {
		t.Fatal("empty image source started a background load")
	}
	if _, ok := lookupImage(key); !ok {
		t.Fatal("completed source result was not stored under its original key")
	}
}
