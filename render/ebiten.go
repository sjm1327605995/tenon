package render

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log/slog"
	"math"
	"strings"
	"unsafe"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Ebiten struct {
	scaleFactor float32
	Screen      *ebiten.Image
}

var whiteImage *ebiten.Image
var testFont *text.GoTextFaceSource
var fonts1 text.Face

func init() {
	testFont, _ = text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	fonts1 = &text.GoTextFace{
		Source: testFont,
		Size:   24,
	}
	img := ebiten.NewImage(3, 3)
	img.Fill(color.White)
	whiteImage = img.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
}
func (e *Ebiten) Render(r RenderCommand) error {
	fullScreen := e.Screen
	boundingBox := r.BoundingBox
	boundingBox.X *= e.scaleFactor
	boundingBox.Y *= e.scaleFactor
	boundingBox.Width *= e.scaleFactor
	boundingBox.Height *= e.scaleFactor
	switch r.CommandType {
	case RENDER_COMMAND_TYPE_RECTANGLE:
		config := &r.RenderData.Rectangle
		if config.CornerRadius.TopLeft > 0 {
			cornerRadius := config.CornerRadius.TopLeft * e.scaleFactor
			if err := renderFillRoundedRect(e.Screen, boundingBox, cornerRadius, config.BackgroundColor); err != nil {
				return err
			}
		} else {
			vector.DrawFilledRect(
				e.Screen,
				boundingBox.X,
				boundingBox.Y,
				boundingBox.Width, boundingBox.Height,
				config.BackgroundColor,
				true,
			)
		}
	case RENDER_COMMAND_TYPE_TEXT:
		config := &r.RenderData.Text
		cloned := strings.Clone(config.StringContents)
		//TODO
		//font := fonts[config.FontId]
		//
		opts := &text.DrawOptions{}
		opts.ColorScale.Scale(
			float32(config.TextColor.R)/255,
			float32(config.TextColor.G)/255,
			float32(config.TextColor.B)/255,
			float32(config.TextColor.A)/255,
		)
		opts.GeoM.Translate(float64(boundingBox.X), float64(boundingBox.Y))
		//TODO
		text.Draw(e.Screen, cloned, fonts1, opts)
	case RENDER_COMMAND_TYPE_SCISSOR_START:
		e.Screen = e.Screen.SubImage(image.Rect(
			int(boundingBox.X), int(boundingBox.Y),
			int(boundingBox.X+boundingBox.Width),
			int(boundingBox.Y+boundingBox.Height),
		)).(*ebiten.Image)
	case RENDER_COMMAND_TYPE_SCISSOR_END:
		e.Screen = fullScreen
	case RENDER_COMMAND_TYPE_IMAGE:
		config := &r.RenderData.Image
		img := (*ebiten.Image)(config.ImageData.(unsafe.Pointer))
		opts := &ebiten.DrawImageOptions{}
		bounds := img.Bounds()
		opts.GeoM.Scale(float64(boundingBox.Width/float32(bounds.Dx())), float64(boundingBox.Height/float32(bounds.Dy())))
		opts.GeoM.Translate(float64(boundingBox.X), float64(boundingBox.Y))
		e.Screen.DrawImage(img, opts)
	case RENDER_COMMAND_TYPE_BORDER:
		config := &r.RenderData.Border
		config.Width.Top = uint16(float32(config.Width.Top) * e.scaleFactor)
		config.Width.Bottom = uint16(float32(config.Width.Bottom) * e.scaleFactor)
		config.Width.Left = uint16(float32(config.Width.Left) * e.scaleFactor)
		config.Width.Right = uint16(float32(config.Width.Right) * e.scaleFactor)

		config.CornerRadius.TopLeft *= e.scaleFactor
		config.CornerRadius.BottomLeft *= e.scaleFactor
		config.CornerRadius.TopRight *= e.scaleFactor
		config.CornerRadius.BottomRight *= e.scaleFactor
		if boundingBox.Width > 0 && boundingBox.Height > 0 {
			maxRadius := min(boundingBox.Width, boundingBox.Height) / 2.0

			if config.Width.Left > 0 {
				clampedRadiusTop := min(config.CornerRadius.TopLeft, maxRadius)
				clampedRadiusBottom := min(config.CornerRadius.BottomLeft, maxRadius)
				vector.DrawFilledRect(
					e.Screen,
					boundingBox.X,
					boundingBox.Y+clampedRadiusTop,
					float32(config.Width.Left), boundingBox.Height-clampedRadiusTop-clampedRadiusBottom,
					config.Color, true,
				)
			}

			if config.Width.Right > 0 {
				clampedRadiusTop := min(config.CornerRadius.TopRight, maxRadius)
				clampedRadiusBottom := min(config.CornerRadius.BottomRight, maxRadius)
				vector.DrawFilledRect(
					e.Screen,
					boundingBox.X+boundingBox.Width-float32(config.Width.Right),
					boundingBox.Y+clampedRadiusTop,
					float32(config.Width.Right),
					boundingBox.Height-clampedRadiusTop-clampedRadiusBottom,
					config.Color, true,
				)
			}

			if config.Width.Top > 0 {
				clampedRadiusLeft := min(config.CornerRadius.TopLeft, maxRadius)
				clampedRadiusRight := min(config.CornerRadius.TopRight, maxRadius)
				vector.DrawFilledRect(
					e.Screen,
					boundingBox.X+clampedRadiusLeft,
					boundingBox.Y,
					boundingBox.Width-clampedRadiusLeft-clampedRadiusRight,
					float32(config.Width.Top),
					config.Color, true,
				)
			}

			if config.Width.Bottom > 0 {
				clampedRadiusLeft := min(config.CornerRadius.BottomLeft, maxRadius)
				clampedRadiusRight := min(config.CornerRadius.BottomRight, maxRadius)
				vector.DrawFilledRect(
					e.Screen,
					boundingBox.X+clampedRadiusLeft,
					boundingBox.Y+boundingBox.Height-float32(config.Width.Bottom),
					boundingBox.Width-clampedRadiusLeft-clampedRadiusRight,
					float32(config.Width.Bottom),
					config.Color, true,
				)
			}

			// corner index: 0->3 topLeft -> CW -> bottonLeft
			if config.Width.Top > 0 && config.CornerRadius.TopLeft > 0 {
				renderCornerBorder(e.Screen, &boundingBox, config, 0, config.Color)
			}
			if config.Width.Top > 0 && config.CornerRadius.TopRight > 0 {
				renderCornerBorder(e.Screen, &boundingBox, config, 1, config.Color)
			}
			if config.Width.Bottom > 0 && config.CornerRadius.BottomLeft > 0 {
				renderCornerBorder(e.Screen, &boundingBox, config, 2, config.Color)
			}
			if config.Width.Bottom > 0 && config.CornerRadius.BottomLeft > 0 {
				renderCornerBorder(e.Screen, &boundingBox, config, 3, config.Color)
			}
		}
	case RENDER_COMMAND_TYPE_NONE:
	case RENDER_COMMAND_TYPE_CUSTOM:
	default:
		slog.Warn("Unknown command type", "type", r.CommandType)
	}
	return nil
}

const numCircleSegments = 16

func renderFillRoundedRect(screen *ebiten.Image, rect BoundingBox, cornerRadius float32, _color color.RGBA) error {
	r := float32(_color.R) / 255
	g := float32(_color.G) / 255
	b := float32(_color.B) / 255
	a := float32(_color.A) / 255

	indexCount, vertexCount := 0, uint16(0)

	minRadius := min(rect.Width, rect.Height) / 2.0
	clampedRadius := min(cornerRadius, minRadius)

	numCircleSegments := max(numCircleSegments, int(clampedRadius*0.5)) // check if it needs to be clamped

	var vertices [512]ebiten.Vertex
	var indices [512]uint16

	// define center rectangle
	// 0 Center TL
	vertices[vertexCount] = ebiten.Vertex{
		DstX:   rect.X + clampedRadius,
		DstY:   rect.Y + clampedRadius,
		SrcX:   1,
		SrcY:   1,
		ColorR: r,
		ColorG: g,
		ColorB: b,
		ColorA: a,
	}
	vertexCount++
	// 1 Center TR
	vertices[vertexCount] = ebiten.Vertex{
		DstX:   rect.X + rect.Width - clampedRadius,
		DstY:   rect.Y + clampedRadius,
		SrcX:   2,
		SrcY:   1,
		ColorR: r,
		ColorG: g,
		ColorB: b,
		ColorA: a,
	}
	vertexCount++
	// 2 Center BR
	vertices[vertexCount] = ebiten.Vertex{
		DstX:   rect.X + rect.Width - clampedRadius,
		DstY:   rect.Y + rect.Height - clampedRadius,
		SrcX:   2,
		SrcY:   2,
		ColorR: r,
		ColorG: g,
		ColorB: b,
		ColorA: a,
	}
	vertexCount++
	// 3 Center BL
	vertices[vertexCount] = ebiten.Vertex{
		DstX:   rect.X + clampedRadius,
		DstY:   rect.Y + rect.Height - clampedRadius,
		SrcX:   1,
		SrcY:   2,
		ColorR: r,
		ColorG: g,
		ColorB: b,
		ColorA: a,
	}
	vertexCount++

	indices[indexCount] = 0
	indexCount++
	indices[indexCount] = 1
	indexCount++
	indices[indexCount] = 3
	indexCount++
	indices[indexCount] = 1
	indexCount++
	indices[indexCount] = 2
	indexCount++
	indices[indexCount] = 3
	indexCount++

	// define rounded corners as triangle fans
	step := (math.Pi / 2) / float32(numCircleSegments)
	for i := 0; i < numCircleSegments; i++ {
		angle1 := float32(i) * step
		angle2 := (float32(i) + 1) * step

		for j := uint16(0); j < 4; j++ {
			var cx, cy, signX, signY float32

			switch j {
			case 0:
				cx = rect.X + clampedRadius
				cy = rect.Y + clampedRadius
				signX = -1
				signY = -1
			case 1:
				cx = rect.X + rect.Width - clampedRadius
				cy = rect.Y + clampedRadius
				signX = 1
				signY = -1
			case 2:
				cx = rect.X + rect.Width - clampedRadius
				cy = rect.Y + rect.Height - clampedRadius
				signX = 1
				signY = 1
			case 3:
				cx = rect.X + clampedRadius
				cy = rect.Y + rect.Height - clampedRadius
				signX = -1
				signY = 1
			default:
				return fmt.Errorf("out of bounds index: %d", j)
			}

			vertices[vertexCount] = ebiten.Vertex{
				DstX:   cx + float32(math.Cos(float64(angle1)))*clampedRadius*signX,
				DstY:   cy + float32(math.Sin(float64(angle1)))*clampedRadius*signY,
				SrcX:   1,
				SrcY:   1,
				ColorR: r,
				ColorG: g,
				ColorB: b,
				ColorA: a,
			}
			vertexCount++
			vertices[vertexCount] = ebiten.Vertex{
				DstX:   cx + float32(math.Cos(float64(angle2)))*clampedRadius*signX,
				DstY:   cy + float32(math.Sin(float64(angle2)))*clampedRadius*signY,
				SrcX:   1,
				SrcY:   1,
				ColorR: r,
				ColorG: g,
				ColorB: b,
				ColorA: a,
			}
			vertexCount++

			indices[indexCount] = j // Connect to corresponding central rectangle vertex
			indexCount++
			indices[indexCount] = vertexCount - 2
			indexCount++
			indices[indexCount] = vertexCount - 1
			indexCount++
		}
	}
	// Define edge rectangles
	// Top edge
	// TL
	vertices[vertexCount] = ebiten.Vertex{
		DstX:   rect.X + clampedRadius,
		DstY:   rect.Y,
		SrcX:   0,
		SrcY:   0,
		ColorR: r,
		ColorG: g,
		ColorB: b,
		ColorA: a,
	}
	vertexCount++
	// TR
	vertices[vertexCount] = ebiten.Vertex{
		DstX:   rect.X + rect.Width - clampedRadius,
		DstY:   rect.Y,
		SrcX:   1,
		SrcY:   0,
		ColorR: r,
		ColorG: g,
		ColorB: b,
		ColorA: a,
	}
	vertexCount++

	indices[indexCount] = 0
	indexCount++
	indices[indexCount] = vertexCount - 2 // TL
	indexCount++
	indices[indexCount] = vertexCount - 1 // TR
	indexCount++
	indices[indexCount] = 1
	indexCount++
	indices[indexCount] = 0
	indexCount++
	indices[indexCount] = vertexCount - 1 // TR
	indexCount++
	// Right edge
	// RT
	vertices[vertexCount] = ebiten.Vertex{
		DstX:   rect.X + rect.Width,
		DstY:   rect.Y + clampedRadius,
		SrcX:   1,
		SrcY:   0,
		ColorR: r,
		ColorG: g,
		ColorB: b,
		ColorA: a,
	}
	vertexCount++
	// RB
	vertices[vertexCount] = ebiten.Vertex{
		DstX:   rect.X + rect.Width,
		DstY:   rect.Y + rect.Height - clampedRadius,
		SrcX:   1,
		SrcY:   1,
		ColorR: r,
		ColorG: g,
		ColorB: b,
		ColorA: a,
	}
	vertexCount++

	indices[indexCount] = 1
	indexCount++
	indices[indexCount] = vertexCount - 2 // RT
	indexCount++
	indices[indexCount] = vertexCount - 1 // RB
	indexCount++
	indices[indexCount] = 2
	indexCount++
	indices[indexCount] = 1
	indexCount++
	indices[indexCount] = vertexCount - 1 // RB
	indexCount++
	// Bottom edge
	// BR
	vertices[vertexCount] = ebiten.Vertex{
		DstX:   rect.X + rect.Width - clampedRadius,
		DstY:   rect.Y + rect.Height,
		SrcX:   1,
		SrcY:   1,
		ColorR: r,
		ColorG: g,
		ColorB: b,
		ColorA: a,
	}
	vertexCount++
	// BL
	vertices[vertexCount] = ebiten.Vertex{
		DstX:   rect.X + clampedRadius,
		DstY:   rect.Y + rect.Height,
		SrcX:   0,
		SrcY:   1,
		ColorR: r,
		ColorG: g,
		ColorB: b,
		ColorA: a,
	}
	vertexCount++

	indices[indexCount] = 2
	indexCount++
	indices[indexCount] = vertexCount - 2 // BR
	indexCount++
	indices[indexCount] = vertexCount - 1 // BL
	indexCount++
	indices[indexCount] = 3
	indexCount++
	indices[indexCount] = 2
	indexCount++
	indices[indexCount] = vertexCount - 1 // BL
	indexCount++
	// Left edge
	// LB
	vertices[vertexCount] = ebiten.Vertex{
		DstX:   rect.X,
		DstY:   rect.Y + rect.Height - clampedRadius,
		SrcX:   0,
		SrcY:   1,
		ColorR: r,
		ColorG: g,
		ColorB: b,
		ColorA: a,
	}
	vertexCount++
	// LT
	vertices[vertexCount] = ebiten.Vertex{
		DstX:   rect.X,
		DstY:   rect.Y + clampedRadius,
		SrcX:   0,
		SrcY:   0,
		ColorR: r,
		ColorG: g,
		ColorB: b,
		ColorA: a,
	}
	vertexCount++

	indices[indexCount] = 3
	indexCount++
	indices[indexCount] = vertexCount - 2 // LB
	indexCount++
	indices[indexCount] = vertexCount - 1 // LT
	indexCount++
	indices[indexCount] = 0
	indexCount++
	indices[indexCount] = 3
	indexCount++
	indices[indexCount] = vertexCount - 1 // LT
	indexCount++

	// Render everything
	screen.DrawTriangles(vertices[:vertexCount], indices[:indexCount], whiteImage, &ebiten.DrawTrianglesOptions{
		AntiAlias: true,
	})

	return nil
}

// all rendering is performed by a single ebitengine call, using two sets of arcing triangles, inner and outer, that fit together; along with two triangles to fill the end gaps.
func renderCornerBorder(screen *ebiten.Image, boundingBox *BoundingBox, config *BorderRenderData, cornerIndex int, _color color.RGBA) {
	/////////////////////////////////
	//The arc is constructed of outer triangles and inner triangles (if needed).
	//First three vertices are first outer triangle's vertices
	//Each two vertices after that are the inner-middle and second-outer vertex of
	//each outer triangle after the first, because there first-outer vertex is equal to the
	//second-outer vertex of the previous triangle. Indices set accordingly.
	//The final two vertices are the missing vertices for the first and last inner triangles (if needed)
	//Everything is in clockwise order (CW).
	/////////////////////////////////

	r := float32(_color.R) / 255
	g := float32(_color.G) / 255
	b := float32(_color.B) / 255
	a := float32(_color.A) / 255

	var centerX, centerY, outerRadius, startAngle, borderWidth float32
	maxRadius := min(boundingBox.Width, boundingBox.Height) / 2.0

	var vertices [512]ebiten.Vertex
	var indices [512]uint16
	indexCount, vertexCount := uint16(0), uint16(0)

	switch cornerIndex {
	case 0:
		startAngle = math.Pi
		outerRadius = min(config.CornerRadius.TopLeft, maxRadius)
		centerX = boundingBox.X + outerRadius
		centerY = boundingBox.Y + outerRadius
		borderWidth = float32(config.Width.Top)
	case 1:
		startAngle = 3 * math.Pi / 2
		outerRadius = min(config.CornerRadius.TopRight, maxRadius)
		centerX = boundingBox.X + boundingBox.Width - outerRadius
		centerY = boundingBox.Y + outerRadius
		borderWidth = float32(config.Width.Top)
	case 2:
		startAngle = 0
		outerRadius = min(config.CornerRadius.BottomRight, maxRadius)
		centerX = boundingBox.X + boundingBox.Width - outerRadius
		centerY = boundingBox.Y + boundingBox.Height - outerRadius
		borderWidth = float32(config.Width.Bottom)
	case 3:
		startAngle = math.Pi / 2
		outerRadius = min(config.CornerRadius.BottomLeft, maxRadius)
		centerX = boundingBox.X + outerRadius
		centerY = boundingBox.Y + boundingBox.Height - outerRadius
		borderWidth = float32(config.Width.Bottom)
	default:
		panic("invalid corner index")
	}

	innerRadius := outerRadius - borderWidth
	minNumOuterTriangles := numCircleSegments
	numOuterTriangles := max(minNumOuterTriangles, int(math.Ceil(float64(outerRadius*0.5))))
	angleStep := math.Pi / (2.0 * float32(numOuterTriangles))

	// outer triangles, in CW order
	for i := 0; i < numOuterTriangles; i++ {
		angle1 := startAngle + float32(i)*angleStep       // first-outer vertex angle
		angle2 := startAngle + (float32(i)+0.5)*angleStep // inner-middle vertex angle
		angle3 := startAngle + float32(i+1)*angleStep     // second-outer vertex angle

		if i == 0 { // first outer triangle
			vertices[vertexCount] = ebiten.Vertex{
				DstX:   centerX + float32(math.Cos(float64(angle1)))*outerRadius,
				DstY:   centerY + float32(math.Sin(float64(angle1)))*outerRadius,
				SrcX:   0,
				SrcY:   0,
				ColorR: r,
				ColorG: g,
				ColorB: b,
				ColorA: a,
			} // vertex index = 0
			vertexCount++
		}
		indices[indexCount] = vertexCount - 1 // will be second-outer vertex of last outer triangle if not first outer triangle.
		indexCount++
		if innerRadius > 0 {
			vertices[vertexCount] = ebiten.Vertex{
				DstX:   centerX + float32(math.Cos(float64(angle2)))*(innerRadius),
				DstY:   centerY + float32(math.Sin(float64(angle2)))*(innerRadius),
				SrcX:   0,
				SrcY:   0,
				ColorR: r,
				ColorG: g,
				ColorB: b,
				ColorA: a,
			}
			vertexCount++
		} else {
			vertices[vertexCount] = ebiten.Vertex{
				DstX:   centerX,
				DstY:   centerY,
				SrcX:   0,
				SrcY:   0,
				ColorR: r,
				ColorG: g,
				ColorB: b,
				ColorA: a,
			}
			vertexCount++
		}

		indices[indexCount] = vertexCount - 1
		indexCount++
		vertices[vertexCount] = ebiten.Vertex{
			DstX:   centerX + float32(math.Cos(float64(angle3)))*outerRadius,
			DstY:   centerY + float32(math.Sin(float64(angle3)))*outerRadius,
			SrcX:   0,
			SrcY:   0,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		}
		vertexCount++
		indices[indexCount] = vertexCount - 1
		indexCount++
	}

	if innerRadius > 0 {
		// inner triangles in CW order (except the first and last)
		for i := 0; i < numOuterTriangles-1; i++ { // skip the last outer triangle
			if i == 0 { // first outer triangle -> second inner triangle
				indices[indexCount] = 1 // inner-middle vertex of first outer triangle
				indexCount++
				indices[indexCount] = 2 // second-outer vertex of first outer triangle
				indexCount++
				indices[indexCount] = 3 // innder-middle vertex of second-outer triangle
				indexCount++
			} else {
				baseIndex := 3                                    // skip first outer triangle
				indices[indexCount] = uint16(baseIndex + (i-1)*2) // inner-middle vertex of current outer triangle
				indexCount++
				indices[indexCount] = uint16(baseIndex + (i-1)*2 + 1) // second-outer vertex of current outer triangle
				indexCount++
				indices[indexCount] = uint16(baseIndex + (i-1)*2 + 2) // inner-middle vertex of next outer triangle
				indexCount++
			}
		}

		endAngle := startAngle + math.Pi/2.0

		// last inner triangle
		indices[indexCount] = vertexCount - 2 // inner-middle vertex of last outer triangle
		indexCount++
		indices[indexCount] = vertexCount - 1 // second-outer vertex of last outer triangle
		indexCount++
		vertices[vertexCount] = ebiten.Vertex{
			DstX:   centerX + float32(math.Cos(float64(endAngle)))*innerRadius,
			DstY:   centerY + float32(math.Sin(float64(endAngle)))*innerRadius,
			SrcX:   0,
			SrcY:   0,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		} // missing vertex
		vertexCount++
		indices[indexCount] = vertexCount - 1
		indexCount++

		// //first inner triangle
		indices[indexCount] = 0 // first-outer vertex of first outer triangle
		indexCount++
		indices[indexCount] = 1 // inner-middle vertex of first outer triangle
		indexCount++
		vertices[vertexCount] = ebiten.Vertex{
			DstX:   centerX + float32(math.Cos(float64(startAngle)))*innerRadius,
			DstY:   centerY + float32(math.Sin(float64(startAngle)))*innerRadius,
			SrcX:   0,
			SrcY:   0,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		} // missing vertex
		vertexCount++
		indices[indexCount] = vertexCount - 1
		indexCount++
	}

	screen.DrawTriangles(vertices[:vertexCount], indices[:indexCount], whiteImage, &ebiten.DrawTrianglesOptions{
		AntiAlias: true,
	})
}
