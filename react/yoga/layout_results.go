package yoga

const MaxCachedMeasurements int32 = 8

type LayoutResults struct {
	position_           [4]float32
	margin_             [4]float32
	border_             [4]float32
	padding_            [4]float32
	direction_          Direction
	hadOverflow_        bool
	dimension_          [2]float32
	measuredDimensions_ [2]float32

	computedFlexBasisGeneration uint32
	computedFlexBasis           FloatOptional
	generationCount             uint32
	lastOwnerDirection          Direction
	nextCachedMeasurementsIndex uint32
	cachedMeasurements          [MaxCachedMeasurements]CachedMeasurement
	cachedLayout                CachedMeasurement
}

func NewLayoutResults() *LayoutResults {
	return &LayoutResults{
		position_:                   [4]float32{Undefined, Undefined, Undefined, Undefined},
		margin_:                     [4]float32{Undefined, Undefined, Undefined, Undefined},
		border_:                     [4]float32{Undefined, Undefined, Undefined, Undefined},
		padding_:                    [4]float32{Undefined, Undefined, Undefined, Undefined},
		direction_:                  DirectionInherit,
		hadOverflow_:                false,
		dimension_:                  [2]float32{Undefined, Undefined},
		measuredDimensions_:         [2]float32{Undefined, Undefined},
		computedFlexBasisGeneration: 0,
		computedFlexBasis:           undefinedFloatOptional,
		lastOwnerDirection:          DirectionInherit,
		cachedMeasurements:          [MaxCachedMeasurements]CachedMeasurement{},
		cachedLayout:                CachedMeasurement{},
	}
}

func (l *LayoutResults) direction() Direction {
	return l.direction_
}

func (l *LayoutResults) setDirection(direction Direction) {
	l.direction_ = direction
}

func (l *LayoutResults) hadOverflow() bool {
	return l.hadOverflow_
}

func (l *LayoutResults) setHadOverflow(hadOverflow bool) {
	l.hadOverflow_ = hadOverflow
}

func (l *LayoutResults) dimension(axis Dimension) float32 {
	return l.dimension_[axis]
}

func (l *LayoutResults) setDimension(axis Dimension, value float32) {
	l.dimension_[axis] = value
}

func (l *LayoutResults) measuredDimension(axis Dimension) float32 {
	return l.measuredDimensions_[axis]
}

func (l *LayoutResults) setMeasuredDimension(axis Dimension, value float32) {
	l.measuredDimensions_[axis] = value
}

func (l *LayoutResults) position(edge Edge) float32 {
	assertCardinalEdge(edge)
	return l.position_[edge]
}

func (l *LayoutResults) setPosition(edge Edge, value float32) {
	assertCardinalEdge(edge)
	l.position_[edge] = value
}

func (l *LayoutResults) margin(edge Edge) float32 {
	assertCardinalEdge(edge)
	return l.margin_[edge]
}

func (l *LayoutResults) setMargin(edge Edge, value float32) {
	assertCardinalEdge(edge)
	l.margin_[edge] = value
}

func (l *LayoutResults) border(edge Edge) float32 {
	assertCardinalEdge(edge)
	return l.border_[edge]
}

func (l *LayoutResults) setBorder(edge Edge, value float32) {
	assertCardinalEdge(edge)
	l.border_[edge] = value
}

func (l *LayoutResults) padding(edge Edge) float32 {
	assertCardinalEdge(edge)
	return l.padding_[edge]
}

func (l *LayoutResults) setPadding(edge Edge, value float32) {
	assertCardinalEdge(edge)
	l.padding_[edge] = value
}

func (l *LayoutResults) isEqual(layout *LayoutResults) bool {
	isEqual := inexactEquals(l.position_[:], layout.position_[:]) &&
		inexactEquals(l.margin_[:], layout.margin_[:]) &&
		inexactEquals(l.border_[:], layout.border_[:]) &&
		inexactEquals(l.padding_[:], layout.padding_[:]) &&
		l.direction() == layout.direction() &&
		l.hadOverflow() == layout.hadOverflow() &&
		l.lastOwnerDirection == layout.lastOwnerDirection &&
		l.nextCachedMeasurementsIndex == layout.nextCachedMeasurementsIndex &&
		l.cachedLayout == layout.cachedLayout &&
		l.computedFlexBasis == layout.computedFlexBasis
	for i := uint32(0); i < uint32(MaxCachedMeasurements) && isEqual; i++ {
		isEqual = isEqual && l.cachedMeasurements[i] == layout.cachedMeasurements[i]
	}
	if !IsUndefined(l.measuredDimensions_[0]) || !IsUndefined(layout.measuredDimensions_[0]) {
		isEqual = isEqual && (l.measuredDimensions_[0] == layout.measuredDimensions_[0])
	}
	if !IsUndefined(l.measuredDimensions_[1]) || !IsUndefined(layout.measuredDimensions_[1]) {
		isEqual = isEqual && (l.measuredDimensions_[1] == layout.measuredDimensions_[1])
	}
	return isEqual
}

func assertCardinalEdge(edge Edge) {
	if edge > EdgeBottom {
		panic("Edge must be top/left/bottom/right")
	}
}
