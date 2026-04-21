package ui

import (
	"fmt"
	"sync"
	"time"

	"github.com/sjm1327605995/tenon/yoga"
)

type LayoutEvent int

const (
	LayoutChanged LayoutEvent = iota
	LoadChanged
	StyleChanged
	NodeAdded
	NodeRemoved
)

type LayoutEventData struct {
	Event     LayoutEvent
	Node      *yoga.Node
	OldLayout LayoutResults
	NewLayout LayoutResults
	Timestamp time.Time
}

type LayoutEventHandler func(LayoutEventData)
type ChangeHandler func(interface{})

type LayoutObserver struct {
	id      string
	handler LayoutEventHandler
	events  []LayoutEvent
}

type LayoutMonitor struct {
	mu         sync.RWMutex
	observers  map[string]*LayoutObserver
	eventQueue chan LayoutEventData
	isRunning  Bool
	stopChan   chan struct{}
	runningWG  sync.WaitGroup
}

var layoutMonitor *LayoutMonitor

func init() {
	layoutMonitor = &LayoutMonitor{
		observers:  make(map[string]*LayoutObserver),
		eventQueue: make(chan LayoutEventData, 1000),
		stopChan:   make(chan struct{}),
	}
}

func GetLayoutMonitor() *LayoutMonitor {
	return layoutMonitor
}

func (m *LayoutMonitor) Subscribe(id string, handler LayoutEventHandler, events []LayoutEvent) string {
	m.mu.Lock()
	defer m.mu.Unlock()

	observer := &LayoutObserver{
		id:      id,
		handler: handler,
		events:  events,
	}
	m.observers[id] = observer
	return id
}

func (m *LayoutMonitor) Unsubscribe(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.observers, id)
}

func (m *LayoutMonitor) Notify(eventData LayoutEventData) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, observer := range m.observers {
		for _, event := range observer.events {
			if event == eventData.Event || event == 0 {
				observer.handler(eventData)
				break
			}
		}
	}
}

func (m *LayoutMonitor) EmitLayoutChange(node *yoga.Node, oldLayout, newLayout LayoutResults) {
	eventData := LayoutEventData{
		Event:     LayoutChanged,
		Node:      node,
		OldLayout: oldLayout,
		NewLayout: newLayout,
		Timestamp: time.Now(),
	}
	m.Notify(eventData)
}

func (m *LayoutMonitor) EmitLoadChange(node *yoga.Node) {
	eventData := LayoutEventData{
		Event:     LoadChanged,
		Node:      node,
		Timestamp: time.Now(),
	}
	m.Notify(eventData)
}

func (m *LayoutMonitor) StartMonitoring(interval time.Duration) {
	if m.isRunning.Load() {
		return
	}

	m.isRunning.Store(true)
	m.runningWG.Add(1)

	go m.eventLoop(interval)
}

func (m *LayoutMonitor) StopMonitoring() {
	if !m.isRunning.Load() {
		return
	}

	m.isRunning.Store(false)
	close(m.stopChan)
	m.runningWG.Wait()
}

func (m *LayoutMonitor) eventLoop(interval time.Duration) {
	defer m.runningWG.Done()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-m.stopChan:
			return
		case eventData := <-m.eventQueue:
			m.processEvent(eventData)
		case <-ticker.C:
			m.checkForChanges()
		}
	}
}

func (m *LayoutMonitor) processEvent(eventData LayoutEventData) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, observer := range m.observers {
		for _, event := range observer.events {
			if event == eventData.Event || event == 0 {
				observer.handler(eventData)
				break
			}
		}
	}
}

func (m *LayoutMonitor) checkForChanges() {
	m.mu.RLock()
	defer m.mu.RUnlock()
}

type LayoutChangeListener struct {
	monitor       *LayoutMonitor
	node          *yoga.Node
	oldLayout     LayoutResults
	changeHandler ChangeHandler
	isListening   Bool
}

func NewLayoutChangeListener(node *yoga.Node, handler ChangeHandler) *LayoutChangeListener {
	return &LayoutChangeListener{
		monitor:       layoutMonitor,
		node:          node,
		changeHandler: handler,
		oldLayout:     getNodeLayout(node),
	}
}

func (l *LayoutChangeListener) Start() {
	if l.isListening.Load() {
		return
	}
	l.isListening.Store(true)
}

func (l *LayoutChangeListener) Stop() {
	if !l.isListening.Load() {
		return
	}
	l.isListening.Store(false)
}

func (l *LayoutChangeListener) CheckAndNotify() {
	if !l.isListening.Load() {
		return
	}

	newLayout := getNodeLayout(l.node)

	if !l.layoutsEqual(l.oldLayout, newLayout) {
		eventData := LayoutEventData{
			Event:     LayoutChanged,
			Node:      l.node,
			OldLayout: l.oldLayout,
			NewLayout: newLayout,
			Timestamp: time.Now(),
		}

		l.monitor.Notify(eventData)

		if l.changeHandler != nil {
			l.changeHandler(newLayout)
		}

		l.oldLayout = newLayout
	}
}

func (l *LayoutChangeListener) layoutsEqual(a, b LayoutResults) bool {
	return a.Width == b.Width && a.Height == b.Height
}

func (l *LayoutChangeListener) GetNode() *yoga.Node {
	return l.node
}

func (l *LayoutChangeListener) GetOldLayout() LayoutResults {
	return l.oldLayout
}

func (l *LayoutChangeListener) GetCurrentLayout() LayoutResults {
	return getNodeLayout(l.node)
}

func StartLayoutEventLoop(rootNode *yoga.Node, interval time.Duration, callback func(LayoutEventData)) {
	monitor := GetLayoutMonitor()

	monitor.Subscribe("layout-loop-"+fmt.Sprint(time.Now().UnixNano()), func(eventData LayoutEventData) {
		callback(eventData)
	}, []LayoutEvent{LayoutChanged, LoadChanged, StyleChanged, NodeAdded, NodeRemoved})

	monitor.StartMonitoring(interval)

	listener := NewLayoutChangeListener(rootNode, func(layout interface{}) {
		_ = layout
	})
	listener.Start()

	go func() {
		for {
			if !monitor.isRunning.Load() {
				break
			}
			listener.CheckAndNotify()
			time.Sleep(interval)
		}
	}()
}

func PrintLayoutInfo(node *yoga.Node, indent int) {
	layout := getNodeLayout(node)
	style := node.StyleGetFlexDirection()

	prefix := ""
	for i := 0; i < indent; i++ {
		prefix += "  "
	}

	fmt.Printf("%sNode Layout: %.2f x %.2f at (%.2f, %.2f)\n",
		prefix, layout.Width, layout.Height, layout.Left, layout.Top)

	fmt.Printf("%s  Flex Direction: %v\n", prefix, style)

	for i, child := range node.GetChildren() {
		fmt.Printf("%s  Child %d:\n", prefix, i)
		PrintLayoutInfo(child, indent+2)
	}
}

type ChangeTracker struct {
	mu         sync.RWMutex
	changes    []LayoutEventData
	maxChanges int
	isTracking Bool
}

func NewChangeTracker(maxChanges int) *ChangeTracker {
	return &ChangeTracker{
		changes:    make([]LayoutEventData, 0, maxChanges),
		maxChanges: maxChanges,
	}
}

func (t *ChangeTracker) Start() {
	t.isTracking.Store(true)
}

func (t *ChangeTracker) Stop() {
	t.isTracking.Store(false)
}

func (t *ChangeTracker) Record(eventData LayoutEventData) {
	if !t.isTracking.Load() {
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	t.changes = append(t.changes, eventData)

	if len(t.changes) > t.maxChanges {
		t.changes = t.changes[1:]
	}
}

func (t *ChangeTracker) GetChanges() []LayoutEventData {
	t.mu.RLock()
	defer t.mu.RUnlock()

	result := make([]LayoutEventData, len(t.changes))
	copy(result, t.changes)
	return result
}

func (t *ChangeTracker) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.changes = t.changes[:0]
}

type LayoutWatcher struct {
	monitor   *LayoutMonitor
	listeners map[*yoga.Node]*LayoutChangeListener
	interval  time.Duration
	stopChan  chan struct{}
	isRunning Bool
	runningWG sync.WaitGroup
	mu        sync.RWMutex
}

func NewLayoutWatcher(interval time.Duration) *LayoutWatcher {
	return &LayoutWatcher{
		monitor:   GetLayoutMonitor(),
		listeners: make(map[*yoga.Node]*LayoutChangeListener),
		interval:  interval,
		stopChan:  make(chan struct{}),
	}
}

func (w *LayoutWatcher) AddNode(node *yoga.Node, handler ChangeHandler) {
	w.mu.Lock()
	defer w.mu.Unlock()

	listener := NewLayoutChangeListener(node, handler)
	w.listeners[node] = listener
	listener.Start()
}

func (w *LayoutWatcher) RemoveNode(node *yoga.Node) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if listener, exists := w.listeners[node]; exists {
		listener.Stop()
		delete(w.listeners, node)
	}
}

func (w *LayoutWatcher) Start() {
	if w.isRunning.Load() {
		return
	}
	w.isRunning.Store(true)

	w.runningWG.Add(1)
	go w.watchLoop()
}

func (w *LayoutWatcher) Stop() {
	if !w.isRunning.Load() {
		return
	}
	w.isRunning.Store(false)
	close(w.stopChan)
	w.runningWG.Wait()
}

func (w *LayoutWatcher) watchLoop() {
	defer w.runningWG.Done()

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-w.stopChan:
			return
		case <-ticker.C:
			w.checkAllNodes()
		}
	}
}

func (w *LayoutWatcher) checkAllNodes() {
	w.mu.RLock()
	defer w.mu.RUnlock()

	for _, listener := range w.listeners {
		listener.CheckAndNotify()
	}
}

var defaultLayoutWatcher *LayoutWatcher

func init() {
	defaultLayoutWatcher = NewLayoutWatcher(100 * time.Millisecond)
}

func GetDefaultLayoutWatcher() *LayoutWatcher {
	return defaultLayoutWatcher
}


