package yoga

import (
	"fmt"
	"strings"
)

const (
	gPrintChanges = false
	gPrintSkips   = false
	gDebuging     = false
)

var (
	gCurrentDebugCount uint32 = 0
)

func vlog(config *Config,
	node *Node,
	level LogLevel,
	format string,
	args ...interface{}) {
	if config != nil {
		config.log(node, level, format, args...)
	} else {
		DefaultLogger(config, node, level, format, args...)
	}
}

func indent(base *strings.Builder, level uint32) {
	for i := uint32(0); i < level; i++ {
		base.WriteString("  ")
	}
}

func appendFloatOptionalIfDefined(base *strings.Builder, str string, value FloatOptional) {
	if value.isDefined() {
		base.WriteString(fmt.Sprintf("%s: %g; ", str, value.unwrap()))
	}
}

func appendNumberIfNotAuto(base *strings.Builder, str string, value Value) {
	if value.unit != UnitAuto {
		appendNumberIfNotUndefined(base, str, value)
	}
}

func appendNumberIfNotZero(base *strings.Builder, str string, number Value) {
	if number.unit == UnitAuto {
		base.WriteString(str)
		base.WriteString(": auto; ")
	} else if !number.IsUndefined() && number.value != 0 {
		appendNumberIfNotUndefined(base, str, number)
	}
}

func areFourValuesEqual(four [EdgeCount]CompactValue) bool {
	return inexactEquals(four[0], four[1]) &&
		inexactEquals(four[0], four[2]) &&
		inexactEquals(four[0], four[3])
}

func appendEdges(base *strings.Builder, key string, edges [EdgeCount]CompactValue) {
	if areFourValuesEqual(edges) {
		edgeValue := (&nodeDefaults).computeEdgeValueForColumn(edges, EdgeLeft)
		appendNumberIfNotUndefined(base, key, edgeValue.Value())
	} else {
		for edge := EdgeLeft; edge < EdgeCount; edge++ {
			appendNumberIfNotZero(base, fmt.Sprintf("%s-%s", key, edge.String()), edges[edge].Value())
		}
	}
}

func appendNumberIfNotUndefined(base *strings.Builder, str string, number Value) {
	if number.unit != UnitUndefined {
		if number.unit == UnitAuto {
			base.WriteString(str)
			base.WriteString(": auto; ")
		} else {
			unit := If(number.unit == UnitPoint, "px", "%")
			base.WriteString(fmt.Sprintf("%s: %g%s; ", str, number.value, unit))
		}
	}
}

func nodeToString(str *strings.Builder, node *Node, options PrintOptions, level uint32) {
	if node == nil {
		return
	}
	indent(str, level)
	str.WriteString("<div ")
	if options&PrintOptionsLayout == PrintOptionsLayout {
		str.WriteString("layout=\"")
		str.WriteString(fmt.Sprintf("width: %g; ", node.getLayout().dimension(DimensionWidth)))
		str.WriteString(fmt.Sprintf("height: %g; ", node.getLayout().dimension(DimensionHeight)))
		str.WriteString(fmt.Sprintf("top: %g;", node.getLayout().position(EdgeTop)))
		str.WriteString(fmt.Sprintf("left: %g; ", node.getLayout().position(EdgeLeft)))
		str.WriteString("\" ")
	}

	if options&PrintOptionsStyle == PrintOptionsStyle {
		str.WriteString("style=\"")
		style := node.getStyle()
		oriStyle := NewNode().getStyle()
		if style.flexDirection() != oriStyle.flexDirection() {
			str.WriteString(fmt.Sprintf("flex-direction: %s; ", style.flexDirection().String()))
		}
		if style.justifyContent() != oriStyle.justifyContent() {
			str.WriteString(fmt.Sprintf("justify-content: %s; ", style.justifyContent().String()))
		}
		if style.alignItems() != oriStyle.alignItems() {
			str.WriteString(fmt.Sprintf("align-items: %s; ", style.alignItems().String()))
		}
		if style.alignContent() != oriStyle.alignContent() {
			str.WriteString(fmt.Sprintf("align-content: %s; ", style.alignContent().String()))
		}
		if style.alignSelf() != oriStyle.alignSelf() {
			str.WriteString(fmt.Sprintf("align-self: %s; ", style.alignSelf().String()))
		}
		appendFloatOptionalIfDefined(str, "flex-grow", style.flexGrow())
		appendFloatOptionalIfDefined(str, "flex-shrink", style.flexShrink())
		appendNumberIfNotAuto(str, "flex-basis", style.flexBasis().Value())
		appendFloatOptionalIfDefined(str, "flex", style.flex())

		if style.flexWrap() != oriStyle.flexWrap() {
			str.WriteString(fmt.Sprintf("flex-wrap: %s; ", style.flexWrap().String()))
		}

		if style.overflow() != oriStyle.overflow() {
			str.WriteString(fmt.Sprintf("overflow: %s; ", style.overflow().String()))
		}

		if style.display() != oriStyle.display() {
			str.WriteString(fmt.Sprintf("display: %s; ", style.display().String()))
		}

		appendEdges(str, "margin", style.margin_)
		appendEdges(str, "padding", style.padding_)
		appendEdges(str, "border", style.border_)

		if style.gap(GutterAll).IsDefined() {
			appendNumberIfNotUndefined(str, "gap", style.gap(GutterAll).Value())
		} else {
			appendNumberIfNotUndefined(str, "column-gap", style.gap(GutterColumn).Value())
			appendNumberIfNotUndefined(str, "row-gap", style.gap(GutterRow).Value())
		}

		appendNumberIfNotAuto(str, "width", style.dimension(DimensionWidth).Value())
		appendNumberIfNotAuto(str, "height", style.dimension(DimensionHeight).Value())
		appendNumberIfNotAuto(str, "max-width", style.maxDimension(DimensionWidth).Value())
		appendNumberIfNotAuto(str, "max-height", style.maxDimension(DimensionHeight).Value())
		appendNumberIfNotAuto(str, "min-width", style.minDimension(DimensionWidth).Value())
		appendNumberIfNotAuto(str, "min-height", style.minDimension(DimensionHeight).Value())

		if style.positionType() != oriStyle.positionType() {
			str.WriteString(fmt.Sprintf("position: %s; ", style.positionType().String()))
		}

		appendNumberIfNotUndefined(str, "left", style.position(EdgeLeft).Value())
		appendNumberIfNotUndefined(str, "right", style.position(EdgeRight).Value())
		appendNumberIfNotUndefined(str, "top", style.position(EdgeTop).Value())
		appendNumberIfNotUndefined(str, "bottom", style.position(EdgeBottom).Value())

		if node.HasMeasureFunc() {
			str.WriteString(fmt.Sprintf("has-custom-measure-func: true; "))
		}
	}
	str.WriteString("\">")

	childCount := node.GetChildCount()
	if options&PrintOptionsChildren == PrintOptionsChildren && childCount > 0 {
		str.WriteString("\n")
		for i := uint32(0); i < childCount; i++ {
			nodeToString(str, node.GetChild(i), options, level+1)
		}
		indent(str, level)
		str.WriteString("</div>")
	}
}

func vprint(node *Node, printOptions PrintOptions) {
	var str strings.Builder
	str.Reset()
	nodeToString(&str, node, printOptions, 0)
	vlog(node.GetConfig(), node, LogLevelDebug, str.String())
}
