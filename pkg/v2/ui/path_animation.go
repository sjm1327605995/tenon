package ui

import (
	"math"
	"time"

	"github.com/sjm1327605995/tenon/pkg/v2/render"
)

// Path 定义二维参数化曲线，t ∈ [0,1]。
type Path interface {
	// Evaluate 返回路径上参数 t 对应的点。
	Evaluate(t float64) (x, y float64)
}

// LinePath 是直线路径。
type LinePath struct {
	X0, Y0 float64
	X1, Y1 float64
}

func (p LinePath) Evaluate(t float64) (x, y float64) {
	return lerpF64(p.X0, p.X1, t), lerpF64(p.Y0, p.Y1, t)
}

// QuadBezierPath 二次贝塞尔曲线（一个控制点）。
type QuadBezierPath struct {
	X0, Y0 float64 // 起点
	CX, CY float64 // 控制点
	X1, Y1 float64 // 终点
}

func (p QuadBezierPath) Evaluate(t float64) (x, y float64) {
	u := 1 - t
	x = u*u*p.X0 + 2*u*t*p.CX + t*t*p.X1
	y = u*u*p.Y0 + 2*u*t*p.CY + t*t*p.Y1
	return
}

// CubicBezierPath 三次贝塞尔曲线（两个控制点）。
type CubicBezierPath struct {
	X0, Y0     float64
	CX1, CY1   float64
	CX2, CY2   float64
	X1, Y1     float64
}

func (p CubicBezierPath) Evaluate(t float64) (x, y float64) {
	u := 1 - t
	u2 := u * u
	t2 := t * t
	x = u2*u*p.X0 + 3*u2*t*p.CX1 + 3*u*t2*p.CX2 + t2*t*p.X1
	y = u2*u*p.Y0 + 3*u2*t*p.CY1 + 3*u*t2*p.CY2 + t2*t*p.Y1
	return
}

// ArcPath 圆弧路径（圆心 + 起始角 + 扫过角度）。
type ArcPath struct {
	CX, CY     float64 // 圆心
	Radius     float64
	StartAngle float64 // 弧度
	SweepAngle float64 // 弧度（正值逆时针，负值顺时针）
}

func (p ArcPath) Evaluate(t float64) (x, y float64) {
	angle := p.StartAngle + p.SweepAngle*t
	return p.CX + p.Radius*math.Cos(angle), p.CY + p.Radius*math.Sin(angle)
}

// PathAnimation 沿路径动画，同时可附加旋转、缩放跟随。
type PathAnimation struct {
	controller *AnimationController
	path       Path

	// 跟随选项
	FollowRotation bool   // 是否沿路径切线旋转
	BaseRotation   float64 // 额外旋转偏移（弧度）

	// 插值目标
	BeginScaleX, EndScaleX float32
	BeginScaleY, EndScaleY float32
	BeginAlpha, EndAlpha   float32

	// 回调
	OnUpdate func(pos render.Offset, rotation float32, scaleX, scaleY, alpha float32)
}

// NewPathAnimation 创建路径动画，默认 300ms。
func NewPathAnimation(path Path, duration time.Duration) *PathAnimation {
	pa := &PathAnimation{
		path: path,
		controller: &AnimationController{
			Duration:   duration,
			LowerBound: 0,
			UpperBound: 1,
		},
		BeginScaleX: 1, EndScaleX: 1,
		BeginScaleY: 1, EndScaleY: 1,
		BeginAlpha: 1, EndAlpha: 1,
	}
	pa.controller.AddListener(pa.tick)
	return pa
}

func (pa *PathAnimation) tick() {
	if pa.OnUpdate == nil {
		return
	}
	t := pa.controller.Value // 0 → 1
	x, y := pa.path.Evaluate(t)
	pos := render.Offset{X: float32(x), Y: float32(y)}

	var rot float32
	if pa.FollowRotation {
		// 用数值微分估算切线方向
		dt := 0.001
		if t > 1-dt {
			dt = -0.001
		}
		x2, y2 := pa.path.Evaluate(t + dt)
		rot = float32(math.Atan2(y2-float64(y), x2-float64(x)) + pa.BaseRotation)
	}

	sx := pa.BeginScaleX + (pa.EndScaleX-pa.BeginScaleX)*float32(t)
	sy := pa.BeginScaleY + (pa.EndScaleY-pa.BeginScaleY)*float32(t)
	alpha := pa.BeginAlpha + (pa.EndAlpha-pa.BeginAlpha)*float32(t)

	pa.OnUpdate(pos, rot, sx, sy, alpha)
}

// Start 启动路径动画。
func (pa *PathAnimation) Start() {
	pa.controller.Value = 0
	pa.controller.Status = AnimationDismissed
	pa.controller.Forward()
	if defaultEngine != nil {
		defaultEngine.RegisterAnimation(pa.controller)
	}
}

// Stop 停止并注销。
func (pa *PathAnimation) Stop() {
	pa.controller.Stop()
	if defaultEngine != nil {
		defaultEngine.UnregisterAnimation(pa.controller)
	}
}

func lerpF64(a, b, t float64) float64 {
	return a + (b-a)*t
}
