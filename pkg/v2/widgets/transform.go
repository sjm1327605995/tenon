package widgets

import (
	"github.com/sjm1327605995/tenon/pkg/v2/render"
	"github.com/sjm1327605995/tenon/pkg/v2/ui"
)

// TransformWidget 对子 Widget 应用几何变换。
// 变换在绘制阶段生效，不影响 Yoga 布局计算（bounds 保持不变）。
// 适用于卡牌旋转、缩放、翻转等视觉效果。
type TransformWidget struct {
	ui.BaseWidget
	Child    ui.Widget
	Rotation float32 // 角度（度），正数为顺时针
	ScaleX   float32
	ScaleY   float32
	SkewX    float32
	SkewY    float32
	OriginX  float32 // 0~1，0.5 为几何中心
	OriginY  float32
	Alpha    float32 // 0~1 透明度
}

// Transform 创建变换包装器，默认锚点为中心、不缩放、不透明。
func Transform(child ui.Widget) TransformWidget {
	return TransformWidget{
		Child:   child,
		ScaleX:  1,
		ScaleY:  1,
		OriginX: 0.5,
		OriginY: 0.5,
		Alpha:   1,
	}
}

// 链式配置 ---------------------------------------------------------------

func (t TransformWidget) Rotate(deg float32) TransformWidget {
	t.Rotation = deg
	return t
}

func (t TransformWidget) Scale(x, y float32) TransformWidget {
	t.ScaleX = x
	t.ScaleY = y
	return t
}

func (t TransformWidget) ScaleUniform(s float32) TransformWidget {
	t.ScaleX = s
	t.ScaleY = s
	return t
}

func (t TransformWidget) Skew(x, y float32) TransformWidget {
	t.SkewX = x
	t.SkewY = y
	return t
}

func (t TransformWidget) Anchor(ox, oy float32) TransformWidget {
	t.OriginX = ox
	t.OriginY = oy
	return t
}

func (t TransformWidget) Opacity(a float32) TransformWidget {
	if a < 0 {
		a = 0
	}
	if a > 1 {
		a = 1
	}
	t.Alpha = a
	return t
}

// Widget 接口 -----------------------------------------------------------

func (t TransformWidget) CreateElement() ui.Element {
	e := &TransformElement{}
	e.SingleChildComponentElement.ComponentElement.BaseElement.Init(e, t)
	return e
}

// TransformElement 负责将变换属性同步到子树最近的 RenderObject。
type TransformElement struct {
	ui.SingleChildComponentElement
}

func (e *TransformElement) PerformRebuild(oldWidget ui.Widget) {
	w := e.GetWidget().(TransformWidget)
	e.Child = ui.UpdateChild(e, e.Child, w.Child)

	t := render.Transform{
		Rotation: w.Rotation,
		ScaleX:   w.ScaleX,
		ScaleY:   w.ScaleY,
		SkewX:    w.SkewX,
		SkewY:    w.SkewY,
		OriginX:  w.OriginX,
		OriginY:  w.OriginY,
		Alpha:    w.Alpha,
	}
	if ro := e.Child.FindRenderObject(); ro != nil {
		ro.SetTransform(t)
	}
}

func (e *TransformElement) Mount(parent ui.Element, slot int) {
	e.SingleChildComponentElement.Mount(parent, slot)
	w := e.GetWidget().(TransformWidget)
	if w.Child != nil {
		e.Child = ui.UpdateChild(e, nil, w.Child)
	}
	// 首次挂载也同步 transform
	if ro := e.Child.FindRenderObject(); ro != nil {
		ro.SetTransform(render.Transform{
			Rotation: w.Rotation,
			ScaleX:   w.ScaleX,
			ScaleY:   w.ScaleY,
			SkewX:    w.SkewX,
			SkewY:    w.SkewY,
			OriginX:  w.OriginX,
			OriginY:  w.OriginY,
			Alpha:    w.Alpha,
		})
	}
}
