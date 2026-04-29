package debug

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type wsClient struct {
	conn     *websocket.Conn
	send     chan []byte
	debugger *Debugger
}

type WebSocketHub struct {
	clients    map[*wsClient]bool
	broadcast  chan []byte
	register   chan *wsClient
	unregister chan *wsClient
	mu         sync.RWMutex
}

func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		clients:    make(map[*wsClient]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *wsClient),
		unregister: make(chan *wsClient),
	}
}

func (h *WebSocketHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("[WSHub] Client connected. Total: %d", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("[WSHub] Client disconnected. Total: %d", len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *WebSocketHub) SendJSON(msgType string, data interface{}) {
	msg := map[string]interface{}{
		"type": msgType,
		"data": data,
	}
	payload, err := safeJSONMarshal(msg)
	if err != nil {
		log.Printf("[WSHub] Failed to marshal JSON: %v", err)
		return
	}
	h.broadcast <- payload
}

func (h *WebSocketHub) SendText(text string) {
	h.broadcast <- []byte(text)
}

func (d *Debugger) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS] Upgrade error: %v", err)
		return
	}

	client := &wsClient{
		conn:     conn,
		send:     make(chan []byte, 256),
		debugger: d,
	}

	d.hub.register <- client

	go d.writePump(client)
	go d.readPump(client)
}

func (d *Debugger) writePump(client *wsClient) {
	defer func() {
		client.conn.Close()
	}()

	for message := range client.send {
		if err := client.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("[WS] Write error: %v", err)
			return
		}
	}
}

func (d *Debugger) readPump(client *wsClient) {
	defer func() {
		d.hub.unregister <- client
		client.conn.Close()
	}()

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[WS] Read error: %v", err)
			}
			break
		}

		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		msgType, _ := msg["type"].(string)
		switch msgType {
		case "getTree":
			root := d.engine.GetRootElement()
			if root != nil {
				info := root.DebugInfo()
				nodeIDCounter = 0
				assignIDs(&info)
				enrichEventCounts(&info, d.engine.GetEventRegistryDebugInfo())
				client.debugger.hub.SendJSON("tree", info)
			}

		case "getPerf":
			perf := d.engine.GetPerfMetrics()
			client.debugger.hub.SendJSON("perf", perf)

		case "getEvents":
			d.mu.RLock()
			limit := 100
			start := len(d.eventLogs) - limit
			if start < 0 {
				start = 0
			}
			events := d.eventLogs[start:]
			d.mu.RUnlock()
			client.debugger.hub.SendJSON("events", events)

		case "getListeners":
			listeners := d.engine.GetEventRegistryDebugInfo()
			client.debugger.hub.SendJSON("listeners", listeners)

		case "getLifecycle":
			logs := d.engine.GetLifecycleLogs()
			client.debugger.hub.SendJSON("lifecycle", logs)

		case "getState":
			d.mu.RLock()
			limit := 100
			start := len(d.stateLogs) - limit
			if start < 0 {
				start = 0
			}
			states := d.stateLogs[start:]
			d.mu.RUnlock()
			client.debugger.hub.SendJSON("state", states)

		case "highlight":
			if data, ok := msg["data"].(map[string]interface{}); ok {
				if path, ok := data["path"].([]interface{}); ok {
					intPath := make([]int, len(path))
					for i, p := range path {
						if f, ok := p.(float64); ok {
							intPath[i] = int(f)
						}
					}
					d.highlightPath = intPath
				}
			}

		case "getSnapshot":
			d.mu.RLock()
			if len(d.snapshots) > 0 {
				client.debugger.hub.SendJSON("snapshot", d.snapshots[len(d.snapshots)-1])
			}
			d.mu.RUnlock()
		}
	}
}

func (d *Debugger) sendWS(msgType string, data interface{}) {
	if d.hub == nil {
		return
	}
	d.hub.SendJSON(msgType, data)
}

func (d *Debugger) startWSPushLoop() {
	go d.hub.Run()

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if !d.IsEnabled() {
				continue
			}

			d.hub.mu.RLock()
			hasClients := len(d.hub.clients) > 0
			d.hub.mu.RUnlock()

			if !hasClients {
				continue
			}

			perf := d.engine.GetPerfMetrics()
			d.hub.SendJSON("perf", perf)
		}
	}()
}
