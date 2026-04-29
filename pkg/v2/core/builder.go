package core

// BuildEngine 负责 Widget → Element 的构建与 diff。
type BuildEngine struct {
	engine *Engine
}

func newBuildEngine(e *Engine) *BuildEngine {
	return &BuildEngine{engine: e}
}

func (b *BuildEngine) scheduleBuild(w Widget) {
	b.engine.buildQueue = append(b.engine.buildQueue, w)
}

func (b *BuildEngine) flushBuildQueue() {
	if len(b.engine.buildQueue) == 0 {
		return
	}
	queue := b.engine.buildQueue
	b.engine.buildQueue = b.engine.buildQueue[:0]

	for _, w := range queue {
		var newRoot Element
		if bw, ok := w.(interface{ RenderWithTracking() Element }); ok {
			newRoot = bw.RenderWithTracking()
		} else {
			newRoot = w.Render()
		}
		if newRoot == nil {
			continue
		}

		var oldRoot Element
		if w == b.engine.rootWidget {
			oldRoot = b.engine.rootElement
		} else if bw, ok := w.(interface{ GetRootElement() Element }); ok {
			oldRoot = bw.GetRootElement()
		}

		if oldRoot == nil {
			if w == b.engine.rootWidget {
				b.engine.rootElement = newRoot
			}
			if bw, ok := w.(interface{ SetRootElement(Element) }); ok {
				bw.SetRootElement(newRoot)
			}
			b.engine.onElementMounted(newRoot)
		} else {
			finalRoot := b.patchElement(oldRoot, newRoot)
			if w == b.engine.rootWidget {
				b.engine.rootElement = finalRoot
			}
			if bw, ok := w.(interface{ SetRootElement(Element) }); ok {
				bw.SetRootElement(finalRoot)
			}
		}
	}

	// 重建后布局可能变化
	if b.engine.hasLayoutDirty() {
		b.engine.calculateLayout()
	}
}

// patchElement 对新旧 Element 做同级浅对比。
// 同类型复用旧节点并递归同步子节点，不同类型替换根节点。
func (b *BuildEngine) patchElement(oldEl, newEl Element) Element {
	if oldEl.ElementType() == newEl.ElementType() {
		// 同类型：复用旧节点，同步 Yoga 样式
		oldLayout, oldOk := oldEl.(LayoutElement)
		newLayout, newOk := newEl.(LayoutElement)
		if oldOk && newOk {
			if oldYoga, newYoga := oldLayout.GetYoga(), newLayout.GetYoga(); oldYoga != nil && newYoga != nil {
				oldYoga.CopyStyleFrom(newYoga)
			}
		}
		// 同步组件属性（声明式重建的关键）
		if syncable, ok := oldEl.(PropertySyncable); ok {
			syncable.SyncFrom(newEl)
		}
		b.patchChildren(oldEl, newEl)
		oldEl.Mark(FlagNeedLayout | FlagNeedDraw)
		return oldEl
	}
	// 类型不同：替换根节点
	if oldEl.GetParent() != nil {
		parent := oldEl.GetParent()
		parent.RemoveChild(oldEl)
		parent.AppendChild(newEl)
	} else {
		b.engine.rootElement = newEl
		b.engine.onElementMounted(newEl)
	}
	return newEl
}

// patchChildren 对两个同类型 Element 的子节点做轻量同级对比。
// 前 minLen 个递归 patch，多余旧节点移除，多余新节点挂载。
func (b *BuildEngine) patchChildren(oldParent, newParent Element) {
	oldChildren := oldParent.GetChildren()
	newChildren := newParent.GetChildren()

	minLen := len(oldChildren)
	if len(newChildren) < minLen {
		minLen = len(newChildren)
	}

	// 1. 同步前 minLen 个子节点（递归复用）
	for i := 0; i < minLen; i++ {
		b.patchElement(oldChildren[i], newChildren[i])
	}

	// 2. 移除多余的旧节点
	for i := len(newChildren); i < len(oldChildren); i++ {
		oldParent.RemoveChild(oldChildren[i])
	}

	// 3. 挂载多余的新节点（先从 newParent 解除 Yoga 关系，再挂到 oldParent）
	if len(newChildren) > len(oldChildren) {
		toAdd := newChildren[len(oldChildren):]
		// 倒序从 newParent 移除，避免索引漂移
		for i := len(toAdd) - 1; i >= 0; i-- {
			newParent.RemoveChild(toAdd[i])
		}
		for _, child := range toAdd {
			oldParent.AppendChild(child)
		}
	}
}
