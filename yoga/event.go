package yoga

import (
	"sync"
	"sync/atomic"
)

type LayoutType int

const (
	LayoutTypeLayout LayoutType = iota
	LayoutTypeMeasure
	LayoutTypeCachedLayout
	LayoutTypeCachedMeasure
)

type LayoutPassReason uint8

const (
	LayoutPassReasonInitial LayoutPassReason = iota
	LayoutPassReasonAbsLayout
	LayoutPassReasonStretch
	LayoutPassReasonMultilineStretch
	LayoutPassReasonFlexLayout
	LayoutPassReasonMeasureChild
	LayoutPassReasonAbsMeasureChild
	LayoutPassReasonFlexMeasure
	LayoutPassReasonCount
)

func (r LayoutPassReason) String() string {
	switch r {
	case LayoutPassReasonInitial:
		return "initial"
	case LayoutPassReasonAbsLayout:
		return "abs_layout"
	case LayoutPassReasonStretch:
		return "stretch"
	case LayoutPassReasonMultilineStretch:
		return "multiline_stretch"
	case LayoutPassReasonFlexLayout:
		return "flex_layout"
	case LayoutPassReasonMeasureChild:
		return "measure"
	case LayoutPassReasonAbsMeasureChild:
		return "abs_measure"
	case LayoutPassReasonFlexMeasure:
		return "flex_measure"
	default:
		return "unknown"
	}
}

// LayoutData is a snapshot of work performed by one layout pass.
type LayoutData struct {
	layouts                     int32
	measures                    int32
	maxMeasureCache             uint32
	cachedLayouts               int32
	cachedMeasures              int32
	measureCallbacks            int32
	measureCallbackReasonsCount [LayoutPassReasonCount]uint32
}

func (d LayoutData) Layouts() int32          { return d.layouts }
func (d LayoutData) Measures() int32         { return d.measures }
func (d LayoutData) MaxMeasureCache() uint32 { return d.maxMeasureCache }
func (d LayoutData) CachedLayouts() int32    { return d.cachedLayouts }
func (d LayoutData) CachedMeasures() int32   { return d.cachedMeasures }
func (d LayoutData) MeasureCallbacks() int32 { return d.measureCallbacks }
func (d LayoutData) MeasureCallbacksFor(reason LayoutPassReason) uint32 {
	if reason >= LayoutPassReasonCount {
		return 0
	}
	return d.measureCallbackReasonsCount[reason]
}

// EventType identifies a node lifecycle or layout event.
type EventType uint8

const (
	EventNodeAllocation EventType = iota
	EventNodeDeallocation
	EventNodeLayout
	EventLayoutPassStart
	EventLayoutPassEnd
	EventMeasureCallbackStart
	EventMeasureCallbackEnd
	EventNodeBaselineStart
	EventNodeBaselineEnd
)

func (t EventType) String() string {
	switch t {
	case EventNodeAllocation:
		return "node_allocation"
	case EventNodeDeallocation:
		return "node_deallocation"
	case EventNodeLayout:
		return "node_layout"
	case EventLayoutPassStart:
		return "layout_pass_start"
	case EventLayoutPassEnd:
		return "layout_pass_end"
	case EventMeasureCallbackStart:
		return "measure_callback_start"
	case EventMeasureCallbackEnd:
		return "measure_callback_end"
	case EventNodeBaselineStart:
		return "node_baseline_start"
	case EventNodeBaselineEnd:
		return "node_baseline_end"
	default:
		return "unknown"
	}
}

// MeasureCallbackData contains the constraints and result of a custom measurement.
type MeasureCallbackData struct {
	Width, Height                 float32
	WidthMode, HeightMode         MeasureMode
	MeasuredWidth, MeasuredHeight float32
	Reason                        LayoutPassReason
}

// Event is an immutable notification delivered synchronously to subscribers.
type Event struct {
	Type       EventType
	Node       *Node
	Config     *Config
	LayoutType LayoutType
	LayoutData LayoutData
	Measure    MeasureCallbackData
	Baseline   float32
}

// Subscriber handles an event. It may subscribe or unsubscribe from inside the callback.
type Subscriber func(Event)

var eventSubscriptions = struct {
	sync.RWMutex
	nextID uint64
	typed  map[EventType]map[uint64]Subscriber
	all    map[uint64]Subscriber
}{
	typed: make(map[EventType]map[uint64]Subscriber),
	all:   make(map[uint64]Subscriber),
}

var eventSubscriberCount atomic.Int64

func addEventSubscriber(eventType EventType, all bool, subscriber Subscriber) func() {
	if subscriber == nil {
		return func() {}
	}
	eventSubscriptions.Lock()
	eventSubscriptions.nextID++
	id := eventSubscriptions.nextID
	if all {
		eventSubscriptions.all[id] = subscriber
	} else {
		subscribers := eventSubscriptions.typed[eventType]
		if subscribers == nil {
			subscribers = make(map[uint64]Subscriber)
			eventSubscriptions.typed[eventType] = subscribers
		}
		subscribers[id] = subscriber
	}
	eventSubscriptions.Unlock()
	eventSubscriberCount.Add(1)

	var once sync.Once
	return func() {
		once.Do(func() {
			removed := false
			eventSubscriptions.Lock()
			if all {
				if _, ok := eventSubscriptions.all[id]; ok {
					delete(eventSubscriptions.all, id)
					removed = true
				}
			} else if subscribers := eventSubscriptions.typed[eventType]; subscribers != nil {
				if _, ok := subscribers[id]; ok {
					delete(subscribers, id)
					removed = true
				}
			}
			eventSubscriptions.Unlock()
			if removed {
				eventSubscriberCount.Add(-1)
			}
		})
	}
}

// Subscribe registers for one event type and returns an idempotent unsubscribe function.
func Subscribe(eventType EventType, subscriber Subscriber) func() {
	return addEventSubscriber(eventType, false, subscriber)
}

// SubscribeAll registers for every Yoga event.
func SubscribeAll(subscriber Subscriber) func() {
	return addEventSubscriber(0, true, subscriber)
}

// ResetEventSubscribers removes all subscriptions. It is intended for tests and shutdown.
func ResetEventSubscribers() {
	eventSubscriptions.Lock()
	eventSubscriptions.typed = make(map[EventType]map[uint64]Subscriber)
	eventSubscriptions.all = make(map[uint64]Subscriber)
	eventSubscriptions.Unlock()
	eventSubscriberCount.Store(0)
}

func publishEvent(event Event) {
	if eventSubscriberCount.Load() == 0 {
		return
	}
	eventSubscriptions.RLock()
	subscribers := make([]Subscriber, 0, len(eventSubscriptions.typed[event.Type])+len(eventSubscriptions.all))
	for _, subscriber := range eventSubscriptions.typed[event.Type] {
		subscribers = append(subscribers, subscriber)
	}
	for _, subscriber := range eventSubscriptions.all {
		subscribers = append(subscribers, subscriber)
	}
	eventSubscriptions.RUnlock()
	for _, subscriber := range subscribers {
		subscriber(event)
	}
}
