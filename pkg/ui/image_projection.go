package ui

import "image"

// WarpImageToQuad maps an image to an arbitrary projected quadrilateral using
// the same homography as the Gio backend. The returned bounds locate the warped
// image on screen. Corners are ordered top-left, top-right, bottom-left, bottom-right.
func WarpImageToQuad(source image.Image, quad [4]Point) (image.Image, Rect) {
	if source == nil {
		return nil, Rect{}
	}
	box := quadBBox(quad[0], quad[1], quad[2], quad[3])
	if box.W <= 0 || box.H <= 0 {
		return nil, Rect{}
	}
	img := warpQuad(source, quad[0], quad[1], quad[2], quad[3], box)
	if img == nil {
		return nil, Rect{}
	}
	return img, box
}
