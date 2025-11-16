package yoga

type LayoutType int

const (
	LayoutTypeLayout LayoutType = iota
	LayoutTypeMeasure
	LayoutTypeCachedLayout
	LayoutTypeCachedMeasure
)

type LayoutPassReason uint8

const (
	LayoutPassReasonInitial LayoutPassReason = iota
	LayoutPassReasonAbsLayout
	LayoutPassReasonStretch
	LayoutPassReasonMultilineStretch
	LayoutPassReasonFlexLayout
	LayoutPassReasonMeasureChild
	LayoutPassReasonAbsMeasureChild
	LayoutPassReasonFlexMeasure
	LayoutPassReasonCount
)

func (r LayoutPassReason) String() string {
	switch r {
	case LayoutPassReasonInitial:
		return "initial"
	case LayoutPassReasonAbsLayout:
		return "abs_layout"
	case LayoutPassReasonStretch:
		return "stretch"
	case LayoutPassReasonMultilineStretch:
		return "multiline_stretch"
	case LayoutPassReasonFlexLayout:
		return "flex_layout"
	case LayoutPassReasonMeasureChild:
		return "measure"
	case LayoutPassReasonAbsMeasureChild:
		return "abs_measure"
	case LayoutPassReasonFlexMeasure:
		return "flex_measure"
	default:
		return "unknown"
	}
}

type LayoutData struct {
	layouts                     int32
	measures                    int32
	maxMeasureCache             uint32
	cachedLayouts               int32
	cachedMeasures              int32
	measureCallbacks            int32
	measureCallbackReasonsCount [LayoutPassReasonCount]uint8
}

var defaultLayoutData = LayoutData{
	layouts:                     0,
	measures:                    0,
	maxMeasureCache:             0,
	cachedLayouts:               0,
	cachedMeasures:              0,
	measureCallbacks:            0,
	measureCallbackReasonsCount: [LayoutPassReasonCount]uint8{},
}
