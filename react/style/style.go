package style

import (
	"github.com/sjm1327605995/tenon/react/common"
	"github.com/sjm1327605995/tenon/react/style/unit"
)

type Option struct {
	Apply func(common.Element)
}

func Height(val unit.Unit) Option {
	return Option{Apply: func(node common.Element) {
		switch val.Type() {
		case "percent":
			node.Yoga().StyleSetHeightPercent(val.Value())
		case "pt":
			node.Yoga().StyleSetHeight(val.Value())
		case "auto":
			node.Yoga().StyleSetHeightAuto()
		}
	}}
}

func Width(val unit.Unit) Option {
	return Option{func(node common.Element) {
		switch val.Type() {
		case "percent":
			node.Yoga().StyleSetWidthPercent(val.Value())
		case "pt":
			node.Yoga().StyleSetWidth(val.Value())
		case "auto":
			node.Yoga().StyleSetWidthAuto()
		}
	}}
}
