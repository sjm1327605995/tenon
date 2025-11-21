package style

import "github.com/sjm1327605995/tenon/react/yoga"

type Yoga interface {
	Yoga() *yoga.Node
}

type Option struct {
	apply func(opts Yoga)

	implSpecificOptFn any
}

func ApplyOptions[T any](base *T, opts ...Option) {

	for i := range opts {
		opt := opts[i]
		if opt.implSpecificOptFn != nil {
			optFn, ok := opt.implSpecificOptFn.(func(*T))
			if ok {
				optFn(base)
			}
		}
	}

}

func Direction(direction yoga.Direction) Option {
	return Option{
		apply: func(opts Yoga) {
			opts.Yoga().StyleSetDirection(direction)
		},
	}
}

func FlexDirection(flexDirection yoga.FlexDirection) Option {
	return Option{
		apply: func(opts Yoga) {
			opts.Yoga().StyleSetFlexDirection(flexDirection)
		},
	}
}

func JustifyContent(justifyContent yoga.Justify) Option {
	return Option{
		apply: func(opts Yoga) {
			opts.Yoga().StyleSetJustifyContent(justifyContent)
		},
	}
}

func AlignContent(alignContent yoga.Align) Option {
	return Option{
		apply: func(opts Yoga) {
			opts.Yoga().StyleSetAlignContent(alignContent)
		},
	}
}

func AlignItems(alignItems yoga.Align) Option {
	return Option{
		apply: func(opts Yoga) {
			opts.Yoga().StyleSetAlignItems(alignItems)
		},
	}
}

func AlignSelf(alignSelf yoga.Align) Option {
	return Option{
		apply: func(opts Yoga) {
			opts.Yoga().StyleSetAlignSelf(alignSelf)
		},
	}
}

func FlexWrap(flexWrap yoga.Wrap) Option {
	return Option{
		apply: func(opts Yoga) {
			opts.Yoga().StyleSetFlexWrap(flexWrap)
		},
	}
}

func Overflow(overflow yoga.Overflow) Option {
	return Option{
		apply: func(opts Yoga) {
			opts.Yoga().StyleSetOverflow(overflow)
		},
	}
}

func Display(display yoga.Display) Option {
	return Option{
		apply: func(opts Yoga) {
			opts.Yoga().StyleSetDisplay(display)
		},
	}
}

func Flex(flex float32) Option {
	return Option{
		apply: func(opts Yoga) {
			opts.Yoga().StyleSetFlex(flex)
		},
	}
}

func FlexGrow(flexGrow float32) Option {
	return Option{
		apply: func(opts Yoga) {
			opts.Yoga().StyleSetFlexGrow(flexGrow)
		},
	}
}

func FlexShrink(flexShrink float32) Option {
	return Option{
		apply: func(opts Yoga) {
			opts.Yoga().StyleSetFlexShrink(flexShrink)
		},
	}
}

func FlexBasis(flexBasis float32) Option {
	return Option{
		apply: func(opts Yoga) {
			opts.Yoga().StyleSetFlexBasis(flexBasis)
		},
	}
}

func FlexBasisPercent(flexBasis float32) Option {
	return Option{
		apply: func(opts Yoga) {
			opts.Yoga().StyleSetFlexBasisPercent(flexBasis)
		},
	}
}

func FlexBasisAuto() Option {
	return Option{
		apply: func(opts Yoga) {
			opts.Yoga().StyleSetFlexBasisAuto()
		},
	}
}

func Width(points float32) Option {
	return Option{
		apply: func(opts Yoga) {
			opts.Yoga().StyleSetWidth(points)
		},
	}
}

func WidthPercent(percent float32) Option {
	return Option{
		apply: func(opts Yoga) {
			opts.Yoga().StyleSetWidthPercent(percent)
		},
	}
}

func WidthAuto() Option {
	return Option{
		apply: func(opts Yoga) {
			opts.Yoga().StyleSetWidthAuto()
		},
	}
}

func Height(height float32) Option {
	return Option{
		apply: func(opts Yoga) {
			opts.Yoga().StyleSetHeight(height)
		},
	}
}

func HeightPercent(height float32) Option {
	return Option{
		apply: func(opts Yoga) {
			opts.Yoga().StyleSetHeightPercent(height)
		},
	}
}

func HeightAuto() Option {
	return Option{
		apply: func(opts Yoga) {
			opts.Yoga().StyleSetHeightAuto()
		},
	}
}

func PositionType(positionType yoga.PositionType) Option {
	return Option{
		apply: func(opts Yoga) {
			opts.Yoga().StyleSetPositionType(positionType)
		},
	}
}

func Position(edge yoga.Edge, position float32) Option {
	return Option{
		apply: func(opts Yoga) {
			opts.Yoga().StyleSetPosition(edge, position)
		},
	}
}

func PositionPercent(edge yoga.Edge, position float32) Option {
	return Option{
		apply: func(opts Yoga) {
			opts.Yoga().StyleSetPositionPercent(edge, position)
		},
	}
}

func Margin(edge yoga.Edge, margin float32) Option {
	return Option{
		apply: func(opts Yoga) {
			opts.Yoga().StyleSetMargin(edge, margin)
		},
	}
}

func MarginPercent(edge yoga.Edge, margin float32) Option {
	return Option{
		apply: func(opts Yoga) {
			opts.Yoga().StyleSetMarginPercent(edge, margin)
		},
	}
}

func MarginAuto(edge yoga.Edge) Option {
	return Option{
		apply: func(opts Yoga) {
			opts.Yoga().StyleSetMarginAuto(edge)
		},
	}
}

func Padding(edge yoga.Edge, padding float32) Option {
	return Option{
		apply: func(opts Yoga) {
			opts.Yoga().StyleSetPadding(edge, padding)
		},
	}
}

func PaddingPercent(edge yoga.Edge, padding float32) Option {
	return Option{
		apply: func(opts Yoga) {
			opts.Yoga().StyleSetPaddingPercent(edge, padding)
		},
	}
}

func BorderAll(border float32) Option {
	return Option{
		apply: func(opts Yoga) {
			opts.Yoga().StyleSetBorder(yoga.EdgeAll, border)
		},
	}
}

func Gap(gutter yoga.Gutter, gapLength float32) Option {
	return Option{
		apply: func(opts Yoga) {
			opts.Yoga().StyleSetGap(gutter, gapLength)
		},
	}
}
