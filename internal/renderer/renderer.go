package renderer

import (
	"fmt"
	"strings"

	"github.com/sjm1327605995/tenon/pkg/types"
)

type Renderer struct {
	width  int
	height int
	buffer [][]rune
}

func NewRenderer(width, height int) *Renderer {
	buffer := make([][]rune, height)
	for i := range buffer {
		buffer[i] = make([]rune, width)
		for j := range buffer[i] {
			buffer[i][j] = ' '
		}
	}

	return &Renderer{
		width:  width,
		height: height,
		buffer: buffer,
	}
}

func (r *Renderer) RenderElement(element types.Element, x, y int) {
	if element == nil {
		return
	}

	layout := element.GetLayout()
	props := element.GetProps()

	switch props.(type) {
	case *types.ViewProps:
		r.renderView(element, x, y)
	case *types.TextProps:
		r.renderText(element, x, y)
	case *types.ImageProps:
		r.renderImage(element, x, y)
	}

	childX := x + int(layout.X)
	childY := y + int(layout.Y)
	for _, child := range element.GetChildren() {
		r.RenderElement(child, childX, childY)
	}
}

func (r *Renderer) renderView(element types.Element, x, y int) {
	layout := element.GetLayout()
	startX := x + int(layout.X)
	startY := y + int(layout.Y)
	width := int(layout.Width)
	height := int(layout.Height)

	if startX < 0 || startY < 0 || width <= 0 || height <= 0 {
		return
	}

	r.drawRect(startX, startY, width, height, '█')

	innerX := startX + 1
	innerY := startY + 1
	innerWidth := width - 2
	innerHeight := height - 2

	if innerWidth > 0 && innerHeight > 0 {
		for i := 0; i < innerHeight; i++ {
			currentY := innerY + i
			if currentY >= 0 && currentY < len(r.buffer) {
				for j := 0; j < innerWidth; j++ {
					currentX := innerX + j
					if currentX >= 0 && currentX < len(r.buffer[currentY]) {
						r.buffer[currentY][currentX] = '░'
					}
				}
			}
		}
	}
}

func (r *Renderer) renderText(element types.Element, x, y int) {
	layout := element.GetLayout()
	startX := x + int(layout.X)
	startY := y + int(layout.Y)

	var content string
	if props, ok := element.GetProps().(*types.TextProps); ok {
		content = props.Content
	} else {
		content = "Text"
	}

	for i, char := range content {
		if startX+i < r.width && startY < r.height {
			r.buffer[startY][startX+i] = char
		}
	}
}

func (r *Renderer) renderImage(element types.Element, x, y int) {
	layout := element.GetLayout()
	startX := x + int(layout.X)
	startY := y + int(layout.Y)
	width := int(layout.Width)
	height := int(layout.Height)

	if startX < 0 || startY < 0 || width <= 0 || height <= 0 {
		return
	}

	r.drawRect(startX, startY, width, height, '▒')
}

func (r *Renderer) drawRect(x, y, w, h int, ch rune) {
	if x < 0 || y < 0 || w <= 0 || h <= 0 {
		return
	}

	for i := 0; i < w && x+i < r.width; i++ {
		if y >= 0 && y < r.height {
			r.buffer[y][x+i] = ch
		}
		if y+h-1 >= 0 && y+h-1 < r.height && h > 1 {
			r.buffer[y+h-1][x+i] = ch
		}
	}

	for i := 0; i < h && y+i < r.height; i++ {
		if x >= 0 && x < r.width {
			r.buffer[y+i][x] = ch
		}
		if x+w-1 >= 0 && x+w-1 < r.width && w > 1 {
			r.buffer[y+i][x+w-1] = ch
		}
	}
}

func (r *Renderer) Render(element types.Element) string {
	r.clearBuffer()

	if element != nil {
		r.RenderElement(element, 0, 0)
	}

	var sb strings.Builder
	for _, row := range r.buffer {
		sb.WriteString(string(row))
		sb.WriteString("\n")
	}
	return sb.String()
}

func (r *Renderer) clearBuffer() {
	for i := range r.buffer {
		for j := range r.buffer[i] {
			r.buffer[i][j] = ' '
		}
	}
}

func (r *Renderer) Print(element types.Element) {
	fmt.Print(r.Render(element))
}