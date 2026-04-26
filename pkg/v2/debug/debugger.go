package debug

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/sjm1327605995/tenon/pkg/v2/core"
)

type Debugger struct {
	engine     *core.Engine
	server     *http.Server
	port       int
	mu         sync.RWMutex
	snapshots  []*LayoutSnapshot
	maxHistory int
	enabled    bool
}

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
	mux.HandleFunc("/", d.handleIndex)

	d.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", d.port),
		Handler: mux,
	}

	go func() {
		fmt.Printf("[Debugger] HTTP server started at http://localhost:%d\n", d.port)
		if err := d.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("[Debugger] Server error: %v\n", err)
		}
	}()

	return nil
}

func (d *Debugger) Stop() error {
	if d.server != nil {
		return d.server.Close()
	}
	return nil
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
	json.NewEncoder(w).Encode(info)
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
	json.NewEncoder(w).Encode(summary)
}

func (d *Debugger) handleSnapshot(w http.ResponseWriter, r *http.Request) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if len(d.snapshots) == 0 {
		http.Error(w, "No snapshots available", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d.snapshots[len(d.snapshots)-1])
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
	json.NewEncoder(w).Encode(diff)
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
		diff.Current = a.Bounds
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

type DiffNode struct {
	Type     string            `json:"type"`
	Key      string            `json:"key,omitempty"`
	Changed  bool              `json:"changed"`
	Changes  []string          `json:"changes"`
	Previous core.LayoutBounds `json:"previous,omitempty"`
	Current  core.LayoutBounds `json:"current,omitempty"`
	Children []*DiffNode       `json:"children,omitempty"`
}

func (d *Debugger) handleHTML(w http.ResponseWriter, r *http.Request) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	root := d.engine.GetRootElement()
	if root == nil {
		http.Error(w, "No root element", http.StatusNotFound)
		return
	}

	html := d.generateHTML(root)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func (d *Debugger) handleIndex(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Tenon Layout Debugger</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; margin: 20px; background: #1a1a2e; color: #eee; }
        h1 { color: #00d9ff; }
        .endpoint { background: #16213e; padding: 15px; margin: 10px 0; border-radius: 8px; border-left: 4px solid #00d9ff; }
        .endpoint h3 { margin: 0 0 10px 0; color: #00d9ff; }
        .endpoint code { background: #0f0f23; padding: 2px 8px; border-radius: 4px; color: #7ee787; }
        .endpoint p { margin: 5px 0; color: #aaa; }
        a { color: #00d9ff; text-decoration: none; }
        a:hover { text-decoration: underline; }
    </style>
</head>
<body>
    <h1>Tenon Layout Debugger</h1>
    <p>AI-assisted debugging interface for Yoga layout tree visualization</p>
    <div class="endpoint">
        <h3><a href="/debug/html">/debug/html</a></h3>
        <code>GET</code>
        <p>Visualize the current layout tree as HTML with CSS absolute positioning</p>
    </div>
    <div class="endpoint">
        <h3><a href="/debug/tree">/debug/tree</a></h3>
        <code>GET</code>
        <p>Get the current layout tree as JSON (full DebugInfo)</p>
    </div>
    <div class="endpoint">
        <h3><a href="/debug/snapshot">/debug/snapshot</a></h3>
        <code>GET</code>
        <p>Get the latest layout snapshot with metadata</p>
    </div>
    <div class="endpoint">
        <h3><a href="/debug/history">/debug/history</a></h3>
        <code>GET</code>
        <p>Get list of all captured layout snapshots</p>
    </div>
    <div class="endpoint">
        <h3><a href="/debug/compare">/debug/compare</a></h3>
        <code>GET</code>
        <p>Compare the last two snapshots and show differences</p>
    </div>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func (d *Debugger) generateHTML(root core.Element) string {
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

	var elementsHTML strings.Builder
	d.renderDebugNodeHTML(&info, &elementsHTML, 0, 0, 0)

	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>Tenon Layout Debugger</title>
    <style>
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body {
            background: #1a1a2e;
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', monospace;
            padding: 20px;
            color: #eee;
        }
        .toolbar {
            position: fixed; top: 0; left: 0; right: 0;
            background: #16213e; padding: 10px 20px;
            border-bottom: 1px solid #333;
            z-index: 10000;
            display: flex; align-items: center; gap: 20px;
        }
        .toolbar h2 { color: #00d9ff; font-size: 16px; margin: 0; }
        .toolbar span { color: #888; font-size: 12px; }
        .toolbar button {
            background: #0f3460; color: #00d9ff; border: 1px solid #00d9ff;
            padding: 4px 12px; border-radius: 4px; cursor: pointer; font-size: 12px;
        }
        .toolbar button:hover { background: #00d9ff; color: #1a1a2e; }
        .main { margin-top: 50px; }
        .canvas-wrapper {
            overflow: auto;
            border: 2px solid #333;
            background: #0f0f23;
            margin: 0 auto;
            position: relative;
        }
        .canvas {
            position: relative;
            width: %dpx;
            height: %dpx;
        }
        .el {
            position: absolute;
            overflow: visible;
            cursor: pointer;
            transition: outline 0.1s;
        }
        .el:hover {
            outline: 2px solid #ff6b6b !important;
            outline-offset: -1px;
            z-index: 9999 !important;
        }
        .el.selected {
            outline: 2px solid #ffd700 !important;
            outline-offset: -1px;
            z-index: 9998 !important;
        }
        .el-label {
            position: absolute; top: 0; left: 0;
            font-size: 9px; line-height: 1.2;
            padding: 1px 4px;
            background: rgba(0,0,0,0.75);
            white-space: nowrap;
            pointer-events: none;
            max-width: 100%%;
            overflow: hidden;
            text-overflow: ellipsis;
        }
        .el-text {
            position: absolute;
            overflow: hidden;
            pointer-events: none;
        }
        .el.hidden-el {
            border-style: dashed !important;
            opacity: 0.4;
        }
        .el.clip {
            overflow: hidden;
        }
        .sidebar {
            position: fixed; top: 50px; right: 0; bottom: 0;
            width: 380px; background: #16213e;
            border-left: 1px solid #333;
            overflow-y: auto; padding: 15px;
            font-size: 12px;
            display: none;
        }
        .sidebar.open { display: block; }
        .sidebar h3 { color: #00d9ff; margin: 0 0 10px 0; font-size: 14px; }
        .sidebar .section { margin-bottom: 15px; }
        .sidebar .section-title { color: #888; font-size: 11px; text-transform: uppercase; margin-bottom: 5px; }
        .sidebar table { width: 100%%; border-collapse: collapse; }
        .sidebar td { padding: 2px 6px; border-bottom: 1px solid #1a1a2e; }
        .sidebar td:first-child { color: #888; width: 120px; }
        .sidebar td:last-child { color: #7ee787; word-break: break-all; }
        .sidebar .color-swatch {
            display: inline-block; width: 12px; height: 12px;
            border-radius: 2px; vertical-align: middle; margin-right: 4px;
            border: 1px solid rgba(255,255,255,0.3);
        }
        .legend {
            position: fixed; bottom: 20px; left: 20px;
            background: #16213e; padding: 12px; border-radius: 8px;
            font-size: 11px; z-index: 100;
        }
        .legend h4 { color: #00d9ff; margin: 0 0 8px 0; }
        .legend-item { margin: 3px 0; display: flex; align-items: center; gap: 6px; }
        .legend-color {
            width: 16px; height: 10px; border-radius: 2px;
            border: 1px solid rgba(255,255,255,0.2);
        }
    </style>
</head>
<body>
    <div class="toolbar">
        <h2>Tenon Debugger</h2>
        <span>Screen: %dx%d</span>
        <button onclick="toggleSidebar()">Inspect Panel</button>
        <button onclick="refreshPage()">Refresh</button>
    </div>
    <div class="main">
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
        const treeData = %s;

        function toggleSidebar() {
            document.getElementById('sidebar').classList.toggle('open');
        }

        function refreshPage() {
            location.reload();
        }

        function selectElement(id) {
            if (selectedEl) selectedEl.classList.remove('selected');
            const el = document.getElementById(id);
            if (el) {
                el.classList.add('selected');
                selectedEl = el;
            }
            showInspector(id);
        }

        function showInspector(id) {
            const node = findNode(treeData, id);
            if (!node) return;
            const container = document.getElementById('inspector-content');
            let html = '';

            html += '<div class="section"><div class="section-title">Element</div>';
            html += '<table>';
            html += tr('Type', node.type);
            html += tr('Key', node.key || '-');
            html += tr('Tag', node.tag || '-');
            html += tr('Visible', node.visible ? 'true' : '<span style="color:#ff6b6b">false</span>');
            html += tr('ClipChildren', node.clipChildren ? 'true' : 'false');
            html += '</table></div>';

            html += '<div class="section"><div class="section-title">Bounds</div>';
            html += '<table>';
            html += tr('X', node.bounds.x.toFixed(1));
            html += tr('Y', node.bounds.y.toFixed(1));
            html += tr('Width', node.bounds.width.toFixed(1));
            html += tr('Height', node.bounds.height.toFixed(1));
            html += '</table></div>';

            html += '<div class="section"><div class="section-title">Yoga Style</div>';
            html += '<table>';
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
                html += '<div class="section"><div class="section-title">Transform</div>';
                html += '<table>';
                if (node.transform.rotation) html += tr('rotation', node.transform.rotation);
                if (node.transform.scaleX !== 1) html += tr('scaleX', node.transform.scaleX);
                if (node.transform.scaleY !== 1) html += tr('scaleY', node.transform.scaleY);
                if (node.transform.alpha !== 1) html += tr('alpha', node.transform.alpha);
                html += '</table></div>';
            }

            if (node.props && Object.keys(node.props).length > 0) {
                html += '<div class="section"><div class="section-title">Properties</div>';
                html += '<table>';
                for (const [k, v] of Object.entries(node.props)) {
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
            if (node.children) {
                for (const child of node.children) {
                    const found = findNode(child, id);
                    if (found) return found;
                }
            }
            return null;
        }

        function tr(key, val) {
            return '<tr><td>' + key + '</td><td>' + val + '</td></tr>';
        }

        function fmtVal(v) {
            if (v === undefined || v === null || isNaN(v)) return 'auto';
            return v.toFixed(1);
        }

        function fmtEdges(t, r, b, l) {
            if (!t && !r && !b && !l) return '0';
            return t.toFixed(0) + ' ' + r.toFixed(0) + ' ' + b.toFixed(0) + ' ' + l.toFixed(0);
        }
    </script>
</body>
</html>`, screenW, screenH, screenW, screenH, elementsHTML.String(), d.buildTreeJSON(root))
}

var nodeIDCounter int

func (d *Debugger) renderDebugNodeHTML(node *core.DebugNode, html *strings.Builder, depth int, parentX, parentY float32) {
	if node == nil {
		return
	}

	bounds := node.Bounds
	if bounds.Width <= 0 || bounds.Height <= 0 {
		for _, child := range node.Children {
			d.renderDebugNodeHTML(child, html, depth+1, parentX, parentY)
		}
		return
	}

	relX := bounds.X - parentX
	relY := bounds.Y - parentY

	nodeIDCounter++
	id := fmt.Sprintf("n%d", nodeIDCounter)

	color := getElementColor(node.Type)
	borderColor := color + "88"
	bgColor := color + "1a"

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
		bgKeys := []string{"backgroundColor", "bgColor", "normalColor"}
		for _, key := range bgKeys {
			if bg, ok := node.Props[key].(string); ok {
				cssStyles = append(cssStyles, fmt.Sprintf("background:%s", bg))
				break
			}
		}
		borderKeys := []string{"borderColor"}
		for _, key := range borderKeys {
			if bc, ok := node.Props[key].(string); ok {
				cssStyles = append(cssStyles, fmt.Sprintf("border:1px solid %s", bc))
				break
			}
		}
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

	textContent := ""
	if node.Type == "Text" {
		if node.Props != nil {
			if content, ok := node.Props["content"].(string); ok {
				textContent = content
			}
		}
	}

	fmt.Fprintf(html, `<div class="el%s%s" id="%s" style="%s" onclick="selectElement('%s')">`, visibleClass, clipClass, id, strings.Join(cssStyles, ";"), id)
	fmt.Fprintf(html, `<div class="el-label" style="color:%s">%s</div>`, color, label)

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
		textWidth := bounds.Width - padLeft - padRight
		if textWidth < 0 {
			textWidth = 0
		}
		fmt.Fprintf(html, `<div class="el-text" style="top:%.0fpx;left:%.0fpx;width:%.0fpx;font-size:%.0fpx;color:%s;overflow:hidden;text-overflow:ellipsis;white-space:nowrap">%s</div>`, padTop, padLeft, textWidth, fontSize, textColor, escapeHTML(textContent))
	}

	for _, child := range node.Children {
		d.renderDebugNodeHTML(child, html, depth+1, bounds.X, bounds.Y)
	}

	html.WriteString(`</div>`)
}

func (d *Debugger) buildTreeJSON(root core.Element) string {
	info := root.DebugInfo()
	nodeIDCounter = 0
	assignIDs(&info)
	data, err := json.Marshal(info)
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
	default:
		return "#c792ea"
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
