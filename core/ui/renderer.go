package ui

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sjm1327605995/tenon/yoga"
)

type Bool struct {
	v atomic.Value
}

func (b *Bool) Load() bool {
	val := b.v.Load()
	if val == nil {
		return false
	}
	return val.(bool)
}

func (b *Bool) Store(val bool) {
	b.v.Store(val)
}

type RenderEvent int

const (
	RenderEventMount RenderEvent = iota
	RenderEventUnmount
	RenderEventLayoutChange
	RenderEventStyleChange
	RenderEventChildrenChange
	RenderEventPropsChange
)

type RenderEventData struct {
	Event       RenderEvent
	Node        *RenderNode
	OldLayout   LayoutResults
	NewLayout   LayoutResults
	ChangeType  string
	Timestamp   time.Time
}

type RenderHandler func(RenderEventData)

type RenderNode struct {
	element   Element
	yogaNode *yoga.Node
	children  []*RenderNode
	parent    *RenderNode
	depth     int
	dirty     Bool
	layout    LayoutResults
	styleHash uint64
	index     int
}

func NewRenderNode(element Element) *RenderNode {
	return &RenderNode{
		element:  element,
		yogaNode: yoga.NewNode(),
		children: make([]*RenderNode, 0),
		depth:    0,
	}
}

func (r *RenderNode) IsDirty() bool {
	return r.dirty.Load()
}

func (r *RenderNode) MarkDirty() {
	r.dirty.Store(true)
}

func (r *RenderNode) ClearDirty() {
	r.dirty.Store(false)
}

func (r *RenderNode) GetLayout() LayoutResults {
	return getNodeLayout(r.yogaNode)
}

func (r *RenderNode) String() string {
	layout := r.layout
	if r.yogaNode != nil {
		layout = getNodeLayout(r.yogaNode)
	}
	return fmt.Sprintf("RenderNode{depth=%d, layout=%.2fx%.2f@%.2f,%.2f}",
		r.depth, layout.Width, layout.Height, layout.Left, layout.Top)
}

type RendererConfig struct {
	Width         float64
	Height        float64
	Direction     Direction
	FontScale     float32
	UseWebDefaults bool
}

func DefaultRendererConfig() *RendererConfig {
	return &RendererConfig{
		Width:         375,
		Height:        812,
		Direction:     DirectionLTR,
		FontScale:     1.0,
		UseWebDefaults: false,
	}
}

type ElementMapper struct {
	config   *RendererConfig
	rootNode *RenderNode
	mu       sync.RWMutex
	nodeMap  map[Element]*RenderNode
}

func NewElementMapper(config *RendererConfig) *ElementMapper {
	return &ElementMapper{
		config:  config,
		nodeMap: make(map[Element]*RenderNode),
	}
}

func (m *ElementMapper) MapElement(element Element) *RenderNode {
	m.mu.Lock()
	defer m.mu.Unlock()

	if node, exists := m.nodeMap[element]; exists {
		return node
	}

	renderNode := m.createRenderNode(element, nil, 0)
	m.nodeMap[element] = renderNode
	return renderNode
}

func (m *ElementMapper) createRenderNode(element Element, parent *RenderNode, depth int) *RenderNode {
	node := NewRenderNode(element)
	node.parent = parent
	node.depth = depth

	if parent != nil {
		parent.children = append(parent.children, node)
	}

	if children := element.Children(); len(children) > 0 {
		for i, child := range children {
			if ve, ok := child.(*ViewElement); ok {
				childNode := m.createRenderNode(ve, node, depth+1)
				childNode.index = i
				m.nodeMap[child] = childNode
			} else if te, ok := child.(*TextElement); ok {
				childNode := m.createRenderNode(te, node, depth+1)
				childNode.index = i
				m.nodeMap[child] = childNode
			}
		}
	}

	return node
}

func (m *ElementMapper) BuildTree(element Element) *RenderNode {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.nodeMap = make(map[Element]*RenderNode)
	rootNode := m.createRenderNode(element, nil, 0)
	m.nodeMap[element] = rootNode
	m.rootNode = rootNode
	return rootNode
}

func (m *ElementMapper) GetRootNode() *RenderNode {
	return m.rootNode
}

func (m *ElementMapper) GetNode(element Element) *RenderNode {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.nodeMap[element]
}

func (m *ElementMapper) UpdateElement(element Element) *RenderNode {
	m.mu.Lock()
	defer m.mu.Unlock()

	node := m.nodeMap[element]
	if node == nil {
		return nil
	}

	node.MarkDirty()
	return node
}

func (m *ElementMapper) RemoveElement(element Element) {
	m.mu.Lock()
	defer m.mu.Unlock()

	node := m.nodeMap[element]
	if node == nil {
		return
	}

	if node.parent != nil {
		for i, child := range node.parent.children {
			if child == node {
				node.parent.children = append(node.parent.children[:i], node.parent.children[i+1:]...)
				break
			}
		}
	}

	delete(m.nodeMap, element)

	for _, child := range node.children {
		m.RemoveElement(child.element)
	}
}

func (m *ElementMapper) CalculateLayout() {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.rootNode == nil || m.rootNode.yogaNode == nil {
		return
	}

	m.rootNode.yogaNode.CalculateLayout(
		float32(m.config.Width),
		float32(m.config.Height),
		yoga.Direction(m.config.Direction),
	)
}

func (m *ElementMapper) PrintTree(node *RenderNode, prefix string) {
	if node == nil {
		return
	}

	layout := getNodeLayout(node.yogaNode)
	elementType := fmt.Sprintf("%T", node.element)

	fmt.Printf("%s%s [%.2f x %.2f at (%.2f, %.2f)]\n",
		prefix, elementType, layout.Width, layout.Height, layout.Left, layout.Top)

	for _, child := range node.children {
		m.PrintTree(child, prefix+"  ")
	}
}

type Renderer struct {
	mapper      *ElementMapper
	config      *RendererConfig
	isRunning   Bool
	stopChan    chan struct{}
	eventChan   chan RenderEventData
	handlers    map[RenderEvent][]RenderHandler
	ticker      *time.Ticker
	interval    time.Duration
	rootElement Element
	oldTree    string
	newTree    string
	treeMu     sync.RWMutex
	runMu      sync.Mutex
}

func NewRenderer(config *RendererConfig) *Renderer {
	if config == nil {
		config = DefaultRendererConfig()
	}

	return &Renderer{
		mapper:    NewElementMapper(config),
		config:    config,
		stopChan:  make(chan struct{}),
		eventChan: make(chan RenderEventData, 1000),
		handlers:  make(map[RenderEvent][]RenderHandler),
		interval:  16 * time.Millisecond,
	}
}

func (r *Renderer) SetConfig(config *RendererConfig) {
	r.config = config
}

func (r *Renderer) GetMapper() *ElementMapper {
	return r.mapper
}

func (r *Renderer) RegisterHandler(event RenderEvent, handler RenderHandler) {
	r.handlers[event] = append(r.handlers[event], handler)
}

func (r *Renderer) UnregisterHandler(event RenderEvent, handler RenderHandler) {
	handlers := r.handlers[event]
	for i, h := range handlers {
		if fmt.Sprintf("%p", h) == fmt.Sprintf("%p", handler) {
			r.handlers[event] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
}

func (r *Renderer) emitEvent(eventData RenderEventData) {
	for _, handler := range r.handlers[eventData.Event] {
		handler(eventData)
	}

	select {
	case r.eventChan <- eventData:
	default:
	}
}

func (r *Renderer) SetRoot(element Element) {
	r.runMu.Lock()
	defer r.runMu.Unlock()

	r.rootElement = element

	oldRoot := r.mapper.GetRootNode()
	r.mapper.BuildTree(element)

	newRoot := r.mapper.GetRootNode()

	if oldRoot == nil {
		r.emitEvent(RenderEventData{
			Event:     RenderEventMount,
			Node:      newRoot,
			Timestamp: time.Now(),
		})
	} else {
		r.compareAndEmit(oldRoot, newRoot)
	}
}

func (r *Renderer) compareAndEmit(old, new *RenderNode) {
	oldLayout := old.GetLayout()
	newLayout := new.GetLayout()

	if oldLayout.Width != newLayout.Width ||
		oldLayout.Height != newLayout.Height ||
		oldLayout.Left != newLayout.Left ||
		oldLayout.Top != newLayout.Top {

		r.emitEvent(RenderEventData{
			Event:     RenderEventLayoutChange,
			Node:      new,
			OldLayout: oldLayout,
			NewLayout: newLayout,
			Timestamp: time.Now(),
		})
	}

	for i := range old.children {
		if i < len(new.children) {
			r.compareAndEmit(old.children[i], new.children[i])
		}
	}
}

func (r *Renderer) Start() {
	r.runMu.Lock()
	defer r.runMu.Unlock()

	if r.isRunning.Load() {
		return
	}

	r.isRunning.Store(true)
	r.ticker = time.NewTicker(r.interval)

	go r.renderLoop()
}

func (r *Renderer) Stop() {
	if !r.isRunning.Load() {
		return
	}

	r.isRunning.Store(false)
	close(r.stopChan)
	r.ticker.Stop()
}

func (r *Renderer) renderLoop() {
	for {
		select {
		case <-r.stopChan:
			return
		case <-r.ticker.C:
			r.calculateAndNotify()
		}
	}
}

func (r *Renderer) calculateAndNotify() {
	if r.rootElement == nil {
		return
	}

	oldRoot := r.mapper.GetRootNode()
	r.mapper.CalculateLayout()
	newRoot := r.mapper.GetRootNode()

	if oldRoot != nil && newRoot != nil {
		r.detectAndEmitChanges(oldRoot, newRoot)
	}
}

func (r *Renderer) detectAndEmitChanges(old, new *RenderNode) {
	oldLayout := old.GetLayout()
	newLayout := new.GetLayout()

	if !r.layoutsEqual(oldLayout, newLayout) {
		r.emitEvent(RenderEventData{
			Event:     RenderEventLayoutChange,
			Node:      new,
			OldLayout: oldLayout,
			NewLayout: newLayout,
			Timestamp: time.Now(),
		})
	}

	if len(old.children) != len(new.children) {
		r.emitEvent(RenderEventData{
			Event:      RenderEventChildrenChange,
			Node:       new,
			ChangeType: fmt.Sprintf("children count: %d -> %d", len(old.children), len(new.children)),
			Timestamp:  time.Now(),
		})
	}

	maxChildren := len(old.children)
	if len(new.children) > maxChildren {
		maxChildren = len(new.children)
	}

	for i := 0; i < maxChildren; i++ {
		if i >= len(old.children) {
			r.emitEvent(RenderEventData{
				Event:      RenderEventMount,
				Node:       new.children[i],
				ChangeType: fmt.Sprintf("added child at index %d", i),
				Timestamp:  time.Now(),
			})
		} else if i >= len(new.children) {
			r.emitEvent(RenderEventData{
				Event:      RenderEventUnmount,
				Node:       old.children[i],
				ChangeType: fmt.Sprintf("removed child at index %d", i),
				Timestamp:  time.Now(),
			})
		} else {
			r.detectAndEmitChanges(old.children[i], new.children[i])
		}
	}
}

func (r *Renderer) layoutsEqual(a, b LayoutResults) bool {
	return a.Width == b.Width && a.Height == b.Height &&
		a.Left == b.Left && a.Top == b.Top
}

func (r *Renderer) ForceRender() {
	r.calculateAndNotify()
}

func (r *Renderer) GetLayout() LayoutResults {
	if r.mapper.GetRootNode() == nil {
		return LayoutResults{}
	}
	return r.mapper.GetRootNode().GetLayout()
}

func (r *Renderer) PrintTree() {
	fmt.Println("\n========== Element Tree ==========")
	root := r.mapper.GetRootNode()
	if root != nil {
		r.mapper.PrintTree(root, "")
	} else {
		fmt.Println("(empty tree)")
	}
	fmt.Println("===================================")
}

type PrinterRenderer struct {
	renderer   *Renderer
	outputChan chan string
	isRunning  Bool
	stopChan   chan struct{}
}

func NewPrinterRenderer(config *RendererConfig) *PrinterRenderer {
	pr := &PrinterRenderer{
		renderer:   NewRenderer(config),
		outputChan: make(chan string, 100),
		stopChan:   make(chan struct{}),
	}

	pr.renderer.RegisterHandler(RenderEventMount, func(e RenderEventData) {
		pr.outputChan <- fmt.Sprintf("[MOUNT] %s at depth %d", pr.elementType(e.Node), e.Node.depth)
	})

	pr.renderer.RegisterHandler(RenderEventUnmount, func(e RenderEventData) {
		pr.outputChan <- fmt.Sprintf("[UNMOUNT] %s at depth %d", pr.elementType(e.Node), e.Node.depth)
	})

	pr.renderer.RegisterHandler(RenderEventLayoutChange, func(e RenderEventData) {
		pr.outputChan <- fmt.Sprintf("[LAYOUT] %s: %.2fx%.2f@%.2f,%.2f -> %.2fx%.2f@%.2f,%.2f",
			pr.elementType(e.Node),
			e.OldLayout.Width, e.OldLayout.Height, e.OldLayout.Left, e.OldLayout.Top,
			e.NewLayout.Width, e.NewLayout.Height, e.NewLayout.Left, e.NewLayout.Top)
	})

	pr.renderer.RegisterHandler(RenderEventChildrenChange, func(e RenderEventData) {
		pr.outputChan <- fmt.Sprintf("[CHILDREN] %s: %s", pr.elementType(e.Node), e.ChangeType)
	})

	pr.renderer.RegisterHandler(RenderEventStyleChange, func(e RenderEventData) {
		pr.outputChan <- fmt.Sprintf("[STYLE] %s: %s", pr.elementType(e.Node), e.ChangeType)
	})

	pr.renderer.RegisterHandler(RenderEventPropsChange, func(e RenderEventData) {
		pr.outputChan <- fmt.Sprintf("[PROPS] %s: %s", pr.elementType(e.Node), e.ChangeType)
	})

	return pr
}

func (pr *PrinterRenderer) elementType(node *RenderNode) string {
	if node == nil || node.element == nil {
		return "nil"
	}
	switch node.element.(type) {
	case *ViewElement:
		return "View"
	case *TextElement:
		return "Text"
	default:
		return fmt.Sprintf("%T", node.element)
	}
}

func (pr *PrinterRenderer) SetRoot(element Element) {
	pr.renderer.SetRoot(element)
}

func (pr *PrinterRenderer) Start() {
	pr.renderer.Start()
	pr.isRunning.Store(true)

	go pr.printLoop()
}

func (pr *PrinterRenderer) Stop() {
	if !pr.isRunning.Load() {
		return
	}

	pr.isRunning.Store(false)
	close(pr.stopChan)
	pr.renderer.Stop()
}

func (pr *PrinterRenderer) printLoop() {
	for {
		select {
		case <-pr.stopChan:
			return
		case msg := <-pr.outputChan:
			fmt.Println(msg)
		case <-time.After(100 * time.Millisecond):
			pr.renderer.PrintTree()
		}
	}
}

func (pr *PrinterRenderer) GetRenderer() *Renderer {
	return pr.renderer
}

func (pr *PrinterRenderer) ForceRender() {
	pr.renderer.ForceRender()
}

func StartLayoutMonitorLoop(rootElement Element, interval time.Duration) *Renderer {
	config := DefaultRendererConfig()
	renderer := NewRenderer(config)
	renderer.SetRoot(rootElement)

	mapper := renderer.GetMapper()
	mapper.CalculateLayout()

	renderer.Start()

	go func() {
		for {
			time.Sleep(interval)
			if !renderer.isRunning.Load() {
				break
			}

			oldLayouts := collectLayouts(mapper.GetRootNode())
			mapper.CalculateLayout()
			newLayouts := collectLayouts(mapper.GetRootNode())

			changes := compareLayouts(oldLayouts, newLayouts)
			if len(changes) > 0 {
				fmt.Printf("\n[Layout Change at %s]\n", time.Now().Format("15:04:05.000"))
				for _, change := range changes {
					fmt.Printf("  %s\n", change)
				}
				renderer.PrintTree()
			}
		}
	}()

	return renderer
}

func collectLayouts(node *RenderNode) map[*RenderNode]LayoutResults {
	layouts := make(map[*RenderNode]LayoutResults)

	var walk func(n *RenderNode)
	walk = func(n *RenderNode) {
		if n == nil {
			return
		}
		layouts[n] = n.GetLayout()
		for _, child := range n.children {
			walk(child)
		}
	}

	walk(node)
	return layouts
}

func compareLayouts(old, new map[*RenderNode]LayoutResults) []string {
	changes := make([]string, 0)

	for node, oldLayout := range old {
		if newLayout, exists := new[node]; exists {
			if oldLayout.Width != newLayout.Width ||
				oldLayout.Height != newLayout.Height ||
				oldLayout.Left != newLayout.Left ||
				oldLayout.Top != newLayout.Top {

				changes = append(changes, fmt.Sprintf("%s: %.2fx%.2f@%.2f,%.2f -> %.2fx%.2f@%.2f,%.2f",
					elementTypeName(node),
					oldLayout.Width, oldLayout.Height, oldLayout.Left, oldLayout.Top,
					newLayout.Width, newLayout.Height, newLayout.Left, newLayout.Top))
			}
		}
	}

	return changes
}

func elementTypeName(node *RenderNode) string {
	if node == nil {
		return "nil"
	}
	switch node.element.(type) {
	case *ViewElement:
		return "View"
	case *TextElement:
		if t, ok := node.element.(*TextElement); ok {
			if len(t.Text()) > 20 {
				return fmt.Sprintf("Text(%q...)", t.Text()[:20])
			}
			return fmt.Sprintf("Text(%q)", t.Text())
		}
		return "Text"
	default:
		return fmt.Sprintf("%T", node.element)
	}
}


