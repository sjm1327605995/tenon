package ui

import (
	"testing"
	"time"

	"github.com/sjm1327605995/tenon/yoga"
)

func TestLayoutMonitorCreation(t *testing.T) {
	monitor := GetLayoutMonitor()
	if monitor == nil {
		t.Fatal("GetLayoutMonitor returned nil")
	}
}

func TestSubscribe(t *testing.T) {
	monitor := GetLayoutMonitor()

	handler := func(e LayoutEventData) {}

	id := monitor.Subscribe("test-observer", handler, []LayoutEvent{LayoutChanged})
	if id == "" {
		t.Errorf("expected non-empty id")
	}

	monitor.Unsubscribe(id)
}

func TestEmitLayoutChange(t *testing.T) {
	monitor := GetLayoutMonitor()

	layoutChanged := false
	handler := func(e LayoutEventData) {
		if e.Event == LayoutChanged {
			layoutChanged = true
		}
	}

	id := monitor.Subscribe("test", handler, []LayoutEvent{LayoutChanged})
	defer monitor.Unsubscribe(id)

	monitor.EmitLayoutChange(nil, LayoutResults{}, LayoutResults{})

	if !layoutChanged {
		t.Errorf("expected LayoutChanged event")
	}
}

func TestEmitLoadChange(t *testing.T) {
	monitor := GetLayoutMonitor()

	loadChanged := false
	handler := func(e LayoutEventData) {
		if e.Event == LoadChanged {
			loadChanged = true
		}
	}

	id := monitor.Subscribe("test", handler, []LayoutEvent{LoadChanged})
	defer monitor.Unsubscribe(id)

	monitor.EmitLoadChange(nil)

	if !loadChanged {
		t.Errorf("expected LoadChanged event")
	}
}

func TestUnsubscribe(t *testing.T) {
	monitor := GetLayoutMonitor()

	handler := func(e LayoutEventData) {}
	id := monitor.Subscribe("test", handler, []LayoutEvent{LayoutChanged})
	monitor.Unsubscribe(id)
}

func TestLayoutChangeListenerCreation(t *testing.T) {
	node := yoga.NewNode()
	node.StyleSetWidth(100)
	node.StyleSetHeight(100)

	listener := NewLayoutChangeListener(node, nil)
	if listener == nil {
		t.Fatal("NewLayoutChangeListener returned nil")
	}
	if listener.GetNode() != node {
		t.Errorf("expected node to match")
	}
}

func TestLayoutChangeListenerStartStop(t *testing.T) {
	node := yoga.NewNode()
	node.StyleSetWidth(100)

	listener := NewLayoutChangeListener(node, nil)

	listener.Start()
	if !listener.isListening.Load() {
		t.Errorf("listener should be listening after Start()")
	}

	listener.Stop()
	if listener.isListening.Load() {
		t.Errorf("listener should not be listening after Stop()")
	}
}

func TestLayoutChangeListenerCheckAndNotify(t *testing.T) {
	node := yoga.NewNode()
	node.StyleSetWidth(100)
	node.StyleSetHeight(100)
	node.CalculateLayout(375, 812, yoga.DirectionLTR)

	changed := false
	listener := NewLayoutChangeListener(node, func(layout interface{}) {
		changed = true
	})

	listener.Start()
	listener.CheckAndNotify()

	node.StyleSetWidth(200)
	node.MarkDirty()
	node.CalculateLayout(375, 812, yoga.DirectionLTR)

	listener.CheckAndNotify()

	if !changed {
		t.Errorf("expected layout change callback to be called")
	}
}

func TestLayoutChangeListenerLayoutsEqual(t *testing.T) {
	node := yoga.NewNode()
	listener := NewLayoutChangeListener(node, nil)

	layoutA := LayoutResults{Width: 100, Height: 100}
	layoutB := LayoutResults{Width: 100, Height: 100}
	layoutC := LayoutResults{Width: 200, Height: 100}

	if !listener.layoutsEqual(layoutA, layoutB) {
		t.Errorf("expected layoutA and layoutB to be equal")
	}
	if listener.layoutsEqual(layoutA, layoutC) {
		t.Errorf("expected layoutA and layoutC to not be equal")
	}
}

func TestChangeTrackerCreation(t *testing.T) {
	tracker := NewChangeTracker(10)
	if tracker == nil {
		t.Fatal("NewChangeTracker returned nil")
	}
}

func TestChangeTrackerStartStop(t *testing.T) {
	tracker := NewChangeTracker(10)

	tracker.Start()
	if !tracker.isTracking.Load() {
		t.Errorf("tracker should be tracking after Start()")
	}

	tracker.Stop()
	if tracker.isTracking.Load() {
		t.Errorf("tracker should not be tracking after Stop()")
	}
}

func TestChangeTrackerRecord(t *testing.T) {
	tracker := NewChangeTracker(10)
	tracker.Start()

	event := LayoutEventData{
		Event:     LayoutChanged,
		Timestamp: time.Now(),
	}

	tracker.Record(event)

	changes := tracker.GetChanges()
	if len(changes) != 1 {
		t.Errorf("expected 1 change, got %d", len(changes))
	}
}

func TestChangeTrackerMaxChanges(t *testing.T) {
	tracker := NewChangeTracker(3)
	tracker.Start()

	for i := 0; i < 5; i++ {
		tracker.Record(LayoutEventData{Event: LayoutChanged})
	}

	changes := tracker.GetChanges()
	if len(changes) != 3 {
		t.Errorf("expected max 3 changes, got %d", len(changes))
	}
}

func TestChangeTrackerClear(t *testing.T) {
	tracker := NewChangeTracker(10)
	tracker.Start()

	tracker.Record(LayoutEventData{Event: LayoutChanged})
	tracker.Clear()

	if len(tracker.GetChanges()) != 0 {
		t.Errorf("expected 0 changes after Clear()")
	}
}

func TestLayoutWatcherCreation(t *testing.T) {
	watcher := NewLayoutWatcher(50 * time.Millisecond)
	if watcher == nil {
		t.Fatal("NewLayoutWatcher returned nil")
	}
}

func TestLayoutWatcherAddRemoveNode(t *testing.T) {
	watcher := NewLayoutWatcher(50 * time.Millisecond)

	node1 := yoga.NewNode()
	node2 := yoga.NewNode()

	watcher.AddNode(node1, nil)
	watcher.AddNode(node2, nil)

	if len(watcher.listeners) != 2 {
		t.Errorf("expected 2 listeners")
	}

	watcher.RemoveNode(node1)
	if len(watcher.listeners) != 1 {
		t.Errorf("expected 1 listener after removal")
	}
}

func TestLayoutWatcherStartStop(t *testing.T) {
	watcher := NewLayoutWatcher(50 * time.Millisecond)

	node := yoga.NewNode()
	watcher.AddNode(node, nil)

	watcher.Start()
	if !watcher.isRunning.Load() {
		t.Errorf("watcher should be running after Start()")
	}

	watcher.Stop()
	if watcher.isRunning.Load() {
		t.Errorf("watcher should not be running after Stop()")
	}
}

func TestGetDefaultLayoutWatcher(t *testing.T) {
	watcher := GetDefaultLayoutWatcher()
	if watcher == nil {
		t.Fatal("GetDefaultLayoutWatcher returned nil")
	}
}

func TestBoolAtomic(t *testing.T) {
	b := &Bool{}
	b.Store(true)
	if !b.Load() {
		t.Errorf("expected true")
	}
	b.Store(false)
	if b.Load() {
		t.Errorf("expected false")
	}
}


