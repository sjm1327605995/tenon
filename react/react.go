package react

import (
	"github.com/sjm1327605995/tenon/react/api"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/unit"
	"github.com/sjm1327605995/tenon/react/core"
	"github.com/sjm1327605995/tenon/react/dom"
	"github.com/sjm1327605995/tenon/react/yoga"
	"image"
)

// ReactDOM manages the lifecycle of a Tenon application.
type ReactDOM struct {
	window *app.Window
	root   api.Component
	vdom   *core.VNode
	dom    dom.INode
}

// NewReactDOM creates a new ReactDOM instance.
func NewReactDOM() *ReactDOM {
	return &ReactDOM{}
}

// Render mounts the root component and starts the application event loop.
// This function will block until the application window is closed.
func (r *ReactDOM) Render(root api.Component) {
	r.root = root
	r.vdom = r.root.Render()
	r.dom = mount(r.vdom)

	go func() {
		r.window = new(app.Window)
		if err := r.loop(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

// loop is the main event loop for the application.
func (r *ReactDOM) loop() error {
	var ops op.Ops
	for {
		e := r.window.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			if dom.Metric != gtx.Metric {
				dom.Metric = gtx.Metric
			}

			// Set root size in Dp and calculate layout
			rootNode := r.dom.GetYogaNode()
			widthDp := gtx.Metric.PxToDp(gtx.Constraints.Max.X)
			heightDp := gtx.Metric.PxToDp(gtx.Constraints.Max.Y)
			rootNode.StyleSetWidth(float32(widthDp))
			rootNode.StyleSetHeight(float32(heightDp))
			rootNode.CalculateLayout(float32(widthDp), float32(heightDp), yoga.DirectionLTR)

			// Draw the real DOM tree.
			draw(gtx, r.dom)

			e.Frame(gtx.Ops)
		}
	}
}

// mount recursively transforms a VNode tree into a real DOM tree (INode).
func mount(vnode *core.VNode) dom.INode {
	if vnode == nil {
		return nil
	}

	var node dom.INode
	switch vnode.Type {
	case "View":
		node = dom.NewView()
	case "Text":
		node = dom.NewText()
	case "Image":
		node = dom.NewImage()
	default:
		panic("unknown VNode type: " + vnode.Type)
	}

	node.ApplyProps(vnode)

	var children []dom.INode
	for i, childVNode := range vnode.Children {
		childNode := mount(childVNode)
		if childNode != nil {
			children = append(children, childNode)
			node.GetYogaNode().InsertChild(childNode.GetYogaNode(), uint32(i))
		}
	}
	node.SetChildren(children)

	return node
}

// draw recursively paints the real DOM tree.
func draw(ctx layout.Context, node dom.INode) {
	if node == nil {
		return
	}

	yogaNode := node.GetYogaNode()
	// Convert Yoga's Dp layout to Px for Gio
	x := ctx.Metric.Dp(unit.Dp(yogaNode.LayoutLeft()))
	y := ctx.Metric.Dp(unit.Dp(yogaNode.LayoutTop()))
	w := ctx.Metric.Dp(unit.Dp(yogaNode.LayoutWidth()))
	h := ctx.Metric.Dp(unit.Dp(yogaNode.LayoutHeight()))

	offset := image.Pt(x, y)
	size := image.Pt(w, h)

	t := op.Offset(offset).Push(ctx.Ops)

	// Create a new context for the child with correct constraints and clipping
	childGtx := ctx
	childGtx.Constraints = layout.Exact(size)
	defer clip.Rect(image.Rectangle{Max: size}).Push(childGtx.Ops).Pop()

	node.Paint(childGtx)

	for _, child := range node.GetChildren() {
		draw(childGtx, child)
	}
	t.Pop()
}
