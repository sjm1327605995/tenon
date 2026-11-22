package ui

// measureHook 保存某个元素最近一次布局后的逻辑坐标矩形，供锚定浮层定位。
type measureHook struct {
	rect  Rect
	fiber *Fiber
}

// Viewport 返回当前视口的逻辑尺寸（用于浮层贴边翻转等）。
// 在组件渲染中读取时，该组件会在窗口尺寸变化后自动重渲染。
func Viewport() Rect {
	if activeGame == nil {
		return Rect{}
	}
	if currentFiber != nil {
		if activeGame.viewportReaders == nil {
			activeGame.viewportReaders = make(map[*Fiber]struct{})
		}
		activeGame.viewportReaders[currentFiber] = struct{}{}
	}
	return Rect{W: float32(activeGame.w) / uiScale, H: float32(activeGame.h) / uiScale}
}

// UseMeasure 返回一个属性和该元素最近测得的屏幕视觉矩形（逻辑像素）。
// 坐标包含自身及祖先的缩放、旋转、平移和伪 3D 变换，用于 Popover/Tooltip/Select 等锚定浮层。
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

// ScrollInfo 是某个 ScrollView 最近一次布局后的滚动状态（逻辑像素）。
type ScrollInfo struct {
	Offset   float32 // 已向下滚动的距离
	Viewport float32 // 可视区高度
}

type scrollHook struct {
	info  ScrollInfo
	fiber *Fiber
}

// UseScroll 返回一个属性和某个 ScrollView 最近的滚动状态。把属性挂到 ScrollView 上，
// 即可在组件里读到当前滚动偏移与视口高度——用于列表虚拟化（VirtualList）等。
func UseScroll() (*Node, ScrollInfo) {
	f := currentFiber
	_, raw := nextHook(f, func() any { return &scrollHook{fiber: f} })
	h := raw.(*scrollHook)
	h.fiber = f
	ref := &Node{typ: typeAttr, applyAttr: func(hp *hostProps) { hp.scrollRef = h }}
	return ref, h.info
}

// syncMeasures 布局后把带 measure 的节点的 bounds（物理）换算为逻辑矩形写回 hook，
// 并把带 UseScroll 的 ScrollView 的滚动状态写回；变化时标记其组件重渲染。
func measuredScreenRect(rn *renderNode) Rect {
	bounds := rn.bounds
	points := [4]pt{
		{X: bounds.X, Y: bounds.Y}, {X: bounds.X + bounds.W, Y: bounds.Y},
		{X: bounds.X, Y: bounds.Y + bounds.H}, {X: bounds.X + bounds.W, Y: bounds.Y + bounds.H},
	}
	for node := rn; node != nil; node = node.parent {
		if node.scene3D {
			continue
		}
		var cam *camera3D
		if node.parent != nil && node.parent.scene3D {
			cam = cameraOf(node.parent)
		}
		if node.hasTransform() || cam != nil {
			transform := contentAffine(layerOf(node, cam))
			for index, point := range points {
				points[index] = transform.transform(point)
			}
		}
	}
	minX, minY := points[0].X, points[0].Y
	maxX, maxY := minX, minY
	for _, point := range points[1:] {
		minX, minY = min(minX, point.X), min(minY, point.Y)
		maxX, maxY = max(maxX, point.X), max(maxY, point.Y)
	}
	return Rect{X: minX / uiScale, Y: minY / uiScale, W: (maxX - minX) / uiScale, H: (maxY - minY) / uiScale}
}

func syncMeasures(rn *renderNode) {
	if rn.measure != nil {
		lg := measuredScreenRect(rn)
		if rn.measure.rect != lg {
			rn.measure.rect = lg
			if activeGame != nil {
				activeGame.markDirty(rn.measure.fiber)
			}
		}
	}
	if rn.scrollRef != nil && rn.scroll {
		info := ScrollInfo{Offset: rn.scrollY / uiScale, Viewport: rn.bounds.H / uiScale}
		if rn.scrollRef.info != info {
			rn.scrollRef.info = info
			if activeGame != nil {
				activeGame.markDirty(rn.scrollRef.fiber)
			}
		}
	}
	for _, c := range rn.children {
		syncMeasures(c)
	}
}
