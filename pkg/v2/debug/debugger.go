package debug

import (
	"encoding"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/sjm1327605995/tenon/pkg/v2/core"
)

func safeJSONMarshal(v interface{}) ([]byte, error) {
	v = preprocessJSON(v)
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func preprocessJSON(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	val := reflect.ValueOf(v)

	// Dereference pointers and interfaces first
	for val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}

	// If it implements json.Marshaler or encoding.TextMarshaler, let json.Marshal handle it
	if val.CanInterface() {
		iface := val.Interface()
		if _, ok := iface.(json.Marshaler); ok {
			return iface
		}
		if _, ok := iface.(encoding.TextMarshaler); ok {
			return iface
		}
	}

	switch val.Kind() {
	case reflect.Float64:
		f := val.Float()
		if math.IsNaN(f) || math.IsInf(f, 0) {
			return nil
		}
		return f
	case reflect.Float32:
		f := float64(val.Float())
		if math.IsNaN(f) || math.IsInf(f, 0) {
			return nil
		}
		return f
	case reflect.Slice:
		result := make([]interface{}, val.Len())
		for i := 0; i < val.Len(); i++ {
			result[i] = preprocessJSON(val.Index(i).Interface())
		}
		return result
	case reflect.Map:
		result := make(map[string]interface{})
		for _, key := range val.MapKeys() {
			result[key.String()] = preprocessJSON(val.MapIndex(key).Interface())
		}
		return result
	case reflect.Struct:
		result := make(map[string]interface{})
		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			fieldValue := val.Field(i)
			if !fieldValue.CanInterface() {
				continue
			}
			// Use json tag as key if available
			name := field.Name
			if tag := field.Tag.Get("json"); tag != "" {
				if idx := strings.Index(tag, ","); idx != -1 {
					name = tag[:idx]
				} else {
					name = tag
				}
				if name == "-" {
					continue
				}
				if name == "" {
					name = field.Name
				}
			}
			result[name] = preprocessJSON(fieldValue.Interface())
		}
		return result
	default:
		return val.Interface()
	}
}

func jsonMarshal(v interface{}) ([]byte, error) {
	return safeJSONMarshal(v)
}

type Debugger struct {
	engine     *core.Engine
	server     *http.Server
	port       int
	mu         sync.RWMutex
	snapshots  []*LayoutSnapshot
	maxHistory int
	enabled    bool

	eventLogs []EventLog
	maxEvents int

	stateLogs []StateChangeLog
	maxStates int

	hub           *WebSocketHub
	highlightPath []int
}

type StateChangeLog struct {
	Timestamp   time.Time   `json:"timestamp"`
	ElementType string      `json:"elementType"`
	ElementKey  string      `json:"elementKey,omitempty"`
	StateKey    string      `json:"stateKey,omitempty"`
	OldValue    interface{} `json:"oldValue"`
	NewValue    interface{} `json:"newValue"`
}

type EventLog struct {
	Timestamp time.Time              `json:"timestamp"`
	Type      string                 `json:"type"`
	Target    string                 `json:"target,omitempty"`
	X         float32                `json:"x,omitempty"`
	Y         float32                `json:"y,omitempty"`
	DeltaX    float32                `json:"deltaX,omitempty"`
	DeltaY    float32                `json:"deltaY,omitempty"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

func (e EventLog) GetType() string    { return e.Type }
func (e EventLog) GetTarget() string  { return e.Target }
func (e EventLog) GetX() float32      { return e.X }
func (e EventLog) GetY() float32      { return e.Y }
func (e EventLog) GetDeltaX() float32 { return e.DeltaX }
func (e EventLog) GetDeltaY() float32 { return e.DeltaY }

type LayoutSnapshot struct {
	Timestamp time.Time       `json:"timestamp"`
	Root      *core.DebugNode `json:"root"`
	Trigger   string          `json:"trigger"`
	FrameID   int             `json:"frameId"`
}

var (
	globalDebugger *Debugger
	frameCounter   int
)

func NewDebugger(engine *core.Engine, port int) *Debugger {
	d := &Debugger{
		engine:     engine,
		port:       port,
		maxHistory: 100,
		snapshots:  make([]*LayoutSnapshot, 0),
		enabled:    true,
		maxEvents:  500,
		eventLogs:  make([]EventLog, 0),
		maxStates:  500,
		stateLogs:  make([]StateChangeLog, 0),
		hub:        NewWebSocketHub(),
	}
	globalDebugger = d
	return d
}

func GetDebugger() *Debugger {
	return globalDebugger
}

func (d *Debugger) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/tree", d.handleTree)
	mux.HandleFunc("/debug/html", d.handleHTML)
	mux.HandleFunc("/debug/history", d.handleHistory)
	mux.HandleFunc("/debug/snapshot", d.handleSnapshot)
	mux.HandleFunc("/debug/compare", d.handleCompare)
	mux.HandleFunc("/debug/events", d.handleEvents)
	mux.HandleFunc("/debug/live", d.handleLive)
	mux.HandleFunc("/debug/perf", d.handlePerf)
	mux.HandleFunc("/debug/listeners", d.handleListeners)
	mux.HandleFunc("/debug/lifecycle", d.handleLifecycle)
	mux.HandleFunc("/debug/state", d.handleStates)
	mux.HandleFunc("/debug/ws", d.handleWS)
	mux.HandleFunc("/debug/devtools", d.handleDevTools)
	mux.HandleFunc("/", d.handleIndex)

	d.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", d.port),
		Handler: mux,
	}

	go func() {
		fmt.Printf("[Debugger] HTTP server started at http://localhost:%d\n", d.port)
		fmt.Printf("[Debugger] DevTools: http://localhost:%d/debug/devtools\n", d.port)
		if err := d.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("[Debugger] Server error: %v\n", err)
		}
	}()

	d.startWSPushLoop()

	return nil
}

func (d *Debugger) Stop() error {
	if d.server != nil {
		return d.server.Close()
	}
	return nil
}

func (d *Debugger) AddEventLog(evt interface {
	GetType() string
	GetTarget() string
	GetX() float32
	GetY() float32
	GetDeltaX() float32
	GetDeltaY() float32
}) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if !d.enabled {
		return
	}
	details := make(map[string]interface{})
	if evt.GetX() != 0 || evt.GetY() != 0 {
		details["x"] = evt.GetX()
		details["y"] = evt.GetY()
	}
	if evt.GetDeltaX() != 0 {
		details["deltaX"] = evt.GetDeltaX()
	}
	if evt.GetDeltaY() != 0 {
		details["deltaY"] = evt.GetDeltaY()
	}
	entry := EventLog{
		Timestamp: time.Now(),
		Type:      evt.GetType(),
		Target:    evt.GetTarget(),
		X:         evt.GetX(),
		Y:         evt.GetY(),
		DeltaX:    evt.GetDeltaX(),
		DeltaY:    evt.GetDeltaY(),
		Details:   details,
	}
	d.eventLogs = append(d.eventLogs, entry)
	if len(d.eventLogs) > d.maxEvents {
		d.eventLogs = d.eventLogs[len(d.eventLogs)-d.maxEvents:]
	}
	d.sendWS("event", entry)
}

func (d *Debugger) AddStateLog(elementType, elementKey, stateKey string, oldValue, newValue interface{}) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if !d.enabled {
		return
	}
	entry := StateChangeLog{
		Timestamp:   time.Now(),
		ElementType: elementType,
		ElementKey:  elementKey,
		StateKey:    stateKey,
		OldValue:    oldValue,
		NewValue:    newValue,
	}
	d.stateLogs = append(d.stateLogs, entry)
	if len(d.stateLogs) > d.maxStates {
		d.stateLogs = d.stateLogs[len(d.stateLogs)-d.maxStates:]
	}
	d.sendWS("state", entry)
}

func (d *Debugger) CaptureLayout(trigger string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.enabled || d.engine == nil {
		return
	}

	root := d.engine.GetRootElement()
	if root == nil {
		return
	}

	frameCounter++
	info := root.DebugInfo()
	snapshot := &LayoutSnapshot{
		Timestamp: time.Now(),
		Root:      &info,
		Trigger:   trigger,
		FrameID:   frameCounter,
	}

	d.snapshots = append(d.snapshots, snapshot)
	if len(d.snapshots) > d.maxHistory {
		d.snapshots = d.snapshots[1:]
	}
}

func (d *Debugger) SetEnabled(enabled bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.enabled = enabled
}

func (d *Debugger) IsEnabled() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.enabled
}

func CaptureLayout(trigger string) {
	if d := GetDebugger(); d != nil {
		d.CaptureLayout(trigger)
	}
}

func LogEvent(evtType string, target string, details map[string]interface{}) {
	if d := GetDebugger(); d != nil {
		if details == nil {
			details = make(map[string]interface{})
		}
		d.AddEventLog(EventLog{
			Timestamp: time.Now(),
			Type:      evtType,
			Target:    target,
			Details:   details,
		})
	}
}

// ==================== HTTP Handlers ====================

func (d *Debugger) handleTree(w http.ResponseWriter, r *http.Request) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	root := d.engine.GetRootElement()
	if root == nil {
		http.Error(w, "No root element", http.StatusNotFound)
		return
	}

	info := root.DebugInfo()
	w.Header().Set("Content-Type", "application/json")
	data, err := jsonMarshal(info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (d *Debugger) handleHistory(w http.ResponseWriter, r *http.Request) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	summary := make([]map[string]interface{}, len(d.snapshots))
	for i, s := range d.snapshots {
		summary[i] = map[string]interface{}{
			"timestamp": s.Timestamp,
			"trigger":   s.Trigger,
			"frameId":   s.FrameID,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	data, err := jsonMarshal(summary)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (d *Debugger) handleSnapshot(w http.ResponseWriter, r *http.Request) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if len(d.snapshots) == 0 {
		http.Error(w, "No snapshots available", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	data, err := jsonMarshal(d.snapshots[len(d.snapshots)-1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (d *Debugger) handleCompare(w http.ResponseWriter, r *http.Request) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if len(d.snapshots) < 2 {
		http.Error(w, "Need at least 2 snapshots to compare", http.StatusBadRequest)
		return
	}

	prev := d.snapshots[len(d.snapshots)-2]
	curr := d.snapshots[len(d.snapshots)-1]

	diff := d.compareNodes(prev.Root, curr.Root)

	w.Header().Set("Content-Type", "application/json")
	data, err := jsonMarshal(diff)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (d *Debugger) handleEvents(w http.ResponseWriter, r *http.Request) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	limit := 100
	if l := r.URL.Query().Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if limit > len(d.eventLogs) {
		limit = len(d.eventLogs)
	}
	start := len(d.eventLogs) - limit
	if start < 0 {
		start = 0
	}

	w.Header().Set("Content-Type", "application/json")
	data, err := jsonMarshal(d.eventLogs[start:])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (d *Debugger) handlePerf(w http.ResponseWriter, r *http.Request) {
	perf := d.engine.GetPerfMetrics()
	w.Header().Set("Content-Type", "application/json")
	data, err := jsonMarshal(perf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (d *Debugger) handleListeners(w http.ResponseWriter, r *http.Request) {
	info := d.engine.GetEventRegistryDebugInfo()
	w.Header().Set("Content-Type", "application/json")
	data, err := jsonMarshal(info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (d *Debugger) handleLifecycle(w http.ResponseWriter, r *http.Request) {
	logs := d.engine.GetLifecycleLogs()
	limit := 100
	if l := r.URL.Query().Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	start := len(logs) - limit
	if start < 0 {
		start = 0
	}
	w.Header().Set("Content-Type", "application/json")
	data, err := jsonMarshal(logs[start:])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (d *Debugger) handleStates(w http.ResponseWriter, r *http.Request) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	limit := 100
	if l := r.URL.Query().Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if limit > len(d.stateLogs) {
		limit = len(d.stateLogs)
	}
	start := len(d.stateLogs) - limit
	if start < 0 {
		start = 0
	}

	w.Header().Set("Content-Type", "application/json")
	data, err := jsonMarshal(d.stateLogs[start:])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (d *Debugger) handleLive(w http.ResponseWriter, r *http.Request) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	root := d.engine.GetRootElement()
	if root == nil {
		http.Error(w, "No root element", http.StatusNotFound)
		return
	}

	info := root.DebugInfo()
	nodeIDCounter = 0
	assignIDs(&info)
	enrichEventCounts(&info, d.engine.GetEventRegistryDebugInfo())

	w.Header().Set("Content-Type", "application/json")
	data, err := jsonMarshal(info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (d *Debugger) handleHTML(w http.ResponseWriter, r *http.Request) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	root := d.engine.GetRootElement()
	if root == nil {
		http.Error(w, "No root element", http.StatusNotFound)
		return
	}

	html := d.GenerateHTML(root)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func (d *Debugger) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>Tenon Debugger</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; margin: 20px; background: #1a1a2e; color: #eee; }
        h1 { color: #00d9ff; }
        .endpoint { background: #16213e; padding: 15px; margin: 10px 0; border-radius: 8px; border-left: 4px solid #00d9ff; }
        .endpoint h3 { margin: 0 0 10px 0; color: #00d9ff; }
        .endpoint code { background: #0f0f23; padding: 2px 8px; border-radius: 4px; color: #7ee787; }
        .endpoint p { margin: 5px 0; color: #aaa; }
        a { color: #00d9ff; text-decoration: none; }
        a:hover { text-decoration: underline; }
        .devtools-link { display: inline-block; background: #00d9ff; color: #1a1a2e; padding: 12px 24px; border-radius: 8px; font-size: 18px; font-weight: bold; margin: 20px 0; }
        .devtools-link:hover { text-decoration: none; background: #33e5ff; }
    </style>
</head>
<body>
    <h1>Tenon Debugger</h1>
    <p>Interactive debugging interface for Tenon GUI framework</p>
    <a href="/debug/devtools" class="devtools-link">Open DevTools</a>
    <h2>API Endpoints</h2>
    <div class="endpoint">
        <h3><a href="/debug/devtools">/debug/devtools</a></h3>
        <code>GET</code>
        <p>Interactive DevTools panel with real-time tree view, event inspector, and performance monitor</p>
    </div>
    <div class="endpoint">
        <h3><a href="/debug/tree">/debug/tree</a></h3>
        <code>GET</code>
        <p>Current layout tree as JSON (full DebugInfo)</p>
    </div>
    <div class="endpoint">
        <h3><a href="/debug/perf">/debug/perf</a></h3>
        <code>GET</code>
        <p>Performance metrics (FPS, frame time, layout time, draw time, element count)</p>
    </div>
    <div class="endpoint">
        <h3><a href="/debug/listeners">/debug/listeners</a></h3>
        <code>GET</code>
        <p>Event listener registry - see which elements have which event callbacks</p>
    </div>
    <div class="endpoint">
        <h3><a href="/debug/lifecycle">/debug/lifecycle</a></h3>
        <code>GET</code>
        <p>Component lifecycle events (mount/unmount log)</p>
    </div>
    <div class="endpoint">
        <h3><a href="/debug/events">/debug/events</a></h3>
        <code>GET</code>
        <p>Recent UI event logs (click, scroll, keyboard, etc.)</p>
    </div>
    <div class="endpoint">
        <h3><a href="/debug/live">/debug/live</a></h3>
        <code>GET</code>
        <p>Live tree state as JSON (includes scroll offsets, real-time bounds)</p>
    </div>
    <div class="endpoint">
        <h3><a href="/debug/html">/debug/html</a></h3>
        <code>GET</code>
        <p>Visualize layout tree as HTML with CSS absolute positioning</p>
    </div>
    <div class="endpoint">
        <h3><a href="/debug/snapshot">/debug/snapshot</a></h3>
        <code>GET</code>
        <p>Latest layout snapshot with metadata</p>
    </div>
    <div class="endpoint">
        <h3><a href="/debug/compare">/debug/compare</a></h3>
        <code>GET</code>
        <p>Compare last two snapshots and show differences</p>
    </div>
    <div class="endpoint">
        <h3>/debug/ws</h3>
        <code>WebSocket</code>
        <p>Real-time push: tree, perf, events. Send commands: getTree, getPerf, getEvents, getListeners, getLifecycle, highlight</p>
    </div>
</body>
</html>`)))
}

// ==================== Diff ====================

type DiffNode struct {
	Type     string            `json:"type"`
	Key      string            `json:"key,omitempty"`
	Changed  bool              `json:"changed"`
	Changes  []string          `json:"changes"`
	Previous core.LayoutBounds `json:"previous,omitempty"`
	Current  core.LayoutBounds `json:"current,omitempty"`
	Children []*DiffNode       `json:"children,omitempty"`
}

func (d *Debugger) compareNodes(a, b *core.DebugNode) *DiffNode {
	if a == nil && b == nil {
		return nil
	}

	diff := &DiffNode{
		Type:    b.Type,
		Key:     b.Key,
		Changed: false,
		Changes: []string{},
	}

	if a == nil {
		diff.Changed = true
		diff.Changes = append(diff.Changes, "added")
		diff.Current = b.Bounds
		return diff
	}

	if b == nil {
		diff.Changed = true
		diff.Changes = append(diff.Changes, "removed")
		diff.Previous = a.Bounds
		return diff
	}

	if a.Bounds != b.Bounds {
		diff.Changed = true
		diff.Changes = append(diff.Changes, "bounds_changed")
		diff.Previous = a.Bounds
		diff.Current = b.Bounds
	}

	if a.Yoga.FlexGrow != b.Yoga.FlexGrow {
		diff.Changed = true
		diff.Changes = append(diff.Changes, fmt.Sprintf("flexGrow: %v -> %v", a.Yoga.FlexGrow, b.Yoga.FlexGrow))
	}

	if a.Yoga.FlexShrink != b.Yoga.FlexShrink {
		diff.Changed = true
		diff.Changes = append(diff.Changes, fmt.Sprintf("flexShrink: %v -> %v", a.Yoga.FlexShrink, b.Yoga.FlexShrink))
	}

	maxLen := len(a.Children)
	if len(b.Children) > maxLen {
		maxLen = len(b.Children)
	}

	for i := 0; i < maxLen; i++ {
		var childA, childB *core.DebugNode
		if i < len(a.Children) {
			childA = a.Children[i]
		}
		if i < len(b.Children) {
			childB = b.Children[i]
		}

		childDiff := d.compareNodes(childA, childB)
		if childDiff != nil && childDiff.Changed {
			diff.Children = append(diff.Children, childDiff)
		}
	}

	return diff
}

// ==================== HTML Visualization ====================

var nodeIDCounter int


func (d *Debugger) GenerateHTML(root core.Element) string {
	info := root.DebugInfo()
	bounds := info.Bounds
	screenW := int(bounds.Width)
	screenH := int(bounds.Height)
	if screenW == 0 {
		screenW = 1280
	}
	if screenH == 0 {
		screenH = 720
	}

	nodeIDCounter = 0

	var elementsHTML strings.Builder
	d.renderDebugNodeHTML(&info, &elementsHTML, 0, 0, 0)

	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>Tenon Layout Debugger</title>
    <style>
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body { background: #1a1a2e; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', monospace; padding: 20px; color: #eee; }
        .toolbar { position: fixed; top: 0; left: 0; right: 0; background: #16213e; padding: 10px 20px; border-bottom: 1px solid #333; z-index: 10000; display: flex; align-items: center; gap: 20px; }
        .toolbar h2 { color: #00d9ff; font-size: 16px; margin: 0; }
        .toolbar span { color: #888; font-size: 12px; }
        .toolbar button { background: #0f3460; color: #00d9ff; border: 1px solid #00d9ff; padding: 4px 12px; border-radius: 4px; cursor: pointer; font-size: 12px; }
        .toolbar button:hover { background: #00d9ff; color: #1a1a2e; }
        .toolbar button.active { background: #00d9ff; color: #1a1a2e; }
        .main { margin-top: 50px; display: flex; }
        .canvas-wrapper { overflow: auto; border: 2px solid #333; background: #0f0f23; margin: 0 auto; position: relative; flex: 1; }
        .canvas { position: relative; width: %dpx; height: %dpx; }
        .el { position: absolute; overflow: visible; cursor: pointer; transition: outline 0.1s; }
        .el:hover { outline: 2px solid #ff6b6b !important; outline-offset: -1px; z-index: 9999 !important; }
        .el.selected { outline: 2px solid #ffd700 !important; outline-offset: -1px; z-index: 9998 !important; }
        .el-label { position: absolute; top: 0; left: 0; font-size: 9px; line-height: 1.2; padding: 1px 4px; background: rgba(0,0,0,0.75); white-space: nowrap; pointer-events: none; max-width: 100%%; overflow: hidden; text-overflow: ellipsis; }
        .el-text { position: absolute; overflow: hidden; pointer-events: none; }
        .el.hidden-el { border-style: dashed !important; opacity: 0.4; }
        .el.clip { overflow: hidden; }
        .scroll-indicator { position: absolute; right: 2px; top: 2px; background: rgba(0,0,0,0.7); color: #ffcb6b; font-size: 9px; padding: 1px 4px; border-radius: 2px; pointer-events: none; z-index: 100; }
        .sidebar { position: fixed; top: 50px; right: 0; bottom: 0; width: 380px; background: #16213e; border-left: 1px solid #333; overflow-y: auto; padding: 15px; font-size: 12px; display: none; }
        .sidebar.open { display: block; }
        .sidebar h3 { color: #00d9ff; margin: 0 0 10px 0; font-size: 14px; }
        .sidebar .section { margin-bottom: 15px; }
        .sidebar .section-title { color: #888; font-size: 11px; text-transform: uppercase; margin-bottom: 5px; }
        .sidebar table { width: 100%%; border-collapse: collapse; }
        .sidebar td { padding: 2px 6px; border-bottom: 1px solid #1a1a2e; }
        .sidebar td:first-child { color: #888; width: 120px; }
        .sidebar td:last-child { color: #7ee787; word-break: break-all; }
        .sidebar .color-swatch { display: inline-block; width: 12px; height: 12px; border-radius: 2px; vertical-align: middle; margin-right: 4px; border: 1px solid rgba(255,255,255,0.3); }
        .event-log { position: fixed; top: 50px; left: 0; bottom: 0; width: 280px; background: #16213e; border-right: 1px solid #333; overflow-y: auto; padding: 10px; font-size: 11px; display: none; }
        .event-log.open { display: block; }
        .event-log h3 { color: #00d9ff; margin: 0 0 8px 0; font-size: 13px; }
        .event-item { padding: 4px; margin: 2px 0; border-radius: 3px; background: #0f0f23; }
        .event-item .evt-type { color: #ffcb6b; font-weight: bold; }
        .event-item .evt-target { color: #7ee787; }
        .event-item .evt-time { color: #666; font-size: 10px; }
        .legend { position: fixed; bottom: 20px; left: 20px; background: #16213e; padding: 12px; border-radius: 8px; font-size: 11px; z-index: 100; }
        .legend h4 { color: #00d9ff; margin: 0 0 8px 0; }
        .legend-item { margin: 3px 0; display: flex; align-items: center; gap: 6px; }
        .legend-color { width: 16px; height: 10px; border-radius: 2px; border: 1px solid rgba(255,255,255,0.2); }
    </style>
</head>
<body>
    <div class="toolbar">
        <h2>Tenon Debugger</h2>
        <span id="status-bar">Screen: %dx%d | Auto-refresh: ON</span>
        <button onclick="toggleAutoRefresh()" id="btn-auto" class="active">Auto Refresh</button>
        <button onclick="toggleSidebar()">Inspect Panel</button>
        <button onclick="toggleEventLog()">Event Log</button>
        <button onclick="refreshPage()">Refresh Now</button>
        <a href="/debug/devtools" style="color:#00d9ff;margin-left:auto">DevTools</a>
    </div>
    <div class="main">
        <div class="event-log" id="event-log">
            <h3>Event Log</h3>
            <div id="event-list">Loading...</div>
        </div>
        <div class="canvas-wrapper">
            <div class="canvas" id="canvas">
%s
            </div>
        </div>
    </div>
    <div class="sidebar" id="sidebar">
        <h3>Element Inspector</h3>
        <div id="inspector-content">Click an element to inspect</div>
    </div>
    <div class="legend">
        <h4>Element Types</h4>
        <div class="legend-item"><div class="legend-color" style="background:rgba(0,217,255,0.3);border-color:#00d9ff;"></div>View</div>
        <div class="legend-item"><div class="legend-color" style="background:rgba(126,231,135,0.3);border-color:#7ee787;"></div>Text</div>
        <div class="legend-item"><div class="legend-color" style="background:rgba(255,107,107,0.3);border-color:#ff6b6b;"></div>Button</div>
        <div class="legend-item"><div class="legend-color" style="background:rgba(255,203,107,0.3);border-color:#ffcb6b;"></div>ScrollView</div>
        <div class="legend-item"><div class="legend-color" style="background:rgba(195,232,141,0.3);border-color:#c3e88d;"></div>TextInput</div>
        <div class="legend-item"><div class="legend-color" style="background:rgba(130,170,255,0.3);border-color:#82aaff;"></div>Image</div>
        <div class="legend-item"><div class="legend-color" style="background:rgba(199,146,234,0.3);border-color:#c792ea;"></div>Other</div>
    </div>
    <script>
        let selectedEl = null;
        let autoRefresh = true;
        let refreshInterval = null;
        const treeData = %s;

        function startAutoRefresh() {
            if (refreshInterval) clearInterval(refreshInterval);
            refreshInterval = setInterval(() => { if (autoRefresh) location.reload(); }, 500);
        }
        function stopAutoRefresh() { if (refreshInterval) { clearInterval(refreshInterval); refreshInterval = null; } }
        function toggleAutoRefresh() {
            autoRefresh = !autoRefresh;
            const btn = document.getElementById('btn-auto');
            btn.textContent = autoRefresh ? 'Auto Refresh' : 'Auto Refresh (OFF)';
            btn.classList.toggle('active', autoRefresh);
            document.getElementById('status-bar').textContent = 'Screen: %dx%d | Auto-refresh: ' + (autoRefresh ? 'ON' : 'OFF');
            if (autoRefresh) startAutoRefresh(); else stopAutoRefresh();
        }
        startAutoRefresh();
        function toggleSidebar() { document.getElementById('sidebar').classList.toggle('open'); }
        function toggleEventLog() { document.getElementById('event-log').classList.toggle('open'); }
        function refreshPage() { location.reload(); }
        function selectElement(id) {
            if (selectedEl) selectedEl.classList.remove('selected');
            const el = document.getElementById(id);
            if (el) { el.classList.add('selected'); selectedEl = el; }
            showInspector(id);
        }
        function showInspector(id) {
            const node = findNode(treeData, id);
            if (!node) return;
            const container = document.getElementById('inspector-content');
            let html = '';
            html += '<div class="section"><div class="section-title">Element</div><table>';
            html += tr('Type', node.type);
            html += tr('Key', node.key || '-');
            html += tr('Tag', node.tag || '-');
            html += tr('Visible', node.visible ? 'true' : '<span style="color:#ff6b6b">false</span>');
            html += tr('ClipChildren', node.clipChildren ? 'true' : 'false');
            html += '</table></div>';
            html += '<div class="section"><div class="section-title">Bounds</div><table>';
            html += tr('X', node.bounds.x.toFixed(1));
            html += tr('Y', node.bounds.y.toFixed(1));
            html += tr('Width', node.bounds.width.toFixed(1));
            html += tr('Height', node.bounds.height.toFixed(1));
            html += '</table></div>';
            html += '<div class="section"><div class="section-title">Yoga Style</div><table>';
            const y = node.yoga;
            if (y.flexDirection) html += tr('flexDirection', y.flexDirection);
            if (y.justifyContent) html += tr('justifyContent', y.justifyContent);
            if (y.alignItems) html += tr('alignItems', y.alignItems);
            html += tr('flexGrow', y.flexGrow);
            html += tr('flexShrink', y.flexShrink);
            if (y.flexWrap) html += tr('flexWrap', y.flexWrap);
            if (y.positionType) html += tr('positionType', y.positionType);
            if (y.display) html += tr('display', y.display);
            html += tr('width', fmtVal(y.width));
            html += tr('height', fmtVal(y.height));
            html += tr('padding', fmtEdges(y.paddingTop, y.paddingRight, y.paddingBottom, y.paddingLeft));
            html += tr('margin', fmtEdges(y.marginTop, y.marginRight, y.marginBottom, y.marginLeft));
            html += tr('border', fmtEdges(y.borderTop, y.borderRight, y.borderBottom, y.borderLeft));
            html += tr('gap', y.gap);
            html += tr('aspectRatio', y.aspectRatio);
            html += '</table></div>';
            if (node.transform && (node.transform.rotation || node.transform.scaleX !== 1 || node.transform.alpha !== 1)) {
                html += '<div class="section"><div class="section-title">Transform</div><table>';
                if (node.transform.rotation) html += tr('rotation', node.transform.rotation);
                if (node.transform.scaleX !== 1) html += tr('scaleX', node.transform.scaleX);
                if (node.transform.scaleY !== 1) html += tr('scaleY', node.transform.scaleY);
                if (node.transform.alpha !== 1) html += tr('alpha', node.transform.alpha);
                html += '</table></div>';
            }
            if (node.props && Object.keys(node.props).length > 0) {
                html += '<div class="section"><div class="section-title">Properties</div><table>';
                for (const [k, v] of Object.entries(node.props)) {
                    if (k === '_id') continue;
                    if (typeof v === 'string' && v.startsWith('rgba')) {
                        html += '<tr><td>' + k + '</td><td><span class="color-swatch" style="background:' + v + '"></span>' + v + '</td></tr>';
                    } else if (typeof v === 'object') {
                        html += tr(k, JSON.stringify(v));
                    } else {
                        html += tr(k, String(v));
                    }
                }
                html += '</table></div>';
            }
            container.innerHTML = html;
        }
        function findNode(node, id) {
            if (!node) return null;
            if (node._id === id) return node;
            if (node.children) { for (const child of node.children) { const found = findNode(child, id); if (found) return found; } }
            return null;
        }
        function tr(key, val) { return '<tr><td>' + key + '</td><td>' + val + '</td></tr>'; }
        function fmtVal(v) { if (v === undefined || v === null || isNaN(v)) return 'auto'; return v.toFixed(1); }
        function fmtEdges(t, r, b, l) { if (!t && !r && !b && !l) return '0'; return t.toFixed(0) + ' ' + r.toFixed(0) + ' ' + b.toFixed(0) + ' ' + l.toFixed(0); }
        async function loadEvents() {
            try {
                const res = await fetch('/debug/events?limit=30');
                const data = await res.json();
                const list = document.getElementById('event-list');
                if (!data || data.length === 0) { list.innerHTML = '<div style="color:#666">No events yet</div>'; return; }
                list.innerHTML = data.slice().reverse().map(e => {
                    const time = new Date(e.timestamp).toLocaleTimeString().split(' ')[0];
                    let extra = '';
                    if (e.deltaX !== undefined && e.deltaX !== 0) extra += ' dx:' + e.deltaX.toFixed(2);
                    if (e.deltaY !== undefined && e.deltaY !== 0) extra += ' dy:' + e.deltaY.toFixed(2);
                    if (e.x !== undefined) extra += ' @(' + e.x.toFixed(0) + ',' + e.y.toFixed(0) + ')';
                    return '<div class="event-item"><span class="evt-time">' + time + '</span> <span class="evt-type">' + e.type + '</span>' + (e.target ? ' <span class="evt-target">' + e.target + '</span>' : '') + extra + '</div>';
                }).join('');
            } catch (err) { document.getElementById('event-list').innerHTML = '<div style="color:#ff6b6b">Failed to load events</div>'; }
        }
        loadEvents();
        setInterval(loadEvents, 500);
    </script>
</body>
</html>`, screenW, screenH, screenW, screenH, elementsHTML.String(), d.buildTreeJSON(root), screenW, screenH)
}

func (d *Debugger) renderDebugNodeHTML(node *core.DebugNode, html *strings.Builder, depth int, parentX, parentY float32) {
	if node == nil {
		return
	}

	bounds := node.Bounds
	renderX := bounds.X
	renderY := bounds.Y
	relX := renderX - parentX
	relY := renderY - parentY

	if bounds.Width <= 0 || bounds.Height <= 0 {
		for _, child := range node.Children {
			d.renderDebugNodeHTML(child, html, depth+1, parentX, parentY)
		}
		return
	}

	nodeIDCounter++
	id := fmt.Sprintf("n%d", nodeIDCounter)

	typeColor := getElementColor(node.Type)
	borderColor := typeColor + "88"
	bgColor := typeColor + "1a"
	labelColor := typeColor

	if node.Props != nil {
		bgKeys := []string{"backgroundColor", "bgColor", "normalColor"}
		for _, key := range bgKeys {
			if bg, ok := node.Props[key].(string); ok {
				bgColor = bg
				break
			}
		}
		borderKeys := []string{"borderColor"}
		for _, key := range borderKeys {
			if bc, ok := node.Props[key].(string); ok {
				borderColor = bc
				break
			}
		}
		if c, ok := node.Props["color"].(string); ok {
			labelColor = c
		}
	}

	visibleClass := ""
	if !node.Visible {
		visibleClass = " hidden-el"
	}

	clipClass := ""
	if node.ClipChildren {
		clipClass = " clip"
	}

	cssStyles := []string{
		fmt.Sprintf("left:%.1fpx", relX),
		fmt.Sprintf("top:%.1fpx", relY),
		fmt.Sprintf("width:%.1fpx", bounds.Width),
		fmt.Sprintf("height:%.1fpx", bounds.Height),
		fmt.Sprintf("border:1px solid %s", borderColor),
		fmt.Sprintf("background:%s", bgColor),
		fmt.Sprintf("z-index:%d", 1000-depth),
	}

	if node.Props != nil {
		if br, ok := node.Props["borderRadius"].(map[string]interface{}); ok {
			tl := br["topLeft"]
			tr := br["topRight"]
			brr := br["bottomRight"]
			bl := br["bottomLeft"]
			cssStyles = append(cssStyles, fmt.Sprintf("border-radius:%vpx %vpx %vpx %vpx", tl, tr, brr, bl))
		}
		if br, ok := node.Props["borderRadius"].(float32); ok {
			cssStyles = append(cssStyles, fmt.Sprintf("border-radius:%.1fpx", br))
		}
		if br, ok := node.Props["borderRadius"].(float64); ok {
			cssStyles = append(cssStyles, fmt.Sprintf("border-radius:%.1fpx", br))
		}
	}

	if node.Transform.Alpha < 1 {
		cssStyles = append(cssStyles, fmt.Sprintf("opacity:%.2f", node.Transform.Alpha))
	}

	label := node.Type
	if node.Key != "" {
		label = node.Key
	}

	fmt.Fprintf(html, `<div class="el%s%s" id="%s" style="%s" onclick="selectElement('%s')">`, visibleClass, clipClass, id, strings.Join(cssStyles, ";"), id)
	fmt.Fprintf(html, `<div class="el-label" style="color:%s">%s</div>`, labelColor, label)

	textContent := ""
	if node.Type == "Text" {
		if node.Props != nil {
			if content, ok := node.Props["content"].(string); ok {
				textContent = content
			}
		}
	}

	if textContent != "" {
		fontSize := 16.0
		if node.Props != nil {
			if fs, ok := node.Props["fontSize"].(float64); ok {
				fontSize = fs
			}
		}
		textColor := "#fff"
		if node.Props != nil {
			if tc, ok := node.Props["color"].(string); ok {
				textColor = tc
			}
		}
		yoga := node.Yoga
		padTop := yoga.PaddingTop
		padLeft := yoga.PaddingLeft
		padRight := yoga.PaddingRight
		// Guard against NaN from unset yoga values
		if padTop != padTop { padTop = 0 }
		if padLeft != padLeft { padLeft = 0 }
		if padRight != padRight { padRight = 0 }
		textWidth := bounds.Width - padLeft - padRight
		if textWidth < 0 {
			textWidth = 0
		}
		fmt.Fprintf(html, `<div class="el-text" style="top:%.0fpx;left:%.0fpx;width:%.0fpx;font-size:%.0fpx;color:%s;overflow:hidden;text-overflow:ellipsis;white-space:nowrap">%s</div>`, padTop, padLeft, textWidth, fontSize, textColor, escapeHTML(textContent))
	}

	for _, child := range node.Children {
		d.renderDebugNodeHTML(child, html, depth+1, renderX, renderY)
	}

	html.WriteString(`</div>`)
}

func (d *Debugger) buildTreeJSON(root core.Element) string {
	info := root.DebugInfo()
	nodeIDCounter = 0
	assignIDs(&info)
	data, err := jsonMarshal(info)
	if err != nil {
		return "{}"
	}
	return string(data)
}

func assignIDs(node *core.DebugNode) {
	if node == nil {
		return
	}
	nodeIDCounter++
	if node.Props == nil {
		node.Props = make(map[string]interface{})
	}
	node.Props["_id"] = fmt.Sprintf("n%d", nodeIDCounter)
	for _, child := range node.Children {
		if child == nil {
			continue
		}
		assignIDs(child)
	}
}

func enrichEventCounts(node *core.DebugNode, listeners []core.DebugListenerInfo) {
	if node == nil {
		return
	}
	count := 0
	for _, l := range listeners {
		if l.Target == node.Type || l.Target == node.Type+" (capture)" {
			count += l.Count
		}
	}
	if node.Props == nil {
		node.Props = make(map[string]interface{})
	}
	node.Props["_eventCount"] = count
	for _, child := range node.Children {
		enrichEventCounts(child, listeners)
	}
}

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}

func getElementColor(elementType string) string {
	switch elementType {
	case "View":
		return "#00d9ff"
	case "Text":
		return "#7ee787"
	case "Button":
		return "#ff6b6b"
	case "ScrollView":
		return "#ffcb6b"
	case "TextInput":
		return "#c3e88d"
	case "Image":
		return "#82aaff"
	case "Pagination":
		return "#ff9cac"
	case "Card":
		return "#89ddff"
	case "Calendar":
		return "#c792ea"
	case "Checkbox":
		return "#82aaff"
	case "Switch":
		return "#c3e88d"
	case "Slider":
		return "#f78c6c"
	case "Tab":
		return "#89ddff"
	case "Table":
		return "#ffcb6b"
	case "Sidebar":
		return "#c792ea"
	case "Modal":
		return "#ff5370"
	case "Tooltip":
		return "#c3e88d"
	default:
		return "#c792ea"
	}
}

func (d *Debugger) handleDevTools(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(devToolsHTML))
}

const devToolsHTML = `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>Tenon AI Debugger</title>
<style>
* { box-sizing: border-box; margin: 0; padding: 0; }
body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, monospace; background: #0d1117; color: #c9d1d9; height: 100vh; overflow: hidden; display: flex; flex-direction: column; }
/* ---- toolbar ---- */
.toolbar { background: #161b22; border-bottom: 1px solid #30363d; padding: 6px 16px; display: flex; align-items: center; gap: 12px; flex-shrink: 0; }
.toolbar .title { color: #58a6ff; font-weight: 700; font-size: 14px; }
.toolbar .sep { color: #30363d; margin: 0 4px; }
.toolbar .ws-status { font-size: 11px; padding: 2px 8px; border-radius: 10px; }
.toolbar .ws-status.connected { background: #238636; color: #fff; }
.toolbar .ws-status.disconnected { background: #da3633; color: #fff; }
.toolbar .perf-metrics { display: flex; gap: 16px; font-size: 11px; color: #8b949e; margin-left: auto; }
.toolbar .perf-metrics span { white-space: nowrap; }
.toolbar .perf-metrics b { color: #7ee787; }
.toolbar .tab-btn { background: transparent; color: #8b949e; border: none; padding: 4px 12px; cursor: pointer; font-size: 12px; border-radius: 4px; }
.toolbar .tab-btn:hover { background: #21262d; color: #c9d1d9; }
.toolbar .tab-btn.active { background: #1f6feb; color: #fff; }
/* ---- main ---- */
.main { display: flex; flex: 1; overflow: hidden; }
/* ---- left panel: tree ---- */
.tree-panel { width: 320px; min-width: 200px; background: #161b22; border-right: 1px solid #30363d; display: flex; flex-direction: column; overflow: hidden; resize: horizontal; }
.tree-panel .panel-header { padding: 8px 12px; background: #0d1117; border-bottom: 1px solid #30363d; font-size: 11px; font-weight: 700; text-transform: uppercase; color: #8b949e; display: flex; justify-content: space-between; align-items: center; flex-shrink: 0; }
.tree-panel .panel-header button { background: #21262d; border: 1px solid #30363d; color: #c9d1d9; padding: 2px 8px; border-radius: 4px; cursor: pointer; font-size: 11px; }
.tree-panel .panel-header button:hover { background: #30363d; }
.tree-panel .tree { flex: 1; overflow: auto; padding: 4px 0; }
.tree-node { user-select: none; }
.tree-node .row { display: flex; align-items: center; padding: 2px 8px; cursor: pointer; font-size: 12px; line-height: 20px; white-space: nowrap; }
.tree-node .row:hover { background: #1c2129; }
.tree-node .row.selected { background: #1f6feb33; border-left: 2px solid #58a6ff; }
.tree-node .arrow { width: 16px; text-align: center; color: #8b949e; font-size: 10px; flex-shrink: 0; }
.tree-node .arrow.expanded { transform: rotate(90deg); display: inline-block; }
.tree-node .tag { color: #8b949e; font-size: 10px; margin-right: 4px; }
.tree-node .type-name { font-weight: 600; }
.tree-node .key-name { color: #8b949e; margin-left: 6px; font-size: 11px; }
.tree-node .children { display: none; }
.tree-node .children.open { display: block; }
.tree-node .evt-dot { width: 6px; height: 6px; border-radius: 50%; margin-left: 6px; flex-shrink: 0; }
.tree-node .evt-dot.has-events { background: #d2a8ff; }
/* ---- center: preview ---- */
.preview-panel { flex: 1; background: #0d1117; display: flex; flex-direction: column; overflow: hidden; }
.preview-panel .preview-toolbar { padding: 4px 12px; background: #161b22; border-bottom: 1px solid #30363d; display: flex; align-items: center; gap: 8px; flex-shrink: 0; }
.preview-panel .preview-toolbar button { background: #21262d; border: 1px solid #30363d; color: #c9d1d9; padding: 3px 10px; border-radius: 4px; cursor: pointer; font-size: 11px; }
.preview-panel .preview-toolbar button:hover { background: #30363d; }
.preview-panel .preview-toolbar button.active { background: #1f6feb; color: #fff; }
.preview-panel .preview-wrap { flex: 1; overflow: auto; position: relative; background: #1a1f2e; }
.preview-panel .preview-root { position: relative; background: #0d1117; min-width: 100%; min-height: 100%; }
.gui-el { position: absolute; cursor: pointer; overflow: visible; font-size: 0; box-sizing: border-box; }
/* hover outline removed — preview should match GUI exactly */
.gui-el.selected { outline: 2px solid #58a6ff; outline-offset: -1px; z-index: 998; }
.gui-el.highlight-flash { outline: 3px solid #d2a8ff; outline-offset: -1px; z-index: 1000; animation: flash 0.6s ease-out; }
@keyframes flash { 0% { outline-color: #d2a8ff; } 100% { outline-color: transparent; } }
.gui-el .el-tag { position: absolute; top: 0; left: 0; font-size: 9px; line-height: 1.2; padding: 1px 4px; background: rgba(0,0,0,0.85); white-space: nowrap; pointer-events: none; max-width: 100%; overflow: hidden; text-overflow: ellipsis; font-weight: 600; z-index: 1; }
.gui-el .el-text { position: absolute; overflow: hidden; pointer-events: none; font-size: 12px; }
.gui-el.hidden-el { opacity: 0.3; }
.gui-el.clip { overflow: hidden; }
/* ---- right panel: inspector ---- */
.inspector { width: 360px; min-width: 240px; background: #161b22; border-left: 1px solid #30363d; display: flex; flex-direction: column; overflow: hidden; resize: horizontal; }
.inspector .panel-header { padding: 8px 12px; background: #0d1117; border-bottom: 1px solid #30363d; font-size: 11px; font-weight: 700; text-transform: uppercase; color: #8b949e; flex-shrink: 0; }
.inspector .content { flex: 1; overflow: auto; padding: 8px 12px; }
.inspector .section { margin-bottom: 14px; }
.inspector .section-title { font-size: 10px; text-transform: uppercase; color: #8b949e; margin-bottom: 4px; border-bottom: 1px solid #21262d; padding-bottom: 2px; }
.inspector .kv { display: flex; padding: 1px 0; font-size: 11px; }
.inspector .kv .k { color: #8b949e; width: 100px; flex-shrink: 0; }
.inspector .kv .v { color: #7ee787; word-break: break-all; font-family: monospace; }
.inspector .no-select { color: #484f58; padding: 20px; text-align: center; font-size: 13px; }
.inspector .evt-item { background: #0d1117; padding: 4px 8px; margin: 2px 0; border-radius: 3px; font-size: 11px; display: flex; align-items: center; gap: 6px; }
.inspector .evt-item .evt-type-tag { padding: 1px 5px; border-radius: 3px; font-size: 10px; font-weight: 700; }
/* ---- bottom: logs ---- */
.logs-panel { height: 180px; min-height: 60px; background: #161b22; border-top: 1px solid #30363d; display: flex; flex-direction: column; overflow: hidden; resize: vertical; }
.logs-panel .log-tabs { padding: 4px 12px; background: #0d1117; border-bottom: 1px solid #30363d; display: flex; gap: 4px; flex-shrink: 0; }
.logs-panel .log-tab { background: transparent; color: #8b949e; border: none; padding: 3px 10px; cursor: pointer; font-size: 11px; border-radius: 3px; }
.logs-panel .log-tab:hover { background: #21262d; }
.logs-panel .log-tab.active { background: #21262d; color: #c9d1d9; }
.logs-panel .log-list { flex: 1; overflow: auto; padding: 4px 8px; font-size: 11px; font-family: monospace; }
.logs-panel .log-item { padding: 2px 4px; border-bottom: 1px solid #1c2129; display: flex; gap: 8px; }
.logs-panel .log-item .log-time { color: #484f58; white-space: nowrap; }
.logs-panel .log-item .log-type { color: #d2a8ff; font-weight: 700; min-width: 70px; }
.logs-panel .log-item .log-target { color: #7ee787; }
.logs-panel .log-item .log-detail { color: #8b949e; }
.logs-panel .log-item .log-change { color: #f0883e; }
/* ---- resize handles ---- */
.resize-handle { width: 3px; cursor: col-resize; background: transparent; flex-shrink: 0; }
.resize-handle:hover { background: #58a6ff; }
/* ---- toast ---- */
.toast { position: fixed; bottom: 20px; left: 50%; transform: translateX(-50%); background: #238636; color: #fff; padding: 8px 20px; border-radius: 6px; font-size: 12px; z-index: 9999; opacity: 0; transition: opacity 0.2s; pointer-events: none; }
.toast.show { opacity: 1; }
</style>
</head>
<body>
<div class="toolbar">
    <span class="title">Tenon AI Debugger</span>
    <span class="sep">|</span>
    <span class="ws-status disconnected" id="ws-status">Disconnected</span>
    <span class="sep">|</span>
    <button class="tab-btn active" data-tab="inspect" onclick="switchMode('inspect',this)">Inspect</button>
    <button class="tab-btn" data-tab="highlight" onclick="switchMode('highlight',this)">Highlight</button>
    <div class="perf-metrics">
        <span>FPS: <b id="perf-fps">--</b></span>
        <span>Frame: <b id="perf-frame">--</b></span>
        <span>Layout: <b id="perf-layout">--</b></span>
        <span>Draw: <b id="perf-draw">--</b></span>
        <span>Elements: <b id="perf-elems">--</b></span>
    </div>
</div>
<div class="main">
    <div class="tree-panel" id="tree-panel">
        <div class="panel-header">
            <span>Component Tree</span>
            <div>
                <button onclick="collapseAll()" title="Collapse All">-</button>
                <button onclick="expandAll()" title="Expand All">+</button>
            </div>
        </div>
        <div class="tree" id="tree"></div>
    </div>
    <div class="preview-panel">
        <div class="preview-toolbar">
            <button id="btn-auto-refresh" class="active" onclick="toggleAutoRefresh()">Auto Refresh</button>
            <button onclick="refreshPreview()">Refresh</button>
            <span style="font-size:11px;color:#8b949e;" id="preview-size"></span>
        </div>
        <div class="preview-wrap">
            <div class="preview-root" id="preview"></div>
        </div>
    </div>
    <div class="inspector" id="inspector">
        <div class="panel-header">Inspector</div>
        <div class="content" id="inspector-content">
            <div class="no-select">Select an element to inspect</div>
        </div>
    </div>
</div>
<div class="logs-panel" id="logs-panel">
    <div class="log-tabs">
        <button class="log-tab active" onclick="switchLogTab('events',this)">Events</button>
        <button class="log-tab" onclick="switchLogTab('state',this)">State Changes</button>
        <button class="log-tab" onclick="switchLogTab('lifecycle',this)">Lifecycle</button>
    </div>
    <div class="log-list" id="log-list"></div>
</div>
<div class="toast" id="toast"></div>

<script>
// ==================== State ====================
let treeData = null;
let selectedNodeId = null;
let selectedNodeData = null;
let eventLogs = [];
let stateLogs = [];
let lifecycleLogs = [];
let currentLogTab = 'events';
let perfData = {};
let ws = null;
let autoRefresh = true;
let refreshTimer = null;
let mode = 'inspect'; // 'inspect' | 'highlight'
let expandedNodes = new Set();

// ==================== WebSocket ====================
function connectWS() {
    const proto = location.protocol === 'https:' ? 'wss:' : 'ws:';
    ws = new WebSocket(proto + '//' + location.host + '/debug/ws');
    ws.onopen = () => {
        document.getElementById('ws-status').textContent = 'Connected';
        document.getElementById('ws-status').className = 'ws-status connected';
        ws.send(JSON.stringify({type:'getTree'}));
        ws.send(JSON.stringify({type:'getPerf'}));
        ws.send(JSON.stringify({type:'getEvents'}));
        ws.send(JSON.stringify({type:'getListeners'}));
        ws.send(JSON.stringify({type:'getLifecycle'}));
    };
    ws.onclose = () => {
        document.getElementById('ws-status').textContent = 'Disconnected';
        document.getElementById('ws-status').className = 'ws-status disconnected';
        setTimeout(connectWS, 2000);
    };
    ws.onerror = (e) => { console.error('[WS] Error:', e); ws.close(); };
    ws.onmessage = (e) => {
        let msg;
        try {
            msg = JSON.parse(e.data);
        } catch (err) {
            console.error('[WS] JSON parse error:', err, 'data:', e.data.slice(0, 200));
            return;
        }
        console.log('[WS] Received:', msg.type, msg.data ? '(data present)' : '(no data)');
        switch(msg.type) {
        case 'tree': treeData = msg.data; console.log('[DevTools] Tree data keys:', treeData ? Object.keys(treeData) : 'null'); normalizeTree(treeData); if(expandedNodes.size === 0) autoExpand(treeData, 2); renderTree(); console.log('[DevTools] After normalize, treeData.bounds:', treeData ? treeData.bounds : 'null'); if(autoRefresh) renderPreview(); break;
        case 'perf': perfData = msg.data; updatePerf(); break;
        case 'events':
            if (Array.isArray(msg.data)) {
                eventLogs = msg.data;
                if (currentLogTab === 'events') renderLogs();
            }
            break;
        case 'state':
            if (msg.data) {
                if (Array.isArray(msg.data)) {
                    stateLogs.push(...msg.data);
                } else {
                    stateLogs.push(msg.data);
                }
                if (stateLogs.length > 500) stateLogs = stateLogs.slice(-500);
                if (currentLogTab === 'state') renderLogs();
            }
            break;
        case 'lifecycle':
            if (Array.isArray(msg.data)) {
                lifecycleLogs = msg.data;
                if (currentLogTab === 'lifecycle') renderLogs();
            }
            break;
        case 'listeners': break;
        }
    };
}
connectWS();

function escapeHTML(s) {
    if (!s) return '';
    return String(s).replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
}

function normalizeTree(node) {
    if (!node) return;
    if (node.props) {
        node._id = node.props._id || '';
        node._eventCount = node.props._eventCount || 0;
    } else {
        node._id = '';
        node._eventCount = 0;
    }
    if (node.children) node.children.forEach(normalizeTree);
}

function autoExpand(node, maxDepth, depth) {
    if (!node) return;
    depth = depth || 0;
    if (node._id && depth < maxDepth) {
        expandedNodes.add(node._id);
    }
    if (node.children) {
        node.children.forEach(child => autoExpand(child, maxDepth, depth + 1));
    }
}

// ==================== Tree ====================
function renderTree() {
    const el = document.getElementById('tree');
    if (!treeData) { el.innerHTML = '<div style="padding:12px;color:#484f58">No tree data</div>'; return; }
    el.innerHTML = buildTreeNode(treeData, 0);
    if (selectedNodeId) {
        const sel = el.querySelector('[data-node-id="' + selectedNodeId + '"]');
        if (sel) sel.classList.add('selected');
    }
}

function buildTreeNode(node, depth) {
    if (!node) return '';
    const hasChildren = node.children && node.children.length > 0;
    const isExpanded = expandedNodes.has(node._id);
    const color = getTypeColor(node.type);
    const hasEvt = node._eventCount > 0;
    let html = '<div class="tree-node">';
    html += '<div class="row' + (selectedNodeId === node._id ? ' selected' : '') + '" data-node-id="' + (node._id||'') + '" data-depth="' + depth + '" onclick="selectTree(\'' + (node._id||'') + '\',event)" style="padding-left:' + (8+depth*14) + 'px">';
    html += '<span class="arrow' + (isExpanded ? ' expanded' : '') + '" onclick="event.stopPropagation();toggleNode(\'' + (node._id||'') + '\')">' + (hasChildren ? '▶' : ' ') + '</span>';
    html += '<span class="type-name" style="color:' + color + '">&lt;' + (node.type||'?') + '&gt;</span>';
    if (node.key) html += '<span class="key-name">#' + node.key + '</span>';
    if (node.tag) html += '<span class="tag">' + node.tag + '</span>';
    if (hasEvt) html += '<span class="evt-dot has-events" title="' + node._eventCount + ' event listeners"></span>';
    html += '</div>';
    if (hasChildren) {
        html += '<div class="children' + (isExpanded ? ' open' : '') + '">';
        for (const child of node.children) {
            html += buildTreeNode(child, depth + 1);
        }
        html += '</div>';
    }
    html += '</div>';
    return html;
}

function toggleNode(id) {
    if (expandedNodes.has(id)) expandedNodes.delete(id);
    else expandedNodes.add(id);
    renderTree();
}

function expandAll() {
    function collect(node) {
        if (!node) return;
        if (node._id) expandedNodes.add(node._id);
        if (node.children) node.children.forEach(collect);
    }
    collect(treeData);
    renderTree();
}

function collapseAll() {
    expandedNodes.clear();
    renderTree();
}

function selectTree(id, evt) {
    selectedNodeId = id;
    renderTree();
    findAndInspect(id);
    highlightOnPreview(id);
}

function findAndInspect(id) {
    function find(node, id) {
        if (!node) return null;
        if (node._id === id) return node;
        if (node.children) {
            for (const c of node.children) { const r = find(c, id); if (r) return r; }
        }
        return null;
    }
    const node = find(treeData, id);
    if (node) {
        selectedNodeData = node;
        renderInspector(node);
    }
}

// ==================== Preview (DOM-based GUI visualization) ====================
function renderPreview() {
    const wrap = document.getElementById('preview');
    if (!treeData) {
        wrap.innerHTML = '<div style="padding:20px;color:#8b949e">Waiting for tree data... (WS should auto-connect)</div>';
        return;
    }
    if (!treeData.bounds) {
        wrap.innerHTML = '<div style="padding:20px;color:#f85149">Error: treeData has no bounds. Data: ' + escapeHTML(JSON.stringify(treeData).slice(0,200)) + '</div>';
        return;
    }
    const w = treeData.bounds.width || 1280;
    const h = treeData.bounds.height || 720;
    const start = Date.now();
    _previewElCount = 0;
    wrap.style.width = w + 'px';
    wrap.style.height = h + 'px';
    // Match preview background to GUI root background if available
    const rootBg = treeData.props && treeData.props.backgroundColor;
    wrap.style.background = rootBg || '#0d1117';
    try {
        const html = buildPreviewEl(treeData, 0, 0);
        if (!html || html.length === 0) {
            wrap.innerHTML = '<div style="padding:20px;color:#f85149">buildPreviewEl returned empty. treeData: ' + escapeHTML(JSON.stringify(treeData).slice(0, 300)) + '</div>';
            return;
        }
        wrap.innerHTML = html;
    } catch(e) {
        wrap.innerHTML = '<div style="padding:20px;color:#f85149">Preview render error: ' + escapeHTML(e.message) + '<br><pre style="font-size:10px;color:#8b949e;margin-top:8px">' + escapeHTML(e.stack || '') + '</pre></div>';
        return;
    }
    const elapsed = Date.now() - start;
    document.getElementById('preview-size').textContent = w + 'x' + h + ' | ' + _previewElCount + ' elems | ' + elapsed + 'ms';
    if (selectedNodeId) highlightOnPreview(selectedNodeId);
}

let _previewElCount = 0;
function buildPreviewEl(node, parentX, parentY) {
    if (!node) return '';
    const b = node.bounds || {x:0,y:0,width:0,height:0};
    const relX = b.x - parentX;
    const relY = b.y - parentY;
    const color = getTypeColor(node.type);
    const label = node.key || node.type || '?';

    let cls = 'gui-el';
    if (!node.visible) cls += ' hidden-el';
    if (node.clipChildren) cls += ' clip';
    if (selectedNodeId === node._id) cls += ' selected';

    const debugInfo = {
        type: node.type,
        key: node.key || '',
        tag: node.tag || '',
        id: node._id || '',
        bounds: b,
        visible: node.visible,
        clipChildren: node.clipChildren,
        yoga: node.yoga || {},
        transform: node.transform || {},
        props: node.props || {},
        eventCount: node._eventCount || 0,
        classes: node.classes || []
    };

    // 读取实际样式属性（与 GUI 一致），回退到淡色类型色
    const p = node.props || {};
    let bgColor = '';
    let borderColor = color + '88';
    let borderRadius = '';

    if (p.backgroundColor) bgColor = p.backgroundColor;
    else if (p.bgColor) bgColor = p.bgColor;
    else if (p.normalColor) bgColor = p.normalColor;

    if (p.borderColor) borderColor = p.borderColor;

    if (p.borderRadius) {
        if (typeof p.borderRadius === 'object') {
            const br = p.borderRadius;
            borderRadius = (br.topLeft||0) + 'px ' + (br.topRight||0) + 'px ' + (br.bottomRight||0) + 'px ' + (br.bottomLeft||0) + 'px';
        } else if (typeof p.borderRadius === 'number') {
            borderRadius = p.borderRadius + 'px';
        }
    }

    let h = '<div class="' + cls + '"';
    h += ' data-gui-id="' + (node._id||'') + '"';
    h += ' data-gui-type="' + (node.type||'') + '"';
    h += ' data-gui-key="' + (node.key||'') + '"';
    h += ' data-ai-debug=\'' + JSON.stringify(debugInfo).replace(/'/g, "&#39;") + '\'';
    h += ' style="left:' + relX + 'px;top:' + relY + 'px;';
    if (b.width > 0) h += 'width:' + b.width + 'px;';
    if (b.height > 0) h += 'height:' + b.height + 'px;';
    if (bgColor) h += 'background:' + bgColor + ';';
    h += 'border:1px solid ' + borderColor + ';';
    if (borderRadius) h += 'border-radius:' + borderRadius + ';';
    h += 'min-width:2px;min-height:2px;';
    if (b.width <= 0 || b.height <= 0) {
        h += 'width:4px;height:4px;background:' + (bgColor || borderColor) + ';';
    }
    if (node.transform && node.transform.alpha !== undefined && node.transform.alpha < 1) {
        h += 'opacity:' + node.transform.alpha + ';';
    }
    h += '"';
    h += ' onclick="previewClick(\'' + (node._id||'') + '\',event)"';
    h += '>';

    if (b.width > 0 && b.height > 0) {
        _previewElCount++;
        h += '<div class="el-tag" style="color:' + color + '">' + escapeHTML(label) + '</div>';

        if (p.content) {
            const fs = p.fontSize || 16;
            const tc = p.color || '#000';
            const y = node.yoga || {};
            const padLeft = (y.paddingLeft === y.paddingLeft) ? (y.paddingLeft || 0) : 0;
            const padRight = (y.paddingRight === y.paddingRight) ? (y.paddingRight || 0) : 0;
            const textWidth = Math.max(0, b.width - padLeft - padRight);
            // Center text vertically within Text bounds (closer to user expectation than baseline offset)
            h += '<div class="el-text" style="top:0;left:' + padLeft + 'px;width:' + textWidth + 'px;height:' + b.height + 'px;display:flex;align-items:center;font-size:' + fs + 'px;color:' + tc + '">' + escapeHTML(String(p.content)) + '</div>';
        }

        // ProgressBar: show fill portion
        if (node.type === 'ProgressBar' && p.progress > 0 && p.fillColor) {
            const fillW = b.width * p.progress;
            h += '<div style="position:absolute;top:0;left:0;width:' + fillW + 'px;height:100%;background:' + p.fillColor + ';border-radius:inherit;opacity:0.9;"></div>';
        }

        // Checkbox: show box outline (always) + checkmark when checked
        if (node.type === 'Checkbox') {
            const boxSize = p.boxSize || 18;
            const bc = p.borderColor || '#888';
            h += '<div style="position:absolute;left:0;top:50%;transform:translateY(-50%);width:' + boxSize + 'px;height:' + boxSize + 'px;border:1.5px solid ' + bc + ';border-radius:3px;background:' + (p.backgroundColor || 'transparent') + ';"></div>';
            if (p.checked) {
                h += '<div style="position:absolute;left:3px;top:50%;transform:translateY(-50%);width:12px;height:12px;background:' + (p.checkColor || '#fff') + ';border-radius:2px;display:flex;align-items:center;justify-content:center;font-size:9px;color:' + (p.checkColor || '#fff') + ';">✓</div>';
            }
        }

        // Radio: show circle outline (always) + inner dot when selected
        if (node.type === 'Radio') {
            const boxSize = p.boxSize || 18;
            const bc = p.borderColor || '#888';
            h += '<div style="position:absolute;left:0;top:50%;transform:translateY(-50%);width:' + boxSize + 'px;height:' + boxSize + 'px;border:1.5px solid ' + bc + ';border-radius:50%;background:' + (p.backgroundColor || 'transparent') + ';"></div>';
            if (p.selected) {
                h += '<div style="position:absolute;left:5px;top:50%;transform:translateY(-50%);width:8px;height:8px;background:' + (p.innerColor || '#fff') + ';border-radius:50%;"></div>';
            }
        }

        // Switch: show thumb circle
        if (node.type === 'Switch' && p.thumbColor) {
            const thumbR = (b.height / 2) - 2;
            const thumbY = b.height / 2;
            const leftX = thumbR + 2;
            const rightX = b.width - thumbR - 2;
            const progress = p.checked ? 1 : 0;
            const thumbX = leftX + (rightX - leftX) * progress;
            h += '<div style="position:absolute;left:' + thumbX + 'px;top:' + thumbY + 'px;transform:translate(-50%,-50%);width:' + (thumbR*2) + 'px;height:' + (thumbR*2) + 'px;background:' + p.thumbColor + ';border-radius:50%;box-shadow:0 1px 3px rgba(0,0,0,0.3);"></div>';
        }
    }

    if (node.children) {
        for (const c of node.children) h += buildPreviewEl(c, b.x, b.y);
    }
    h += '</div>';
    return h;
}

function previewClick(id, evt) {
    evt.stopPropagation();
    if (mode === 'inspect') {
        selectedNodeId = id;
        renderTree();
        findAndInspect(id);
        highlightOnPreview(id);
    } else if (mode === 'highlight') {
        if (ws && ws.readyState === WebSocket.OPEN) {
            const path = getPathToNode(treeData, id);
            ws.send(JSON.stringify({type:'highlight', data:{path:path}}));
            showToast('Highlight sent to app');
        }
    }
}

function highlightOnPreview(id) {
    const wrap = document.getElementById('preview');
    wrap.querySelectorAll('.gui-el.selected').forEach(el => el.classList.remove('selected'));
    const el = wrap.querySelector('[data-gui-id="' + id + '"]');
    if (el) {
        el.classList.add('selected');
        el.scrollIntoView({behavior:'smooth',block:'nearest'});
    }
}

function getPathToNode(node, id) {
    function find(node, id, path) {
        if (!node) return null;
        path.push(node._id || '');
        if (node._id === id) return path;
        if (node.children) {
            for (const c of node.children) {
                const r = find(c, id, [...path]);
                if (r) return r;
            }
        }
        return null;
    }
    const path = find(node, id, []);
    return path || [];
}

// ==================== Inspector ====================
function renderInspector(node) {
    const el = document.getElementById('inspector-content');
    const b = node.bounds || {};
    const y = node.yoga || {};
    const t = node.transform || {};
    let h = '';

    h += '<div class="section"><div class="section-title">Element</div>';
    h += kv('Type', node.type);
    h += kv('Key', node.key || '-');
    h += kv('Tag', node.tag || '-');
    h += kv('Visible', node.visible ? 'true' : '<span style="color:#da3633">false</span>');
    h += hcls(node.classes);
    h += '</div>';

    h += '<div class="section"><div class="section-title">Bounds</div>';
    h += kv('X', fmt(b.x));
    h += kv('Y', fmt(b.y));
    h += kv('Width', fmt(b.width));
    h += kv('Height', fmt(b.height));
    h += '</div>';

    h += '<div class="section"><div class="section-title">Yoga Layout</div>';
    if (y.flexDirection) h += kv('flexDirection', y.flexDirection);
    if (y.justifyContent) h += kv('justifyContent', y.justifyContent);
    if (y.alignItems) h += kv('alignItems', y.alignItems);
    h += kv('flexGrow', y.flexGrow);
    h += kv('flexShrink', y.flexShrink);
    if (y.flexWrap) h += kv('flexWrap', y.flexWrap);
    if (y.positionType) h += kv('positionType', y.positionType);
    h += kv('width', fmtVal(y.width));
    h += kv('height', fmtVal(y.height));
    h += kv('padding', fmt4(y.paddingTop,y.paddingRight,y.paddingBottom,y.paddingLeft));
    h += kv('margin', fmt4(y.marginTop,y.marginRight,y.marginBottom,y.marginLeft));
    h += kv('border', fmt4(y.borderTop,y.borderRight,y.borderBottom,y.borderLeft));
    if (y.gap) h += kv('gap', y.gap);
    if (y.aspectRatio) h += kv('aspectRatio', y.aspectRatio);
    h += '</div>';

    if (t.rotation || t.scaleX !== 1 || t.scaleY !== 1 || t.alpha !== 1) {
        h += '<div class="section"><div class="section-title">Transform</div>';
        if (t.rotation) h += kv('rotation', t.rotation + '°');
        if (t.scaleX !== 1) h += kv('scaleX', t.scaleX);
        if (t.scaleY !== 1) h += kv('scaleY', t.scaleY);
        if (t.alpha !== 1) h += kv('alpha', t.alpha);
        h += '</div>';
    }

    if (node.props && Object.keys(node.props).length > 0) {
        h += '<div class="section"><div class="section-title">Properties</div>';
        for (const [k, v] of Object.entries(node.props)) {
            if (k === '_id' || k === '_eventCount') continue;
            if (typeof v === 'string' && (v.startsWith('rgba') || v.startsWith('rgb') || v.startsWith('#'))) {
                h += '<div class="kv"><span class="k">' + k + '</span><span class="v"><span style="display:inline-block;width:12px;height:12px;background:' + v + ';border-radius:2px;vertical-align:middle;margin-right:4px;border:1px solid rgba(255,255,255,0.2)"></span>' + v + '</span></div>';
            } else if (typeof v === 'object') {
                h += kv(k, JSON.stringify(v));
            } else {
                h += kv(k, String(v));
            }
        }
        h += '</div>';
    }

    if (node._eventCount > 0) {
        h += '<div class="section"><div class="section-title">Event Listeners (' + node._eventCount + ')</div>';
        h += '<div style="color:#8b949e;font-size:11px">Loading listener details...</div>';
        h += '<div id="evt-detail-' + (node._id||'') + '"></div>';
        h += '</div>';
        setTimeout(() => loadListenersForNode(node), 100);
    }

    el.innerHTML = h;
    function fmt(v) { return v !== undefined ? v.toFixed(1) : '-'; }
    function kv(k, v) { return '<div class="kv"><span class="k">' + k + '</span><span class="v">' + (v !== undefined && v !== null ? v : '-') + '</span></div>'; }
    function fmtVal(v) { return v !== undefined && v !== null && !isNaN(v) ? v.toFixed(1) : 'auto'; }
    function fmt4(a,b,c,d) { return (a||0).toFixed(0)+' '+(b||0).toFixed(0)+' '+(c||0).toFixed(0)+' '+(d||0).toFixed(0); }
    function hcls(cs) { if (!cs || cs.length===0) return ''; return '<div class="kv"><span class="k">classes</span><span class="v">' + cs.join(' ') + '</span></div>'; }
}

async function loadListenersForNode(node) {
    try {
        const res = await fetch('/debug/listeners');
        const data = await res.json();
        const el = document.getElementById('evt-detail-' + (node._id||''));
        if (!el) return;
        const typeColors = {
            'Click': '#ff7b72', 'MouseDown': '#d2a8ff', 'MouseUp': '#a5d6ff',
            'MouseMove': '#79c0ff', 'Scroll': '#7ee787', 'KeyDown': '#f0883e',
            'KeyUp': '#f0883e', 'FocusIn': '#ffa657', 'FocusOut': '#ffa657',
            'MouseEnter': '#56d364', 'MouseLeave': '#56d364'
        };
        const relevant = data.filter(d => d.target === node.type || d.target === (node.type + ' (capture)'));
        if (relevant.length === 0) {
            el.innerHTML = '<div style="color:#484f58;font-size:11px">No detailed info available</div>';
            return;
        }
        el.innerHTML = relevant.map(d => {
            const isCapture = d.target.includes('(capture)');
            return '<div class="evt-item"><span class="evt-type-tag" style="background:' + (typeColors[d.eventType]||'#30363d') + ';color:#fff">' + d.eventType + (isCapture ? ' ↓' : ' ↑') + '</span><span style="color:#8b949e">' + d.count + ' callback' + (d.count>1?'s':'') + '</span></div>';
        }).join('');
    } catch(e) {}
}

// ==================== Logs ====================
function switchLogTab(tab, btn) {
    currentLogTab = tab;
    document.querySelectorAll('.log-tab').forEach(b => b.classList.remove('active'));
    if (btn) btn.classList.add('active');
    renderLogs();
}

function renderLogs() {
    const el = document.getElementById('log-list');
    if (currentLogTab === 'events') return renderEventLogs(el);
    if (currentLogTab === 'state') return renderStateLogs(el);
    if (currentLogTab === 'lifecycle') return renderLifecycleLogs(el);
}

function renderEventLogs(el) {
    if (!eventLogs.length) { el.innerHTML = '<div style="color:#484f58;padding:8px">No events</div>'; return; }
    el.innerHTML = eventLogs.slice().reverse().slice(0,100).map(e => {
        const t = new Date(e.timestamp).toLocaleTimeString();
        let extra = '';
        if (e.x!==undefined) extra += ' @(' + e.x.toFixed(0) + ',' + e.y.toFixed(0) + ')';
        if (e.deltaX) extra += ' dx:' + e.deltaX.toFixed(1);
        if (e.deltaY) extra += ' dy:' + e.deltaY.toFixed(1);
        return '<div class="log-item"><span class="log-time">' + t + '</span><span class="log-type">' + e.type + '</span><span class="log-target">' + (e.target||'') + '</span><span class="log-detail">' + extra + '</span></div>';
    }).join('');
}

function renderStateLogs(el) {
    if (!stateLogs.length) { el.innerHTML = '<div style="color:#484f58;padding:8px">No state changes</div>'; return; }
    el.innerHTML = stateLogs.slice().reverse().slice(0,100).map(s => {
        const t = new Date(s.timestamp).toLocaleTimeString();
        const ov = JSON.stringify(s.oldValue);
        const nv = JSON.stringify(s.newValue);
        return '<div class="log-item"><span class="log-time">' + t + '</span><span class="log-type">State</span><span class="log-target">' + (s.elementType||'') + (s.elementKey?'#'+s.elementKey:'') + '</span><span class="log-change">' + ov + ' → ' + nv + '</span></div>';
    }).join('');
}

function renderLifecycleLogs(el) {
    if (!lifecycleLogs.length) { el.innerHTML = '<div style="color:#484f58;padding:8px">No lifecycle events</div>'; return; }
    el.innerHTML = lifecycleLogs.slice().reverse().slice(0,100).map(l => {
        const t = new Date(l.timestamp).toLocaleTimeString();
        const actColor = l.action === 'mount' ? '#56d364' : '#f85149';
        return '<div class="log-item"><span class="log-time">' + t + '</span><span class="log-type" style="color:' + actColor + '">' + l.action + '</span><span class="log-target">' + l.elementType + (l.elementKey?'#'+l.elementKey:'') + '</span></div>';
    }).join('');
}

// ==================== Perf ====================
function updatePerf() {
    const p = perfData;
    document.getElementById('perf-fps').textContent = p.lastFrameTime ? (1e9/p.lastFrameTime).toFixed(0) : '--';
    document.getElementById('perf-frame').textContent = p.lastFrameTime ? (p.lastFrameTime/1e6).toFixed(2)+'ms' : '--';
    document.getElementById('perf-layout').textContent = p.lastLayoutTime ? (p.lastLayoutTime/1e6).toFixed(2)+'ms' : '--';
    document.getElementById('perf-draw').textContent = p.lastDrawTime ? (p.lastDrawTime/1e6).toFixed(2)+'ms' : '--';
    document.getElementById('perf-elems').textContent = p.elementCount || '--';
}

// ==================== Mode ====================
function switchMode(m, btn) {
    mode = m;
    document.querySelectorAll('.tab-btn').forEach(b => b.classList.remove('active'));
    if (btn) btn.classList.add('active');
    showToast(m === 'inspect' ? 'Inspect mode: click elements to inspect' : 'Highlight mode: click elements to highlight in app');
}

// ==================== Utils ====================
function toggleAutoRefresh() {
    autoRefresh = !autoRefresh;
    const btn = document.getElementById('btn-auto-refresh');
    btn.textContent = autoRefresh ? 'Auto Refresh' : 'Auto Refresh (OFF)';
    btn.classList.toggle('active', autoRefresh);
    if (autoRefresh) {
        refreshTimer = setInterval(() => {
            if (ws && ws.readyState === WebSocket.OPEN) ws.send(JSON.stringify({type:'getTree'}));
        }, 500);
    } else {
        if (refreshTimer) clearInterval(refreshTimer);
    }
}
refreshTimer = setInterval(() => {
    if (ws && ws.readyState === WebSocket.OPEN) ws.send(JSON.stringify({type:'getTree'}));
}, 500);

function refreshPreview() {
    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({type:'getTree'}));
        ws.send(JSON.stringify({type:'getPerf'}));
        ws.send(JSON.stringify({type:'getEvents'}));
        ws.send(JSON.stringify({type:'getLifecycle'}));
    }
}

function getTypeColor(type) {
    const m = {
        'View':'#58a6ff','Text':'#7ee787','Button':'#f0883e','ScrollView':'#d2a8ff',
        'TextInput':'#a5d6ff','Image':'#79c0ff','Pagination':'#ff7b72',
        'Card':'#56d364','Calendar':'#d2a8ff','Checkbox':'#79c0ff',
        'Switch':'#a5d6ff','Slider':'#f0883e','Tab':'#56d364','Table':'#d2a8ff',
        'Sidebar':'#d2a8ff','Modal':'#f85149','Tooltip':'#a5d6ff'
    };
    return m[type] || '#c792ea';
}

function showToast(msg) {
    const t = document.getElementById('toast');
    t.textContent = msg;
    t.classList.add('show');
    setTimeout(() => t.classList.remove('show'), 2000);
}
</script>
</body>
</html>`
