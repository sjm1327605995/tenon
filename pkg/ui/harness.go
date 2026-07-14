package ui

import "unicode/utf8"

// Headless test harness — mount a component tree and drive real interactions
// (click / hover / press / drag / keyboard-style text input) without opening an
// Ebiten window, then assert on the resulting render tree.
//
// It exists so packages built on top of pkg/ui (e.g. pkg/shadcn) and application
// code can test *behavior*, not just that a constructor returns a non-nil Node:
//
//	h := ui.Mount(App(), 400, 300)
//	h.Root().ByText("+").Click()
//	if got := h.Root().ByText("1"); !got.Exists() {
//		t.Fatal("counter did not increment")
//	}
//
// The harness reuses the exact reconcile → layout → hit-test → event path the
// real engine uses; the only thing it replaces is Ebiten's global input polling,
// which it drives directly against render nodes.
//
// Query handles are snapshots of the tree at the moment they were taken. Host
// nodes are reused across re-renders, so a handle stays valid while structure is
// stable; after an action that adds, removes, or replaces nodes, re-query from
// Root (or Overlays) rather than reusing an old handle.

// Harness drives a mounted component tree for tests. Not safe for concurrent use.
type Harness struct {
	g *game
}

// Mount reconciles root into a virtual w×h window and settles it (layout +
// mount effects run). Coordinates are logical pixels (device scale is fixed at
// 1). root is typically a Use(...) or Memo(...) component node.
func Mount(root *Node, w, h int) *Harness {
	if w <= 0 {
		w = 800
	}
	if h <= 0 {
		h = 600
	}
	initFont()
	uiScale = 1
	g := &game{root: root, w: w, h: h}
	activeGame = g
	g.rootFiber = reconcile(nil, nil, root)
	hn := &Harness{g: g}
	hn.settle()
	return hn
}

// MountDefault mounts root in a default 800×600 window.
func MountDefault(root *Node) *Harness { return Mount(root, 800, 600) }

// Root returns a Query for the root render node.
func (h *Harness) Root() *Query {
	return &Query{rn: rootRenderNode(h.g.rootFiber), h: h}
}

// Paint runs a full paint pass through a recording backend and returns the
// ordered draw operations (main tree first, then Portal overlays). It renders
// nothing to a GPU — use it to assert *what* a component draws (fills, text,
// borders, clips, focus ring…) headlessly, i.e. golden-style paint tests.
func (h *Harness) Paint() []PaintOp {
	rp := &recordPainter{}
	if rn := rootRenderNode(h.g.rootFiber); rn != nil {
		paint(rp, rn)
	}
	for _, pf := range h.g.portals {
		if pf.overlayRoot != nil {
			paint(rp, pf.overlayRoot)
		}
	}
	return rp.ops
}

// Overlays returns a Query per live Portal overlay (dialogs, dropdowns,
// tooltips…), in paint order (last is topmost). Portals live outside the main
// tree, so Root().Find does not reach them — search these instead.
func (h *Harness) Overlays() []*Query {
	out := make([]*Query, 0, len(h.g.portals))
	for _, pf := range h.g.portals {
		if pf.overlayRoot != nil {
			out = append(out, &Query{rn: pf.overlayRoot, h: h})
		}
	}
	return out
}

// Size reports the current virtual window size in logical pixels.
func (h *Harness) Size() (w, hgt int) { return h.g.w, h.g.h }

// Resize changes the virtual window size and re-lays-out, mirroring a real
// window resize (only size-dependent subtrees recompute).
func (h *Harness) Resize(w, hgt int) {
	h.g.w, h.g.h = w, hgt
	h.settle()
}

// Step advances animations and continuous loops by dtMs milliseconds, then
// settles. Use it to drive UseTween / UseTransition / UseElapsed to a point.
func (h *Harness) Step(dtMs float32) {
	if dtMs > 0 {
		h.g.tickAnims(dtMs)
		h.g.tickLoops(dtMs)
		h.g.tickLayoutAnim(dtMs)
	}
	h.settle()
}

// ClickAt hit-tests at a screen point (logical px) and fires the nearest
// onClick up the chain, exactly like a real left-click. Returns whether a
// handler fired.
func (h *Harness) ClickAt(x, y float32) bool {
	n := h.g.hitTop(x, y)
	for c := n; c != nil; c = c.parent {
		if c.onClick != nil {
			c.onClick()
			h.settle()
			return true
		}
	}
	return false
}

// RightClickAt hit-tests at a screen point and fires the nearest onContextMenu
// up the chain (right-click), passing the point. Returns whether a handler fired.
func (h *Harness) RightClickAt(x, y float32) bool {
	for c := h.g.hitTop(x, y); c != nil; c = c.parent {
		if c.onContextMenu != nil {
			c.onContextMenu(x, y)
			h.settle()
			return true
		}
	}
	return false
}

// Focused returns a Query for the currently focused node, or an empty Query if
// nothing is focused.
func (h *Harness) Focused() *Query {
	if h.g.focusedFiber == nil {
		return &Query{h: h}
	}
	return &Query{rn: h.g.focusedFiber.rnode, h: h}
}

// Tab moves focus to the next focusable element (Tab key), wrapping around, and
// returns the newly focused node. Focusables are buttons and inputs, in tree
// order, including Portal overlays.
func (h *Harness) Tab() *Query {
	h.g.focusNext(true)
	h.settle()
	return h.Focused()
}

// ShiftTab moves focus to the previous focusable element (Shift+Tab).
func (h *Harness) ShiftTab() *Query {
	h.g.focusNext(false)
	h.settle()
	return h.Focused()
}

// Arrow moves focus within the nearest matching ArrowNav group (roving focus),
// like pressing an arrow key: orient picks vertical (Up/Down) or horizontal
// (Left/Right); forward = Down/Right. Returns the newly focused node (unchanged
// if no matching group contains focus). Settles afterward.
func (h *Harness) Arrow(orient NavOrient, forward bool) *Query {
	h.g.moveFocusInGroup(forward, orient)
	h.settle()
	return h.Focused()
}

// Enter activates the focused element (Enter/Space) — fires its onClick unless
// it is an input. Returns whether a handler fired. Settles afterward.
func (h *Harness) Enter() bool {
	ok := h.g.activateFocused()
	h.settle()
	return ok
}

// Escape fires the Esc action: the topmost UseEscape handler if any is active,
// otherwise it clears focus. Settles afterward.
func (h *Harness) Escape() {
	h.g.fireEscape()
	h.settle()
}

// settle drives the tree to a fixed point: drain the re-render queue, lay out,
// run pending effects, and repeat until nothing is dirty and no effects remain.
func (h *Harness) settle() {
	g := h.g
	for i := 0; i < 100; i++ {
		for guard := 0; len(g.dirty) > 0 && guard < 100; guard++ {
			g.flushDirty()
		}
		g.rootRN = rootRenderNode(g.rootFiber)
		g.layout()
		if len(g.pendingEffects) == 0 {
			if len(g.dirty) == 0 {
				return
			}
			continue
		}
		g.flushEffects()
	}
}

// Query is a handle to a render node in a mounted tree. A nil or empty Query is
// safe to call read methods on (they return zero values), so finder chains like
// q.ByText("x").Click() never panic on a miss — guard with Exists first.
type Query struct {
	rn *renderNode
	h  *Harness
}

// Exists reports whether this handle points at a real node.
func (q *Query) Exists() bool { return q != nil && q.rn != nil }

// Kind returns the node kind: "box", "text", "input", "image", or "scroll".
func (q *Query) Kind() string {
	if !q.Exists() {
		return ""
	}
	switch q.rn.kind {
	case rnText:
		return "text"
	case rnInput:
		return "input"
	case rnImage:
		return "image"
	case rnScroll:
		return "scroll"
	default:
		return "box"
	}
}

// Text returns this node's own text ("" if it is not a text node).
func (q *Query) Text() string {
	if q.Exists() && q.rn.kind == rnText {
		return q.rn.text
	}
	return ""
}

// Value returns an input node's current value ("" if not an input).
func (q *Query) Value() string {
	if q.Exists() && q.rn.kind == rnInput {
		return q.rn.value
	}
	return ""
}

// Bounds returns the node's absolute layout rectangle (logical px).
func (q *Query) Bounds() Rect {
	if !q.Exists() {
		return Rect{}
	}
	return q.rn.bounds
}

// Focusable reports whether the node participates in Tab focus navigation.
func (q *Query) Focusable() bool { return q.Exists() && q.rn.focusable }

// Clickable reports whether the node itself carries an onClick handler.
func (q *Query) Clickable() bool { return q.Exists() && q.rn.onClick != nil }

// Count returns the number of direct child render nodes.
func (q *Query) Count() int {
	if !q.Exists() {
		return 0
	}
	return len(q.rn.children)
}

// Child returns the i-th direct child (empty Query if out of range).
func (q *Query) Child(i int) *Query {
	if !q.Exists() || i < 0 || i >= len(q.rn.children) {
		return &Query{h: q.h}
	}
	return &Query{rn: q.rn.children[i], h: q.h}
}

// Children returns all direct children as Queries.
func (q *Query) Children() []*Query {
	if !q.Exists() {
		return nil
	}
	out := make([]*Query, len(q.rn.children))
	for i, c := range q.rn.children {
		out[i] = &Query{rn: c, h: q.h}
	}
	return out
}

// Texts returns the text of every text node in this subtree, in tree order.
func (q *Query) Texts() []string {
	var out []string
	q.walk(func(n *renderNode) {
		if n.kind == rnText {
			out = append(out, n.text)
		}
	})
	return out
}

// AllText joins the subtree's text nodes with a single space — handy for
// asserting a component rendered some copy without pinning its node structure.
func (q *Query) AllText() string {
	texts := q.Texts()
	out := ""
	for i, s := range texts {
		if i > 0 {
			out += " "
		}
		out += s
	}
	return out
}

func (q *Query) walk(fn func(*renderNode)) {
	if !q.Exists() {
		return
	}
	var rec func(*renderNode)
	rec = func(n *renderNode) {
		fn(n)
		for _, c := range n.children {
			rec(c)
		}
	}
	rec(q.rn)
}

// Find returns the first node in this subtree (self included, depth-first) that
// matches pred, or an empty Query.
func (q *Query) Find(pred func(*Query) bool) *Query {
	var found *renderNode
	var rec func(*renderNode) bool
	rec = func(n *renderNode) bool {
		if pred(&Query{rn: n, h: q.h}) {
			found = n
			return true
		}
		for _, c := range n.children {
			if rec(c) {
				return true
			}
		}
		return false
	}
	if q.Exists() {
		rec(q.rn)
	}
	return &Query{rn: found, h: q.h}
}

// FindAll returns every node in this subtree (self included) matching pred.
func (q *Query) FindAll(pred func(*Query) bool) []*Query {
	var out []*Query
	q.walk(func(n *renderNode) {
		sub := &Query{rn: n, h: q.h}
		if pred(sub) {
			out = append(out, sub)
		}
	})
	return out
}

// ByText finds the first text node whose text equals s.
func (q *Query) ByText(s string) *Query {
	return q.Find(func(n *Query) bool { return n.rn.kind == rnText && n.rn.text == s })
}

// AllByText finds every text node whose text equals s.
func (q *Query) AllByText(s string) []*Query {
	return q.FindAll(func(n *Query) bool { return n.rn.kind == rnText && n.rn.text == s })
}

// ByKind finds the first node of the given kind ("box"/"text"/"input"/"image"/"scroll").
func (q *Query) ByKind(kind string) *Query {
	return q.Find(func(n *Query) bool { return n.Kind() == kind })
}

// ByPlaceholder finds the first input node with the given placeholder text.
func (q *Query) ByPlaceholder(s string) *Query {
	return q.Find(func(n *Query) bool { return n.rn.kind == rnInput && n.rn.placeholder == s })
}

// Placeholder returns an input node's placeholder text ("" if not an input).
func (q *Query) Placeholder() string {
	if q.Exists() && q.rn.kind == rnInput {
		return q.rn.placeholder
	}
	return ""
}

// IsFocused reports whether this node is the currently focused element.
func (q *Query) IsFocused() bool {
	return q.Exists() && q.h.g.focusedFiber != nil && q.h.g.focusedFiber.rnode == q.rn
}

// Clickables returns every node in this subtree that carries an onClick handler.
func (q *Query) Clickables() []*Query {
	return q.FindAll(func(n *Query) bool { return n.rn.onClick != nil })
}

// Click hit-tests at this node's center and fires the nearest onClick up the
// chain — the same path a real left-click takes (so it respects transforms and
// overlapping overlays). Returns whether a handler fired. Settles afterward.
func (q *Query) Click() bool {
	if !q.Exists() {
		return false
	}
	b := q.rn.bounds
	return q.h.ClickAt(b.X+b.W/2, b.Y+b.H/2)
}

// Hover fires onHover(on) on this node and its ancestors that carry a hover
// handler, matching the engine's hover-chain behavior. Settles afterward.
func (q *Query) Hover(on bool) {
	if !q.Exists() {
		return
	}
	for c := q.rn; c != nil; c = c.parent {
		if c.onHover != nil {
			c.onHover(on)
		}
	}
	q.h.settle()
}

// Press fires the nearest onPress(down) up the chain (button press/release
// visual state). Returns whether a handler fired. Settles afterward.
func (q *Query) Press(down bool) bool {
	if !q.Exists() {
		return false
	}
	for c := q.rn; c != nil; c = c.parent {
		if c.onPress != nil {
			c.onPress(down)
			q.h.settle()
			return true
		}
	}
	return false
}

// Drag fires the nearest onDrag(dx, dy) up the chain with a logical-pixel delta
// (sliders, reorderable lists…). Returns whether a handler fired. Settles.
func (q *Query) Drag(dx, dy float32) bool {
	if !q.Exists() {
		return false
	}
	for c := q.rn; c != nil; c = c.parent {
		if c.onDrag != nil {
			c.onDrag(dx, dy)
			q.h.settle()
			return true
		}
	}
	return false
}

// Focus makes this node the focused element (as a click on an input would).
// Returns the same Query for chaining.
func (q *Query) Focus() *Query {
	if q.Exists() && q.rn.owner != nil {
		q.h.g.focusedFiber = q.rn.owner
		if q.rn.kind == rnInput {
			q.rn.caretPos = len(q.rn.value)
			q.rn.selAnchor = q.rn.caretPos
		}
	}
	return q
}

// Type inserts text at the input's caret and fires onChange, mirroring the
// engine's edit path (in-place value update + change notification). Returns
// false if this is not an input node. Settles afterward; re-query for the
// updated node when the input is controlled.
func (q *Query) Type(text string) bool {
	if !q.Exists() || q.rn.kind != rnInput {
		return false
	}
	rn := q.rn
	val := rn.value
	caret := clampi(rn.caretPos, 0, len(val))
	val = val[:caret] + text + val[caret:]
	caret += len(text)
	rn.value = val
	rn.caretPos, rn.selAnchor = caret, caret
	if rn.onChange != nil {
		rn.onChange(val)
	}
	q.h.settle()
	return true
}

// SetValue replaces an input's entire value and fires onChange. Returns false
// if this is not an input node. Settles afterward.
func (q *Query) SetValue(text string) bool {
	if !q.Exists() || q.rn.kind != rnInput {
		return false
	}
	rn := q.rn
	rn.value = text
	rn.caretPos, rn.selAnchor = len(text), len(text)
	if rn.onChange != nil {
		rn.onChange(text)
	}
	q.h.settle()
	return true
}

// Backspace deletes up to n runes before the caret and fires onChange, mirroring
// the engine's Backspace handling. Returns false if this is not an input node.
// Settles afterward; re-query for the updated node when the input is controlled.
func (q *Query) Backspace(n int) bool {
	if !q.Exists() || q.rn.kind != rnInput {
		return false
	}
	rn := q.rn
	val := rn.value
	caret := clampi(rn.caretPos, 0, len(val))
	for i := 0; i < n && caret > 0; i++ {
		_, sz := utf8.DecodeLastRuneInString(val[:caret])
		val = val[:caret-sz] + val[caret:]
		caret -= sz
	}
	rn.value = val
	rn.caretPos, rn.selAnchor = caret, caret
	if rn.onChange != nil {
		rn.onChange(val)
	}
	q.h.settle()
	return true
}

// Clear empties an input's value and fires onChange. Returns false if this is
// not an input node. Settles afterward.
func (q *Query) Clear() bool { return q.SetValue("") }

// ScrollBy scrolls the nearest scrollable ancestor (a ScrollView) by dy logical
// pixels — positive scrolls content upward, revealing lower rows — clamped to
// the content bounds, exactly like a wheel scroll. Returns whether a scroll
// container was found. Settles afterward.
func (q *Query) ScrollBy(dy float32) bool {
	if !q.Exists() {
		return false
	}
	for c := q.rn; c != nil; c = c.parent {
		if c.scroll {
			c.scrollY += dy
			q.h.g.boundsDirty = true
			q.h.settle()
			return true
		}
	}
	return false
}
