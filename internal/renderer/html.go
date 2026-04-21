package renderer

import (
	"fmt"
	"strings"

	"github.com/sjm1327605995/tenon/pkg/types"
)

type HTMLRenderer struct {
	builder strings.Builder
	indent  int
}

func NewHTMLRenderer() *HTMLRenderer {
	return &HTMLRenderer{
		builder: strings.Builder{},
	}
}

func (r *HTMLRenderer) Render(element types.Element) string {
	r.builder.Reset()
	r.writeHeader()
	if element != nil {
		r.renderElement(element)
	}
	r.writeFooter()
	return r.builder.String()
}

func (r *HTMLRenderer) writeHeader() {
	r.builder.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Tenon Debug</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            padding: 20px;
            background: #f5f5f5;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
        }
        h1 {
            color: #333;
            margin-bottom: 20px;
            font-size: 24px;
        }
        .debug-panel {
            background: white;
            border-radius: 8px;
            padding: 20px;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
        }
        .view {
            display: flex;
        }
        .text {
            padding: 4px 8px;
        }
        .image {
            background: #e0e0e0;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Tenon Debug Panel</h1>
        <div class="debug-panel">
`)
}

func (r *HTMLRenderer) writeFooter() {
	r.builder.WriteString(`
        </div>
    </div>
</body>
</html>
`)
}

func (r *HTMLRenderer) renderElement(element types.Element) {
	indentStr := strings.Repeat("    ", r.indent+2)

	switch props := element.GetProps().(type) {
	case *types.ViewProps:
		r.renderView(props, indentStr)
	case *types.TextProps:
		r.renderText(props, indentStr)
	case *types.ImageProps:
		r.renderImage(props, indentStr)
	}

	r.indent++
	for _, child := range element.GetChildren() {
		if child != nil {
			r.renderElement(child)
		}
	}
	r.indent--
}

func (r *HTMLRenderer) renderView(props *types.ViewProps, indentStr string) {
	style := props.Style

	r.builder.WriteString(indentStr)
	r.builder.WriteString(`<div class="view" style="`)

	if style != nil {
		r.renderFlexStyle(style)
		r.renderDimensionStyle(style)
		r.renderSpacingStyle(style)
		if style.FlexGrow > 0 {
			r.builder.WriteString(fmt.Sprintf("flex-grow: %.0f; ", style.FlexGrow))
		}
	}

	r.builder.WriteString(`border: 1px dashed #ccc; `)
	r.builder.WriteString(`">` + "\n")
}

func (r *HTMLRenderer) renderFlexStyle(style *types.ViewStyle) {
	flexDirection := "row"
	switch style.FlexDirection {
	case 0:
		flexDirection = "column"
	case 1:
		flexDirection = "column-reverse"
	case 2:
		flexDirection = "row"
	case 3:
		flexDirection = "row-reverse"
	}
	r.builder.WriteString(fmt.Sprintf("flex-direction: %s; ", flexDirection))

	justifyContent := "flex-start"
	switch style.JustifyContent {
	case 0:
		justifyContent = "flex-start"
	case 1:
		justifyContent = "center"
	case 2:
		justifyContent = "flex-end"
	case 3:
		justifyContent = "space-between"
	case 4:
		justifyContent = "space-around"
	case 5:
		justifyContent = "space-evenly"
	}
	r.builder.WriteString(fmt.Sprintf("justify-content: %s; ", justifyContent))

	alignItems := "stretch"
	switch style.AlignItems {
	case 0:
		alignItems = "stretch"
	case 1:
		alignItems = "flex-start"
	case 2:
		alignItems = "center"
	case 3:
		alignItems = "flex-end"
	case 4:
		alignItems = "stretch"
	case 5:
		alignItems = "baseline"
	}
	r.builder.WriteString(fmt.Sprintf("align-items: %s; ", alignItems))
}

func (r *HTMLRenderer) renderDimensionStyle(style *types.ViewStyle) {
	if style.Width.Unit == types.UnitPx {
		r.builder.WriteString(fmt.Sprintf("width: %.0fpx; ", style.Width.Value))
	}
	if style.Height.Unit == types.UnitPx {
		r.builder.WriteString(fmt.Sprintf("height: %.0fpx; ", style.Height.Value))
	}
	if style.Background != "" {
		r.builder.WriteString(fmt.Sprintf("background: %s; ", style.Background))
	}
}

func (r *HTMLRenderer) renderSpacingStyle(style *types.ViewStyle) {
	if style.Padding.Unit == types.UnitPx {
		r.builder.WriteString(fmt.Sprintf("padding: %.0fpx; ", style.Padding.Value))
	}
	if style.Margin.Unit == types.UnitPx {
		r.builder.WriteString(fmt.Sprintf("margin: %.0fpx; ", style.Margin.Value))
	}
	if style.MarginBottom.Unit == types.UnitPx {
		r.builder.WriteString(fmt.Sprintf("margin-bottom: %.0fpx; ", style.MarginBottom.Value))
	}
	if style.MarginRight.Unit == types.UnitPx {
		r.builder.WriteString(fmt.Sprintf("margin-right: %.0fpx; ", style.MarginRight.Value))
	}
}

func (r *HTMLRenderer) renderText(props *types.TextProps, indentStr string) {
	style := props.Style

	r.builder.WriteString(indentStr)
	r.builder.WriteString(`<span class="text" style="`)

	if style != nil {
		if style.FontSize.Unit == types.UnitPx {
			r.builder.WriteString(fmt.Sprintf("font-size: %.0fpx; ", style.FontSize.Value))
		}
		if style.Color != "" {
			r.builder.WriteString(fmt.Sprintf("color: %s; ", style.Color))
		}
	}

	r.builder.WriteString(`">`)
	r.builder.WriteString(escapeHTML(props.Content))
	r.builder.WriteString(`</span>` + "\n")
}

func (r *HTMLRenderer) renderImage(props *types.ImageProps, indentStr string) {
	style := props.Style

	r.builder.WriteString(indentStr)
	r.builder.WriteString(`<div class="image" style="`)

	if style != nil {
		if style.Width.Unit == types.UnitPx {
			r.builder.WriteString(fmt.Sprintf("width: %.0fpx; ", style.Width.Value))
		}
		if style.Height.Unit == types.UnitPx {
			r.builder.WriteString(fmt.Sprintf("height: %.0fpx; ", style.Height.Value))
		}
	}

	r.builder.WriteString(`border: 1px solid #999; `)
	r.builder.WriteString(`">`)
	if props.Source != "" {
		r.builder.WriteString(fmt.Sprintf(`<img src="%s" style="width:100%%;height:100%%;"/>`, props.Source))
	}
	r.builder.WriteString(`</div>` + "\n")
}

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	s = strings.ReplaceAll(s, `'`, "&#39;")
	return s
}

func RenderToHTML(element types.Element) string {
	renderer := NewHTMLRenderer()
	return renderer.Render(element)
}

func SaveHTML(element types.Element, filename string) error {
	html := RenderToHTML(element)
	return writeFile(filename, []byte(html))
}

func writeFile(filename string, data []byte) error {
	return nil
}