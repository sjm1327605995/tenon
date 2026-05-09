package widgets

import (
	"time"

	"github.com/sjm1327605995/tenon/pkg/engine"
)

// FlipCardWidget 是一个带翻转动画的层叠卡片。
// 正面和背面两个子 Widget 绝对重叠，通过 scaleX 动画模拟 2D 翻转。
// 用法：
//   card := FlipCard(frontWidget, backWidget).Duration(300 * time.Millisecond)
//   // 在事件中通过 state.Toggle() 触发翻转
type FlipCardWidget struct {
	engine.BaseWidget
	Front    engine.Widget
	Back     engine.Widget
	Duration time.Duration
}

// FlipCard 创建翻转卡片，默认动画时长 300ms。
func FlipCard(front, back engine.Widget) FlipCardWidget {
	return FlipCardWidget{Front: front, Back: back, Duration: 300 * time.Millisecond}
}

func (f FlipCardWidget) WithDuration(d time.Duration) FlipCardWidget {
	f.Duration = d
	return f
}

func (f FlipCardWidget) CreateElement() engine.Element {
	return engine.NewStatefulElement(f)
}

func (f FlipCardWidget) CreateState() engine.State {
	return &FlipCardState{}
}

// FlipCardState 管理翻转动画和当前朝向。
type FlipCardState struct {
	engine.BaseStateOf[FlipCardWidget]
	flipped    bool
	controller *engine.AnimationController
}

func (s *FlipCardState) InitState() {
	w := s.Widget()
	s.controller = &engine.AnimationController{
		Duration:   w.Duration,
		LowerBound: 0,
		UpperBound: 1,
	}
	if engine.DefaultEngine() != nil {
		engine.DefaultEngine().RegisterAnimation(s.controller)
	}
	s.controller.AddListener(s.onTick)
}

func (s *FlipCardState) Dispose() {
	if engine.DefaultEngine() != nil {
		engine.DefaultEngine().UnregisterAnimation(s.controller)
	}
	s.controller.Stop()
	s.controller.RemoveListener(s.onTick)
}

func (s *FlipCardState) DidUpdateWidget(old engine.Widget) {
	oldW := old.(FlipCardWidget)
	newW := s.Widget()
	if oldW.Duration != newW.Duration {
		s.controller.Duration = newW.Duration
	}
}

// Toggle 触发翻转动画。可以在外部通过 key 调用：
//   ctx.FindAncestorWidgetOfExactType(...) // 或 GlobalKey 获取 state
func (s *FlipCardState) Toggle() {
	s.flipped = !s.flipped
	s.controller.Value = 0
	s.controller.Status = engine.AnimationDismissed
	s.controller.Forward()
}

// IsFlipped 返回当前是否显示背面。
func (s *FlipCardState) IsFlipped() bool {
	return s.flipped
}

func (s *FlipCardState) onTick() {
	s.SetState(nil)
}

func (s *FlipCardState) Build(ctx engine.BuildContext) engine.Widget {
	w := s.Widget()
	v := s.controller.Value // 0 → 1

	var frontScaleX, backScaleX float32
	if !s.flipped {
		// 从正面翻到背面（或停在正面）
		if v < 0.5 {
			frontScaleX = float32(1 - v*2)
			backScaleX = 0
		} else {
			frontScaleX = 0
			backScaleX = float32(-(v - 0.5) * 2)
		}
	} else {
		// 从背面翻到正面（或停在背面）
		if v < 0.5 {
			frontScaleX = 0
			backScaleX = float32(-(1 - v*2))
		} else {
			frontScaleX = float32((v - 0.5) * 2)
			backScaleX = 0
		}
	}

	return Stack(
		Positioned(Transform(w.Front).Scale(frontScaleX, 1)).L(0).T(0).R(0).B(0),
		Positioned(Transform(w.Back).Scale(backScaleX, 1)).L(0).T(0).R(0).B(0),
	)
}
