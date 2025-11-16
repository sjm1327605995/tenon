package yoga

const (
	AlignCount               = 9
	DimensionCount           = 2
	DirectionCount           = 3
	DisplayCount             = 2
	EdgeCount                = 9
	ExperimentalFeatureCount = 3
	FlexDirectionCount       = 4
	JustifyCount             = 6
	LogLevelCount            = 6
	MeasureModeCount         = 3
	NodeTypeCount            = 2
	OverflowCount            = 3
	PositionTypeCount        = 3
	PrintOptionsCount        = 3
	UnitCount                = 4
	WrapCount                = 3
	GutterCount              = 3
)

type Align uint8

const (
	AlignAuto Align = iota
	AlignFlexStart
	AlignCenter
	AlignFlexEnd
	AlignStretch
	AlignBaseline
	AlignSpaceBetween
	AlignSpaceAround
	AlignSpaceEvenly
)

func (a Align) String() string {
	switch a {
	case AlignAuto:
		return "auto"
	case AlignFlexStart:
		return "flex-start"
	case AlignCenter:
		return "center"
	case AlignFlexEnd:
		return "flex-end"
	case AlignStretch:
		return "stretch"
	case AlignBaseline:
		return "baseline"
	case AlignSpaceBetween:
		return "space-between"
	case AlignSpaceAround:
		return "space-around"
	case AlignSpaceEvenly:
		return "space-evenly"
	}
	return "unknown"
}

type Dimension uint8

const (
	DimensionWidth Dimension = iota
	DimensionHeight
)

func (d Dimension) String() string {
	switch d {
	case DimensionWidth:
		return "width"
	case DimensionHeight:
		return "height"
	}
	return "unknown"
}

type Direction uint8

const (
	DirectionInherit Direction = iota
	DirectionLTR
	DirectionRTL
)

func (d Direction) String() string {
	switch d {
	case DirectionInherit:
		return "inherit"
	case DirectionLTR:
		return "ltr"
	case DirectionRTL:
		return "rtl"
	}
	return "unknown"
}

type Display uint8

const (
	DisplayFlex Display = iota
	DisplayNone
)

func (d Display) String() string {
	switch d {
	case DisplayFlex:
		return "flex"
	case DisplayNone:
		return "none"
	}
	return "unknown"
}

type Edge uint8

const (
	EdgeLeft Edge = iota
	EdgeTop
	EdgeRight
	EdgeBottom
	EdgeStart
	EdgeEnd
	EdgeHorizontal
	EdgeVertical
	EdgeAll
)

func (e Edge) String() string {
	switch e {
	case EdgeLeft:
		return "left"
	case EdgeTop:
		return "top"
	case EdgeRight:
		return "right"
	case EdgeBottom:
		return "bottom"
	case EdgeStart:
		return "start"
	case EdgeEnd:
		return "end"
	case EdgeHorizontal:
		return "horizontal"
	case EdgeVertical:
		return "vertical"
	case EdgeAll:
		return "all"
	}
	return "unknown"
}

type Errata uint32

const (
	ErrataNone                                Errata = iota
	ErrataStretchFlexBasis                           = 1
	ErrataStartingEndingEdgeFromFlexDirection        = 2
	ErrataPositionStaticBehavesLikeRelative          = 4
	ErrataAll                                        = 2147483647
	ErrataClassic                                    = 2147483646
)

func (e Errata) String() string {
	switch e {
	case ErrataNone:
		return "none"
	case ErrataStretchFlexBasis:
		return "stretch-flex-basis"
	case ErrataStartingEndingEdgeFromFlexDirection:
		return "starting-ending-edge-from-flex-direction"
	case ErrataPositionStaticBehavesLikeRelative:
		return "position-static-behaves-like-relative"
	case ErrataAll:
		return "all"
	case ErrataClassic:
		return "classic"
	}
	return "unknown"
}

type ExperimentalFeature uint8

const (
	ExperimentalFeatureWebFlexBasis ExperimentalFeature = iota
	ExperimentalFeatureAbsolutePercentageAgainstPaddingEdge
	ExperimentalFeatureFixJNILocalRefOverflows
)

func (e ExperimentalFeature) String() string {
	switch e {
	case ExperimentalFeatureWebFlexBasis:
		return "web-flex-basis"
	case ExperimentalFeatureAbsolutePercentageAgainstPaddingEdge:
		return "absolute-percentage-against-padding-edge"
	case ExperimentalFeatureFixJNILocalRefOverflows:
		return "fix-jni-local-ref-overflows"
	}
	return "unknown"
}

type FlexDirection uint8

const (
	FlexDirectionColumn FlexDirection = iota
	FlexDirectionColumnReverse
	FlexDirectionRow
	FlexDirectionRowReverse
)

func (f FlexDirection) String() string {
	switch f {
	case FlexDirectionColumn:
		return "column"
	case FlexDirectionColumnReverse:
		return "column-reverse"
	case FlexDirectionRow:
		return "row"
	case FlexDirectionRowReverse:
		return "row-reverse"
	}
	return "unknown"
}

type Gutter uint8

const (
	GutterColumn Gutter = iota
	GutterRow
	GutterAll
)

func (g Gutter) String() string {
	switch g {
	case GutterColumn:
		return "column"
	case GutterRow:
		return "row"
	case GutterAll:
		return "all"
	}
	return "unknown"
}

type Justify uint8

const (
	JustifyFlexStart Justify = iota
	JustifyCenter
	JustifyFlexEnd
	JustifySpaceBetween
	JustifySpaceAround
	JustifySpaceEvenly
)

func (j Justify) String() string {
	switch j {
	case JustifyCenter:
		return "center"
	case JustifyFlexEnd:
		return "flex-end"
	case JustifySpaceBetween:
		return "space-between"
	case JustifySpaceAround:
		return "space-around"
	case JustifyFlexStart:
		return "flex-start"
	case JustifySpaceEvenly:
		return "space-evenly"
	}
	return "unknown"
}

type LogLevel uint8

const (
	LogLevelError LogLevel = iota
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
	LogLevelVerbose
	LogLevelFatal
)

func (l LogLevel) String() string {
	switch l {
	case LogLevelError:
		return "error"
	case LogLevelWarn:
		return "warn"
	case LogLevelInfo:
		return "info"
	case LogLevelDebug:
		return "debug"
	case LogLevelVerbose:
		return "verbose"
	case LogLevelFatal:
		return "fatal"
	}
	return "unknown"
}

type MeasureMode uint8

const (
	MeasureModeUndefined MeasureMode = iota
	MeasureModeExactly
	MeasureModeAtMost
)

func (m MeasureMode) String() string {
	switch m {
	case MeasureModeUndefined:
		return "undefined"
	case MeasureModeExactly:
		return "exactly"
	case MeasureModeAtMost:
		return "at-most"
	}
	return "unknown"
}

type NodeType uint8

const (
	NodeTypeDefault NodeType = iota
	NodeTypeText
)

func (n NodeType) String() string {
	switch n {
	case NodeTypeDefault:
		return "default"
	case NodeTypeText:
		return "text"
	}
	return "unknown"
}

type Overflow uint8

const (
	OverflowVisible Overflow = iota
	OverflowHidden
	OverflowScroll
)

func (o Overflow) String() string {
	switch o {
	case OverflowVisible:
		return "visible"
	case OverflowHidden:
		return "hidden"
	case OverflowScroll:
		return "scroll"
	}
	return "unknown"
}

type PositionType uint8

const (
	PositionTypeStatic PositionType = iota
	PositionTypeRelative
	PositionTypeAbsolute
)

func (p PositionType) String() string {
	switch p {
	case PositionTypeStatic:
		return "static"
	case PositionTypeRelative:
		return "relative"
	case PositionTypeAbsolute:
		return "absolute"
	}
	return "unknown"
}

type PrintOptions uint8

const (
	PrintOptionsLayout PrintOptions = 1 << iota
	PrintOptionsStyle
	PrintOptionsChildren
)

func (p PrintOptions) String() string {
	switch p {
	case PrintOptionsLayout:
		return "layout"
	case PrintOptionsStyle:
		return "style"
	case PrintOptionsChildren:
		return "children"
	}
	return "unknown"
}

type Unit uint8

const (
	UnitUndefined Unit = iota
	UnitPoint
	UnitPercent
	UnitAuto
)

func (u Unit) String() string {
	switch u {
	case UnitUndefined:
		return "undefined"
	case UnitPoint:
		return "point"
	case UnitPercent:
		return "percent"
	case UnitAuto:
		return "auto"
	}
	return "unknown"
}

type Wrap uint8

const (
	WrapNoWrap Wrap = iota
	WrapWrap
	WrapWrapReverse
)

func (w Wrap) String() string {
	switch w {
	case WrapNoWrap:
		return "no-wrap"
	case WrapWrap:
		return "wrap"
	case WrapWrapReverse:
		return "wrap-reverse"
	}
	return "unknown"
}
