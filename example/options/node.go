package tenon

import (
	"github.com/sjm1327605995/tenon/style"
)

type INode interface {
	style.Yoga
	Body(node ...INode) INode
}
