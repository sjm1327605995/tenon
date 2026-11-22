package g3n

import (
	"time"

	"github.com/g3n/engine/core"
	"github.com/sjm1327605995/tenon/pkg/value"
)

// TransformProps controls a native node transform. Nil Position/Rotation default
// to zero, nil Scale defaults to (1,1,1), and Spin is radians per second.
type TransformProps struct {
	Position value.Value[Vec3]
	Rotation value.Value[Vec3]
	Scale    value.Value[Vec3]
	Spin     value.Value[Vec3]
}

type nodeComponent struct {
	node   core.INode
	unsubs []func()
	spin   value.Value[Vec3]
}

func (c *nodeComponent) bindTransform(props TransformProps) {
	n := c.node.GetNode()
	if props.Position != nil {
		watch(&c.unsubs, props.Position, func(v Vec3) { n.SetPosition(v.X, v.Y, v.Z) })
	}
	if props.Rotation != nil {
		watch(&c.unsubs, props.Rotation, func(v Vec3) { n.SetRotation(v.X, v.Y, v.Z) })
	}
	if props.Scale != nil {
		watch(&c.unsubs, props.Scale, func(v Vec3) { n.SetScale(v.X, v.Y, v.Z) })
	}
	c.spin = props.Spin
}

func (c *nodeComponent) Update(dt time.Duration) {
	if c.node == nil || c.spin == nil {
		return
	}
	spin := c.spin.Get()
	seconds := float32(dt.Seconds())
	n := c.node.GetNode()
	if spin.X != 0 {
		n.RotateX(spin.X * seconds)
	}
	if spin.Y != 0 {
		n.RotateY(spin.Y * seconds)
	}
	if spin.Z != 0 {
		n.RotateZ(spin.Z * seconds)
	}
}

func (c *nodeComponent) Close() {
	for _, unsubscribe := range c.unsubs {
		unsubscribe()
	}
	c.unsubs = nil
	if c.node != nil {
		if parent := c.node.GetNode().Parent(); parent != nil {
			if removable, ok := parent.(interface{ Remove(core.INode) bool }); ok {
				removable.Remove(c.node)
			}
		}
	}
	c.node = nil
}
