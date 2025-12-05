package yoga

type Style struct {
	direction_      Direction
	flexDirection_  FlexDirection
	justifyContent_ Justify
	alignContent_   Align
	alignItems_     Align
	alignSelf_      Align
	positionType_   PositionType
	flexWrap_       Wrap
	overflow_       Overflow
	display_        Display

	flex_          FloatOptional
	flexGrow_      FloatOptional
	flexShrink_    FloatOptional
	flexBasis_     CompactValue
	margin_        [EdgeCount]CompactValue
	position_      [EdgeCount]CompactValue
	padding_       [EdgeCount]CompactValue
	border_        [EdgeCount]CompactValue
	gap_           [GutterCount]CompactValue
	dimensions_    [DimensionCount]CompactValue
	minDimensions_ [DimensionCount]CompactValue
	maxDimensions_ [DimensionCount]CompactValue
	// Yoga specific properties, not compatible with flexbox specification
	aspectRatio_ FloatOptional
}

const (
	DefaultFlexGrow      float32 = 0.0
	DefaultFlexShrink    float32 = 0.0
	WebDefaultFlexShrink float32 = 1.0
)

var (
	edgeUndefinedCompactValue      = [EdgeCount]CompactValue{CompactValueOfUndefined(), CompactValueOfUndefined(), CompactValueOfUndefined(), CompactValueOfUndefined(), CompactValueOfUndefined(), CompactValueOfUndefined(), CompactValueOfUndefined(), CompactValueOfUndefined(), CompactValueOfUndefined()}
	gutterUndefinedCompactValue    = [GutterCount]CompactValue{CompactValueOfUndefined(), CompactValueOfUndefined(), CompactValueOfUndefined()}
	dimensionUndefinedCompactValue = [DimensionCount]CompactValue{CompactValueOfUndefined(), CompactValueOfUndefined()}
	dimensionAutoCompactValue      = [DimensionCount]CompactValue{CompactValueOfAuto(), CompactValueOfAuto()}

	defaultStyle = Style{
		direction_:      DirectionInherit,
		flexDirection_:  FlexDirectionColumn,
		justifyContent_: JustifyFlexStart,
		alignContent_:   AlignFlexStart,
		alignItems_:     AlignStretch,
		alignSelf_:      AlignAuto,
		positionType_:   PositionTypeRelative,
		flexWrap_:       WrapNoWrap,
		overflow_:       OverflowVisible,
		display_:        DisplayFlex,
		flex_:           undefinedFloatOptional,
		flexGrow_:       undefinedFloatOptional,
		flexShrink_:     undefinedFloatOptional,
		flexBasis_:      CompactValueOfAuto(),
		margin_:         edgeUndefinedCompactValue,
		position_:       edgeUndefinedCompactValue,
		padding_:        edgeUndefinedCompactValue,
		border_:         edgeUndefinedCompactValue,
		gap_:            gutterUndefinedCompactValue,
		dimensions_:     dimensionAutoCompactValue,
		minDimensions_:  dimensionUndefinedCompactValue,
		maxDimensions_:  dimensionUndefinedCompactValue,
		aspectRatio_:    undefinedFloatOptional,
	}
)

func NewStyle() *Style {
	return &Style{
		flex_:          undefinedFloatOptional,
		flexGrow_:      undefinedFloatOptional,
		flexShrink_:    undefinedFloatOptional,
		margin_:        edgeUndefinedCompactValue,
		position_:      edgeUndefinedCompactValue,
		padding_:       edgeUndefinedCompactValue,
		border_:        edgeUndefinedCompactValue,
		gap_:           gutterUndefinedCompactValue,
		dimensions_:    dimensionUndefinedCompactValue,
		minDimensions_: dimensionUndefinedCompactValue,
		maxDimensions_: dimensionUndefinedCompactValue,
		aspectRatio_:   undefinedFloatOptional,
	}
}

func (s *Style) direction() Direction {
	return s.direction_
}

func (s *Style) setDirection(direction Direction) {
	s.direction_ = direction
}

func (s *Style) flexDirection() FlexDirection {
	return s.flexDirection_
}

func (s *Style) setFlexDirection(flexDirection FlexDirection) {
	s.flexDirection_ = flexDirection
}

func (s *Style) justifyContent() Justify {
	return s.justifyContent_
}

func (s *Style) setJustifyContent(justifyContent Justify) {
	s.justifyContent_ = justifyContent
}

func (s *Style) alignContent() Align {
	return s.alignContent_
}

func (s *Style) setAlignContent(alignContent Align) {
	s.alignContent_ = alignContent
}

func (s *Style) alignItems() Align {
	return s.alignItems_
}

func (s *Style) setAlignItems(alignItems Align) {
	s.alignItems_ = alignItems
}

func (s *Style) alignSelf() Align {
	return s.alignSelf_
}

func (s *Style) setAlignSelf(alignSelf Align) {
	s.alignSelf_ = alignSelf
}

func (s *Style) positionType() PositionType {
	return s.positionType_
}

func (s *Style) setPositionType(positionType PositionType) {
	s.positionType_ = positionType
}

func (s *Style) flexWrap() Wrap {
	return s.flexWrap_
}

func (s *Style) setFlexWrap(flexWrap Wrap) {
	s.flexWrap_ = flexWrap
}

func (s *Style) overflow() Overflow {
	return s.overflow_
}

func (s *Style) setOverflow(overflow Overflow) {
	s.overflow_ = overflow
}

func (s *Style) display() Display {
	return s.display_
}

func (s *Style) setDisplay(display Display) {
	s.display_ = display
}

func (s *Style) flex() FloatOptional {
	return s.flex_
}

func (s *Style) flexGrow() FloatOptional {
	return s.flexGrow_
}

func (s *Style) flexShrink() FloatOptional {
	return s.flexShrink_
}

func (s *Style) flexBasis() CompactValue {
	return s.flexBasis_
}

func (s *Style) margin(edge Edge) CompactValue {
	return s.margin_[edge]
}

// setMargin
func (s *Style) setMargin(edge Edge, value CompactValue) {
	s.margin_[edge] = value
}

func (s *Style) position(edge Edge) CompactValue {
	return s.position_[edge]
}

// setPosition
func (s *Style) setPosition(edge Edge, value CompactValue) {
	s.position_[edge] = value
}

func (s *Style) padding(edge Edge) CompactValue {
	return s.padding_[edge]
}

// setPadding
func (s *Style) setPadding(edge Edge, value CompactValue) {
	s.padding_[edge] = value
}

func (s *Style) border(edge Edge) CompactValue {
	return s.border_[edge]
}

// setBorder
func (s *Style) setBorder(edge Edge, value CompactValue) {
	s.border_[edge] = value
}

func (s *Style) gap(gutter Gutter) CompactValue {
	return s.gap_[gutter]
}

// setGap
func (s *Style) setGap(gutter Gutter, value CompactValue) {
	s.gap_[gutter] = value
}

// dimension
func (s *Style) dimension(dimension Dimension) CompactValue {
	return s.dimensions_[dimension]
}

// setDimension
func (s *Style) setDimension(dimension Dimension, value CompactValue) {
	s.dimensions_[dimension] = value
}

// minDimension
func (s *Style) minDimension(dimension Dimension) CompactValue {
	return s.minDimensions_[dimension]
}

// setMinDimension
func (s *Style) setMinDimension(dimension Dimension, value CompactValue) {
	s.minDimensions_[dimension] = value
}

// maxDimension
func (s *Style) maxDimension(dimension Dimension) CompactValue {
	return s.maxDimensions_[dimension]
}

// setMaxDimension
func (s *Style) setMaxDimension(dimension Dimension, value CompactValue) {
	s.maxDimensions_[dimension] = value
}

// aspectRatio
func (s *Style) aspectRatio() FloatOptional {
	return s.aspectRatio_
}

// resolveColumnGap
func (s *Style) resolveColumnGap() CompactValue {
	if s.gap_[GutterColumn].IsDefined() {
		return s.gap_[GutterColumn]
	} else {
		return s.gap_[GutterAll]
	}
}

// resolveRowGap
func (s *Style) resolveRowGap() CompactValue {
	if s.gap_[GutterRow].IsDefined() {
		return s.gap_[GutterRow]
	} else {
		return s.gap_[GutterAll]
	}
}

// equal
func (s *Style) equal(other *Style) bool {
	return s.direction_ == other.direction_ &&
		s.flexDirection_ == other.flexDirection_ &&
		s.justifyContent_ == other.justifyContent_ &&
		s.alignContent_ == other.alignContent_ &&
		s.alignItems_ == other.alignItems_ && s.alignSelf_ == other.alignSelf_ &&
		s.positionType_ == other.positionType_ && s.flexWrap_ == other.flexWrap_ &&
		s.overflow_ == other.overflow_ && s.display_ == other.display_ &&
		inexactEqual(s.flex_.unwrap(), other.flex_.unwrap()) &&
		inexactEquals(s.flexGrow_, other.flexGrow_) &&
		inexactEquals(s.flexShrink_, other.flexShrink_) &&
		inexactEquals(s.flexBasis_, other.flexBasis_) &&
		inexactEquals(s.margin_[:], other.margin_[:]) &&
		inexactEquals(s.position_[:], other.position_[:]) &&
		inexactEquals(s.padding_[:], other.padding_[:]) &&
		inexactEquals(s.border_[:], other.border_[:]) &&
		inexactEquals(s.gap_[:], other.gap_[:]) &&
		inexactEquals(s.dimensions_[:], other.dimensions_[:]) &&
		inexactEquals(s.minDimensions_[:], other.minDimensions_[:]) &&
		inexactEquals(s.maxDimensions_[:], other.maxDimensions_[:]) &&
		inexactEquals(s.aspectRatio_, other.aspectRatio_)
}
