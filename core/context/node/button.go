package node

type Button struct {
	State int
	INode
}

func NewButton(node INode) *Button {
	return &Button{
		INode: node,
	}
}
func (b *Button) OnClick() {

}
func (b *Button) OnHover() {

}
