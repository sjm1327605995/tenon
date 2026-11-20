package yoga

import (
	"sync"
	"sync/atomic"
	"testing"
)

func eventIndex(events []EventType, target EventType) int {
	for i, event := range events {
		if event == target {
			return i
		}
	}
	return -1
}

func TestEventLifecycleAndLayoutData(t *testing.T) {
	ResetEventSubscribers()
	t.Cleanup(ResetEventSubscribers)
	var events []EventType
	var endData LayoutData
	var measure MeasureCallbackData
	unsubscribe := SubscribeAll(func(event Event) {
		events = append(events, event.Type)
		switch event.Type {
		case EventLayoutPassEnd:
			endData = event.LayoutData
		case EventMeasureCallbackEnd:
			measure = event.Measure
		}
	})
	defer unsubscribe()

	config := NewConfig(nil)
	root := NewNodeWithConfig(config)
	root.StyleSetWidth(100)
	root.StyleSetHeight(20)
	leaf := NewNodeWithConfig(config)
	leaf.StyleSetAlignSelf(AlignFlexStart)
	leaf.SetMeasureFunc(func(_ *Node, _ float32, _ MeasureMode, _ float32, _ MeasureMode) Size {
		return Size{Width: 40, Height: 10}
	})
	root.InsertChild(leaf, 0)
	root.CalculateLayout(100, 20, DirectionLTR)
	root.FreeRecursive()

	ordered := []EventType{
		EventNodeAllocation,
		EventLayoutPassStart,
		EventMeasureCallbackStart,
		EventMeasureCallbackEnd,
		EventNodeLayout,
		EventLayoutPassEnd,
		EventNodeDeallocation,
	}
	previous := -1
	for _, eventType := range ordered {
		index := eventIndex(events, eventType)
		if index <= previous {
			t.Fatalf("event %v is missing or out of order: %v", eventType, events)
		}
		previous = index
	}
	if measure.MeasuredWidth != 40 || measure.MeasuredHeight != 10 {
		t.Fatalf("measure event data=%+v, want 40x10", measure)
	}
	if endData.MeasureCallbacks() == 0 || endData.Layouts() == 0 {
		t.Fatalf("layout event stats=%+v, want non-zero callbacks and layouts", endData)
	}
}

func TestBaselineEvents(t *testing.T) {
	ResetEventSubscribers()
	t.Cleanup(ResetEventSubscribers)
	var events []EventType
	var baseline float32
	n := NewNode()
	SubscribeAll(func(event Event) {
		events = append(events, event.Type)
		if event.Type == EventNodeBaselineEnd {
			baseline = event.Baseline
		}
	})
	n.SetBaselineFunc(func(_ *Node, _, _ float32) float32 { return 12 })

	if got := calculateBaseline(n); got != 12 {
		t.Fatalf("baseline=%v, want 12", got)
	}
	if baseline != 12 || len(events) < 2 ||
		events[0] != EventNodeBaselineStart || events[1] != EventNodeBaselineEnd {
		t.Fatalf("baseline events=%v, value=%v", events, baseline)
	}
}

func TestSubscriptionFilteringAndIdempotentUnsubscribe(t *testing.T) {
	ResetEventSubscribers()
	t.Cleanup(ResetEventSubscribers)
	calls := 0
	unsubscribe := Subscribe(EventNodeAllocation, func(Event) { calls++ })
	NewNode()
	NewNode().CalculateLayout(10, 10, DirectionLTR)
	if calls != 2 {
		t.Fatalf("allocation subscriber calls=%d, want 2", calls)
	}
	unsubscribe()
	unsubscribe()
	NewNode()
	if calls != 2 {
		t.Fatal("unsubscribed callback was invoked")
	}
	ResetEventSubscribers()
	unsubscribe()
	if got := eventSubscriberCount.Load(); got != 0 {
		t.Fatalf("subscriber count=%d after reset and repeated unsubscribe", got)
	}
}

func TestConcurrentSubscriptionsAndLayouts(t *testing.T) {
	ResetEventSubscribers()
	t.Cleanup(ResetEventSubscribers)
	var calls atomic.Int64
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 25; j++ {
				unsubscribe := Subscribe(EventLayoutPassEnd, func(Event) { calls.Add(1) })
				root := sizedNode(20, 20)
				root.CalculateLayout(20, 20, DirectionLTR)
				unsubscribe()
			}
		}()
	}
	wg.Wait()
	if calls.Load() == 0 {
		t.Fatal("concurrent subscribers received no layout events")
	}
	if got := eventSubscriberCount.Load(); got != 0 {
		t.Fatalf("subscriber count=%d, want 0", got)
	}
}
