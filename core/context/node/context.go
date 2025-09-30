package node

import (
	"github.com/dhconnelly/rtreego"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/millken/yoga"
)

type Context struct {
	Rtree *rtreego.Rtree
	root  *View
}

func NewContext() *Context {
	return &Context{
		Rtree: rtreego.NewTree(2, 25, 50),
		root:  NewView(),
	}
}

func (c *Context) nearestNeighbor(point rtreego.Point) {
	rect, _ := rtreego.NewRect(point, []float64{1, 1})
	p := c.Rtree.SearchIntersect(rect)
	if p == nil {
		return
	}
	for _, v := range p {
		if inode, ok := v.(INode); ok {
			inode.OnEvent()
		}
	}
}
func (c *Context) Update() {
	x, y := ebiten.CursorPosition()
	if x > 0 && y > 0 {
		c.ListenMouse(float32(x), float32(y))
	}
	c.root.Measure()
	c.root.OnLayout()

}
func (c *Context) SetLayout(outsideWidth, outsideHeight float32) {
	c.root.SetWidth(outsideWidth)
	c.root.SetHeight(outsideHeight)
	c.root.SetDirtiedFunc(func(node *yoga.Node) {
		c.root.SetRtreeRect(c.Rtree)
	})

}
func (c *Context) View(f func(v *View)) *Context {
	f(c.root)
	return c
}
func (c *Context) Render(screen *ebiten.Image) {
	c.root.OnDraw(screen, c.Rtree)
}
func (c *Context) ListenMouse(x, y float32) {
	c.nearestNeighbor(rtreego.Point{float64(x), float64(y)})
}
