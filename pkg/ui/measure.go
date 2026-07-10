package ui

// measureHook 保存某个元素最近一次布局后的逻辑坐标矩形，供锚定浮层定位。
type measureHook struct {
	rect  Rect
	fiber *Fiber
}

// Viewport 返回当前视口的逻辑尺寸（用于浮层贴边翻转等）。
func Viewport() Rect {
	if activeGame == nil {
		return Rect{}
	}
	return Rect{W: float32(activeGame.w) / uiScale, H: float32(activeGame.h) / uiScale}
}

// UseMeasure 返回一个属性和该元素最近测得的矩形（逻辑像素）。
// 把返回的属性挂到元素上，即可在下一帧读到它的屏幕位置——用于 Popover/Tooltip/Select 等锚定浮层。
//
//	ref, rect := ui.UseMeasure()
//	ui.Div(ref, trigger)           // 挂到锚点
//	ui.If(open, panelAt(rect))     // 用 rect 定位浮层
func UseMeasure() (*Node, Rect) {
	f := currentFiber
	_, raw := nextHook(f, func() any { return &measureHook{fiber: f} })
	h := raw.(*measureHook)
	h.fiber = f
	ref := &Node{typ: typeAttr, applyAttr: func(hp *hostProps) { hp.measure = h }}
	return ref, h.rect
}

// syncMeasures 布局后把带 measure 的节点的 bounds（物理）换算为逻辑矩形写回 hook；
// 变化时标记其组件重渲染，让锚定浮层跟随。
func syncMeasures(rn *renderNode) {
	if rn.measure != nil {
		lg := Rect{
			X: rn.bounds.X / uiScale, Y: rn.bounds.Y / uiScale,
			W: rn.bounds.W / uiScale, H: rn.bounds.H / uiScale,
		}
		if rn.measure.rect != lg {
			rn.measure.rect = lg
			if activeGame != nil {
				activeGame.markDirty(rn.measure.fiber)
			}
		}
	}
	for _, c := range rn.children {
		syncMeasures(c)
	}
}
