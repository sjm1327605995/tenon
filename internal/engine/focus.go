package engine

import (
	"github.com/sjm1327605995/tenon/internal/render"
)

// popupDismissers 存储所有需要点击外部时关闭的 popup 回调。
var popupDismissers = make(map[int]func())
var nextDismisserID int

// RegisterPopupDismisser 注册一个 popup 关闭回调，返回唯一 ID。
func RegisterPopupDismisser(fn func()) int {
	nextDismisserID++
	popupDismissers[nextDismisserID] = fn
	return nextDismisserID
}

// UnregisterPopupDismisser 注销指定 ID 的 popup 关闭回调。
func UnregisterPopupDismisser(id int) {
	delete(popupDismissers, id)
}

// DismissAllPopups 关闭所有已注册的 popup。
// 调用前先清空列表，防止回调中注册的新 popup 被误关。
func DismissAllPopups() {
	listeners := make(map[int]func())
	for k, v := range popupDismissers {
		listeners[k] = v
	}
	popupDismissers = make(map[int]func())
	for _, fn := range listeners {
		fn()
	}
}

// FocusNode 表示一个可聚焦的节点。
type FocusNode struct {
	// CanFocus 是否可聚焦。
	CanFocus bool
	// Focusable 是否可通过 Tab 键到达。
	Focusable bool
	// OnFocus 获得焦点时回调。
	OnFocus func()
	// OnBlur 失去焦点时回调。
	OnBlur func()

	// owner 持有该节点的 RenderObject
	owner render.RenderObject
}

// FocusManager 管理焦点树，处理 Tab 键导航和焦点切换。
type FocusManager struct {
	nodes    []*FocusNode
	focused  *FocusNode
	tabIndex int
}

// NewFocusManager 创建焦点管理器。
func NewFocusManager() *FocusManager {
	return &FocusManager{}
}

// Register 注册一个可聚焦的节点。
func (fm *FocusManager) Register(node *FocusNode) {
	fm.nodes = append(fm.nodes, node)
}

// Unregister 移除一个焦点节点。
func (fm *FocusManager) Unregister(node *FocusNode) {
	for i, n := range fm.nodes {
		if n == node {
			fm.nodes = append(fm.nodes[:i], fm.nodes[i+1:]...)
			if fm.focused == node {
				fm.focused = nil
			}
			break
		}
	}
}

// Focus 将焦点设置到指定节点。
func (fm *FocusManager) Focus(node *FocusNode) {
	if fm.focused == node {
		return
	}
	if fm.focused != nil && fm.focused.OnBlur != nil {
		fm.focused.OnBlur()
	}
	fm.focused = node
	if node != nil && node.OnFocus != nil {
		node.OnFocus()
	}
}

// Unfocus 移除当前焦点。
func (fm *FocusManager) Unfocus() {
	if fm.focused != nil && fm.focused.OnBlur != nil {
		fm.focused.OnBlur()
	}
	fm.focused = nil
}

// GetFocused 返回当前聚焦的节点。
func (fm *FocusManager) GetFocused() *FocusNode {
	return fm.focused
}

// NextFocus 将焦点移到下一个可聚焦节点（Tab 键）。
func (fm *FocusManager) NextFocus() {
	focusable := fm.getFocusableNodes()
	if len(focusable) == 0 {
		return
	}
	if fm.focused == nil {
		fm.Focus(focusable[0])
		return
	}
	for i, n := range focusable {
		if n == fm.focused {
			next := (i + 1) % len(focusable)
			fm.Focus(focusable[next])
			return
		}
	}
	fm.Focus(focusable[0])
}

// PreviousFocus 将焦点移到上一个可聚焦节点（Shift+Tab 键）。
func (fm *FocusManager) PreviousFocus() {
	focusable := fm.getFocusableNodes()
	if len(focusable) == 0 {
		return
	}
	if fm.focused == nil {
		fm.Focus(focusable[len(focusable)-1])
		return
	}
	for i, n := range focusable {
		if n == fm.focused {
			prev := (i - 1 + len(focusable)) % len(focusable)
			fm.Focus(focusable[prev])
			return
		}
	}
	fm.Focus(focusable[len(focusable)-1])
}

// Clear 清除所有焦点节点。
func (fm *FocusManager) Clear() {
	fm.Unfocus()
	fm.nodes = nil
}

func (fm *FocusManager) getFocusableNodes() []*FocusNode {
	var result []*FocusNode
	for _, n := range fm.nodes {
		if n.CanFocus && n.Focusable {
			result = append(result, n)
		}
	}
	return result
}

// SetOwner 设置节点关联的 RenderObject。
func (fn *FocusNode) SetOwner(ro render.RenderObject) {
	fn.owner = ro
}

// GetOwner 返回节点关联的 RenderObject。
func (fn *FocusNode) GetOwner() render.RenderObject {
	return fn.owner
}
